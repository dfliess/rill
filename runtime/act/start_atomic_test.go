package act_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// failAfterEnqueue lets the real store create the run and enqueue its workflow in one transaction, then fails before
// the commit: the moment a crash between the two writes used to leave a queued run with no workflow behind it.
type failAfterEnqueue struct {
	*act.PostgresRunStore
	err   error
	runID string // the ID Start tried to create the run under
}

func (s *failAfterEnqueue) CreateRunAndEnqueue(ctx context.Context, run act.NewRun, enqueue func(context.Context, pgx.Tx) error) (string, error) {
	s.runID = run.RunID
	return s.PostgresRunStore.CreateRunAndEnqueue(ctx, run, func(ctx context.Context, tx pgx.Tx) error {
		if err := enqueue(ctx, tx); err != nil {
			return err
		}
		return s.err
	})
}

// TestStartCreatesRunAndWorkflowTogether: Start records the run and launches its workflow as one transaction
// (ADR-0032). A failure between the two leaves neither: no run the API would show forever as queued, and no workflow
// the store knows nothing about.
func TestStartCreatesRunAndWorkflowTogether(t *testing.T) {
	store := newRunStore(t)
	errCrash := errors.New("process died between the run and its workflow")
	failing := &failAfterEnqueue{PostgresRunStore: store, err: errCrash}
	e := newStoreExecutor(t, gatewayFlowRunner{}, failing)
	ctx := t.Context()

	const inst = "inst-start-atomic"
	key := "start-atomic-" + uuid.NewString()
	_, err := e.Start(ctx, act.AgentRunInput{InstanceID: inst, AgentName: "triage", Prompt: "alerta", IdempotencyKey: key})
	require.ErrorIs(t, err, errCrash)

	runID := failing.runID
	require.NotEmpty(t, runID)
	_, err = store.GetRun(ctx, inst, runID)
	require.ErrorIs(t, err, act.ErrRunNotFound, "a run whose workflow was never launched must not exist")

	dsn, schema := requirePostgres(t)
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()
	var n int
	require.NoError(t, pool.QueryRow(ctx, fmt.Sprintf(`SELECT count(*) FROM %s.workflow_status WHERE workflow_uuid = $1`, schema), runID).Scan(&n))
	require.Zero(t, n, "a workflow whose run was never recorded must not exist")
}

// TestParkedRunsDoNotBlockTheNextOne: a run waiting on a human counts as PENDING for DBOS, so a queue with a
// concurrency limit would stop dispatching once that many runs were parked (ADR-0022, decision 3). The Act queue has
// no limit: with several runs parked, the next one still starts.
func TestParkedRunsDoNotBlockTheNextOne(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-Q"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})
	ctx := t.Context()

	const parked = 12
	for i := 0; i <= parked; i++ {
		runID, err := e.Start(ctx, act.AgentRunInput{
			InstanceID: instanceID, AgentName: "triage", Prompt: "alerta",
			IdempotencyKey: fmt.Sprintf("queue-%d-%s", i, uuid.NewString()), Actor: act.Actor{Subject: "user:alice"},
		})
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			run, gErr := store.GetRun(ctx, instanceID, runID)
			return gErr == nil && run.Status == act.RunStatusWaitingApproval
		}, 20*time.Second, 50*time.Millisecond,
			"run %d never started: runs parked on approval blocked the queue (a concurrency limit counts every PENDING run on it, of any version)", i)
	}
}
