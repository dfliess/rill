package act_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// panickingRunner panics inside RunSegment when the prompt asks for it, and otherwise delegates to the real runner.
// RunSegment runs inside a checkpointed step, so this is a panic in a step of the durable workflow.
type panickingRunner struct {
	act.Runner
}

func (r *panickingRunner) RunSegment(ctx context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	if strings.Contains(in.Prompt, "PANIC") {
		panic("boom in a step")
	}
	return r.Runner.RunSegment(ctx, in)
}

// TestExecutorSurvivesPanicInStep verifies kairos-cloud#215 and #216: a panic inside a step fails that run, records
// it as failed in the store, and leaves the process, and every other run in it, alive. Up to DBOS v1.3.0 the panic
// escaped and killed the whole runtime, runs parked on human approval included (#215; with v1.3.0 this test binary
// crashes on the first run). DBOS v1.5.0 recovers it outside the workflow, so without our own recover the run stayed
// running in the store forever (#216; verified by mutation: without the recover in runAgentWorkflow the status
// assertion below fails with "running").
func TestExecutorSurvivesPanicInStep(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	base := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, &scriptedAI{response: "todo en orden"}),
	}
	store := newRunStore(t)
	e := newStoreExecutor(t, &panickingRunner{Runner: base}, store)
	ctx := t.Context()

	panicID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "PANIC",
		IdempotencyKey: "run-panic-" + uuid.NewString(),
	})
	require.NoError(t, err)
	_, err = e.Result(panicID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "boom in a step", "the run's error carries the panic")

	run, err := store.GetRun(ctx, instanceID, panicID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, run.Status, "a panicked run is closed as failed, not left running")
	events, err := store.ListRunEvents(ctx, instanceID, panicID, 0, 100)
	require.NoError(t, err)
	require.Contains(t, eventTypes(events), act.EventTypeFailed)

	okID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "alerta normal",
		IdempotencyKey: "run-after-panic-" + uuid.NewString(),
	})
	require.NoError(t, err)
	res, err := e.Result(okID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status, "the executor keeps serving runs after a panic")
}

// panickingExecutor panics while performing the external write, after the gateway claimed the action as executing.
type panickingExecutor struct {
	*fakeExecutor
}

func (p *panickingExecutor) Execute(context.Context, act.ExecuteRequest) (act.ExecuteResult, error) {
	panic("boom mid-write")
}

// TestExecutorPanicMidWriteLeavesActionIndeterminate verifies kairos-cloud#216: a panic while an external write is in
// flight closes the run as failed and hands the action to a human as indeterminate, since nobody can tell whether the
// write landed. Left executing, nothing revisits it: the workflow errored, so DBOS never recovers it, and the lease
// reclaim only runs inside the run. Verified by mutation: without the recover the action stays executing.
func TestExecutorPanicMidWriteLeavesActionIndeterminate(t *testing.T) {
	store := newRunStore(t)
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassUnknown)),
		Ledger:   store,
		Executor: &panickingExecutor{fakeExecutor: &fakeExecutor{}},
	}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: autoApproveRunner{}, Store: store, Gateway: gateway, Proposer: act.NewCapturedProposer(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-panic-midwrite"
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "panic-midwrite-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)
	_, err = e.Result(runID)
	require.Error(t, err)

	action, err := store.GetAction(t.Context(), inst, runID, "jira.create_issue")
	require.NoError(t, err)
	require.Equal(t, act.ActionIndeterminate, action.Status, "an interrupted write goes to a human, it is not left executing")
	run, err := store.GetRun(t.Context(), inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, run.Status)
}
