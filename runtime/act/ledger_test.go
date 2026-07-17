package act_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// proposeAction seeds a proposed action on the ledger and returns its identifiers, so a test can focus on the
// transition it exercises rather than the boilerplate of a proposal.
func proposeAction(t *testing.T, ledger act.ActionLedger, instanceID, runID, toolCallID string) {
	t.Helper()
	argsHash := act.HashArgs("args-" + toolCallID)
	require.NoError(t, ledger.ProposeAction(t.Context(), act.NewAction{
		ToolCallID:     toolCallID,
		RunID:          runID,
		InstanceID:     instanceID,
		AgentName:      "triage",
		Tool:           "jira.create_issue",
		Connector:      "jira_ops",
		Class:          act.ClassIdempotentNative,
		ArgsHash:       argsHash,
		IdempotencyKey: act.DeriveIdempotencyKey(runID, toolCallID, argsHash),
		RedactedArgs:   map[string]any{"summary": "coste alto"},
		PolicyDecision: act.PolicyApprovalRequired,
		Proposal:       "crear ticket",
		RequestedBy:    "user:alice",
	}))
}

func transition(t *testing.T, ledger act.ActionLedger, instanceID, runID, toolCallID string, status act.ActionStatus) (*act.Action, error) {
	t.Helper()
	return ledger.TransitionAction(t.Context(), act.ActionTransition{
		InstanceID: instanceID, RunID: runID, ToolCallID: toolCallID, Status: status,
	})
}

// TestLedgerHappyPath walks the full action state machine (§11.2) proposed → approval_pending → approved → executing
// → succeeded → verified, and checks the terminal detail (external reference, executed/verified timestamps) lands.
func TestLedgerHappyPath(t *testing.T) {
	ledger := newRunStore(t)
	ctx := t.Context()
	const inst, run, call = "inst-ledger", "run-1", "call-1"
	proposeAction(t, ledger, inst, run, call)

	got, err := ledger.GetAction(ctx, inst, run, call)
	require.NoError(t, err)
	require.Equal(t, act.ActionProposed, got.Status)
	require.Nil(t, got.ExecutedOn)

	for _, s := range []act.ActionStatus{act.ActionApprovalPending, act.ActionApproved, act.ActionExecuting} {
		_, err := transition(t, ledger, inst, run, call, s)
		require.NoError(t, err)
	}

	// Succeed with an external reference and a redacted result.
	_, err = ledger.TransitionAction(ctx, act.ActionTransition{
		InstanceID: inst, RunID: run, ToolCallID: call, Status: act.ActionSucceeded,
		ExternalReference: "PROJ-123", RedactedResult: map[string]any{"key": "PROJ-123"},
	})
	require.NoError(t, err)

	verified, err := transition(t, ledger, inst, run, call, act.ActionVerified)
	require.NoError(t, err)
	require.Equal(t, act.ActionVerified, verified.Status)
	require.Equal(t, "PROJ-123", verified.ExternalReference)
	require.Equal(t, "PROJ-123", verified.RedactedResult["key"])
	require.NotNil(t, verified.ExecutedOn, "executed_on latches when the action enters executing")
	require.NotNil(t, verified.VerifiedOn)
}

// TestLedgerProposeIdempotent verifies a replayed proposal neither errors nor overwrites the action: re-proposing an
// action already advanced past proposed leaves its status intact (ON CONFLICT DO NOTHING). This is what makes the
// gateway's Propose step safe under at-least-once (§12).
func TestLedgerProposeIdempotent(t *testing.T) {
	ledger := newRunStore(t)
	ctx := t.Context()
	const inst, run, call = "inst-idem", "run-1", "call-1"
	proposeAction(t, ledger, inst, run, call)
	_, err := transition(t, ledger, inst, run, call, act.ActionApprovalPending)
	require.NoError(t, err)

	// Replay the proposal: the row already exists and is past proposed, so the insert is a no-op and the status holds.
	proposeAction(t, ledger, inst, run, call)
	got, err := ledger.GetAction(ctx, inst, run, call)
	require.NoError(t, err)
	require.Equal(t, act.ActionApprovalPending, got.Status)
}

// TestLedgerRejectsIllegalTransition verifies the state machine refuses an edge it does not permit (proposed straight
// to executing), so a bug or a hostile caller cannot skip the approval gate at the ledger level.
func TestLedgerRejectsIllegalTransition(t *testing.T) {
	ledger := newRunStore(t)
	const inst, run, call = "inst-illegal", "run-1", "call-1"
	proposeAction(t, ledger, inst, run, call)

	_, err := transition(t, ledger, inst, run, call, act.ActionExecuting)
	require.ErrorIs(t, err, act.ErrActionTransitionInvalid)
}

// TestLedgerSameStatusIsNoop verifies re-applying the current status is an idempotent no-op — a replayed step that
// re-marks executing after a crash must not error.
func TestLedgerSameStatusIsNoop(t *testing.T) {
	ledger := newRunStore(t)
	const inst, run, call = "inst-noop", "run-1", "call-1"
	proposeAction(t, ledger, inst, run, call)
	for _, s := range []act.ActionStatus{act.ActionApprovalPending, act.ActionApproved, act.ActionExecuting} {
		_, err := transition(t, ledger, inst, run, call, s)
		require.NoError(t, err)
	}
	got, err := transition(t, ledger, inst, run, call, act.ActionExecuting) // replay same edge
	require.NoError(t, err)
	require.Equal(t, act.ActionExecuting, got.Status)
}

// TestLedgerTerminalGuard verifies a terminal action never transitions again: a stray edge after policy_rejected (or
// after indeterminate) is a no-op that leaves the recorded outcome intact. This keeps a replayed step from regressing
// a finished action.
func TestLedgerTerminalGuard(t *testing.T) {
	ledger := newRunStore(t)
	const inst, run, call = "inst-terminal", "run-1", "call-1"
	proposeAction(t, ledger, inst, run, call)
	rejected, err := transition(t, ledger, inst, run, call, act.ActionPolicyRejected)
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, rejected.Status)

	// A stray attempt to advance a terminal action is a no-op, not an error, and does not change the status.
	got, err := transition(t, ledger, inst, run, call, act.ActionApproved)
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, got.Status)
}

// TestLedgerScopingIsolatesTenants verifies an action is unreadable and immutable across instances (§15.2).
func TestLedgerScopingIsolatesTenants(t *testing.T) {
	ledger := newRunStore(t)
	ctx := t.Context()
	const instA, instB, run, call = "inst-a", "inst-b", "run-1", "call-1"
	proposeAction(t, ledger, instA, run, call)

	_, err := ledger.GetAction(ctx, instB, run, call)
	require.ErrorIs(t, err, act.ErrActionNotFound)

	_, err = transition(t, ledger, instB, run, call, act.ActionApprovalPending)
	require.ErrorIs(t, err, act.ErrActionNotFound)

	actions, err := ledger.ListActions(ctx, instA, run)
	require.NoError(t, err)
	require.Len(t, actions, 1)
}
