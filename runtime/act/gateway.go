package act

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/google/uuid"
)

// defaultActionTimeout bounds a single external write when the gateway is not configured otherwise. defaultMaxOutput
// bounds how much of a tool result is captured, so a hostile or verbose connector cannot flood the ledger or logs.
const (
	defaultActionTimeout = 30 * time.Second
	defaultMaxOutput     = 64 * 1024
)

// ToolProposal is the concrete action the run proposes: the tool, its connector and its arguments. It is the input
// to the gateway. In production the runtime/ai loop emits it as a structured tool call; this frente models it
// explicitly so the gateway, ledger and executor path can be wired and exercised end to end independently of the
// loop refactor and of the outbound MCP client (built by a parallel frente).
type ToolProposal struct {
	// ToolCallID is the model's identifier for the call. It is the action's identity within the run (§6.4): the ledger
	// keys on it and the idempotency key derives from it, so it must be stable for a given proposed effect.
	ToolCallID string
	Tool       string
	Connector  string
	// Args are the proposed arguments as the model produced them, with secret references (never resolved secrets) in
	// place. They are validated against the tool's schema, hashed for the approval binding, and persisted redacted.
	Args map[string]any
	// Summary is a short human-readable description of the intended effect, shown in the approval inbox.
	Summary string
}

// ToolDescriptor is the gateway's fixed knowledge of an action tool: its connector, pinned version, idempotency
// classification, input schema and read-only hint. It is the tool manifest entry the reconciler compiled (§8.2) and
// the MCP tools/list pinned at snapshot time (§16.3); the model never influences it. A tool with no descriptor is
// unknown and is denied.
type ToolDescriptor struct {
	Name      string
	Connector string
	Version   string
	// Class is the idempotency classification (§12). It must come from the tool's contract, never from an untrusted
	// annotation. An unknown class can never be auto-executed and is never auto-retried on an uncertain result.
	Class ToolClass
	// InputSchema is the tool's fixed argument schema. The gateway validates the proposed arguments against it before
	// anything else (§16.2). A nil schema denies: an action whose arguments cannot be validated is not executed.
	InputSchema *jsonschema.Schema
	// ReadOnlyHint is the server-announced read-only annotation; TrustReadOnly is whether the connector opted to trust
	// it. An untrusted hint is ignored (§16.3).
	ReadOnlyHint  bool
	TrustReadOnly bool
	// AuthSecret is the name of the connector's bearer-credential secret, if it has one. The gateway resolves it only
	// to add its plaintext to the redaction set, so a tool result that echoes the connector token back is scrubbed
	// before it reaches the ledger or a log (§16.2, §17.1). It is a reference, never the secret; the value never leaves
	// the resolution/scrub path.
	AuthSecret string
	// Verifiable reports whether the tool supports post-write verification (§16.1). When true and the write succeeds,
	// the gateway calls ActionExecutor.Verify to confirm the effect.
	Verifiable bool
}

// ToolRegistry resolves the descriptor for a tool available to an agent. Because the alignment model is default-open
// (connecting an MCP server exposes all its tools, filtered only by optional include/exclude), registry membership is
// the allowlist for action tools: a tool the registry does not return is not available to the agent and is denied.
type ToolRegistry interface {
	// Lookup returns the descriptor for tool, or ok=false if the tool is not available to the agent. An error is a
	// transport/config failure distinct from a missing tool; the gateway treats it fail-closed.
	Lookup(ctx context.Context, instanceID, agentName, tool string) (desc ToolDescriptor, ok bool, err error)
}

// KillScope identifies what a kill-switch query is about, so a switch can be engaged at the platform, tenant
// (instance) or agent level, optionally narrowed to a tool or connector (§17.2).
type KillScope struct {
	InstanceID string
	AgentName  string
	Tool       string
	Connector  string
}

// KillSwitch reports whether actions are disabled for a scope. It is the simple, UI-less control the design asks for:
// one check in the gateway that, when engaged, denies (§17.2). A nil KillSwitch on a Gateway means no switch is wired
// and actions proceed normally; a configured switch that errors is treated as engaged (fail closed).
type KillSwitch interface {
	Disabled(ctx context.Context, scope KillScope) (disabled bool, reason string, err error)
}

// TraceContext carries correlation identifiers into an ActionExecutor call. It never carries secrets or arguments:
// it exists so an external call can be tied back to its run and trace in observability without leaking data (§18.1).
type TraceContext struct {
	RunID     string
	AgentName string
	TraceID   string
}

// ExecuteOutcome classifies the result of an external write (§12). It is the gateway's authoritative signal for how
// the run and the ledger proceed, and it is what keeps an uncertain non-idempotent write from being retried.
type ExecuteOutcome string

const (
	// OutcomeSucceeded means the write is confirmed to have happened; ExternalReference holds its handle.
	OutcomeSucceeded ExecuteOutcome = "succeeded"
	// OutcomeFailed means the write definitively did not happen; it is safe (though not automatic in v1) to retry.
	OutcomeFailed ExecuteOutcome = "failed"
	// OutcomeIndeterminate means the executor could not tell whether the write landed. The ledger records it and the
	// run does NOT retry automatically (§12); a human resolves it.
	OutcomeIndeterminate ExecuteOutcome = "indeterminate"
)

// ExecuteRequest is the fully-prepared external write handed to an ActionExecutor. The gateway has already validated
// arguments, resolved policy, resolved secrets and derived the idempotency key by the time this is built: an executor
// implementation performs exactly one already-authorized call and reports a classifiable result.
type ExecuteRequest struct {
	RunID string
	// InstanceID scopes the write to its tenant instance. A runtime-backed executor resolves the connector and its
	// secret against this instance; the in-process mcpExecutor (bound to a fixed client) ignores it.
	InstanceID  string
	ToolCallID  string
	Tool        string
	Connector   string
	ToolVersion string
	// ReadOnly is true only for a tool the gateway treats as a trusted read-only tool (its server hint is trusted by
	// the connector). An executor uses it to classify a tool-reported failure: a write tool's IsError is indeterminate
	// (it may have written before failing), a read-only tool's is a safe failed (§12).
	ReadOnly bool
	// Args are the fully-resolved arguments, secret references replaced by their plaintext values. They exist only for
	// the duration of the call; the gateway persists only the redacted form. An implementation MUST NOT log them in
	// the clear (§16.2).
	Args map[string]any
	// IdempotencyKey is the stable dedup key for this action. An implementation MUST pass it to the target system when
	// the system supports an idempotency key, so an at-least-once re-issue does not double-write (§12).
	IdempotencyKey string
	Trace          TraceContext
	// Timeout bounds the call; MaxOutputBytes bounds how much of the result to return. Both come from the gateway.
	Timeout        time.Duration
	MaxOutputBytes int
}

// ExecuteResult is an ActionExecutor's classified report of a write. The executor is authoritative on Outcome: it is
// the only component that knows whether the target system confirmed, refused or left the result ambiguous.
type ExecuteResult struct {
	Outcome ExecuteOutcome
	// ExternalReference is the durable handle the target system returned (e.g. a Jira issue key). Persisted.
	ExternalReference string
	// Output is the raw tool result to persist after the gateway redacts it; it may be nil.
	Output map[string]any
	// Message is a short human summary (e.g. "created PROJ-123"), persisted after redaction and surfaced on the run.
	Message string
}

// VerifyRequest asks an ActionExecutor to confirm a prior write exists in the target system, by external reference.
type VerifyRequest struct {
	RunID             string
	ToolCallID        string
	Tool              string
	Connector         string
	ExternalReference string
	Trace             TraceContext
	Timeout           time.Duration
}

// VerifyResult is an ActionExecutor's verification report. Confirmed=false leaves the action unverified, not failed:
// the write already succeeded; verification is a stronger, best-effort confirmation (§16.1).
type VerifyResult struct {
	Confirmed bool
	Detail    string
}

// ActionExecutor performs the actual external write behind a tool call, and later verifies its effect. It is the seam
// the outbound MCP client implements (built by a parallel frente): the gateway resolves policy, secrets, idempotency
// and redaction around it, so an implementation only performs one already-authorized call and reports a classifiable
// outcome. Tests substitute a fake.
//
// This is the boundary the design's action-tool contract (prepare/execute/verify/compensate, §12) reduces to for v1:
// prepare is the gateway's Authorize; execute and verify are here; compensate is deferred (v1 gates every write on
// human approval and never auto-retries, so no automatic compensation is required yet).
type ActionExecutor interface {
	// Execute performs the external write. It MUST pass req.IdempotencyKey to the target system when supported, and
	// MUST classify the result via ExecuteResult.Outcome — in particular returning OutcomeIndeterminate (or a non-nil
	// error) when it cannot tell whether the write landed, so the gateway does not treat an ambiguous result as
	// success (§12).
	Execute(ctx context.Context, req ExecuteRequest) (ExecuteResult, error)
	// Verify confirms or refutes that a prior Execute's effect exists in the target system. It is called only for a
	// succeeded action of a verifiable tool.
	Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error)
}

// PolicyConfig holds the deterministic override layers the gateway applies on top of an agent's own approval config.
// Each layer may only harden a decision (§8.1); neither can relax the agent's rule. These are the platform-wide and
// per-tenant floors (§17.2, open question 8).
type PolicyConfig struct {
	// Global is the platform floor no project can relax. Tenant is the per-tenant floor. Both key by tool name.
	Global map[string]PolicyDecision
	Tenant map[string]PolicyDecision
}

// Gateway is the single choke point every action tool call passes through (§16.2). The LLM never reaches a connector
// directly: it proposes, and the gateway checks the allowlist (registry membership), validates arguments against the
// tool schema, classifies risk and evaluates deterministic policy, resolves secrets server-side, derives the
// idempotency key and trace, applies the timeout and size caps, records the action ledger, and redacts inputs and
// outputs before anything is persisted.
type Gateway struct {
	// Registry resolves tool descriptors (the allowlist and schema source). Required.
	Registry ToolRegistry
	// Ledger records the action's identity and lifecycle (§12). Required.
	Ledger ActionLedger
	// Executor performs the external write and verification. Required.
	Executor ActionExecutor
	// Secrets resolves secret references server-side. May be nil if no proposal references secrets; a reference with a
	// nil resolver is a hard error (fail closed).
	Secrets SecretResolver
	// KillSwitch, when set, disables actions per scope (§17.2). Nil means no switch is wired.
	KillSwitch KillSwitch
	// Policy holds the global/tenant hardening layers.
	Policy PolicyConfig
	// Timeout and MaxOutputBytes bound each external write. Zero values fall back to package defaults.
	Timeout        time.Duration
	MaxOutputBytes int
	Logger         *slog.Logger
}

// Authorization is the gateway's ruling on a proposal: the deterministic decision plus the identity fields the ledger
// and the approval bind to. It is produced even for a denied proposal, so a rejected action still has a hash and a
// dedup key on record for audit.
type Authorization struct {
	Decision       PolicyDecision
	Reason         string
	Class          ToolClass
	ArgsHash       string
	IdempotencyKey string
	ToolVersion    string
	Connector      string
	Verifiable     bool
}

// AuthorizeInput is the fully-resolved, side-effect-free input to Authorize. The gateway resolves the descriptor, the
// kill switch and the policy layers before calling it, so Authorize is a pure function that cannot be swayed by a
// prompt and is trivially testable in isolation.
type AuthorizeInput struct {
	RunID       string
	Proposal    ToolProposal
	Descriptor  ToolDescriptor
	Found       bool
	AutoApprove map[string]bool
	Global      map[string]PolicyDecision
	Tenant      map[string]PolicyDecision
	KillSwitch  bool
	KillReason  string
}

// Authorize applies validation and deterministic policy to a proposal, purely and fail-closed (§6.3, §6.5). It always
// computes the args hash and idempotency key first — even a denied proposal gets an identity for the ledger — then
// resolves the decision in strict precedence: kill switch, then unknown tool, then argument validation, then policy.
// Any ambiguity resolves toward deny.
func Authorize(in AuthorizeInput) Authorization {
	auth := Authorization{
		Class:       in.Descriptor.Class,
		ToolVersion: in.Descriptor.Version,
		Connector:   in.Proposal.Connector,
		Verifiable:  in.Descriptor.Verifiable,
	}
	argsHash, hashErr := HashCanonicalArgs(in.Proposal.Args)
	auth.ArgsHash = argsHash
	if hashErr == nil {
		auth.IdempotencyKey = DeriveIdempotencyKey(in.RunID, in.Proposal.ToolCallID, argsHash)
	}

	switch {
	case in.KillSwitch:
		// The operator's stop wins over everything (§17.2).
		auth.Decision, auth.Reason = PolicyDeny, killReasonOr(in.KillReason)
	case !in.Found:
		// No descriptor: the tool is not available to the agent, so it cannot be validated or classified. Deny rather
		// than surface an unknown tool to a human as a legitimate proposal (§6.5).
		auth.Decision, auth.Reason = PolicyDeny, "unknown tool: not available to agent"
	case hashErr != nil:
		auth.Decision, auth.Reason = PolicyDeny, "invalid arguments: not canonicalizable"
	default:
		if err := validateArgs(in.Descriptor.InputSchema, in.Proposal.Args); err != nil {
			auth.Decision, auth.Reason = PolicyDeny, "invalid arguments: "+err.Error()
			break
		}
		res := Evaluate(PolicyInput{
			Tool:      in.Proposal.Tool,
			Connector: in.Proposal.Connector,
			// Registry membership is the allowlist for action tools (default-open model): a found descriptor means the
			// agent may use the tool, subject to policy.
			Allowed:       map[string]bool{in.Proposal.Tool: true},
			AutoApprove:   in.AutoApprove,
			Class:         in.Descriptor.Class,
			ReadOnlyHint:  in.Descriptor.ReadOnlyHint,
			TrustReadOnly: in.Descriptor.TrustReadOnly,
			GlobalPolicy:  in.Global,
			TenantPolicy:  in.Tenant,
		})
		auth.Decision, auth.Reason = res.Decision, res.Reason
	}
	return auth
}

// Propose is the gateway's first durable step for an action: it resolves the descriptor and kill switch, authorizes
// the proposal, and records the action ledger in one place — proposed, then the policy decision (policy_rejected /
// approval_pending / approved). Both ledger writes are idempotent, so a replayed Propose step neither duplicates the
// action nor advances it twice. It returns the Authorization for the workflow to route on.
func (g *Gateway) Propose(ctx context.Context, in ProposeActionInput) (Authorization, error) {
	// Freeze the proposed arguments before anything reads them: the args hash, the persisted ledger row and the
	// arguments eventually executed must all bind to one immutable snapshot, so a Proposer that mutates its own map
	// after proposing cannot desynchronize what was approved from what runs (§17.1). in is a value, so this rebinds only
	// this call's copy, and Authorize and ProposeAction below both see the frozen args.
	in.Proposal.Args = freezeArgs(in.Proposal.Args)

	desc, found, err := g.Registry.Lookup(ctx, in.InstanceID, in.AgentName, in.Proposal.Tool)
	if err != nil {
		// A registry failure is fail-closed: without the manifest the gateway cannot validate or classify, so it must
		// not proceed. Surfacing the error fails the run rather than silently allowing an unvalidated write.
		return Authorization{}, fmt.Errorf("act: lookup tool %q: %w", in.Proposal.Tool, err)
	}
	killed, killReason := g.killDisabled(ctx, KillScope{
		InstanceID: in.InstanceID, AgentName: in.AgentName, Tool: in.Proposal.Tool, Connector: in.Proposal.Connector,
	})

	auth := Authorize(AuthorizeInput{
		RunID:       in.RunID,
		Proposal:    in.Proposal,
		Descriptor:  desc,
		Found:       found,
		AutoApprove: in.AutoApprove,
		Global:      g.Policy.Global,
		Tenant:      g.Policy.Tenant,
		KillSwitch:  killed,
		KillReason:  killReason,
	})

	// Record the action's identity in the proposed state before any external effect (§12). Idempotent on the tool call.
	if err := g.Ledger.ProposeAction(ctx, NewAction{
		ToolCallID:     in.Proposal.ToolCallID,
		RunID:          in.RunID,
		InstanceID:     in.InstanceID,
		AgentName:      in.AgentName,
		Tool:           in.Proposal.Tool,
		Connector:      in.Proposal.Connector,
		ToolVersion:    auth.ToolVersion,
		Class:          auth.Class,
		ArgsHash:       auth.ArgsHash,
		IdempotencyKey: auth.IdempotencyKey,
		RedactedArgs:   in.Proposal.Args, // proposal args carry secret references, not values: safe to persist as-is
		PolicyDecision: auth.Decision,
		PolicyReason:   auth.Reason,
		Proposal:       in.Proposal.Summary,
		RequestedBy:    in.Actor.Subject,
	}); err != nil {
		return Authorization{}, fmt.Errorf("act: propose action: %w", err)
	}

	// Record the policy decision as the next ledger state. Deny and allow are recorded directly; approval_required
	// parks the action pending a human, which the workflow then drives via the existing approval mechanism.
	var next ActionStatus
	switch auth.Decision {
	case PolicyDeny:
		next = ActionPolicyRejected
	case PolicyApprovalRequired:
		next = ActionApprovalPending
	case PolicyAllow:
		next = ActionApproved
	default:
		// Fail closed on any unrecognized decision.
		next = ActionPolicyRejected
	}
	if _, err := g.Ledger.TransitionAction(ctx, ActionTransition{
		InstanceID: in.InstanceID, RunID: in.RunID, ToolCallID: in.Proposal.ToolCallID, Status: next,
	}); err != nil {
		return Authorization{}, fmt.Errorf("act: record policy decision: %w", err)
	}
	return auth, nil
}

// ProposeActionInput is the input to Gateway.Propose.
type ProposeActionInput struct {
	InstanceID string
	RunID      string
	AgentName  string
	Actor      Actor
	Proposal   ToolProposal
	// AutoApprove is the agent's approval.automatic set. In v1 all writes require approval, so this is typically empty.
	AutoApprove map[string]bool
}

// RecordApproved advances an approval-pending action to approved once a human has decided, binding the decision to the
// exact args hash and recording the decider. It rejects a decision whose args hash does not match the action's — a
// changed argument after approval is a different effect and must be re-approved (§11.2). It is idempotent: re-applying
// approved is a no-op.
func (g *Gateway) RecordApproved(ctx context.Context, in ApprovedInput) error {
	action, err := g.Ledger.GetAction(ctx, in.InstanceID, in.RunID, in.ToolCallID)
	if err != nil {
		return err
	}
	if in.ArgsHash != action.ArgsHash {
		// The approval was for different arguments than the action now holds: fail closed, do not execute (§17.2).
		return fmt.Errorf("act: approval args hash mismatch for %q: approved %s, action %s: %w",
			in.ToolCallID, in.ArgsHash, action.ArgsHash, ErrActionTransitionInvalid)
	}
	if _, err := g.Ledger.TransitionAction(ctx, ActionTransition{
		InstanceID: in.InstanceID, RunID: in.RunID, ToolCallID: in.ToolCallID, Status: ActionApproved, DecidedBy: in.DecidedBy,
	}); err != nil {
		return fmt.Errorf("act: record approval: %w", err)
	}
	return nil
}

// ApprovedInput is the input to Gateway.RecordApproved.
type ApprovedInput struct {
	InstanceID string
	RunID      string
	ToolCallID string
	// ArgsHash is the hash the human approved; it must equal the action's stored args hash or the decision is rejected.
	ArgsHash  string
	DecidedBy string
}

// Execute performs the approved external write exactly once under at-least-once delivery (§12). It is the idempotency
// heart of the gateway: it reads the action first and short-circuits on any state that means the write already
// happened or must not be retried, so a replayed step (DBOS re-running after a crash, or a lost checkpoint) never
// drives a second external effect. Only then does it resolve secrets, call the executor with the stable idempotency
// key, redact the result, and record the terminal ledger state. It returns nil error for a business failure or an
// indeterminate result — those are outcomes to be recorded, not step errors to be retried — reserving errors for
// infrastructure faults.
func (g *Gateway) Execute(ctx context.Context, in ExecuteActionInput) (ExecuteReport, error) {
	action, err := g.Ledger.GetAction(ctx, in.InstanceID, in.RunID, in.Proposal.ToolCallID)
	if err != nil {
		return ExecuteReport{}, err
	}

	switch action.Status {
	case ActionSucceeded, ActionVerified:
		// Already done on a prior attempt: return the recorded effect without re-executing (§12).
		return ExecuteReport{Outcome: OutcomeSucceeded, ExternalReference: action.ExternalReference}, nil
	case ActionIndeterminate:
		// A prior attempt could not be confirmed: never auto-retry (§12).
		return ExecuteReport{Outcome: OutcomeIndeterminate, ExternalReference: action.ExternalReference}, nil
	case ActionFailed:
		return ExecuteReport{Outcome: OutcomeFailed}, nil
	case ActionExecuting:
		// A prior or concurrent attempt owns this write. Do NOT blindly mark it indeterminate: that would clobber a
		// still-live attempt mid-write and let the owner then report success over an indeterminate ledger. Reclaim it to
		// indeterminate only if its lease has expired (owner presumed dead), leaving a live owner to finish (§12). In v1
		// no class is auto-retried, so recovery never re-issues the effect.
		return g.recoverExecuting(ctx, in, action)
	case ActionApproved:
		// Normal path: reauthorize, rebind, then claim below.
	default:
		// proposed / approval_pending / policy_rejected: not authorized to execute. Fail closed.
		return ExecuteReport{}, fmt.Errorf("act: action %q not approved for execution (status %s)", in.Proposal.ToolCallID, action.Status)
	}

	// Reauthorize immediately before the write (§17.2): an approval decided minutes ago must not execute if the tool was
	// revoked from the registry or the policy hardened to deny during the wait. Reload the descriptor and re-evaluate;
	// a tool no longer available, or a fresh decision stricter than the one this action was approved under, refuses.
	desc, found, err := g.Registry.Lookup(ctx, in.InstanceID, in.AgentName, in.Proposal.Tool)
	if err != nil {
		return ExecuteReport{}, fmt.Errorf("act: reauthorize lookup %q: %w", in.Proposal.Tool, err)
	}
	if reason, ok := g.reauthorize(in, action, desc, found); !ok {
		return g.refuse(ctx, in, "reauthorization denied: "+reason)
	}

	// Rebind to the approved action (§11.2, §17.1): the write must be exactly what was proposed and approved. Recompute
	// the args hash and compare it — along with tool, connector and pinned version — against the ledger row, so a
	// Proposer that mutated the map, or an argument changed after approval, cannot execute a different effect.
	if reason, ok := rebindsToAction(in.Proposal, desc, action); !ok {
		return g.refuse(ctx, in, "proposal does not match approved action: "+reason)
	}

	// Re-check the kill switch immediately before the write: it may have been engaged during a long approval wait, and
	// fail-closed means an engaged switch stops the action even now (§17.2). The action stays approved (not executed).
	if killed, reason := g.killDisabled(ctx, KillScope{
		InstanceID: in.InstanceID, AgentName: in.AgentName, Tool: in.Proposal.Tool, Connector: in.Proposal.Connector,
	}); killed {
		return ExecuteReport{}, fmt.Errorf("act: action %q blocked by kill switch: %s", in.Proposal.ToolCallID, reason)
	}

	// Claim the write exclusively (§12): the approved->executing transition is an atomic conditional UPDATE, so of any
	// number of concurrent callers only the one that flips the row owns the write. A caller that loses the claim never
	// executes; it reports the outcome the winner recorded. This is what keeps the write single-effect under
	// at-least-once, beyond DBOS's own one-workflow-per-run guarantee.
	attemptID := uuid.NewString()
	claimed, claimedAction, err := g.Ledger.ClaimExecuting(ctx, in.InstanceID, in.RunID, in.Proposal.ToolCallID, attemptID, g.leaseDuration())
	if err != nil {
		return ExecuteReport{}, fmt.Errorf("act: claim executing: %w", err)
	}
	if !claimed {
		// Lost the claim to a concurrent attempt: reconcile without a second effect. Report the recorded outcome, and if
		// the winner's lease has since expired without finalizing, reclaim it to indeterminate rather than clobber a live
		// owner (§12).
		return g.recoverExecuting(ctx, in, claimedAction)
	}
	action = claimedAction

	// Resolve secrets into a call-only copy from the FROZEN arguments the ledger recorded at propose, never from the
	// caller's live map: the executed arguments are then exactly the ones that were hashed and approved, even if a
	// Proposer mutated its map after proposing (§16.2, §17.1). The resolved values never touch the ledger; they build
	// the executor call and scrub any that echo back in the result. A resolution failure fails closed: no write.
	resolvedArgs, secretValues, err := resolveSecrets(ctx, g.Secrets, action.RedactedArgs)
	if err != nil {
		return g.finalizeOwned(ctx, in, attemptID,
			ActionTransition{Status: ActionFailed, Error: "secret resolution failed"},
			ExecuteReport{Outcome: OutcomeFailed})
	}
	// The connector's bearer token must be in the scrub set before ANY result is persisted or logged: it is not one of
	// the args, but a hostile connector could echo it back in any result field (§17.1). On the write path this is STRICT
	// (C1): if a declared credential cannot be resolved we cannot guarantee its redaction, so we fail closed rather than
	// run a write whose result we could not scrub. The failure is pre-dispatch, so no external effect ran.
	secretValues, secErr := g.connectorRedactionSecret(ctx, desc, secretValues)
	if secErr != nil {
		if g.Logger != nil {
			g.Logger.Warn("act: connector credential unavailable for redaction; refusing write",
				"run", in.RunID, "tool_call", in.Proposal.ToolCallID, "connector", desc.Connector)
		}
		return g.finalizeOwned(ctx, in, attemptID,
			ActionTransition{Status: ActionFailed, Error: "connector credential unavailable for redaction"},
			ExecuteReport{Outcome: OutcomeFailed})
	}

	res, execErr := g.Executor.Execute(ctx, ExecuteRequest{
		RunID:          in.RunID,
		InstanceID:     in.InstanceID,
		ToolCallID:     in.Proposal.ToolCallID,
		Tool:           in.Proposal.Tool,
		Connector:      in.Proposal.Connector,
		ToolVersion:    action.ToolVersion,
		ReadOnly:       desc.ReadOnlyHint && desc.TrustReadOnly,
		Args:           resolvedArgs,
		IdempotencyKey: action.IdempotencyKey,
		Trace:          TraceContext{RunID: in.RunID, AgentName: in.AgentName, TraceID: in.TraceID},
		Timeout:        g.timeout(),
		MaxOutputBytes: g.maxOutput(),
	})

	outcome := g.classify(res, execErr)
	switch outcome {
	case OutcomeSucceeded:
		// Redact, then cap both the structured output and the message at the gateway boundary before persisting: a
		// hostile or runaway executor must not flood the ledger even if its own size cap was absent or bypassed
		// (defense in depth, §16.2, D2).
		redacted := g.capOutput(RedactMap(res.Output, secretValues))
		message := truncate(redactString(res.Message, secretValues), g.maxOutput())
		// The external reference is attacker-controlled: a hostile server can return a secret or arbitrary junk as the
		// handle. Redact it, require it to match a bounded reference shape, and drop it otherwise — a rejected reference
		// does not undo the confirmed write, it just stores no handle (§17.1).
		ref, refOK := g.sanitizeExternalReference(res.ExternalReference, secretValues)
		if !refOK && g.Logger != nil {
			g.Logger.Warn("act: external reference rejected (malformed or leaked)", "run", in.RunID, "tool_call", in.Proposal.ToolCallID)
		}
		return g.finalizeOwned(ctx, in, attemptID,
			ActionTransition{Status: ActionSucceeded, ExternalReference: ref, RedactedResult: redacted, RedactedMessage: message},
			ExecuteReport{Outcome: OutcomeSucceeded, ExternalReference: ref, Message: message})
	case OutcomeIndeterminate:
		return g.finalizeOwned(ctx, in, attemptID,
			ActionTransition{Status: ActionIndeterminate, Error: truncate(redactExecErr(execErr, secretValues), g.maxOutput())},
			ExecuteReport{Outcome: OutcomeIndeterminate})
	default: // OutcomeFailed
		return g.finalizeOwned(ctx, in, attemptID,
			ActionTransition{Status: ActionFailed, Error: truncate(redactExecErr(execErr, secretValues), g.maxOutput())},
			ExecuteReport{Outcome: OutcomeFailed})
	}
}

// ExecuteActionInput is the input to Gateway.Execute.
type ExecuteActionInput struct {
	InstanceID string
	RunID      string
	AgentName  string
	Proposal   ToolProposal
	TraceID    string
	// AutoApprove is the agent's approval.automatic set, the SAME input Propose authorized under. Reauthorization
	// re-evaluates the proposal against the current policy with these inputs (§17.2), so an action auto-approved at
	// propose is not spuriously re-flagged as newly approval-required and rejected during the pre-write recheck. Empty
	// in v1 (all writes require approval); wired for the auto-approve path.
	AutoApprove map[string]bool
}

// ExecuteReport is the gateway's account of a write for the workflow: the classified outcome and, on success, the
// external reference and a redacted human message. The workflow maps the outcome to the run's terminal status.
type ExecuteReport struct {
	Outcome           ExecuteOutcome
	ExternalReference string
	Message           string
}

// finalizeOwned records the owning attempt's terminal outcome, guarded on still owning the executing row (§12). The
// transition carries attemptID: if recovery reclaimed the action while this attempt was mid-write (its lease expired),
// the guarded update affects no row and the reloaded action reveals the recovered state. In that case the owner does
// NOT report its own outcome — that would contradict the ledger — it reconciles to the recorded outcome and reports
// that. A confirmed-but-reclaimed write is logged (sanitized) as an audit signal rather than silently flipping the
// recovered indeterminate back to success: v1 leaves a human to resolve an indeterminate action. want is the report to
// return when the owner's outcome is the one that stuck. The reason on an indeterminate/failed transition is already
// capped by the caller (§16.2).
func (g *Gateway) finalizeOwned(ctx context.Context, in ExecuteActionInput, attemptID string, tr ActionTransition, want ExecuteReport) (ExecuteReport, error) {
	tr.InstanceID, tr.RunID, tr.ToolCallID, tr.AttemptID = in.InstanceID, in.RunID, in.Proposal.ToolCallID, attemptID
	updated, err := g.Ledger.TransitionAction(ctx, tr)
	if err != nil {
		return ExecuteReport{}, fmt.Errorf("act: record %s: %w", tr.Status, err)
	}
	if updated.Status != tr.Status {
		if g.Logger != nil {
			g.Logger.Warn("act: finalize lost ownership; reconciling to the recorded outcome",
				"run", in.RunID, "tool_call", in.Proposal.ToolCallID,
				"attempted", string(tr.Status), "recorded", string(updated.Status))
		}
		return reportRecorded(updated), nil
	}
	return want, nil
}

// ErrLeaseHeld is returned by Execute when the action is executing under a live lease held by another attempt: the
// owning attempt is presumed alive and may still finalize. It is an infrastructure signal (the caller should retry the
// step after a delay) rather than a terminal business outcome, so the workflow does not permanently fail the run while
// the owner may still succeed. A subsequent retry after the owner finalizes sees the recorded outcome; a retry after
// the lease expires reclaims the action to indeterminate (§12).
var ErrLeaseHeld = errors.New("act: action is executing under a live lease; retry later")

// recoverExecuting handles an action found (or reloaded) in the executing state that this caller does not own: it is
// either a live attempt still writing or the remains of a crashed one. It reclaims to indeterminate ONLY when the
// lease has expired (owner presumed dead) or was never set (a crash that reached executing without a claim); a live
// lease is left untouched so the owner is never clobbered mid-write (§12).
//
// When the lease is still live (the reclaim attempt leaves the action executing), the function returns ErrLeaseHeld
// rather than a terminal business outcome: reporting OutcomeIndeterminate here would terminalize the run permanently
// while the owner may still succeed, and nothing would revisit the action after the lease expires. By returning an
// infrastructure error, the workflow/DBOS retries the step — on the next attempt the owner will have finalized
// (reporting its recorded outcome) or the lease will have expired (reclaiming to indeterminate).
func (g *Gateway) recoverExecuting(ctx context.Context, in ExecuteActionInput, action *Action) (ExecuteReport, error) {
	if action.Status != ActionExecuting {
		// The owner finalized between the read and here (or the row is otherwise past executing): report what it recorded.
		return reportRecorded(action), nil
	}
	_, cur, err := g.Ledger.ReclaimExpiredLease(ctx, in.InstanceID, in.RunID, in.Proposal.ToolCallID,
		"interrupted mid-write; execution lease expired before the effect was confirmed")
	if err != nil {
		return ExecuteReport{}, fmt.Errorf("act: reclaim expired lease: %w", err)
	}
	if cur.Status == ActionExecuting {
		// The lease is still live: the owning attempt is presumed alive. Do not report an indeterminate outcome, which
		// would terminalize the run while the owner may still finalize successfully or the lease may still expire for a
		// proper reclaim. Return a retryable infrastructure error so DBOS retries the step.
		return ExecuteReport{}, fmt.Errorf("act: action %q executing under live lease (attempt %s): %w",
			in.Proposal.ToolCallID, cur.AttemptID, ErrLeaseHeld)
	}
	return reportRecorded(cur), nil
}

// leaseDuration is how long a claimed attempt is presumed alive: the action's own timeout plus a margin for the
// finalize round-trip, so a healthy attempt always records its outcome inside the lease and never reclaims itself.
func (g *Gateway) leaseDuration() time.Duration {
	return g.timeout() + defaultLeaseMargin
}

// classify decides the outcome from the executor's report and error, fail-closed on ambiguity. An explicit outcome
// from the executor is authoritative; a bare error is uncertain and, in v1, always indeterminate: the request may
// have been dispatched, so the effect is unknown and the write is never auto-retried for any class (§12). See the
// RetriableOnUncertain TODO for the prerequisites to reintroduce class-based retry.
func (g *Gateway) classify(res ExecuteResult, execErr error) ExecuteOutcome {
	if execErr != nil {
		// An executor that returns OutcomeFailed with an error is reporting a definitive pre-dispatch failure: the call
		// provably never went out, so the effect is known (no effect). Respect the authoritative classification rather
		// than degrading to indeterminate, which would leave a human to resolve what the executor already resolved.
		// Any other error without an explicit failed outcome is genuinely uncertain and stays indeterminate (§12).
		if res.Outcome == OutcomeFailed {
			return OutcomeFailed
		}
		return OutcomeIndeterminate
	}
	switch res.Outcome {
	case OutcomeSucceeded, OutcomeFailed, OutcomeIndeterminate:
		return res.Outcome
	default:
		// An executor that returns no error and no classification is treated as uncertain, never as an implicit
		// success (§6.5).
		return OutcomeIndeterminate
	}
}

// reauthorize re-evaluates an approved action against the current registry and policy immediately before the write
// (§17.2). It returns ok=false with a sanitized reason when the tool is no longer available to the agent, or when the
// fresh deterministic decision is stricter than the one the action was approved under — a revoked tool or a policy
// hardened to deny during the approval wait must not execute. It re-evaluates with the SAME authorization inputs as
// the original Propose, including AutoApprove, so an action still permitted under the current policy is not spuriously
// rejected (an auto-approved action must not be read as newly approval-required). The kill switch is checked
// separately by the caller, so it is not read here (KillSwitch stays false) to keep that path's distinct behavior.
func (g *Gateway) reauthorize(in ExecuteActionInput, action *Action, desc ToolDescriptor, found bool) (reason string, ok bool) {
	if !found {
		return "tool no longer available to agent", false
	}
	fresh := Authorize(AuthorizeInput{
		RunID:       in.RunID,
		Proposal:    in.Proposal,
		Descriptor:  desc,
		Found:       found,
		AutoApprove: in.AutoApprove,
		Global:      g.Policy.Global,
		Tenant:      g.Policy.Tenant,
		KillSwitch:  false,
	})
	// harden returns the stricter of the two; if the fresh decision is stricter than what was approved, harden changes
	// the recorded decision, which is the signal that authority narrowed since approval.
	if harden(action.PolicyDecision, fresh.Decision) != action.PolicyDecision {
		return "policy is now " + string(fresh.Decision) + ": " + fresh.Reason, false
	}
	return "", true
}

// refuse records an approved action as policy_rejected with a sanitized reason and reports a failed (definitively
// no-op) outcome. It is the terminal path for a reauthorization denial or a rebind mismatch: no external effect ran,
// so the run fails cleanly and the ledger stays truthful. It returns nil error because the refusal is a business
// outcome, not an infrastructure fault to be retried.
func (g *Gateway) refuse(ctx context.Context, in ExecuteActionInput, reason string) (ExecuteReport, error) {
	if _, err := g.Ledger.TransitionAction(ctx, ActionTransition{
		InstanceID: in.InstanceID, RunID: in.RunID, ToolCallID: in.Proposal.ToolCallID, Status: ActionPolicyRejected,
		Error: reason,
	}); err != nil {
		return ExecuteReport{}, fmt.Errorf("act: record refusal: %w", err)
	}
	return ExecuteReport{Outcome: OutcomeFailed}, nil
}

// rebindsToAction reports whether the proposal about to execute is identical to the approved action on the ledger
// (§11.2, §17.1): the tool, connector and pinned version must match the row, and the recomputed args hash must equal
// the approved hash. Any divergence means a different effect than the human approved and must not execute.
func rebindsToAction(p ToolProposal, desc ToolDescriptor, action *Action) (reason string, ok bool) {
	if p.Tool != action.Tool {
		return "tool changed", false
	}
	if p.Connector != action.Connector {
		return "connector changed", false
	}
	if desc.Version != action.ToolVersion {
		return "tool version changed", false
	}
	argsHash, err := HashCanonicalArgs(p.Args)
	if err != nil {
		return "arguments not canonicalizable", false
	}
	if argsHash != action.ArgsHash {
		return "arguments changed", false
	}
	return "", true
}

// reportRecorded maps an action's recorded state into the outcome a non-owning caller reports without re-executing: a
// finished write reports its recorded outcome; a still-executing write (owned by another live attempt) or anything
// else cannot be confirmed from here and is reported indeterminate, never re-driven (§12).
func reportRecorded(action *Action) ExecuteReport {
	switch action.Status {
	case ActionSucceeded, ActionVerified:
		return ExecuteReport{Outcome: OutcomeSucceeded, ExternalReference: action.ExternalReference}
	case ActionFailed:
		return ExecuteReport{Outcome: OutcomeFailed}
	default:
		return ExecuteReport{Outcome: OutcomeIndeterminate, ExternalReference: action.ExternalReference}
	}
}

// Verify confirms a succeeded action's effect when its tool supports verification (§16.1). It is best-effort: an
// executor error or an unconfirmed result leaves the action succeeded (not failed) and does not fail the run. Only a
// confirmed result advances the action to verified.
func (g *Gateway) Verify(ctx context.Context, in ExecuteActionInput) (VerifyResult, error) {
	action, err := g.Ledger.GetAction(ctx, in.InstanceID, in.RunID, in.Proposal.ToolCallID)
	if err != nil {
		return VerifyResult{}, err
	}
	if action.Status == ActionVerified {
		return VerifyResult{Confirmed: true}, nil // already verified: idempotent no-op
	}
	if action.Status != ActionSucceeded {
		return VerifyResult{}, nil // nothing to verify
	}

	// Build the redaction set for the verification result the same way the write path does: arg secrets plus the
	// connector token, so a detail or error that echoes a credential back is scrubbed (§17.1). The connector credential
	// must be GUARANTEED present before any server-controlled detail is persisted or returned; if it cannot be, the
	// detail is dropped rather than risk a token in the clear (C1). Verification stays best-effort: this drops the
	// detail, it never fails the run.
	secretValues, connectorRedactable := g.verifyRedactionSet(ctx, in)

	res, verr := g.Executor.Verify(ctx, VerifyRequest{
		RunID:             in.RunID,
		ToolCallID:        in.Proposal.ToolCallID,
		Tool:              in.Proposal.Tool,
		Connector:         in.Proposal.Connector,
		ExternalReference: action.ExternalReference,
		Trace:             TraceContext{RunID: in.RunID, AgentName: in.AgentName, TraceID: in.TraceID},
		Timeout:           g.timeout(),
	})
	if verr != nil {
		// Best-effort: a verification error does not undo a confirmed write. Log it redacted (it may quote a tool
		// response) and leave the action succeeded.
		if g.Logger != nil {
			g.Logger.Warn("act: verify failed", "run", in.RunID, "tool_call", in.Proposal.ToolCallID,
				"err", redactString(verr.Error(), secretValues))
		}
		return VerifyResult{}, nil
	}
	// The verification detail is server-controlled and is checkpointed by DBOS when it returns to the workflow, so it
	// must be redacted and length-capped on the returned VerifyResult too, not only on the ledger copy (§16.2, C1). If
	// the connector credential could not be guaranteed in the redaction set, drop the detail entirely.
	if connectorRedactable {
		res.Detail = truncate(redactString(res.Detail, secretValues), g.maxOutput())
	} else {
		res.Detail = ""
	}
	if !res.Confirmed {
		return res, nil
	}
	if _, err := g.Ledger.TransitionAction(ctx, ActionTransition{
		InstanceID: in.InstanceID, RunID: in.RunID, ToolCallID: in.Proposal.ToolCallID, Status: ActionVerified,
		RedactedVerification: map[string]any{"detail": res.Detail},
	}); err != nil {
		return VerifyResult{}, fmt.Errorf("act: record verified: %w", err)
	}
	return res, nil
}

// killDisabled reads the kill switch fail-closed: a nil switch is not engaged (no switch wired), and a switch that
// errors is treated as engaged so the gateway never proceeds on an unknown switch state (§17.2).
func (g *Gateway) killDisabled(ctx context.Context, scope KillScope) (bool, string) {
	if g.KillSwitch == nil {
		return false, ""
	}
	disabled, reason, err := g.KillSwitch.Disabled(ctx, scope)
	if err != nil {
		return true, "kill switch state unavailable"
	}
	return disabled, reason
}

// connectorRedactionSecret resolves the connector's bearer credential (a reference on the descriptor) and adds its
// plaintext to the redaction set, so a result echoing the token back is scrubbed (§17.1). On the write path it is
// STRICT: a declared credential that cannot be resolved is an error, because running the write would persist a result
// this gateway could not guarantee to redact (C1). A connector with no credential, or one that resolves empty, adds
// nothing and is not an error.
func (g *Gateway) connectorRedactionSecret(ctx context.Context, desc ToolDescriptor, secretValues []string) ([]string, error) {
	if desc.AuthSecret == "" {
		return secretValues, nil
	}
	if g.Secrets == nil {
		return nil, fmt.Errorf("act: connector credential %q declared but no secret resolver is configured", desc.AuthSecret)
	}
	token, err := g.Secrets.Resolve(ctx, desc.AuthSecret)
	if err != nil {
		return nil, fmt.Errorf("act: resolve connector credential for redaction: %w", err)
	}
	if token == "" {
		return secretValues, nil
	}
	return append(secretValues, token), nil
}

// verifyRedactionSet builds the redaction set for the best-effort verification path: the resolved argument secrets
// plus the connector bearer token. connectorRedactable reports whether the connector credential is guaranteed to be in
// the set (or that there is none); it is false when the descriptor could not be loaded or a declared credential could
// not be resolved, in which case the caller drops any server-controlled detail rather than risk persisting or
// returning a token in the clear (C1).
func (g *Gateway) verifyRedactionSet(ctx context.Context, in ExecuteActionInput) (secretValues []string, connectorRedactable bool) {
	secretValues = g.redactionSecrets(ctx, in.Proposal.Args)
	desc, found, err := g.Registry.Lookup(ctx, in.InstanceID, in.AgentName, in.Proposal.Tool)
	if err != nil || !found {
		// Without the descriptor we cannot know whether a connector credential exists, so we cannot guarantee it is
		// redactable: fail safe by declaring server-controlled detail non-redactable.
		return secretValues, false
	}
	if desc.AuthSecret == "" {
		return secretValues, true // no connector credential to leak
	}
	if g.Secrets == nil {
		return secretValues, false
	}
	token, rerr := g.Secrets.Resolve(ctx, desc.AuthSecret)
	if rerr != nil {
		return secretValues, false
	}
	if token == "" {
		return secretValues, true // resolves to empty: nothing to redact
	}
	return append(secretValues, token), true
}

// redactionSecrets resolves the arg secrets to build a redaction set, best-effort: on any resolution failure it
// returns what it has rather than blocking a best-effort path like verification. It never persists or returns the
// values for anything but scrubbing.
func (g *Gateway) redactionSecrets(ctx context.Context, args map[string]any) []string {
	_, values, err := resolveSecrets(ctx, g.Secrets, args)
	if err != nil {
		return nil
	}
	return values
}

// externalRefPattern bounds an acceptable external reference to a single line of printable ASCII, up to 256 bytes: a
// Jira key, a URL or a similar handle. Anything else — whitespace, control characters, non-ASCII or an over-long blob
// — is a malformed handle a hostile server may have crafted, and is rejected.
var externalRefPattern = regexp.MustCompile(`^[\x21-\x7e]{1,256}$`)

// sanitizeExternalReference redacts and validates the handle a server returned. It returns ok=false (and an empty
// reference) when the handle echoed a resolved secret or does not match the expected shape, so a leaked credential or
// arbitrary junk is never persisted as the action's reference (§17.1, C2 caps its length via the pattern).
func (g *Gateway) sanitizeExternalReference(ref string, secretValues []string) (string, bool) {
	if ref == "" {
		return "", true
	}
	if redactString(ref, secretValues) != ref {
		return "", false // the reference echoed a secret
	}
	if !externalRefPattern.MatchString(ref) {
		return "", false
	}
	return ref, true
}

func (g *Gateway) timeout() time.Duration {
	if g.Timeout > 0 {
		return g.Timeout
	}
	return defaultActionTimeout
}

func (g *Gateway) maxOutput() int {
	if g.MaxOutputBytes > 0 {
		return g.MaxOutputBytes
	}
	return defaultMaxOutput
}

// capOutput bounds the serialized size of a result map before it is persisted to the ledger, so a hostile or runaway
// executor cannot flood the row even if its own output cap was absent or bypassed (defense in depth, §16.2, D2). Over
// the cap the whole map is replaced by a small marker rather than truncated mid-structure, which would corrupt the
// JSON. A nil map stays nil.
func (g *Gateway) capOutput(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	if b, err := json.Marshal(m); err != nil || len(b) > g.maxOutput() {
		return map[string]any{"_dropped": "output exceeded size limit"}
	}
	return m
}

// validateArgs validates proposed arguments against a tool's fixed input schema (§16.2). A nil schema denies: an
// action whose arguments cannot be validated is not executed. A nil args map validates as an empty object.
func validateArgs(schema *jsonschema.Schema, args map[string]any) error {
	if schema == nil {
		return errors.New("no input schema")
	}
	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("resolve schema: %w", err)
	}
	var instance any = args
	if args == nil {
		instance = map[string]any{}
	}
	return resolved.Validate(instance)
}

// redactExecErr renders an executor error to a sanitized, secret-scrubbed message for the ledger. A nil error yields
// a generic label so the ledger still records why the action is not a success.
func redactExecErr(err error, secretValues []string) string {
	if err == nil {
		return "uncertain result"
	}
	return redactString(err.Error(), secretValues)
}

func killReasonOr(reason string) string {
	if reason == "" {
		return "kill switch engaged"
	}
	return reason
}

// MapToolRegistry is an in-memory ToolRegistry keyed by tool name, for wiring and tests. Production replaces it with
// a registry backed by the agent snapshot's compiled tool manifest and the pinned MCP tools/list (§16.3).
type MapToolRegistry struct {
	tools map[string]ToolDescriptor
}

var _ ToolRegistry = (*MapToolRegistry)(nil)

// NewMapToolRegistry builds a registry from the given descriptors, keyed by Name.
func NewMapToolRegistry(descriptors ...ToolDescriptor) *MapToolRegistry {
	tools := make(map[string]ToolDescriptor, len(descriptors))
	for _, d := range descriptors {
		tools[d.Name] = d
	}
	return &MapToolRegistry{tools: tools}
}

func (r *MapToolRegistry) Lookup(_ context.Context, _, _, tool string) (ToolDescriptor, bool, error) {
	d, ok := r.tools[tool]
	return d, ok, nil
}

// MapKillSwitch is a simple in-memory KillSwitch: a switch can be engaged for the whole platform, for an instance, or
// for a specific agent within an instance (§17.2). It is safe for concurrent use. A more elaborate, persisted switch
// replaces it later without changing the gateway.
type MapKillSwitch struct {
	mu        sync.RWMutex
	platform  string // non-empty reason means the whole platform is disabled
	instances map[string]string
	agents    map[string]string // key: instanceID + "\x00" + agentName
}

var _ KillSwitch = (*MapKillSwitch)(nil)

// NewMapKillSwitch returns an empty (nothing engaged) switch.
func NewMapKillSwitch() *MapKillSwitch {
	return &MapKillSwitch{instances: map[string]string{}, agents: map[string]string{}}
}

// DisablePlatform engages the platform-wide switch with a reason. DisableInstance and DisableAgent engage narrower
// scopes. An empty reason disengages the scope.
func (k *MapKillSwitch) DisablePlatform(reason string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.platform = reason
}

func (k *MapKillSwitch) DisableInstance(instanceID, reason string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if reason == "" {
		delete(k.instances, instanceID)
		return
	}
	k.instances[instanceID] = reason
}

func (k *MapKillSwitch) DisableAgent(instanceID, agentName, reason string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	key := instanceID + "\x00" + agentName
	if reason == "" {
		delete(k.agents, key)
		return
	}
	k.agents[key] = reason
}

func (k *MapKillSwitch) Disabled(_ context.Context, scope KillScope) (bool, string, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	if k.platform != "" {
		return true, k.platform, nil
	}
	if reason := k.instances[scope.InstanceID]; reason != "" {
		return true, reason, nil
	}
	if reason := k.agents[scope.InstanceID+"\x00"+scope.AgentName]; reason != "" {
		return true, reason, nil
	}
	return false, "", nil
}
