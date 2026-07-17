package act_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// baseInput is a minimal allowed, classified write: a tool the agent may call, with a known idempotency class, no
// auto-approval and no policy overrides. Individual tests tweak one field to isolate a rule.
func baseInput() act.PolicyInput {
	return act.PolicyInput{
		Tool:        "jira.create_issue",
		Connector:   "jira_ops",
		Allowed:     map[string]bool{"jira.create_issue": true},
		AutoApprove: map[string]bool{},
		Class:       act.ClassIdempotentNative,
	}
}

// TestPolicyDefaultWriteRequiresApproval: a known, allowed write that is not auto-approved needs a human.
func TestPolicyDefaultWriteRequiresApproval(t *testing.T) {
	got := act.Evaluate(baseInput())
	require.Equal(t, act.PolicyApprovalRequired, got.Decision)
}

// TestPolicyAutoApproveAllows: an explicitly auto-approved tool may skip human approval.
func TestPolicyAutoApproveAllows(t *testing.T) {
	in := baseInput()
	in.AutoApprove = map[string]bool{"jira.create_issue": true}
	got := act.Evaluate(in)
	require.Equal(t, act.PolicyAllow, got.Decision)
}

// TestPolicyUnknownToolDenies is the core fail-closed rule (§6.5): a tool outside the agent allowlist is denied, never
// surfaced to a human as a legitimate proposal.
func TestPolicyUnknownToolDenies(t *testing.T) {
	in := baseInput()
	in.Tool = "jira.delete_project" // not in Allowed
	got := act.Evaluate(in)
	require.Equal(t, act.PolicyDeny, got.Decision)
	require.Contains(t, got.Reason, "allowlist")
}

// TestPolicyKillSwitchDenies: an engaged kill switch denies unconditionally, even an otherwise auto-approved tool.
func TestPolicyKillSwitchDenies(t *testing.T) {
	in := baseInput()
	in.AutoApprove = map[string]bool{"jira.create_issue": true} // would be allow…
	in.KillSwitch = true                                        // …but the switch wins.
	got := act.Evaluate(in)
	require.Equal(t, act.PolicyDeny, got.Decision)
	require.Contains(t, got.Reason, "kill switch")
}

// TestPolicyAutoApproveAllowsUnclassified: an explicitly auto-approved tool auto-executes even when its idempotency
// class is unknown. The class only governs auto-RETRY of an uncertain result (disabled in v1), not whether a human
// gates the single execution; a generic MCP tool is always unclassified, so this is the path every auto-approved
// connector tool actually takes (§12). Without the exemption a connector's approval.auto / auto_approve posture would
// silently no-op.
func TestPolicyAutoApproveAllowsUnclassified(t *testing.T) {
	in := baseInput()
	in.AutoApprove = map[string]bool{"jira.create_issue": true}
	in.Class = act.ClassUnknown
	got := act.Evaluate(in)
	require.Equal(t, act.PolicyAllow, got.Decision)
	require.Contains(t, got.Reason, "auto-approved")
}

// TestPolicyUnclassifiedNonAutoAllowNeedsApproval: an allow that was NOT an explicit auto-approve (here a trusted
// read-only hint) still falls back to approval when the class is unknown. The fail-closed default holds for any path
// that did not deliberately opt a tool into unsupervised execution (§6.5).
func TestPolicyUnclassifiedNonAutoAllowNeedsApproval(t *testing.T) {
	in := baseInput()
	in.ReadOnlyHint = true
	in.TrustReadOnly = true
	in.Class = act.ClassUnknown
	got := act.Evaluate(in)
	require.Equal(t, act.PolicyApprovalRequired, got.Decision)
	require.Contains(t, got.Reason, "unclassified")
}

// TestPolicyTrustedReadOnlyAllows: a read-only tool the connector trusts may auto-execute; an untrusted hint is
// ignored and the write still needs approval (§16.3).
func TestPolicyTrustedReadOnlyAllows(t *testing.T) {
	in := baseInput()
	in.ReadOnlyHint = true

	in.TrustReadOnly = false
	require.Equal(t, act.PolicyApprovalRequired, act.Evaluate(in).Decision, "untrusted hint is ignored")

	in.TrustReadOnly = true
	require.Equal(t, act.PolicyAllow, act.Evaluate(in).Decision, "trusted read-only tool auto-executes")
}

// TestPolicyGlobalAndTenantOnlyHarden: an override may tighten a decision but never loosen it (§8.1).
func TestPolicyGlobalAndTenantOnlyHarden(t *testing.T) {
	// Global hardens an auto-approved allow up to deny.
	in := baseInput()
	in.AutoApprove = map[string]bool{"jira.create_issue": true}
	in.GlobalPolicy = map[string]act.PolicyDecision{"jira.create_issue": act.PolicyDeny}
	require.Equal(t, act.PolicyDeny, act.Evaluate(in).Decision)

	// Tenant hardens a base approval up to deny.
	in = baseInput()
	in.TenantPolicy = map[string]act.PolicyDecision{"jira.create_issue": act.PolicyDeny}
	require.Equal(t, act.PolicyDeny, act.Evaluate(in).Decision)

	// A looser override is ignored: a base approval stays approval even if global says allow.
	in = baseInput()
	in.GlobalPolicy = map[string]act.PolicyDecision{"jira.create_issue": act.PolicyAllow}
	require.Equal(t, act.PolicyApprovalRequired, act.Evaluate(in).Decision)
}

// TestPolicyCorruptOverrideFailsClosed: an unrecognized policy value hardens to deny rather than granting access.
func TestPolicyCorruptOverrideFailsClosed(t *testing.T) {
	in := baseInput()
	in.AutoApprove = map[string]bool{"jira.create_issue": true}
	in.GlobalPolicy = map[string]act.PolicyDecision{"jira.create_issue": act.PolicyDecision("garbage")}
	require.Equal(t, act.PolicyDeny, act.Evaluate(in).Decision)
}

// TestToolClassRetriable pins the v1 retry classification (§12): auto-retry on an uncertain result is disabled for
// EVERY class, so no class is ever re-driven automatically. The dedup story a class describes is not yet safe to act
// on because CallTool does not propagate the idempotency key and no preflight lookup exists (see RetriableOnUncertain).
func TestToolClassRetriable(t *testing.T) {
	require.False(t, act.ClassIdempotentNative.RetriableOnUncertain())
	require.False(t, act.ClassIdempotentLookup.RetriableOnUncertain())
	require.False(t, act.ClassCompensable.RetriableOnUncertain())
	require.False(t, act.ClassNonIdempotent.RetriableOnUncertain())
	require.False(t, act.ClassUnknown.RetriableOnUncertain())
}
