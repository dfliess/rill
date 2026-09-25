# runtime/act — Kairos Act durable executor

Kairos Act (issue kairosagentica/kairos-cloud#123) runs a declarative `Agent` resource as a **durable** workflow, so a
run survives process crashes, deploys and long human approval waits, and every external write it proposes is governed
and executed exactly once. Design: `docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md` (the decisions banner
prevails).

## Where it sits

- The **`Agent` resource**: proto `AgentSpec`/`AgentState`, `runtime/parser/parse_agent.go`, and a validating
  reconciler (`runtime/reconcilers/agent.go`) that never executes anything.
- **Dynamic execution**: `runtime/ai` runs an `AgentSnapshot` fail-closed (`dynamic_agent.go`, `agent_snapshot.go`).
- **This package** drives that execution durably. Reconciliation never executes an agent; execution enters **only**
  through `AgentExecutor`. `runtime/server/act_bootstrap.go` wires it in-process when an Act Postgres is configured.

## API

- `AgentExecutor` (`executor.go`) — `Start` / `Resume` / `ResumeInTx` / `Resumable` / `Cancel`. The reversibility seam:
  DBOS lives entirely behind it, so swapping the durable backend never touches `AgentSpec`, the APIs or the UI.
- `DBOSExecutor` (`dbos_executor.go`) — the DBOS implementation. One generic `RunAgentWorkflow` is registered under a
  stable name before `Launch`; every agent is data flowing through it, not a workflow per resource.
- `CatalogAgentProvider` (`provider.go`) — `ai.AgentDefinitionProvider` backed by the reconciled catalog. It reads
  `AgentState.ValidSpec`, so an agent that failed validation looks absent.
- `Runner` / `SessionRunner` (`runner.go`) — the model loop, one segment at a time, over `runtime/ai`.
- `ApprovalNotifier` (`executor.go`) — told once per **pause** of a run, with the actions awaiting a decision. Rill
  Cloud implements it in `runtime/server/act_approval_push.go` (web push, kairos-cloud#143).

## Workflow shape

```
step load_snapshot            → capture the immutable snapshot (the run is bound to it from here on)
step emit:run.running
loop per segment <seg>:
  step run_agent_<seg>        → run the model/tool loop until it answers or pauses on governed writes
  step persist_conversation   → first segment only: link the run to its AI session
  (no proposals)              → step emit:run.succeeded, done
  per proposed write <i>:
    step propose_action_<seg>_<i>, step gateway_propose_<seg>_<i>   → policy decision + ledger row
  step setup_approvals_<seg>  → every manual approval of the batch, at once
  step notify_approvals_<seg>
  per write <i>, in order:
    manual: Recv "approval:<tool_call_id>"                          → durable wait, as long as it takes
            step load_approver_<seg>_<i> (or load_denier_…), step gateway_record_approved_<seg>_<i>, step emit:run.resumed
    step gateway_execute_<seg>_<i>                                  → the write, exactly once
    step gateway_verify_<seg>_<i>                                   → when the tool is verifiable
  the results go back to the model in the next segment
```

**Changing this shape is dangerous.** The DBOS application version is pinned (`defaultApplicationVersion`), so a deploy
replays checkpoints recorded by older code. Every change to the step sequence goes behind a named `dbos.Patch`, read
**before the first step**, not at the point of the change: see the comment on `defaultApplicationVersion` and
`TestPatchGuardsShapeChangeAgainstParkedRuns`. Step options such as retries take no step slot and are safe.

**A cancelled workflow context is not a failure.** On a deploy DBOS cancels the context, keeps the run `PENDING` and
re-runs the interrupted step in the next process, so the workflow records nothing then (`fail`, kairos-cloud#170).

## Tool gateway and action ledger

Pieces (`policy.go`, `args.go`, `redact.go`, `ledger.go`/`pgledger.go`, `gateway.go`, `gateway_wiring.go`):

- **Policy engine** (`Evaluate`) — pure, fail-closed: unknown tool, unclassified tool, or corrupt or ambiguous policy
  never yield an implicit allow.
- **Action ledger** (`ActionLedger`, table `agent_actions`) — one row per proposed write, keyed by
  `(instance, run, tool_call_id)`, with the canonical args hash, idempotency key, policy decision and redacted
  args/result. Idempotent under at-least-once.
- **Tool Gateway** (`Gateway`) — the choke point every action passes through: allowlist, JSON-Schema validation,
  policy, server-side secrets, idempotency key, timeout/size caps, redaction. `Execute` reads the ledger first and
  claims the write atomically under a lease, so an approved write runs exactly once even if the step re-runs. A
  recovered step that meets another attempt's live lease waits it out (`ErrLeaseHeld`, retried by the step).
- **Kill switch** (`killswitch.go`) — the `KILL_ACT` project variable, checked at propose time and again right before
  the write.
- **Decisions** — the API records a decision and delivers it to the run in one transaction
  (`ResolveApprovalAndDeliver` + `ResumeInTx`, DBOS `WithSendTransaction`, kairos-cloud#194).

## Tests

Tests need PostgreSQL for the DBOS system tables. They default to `postgres://postgres:dbos@localhost:55432/dbos_spike`
(schema `act_dbos`; `docker start dbos-spike-pg`) and **skip** if it is unreachable, unless `ACT_TEST_REQUIRE_POSTGRES`
is set, which CI does: then a missing database fails. Override with `ACT_TEST_DBOS_URL` / `ACT_TEST_DBOS_SCHEMA`. The
LLM is always scripted; no real model is called. The crash tests (`crash_test.go`, `crash_gateway_test.go`) re-exec the
test binary as a worker subprocess (`ACT_TEST_WORKER=1`) to kill it mid-run and recover in a fresh process.
