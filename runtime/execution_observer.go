package runtime

import (
	"context"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
)

// ExecutionObserver is notified once a reconciler has persisted an execution of a resource that produces
// executions (currently Alert and Report). It is a generic, upstreamable seam: the reconcilers announce that an
// execution happened and contain no further logic; a registered observer derives higher-level transition events,
// deduplicates them and dispatches work. In Kairos this is where Act hangs its trigger machinery
// (runtime/act/trigger); with no observer registered the seam is inert.
//
// OnExecution runs synchronously in the reconciler's goroutine, so it must be quick and must never fail
// reconciliation: it returns no error and is expected to handle and log its own failures. A slow or unavailable
// downstream is the observer's problem to bound, not the reconciler's.
type ExecutionObserver interface {
	OnExecution(ctx context.Context, event ExecutionEvent)
}

// ExecutionEvent is the raw signal handed to an ExecutionObserver. It carries the resource exactly as persisted,
// including the execution just appended to its history: for both Alert and Report, State.ExecutionHistory[0] is
// that execution and State.ExecutionHistory[1], if present, is the one before it. The observer derives transition
// events (§9.3/§9.4) from that history; the reconciler holds none of that logic.
type ExecutionEvent struct {
	InstanceID string
	Resource   *runtimev1.Resource
}

// executionObserverHolder boxes the observer interface so it can live in an atomic.Pointer (which needs a concrete
// pointer type). It lets SetExecutionObserver publish an observer race-free relative to reconcilers calling
// ObserveExecution.
type executionObserverHolder struct {
	obs ExecutionObserver
}

// SetExecutionObserver registers (or, with nil, clears) the observer notified after alert/report executions. It is
// intended to be called once during application wiring, before reconciliation starts.
func (r *Runtime) SetExecutionObserver(obs ExecutionObserver) {
	r.executionObserver.Store(&executionObserverHolder{obs: obs})
}

// ObserveExecution forwards a just-persisted execution to the registered observer, if any. It is safe to call with
// no observer registered (the common, upstream case), where it is a no-op.
func (r *Runtime) ObserveExecution(ctx context.Context, instanceID string, self *runtimev1.Resource) {
	h := r.executionObserver.Load()
	if h == nil || h.obs == nil {
		return
	}
	h.obs.OnExecution(ctx, ExecutionEvent{InstanceID: instanceID, Resource: self})
}
