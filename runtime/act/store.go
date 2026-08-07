package act

import (
	"context"
	"errors"
	"time"
)

// ErrRunNotFound is returned when a run (or an event/approval addressed through one) does not exist for the
// requested instance. It is distinct from a store/transport error so the API can map it to NotFound rather than
// Internal, and so scoping failures (a run that exists under another instance) look identical to a missing run.
var ErrRunNotFound = errors.New("act: run not found")

// ErrApprovalNotFound is returned when an approval does not exist for the requested instance.
var ErrApprovalNotFound = errors.New("act: approval not found")

// Event types emitted onto agent_run_events as a run moves through its lifecycle. They are stable strings (not an
// enum) so the store, the proto API and the frontend agree on the wire value without a shared code generator, and
// so a new type can be added without a migration. The lifecycle edges (queued, running, terminal states) fire at
// most once per run; the approval edges (waiting_approval, resumed, rejected) fire once per governed
// action, so the segmented loop repeats them. Idempotency is keyed by RunTransition.DedupeKey (see
// PostgresRunStore.RecordTransition): the executor scopes the approval edges per segment and leaves the rest keyed
// by their type alone.
const (
	EventTypeQueued          = "run.queued"
	EventTypeRunning         = "run.running"
	EventTypeWaitingApproval = "run.waiting_approval"
	EventTypeResumed         = "run.resumed"
	EventTypeSucceeded       = "run.succeeded"
	EventTypeFailed          = "run.failed"
	EventTypeRejected        = "run.rejected"
	EventTypeCancelled       = "run.cancelled"
)

// Event visibility tiers (§15.1). A run event is addressed to an audience: the initiating user, a platform
// operator, or an auditor. Phase 1A writes everything as VisibilityUser; the column exists so later redaction can
// widen or narrow an event's audience without reshaping the table.
const (
	VisibilityUser     = "user"
	VisibilityOperator = "operator"
	VisibilityAuditor  = "auditor"
)

// Approval states (§11.2, projected to the approval record). Pending is the only non-terminal state; the API and
// the workflow both resolve it, so the resolution is guarded to a single winner.
const (
	ApprovalStatusPending  = "pending"
	ApprovalStatusApproved = "approved"
	ApprovalStatusDenied   = "denied"
	// ApprovalStatusCancelled marks an approval withdrawn because its run reached a terminal state without deciding it
	// (cancelled, failed, or a sibling action that terminated the run). It is terminal like approved/denied, so the
	// inbox (which shows only pending approvals) stops offering a decision on a run that will never resume.
	ApprovalStatusCancelled = "cancelled"
)

// NewRun is the immutable identity and provenance of a run, written once when the run is created (status queued).
// Everything here is fixed for the life of the run; mutable lifecycle state (status, timestamps, error) is written
// separately by RecordTransition so the create path and the transition path never contend on the same columns.
type NewRun struct {
	// RunID is the durable workflow ID (ComposeRunID output). It is the primary key: creating the same run twice
	// (a replayed Start) is a no-op, which is what keeps Start idempotent end-to-end, not just inside DBOS.
	RunID      string
	InstanceID string
	// OrganizationID and ProjectID are the outer tenancy scope (§15.2). They are nullable in Phase 1A because the
	// runtime addresses tenants by instance; they are carried now so a later identity model can backfill and index
	// on them without a schema change.
	OrganizationID string
	ProjectID      string
	AgentName      string
	// SpecHash binds the run to the exact agent definition it started from (§8.3). Empty until the snapshot step.
	SpecHash string
	// Trigger records what started the run: "manual", "alert", etc. (§9.5). The spike only starts manual runs.
	Trigger string
	// TriggerRef names the source resource that fired an automatic trigger (e.g. the alert's resource name). Empty for
	// manual and schedule triggers; carried so the run's provenance links back to the resource that started it.
	TriggerRef string
	// IdempotencyKey is the caller's stable key for this run, retained for audit even though RunID already encodes it.
	IdempotencyKey string
	// ConversationID links the run to a Rill AI conversation, when one exists (§13.1). Nullable.
	ConversationID string
	Actor          Actor
}

// Run is the stored, queryable state of a run: its immutable identity plus its current lifecycle state. It is the
// API's source of truth, populated entirely from Postgres, so a reader never touches the workflow engine (§7.2).
type Run struct {
	RunID          string
	InstanceID     string
	OrganizationID string
	ProjectID      string
	AgentName      string
	SpecHash       string
	Trigger        string
	TriggerRef     string
	IdempotencyKey string
	ConversationID string
	Actor          Actor
	Status         RunStatus
	// Error is a sanitized failure message, set only for a failed run. It never carries secrets or raw provider
	// output (§17.2); the executor is responsible for redacting before it reaches the store.
	Error      string
	CreatedOn  time.Time
	UpdatedOn  time.Time
	StartedOn  *time.Time // first transition into running; nil while queued
	FinishedOn *time.Time // transition into a terminal status; nil while non-terminal
}

// RunTransition atomically moves a run to a new status and records the matching event. The executor emits one per
// lifecycle edge; the store applies the status update and the event append in a single transaction so a reader can
// never observe a status without its event or vice versa.
type RunTransition struct {
	InstanceID string
	RunID      string
	Status     RunStatus
	// EventType is the event to append for this edge (an EventType* constant).
	EventType string
	// DedupeKey is the idempotency key for this edge within the run: re-applying a transition with a key already
	// recorded is a no-op, which is what stops a mid-step crash replay from moving the status twice. Empty defaults
	// to EventType (one edge per run, the strict historical behavior); the segmented executor scopes the repeatable
	// approval edges as "<event_type>:<segment>" so a run that pauses on several governed actions records each pause.
	DedupeKey string
	// Error is the sanitized error to store on the run, set only when transitioning to a failed status.
	Error string
	// Payload is optional structured detail for the event (redacted). Stored as JSON.
	Payload map[string]any
	// SpecHash, when non-empty, is written onto the run as part of this transition. It lets the running transition
	// (which follows the snapshot step) record which definition the run bound to without a separate write.
	SpecHash string
}

// RunEvent is a stored lifecycle event: one row of agent_run_events. ID is the global, monotonic cursor a stream
// consumer pages by; Seq is the per-run monotonic position for display.
type RunEvent struct {
	ID         int64
	RunID      string
	InstanceID string
	Seq        int64
	EventType  string
	Status     RunStatus
	Payload    map[string]any
	Visibility string
	CreatedOn  time.Time
}

// NewApproval is a human-approval request created when a run suspends at the approval gate (§11.2). It captures the
// exact proposal so a later ApproveAgentApproval can verify the decision applies to the arguments the human saw
// (§6.4): the run is resumed only if the caller-supplied hash matches ArgsHash.
type NewApproval struct {
	ApprovalID string
	RunID      string
	InstanceID string
	// ToolName/Connector/ToolCallID identify the proposed effect. The spike proposes a single simulated action, so
	// these are synthetic; in production they come from the tool gateway's concrete proposal.
	ToolName   string
	Connector  string
	ToolCallID string
	// ArgsHash is the canonical hash of the proposed arguments. The approval is bound to it: changing an argument
	// after approval invalidates the decision and forces a fresh request (§11.2).
	ArgsHash string
	// Proposal is the normalized, human-readable description of what will happen, shown in the approval inbox.
	Proposal string
	// Policy records the deterministic policy decision that routed this to a human (§6.3). Free-form in Phase 1A.
	Policy string
	// RequestedBy is the actor subject on whose behalf the run proposed the action.
	RequestedBy string
	// Position is the 1-based index of this action within its batch (e.g. 1 of 3). Zero when the action is the only
	// one in the turn or when the caller does not track ordering.
	Position int
	// Total is the number of actions in the batch (e.g. 3 when three writes were proposed in one turn).
	Total int
}

// Approval is the stored, queryable state of an approval request.
type Approval struct {
	ApprovalID  string
	RunID       string
	InstanceID  string
	ToolName    string
	Connector   string
	ToolCallID  string
	ArgsHash    string
	Proposal    string
	Policy      string
	Status      string
	RequestedBy string
	DecidedBy   string
	CreatedOn   time.Time
	DecidedOn   *time.Time
	Position    int
	Total       int
}

// ListRunsFilter scopes and narrows a run listing. InstanceID is mandatory: every query is scoped to one instance
// as the first isolation barrier (§15.2), so a caller can never list across tenants.
type ListRunsFilter struct {
	InstanceID string
	AgentName  string    // optional: only runs of this agent
	Status     RunStatus // optional: only runs in this status
	Limit      int       // optional: caps the page size (a store default applies when zero)
}

// ListApprovalsFilter scopes and narrows an approval listing. Like runs, InstanceID is mandatory.
type ListApprovalsFilter struct {
	InstanceID string
	RunID      string // optional: only approvals for this run
	Status     string // optional: only approvals in this status (e.g. "pending" for the inbox)
	Limit      int
}

// RunStore is the durable, queryable home of runs, run events and approvals: the product's state plane, separate
// from the workflow engine's own tables (§7.2, §15). The executor writes to it as runs progress; the API reads from
// it without ever touching DBOS. Every method is scoped by instance so tenant isolation is enforced in the store,
// not left to callers.
type RunStore interface {
	// Migrate creates the store's schema, tables and indexes idempotently. It is safe to call on every startup.
	Migrate(ctx context.Context) error

	// CreateRun inserts a run in the queued state and appends its queued event, atomically. It is idempotent on
	// RunID: a second create for the same run is a no-op (not an error), so a replayed Start does not duplicate.
	CreateRun(ctx context.Context, run NewRun) error
	// RecordTransition applies a status change and appends its event in one transaction. It is idempotent per edge:
	// re-applying an already-recorded transition leaves the row unchanged and appends no duplicate event.
	RecordTransition(ctx context.Context, t RunTransition) error
	// SetRunConversationID links a run to the AI session it executed in, overwriting the conversation_id seeded at
	// create time (normally empty). It is a simple idempotent UPDATE scoped by instance; re-applying it with the same
	// session ID is a no-op, and a run that does not exist for the instance yields ErrRunNotFound.
	SetRunConversationID(ctx context.Context, instanceID, runID, conversationID string) error
	// GetRun returns a run scoped to instanceID, or ErrRunNotFound.
	GetRun(ctx context.Context, instanceID, runID string) (*Run, error)
	// ListRuns returns runs matching the filter, newest first.
	ListRuns(ctx context.Context, f ListRunsFilter) ([]*Run, error)
	// ListRunEvents returns a run's events with ID greater than afterID, oldest first, capped by limit (a store
	// default applies when limit is zero). It is the paging primitive the SSE stream and polling both build on:
	// passing the last seen ID resumes exactly after it.
	ListRunEvents(ctx context.Context, instanceID, runID string, afterID int64, limit int) ([]*RunEvent, error)

	// CreateApproval inserts a pending approval. It is idempotent on ApprovalID so the workflow step that creates it
	// is safe to replay.
	CreateApproval(ctx context.Context, a NewApproval) error
	// GetApproval returns an approval scoped to instanceID, or ErrApprovalNotFound.
	GetApproval(ctx context.Context, instanceID, approvalID string) (*Approval, error)
	// ListApprovals returns approvals matching the filter, newest first.
	ListApprovals(ctx context.Context, f ListApprovalsFilter) ([]*Approval, error)
	// ResolveApproval atomically claims a pending approval, moving it to status (approved/denied) and
	// recording the decider. It transitions only from pending, so two concurrent decisions on one approval yield a
	// single winner: the loser gets ErrApprovalNotResolvable. This is the guard behind double-submit safety.
	ResolveApproval(ctx context.Context, instanceID, approvalID, status, decidedBy string) (*Approval, error)
	// CancelPendingApprovals moves every still-pending approval of a run to cancelled, in one statement. Cancelling a
	// run withdraws its outstanding approval: unlike ResolveApproval it is not a human decision, so it records no
	// decider and is not guarded against double-submit (it is a bulk, idempotent sweep). It returns the number of
	// approvals it cancelled, which is normally zero or one.
	CancelPendingApprovals(ctx context.Context, instanceID, runID string) (int64, error)
}

// ErrApprovalNotResolvable is returned by ResolveApproval when the approval is not pending (already decided or
// cancelled). It lets the API reject a double-submit with FailedPrecondition rather than silently re-deciding.
var ErrApprovalNotResolvable = errors.New("act: approval is not pending")
