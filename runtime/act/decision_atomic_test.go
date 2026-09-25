package act_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// TestDecisionAndDeliveryCommitTogether: recording a decision and delivering it to the run are one transaction
// (kairos-cloud#194). If anything fails after the Send but before the commit, neither happened: the approval is still
// pending and the run never receives the message. Before, the claim committed on its own, so a crash between the two
// left a decided approval and a run waiting forever.
func TestDecisionAndDeliveryCommitTogether(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-194"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})
	ctx := t.Context()

	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste subio +40%.",
		IdempotencyKey: "atomic-decision-" + uuid.NewString(), Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)
	approvalID := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		a, gErr := store.GetApproval(ctx, instanceID, approvalID)
		return gErr == nil && a.Status == act.ApprovalStatusPending
	}, 15*time.Second, 50*time.Millisecond)

	errCrash := errors.New("process died between the send and the commit")
	_, err = store.ResolveApprovalAndDeliver(ctx, instanceID, approvalID, act.ApprovalStatusApproved, "user:bob",
		func(ctx context.Context, tx pgx.Tx) error {
			if err := e.ResumeInTx(ctx, tx, runID, act.ApprovalApproved, "call-1"); err != nil {
				return err
			}
			return errCrash
		})
	require.ErrorIs(t, err, errCrash)

	a, err := store.GetApproval(ctx, instanceID, approvalID)
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusPending, a.Status, "a rolled-back delivery leaves the decision unrecorded")
	require.Never(t, func() bool { return exec.executeCount() > 0 }, 2*time.Second, 100*time.Millisecond,
		"the run must not receive a decision whose transaction rolled back")

	_, err = store.ResolveApprovalAndDeliver(ctx, instanceID, approvalID, act.ApprovalStatusApproved, "user:bob",
		func(ctx context.Context, tx pgx.Tx) error {
			return e.ResumeInTx(ctx, tx, runID, act.ApprovalApproved, "call-1")
		})
	require.NoError(t, err)
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, exec.executeCount())
}
