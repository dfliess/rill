package act

// PolicyDecision is the deterministic ruling on a proposed action tool call (§6.3 of
// docs/propuesta-diseno-act-rill-agents-go-dbos-2026-07-14.md). It is produced by the policy engine, never by the
// model: the model may propose a tool and arguments, but whether that call is allowed, needs a human, or is refused
// is decided here.
type PolicyDecision string

const (
	// PolicyAllow lets the action auto-execute without a human decision. It is reached by an explicitly auto-approved
	// tool (a connector's approval.auto / auto_approve posture) or a trusted read-only tool. The default for a write is
	// still approval; allow is opt-in per connector, and it never implies auto-retry (an uncertain result still parks at
	// indeterminate for a human).
	PolicyAllow PolicyDecision = "allow"
	// PolicyApprovalRequired suspends the run until a human approves the exact proposed arguments (§11.2).
	PolicyApprovalRequired PolicyDecision = "approval_required"
	// PolicyDeny refuses the action outright; no external effect runs and the run ends without touching the outside
	// world. It is the fail-closed default for anything ambiguous (§6.5).
	PolicyDeny PolicyDecision = "deny"
)

// strictness orders decisions from most permissive to most restrictive. Hardening a decision means moving up this
// order; a policy layer may only ever make a decision stricter, never looser (§8.1: a global or tenant policy may
// harden the agent's rule, never relax it). An unrecognized decision string sorts as deny, so a corrupt policy value
// fails closed rather than granting access.
func (d PolicyDecision) strictness() int {
	switch d {
	case PolicyAllow:
		return 0
	case PolicyApprovalRequired:
		return 1
	case PolicyDeny:
		return 2
	default:
		return 2
	}
}

// valid reports whether d is one of the three recognized decisions. An unrecognized value is never trusted as an
// override.
func (d PolicyDecision) valid() bool {
	switch d {
	case PolicyAllow, PolicyApprovalRequired, PolicyDeny:
		return true
	default:
		return false
	}
}

// harden returns the stricter of two decisions. A looser override leaves the base decision unchanged, which is what
// prevents any policy layer from relaxing the agent's own rule. An unrecognized override is canonicalized to deny
// before comparison, so a corrupt policy value fails closed rather than being propagated verbatim.
func harden(base, override PolicyDecision) PolicyDecision {
	if !override.valid() {
		override = PolicyDeny
	}
	if override.strictness() > base.strictness() {
		return override
	}
	return base
}

// ToolClass is a tool's idempotency classification (§12). It governs whether an action whose external result is
// uncertain may be safely retried, so it must be derived from the tool's contract — never from an untrusted MCP
// annotation (§16.3). The zero value is the unclassified class, which is treated fail-closed.
type ToolClass string

const (
	// ClassIdempotentNative names a tool that accepts an external idempotency key, so re-issuing the same call is
	// safe: the target system dedups it.
	ClassIdempotentNative ToolClass = "idempotent_native"
	// ClassIdempotentLookup names a tool with no native key but a cheap pre-check: a reference can be looked up
	// before creating, so a retry can detect the prior effect instead of duplicating it.
	ClassIdempotentLookup ToolClass = "idempotent_lookup"
	// ClassCompensable names a tool whose effect has a safe inverse, so a duplicate can be compensated.
	ClassCompensable ToolClass = "compensable"
	// ClassNonIdempotent names a tool with no dedup, lookup or compensation story: a call whose result is uncertain
	// cannot be retried safely and must go to indeterminate for human resolution (§12).
	ClassNonIdempotent ToolClass = "non_idempotent"
	// ClassUnknown is the zero value: a tool whose class the gateway could not establish. It is never auto-allowed and
	// never auto-retried.
	ClassUnknown ToolClass = ""
)

// known reports whether the class is one of the four defined classifications. An unknown class fails closed.
func (c ToolClass) known() bool {
	switch c {
	case ClassIdempotentNative, ClassIdempotentLookup, ClassCompensable, ClassNonIdempotent:
		return true
	default:
		return false
	}
}

// RetriableOnUncertain reports whether an action whose result could not be confirmed may be re-executed with the same
// idempotency key. It is DISABLED in v1: it always returns false, so an uncertain write of ANY class goes to
// indeterminate for human resolution rather than being auto-retried (§12).
//
// TODO(act, phase 2): re-enable class-based retry (idempotent_native, idempotent_lookup) only once two prerequisites
// hold, both missing today: (a) MCPClient.CallTool must propagate the idempotency key to the target system so an
// idempotent_native re-issue actually dedups there rather than duplicating the write; and (b) a real preflight lookup
// must exist for idempotent_lookup so a retry can detect the prior effect before creating a second one. Until both
// land, no class is safe to auto-retry and this must stay false.
func (c ToolClass) RetriableOnUncertain() bool {
	return false
}

// PolicyInput is the complete, self-contained input to a policy decision. It is deliberately a plain value with no
// I/O: every field is resolved by the gateway before evaluation (allowlist intersected with claims, kill switch read
// from its store, classification from the tool manifest), so Evaluate is a pure function that is trivially testable
// and cannot be influenced by a prompt.
type PolicyInput struct {
	Tool      string
	Connector string
	// Allowed is the effective set of tool names the agent may call: the snapshot allowlist intersected with what the
	// initiating actor's claims permit. A proposed tool outside it is denied — the model can never reach a tool the
	// agent did not declare or the actor cannot use.
	Allowed map[string]bool
	// AutoApprove is the agent's approval.automatic set: tools whose base decision is allow rather than approval.
	AutoApprove map[string]bool
	// Class is the tool's idempotency classification (§12). An unclassified tool can never be auto-allowed.
	Class ToolClass
	// ReadOnlyHint is the server-announced read-only annotation for the tool; TrustReadOnly is whether the connector
	// opted to trust that hint (off by default). An untrusted hint is ignored: MCP annotations are not a security
	// boundary (§16.3).
	ReadOnlyHint  bool
	TrustReadOnly bool
	// KillSwitch, when true, denies every action: the platform, tenant or agent switch is engaged (§17.2).
	KillSwitch bool
}

// PolicyResult is the outcome of a policy evaluation: the decision plus a short, non-secret reason for the audit
// trail and the approval inbox.
type PolicyResult struct {
	Decision PolicyDecision
	Reason   string
}

// Evaluate applies the deterministic action policy (§6.3, §6.5). It is pure and fail-closed: every ambiguous or
// unrecognized condition resolves toward deny or approval, never toward an implicit allow. The order encodes the
// precedence: an engaged kill switch or an out-of-allowlist tool denies before anything else is considered, an
// unclassified tool can never auto-execute, and global then tenant policy may only tighten the result.
func Evaluate(in PolicyInput) PolicyResult {
	// The kill switch denies unconditionally: an engaged switch is the operator's stop, and it must win over any
	// auto-approval or hint (§17.2).
	if in.KillSwitch {
		return PolicyResult{Decision: PolicyDeny, Reason: "kill switch engaged"}
	}
	// A tool the agent did not declare — or that the initiating actor's claims do not permit — is unreachable. Denying
	// here (rather than falling through to approval) is the core fail-closed rule: an unknown tool is never presented
	// to a human as if it were a legitimate proposal (§6.5).
	if !in.Allowed[in.Tool] {
		return PolicyResult{Decision: PolicyDeny, Reason: "tool not in agent allowlist"}
	}

	// Base decision: a write requires approval unless the agent explicitly auto-approved the tool, or the tool is a
	// trusted read-only tool. Nothing else earns an allow.
	decision := PolicyApprovalRequired
	reason := "write requires approval"
	autoApproved := in.AutoApprove[in.Tool]
	switch {
	case autoApproved:
		decision, reason = PolicyAllow, "auto-approved by agent policy"
	case in.ReadOnlyHint && in.TrustReadOnly:
		decision, reason = PolicyAllow, "trusted read-only tool"
	}

	// An unclassified tool cannot auto-execute UNLESS it was explicitly auto-approved. The class governs auto-RETRY of an
	// uncertain result, not whether a human must gate the single execution: retry stays disabled in v1 for every class
	// (RetriableOnUncertain), so an auto-approved write still runs exactly once and an uncertain result parks at
	// indeterminate for a human rather than being re-issued. An explicit auto-approve is the operator's deliberate choice
	// to run this tool unsupervised, so it is exempt; any OTHER allow of an unclassified tool still falls back to
	// approval, keeping the fail-closed default for a path that did not opt in (§6.5, §12). A generic MCP tool is always
	// unclassified, so this exemption is what makes a connector's approval.auto / auto_approve posture take effect.
	if decision == PolicyAllow && !autoApproved && !in.Class.known() {
		decision, reason = PolicyApprovalRequired, "unclassified tool cannot auto-execute"
	}

	return PolicyResult{Decision: decision, Reason: reason}
}
