package trigger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
)

// requireStore returns a migrated Store in a throwaway schema, skipping the test if Postgres is unreachable. It
// reuses WS1's local container (docker start dbos-spike-pg) unless overridden by env. Each test gets its own schema
// so tests never see each other's rows.
func requireStore(t *testing.T) *Store {
	t.Helper()
	dsn := os.Getenv("ACT_TEST_DBOS_URL")
	if dsn == "" {
		dsn = "postgres://postgres:dbos@localhost:55432/dbos_spike?sslmode=disable"
	}

	u, err := url.Parse(dsn)
	require.NoError(t, err)
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		t.Skipf("act/trigger tests need Postgres at %s (start it with: docker start dbos-spike-pg): %v", u.Host, err)
	}
	_ = conn.Close()

	pool, err := pgxpool.New(t.Context(), dsn)
	require.NoError(t, err)

	schema := "act_trigger_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	store := NewStore(pool, schema)
	require.NoError(t, store.Migrate(t.Context()))

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", store.schema))
		pool.Close()
	})
	return store
}

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// fakeExecutor is an act.AgentExecutor that records every Start and, like the real DBOS executor, treats a repeated
// composed run ID as the same run (OAOO). A test asserts on the number of distinct runs it actually created.
type fakeExecutor struct {
	mu     sync.Mutex
	starts []act.AgentRunInput
	runs   map[string]bool
}

var _ act.AgentExecutor = (*fakeExecutor)(nil)

func (f *fakeExecutor) Start(_ context.Context, in act.AgentRunInput) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	runID := act.ComposeRunID(in.InstanceID, in.AgentName, in.IdempotencyKey)
	f.starts = append(f.starts, in)
	if f.runs == nil {
		f.runs = map[string]bool{}
	}
	f.runs[runID] = true
	return runID, nil
}

func (f *fakeExecutor) Resume(context.Context, string, act.ApprovalDecision, string) error {
	return nil
}
func (f *fakeExecutor) Resumable(context.Context, string) (bool, error) { return true, nil }
func (f *fakeExecutor) Cancel(context.Context, string, string) error    { return nil }

func (f *fakeExecutor) startCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.starts)
}

func (f *fakeExecutor) distinctRuns() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.runs)
}

func (f *fakeExecutor) firstStart() act.AgentRunInput {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.starts[0]
}

// fakeCatalog is a static TriggerCatalog keyed by instance, so the dispatcher's match/dedup logic can be tested
// without a running runtime.
type fakeCatalog struct {
	byInstance map[string][]TriggerDef
}

func (f *fakeCatalog) ListTriggers(_ context.Context, instanceID string) ([]TriggerDef, error) {
	return f.byInstance[instanceID], nil
}

// enqueueAlertEvent derives-and-enqueues one alert event directly (bypassing the reconciler) so a dispatcher test
// controls exactly what is in the outbox. It returns the event.
func enqueueAlertEvent(t *testing.T, store *Store, instanceID, alertName, eventType string, execTime time.Time) Event {
	t.Helper()
	e := newEvent(instanceID, resourceKindAlert, alertName, eventType, "fail", execTime, execTime)
	_, err := store.EnqueueEvent(t.Context(), e)
	require.NoError(t, err)
	return e
}

func resetOutboxPending(t *testing.T, store *Store, instanceID string) {
	t.Helper()
	_, err := store.pool.Exec(t.Context(),
		fmt.Sprintf("UPDATE %s.agent_trigger_outbox SET state='pending', leased_until=NULL, available_on=now() WHERE instance_id=$1", store.schema),
		instanceID)
	require.NoError(t, err)
}

func countEvents(t *testing.T, store *Store, instanceID, triggerID, outcome string) int {
	t.Helper()
	var n int
	err := store.pool.QueryRow(t.Context(),
		fmt.Sprintf("SELECT count(*) FROM %s.agent_trigger_events WHERE instance_id=$1 AND trigger_id=$2 AND outcome=$3", store.schema),
		instanceID, triggerID, outcome).Scan(&n)
	require.NoError(t, err)
	return n
}

func countOutbox(t *testing.T, store *Store, instanceID, state string) int {
	t.Helper()
	var n int
	err := store.pool.QueryRow(t.Context(),
		fmt.Sprintf("SELECT count(*) FROM %s.agent_trigger_outbox WHERE instance_id=$1 AND state=$2", store.schema),
		instanceID, state).Scan(&n)
	require.NoError(t, err)
	return n
}

func alertTrigger(name, agent, alertName string, events ...string) TriggerDef {
	return TriggerDef{
		Name:         name,
		Agent:        agent,
		SourceKind:   resourceKindAlert,
		SourceName:   alertName,
		Events:       events,
		ActorService: true,
		ActorSubject: "svc_act",
		Prompt:       "investiga la alerta",
	}
}

// TestTriggerMatches covers the pure match predicate, including the ADR-0016 wildcard: an empty SourceName matches
// any resource of the trigger's kind, while a named source matches only that resource. It needs no Postgres.
func TestTriggerMatches(t *testing.T) {
	alertEvent := OutboxEvent{ResourceKind: resourceKindAlert, ResourceName: "revenue_drop", EventType: EventAlertEnteredFail}
	short := shortEventName(alertEvent.EventType) // "entered_fail"

	cases := []struct {
		name string
		tr   TriggerDef
		want bool
	}{
		{
			name: "named source matches its own alert",
			tr:   TriggerDef{SourceKind: resourceKindAlert, SourceName: "revenue_drop", Events: []string{"entered_fail"}},
			want: true,
		},
		{
			name: "named source does not match a different alert",
			tr:   TriggerDef{SourceKind: resourceKindAlert, SourceName: "signups_drop", Events: []string{"entered_fail"}},
			want: false,
		},
		{
			name: "wildcard source matches any alert of that kind",
			tr:   TriggerDef{SourceKind: resourceKindAlert, SourceName: "", Events: []string{"entered_fail"}},
			want: true,
		},
		{
			name: "wildcard still requires the kind to match",
			tr:   TriggerDef{SourceKind: resourceKindReport, SourceName: "", Events: []string{"entered_fail"}},
			want: false,
		},
		{
			name: "wildcard still requires a subscribed event",
			tr:   TriggerDef{SourceKind: resourceKindAlert, SourceName: "", Events: []string{"recovered"}},
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, triggerMatches(tc.tr, alertEvent, short))
		})
	}
}

// TestTriggerMatches_Schedule covers the schedule trigger's self-referential match: it fires only on its own tick
// (by trigger name), and two schedules with the same cron never cross-fire. It needs no Postgres.
func TestTriggerMatches_Schedule(t *testing.T) {
	tick := OutboxEvent{ResourceKind: resourceKindSchedule, ResourceName: "nightly", EventType: EventScheduleTick}
	short := shortEventName(tick.EventType) // "tick"

	nightly := TriggerDef{Name: "nightly", SourceKind: resourceKindSchedule, Cron: "0 3 * * *"}
	other := TriggerDef{Name: "hourly", SourceKind: resourceKindSchedule, Cron: "0 * * * *"}

	require.True(t, triggerMatches(nightly, tick, short), "a schedule matches its own tick")
	require.False(t, triggerMatches(other, tick, short), "a different schedule with the same cron does not cross-fire")

	// A schedule trigger must not be activated by an alert event, even a wildcard-shaped one.
	alertTick := OutboxEvent{ResourceKind: resourceKindAlert, ResourceName: "nightly", EventType: EventAlertEnteredFail}
	require.False(t, triggerMatches(nightly, alertTick, shortEventName(alertTick.EventType)))
}

// TestEnqueueDueSchedules drives the schedule ticker end to end against Postgres: the first scan of an instance
// backfills nothing, a scan whose window crosses a cron boundary enqueues exactly one tick, re-scanning does not
// double-enqueue, and Dispatch turns the tick into exactly one run whose provenance is truthfully "schedule".
func TestEnqueueDueSchedules(t *testing.T) {
	store := requireStore(t)
	const instanceID = "inst_sched"

	tr := TriggerDef{
		Name:       "nightly",
		Agent:      "revenue-agent",
		SourceKind: resourceKindSchedule,
		Cron:       "0 3 * * *", // 03:00 every day
		Prompt:     "corre el resumen diario",
	}
	cat := &fakeCatalog{byInstance: map[string][]TriggerDef{instanceID: {tr}}}
	exec := &fakeExecutor{}

	// A hand-cranked clock so the test controls exactly which cron boundaries fall in each scan window.
	var now time.Time
	d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger(), Now: func() time.Time { return now }})

	// First scan just before the boundary: it only sets the high-water mark, so no past boundary is backfilled.
	now = time.Date(2026, 1, 1, 2, 59, 0, 0, time.UTC)
	n, err := d.EnqueueDueSchedules(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 0, n, "first scan backfills nothing")
	require.Equal(t, 0, countOutbox(t, store, instanceID, "pending"))

	// Scan whose window (02:59, 03:00:30] crosses the 03:00 boundary: exactly one tick.
	now = time.Date(2026, 1, 1, 3, 0, 30, 0, time.UTC)
	n, err = d.EnqueueDueSchedules(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 1, n, "the crossed boundary enqueues one tick")
	require.Equal(t, 1, countOutbox(t, store, instanceID, "pending"))

	// Re-scanning at the same instant enqueues nothing new: the window has advanced past the boundary and the
	// event_id is stable per boundary, so even an overlapping re-scan would collapse on (instance, event_id).
	n, err = d.EnqueueDueSchedules(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 0, n, "no double-enqueue for the same boundary")
	require.Equal(t, 1, countOutbox(t, store, instanceID, "pending"))

	// Dispatch turns the tick into exactly one run, recorded truthfully as a schedule-driven run.
	started, err := d.Dispatch(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 1, started)
	require.Equal(t, 1, exec.distinctRuns())
	got := exec.firstStart()
	require.Equal(t, "revenue-agent", got.AgentName)
	require.Equal(t, act.TriggerSchedule, got.Trigger)
	require.Empty(t, got.TriggerRef, "a schedule run has no source resource to point at")
	require.Equal(t, "corre el resumen diario", got.Prompt)

	// Redelivering the tick (lease expiry) starts no second run: OAOO holds through ReserveDispatch and the executor.
	resetOutboxPending(t, store, instanceID)
	started, err = d.Dispatch(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 0, started)
	require.Equal(t, 1, exec.distinctRuns())
}

// TestDispatch_TriggerActorAttributesBecomeRunClaims verifies the run's governed identity: a trigger's run-as
// attributes are handed to the started run as SecurityClaims (mirroring an alert's query_for_attributes), and a
// trigger with no attributes yields nil claims so the run fails closed rather than running privileged.
func TestDispatch_TriggerActorAttributesBecomeRunClaims(t *testing.T) {
	store := requireStore(t)

	t.Run("attributes become claims", func(t *testing.T) {
		instanceID := uuid.NewString()
		tr := alertTrigger("on_fail", "incident", "revenue_drop", "entered_fail")
		tr.ActorService = true
		tr.ActorSubject = "act_ops"
		tr.ActorAttributes = map[string]any{"email": "act@kairosagentica.com", "admin": false}
		cat := &fakeCatalog{byInstance: map[string][]TriggerDef{instanceID: {tr}}}
		exec := &fakeExecutor{}
		d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

		enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredFail, time.Now())
		started, err := d.Dispatch(t.Context(), instanceID)
		require.NoError(t, err)
		require.Equal(t, 1, started)

		got := exec.firstStart()
		require.True(t, got.Actor.ServicePrincipal)
		require.NotNil(t, got.Actor.Claims, "run-as attributes must give the run governed claims")
		require.Equal(t, "act@kairosagentica.com", got.Actor.Claims.UserAttributes["email"])
		require.Equal(t, false, got.Actor.Claims.UserAttributes["admin"])
		require.False(t, got.Actor.Claims.SkipChecks, "governed claims must never skip security checks")
		for _, p := range []runtime.Permission{runtime.ReadObjects, runtime.ReadMetrics, runtime.UseAI} {
			require.True(t, got.Actor.Claims.Can(p),
				"governed run claims must carry %v: without it CheckAccess fails and the agent's metrics tools are silently dropped from its callable set", p)
		}
	})

	t.Run("no attributes fails closed with nil claims", func(t *testing.T) {
		instanceID := uuid.NewString()
		// alertTrigger sets no ActorAttributes, so the run must carry no claims.
		cat := &fakeCatalog{byInstance: map[string][]TriggerDef{
			instanceID: {alertTrigger("on_fail", "incident", "revenue_drop", "entered_fail")},
		}}
		exec := &fakeExecutor{}
		d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

		enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredFail, time.Now())
		_, err := d.Dispatch(t.Context(), instanceID)
		require.NoError(t, err)
		require.Nil(t, exec.firstStart().Actor.Claims, "no run-as attributes => nil claims => the worker fails closed")
	})
}

// TestDispatch_OneRunPerTransition verifies the core at-least-once guarantee: a transition event starts exactly one
// run, and re-delivering the same event (lease expiry, redelivery) starts no second run and records no second
// dispatch. The executor is idempotent on its run ID and the dispatch ledger is unique on (trigger, event).
func TestDispatch_OneRunPerTransition(t *testing.T) {
	store := requireStore(t)
	instanceID := uuid.NewString()
	cat := &fakeCatalog{byInstance: map[string][]TriggerDef{
		instanceID: {alertTrigger("triage_on_fail", "triage", "revenue_drop", "entered_fail")},
	}}
	exec := &fakeExecutor{}
	d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

	enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredFail, time.Now())

	started, err := d.Dispatch(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 1, started)
	require.Equal(t, 1, exec.distinctRuns())
	require.Equal(t, 1, countEvents(t, store, instanceID, "triage_on_fail", "dispatched"))
	require.Equal(t, 0, countOutbox(t, store, instanceID, "pending"))
	require.Equal(t, 1, countOutbox(t, store, instanceID, "done"))

	// Provenance: an alert-driven run records both its source kind and the alert's resource name, so the run detail
	// can link back to the alert that fired it.
	require.Equal(t, act.TriggerAlert, exec.firstStart().Trigger)
	require.Equal(t, "revenue_drop", exec.firstStart().TriggerRef, "the run must reference the source alert by name")

	// Re-deliver the same event: no new run, no new dispatch row.
	resetOutboxPending(t, store, instanceID)
	started, err = d.Dispatch(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 0, started)
	require.Equal(t, 1, exec.distinctRuns(), "a replayed event must not create a second run")
	require.Equal(t, 1, countEvents(t, store, instanceID, "triage_on_fail", "dispatched"))
}

// TestDispatch_RecoveryUsesFrozenAgent verifies the run identity is frozen at reservation. If a crash leaves a slot
// reserved but not yet dispatched and the trigger is then re-pointed to a different agent, completing the dispatch
// starts the run under the FROZEN agent, so the same (instance, trigger, event) never yields a second run — the bug
// where the mutable agent leaked into the run id.
func TestDispatch_RecoveryUsesFrozenAgent(t *testing.T) {
	store := requireStore(t)
	instanceID := uuid.NewString()

	// Simulate a crash after reservation but before the run started: reserve the slot under agent_a.
	ev := enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredFail, time.Now())
	reserved, err := store.ReserveDispatch(t.Context(), instanceID, "trig", ev.EventID, ev.EventType, "agent_a", 0)
	require.NoError(t, err)
	require.True(t, reserved.New)
	require.Equal(t, "agent_a", reserved.AgentName)

	// The trigger now points at a different agent.
	cat := &fakeCatalog{byInstance: map[string][]TriggerDef{
		instanceID: {alertTrigger("trig", "agent_b", "revenue_drop", "entered_fail")},
	}}
	exec := &fakeExecutor{}
	d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

	_, err = d.Dispatch(t.Context(), instanceID)
	require.NoError(t, err)

	// Exactly one run, started under the frozen agent_a rather than the trigger's current agent_b.
	require.Equal(t, 1, exec.startCount())
	require.Equal(t, "agent_a", exec.firstStart().AgentName)
	require.Equal(t, 1, exec.distinctRuns())
	require.Equal(t, 1, countEvents(t, store, instanceID, "trig", "dispatched"))
}

// TestDispatch_NoMatchIsDoneNotLost verifies an event no trigger subscribes to is marked done (not retried forever)
// and starts no run.
func TestDispatch_NoMatchIsDoneNotLost(t *testing.T) {
	store := requireStore(t)
	instanceID := uuid.NewString()
	cat := &fakeCatalog{byInstance: map[string][]TriggerDef{
		instanceID: {alertTrigger("triage_on_fail", "triage", "revenue_drop", "entered_fail")},
	}}
	exec := &fakeExecutor{}
	d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

	// recovered is not subscribed by the only trigger.
	enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertRecovered, time.Now())
	started, err := d.Dispatch(t.Context(), instanceID)
	require.NoError(t, err)
	require.Equal(t, 0, started)
	require.Equal(t, 0, exec.startCount())
	require.Equal(t, 1, countOutbox(t, store, instanceID, "done"))
}

// TestDispatch_WindowDedup verifies §9.2 window deduplication: two distinct events for the same trigger within the
// window collapse onto one run, while a zero window lets each event start its own run.
func TestDispatch_WindowDedup(t *testing.T) {
	t.Run("within window collapses to one run", func(t *testing.T) {
		store := requireStore(t)
		instanceID := uuid.NewString()
		tr := alertTrigger("triage", "triage", "revenue_drop", "entered_fail", "entered_error")
		tr.DedupWindow = 24 * time.Hour
		cat := &fakeCatalog{byInstance: map[string][]TriggerDef{instanceID: {tr}}}
		exec := &fakeExecutor{}
		d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

		base := time.Now()
		enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredFail, base)
		enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredError, base.Add(time.Minute))

		started, err := d.Dispatch(t.Context(), instanceID)
		require.NoError(t, err)
		require.Equal(t, 1, started, "second event within the window must not start a run")
		require.Equal(t, 1, exec.distinctRuns())
		require.Equal(t, 1, countEvents(t, store, instanceID, "triage", "dispatched"))
		require.Equal(t, 1, countEvents(t, store, instanceID, "triage", "deduplicated"))
	})

	t.Run("zero window starts a run per event", func(t *testing.T) {
		store := requireStore(t)
		instanceID := uuid.NewString()
		tr := alertTrigger("triage", "triage", "revenue_drop", "entered_fail", "entered_error")
		cat := &fakeCatalog{byInstance: map[string][]TriggerDef{instanceID: {tr}}}
		exec := &fakeExecutor{}
		d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

		base := time.Now()
		enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredFail, base)
		enqueueAlertEvent(t, store, instanceID, "revenue_drop", EventAlertEnteredError, base.Add(time.Minute))

		started, err := d.Dispatch(t.Context(), instanceID)
		require.NoError(t, err)
		require.Equal(t, 2, started)
		require.Equal(t, 2, exec.distinctRuns())
	})
}

// TestDispatch_MultiTenantScoping verifies a dispatch for one instance never touches another's events, even when
// both have an identically-named trigger and alert. The started run is scoped to the dispatched instance.
func TestDispatch_MultiTenantScoping(t *testing.T) {
	store := requireStore(t)
	instA := uuid.NewString()
	instB := uuid.NewString()
	cat := &fakeCatalog{byInstance: map[string][]TriggerDef{
		instA: {alertTrigger("triage", "triage", "revenue_drop", "entered_fail")},
		instB: {alertTrigger("triage", "triage", "revenue_drop", "entered_fail")},
	}}
	exec := &fakeExecutor{}
	d := NewDispatcher(store, cat, exec, DispatcherOptions{Logger: quietLogger()})

	enqueueAlertEvent(t, store, instA, "revenue_drop", EventAlertEnteredFail, time.Now())
	enqueueAlertEvent(t, store, instB, "revenue_drop", EventAlertEnteredFail, time.Now())

	started, err := d.Dispatch(t.Context(), instA)
	require.NoError(t, err)
	require.Equal(t, 1, started)
	require.Equal(t, 1, exec.startCount())
	require.Equal(t, instA, exec.firstStart().InstanceID, "the run must be scoped to the dispatched instance")

	// Instance B is untouched: its event is still pending and it created no dispatch rows.
	require.Equal(t, 1, countOutbox(t, store, instB, "pending"))
	require.Equal(t, 0, countEvents(t, store, instB, "triage", "dispatched"))
}

// TestObserver_EnqueueIdempotent verifies the observer path end to end: it derives the entered_fail event from a
// persisted alert execution and enqueues it, and re-observing the same execution (a source replay) enqueues nothing
// new. This is the at-least-once boundary from the source (§9.6).
func TestObserver_EnqueueIdempotent(t *testing.T) {
	store := requireStore(t)
	instanceID := uuid.NewString()
	obs := NewObserver(store, quietLogger())
	t.Cleanup(obs.Close)

	t0 := time.Now().Truncate(time.Second)
	res := &runtimev1.Resource{
		Meta: &runtimev1.ResourceMeta{Name: &runtimev1.ResourceName{Kind: runtime.ResourceKindAlert, Name: "revenue_drop"}},
		Resource: &runtimev1.Resource_Alert{Alert: &runtimev1.Alert{
			State: &runtimev1.AlertState{ExecutionHistory: []*runtimev1.AlertExecution{
				{Result: &runtimev1.AssertionResult{Status: runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL}, ExecutionTime: timestamppb.New(t0), FinishedOn: timestamppb.New(t0)},
				{Result: &runtimev1.AssertionResult{Status: runtimev1.AssertionStatus_ASSERTION_STATUS_PASS}, ExecutionTime: timestamppb.New(t0.Add(-time.Minute)), FinishedOn: timestamppb.New(t0.Add(-time.Minute))},
			}},
		}},
	}
	ev := runtime.ExecutionEvent{InstanceID: instanceID, Resource: res}

	// The observer publishes asynchronously, so wait for the background goroutine to drain the derived event to the
	// outbox. The inline query avoids failing an assertion from the polling goroutine.
	pendingCount := func() int {
		var n int
		if err := store.pool.QueryRow(t.Context(),
			fmt.Sprintf("SELECT count(*) FROM %s.agent_trigger_outbox WHERE instance_id=$1 AND state='pending'", store.schema),
			instanceID).Scan(&n); err != nil {
			return -1
		}
		return n
	}

	obs.OnExecution(t.Context(), ev)
	require.Eventually(t, func() bool { return pendingCount() == 1 }, 5*time.Second, 20*time.Millisecond)

	// Re-observe the same execution: the outbox stays at one (idempotent on instance + event_id), never growing to
	// two however the async publish is timed.
	obs.OnExecution(t.Context(), ev)
	require.Never(t, func() bool { return pendingCount() != 1 }, 500*time.Millisecond, 50*time.Millisecond)
}
