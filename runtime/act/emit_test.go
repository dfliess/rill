package act_test

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// newStoreExecutor builds an in-process executor wired to a real RunStore, so a test can observe the run state and
// event stream the executor emits as a run progresses. The DBOS system tables and the product store live in separate
// schemas of the same Postgres, exactly as production separates them (§15.2).
func newStoreExecutor(t *testing.T, runner act.Runner, store act.RunStore) *act.DBOSExecutor {
	t.Helper()
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL:        dsn,
		DatabaseSchema:     schema,
		ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner:             runner,
		Store:              store,
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })
	return e
}

func eventTypes(events []*act.RunEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.EventType
	}
	return out
}

// TestExecutorEmitsLifecycleEventsOnApproval walks the happy path and asserts the run's state plane reflects every
// edge: queued at Start, then the full running -> waiting_approval -> resumed -> succeeded sequence, with a matching
// approval persisted for the human to act on.
func TestExecutorEmitsLifecycleEventsOnApproval(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon un ticket.")
	store := newRunStore(t)
	llm := &scriptedAI{response: "Propongo crear un ticket P2 por el pico de coste."}
	sink := &recordingSink{}
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, llm),
		Actions:  sink,
	}
	e := newStoreExecutor(t, runner, store)
	ctx := t.Context()

	key := "run-emit-" + uuid.NewString()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "El coste diario subio +40%.",
		IdempotencyKey: key,
		Actor:          act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Queryable the instant Start returns, before any decision: the run exists and is at least queued.
	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Contains(t, []act.RunStatus{act.RunStatusQueued, act.RunStatusRunning, act.RunStatusWaitingApproval}, run.Status)

	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, ""))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	// The store's terminal view agrees with the workflow result.
	run, err = store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, run.Status)
	require.NotNil(t, run.FinishedOn)
	// The run is bound to the exact agent version it executed: the snapshot's spec hash was recorded on the running
	// transition (§8.3).
	require.NotEmpty(t, run.SpecHash)

	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Equal(t, []string{
		act.EventTypeQueued,
		act.EventTypeRunning,
		act.EventTypeWaitingApproval,
		act.EventTypeResumed,
		act.EventTypeSucceeded,
	}, eventTypes(events))

	// The approval was persisted for the human, bound to the exact proposal the agent produced.
	approvals, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: instanceID, RunID: runID})
	require.NoError(t, err)
	require.Len(t, approvals, 1)
	require.Equal(t, act.ApprovalIDForRun(runID), approvals[0].ApprovalID)
	require.Equal(t, act.HashArgs(res.Response), approvals[0].ArgsHash)
	require.Equal(t, "user:alice", approvals[0].RequestedBy)
}

// TestExecutorEmitsRejection verifies a denied run records the rejected transition, attributes it to whoever denied it
// (Bob, not the run's initiator), and performs no external effect.
func TestExecutorEmitsRejection(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	store := newRunStore(t)
	llm := &scriptedAI{response: "propuesta a rechazar"}
	sink := &recordingSink{}
	runner := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, llm),
		Actions:  sink,
	}
	e := newStoreExecutor(t, runner, store)
	ctx := t.Context()

	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "propuesta arriesgada",
		IdempotencyKey: "run-reject-" + uuid.NewString(),
		Actor:          act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Deny as Bob in the API's order: claim the approval first, then deliver the durable decision.
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, instanceID, act.ApprovalIDForRun(runID))
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(ctx, instanceID, act.ApprovalIDForRun(runID), act.ApprovalStatusDenied, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalRejected, ""))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusRejected, res.Status)
	require.Empty(t, sink.snapshot())

	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	rejected := events[len(events)-1]
	require.Equal(t, act.EventTypeRejected, rejected.EventType)
	require.Equal(t, "admin:bob", rejected.Payload["decided_by"], "the timeline attributes the denial to the denier")
}
