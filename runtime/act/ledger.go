package act

import (
	"context"
	"errors"
	"time"
)

// defaultLeaseMargin is added to an action's timeout to derive its execution lease (§12): the extra headroom covers
// the finalize round-trip after the external call returns, so a healthy attempt always records its terminal outcome
// well inside the lease and never reclaims itself. Recovery only reclaims an attempt whose lease has fully elapsed.
const defaultLeaseMargin = 60 * time.Second

// ErrActionNotFound is returned when no action exists for the requested (instance, run, tool call).
var ErrActionNotFound = errors.New("act: action not found")

// ErrActionTransitionInvalid is returned when a caller attempts a ledger transition that the action state machine
// does not permit from the action's current state (§11.2). It is distinct from a store error so the gateway can
// treat it as a programming/replay fault rather than a transient failure.
var ErrActionTransitionInvalid = errors.New("act: invalid action transition")

// ActionStatus is the state of a single external write in the action ledger (§11.2). The ledger is the durable,
// auditable record of every effect the run proposes: it exists so a write has identity (§6.4) and so an at-least-once
// retry can tell whether the effect already happened before re-issuing it (§12).
type ActionStatus string

const (
	// ActionProposed is the initial state: the model proposed the tool call and the gateway recorded it, before any
	// policy decision or execution. Recording it here is what gives the write identity prior to any external effect.
	ActionProposed ActionStatus = "proposed"
	// ActionPolicyRejected is terminal: the deterministic policy denied the action (unknown tool, kill switch, deny
	// override). No external effect ran.
	ActionPolicyRejected ActionStatus = "policy_rejected"
	// ActionApprovalPending marks an action suspended waiting for a human decision on the exact proposed arguments.
	ActionApprovalPending ActionStatus = "approval_pending"
	// ActionApproved marks a human-approved action, bound to the exact args hash the human saw, ready to execute.
	ActionApproved ActionStatus = "approved"
	// ActionExecuting marks an action whose external write is in flight. Recorded before the call so a crash mid-write
	// is detectable on recovery: an action found still executing was interrupted, and only a retriable class may be
	// safely re-issued (§12).
	ActionExecuting ActionStatus = "executing"
	// ActionSucceeded marks a confirmed external write; ExternalReference holds the durable handle it produced.
	ActionSucceeded ActionStatus = "succeeded"
	// ActionFailed is terminal: the external write definitively failed. It is not retried automatically in v1.
	ActionFailed ActionStatus = "failed"
	// ActionIndeterminate is terminal-until-human: the external result could not be confirmed (a timeout or dropped
	// connection after send). A non-idempotent action here must NOT be retried automatically (§12); a human resolves it.
	ActionIndeterminate ActionStatus = "indeterminate"
	// ActionVerified marks a succeeded action whose effect a verification tool subsequently confirmed in the target
	// system (§16.1). It is the strongest terminal state.
	ActionVerified ActionStatus = "verified"
)

// IsTerminal reports whether an action can no longer transition. A terminal action's row is never mutated again, so a
// replayed transition targeting it is a no-op.
func (s ActionStatus) IsTerminal() bool {
	switch s {
	case ActionPolicyRejected, ActionFailed, ActionIndeterminate, ActionVerified:
		return true
	default:
		return false
	}
}

// actionEdges is the action state machine (§11.2): the set of transitions the ledger permits from each state. Any
// transition not listed here (and not a same-state replay) is rejected with ErrActionTransitionInvalid. The gateway
// never issues an illegal edge because it reads the action before executing and short-circuits on an advanced state;
// this map is the defensive backstop that keeps the ledger's history well-formed even so.
var actionEdges = map[ActionStatus]map[ActionStatus]bool{
	ActionProposed: {
		ActionPolicyRejected:  true,
		ActionApprovalPending: true,
		ActionApproved:        true,
	},
	ActionApprovalPending: {
		ActionApproved: true,
	},
	ActionApproved: {
		ActionExecuting: true,
		// Reauthorization or rebinding can refuse an already-approved action at execution time (§17.2): the tool was
		// revoked, the policy hardened to deny, or the proposal no longer matches what was approved. Recording that as
		// policy_rejected keeps the ledger truthful that no external effect ran.
		ActionPolicyRejected: true,
	},
	ActionExecuting: {
		ActionSucceeded:     true,
		ActionFailed:        true,
		ActionIndeterminate: true,
	},
	ActionSucceeded: {
		ActionVerified: true,
	},
}

// NewAction is the immutable identity and proposal of an external write, recorded once in the proposed state before
// any effect (§12). Everything here is fixed for the life of the action; mutable execution state (status, external
// reference, result) is written by TransitionAction.
type NewAction struct {
	// ToolCallID is the model's identifier for this tool call. Together with RunID and InstanceID it is the action's
	// primary key: proposing the same call twice (a replay) is a no-op, which is what makes the write single-effect.
	ToolCallID string
	RunID      string
	InstanceID string
	AgentName  string
	// Tool/Connector/ToolVersion identify the concrete effect. ToolVersion pins the exact tool schema the snapshot
	// captured, so an audit can tell which contract produced the write even if the connector later changes (§16.3).
	Tool        string
	Connector   string
	ToolVersion string
	// Class is the idempotency classification that governs whether an uncertain result may be retried (§12).
	Class ToolClass
	// ArgsHash is the canonical hash of the proposed arguments; the approval binds to it (§11.2). IdempotencyKey is
	// the stable dedup key derived from (run, tool call, args hash) (§12).
	ArgsHash       string
	IdempotencyKey string
	// RedactedArgs is the proposed arguments as persisted: secret references are preserved, resolved secret values are
	// never stored (§16.2). It is what a human sees in the approval inbox.
	RedactedArgs map[string]any
	// PolicyDecision and PolicyReason record the deterministic decision that routed this action (§6.3).
	PolicyDecision PolicyDecision
	PolicyReason   string
	// Proposal is a short human-readable description of the intended effect.
	Proposal string
	// RequestedBy is the actor subject on whose behalf the run proposed the action.
	RequestedBy string
}

// Action is the stored, queryable state of one external write: its immutable identity and proposal plus its current
// execution state. It is the audit source of truth for who proposed, approved and performed each effect (§18.3).
type Action struct {
	ToolCallID  string
	RunID       string
	InstanceID  string
	AgentName   string
	Tool        string
	Connector   string
	ToolVersion string
	Class       ToolClass

	ArgsHash       string
	IdempotencyKey string
	RedactedArgs   map[string]any

	PolicyDecision PolicyDecision
	PolicyReason   string
	Proposal       string

	Status ActionStatus
	// AttemptID identifies the live execution attempt that claimed this action's write (§12). The claim from approved
	// to executing is an atomic conditional update that stamps a fresh AttemptID alongside a lease_expiry, so exactly
	// one caller ever owns the write; the attempt id gates the owner's finalize (a reclaimed attempt cannot overwrite
	// the recovered outcome) and distinguishes a live attempt from a crashed one whose lease has expired.
	AttemptID string
	// ExternalReference is the durable handle the target system returned (e.g. a Jira issue key), set on success.
	ExternalReference string
	// RedactedResult and RedactedVerification are the tool result and verification, scrubbed of secrets before
	// storage (§17.2). Nil until the corresponding step runs.
	RedactedResult map[string]any
	// RedactedMessage is the tool's text result (e.g. a greeting or a "created PROJ-123" summary), scrubbed of secrets.
	// It is distinct from RedactedResult (the structured output): a tool that returns only text lands here. Empty until
	// the write succeeds. It is what the durable loop feeds back to the model so it can compose a closing turn.
	RedactedMessage      string
	RedactedVerification map[string]any
	// Error is a sanitized failure message, set only for a failed or indeterminate action. It never carries secrets.
	Error string

	RequestedBy string
	DecidedBy   string

	CreatedOn  time.Time
	UpdatedOn  time.Time
	ExecutedOn *time.Time // latches when the action first enters executing
	VerifiedOn *time.Time // set when the action is verified
}

// ActionTransition moves an action to a new status, carrying the fields that become known at that edge. Only the
// fields relevant to the target status are read: ExternalReference/RedactedResult/RedactedMessage on succeeded, Error
// on failed/indeterminate, RedactedVerification on verified, DecidedBy on approved. A same-status transition is an
// idempotent no-op so a replayed step is safe.
type ActionTransition struct {
	InstanceID string
	RunID      string
	ToolCallID string
	Status     ActionStatus

	// AttemptID, when set, guards the transition on ownership: the underlying update applies only while the row is still
	// executing under this exact attempt (§12). A finalize whose attempt no longer owns the row (its lease expired and
	// recovery reclaimed it) affects no row and does NOT overwrite the recovered outcome, so a late owner never
	// contradicts the ledger. Empty means an unconditional transition (the pre-execution edges carry no attempt).
	AttemptID string

	ExternalReference    string
	RedactedResult       map[string]any
	RedactedMessage      string
	RedactedVerification map[string]any
	Error                string
	DecidedBy            string
}

// ActionLedger is the durable home of the action ledger (table agent_actions, §15.1). It records every proposed
// external write and its lifecycle so a write has identity, an at-least-once retry can detect a prior effect, and an
// audit can reconstruct who did what (§12, §18.3). It is scoped by instance so tenant isolation is enforced in the
// store, and every write is idempotent so a replayed DBOS step never duplicates a row or an effect.
//
// It is deliberately separate from ActionExecutor: the ledger records identity and state, while ActionExecutor
// performs the outside-world call. The gateway drives both.
type ActionLedger interface {
	// ProposeAction inserts an action in the proposed state, idempotent on (instance, run, tool call): a replayed
	// proposal is a no-op, not a duplicate. Two DIFFERENT proposals colliding on one tool call are rejected with
	// ErrActionTransitionInvalid, so a mutated proposal cannot silently ride an existing row (§12, §17.1).
	ProposeAction(ctx context.Context, a NewAction) error
	// TransitionAction applies a state-machine edge (§11.2) and returns the resulting action. A same-status transition
	// is a no-op; an edge the machine does not permit returns ErrActionTransitionInvalid.
	TransitionAction(ctx context.Context, t ActionTransition) (*Action, error)
	// ClaimExecuting atomically transitions an approved action to executing, stamping attemptID as its owner and a lease
	// that bounds how long the attempt is presumed alive. It is the exclusive claim behind single-effect execution
	// (§12): the UPDATE is conditional on the row still being approved, so of any number of concurrent callers only the
	// one whose UPDATE flips approved->executing gets claimed=true and may perform the write. A caller that finds the row
	// already past approved gets claimed=false and MUST NOT execute; the returned action carries the current state so it
	// can report the recorded outcome instead. lease_expiry is set to now()+lease using the database clock, so it is
	// compared against the same clock on recovery.
	ClaimExecuting(ctx context.Context, instanceID, runID, toolCallID, attemptID string, lease time.Duration) (claimed bool, action *Action, err error)
	// ReclaimExpiredLease moves an executing action to indeterminate ONLY if its lease has expired (or was never set, as
	// after a crash that reached executing without a claim): an expired lease means the owning attempt is presumed dead,
	// so its interrupted write is handed to a human rather than re-issued (§12). It is the recovery-safe counterpart to
	// ClaimExecuting: an executing action under a LIVE lease is left untouched (reclaimed=false) so a still-working owner
	// is never clobbered mid-write. The returned action is the current row either way, to report against.
	ReclaimExpiredLease(ctx context.Context, instanceID, runID, toolCallID, reason string) (reclaimed bool, action *Action, err error)
	// GetAction returns the action for (instance, run, tool call), or ErrActionNotFound. It is the read the gateway
	// performs before executing, to short-circuit an already-completed or interrupted write (§12).
	GetAction(ctx context.Context, instanceID, runID, toolCallID string) (*Action, error)
	// ListActions returns a run's actions, newest first, for the audit view and the approval inbox.
	ListActions(ctx context.Context, instanceID, runID string) ([]*Action, error)
}
