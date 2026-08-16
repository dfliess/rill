// Package act is the Kairos Act spike: a durable executor for declarative agents (issue kairosagentica/kairos-cloud#123).
//
// The Agent resource (proto AgentSpec, parser, validating reconciler) describes an agent declaratively; runtime/ai
// executes a snapshot of it fail-closed. This package is the missing third piece: it drives that execution as a
// durable workflow so a run survives process crashes, deploys and long human-approval waits.
//
// The seam is AgentExecutor. Everything durable lives behind it (DBOS Transact for the spike), so the rest of the
// runtime, the Agent API and the UI never depend on the workflow engine and the choice stays reversible (§19.2 of
// docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md).
package act

import (
	"context"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
)

// AgentRunInput is the immutable request to start an agent run.
//
// IdempotencyKey is the caller's stable identity for this run: an alert firing once, a form submitted once. It is
// combined with InstanceID and AgentName to form the durable workflow ID, which is what makes Start idempotent
// (OAOO): a second Start with the same key attaches to the same run rather than starting a new one. The design
// derives the key from (instance, resource, execution_time); the spike takes it verbatim from the caller.
type AgentRunInput struct {
	InstanceID     string
	AgentName      string
	Prompt         string
	IdempotencyKey string
	// Trigger records what started the run: "manual", "alert", etc. (§9.5). Empty defaults to "manual", the only
	// trigger the spike starts. It is carried into the run record for audit and filtering.
	Trigger string
	// TriggerRef names the source resource that fired an automatic trigger (e.g. the alert's resource name), so the
	// run's provenance links back to it. Empty for manual and schedule triggers.
	TriggerRef string
	// ConversationID links the run to a Rill AI conversation when the run continues one (§13.1). Optional.
	ConversationID string
	// Actor records who authorized this run, including the security identity the run executes with. It is
	// checkpointed with the run (as part of the workflow input), so the durable, auditable record — and the tool
	// authority — are bound to the initiator rather than to whichever session the worker happens to open.
	Actor Actor
}

// Actor identifies who authorized a run: a user for a manual run, or a declared service principal for an automatic
// trigger (§17.3 of the design).
//
// The run executes with Claims: the fail-closed tool intersection (ai.DynamicAgent.allowedTools) resolves
// CheckAccess against them, so a run can never exceed the permissions of whoever (or whatever) started it. A manual
// run carries the initiating user's claims.
//
// An automatic trigger runs under a governed identity declared the alert way: the trigger names a user
// (AgentTriggerActor.user_id/user_email, resolved to attributes via the admin service) or gives explicit
// attributes, and the dispatcher builds Claims from those (§17.3). A trigger that declares no identity yields nil
// Claims, and the worker fails closed (no privileged fallback), so an automatic run is safe-but-inert until an
// identity is declared.
type Actor struct {
	// Subject is the stable identifier of the initiator: a user subject, or a service principal name.
	Subject string
	// ServicePrincipal is true when an automatic trigger started the run under a declared service identity rather
	// than a human. Manual runs leave it false. An automatic run must never implicitly inherit its author's rights.
	ServicePrincipal bool
	// Claims is the initiator's security identity, sufficient to rebuild the SecurityClaims the run's Session is
	// opened with. It is JSON-serializable (runtime.SecurityClaims has custom (Un)MarshalJSON), so it is
	// checkpointed with the run and reconstructed on recovery, and the tool intersection is resolved against it: a
	// restricted user cannot read metrics or take an action through a run that the same user could not perform
	// directly. Nil only for a run started without an initiator identity (e.g. a not-yet-wired automatic trigger).
	Claims *runtime.SecurityClaims
}

// Trigger sources that may start a run (§9.5). A run's recorded trigger must be one of these; the manual-start API
// forces TriggerManual so a caller cannot forge automatic provenance, and the executor rejects any other unknown
// value so an internal caller cannot either.
const (
	TriggerManual   = "manual"
	TriggerAlert    = "alert"
	TriggerReport   = "report"
	TriggerSchedule = "schedule"
	TriggerWebhook  = "webhook"
	TriggerAPI      = "api"
)

var validTriggers = map[string]bool{
	TriggerManual:   true,
	TriggerAlert:    true,
	TriggerReport:   true,
	TriggerSchedule: true,
	TriggerWebhook:  true,
	TriggerAPI:      true,
}

// ApprovalDecision is a human's ruling on a run that is suspended waiting for approval.
type ApprovalDecision string

const (
	ApprovalApproved ApprovalDecision = "approved"
	ApprovalRejected ApprovalDecision = "rejected"
)

// RunStatus is the terminal (or current) state of a run (§11.1).
type RunStatus string

const (
	// RunStatusQueued is the initial state: the run has been created and enqueued but the worker has not started it.
	RunStatusQueued RunStatus = "queued"
	// RunStatusRunning is an actively executing run (model/tool loop).
	RunStatusRunning RunStatus = "running"
	// RunStatusWaitingApproval is a run suspended on a durable wait for a human approval decision.
	RunStatusWaitingApproval RunStatus = "waiting_approval"
	// RunStatusSucceeded is a run that finished its work (including any approved action) without error.
	RunStatusSucceeded RunStatus = "succeeded"
	// RunStatusFailed is a run that ended with an error.
	RunStatusFailed RunStatus = "failed"
	// RunStatusRejected is a run a human denied at the approval gate; no external effect ran.
	RunStatusRejected RunStatus = "rejected"
	// RunStatusCancelled is a run stopped by an operator.
	RunStatusCancelled RunStatus = "cancelled"
)

// IsTerminal reports whether s is a final state a run cannot leave. A streaming consumer stops once the run reaches
// one; the API never transitions a run out of one.
func (s RunStatus) IsTerminal() bool {
	switch s {
	case RunStatusSucceeded, RunStatusFailed, RunStatusRejected, RunStatusCancelled:
		return true
	default:
		return false
	}
}

// AgentExecutor starts, resumes and cancels durable agent runs. It is the reversibility seam: DBOS is entirely
// behind this interface, so swapping the durable backend (e.g. to Temporal) changes only the implementation and
// the worker deployment, never the AgentSpec, the APIs or the UI.
type AgentExecutor interface {
	// Start launches a run for in.IdempotencyKey and returns its run ID. Calling it again with the same
	// IdempotencyKey is a no-op that returns the same run ID (idempotent by construction).
	Start(ctx context.Context, in AgentRunInput) (runID string, err error)
	// Resume delivers a human approval decision to a run that is waiting for one. toolCallID identifies the specific
	// action being decided: when non-empty, the decision is routed to the per-action DBOS topic so the executor
	// matches it to the right governed call; when empty, it falls back to the legacy single-topic path (Fase 1 runs
	// and old segment-based approvals).
	Resume(ctx context.Context, runID string, decision ApprovalDecision, toolCallID string) error
	// Resumable reports whether a run's durable execution is still able to act on a decision, as opposed to having
	// been given up on. It exists because delivering a decision cannot fail loudly: the durable message is simply
	// filed for a run that will never read it, so an approval recorded against an abandoned run tells its approver
	// the action was authorized while nothing will ever perform it. Callers ask BEFORE recording the decision, and a
	// false answer means refuse to record it at all (kairos-cloud#135).
	//
	// It answers about the execution, not the approval: whether the approval is still open is the store's question.
	// An error means the executor could not tell, which callers treat as a refusal rather than a licence to sign.
	Resumable(ctx context.Context, runID string) (bool, error)
	// Cancel stops a run in the given instance. instanceID scopes the cancellation so a caller can only stop a run in
	// the instance it addresses, and lets the executor record the terminal state against the right tenant.
	Cancel(ctx context.Context, instanceID, runID string) error
}

// AgentRunResult is the durable outcome of a run. DBOS checkpoints it as the workflow result, so a Start replaying
// a finished run returns this without re-executing anything.
type AgentRunResult struct {
	RunID       string    `json:"run_id"`
	AgentName   string    `json:"agent_name"`
	Status      RunStatus `json:"status"`
	Response    string    `json:"response"`     // the agent's final model response
	ActionTaken bool      `json:"action_taken"` // whether the (simulated) external write ran
	ActionRef   string    `json:"action_ref"`   // observable reference produced by the action, empty if none
}

// ActionRequest is what the workflow hands to the action step after approval. In the spike the action is simulated;
// in production this is the input to the tool gateway (§16).
type ActionRequest struct {
	RunID     string
	AgentName string
	// Actor is the identity that authorized the run, carried through so the action ledger can record who the
	// external write was performed on behalf of (§18.3). It is the same value checkpointed on AgentRunInput.
	Actor Actor
	// SnapshotMark is a value derived from the immutable snapshot captured at run start. It exists so a test can
	// prove the action used the original snapshot even if the agent definition changed mid-run.
	SnapshotMark string
	// Response is the agent's proposal that motivated the action.
	Response string
}

// ActionResult is the outcome of the (simulated) external write.
type ActionResult struct {
	Ref string // e.g. an external ticket reference
}

// RunSegmentInput is the immutable request for one segment of a run's model/tool loop (b2). A run advances one
// segment at a time: each segment runs the model until it either produces a final answer or pauses on a governed
// write. Fresh on the first segment (SessionID empty, seeded by Prompt); on resume it reopens the same session and
// injects the prior action's result.
type RunSegmentInput struct {
	InstanceID string
	Claims     *runtime.SecurityClaims
	// Snapshot is the immutable agent definition captured at run start; every segment runs from the same one.
	Snapshot *ai.AgentSnapshot
	// Prompt seeds the conversation on the first segment. Ignored on resume (the conversation is reconstructed from
	// the reopened session's message tree).
	Prompt string
	// SessionID is the AI session to continue: empty opens a fresh session (first segment) and the caller reads the
	// opened ID back from RunSegmentResult; non-empty reopens that session so the resumed segment reconstructs the
	// conversation from its persisted message tree, under the same checkpointed claims.
	SessionID string
	// Resume carries the executed actions' results to inject as the turns the model reacts to before it runs again.
	// Nil on the first segment (there are no prior actions).
	Resume []*ai.InjectedResult
}

// RunSegmentResult is the checkpointed outcome of one segment.
type RunSegmentResult struct {
	// SessionID is the AI session the segment ran in. On the first segment it is the freshly opened session the
	// executor persists as the run's conversation and threads into subsequent segments; it is stable thereafter.
	SessionID string
	// Response is the segment's final model text (a wrap-up when the segment paused on a write, the closing answer
	// when it did not).
	Response string
	// Proposed is the list of write actions the segment paused on, empty when the segment produced a final answer and
	// the run is done. A non-empty slice is the pause point: the workflow governs and executes all of them in order.
	Proposed []*ai.ProposedAction
}

// Runner performs the non-deterministic, side-effecting work of a run. The durable workflow invokes each method
// inside a DBOS step so its result is checkpointed and replayed on recovery instead of re-executed.
//
// Splitting this out keeps the workflow body deterministic (it only branches on already-checkpointed values) and
// lets tests substitute a scripted LLM and an observable action sink for the real runtime wiring.
type Runner interface {
	// LoadSnapshot resolves the agent's immutable definition. It runs once, at the start of a run: the snapshot is
	// then fixed for the life of that run (§8.3), so editing the agent afterwards does not affect it.
	LoadSnapshot(ctx context.Context, instanceID, agentName string) (*ai.AgentSnapshot, error)
	// RunSegment executes ONE segment of the agent's model/tool loop from a fixed snapshot, under the initiating
	// actor's security claims, and returns its final response, the AI session it ran in, and the write it paused on (if
	// any). The claims bound the callable tool set: a tool the initiator cannot access is dropped from the run
	// (fail-closed intersection, §17.3), so a run never acts beyond its initiator's authority regardless of what the
	// snapshot declares.
	//
	// A fresh segment (in.SessionID empty) opens a new session seeded by the prompt; a resumed segment (in.SessionID
	// set, in.Resume carrying the prior action's result) reopens the same session and reconstructs the conversation
	// from its persisted message tree. RunSegmentResult.Proposed is the write the segment paused on (captured, not
	// executed), or nil when the segment produced a final answer: the executor governs a non-nil proposal through the
	// tool gateway and injects the result into the next segment, while a nil one ends the run with just its response.
	RunSegment(ctx context.Context, in RunSegmentInput) (RunSegmentResult, error)
	// ApplyAction performs the (simulated) external write. It must be safe under at-least-once delivery: a DBOS
	// step whose process dies after the side effect but before the checkpoint is re-run on recovery (§12).
	ApplyAction(ctx context.Context, req ActionRequest) (ActionResult, error)
}
