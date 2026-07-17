package act_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// blockingExecutor is an ActionExecutor that parks inside Execute until released, so a test can hold an action in the
// executing state under a live lease while it drives a second, concurrent Execute. It records how many times it was
// entered so a test can assert the external write ran exactly once.
type blockingExecutor struct {
	entered chan struct{}
	release chan struct{}
	result  act.ExecuteResult

	mu    sync.Mutex
	calls int
}

var _ act.ActionExecutor = (*blockingExecutor)(nil)

func (b *blockingExecutor) Execute(_ context.Context, _ act.ExecuteRequest) (act.ExecuteResult, error) {
	b.mu.Lock()
	b.calls++
	b.mu.Unlock()
	b.entered <- struct{}{}
	<-b.release
	return b.result, nil
}

func (b *blockingExecutor) Verify(_ context.Context, _ act.VerifyRequest) (act.VerifyResult, error) {
	return act.VerifyResult{}, nil
}

func (b *blockingExecutor) callCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.calls
}

// proposeAndApprove drives a proposal through propose and human approval without executing it, leaving the action in
// the approved state ready for an Execute.
func proposeAndApprove(t *testing.T, g *act.Gateway, instanceID, runID string, proposal act.ToolProposal) act.Authorization {
	t.Helper()
	ctx := t.Context()
	auth, err := g.Propose(ctx, act.ProposeActionInput{
		InstanceID: instanceID, RunID: runID, AgentName: "triage", Actor: act.Actor{Subject: "u"}, Proposal: proposal,
	})
	require.NoError(t, err)
	require.NoError(t, g.RecordApproved(ctx, act.ApprovedInput{
		InstanceID: instanceID, RunID: runID, ToolCallID: proposal.ToolCallID, ArgsHash: auth.ArgsHash, DecidedBy: "admin",
	}))
	return auth
}

// TestGatewayLiveLeaseNotReclaimedSingleEffect proves the lease-based ownership guarantee under concurrency (§12): with
// one worker mid-write under a live lease, a second worker on the SAME approved action neither reclaims the row nor
// runs a second external effect, and the owner alone finalizes its claim. The second worker gets a retryable error
// (ErrLeaseHeld) rather than a terminal indeterminate outcome, so the workflow retries the step instead of permanently
// failing the run while the owner may still succeed. It covers A5 (a), (b) and (d).
func TestGatewayLiveLeaseNotReclaimedSingleEffect(t *testing.T) {
	exec := &blockingExecutor{
		entered: make(chan struct{}),
		release: make(chan struct{}),
		result:  act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"},
	}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil)
	const inst, run = "inst-live-lease", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()
	proposeAndApprove(t, g, inst, run, proposal)

	// Worker 1 claims the write and parks inside the executor, holding a live lease.
	type result struct {
		report act.ExecuteReport
		err    error
	}
	done := make(chan result, 1)
	go func() {
		r, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
		done <- result{r, err}
	}()
	<-exec.entered // worker 1 has claimed approved->executing and is mid-write

	// Worker 2 finds the action executing under a LIVE foreign lease: it must not reclaim it and must not run a second
	// effect. It returns a retryable error so the workflow retries later rather than terminalizing the run.
	_, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.ErrorIs(t, err, act.ErrLeaseHeld, "a live foreign lease signals retry, not a terminal outcome")
	require.Equal(t, 1, exec.callCount(), "worker 2 does not run a second external effect")

	mid, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionExecuting, mid.Status, "a live lease is not reclaimed to indeterminate")
	require.NotEmpty(t, mid.AttemptID)

	// Release worker 1: it finalizes succeeded under its own attempt, having been the only one to execute.
	close(exec.release)
	w1 := <-done
	require.NoError(t, w1.err)
	require.Equal(t, act.OutcomeSucceeded, w1.report.Outcome)
	require.Equal(t, 1, exec.callCount(), "the external write ran exactly once")

	final, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, final.Status)
	require.Equal(t, mid.AttemptID, final.AttemptID, "the owning attempt finalized its own claim")

	// Worker 2 retries after the owner finalized: it now reconciles to the succeeded outcome without re-executing.
	r2, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, r2.Outcome, "a retry after the owner finalized reports the recorded outcome")
	require.Equal(t, 1, exec.callCount(), "the external write still ran exactly once")
}

// TestGatewayExpiredLeaseReclaimedToIndeterminate proves recovery of a presumed-dead worker (§12): an action left
// executing under an EXPIRED lease is reclaimed to indeterminate on the next Execute, and the interrupted write is
// never re-issued. It covers A5 (c).
func TestGatewayExpiredLeaseReclaimedToIndeterminate(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassNonIdempotent), nil)
	const inst, run = "inst-expired-lease", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()
	proposeAndApprove(t, g, inst, run, proposal)

	// Simulate a worker that claimed the write and then died: it holds an already-expired lease under attempt "dead-1".
	claimed, _, err := ledger.ClaimExecuting(ctx, inst, run, "call-1", "dead-1", -5*time.Second)
	require.NoError(t, err)
	require.True(t, claimed)

	// Recovery finds the executing action with an expired lease: it reclaims to indeterminate, never re-issuing the write.
	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeIndeterminate, report.Outcome)
	require.Equal(t, 0, exec.executeCount(), "a reclaimed interrupted write is never re-issued")

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionIndeterminate, action.Status)
}

// TestGatewayLiveLeaseDoesNotTerminalizeRun proves that a crash-and-restart within the lease window does not
// permanently fail the run (recovery inside a still-live lease): the restarted worker gets a retryable error (not an
// indeterminate business outcome), the action stays executing under the live lease, and a subsequent retry after the
// owner finalizes reports the correct outcome without a second external effect.
func TestGatewayLiveLeaseDoesNotTerminalizeRun(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	g, ledger := newGateway(t, exec, createIssueDescriptor(act.ClassNonIdempotent), nil)
	const inst, run = "inst-live-lease-no-terminal", "run-1"
	proposal := createIssueProposal()
	ctx := t.Context()
	proposeAndApprove(t, g, inst, run, proposal)

	// Simulate a worker that claimed the write and is still alive (lease far from expired).
	claimed, _, err := ledger.ClaimExecuting(ctx, inst, run, "call-1", "alive-1", 5*time.Minute)
	require.NoError(t, err)
	require.True(t, claimed)

	// A restarted worker (or DBOS retry) calls Execute while the lease is still live: it must get a retryable error,
	// not a terminal indeterminate outcome that would permanently fail the run.
	_, err = g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.ErrorIs(t, err, act.ErrLeaseHeld)

	action, err := ledger.GetAction(ctx, inst, run, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionExecuting, action.Status, "the action stays executing under the live lease")
	require.Equal(t, "alive-1", action.AttemptID, "the owner's attempt is untouched")
	require.Equal(t, 0, exec.executeCount(), "no external effect was driven")

	// Simulate the owner finalizing successfully.
	_, err = ledger.TransitionAction(ctx, act.ActionTransition{
		InstanceID: inst, RunID: run, ToolCallID: "call-1", Status: act.ActionSucceeded,
		ExternalReference: "PROJ-1", AttemptID: "alive-1",
	})
	require.NoError(t, err)

	// The next DBOS retry finds the action succeeded: the run recovers without permanent failure.
	report, err := g.Execute(ctx, act.ExecuteActionInput{InstanceID: inst, RunID: run, AgentName: "triage", Proposal: proposal})
	require.NoError(t, err)
	require.Equal(t, act.OutcomeSucceeded, report.Outcome)
	require.Equal(t, "PROJ-1", report.ExternalReference)
	require.Equal(t, 0, exec.executeCount(), "no second external effect was ever driven")
}
