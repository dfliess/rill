package act_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// TestStartGivesARunAnOpaqueID: a run is addressed by a UUID, not by an ID composed of its instance, agent and key
// (ADR-0034). The composed ID carried slashes, so every REST route with the run in its path answered 404, and it put
// the instance and the caller's user ID in every URL and log line.
func TestStartGivesARunAnOpaqueID(t *testing.T) {
	store := newRunStore(t)
	e := newStoreExecutor(t, gatewayFlowRunner{}, store)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: "inst-opaque", AgentName: "triage", Prompt: "alerta", IdempotencyKey: "opaque-" + uuid.NewString(),
	})
	require.NoError(t, err)

	id, err := uuid.Parse(runID)
	require.NoError(t, err, "the run ID must be a UUID, got %q", runID)
	require.Equal(t, uuid.Version(7), id.Version())
}

// TestStartAttachesToTheRunOfItsKey: the run store, not the ID, decides what "the same run" is. A repeated Start
// with the same (instance, agent, key) returns the run that key already has, whatever the format of its ID: a run
// recorded before the IDs became opaque keeps its old ID, and a retry across the deploy must attach to it rather than
// start a duplicate (ADR-0034).
func TestStartAttachesToTheRunOfItsKey(t *testing.T) {
	store := newRunStore(t)
	e := newStoreExecutor(t, gatewayFlowRunner{}, store)
	ctx := t.Context()
	const inst = "inst-attach"

	t.Run("a run with an ID of the old format", func(t *testing.T) {
		// Any ID the current code would not produce: what matters is that the key, not the ID, finds the run.
		key := "legacy-" + uuid.NewString()
		legacyID := "legacy/" + key
		require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: legacyID, InstanceID: inst, AgentName: "triage", IdempotencyKey: key}))

		runID, err := e.Start(ctx, act.AgentRunInput{InstanceID: inst, AgentName: "triage", Prompt: "alerta", IdempotencyKey: key})
		require.NoError(t, err)
		require.Equal(t, legacyID, runID)
	})

	t.Run("a run started by this code", func(t *testing.T) {
		in := act.AgentRunInput{InstanceID: inst, AgentName: "triage", Prompt: "alerta", IdempotencyKey: "again-" + uuid.NewString()}
		first, err := e.Start(ctx, in)
		require.NoError(t, err)
		second, err := e.Start(ctx, in)
		require.NoError(t, err)
		require.Equal(t, first, second)
	})

	runs, err := store.ListRuns(ctx, act.ListRunsFilter{InstanceID: inst})
	require.NoError(t, err)
	require.Len(t, runs, 2, "each key must have exactly one run")
}

// TestStartRepairsARunWithoutWorkflow: before ADR-0032 a crash between recording a run and launching its workflow
// left a queued run with no workflow behind it. A retry with the same key must launch that run's workflow, not just
// return its ID, or the dispatcher would mark the event dispatched and the run would stay queued forever.
func TestStartRepairsARunWithoutWorkflow(t *testing.T) {
	store := newRunStore(t)
	e := newStoreExecutor(t, gatewayFlowRunner{}, store)
	ctx := t.Context()
	const inst = "inst-orphan"
	key := "orphan-" + uuid.NewString()
	orphan := newRunID()
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: orphan, InstanceID: inst, AgentName: "triage", IdempotencyKey: key}))

	runID, err := e.Start(ctx, act.AgentRunInput{InstanceID: inst, AgentName: "triage", Prompt: "alerta", IdempotencyKey: key})
	require.NoError(t, err)
	require.Equal(t, orphan, runID)

	require.Eventually(t, func() bool {
		run, gErr := store.GetRun(ctx, inst, orphan)
		return gErr == nil && run.Status != act.RunStatusQueued
	}, 20*time.Second, 50*time.Millisecond, "the retry must launch the orphan's workflow")
}

// TestStartKeepsTenantsAndAgentsApart: the key deduplicates only within one instance and one agent. Two tenants, or
// two agents, that reuse a key get separate runs, or a Start could attach to another tenant's run and an approval land
// on the wrong one.
func TestStartKeepsTenantsAndAgentsApart(t *testing.T) {
	store := newRunStore(t)
	e := newStoreExecutor(t, gatewayFlowRunner{}, store)
	key := "shared-" + uuid.NewString()
	start := func(inst, agent string) string {
		runID, err := e.Start(t.Context(), act.AgentRunInput{InstanceID: inst, AgentName: agent, Prompt: "alerta", IdempotencyKey: key})
		require.NoError(t, err)
		return runID
	}

	a := start("inst-a", "triage")
	require.NotEqual(t, a, start("inst-b", "triage"))
	require.NotEqual(t, a, start("inst-a", "cobranza"))
}

// TestApprovalIDIsOpaqueAndStable: an approval is addressed by a UUID derived from its run and tool call. It must be
// the same every time, because the workflow derives it again after the decision to learn who signed; and it must be
// safe in a URL whatever the model provider puts in the tool call ID.
func TestApprovalIDIsOpaqueAndStable(t *testing.T) {
	runID := uuid.Must(uuid.NewV7()).String()
	id := act.ApprovalIDForToolCall(runID, "tool/42")

	parsed, err := uuid.Parse(id)
	require.NoError(t, err, "the approval ID must be a UUID, got %q", id)
	require.Equal(t, uuid.Version(5), parsed.Version())
	require.Equal(t, id, act.ApprovalIDForToolCall(runID, "tool/42"))
	require.NotEqual(t, id, act.ApprovalIDForToolCall(runID, "tool/43"))
	require.NotEqual(t, id, act.ApprovalIDForToolCall(uuid.Must(uuid.NewV7()).String(), "tool/42"))
}

// TestOneApprovalPerToolCall: the store refuses a second approval for the same tool call of a run. Deriving the ID
// from the pair made that true by construction; the constraint states it in the schema.
func TestOneApprovalPerToolCall(t *testing.T) {
	store := newRunStore(t)
	ctx := t.Context()
	const inst = "inst-approval-unique"
	runID := uuid.Must(uuid.NewV7()).String()
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: runID, InstanceID: inst, AgentName: "triage", IdempotencyKey: "k"}))

	approval := func(id string) act.NewApproval {
		return act.NewApproval{
			ApprovalID: id, RunID: runID, InstanceID: inst, ToolCallID: "call-1",
			ArgsHash: act.HashArgs("x"), CanonicalArgs: "x", Proposal: "x", RequestedBy: "user:alice",
		}
	}
	require.NoError(t, store.CreateApproval(ctx, approval(uuid.NewString())))
	require.Error(t, store.CreateApproval(ctx, approval(uuid.NewString())))
}
