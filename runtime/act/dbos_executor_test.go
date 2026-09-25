package act_test

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// TestExecutorDedupsByIdempotencyKey verifies OAOO: two Starts with the same IdempotencyKey are one run. The
// second Start attaches to the first run rather than starting a new one, so the agent runs once.
func TestExecutorDedupsByIdempotencyKey(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	llm := &scriptedAI{response: "respuesta unica"}
	e := newExecutor(t, rt, llm)

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

	res, err := e.Result(id1)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, llm.callCount(), "despite two Starts, the agent ran once")
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
// this proves the persist step wrote the real session, not a leftover input value. The scripted model answers in
// text, so the run investigates and finishes without a write.
func TestExecutorLinksRunToConversation(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon un ticket.")
	store := newRunStore(t)
	llm := &scriptedAI{response: "Propongo crear un ticket."}
	base := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, llm),
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

// snapshotProbeRunner resolves the agent from the real catalog but runs a scripted loop: the first segment blocks
// until the test has edited the agent, then proposes one write; the resumed segment closes. Every segment records the
// instructions it ran with, so a test can see which definition each one used.
type snapshotProbeRunner struct {
	provider act.Runner
	started  chan struct{}
	release  chan struct{}

	mu   sync.Mutex
	seen []string
}

func (r *snapshotProbeRunner) LoadSnapshot(ctx context.Context, instanceID, agentName string) (*ai.AgentSnapshot, error) {
	return r.provider.LoadSnapshot(ctx, instanceID, agentName)
}

func (r *snapshotProbeRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	r.mu.Lock()
	r.seen = append(r.seen, in.Snapshot.Instructions)
	r.mu.Unlock()
	if len(in.Resume) > 0 {
		return act.RunSegmentResult{Response: "hecho", SessionID: in.SessionID}, nil
	}
	close(r.started)
	<-r.release
	return act.RunSegmentResult{Response: "propongo un ticket", SessionID: "sess-snap", Proposed: []*ai.ProposedAction{{
		ToolCallID: "call-1", Tool: "jira.create_issue", Connector: "jira_ops", Args: map[string]any{"summary": "P2"},
	}}}, nil
}

func (r *snapshotProbeRunner) instructionsSeen() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.seen...)
}

// TestExecutorUsesSnapshotFromRunStart verifies §8.3: a run is bound to the agent snapshot captured when it starts.
// The agent is edited while the run is inside its first segment; the segment resumed after the approval, long after
// the edit, still runs the original definition.
func TestExecutorUsesSnapshotFromRunStart(t *testing.T) {
	const original = "ORIGINAL: propon un ticket P2."
	const edited = "EDITED: reinicia el cluster de inmediato."

	rt, instanceID := newInstanceWithAgent(t, original)
	runner := &snapshotProbeRunner{
		provider: &act.SessionRunner{Provider: act.NewCatalogAgentProvider(rt)},
		started:  make(chan struct{}),
		release:  make(chan struct{}),
	}
	e := newExecutorWithRunner(t, runner)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "alerta",
		IdempotencyKey: "run-snap-" + uuid.NewString(),
	})
	require.NoError(t, err)

	// The snapshot is captured before the first segment starts. Edit the agent and re-reconcile while it runs.
	<-runner.started
	testruntime.PutFiles(t, rt, instanceID, map[string]string{"triage.yaml": agentYAML(edited)})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	snapNow, err := act.NewCatalogAgentProvider(rt).GetAgent(t.Context(), instanceID, "triage")
	require.NoError(t, err)
	require.Contains(t, snapNow.Instructions, "EDITED", "the edit really took effect")
	close(runner.release)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	seen := runner.instructionsSeen()
	require.Len(t, seen, 2, "a first segment and one resumed after the approval")
	for _, instructions := range seen {
		require.Contains(t, instructions, "ORIGINAL")
		require.NotContains(t, instructions, "EDITED")
	}
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

// TestExecutorBindsActorToRun verifies the run is durably bound to the identity that authorized it: the actor passed
// to Start is checkpointed with the run and recorded on the action it proposed, which is what the audit reads.
func TestExecutorBindsActorToRun(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-3"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})

	actor := act.Actor{Subject: "user:alice"}
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio.",
		IdempotencyKey: "run-actor-" + uuid.NewString(),
		Actor:          actor,
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	run, err := store.GetRun(t.Context(), instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, actor.Subject, run.Actor.Subject)
	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, actor.Subject, action.RequestedBy, "the write is recorded on behalf of the run's initiator")
}

// TestExecutorTagsRunsAsActTraffic pins where an agent's spend gets charged. The tokens and tool calls are recorded
// from inside the durable loop, which does not run in the request that started it, so the source has to come from the
// executor's own context rather than from the API handler that tagged its own.
func TestExecutorTagsRunsAsActTraffic(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon un ticket.")
	llm := &scriptedAI{response: "propuesta"}
	e := newExecutor(t, rt, llm)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio.",
		IdempotencyKey: "run-source-" + uuid.NewString(),
	})
	require.NoError(t, err)

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
