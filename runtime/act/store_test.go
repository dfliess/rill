package act_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// seedRun inserts a minimal parent run so child-table operations (actions, events, approvals) satisfy the FK
// constraint on run_id. CreateRun is idempotent (ON CONFLICT DO NOTHING), so calling it twice is harmless.
func seedRun(t *testing.T, store *act.PostgresRunStore, instanceID, runID string) {
	t.Helper()
	require.NoError(t, store.CreateRun(t.Context(), act.NewRun{
		RunID: runID, InstanceID: instanceID, AgentName: "seed",
	}))
}

// newRunStore returns a migrated PostgresRunStore on a schema unique to this test, so tests never see each other's
// rows and a rerun starts clean. It reuses the same Postgres as the DBOS spike (requirePostgres skips if it is down).
func newRunStore(t *testing.T) *act.PostgresRunStore {
	t.Helper()
	dsn, _ := requirePostgres(t)
	// A schema per test: sanitized UUID (bare identifier), dropped on cleanup so the dev database does not accrete.
	schema := "act_test_" + uuid.New().String()[:8]
	store, err := act.NewPostgresRunStore(t.Context(), act.StoreConfig{DatabaseURL: dsn, Schema: schema})
	require.NoError(t, err)
	require.NoError(t, store.Migrate(t.Context()))
	t.Cleanup(func() { store.DropSchemaForTest(t.Context()); store.Close() })
	return store
}

// TestRunStoreRoundTrip covers the create -> transition -> read cycle: a run is created queued, walked through
// running to succeeded, and every read reflects the latest state with a matching, ordered event stream.
func TestRunStoreRoundTrip(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-round-trip"
	runID := act.ComposeRunID(instanceID, "triage", "key-1")

	require.NoError(t, store.CreateRun(ctx, act.NewRun{
		RunID:          runID,
		InstanceID:     instanceID,
		AgentName:      "triage",
		Trigger:        "manual",
		IdempotencyKey: "key-1",
		Actor:          act.Actor{Subject: "user:alice"},
	}))

	// Immediately queryable after create: this is what lets StartAgentRun answer 202 with a real, readable run.
	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusQueued, run.Status)
	require.Equal(t, "triage", run.AgentName)
	require.Equal(t, "manual", run.Trigger)
	require.Empty(t, run.TriggerRef, "a manual run has no source resource to reference")
	require.Equal(t, "user:alice", run.Actor.Subject)
	require.Nil(t, run.StartedOn)
	require.Nil(t, run.FinishedOn)

	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning,
	}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusSucceeded, EventType: act.EventTypeSucceeded,
	}))

	run, err = store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, run.Status)
	require.NotNil(t, run.StartedOn, "started_on latches on the first running transition")
	require.NotNil(t, run.FinishedOn, "finished_on latches on the terminal transition")

	// Events are monotonic per run and ordered by the global cursor id.
	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Len(t, events, 3)
	require.Equal(t, []string{act.EventTypeQueued, act.EventTypeRunning, act.EventTypeSucceeded},
		[]string{events[0].EventType, events[1].EventType, events[2].EventType})
	require.Equal(t, []int64{1, 2, 3}, []int64{events[0].Seq, events[1].Seq, events[2].Seq})
	require.Less(t, events[0].ID, events[1].ID)
	require.Less(t, events[1].ID, events[2].ID)
}

// TestRunStoreTriggerRefRoundTrip verifies trigger_ref, the source resource an automatic trigger fired from,
// survives the create -> read cycle through both GetRun and ListRuns.
func TestRunStoreTriggerRefRoundTrip(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-trigger-ref"
	runID := act.ComposeRunID(instanceID, "triage", "alert-1")

	require.NoError(t, store.CreateRun(ctx, act.NewRun{
		RunID:      runID,
		InstanceID: instanceID,
		AgentName:  "triage",
		Trigger:    act.TriggerAlert,
		TriggerRef: "revenue_drop_alert",
	}))

	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.TriggerAlert, run.Trigger)
	require.Equal(t, "revenue_drop_alert", run.TriggerRef)

	runs, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instanceID})
	require.NoError(t, err)
	require.Len(t, runs, 1)
	require.Equal(t, "revenue_drop_alert", runs[0].TriggerRef)
}

// TestRunStoreCreateRunIdempotent verifies a replayed create neither errors nor duplicates the run or its event.
func TestRunStoreCreateRunIdempotent(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-idem"
	runID := act.ComposeRunID(instanceID, "triage", "key-1")
	newRun := act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: "triage", IdempotencyKey: "key-1"}

	require.NoError(t, store.CreateRun(ctx, newRun))
	require.NoError(t, store.CreateRun(ctx, newRun)) // replay: no error, no duplicate

	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Len(t, events, 1, "the queued event is appended once despite two creates")
}

// TestRunStoreTransitionIdempotent verifies re-applying the same transition is a no-op: no duplicate event, no error.
func TestRunStoreTransitionIdempotent(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-idem-tr"
	runID := act.ComposeRunID(instanceID, "triage", "key-1")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: "triage"}))

	tr := act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning}
	require.NoError(t, store.RecordTransition(ctx, tr))
	require.NoError(t, store.RecordTransition(ctx, tr)) // replay of the same edge

	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Len(t, events, 2, "queued + running, running not duplicated")
}

// TestRunStoreTransitionGuardsTerminal verifies a terminal run never transitions again: re-applying a stale earlier
// edge does not regress its status, and a different terminal edge does not overwrite its recorded outcome. This is
// what makes a replayed DBOS transition step safe once the run has finished.
func TestRunStoreTransitionGuardsTerminal(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-terminal"
	runID := act.ComposeRunID(instanceID, "triage", "k")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: "triage"}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusSucceeded, EventType: act.EventTypeSucceeded}))

	// A replayed earlier edge after the run is terminal is a no-op: the status stays succeeded rather than regressing
	// back to running.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning}))
	// A different terminal edge does not overwrite the recorded outcome or its (absent) error either.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusFailed, EventType: act.EventTypeFailed, Error: "should be ignored"}))

	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, run.Status)
	require.Empty(t, run.Error)

	// The guarded transitions appended no stray events.
	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Equal(t, []string{act.EventTypeQueued, act.EventTypeRunning, act.EventTypeSucceeded}, eventTypes(events))
}

// TestRunStoreTransitionSegmentScopedEdgesRepeat covers the store half of kairos-cloud#129: a run that pauses on a
// SECOND governed action must record that pause and move back to waiting_approval. The approval edges carry a
// segment-scoped DedupeKey, so a new segment's edge inserts and updates the status, while a replay of an
// already-recorded segment (same key) stays a no-op exactly like before.
func TestRunStoreTransitionSegmentScopedEdgesRepeat(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-seg-edges"
	runID := act.ComposeRunID(instanceID, "triage", "k")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: "triage"}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning}))

	// Segment 0: pause, then resume.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusWaitingApproval,
		EventType: act.EventTypeWaitingApproval, DedupeKey: act.EventTypeWaitingApproval + ":0",
	}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning,
		EventType: act.EventTypeResumed, DedupeKey: act.EventTypeResumed + ":0",
	}))

	// Segment 1: the second pause — the edge the old (run_id, event_type) uniqueness silently swallowed, leaving the
	// run stuck on running while it waited for a decision nobody could see.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusWaitingApproval,
		EventType: act.EventTypeWaitingApproval, DedupeKey: act.EventTypeWaitingApproval + ":1",
	}))
	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status, "the second pause must move the status, not conflict into a no-op")

	// A replay of the SAME segment's pause (a re-applied DBOS step) is still idempotent: no event, no status churn.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusWaitingApproval,
		EventType: act.EventTypeWaitingApproval, DedupeKey: act.EventTypeWaitingApproval + ":1",
	}))

	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning,
		EventType: act.EventTypeResumed, DedupeKey: act.EventTypeResumed + ":1",
	}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusSucceeded, EventType: act.EventTypeSucceeded}))

	// The timeline shows both pauses and both resumptions, once each, in order.
	events, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Equal(t, []string{
		act.EventTypeQueued,
		act.EventTypeRunning,
		act.EventTypeWaitingApproval,
		act.EventTypeResumed,
		act.EventTypeWaitingApproval,
		act.EventTypeResumed,
		act.EventTypeSucceeded,
	}, eventTypes(events))
}

// TestRunStoreMigrateMovesEventUniquenessToDedupeKey exercises the startup migration against a table with the
// PRE-dedupe_key shape (idempotency on UNIQUE (run_id, event_type)): the backfill must key existing rows without
// losing any — suffixing the approval edges with ":0", the segment every pre-migration run implicitly was on — and
// the moved uniqueness must admit a second segment's pause while still deduplicating a replayed edge.
func TestRunStoreMigrateMovesEventUniquenessToDedupeKey(t *testing.T) {
	dsn, _ := requirePostgres(t)
	ctx := t.Context()
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// Recreate agent_run_events exactly as an earlier build's Migrate left it, seeded with a run parked on its first
	// approval gate. Postgres names the inline constraint agent_run_events_run_id_event_type_key, the name the
	// migration drops. agent_runs is also created (any real older build had it) so the events are not orphans — the FK
	// migration cleans up orphan rows before validating the constraint.
	schema := "act_test_" + uuid.New().String()[:8]
	runs := schema + ".agent_runs"
	events := schema + ".agent_run_events"
	_, err = pool.Exec(ctx, `CREATE SCHEMA `+schema)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `CREATE TABLE `+runs+` (
		run_id text PRIMARY KEY,
		instance_id text NOT NULL,
		organization_id text,
		project_id text,
		agent_name text NOT NULL,
		spec_hash text NOT NULL DEFAULT '',
		trigger text NOT NULL DEFAULT '',
		trigger_ref text NOT NULL DEFAULT '',
		idempotency_key text NOT NULL DEFAULT '',
		conversation_id text NOT NULL DEFAULT '',
		actor_subject text NOT NULL DEFAULT '',
		actor_service_principal boolean NOT NULL DEFAULT false,
		status text NOT NULL,
		error text NOT NULL DEFAULT '',
		created_on timestamptz NOT NULL DEFAULT now(),
		updated_on timestamptz NOT NULL DEFAULT now(),
		started_on timestamptz,
		finished_on timestamptz
	)`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `CREATE TABLE `+events+` (
		id bigserial PRIMARY KEY,
		run_id text NOT NULL,
		instance_id text NOT NULL,
		seq bigint NOT NULL,
		event_type text NOT NULL,
		status text NOT NULL DEFAULT '',
		payload jsonb,
		visibility text NOT NULL DEFAULT 'user',
		created_on timestamptz NOT NULL DEFAULT now(),
		UNIQUE (run_id, seq),
		UNIQUE (run_id, event_type)
	)`)
	require.NoError(t, err)

	const instanceID = "inst-migrate"
	runID := act.ComposeRunID(instanceID, "triage", "k")
	_, err = pool.Exec(ctx, `INSERT INTO `+runs+` (run_id, instance_id, agent_name, status) VALUES ($1, $2, 'triage', $3)`,
		runID, instanceID, string(act.RunStatusWaitingApproval))
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO `+events+` (run_id, instance_id, seq, event_type, status) VALUES
		($1, $2, 1, $3, 'queued'), ($1, $2, 2, $4, 'running'), ($1, $2, 3, $5, 'waiting_approval')`,
		runID, instanceID, act.EventTypeQueued, act.EventTypeRunning, act.EventTypeWaitingApproval)
	require.NoError(t, err)

	store, err := act.NewPostgresRunStore(ctx, act.StoreConfig{DatabaseURL: dsn, Schema: schema})
	require.NoError(t, err)
	t.Cleanup(func() { store.DropSchemaForTest(ctx); store.Close() })
	require.NoError(t, store.Migrate(ctx))
	require.NoError(t, store.Migrate(ctx), "the migration runs on every startup, so it must be re-runnable")

	// Every pre-existing event survived, keyed by its type — the approval edge scoped to segment 0 so a replayed old
	// edge deduplicates against the keys the new executor emits.
	rows, err := pool.Query(ctx, `SELECT event_type, dedupe_key FROM `+events+` WHERE run_id=$1 ORDER BY seq`, runID)
	require.NoError(t, err)
	defer rows.Close()
	backfilled := map[string]string{}
	for rows.Next() {
		var eventType, dedupeKey string
		require.NoError(t, rows.Scan(&eventType, &dedupeKey))
		backfilled[eventType] = dedupeKey
	}
	require.NoError(t, rows.Err())
	require.Equal(t, map[string]string{
		act.EventTypeQueued:          act.EventTypeQueued,
		act.EventTypeRunning:         act.EventTypeRunning,
		act.EventTypeWaitingApproval: act.EventTypeWaitingApproval + ":0",
	}, backfilled)

	// A replayed pre-migration edge still deduplicates: same segment-0 key, no new event.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusWaitingApproval,
		EventType: act.EventTypeWaitingApproval, DedupeKey: act.EventTypeWaitingApproval + ":0",
	}))
	migrated, err := store.ListRunEvents(ctx, instanceID, runID, 0, 100)
	require.NoError(t, err)
	require.Len(t, migrated, 3, "a replayed old edge must not duplicate after the migration")

	// The run continues where it parked: resume segment 0, then pause on a SECOND action. The second waiting_approval
	// shares its event_type with the backfilled row, so this insert also proves the old constraint is gone.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning,
		EventType: act.EventTypeResumed, DedupeKey: act.EventTypeResumed + ":0",
	}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusWaitingApproval,
		EventType: act.EventTypeWaitingApproval, DedupeKey: act.EventTypeWaitingApproval + ":1",
	}))
	run, err := store.GetRun(ctx, instanceID, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status)
}

// TestRunStoreCancelPendingApprovals verifies that cancelling a run withdraws its pending approval (so the inbox stops
// offering a decision) while leaving an already-decided approval, and another run's approval, untouched.
func TestRunStoreCancelPendingApprovals(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-appr-cancel"
	pendingRun := act.ComposeRunID(instanceID, "triage", "pending")
	decidedRun := act.ComposeRunID(instanceID, "triage", "decided")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: pendingRun, InstanceID: instanceID, AgentName: "triage"}))
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: decidedRun, InstanceID: instanceID, AgentName: "triage"}))

	pendingApproval := act.ApprovalIDForRun(pendingRun)
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: pendingApproval, RunID: pendingRun, InstanceID: instanceID,
		ArgsHash: act.HashArgs("x"), Proposal: "x", RequestedBy: "user:alice",
	}))
	// A second run whose approval was already approved: cancelling the first run must not touch it.
	decidedApproval := act.ApprovalIDForRun(decidedRun)
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: decidedApproval, RunID: decidedRun, InstanceID: instanceID,
		ArgsHash: act.HashArgs("y"), Proposal: "y", RequestedBy: "user:alice",
	}))
	_, err := store.ResolveApproval(ctx, instanceID, decidedApproval, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)

	// Cancelling the first run withdraws exactly its one pending approval.
	n, err := store.CancelPendingApprovals(ctx, instanceID, pendingRun)
	require.NoError(t, err)
	require.Equal(t, int64(1), n)

	cancelled, err := store.GetApproval(ctx, instanceID, pendingApproval)
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusCancelled, cancelled.Status)

	// The already-approved approval is left alone.
	untouched, err := store.GetApproval(ctx, instanceID, decidedApproval)
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, untouched.Status)

	// Idempotent: with nothing pending, a second sweep cancels zero.
	n, err = store.CancelPendingApprovals(ctx, instanceID, pendingRun)
	require.NoError(t, err)
	require.Equal(t, int64(0), n)
}

// TestRunStoreScopingIsolatesTenants is the isolation test (§15.2): two instances create runs, and neither can see
// nor read the other's, whether listing or fetching by ID.
func TestRunStoreScopingIsolatesTenants(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instA, instB = "inst-a", "inst-b"
	runA := act.ComposeRunID(instA, "triage", "k")
	runB := act.ComposeRunID(instB, "triage", "k")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runA, InstanceID: instA, AgentName: "triage"}))
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runB, InstanceID: instB, AgentName: "triage"}))

	// Each instance lists only its own run.
	listA, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instA})
	require.NoError(t, err)
	require.Len(t, listA, 1)
	require.Equal(t, runA, listA[0].RunID)

	listB, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instB})
	require.NoError(t, err)
	require.Len(t, listB, 1)
	require.Equal(t, runB, listB[0].RunID)

	// A cross-instance fetch fails as not-found, indistinguishable from a truly missing run: instance B cannot read
	// instance A's run even with the exact run ID.
	_, err = store.GetRun(ctx, instB, runA)
	require.ErrorIs(t, err, act.ErrRunNotFound)

	// A cross-instance transition also fails closed, so state cannot be mutated across tenants.
	err = store.RecordTransition(ctx, act.RunTransition{InstanceID: instB, RunID: runA, Status: act.RunStatusRunning, EventType: act.EventTypeRunning})
	require.ErrorIs(t, err, act.ErrRunNotFound)

	// Events are unreadable across tenants too.
	events, err := store.ListRunEvents(ctx, instB, runA, 0, 100)
	require.NoError(t, err)
	require.Empty(t, events)
}

// TestRunStoreListFilters verifies the agent and status filters narrow within an instance.
func TestRunStoreListFilters(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-filters"
	mk := func(agent, key string, status act.RunStatus) string {
		runID := act.ComposeRunID(instanceID, agent, key)
		require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: agent}))
		if status != act.RunStatusQueued {
			require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: status, EventType: string(status)}))
		}
		return runID
	}
	mk("triage", "1", act.RunStatusQueued)
	mk("triage", "2", act.RunStatusSucceeded)
	mk("summary", "3", act.RunStatusQueued)

	byAgent, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instanceID, AgentName: "triage"})
	require.NoError(t, err)
	require.Len(t, byAgent, 2)

	byStatus, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instanceID, Status: act.RunStatusQueued})
	require.NoError(t, err)
	require.Len(t, byStatus, 2)
}

// TestRunStoreAccessFilterInWhere pins that the access allow-list restricts inside the WHERE clause, not by
// filtering a fetched page: with more inaccessible candidates than the page size, the page still comes back
// full of accessible rows. In-memory filtering would return the newest (inaccessible) rows and then drop them,
// yielding an empty page and lying pagination.
func TestRunStoreAccessFilterInWhere(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-access-filter"
	mk := func(agent, key string) string {
		runID := act.ComposeRunID(instanceID, agent, key)
		require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: agent}))
		return runID
	}
	// The accessible agent's runs are the OLDEST: a newest-first page of 2 with no WHERE restriction would
	// contain only "triage" rows.
	cob1 := mk("cobranza", "1")
	cob2 := mk("cobranza", "2")
	time.Sleep(20 * time.Millisecond) // separate created_on so the newest-first order is deterministic
	mk("triage", "3")
	mk("triage", "4")
	mk("triage", "5")

	// The page comes back full of accessible rows despite three newer inaccessible candidates.
	page, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instanceID, AccessibleAgents: []string{"cobranza"}, Limit: 2})
	require.NoError(t, err)
	require.Len(t, page, 2)
	require.ElementsMatch(t, []string{cob1, cob2}, []string{page[0].RunID, page[1].RunID})

	// A non-nil empty allow-list matches nothing; nil applies no restriction.
	none, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instanceID, AccessibleAgents: []string{}})
	require.NoError(t, err)
	require.Empty(t, none)
	all, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: instanceID})
	require.NoError(t, err)
	require.Len(t, all, 5)

	// Approvals resolve the allow-list through their run.
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: act.ApprovalIDForRun(cob1), RunID: cob1, InstanceID: instanceID, ArgsHash: act.HashArgs("x"),
	}))
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: act.ApprovalIDForRun(act.ComposeRunID(instanceID, "triage", "3")), RunID: act.ComposeRunID(instanceID, "triage", "3"), InstanceID: instanceID, ArgsHash: act.HashArgs("y"),
	}))
	approvals, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: instanceID, AccessibleAgents: []string{"cobranza"}})
	require.NoError(t, err)
	require.Len(t, approvals, 1)
	require.Equal(t, cob1, approvals[0].RunID)
	noApprovals, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: instanceID, AccessibleAgents: []string{}})
	require.NoError(t, err)
	require.Empty(t, noApprovals)
}

// TestRunStoreApprovalLifecycle covers create -> read -> resolve, and the double-submit guard: only the first
// resolution of a pending approval wins.
func TestRunStoreApprovalLifecycle(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()

	const instanceID = "inst-appr"
	runID := act.ComposeRunID(instanceID, "triage", "k")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: instanceID, AgentName: "triage"}))

	approvalID := act.ApprovalIDForRun(runID)
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID:  approvalID,
		RunID:       runID,
		InstanceID:  instanceID,
		ToolName:    "act.propose_action",
		ArgsHash:    act.HashArgs("crear ticket P2"),
		Proposal:    "crear ticket P2",
		RequestedBy: "user:alice",
	}))

	got, err := store.GetApproval(ctx, instanceID, approvalID)
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusPending, got.Status)
	require.Equal(t, act.HashArgs("crear ticket P2"), got.ArgsHash)
	// The projection joins the approval's run: the agent and the run's actor ride along, so a reader can
	// resolve the approve policy for this concrete approval without a second lookup.
	require.Equal(t, "triage", got.AgentName)
	require.Empty(t, got.RunActorSubject)

	// Only pending approvals appear in the inbox filter.
	pending, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Len(t, pending, 1)

	// First resolution wins; the second sees a non-pending approval.
	resolved, err := store.ResolveApproval(ctx, instanceID, approvalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, resolved.Status)
	require.Equal(t, "admin:bob", resolved.DecidedBy)
	require.NotNil(t, resolved.DecidedOn)

	_, err = store.ResolveApproval(ctx, instanceID, approvalID, act.ApprovalStatusDenied, "admin:carol")
	require.ErrorIs(t, err, act.ErrApprovalNotResolvable)

	// Cross-tenant resolution fails closed.
	_, err = store.ResolveApproval(ctx, "other-inst", approvalID, act.ApprovalStatusApproved, "x")
	require.ErrorIs(t, err, act.ErrApprovalNotFound)
}
