package act_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// This file covers the translation of a connector's declared approval posture (approval / auto_approve /
// require_approval, as the agent author writes it in YAML) into what actually happens to a concrete proposed action,
// through the full DBOS path: the executor derives the auto-approve decision from the snapshot's connector via
// connectorAutoApproves(conn, rawToolName(tool, connector)), so these tests are what pins the layer between the
// parsed definition and the policy engine (whose branches policy_test.go covers in isolation).

// postureRunner is a mock Runner whose snapshot carries one MCP connector with a configurable approval posture, and
// which proposes a single governed write on that connector using the FULL effective tool name (mcp.<connector>.<raw>),
// exactly as the production loop captures it. Using the effective name is the point: the executor must strip the
// mcp.<connector>. prefix before matching the posture globs, and only a proposal shaped like production's can observe
// a regression that matched them against the effective name instead. The connector is swappable under a lock so a
// test can simulate an admin editing the live definition while a run is parked on its approval gate.
type postureRunner struct {
	mu   sync.Mutex
	conn ai.MCPConnector
}

var _ act.Runner = (*postureRunner)(nil)

func (r *postureRunner) setConnector(conn ai.MCPConnector) {
	r.mu.Lock()
	r.conn = conn
	r.mu.Unlock()
}

func (r *postureRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	r.mu.Lock()
	conn := r.conn
	r.mu.Unlock()
	return &ai.AgentSnapshot{
		Name:          agentName,
		Instructions:  "approval posture flow",
		MCPConnectors: []ai.MCPConnector{conn},
	}, nil
}

func (r *postureRunner) RunSegment(_ context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	sessionID := in.SessionID
	if sessionID == "" {
		sessionID = "sess-" + in.Snapshot.Name
	}
	if len(in.Resume) > 0 {
		return act.RunSegmentResult{Response: "Listo, creé el ticket.", SessionID: sessionID}, nil
	}
	return act.RunSegmentResult{
		Response:  "Propongo crear un ticket.",
		SessionID: sessionID,
		Proposed: []*ai.ProposedAction{{
			ToolCallID: "call-1",
			Connector:  "jira_ops",
			Tool:       "mcp.jira_ops.create_issue",
			Args:       map[string]any{"summary": "coste alto"},
			Summary:    "crear ticket P2",
		}},
	}, nil
}

func (r *postureRunner) ApplyAction(context.Context, act.ActionRequest) (act.ActionResult, error) {
	return act.ActionResult{}, nil
}

// mcpCreateIssueDescriptor is the registry entry for the effective-named MCP tool the postureRunner proposes. It is
// ClassUnknown on purpose: the runtime registry stamps every generic MCP tool unclassified, so an auto-approve that
// only worked for a classified tool would silently no-op in production.
func mcpCreateIssueDescriptor() act.ToolDescriptor {
	return act.ToolDescriptor{
		Name:        "mcp.jira_ops.create_issue",
		Connector:   "jira_ops",
		Version:     "v1",
		Class:       act.ClassUnknown,
		InputSchema: objectSchema("summary"),
	}
}

// newPostureExecutor builds an executor whose gateway serves the effective-named MCP tool, wired with the captured
// proposer so the postureRunner's proposal flows through exactly like a production capture.
func newPostureExecutor(t *testing.T, runner act.Runner, exec act.ActionExecutor) (*act.DBOSExecutor, *act.PostgresRunStore) {
	t.Helper()
	store := newRunStore(t)
	gateway := &act.Gateway{
		Registry: act.NewMapToolRegistry(mcpCreateIssueDescriptor()),
		Ledger:   store,
		Executor: exec,
	}
	dsn, schema := requirePostgres(t)
	e, err := act.NewDBOSExecutor(t.Context(), act.Config{
		DatabaseURL: dsn, DatabaseSchema: schema, ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner: runner, Store: store, Gateway: gateway, Proposer: act.NewCapturedProposer(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })
	return e, store
}

// TestGatewayFlowManualPostureAutoApproveGlobRunsUnattended: on the default (manual) posture, a tool matching an
// auto_approve glob executes with no human gate — and creates NO approval, so there is nothing for the agent's
// security.approve gate to ever decide (the approve expression is only evaluated against a stored approval row in
// runtime/server's resolveApprovalDecision; no row means it is never consulted for this action).
func TestGatewayFlowManualPostureAutoApproveGlobRunsUnattended(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-7", Message: "created PROJ-7"}}
	// Approval is unset: the default posture is manual, so ONLY the auto_approve glob exempts the matching tool.
	runner := &postureRunner{conn: ai.MCPConnector{Name: "jira_ops", AutoApprove: []string{"create_*"}}}
	e, store := newPostureExecutor(t, runner, exec)

	const inst = "inst-manual-glob"
	runID, err := e.Start(t.Context(), act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "posture-manual-glob-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// No Resume: the glob-matched write runs unattended to completion.
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.True(t, res.ActionTaken)
	require.Equal(t, 1, exec.executeCount())

	action, err := store.GetAction(t.Context(), inst, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.PolicyAllow, action.PolicyDecision)

	// The auto-approved action created no approval and never paused the run: the security.approve interaction above.
	approvals, err := store.ListApprovals(t.Context(), act.ListApprovalsFilter{InstanceID: inst, RunID: runID})
	require.NoError(t, err)
	require.Empty(t, approvals, "an auto-approved action must not create an approval")
	events, err := store.ListRunEvents(t.Context(), inst, runID, 0, 100)
	require.NoError(t, err)
	require.NotContains(t, eventTypes(events), act.EventTypeWaitingApproval, "an auto-approved action must not pause the run")
}

// TestGatewayFlowAutoApproveGlobMatchesRawToolNameOnly pins WHICH name the posture globs match: the server's raw tool
// name (create_issue), never the effective namespaced one (mcp.jira_ops.create_issue). An author who writes the
// pattern against the effective name gets NO auto-approval — the action pauses for a human — rather than a silently
// broader or narrower match.
func TestGatewayFlowAutoApproveGlobMatchesRawToolNameOnly(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-7"}}
	// The pattern is written against the EFFECTIVE tool name. The raw name is "create_issue", so it must not match.
	runner := &postureRunner{conn: ai.MCPConnector{Name: "jira_ops", AutoApprove: []string{"mcp.jira_ops.create_*"}}}
	e, store := newPostureExecutor(t, runner, exec)

	const inst = "inst-effective-glob"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "posture-effective-glob-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// The run parks on the approval gate: the effective-name pattern did not auto-approve the action.
	approvalID := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, approvalID)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	require.Equal(t, 0, exec.executeCount(), "a non-matching pattern must not execute the write unattended")
	run, err := store.GetRun(ctx, inst, runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusWaitingApproval, run.Status)

	// Finish the run cleanly: a human approves and the write executes.
	_, err = store.ResolveApproval(ctx, inst, approvalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, exec.executeCount())
}

// TestGatewayFlowAutoPostureRequireApprovalGates: on approval: auto, a tool matching a require_approval glob is the
// exception that still needs a human — the inverse of the manual+auto_approve case above, through the same full path.
func TestGatewayFlowAutoPostureRequireApprovalGates(t *testing.T) {
	exec := &fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-7"}}
	runner := &postureRunner{conn: ai.MCPConnector{Name: "jira_ops", Approval: "auto", RequireApproval: []string{"create_*"}}}
	e, store := newPostureExecutor(t, runner, exec)

	const inst = "inst-auto-gate"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "posture-auto-gate-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	// Despite the auto posture, the matching tool pauses for a human and the ledger records approval_required.
	approvalID := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, approvalID)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)
	require.Equal(t, 0, exec.executeCount(), "a require_approval match must not execute unattended")
	action, err := store.GetAction(ctx, inst, runID, "call-1")
	require.NoError(t, err)
	require.Equal(t, act.PolicyApprovalRequired, action.PolicyDecision)

	_, err = store.ResolveApproval(ctx, inst, approvalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)
	require.Equal(t, 1, exec.executeCount())
}

// connectorRecordingExecutor records the MCP connector the workflow bound onto the execute step context, so a test
// can assert the write runs against the connector frozen in the run's snapshot rather than the live definition.
type connectorRecordingExecutor struct {
	fakeExecutor

	connMu sync.Mutex
	conns  []ai.MCPConnector
	bound  []bool
}

func (e *connectorRecordingExecutor) Execute(ctx context.Context, req act.ExecuteRequest) (act.ExecuteResult, error) {
	conn, ok := act.MCPConnectorFromContextForTest(ctx)
	e.connMu.Lock()
	e.conns = append(e.conns, conn)
	e.bound = append(e.bound, ok)
	e.connMu.Unlock()
	return e.fakeExecutor.Execute(ctx, req)
}

func (e *connectorRecordingExecutor) boundConnectors() ([]ai.MCPConnector, []bool) {
	e.connMu.Lock()
	defer e.connMu.Unlock()
	return append([]ai.MCPConnector(nil), e.conns...), append([]bool(nil), e.bound...)
}

// TestGatewayFlowExecuteRunsAgainstConnectorFrozenAtRunStart is the workflow-level half of the §8.3 freeze (the
// resolver-level half lives in gateway_wiring_internal_test.go): the execute step must bind the connector captured in
// the run's snapshot onto the step context, so editing the connector between the proposal and the approval cannot
// change what the already-approved action executes against. The runner's "live" definition is swapped to a hostile
// config during the approval wait; the executor must still see the config the approver reviewed.
func TestGatewayFlowExecuteRunsAgainstConnectorFrozenAtRunStart(t *testing.T) {
	frozen := ai.MCPConnector{Name: "jira_ops", URL: "https://jira.example/mcp", AllowedHosts: []string{"jira.example"}}
	runner := &postureRunner{conn: frozen} // manual posture: the run pauses, leaving a window to edit the connector
	exec := &connectorRecordingExecutor{fakeExecutor: fakeExecutor{result: act.ExecuteResult{Outcome: act.OutcomeSucceeded, ExternalReference: "PROJ-7"}}}
	e, store := newPostureExecutor(t, runner, exec)

	const inst = "inst-frozen-conn"
	ctx := t.Context()
	runID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID: inst, AgentName: "triage", Prompt: "actua", IdempotencyKey: "posture-frozen-" + uuid.NewString(),
		Actor: act.Actor{Subject: "user:alice"},
	})
	require.NoError(t, err)

	approvalID := act.ApprovalIDForToolCall(runID, "call-1")
	require.Eventually(t, func() bool {
		_, gErr := store.GetApproval(ctx, inst, approvalID)
		return gErr == nil
	}, 15*time.Second, 50*time.Millisecond)

	// The admin edits the connector while the run waits: new target, egress allowlist and posture.
	runner.setConnector(ai.MCPConnector{Name: "jira_ops", URL: "https://evil.example/mcp", AllowedHosts: []string{"evil.example"}, Approval: "auto"})

	_, err = store.ResolveApproval(ctx, inst, approvalID, act.ApprovalStatusApproved, "admin:bob")
	require.NoError(t, err)
	require.NoError(t, e.Resume(ctx, runID, act.ApprovalApproved, "call-1"))
	res, err := e.Result(runID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status)

	conns, bound := exec.boundConnectors()
	require.Len(t, conns, 1)
	require.True(t, bound[0], "the execute step must bind the captured connector onto the step context")
	require.Equal(t, "https://jira.example/mcp", conns[0].URL, "the approved write executes against the connector frozen at run start")
	require.Equal(t, []string{"jira.example"}, conns[0].AllowedHosts)
}
