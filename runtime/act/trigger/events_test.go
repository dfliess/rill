package trigger

import (
	"testing"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// alertResource builds an Alert resource with the given (newest-first) execution history and optional spec.
func alertResource(name string, spec *runtimev1.AlertSpec, history ...*runtimev1.AlertExecution) *runtimev1.Resource {
	return &runtimev1.Resource{
		Meta: &runtimev1.ResourceMeta{Name: &runtimev1.ResourceName{Kind: runtime.ResourceKindAlert, Name: name}},
		Resource: &runtimev1.Resource_Alert{Alert: &runtimev1.Alert{
			Spec:  spec,
			State: &runtimev1.AlertState{ExecutionHistory: history},
		}},
	}
}

// alertExec builds one alert execution with a status and an execution time.
func alertExec(status runtimev1.AssertionStatus, at time.Time) *runtimev1.AlertExecution {
	return &runtimev1.AlertExecution{
		Result:        &runtimev1.AssertionResult{Status: status},
		ExecutionTime: timestamppb.New(at),
		FinishedOn:    timestamppb.New(at),
	}
}

func deriveTypes(ev runtime.ExecutionEvent) []string {
	events := DeriveEvents(ev)
	types := make([]string, 0, len(events))
	for _, e := range events {
		types = append(types, e.EventType)
	}
	return types
}

const (
	statusPass  = runtimev1.AssertionStatus_ASSERTION_STATUS_PASS
	statusFail  = runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL
	statusError = runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR
)

// TestDeriveAlertEvents_PassFailFailPass is the core §9.3 assertion: a pass->fail->fail->pass sequence emits
// entered_fail exactly once (on the transition, not on the continued fail) and recovered exactly once. This is what
// stops a fail that stays fail from re-triggering an investigation on every tick.
func TestDeriveAlertEvents_PassFailFailPass(t *testing.T) {
	instanceID := "inst"
	t0 := time.Date(2026, 7, 14, 8, 0, 0, 0, time.UTC)
	tick := func(d time.Duration) time.Time { return t0.Add(d) }

	pass1 := alertExec(statusPass, tick(0))
	fail2 := alertExec(statusFail, tick(5*time.Minute))
	fail3 := alertExec(statusFail, tick(10*time.Minute))
	pass4 := alertExec(statusPass, tick(15*time.Minute))

	ev := func(history ...*runtimev1.AlertExecution) runtime.ExecutionEvent {
		return runtime.ExecutionEvent{InstanceID: instanceID, Resource: alertResource("revenue_drop", nil, history...)}
	}

	// tick 1: first pass, no history -> no recovered (§9.3: not emitted for the first pass).
	require.Empty(t, deriveTypes(ev(pass1)))
	// tick 2: pass -> fail transition -> entered_fail.
	require.Equal(t, []string{EventAlertEnteredFail}, deriveTypes(ev(fail2, pass1)))
	// tick 3: fail -> fail continuation, renotify off -> nothing.
	require.Empty(t, deriveTypes(ev(fail3, fail2, pass1)))
	// tick 4: fail -> pass transition -> recovered.
	require.Equal(t, []string{EventAlertRecovered}, deriveTypes(ev(pass4, fail3, fail2, pass1)))
}

// TestDeriveAlertEvents_EnteredError verifies error transitions: the first error and an error from another state
// emit entered_error, while a continued error without renotify emits nothing.
func TestDeriveAlertEvents_EnteredError(t *testing.T) {
	t0 := time.Date(2026, 7, 14, 8, 0, 0, 0, time.UTC)
	err1 := alertExec(statusError, t0)
	err2 := alertExec(statusError, t0.Add(5*time.Minute))
	pass0 := alertExec(statusPass, t0.Add(-5*time.Minute))

	ev := func(history ...*runtimev1.AlertExecution) runtime.ExecutionEvent {
		return runtime.ExecutionEvent{InstanceID: "inst", Resource: alertResource("a", nil, history...)}
	}

	// First execution is an error -> entered_error.
	require.Equal(t, []string{EventAlertEnteredError}, deriveTypes(ev(err1)))
	// Error from pass -> entered_error.
	require.Equal(t, []string{EventAlertEnteredError}, deriveTypes(ev(err1, pass0)))
	// Continued error, renotify off -> nothing.
	require.Empty(t, deriveTypes(ev(err2, err1, pass0)))
}

// TestDeriveAlertEvents_RenotifyDue verifies the history-derived renotify clock: a continued fail re-triggers once
// per renotify_after measured from the transition, independent of notification delivery.
func TestDeriveAlertEvents_RenotifyDue(t *testing.T) {
	spec := &runtimev1.AlertSpec{Renotify: true, RenotifyAfterSeconds: 3600} // 1h
	t0 := time.Date(2026, 7, 14, 8, 0, 0, 0, time.UTC)
	at := func(m int) time.Time { return t0.Add(time.Duration(m) * time.Minute) }

	fail0 := alertExec(statusFail, at(0))     // transition (entered_fail)
	fail30 := alertExec(statusFail, at(30))   // 30m into the run
	fail65 := alertExec(statusFail, at(65))   // crosses the 1h boundary
	fail90 := alertExec(statusFail, at(90))   // still in the same hour bucket as 65m
	fail125 := alertExec(statusFail, at(125)) // crosses the 2h boundary

	ev := func(history ...*runtimev1.AlertExecution) runtime.ExecutionEvent {
		return runtime.ExecutionEvent{InstanceID: "inst", Resource: alertResource("a", spec, history...)}
	}

	require.Equal(t, []string{EventAlertEnteredFail}, deriveTypes(ev(fail0)))
	require.Empty(t, deriveTypes(ev(fail30, fail0))) // 30m < 1h -> no renotify
	require.Equal(t, []string{EventAlertRenotifyDue}, deriveTypes(ev(fail65, fail30, fail0)))
	require.Empty(t, deriveTypes(ev(fail90, fail65, fail30, fail0))) // same hour bucket -> no renotify
	require.Equal(t, []string{EventAlertRenotifyDue}, deriveTypes(ev(fail125, fail90, fail65, fail30, fail0)))
}

// TestDeriveReportEvents verifies a report emits completed or failed per its recorded error (§9.4).
func TestDeriveReportEvents(t *testing.T) {
	t0 := time.Date(2026, 7, 14, 8, 0, 0, 0, time.UTC)
	ok := &runtimev1.ReportExecution{ReportTime: timestamppb.New(t0), FinishedOn: timestamppb.New(t0)}
	bad := &runtimev1.ReportExecution{ErrorMessage: "delivery failed", ReportTime: timestamppb.New(t0), FinishedOn: timestamppb.New(t0)}

	reportRes := func(exec *runtimev1.ReportExecution) *runtimev1.Resource {
		return &runtimev1.Resource{
			Meta:     &runtimev1.ResourceMeta{Name: &runtimev1.ResourceName{Kind: runtime.ResourceKindReport, Name: "weekly"}},
			Resource: &runtimev1.Resource_Report{Report: &runtimev1.Report{State: &runtimev1.ReportState{ExecutionHistory: []*runtimev1.ReportExecution{exec}}}},
		}
	}

	require.Equal(t, []string{EventReportCompleted}, deriveTypes(runtime.ExecutionEvent{InstanceID: "inst", Resource: reportRes(ok)}))
	require.Equal(t, []string{EventReportFailed}, deriveTypes(runtime.ExecutionEvent{InstanceID: "inst", Resource: reportRes(bad)}))
}

// TestComputeEventID_StableAndInjective guards the event identity: it is stable for the same inputs (so a replay
// deduplicates) and distinguishes inputs that a naive join would collide, including the (resource, event_type)
// boundary.
func TestComputeEventID_StableAndInjective(t *testing.T) {
	tm := time.Date(2026, 7, 14, 8, 0, 0, 0, time.UTC)

	require.Equal(t, ComputeEventID("inst", "a", tm, "alert.entered_fail"), ComputeEventID("inst", "a", tm, "alert.entered_fail"))
	// Different event type on the same execution -> different id.
	require.NotEqual(t, ComputeEventID("inst", "a", tm, "alert.entered_fail"), ComputeEventID("inst", "a", tm, "alert.recovered"))
	// Different execution time -> different id.
	require.NotEqual(t, ComputeEventID("inst", "a", tm, "alert.entered_fail"), ComputeEventID("inst", "a", tm.Add(time.Second), "alert.entered_fail"))
	// Boundary between resource name and event type must not collide.
	require.NotEqual(t, ComputeEventID("inst", "a", tm, "b.c"), ComputeEventID("inst", "ab", tm, ".c"))
	// Different instance -> different id (multi-tenant).
	require.NotEqual(t, ComputeEventID("inst1", "a", tm, "alert.entered_fail"), ComputeEventID("inst2", "a", tm, "alert.entered_fail"))
}

// TestComposeIdempotencyKey_Injective guards the executor idempotency key: two triggers must not collide onto one
// run for the same event, even when their names and the event id would join ambiguously.
func TestComposeIdempotencyKey_Injective(t *testing.T) {
	require.NotEqual(t, ComposeIdempotencyKey("a", "b/c"), ComposeIdempotencyKey("a/b", "c"))
	require.Equal(t, ComposeIdempotencyKey("triage", "e1"), ComposeIdempotencyKey("triage", "e1"))
}

// TestShortEventName strips the kind namespace a trigger does not spell out.
func TestShortEventName(t *testing.T) {
	require.Equal(t, "entered_fail", shortEventName("alert.entered_fail"))
	require.Equal(t, "completed", shortEventName("report.completed"))
	require.Equal(t, "bare", shortEventName("bare"))
}
