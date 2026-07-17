package mcpconn

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Default guard limits applied when Options leaves them zero. They are conservative on purpose: the
// gateway can raise them per connector, but an unconfigured client is already fail-closed.
const (
	defaultDialTimeout     = 10 * time.Second
	defaultRequestTimeout  = 30 * time.Second
	defaultMaxPayloadBytes = 4 << 20 // 4 MiB
	clientName             = "rill-act-mcp-client"
	clientVersion          = "0.1.0"
)

// Sentinel errors the gateway can branch on. ErrSchemaChanged in particular is load-bearing: it means a
// call was stopped because the remote schema no longer matches what discovery pinned, not that the call
// failed on the wire.
var (
	// ErrToolNotDiscovered is returned when CallTool is given an effective name outside this connector's
	// namespace, or with an empty expected schema hash (an action carrying no pin must not execute).
	ErrToolNotDiscovered = errors.New("mcpconn: tool not discovered; run ListTools first")
	// ErrSchemaChanged is returned when a tool's current remote schema differs from the caller's expected
	// pin. The call is not executed: the connector must be re-validated before the tool can be trusted again.
	ErrSchemaChanged = errors.New("mcpconn: remote tool schema changed since discovery")
	// ErrToolGone is returned when a pinned tool is no longer offered by the server at call time.
	ErrToolGone = errors.New("mcpconn: tool no longer offered by server")
)

// ToolHints are the server-advertised MCP annotations for a tool. THEY ARE NOT TRUSTED: a server may lie
// about them, so they never decide policy here. They are surfaced so the gateway can apply them under an
// explicit opt-in (ConnectorConfig.TrustReadOnlyHint), per §16.3. Defaults follow the MCP spec: a missing
// destructive/open-world hint is assumed true, a missing read-only/idempotent hint is assumed false.
type ToolHints struct {
	ReadOnly    bool
	Destructive bool
	Idempotent  bool
	OpenWorld   bool
}

// RemoteTool is one tool discovered on a remote MCP server, in the shape the Tool Gateway consumes. Name
// is the effective, namespaced identifier the gateway and model use; RawName is what the server knows.
type RemoteTool struct {
	// Name is the namespaced effective name: mcp.<connector>.<rawName>.
	Name string
	// RawName is the tool name as exposed by the remote server.
	RawName string
	// Connector is the connector name this tool came from.
	Connector string
	// Title and Description are display/hint text passed through unmodified.
	Title       string
	Description string
	// InputSchema and OutputSchema are the server's JSON Schemas, marshaled deterministically. OutputSchema
	// is nil when the server does not declare one.
	InputSchema  json.RawMessage
	OutputSchema json.RawMessage
	// SchemaHash is the pin: a sha256 over the tool name and its input/output schema. A later call whose
	// recomputed hash differs is stopped with ErrSchemaChanged.
	SchemaHash string
	// Hints are the untrusted server annotations. Do not make trust decisions from these here.
	Hints ToolHints
}

// RemoteToolResult is the untrusted result of a remote tool call, in the shape the gateway feeds back to
// the model (after its own validation and redaction).
type RemoteToolResult struct {
	// Text is the concatenation of the text content blocks the tool returned.
	Text string
	// StructuredContent is the tool's structured output, if any.
	StructuredContent json.RawMessage
	// IsError reports a tool-level error (the tool ran and reported failure), distinct from a transport or
	// guard error, which are returned as a Go error instead.
	IsError bool
}

// MCPClient is the seam the Tool Gateway consumes. ListTools performs discovery and returns each tool's
// schema hash for the caller to pin in the run snapshot; CallTool executes a tool by its effective name,
// comparing the remote schema against the caller-supplied expected hash and refusing to run on drift. The
// expected hash comes from the action's own snapshot (ExecuteRequest.ToolVersion), so an in-flight run is
// pinned to an immutable value rather than to mutable client state: a later ListTools on the same client
// cannot repin a call that is already approved. Implementations are safe for concurrent use.
type MCPClient interface {
	ListTools(ctx context.Context) ([]RemoteTool, error)
	CallTool(ctx context.Context, name, expectedSchemaHash string, argsJSON json.RawMessage) (RemoteToolResult, error)
}

// Options configures a Client. All fields are optional; zero values fall back to the package defaults.
type Options struct {
	// BearerToken is the resolved credential to present as "Authorization: Bearer ...". The caller resolves
	// it from the secret referenced by ConnectorConfig.Auth; this package never touches the secret manager.
	BearerToken string
	// AllowPrivateNetworks is the out-of-band trust flag that lets this client reach loopback, private and
	// special-purpose ranges (and use plaintext http). It is a BUILD-TIME decision set only by trusted Go
	// code — packaged/local development connectors and the in-process test server — and MUST NEVER be derived
	// from a tenant connector definition. Off by default: a Tenant connector is always fail-closed to public
	// hosts over https.
	AllowPrivateNetworks bool
	// DialTimeout bounds the TCP/TLS connect; RequestTimeout bounds a single ListTools or CallTool operation.
	DialTimeout    time.Duration
	RequestTimeout time.Duration
	// MaxPayloadBytes caps each response body.
	MaxPayloadBytes int64
}

func (o Options) withDefaults() Options {
	if o.DialTimeout <= 0 {
		o.DialTimeout = defaultDialTimeout
	}
	if o.RequestTimeout <= 0 {
		o.RequestTimeout = defaultRequestTimeout
	}
	if o.MaxPayloadBytes <= 0 {
		o.MaxPayloadBytes = defaultMaxPayloadBytes
	}
	return o
}

// Client is the concrete MCPClient over a single MCP connector. It connects lazily on first use and reuses
// the session across calls. The gateway is expected to create one Client per connector per run: ListTools
// at snapshot creation, then CallTool for each dispatched tool. It holds NO per-tool pin state: the schema
// pin an in-flight action verifies against is supplied by the caller on each CallTool, so discovery on this
// client never repins an already-approved call.
type Client struct {
	cfg     *ConnectorConfig
	opts    Options
	sdk     *mcp.Client
	tport   *mcp.StreamableClientTransport
	httpcli *http.Client

	mu      sync.Mutex
	session *mcp.ClientSession
}

var _ MCPClient = (*Client)(nil)

// NewClient builds a Client for a validated connector config. It does not connect: the first ListTools or
// CallTool establishes the session, so a blocked host surfaces as an error from the operation, not the
// constructor.
func NewClient(cfg *ConnectorConfig, opts Options) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("mcpconn: nil config")
	}
	// Validate with the caller's trust flag: only trusted construction (AllowPrivateNetworks) may relax the
	// https requirement for a loopback dev endpoint. A tenant connector is validated strictly (flag off).
	if err := cfg.validateWith(opts.AllowPrivateNetworks); err != nil {
		return nil, err
	}
	opts = opts.withDefaults()

	httpcli := guardedHTTPClient(cfg.Network, opts.AllowPrivateNetworks, opts.DialTimeout, opts.MaxPayloadBytes, opts.BearerToken)

	tport := &mcp.StreamableClientTransport{
		Endpoint:   cfg.URL,
		HTTPClient: httpcli,
		MaxRetries: 0,
		// We do not need server-initiated messages for discovery and calls, and dropping the standalone SSE
		// stream keeps every response a bounded request/response pair for the payload cap.
		DisableStandaloneSSE: true,
		OAuthHandler:         nil,
	}

	sdk := mcp.NewClient(&mcp.Implementation{
		Name:       clientName,
		Title:      "",
		Version:    clientVersion,
		WebsiteURL: "",
		Icons:      nil,
	}, nil)

	return &Client{
		cfg:     cfg,
		opts:    opts,
		sdk:     sdk,
		tport:   tport,
		httpcli: httpcli,
		mu:      sync.Mutex{},
		session: nil,
	}, nil
}

// ensureSession connects on first use and memoizes the session. It holds c.mu, so callers must not.
func (c *Client) ensureSession(ctx context.Context) (*mcp.ClientSession, error) {
	if c.session != nil {
		return c.session, nil
	}
	sess, err := c.sdk.Connect(ctx, c.tport, nil)
	if err != nil {
		return nil, fmt.Errorf("mcpconn: connect %s: %w", c.cfg.Name, err)
	}
	c.session = sess
	return sess, nil
}

// ListTools discovers the remote tools and returns them with effective, namespaced names and a schema hash
// each. It is the snapshot the gateway takes at the start of a run; the gateway (not this client) stores the
// hashes, so nothing here is retained across calls.
func (c *Client) ListTools(ctx context.Context) ([]RemoteTool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.RequestTimeout)
	defer cancel()

	c.mu.Lock()
	defer c.mu.Unlock()

	sess, err := c.ensureSession(ctx)
	if err != nil {
		return nil, err
	}

	var tools []RemoteTool
	// Tools ranges over the paginated tools/list, so a server that pages its catalog is handled transparently.
	for t, err := range sess.Tools(ctx, nil) {
		if err != nil {
			return nil, fmt.Errorf("mcpconn: list tools %s: %w", c.cfg.Name, err)
		}
		rt, err := c.toRemoteTool(t)
		if err != nil {
			return nil, err
		}
		tools = append(tools, rt)
	}

	return tools, nil
}

// CallTool executes a tool by its effective name, verifying the live remote schema against expectedSchemaHash
// before dispatching. expectedSchemaHash is the pin the action carries from its own snapshot, NOT client state,
// so a schema that drifted since that snapshot stops the call with ErrSchemaChanged and never dispatches (§16.3).
// An empty expected hash refuses the call: an action with no pin must not execute. argsJSON is the
// already-validated argument object as JSON.
func (c *Client) CallTool(ctx context.Context, name, expectedSchemaHash string, argsJSON json.RawMessage) (RemoteToolResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.RequestTimeout)
	defer cancel()

	if expectedSchemaHash == "" {
		return RemoteToolResult{}, fmt.Errorf("%w: %q has no pinned schema hash", ErrToolNotDiscovered, name)
	}
	prefix := EffectiveName(c.cfg.Name, "")
	if !strings.HasPrefix(name, prefix) {
		return RemoteToolResult{}, fmt.Errorf("%w: %q is not a tool of connector %q", ErrToolNotDiscovered, name, c.cfg.Name)
	}
	rawName := strings.TrimPrefix(name, prefix)

	c.mu.Lock()
	defer c.mu.Unlock()

	sess, err := c.ensureSession(ctx)
	if err != nil {
		return RemoteToolResult{}, err
	}

	// Re-validate the caller's pin against the live schema before dispatching. This is the fail-closed
	// enforcement of §16.3: a remote schema different from the action's snapshot stops the call.
	//
	// TODO(act, #8): this is a TOCTOU that the MCP protocol cannot fully close. verifyPin re-lists tools (one remote
	// request) and CallTool then invokes the tool (a second remote request), so a MALICIOUS server can advertise the
	// pinned schema V1 on the re-list and serve a different schema V2 on the call: there is no single round trip that
	// both validates and invokes, so an absolute V1->V1 guarantee does not exist over MCP. The pin still defends against
	// benign drift and a server that changes its catalog between snapshots; a hostile server is mitigated by the layers
	// around this client, not by the pin: connector trust tiers, human approval of the exact args, and the gateway
	// gating every write. Do not represent verifyPin as a security boundary against a hostile server on its own.
	if err := c.verifyPin(ctx, sess, rawName, expectedSchemaHash); err != nil {
		return RemoteToolResult{}, err
	}

	var args any
	if len(argsJSON) > 0 {
		if err := json.Unmarshal(argsJSON, &args); err != nil {
			return RemoteToolResult{}, fmt.Errorf("mcpconn: invalid arguments for %q: %w", name, err)
		}
	}

	res, err := sess.CallTool(ctx, &mcp.CallToolParams{
		Meta:      nil,
		Name:      rawName,
		Arguments: args,
	})
	if err != nil {
		return RemoteToolResult{}, fmt.Errorf("mcpconn: call %q: %w", name, err)
	}
	return toRemoteResult(res)
}

// verifyPin fetches the current schema for rawName and compares its hash against expectedHash. It errors with
// ErrSchemaChanged on drift and ErrToolGone if the tool has vanished. Holds nothing extra: caller holds c.mu.
func (c *Client) verifyPin(ctx context.Context, sess *mcp.ClientSession, rawName, expectedHash string) error {
	for t, err := range sess.Tools(ctx, nil) {
		if err != nil {
			return fmt.Errorf("mcpconn: re-list tools %s: %w", c.cfg.Name, err)
		}
		if t.Name != rawName {
			continue
		}
		hash, _, _, err := hashToolSchema(t)
		if err != nil {
			return err
		}
		if hash != expectedHash {
			return fmt.Errorf("%w: %s", ErrSchemaChanged, rawName)
		}
		return nil
	}
	return fmt.Errorf("%w: %s", ErrToolGone, rawName)
}

// Close tears down the session if one is open. It is safe to call more than once.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.session == nil {
		return nil
	}
	err := c.session.Close()
	c.session = nil
	return err
}

// toRemoteTool converts an SDK tool into a RemoteTool, computing its schema hash and namespaced name.
func (c *Client) toRemoteTool(t *mcp.Tool) (RemoteTool, error) {
	hash, inJSON, outJSON, err := hashToolSchema(t)
	if err != nil {
		return RemoteTool{}, err
	}
	return RemoteTool{
		Name:         EffectiveName(c.cfg.Name, t.Name),
		RawName:      t.Name,
		Connector:    c.cfg.Name,
		Title:        t.Title,
		Description:  t.Description,
		InputSchema:  inJSON,
		OutputSchema: outJSON,
		SchemaHash:   hash,
		Hints:        toHints(t.Annotations),
	}, nil
}

// EffectiveName builds the namespaced tool name used everywhere outside the wire: mcp.<connector>.<tool>.
func EffectiveName(connector, tool string) string {
	return "mcp." + connector + "." + tool
}

// hashToolSchema marshals the tool's input and output schemas deterministically and returns their pin hash
// alongside the marshaled forms. map[string]any (the shape the client receives) marshals with sorted keys,
// so the same schema always hashes the same way.
func hashToolSchema(t *mcp.Tool) (hash string, inJSON, outJSON json.RawMessage, err error) {
	inBytes, err := json.Marshal(t.InputSchema)
	if err != nil {
		return "", nil, nil, fmt.Errorf("mcpconn: marshal input schema for %q: %w", t.Name, err)
	}
	if t.OutputSchema != nil {
		outBytes, err := json.Marshal(t.OutputSchema)
		if err != nil {
			return "", nil, nil, fmt.Errorf("mcpconn: marshal output schema for %q: %w", t.Name, err)
		}
		outJSON = outBytes
	}
	envelope, err := json.Marshal(struct {
		Name   string          `json:"name"`
		Input  json.RawMessage `json:"input"`
		Output json.RawMessage `json:"output,omitempty"`
	}{Name: t.Name, Input: inBytes, Output: outJSON})
	if err != nil {
		return "", nil, nil, fmt.Errorf("mcpconn: marshal schema envelope for %q: %w", t.Name, err)
	}
	sum := sha256.Sum256(envelope)
	return fmt.Sprintf("%x", sum), inBytes, outJSON, nil
}

// toHints copies the SDK annotations into our untrusted ToolHints, applying the MCP spec defaults for the
// pointer fields (a nil destructive/open-world hint means true).
func toHints(a *mcp.ToolAnnotations) ToolHints {
	if a == nil {
		return ToolHints{ReadOnly: false, Destructive: true, Idempotent: false, OpenWorld: true}
	}
	destructive := true
	if a.DestructiveHint != nil {
		destructive = *a.DestructiveHint
	}
	openWorld := true
	if a.OpenWorldHint != nil {
		openWorld = *a.OpenWorldHint
	}
	return ToolHints{
		ReadOnly:    a.ReadOnlyHint,
		Destructive: destructive,
		Idempotent:  a.IdempotentHint,
		OpenWorld:   openWorld,
	}
}

// toRemoteResult maps an SDK CallToolResult into the gateway-facing RemoteToolResult, joining text blocks
// and carrying structured output through.
func toRemoteResult(res *mcp.CallToolResult) (RemoteToolResult, error) {
	var b strings.Builder
	for _, content := range res.Content {
		if tc, ok := content.(*mcp.TextContent); ok {
			b.WriteString(tc.Text)
		}
	}
	out := RemoteToolResult{
		Text:              b.String(),
		StructuredContent: nil,
		IsError:           res.IsError,
	}
	if res.StructuredContent != nil {
		sc, err := json.Marshal(res.StructuredContent)
		if err != nil {
			return RemoteToolResult{}, fmt.Errorf("mcpconn: marshal structured content: %w", err)
		}
		out.StructuredContent = sc
	}
	return out, nil
}
