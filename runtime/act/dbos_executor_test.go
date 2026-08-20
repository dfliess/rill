package act_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// TestExecutorRejectsIncompleteModelSnapshotBeforeAnyRunSideEffect pins the upgrade guard's placement. A snapshot
// checkpointed by a pre-routing build has no effective connector/driver; recovery must reject it before opening a
// session or consuming an approval that could lead to an external write.
func TestExecutorRejectsIncompleteModelSnapshotBeforeAnyRunSideEffect(t *testing.T) {
	var sessionCalls atomic.Int32
	sink := &recordingSink{}
	runner := &act.SessionRunner{
		Provider: ai.NewStaticAgentProvider(&ai.AgentSnapshot{Name: "legacy", Instructions: "old checkpoint"}),
		Sessions: func(context.Context, string, string, *runtime.SecurityClaims, *ai.AgentSnapshot) (*ai.Session, func(), error) {
			sessionCalls.Add(1)
			return nil, nil, errors.New("session must not open")
		},
		Actions: sink,
	}
	e := newExecutorWithRunner(t, runner)
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: "legacy-instance", AgentName: "legacy", Prompt: "resume",
		IdempotencyKey: "legacy-" + uuid.NewString(),
	})
	require.NoError(t, err)

	_, err = e.Result(runID)
	require.ErrorContains(t, err, "durable model connector snapshot")
	require.Zero(t, sessionCalls.Load(), "validation must run before the model session opens")
	require.Empty(t, sink.snapshot(), "validation must run before an external action")
}

// TestExecutorStartApproveRunsActionOnce is the happy path: a run investigates, waits for approval, and on
// "approved" performs its (simulated) write exactly once, ending in the succeeded state.
func TestExecutorStartApproveRunsActionOnce(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon un ticket de remediacion.")
	llm := &scriptedAI{response: "Propongo crear un ticket P2 por el pico de coste."}
	sink := &recordingSink{}
	e := newExecutor(t, rt, llm, sink)

	key := "run-a-" + uuid.NewString()
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio +40% en el proyecto X.",
		IdempotencyKey: key,
	})
	require.NoError(t, err)
	require.Equal(t, act.ComposeRunID(instanceID, "triage", key), runID)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, ""))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.True(t, res.ActionTaken)
	require.NotEmpty(t, res.ActionRef)
	require.Equal(t, "Propongo crear un ticket P2 por el pico de coste.", res.Response)

	// The write ran exactly once, and the agent loop was driven exactly once.
	require.Len(t, sink.snapshot(), 1)
	require.Equal(t, 1, llm.callCount())
}

// TestExecutorRejectionSkipsAction verifies the executor never touches the outside world without approval.
func TestExecutorRejectionSkipsAction(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon un ticket.")
	llm := &scriptedAI{response: "Propongo reiniciar el pipeline."}
	sink := &recordingSink{}
	e := newExecutor(t, rt, llm, sink)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "Fallo en el pipeline nocturno.",
		IdempotencyKey: "run-reject-" + uuid.NewString(),
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalRejected, ""))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusRejected, res.Status)
	require.False(t, res.ActionTaken)
	require.Empty(t, sink.snapshot())
}

// TestExecutorDedupsByIdempotencyKey verifies OAOO: two Starts with the same IdempotencyKey are one run. The
// second Start attaches to the first run rather than starting a new one, so the agent and the action each run once.
func TestExecutorDedupsByIdempotencyKey(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	llm := &scriptedAI{response: "propuesta unica"}
	sink := &recordingSink{}
	e := newExecutor(t, rt, llm, sink)

	in := act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "alerta duplicada",
		IdempotencyKey: "run-dedup-" + uuid.NewString(),
	}

	id1, err := e.Start(t.Context(), in)
	require.NoError(t, err)
	id2, err := e.Start(t.Context(), in) // same key: must attach to the same run
	require.NoError(t, err)
	require.Equal(t, id1, id2)

	require.NoError(t, e.Resume(t.Context(), id1, act.ApprovalApproved, ""))

	res, err := e.Result(id1)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	// Despite two Starts, exactly one execution: one agent run, one action.
	require.Equal(t, 1, llm.callCount())
	require.Len(t, sink.snapshot(), 1)
}

// gatedRunner wraps a Runner and closes ranAgent the first time RunSegment is called. By then LoadSnapshot has
// already returned the snapshot to the workflow, so a test can safely mutate the agent definition afterwards and
// know the run is already bound to the original snapshot.
type gatedRunner struct {
	act.Runner
	ranAgent chan struct{}
	once     sync.Once
}

func (g *gatedRunner) RunSegment(ctx context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	g.once.Do(func() { close(g.ranAgent) })
	return g.Runner.RunSegment(ctx, in)
}

// recordingRunner wraps a Runner and captures the session ID its RunSegment returned, so a test can assert the run's
// persisted conversation_id matches the session the run actually executed in.
type recordingRunner struct {
	act.Runner
	mu        sync.Mutex
	sessionID string
}

func (r *recordingRunner) RunSegment(ctx context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	out, err := r.Runner.RunSegment(ctx, in)
	r.mu.Lock()
	if out.SessionID != "" {
		r.sessionID = out.SessionID
	}
	r.mu.Unlock()
	return out, err
}

func (r *recordingRunner) capturedSessionID() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.sessionID
}

// TestExecutorLinksRunToConversation verifies §13.1: a finished run is linked to the AI session it executed in. The
// store's conversation_id ends up equal to the session ID the runner opened (and non-empty), so the API and UI can
// render the run as its conversation. CreateRun seeds conversation_id empty here (no ConversationID on the input), so
// this proves the persist step wrote the real session, not a leftover input value.
func TestExecutorLinksRunToConversation(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon un ticket.")
	store := newRunStore(t)
	llm := &scriptedAI{response: "Propongo crear un ticket."}
	base := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, llm),
		Actions:  &recordingSink{},
	}
	runner := &recordingRunner{Runner: base}
	e := newStoreExecutor(t, runner, store)
	ctx := t.Context()

	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio +40%.",
		IdempotencyKey: "run-conv-" + uuid.NewString(),
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, ""))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	sessionID := runner.capturedSessionID()
	require.NotEmpty(t, sessionID, "the runner opened a session for the run")

	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.NotEmpty(t, run.ConversationID, "the run must be linked to a conversation")
	require.Equal(t, sessionID, run.ConversationID, "the run is linked to the exact session it executed in")
}

// TestExecutorUsesSnapshotFromRunStart verifies §8.3: a run is bound to the agent snapshot captured when it
// starts. Editing the agent mid-run changes what the provider returns, but not what the in-flight run does.
func TestExecutorUsesSnapshotFromRunStart(t *testing.T) {
	const original = "ORIGINAL: propon un ticket P2."
	const edited = "EDITED: reinicia el cluster de inmediato."

	rt, instanceID := newInstanceWithAgent(t, original)
	llm := &scriptedAI{response: "propuesta"}
	sink := &recordingSink{}

	base := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, llm),
		Actions:  sink,
	}
	gate := &gatedRunner{Runner: base, ranAgent: make(chan struct{})}
	e := newExecutorWithRunner(t, gate)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "alerta",
		IdempotencyKey: "run-snap-" + uuid.NewString(),
	})
	require.NoError(t, err)

	// Wait until the snapshot has been captured, then edit the agent and re-reconcile.
	<-gate.ranAgent
	testruntime.PutFiles(t, rt, instanceID, map[string]string{"triage.yaml": agentYAML(edited)})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	// The provider now serves the edited definition: proof the edit really took effect.
	snapNow, err := act.NewCatalogAgentProvider(rt).GetAgent(t.Context(), instanceID, "triage")
	require.NoError(t, err)
	require.Contains(t, snapNow.Instructions, "EDITED")

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, ""))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	// The action carried the ORIGINAL snapshot, not the edited one: the run used the definition from its start.
	applied := sink.snapshot()
	require.Len(t, applied, 1)
	require.Contains(t, applied[0].SnapshotMark, "ORIGINAL")
	require.NotContains(t, applied[0].SnapshotMark, "EDITED")
}

// TestComposeRunIDIsInjective guards the run-ID encoding against delimiter ambiguity: two logically distinct
// (instance, agent, key) tuples that would collide under naive "/" joining must produce different run IDs, or a
// Start could attach to the wrong run and cross approvals between tenants or agents.
func TestComposeRunIDIsInjective(t *testing.T) {
	a := act.ComposeRunID("inst", "a/b", "c")
	b := act.ComposeRunID("inst", "a", "b/c")
	require.NotEqual(t, a, b)

	// Same tuple is stable (idempotency depends on it).
	require.Equal(t, act.ComposeRunID("inst", "triage", "k"), act.ComposeRunID("inst", "triage", "k"))
}

// TestExecutorBindsActorToRun verifies the run is durably bound to the identity that authorized it: the actor
// passed to Start is checkpointed with the run and reaches the action, where the audit ledger will record it.
func TestExecutorBindsActorToRun(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon un ticket.")
	llm := &scriptedAI{response: "propuesta con actor"}
	sink := &recordingSink{}
	e := newExecutor(t, rt, llm, sink)

	actor := act.Actor{Subject: "user:alice", ServicePrincipal: false}
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio.",
		IdempotencyKey: "run-actor-" + uuid.NewString(),
		Actor:          actor,
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, ""))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	// The authorizing identity survived the durable round-trip (Start input, checkpoint, action step).
	applied := sink.snapshot()
	require.Len(t, applied, 1)
	require.Equal(t, actor, applied[0].Actor)
}

// TestExecutorTagsRunsAsActTraffic pins where an agent's spend gets charged. The tokens and tool calls are recorded
// from inside the durable loop, which does not run in the request that started it, so the source has to come from the
// executor's own context rather than from the API handler that tagged its own.
func TestExecutorTagsRunsAsActTraffic(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon un ticket.")
	llm := &scriptedAI{response: "propuesta"}
	sink := &recordingSink{}
	e := newExecutor(t, rt, llm, sink)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio.",
		IdempotencyKey: "run-source-" + uuid.NewString(),
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, ""))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	require.Equal(t, runtime.RequestSourceAct, llm.requestSource())
}

// TestRunAgentUsesInitiatorClaimsForToolAccess proves the fix for the confused-deputy hole (§17.3): a run executes
// tools with the initiating actor's claims, not the worker's. The same agent, snapshot and model, run under two
// different identities, produce opposite outcomes purely because of the claims the run carries — an actor without
// ReadObjects cannot execute a read tool the snapshot declares, while a privileged actor can.
func TestRunAgentUsesInitiatorClaimsForToolAccess(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "irrelevante: usamos un provider estático")

	// A dynamic agent whose only tool is a read-only metrics tool gated by the ReadObjects permission.
	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "reader",
		Instructions: "Lista las vistas de métricas disponibles.",
		Tools:        []string{ai.ListMetricsViewsName},
		MaxSteps:     3,
	})
	snap, err := provider.GetAgent(t.Context(), instanceID, "reader")
	require.NoError(t, err)

	newRunner := func() *act.SessionRunner {
		return &act.SessionRunner{
			Provider: provider,
			Sessions: scriptedSessionFactory(rt, &toolForcingAI{toolName: ai.ListMetricsViewsName}),
			Actions:  &recordingSink{},
		}
	}

	// Restricted actor: no ReadObjects, so CheckAccess drops the tool from the run. The forced call then hits a tool
	// that is not in the allowed set and the run fails closed — the actor cannot read metrics through the run.
	restricted := &runtime.SecurityClaims{UserID: "user:restricted"}
	_, err = newRunner().RunSegment(t.Context(), act.RunSegmentInput{
		InstanceID: instanceID, Claims: restricted, Snapshot: snap, Prompt: "listalas",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), ai.ListMetricsViewsName)
	require.Contains(t, err.Error(), "not in the allowed set")

	// Privileged actor: the same tool is available, so the same forced call executes and the run completes. The only
	// difference from the failing case is the initiator's claims.
	privileged := &runtime.SecurityClaims{UserID: "user:privileged", SkipChecks: true}
	out, err := newRunner().RunSegment(t.Context(), act.RunSegmentInput{
		InstanceID: instanceID, Claims: privileged, Snapshot: snap, Prompt: "listalas",
	})
	require.NoError(t, err)
	require.Equal(t, "listo", out.Response)
	require.NotEmpty(t, out.SessionID, "RunSegment returns the session it opened so the run can link its conversation")
}
