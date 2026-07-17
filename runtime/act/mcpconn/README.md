# runtime/act/mcpconn — Kairos Act outbound MCP client (Fase 2A)

Experimental. This package gives Act the **client** role for the Model Context Protocol: it connects a
Kairos agent to a remote MCP server, discovers its tools, and calls them under egress and schema-pinning
guards (issue kairosagentica/kairos-cloud#123). Design:
`docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md`, §16.3 (the decisions banner prevails:
default-open tools, optional include/exclude, opt-in `trust_read_only_hint`).

The fork already uses the official MCP Go SDK as a **server** (`runtime/ai/mcp.go`,
`runtime/server/mcp.go`). This package is the mirror image: the same SDK, used only as a client. It never
decides whether a tool may run — that is the Tool Gateway's job. It decides whether a server can be
**reached** and whether a tool's schema still **matches** what discovery pinned.

## The seam the gateway consumes

```go
type MCPClient interface {
    ListTools(ctx context.Context) ([]RemoteTool, error)
    CallTool(ctx context.Context, name, expectedSchemaHash string, argsJSON json.RawMessage) (RemoteToolResult, error)
}
```

- `ListTools` performs `tools/list`, hashes each tool's input/output schema (the **pin**), and returns the
  tools with effective, namespaced names `mcp.<connector>.<tool>`. It is the discovery snapshot the gateway
  takes at run start; the gateway (not the client) stores each hash.
- `CallTool` takes an effective name, the **expected schema hash** the action pinned in its own snapshot, and
  validated argument JSON. It re-reads the tool's live schema and refuses to dispatch (`ErrSchemaChanged`) if
  it drifted from the expected hash. The expected hash is supplied per call rather than read from client
  state, so a later `ListTools` on the same client cannot repin a call that is already approved. An empty
  expected hash (an action with no pin) is refused (`ErrToolNotDiscovered`); a tool that vanished is
  `ErrToolGone`.

`RemoteTool` carries the untrusted server annotations in `Hints` (read-only, destructive, idempotent,
open-world). **These do not decide policy here.** The gateway applies them only under the connector's
explicit `TrustReadOnlyHint` opt-in (§16.3).

## Why a Config type, not a formal Rill driver

The connector is modeled as `ConnectorConfig` (a plain value with `Validate`) parsed by
`ParseConnectorConfig`, **not** as a registered driver in `runtime/drivers`. Wiring a formal driver would
mean touching the driver registry, the connector spec plumbing, and the reconciler surface for a v1 that
only needs "parse, validate, connect, guard". A Config type owned by the act package is enough for the
gateway to consume and keeps the blast radius inside `runtime/act`. Promoting it to a real driver (so a
project can declare `type: connector`/`driver: mcp` in YAML and have the reconciler validate it) is a
follow-up, tracked for the gateway/integration frente.

## Egress protection (SSRF)

Every request goes through a guarded `*http.Client` (`ssrf.go`):

- **Host allowlist + scheme** are checked per request, including redirect hops (redirects are refused
  outright with a distinct `ErrRedirectBlocked`, since a refused redirect fires only after a response).
  `https` is mandatory unless the trusted Go flag `Options.AllowPrivateNetworks` is set (dev/local, the
  in-process mock). `private_ranges` is not a tenant knob: only `deny` is accepted in config.
- **Resolved IP** is checked at dial time via `net.Dialer.Control`, which runs after DNS resolution with
  the concrete address. Unless `AllowPrivateNetworks` is set, loopback, RFC1918/unique-local, link-local,
  multicast, unspecified and special-purpose ranges (CGNAT, benchmarking, protocol assignments, the
  TEST-NET blocks, reserved space, IPv6 documentation and NAT64) are refused (`isBlockedIP`). Because the
  check runs on the dialed IP, DNS-rebinding cannot slip a private address past the hostname check.
- **Payload cap**: each response body is capped (`ErrPayloadTooLarge`); the standalone SSE stream is
  disabled so every response is a bounded request/response pair.
- **Timeouts**: a dial timeout on the dialer and a per-operation deadline on `ListTools`/`CallTool`.

The credential is a **reference** in the config (`auth.secret`); the caller resolves it server-side and
passes the token via `Options.BearerToken`. This package never reads the secret manager. The token is
attached below the egress guard, so it is only ever set on a request already bound for an allowed host.

## Schema pinning

`hashToolSchema` is `sha256` over the tool name and its JSON-marshaled input/output schemas. `map[string]any`
(the shape the client receives over the wire) marshals with sorted keys, so an unchanged schema always
hashes the same. `ListTools` returns one hash per effective name for the gateway to pin in the run snapshot;
`CallTool` re-reads the live schema and compares it against the hash the caller passes (the action's pin),
holding no per-tool state of its own. A `notifications/tools/list_changed` invalidation is not wired for v1: the authoritative check is
the per-call re-read, which is strictly stronger (it catches drift even without a notification). Wiring the
notification as a cheaper early-invalidation signal is a follow-up.

## Usage

```go
cfg, err := mcpconn.ParseConnectorConfig("jira_mcp", props) // props: decoded connector YAML body
client, err := mcpconn.NewClient(cfg, mcpconn.Options{BearerToken: resolvedSecret})
defer client.Close()

tools, err := client.ListTools(ctx)                                        // discovery + hashes, at run snapshot
res, err := client.CallTool(ctx, tools[0].Name, tools[0].SchemaHash, args)  // re-validates the pin, then dispatches
```

One `Client` per connector per run: `ListTools` once, then `CallTool` per dispatched tool.
