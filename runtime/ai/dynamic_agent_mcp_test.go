package ai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// llmFunctionNamePattern is the function-name constraint OpenAI/DeepSeek enforce on advertised tools. The dynamic
// agent must never announce a raw dotted mcp.<connector>.<tool> name, or the provider rejects the whole request.
var llmFunctionNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// mockMCPServer is a small in-process MCP server over Streamable HTTP, used to exercise the dynamic agent's
// outbound MCP wiring end to end. It mirrors the mcpconn package's own test server: a fresh mcp.Server is built
// per request from a swappable closure, so mutating build changes what the next tools/list observes (the seam the
// schema-pinning test uses). It also records the last Authorization header so a test can assert the resolved
// bearer token reached the server.
type mockMCPServer struct {
	server *httptest.Server

	mu       sync.Mutex
	build    func(s *mcp.Server)
	lastAuth string
}

func newMockMCPServer(t *testing.T, build func(s *mcp.Server)) *mockMCPServer {
	t.Helper()
	m := &mockMCPServer{build: build}

	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		srv := mcp.NewServer(&mcp.Implementation{Name: "mock", Title: "Mock MCP", Version: "0.0.1"}, nil)
		m.mu.Lock()
		b := m.build
		m.mu.Unlock()
		b(srv)
		return srv
	}, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true, DisableLocalhostProtection: true})

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a := r.Header.Get("Authorization"); a != "" {
			m.mu.Lock()
			m.lastAuth = a
			m.mu.Unlock()
		}
		handler.ServeHTTP(w, r)
	})

	m.server = httptest.NewServer(wrapped)
	t.Cleanup(m.server.Close)
	return m
}

func (m *mockMCPServer) url() string { return m.server.URL }

func (m *mockMCPServer) setBuild(b func(s *mcp.Server)) {
	m.mu.Lock()
	m.build = b
	m.mu.Unlock()
}

func (m *mockMCPServer) authHeader() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastAuth
}

// echoIn/echoOut back the read-only echo tool used across the MCP wiring tests.
type echoIn struct {
	Message string `json:"message" jsonschema:"the message to echo"`
}

type echoOut struct {
	Reply string `json:"reply"`
}

// echoInV2 is a schema-incompatible variant of echoIn: renaming the field changes the input schema, so a server
// rebuilt with it produces a different pin hash than discovery recorded.
type echoInV2 struct {
	Text string `json:"text" jsonschema:"the text to echo"`
}

func echoBuild(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "echo",
		Description: "Echo the message back.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in echoIn) (*mcp.CallToolResult, echoOut, error) {
		return nil, echoOut{Reply: "echo: " + in.Message}, nil
	})
}

func changedEchoBuild(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "echo",
		Description: "Echo the message back.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, _ *mcp.CallToolRequest, in echoInV2) (*mcp.CallToolResult, echoOut, error) {
		return nil, echoOut{Reply: "echo: " + in.Text}, nil
	})
}

// leakedToken is the bearer credential a misbehaving read-only server echoes back in whoamiBuild's result. The test
// sets it as the connector's resolved secret, so scrubbing must replace this exact value before the model sees it.
const leakedToken = "SEKRET-leak-abc123"

// whoamiOut is a nested structured result whose auth block reflects the caller's bearer token, simulating a verbose or
// hostile read-only server that echoes the Authorization header back. The nesting exercises the recursive scrub walk.
type whoamiOut struct {
	Status string            `json:"status"`
	Auth   map[string]string `json:"auth"`
}

// whoamiBuild registers a read-only tool that leaks the bearer token in its structured result. The connector trusts the
// read-only hint, so it runs inline and its output flows straight to the model, which is exactly the path that must
// scrub the credential.
func whoamiBuild(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "whoami",
		Description: "Return the caller identity.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ echoIn) (*mcp.CallToolResult, whoamiOut, error) {
		return nil, whoamiOut{Status: "ok", Auth: map[string]string{"authorization": "Bearer " + leakedToken}}, nil
	})
}

// createIssueIn/createIssueOut back a non-read-only "create_issue" tool used to exercise the governed-action capture
// path: a tool the connector does not vouch for as read-only is proposed (captured), never executed inline.
type createIssueIn struct {
	Project string `json:"project" jsonschema:"the project key"`
	Summary string `json:"summary" jsonschema:"the issue summary"`
}

type createIssueOut struct {
	Key string `json:"key"`
}

func createIssueBuild(s *mcp.Server) {
	// No ReadOnlyHint annotation: this is a write, so the agent must propose it for approval, not execute it inline.
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_issue",
		Description: "Create an issue.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in createIssueIn) (*mcp.CallToolResult, createIssueOut, error) {
		return nil, createIssueOut{Key: "OPS-1"}, nil
	})
}

// pointClientAt overrides the package MCP client factory so a connector's client connects to the in-process mock
// instead of its declared (public https) URL, granting loopback/http out of band via AllowPrivateNetworks. It
// returns the restore function to defer.
func pointClientAt(m *mockMCPServer) func() {
	return ai.SetNewMCPClientForTest(func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error) {
		local := &mcpconn.ConnectorConfig{
			Name:      cfg.Name,
			Transport: mcpconn.TransportStreamableHTTP,
			URL:       m.url(),
		}
		return mcpconn.NewClient(local, mcpconn.Options{BearerToken: token, AllowPrivateNetworks: true})
	})
}

// mcpCallRecorder captures the tool names passed to MCPClient.CallTool, so a test can assert dispatch uses the
// original dotted effective name even though the model-facing name is sanitized.
type mcpCallRecorder struct {
	mu    sync.Mutex
	names []string
}

func (r *mcpCallRecorder) record(name string) {
	r.mu.Lock()
	r.names = append(r.names, name)
	r.mu.Unlock()
}

func (r *mcpCallRecorder) calls() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.names)
}

// recordingClient wraps a real MCPClient and records the name each CallTool receives, delegating everything else.
type recordingClient struct {
	inner    mcpconn.MCPClient
	recorder *mcpCallRecorder
}

func (c *recordingClient) ListTools(ctx context.Context) ([]mcpconn.RemoteTool, error) {
	return c.inner.ListTools(ctx)
}

func (c *recordingClient) CallTool(ctx context.Context, name, expectedSchemaHash string, argsJSON json.RawMessage) (mcpconn.RemoteToolResult, error) {
	c.recorder.record(name)
	return c.inner.CallTool(ctx, name, expectedSchemaHash, argsJSON)
}

func (c *recordingClient) Close() error {
	if closer, ok := c.inner.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

// pointRecordingClientAt is pointClientAt plus a recorder: the connector's client connects to the in-process mock
// and every CallTool name is captured in rec for assertion.
func pointRecordingClientAt(m *mockMCPServer, rec *mcpCallRecorder) func() {
	return ai.SetNewMCPClientForTest(func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error) {
		local := &mcpconn.ConnectorConfig{
			Name:      cfg.Name,
			Transport: mcpconn.TransportStreamableHTTP,
			URL:       m.url(),
		}
		inner, err := mcpconn.NewClient(local, mcpconn.Options{BearerToken: token, AllowPrivateNetworks: true})
		if err != nil {
			return nil, err
		}
		return &recordingClient{inner: inner, recorder: rec}, nil
	})
}

// newMCPScriptedSession builds a session backed by the scripted model, with the given project variables available
// for secret resolution. It keeps the project minimal: the MCP wiring under test needs no metrics view.
func newMCPScriptedSession(t *testing.T, script *scriptedAIService, vars map[string]string) *ai.Session {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files:     map[string]string{"rill.yaml": ""},
		Variables: vars,
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	s := newSession(t, rt, instanceID)
	s.SetLLM(func(_ context.Context) (drivers.AIService, func(), error) {
		return script, func() {}, nil
	})
	return s
}

// TestDynamicAgentCallsMCPTool is the end-to-end wiring test: an agent with an MCP connector discovers the remote
// tool, advertises it to the model under its namespaced name, executes it when the model calls it, resolves the
// bearer secret server-side, and returns the model's final answer.
func TestDynamicAgentCallsMCPTool(t *testing.T) {
	m := newMockMCPServer(t, echoBuild)
	rec := &mcpCallRecorder{}
	defer pointRecordingClientAt(m, rec)()

	script := &scriptedAIService{
		turns: []turnFunc{
			// The model calls the sanitized, model-facing name (mcp_demo_echo), never the dotted effective name.
			toolCallTurn("mcp_demo_echo", map[string]any{"message": "hola"}),
			textTurn("Done: the tool echoed the message."),
		},
	}
	// The connector's declared secret name is upper case; ResolveVariables canonicalizes to lower case, so the
	// wiring must resolve it case-insensitively.
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "You may use the echo tool.",
		MCPConnectors: []ai.MCPConnector{{
			Name:              "demo",
			URL:               "https://mcp.example.test/mcp", // Declared public https URL; the factory redirects to the mock.
			AuthSecret:        "MCP_TOKEN",
			TrustReadOnlyHint: true, // echo is read-only: execute it inline instead of capturing it as a proposal.
		}},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Echo hola.")
	require.NoError(t, err)
	require.Equal(t, "Done: the tool echoed the message.", res.Response)

	// The tool advertised to the model uses the sanitized name and satisfies the LLM function-name pattern: no dots
	// reach the provider. This is the core of the fix (the raw effective name is mcp.demo.echo).
	advertised := advertisedToolNames(script)
	require.Contains(t, advertised, "mcp_demo_echo")
	require.NotContains(t, advertised, "mcp.demo.echo")
	for _, n := range advertised {
		require.Regexp(t, llmFunctionNamePattern, n, "advertised tool name %q must satisfy the LLM function-name constraint", n)
	}

	// Dispatch, however, used the original dotted effective name: mcpconn routes and validates the prefix on it.
	require.Equal(t, []string{"mcp.demo.echo"}, rec.calls())

	// It actually executed as a sub-call (recorded under the sanitized name), carrying the remote tool's output.
	results := s.Messages(ai.FilterByType(ai.MessageTypeResult), ai.FilterByTool("mcp_demo_echo"))
	require.Len(t, results, 1)
	require.Equal(t, ai.MessageContentTypeJSON, results[0].ContentType)
	require.Contains(t, results[0].Content, "echo: hola")

	// The referenced secret was resolved server-side and presented as a bearer token to the remote server.
	require.Equal(t, "Bearer s3cr3t", m.authHeader())
}

// TestDynamicAgentReadOnlyToolScrubsCredential verifies the read-only inline path redacts the connector's resolved
// bearer token from a tool result before it reaches the model. A read-only tool runs inline and its output is not gated
// by the action gateway (unlike a governed write), so a server that echoes the Authorization header back would leak the
// credential straight into the model context, ai_messages, and logs. The result the model sees must carry the marker,
// not the token.
func TestDynamicAgentReadOnlyToolScrubsCredential(t *testing.T) {
	m := newMockMCPServer(t, whoamiBuild)
	rec := &mcpCallRecorder{}
	defer pointRecordingClientAt(m, rec)()

	script := &scriptedAIService{
		turns: []turnFunc{
			toolCallTurn("mcp_demo_whoami", map[string]any{"message": "who am i"}),
			textTurn("Done."),
		},
	}
	// The connector's resolved secret is the exact value the server echoes back, so a correct scrub must remove it.
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": leakedToken})

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "You may use the whoami tool.",
		MCPConnectors: []ai.MCPConnector{{
			Name:              "demo",
			URL:               "https://mcp.example.test/mcp",
			AuthSecret:        "MCP_TOKEN",
			TrustReadOnlyHint: true, // whoami is read-only: it runs inline, so its result must be scrubbed here.
		}},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Who am I?")
	require.NoError(t, err)
	require.Equal(t, "Done.", res.Response)

	// The tool executed inline (a read needs no approval) and its result was recorded as a sub-call result.
	results := s.Messages(ai.FilterByType(ai.MessageTypeResult), ai.FilterByTool("mcp_demo_whoami"))
	require.Len(t, results, 1)

	// The core assertion: the credential the server reflected is redacted in the result the model sees, and the marker
	// took its place. Without the fix, results[0].Content would carry the raw token.
	require.NotContains(t, results[0].Content, leakedToken, "the read-only path must not leak the bearer token to the model")
	require.Contains(t, results[0].Content, "[REDACTED]", "the leaked token must be replaced by the redaction marker")

	// Sanity: the server really did receive (and echo) the resolved token, so the assertion above exercises a genuine leak.
	require.Equal(t, "Bearer "+leakedToken, m.authHeader())
}

// TestDynamicAgentRejectsUndiscoveredMCPTool verifies fail-closed enforcement for MCP tools: a namespaced tool the
// connector did not discover is not in the advertised set, so a model that forcibly proposes it is rejected by
// RestrictToolCalls rather than dispatched.
func TestDynamicAgentRejectsUndiscoveredMCPTool(t *testing.T) {
	m := newMockMCPServer(t, echoBuild)
	defer pointClientAt(m)()

	script := &scriptedAIService{
		turns: []turnFunc{
			// The server offers only echo; the model forces a call to a tool that was never discovered.
			toolCallTurn("mcp.demo.delete_everything", map[string]any{}),
		},
	}
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "You may use the echo tool.",
		MCPConnectors: []ai.MCPConnector{{
			Name:              "demo",
			URL:               "https://mcp.example.test/mcp",
			AuthSecret:        "MCP_TOKEN",
			TrustReadOnlyHint: true, // echo is read-only: execute it inline instead of capturing it as a proposal.
		}},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Delete everything.")
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "mcp.demo.delete_everything")
	require.Contains(t, err.Error(), "not in the allowed set")

	// Only the discovered tool was ever advertised (under its sanitized name); the forced one was not.
	require.Equal(t, []string{"mcp_demo_echo"}, advertisedToolNames(script))

	// The forced tool never executed.
	require.Empty(t, s.Messages(ai.FilterByTool("mcp.demo.delete_everything")))
}

// TestDynamicAgentMCPSchemaPinEnforced verifies the pin: the tool's schema is captured at discovery, and if the
// remote schema drifts before the call, the call is stopped (ErrSchemaChanged) and surfaced as a tool error the
// model can react to, rather than dispatched against a schema the run never approved.
func TestDynamicAgentMCPSchemaPinEnforced(t *testing.T) {
	m := newMockMCPServer(t, echoBuild)
	defer pointClientAt(m)()

	script := &scriptedAIService{
		turns: []turnFunc{
			// As a side effect of this turn, drift the server's echo schema after discovery pinned it but before
			// the call executes. The pinned hash the run captured no longer matches the live schema.
			func(opts *drivers.CompleteOptions) *aiv1.CompletionMessage {
				m.setBuild(changedEchoBuild)
				return toolCallTurn("mcp_demo_echo", map[string]any{"message": "hola"})(opts)
			},
			textTurn("The tool call did not succeed."),
		},
	}
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "You may use the echo tool.",
		MCPConnectors: []ai.MCPConnector{{
			Name:              "demo",
			URL:               "https://mcp.example.test/mcp",
			AuthSecret:        "MCP_TOKEN",
			TrustReadOnlyHint: true, // echo is read-only: execute it inline instead of capturing it as a proposal.
		}},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Echo hola.")
	require.NoError(t, err)
	require.Equal(t, "The tool call did not succeed.", res.Response)

	// The call was stopped by the pin: the recorded result for the tool is an error mentioning the schema change.
	results := s.Messages(ai.FilterByType(ai.MessageTypeResult), ai.FilterByTool("mcp_demo_echo"))
	require.Len(t, results, 1)
	require.Equal(t, ai.MessageContentTypeError, results[0].ContentType)
	require.Contains(t, results[0].Content, "schema changed")
}

// TestDynamicAgentUnregistersMCPToolsAfterRun guards against the run's MCP tools leaking onto the shared session:
// a discovered tool is callable during the run, but once Run returns it is gone from the session's per-run
// registry, so a later sub-agent or run on the same (still live) session cannot reach it.
func TestDynamicAgentUnregistersMCPToolsAfterRun(t *testing.T) {
	m := newMockMCPServer(t, echoBuild)
	defer pointClientAt(m)()

	script := &scriptedAIService{
		turns: []turnFunc{
			toolCallTurn("mcp_demo_echo", map[string]any{"message": "hola"}),
			textTurn("Done."),
		},
	}
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})

	// The tool is not registered before the run.
	_, ok := s.Tool("mcp_demo_echo")
	require.False(t, ok)

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "You may use the echo tool.",
		MCPConnectors: []ai.MCPConnector{{
			Name:              "demo",
			URL:               "https://mcp.example.test/mcp",
			AuthSecret:        "MCP_TOKEN",
			TrustReadOnlyHint: true, // echo is read-only: execute it inline instead of capturing it as a proposal.
		}},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Echo hola.")
	require.NoError(t, err)
	require.Equal(t, "Done.", res.Response)

	// It executed during the run (so it really was registered), but the session — still alive — no longer exposes
	// it: the per-run registration was torn down when Run returned.
	require.Len(t, s.Messages(ai.FilterByType(ai.MessageTypeResult), ai.FilterByTool("mcp_demo_echo")), 1)
	_, ok = s.Tool("mcp_demo_echo")
	require.False(t, ok, "MCP tool must not survive the run on the shared session")
}

// TestDynamicAgentCleansUpMCPToolsOnPartialFailure verifies the fail-closed cleanup: when a later connector fails
// setup after an earlier one already registered its tools, the failed run leaves no residual tools behind.
func TestDynamicAgentCleansUpMCPToolsOnPartialFailure(t *testing.T) {
	m := newMockMCPServer(t, echoBuild)
	defer pointClientAt(m)()

	script := &scriptedAIService{}
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "You may use the echo tool.",
		MCPConnectors: []ai.MCPConnector{
			// The first connector resolves and registers mcp.demo.echo...
			{Name: "demo", URL: "https://mcp.example.test/mcp", AuthSecret: "MCP_TOKEN"},
			// ...the second fails fast because its referenced secret is not set, aborting the run.
			{Name: "broken", URL: "https://mcp.example.test/mcp", AuthSecret: "MISSING_TOKEN"},
		},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Echo hola.")
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "broken")

	// The model was never invoked (setup failed before Complete), and the first connector's tool was cleaned up.
	script.mu.Lock()
	calls := script.calls
	script.mu.Unlock()
	require.Equal(t, 0, calls)
	_, ok := s.Tool("mcp_demo_echo")
	require.False(t, ok, "a failed run must not leave partially-registered MCP tools behind")
}

// TestDynamicAgentProposesWriteTool verifies the governed-action capture: when the model calls a write (non-read-only)
// MCP tool, the loop captures it as the agent's ProposedAction and returns a pending-approval placeholder rather than
// executing it. The run ends with the proposal surfaced, for the durable workflow to route through human approval.
func TestDynamicAgentProposesWriteTool(t *testing.T) {
	m := newMockMCPServer(t, createIssueBuild)
	rec := &mcpCallRecorder{}
	defer pointRecordingClientAt(m, rec)()

	script := &scriptedAIService{
		turns: []turnFunc{
			// One turn: the model calls the governed write and the loop pauses on it. No second turn runs, because the
			// segment does not continue past the proposal (nothing is said before approval).
			toolCallTurn("mcp_demo_create_issue", map[string]any{"project": "OPS", "summary": "Disk full"}),
		},
	}
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "Create an issue when asked.",
		MCPConnectors: []ai.MCPConnector{{
			Name:       "demo",
			URL:        "https://mcp.example.test/mcp",
			AuthSecret: "MCP_TOKEN",
			// No TrustReadOnlyHint: every tool on this connector is a governed write.
		}},
		MaxSteps: 5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "mcp_agent", "Open a ticket: disk full on OPS.")
	require.NoError(t, err)
	require.Empty(t, res.Response, "the segment pauses at the governed write, answering nothing before approval")

	// The write tool did NOT execute: no CallTool ever reached the client (only discovery's ListTools did).
	require.Empty(t, rec.calls(), "a proposed write must not be executed inline")

	// The proposal was captured with the dotted effective name and the model's arguments.
	require.NotNil(t, res.Proposed)
	require.Equal(t, "demo", res.Proposed.Connector)
	require.Equal(t, "mcp.demo.create_issue", res.Proposed.Tool)
	require.Equal(t, "OPS", res.Proposed.Args["project"])
	require.Equal(t, "Disk full", res.Proposed.Args["summary"])
	require.Contains(t, res.Proposed.Summary, "create_issue")
	require.NotEmpty(t, res.Proposed.SchemaHash, "the proposal pins the tool schema discovered at run start")
}

// TestDynamicAgentResumesAfterAction is the core of the segmented durable loop: a first segment proposes a governed
// write and pauses; a second segment resumes with the executed action's result injected, and the model — now seeing
// the outcome of the action it proposed — composes a closing answer without re-executing the write. The second segment
// runs on the same session, and reconstructs the conversation from the session's persisted message tree (in production
// the workflow reopens that session by ID across the pause), so the model keeps full context without threading it.
func TestDynamicAgentResumesAfterAction(t *testing.T) {
	m := newMockMCPServer(t, createIssueBuild)
	rec := &mcpCallRecorder{}
	defer pointRecordingClientAt(m, rec)()

	script := &scriptedAIService{
		turns: []turnFunc{
			// Segment 1: propose the write; the loop pauses on it without a further turn (nothing is said before approval).
			toolCallTurn("mcp_demo_create_issue", map[string]any{"project": "OPS", "summary": "Disk full"}),
			// Segment 2 (resume): having seen the executed result injected, close.
			textTurn("Done! I created the issue: PROJ-42."),
		},
	}
	s := newMCPScriptedSession(t, script, map[string]string{"MCP_TOKEN": "s3cr3t"})
	snap := &ai.AgentSnapshot{
		Name:         "mcp_agent",
		Instructions: "Create an issue when asked.",
		MCPConnectors: []ai.MCPConnector{{
			Name:       "demo",
			URL:        "https://mcp.example.test/mcp",
			AuthSecret: "MCP_TOKEN",
			// No TrustReadOnlyHint: create_issue is a governed write.
		}},
		MaxSteps: 5,
	}

	// Segment 1: the model proposes the write; the run pauses with the proposal captured and no inline effect.
	seg1, err := (&ai.DynamicAgent{Snapshot: snap}).Run(t.Context(), s, "Open a ticket: disk full on OPS.", nil)
	require.NoError(t, err)
	require.NotNil(t, seg1.Proposed, "segment 1 pauses on a governed write")
	require.Equal(t, "mcp.demo.create_issue", seg1.Proposed.Tool)
	require.NotEmpty(t, seg1.Proposed.ToolCallID, "the proposal carries the model's tool-call identity")
	require.Empty(t, seg1.Response, "segment 1 pauses at the tool call and answers nothing before approval")
	require.Empty(t, rec.calls(), "a proposed write is never executed inline")

	// Segment 2: resume on the same session with the executed result injected. A fresh DynamicAgent (its captured-proposal
	// state starts empty) reconstructs the conversation from the session's tree, so the model sees the outcome and closes,
	// proposing nothing further.
	seg2, err := (&ai.DynamicAgent{Snapshot: snap}).Run(t.Context(), s, "", &ai.InjectedResult{
		ToolCallID: seg1.Proposed.ToolCallID,
		Tool:       seg1.Proposed.Tool,
		Message:    "created PROJ-42",
	})
	require.NoError(t, err)
	require.Nil(t, seg2.Proposed, "the resumed segment closes without proposing another action")
	require.Equal(t, "Done! I created the issue: PROJ-42.", seg2.Response)
	require.Empty(t, rec.calls(), "resuming must not re-execute the write inline")

	// The executed result was injected into the run's conversation (so the model saw it and the UI can show it).
	var sawInjected bool
	for _, msg := range s.Messages(ai.FilterByType(ai.MessageTypeText)) {
		if strings.Contains(msg.Content, "created PROJ-42") {
			sawInjected = true
		}
	}
	require.True(t, sawInjected, "the resume injects the executed action's result as a turn in the session")
}
