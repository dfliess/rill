package act_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// shutdownInsideApprovalSetup holds the approval setup step open until the process drains, so a shutdown lands inside
// the step deterministically instead of by timing (kairos-cloud#170, problem 1). It stands in for a deploy arriving
// while a run is entering its approval wait.
type shutdownInsideApprovalSetup struct {
	act.RunStore
	entered chan struct{}
}

func (s *shutdownInsideApprovalSetup) CreateApproval(ctx context.Context, a act.NewApproval) error {
	close(s.entered)
	<-ctx.Done()
	return ctx.Err()
}

func newShutdownExecutor(t *testing.T, store act.RunStore, ledger act.ActionLedger, version string) *act.DBOSExecutor {
	t.Helper()
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL:        dsn,
		DatabaseSchema:     schema,
		ApplicationVersion: version,
		Runner:             gatewayFlowRunner{},
		Store:              store,
		Gateway: &act.Gateway{
			Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
			Ledger:   ledger,
			Executor: &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-9"}},
		},
		Proposer: fixedProposer{proposal: flowProposal(), ok: true},
		Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	return e
}

// TestShutdownDuringApprovalSetupLeavesRunRecoverable: a deploy that interrupts the step creating a run's approvals
// must leave the run for the next process to recover, not file it as failed. DBOS does not checkpoint an interrupted
// step and re-runs it on recovery, so a run marked failed on the way out can never ask for its signature again: the
// recovered step's CreateApproval refuses to attach to a terminal run.
func TestShutdownDuringApprovalSetupLeavesRunRecoverable(t *testing.T) {
	_, inst := newInstanceWithAgent(t, "Investiga la alerta y propon una remediacion.")
	store := newRunStore(t)
	version := "act-test-" + uuid.NewString()
	ctx := t.Context()

	blocking := &shutdownInsideApprovalSetup{RunStore: store, entered: make(chan struct{})}
	first := newShutdownExecutor(t, blocking, store, version)
	runID, err := first.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "El coste subio +40%.",
		IdempotencyKey: "shutdown-setup-" + uuid.NewString(), Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	select {
	case <-blocking.entered:
	case <-time.After(15 * time.Second):
		t.Fatal("the run never reached its approval setup")
	}
	first.Close(10 * time.Second)

	require.Never(t, func() bool {
		run, gErr := store.GetRun(ctx, inst, runID)
		return gErr != nil || run.Status.IsTerminal()
	}, 2*time.Second, 200*time.Millisecond, "a shutdown is not a failure: the run must stay open for the next process")

	// The next process recovers the run into the same step, which now completes and parks on the approval.
	second := newShutdownExecutor(t, store, store, version)
	t.Cleanup(func() { second.Close(5 * time.Second) })

	approvalID := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		a, gErr := store.GetApproval(ctx, inst, approvalID)
		return gErr == nil && a.Status == act.ApprovalStatusPending
	}, 15*time.Second, 100*time.Millisecond, "the recovered run must be able to ask for its signature")

	require.NoError(t, second.Resume(ctx, runID, act.ApprovalApproved, "call-1"))
	res, err := second.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, "PROJ-9", res.ActionRef)
}
