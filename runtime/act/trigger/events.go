// Package trigger is the Kairos Act trigger path (issue kairosagentica/kairos-cloud#123): it turns persisted
// alert/report executions into durable events and dispatches them to agent runs, in shadow mode.
//
// The flow has three moving parts, all behind the generic runtime.ExecutionObserver seam so the reconcilers stay
// free of Act logic (§9 of docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md):
//
//   - Observer derives transition events (§9.3/§9.4) from an execution's state history and publishes each to a
//     durable outbox. It never calls an LLM and never blocks reconciliation.
//   - The outbox (agent_trigger_outbox) is the at-least-once queue; a re-published event is idempotent on
//     (instance_id, event_id).
//   - Dispatcher polls the outbox, matches events against the valid AgentTrigger resources in the catalog,
//     deduplicates by window, and in shadow mode starts an AgentExecutor run with an idempotency key derived from
//     the event. The dispatch ledger (agent_trigger_events) is unique on (instance_id, trigger_id, event_id), so a
//     replayed event never creates two runs.
package trigger

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
)

// schemaVersion is the version of the event envelope written to the outbox (§9.2). Bump it when the envelope shape
// changes so a consumer can tell old rows from new ones.
const schemaVersion = 1

// Resource kinds as they appear in an event's short namespace (and in an AgentTrigger's source.kind).
const (
	resourceKindAlert    = "alert"
	resourceKindReport   = "report"
	resourceKindSchedule = "schedule"
)

// Fully-qualified event types (§9.3/§9.4). The short name after the "<kind>." prefix is what an AgentTrigger
// subscribes to in source.events. renotify_due and evaluated are intentionally not emitted yet in shadow mode
// (see the deriveAlert* comments).
const (
	EventAlertEnteredFail  = "alert.entered_fail"
	EventAlertRecovered    = "alert.recovered"
	EventAlertEnteredError = "alert.entered_error"
	EventAlertRenotifyDue  = "alert.renotify_due"
	EventReportCompleted   = "report.completed"
	EventReportFailed      = "report.failed"
	// EventScheduleTick is emitted by the scheduler's ticker when a schedule trigger's cron boundary passes. Unlike
	// alert/report events it has no upstream execution: the trigger is its own source (§9.5).
	EventScheduleTick = "schedule.tick"
)

// Event is a derived transition event: a fact about one execution of one resource. Its identity, EventID, is a
// pure function of (instance, resource, execution_time, event_type), so it is stable across replays and does not
// depend on when the dispatcher happens to see it.
type Event struct {
	InstanceID    string
	EventID       string
	EventType     string
	ResourceKind  string
	ResourceName  string
	ExecutionTime time.Time
	OccurredOn    time.Time
	Status        string
}

// Envelope is the versioned, serializable form of an Event stored in the outbox payload (§9.2). It carries
// references and bounded evidence, never secrets or full datasets: the worker re-queries Rill for live context.
type Envelope struct {
	EventID       string   `json:"event_id"`
	EventType     string   `json:"event_type"`
	InstanceID    string   `json:"instance_id"`
	ResourceKind  string   `json:"resource_kind"`
	ResourceName  string   `json:"resource_name"`
	ExecutionTime string   `json:"execution_time"`
	OccurredOn    string   `json:"occurred_on"`
	Status        string   `json:"status"`
	ContextRefs   []string `json:"context_refs"`
	SchemaVersion int      `json:"schema_version"`
}

// Envelope returns the serializable envelope for the event.
func (e Event) Envelope() Envelope {
	return Envelope{
		EventID:       e.EventID,
		EventType:     e.EventType,
		InstanceID:    e.InstanceID,
		ResourceKind:  e.ResourceKind,
		ResourceName:  e.ResourceName,
		ExecutionTime: e.ExecutionTime.UTC().Format(time.RFC3339Nano),
		OccurredOn:    e.OccurredOn.UTC().Format(time.RFC3339Nano),
		Status:        e.Status,
		ContextRefs:   []string{},
		SchemaVersion: schemaVersion,
	}
}

// ComputeEventID derives an event's stable identity from (instance_id, resource_name, execution_time, event_type).
// It length-prefixes the strings and encodes the time as Unix nanoseconds so the encoding is injective, and it
// depends only on values persisted with the execution so a replay of the same execution yields the same ID.
// event_type is namespaced by resource kind ("alert.entered_fail"), so resource kind need not be a separate input.
func ComputeEventID(instanceID, resourceName string, executionTime time.Time, eventType string) string {
	h := sha256.New()
	writeField(h, instanceID)
	writeField(h, resourceName)
	_ = binary.Write(h, binary.BigEndian, executionTime.UTC().UnixNano())
	writeField(h, eventType)
	return hex.EncodeToString(h.Sum(nil))
}

// writeField length-prefixes a string into the hash so field boundaries are unambiguous.
func writeField(h interface{ Write([]byte) (int, error) }, s string) {
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(s)))
	_, _ = h.Write(lenBuf[:])
	_, _ = h.Write([]byte(s))
}

// DeriveEvents turns a persisted execution into zero or more transition events. It reads only the resource's state
// history, so it is a pure function that is trivial to unit-test; all the trigger semantics of §9.3/§9.4 live here,
// not in the reconcilers.
func DeriveEvents(ev runtime.ExecutionEvent) []Event {
	res := ev.Resource
	if res == nil || res.Meta == nil || res.Meta.Name == nil {
		return nil
	}
	name := res.Meta.Name.Name
	switch {
	case res.GetAlert() != nil:
		return deriveAlertEvents(ev.InstanceID, name, res.GetAlert())
	case res.GetReport() != nil:
		return deriveReportEvents(ev.InstanceID, name, res.GetReport())
	default:
		return nil
	}
}

// deriveAlertEvents derives an alert's transition event from its execution history (newest first, [0] is the
// execution just persisted). It emits at most one event: the transition into the current status, or a renotify on a
// continued fail/error. It never emits on a boring continuation, which is what keeps a fail that stays fail from
// re-triggering an investigation on every tick (§9.3).
func deriveAlertEvents(instanceID, name string, a *runtimev1.Alert) []Event {
	hist := a.GetState().GetExecutionHistory()
	if len(hist) == 0 || hist[0].GetResult() == nil {
		return nil
	}
	cur := hist[0]
	var prev *runtimev1.AlertExecution
	if len(hist) > 1 {
		prev = hist[1]
	}
	prevStatus := runtimev1.AssertionStatus_ASSERTION_STATUS_UNSPECIFIED
	if prev != nil && prev.GetResult() != nil {
		prevStatus = prev.Result.Status
	}

	var eventType, status string
	switch cur.Result.Status {
	case runtimev1.AssertionStatus_ASSERTION_STATUS_PASS:
		status = "pass"
		// recovered: a change into pass from a non-pass state; never the first pass without history (§9.3).
		if prev != nil && prevStatus != runtimev1.AssertionStatus_ASSERTION_STATUS_PASS {
			eventType = EventAlertRecovered
		}
	case runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL:
		status = "fail"
		if prevStatus != runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL {
			eventType = EventAlertEnteredFail
		} else if alertRenotifyDue(a, hist) {
			eventType = EventAlertRenotifyDue
		}
	case runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR:
		status = "error"
		if prevStatus != runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR {
			eventType = EventAlertEnteredError
		} else if alertRenotifyDue(a, hist) {
			eventType = EventAlertRenotifyDue
		}
	default:
		return nil
	}
	if eventType == "" {
		return nil
	}
	return []Event{newEvent(instanceID, resourceKindAlert, name, eventType, status, alertExecTime(cur), alertOccurredOn(cur))}
}

// TODO(act, phase 1): known limit #10 — renotify_due is derived purely from the execution history window, which is
// bounded (Rill keeps only the last ~25 executions). Once the history is long enough that the transition execution
// which opened the current same-status run falls off the end, the walk-back below can no longer find t0, so a
// continued fail/error may re-emit renotify_due on every tick instead of once per window. The fix is to persist the
// origin of the current transition (when the run entered its current status) rather than recomputing it from a
// truncatable history, so the renotify boundary math is stable regardless of how many executions are retained.

// alertRenotifyDue reports whether a continued fail/error should re-trigger, reusing the alert's own renotify clock
// (§9.3: renotify_due "reutiliza renotify_after"). The clock is measured purely from the execution history — the
// time since the transition that opened the current same-status run — so it is decoupled from whether an email or
// Slack notification was actually delivered, and from the on_fail/on_error notification flags.
func alertRenotifyDue(a *runtimev1.Alert, hist []*runtimev1.AlertExecution) bool {
	if a.GetSpec() == nil || !a.Spec.Renotify || len(hist) < 2 {
		return false
	}
	if a.Spec.RenotifyAfterSeconds == 0 {
		// No suppression period: renotify on every continuation tick.
		return true
	}
	cur := hist[0]
	curStatus := cur.Result.Status
	// Walk back to the transition execution that opened the current same-status run.
	t0 := alertExecTime(cur)
	for i := 1; i < len(hist); i++ {
		if hist[i].GetResult() == nil || hist[i].Result.Status != curStatus {
			break
		}
		t0 = alertExecTime(hist[i])
	}
	window := time.Duration(a.Spec.RenotifyAfterSeconds) * time.Second
	tc := alertExecTime(cur)
	tp := alertExecTime(hist[1])
	// Fire once per window: only when the current tick crosses a renotify boundary the previous tick had not.
	return tc.Sub(t0)/window > tp.Sub(t0)/window
}

// deriveReportEvents derives a report's completion event: reports have no transitions, so each finished execution
// emits exactly completed or failed depending on whether it recorded an error (§9.4).
func deriveReportEvents(instanceID, name string, rep *runtimev1.Report) []Event {
	hist := rep.GetState().GetExecutionHistory()
	if len(hist) == 0 {
		return nil
	}
	cur := hist[0]
	eventType, status := EventReportCompleted, "completed"
	if cur.ErrorMessage != "" {
		eventType, status = EventReportFailed, "failed"
	}
	return []Event{newEvent(instanceID, resourceKindReport, name, eventType, status, reportExecTime(cur), reportOccurredOn(cur))}
}

// newScheduleEvent assembles the synthetic tick a schedule trigger emits for one cron boundary. The trigger is its
// own source, so the event's resource name is the trigger's own name and its identity time is the boundary: the
// EventID is then stable per (instance, trigger, boundary), which is exactly what makes re-enqueuing the same
// boundary (across overlapping scan windows or a second poller) idempotent, and what keeps at most one run per tick.
func newScheduleEvent(instanceID, triggerName string, boundary, occurredOn time.Time) Event {
	return newEvent(instanceID, resourceKindSchedule, triggerName, EventScheduleTick, "scheduled", boundary, occurredOn)
}

// newEvent assembles an Event and computes its stable EventID.
func newEvent(instanceID, kind, name, eventType, status string, execTime, occurredOn time.Time) Event {
	return Event{
		InstanceID:    instanceID,
		EventID:       ComputeEventID(instanceID, name, execTime, eventType),
		EventType:     eventType,
		ResourceKind:  kind,
		ResourceName:  name,
		ExecutionTime: execTime,
		OccurredOn:    occurredOn,
		Status:        status,
	}
}

// alertExecTime returns the stable identity time of an alert execution: its execution_time, or its finished_on as a
// fallback. Both are persisted with the execution, so either survives a replay (§9.2).
func alertExecTime(e *runtimev1.AlertExecution) time.Time {
	if e.GetExecutionTime() != nil {
		return e.ExecutionTime.AsTime()
	}
	if e.GetFinishedOn() != nil {
		return e.FinishedOn.AsTime()
	}
	return time.Time{}
}

func alertOccurredOn(e *runtimev1.AlertExecution) time.Time {
	if e.GetFinishedOn() != nil {
		return e.FinishedOn.AsTime()
	}
	return alertExecTime(e)
}

// reportExecTime returns the stable identity time of a report execution: its report_time, or finished_on/started_on
// as fallbacks for executions that failed before a report time was assigned.
func reportExecTime(e *runtimev1.ReportExecution) time.Time {
	if e.GetReportTime() != nil {
		return e.ReportTime.AsTime()
	}
	if e.GetFinishedOn() != nil {
		return e.FinishedOn.AsTime()
	}
	if e.GetStartedOn() != nil {
		return e.StartedOn.AsTime()
	}
	return time.Time{}
}

func reportOccurredOn(e *runtimev1.ReportExecution) time.Time {
	if e.GetFinishedOn() != nil {
		return e.FinishedOn.AsTime()
	}
	return reportExecTime(e)
}
