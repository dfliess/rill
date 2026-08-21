package ai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/email"
	"github.com/rilldata/rill/runtime/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"

	_ "github.com/rilldata/rill/runtime/drivers/admin"
	_ "github.com/rilldata/rill/runtime/drivers/duckdb"
	_ "github.com/rilldata/rill/runtime/drivers/file"
	_ "github.com/rilldata/rill/runtime/drivers/openai"
	_ "github.com/rilldata/rill/runtime/drivers/sqlite"
	_ "github.com/rilldata/rill/runtime/reconcilers"
	_ "github.com/rilldata/rill/runtime/resolvers"
)

type restartChatRequest struct {
	Authorization string
	Body          map[string]any
}

// restartChatServer emulates the DeepSeek failure contract that motivated this regression test: replaying an
// assistant tool call without provider-private reasoning is rejected unless thinking was disabled on the request.
// The first turn performs an inline read, the second proposes a governed MCP write, and the third closes after the
// approved result is injected. This matches the real failure shape: more than one reasoning-bearing tool turn exists
// before the process restart, and reopening must preserve both the original user prompt and the injected action result
// as user messages.
type restartChatServer struct {
	server           *httptest.Server
	expectedThinking string
	mu               sync.Mutex
	reqs             []restartChatRequest
}

const (
	restartReasoningContent      = "durable provider reasoning for the read: keep this exact string across the restart"
	restartWriteReasoningContent = "durable provider reasoning for the write: keep this exact string across the restart"
	restartReadToolName          = "mcp_demo_inspect_disk"
	restartWriteToolName         = "mcp_demo_create_issue"
)

func newRestartChatServer(t *testing.T, expectedThinking ...string) *restartChatServer {
	t.Helper()
	s := &restartChatServer{}
	if len(expectedThinking) > 0 {
		s.expectedThinking = expectedThinking[0]
	}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v1/chat/completions" {
			http.Error(w, "unexpected path "+req.URL.Path, http.StatusNotFound)
			return
		}
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON request: "+err.Error(), http.StatusBadRequest)
			return
		}

		s.mu.Lock()
		turn := len(s.reqs)
		s.reqs = append(s.reqs, restartChatRequest{Authorization: req.Header.Get("Authorization"), Body: body})
		s.mu.Unlock()

		thinkingType := ""
		if thinking, ok := body["thinking"].(map[string]any); ok {
			thinkingType, _ = thinking["type"].(string)
		}
		if s.expectedThinking != "" && thinkingType != s.expectedThinking {
			http.Error(w, fmt.Sprintf("expected thinking %q, got %q", s.expectedThinking, thinkingType), http.StatusBadRequest)
			return
		}
		if err := validateRestartRequest(body, turn, thinkingType); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, `{"error":{"message":%q}}`, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch turn {
		case 0:
			reasoningField := ""
			if thinkingType != "disabled" {
				reasoningField = fmt.Sprintf(`"reasoning_content":%q,`, restartReasoningContent)
			}
			_, _ = fmt.Fprintf(w, `{
				"id":"chatcmpl-read","object":"chat.completion","created":1,"model":"agent-model-v1",
				"choices":[{"index":0,"message":{"role":"assistant","content":"",%s"tool_calls":[{
					"id":"provider-call-read","type":"function","function":{"name":"mcp_demo_inspect_disk","arguments":"{\"host\":\"db-01\"}"}
				}]} ,"finish_reason":"tool_calls"}],
				"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}
			}`, reasoningField)
		case 1:
			reasoningField := ""
			if thinkingType != "disabled" {
				reasoningField = fmt.Sprintf(`"reasoning_content":%q,`, restartWriteReasoningContent)
			}
			_, _ = fmt.Fprintf(w, `{
				"id":"chatcmpl-proposal","object":"chat.completion","created":2,"model":"agent-model-v1",
				"choices":[{"index":0,"message":{"role":"assistant","content":"",%s"tool_calls":[{
					"id":"provider-call-write","type":"function","function":{"name":"mcp_demo_create_issue","arguments":"{\"project\":\"OPS\",\"summary\":\"Disk full\"}"}
				}]} ,"finish_reason":"tool_calls"}],
				"usage":{"prompt_tokens":12,"completion_tokens":6,"total_tokens":18}
			}`, reasoningField)
		default:
			_, _ = io.WriteString(w, `{
				"id":"chatcmpl-final","object":"chat.completion","created":3,"model":"agent-model-v1",
				"choices":[{"index":0,"message":{"role":"assistant","content":"Done: created OPS-42."},"finish_reason":"stop"}],
				"usage":{"prompt_tokens":14,"completion_tokens":4,"total_tokens":18}
			}`)
		}
	}))
	t.Cleanup(s.server.Close)
	return s
}

func (s *restartChatServer) url() string { return s.server.URL + "/v1" }

func (s *restartChatServer) requests() []restartChatRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]restartChatRequest(nil), s.reqs...)
}

type restartAssistantToolTurn struct {
	tool             string
	reasoning        string
	hasReasoning     bool
	hasStringContent bool
}

func restartRequestRoles(body map[string]any) []string {
	messages, _ := body["messages"].([]any)
	roles := make([]string, 0, len(messages))
	for _, raw := range messages {
		msg, _ := raw.(map[string]any)
		role, _ := msg["role"].(string)
		roles = append(roles, role)
	}
	return roles
}

func restartAssistantToolTurns(body map[string]any) []restartAssistantToolTurn {
	messages, _ := body["messages"].([]any)
	var turns []restartAssistantToolTurn
	for _, raw := range messages {
		msg, _ := raw.(map[string]any)
		if msg["role"] != "assistant" {
			continue
		}
		calls, ok := msg["tool_calls"].([]any)
		if !ok || len(calls) == 0 {
			continue
		}
		call, _ := calls[0].(map[string]any)
		function, _ := call["function"].(map[string]any)
		tool, _ := function["name"].(string)
		reasoningRaw, hasReasoning := msg["reasoning_content"]
		reasoning, _ := reasoningRaw.(string)
		_, hasStringContent := msg["content"].(string)
		turns = append(turns, restartAssistantToolTurn{
			tool: tool, reasoning: reasoning, hasReasoning: hasReasoning, hasStringContent: hasStringContent,
		})
	}
	return turns
}

func validateRestartRequest(body map[string]any, turn int, thinking string) error {
	expectedRoles := [][]string{
		{"system", "user"},
		{"system", "user", "assistant", "tool"},
		{"system", "user", "assistant", "tool", "assistant", "tool", "user"},
	}
	if turn >= len(expectedRoles) {
		return fmt.Errorf("unexpected model request %d", turn+1)
	}
	roles := restartRequestRoles(body)
	if !slices.Equal(roles, expectedRoles[turn]) {
		return fmt.Errorf("invalid message roles on request %d: got %v, want %v", turn+1, roles, expectedRoles[turn])
	}

	toolTurns := restartAssistantToolTurns(body)
	expectedTools := []string{restartReadToolName, restartWriteToolName}
	expectedReasoning := []string{restartReasoningContent, restartWriteReasoningContent}
	if len(toolTurns) != turn {
		return fmt.Errorf("invalid assistant tool-turn count on request %d: got %d, want %d", turn+1, len(toolTurns), turn)
	}
	for i, toolTurn := range toolTurns {
		if toolTurn.tool != expectedTools[i] {
			return fmt.Errorf("invalid tool on assistant turn %d: got %q, want %q", i+1, toolTurn.tool, expectedTools[i])
		}
		if !toolTurn.hasStringContent {
			return fmt.Errorf("assistant tool turn %d must have non-null string content", i+1)
		}
		if thinking == "disabled" {
			if toolTurn.hasReasoning {
				return fmt.Errorf("assistant tool turn %d unexpectedly replayed reasoning_content with thinking disabled", i+1)
			}
			continue
		}
		if !toolTurn.hasReasoning || toolTurn.reasoning != expectedReasoning[i] {
			return fmt.Errorf("The reasoning_content in thinking mode must be passed back for assistant tool turn %d", i+1)
		}
	}
	return nil
}

type restartMCPCreateIssueIn struct {
	Project string `json:"project" jsonschema:"project key"`
	Summary string `json:"summary" jsonschema:"issue summary"`
}

type restartMCPCreateIssueOut struct {
	Key string `json:"key"`
}

type restartMCPInspectDiskIn struct {
	Host string `json:"host" jsonschema:"host name"`
}

type restartMCPInspectDiskOut struct {
	UsagePercent int `json:"usage_percent"`
}

func newRestartMCPServer(t *testing.T) *httptest.Server {
	t.Helper()
	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		srv := mcp.NewServer(&mcp.Implementation{Name: "restart-mcp", Version: "1.0.0"}, nil)
		mcp.AddTool(srv, &mcp.Tool{
			Name: "inspect_disk", Description: "Inspect disk usage.",
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
		}, func(_ context.Context, _ *mcp.CallToolRequest, _ restartMCPInspectDiskIn) (*mcp.CallToolResult, restartMCPInspectDiskOut, error) {
			return nil, restartMCPInspectDiskOut{UsagePercent: 99}, nil
		})
		mcp.AddTool(srv, &mcp.Tool{Name: "create_issue", Description: "Create an issue."},
			func(_ context.Context, _ *mcp.CallToolRequest, _ restartMCPCreateIssueIn) (*mcp.CallToolResult, restartMCPCreateIssueOut, error) {
				return nil, restartMCPCreateIssueOut{Key: "OPS-42"}, nil
			})
		return srv
	}, &mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true, DisableLocalhostProtection: true})
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func restartConnectorYAML(baseURL, model, thinking string) string {
	return fmt.Sprintf(`
type: connector
driver: openai
api_key: project-value-overridden-by-instance
base_url: %q
model: %q
structured_output_mode: json_object
extra_body:
  thinking:
    type: %s
`, baseURL, model, thinking)
}

func restartAgentYAML(connector, model string) string {
	return fmt.Sprintf(`
type: agent
display_name: Durable DeepSeek agent
instructions: "Create an issue when requested, then report the result."
model:
  connector: %q
  name: %q
tools: []
mcp:
  - name: demo
    url: https://mcp.example.test/mcp
    approval: manual
    trust_read_only_hint: true
limits:
  max_steps: 5
`, connector, model)
}

func writeRestartProjectFile(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o644))
}

func newPersistentActRuntime(t *testing.T, metastoreDSN, storageDir string) *runtime.Runtime {
	t.Helper()
	metastoreConfig, err := structpb.NewStruct(map[string]any{"dsn": metastoreDSN})
	require.NoError(t, err)
	opts := &runtime.Options{
		MetastoreConnector: "metastore",
		SystemConnectors: []*runtimev1.Connector{{
			Type: "sqlite", Name: "metastore", Config: metastoreConfig,
		}},
		ConnectionCacheSize:          100,
		QueryCacheSizeBytes:          int64(100 * datasize.MB),
		SecurityEngineCacheSize:      100,
		ControllerLogBufferCapacity:  1000,
		ControllerLogBufferSizeBytes: int64(4 * datasize.MB),
		AllowHostAccess:              true,
	}
	rt, err := runtime.New(t.Context(), opts, zap.NewNop(), storage.MustNew(storageDir, nil), activity.NewNoopClient(), email.New(email.NewTestSender()))
	require.NoError(t, err)
	return rt
}

func requireRestartPostgres(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("ACT_TEST_DBOS_URL")
	if dsn == "" {
		dsn = "postgres://postgres:dbos@localhost:55432/dbos_spike?sslmode=disable"
	}
	u, err := url.Parse(dsn)
	require.NoError(t, err)
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		if os.Getenv("ACT_TEST_REQUIRE_POSTGRES") == "1" {
			t.Fatalf("required Act restart Postgres is unreachable at %s: %v", u.Host, err)
		}
		t.Skipf("Act restart integration needs Postgres at %s: %v", u.Host, err)
	}
	require.NoError(t, conn.Close())
	return dsn
}

func newRestartRunStore(t *testing.T, dsn string) *act.PostgresRunStore {
	t.Helper()
	schema := "act_restart_" + uuid.NewString()[:8]
	store, err := act.NewPostgresRunStore(t.Context(), act.StoreConfig{DatabaseURL: dsn, Schema: schema})
	require.NoError(t, err)
	require.NoError(t, store.Migrate(t.Context()))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		store.DropSchemaForTest(ctx)
		store.Close()
	})
	return store
}

type restartActionExecutor struct {
	mu    sync.Mutex
	calls int
}

func (e *restartActionExecutor) Execute(_ context.Context, _ act.ExecuteRequest) (act.ExecuteResult, error) {
	e.mu.Lock()
	e.calls++
	e.mu.Unlock()
	return act.ExecuteResult{
		Outcome: act.OutcomeSucceeded, ExternalReference: "OPS-42", Message: "created OPS-42",
	}, nil
}

func (e *restartActionExecutor) Verify(context.Context, act.VerifyRequest) (act.VerifyResult, error) {
	return act.VerifyResult{}, nil
}

func (e *restartActionExecutor) executeCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.calls
}

var _ act.ActionExecutor = (*restartActionExecutor)(nil)

func newRestartGateway(store *act.PostgresRunStore, executor act.ActionExecutor) *act.Gateway {
	return &act.Gateway{
		Registry: act.NewMapToolRegistry(act.ToolDescriptor{
			Name: "mcp.demo.create_issue", Connector: "demo", Version: "1",
			Class: act.ClassIdempotentNative,
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"project": {Type: "string"},
					"summary": {Type: "string"},
				},
				Required: []string{"project", "summary"},
			},
		}),
		Ledger: store, Executor: executor,
	}
}

// TestOpenAIConnectorConfigSurvivesApprovalAndRuntimeRestart exercises the production-shaped path as one composed
// regression: catalog snapshot -> per-agent session -> OpenAI request -> MCP discovery/write capture -> DBOS durable
// approval wait -> worker and Runtime shutdown -> catalog/session reopen -> approved action -> second OpenAI request.
// It deliberately edits the live agent and connector while parked. The resumed segment must keep the checkpointed
// endpoint/model/thinking posture, use only the rotated secret, and never fall through to the instance default.
func TestOpenAIConnectorConfigSurvivesApprovalAndRuntimeRestart(t *testing.T) {
	dsn := requireRestartPostgres(t)
	dbosSchema := os.Getenv("ACT_TEST_DBOS_SCHEMA")
	if dbosSchema == "" {
		dbosSchema = "act_dbos"
	}
	store := newRestartRunStore(t, dsn)
	provider := newRestartChatServer(t)
	forbidden := newRestartChatServer(t)
	mcpServer := newRestartMCPServer(t)
	restoreMCP := ai.SetNewMCPClientForTest(func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error) {
		local := &mcpconn.ConnectorConfig{Name: cfg.Name, Transport: mcpconn.TransportStreamableHTTP, URL: mcpServer.URL}
		return mcpconn.NewClient(local, mcpconn.Options{BearerToken: token, AllowPrivateNetworks: true})
	})
	defer restoreMCP()

	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	storageDir := filepath.Join(root, "storage")
	metastoreDSN := "file:" + filepath.Join(root, "metastore.sqlite") + "?cache=shared"
	catalogDSN := "file:" + filepath.Join(root, "catalog.sqlite") + "?cache=shared"
	require.NoError(t, os.MkdirAll(repoDir, 0o755))
	writeRestartProjectFile(t, repoDir, "rill.yaml", "")
	writeRestartProjectFile(t, repoDir, "connectors/deepseek.yaml", restartConnectorYAML(provider.url(), "connector-model-v1", "disabled"))
	writeRestartProjectFile(t, repoDir, "connectors/default_ai.yaml", restartConnectorYAML(forbidden.url(), "default-model", "enabled"))
	writeRestartProjectFile(t, repoDir, "triage.yaml", restartAgentYAML("deepseek", "agent-model-v1"))

	repoConfig, err := structpb.NewStruct(map[string]any{"dsn": repoDir})
	require.NoError(t, err)
	olapConfig, err := structpb.NewStruct(map[string]any{"dsn": ":memory:", "mode": "readwrite"})
	require.NoError(t, err)
	catalogConfig, err := structpb.NewStruct(map[string]any{"dsn": catalogDSN})
	require.NoError(t, err)

	rt1 := newPersistentActRuntime(t, metastoreDSN, storageDir)
	rt1Closed := false
	t.Cleanup(func() {
		if !rt1Closed {
			_ = rt1.Close()
		}
	})
	instance := &drivers.Instance{
		ID:               "restart-" + uuid.NewString(),
		Environment:      "test",
		OLAPConnector:    "duckdb",
		RepoConnector:    "repo",
		AIConnector:      "default_ai",
		CatalogConnector: "catalog",
		AdminConnector:   "noop_admin",
		Connectors: []*runtimev1.Connector{
			{Type: "file", Name: "repo", Config: repoConfig},
			{Type: "duckdb", Name: "duckdb", Config: olapConfig},
			{Type: "sqlite", Name: "catalog", Config: catalogConfig},
		},
		Variables: map[string]string{
			"connector.deepseek.api_key":   "initial-secret",
			"connector.default_ai.api_key": "forbidden-secret",
		},
	}
	require.NoError(t, rt1.CreateInstance(t.Context(), instance))
	require.NoError(t, rt1.WaitUntilIdle(t.Context(), instance.ID, false))

	initialProvider := act.NewCatalogAgentProvider(rt1)
	snapshot, err := initialProvider.GetAgent(t.Context(), instance.ID, "triage")
	require.NoError(t, err)
	require.Equal(t, "deepseek", snapshot.ModelConnector)
	require.Equal(t, "agent-model-v1", snapshot.ModelProperties["model"])
	require.NotContains(t, snapshot.ModelProperties, "api_key")
	snapshotJSON, err := json.Marshal(snapshot)
	require.NoError(t, err)
	require.NotContains(t, string(snapshotJSON), "initial-secret")

	actionExecutor := &restartActionExecutor{}
	gateway := newRestartGateway(store, actionExecutor)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	applicationVersion := "openai-restart-" + uuid.NewString()
	runner1 := &act.SessionRunner{
		Provider: initialProvider,
		Sessions: act.NewSessionFactory(rt1, activity.NewNoopClient()),
		Actions:  act.NewDenyActionSink(),
	}
	executor1, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: dbosSchema, ApplicationVersion: applicationVersion,
		Runner: runner1, Store: store, Gateway: gateway, Proposer: act.NewCapturedProposer(), Logger: logger,
	})
	require.NoError(t, err)
	executor1Closed := false
	t.Cleanup(func() {
		if !executor1Closed {
			executor1.Close(10 * time.Second)
		}
	})

	runID, err := executor1.Start(t.Context(), act.AgentRunInput{
		InstanceID: instance.ID, AgentName: "triage", Prompt: "Open a ticket for the full disk.",
		IdempotencyKey: "restart-" + uuid.NewString(),
		Actor:          act.Actor{Subject: "user:alice", Claims: &runtime.SecurityClaims{UserID: "alice", SkipChecks: true}},
	})
	require.NoError(t, err)

	var approval *act.Approval
	require.Eventually(t, func() bool {
		approvals, listErr := store.ListApprovals(t.Context(), act.ListApprovalsFilter{
			InstanceID: instance.ID, RunID: runID, Status: act.ApprovalStatusPending,
		})
		if listErr != nil || len(approvals) != 1 {
			return false
		}
		approval = approvals[0]
		return true
	}, 20*time.Second, 50*time.Millisecond)
	require.Len(t, provider.requests(), 2)
	require.Empty(t, forbidden.requests())

	// Drift every non-secret source while the run is parked, but rotate the credential that is intentionally live.
	writeRestartProjectFile(t, repoDir, "connectors/deepseek.yaml", restartConnectorYAML(forbidden.url(), "connector-model-v2", "enabled"))
	writeRestartProjectFile(t, repoDir, "triage.yaml", restartAgentYAML("default_ai", "agent-model-v2"))
	current, err := rt1.Instance(t.Context(), instance.ID)
	require.NoError(t, err)
	edited := *current
	edited.Variables = maps.Clone(current.Variables)
	edited.Variables["connector.deepseek.api_key"] = "rotated-secret"
	require.NoError(t, rt1.EditInstance(t.Context(), &edited, false))

	// A real lifecycle boundary: both the durable worker and Rill Runtime are destroyed before approval arrives.
	executor1.Close(10 * time.Second)
	executor1Closed = true
	require.NoError(t, rt1.Close())
	rt1Closed = true

	rt2 := newPersistentActRuntime(t, metastoreDSN, storageDir)
	t.Cleanup(func() { _ = rt2.Close() })
	require.NoError(t, rt2.WaitUntilIdle(t.Context(), instance.ID, false))
	liveProvider := act.NewCatalogAgentProvider(rt2)
	liveSnapshot, err := liveProvider.GetAgent(t.Context(), instance.ID, "triage")
	require.NoError(t, err)
	require.Equal(t, "default_ai", liveSnapshot.ModelConnector, "the project edit must be live after restart")
	require.Equal(t, "agent-model-v2", liveSnapshot.ModelProperties["model"])

	runner2 := &act.SessionRunner{
		Provider: liveProvider,
		Sessions: act.NewSessionFactory(rt2, activity.NewNoopClient()),
		Actions:  act.NewDenyActionSink(),
	}
	executor2, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: dbosSchema, ApplicationVersion: applicationVersion,
		Runner: runner2, Store: store, Gateway: gateway, Proposer: act.NewCapturedProposer(), Logger: logger,
	})
	require.NoError(t, err)
	t.Cleanup(func() { executor2.Close(10 * time.Second) })

	_, err = store.ResolveApproval(t.Context(), instance.ID, approval.ApprovalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, executor2.Resume(t.Context(), runID, act.ApprovalApproved, approval.ToolCallID))
	result, err := executor2.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, result.Status)
	require.Equal(t, "Done: created OPS-42.", result.Response)
	require.True(t, result.ActionTaken)
	require.Equal(t, "OPS-42", result.ActionRef)
	require.Equal(t, 1, actionExecutor.executeCount())

	reqs := provider.requests()
	require.Len(t, reqs, 3)
	require.Equal(t, "Bearer initial-secret", reqs[0].Authorization)
	require.Equal(t, "Bearer initial-secret", reqs[1].Authorization)
	require.Equal(t, "Bearer rotated-secret", reqs[2].Authorization)
	for turn, req := range reqs {
		require.Equal(t, "agent-model-v1", req.Body["model"])
		require.Equal(t, map[string]any{"type": "disabled"}, req.Body["thinking"])
		require.NoError(t, validateRestartRequest(req.Body, turn, "disabled"))
	}
	require.Empty(t, forbidden.requests(), "neither the default nor drifted endpoint may serve the checkpointed run")
}

const restartHelperPhaseEnv = "RILL_TEST_OPENAI_RESTART_PHASE"

type restartSubprocessResult struct {
	Status       act.RunStatus `json:"status"`
	Response     string        `json:"response"`
	ActionRef    string        `json:"action_ref"`
	ExecuteCount int           `json:"execute_count"`
}

func restartTemplatedConnectorYAML(baseURL, model, thinking string) string {
	return fmt.Sprintf(`
type: connector
driver: openai
api_key: "{{ .env.DEEPSEEK_API_KEY }}"
base_url: %q
model: %q
structured_output_mode: json_object
extra_body:
  thinking:
    type: %s
`, baseURL, model, thinking)
}

func openRestartStore(t *testing.T, dsn, schema string) *act.PostgresRunStore {
	t.Helper()
	store, err := act.NewPostgresRunStore(t.Context(), act.StoreConfig{DatabaseURL: dsn, Schema: schema})
	require.NoError(t, err)
	require.NoError(t, store.Migrate(t.Context()))
	return store
}

func restartInstanceFromEnv(t *testing.T, rt *runtime.Runtime) *drivers.Instance {
	t.Helper()
	repoConfig, err := structpb.NewStruct(map[string]any{"dsn": os.Getenv("RILL_TEST_RESTART_REPO")})
	require.NoError(t, err)
	olapConfig, err := structpb.NewStruct(map[string]any{"dsn": ":memory:", "mode": "readwrite"})
	require.NoError(t, err)
	catalogConfig, err := structpb.NewStruct(map[string]any{"dsn": os.Getenv("RILL_TEST_RESTART_CATALOG_DSN")})
	require.NoError(t, err)
	return &drivers.Instance{
		ID:               os.Getenv("RILL_TEST_RESTART_INSTANCE"),
		Environment:      "test",
		OLAPConnector:    "duckdb",
		RepoConnector:    "repo",
		AIConnector:      "default_ai",
		CatalogConnector: "catalog",
		AdminConnector:   "noop_admin",
		Connectors: []*runtimev1.Connector{
			{Type: "file", Name: "repo", Config: repoConfig},
			{Type: "duckdb", Name: "duckdb", Config: olapConfig},
			{Type: "sqlite", Name: "catalog", Config: catalogConfig},
		},
	}
}

func runOpenAIRestartHelper(t *testing.T, phase string) {
	t.Helper()
	mcpURL := os.Getenv("RILL_TEST_RESTART_MCP_URL")
	restoreMCP := ai.SetNewMCPClientForTest(func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error) {
		local := &mcpconn.ConnectorConfig{Name: cfg.Name, Transport: mcpconn.TransportStreamableHTTP, URL: mcpURL}
		return mcpconn.NewClient(local, mcpconn.Options{BearerToken: token, AllowPrivateNetworks: true})
	})
	defer restoreMCP()

	dsn := os.Getenv("ACT_TEST_DBOS_URL")
	store := openRestartStore(t, dsn, os.Getenv("RILL_TEST_RESTART_STORE_SCHEMA"))
	defer store.Close()
	rt := newPersistentActRuntime(t, os.Getenv("RILL_TEST_RESTART_METASTORE_DSN"), os.Getenv("RILL_TEST_RESTART_STORAGE"))

	instanceID := os.Getenv("RILL_TEST_RESTART_INSTANCE")
	if phase == "phase1" {
		require.NoError(t, rt.CreateInstance(t.Context(), restartInstanceFromEnv(t, rt)))
	}
	require.NoError(t, rt.WaitUntilIdle(t.Context(), instanceID, false))

	actionExecutor := &restartActionExecutor{}
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: act.NewSessionFactory(rt, activity.NewNoopClient()),
		Actions:  act.NewDenyActionSink(),
	}
	executor, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: os.Getenv("RILL_TEST_RESTART_DBOS_SCHEMA"),
		ApplicationVersion: os.Getenv("RILL_TEST_RESTART_APP_VERSION"),
		Runner:             runner, Store: store, Gateway: newRestartGateway(store, actionExecutor),
		Proposer: act.NewCapturedProposer(), Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)

	runID := act.ComposeRunID(instanceID, "triage", os.Getenv("RILL_TEST_RESTART_IDEMPOTENCY"))
	if phase == "phase1" {
		started, startErr := executor.Start(t.Context(), act.AgentRunInput{
			InstanceID: instanceID, AgentName: "triage", Prompt: "Open a ticket for the full disk.",
			IdempotencyKey: os.Getenv("RILL_TEST_RESTART_IDEMPOTENCY"),
			Actor:          act.Actor{Subject: "user:alice", Claims: &runtime.SecurityClaims{UserID: "alice", SkipChecks: true}},
		})
		require.NoError(t, startErr)
		require.Equal(t, runID, started)
		require.Eventually(t, func() bool {
			approvals, listErr := store.ListApprovals(t.Context(), act.ListApprovalsFilter{
				InstanceID: instanceID, RunID: runID, Status: act.ApprovalStatusPending,
			})
			return listErr == nil && len(approvals) == 1
		}, 20*time.Second, 50*time.Millisecond)

		// This is deliberately not a graceful Close: no deferred cleanup runs and all package globals disappear with
		// the process. The parent starts phase2 in a fresh test binary against the same DBOS/catalog files.
		os.Exit(0)
	}

	require.Equal(t, "phase2", phase)
	defer func() {
		executor.Close(10 * time.Second)
		require.NoError(t, rt.Close())
	}()
	var approval *act.Approval
	require.Eventually(t, func() bool {
		approvals, listErr := store.ListApprovals(t.Context(), act.ListApprovalsFilter{
			InstanceID: instanceID, RunID: runID, Status: act.ApprovalStatusPending,
		})
		if listErr != nil || len(approvals) != 1 {
			return false
		}
		approval = approvals[0]
		return true
	}, 20*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(t.Context(), instanceID, approval.ApprovalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, executor.Resume(t.Context(), runID, act.ApprovalApproved, approval.ToolCallID))
	result, err := executor.Result(runID)
	require.NoError(t, err)

	encoded, err := json.Marshal(restartSubprocessResult{
		Status: result.Status, Response: result.Response, ActionRef: result.ActionRef,
		ExecuteCount: actionExecutor.executeCount(),
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(os.Getenv("RILL_TEST_RESTART_RESULT"), encoded, 0o600))
}

// TestOpenAIRestartSubprocessHelper is invoked only by TestOpenAIConnectorConfigSurvivesProcessRestart. Phase1 exits
// the process abruptly while DBOS is waiting; phase2 is a separately launched binary that recovers and approves it.
func TestOpenAIRestartSubprocessHelper(t *testing.T) {
	phase := os.Getenv(restartHelperPhaseEnv)
	if phase == "" {
		t.Skip("subprocess helper")
	}
	runOpenAIRestartHelper(t, phase)
}

func runRestartSubprocess(t *testing.T, phase string, environment []string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestOpenAIRestartSubprocessHelper$", "-test.v", "-test.timeout=55s")
	cmd.Env = append(os.Environ(), environment...)
	cmd.Env = append(cmd.Env, restartHelperPhaseEnv+"="+phase)
	output, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "%s subprocess failed:\n%s", phase, output)
}

// TestOpenAIConnectorConfigSurvivesProcessRestart repeats the composed regression across two OS processes. It is the
// guard against the original class of bug: any future cache or provider state held only in package globals is gone
// before the approval is signed, while DBOS, the agent snapshot and the catalog conversation remain.
func TestOpenAIConnectorConfigSurvivesProcessRestart(t *testing.T) {
	for _, thinking := range []string{"disabled", "enabled"} {
		t.Run(thinking, func(t *testing.T) {
			testOpenAIConnectorConfigSurvivesProcessRestart(t, thinking)
		})
	}
}

func testOpenAIConnectorConfigSurvivesProcessRestart(t *testing.T, thinking string) {
	dsn := requireRestartPostgres(t)
	provider := newRestartChatServer(t, thinking)
	forbidden := newRestartChatServer(t)
	mcpServer := newRestartMCPServer(t)

	root := t.TempDir()
	repoDir := filepath.Join(root, "repo")
	require.NoError(t, os.MkdirAll(repoDir, 0o755))
	writeRestartProjectFile(t, repoDir, "rill.yaml", "")
	writeRestartProjectFile(t, repoDir, ".env", "DEEPSEEK_API_KEY=initial-secret\n")
	writeRestartProjectFile(t, repoDir, "connectors/deepseek.yaml", restartTemplatedConnectorYAML(provider.url(), "connector-model-v1", thinking))
	writeRestartProjectFile(t, repoDir, "connectors/default_ai.yaml", restartTemplatedConnectorYAML(forbidden.url(), "default-model", "enabled"))
	writeRestartProjectFile(t, repoDir, "triage.yaml", restartAgentYAML("deepseek", "agent-model-v1"))

	instanceID := "restart-process-" + uuid.NewString()
	storeSchema := "act_restart_process_" + uuid.NewString()[:8]
	store := openRestartStore(t, dsn, storeSchema)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		store.DropSchemaForTest(ctx)
		store.Close()
	})

	dbosSchema := os.Getenv("ACT_TEST_DBOS_SCHEMA")
	if dbosSchema == "" {
		dbosSchema = "act_dbos"
	}
	resultPath := filepath.Join(root, "result.json")
	environment := []string{
		"ACT_TEST_DBOS_URL=" + dsn,
		"ACT_TEST_REQUIRE_POSTGRES=1",
		"RILL_TEST_RESTART_REPO=" + repoDir,
		"RILL_TEST_RESTART_STORAGE=" + filepath.Join(root, "storage"),
		"RILL_TEST_RESTART_METASTORE_DSN=file:" + filepath.Join(root, "metastore.sqlite") + "?cache=shared",
		"RILL_TEST_RESTART_CATALOG_DSN=file:" + filepath.Join(root, "catalog.sqlite") + "?cache=shared",
		"RILL_TEST_RESTART_INSTANCE=" + instanceID,
		"RILL_TEST_RESTART_STORE_SCHEMA=" + storeSchema,
		"RILL_TEST_RESTART_DBOS_SCHEMA=" + dbosSchema,
		"RILL_TEST_RESTART_APP_VERSION=openai-process-restart-" + uuid.NewString(),
		"RILL_TEST_RESTART_IDEMPOTENCY=process-" + uuid.NewString(),
		"RILL_TEST_RESTART_MCP_URL=" + mcpServer.URL,
		"RILL_TEST_RESTART_RESULT=" + resultPath,
	}

	runRestartSubprocess(t, "phase1", environment)
	require.Eventually(t, func() bool { return len(provider.requests()) == 2 }, 5*time.Second, 20*time.Millisecond)
	require.Empty(t, forbidden.requests())

	writeRestartProjectFile(t, repoDir, ".env", "DEEPSEEK_API_KEY=rotated-secret\n")
	driftThinking := "disabled"
	if thinking == "disabled" {
		driftThinking = "enabled"
	}
	writeRestartProjectFile(t, repoDir, "connectors/deepseek.yaml", restartTemplatedConnectorYAML(forbidden.url(), "connector-model-v2", driftThinking))
	writeRestartProjectFile(t, repoDir, "triage.yaml", restartAgentYAML("default_ai", "agent-model-v2"))
	runRestartSubprocess(t, "phase2", environment)

	resultBytes, err := os.ReadFile(resultPath)
	require.NoError(t, err)
	var result restartSubprocessResult
	require.NoError(t, json.Unmarshal(resultBytes, &result))
	require.Equal(t, act.RunStatusSucceeded, result.Status)
	require.Equal(t, "Done: created OPS-42.", result.Response)
	require.Equal(t, "OPS-42", result.ActionRef)
	require.Equal(t, 1, result.ExecuteCount)

	reqs := provider.requests()
	require.Len(t, reqs, 3)
	require.Equal(t, "Bearer initial-secret", reqs[0].Authorization)
	require.Equal(t, "Bearer initial-secret", reqs[1].Authorization)
	require.Equal(t, "Bearer rotated-secret", reqs[2].Authorization)
	for turn, req := range reqs {
		require.Equal(t, "agent-model-v1", req.Body["model"])
		require.Equal(t, map[string]any{"type": thinking}, req.Body["thinking"])
		require.NoError(t, validateRestartRequest(req.Body, turn, thinking))
	}
	require.Empty(t, forbidden.requests())
}
