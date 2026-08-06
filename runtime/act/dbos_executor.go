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

	// payloadKeyDecidedBy carries the subject that decided an approval on the run.resumed / run.rejected events. The
	// same subject is written to the action ledger's decided_by; the event is what makes it readable through the API,
	// which does not expose the ledger.
	payloadKeyDecidedBy = "decided_by"

	// defaultApplicationVersion pins the DBOS application version. It MUST be stable across process restarts:
	// otherwise every rebuild is a fresh "version" and a restarted worker ignores runs the previous binary started,
	// which hangs recovery. In production this is the deploy SHA (§19.1).
	defaultApplicationVersion = "act-spike-v1"

	// approvalWaitForever is the duration a run waits for a human decision when no approval deadline is configured
	// (the default). dbos.Recv has no "wait forever" argument and a zero timeout fires immediately, so we pass this
	// effectively-infinite duration (~100 years, comfortably within int64 nanoseconds). The wait is durable and cheap
	// (a parked workflow blocked on Recv), and Cancel is the human's way to abort; expiry is opt-in via ApprovalTimeout.
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
	// ApprovalTimeout, when positive, is how long a run waits for a human decision before the approval expires and the
	// run ends. Zero or unset (the default) means no deadline: the run waits until a decision arrives or it is cancelled.
	ApprovalTimeout time.Duration
	Logger          *slog.Logger
}

// DBOSExecutor is the DBOS-backed AgentExecutor. It both embeds the worker (it registers the workflow and runs
// recovered/started runs) and drives runs from the outside (Start/Resume/Cancel), which is the single-process
// dev topology (§20.1). In production the worker and the runtime that drives it are separate processes sharing the
// same Postgres; the split does not change this code, only which side holds the dbos.Client.
type DBOSExecutor struct {
	ctx             dbos.DBOSContext
	client          dbos.Client
	runner          Runner
	store           RunStore
	gateway         *Gateway
	proposer        Proposer
	version         string
	approvalTimeout time.Duration
	logger          *slog.Logger
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
	// A non-positive ApprovalTimeout means no deadline (the default): a run waits for a human decision indefinitely
	// (or until cancelled), and an approval is never auto-expired. Expiry is opt-in by setting a positive timeout.
	approvalTimeout := cfg.ApprovalTimeout

	dctx, err := dbos.NewDBOSContext(ctx, dbos.Config{
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
		ctx:             dctx,
		runner:          cfg.Runner,
		store:           cfg.Store,
		gateway:         cfg.Gateway,
		proposer:        cfg.Proposer,
		version:         version,
		approvalTimeout: approvalTimeout,
		logger:          cfg.Logger,
	}

	// Register the ONE generic workflow. It is a bound method so the workflow body reaches the runner (and, in
	// production, the runtime) without those crossing the serialized workflow boundary. Recovery re-invokes this
	// same registered closure by name, so a recovered run uses the recovering process's runner.
	dbos.RegisterWorkflow(dctx, e.runAgentWorkflow, dbos.WithWorkflowName(workflowName))

	if err := dbos.Launch(dctx); err != nil {
		dbos.Shutdown(dctx, 5*time.Second)
		return nil, fmt.Errorf("act: launch DBOS: %w", err)
	}

	client, err := dbos.NewClient(ctx, dbos.ClientConfig{
		DatabaseURL:    cfg.DatabaseURL,
		DatabaseSchema: cfg.DatabaseSchema,
		Logger:         cfg.Logger,
	})
	if err != nil {
		dbos.Shutdown(dctx, 5*time.Second)
		return nil, fmt.Errorf("act: init DBOS client: %w", err)
	}
	e.client = client

	return e, nil
}

// Close drains the worker and closes the client. Callers should allow in-flight steps to finish within timeout.
func (e *DBOSExecutor) Close(timeout time.Duration) {
	if e.client != nil {
		e.client.Shutdown(timeout)
	}
	dbos.Shutdown(e.ctx, timeout)
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
func (e *DBOSExecutor) Resume(_ context.Context, runID string, decision ApprovalDecision) error {
	if err := e.client.Send(runID, string(decision), approvalTopic); err != nil {
		return fmt.Errorf("act: resume run %q: %w", runID, err)
	}
	return nil
}

// Cancel stops a run. After the workflow is cancelled it records the cancelled transition directly (not as a step,
// since the cancelled workflow will not run more steps) so the API reflects the terminal state. The store write is
// best-effort: a cancelled workflow is already stopped, so a failed state write must not turn Cancel into an error.
//
// The transition is scoped by instance, so a caller can only cancel a run in the instance it addresses. Overwriting
// an already-terminal run's status is a Phase 1 loose end (RecordTransition does not yet guard terminal states).
func (e *DBOSExecutor) Cancel(ctx context.Context, instanceID, runID string) error {
	if err := e.client.CancelWorkflow(runID); err != nil {
		return fmt.Errorf("act: cancel run %q: %w", runID, err)
	}
	if e.store == nil {
		return nil
	}
	// The workflow is already cancelled; the store writes below only reflect that, so a failure to record them must not
	// fail Cancel. They are independent: withdrawing the pending approval matters even if the run transition did not
	// land, so the inbox stops offering approve/deny on a run that will never resume (the panel is gated on pending).
	if err := e.store.RecordTransition(ctx, RunTransition{
		InstanceID: instanceID,
		RunID:      runID,
		Status:     RunStatusCancelled,
		EventType:  EventTypeCancelled,
	}); err != nil && e.logger != nil {
		e.logger.Warn("act: cancel: recording cancelled transition failed", "run", runID, "err", err)
	}
	if _, err := e.store.CancelPendingApprovals(ctx, instanceID, runID); err != nil && e.logger != nil {
		e.logger.Warn("act: cancel: withdrawing pending approvals failed", "run", runID, "err", err)
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
func (e *DBOSExecutor) runAgentWorkflow(ctx dbos.DBOSContext, in AgentRunInput) (AgentRunResult, error) {
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
		return runSegmentStepResult{Response: r.Response, SessionID: r.SessionID, Proposal: r.Proposed}, runErr
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

	// Wait durably for a decision. A timeout surfaces as a DBOS TimeoutError (not a zero value), which we read as an
	// expiration; any decision other than "approved" ends the run without touching the outside world.
	decision, err := dbos.Recv[string](ctx, approvalTopic, e.approvalRecvTimeout())
	if err != nil {
		if errors.Is(err, &dbos.DBOSError{Code: dbos.TimeoutError}) {
			res.Status = RunStatusExpired
			if err := e.recordExpiry(ctx, in, runID, ApprovalIDForRun(runID)); err != nil {
				return res, err
			}
			return res, nil
		}
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
		if err := e.emitWithPayload(ctx, in, runID, RunStatusRejected, EventTypeRejected, deciderPayload(denier)); err != nil {
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
	if err := e.emitWithPayload(ctx, in, runID, RunStatusRunning, EventTypeResumed, deciderPayload(approver)); err != nil {
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
func (e *DBOSExecutor) emitRunning(ctx dbos.DBOSContext, in AgentRunInput, runID, specHash string) error {
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
func (e *DBOSExecutor) persistConversation(ctx dbos.DBOSContext, in AgentRunInput, runID, sessionID string) error {
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
func (e *DBOSExecutor) emit(ctx dbos.DBOSContext, in AgentRunInput, runID string, status RunStatus, eventType string) error {
	return e.emitWithPayload(ctx, in, runID, status, eventType, nil)
}

// emitWithPayload is emit with structured detail attached to the event. The payload must already be redacted: it is
// stored verbatim and read by the run timeline, which is user-visible (§17.2). A nil payload behaves exactly like emit,
// so the step name and shape are unchanged for the transitions that carry no detail.
func (e *DBOSExecutor) emitWithPayload(ctx dbos.DBOSContext, in AgentRunInput, runID string, status RunStatus, eventType string, payload map[string]any) error {
	if e.store == nil {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		return true, e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     status,
			EventType:  eventType,
			Payload:    payload,
		})
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
// delivered, so it is authoritative here (§11.2, §18.3). Falls back to the initiator when no store is wired (the
// executor's own unit tests) or when a non-API resume path left decided_by unset — see the Resume TODO above.
func (e *DBOSExecutor) loadDecider(ctx dbos.DBOSContext, in AgentRunInput, approvalID, stepName string) (string, error) {
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
func (e *DBOSExecutor) fail(ctx dbos.DBOSContext, in AgentRunInput, runID string, res AgentRunResult, what string, err error) (AgentRunResult, error) {
	res.Status = RunStatusFailed
	if e.store != nil {
		_, _ = dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
			return true, e.store.RecordTransition(stepCtx, RunTransition{
				InstanceID: in.InstanceID,
				RunID:      runID,
				Status:     RunStatusFailed,
				EventType:  EventTypeFailed,
				Error:      err.Error(),
			})
		}, dbos.WithStepName("emit:"+EventTypeFailed))
	}
	return res, fmt.Errorf("act: %s: %w", what, err)
}

// interruptedAwaitingApproval wraps a cancelled-context error from the approval wait. A cancelled context there means
// the process is shutting down (a deploy) or the run was cancelled (CancelWorkflow), NOT that the run failed: so
// unlike fail it records no failed transition, leaving the store row truthfully waiting_approval. That keeps the run
// cancellable to clean it up, and a Cancel that raced the shutdown still resolves cleanly.
//
// Known limitation (Fase 1): on shutdown DBOS v0.19 still terminalizes the interrupted workflow (it records a terminal
// status via an uncancellable write, and recovery only re-runs PENDING workflows), so the parked approval does not
// survive the restart on its own. Making it resume across a deploy needs the split worker topology; this only stops the
// interrupt from masquerading as a failure in the meantime.
func interruptedAwaitingApproval(err error) error {
	return fmt.Errorf("act: approval wait interrupted (context cancelled): %w", err)
}

// setupApproval records the waiting_approval transition and persists the approval request atomically, in one durable
// step. The approval identity and the args hash are derived deterministically from the run and proposal, so a replay
// creates the same approval (idempotent) rather than a second one.
func (e *DBOSExecutor) setupApproval(ctx dbos.DBOSContext, in AgentRunInput, runID, proposal string) error {
	if e.store == nil {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		if err := e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusWaitingApproval,
			EventType:  EventTypeWaitingApproval,
		}); err != nil {
			return false, err
		}
		return true, e.store.CreateApproval(stepCtx, NewApproval{
			ApprovalID:  ApprovalIDForRun(runID),
			RunID:       runID,
			InstanceID:  in.InstanceID,
			ToolName:    proposedActionTool,
			ArgsHash:    HashArgs(proposal),
			Proposal:    proposal,
			Policy:      "approval_required",
			RequestedBy: in.Actor.Subject,
			ExpiresOn:   e.approvalDeadline(),
		})
	}, dbos.WithStepName("setup_approval"))
	if err != nil {
		return fmt.Errorf("act: setup approval: %w", err)
	}
	return nil
}

// approvalRecvTimeout is the duration a run's approval wait passes to dbos.Recv. With no deadline configured (the
// default, approvalTimeout <= 0) it returns approvalWaitForever, so the run waits until a decision arrives or it is
// cancelled rather than auto-expiring; a positive ApprovalTimeout opts into a real deadline.
func (e *DBOSExecutor) approvalRecvTimeout() time.Duration {
	if e.approvalTimeout <= 0 {
		return approvalWaitForever
	}
	return e.approvalTimeout
}

// approvalDeadline is the wall-clock expiry stored on a pending approval, or nil when no deadline is configured (the
// default). A nil deadline is what makes the store never reject a late decision and the UI never mark an approval
// expired: expiry only exists when a deployment opts into it with a positive ApprovalTimeout.
func (e *DBOSExecutor) approvalDeadline() *time.Time {
	if e.approvalTimeout <= 0 {
		return nil
	}
	t := time.Now().Add(e.approvalTimeout)
	return &t
}

// recordExpiry records the expired transition and marks the given pending approval expired, in one durable step.
// Marking the approval tolerates it already being resolved (a decision that raced the timeout) or absent, since the run
// is ending regardless. approvalID identifies the action that timed out (per segment on the gateway path).
func (e *DBOSExecutor) recordExpiry(ctx dbos.DBOSContext, in AgentRunInput, runID, approvalID string) error {
	if e.store == nil {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		if err := e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusExpired,
			EventType:  EventTypeExpired,
		}); err != nil {
			return false, err
		}
		_, err := e.store.ResolveApproval(stepCtx, in.InstanceID, approvalID, ApprovalStatusExpired, "")
		if err != nil && !errors.Is(err, ErrApprovalNotResolvable) && !errors.Is(err, ErrApprovalNotFound) {
			return false, err
		}
		return true, nil
	}, dbos.WithStepName("emit:"+EventTypeExpired))
	if err != nil {
		return fmt.Errorf("act: emit expired: %w", err)
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
// AI session it ran in, and the write it paused on (if any). Bundling them in one named struct means DBOS checkpoints
// them together, so a replay yields the same session ID and proposal rather than recomputing them against a freshly
// re-driven loop.
type runSegmentStepResult struct {
	Response  string
	SessionID string
	// Proposal is the write action the segment captured (nil for a segment that produced a final answer). It is
	// checkpointed alongside the response, so a replay after a crash yields the same proposal instead of re-driving the
	// loop.
	Proposal *ai.ProposedAction
}

// runSegmentedLoop drives a run as a segmented durable loop (b2). Each iteration runs one segment of the model/tool
// loop as a durable step; when a segment pauses on a governed write, the action flows through the Tool Gateway (policy,
// action ledger, human approval or auto-approve, idempotent execution) and its result is injected into the next
// segment, so the model sees the outcome of the action it proposed and either proposes another or composes a final
// answer. The body branches only on checkpointed values, so a crash resumes at the last completed step, and — because
// the gateway's Execute is idempotent on the ledger — an approved write is never driven twice under at-least-once
// delivery (§12). A crash mid-segment re-drives that segment's model turns (the accepted spike granularity: replaying
// a segment re-flushes its trace but never re-executes a committed external write).
func (e *DBOSExecutor) runSegmentedLoop(ctx dbos.DBOSContext, in AgentRunInput, runID string, snapshot *ai.AgentSnapshot, res AgentRunResult) (AgentRunResult, error) {
	var sessionID string
	var resume *ai.InjectedResult
	for seg := 0; ; seg++ {
		// Step: run one segment. Fresh on seg 0 (seeded by the prompt); on resume it reopens the same session and injects
		// the prior action's result. The result is a named struct so DBOS checkpoints the response, session ID and
		// proposal together: a replay yields the same session and proposal instead of re-driving the loop.
		out, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (runSegmentStepResult, error) {
			r, runErr := e.runner.RunSegment(stepCtx, RunSegmentInput{
				InstanceID: in.InstanceID, Claims: in.Actor.Claims, Snapshot: snapshot,
				Prompt: in.Prompt, SessionID: sessionID, Resume: resume,
			})
			return runSegmentStepResult{Response: r.Response, SessionID: r.SessionID, Proposal: r.Proposed}, runErr
		}, dbos.WithStepName(fmt.Sprintf("%s_%d", stepRunAgent, seg)))
		if err != nil {
			return e.fail(ctx, in, runID, res, "run segment", err)
		}
		res.Response = out.Response

		// Bind the run to the session the first segment opened, so the API/UI render it as a chat (§13.1) and every later
		// segment reopens the same conversation. The session ID is checkpointed, so a replay threads the same value.
		if seg == 0 {
			sessionID = out.SessionID
			if err := e.persistConversation(ctx, in, runID, sessionID); err != nil {
				return e.fail(ctx, in, runID, res, "persist conversation", err)
			}
		}

		// No governed write pending: the segment produced a final answer, so the run is done and succeeds with it.
		if out.Proposal == nil {
			res.Status = RunStatusSucceeded
			if err := e.emit(ctx, in, runID, RunStatusSucceeded, EventTypeSucceeded); err != nil {
				return res, err
			}
			return res, nil
		}

		// Govern and execute the proposed write. On success it returns the redacted result to inject into the next
		// segment; on any terminal outcome (denied, rejected, expired, indeterminate, failed) it sets res and ends the run.
		inject, terminal, err := e.governProposedAction(ctx, in, runID, snapshot, &res, out.Proposal, seg)
		if terminal || err != nil {
			return res, err
		}
		resume = inject
	}
}

// governProposedAction governs a single proposed write for segment seg: it derives the concrete tool proposal, applies
// deterministic policy, waits for a human approval (or auto-approves per the connector's snapshot posture), then
// executes the write idempotently through the gateway. On a confirmed write it returns the redacted result to inject
// into the next segment (terminal=false). On any terminal outcome (denied, rejected, expired, indeterminate, failed) it
// sets *res and returns terminal=true. A non-nil error is a workflow error (already recorded via fail); the caller
// returns *res with it. Its steps are named per segment so a run proposing several actions keeps stable, unambiguous
// step and approval identities across replay.
func (e *DBOSExecutor) governProposedAction(ctx dbos.DBOSContext, in AgentRunInput, runID string, snapshot *ai.AgentSnapshot, res *AgentRunResult, captured *ai.ProposedAction, seg int) (*ai.InjectedResult, bool, error) {
	approvalID := ApprovalIDForSegment(runID, seg)

	// Step: derive the concrete tool proposal from the captured action.
	proposed, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (actionProposalStep, error) {
		p, ok, perr := e.proposer.Propose(stepCtx, ProposeInput{
			InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Actor: in.Actor,
			Snapshot: snapshot, Response: res.Response, Captured: captured,
		})
		return actionProposalStep{Proposal: p, OK: ok}, perr
	}, dbos.WithStepName(fmt.Sprintf("propose_action_%d", seg)))
	if err != nil {
		*res, err = e.fail(ctx, in, runID, *res, "propose action", err)
		return nil, true, err
	}
	if !proposed.OK {
		// Defensive: the caller only governs a non-nil captured proposal, and the captured proposer surfaces it, so this
		// is not expected. Treat a declined proposal as no pending action and succeed with the response.
		res.Status = RunStatusSucceeded
		if err := e.emit(ctx, in, runID, RunStatusSucceeded, EventTypeSucceeded); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	}
	proposal := proposed.Proposal

	// Resolve the connector captured in the run's immutable snapshot once. It decides the approval posture below (an
	// auto-approved action skips the human wait, while the gateway still records the ledger and executes idempotently),
	// and it is bound onto the execute/verify step context further down so the governed write runs against exactly the
	// connector config the approver reviewed (§8.3) — an admin editing the connector (URL, allowed_hosts, approval)
	// during the approval wait cannot redirect the approved write. This is deterministic (snapshot + static globs), so
	// it is safe to feed the checkpointed propose step below.
	capturedConn, hasCapturedConn := mcpConnector(snapshot, proposal.Connector)
	autoApprove := map[string]bool{}
	if hasCapturedConn && connectorAutoApproves(capturedConn, rawToolName(proposal.Tool, proposal.Connector)) {
		autoApprove[proposal.Tool] = true
	}

	// Step: authorize the proposal. gateway.Propose validates the tool and arguments, evaluates deterministic policy
	// (auto-approve applies here), and records the action ledger (proposed, then the policy decision).
	auth, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (Authorization, error) {
		return e.gateway.Propose(stepCtx, ProposeActionInput{
			InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Actor: in.Actor,
			Proposal: proposal, AutoApprove: autoApprove,
		})
	}, dbos.WithStepName(fmt.Sprintf("gateway_propose_%d", seg)))
	if err != nil {
		*res, err = e.fail(ctx, in, runID, *res, "authorize action", err)
		return nil, true, err
	}

	switch auth.Decision {
	case PolicyDeny:
		// The deterministic policy refused the action: no external effect runs. The ledger already recorded
		// policy_rejected; the run ends failed with the sanitized reason.
		res.Status = RunStatusFailed
		if err := e.emitFailedReason(ctx, in, runID, "action denied by policy: "+auth.Reason); err != nil {
			return nil, true, err
		}
		return nil, true, nil

	case PolicyApprovalRequired:
		// Persist the exact proposal as an approval (keyed per segment so several actions in a run stay distinct) and
		// suspend on a durable wait.
		if err := e.setupActionApproval(ctx, in, runID, approvalID, auth, proposal); err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "await approval setup", err)
			return nil, true, err
		}
		decision, err := dbos.Recv[string](ctx, approvalTopic, e.approvalRecvTimeout())
		if err != nil {
			if errors.Is(err, &dbos.DBOSError{Code: dbos.TimeoutError}) {
				res.Status = RunStatusExpired
				if err := e.recordExpiry(ctx, in, runID, approvalID); err != nil {
					return nil, true, err
				}
				return nil, true, nil
			}
			if errors.Is(err, context.Canceled) {
				return nil, true, interruptedAwaitingApproval(err)
			}
			*res, err = e.fail(ctx, in, runID, *res, "await approval", err)
			return nil, true, err
		}
		if ApprovalDecision(decision) != ApprovalApproved {
			// A human denied the action; it never advances past approval_pending in the ledger and no write runs. The
			// rejected event carries who denied it, so the timeline attributes the stop as it attributes an approval.
			res.Status = RunStatusRejected
			denier, err := e.loadDecider(ctx, in, approvalID, fmt.Sprintf("%s_%d", stepLoadDenier, seg))
			if err != nil {
				*res, err = e.fail(ctx, in, runID, *res, "load denier", err)
				return nil, true, err
			}
			if err := e.emitWithPayload(ctx, in, runID, RunStatusRejected, EventTypeRejected, deciderPayload(denier)); err != nil {
				return nil, true, err
			}
			return nil, true, nil
		}
		// Step: reload the resolved approval to record the REAL approver on the ledger and on the resumed event.
		decidedBy, err := e.loadDecider(ctx, in, approvalID, fmt.Sprintf("%s_%d", stepLoadApprover, seg))
		if err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "load approver", err)
			return nil, true, err
		}

		// Step: record the approval on the ledger, bound to the exact args hash the human saw (§11.2). A hash mismatch
		// fails closed here rather than executing arguments that were not approved.
		if _, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
			return true, e.gateway.RecordApproved(stepCtx, ApprovedInput{
				InstanceID: in.InstanceID, RunID: runID, ToolCallID: proposal.ToolCallID,
				ArgsHash: auth.ArgsHash, DecidedBy: decidedBy,
			})
		}, dbos.WithStepName(fmt.Sprintf("gateway_record_approved_%d", seg))); err != nil {
			*res, err = e.fail(ctx, in, runID, *res, "record approval", err)
			return nil, true, err
		}
		if err := e.emitWithPayload(ctx, in, runID, RunStatusRunning, EventTypeResumed, deciderPayload(decidedBy)); err != nil {
			return nil, true, err
		}

	case PolicyAllow:
		// Auto-approved: the ledger already recorded approved; proceed straight to the write.

	default:
		// Fail closed on any unrecognized decision.
		res.Status = RunStatusFailed
		if err := e.emitFailedReason(ctx, in, runID, "action blocked: undecidable policy"); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	}

	// Step: perform the external write through the gateway. It is idempotent on the action ledger, so a replay after a
	// crash returns the recorded outcome instead of re-executing. The business outcome is in the report, not an error.
	report, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (ExecuteReport, error) {
		// Bind the instance and the snapshot's captured connector onto the context: the gateway's shared SecretResolver
		// resolves the connector credential against this run's project, and the runtime-backed connector resolver freezes
		// the connector's target/allowed_hosts/approval to what was reviewed (§8.3), so a connector edited during the
		// approval wait cannot redirect the approved write. The bearer secret is still resolved live.
		return e.gateway.Execute(withGovernedConnector(stepCtx, in.InstanceID, capturedConn, hasCapturedConn), ExecuteActionInput{
			InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Proposal: proposal, TraceID: runID,
			// Reauthorization immediately before the write re-evaluates policy with the SAME AutoApprove the proposal was
			// authorized under: without it, an auto-approved unclassified tool (every generic MCP tool is ClassUnknown) is
			// recomputed as approval_required and refused as policy_rejected, so the connector's approval.auto / auto_approve
			// posture would execute in Propose but be rejected here. See TestGatewayReauthorizePreservesAutoApprove.
			AutoApprove: autoApprove,
		})
	}, dbos.WithStepName(fmt.Sprintf("gateway_execute_%d", seg)))
	if err != nil {
		*res, err = e.fail(ctx, in, runID, *res, "apply action", err)
		return nil, true, err
	}

	switch report.Outcome {
	case OutcomeSucceeded:
		// The write is confirmed. Unlike the old linear path, a successful action does NOT end the run: the loop resumes
		// so the model sees the result and closes (or proposes another action). res records the latest action taken; the
		// succeeded transition is emitted when the model finally produces a final answer.
		res.ActionTaken = true
		res.ActionRef = report.ExternalReference
		if auth.Verifiable {
			// Best-effort verification (§16.1): a verify failure never undoes a confirmed write, so its step outcome does
			// not affect the run's success. But the error is not discarded: gateway.Verify returns nil for a connector
			// non-confirmation (expected, best-effort) and an error only for a failure to PERSIST the verification — an
			// audit-integrity signal. Surface that (sanitized, no secrets — the gateway already redacts) rather than
			// masking it under a successful run.
			if _, verifyErr := dbos.RunAsStep(ctx, func(stepCtx context.Context) (VerifyResult, error) {
				// Bind the snapshot's connector here too: Verify's redaction set is built from a connector Lookup, so it
				// must resolve the frozen connector credential, not a live-edited one.
				return e.gateway.Verify(withGovernedConnector(stepCtx, in.InstanceID, capturedConn, hasCapturedConn), ExecuteActionInput{
					InstanceID: in.InstanceID, RunID: runID, AgentName: in.AgentName, Proposal: proposal, TraceID: runID,
				})
			}, dbos.WithStepName(fmt.Sprintf("gateway_verify_%d", seg))); verifyErr != nil && e.logger != nil {
				e.logger.Warn("act: verification persistence failed; write is confirmed but audit record is incomplete",
					"run", runID, "err", verifyErr)
			}
		}
		// Inject the redacted result so the next segment's model reacts to the outcome. Prefer the tool's text result;
		// fall back to the external reference (e.g. a ticket key) when the tool returned only a handle.
		message := report.Message
		if message == "" {
			message = report.ExternalReference
		}
		return &ai.InjectedResult{ToolCallID: proposal.ToolCallID, Tool: proposal.Tool, Message: message}, false, nil

	case OutcomeIndeterminate:
		// The external result could not be confirmed. The ledger holds indeterminate; the run cannot succeed and is
		// not retried automatically (§12). It ends failed with a clear marker for a human.
		res.Status = RunStatusFailed
		if err := e.emitFailedReason(ctx, in, runID, "action result indeterminate: awaiting human resolution"); err != nil {
			return nil, true, err
		}
		return nil, true, nil

	default: // OutcomeFailed
		res.Status = RunStatusFailed
		if err := e.emitFailedReason(ctx, in, runID, "action failed"); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	}
}

// setupActionApproval records the waiting_approval transition and persists the gateway's concrete proposal as an
// approval request, in one durable step. Unlike setupApproval (the Fase 1 simulated variant), the approval carries the
// real tool, connector, tool call and canonical args hash from the authorization, so the inbox binds to exactly what
// will execute (§11.2). Idempotent: a replay creates the same approval rather than a second one.
func (e *DBOSExecutor) setupActionApproval(ctx dbos.DBOSContext, in AgentRunInput, runID, approvalID string, auth Authorization, proposal ToolProposal) error {
	if e.store == nil {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		if err := e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusWaitingApproval,
			EventType:  EventTypeWaitingApproval,
		}); err != nil {
			return false, err
		}
		return true, e.store.CreateApproval(stepCtx, NewApproval{
			ApprovalID:  approvalID,
			RunID:       runID,
			InstanceID:  in.InstanceID,
			ToolName:    proposal.Tool,
			Connector:   proposal.Connector,
			ToolCallID:  proposal.ToolCallID,
			ArgsHash:    auth.ArgsHash,
			Proposal:    proposal.Summary,
			Policy:      string(auth.Decision),
			RequestedBy: in.Actor.Subject,
			ExpiresOn:   e.approvalDeadline(),
		})
	}, dbos.WithStepName("setup_approval"))
	if err != nil {
		return fmt.Errorf("act: setup action approval: %w", err)
	}
	return nil
}

// emitFailedReason records a failed run transition carrying a sanitized reason, as a durable step, and returns nil on
// success. Unlike fail(), it is for a business-terminal failure (policy denied, action failed/indeterminate) that
// ends the run without a workflow error, so DBOS does not retry it.
func (e *DBOSExecutor) emitFailedReason(ctx dbos.DBOSContext, in AgentRunInput, runID, reason string) error {
	if e.store == nil {
		return nil
	}
	_, err := dbos.RunAsStep(ctx, func(stepCtx context.Context) (bool, error) {
		return true, e.store.RecordTransition(stepCtx, RunTransition{
			InstanceID: in.InstanceID,
			RunID:      runID,
			Status:     RunStatusFailed,
			EventType:  EventTypeFailed,
			Error:      reason,
		})
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
// historical ":1" suffix, so existing single-action runs, approvals and tests are unaffected.
func ApprovalIDForSegment(runID string, seg int) string { return runID + ":" + strconv.Itoa(seg+1) }

// HashArgs returns the canonical hash an approval is bound to: the decision is valid only for these exact arguments
// (§6.4, §11.2). The spike's simulated path hashes the proposal text; the gateway path hashes the canonicalized tool
// arguments (HashCanonicalArgs). Both yield the same "sha256:"-prefixed form so a consumer cannot tell them apart.
func HashArgs(args string) string {
	sum := sha256.Sum256([]byte(args))
	return "sha256:" + hex.EncodeToString(sum[:])
}
