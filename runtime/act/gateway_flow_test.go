package act_test

import (
	"context"
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
	if in.Resume != nil {
		return act.RunSegmentResult{Response: "Hecho, creé el ticket.", SessionID: sessionID}, nil
	}
	return act.RunSegmentResult{
		Response:  "Propongo crear un ticket.",
		SessionID: sessionID,
		Proposed:  &ai.ProposedAction{Tool: "jira.create_issue", Connector: "jira_ops"},
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
		ApprovalTimeout:    60 * time.Second,
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

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved))
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
		_, gErr := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForRun(runID))
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(t.Context(), instanceID, act.ApprovalIDForRun(runID), act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved))

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
		_, gErr := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForRun(runID))
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)

	// Cancel the parked run, then — racing a late approver — try to deliver an approval. Neither call should drive the
	// external write. Resume may error (a cancelled workflow refuses the send); that is fine, we only require no write.
	require.NoError(t, e.Cancel(t.Context(), instanceID, runID))
	_ = e.Resume(t.Context(), runID, act.ApprovalApproved)

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
	appr, err := store.GetApproval(t.Context(), instanceID, act.ApprovalIDForRun(runID))
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusCancelled, appr.Status, "cancelling a run must withdraw its pending approval")
}

// TestGatewayFlowRejectionSkipsWrite verifies a human denial ends the run rejected and never reaches the executor; the
// action stays approval_pending on the ledger.
func TestGatewayFlowRejectionSkipsWrite(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1"}}
	e, store, instanceID := newGatewayExecutor(t, exec, createIssueDescriptor(act.ClassIdempotentNative), nil, fixedProposer{proposal: flowProposal(), ok: true})

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "fallo", IdempotencyKey: "flow-reject-" + uuid.NewString(),
	})
	require.NoError(t, err)

	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalRejected))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusRejected, res.Status)
	require.False(t, res.ActionTaken)
	require.Equal(t, 0, exec.executeCount())

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionApprovalPending, action.Status)
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
	if in.Resume != nil {
		return act.RunSegmentResult{Response: "Listo, creé el ticket.", SessionID: sessionID}, nil
	}
	// First segment: propose the governed write on the auto-approved connector.
	return act.RunSegmentResult{
		Response:  "Propongo crear un ticket.",
		SessionID: sessionID,
		Proposed: &ai.ProposedAction{
			Connector: "jira_ops", Tool: "jira.create_issue",
			Args: map[string]any{"summary": "coste alto"}, Summary: "crear ticket P2",
		},
	}, nil
}

// ApplyAction is unused on the gateway path (the gateway performs the write); present only to satisfy act.Runner.
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
	if in.Resume != nil {
		r.resumes = append(r.resumes, *in.Resume)
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
		Proposed:  &ai.ProposedAction{Tool: "jira.create_issue", Connector: "jira_ops", Args: map[string]any{"summary": "coste alto"}},
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
		ApprovalTimeout: 60 * time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
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
		ApprovalTimeout: 60 * time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
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
		ApprovalTimeout: 60 * time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
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
	firstApproval := act.ApprovalIDForSegment(runID, 0)
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, firstApproval)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	_, err = store.ResolveApproval(ctx, inst, firstApproval, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved))

	// Second gate: the run pauses again on its next governed action. The approval row is written in the same durable
	// step as the waiting_approval transition, after it, so once the row is visible the status write has committed.
	secondApproval := act.ApprovalIDForSegment(runID, 1)
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, secondApproval)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)

	// (b) The bug: while the second decision is pending, the run must SAY it is waiting — this status is what the Home
	// inbox and Act's default filter select on, so "running" here is an approval nobody sees.
	run, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status, "the second pause must move the run back to waiting_approval")

	// Carol, not Bob, decides the second action.
	_, err = store.ResolveApproval(ctx, inst, secondApproval, act.ApprovalStatusApproved, "admin:carol")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved))

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
		ApprovalTimeout: 60 * time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
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
		ApprovalTimeout: 60 * time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: instanceID, AgentName: "triage", Prompt: "actua", IdempotencyKey: "flow-deny-" + uuid.NewString(),
	})
	require.NoError(t, err)

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusFailed, res.Status)
	require.False(t, res.ActionTaken)
	require.Equal(t, 0, exec.executeCount())

	action, err := store.GetAction(t.Context(), instanceID, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.ActionPolicyRejected, action.Status)
}
