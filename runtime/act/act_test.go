package act_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	_ "github.com/rilldata/rill/runtime/resolvers"
)

// requirePostgres returns the DSN and schema for the DBOS system tables, skipping the test if Postgres is
// unreachable. The spike reuses WS1's local container (docker start dbos-spike-pg) unless overridden by env.
func requirePostgres(t *testing.T) (dsn, schema string) {
	t.Helper()
	dsn = os.Getenv("ACT_TEST_DBOS_URL")
	if dsn == "" {
		dsn = "postgres://postgres:dbos@localhost:55432/dbos_spike?sslmode=disable"
	}
	schema = os.Getenv("ACT_TEST_DBOS_SCHEMA")
	if schema == "" {
		schema = "act_dbos"
	}

	u, err := url.Parse(dsn)
	require.NoError(t, err)
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		t.Skipf("act tests need Postgres at %s (start it with: docker start dbos-spike-pg): %v", u.Host, err)
	}
	_ = conn.Close()
	return dsn, schema
}

// scriptedAI is a deterministic drivers.AIService: every completion returns the same text, which ends the tool
// loop immediately. It records how many times it was called so a test can assert the agent ran exactly once.
type scriptedAI struct {
	response string

	mu     sync.Mutex
	calls  int
	source runtime.RequestSource
}

var _ drivers.AIService = (*scriptedAI)(nil)

func (s *scriptedAI) Complete(ctx context.Context, _ *drivers.CompleteOptions) (*drivers.CompleteResult, error) {
	s.mu.Lock()
	s.calls++
	// Recorded so a test can pin what the meter will see: the completion runs deep inside the durable loop, far from
	// the request that started the run.
	s.source = runtime.RequestSourceFromContext(ctx)
	s.mu.Unlock()
	return &drivers.CompleteResult{
		Message: &aiv1.CompletionMessage{
			Role:    "assistant",
			Content: []*aiv1.ContentBlock{{BlockType: &aiv1.ContentBlock_Text{Text: s.response}}},
		},
		Provider:     "scripted",
		InputTokens:  1,
		OutputTokens: 1,
	}, nil
}

func (s *scriptedAI) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func (s *scriptedAI) requestSource() runtime.RequestSource {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.source
}

// toolForcingAI is a drivers.AIService that forces a call to toolName on its first completion and then returns
// plain text (ending the loop). It lets a test drive the agent into actually invoking a tool, so tool-access
// enforcement (the claims-based intersection) is exercised rather than bypassed by a model that never calls a tool.
type toolForcingAI struct {
	toolName string

	mu    sync.Mutex
	calls int
}

var _ drivers.AIService = (*toolForcingAI)(nil)

func (s *toolForcingAI) Complete(_ context.Context, _ *drivers.CompleteOptions) (*drivers.CompleteResult, error) {
	s.mu.Lock()
	n := s.calls
	s.calls++
	s.mu.Unlock()

	msg := &aiv1.CompletionMessage{Role: "assistant"}
	if n == 0 {
		args, err := structpb.NewStruct(map[string]any{})
		if err != nil {
			return nil, err
		}
		msg.Content = []*aiv1.ContentBlock{{BlockType: &aiv1.ContentBlock_ToolCall{
			ToolCall: &aiv1.ToolCall{Id: "call_1", Name: s.toolName, Input: args},
		}}}
	} else {
		msg.Content = []*aiv1.ContentBlock{{BlockType: &aiv1.ContentBlock_Text{Text: "listo"}}}
	}
	return &drivers.CompleteResult{Message: msg, Provider: "scripted", InputTokens: 1, OutputTokens: 1}, nil
}

// recordingSink is an in-memory act.ActionSink that records each applied action so a test can assert the action
// ran exactly once and inspect what it was given (e.g. the snapshot marker).
type recordingSink struct {
	mu      sync.Mutex
	applied []act.ActionRequest
}

var _ act.ActionSink = (*recordingSink)(nil)

func (s *recordingSink) Apply(_ context.Context, req act.ActionRequest) (act.ActionResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applied = append(s.applied, req)
	return act.ActionResult{Ref: fmt.Sprintf("ticket-%d", len(s.applied))}, nil
}

func (s *recordingSink) snapshot() []act.ActionRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]act.ActionRequest(nil), s.applied...)
}

// agentYAML renders a minimal, valid Agent resource with the given single-line instructions.
func agentYAML(instructions string) string {
	return fmt.Sprintf(`
type: agent
display_name: Ticket Triage
instructions: %q
tools: []
limits:
  max_steps: 3
`, instructions)
}

// newInstanceWithAgent creates a runtime instance whose project contains one reconciled Agent named "triage".
func newInstanceWithAgent(t *testing.T, instructions string) (*runtime.Runtime, string) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": `ai_connector: mock_ai
ai_instructions: "Menciona la residencia de datos cuando sea relevante."`,
			"triage.yaml": agentYAML(instructions),
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	// Sanity check: the agent reconciled to a valid, resolvable snapshot named "triage".
	snap, err := act.NewCatalogAgentProvider(rt).GetAgent(t.Context(), instanceID, "triage")
	require.NoError(t, err)
	require.Equal(t, "triage", snap.Name)
	return rt, instanceID
}

// scriptedSessionFactory returns a SessionFactory that builds a session bound to the given scripted model, using the
// initiator's claims that the run carries. It does not force SkipChecks: the session runs with exactly the claims
// passed, so a restricted actor's run resolves tool access against the restricted claims (§17.3). A run started
// without an initiator identity (nil claims: the executor's own storeless unit tests) falls back to a local-dev
// SkipChecks session, mirroring Rill Developer where auth is disabled.
func scriptedSessionFactory(rt *runtime.Runtime, llm drivers.AIService) act.SessionFactory {
	return func(ctx context.Context, instanceID, sessionID string, claims *runtime.SecurityClaims, _ *ai.AgentSnapshot) (*ai.Session, func(), error) {
		if claims == nil {
			claims = &runtime.SecurityClaims{UserID: uuid.NewString(), SkipChecks: true}
		}
		runner := ai.NewRunner(rt, activity.NewNoopClient())
		s, err := runner.Session(ctx, &ai.SessionOptions{
			InstanceID: instanceID,
			SessionID:  sessionID,
			Claims:     claims,
			UserAgent:  "act-spike",
		})
		if err != nil {
			return nil, nil, err
		}
		s.SetLLM(func(context.Context) (drivers.AIService, func(), error) { return llm, func() {}, nil })
		return s, func() {}, nil
	}
}

// newExecutor builds an in-process executor running a real SessionRunner (CatalogAgentProvider + scripted model +
// the given sink). Each executor gets a unique application version so it never recovers another test's runs.
func newExecutor(t *testing.T, rt *runtime.Runtime, llm drivers.AIService, sink act.ActionSink) *act.DBOSExecutor {
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, llm),
		Actions:  sink,
	}
	return newExecutorWithRunner(t, runner)
}

func newExecutorWithRunner(t *testing.T, runner act.Runner) *act.DBOSExecutor {
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL:        dsn,
		DatabaseSchema:     schema,
		ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner:             runner,
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })
	return e
}
