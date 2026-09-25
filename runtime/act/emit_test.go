package act_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

func eventTypes(events []*act.RunEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.EventType
	}
	return out
}

// TestExecutorEmitsLifecycleEventsOnApproval walks the happy path and asserts the run's state plane reflects every
// edge: queued at Start, then running -> waiting_approval -> resumed -> succeeded, with the approval persisted for the
// human bound to the exact arguments of the proposed action.
func TestExecutorEmitsLifecycleEventsOnApproval(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-7"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})
	ctx := t.Context()

	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste diario subio +40%.",
		IdempotencyKey: "run-emit-" + uuid.NewString(), Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Queryable the instant Start returns, before any decision: the run exists and is at least queued.
	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Contains(t, []act.RunStatus{act.RunStatusQueued, act.RunStatusRunning, act.RunStatusWaitingApproval}, run.Status)

	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	run, err = store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, run.Status)
	require.NotNil(t, run.FinishedOn)
	// The run is bound to the exact agent version it executed: the snapshot's spec hash, recorded on the running
	// transition (§8.3).
	require.Equal(t, "flow-spec-1", run.SpecHash)

	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Equal(t, []string{
		act.EventTypeQueued,
		act.EventTypeRunning,
		act.EventTypeWaitingApproval,
		act.EventTypeResumed,
		act.EventTypeSucceeded,
	}, eventTypes(events))

	approvals, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: instanceID, RunID: runID})
	require.NoError(t, err)
	require.Len(t, approvals, 1)
	require.Equal(t, act.ApprovalIDForToolCall(runID, "call-1"), approvals[0].ApprovalID)
	canonical, err := act.CanonicalizeArgs(flowProposal().Args)
	require.NoError(t, err)
	require.Equal(t, act.HashArgs(string(canonical)), approvals[0].ArgsHash)
	require.Equal(t, "user:alice", approvals[0].RequestedBy)
}

// TestExecutorEmitsRejection verifies a denied action records the rejected edge, attributes it to whoever denied it
// (Bob, not the run's initiator), and performs no external effect. The rejection is handed back to the model, which
// closes the run: a denied action does not fail the run.
func TestExecutorEmitsRejection(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-8"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})
	ctx := t.Context()

	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "propuesta arriesgada",
		IdempotencyKey: "run-reject-" + uuid.NewString(), Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Deny as Bob in the API's order: claim the approval first, then deliver the durable decision.
	approvalID := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, instanceID, approvalID)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(ctx, instanceID, approvalID, act.ApprovalStatusDenied, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalRejected, "call-1"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.False(t, res.ActionTaken)
	require.Equal(t, 0, exec.executeCount())

	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	var rejected *act.RunEvent
	for _, ev := range events {
		if ev.EventType == act.EventTypeRejected {
			rejected = ev
		}
	}
	require.NotNil(t, rejected, "the timeline records the rejection")
	require.Equal(t, "admin:bob", rejected.Payload["decided_by"], "the timeline attributes the denial to the denier")
}
