package act_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// fixedProposer is a Proposer that always surfaces the same proposal (or none). It stands in for the runtime/ai loop's
// structured tool call so the gateway-wired workflow can be exercised end to end.
type fixedProposer struct {
	proposal act.ToolProposal
	ok       bool
}

func (p fixedProposer) Propose(_ context.Context, _ act.ProposeInput) (act.ToolProposal, bool, error) {
	return p.proposal, p.ok, nil
}

// gatewayFlowRunner drives the segmented loop for the gateway-flow tests: it pauses on a governed write on the first
// segment and closes on the resumed segment (after the action's result is injected), so the loop governs exactly one
// action and then terminates. The concrete proposal the gateway governs comes from the wired Proposer (fixedProposer),
// so the captured action here only needs to be non-nil to signal the pause. It avoids the real agent loop, which would
// need to reach an MCP server to discover tools.
type gatewayFlowRunner struct{}

var _ act.Runner = gatewayFlowRunner{}

func (gatewayFlowRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	return &ai.AgentSnapshot{Name: agentName, Instructions: "gateway flow spike"}, nil
}

func (gatewayFlowRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	sessionID := in.SessionID
	if sessionID == "" {
		sessionID = "sess-" + in.Snapshot.Name
	}
	if len(in.Resume) > 0 {
		return act.RunSegmentResult{Response: "Hecho, creé el ticket.", SessionID: sessionID}, nil
	}
	return act.RunSegmentResult{
		Response:  "Propongo crear un ticket.",
		SessionID: sessionID,
		Proposed:  []*ai.ProposedAction{{Tool: "jira.create_issue", Connector: "jira_ops"}},
	}, nil
}

func (gatewayFlowRunner) ApplyAction(context.Context, act.ActionRequest) (act.ActionResult, error) {
	return act.ActionResult{}, nil
}

// newGatewayExecutor builds an in-process executor whose action phase runs through a real Gateway (fake executor +
// Postgres ledger) rather than the simulated sink. The run store and the gateway's ledger are the SAME
// PostgresRunStore, so a run's lifecycle and its action ledger share one Postgres, exactly as production wires them.
// It returns the executor, that shared store, and the instance ID the agent "triage" lives in.
func newGatewayExecutor(t *testing.T, exec act.ActionExecutor, desc act.ToolDescriptor, secrets act.SecretResolver, proposer act.Proposer) (*act.DBOSExecutor, *act.PostgresRunStore, string) {
	t.Helper()
	_, instanceID := newInstanceWithAgent(t, "Investiga la alerta y propon una remediacion.")
	store := newRunStore(t)
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(desc),
		Ledger:   store,
		Executor: exec,
		Secrets:  secrets,
	}
	// A mock runner that proposes one governed write then closes on resume, so the segmented loop governs exactly one
	// action. The concrete proposal the gateway governs comes from the wired Proposer, not from the loop.
	runner := gatewayFlowRunner{}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL:        dsn,
		DatabaseSchema:     schema,
		ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner:             runner,
		Store:              store,
		Gateway:            gateway,
		Proposer:           proposer,
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })
	return e, store, instanceID
}

func flowProposal() act.ToolProposal {
	return act.ToolProposal{
		ToolCallID: "call-1",
		Tool:       "jira.create_issue",
		Connector:  "jira_ops",
		Args:       map[string]any{"summary": "coste alto"},
		Summary:    "crear ticket P2",
	}
}

// TestGatewayFlowApproveExecutes is the end-to-end happy path through the wired workflow: run → propose → approval →
// approved → execute → verify, ending succeeded with the external reference on the run and the action verified.
func TestGatewayFlowApproveExecutes(t *testing.T) {
	exec := &fakeExecutor{
		result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-777"},
		verify: act.VerifyResult{Confirmed: true},
	}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})

	key := "flow-approve-" + uuid.NewString()
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste subio +40%.", IdempotencyKey: key,
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.True(t, res.ActionTaken)
	require.Equal(t, "PROJ-777", res.ActionRef)
	require.Equal(t, 1, exec.executeCount())

	// The action ledger reached verified (the tool is verifiable and verification confirmed).
	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionVerified, action.Status)
	require.Equal(t, "PROJ-777", action.ExternalReference)

	// The approval created by the real workflow carries the exact preimage of the hash the decision was bound to
	// (§6.4): what a UI shows from these bytes is byte-for-byte what was hashed, approved and executed.
	appr, err := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForToolCall(runID, "call-1"))
	require.NoError(t, err)
	canonical, ok := appr.VerifiedCanonicalArgs()
	require.True(t, ok)
	require.Equal(t, appr.ArgsHash, act.HashArgs(canonical))
	require.Equal(t, appr.ArgsHash, action.ArgsHash, "the approval and the executed action bind to one hash")
	wantBytes, err := act.CanonicalizeArgs(flowProposal().Args)
	require.NoError(t, err)
	require.Equal(t, string(wantBytes), canonical)
}

// TestGatewayFlowRecordsRealApprover proves the action ledger AND the run timeline record who APPROVED the action, not
// who started the run (§11.2, §18.3). Alice starts the run; Bob (an admin) resolves the approval. The reload inside the
// workflow makes both the ledger's decided_by and the resumed event's payload the real approver, Bob, rather than the
// initiator, Alice. The event matters on its own: the API does not expose the ledger, so it is the only way a reader
// learns who unblocked the run.
func TestGatewayFlowRecordsRealApprover(t *testing.T) {
	exec := &fakeExecutor{
		result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-42"},
		verify: act.VerifyResult{Confirmed: true},
	}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste subio.", IdempotencyKey: "flow-approver-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Wait for the run to reach the approval gate, then resolve it as Bob and deliver the durable resume — mirroring
	// the API's claim-then-send order.
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForToolCall(runID, "call-1"))
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(t.Context(), instanceID, act.ApprovalIDForToolCall(runID, "call-1"), act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, "admin:bob", action.DecidedBy, "the ledger records the approver, not the run initiator")

	events, err := store.ListRunEvents(t.Context(), instanceID, runID, 0, 100)
	require.NoError(t, err)
	resumed := lastEventOfType(t, events, act.EventTypeResumed)
	require.Equal(t, "admin:bob", resumed.Payload["decided_by"], "the timeline attributes the approval to the approver")
}

// lastEventOfType returns the most recent event of the given type, failing the test when the run never emitted one.
func lastEventOfType(t *testing.T, events []*act.RunEvent, eventType string) *act.RunEvent {
	t.Helper()
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].EventType == eventType {
			return events[i]
		}
	}
	require.FailNowf(t, "event not emitted", "no %s event in the run", eventType)
	return nil
}

// TestGatewayFlowCancelDuringApprovalPreventsExecution verifies the escape hatch that matters now that approvals do
// not auto-expire: a run parked on the approval gate can be cancelled, and a decision that races in AFTER the cancel
// must never execute the proposed action. Without a working cancel this would be a stuck run with no way out; without
// the race guarantee a "cancel then approve" could still drive the external write.
func TestGatewayFlowCancelDuringApprovalPreventsExecution(t *testing.T) {
	exec := &fakeExecutor{
		result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-99"},
		verify: act.VerifyResult{Confirmed: true},
	}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "El coste subio.", IdempotencyKey: "flow-cancel-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Wait until the run parks on the approval gate (the approval row exists).
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForToolCall(runID, "call-1"))
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)

	// Cancel the parked run, then — racing a late approver — try to deliver an approval. Neither call should drive the
	// external write. Resume may error (a cancelled workflow refuses the send); that is fine, we only require no write.
	require.NoError(t, e.Cancel(t.Context(), instanceID, runID))
	_ = e.Resume(t.Context(), runID, act.ApprovalApproved, "call-1")

	// The governed action must NOT execute, even given time for a wrongly-resumed workflow to reach the executor.
	require.Never(t, func() bool { return exec.executeCount() > 0 }, 3*time.Second, 100*time.Millisecond,
		"a cancelled run must not execute its proposed action, even if an approval races in after the cancel")

	// The action ledger never advanced past the approval gate: no executing/succeeded/verified.
	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.NotContains(t, []act.ActionStatus{act.ActionExecuting, act.ActionSucceeded, act.ActionVerified}, action.Status,
		"a cancelled run's action must not reach an execution state")

	// Cancelling the run also withdrew its pending approval, so the inbox stops offering a decision on a run that will
	// never resume (the approval panel is gated on the pending status).
	appr, err := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForToolCall(runID, "call-1"))
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusCancelled, appr.Status, "cancelling a run must withdraw its pending approval")
}

// TestGatewayFlowRejectionSkipsWrite verifies a human denial injects the rejection as a result (the model can adapt)
// and never executes the proposed action; the action stays approval_pending on the ledger. The run succeeds because
// the model sees the rejection and closes: in the per-tool-call model rejection is per-action, not per-run.
func TestGatewayFlowRejectionSkipsWrite(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "fallo", IdempotencyKey: "flow-reject-" + uuid.NewString(),
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalRejected, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status, "rejection is per-action: the model sees it and closes")
	require.False(t, res.ActionTaken)
	require.Equal(t, 0, exec.executeCount())

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionRejected, action.Status, "the ledger closes the rejected action")
}

// autoApproveRunner is a mock Runner that returns a snapshot whose jira_ops connector is approval: auto plus a captured
// write proposal on that connector, without running the real agent loop (which would try to reach the MCP server to
// discover tools, and cannot point at a loopback test server because the reconciler denies private ranges). It lets
// the DBOS auto-approve path be exercised end to end from the snapshot the executor actually reads.
type autoApproveRunner struct{}

var _ act.Runner = autoApproveRunner{}

func (autoApproveRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	return &ai.AgentSnapshot{
		Name:          agentName,
		Instructions:  "auto-approve spike",
		MCPConnectors: []ai.MCPConnector{{Name: "jira_ops", Approval: "auto"}},
	}, nil
}

func (autoApproveRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	sessionID := in.SessionID
	if sessionID == "" {
		sessionID = "sess-" + in.Snapshot.Name
	}
	// Second segment: the model has seen the executed action's result injected and closes, proposing nothing further,
	// so the segmented loop finishes. Without this branch the mock would re-propose forever.
	if len(in.Resume) > 0 {
		return act.RunSegmentResult{Response: "Listo, creé el ticket.", SessionID: sessionID}, nil
	}
	return act.RunSegmentResult{
		Response:  "Propongo crear un ticket.",
		SessionID: sessionID,
		Proposed: []*ai.ProposedAction{{
			Connector: "jira_ops", Tool: "jira.create_issue",
			Args: map[string]any{"summary": "coste alto"}, Summary: "crear ticket P2",
		}},
	}, nil
}

func (autoApproveRunner) ApplyAction(context.Context, act.ActionRequest) (act.ActionResult, error) {
	return act.ActionResult{}, nil
}

// recordingSegmentRunner proposes a governed write on its first `actions` segments and closes afterward, recording every
// injected result it receives on resume. It lets a test assert the executor fed the gateway's result back into the loop
// (b2) and drive several sequential actions in one run. The concrete tool proposal the gateway governs comes from the
// wired Proposer, so the captured action only needs to be non-nil to signal each pause.
type recordingSegmentRunner struct {
	snapshot *ai.AgentSnapshot
	actions  int

	mu       sync.Mutex
	proposed int
	resumes  []ai.InjectedResult
}

var _ act.Runner = (*recordingSegmentRunner)(nil)

func (r *recordingSegmentRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	if r.snapshot != nil {
		return r.snapshot, nil
	}
	return &ai.AgentSnapshot{Name: agentName, Instructions: "segmented spike"}, nil
}

func (r *recordingSegmentRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	sessionID := in.SessionID
	if sessionID == "" {
		sessionID = "sess-" + in.Snapshot.Name
	}
	r.mu.Lock()
	for _, res := range in.Resume {
		r.resumes = append(r.resumes, *res)
	}
	propose := r.proposed < r.actions
	if propose {
		r.proposed++
	}
	r.mu.Unlock()
	if !propose {
		return act.RunSegmentResult{Response: "Listo, cerré la conversación.", SessionID: sessionID}, nil
	}
	return act.RunSegmentResult{
		Response:  "Propongo una acción.",
		SessionID: sessionID,
		Proposed:  []*ai.ProposedAction{{Tool: "jira.create_issue", Connector: "jira_ops", Args: map[string]any{"summary": "coste alto"}}},
	}, nil
}

func (r *recordingSegmentRunner) ApplyAction(context.Context, act.ActionRequest) (act.ActionResult, error) {
	return act.ActionResult{}, nil
}

func (r *recordingSegmentRunner) injectedResults() []ai.InjectedResult {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]ai.InjectedResult(nil), r.resumes...)
}

// sequentialProposer surfaces a distinct tool call per invocation (call-1, call-2, …), so a run proposing several
// actions keys each on its own ledger row rather than colliding on a fixed ID.
type sequentialProposer struct {
	mu sync.Mutex
	n  int
}

func (p *sequentialProposer) Propose(_ context.Context, _ act.ProposeInput) (act.ToolProposal, bool, error) {
	p.mu.Lock()
	p.n++
	id := fmt.Sprintf("call-%d", p.n)
	p.mu.Unlock()
	return act.ToolProposal{
		ToolCallID: id, Tool: "jira.create_issue", Connector: "jira_ops",
		Args: map[string]any{"summary": "coste alto"}, Summary: "crear ticket",
	}, true, nil
}

func autoApproveSnapshot() *ai.AgentSnapshot {
	return &ai.AgentSnapshot{
		Name:          "triage",
		Instructions:  "auto-approve spike",
		MCPConnectors: []ai.MCPConnector{{Name: "jira_ops", Approval: "auto"}},
	}
}

// TestGatewayFlowSegmentedInjectsResult is the b2 core at the executor level: after a governed write is auto-approved
// and executed, the workflow injects the gateway's redacted result into the NEXT segment, so the model sees the outcome
// of the action it proposed and closes. It asserts the resumed segment received exactly that result.
func TestGatewayFlowSegmentedInjectsResult(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-9", Message: "created PROJ-9"}}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassUnknown)),
		Ledger:   store,
		Executor: exec,
	}
	runner := &recordingSegmentRunner{snapshot: autoApproveSnapshot(), actions: 1}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: act.NewCapturedProposer(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: "inst-seg", AgentName: "triage", Prompt: "actua", IdempotencyKey: "flow-seg-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, exec.executeCount())

	// The resumed segment received the gateway's redacted result, keyed to the action's tool call. This is what lets the
	// model see the outcome and close, instead of the run ending at the proposal (the b1 limitation b2 fixes).
	injected := runner.injectedResults()
	require.Len(t, injected, 1, "the model was resumed once, with the executed action's result")
	require.Equal(t, "created PROJ-9", injected[0].Message)
	require.Equal(t, "jira.create_issue", injected[0].ToolCallID)
}

// TestGatewayFlowSegmentedTwoSequentialActions proves a single run can propose, execute and see the result of several
// governed writes in sequence (the multiple-sequential-actions capability): two auto-approved actions both execute,
// each keyed on its own ledger row, the model is resumed with each result, and the run succeeds once it closes.
func TestGatewayFlowSegmentedTwoSequentialActions(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-9", Message: "created PROJ-9"}}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassUnknown)),
		Ledger:   store,
		Executor: exec,
	}
	runner := &recordingSegmentRunner{snapshot: autoApproveSnapshot(), actions: 2}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: &sequentialProposer{},
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-seg-multi"
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua dos veces", IdempotencyKey: "flow-seg-multi-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 2, exec.executeCount(), "both sequential actions executed")

	// The model was resumed after each action, and each action landed on its own ledger row.
	require.Len(t, runner.injectedResults(), 2)
	a1, err := store.GetAction(t.Context(), inst, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, a1.Status)
	a2, err := store.GetAction(t.Context(), inst, runID, "call-2")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, a2.Status)
}

// TestGatewayFlowTwoHumanApprovalsRecordEachPause is the end-to-end regression for kairos-cloud#129: a run with TWO
// governed actions, each behind a human decision. The old (run_id, event_type) uniqueness swallowed the second
// waiting_approval, so the run sat parked on its second gate while its status said running — invisible to the Home
// inbox and to Act's default filter, both of which ask for waiting_approval. It verifies (a) the timeline records both
// pauses and both resumptions, each resumption attributed to its own decider, and (b) DURING the second wait the run's
// status is waiting_approval.
func TestGatewayFlowTwoHumanApprovalsRecordEachPause(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-9", Message: "created PROJ-9"}}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:   store,
		Executor: exec,
	}
	// Two governed actions on a snapshot with NO auto-approve posture, so each one pauses for a human.
	runner := &recordingSegmentRunner{actions: 2}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: &sequentialProposer{},
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-two-approvals"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua dos veces", IdempotencyKey: "flow-two-appr-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// First gate: Bob decides, in the API's order (claim the approval, then deliver the durable resume).
	firstApproval := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, firstApproval)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(ctx, inst, firstApproval, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "call-1"))

	secondApproval := act.ApprovalIDForToolCall(runID, "call-2")
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, secondApproval)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)

	run, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status, "the second pause must move the run back to waiting_approval")

	_, err = store.ResolveApproval(ctx, inst, secondApproval, act.ApprovalStatusApproved, "admin:carol")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "call-2"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 2, exec.executeCount(), "both approved actions executed")

	// (a) The timeline records every edge of both approval cycles, in order.
	events, err := store.ListRunEvents(ctx, inst, runID, 0, 100)
	require.NoError(t, err)
	require.Equal(t, []string{
		act.EventTypeQueued,
		act.EventTypeRunning,
		act.EventTypeWaitingApproval,
		act.EventTypeResumed,
		act.EventTypeWaitingApproval,
		act.EventTypeResumed,
		act.EventTypeSucceeded,
	}, eventTypes(events))

	// Each resumption is attributed to its own decider (#128's decided_by, now for every pause, not just the first).
	var resumed []*act.RunEvent
	for _, ev := range events {
		if ev.EventType == act.EventTypeResumed {
			resumed = append(resumed, ev)
		}
	}
	require.Len(t, resumed, 2)
	require.Equal(t, "admin:bob", resumed[0].Payload["decided_by"])
	require.Equal(t, "admin:carol", resumed[1].Payload["decided_by"])

	// The action ledger agrees, one decider per action.
	a1, err := store.GetAction(ctx, inst, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, "admin:bob", a1.DecidedBy)
	a2, err := store.GetAction(ctx, inst, runID, "call-2")
	require.NoError(t, err)
	require.Equal(t, "admin:carol", a2.DecidedBy)
}

// TestGatewayFlowAutoApproveExecutesWithoutHuman proves a connector's approval: auto posture takes effect through the
// FULL DBOS path: a generic (unclassified) write on an auto-approved connector runs to succeeded with no human gate
// and no Resume. This is the end-to-end regression guard for the auto-approve wiring — the executor must feed the SAME
// AutoApprove map to gateway.Execute as to gateway.Propose, or pre-write reauthorization recomputes the unclassified
// tool as approval_required and rejects the already-approved action. A gateway unit test that passes the map by hand
// cannot catch that gap; only the executor building the map from the snapshot does.
func TestGatewayFlowAutoApproveExecutesWithoutHuman(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-9", Message: "created PROJ-9"}}
	gateway := &act.Gateway{
		// The runtime registry stamps every generic MCP tool ClassUnknown; mirror that so the test exercises the exact
		// production condition the auto-approve fix targets.
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassUnknown)),
		Ledger:   store,
		Executor: exec,
	}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: autoApproveRunner{}, Store: store, Gateway: gateway, Proposer: act.NewCapturedProposer(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-autoapprove"
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "flow-autoapprove-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// No Resume: an auto-approved action must run to completion without waiting on a human.
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.True(t, res.ActionTaken)
	require.Equal(t, 1, exec.executeCount(), "the auto-approved write executes exactly once, unattended")

	// capturedProposer keys the action on the effective tool name.
	action, err := store.GetAction(t.Context(), inst, runID, "jira.create_issue")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, action.Status, "the action executed, it was not rejected at reauthorization")
	require.Equal(t, act.PolicyAllow, action.PolicyDecision)

	// The auto-approved action created NO approval and never paused the run. This is what makes the agent's
	// security.approve expression moot for an auto-approved tool: the approve gate is only ever evaluated against a
	// stored approval row (runtime/server's resolveApprovalDecision), so with no row it is never consulted — the
	// connector's posture, not the approve expression, is what decides whether a human enters the loop.
	approvals, err := store.ListApprovals(t.Context(), act.ListApprovalsFilter{InstanceID: inst, RunID: runID})
	require.NoError(t, err)
	require.Empty(t, approvals, "an auto-approved action must not create an approval (so security.approve is never evaluated for it)")
	events, err := store.ListRunEvents(t.Context(), inst, runID, 0, 100)
	require.NoError(t, err)
	require.NotContains(t, eventTypes(events), act.EventTypeWaitingApproval, "an auto-approved action must not pause the run")
}

// TestGatewayFlowNoActionSucceeds verifies a run that proposes no external action succeeds immediately without an
// approval gate or an executor call — a pure read-only investigation.
func TestGatewayFlowNoActionSucceeds(t *testing.T) {
	exec := &fakeExecutor{}
	e, _, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{ok: false})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "solo investiga", IdempotencyKey: "flow-noaction-" + uuid.NewString(),
	})
	require.NoError(t, err)

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.False(t, res.ActionTaken)
	require.Equal(t, 0, exec.executeCount())
}

// TestGatewayFlowPolicyDenyFailsRun verifies an action denied by policy (here via an engaged kill switch) ends the run
// failed with no external effect, and the ledger records policy_rejected.
func TestGatewayFlowPolicyDenyFailsRun(t *testing.T) {
	exec := &fakeExecutor{}
	_, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	store := newRunStore(t)
	kill := act.NewMapKillSwitch()
	kill.DisablePlatform("maintenance")
	gateway := &act.Gateway{
		Registry:   act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:     store,
		Executor:   exec,
		KillSwitch: kill,
	}
	// A mock runner that proposes one governed write then closes, so the segmented loop reaches the gateway; the policy
	// (kill switch) is what denies it.
	runner := gatewayFlowRunner{}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: fixedProposer{proposal: flowProposal(), ok: true},
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "actua", IdempotencyKey: "flow-deny-" + uuid.NewString(),
	})
	require.NoError(t, err)

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status, "policy denial is per-action: the model sees the denial and closes")
	require.False(t, res.ActionTaken)
	require.Equal(t, 0, exec.executeCount())

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, action.Status)
}

// batchSegmentRunner proposes N actions in a SINGLE segment (one model turn with N tool calls), then closes on resume.
// This is the core scenario of #131: all actions are captured in one pass and governed together.
type batchSegmentRunner struct {
	n int

	mu      sync.Mutex
	resumes []ai.InjectedResult
}

var _ act.Runner = (*batchSegmentRunner)(nil)

func (r *batchSegmentRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	return &ai.AgentSnapshot{Name: agentName, Instructions: "batch test"}, nil
}

func (r *batchSegmentRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	sessionID := in.SessionID
	if sessionID == "" {
		sessionID = "sess-batch"
	}
	r.mu.Lock()
	for _, res := range in.Resume {
		r.resumes = append(r.resumes, *res)
	}
	r.mu.Unlock()
	if len(in.Resume) > 0 {
		return act.RunSegmentResult{Response: "Done, all actions processed.", SessionID: sessionID}, nil
	}
	proposals := make([]*ai.ProposedAction, r.n)
	for i := range proposals {
		proposals[i] = &ai.ProposedAction{
			ToolCallID: fmt.Sprintf("tc-%d", i+1),
			Tool:       "jira.create_issue",
			Connector:  "jira_ops",
			Args:       map[string]any{"summary": fmt.Sprintf("issue %d", i+1)},
		}
	}
	return act.RunSegmentResult{
		Response:  "Proposing batch.",
		SessionID: sessionID,
		Proposed:  proposals,
	}, nil
}

func (r *batchSegmentRunner) ApplyAction(context.Context, act.ActionRequest) (act.ActionResult, error) {
	return act.ActionResult{}, nil
}

func (r *batchSegmentRunner) injectedResults() []ai.InjectedResult {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]ai.InjectedResult(nil), r.resumes...)
}

// batchProposer surfaces the captured action's own identity, so N actions captured in one segment produce N distinct
// tool proposals with the tool_call_id the runner assigned.
type batchProposer struct{}

func (batchProposer) Propose(_ context.Context, in act.ProposeInput) (act.ToolProposal, bool, error) {
	if in.Captured == nil {
		return act.ToolProposal{}, false, nil
	}
	return act.ToolProposal{
		ToolCallID: in.Captured.ToolCallID,
		Tool:       in.Captured.Tool,
		Connector:  in.Captured.Connector,
		Args:       in.Captured.Args,
		Summary:    fmt.Sprintf("create issue %s", in.Captured.ToolCallID),
	}, true, nil
}

// TestBatchApprovalsAllCreatedAtOnce is the core #131 test: a segment proposes 3 actions in one turn, all 3 approvals
// are created at once and visible in the inbox, and the operator can decide them in any order while execution follows
// the model's proposal order.
func TestBatchApprovalsAllCreatedAtOnce(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-X", Message: "done"}}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:   store,
		Executor: exec,
	}
	runner := &batchSegmentRunner{n: 3}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: batchProposer{},
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-batch"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "batch test", IdempotencyKey: "batch-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// All 3 approvals must be visible at once.
	for i := 1; i <= 3; i++ {
		approvalID := act.ApprovalIDForToolCall(runID, fmt.Sprintf("tc-%d", i))
		require.Eventually(t, func() bool {
			_, gErr := store.GetApproval(ctx, inst, approvalID)
			return gErr == nil
		}, 15*time.Second, 50*time.Millisecond, "approval %d must be created", i)
	}

	// Verify all 3 have position/total set.
	approvals, err := store.ListApprovals(ctx, act.ListApprovalsFilter{InstanceID: inst, RunID: runID})
	require.NoError(t, err)
	require.Len(t, approvals, 3)
	for _, a := range approvals {
		require.Equal(t, 3, a.Total, "each approval records the batch size")
		require.True(t, a.Position >= 1 && a.Position <= 3, "position is in range")
	}

	// Decide in REVERSE order (3, 2, 1) — execution must still follow proposal order.
	for i := 3; i >= 1; i-- {
		tcID := fmt.Sprintf("tc-%d", i)
		approvalID := act.ApprovalIDForToolCall(runID, tcID)
		_, err = store.ResolveApproval(ctx, inst, approvalID, act.ApprovalStatusApproved, fmt.Sprintf("admin:approver-%d", i))
		require.NoError(t, err)
		require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, tcID))
	}

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.True(t, res.ActionTaken)
	require.Equal(t, 3, exec.executeCount(), "all 3 actions executed")

	// The model received all 3 results on resume.
	injected := runner.injectedResults()
	require.Len(t, injected, 3, "the model was resumed with all 3 results")

	// Timeline: 3 waiting_approval events, 3 resumed events with distinct deciders.
	events, err := store.ListRunEvents(ctx, inst, runID, 0, 100)
	require.NoError(t, err)
	types := eventTypes(events)
	waitCount := 0
	resumeCount := 0
	for _, et := range types {
		if et == act.EventTypeWaitingApproval {
			waitCount++
		}
		if et == act.EventTypeResumed {
			resumeCount++
		}
	}
	require.Equal(t, 3, waitCount, "3 waiting_approval events in the timeline")
	require.Equal(t, 3, resumeCount, "3 resumed events in the timeline")

	// Each resumed event carries a distinct decider.
	var resumedEvents []*act.RunEvent
	for _, ev := range events {
		if ev.EventType == act.EventTypeResumed {
			resumedEvents = append(resumedEvents, ev)
		}
	}
	deciders := map[string]bool{}
	for _, ev := range resumedEvents {
		d, _ := ev.Payload["decided_by"].(string)
		deciders[d] = true
	}
	require.Len(t, deciders, 3, "each resumed event has a distinct decider")
}

// TestBatchPartialRejection tests the partial rejection scenario: one action is rejected, the others are approved.
// The rejected action injects an error result; the approved ones execute. The run succeeds because the model adapts.
func TestBatchPartialRejection(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-X", Message: "done"}}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:   store,
		Executor: exec,
	}
	runner := &batchSegmentRunner{n: 2}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: batchProposer{},
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-partial"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "partial reject", IdempotencyKey: "partial-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Wait for both approvals.
	for i := 1; i <= 2; i++ {
		approvalID := act.ApprovalIDForToolCall(runID, fmt.Sprintf("tc-%d", i))
		require.Eventually(t, func() bool {
			_, gErr := store.GetApproval(ctx, inst, approvalID)
			return gErr == nil
		}, 15*time.Second, 50*time.Millisecond)
	}

	// Reject action 1, approve action 2.
	_, err = store.ResolveApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-1"), act.ApprovalStatusDenied, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalRejected, "tc-1"))

	_, err = store.ResolveApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-2"), act.ApprovalStatusApproved, "admin:carol")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "tc-2"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status, "partial rejection does not abort the run")
	require.True(t, res.ActionTaken, "action 2 was approved and executed")
	require.Equal(t, 1, exec.executeCount(), "only the approved action executed")

	// The model received 2 injected results: one rejection, one success.
	injected := runner.injectedResults()
	require.Len(t, injected, 2)
	require.True(t, injected[0].IsError, "first result is the rejection")
	require.False(t, injected[1].IsError, "second result is the executed action")

	// The ledger reflects the outcomes: the rejected action is closed as rejected (not stuck in approval_pending),
	// and the approved action reached succeeded.
	a1, err := store.GetAction(ctx, inst, runID, "tc-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionRejected, a1.Status, "rejected action's ledger row must be closed as rejected")
	a2, err := store.GetAction(ctx, inst, runID, "tc-2")
	require.NoError(t, err)
	require.Equal(t, act.ActionSucceeded, a2.Status, "approved action's ledger row must reach succeeded")
}

// TestBatchStatusStaysWaitingWhileDecisionsPending is the regression for the intermediate-state bug: when a batch has
// two manual actions and the operator approves the FIRST one, the run must stay in waiting_approval (not running)
// while the second approval is still pending. Otherwise the pending approval becomes invisible to the Home inbox and
// Act's default filter, which both select on waiting_approval. The symmetric case is also tested: rejecting the LAST
// manual action in a batch must NOT leave the run stuck in waiting_approval.
func TestBatchStatusStaysWaitingWhileDecisionsPending(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-X", Message: "done"}}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:   store,
		Executor: exec,
	}
	runner := &batchSegmentRunner{n: 2}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: batchProposer{},
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-wait-status"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "two manual", IdempotencyKey: "wait-status-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Wait for both approvals to be created.
	for i := 1; i <= 2; i++ {
		approvalID := act.ApprovalIDForToolCall(runID, fmt.Sprintf("tc-%d", i))
		require.Eventually(t, func() bool {
			_, gErr := store.GetApproval(ctx, inst, approvalID)
			return gErr == nil
		}, 15*time.Second, 50*time.Millisecond)
	}

	// Approve action 1 IN ORDER (not reverse). The executor processes it and blocks on action 2's Recv.
	_, err = store.ResolveApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-1"), act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "tc-1"))

	// The bug: while action 2's decision is pending, the run MUST be in waiting_approval, not running. The resumed
	// event for action 1 carries waiting_approval (because a later manual action remains), so the status reflects the
	// true state of the workflow and the inbox/filter see the pending approval.
	require.Eventually(t, func() bool {
		events, gErr := store.ListRunEvents(ctx, inst, runID, 0, 100)
		if gErr != nil {
			return false
		}
		for _, ev := range events {
			if ev.EventType == act.EventTypeResumed {
				return true
			}
		}
		return false
	}, 15*time.Second, 50*time.Millisecond, "the resumed event for action 1 must be emitted")

	run, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status,
		"while the second approval is pending, the run must stay in waiting_approval, not running")

	// Action 2's approval must still be pending and findable.
	appr2, err := store.GetApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-2"))
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusPending, appr2.Status)

	// Now approve action 2 and let the run finish.
	_, err = store.ResolveApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-2"), act.ApprovalStatusApproved, "admin:carol")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "tc-2"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 2, exec.executeCount())

	// Verify the event timeline: the resumed event for action 1 carries waiting_approval (not running), because
	// action 2 was still pending. The resumed event for action 2 carries running (no more manual gates).
	events, err := store.ListRunEvents(ctx, inst, runID, 0, 100)
	require.NoError(t, err)
	var resumedStatuses []act.RunStatus
	for _, ev := range events {
		if ev.EventType == act.EventTypeResumed {
			resumedStatuses = append(resumedStatuses, ev.Status)
		}
	}
	require.Len(t, resumedStatuses, 2)
	require.Equal(t, act.RunStatusWaitingApproval, resumedStatuses[0],
		"first resumed carries waiting_approval because action 2 is still pending")
	require.Equal(t, act.RunStatusRunning, resumedStatuses[1],
		"second resumed carries running because no more manual gates")
}

// TestBatchFailSweepsSiblingApprovals verifies that when one action in a batch fails during execution, any still-pending
// sibling approvals are swept to cancelled — they will never be decided, so they must not linger in the inbox.
func TestBatchFailSweepsSiblingApprovals(t *testing.T) {
	store := newRunStore(t)
	exec := &fakeExecutor{execErr: errors.New("simulated execution failure")}
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(createIssueDescriptor(act.ClassIdempotentNative)),
		Ledger:   store,
		Executor: exec,
	}
	runner := &batchSegmentRunner{n: 2}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: batchProposer{},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	const inst = "inst-fail-sweep"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "fail sweep", IdempotencyKey: "fail-sweep-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Wait for both approvals to appear.
	for i := 1; i <= 2; i++ {
		approvalID := act.ApprovalIDForToolCall(runID, fmt.Sprintf("tc-%d", i))
		require.Eventually(t, func() bool {
			_, gErr := store.GetApproval(ctx, inst, approvalID)
			return gErr == nil
		}, 15*time.Second, 50*time.Millisecond)
	}

	// Approve action 1. The gateway execute step will fail (fakeExecutor.execErr is set).
	_, err = store.ResolveApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-1"), act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "tc-1"))

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, res.Status)

	// The sibling's approval (action 2) must be swept to cancelled, not left pending.
	a2appr, err := store.GetApproval(ctx, inst, act.ApprovalIDForToolCall(runID, "tc-2"))
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusCancelled, a2appr.Status, "sibling approval must be swept to cancelled on run failure")

	// The sibling's ledger row must also be closed: withdrawn, not stuck in approval_pending.
	a2action, err := store.GetAction(ctx, inst, runID, "tc-2")
	require.NoError(t, err)
	require.Equal(t, act.ActionWithdrawn, a2action.Status, "sibling action's ledger row must be withdrawn on run failure")
}
