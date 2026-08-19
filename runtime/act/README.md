# runtime/act — Kairos Act durable executor (Fase 0 spike)

Experimental. This package is the Fase 0 spike for Kairos Act (issue kairosagentica/kairos-cloud#123). It runs a
declarative `Agent` resource as a **durable** workflow, so a run survives process crashes, deploys and long human
approval waits. Design: `docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md` (the decisions banner prevails).

## Where it sits

Three pieces already existed before this package:

- The **`Agent` resource**: proto `AgentSpec`/`AgentState`, `runtime/parser/parse_agent.go`, and a validating
  reconciler (`runtime/reconcilers/agent.go`) that never executes anything.
- **Dynamic execution**: `runtime/ai` runs an `AgentSnapshot` fail-closed (`dynamic_agent.go`, `agent_snapshot.go`).
- A **DBOS validation spike** (WS1) that proved DBOS Transact Go is a fit.

This package is the missing third leg: it drives that execution durably. Reconciliation still never executes an
agent; execution enters **only** through `AgentExecutor`.

## API

- `AgentExecutor` (`executor.go`) — `Start` / `Resume` / `Cancel`. The reversibility seam: DBOS lives entirely
  behind it, so swapping the durable backend never touches `AgentSpec`, the APIs or the UI.
- `DBOSExecutor` (`dbos_executor.go`) — the DBOS implementation. One generic `RunAgentWorkflow` is registered under
  a stable name before `Launch`; every agent is data flowing through it, not a workflow per resource.
- `CatalogAgentProvider` (`provider.go`) — `ai.AgentDefinitionProvider` backed by the reconciled catalog. It reads
  `AgentState.ValidSpec`, so an agent that failed validation looks absent.
- `Runner` / `SessionRunner` (`runner.go`) — the side-effecting seam. The workflow wraps each `Runner` call in a
  DBOS step; `SessionRunner` implements it over `runtime/ai` and an `ActionSink`.
- `ApprovalNotifier` (`executor.go`) — optional seam told once per **pause** of a run, with the batch's actions
  awaiting a decision. Who to tell needs the project's membership and the agent's approve policy, which this
  package cannot reach, so Rill Cloud implements it in `runtime/server/act_approval_push.go` (web push,
  kairos-cloud#143). Emitted right after the durable step that creates the approvals, never inside it: a replay
  must not skip the notification, and re-notifying is the harmless direction.

## Workflow shape (spike)

```
step load_snapshot   → capture the immutable snapshot (run is bound to it from here on)
step run_agent       → run the whole runtime/ai model/tool loop as ONE step
Recv "approval"      → suspend durably, indefinitely; Cancel is how a run nobody decides is ended
step apply_action    → simulated external write, only if approved, exactly once on the happy path
```

Step-granularity trade-off: `run_agent` is one step for the spike. Fase 1 must split it into one step per model
turn and per tool call (§10.3), pairing each write with an action ledger (§12), so a crash mid-loop resumes at the
last completed turn and external writes are never re-driven.

## Fase 2: tool gateway and action ledger

When a `Gateway` and a `Proposer` are wired into `Config`, the run's action phase runs the real flow instead of the
simulated write above:

```
step propose_action    → derive the proposed tool call (Proposer)
step gateway_propose    → validate + policy + record ledger (proposed → policy decision)
  deny → run fails, no effect        approval_required → suspend        allow → (Fase 3) execute directly
Recv "approval"         → the existing durable approval wait
step gateway_record_approved → bind the decision to the exact args hash
step apply_action       → gateway.Execute: idempotent external write via ActionExecutor
step gateway_verify     → best-effort confirmation (verifiable tools)
```

Pieces (`policy.go`, `args.go`, `redact.go`, `ledger.go`/`pgledger.go`, `gateway.go`):

- **Policy engine** (`Evaluate`, `policy.go`) — pure, fail-closed: `(tool, allowlist, class, approval config, global/
  tenant policy, kill switch) → allow / approval_required / deny`. Unknown tool, unclassified tool, corrupt or
  ambiguous policy never yield an implicit allow; global/tenant layers may only harden.
- **Action ledger** (`ActionLedger`, table `agent_actions`) — one row per proposed write, keyed by
  `(instance, run, tool_call_id)`, with the canonical args hash, idempotency key, policy decision, redacted args/result
  and the §11.2 state machine. Idempotent under at-least-once: a replayed proposal or transition is a no-op.
- **Tool Gateway** (`Gateway`) — the choke point every action tool call passes through: allowlist (registry
  membership), JSON-Schema arg validation, classification, policy, server-side secret resolution, idempotency key +
  trace, timeout/size caps, and redaction of everything persisted. `Gateway.Execute` reads the ledger first and
  short-circuits on any state that means the write already happened or must not be retried, so an approved write runs
  exactly once even if the step re-runs on recovery.
- **`ActionExecutor`** (`gateway.go`) — the seam the outbound MCP client implements: `Execute` + `Verify`. Tests use a
  fake. The gateway resolves policy, secrets, idempotency and redaction around it.
- **Kill switch** (`KillSwitch`, §17.2) — a per-platform/tenant/agent check the gateway consults at propose time and
  again immediately before the write; engaged → deny.

## Tests

Tests need PostgreSQL for the DBOS system tables. They default to WS1's container
(`postgres://postgres:dbos@localhost:55432/dbos_spike`, schema `act_dbos`; `docker start dbos-spike-pg`) and
**skip** if it is unreachable. Override with `ACT_TEST_DBOS_URL` / `ACT_TEST_DBOS_SCHEMA`. The LLM is always the
scripted mock from `runtime/ai`; no real model is called. The crash test re-execs the test binary as a worker
subprocess (`ACT_TEST_WORKER=1`) to SIGKILL a worker mid-run and recover it in a fresh process.
