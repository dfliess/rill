# runtime/act/trigger — Kairos Act trigger path (Fase 1, shadow mode)

Experimental. This package turns persisted alert/report executions into durable events and dispatches them to
agent runs (issue kairosagentica/kairos-cloud#123), in **shadow mode**: it creates the run but the run is the
supervised v1 executor. Design: `docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md`, §9 and §15.1 (the
decisions banner prevails).

## The seam

Reconcilers never contain Act logic. The runtime exposes a generic `ExecutionObserver`
(`runtime/execution_observer.go`); `alert.go` and `report.go` each call `Runtime.ObserveExecution` on exactly one
line, right after they persist an execution. With no observer registered the seam is inert (the upstream case).

## The three pieces

- **`Observer`** (`observer.go`) implements `runtime.ExecutionObserver`. It derives transition events from the
  execution history (`events.go`, §9.3/§9.4) and publishes each to the outbox. No LLM, no matching, never fails
  reconciliation.
- **`Store`** (`store.go`) owns two Postgres tables in its own schema:
  - `agent_trigger_outbox` — the at-least-once queue; unique `(instance_id, event_id)`, so a re-published event is
    a no-op.
  - `agent_trigger_events` — the dispatch ledger; unique `(instance_id, trigger_id, event_id)`, so a replayed event
    never creates two runs. Also the window-dedup lookup.
- **`Dispatcher`** (`dispatcher.go`) leases pending outbox events, matches them against the valid `AgentTrigger`
  resources (`TriggerCatalog`; `RuntimeTriggerCatalog` reads the catalog), deduplicates by window, and starts an
  `act.AgentExecutor` run per matched `(trigger, event)`.

## Event identity

`event_id = sha256(instance_id, resource_name, execution_time, event_type)`, length-prefixed and injective
(`ComputeEventID`). It depends only on values persisted with the execution, so it survives replay and does not
depend on when the dispatcher sees it. The executor idempotency key is `ComposeIdempotencyKey(trigger, event_id)`,
so two triggers on the same agent for the same event still get distinct runs, and a replay attaches to the same run.

## Wiring (integrator)

```go
store := trigger.NewStore(pool, "act_trigger")
_ = store.Migrate(ctx)
obs := trigger.NewObserver(store, logger) // starts a background publisher goroutine
defer obs.Close()                         // stops the goroutine at shutdown
rt.SetExecutionObserver(obs)              // once, before reconciliation starts
d := trigger.NewDispatcher(store, trigger.NewRuntimeTriggerCatalog(rt), executor, trigger.DispatcherOptions{})
// poll d.Dispatch(ctx, instanceID) per instance on an interval
```

`Observer` is non-blocking: `OnExecution` derives events in the reconciler's goroutine but hands them to a bounded
in-memory buffer that a background goroutine drains to Postgres, so a slow or unavailable Postgres never slows or
fails alert/report reconciliation. A full buffer drops the event with a log (shadow mode; §9.6's source re-publish
is the backstop). `Close` stops the publisher and drains what is still buffered.

**Not yet wired (known limit #5).** Nothing drives `Dispatch` in a deployed runtime yet: the observer registration,
dispatcher construction and the polling loop are exercised only by tests. Connecting them to a long-running process
(the `rill-agent-worker`) is a later deployment-increment task, so no automatic trigger currently fires in
production.

## Tests

The pure derivation and identity tests run anywhere. The `Store`/`Observer`/`Dispatcher` tests need PostgreSQL and
default to WS1's container (`postgres://postgres:dbos@localhost:55432/dbos_spike`; `docker start dbos-spike-pg`),
**skipping** if it is unreachable. Override with `ACT_TEST_DBOS_URL`. Each test uses a throwaway schema. The
executor is a fake; the dispatcher's match/dedup needs no DBOS.
