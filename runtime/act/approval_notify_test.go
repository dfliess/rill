package act_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// These tests cover the notification the executor emits when a run pauses on human approvals (kairos-cloud#143):
// one per pause, carrying the batch's actions, and unable to affect the run whatever the notifier does.

// recordingApprovalNotifier records the pauses it is told about, in order.
type recordingApprovalNotifier struct {
	mu      sync.Mutex
	notices []act.PendingApprovals
}

var _ act.ApprovalNotifier = (*recordingApprovalNotifier)(nil)

func (n *recordingApprovalNotifier) NotifyPendingApprovals(_ context.Context, pending act.PendingApprovals) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notices = append(n.notices, pending)
}

func (n *recordingApprovalNotifier) all() []act.PendingApprovals {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]act.PendingApprovals(nil), n.notices...)
}

// panicApprovalNotifier stands in for a broken notifier: the run must survive it.
type panicApprovalNotifier struct{}

func (panicApprovalNotifier) NotifyPendingApprovals(context.Context, act.PendingApprovals) {
	panic("notifier is broken")
}

// notifyRunner proposes a configurable number of actions on each successive segment, then closes. pauses[i] is how
// many actions segment i proposes, so {3} is one pause of three actions and {1, 1} is two pauses of one.
type notifyRunner struct {
	pauses   []int
	tool     string
	snapshot *ai.AgentSnapshot

	mu  sync.Mutex
	seg int
}

var _ act.Runner = (*notifyRunner)(nil)

func (r *notifyRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	if r.snapshot != nil {
		return r.snapshot, nil
	}
	return &ai.AgentSnapshot{Name: agentName, Instructions: "approval notification test"}, nil
}

func (r *notifyRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	sessionID := in.SessionID
	if sessionID == "" {
		sessionID = "sess-notify"
	}

	r.mu.Lock()
	seg := r.seg
	r.seg++
	r.mu.Unlock()

	if seg >= len(r.pauses) {
		return act.RunSegmentResult{Response: "Listo.", SessionID: sessionID}, nil
	}
	proposals := make([]*ai.ProposedAction, r.pauses[seg])
	for i := range proposals {
		proposals[i] = &ai.ProposedAction{
			ToolCallID: fmt.Sprintf("tc-%d-%d", seg+1, i+1),
			Tool:       r.tool,
			Connector:  "jira_ops",
			Args:       map[string]any{"summary": fmt.Sprintf("issue %d", i+1)},
		}
	}
	return act.RunSegmentResult{Response: "Propongo acciones.", SessionID: sessionID, Proposed: proposals}, nil
}

func (r *notifyRunner) ApplyAction(context.Context, act.ActionRequest) (act.ActionResult, error) {
	return act.ActionResult{}, nil
}

// newNotifyExecutor wires an executor over a fresh run store with the given runner and notifier. Its registry serves
// the namespaced MCP tool (mcpCreateIssueDescriptor), so the proposals carry the name the ledger stores and the
// notification has to reduce it to the raw one.
func newNotifyExecutor(t *testing.T, runner act.Runner, notifier act.ApprovalNotifier) (*act.DBOSExecutor, *act.PostgresRunStore, *fakeExecutor) {
	t.Helper()
	store := newRunStore(t)
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-1", Message: "done"}}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Proposer: batchProposer{},
		Gateway: &act.Gateway{
			Registry: act.NewMapToolRegistry(mcpCreateIssueDescriptor()),
			Ledger:   store,
			Executor: exec,
		},
		Approvals: notifier,
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })
	return e, store, exec
}

// approveToolCall waits for an action's approval to exist, resolves it and hands the decision to the run.
func approveToolCall(t *testing.T, e *act.DBOSExecutor, store *act.PostgresRunStore, instanceID, runID, toolCallID string) {
	t.Helper()
	approvalID := act.ApprovalIDForToolCall(runID, toolCallID)
	require.Eventually(t, func() bool {
		_, err := store.GetApproval(t.Context(), instanceID, approvalID)
		return err == nil
	}, 15*time.Second, 50*time.Millisecond, "approval for %s must be created", toolCallID)

	_, err := store.ResolveApproval(t.Context(), instanceID, approvalID, act.ApprovalStatusApproved, "admin:ana")
	require.NoError(t, err)
	require.NoError(t, e.Resume(t.Context(), runID, act.ApprovalApproved, toolCallID))
}

// TestApprovalNotifierBatchIsOneNotification is the decision the grill fixed: a batch is one moment of signing, so
// three actions pausing together produce ONE notification that carries all three, not three notifications.
func TestApprovalNotifierBatchIsOneNotification(t *testing.T) {
	notifier := &recordingApprovalNotifier{}
	e, store, _ := newNotifyExecutor(t, &notifyRunner{pauses: []int{3}, tool: "mcp.jira_ops.create_issue"}, notifier)

	const inst = "inst-notify-batch"
	key := "notify-batch-" + uuid.NewString()
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: key,
		Actor: act.Actor{Subject: "usr_marta"},
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool { return len(notifier.all()) > 0 }, 15*time.Second, 50*time.Millisecond)
	notices := notifier.all()
	require.Len(t, notices, 1, "one pause is one notification, whatever the batch size")
	require.Equal(t, inst, notices[0].InstanceID)
	require.Equal(t, runID, notices[0].RunID)
	require.Equal(t, "triage", notices[0].AgentName)
	require.Equal(t, key, notices[0].IdempotencyKey, "the run's URL is keyed by the idempotency key, not the composite run ID")
	require.Equal(t, "usr_marta", notices[0].RunActor, "an approve policy may discriminate on who the run acts for")
	// The RAW tool name: the approve expression matches the name the connector's own configuration uses, not the
	// namespaced form the ledger stores.
	require.Equal(t, []act.PendingAction{
		{Tool: "create_issue", Connector: "jira_ops"},
		{Tool: "create_issue", Connector: "jira_ops"},
		{Tool: "create_issue", Connector: "jira_ops"},
	}, notices[0].Actions)

	for i := 1; i <= 3; i++ {
		approveToolCall(t, e, store, inst, runID, fmt.Sprintf("tc-1-%d", i))
	}
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Len(t, notifier.all(), 1, "the decisions do not notify again")
}

// TestApprovalNotifierOnePerPause checks the other half of the same decision: a run that pauses twice is announced
// twice, because each pause is a fresh demand on the approver.
func TestApprovalNotifierOnePerPause(t *testing.T) {
	notifier := &recordingApprovalNotifier{}
	e, store, exec := newNotifyExecutor(t, &notifyRunner{pauses: []int{1, 1}, tool: "mcp.jira_ops.create_issue"}, notifier)

	const inst = "inst-notify-pauses"
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua dos veces", IdempotencyKey: "notify-pauses-" + uuid.NewString(),
		Actor: act.Actor{Subject: "usr_marta"},
	})
	require.NoError(t, err)

	approveToolCall(t, e, store, inst, runID, "tc-1-1")
	approveToolCall(t, e, store, inst, runID, "tc-2-1")

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 2, exec.executeCount())

	notices := notifier.all()
	require.Len(t, notices, 2, "each pause of the run is announced")
	for _, n := range notices {
		require.Equal(t, runID, n.RunID)
		require.Len(t, n.Actions, 1)
	}
}

// TestApprovalNotifierSkipsAutoApproved checks that a batch nobody has to sign announces nothing: the connector's
// auto-approve posture executes it without a human, so there is no pause to report.
func TestApprovalNotifierSkipsAutoApproved(t *testing.T) {
	notifier := &recordingApprovalNotifier{}
	runner := &notifyRunner{
		pauses:   []int{1},
		tool:     "mcp.jira_ops.create_issue",
		snapshot: &ai.AgentSnapshot{Name: "triage", Instructions: "auto", MCPConnectors: []ai.MCPConnector{{Name: "jira_ops", Approval: "auto"}}},
	}
	e, _, exec := newNotifyExecutor(t, runner, notifier)

	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: "inst-notify-auto", AgentName: "triage", Prompt: "actua", IdempotencyKey: "notify-auto-" + uuid.NewString(),
		Actor: act.Actor{Subject: "usr_marta"},
	})
	require.NoError(t, err)

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, exec.executeCount(), "the action ran without a human")
	require.Empty(t, notifier.all(), "an auto-approved batch has nobody to notify")
}

// TestApprovalNotifierFailureDoesNotBreakRun pins the contract the emission point rests on: the approvals are
// already committed when the notifier is called, so even a notifier that blows up leaves the run waiting for its
// decision, and the decision still completes it.
func TestApprovalNotifierFailureDoesNotBreakRun(t *testing.T) {
	e, store, exec := newNotifyExecutor(t, &notifyRunner{pauses: []int{1}, tool: "mcp.jira_ops.create_issue"}, panicApprovalNotifier{})

	const inst = "inst-notify-panic"
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "notify-panic-" + uuid.NewString(),
		Actor: act.Actor{Subject: "usr_marta"},
	})
	require.NoError(t, err)

	approveToolCall(t, e, store, inst, runID, "tc-1-1")

	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, exec.executeCount())
}
