package act

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/rilldata/rill/runtime/ai"
)

// This file wires the tool gateway to real, runtime-backed collaborators so a singleton gateway can serve any agent's
// outbound MCP connectors: the registry discovers a proposed tool's descriptor live, the executor performs the write
// through an MCP client built at action time, the secret resolver reads the connector's bearer from the instance's
// variables, and the proposer surfaces the tool call the runtime/ai loop already captured. The gateway itself is
// unchanged; these are the production implementations of its Registry, Executor, Secrets and the executor's Proposer.

// instanceContextKey scopes an instance ID onto a context. The gateway is shared across instances, but its
// SecretResolver.Resolve takes only a secret name, so the executor binds the instance onto the step context before the
// gateway's execute/verify steps, and the resolver reads it back here.
type instanceContextKey struct{}

// withInstance returns ctx carrying instanceID, for the shared SecretResolver to resolve against the right project.
func withInstance(ctx context.Context, instanceID string) context.Context {
	return context.WithValue(ctx, instanceContextKey{}, instanceID)
}

func instanceFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(instanceContextKey{}).(string)
	return v, ok
}

// mcpConnectorContextKey scopes the MCP connector captured in a run's immutable snapshot onto a context. The gateway is
// MCP-agnostic and resolves a connector by name at action time; on the governed execute/verify path the segmented
// workflow binds the snapshot's connector onto the step context (mirroring withInstance), and connectorFor reads it
// back here so the write runs against exactly the connector config the approver reviewed (§8.3).
type mcpConnectorContextKey struct{}

// withMCPConnector returns ctx carrying the connector captured in the run's snapshot, so the runtime-backed resolver
// freezes the connector's target/allowed_hosts/approval to what was reviewed rather than re-reading the live agent.
func withMCPConnector(ctx context.Context, conn ai.MCPConnector) context.Context {
	return context.WithValue(ctx, mcpConnectorContextKey{}, conn)
}

func mcpConnectorFromContext(ctx context.Context) (ai.MCPConnector, bool) {
	conn, ok := ctx.Value(mcpConnectorContextKey{}).(ai.MCPConnector)
	return conn, ok
}

// withGovernedConnector binds the instance and, when one was captured, the snapshot's connector onto ctx for a governed
// gateway step (execute/verify): the shared SecretResolver then targets the run's project, and connectorFor freezes the
// reviewed connector config. A zero-value connector (ok=false) leaves only the instance bound, so a step with no
// captured connector falls back to the live agent definition.
func withGovernedConnector(ctx context.Context, instanceID string, conn ai.MCPConnector, ok bool) context.Context {
	ctx = withInstance(ctx, instanceID)
	if ok {
		ctx = withMCPConnector(ctx, conn)
	}
	return ctx
}

// mcpToolResolver resolves an agent's outbound MCP connector at action time. It loads the agent's immutable snapshot
// to find the connector a proposed tool belongs to, resolves the connector's bearer secret from the instance's
// variables (server-side, never from the model), and builds an mcpconn client below its egress guard. Both the tool
// registry (which discovers a tool's descriptor) and the action executor (which performs the write) go through it, so
// the gateway can serve any agent's connectors without per-run construction.
type mcpToolResolver struct {
	rt        *runtime.Runtime
	provider  ai.AgentDefinitionProvider
	newClient func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error)
}

func newMCPToolResolver(rt *runtime.Runtime) *mcpToolResolver {
	return &mcpToolResolver{
		rt:       rt,
		provider: NewCatalogAgentProvider(rt),
		newClient: func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error) {
			return mcpconn.NewClient(cfg, mcpconn.Options{BearerToken: token})
		},
	}
}

// connectorFor returns the named connector plus its resolved bearer token. It prefers the connector captured in the
// run's immutable snapshot (bound onto the step context by the segmented workflow): freezing the connector's
// target/allowed_hosts/approval to what was reviewed is what makes an approved action execute against exactly the config
// the approver saw (§8.3), so an admin editing the connector (URL, allowed_hosts, approval) between proposal and
// approval cannot redirect the approved write. Only the bearer SECRET is resolved live (inst.ResolveVariables), so a
// rotated token is still picked up and the credential value is never frozen into a snapshot. When no connector is
// captured on the context (a path outside the segmented workflow), it falls back to the LIVE agent definition,
// preserving the pre-existing behavior. It fails closed: an unknown connector, or a declared secret that resolves
// empty, is an error rather than an unauthenticated call. The token is never logged nor wrapped into an error.
func (r *mcpToolResolver) connectorFor(ctx context.Context, instanceID, agentName, connector string) (ai.MCPConnector, string, error) {
	conn, err := r.resolveConnectorConfig(ctx, instanceID, agentName, connector)
	if err != nil {
		return ai.MCPConnector{}, "", err
	}
	var token string
	if conn.AuthSecret != "" {
		inst, err := r.rt.Instance(ctx, instanceID)
		if err != nil {
			return ai.MCPConnector{}, "", fmt.Errorf("act: load instance %q: %w", instanceID, err)
		}
		// ResolveVariables(true) lower-cases keys (Rill variable names are case-insensitive), so look the secret up by
		// its lowercased name rather than the literal as written on the connector.
		token = inst.ResolveVariables(true)[strings.ToLower(conn.AuthSecret)]
		if token == "" {
			return ai.MCPConnector{}, "", fmt.Errorf("act: mcp connector %q secret %q is empty", connector, conn.AuthSecret)
		}
	}
	return conn, token, nil
}

// resolveConnectorConfig returns the connector's frozen configuration (target, allowed_hosts, approval, auth secret
// NAME), preferring the connector captured in the run's snapshot on the context and falling back to the live agent
// definition when none is bound. The returned value carries no secret value: connectorFor resolves that live.
func (r *mcpToolResolver) resolveConnectorConfig(ctx context.Context, instanceID, agentName, connector string) (ai.MCPConnector, error) {
	// Governed execute/verify path: use the connector frozen in the run's immutable snapshot (bound by the workflow),
	// so a connector edit between proposal and approval cannot change what the approved action executes against (§8.3).
	if captured, ok := mcpConnectorFromContext(ctx); ok && captured.Name == connector {
		return captured, nil
	}
	// Fallback: no captured connector on the context (e.g. a read-only lookup outside the segmented workflow). Resolve
	// the live agent definition, matching the pre-existing behavior.
	snap, err := r.provider.GetAgent(ctx, instanceID, agentName)
	if err != nil {
		return ai.MCPConnector{}, fmt.Errorf("act: resolve agent %q: %w", agentName, err)
	}
	for i := range snap.MCPConnectors {
		if snap.MCPConnectors[i].Name == connector {
			return snap.MCPConnectors[i], nil
		}
	}
	return ai.MCPConnector{}, fmt.Errorf("act: agent %q has no mcp connector %q", agentName, connector)
}

// clientFor builds an mcpconn client for the named connector, resolving its secret. The caller must close the client.
func (r *mcpToolResolver) clientFor(ctx context.Context, instanceID, agentName, connector string) (mcpconn.MCPClient, ai.MCPConnector, error) {
	conn, token, err := r.connectorFor(ctx, instanceID, agentName, connector)
	if err != nil {
		return nil, ai.MCPConnector{}, err
	}
	props := map[string]any{
		"transport": mcpconn.TransportStreamableHTTP,
		"url":       conn.URL,
		"auth":      map[string]any{"secret": conn.AuthSecret},
		"network": map[string]any{
			"allowed_hosts":  conn.AllowedHosts,
			"private_ranges": mcpconn.PrivateRangesDeny,
		},
		"trust_read_only_hint": conn.TrustReadOnlyHint,
	}
	cfg, err := mcpconn.ParseConnectorConfig(conn.Name, props)
	if err != nil {
		return nil, ai.MCPConnector{}, fmt.Errorf("act: mcp connector %q config: %w", connector, err)
	}
	client, err := r.newClient(cfg, token)
	if err != nil {
		return nil, ai.MCPConnector{}, fmt.Errorf("act: mcp connector %q client: %w", connector, err)
	}
	return client, conn, nil
}

// connectorOf parses the connector name out of an effective tool name (mcp.<connector>.<raw>).
func connectorOf(tool string) (string, bool) {
	rest, ok := strings.CutPrefix(tool, "mcp.")
	if !ok {
		return "", false
	}
	name, _, ok := strings.Cut(rest, ".")
	if !ok || name == "" {
		return "", false
	}
	return name, true
}

// runtimeToolRegistry is the production ToolRegistry: it resolves a proposed tool's descriptor by discovering it live
// from the agent's MCP connector, pinning the schema exactly as the run's loop did. Membership in the connector's
// discovered tools is the allowlist (default-open model, gated by approval); a tool the connector does not offer, or a
// name that is not an MCP action tool, resolves as unknown and is denied by the gateway.
type runtimeToolRegistry struct {
	resolver *mcpToolResolver
}

var _ ToolRegistry = (*runtimeToolRegistry)(nil)

func (reg *runtimeToolRegistry) Lookup(ctx context.Context, instanceID, agentName, tool string) (ToolDescriptor, bool, error) {
	connector, ok := connectorOf(tool)
	if !ok {
		return ToolDescriptor{}, false, nil
	}
	client, conn, err := reg.resolver.clientFor(ctx, instanceID, agentName, connector)
	if err != nil {
		return ToolDescriptor{}, false, err
	}
	defer closeMCPClient(client)

	tools, err := client.ListTools(ctx)
	if err != nil {
		return ToolDescriptor{}, false, fmt.Errorf("act: discover tools for connector %q: %w", connector, err)
	}
	for i := range tools {
		if tools[i].Name != tool {
			continue
		}
		desc := ToolDescriptor{
			Name:      tools[i].Name,
			Connector: connector,
			Version:   tools[i].SchemaHash,
			// MCP tools declare no reliable idempotency class, so treat every one as unknown: it can never be
			// auto-executed and is never auto-retried on an uncertain result (§12). Approval is the barrier.
			Class:         ClassUnknown,
			ReadOnlyHint:  tools[i].Hints.ReadOnly,
			TrustReadOnly: conn.TrustReadOnlyHint,
			AuthSecret:    conn.AuthSecret,
		}
		if len(tools[i].InputSchema) > 0 {
			var schema jsonschema.Schema
			if err := json.Unmarshal(tools[i].InputSchema, &schema); err != nil {
				return ToolDescriptor{}, false, fmt.Errorf("act: parse input schema for tool %q: %w", tool, err)
			}
			desc.InputSchema = &schema
		}
		return desc, true, nil
	}
	return ToolDescriptor{}, false, nil
}

// runtimeActionExecutor is the production ActionExecutor: it builds an MCP client for the proposal's connector and
// performs the write through the shared mcpExecutor classification logic. The connector's secret is resolved
// server-side; the client lives only for the call.
type runtimeActionExecutor struct {
	resolver *mcpToolResolver
}

var _ ActionExecutor = (*runtimeActionExecutor)(nil)

func (x *runtimeActionExecutor) Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error) {
	client, _, err := x.resolver.clientFor(ctx, req.InstanceID, req.Trace.AgentName, req.Connector)
	if err != nil {
		// Pre-dispatch failure: no request left the client, so the write provably did not land (§12).
		return ExecuteResult{Outcome: OutcomeFailed}, err
	}
	defer closeMCPClient(client)
	return NewMCPActionExecutor(client).Execute(ctx, req)
}

func (x *runtimeActionExecutor) Verify(_ context.Context, _ VerifyRequest) (VerifyResult, error) {
	// Generic MCP has no standard confirm-this-write call, so a generic connector is unverifiable (mcpExecutor.Verify).
	return VerifyResult{Confirmed: false, Detail: "verification not supported for generic MCP tools"}, nil
}

// mcpSecretResolver resolves a secret reference to plaintext from an instance's variables. The gateway uses it only to
// build the executor call and to add credentials to the output scrub set; a resolved value is never persisted, logged
// or returned to the model. The instance comes from the context (bound by the executor before the gateway steps),
// since one resolver is shared across instances.
type mcpSecretResolver struct {
	rt *runtime.Runtime
}

var _ SecretResolver = (*mcpSecretResolver)(nil)

func (s *mcpSecretResolver) Resolve(ctx context.Context, name string) (string, error) {
	instanceID, ok := instanceFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("act: cannot resolve secret %q: no instance bound on context", name)
	}
	inst, err := s.rt.Instance(ctx, instanceID)
	if err != nil {
		return "", fmt.Errorf("act: resolve secret %q: %w", name, err)
	}
	v := inst.ResolveVariables(true)[strings.ToLower(name)]
	if v == "" {
		return "", fmt.Errorf("act: secret %q is not set", name)
	}
	return v, nil
}

// capturedProposer surfaces the write action the runtime/ai loop captured (ProposeInput.Captured) as the run's
// structured tool proposal. It is the production Proposer: the loop already chose the tool and arguments, so this is a
// pure, deterministic mapping, safe to run inside a checkpointed step. A nil capture is a pure investigation.
type capturedProposer struct{}

var _ Proposer = capturedProposer{}

func (capturedProposer) Propose(_ context.Context, in ProposeInput) (ToolProposal, bool, error) {
	if in.Captured == nil {
		return ToolProposal{}, false, nil
	}
	c := in.Captured
	// The identity is the model's own tool-call ID (the captured call message ID): distinct per call, so a run that
	// proposes several actions keys each on its own ledger row. Fall back to the effective tool name only if the capture
	// carried no ID, preserving a non-empty identity for the ledger.
	toolCallID := c.ToolCallID
	if toolCallID == "" {
		toolCallID = c.Tool
	}
	return ToolProposal{
		ToolCallID: toolCallID,
		Tool:       c.Tool,
		Connector:  c.Connector,
		Args:       c.Args,
		Summary:    c.Summary,
	}, true, nil
}

// NewMCPGateway builds the production action gateway: it routes an agent's approved MCP tool call through policy, the
// action ledger, server-side secret resolution and an idempotent MCP write, all resolved from the runtime at action
// time. ledger is the run store (its agent_actions table is the ledger); logger may be nil.
func NewMCPGateway(rt *runtime.Runtime, ledger ActionLedger, logger *slog.Logger) *Gateway {
	resolver := newMCPToolResolver(rt)
	return &Gateway{
		Registry: &runtimeToolRegistry{resolver: resolver},
		Ledger:   ledger,
		Executor: &runtimeActionExecutor{resolver: resolver},
		Secrets:  &mcpSecretResolver{rt: rt},
		Logger:   logger,
	}
}

// NewCapturedProposer returns the production Proposer, which surfaces the write action the runtime/ai loop captured.
func NewCapturedProposer() Proposer { return capturedProposer{} }

// mcpConnector returns the named connector from the snapshot, if present. It is how the executor recovers a proposed
// tool's approval posture (the snapshot is immutable for the run, so this is deterministic).
func mcpConnector(snapshot *ai.AgentSnapshot, name string) (ai.MCPConnector, bool) {
	if snapshot == nil {
		return ai.MCPConnector{}, false
	}
	for i := range snapshot.MCPConnectors {
		c := &snapshot.MCPConnectors[i]
		if c.Name == name {
			return *c, true
		}
	}
	return ai.MCPConnector{}, false
}

// connectorAutoApproves reports whether the connector's approval posture auto-approves an action on rawTool, applying
// the glob exceptions: with approval "auto" every tool auto-approves unless it matches a require_approval pattern;
// otherwise (manual, the default) a tool needs approval unless it matches an auto_approve pattern. It is deliberately
// fail-safe toward approval: an unset or unrecognized posture is manual.
func connectorAutoApproves(conn ai.MCPConnector, rawTool string) bool {
	if conn.Approval == "auto" {
		return !matchesToolGlob(conn.RequireApproval, rawTool)
	}
	return matchesToolGlob(conn.AutoApprove, rawTool)
}

// matchesToolGlob reports whether name matches any of the glob patterns (path.Match semantics: *, ?, [set]). A
// malformed pattern never matches; the parser rejects those at definition time, so this is defense in depth.
func matchesToolGlob(patterns []string, name string) bool {
	for _, p := range patterns {
		if ok, err := path.Match(p, name); err == nil && ok {
			return true
		}
	}
	return false
}

// rawToolName strips the mcp.<connector>. prefix from an effective tool name, returning the server's raw tool name,
// which is what a connector's approval globs match against.
func rawToolName(tool, connector string) string {
	return strings.TrimPrefix(tool, "mcp."+connector+".")
}

// closeMCPClient closes an MCP client if it supports Close (the mcpconn client does); a close error is not actionable.
func closeMCPClient(c mcpconn.MCPClient) {
	if closer, ok := c.(interface{ Close() error }); ok {
		_ = closer.Close()
	}
}
