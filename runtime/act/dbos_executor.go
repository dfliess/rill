package act

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
)

const (
	appName      = "kairos-act"
	workflowName = "RunAgentWorkflow" // stable name: required for cross-process enqueue and recovery
	// runQueue is the queue every run is enqueued on, so that creating the run and launching its workflow commit in one
	// transaction (ADR-0032). It has no concurrency limit, and must never get one: a run parked on a human approval
	// counts as PENDING, so a limit of N would stop dispatching once N runs were waiting (ADR-0022, decision 3).
	runQueue      = "act-runs"
	approvalTopic = "approval"

	stepLoadSnapshot        = "load_snapshot"
	stepRunAgent            = "run_agent"
	stepPersistConversation = "persist_conversation"
	stepLoadApprover        = "load_approver"
	stepLoadDenier          = "load_denier"

	// terminalWriteTimeout bounds the writes that close out a run (recording its terminal status, sweeping its
	// approvals). They run on a context that outlives the workflow's precisely because a shutdown cancels that one,
	// so they need a deadline of their own: an unbounded write here would hold the process open on a dead database.
	terminalWriteTimeout = 10 * time.Second

	// payloadKeyDecidedBy carries the subject that decided an approval on the run.resumed / run.rejected events. The
	// same subject is written to the action ledger's decided_by; the event is what makes it readable through the API,
	// which does not expose the ledger.
	payloadKeyDecidedBy = "decided_by"

	// defaultApplicationVersion pins the DBOS application version to a constant, overriding DBOS's default of
	// versioning the application by the binary's SHA-256. That default is not unstable across restarts — it is
	// stable per binary, by design: only a worker running the exact code that recorded a workflow's checkpoints
	// will recover it, which is what makes replay safe. We override it anyway because a run parked on a human
	// approval must survive a deploy (the premise of the whole HITL design), and this deployment replaces the
	// binary in place: under per-binary versioning no worker with the old hash would ever come back for the run.
	//
	// The price of the constant is that a new binary replays checkpoints written by DIFFERENT code, so every
	// change to the workflow's step sequence MUST be gated behind a named dbos.Patch (armed via EnablePatching
	// below). Deploying an unguarded new step killed 19 parked runs on 2026-09-02: the step ran in the slot the
	// parked Recv had not recorded yet, and the run died with "UnexpectedStep: DBOS.sleep was recorded when
	// DBOS.recv was expected". And the Patch call must sit where every in-flight run already has a checkpoint in
	// the next slot — in practice before the workflow's first step — NOT at the point of the change: a parked
	// Recv records nothing in its own slot (only its timeout sleep, one slot later), so a Patch placed just
	// before the wait reads that empty slot as "new workflow" and steers the parked run into the new branch,
	// reproducing the same incident. TestPatchGuardsShapeChangeAgainstParkedRuns pins both halves of this rule.
	//
	// Retiring a patch follows Kairos ADR-0033, which is DBOS's own procedure: dbos.DeprecatePatch once no unfinished
	// workflow predates the deploy that added the patch, then removal once none predates the deprecation. Both are
	// checked in the database, never by calendar: a parked run can wait for months.
	defaultApplicationVersion = "act-spike-v1"

	// approvalWaitForever is the duration a run waits for a human decision. dbos.Recv has no "wait forever" argument and
	// a zero timeout fires immediately, so we pass this effectively-infinite duration (~100 years, comfortably within
	// int64 nanoseconds). The wait is durable and cheap (a parked workflow blocked on Recv), and Cancel is the human's
	// way to abort a run nobody will decide.
	//
	// Do not read this number as a promise about calendar time. What actually bounds the wait is maxRecoveryAttempts
	// below: every process start spends one attempt of a parked run's budget, so the wait ends at whichever comes
	// first — and it is never the hundred years. At DBOS's default of 100 attempts and the ~0.8 restarts/day measured
	// in production on 2026-09-02, "forever" was about four months.
	approvalWaitForever = 100 * 365 * 24 * time.Hour

	// maxRecoveryAttempts is the recovery budget of the agent workflow, raised far above the DBOS default of 100
	// because that default silently contradicts approvalWaitForever above.
	//
	// To DBOS, a run parked on a human decision is indistinguishable from one that keeps failing: it sits PENDING,
	// and every process start recovers it and spends one attempt (`IncrementAttempts` is set on both the recovery
	// path and the dequeue path, so re-enqueueing does not avoid it either). At 100 the wait is not "forever", it is
	// however long 100 restarts take. Measured in production on 2026-09-02: runs parked since 2026-08-17 had already
	// spent 13 attempts, ~0.8/day, which exhausts the default within months — the run would then be dead-lettered as
	// MAX_RECOVERY_ATTEMPTS_EXCEEDED and, per the Status switch below, treated as terminal. A client must never lose
	// a run because nobody signed it yet, least of all against a deadline the signer cannot see.
	//
	// Raised rather than disabled: 0 switches dead-lettering off entirely, and the budget still has one job — a
	// genuine poison pill, a run that crashes the process on every recovery, must stop instead of looping forever.
	// A crash loop spends this budget in about a day; an approval, at the rate measured above, in about 34 years.
	maxRecoveryAttempts = 10000
)

// Config configures a DBOSExecutor and the DBOS worker it embeds.
type Config struct {
	// DatabaseURL is the Postgres DSN for the DBOS system tables. For Act this is the Platform Postgres instance,
	// logical database "act"; DatabaseSchema keeps the DBOS tables in their own schema.
	DatabaseURL    string
	DatabaseSchema string
	// ApplicationVersion pins the worker's version for blue/green and recovery. Defaults to defaultApplicationVersion.
	ApplicationVersion string
	// Runner performs the side-effecting work of each run. Required.
	Runner Runner
	// Store is the queryable run store the executor writes lifecycle state and events to (§15). It is the product's
	// state plane, separate from the DBOS system tables: the API reads runs from it without touching the workflow
	// engine. Required: it also holds the one run of each idempotency key that makes Start idempotent (ADR-0034).
	Store RunStore
	// Gateway governs every write a run proposes: policy, the action ledger, secrets, the idempotent external call and
	// its verification (§16.2, §12). Proposer surfaces the write the model loop captured. Both are required. The
	// Gateway's Ledger should share the Postgres of Store (its agent_actions table lives in the product schema).
	Gateway  *Gateway
	Proposer Proposer
	// Approvals, when set, is notified every time a run pauses on human approvals, so whoever may decide is told out
	// of band (kairos-cloud#143). Optional and best-effort: a deployment without it simply notifies nobody, and a
	// notification that fails never touches the run.
	Approvals ApprovalNotifier
	Logger    *slog.Logger
}

// DBOSExecutor is the DBOS-backed AgentExecutor. It both embeds the worker (it registers the workflow and runs
// recovered/started runs) and drives runs from the outside (Start/Resume/Cancel), which is the single-process
// dev topology (§20.1). In production the worker and the runtime that drives it are separate processes sharing the
// same Postgres; the split does not change this code, only which side holds the dbos.Client.
type DBOSExecutor struct {
	ctx       dbos.Context
	client    dbos.Client
	runner    Runner
	store     RunStore
	gateway   *Gateway
	proposer  Proposer
	approvals ApprovalNotifier
	version   string
	logger    *slog.Logger
}

var _ AgentExecutor = (*DBOSExecutor)(nil)

// NewDBOSExecutor builds the worker: it registers the single generic RunAgentWorkflow under a stable name (before
// Launch, as DBOS requires), launches the runtime — which recovers this executor's in-flight runs — and opens a
// dbos.Client for delivering approvals and cancellations.
func NewDBOSExecutor(ctx context.Context, cfg Config) (*DBOSExecutor, error) {
	if cfg.Runner == nil {
		return nil, errors.New("act: Config.Runner is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("act: Config.DatabaseURL is required")
	}
	if cfg.Gateway == nil || cfg.Proposer == nil {
		return nil, errors.New("act: Config.Gateway and Config.Proposer are required")
	}
	// The store is what makes Start idempotent: it holds the one run of each idempotency key (ADR-0034).
	if cfg.Store == nil {
		return nil, errors.New("act: Config.Store is required")
	}

	version := cfg.ApplicationVersion
	if version == "" {
		version = defaultApplicationVersion
	}

	// Tag the executor's context as Act traffic. The API handlers already tag theirs, but a run does not execute in the
	// request that started it: RunWorkflow below hands the work to DBOS on this context, and recovery after a restart
	// has no request to inherit from at all. Without this, every token and tool call an agent spends is recorded with
	// an empty source, so the meter cannot tell agent spend from chat spend.
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)

	// Detach the engine from the caller's cancellation. Close is what tears DBOS down, and dbos.Shutdown marks that
	// cancellation with a cause of its own so a workflow unwinding under it is left PENDING for the next process to
	// recover instead of being durably cancelled. That check reads the CAUSE, so it only holds when Shutdown is the
	// thing doing the cancelling. The runtime hands us a context cancelled by SIGTERM (graceful.WithCancelOnTerminate),
	// which reaches the workflows first and carries a plain context.Canceled: the check does not fire, and every run
	// parked on a human approval is durably cancelled on the way out, which is the bug this whole path exists to avoid
	// (kairos-cloud#135). Detaching costs nothing, because the deferred Close already runs after the server stops.
	ctx = context.WithoutCancel(ctx)

	dctx, err := dbos.NewContext(ctx, dbos.Config{
		AppName:            appName,
		DatabaseURL:        cfg.DatabaseURL,
		DatabaseSchema:     cfg.DatabaseSchema,
		ApplicationVersion: version,
		// The pinned ApplicationVersion means a deploy replays checkpoints recorded by older code, so any future
		// change to the workflow's step sequence must hide behind a named dbos.Patch — see the rule (and the
		// 2026-09-02 incident that wrote it) at defaultApplicationVersion. This only arms the mechanism; Patch is
		// independent of the version string, and the existing steps stay unwrapped on purpose: a Patch consumes a
		// step slot of its own, so wrapping an already-deployed step would shift the numbering and be a shape
		// change in itself.
		EnablePatching: true,
		Logger:         cfg.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("act: init DBOS: %w", err)
	}

	e := &DBOSExecutor{
		ctx:       dctx,
		runner:    cfg.Runner,
		store:     cfg.Store,
		gateway:   cfg.Gateway,
		proposer:  cfg.Proposer,
		approvals: cfg.Approvals,
		version:   version,
		logger:    cfg.Logger,
	}

	// Register the ONE generic workflow. It is a bound method so the workflow body reaches the runner (and, in
	// production, the runtime) without those crossing the serialized workflow boundary. Recovery re-invokes this
	// same registered closure by name, so a recovered run uses the recovering process's runner.
	dbos.RegisterWorkflow(dctx, e.runAgentWorkflow,
		dbos.WithWorkflowName(workflowName),
		dbos.WithMaxRecoveryAttempts(maxRecoveryAttempts),
	)

	// The queue's configuration is persisted and reloaded by DBOS. Always overwrite it, so a limit set on it by hand
	// is undone on the next start rather than silently blocking parked runs.
	if _, err := dbos.RegisterQueue(dctx, runQueue, dbos.WithQueueOnConflict(dbos.QueueConflictAlwaysUpdate)); err != nil {
		_ = dbos.Shutdown(dctx, 5*time.Second) // best effort: we are already failing to start
		return nil, fmt.Errorf("act: register queue %q: %w", runQueue, err)
	}

	if err := dbos.Launch(dctx); err != nil {
		_ = dbos.Shutdown(dctx, 5*time.Second) // best effort: we are already failing to start
		return nil, fmt.Errorf("act: launch DBOS: %w", err)
	}

	client, err := dbos.NewClient(ctx, dbos.ClientConfig{
		DatabaseURL:    cfg.DatabaseURL,
		DatabaseSchema: cfg.DatabaseSchema,
		Logger:         cfg.Logger,
	})
	if err != nil {
		_ = dbos.Shutdown(dctx, 5*time.Second) // best effort: we are already failing to start
		return nil, fmt.Errorf("act: init DBOS client: %w", err)
	}
	e.client = client

	return e, nil
}

// Close drains the worker and closes the client. Callers should allow in-flight steps to finish within timeout.
//
// A shutdown error is logged, not returned: Close runs on the way out of the process, so there is no caller left to
// act on it, and the client is closed even if the worker's drain reported a problem.
func (e *DBOSExecutor) Close(timeout time.Duration) {
	if e.client != nil {
		if err := dbos.Shutdown(e.client, timeout); err != nil && e.logger != nil {
			e.logger.Warn("act: closing the DBOS client failed", "err", err)
		}
	}
	if err := dbos.Shutdown(e.ctx, timeout); err != nil && e.logger != nil {
		e.logger.Warn("act: draining the DBOS worker failed", "err", err)
	}
}

// Start launches a run and returns its ID, which Resume/Cancel/Result address the run by. A repeated Start with the
// same (instance, agent, key) attaches to the run that key already has (OAOO) instead of starting a second one, and
// returns that run's ID.
func (e *DBOSExecutor) Start(ctx context.Context, in AgentRunInput) (string, error) {
	if in.InstanceID == "" || in.AgentName == "" {
		return "", errors.New("act: AgentRunInput.InstanceID and AgentName are required")
	}
	if in.IdempotencyKey == "" {
		return "", errors.New("act: AgentRunInput.IdempotencyKey is required")
	}
	// An opaque ID (ADR-0034): the store, not the ID, deduplicates by (instance, agent, key), so two instances
	// (tenants) or two agents that reuse a key still get separate runs. A different payload under the same key is not
	// detected: the second Start attaches to the first run whatever it carries.
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("act: new run ID: %w", err)
	}
	runID := id.String()

	// Normalize and validate the trigger provenance against the known sources (§9.5). Validating here, not just at
	// the API boundary, keeps a run's recorded provenance trustworthy even for internal Go callers: an unknown
	// trigger string can never be persisted.
	trigger := in.Trigger
	if trigger == "" {
		trigger = TriggerManual
	}
	if !validTriggers[trigger] {
		return "", fmt.Errorf("act: unknown trigger %q", trigger)
	}

	// Record the run and enqueue its workflow in ONE transaction, so the run is queryable the instant Start returns
	// (§14.2) and exists if and only if its workflow does (ADR-0032). If the key already has a run, the store hands
	// enqueue that run's ID instead: a no-op for a run with a workflow, the launch for one left without it.
	enqueue := func(ctx context.Context, tx pgx.Tx, runID string) error {
		_, err := dbos.Enqueue[AgentRunResult](e.client, runQueue, workflowName, in,
			dbos.WithEnqueueWorkflowID(runID),
			dbos.WithEnqueueApplicationVersion(e.version),
			dbos.WithEnqueueApplicationName(appName),
			dbos.WithEnqueueTransaction(tx))
		return err
	}
	runID, err = e.store.CreateRunAndEnqueue(ctx, NewRun{
		RunID:          runID,
		InstanceID:     in.InstanceID,
		AgentName:      in.AgentName,
		Trigger:        trigger,
		TriggerRef:     in.TriggerRef,
		IdempotencyKey: in.IdempotencyKey,
		ConversationID: in.ConversationID,
		Actor:          in.Actor,
	}, enqueue)
	if err != nil {
		return "", fmt.Errorf("act: start run %q: %w", runID, err)
	}
	return runID, nil
}

// Resume delivers an approval decision to a suspended run. The message is durable; it is safe to send it before the
// run reaches its wait, in which case the run consumes it when it gets there.
//
// TODO(act, #11): the durable resume message carries only the decision, not the approver's identity. The API resume
// path records the real decider first: ResolveApproval sets the approval row's decided_by before the resume is sent,
// and the workflow reloads it (see runActionViaGateway's load_approver step) so the ledger attributes the approval
// correctly. A DIRECT call to Resume that bypasses the API leaves decided_by unset, so the workflow falls back to the
// run initiator as the approver — an audit gap, not an authorization one (the write is still gated on approval). Any
// non-API resume path added later must set the approval's decided_by first, or thread the approver through this call.
func (e *DBOSExecutor) Resume(_ context.Context, runID string, decision ApprovalDecision, toolCallID string) error {
	return e.send(runID, decision, toolCallID)
}

// ResumeInTx delivers the decision inside tx, so it reaches the run only if tx commits (dbos.WithSendTransaction,
// which needs tx on the same database as the DBOS system tables: BootstrapAct gives both one DSN).
func (e *DBOSExecutor) ResumeInTx(_ context.Context, tx pgx.Tx, runID string, decision ApprovalDecision, toolCallID string) error {
	return e.send(runID, decision, toolCallID, dbos.WithSendTransaction(tx))
}

func (e *DBOSExecutor) send(runID string, decision ApprovalDecision, toolCallID string, opts ...dbos.SendOption) error {
	// Every governed action waits on its own topic. An approval without a tool call comes from the retired simulated
	// path, whose run no longer exists to read it: refuse rather than file a message nobody will consume.
	if toolCallID == "" {
		return fmt.Errorf("act: resume run %q: the approval names no tool call", runID)
	}
	if err := dbos.Send(e.client, runID, string(decision), approvalTopic+":"+toolCallID, opts...); err != nil {
		return fmt.Errorf("act: resume run %q: %w", runID, err)
	}
	return nil
}

// Resumable reports whether the run's workflow can still consume a decision. A workflow in a terminal state cannot:
// its function has returned, so nothing is left to receive the message. Send would still succeed there: it inserts a
// notification row whose only constraint is that the workflow EXISTS, which a terminal one does — so this is the check
// that keeps an approval from being recorded against a run that will never act on it.
//
// Reading the status before the decision is recorded is not a race in the direction that matters: terminal states are
// final, so a "no" cannot become a "yes". A "yes" that dies immediately afterwards leaves the caller exactly where it
// was before this check existed.
func (e *DBOSExecutor) Resumable(_ context.Context, runID string) (bool, error) {
	h, err := dbos.RetrieveWorkflow[AgentRunResult](e.ctx, runID)
	if err != nil {
		return false, fmt.Errorf("act: retrieve run %q: %w", runID, err)
	}
	st, err := h.GetStatus()
	if err != nil {
		return false, fmt.Errorf("act: status of run %q: %w", runID, err)
	}
	switch st.Status {
	case dbos.WorkflowStatusSuccess, dbos.WorkflowStatusError, dbos.WorkflowStatusCancelled, dbos.WorkflowStatusMaxRecoveryAttemptsExceeded:
		return false, nil
	default:
		return true, nil
	}
}

// Cancel stops a run. After the workflow is cancelled it records the cancelled transition directly (not as a step,
// since the cancelled workflow will not run more steps) so the API reflects the terminal state. The store write is
// best-effort: a cancelled workflow is already stopped, so a failed state write must not turn Cancel into an error.
//
// The transition is scoped by instance, so a caller can only cancel a run in the instance it addresses. Cancelling a
// run that already finished leaves its status alone: the store never moves a run out of a terminal state.
func (e *DBOSExecutor) Cancel(ctx context.Context, instanceID, runID string) error {
	if err := dbos.CancelWorkflow(e.client, runID); err != nil {
		return fmt.Errorf("act: cancel run %q: %w", runID, err)
	}
	// The workflow is already cancelled; the store write below only reflects that, so a failure to record it must not
	// fail Cancel. The transition and the approval withdrawal go together in CloseRun rather than as two independent
	// writes: withdrawing the approval of a run whose cancelled transition did NOT land leaves it reading
	// waiting_approval with nothing left to decide, which is worse than the inbox briefly offering a decision on a run
	// that will not resume — that one an operator can retry, the other one needs a database (kairos-cloud#135).
	if _, err := e.store.CloseRun(ctx, RunTransition{
		InstanceID: instanceID,
		RunID:      runID,
		Status:     RunStatusCancelled,
		EventType:  EventTypeCancelled,
	}); err != nil && e.logger != nil {
		e.logger.Warn("act: cancel: closing run failed", "run", runID, "err", err)
	}
	return nil
}

// Result waits for a run to finish and returns its checkpointed outcome. Only tests use it; the product surfaces run
// state through the Agent API (§14), not through the executor.
func (e *DBOSExecutor) Result(runID string) (AgentRunResult, error) {
	h, err := dbos.RetrieveWorkflow[AgentRunResult](e.ctx, runID)
	if err != nil {
		return AgentRunResult{}, fmt.Errorf("act: retrieve run %q: %w", runID, err)
	}
	return h.GetResult()
}

// runAgentWorkflow is the single generic durable workflow that runs any agent. Its body is deterministic: every
// non-deterministic or side-effecting operation is a checkpointed step, and the body only branches on values that
// are already checkpointed. On recovery, completed steps replay from their checkpoints instead of re-executing.
// After binding the run to its snapshot it hands over to the segmented loop (runSegmentedLoop).
func (e *DBOSExecutor) runAgentWorkflow(ctx dbos.Context, in AgentRunInput) (AgentRunResult, error) {
	runID, err := dbos.GetWorkflowID(ctx)
	if err != nil {
		return AgentRunResult{}, err
	}
	res := AgentRunResult{RunID: runID, AgentName: in.AgentName, Status: RunStatusRunning}

	// Step 1: capture the immutable snapshot. After this checkpoint the run is bound to this definition; editing
	// the agent resource afterwards does not change what this run does. It runs while the run is still queued so the
	// running transition below can carry the snapshot's spec hash.
	snapshot, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (*ai.AgentSnapshot, error) {
		return e.runner.LoadSnapshot(stepCtx, in.InstanceID, in.AgentName)
	}, dbos.WithStepName(stepLoadSnapshot))
	if err != nil {
		return e.fail(ctx, in, runID, res, "load snapshot", err)
	}

	// The worker has the run's definition: move it from queued to running and bind it to the agent version it
	// executed by recording the snapshot's spec hash on the same transition. The emit is a durable step, so the
	// state plane advances exactly once per edge and a reader sees a consistent status/event/spec_hash triple.
	if err := e.emitRunning(ctx, in, runID, snapshot.SpecHash); err != nil {
		return res, err
	}

	return e.runSegmentedLoop(ctx, in, runID, snapshot, res)
}

// emitRunning records the queued->running transition and binds the run to the snapshot's spec hash in one durable
// step. It is separate from emit because it is the only transition that carries a spec hash: it fires right after
// the snapshot is captured, so the run row records exactly which agent version it executed (§8.3, §12).
func (e *DBOSExecutor) emitRunning(ctx dbos.Context, in AgentRunInput, runID, specHash string) error {
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		return true, e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusRunning,
			EventType:  EventTypeRunning,
			SpecHash:   specHash,
		})
	}, dbos.WithStepName("emit:"+EventTypeRunning))
	if err != nil {
		return fmt.Errorf("act: emit %s: %w", EventTypeRunning, err)
	}
	return nil
}

// persistConversation links the run to the AI session it executed in, as a durable step. It is a no-op when no session
// was opened (SessionID empty), so it never overwrites an existing link with a blank. The store UPDATE is idempotent, so a replay re-writes the same checkpointed session ID.
func (e *DBOSExecutor) persistConversation(ctx dbos.Context, in AgentRunInput, runID, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		return true, e.store.SetRunConversationID(stepCtx, in.InstanceID, runID, sessionID)
	}, dbos.WithStepName(stepPersistConversation))
	if err != nil {
		return fmt.Errorf("act: persist conversation: %w", err)
	}
	return nil
}

// emit records one lifecycle transition as a durable step.
// The step is idempotent in the store (one event per edge), so a crash between the side effect and the checkpoint replays it harmlessly.
func (e *DBOSExecutor) emit(ctx dbos.Context, in AgentRunInput, runID string, status RunStatus, eventType string) error {
	return e.emitWithPayload(ctx, in, runID, status, eventType, "", nil)
}

// emitWithPayload is emit with structured detail attached to the event, deduplicated by dedupeKey (empty keys the
// edge on its event type alone — once per run). The payload must already be redacted: it is stored verbatim and read
// by the run timeline, which is user-visible (§17.2). A nil payload behaves exactly like emit, so the step name and
// shape are unchanged for the transitions that carry no detail.
func (e *DBOSExecutor) emitWithPayload(ctx dbos.Context, in AgentRunInput, runID string, status RunStatus, eventType, dedupeKey string, payload map[string]any) error {
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		tr := RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     status,
			EventType:  eventType,
			DedupeKey:  dedupeKey,
			Payload:    payload,
		}
		// A terminal status closes the run: its still-pending approvals are withdrawn in the same transaction, so
		// "a terminal run has no pending approvals" holds by construction rather than by every caller remembering to
		// sweep. Forgetting that sweep, or doing it when the run did NOT go terminal, is what stranded runs in
		// waiting_approval with nothing left to decide (kairos-cloud#135).
		if status.IsTerminal() {
			_, closeErr := e.store.CloseRun(stepCtx, tr)
			return true, closeErr
		}
		return true, e.store.RecordTransition(stepCtx, tr)
	}, dbos.WithStepName("emit:"+eventType))
	if err != nil {
		return fmt.Errorf("act: emit %s: %w", eventType, err)
	}
	return nil
}

// deciderPayload is the event detail that attributes an approval decision on the run timeline. It is the only way a
// reader learns who unblocked (or stopped) a run from the events alone: the ledger records the same subject on the
// action, but the ledger is not exposed through the API.
func deciderPayload(subject string) map[string]any {
	if subject == "" {
		return nil
	}
	return map[string]any{payloadKeyDecidedBy: subject}
}

// loadDecider resolves who actually decided an approval, as a durable step. The run initiator (in.Actor.Subject, e.g.
// Alice) is not necessarily who decided (e.g. Bob, an admin), and the durable resume message carries only the decision,
// not the decider; the approval row's decided_by was set by the API's claim (ResolveApproval) before the resume was
// delivered, so it is authoritative here (§11.2, §18.3). Falls back to the initiator when a non-API resume path left
// decided_by unset — see the Resume TODO above.
func (e *DBOSExecutor) loadDecider(ctx dbos.Context, in AgentRunInput, approvalID, stepName string) (string, error) {
	loaded, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (string, error) {
		ap, gErr := e.store.GetApproval(stepCtx, in.InstanceID, approvalID)
		if gErr != nil {
			return "", gErr
		}
		return ap.DecidedBy, nil
	}, dbos.WithStepName(stepName))
	if err != nil {
		return "", err
	}
	if loaded == "" {
		return in.Actor.Subject, nil
	}
	return loaded, nil
}

// fail records the run as failed with a sanitized error and returns the original error unchanged. The failed emit is
// best-effort: an error recording the failure must not mask the failure that caused it. "what" is a short label for
// the failed stage; err.Error() is stored as the sanitized message (the runner is responsible for not leaking
// secrets into it — §17.2).
//
// A cancelled workflow context is not a failure, and fail records nothing for it. It means the process is draining
// (a deploy) or the run was cancelled, and in both cases DBOS owns the outcome: it discards whatever the workflow
// returns, leaves a drained run PENDING for the next process, and re-runs the interrupted step, which it never
// checkpointed. Writing failed here is what killed runs caught by a deploy mid-step: the recovered step then met a
// terminal run and could not ask for its signature again (kairos-cloud#170). DBOS's guidance is the same: once the
// context is cancelled, return early, since any work done after it only computes a result DBOS throws away.
// An operator's Cancel records its own terminal state, so nothing is lost by returning here.
func (e *DBOSExecutor) fail(ctx dbos.Context, in AgentRunInput, runID string, res AgentRunResult, what string, err error) (AgentRunResult, error) {
	if ctx.Err() != nil {
		return res, fmt.Errorf("act: %s: %w", what, err)
	}
	res.Status = RunStatusFailed
	e.closeRun(ctx, in, runID, RunTransition{
		InstanceID: in.InstanceID,
		RunID:      runID,
		Status:     RunStatusFailed,
		EventType:  EventTypeFailed,
		Error:      err.Error(),
	}, what)
	return res, fmt.Errorf("act: %s: %w", what, err)
}

// closeRun durably marks a run terminal and withdraws its pending approvals as one atomic store write, then withdraws
// its ledger actions. It is best-effort: an error closing the run must not mask the failure that caused it.
//
// The store write is retried on a context that outlives the workflow's, because a shutdown arriving during this write
// cancels that one, and a run that genuinely failed must be recorded as failed whether or not the process stays up. A
// shutdown that arrived before the failure never gets here: fail leaves that run to DBOS. Doing the two writes separately is what stranded runs reading
// waiting_approval with every approval cancelled and no decider (kairos-cloud#135), so CloseRun binds them: either
// the run is terminal AND its approvals are withdrawn, or neither happened and the run stays decidable.
func (e *DBOSExecutor) closeRun(ctx dbos.Context, in AgentRunInput, runID string, tr RunTransition, what string) {
	// The step returns bool, not the cancelled count: DBOS checkpoints a step's output by type, so a run interrupted
	// mid-deploy replays its recorded output through the new binary. Widening the type here would fail to deserialize
	// checkpoints written by the old one.
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		_, closeErr := e.store.CloseRun(stepCtx, tr)
		return true, closeErr
	}, dbos.WithStepName("emit:"+tr.EventType))
	if err != nil {
		closeCtx, cancel := context.WithTimeout(context.Background(), terminalWriteTimeout)
		_, err = e.store.CloseRun(closeCtx, tr)
		cancel()
	}
	if err != nil && e.logger != nil {
		e.logger.Error("act: closing run failed; it stays decidable with its approvals pending",
			"run", runID, "stage", what, "status", tr.Status, "err", err)
	}
	e.withdrawPendingActions(in.InstanceID, runID)
}

// withdrawPendingActions closes out the ledger actions of a run that reached a terminal state. It is the ledger
// counterpart of the approval sweep CloseRun performs, kept separate because the ledger may live behind the gateway
// rather than in the run store, so it cannot join that transaction. A ledger action left open is an audit loose end,
// not a run a human sees stuck, so best-effort is enough here. It runs on a context that outlives the workflow's,
// which a shutdown may already have cancelled.
func (e *DBOSExecutor) withdrawPendingActions(instanceID, runID string) {
	ledger := e.ledger()
	if ledger == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), terminalWriteTimeout)
	defer cancel()
	if _, err := ledger.WithdrawPendingActions(ctx, instanceID, runID); err != nil && e.logger != nil {
		e.logger.Warn("act: withdrawing pending actions failed", "run", runID, "err", err)
	}
}

// closeLedgerAction transitions a single ledger action to a terminal status with an optional decider. Best-effort: a
// failure is logged but never fails the run, since the approval row already records the decision.
func (e *DBOSExecutor) closeLedgerAction(ctx context.Context, instanceID, runID, toolCallID string, status ActionStatus, decidedBy string) {
	ledger := e.ledger()
	if ledger == nil {
		return
	}
	if _, err := ledger.TransitionAction(ctx, ActionTransition{
		InstanceID: instanceID, RunID: runID, ToolCallID: toolCallID,
		Status: status, DecidedBy: decidedBy,
	}); err != nil && e.logger != nil {
		e.logger.Warn("act: close ledger action failed", "run", runID, "toolCallID", toolCallID, "status", status, "err", err)
	}
}

// ledger returns the gateway's action ledger (in production the run store, which serves both).
func (e *DBOSExecutor) ledger() ActionLedger {
	return e.gateway.Ledger
}

// errInterruptedAwaitingApproval marks an approval wait cut short by a cancelled context, so callers unwinding the
// stack can tell a shutdown apart from a real failure. It must stay distinguishable all the way up: the difference
// decides whether the run's pending approvals are swept, and sweeping them on a shutdown is what stranded runs in
// waiting_approval with nothing left to decide (kairos-cloud#135).
var errInterruptedAwaitingApproval = errors.New("act: approval wait interrupted")

// interruptedAwaitingApproval wraps a cancelled-context error from the approval wait. A cancelled context there means
// the process is shutting down (a deploy) or the run was cancelled, NOT that the run failed: so unlike fail it records
// no failed transition, leaving the store row truthfully waiting_approval. That keeps the run cancellable to clean it
// up, and a Cancel that raced the shutdown still resolves cleanly.
//
// Returning this error is safe for the run's durability, which was not always true. Until DBOS v1.1 a returning error
// was taken as a claim that the run was unrecoverable, so a deploy wrote off every approval nobody had signed yet;
// that is how 27 of the 29 terminal workflows in production died (kairos-cloud#135, upstream #423). Since v1.1 the
// shutdown cancellation carries its own cause and a run unwinding under it skips the outcome write entirely, so the
// row keeps the PENDING it already had and the next process recovers it back into this same Recv, which reads the
// notifications table before waiting: a decision taken while the platform was down is consumed at once.
func interruptedAwaitingApproval(err error) error {
	return fmt.Errorf("%w (context cancelled): %w", errInterruptedAwaitingApproval, err)
}

// actionProposalStep is the checkpointed result of the propose step: the proposed tool call and whether the run
// proposed one at all. It is a named type because a DBOS step returns a single value, and a bare tuple cannot be
// checkpointed.
type actionProposalStep struct {
	Proposal ToolProposal
	OK       bool
}

// runSegmentStepResult is the checkpointed outcome of one run_segment step: the segment's final response, the ID of the
// AI session it ran in, and the writes it paused on (if any). Bundling them in one named struct means DBOS checkpoints
// them together, so a replay yields the same session ID and proposals rather than recomputing them against a freshly
// re-driven loop.
type runSegmentStepResult struct {
	Response  string
	SessionID string
	// Proposals is the list of write actions the segment captured (empty for a segment that produced a final answer).
	// Checkpointed alongside the response, so a replay after a crash yields the same proposals.
	Proposals []*ai.ProposedAction
}

// runSegmentedLoop drives a run as a segmented durable loop (b2). Each iteration runs one segment of the model/tool
// loop as a durable step; when a segment pauses on a governed write, the action flows through the Tool Gateway (policy,
// action ledger, human approval or auto-approve, idempotent execution) and its result is injected into the next
// segment, so the model sees the outcome of the action it proposed and either proposes another or composes a final
// answer. The body branches only on checkpointed values, so a crash resumes at the last completed step, and — because
// the gateway's Execute is idempotent on the ledger — an approved write is never driven twice under at-least-once
// delivery (§12). A crash mid-segment re-drives that segment's model turns (the accepted granularity: replaying a
// segment re-flushes its trace but never re-executes a committed external write).
func (e *DBOSExecutor) runSegmentedLoop(ctx dbos.Context, in AgentRunInput, runID string, snapshot *ai.AgentSnapshot, res AgentRunResult) (AgentRunResult, error) {
	var sessionID string
	var resume []*ai.InjectedResult
	for seg := 0; ; seg++ {
		out, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (runSegmentStepResult, error) {
			r, runErr := e.runner.RunSegment(stepCtx, RunSegmentInput{
				InstanceID: in.InstanceID, Claims: in.Actor.Claims, Snapshot: snapshot,
				Prompt: in.Prompt, SessionID: sessionID, Resume: resume,
			})
			return runSegmentStepResult{Response: r.Response, SessionID: r.SessionID, Proposals: r.Proposed}, runErr
		}, dbos.WithStepName(fmt.Sprintf("%s_%d", stepRunAgent, seg)))
		if err != nil {
			return e.fail(ctx, in, runID, res, "run segment", err)
		}
		res.Response = out.Response

		if seg == 0 {
			sessionID = out.SessionID
			if err := e.persistConversation(ctx, in, runID, sessionID); err != nil {
				return e.fail(ctx, in, runID, res, "persist conversation", err)
			}
		}

		if len(out.Proposals) == 0 {
			res.Status = RunStatusSucceeded
			if err := e.emit(ctx, in, runID, RunStatusSucceeded, EventTypeSucceeded); err != nil {
				return res, err
			}
			return res, nil
		}

		injected, terminal, err := e.governProposedActions(ctx, in, runID, snapshot, &res, out.Proposals, seg)
		if terminal || err != nil {
			return res, err
		}
		resume = injected
	}
}

// governedAction holds the gateway authorization and connector resolution for one proposed action, produced during the
// batch-propose phase so the per-action processing loop can apply policy, wait for approval, and execute in order.
type governedAction struct {
	proposal     ToolProposal
	auth         Authorization
	capturedConn ai.MCPConnector
	hasConn      bool
	autoApprove  map[string]bool
}

// canonicalArgsForApproval returns the preimage to persist with an action's approval. Normally it is the one the
// gateway captured while hashing, which is the only source that cannot have drifted from the hash.
//
// The fallback exists for one situation: a run whose propose step was checkpointed by a build from before the
// Authorization carried the preimage. On recovery DBOS replays that checkpoint, the field deserializes empty, and
// refusing there would kill a legitimate in-flight run during a deploy. So rebuild it from the proposal's own
// arguments and accept it ONLY if it hashes back to the authorization's hash. That check is what keeps this from
// being a hole: a rebuild that matches IS the preimage, whatever produced it, and one that does not match is
// discarded so the store's own binding check refuses the approval.
func canonicalArgsForApproval(a *governedAction) string {
	if a.auth.CanonicalArgs != "" {
		return a.auth.CanonicalArgs
	}
	rebuilt, err := CanonicalizeArgs(a.proposal.Args)
	if err != nil || hashBytes(rebuilt) != a.auth.ArgsHash {
		return ""
	}
	return string(rebuilt)
}

// governProposedActions governs N proposed writes for segment seg. It derives each proposal, evaluates policy, creates
// all manual approvals at once (visible in the inbox), then processes each action in proposal order: auto-approved
// actions execute immediately; manual actions wait for a per-tool-call DBOS topic. A rejection injects an error result
// for that call (the model sees it on resume and adapts) without ending the run. Truly terminal outcomes (indeterminate,
// workflow errors) stop the run; everything else produces an InjectedResult per action for the next segment.
func (e *DBOSExecutor) governProposedActions(ctx dbos.Context, in AgentRunInput, runID string, snapshot *ai.AgentSnapshot, res *AgentRunResult, captured []*ai.ProposedAction, seg int) ([]*ai.InjectedResult, bool, error) {
	n := len(captured)
	actions := make([]governedAction, n)

	// Phase 1: propose and authorize every action. Each runs as its own durable step so the gateway records the ledger
	// entry and the policy decision per action, and a crash replays only the unfinished proposals.
	for i, cap := range captured {
		proposed, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (actionProposalStep, error) {
			p, ok, perr := e.proposer.Propose(stepCtx, ProposeInput{
				InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Actor: in.Actor,
				Snapshot: snapshot, Response: res.Response, Captured: cap,
			})
			return actionProposalStep{Proposal: p, OK: ok}, perr
		}, dbos.WithStepName(fmt.Sprintf("propose_action_%d_%d", seg, i)))
		if err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "propose action", err)
			return nil, true, err
		}
		if !proposed.OK {
			res.Status = RunStatusSucceeded
			if err := e.emit(ctx, in, runID, RunStatusSucceeded, EventTypeSucceeded); err != nil {
				return nil, true, err
			}
			return nil, true, nil
		}
		proposal := proposed.Proposal

		capturedConn, hasConn := mcpConnector(snapshot, proposal.Connector)
		autoApprove := map[string]bool{}
		if hasConn && connectorAutoApproves(capturedConn, RawToolName(proposal.Tool, proposal.Connector)) {
			autoApprove[proposal.Tool] = true
		}

		auth, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (Authorization, error) {
			return e.gateway.Propose(stepCtx, ProposeActionInput{
				InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Actor: in.Actor,
				Proposal: proposal, AutoApprove: autoApprove,
			})
		}, dbos.WithStepName(fmt.Sprintf("gateway_propose_%d_%d", seg, i)))
		if err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "authorize action", err)
			return nil, true, err
		}

		actions[i] = governedAction{
			proposal: proposal, auth: auth,
			capturedConn: capturedConn, hasConn: hasConn, autoApprove: autoApprove,
		}
	}

	// Phase 2: create all manual approvals at once in one durable step, so the inbox shows every pending decision the
	// moment the run pauses. Auto-approved and policy-denied actions skip approval creation.
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		for i := range actions {
			a := &actions[i]
			if a.auth.Decision != PolicyApprovalRequired {
				continue
			}
			approvalID := ApprovalIDForToolCall(runID, a.proposal.ToolCallID)
			dedupeKey := actionDedupeKey(EventTypeWaitingApproval, a.proposal.ToolCallID)
			if err := e.store.RecordTransition(stepCtx, RunTransition{
				InstanceID: in.InstanceID, RunID: runID,
				Status: RunStatusWaitingApproval, EventType: EventTypeWaitingApproval,
				DedupeKey: dedupeKey,
			}); err != nil {
				return false, err
			}
			if err := e.store.CreateApproval(stepCtx, NewApproval{
				ApprovalID: approvalID,
				RunID:      runID,
				InstanceID: in.InstanceID,
				ToolName:   a.proposal.Tool,
				Connector:  a.proposal.Connector,
				ToolCallID: a.proposal.ToolCallID,
				ArgsHash:   a.auth.ArgsHash,
				// The preimage travels on the Authorization, captured from the same CanonicalizeArgs output the
				// hash was computed over; recomputing it here from a.proposal.Args could diverge from the hash
				// (the gateway hashed its own frozen copy) and would break the "you sign what you see" invariant.
				CanonicalArgs: canonicalArgsForApproval(a),
				Proposal:      a.proposal.Summary,
				Policy:        string(a.auth.Decision),
				RequestedBy:   in.Actor.Subject,
				Position:      i + 1,
				Total:         n,
			}); err != nil {
				return false, err
			}
		}
		return true, nil
	}, dbos.WithStepName(fmt.Sprintf("setup_approvals_%d", seg)))
	if err != nil {
		*res, err = e.fail(ctx, in, runID, *res, "setup approvals", err)
		return nil, true, err
	}

	// Tell whoever may decide that the run is waiting on them (kairos-cloud#143), in a durable step of its own.
	//
	// This used to sit outside any step, so that a recovery would notify again rather than risk staying silent.
	// The duplicate was not as harmless as that assumed: the tag collapses the banner but not the buzz, and once
	// a notification leaves a record (kairos-cloud#148) a replay leaves another row. A run recovered after days
	// would announce itself again as if it were new.
	//
	// A step is where DBOS says this belongs: calling a third-party API from the workflow body is the thing its
	// determinism rules tell you not to do, and it was being done here. A completed step is never re-executed,
	// so recovery no longer re-announces every parked run.
	//
	// Its own step rather than sharing the one above is what keeps the other half: if the process dies between
	// the two, the notification step never completed and recovery does send it.
	//
	// This is at-least-once, not exactly-once, and DBOS is explicit about it: a crash after the push leaves the
	// notifier but before the checkpoint lands will notify again on recovery. That window is narrow and its
	// duplicate is the harmless direction. What it replaces was not narrow at all — every restart re-announced
	// every parked run, so a runtime with ten waiting sent ten notifications to every device.
	_, err = dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		e.notifyPendingApprovals(stepCtx, in, runID, actions)
		// The notifier's contract is that it never fails, and a notification is not worth a run: reporting an
		// error here would only make DBOS retry a step whose work is already done.
		return true, nil
	}, dbos.WithStepName(fmt.Sprintf("notify_approvals_%d", seg)))
	if err != nil {
		// DBOS failed to record the step, not the notifier. The approvals are committed and the run must reach
		// its wait, so this is logged and swallowed like the notification itself.
		if e.logger != nil {
			e.logger.Error("act: recording the approval notification step failed", "run", runID, "err", err)
		}
	}

	// Phase 3: process each action in proposal order. Auto-approved actions execute immediately; manual ones wait for
	// the human's decision on a per-tool-call DBOS topic, so decisions arriving in any order are buffered and consumed
	// correctly. A rejection or policy denial injects an error result for that call; the run continues to the next
	// action. An indeterminate or workflow error is terminal.
	results := make([]*ai.InjectedResult, n)
	for i := range actions {
		// hasMoreManual tells processOneAction whether a later action in the batch still needs a human decision. When
		// true, the post-decision status stays waiting_approval (so the run remains visible in the inbox and the Home
		// pending filter); when false, it moves to running (no more gates ahead). This is a pure function of the batch
		// structure, deterministic, and recovery-safe.
		hasMoreManual := false
		for j := i + 1; j < len(actions); j++ {
			if actions[j].auth.Decision == PolicyApprovalRequired {
				hasMoreManual = true
				break
			}
		}
		// Withdrawing the run's approvals is not done here: it belongs to the store write that marks the run
		// terminal (CloseRun), so the two can never be observed half-applied. Sweeping them from this level is what
		// stranded runs reading waiting_approval with every approval cancelled and no decider, since the paths that
		// end a run WITHOUT marking it terminal — a shutdown interrupting the approval wait above all — swept anyway
		// (kairos-cloud#135). Only the ledger, which may sit behind the gateway and cannot join that transaction,
		// is closed out from here.
		result, terminal, err := e.processOneAction(ctx, in, runID, res, actions[i], seg, i, hasMoreManual)
		if err != nil {
			// A cancelled workflow context leaves the run to DBOS (see fail), so its actions stay as they are for the
			// recovered run to pick up.
			if !errors.Is(err, errInterruptedAwaitingApproval) && ctx.Err() == nil {
				e.withdrawPendingActions(in.InstanceID, runID)
			}
			return nil, true, err
		}
		if terminal {
			e.withdrawPendingActions(in.InstanceID, runID)
			return nil, true, nil
		}
		results[i] = result
	}

	return results, false, nil
}

// notifyPendingApprovals hands the batch's manual approvals to the configured ApprovalNotifier. Auto-approved and
// policy-denied actions are left out: nobody has to sign them, and a batch with none of them notifies nobody.
func (e *DBOSExecutor) notifyPendingApprovals(ctx context.Context, in AgentRunInput, runID string, actions []governedAction) {
	if e.approvals == nil {
		return
	}

	// The approvals are already committed and the run is about to park on them: whatever the notifier does, the run
	// must reach that wait. Its contract says it never fails, so this catches only a bug in it, but a panic here
	// would unwind the workflow and lose a run that has nothing wrong with it.
	defer func() {
		if r := recover(); r != nil && e.logger != nil {
			e.logger.Error("act: notifying the pending approvals panicked", "run", runID, "panic", r)
		}
	}()

	pending := make([]PendingAction, 0, len(actions))
	for i := range actions {
		if actions[i].auth.Decision != PolicyApprovalRequired {
			continue
		}
		pending = append(pending, PendingAction{
			Tool:      RawToolName(actions[i].proposal.Tool, actions[i].proposal.Connector),
			Connector: actions[i].proposal.Connector,
		})
	}
	if len(pending) == 0 {
		return
	}

	e.approvals.NotifyPendingApprovals(ctx, PendingApprovals{
		InstanceID: in.InstanceID,
		RunID:      runID,
		AgentName:  in.AgentName,
		RunActor:   in.Actor.Subject,
		Actions:    pending,
	})
}

// processOneAction handles policy, approval wait, and execution for one action within a governed batch. It returns the
// InjectedResult for the model (terminal=false) or sets *res for a truly terminal outcome (terminal=true). A rejection
// or policy denial injects an error result so the model sees the outcome; only indeterminate/workflow errors end
// the run. hasMoreManual indicates whether a later action in the batch still requires a human decision: when true, the
// post-decision status stays waiting_approval so the run remains visible in the inbox; when false, it moves to running.
func (e *DBOSExecutor) processOneAction(ctx dbos.Context, in AgentRunInput, runID string, res *AgentRunResult, a governedAction, seg, idx int, hasMoreManual bool) (*ai.InjectedResult, bool, error) {
	proposal := a.proposal
	auth := a.auth

	switch auth.Decision {
	case PolicyDeny:
		return &ai.InjectedResult{
			ToolCallID: proposal.ToolCallID, Tool: proposal.Tool, IsError: true,
			Message: "This action was denied by policy: " + auth.Reason,
		}, false, nil

	case PolicyApprovalRequired:
		approvalID := ApprovalIDForToolCall(runID, proposal.ToolCallID)
		topic := approvalTopic + ":" + proposal.ToolCallID

		decision, err := dbos.Recv[string](ctx, topic, approvalWaitForever)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil, true, interruptedAwaitingApproval(err)
			}
			*res, err = e.fail(ctx, in, runID, *res, "await approval", err)
			return nil, true, err
		}
		if ApprovalDecision(decision) != ApprovalApproved {
			denier, err := e.loadDecider(ctx, in, approvalID, fmt.Sprintf("%s_%d_%d", stepLoadDenier, seg, idx))
			if err != nil {
				*res, err = e.fail(ctx, in, runID, *res, "load denier", err)
				return nil, true, err
			}
			postRejectStatus := RunStatusRunning
			if hasMoreManual {
				postRejectStatus = RunStatusWaitingApproval
			}
			if err := e.emitWithPayload(ctx, in, runID, postRejectStatus, EventTypeRejected,
				actionDedupeKey(EventTypeRejected, proposal.ToolCallID), deciderPayload(denier)); err != nil {
				// A real failure, not a shutdown: end the run through fail so it is marked terminal and its remaining
				// approvals are withdrawn with it. Returning the bare error would leave the run reading
				// waiting_approval forever, which is the state this release exists to stop producing.
				*res, err = e.fail(ctx, in, runID, *res, "emit rejected", err)
				return nil, true, err
			}
			e.closeLedgerAction(ctx, in.InstanceID, runID, proposal.ToolCallID, ActionRejected, denier)
			return &ai.InjectedResult{
				ToolCallID: proposal.ToolCallID, Tool: proposal.Tool, IsError: true,
				Message: "This action was rejected by the operator.",
			}, false, nil
		}
		decidedBy, err := e.loadDecider(ctx, in, approvalID, fmt.Sprintf("%s_%d_%d", stepLoadApprover, seg, idx))
		if err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "load approver", err)
			return nil, true, err
		}
		if _, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
			return true, e.gateway.RecordApproved(stepCtx, ApprovedInput{
				InstanceID: in.InstanceID, RunID: runID, ToolCallID: proposal.ToolCallID,
				ArgsHash: auth.ArgsHash, DecidedBy: decidedBy,
			})
		}, dbos.WithStepName(fmt.Sprintf("gateway_record_approved_%d_%d", seg, idx))); err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "record approval", err)
			return nil, true, err
		}
		postResumeStatus := RunStatusRunning
		if hasMoreManual {
			postResumeStatus = RunStatusWaitingApproval
		}
		if err := e.emitWithPayload(ctx, in, runID, postResumeStatus, EventTypeResumed,
			actionDedupeKey(EventTypeResumed, proposal.ToolCallID), deciderPayload(decidedBy)); err != nil {
			// As above: a failure here is not a shutdown, so the run must be closed rather than left waiting.
			*res, err = e.fail(ctx, in, runID, *res, "emit resumed", err)
			return nil, true, err
		}

	case PolicyAllow:
		// Auto-approved: proceed to execution.

	default:
		return &ai.InjectedResult{
			ToolCallID: proposal.ToolCallID, Tool: proposal.Tool, IsError: true,
			Message: "This action was blocked: undecidable policy.",
		}, false, nil
	}

	// Execute the write through the gateway.
	report, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (ExecuteReport, error) {
		return e.gateway.Execute(withGovernedConnector(stepCtx, in.InstanceID, a.capturedConn, a.hasConn), ExecuteActionInput{
			InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Proposal: proposal, TraceID: runID,
			AutoApprove: a.autoApprove,
		})
	}, append(leaseHeldRetries(e.gateway.leaseDuration()), dbos.WithStepName(fmt.Sprintf("gateway_execute_%d_%d", seg, idx)))...)
	if err != nil {
		*res, err = e.fail(ctx, in, runID, *res, "apply action", err)
		return nil, true, err
	}

	switch report.Outcome {
	case OutcomeSucceeded:
		res.ActionTaken = true
		res.ActionRef = report.ExternalReference
		if auth.Verifiable {
			if _, verifyErr := dbos.RunAsStep(ctx, func(stepCtx context.Context) (VerifyResult, error) {
				return e.gateway.Verify(withGovernedConnector(stepCtx, in.InstanceID, a.capturedConn, a.hasConn), ExecuteActionInput{
					InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Proposal: proposal, TraceID: runID,
				})
			}, dbos.WithStepName(fmt.Sprintf("gateway_verify_%d_%d", seg, idx))); verifyErr != nil && e.logger != nil {
				e.logger.Warn("act: verification persistence failed; write is confirmed but audit record is incomplete",
					"run", runID, "err", verifyErr)
			}
		}
		message := report.Message
		if message == "" {
			message = report.ExternalReference
		}
		return &ai.InjectedResult{ToolCallID: proposal.ToolCallID, Tool: proposal.Tool, Message: message}, false, nil

	case OutcomeIndeterminate:
		res.Status = RunStatusFailed
		if err := e.emitFailedReason(ctx, in, runID, "action result indeterminate: awaiting human resolution"); err != nil {
			// The run IS over (the write's outcome is unknown and needs a human), so if recording that failed, close it
			// through fail, which retries on a context that survives a shutdown. Returning here would leave a finished
			// run reading waiting_approval.
			*res, err = e.fail(ctx, in, runID, *res, "emit indeterminate", err)
			return nil, true, err
		}
		return nil, true, nil

	default:
		return &ai.InjectedResult{
			ToolCallID: proposal.ToolCallID, Tool: proposal.Tool, IsError: true,
			Message: "The action failed: " + report.Message,
		}, false, nil
	}
}

// leaseHeldRetries makes the execute step wait out a write another attempt still holds, and retry nothing else.
//
// A recovered run can reach its write while the lease of the attempt that died mid-write is still live: the gateway
// then answers ErrLeaseHeld, meaning "not yet", not "failed". DBOS steps do not retry unless told to (the default is
// zero), so without this the run failed on the first answer and the action stayed executing with nothing left to
// revisit it (kairos-cloud#170). The budget outlasts the lease, so the last retry finds the owner's recorded outcome or
// an expired lease that the gateway reclaims to indeterminate.
//
// Retries run inside the same step and take no step slot of their own, so this does not change the workflow's shape
// and needs no dbos.Patch (dbos executeStepWithRetry, v1.2.0).
func leaseHeldRetries(lease time.Duration) []dbos.StepOption {
	const interval = 5 * time.Second
	return []dbos.StepOption{
		dbos.WithStepRetryPredicate(func(err error) bool { return errors.Is(err, ErrLeaseHeld) }),
		dbos.WithStepBaseInterval(time.Second),
		dbos.WithStepMaxInterval(interval),
		dbos.WithStepMaxRetries(int(lease/interval) + 3),
	}
}

// emitFailedReason records a failed run transition carrying a sanitized reason, as a durable step, and returns nil on
// success. Unlike fail(), it is for a business-terminal failure (policy denied, action failed/indeterminate) that
// ends the run without a workflow error, so DBOS does not retry it. It closes the run rather than only recording the
// transition, so the approvals of a batch whose later actions will never be reached are withdrawn with it.
func (e *DBOSExecutor) emitFailedReason(ctx dbos.Context, in AgentRunInput, runID, reason string) error {
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		_, closeErr := e.store.CloseRun(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusFailed,
			EventType:  EventTypeFailed,
			Error:      reason,
		})
		return true, closeErr
	}, dbos.WithStepName("emit:"+EventTypeFailed))
	if err != nil {
		return fmt.Errorf("act: emit failed: %w", err)
	}
	return nil
}

// approvalIDSpace is the UUID v5 namespace of approval IDs. Changing it changes every approval ID a run derives.
var approvalIDSpace = uuid.MustParse("5f0c7a52-3d1e-4b8a-9c6f-2e4d8b1a7f30")

// ApprovalIDForToolCall derives the approval ID of a tool call within a run: a UUID v5 of the pair (ADR-0034). It is
// deterministic because the workflow derives it again after the decision to learn who signed, and opaque so it is safe
// in a URL whatever the model provider puts in the tool call ID. The run ID is length-prefixed so no two pairs render
// alike. Approvals recorded before ADR-0034 keep the ID they were stored with.
func ApprovalIDForToolCall(runID, toolCallID string) string {
	return uuid.NewSHA1(approvalIDSpace, fmt.Appendf(nil, "%d:%s/%s", len(runID), runID, toolCallID)).String()
}

// actionDedupeKey scopes an approval-cycle event (waiting_approval, resumed, rejected) to the tool call that produced
// it, so a segment that governs N actions records N instances of each edge while a replay of the SAME edge still
// deduplicates (kairos-cloud#129).
func actionDedupeKey(eventType, toolCallID string) string { return eventType + ":" + toolCallID }

// HashArgs returns the canonical hash an approval is bound to: the decision is valid only for these exact arguments
// (§6.4, §11.2): the "sha256:"-prefixed hash of the canonicalized tool arguments. Approval.VerifiedCanonicalArgs
// verifies a stored preimage with it, including the proposal texts that approvals from the retired simulated path
// hashed.
func HashArgs(args string) string {
	sum := sha256.Sum256([]byte(args))
	return "sha256:" + hex.EncodeToString(sum[:])
}
