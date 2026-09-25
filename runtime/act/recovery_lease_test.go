package act_test

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// deadAttemptLedger reproduces what a recovered run finds when the process died mid-write: the approved action is
// already executing, claimed by an attempt that no longer exists, under a lease that is still live. The first time the
// gateway reads the action as approved, a phantom attempt claims it with a short lease before the read returns, so the
// workflow's gateway_execute step meets exactly the state a fast restart leaves behind (kairos-cloud#170).
//
// onClaimed, when set, runs in the background after the phantom claim, standing in for an owner that was not dead
// after all and finalizes its write while the recovered step waits.
type deadAttemptLedger struct {
	act.ActionLedger
	lease     time.Duration
	onClaimed func(instanceID, runID, toolCallID string)
	once      sync.Once
}

const deadAttemptID = "attempt-of-a-dead-process"

func (l *deadAttemptLedger) GetAction(ctx context.Context, instanceID, runID, toolCallID string) (*act.Action, error) {
	a, err := l.ActionLedger.GetAction(ctx, instanceID, runID, toolCallID)
	if err != nil || a.Status != act.ActionApproved {
		return a, err
	}
	var claimErr error
	l.once.Do(func() {
		_, _, claimErr = l.ActionLedger.ClaimExecuting(ctx, instanceID, runID, toolCallID, deadAttemptID, l.lease)
		if claimErr == nil && l.onClaimed != nil {
			go l.onClaimed(instanceID, runID, toolCallID)
		}
	})
	if claimErr != nil {
		return nil, claimErr
	}
	return l.ActionLedger.GetAction(ctx, instanceID, runID, toolCallID)
}

// newLeaseRecoveryExecutor is newGatewayExecutor with the dead-attempt ledger in front of the gateway. The run store
// stays the plain store: only the gateway's reads see the phantom claim, as in production after a crash.
func newLeaseRecoveryExecutor(t *testing.T, exec act.ActionExecutor, ledger func(*act.PostgresRunStore) act.ActionLedger) (*act.DBOSExecutor, *act.PostgresRunStore, string) {
	t.Helper()
	_, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon una remediacion.")
	store := newRunStore(t)
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL:        dsn,
		DatabaseSchema:     schema,
		ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner:             gatewayFlowRunner{},
		Store:              store,
		Gateway: &act.Gateway{
			Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassNonIdempotent)),
			Ledger:   ledger(store),
			Executor: exec,
		},
		Proposer: fixedProposer{proposal: flowProposal(), ok: true},
		Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })
	return e, store, instanceID
}

// TestRecoveredWriteUnderLiveLeaseEndsIndeterminate: a recovered run whose write was claimed by a process that died
// must wait the dead attempt's lease out and hand the write to a human as indeterminate. It must not fail the run on
// the first ErrLeaseHeld and leave the action executing forever, with nothing that will ever revisit it.
func TestRecoveredWriteUnderLiveLeaseEndsIndeterminate(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	e, store, instanceID := newLeaseRecoveryExecutor(t, exec, func(s *act.PostgresRunStore) act.ActionLedger {
		return &deadAttemptLedger{ActionLedger: s, lease: 2 * time.Second}
	})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste subio +40%.",
		IdempotencyKey: "lease-dead-" + uuid.NewString(), Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)
	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1"))

	_, _ = e.Result(runID)

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionIndeterminate, action.Status, "a write interrupted mid-flight goes to a human, it is not left executing")
	require.Equal(t, 0, exec.executeCount(), "recovery never re-issues a write another attempt may have performed")

	run, err := store.GetRun(t.Context(), instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, run.Status)
	require.Contains(t, run.Error, "indeterminate")
}

// TestRecoveredWriteUnderLiveLeaseReportsOwnersOutcome: when the attempt holding the lease is alive after all and
// finalizes while the recovered step waits, the run reports that outcome, without a second write.
func TestRecoveredWriteUnderLiveLeaseReportsOwnersOutcome(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-2"}}
	e, store, instanceID := newLeaseRecoveryExecutor(t, exec, func(s *act.PostgresRunStore) act.ActionLedger {
		return &deadAttemptLedger{ActionLedger: s, lease: 30 * time.Second, onClaimed: func(inst, run, call string) {
			time.Sleep(500 * time.Millisecond)
			_, _ = s.TransitionAction(context.Background(), act.ActionTransition{
				InstanceID: inst, RunID: run, ToolCallID: call, Status: act.ActionSucceeded,
				ExternalReference: "PROJ-OWNER", AttemptID: deadAttemptID,
			})
		}}
	})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste subio +40%.",
		IdempotencyKey: "lease-owner-" + uuid.NewString(), Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)
	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, "PROJ-OWNER", res.ActionRef)
	require.Equal(t, 0, exec.executeCount(), "the owner's write is the only one")

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, action.Status)
}
