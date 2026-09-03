package act_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// TestPatchGuardsShapeChangeAgainstParkedRuns pins the rule that keeps a deploy from killing runs parked on a human
// approval, written after 2026-09-02 did exactly that to 19 of them. Act pins a constant DBOS ApplicationVersion so a
// parked run survives the deploy at all (see defaultApplicationVersion), which means the new binary replays
// checkpoints recorded by older code — and DBOS then insists the step sequence matches. A new step deployed in front
// of the approval Recv landed in the slot the parked Recv had not recorded yet, the Recv itself landed on the
// recorded "DBOS.sleep", and every parked run died with UnexpectedStep.
//
// The rule has two halves, and this test pins both:
//
//  1. A change to the workflow's step sequence goes behind a named dbos.Patch: runs recorded before the patch point
//     keep the old shape, runs born after it take the new one.
//  2. The Patch call sits where every in-flight run already has a checkpoint in the NEXT slot — here, before the
//     workflow's first step — not at the point of the change. Patch decides by reading that next slot (an empty one
//     means "new workflow", and it takes the new branch). A run parked on Recv has recorded nothing in the Recv's own
//     slot — only its timeout sleep, one slot later — so a Patch placed just before the wait reads exactly that empty
//     slot and steers the PARKED run into the new branch, reproducing the incident with the patch in place.
//
// Both halves are pinned by mutation: strip the Patch (the incident's shape) or move it to the point of the change
// (the placement DBOS's own docs suggest) and phase 2 fails with the incident's UnexpectedStep error.
//
// The workflow here is a hand-rolled miniature rather than the executor's real one because the test needs to run two
// SHAPES of the same workflow against one database, and the executor's shape is compiled in. What carries over is the
// exact structure that failed: steps, then a parked Recv, then a deploy that adds a step before the wait.
func TestPatchGuardsShapeChangeAgainstParkedRuns(t *testing.T) {
	dsn, schema := requirePostgres(t)

	const topic = "approval"
	wfName := "patch_rule_probe"
	// Both phases share one version, as production's pinned constant does: that is what makes phase 2 recover
	// phase 1's parked run instead of ignoring it (and what isolates this test from the suite's other runs).
	version := "act-patch-rule-" + uuid.NewString()

	newDBOS := func(register func(dbos.Context)) dbos.Context {
		dctx, err := dbos.NewContext(context.Background(), dbos.Config{
			AppName:            "act-patch-rule-test",
			DatabaseURL:        dsn,
			DatabaseSchema:     schema,
			ApplicationVersion: version,
			EnablePatching:     true, // as NewDBOSExecutor arms it
		})
		require.NoError(t, err)
		register(dctx)
		require.NoError(t, dbos.Launch(dctx))
		return dctx
	}

	prepare := func(context.Context) (string, error) { return "prepared", nil }

	// Shape A — the binary that started the run: one step, then park on the approval Recv.
	shapeA := func(ctx dbos.Context, _ string) (string, error) {
		if _, err := dbos.RunAsStep(ctx, prepare, dbos.WithStepName("prepare")); err != nil {
			return "", err
		}
		decision, err := dbos.Recv[string](ctx, topic, time.Hour)
		if err != nil {
			return "", err
		}
		return "old:" + decision, nil
	}

	// Shape B — the next deploy: an "announce" step now precedes the wait, gated behind a patch. The patch is read
	// before the FIRST step, where a run parked anywhere in this workflow has already recorded a checkpoint in the
	// slot Patch inspects; its verdict then steers the change further down. announced counts executions of the new
	// step so the test can tell the recovered run did not take it.
	var announced atomic.Int32
	shapeB := func(ctx dbos.Context, _ string) (string, error) {
		patched, err := dbos.Patch(ctx, "announce-before-wait")
		if err != nil {
			return "", err
		}
		if _, err := dbos.RunAsStep(ctx, prepare, dbos.WithStepName("prepare")); err != nil {
			return "", err
		}
		if patched {
			if _, err := dbos.RunAsStep(ctx, func(context.Context) (string, error) {
				announced.Add(1)
				return "announced", nil
			}, dbos.WithStepName("announce")); err != nil {
				return "", err
			}
		}
		decision, err := dbos.Recv[string](ctx, topic, time.Hour)
		if err != nil {
			return "", err
		}
		if patched {
			return "new:" + decision, nil
		}
		return "old:" + decision, nil
	}

	// Phase 1: the old binary starts a run, the run parks on its approval Recv, and the process drains gracefully —
	// a deploy, as the executor's shutdown tests establish it, leaving the workflow PENDING for the next process.
	oldProc := newDBOS(func(dctx dbos.Context) {
		dbos.RegisterWorkflow(dctx, shapeA, dbos.WithWorkflowName(wfName))
	})
	parkedID := "patch-rule-parked-" + uuid.NewString()
	_, err := dbos.RunWorkflow(oldProc, shapeA, "investiga", dbos.WithWorkflowID(parkedID))
	require.NoError(t, err)
	awaitParkedOnRecv(t, dsn, schema, parkedID)
	require.NoError(t, dbos.Shutdown(oldProc, 10*time.Second))
	require.Equal(t, "PENDING", workflowStatus(t, dsn, schema, parkedID),
		"the drain must leave the parked run for the next process, or there is nothing to recover")

	// Phase 2: the new binary comes up with shape B under the same workflow name and version; Launch recovers the
	// parked run and replays it against the new code. The approval is durable, so it is safe to send before the
	// replay reaches the wait — exactly how a decision races a recovering runtime in production.
	newProc := newDBOS(func(dctx dbos.Context) {
		dbos.RegisterWorkflow(dctx, shapeB, dbos.WithWorkflowName(wfName))
	})
	t.Cleanup(func() { _ = dbos.Shutdown(newProc, 10*time.Second) })

	require.NoError(t, dbos.Send(newProc, parkedID, "approved", topic))
	parked, err := dbos.RetrieveWorkflow[string](newProc, parkedID)
	require.NoError(t, err)
	got, err := parked.GetResult()
	require.NoError(t, err,
		"a run parked before the deploy must survive the new shape (this is where the 2026-09-02 UnexpectedStep surfaces)")
	require.Equal(t, "old:approved", got, "a pre-patch run stays on the branch its checkpoints were recorded under")
	require.Equal(t, int32(0), announced.Load(), "the recovered run must not execute the step its version never had")

	// A run born after the deploy takes the new branch, which is the other half of Patch's contract — without it the
	// green above could just mean the new code is unreachable.
	freshID := "patch-rule-fresh-" + uuid.NewString()
	fresh, err := dbos.RunWorkflow(newProc, shapeB, "investiga", dbos.WithWorkflowID(freshID))
	require.NoError(t, err)
	require.NoError(t, dbos.Send(newProc, freshID, "approved", topic))
	got, err = fresh.GetResult()
	require.NoError(t, err)
	require.Equal(t, "new:approved", got, "a run born after the deploy takes the patched branch")
	require.Equal(t, int32(1), announced.Load(), "the new step runs for the new run, and only for it")
}

// awaitParkedOnRecv waits until the workflow has checkpointed its Recv timeout as a "DBOS.sleep" step, the durable
// trace Recv writes immediately before suspending. That checkpoint is what makes the parked state reproducible: the
// Recv's own slot stays unrecorded while the sleep sits one slot later, which is the layout the whole patching rule
// is about.
func awaitParkedOnRecv(t *testing.T, dsn, schema, workflowID string) {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	defer pool.Close()

	require.Eventually(t, func() bool {
		var n int
		err := pool.QueryRow(context.Background(),
			fmt.Sprintf("SELECT count(*) FROM %s.operation_outputs WHERE workflow_uuid = $1 AND function_name = 'DBOS.sleep'", schema),
			workflowID).Scan(&n)
		return err == nil && n > 0
	}, 15*time.Second, 50*time.Millisecond, "workflow %q never parked on its Recv", workflowID)
}
