package act

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
)

const (
	appName       = "kairos-act"
	workflowName  = "RunAgentWorkflow" // stable name: required for cross-process enqueue and recovery
	approvalTopic = "approval"

	stepLoadSnapshot        = "load_snapshot"
	stepRunAgent            = "run_agent"
	stepPersistConversation = "persist_conversation"
	stepApplyAction         = "apply_action"
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

	// defaultApplicationVersion pins the DBOS application version. It MUST be stable across process restarts:
	// otherwise every rebuild is a fresh "version" and a restarted worker ignores runs the previous binary started,
	// which hangs recovery. In production this is the deploy SHA (§19.1).
	defaultApplicationVersion = "act-spike-v1"

	// approvalWaitForever is the duration a run waits for a human decision. dbos.Recv has no "wait forever" argument and
	// a zero timeout fires immediately, so we pass this effectively-infinite duration (~100 years, comfortably within
	// int64 nanoseconds). The wait is durable and cheap (a parked workflow blocked on Recv), and Cancel is the human's
	// way to abort a run nobody will decide.
	approvalWaitForever = 100 * 365 * 24 * time.Hour
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
	// Store, when set, is the queryable run store the executor writes lifecycle state and events to (§15). It is the
	// product's state plane, separate from the DBOS system tables: the API reads runs from it without touching the
	// workflow engine. Optional — a nil Store keeps the executor's original spike behavior (no state emission), which
	// is what the executor's own unit tests rely on.
	Store RunStore
	// Gateway and Proposer, when both set, switch the run's action phase from the Fase 1 simulated write (Runner.
	// ApplyAction / ActionSink) to the real Fase 2 flow: the Proposer surfaces the proposed tool call, and the Gateway
	// enforces policy, records the action ledger, resolves secrets, executes via ActionExecutor and verifies (§16.2,
	// §12). When either is nil the executor keeps the simulated path, so the merged Fase 1 tests run unchanged. The
	// Gateway's Ledger should share the same Postgres as Store (its agent_actions table lives in the product schema).
	Gateway  *Gateway
	Proposer Proposer
	Logger   *slog.Logger
}

// DBOSExecutor is the DBOS-backed AgentExecutor. It both embeds the worker (it registers the workflow and runs
// recovered/started runs) and drives runs from the outside (Start/Resume/Cancel), which is the single-process
// dev topology (§20.1). In production the worker and the runtime that drives it are separate processes sharing the
// same Postgres; the split does not change this code, only which side holds the dbos.Client.
type DBOSExecutor struct {
	ctx      dbos.Context
	client   dbos.Client
	runner   Runner
	store    RunStore
	gateway  *Gateway
	proposer Proposer
	version  string
	logger   *slog.Logger
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

	version := cfg.ApplicationVersion
	if version == "" {
		version = defaultApplicationVersion
	}

	// Tag the executor's context as Act traffic. The API handlers already tag theirs, but a run does not execute in the
	// request that started it: RunWorkflow below hands the work to DBOS on this context, and recovery after a restart
	// has no request to inherit from at all. Without this, every token and tool call an agent spends is recorded with
	// an empty source, so the meter cannot tell agent spend from chat spend.
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)

	dctx, err := dbos.NewContext(ctx, dbos.Config{
		AppName:            appName,
		DatabaseURL:        cfg.DatabaseURL,
		DatabaseSchema:     cfg.DatabaseSchema,
		ApplicationVersion: version,
		Logger:             cfg.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("act: init DBOS: %w", err)
	}

	e := &DBOSExecutor{
		ctx:      dctx,
		runner:   cfg.Runner,
		store:    cfg.Store,
		gateway:  cfg.Gateway,
		proposer: cfg.Proposer,
		version:  version,
		logger:   cfg.Logger,
	}

	// Register the ONE generic workflow. It is a bound method so the workflow body reaches the runner (and, in
	// production, the runtime) without those crossing the serialized workflow boundary. Recovery re-invokes this
	// same registered closure by name, so a recovered run uses the recovering process's runner.
	dbos.RegisterWorkflow(dctx, e.runAgentWorkflow, dbos.WithWorkflowName(workflowName))

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

// ComposeRunID builds the durable workflow ID from (instance, agent, idempotency key). It length-prefixes the two
// namespacing components so the encoding is injective: without the prefixes, ("a/b","c") and ("a","b/c") would
// both render as "a/b/c" and collide onto one run. The idempotency key is last and unconstrained, so it needs no
// prefix. The result stays readable in logs while remaining unambiguous whatever characters the fields contain.
//
// It is exported so a caller that needs to address a run it did not start (the trigger dispatcher deduplicating an
// event, for example) can compute the same ID from its parts.
func ComposeRunID(instanceID, agentName, idempotencyKey string) string {
	return fmt.Sprintf("%d:%s/%d:%s/%s", len(instanceID), instanceID, len(agentName), agentName, idempotencyKey)
}

// Start launches a run. The workflow ID is the idempotency key namespaced by instance and agent, so a repeated
// Start with the same (instance, agent, key) attaches to the existing run (OAOO) instead of starting a second one;
// after the run finishes it returns the cached result. It returns that composed run ID, which Resume/Cancel/Result
// address the run by.
func (e *DBOSExecutor) Start(ctx context.Context, in AgentRunInput) (string, error) {
	if in.InstanceID == "" || in.AgentName == "" {
		return "", errors.New("act: AgentRunInput.InstanceID and AgentName are required")
	}
	if in.IdempotencyKey == "" {
		return "", errors.New("act: AgentRunInput.IdempotencyKey is required")
	}
	// Namespace the workflow ID by instance and agent: two instances (tenants) or two agents that reuse the same
	// idempotency key must not collide onto one run, or an approval could land on the wrong run. Detecting a
	// different payload under the same key is a Phase 1 loose end: it requires reading back the stored input.
	runID := ComposeRunID(in.InstanceID, in.AgentName, in.IdempotencyKey)

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

	// Record the run as queued BEFORE enqueuing the workflow, so it is queryable the instant Start returns 202
	// (§14.2). CreateRun is idempotent on the run ID, so a repeated Start (same idempotency key) does not duplicate
	// it, matching the workflow's own OAOO enqueue below.
	if e.store != nil {
		if err := e.store.CreateRun(ctx, NewRun{
			RunID:          runID,
			InstanceID:     in.InstanceID,
			AgentName:      in.AgentName,
			Trigger:        trigger,
			TriggerRef:     in.TriggerRef,
			IdempotencyKey: in.IdempotencyKey,
			ConversationID: in.ConversationID,
			Actor:          in.Actor,
		}); err != nil {
			return "", fmt.Errorf("act: record run %q: %w", runID, err)
		}
	}

	_, err := dbos.RunWorkflow(e.ctx, e.runAgentWorkflow, in, dbos.WithWorkflowID(runID))
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
	topic := approvalTopic
	if toolCallID != "" {
		topic = approvalTopic + ":" + toolCallID
	}
	if err := dbos.Send(e.client, runID, string(decision), topic); err != nil {
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
// The transition is scoped by instance, so a caller can only cancel a run in the instance it addresses. Overwriting
// an already-terminal run's status is a Phase 1 loose end (RecordTransition does not yet guard terminal states).
func (e *DBOSExecutor) Cancel(ctx context.Context, instanceID, runID string) error {
	if err := dbos.CancelWorkflow(e.client, runID); err != nil {
		return fmt.Errorf("act: cancel run %q: %w", runID, err)
	}
	if e.store == nil {
		return nil
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

// Result waits for a run to finish and returns its checkpointed outcome. It is a spike-only observability hook; the
// product surfaces run state through the Agent API (§14), not through the executor.
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
//
// Step granularity trade-off (spike): the whole model/tool loop runs as ONE step (run_agent). That is the pragmatic
// choice for the spike — it keeps the durable surface small and reuses runtime/ai unchanged — but it means a crash
// inside the loop replays the entire loop, not just the last model turn, and re-drives the LLM. Phase 1 must split
// this into one durable step per model turn and per tool call (§10.3), pairing each write step with an action
// ledger (§12), so the loop resumes at the last completed turn and external writes are not re-driven.
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

	// Fase 2 action phase: when a gateway and proposer are wired, drive the run as a segmented durable loop — each
	// segment runs the model/tool loop until it pauses on a governed write, the write is governed and executed, and its
	// result is injected so the next segment's model sees the outcome and continues (or closes). Otherwise fall back to
	// the Fase 1 single-segment path below (simulated write), so the merged executor tests run unchanged.
	if e.gateway != nil && e.proposer != nil {
		return e.runSegmentedLoop(ctx, in, runID, snapshot, res)
	}

	// Step 2 (Fase 1): run the agent's model/tool loop from the fixed snapshot as ONE durable step (see the granularity
	// trade-off above), under the initiating actor's checkpointed claims so the tool intersection reflects the
	// initiator's authority, not the worker's. The step returns a named struct so DBOS checkpoints the response and the
	// session ID together: on replay the session ID is the same deterministic value, not a fresh session opened by
	// re-execution.
	out, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (runSegmentStepResult, error) {
		r, runErr := e.runner.RunSegment(stepCtx, RunSegmentInput{
			InstanceID: in.InstanceID, Claims: in.Actor.Claims, Snapshot: snapshot, Prompt: in.Prompt,
		})
		return runSegmentStepResult{Response: r.Response, SessionID: r.SessionID, Proposals: r.Proposed}, runErr
	}, dbos.WithStepName(stepRunAgent))
	if err != nil {
		return e.fail(ctx, in, runID, res, "run agent", err)
	}
	res.Response = out.Response

	// Persist the AI session the run executed in as the run's conversation, so the API and UI can render the run as a
	// chat (§13.1). This overwrites the conversation_id CreateRun seeded (normally empty) with the real session. It is
	// a durable, idempotent step, so a replay re-persists the same checkpointed session ID.
	if err := e.persistConversation(ctx, in, runID, out.SessionID); err != nil {
		return e.fail(ctx, in, runID, res, "persist conversation", err)
	}

	// Suspend for a durable human approval. Before waiting, record the run as waiting_approval and persist the exact
	// proposal as an approval request (§11.2), so the approval inbox and the API can see and act on it while the
	// workflow is parked. The setup is a durable, idempotent step, so recovery does not duplicate the approval.
	if err := e.setupApproval(ctx, in, runID, out.Response); err != nil {
		return e.fail(ctx, in, runID, res, "await approval setup", err)
	}

	// Wait durably for a decision. The wait is indefinite: Cancel is the human's way to abort a run nobody will decide.
	decision, err := dbos.Recv[string](ctx, approvalTopic, approvalWaitForever)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return res, interruptedAwaitingApproval(err)
		}
		return e.fail(ctx, in, runID, res, "await approval", err)
	}
	if ApprovalDecision(decision) != ApprovalApproved {
		res.Status = RunStatusRejected
		denier, err := e.loadDecider(ctx, in, ApprovalIDForRun(runID), stepLoadDenier)
		if err != nil {
			return e.fail(ctx, in, runID, res, "load denier", err)
		}
		if err := e.emitWithPayload(ctx, in, runID, RunStatusRejected, EventTypeRejected, segmentDedupeKey(EventTypeRejected, 0), deciderPayload(denier)); err != nil {
			return res, err
		}
		return res, nil
	}

	// Approved: the run resumes and proceeds to the (approved) external write. The resumed event carries the real
	// approver so the timeline attributes the decision without a ledger read.
	approver, err := e.loadDecider(ctx, in, ApprovalIDForRun(runID), stepLoadApprover)
	if err != nil {
		return e.fail(ctx, in, runID, res, "load approver", err)
	}
	if err := e.emitWithPayload(ctx, in, runID, RunStatusRunning, EventTypeResumed, segmentDedupeKey(EventTypeResumed, 0), deciderPayload(approver)); err != nil {
		return res, err
	}

	// Step 3: the external write, executed exactly once on the happy path and only after approval.
	action, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (ActionResult, error) {
		return e.runner.ApplyAction(stepCtx, ActionRequest{
			RunID:        runID,
			AgentName:    in.AgentName,
			Actor:        in.Actor,
			SnapshotMark: snapshot.Instructions,
			Response:     out.Response,
		})
	}, dbos.WithStepName(stepApplyAction))
	if err != nil {
		return e.fail(ctx, in, runID, res, "apply action", err)
	}
	res.Status = RunStatusSucceeded
	res.ActionTaken = true
	res.ActionRef = action.Ref
	if err := e.emit(ctx, in, runID, RunStatusSucceeded, EventTypeSucceeded); err != nil {
		return res, err
	}
	return res, nil
}

// emitRunning records the queued->running transition and binds the run to the snapshot's spec hash in one durable
// step. It is separate from emit because it is the only transition that carries a spec hash: it fires right after
// the snapshot is captured, so the run row records exactly which agent version it executed (§8.3, §12).
func (e *DBOSExecutor) emitRunning(ctx dbos.Context, in AgentRunInput, runID, specHash string) error {
	if e.store == nil {
		return nil
	}
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

// persistConversation links the run to the AI session it executed in, as a durable step. It is a no-op when no store
// is configured (the storeless unit tests) or when no session was opened (SessionID empty), so it never overwrites an
// existing link with a blank. The store UPDATE is idempotent, so a replay re-writes the same checkpointed session ID.
func (e *DBOSExecutor) persistConversation(ctx dbos.Context, in AgentRunInput, runID, sessionID string) error {
	if e.store == nil || sessionID == "" {
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

// emit records one lifecycle transition as a durable step. It is a no-op when no store is configured (the executor's
// own unit tests run storeless), so wiring a store adds state emission without changing the storeless step sequence.
// The step is idempotent in the store (one event per edge), so a crash between the side effect and the checkpoint
// replays it harmlessly.
func (e *DBOSExecutor) emit(ctx dbos.Context, in AgentRunInput, runID string, status RunStatus, eventType string) error {
	return e.emitWithPayload(ctx, in, runID, status, eventType, "", nil)
}

// emitWithPayload is emit with structured detail attached to the event, deduplicated by dedupeKey (empty keys the
// edge on its event type alone — once per run). The payload must already be redacted: it is stored verbatim and read
// by the run timeline, which is user-visible (§17.2). A nil payload behaves exactly like emit, so the step name and
// shape are unchanged for the transitions that carry no detail.
func (e *DBOSExecutor) emitWithPayload(ctx dbos.Context, in AgentRunInput, runID string, status RunStatus, eventType, dedupeKey string, payload map[string]any) error {
	if e.store == nil {
		return nil
	}
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

// segmentDedupeKey scopes an approval-cycle event (waiting_approval, resumed, rejected) to the segment whose
// governed action produced it. The segmented loop pauses once per action, so these edges legitimately repeat within a
// run; keying them per segment lets the store record each occurrence — moving the status with it — while a replay of
// the SAME segment's edge still deduplicates. The once-per-run edges keep the strict key (their bare event type), so
// a stale replay can never regress the run's status. Fixes kairos-cloud#129.
func segmentDedupeKey(eventType string, seg int) string { return eventType + ":" + strconv.Itoa(seg) }

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
// delivered, so it is authoritative here (§11.2, §18.3). Falls back to the initiator when no store is wired (the
// executor's own unit tests) or when a non-API resume path left decided_by unset — see the Resume TODO above.
func (e *DBOSExecutor) loadDecider(ctx dbos.Context, in AgentRunInput, approvalID, stepName string) (string, error) {
	if e.store == nil {
		return in.Actor.Subject, nil
	}
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
func (e *DBOSExecutor) fail(ctx dbos.Context, in AgentRunInput, runID string, res AgentRunResult, what string, err error) (AgentRunResult, error) {
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
// The store write is retried on a context that outlives the workflow's, because a shutdown cancels that one and the
// run would otherwise be left mid-flight. Doing the two writes separately is what stranded runs reading
// waiting_approval with every approval cancelled and no decider (kairos-cloud#135), so CloseRun binds them: either
// the run is terminal AND its approvals are withdrawn, or neither happened and the run stays decidable.
func (e *DBOSExecutor) closeRun(ctx dbos.Context, in AgentRunInput, runID string, tr RunTransition, what string) {
	if e.store == nil {
		return
	}
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

// ledger returns the action ledger for this executor: the gateway's ledger when wired, or the store itself when it
// implements ActionLedger (production: PostgresRunStore serves both). Returns nil when neither is available.
func (e *DBOSExecutor) ledger() ActionLedger {
	if e.gateway != nil {
		return e.gateway.Ledger
	}
	if l, ok := e.store.(ActionLedger); ok {
		return l
	}
	return nil
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

// setupApproval records the waiting_approval transition and persists the approval request atomically, in one durable
// step. The approval identity and the args hash are derived deterministically from the run and proposal, so a replay
// creates the same approval (idempotent) rather than a second one.
func (e *DBOSExecutor) setupApproval(ctx dbos.Context, in AgentRunInput, runID, proposal string) error {
	if e.store == nil {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		if err := e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusWaitingApproval,
			EventType:  EventTypeWaitingApproval,
			DedupeKey:  segmentDedupeKey(EventTypeWaitingApproval, 0), // the Fase 1 path has exactly one action: segment 0
		}); err != nil {
			return false, err
		}
		return true, e.store.CreateApproval(stepCtx, NewApproval{
			ApprovalID: ApprovalIDForRun(runID),
			RunID:      runID,
			InstanceID: in.InstanceID,
			ToolName:   proposedActionTool,
			ArgsHash:   HashArgs(proposal),
			// This simulated path hashes the proposal text itself, so the proposal IS the preimage: persisting it
			// as the canonical args keeps the invariant (stored bytes hash to ArgsHash) uniform across both paths.
			CanonicalArgs: proposal,
			Proposal:      proposal,
			Policy:        "approval_required",
			RequestedBy:   in.Actor.Subject,
		})
	}, dbos.WithStepName("setup_approval"))
	if err != nil {
		return fmt.Errorf("act: setup approval: %w", err)
	}
	return nil
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
// delivery (§12). A crash mid-segment re-drives that segment's model turns (the accepted spike granularity: replaying
// a segment re-flushes its trace but never re-executes a committed external write).
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
	if e.store != nil {
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
			if !errors.Is(err, errInterruptedAwaitingApproval) {
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
	}, dbos.WithStepName(fmt.Sprintf("gateway_execute_%d_%d", seg, idx)))
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

// emitFailedReason records a failed run transition carrying a sanitized reason, as a durable step, and returns nil on
// success. Unlike fail(), it is for a business-terminal failure (policy denied, action failed/indeterminate) that
// ends the run without a workflow error, so DBOS does not retry it. It closes the run rather than only recording the
// transition, so the approvals of a batch whose later actions will never be reached are withdrawn with it.
func (e *DBOSExecutor) emitFailedReason(ctx dbos.Context, in AgentRunInput, runID, reason string) error {
	if e.store == nil {
		return nil
	}
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

// proposedActionTool is the synthetic tool name the spike's single simulated action is proposed under. Production
// replaces it with the concrete tool the gateway proposes.
const proposedActionTool = "act.propose_action"

// ApprovalIDForRun derives the deterministic approval ID for a run's FIRST action. It is kept for the Fase 1 path and
// the single-action callers/tests; the gateway loop keys per segment via ApprovalIDForSegment.
func ApprovalIDForRun(runID string) string { return ApprovalIDForSegment(runID, 0) }

// ApprovalIDForSegment derives the deterministic approval ID for the action a given segment proposed, so a run that
// proposes several sequential actions gets a distinct approval (and inbox row) per action. Segment 0 keeps the
// historical ":1" suffix, so existing single-action runs, approvals and tests are unaffected. Retained for the
// Fase 1 (simulated) path and backward compatibility; the gateway path uses ApprovalIDForToolCall.
func ApprovalIDForSegment(runID string, seg int) string { return runID + ":" + strconv.Itoa(seg+1) }

// ApprovalIDForToolCall derives the deterministic approval ID for a specific tool call within a run. It uses a ":tc:"
// infix so old segment-based IDs (":1", ":2") and new tool-call-based IDs (":tc:call-1") never collide, and an
// existing pending approval with a segment-based ID remains addressable after the upgrade.
func ApprovalIDForToolCall(runID, toolCallID string) string { return runID + ":tc:" + toolCallID }

// actionDedupeKey scopes an approval-cycle event (waiting_approval, resumed, rejected) to the tool call that produced
// it, so a segment that governs N actions records N instances of each edge while a replay of the SAME edge still
// deduplicates. This is the per-tool-call generalization of segmentDedupeKey (#129), keyed by tool_call_id instead of
// segment index.
func actionDedupeKey(eventType, toolCallID string) string { return eventType + ":" + toolCallID }

// HashArgs returns the canonical hash an approval is bound to: the decision is valid only for these exact arguments
// (§6.4, §11.2). The spike's simulated path hashes the proposal text; the gateway path hashes the canonicalized tool
// arguments (HashCanonicalArgs). Both yield the same "sha256:"-prefixed form so a consumer cannot tell them apart,
// which is also what lets Approval.VerifiedCanonicalArgs verify a stored preimage from either path with this one
// function.
func HashArgs(args string) string {
	sum := sha256.Sum256([]byte(args))
	return "sha256:" + hex.EncodeToString(sum[:])
}
