package act

import (
	"context"
	"errors"
	"testing"

	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

func TestConnectorOf(t *testing.T) {
	cases := []struct {
		tool      string
		connector string
		ok        bool
	}{
		{"mcp.mockmcp.HelloWorld", "mockmcp", true},
		{"mcp.jira.create_issue", "jira", true},
		{"list_metrics_views", "", false}, // a built-in analytical tool is not an MCP action tool
		{"mcp.only", "", false},           // connector present but no tool segment
		{"mcp.", "", false},               // empty connector
	}
	for _, c := range cases {
		got, ok := connectorOf(c.tool)
		require.Equal(t, c.ok, ok, c.tool)
		require.Equal(t, c.connector, got, c.tool)
	}
}

func TestCapturedProposerSurfacesCapture(t *testing.T) {
	// No capture: a pure read-only investigation proposes nothing, so the run succeeds without an action.
	_, ok, err := capturedProposer{}.Propose(context.Background(), ProposeInput{})
	require.NoError(t, err)
	require.False(t, ok)

	// A captured write is surfaced verbatim as the structured proposal the gateway will govern. Its identity is the
	// model's own tool-call ID (the captured call message ID), so two calls in one run key distinct ledger rows.
	captured := &ai.ProposedAction{
		ToolCallID: "call-msg-7f3a",
		Connector:  "mockmcp",
		Tool:       "mcp.mockmcp.HelloWorld",
		Args:       map[string]any{"name": "Marta"},
		Summary:    "HelloWorld(name=Marta)",
	}
	prop, ok, err := capturedProposer{}.Propose(context.Background(), ProposeInput{Captured: captured})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "mcp.mockmcp.HelloWorld", prop.Tool)
	require.Equal(t, "mockmcp", prop.Connector)
	require.Equal(t, "Marta", prop.Args["name"])
	require.Equal(t, "HelloWorld(name=Marta)", prop.Summary)
	require.Equal(t, "call-msg-7f3a", prop.ToolCallID, "the proposal keys on the model's own tool-call ID")

	// A capture with no ID falls back to the effective tool name so the ledger identity is never empty.
	fallback, ok, err := capturedProposer{}.Propose(context.Background(), ProposeInput{Captured: &ai.ProposedAction{Tool: "mcp.mockmcp.HelloWorld"}})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "mcp.mockmcp.HelloWorld", fallback.ToolCallID)
}

func TestConnectorAutoApproves(t *testing.T) {
	cases := []struct {
		name string
		conn ai.MCPConnector
		raw  string
		want bool
	}{
		{"manual default needs approval", ai.MCPConnector{}, "create_issue", false},
		{"manual + auto_approve exact", ai.MCPConnector{Approval: "manual", AutoApprove: []string{"add_note"}}, "add_note", true},
		{"manual + auto_approve glob hit", ai.MCPConnector{AutoApprove: []string{"create_*"}}, "create_issue", true},
		{"manual + auto_approve glob miss", ai.MCPConnector{AutoApprove: []string{"create_*"}}, "delete_issue", false},
		{"auto approves everything", ai.MCPConnector{Approval: "auto"}, "delete_issue", true},
		{"auto + require_approval gates the match", ai.MCPConnector{Approval: "auto", RequireApproval: []string{"delete_*"}}, "delete_issue", false},
		{"auto + require_approval lets others through", ai.MCPConnector{Approval: "auto", RequireApproval: []string{"delete_*"}}, "create_issue", true},
	}
	for _, c := range cases {
		require.Equal(t, c.want, connectorAutoApproves(c.conn, c.raw), c.name)
	}
}

// TestClassifyRespectsAuthoritativePreDispatchFailed proves that classify honours an executor's explicit OutcomeFailed
// even when the call returns an error: the executor is authoritative on pre-dispatch failures (schema drift, tool
// gone) where the request provably never left the client, so degrading to indeterminate would force a human to resolve
// what the executor already resolved (§12).
func TestClassifyRespectsAuthoritativePreDispatchFailed(t *testing.T) {
	g := &Gateway{} // classify needs no wiring: it is a pure classification of the executor's report

	cases := []struct {
		name    string
		res     ExecuteResult
		err     error
		outcome ExecuteOutcome
	}{
		{
			name:    "pre-dispatch failed with error stays failed",
			res:     ExecuteResult{Outcome: OutcomeFailed},
			err:     errors.New("act: mcp call rejected before dispatch: tool gone"),
			outcome: OutcomeFailed,
		},
		{
			name:    "error without explicit outcome is indeterminate",
			res:     ExecuteResult{},
			err:     errors.New("act: mcp call failed: connection reset"),
			outcome: OutcomeIndeterminate,
		},
		{
			name:    "error with empty outcome string is indeterminate",
			res:     ExecuteResult{Outcome: ""},
			err:     errors.New("transport error"),
			outcome: OutcomeIndeterminate,
		},
		{
			name:    "no error succeeded",
			res:     ExecuteResult{Outcome: OutcomeSucceeded},
			err:     nil,
			outcome: OutcomeSucceeded,
		},
		{
			name:    "no error failed",
			res:     ExecuteResult{Outcome: OutcomeFailed},
			err:     nil,
			outcome: OutcomeFailed,
		},
		{
			name:    "no error indeterminate",
			res:     ExecuteResult{Outcome: OutcomeIndeterminate},
			err:     nil,
			outcome: OutcomeIndeterminate,
		},
		{
			name:    "no error unrecognised outcome is indeterminate",
			res:     ExecuteResult{Outcome: "bogus"},
			err:     nil,
			outcome: OutcomeIndeterminate,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.outcome, g.classify(c.res, c.err))
		})
	}
}

func TestRawToolName(t *testing.T) {
	require.Equal(t, "HelloWorld", rawToolName("mcp.mockmcp.HelloWorld", "mockmcp"))
	require.Equal(t, "create_issue", rawToolName("mcp.jira.create_issue", "jira"))
}

// editableAgentProvider is an AgentDefinitionProvider whose live snapshot a test can set to simulate an admin editing
// the agent's MCP connector between a run's proposal and its approval. It counts GetAgent calls so a test can assert the
// governed execute path never re-read the live definition.
type editableAgentProvider struct {
	snap  *ai.AgentSnapshot
	calls int
}

var _ ai.AgentDefinitionProvider = (*editableAgentProvider)(nil)

func (p *editableAgentProvider) GetAgent(context.Context, string, string) (*ai.AgentSnapshot, error) {
	p.calls++
	return p.snap, nil
}

func (p *editableAgentProvider) ListAgents(context.Context, string) ([]*ai.AgentSnapshot, error) {
	p.calls++
	return []*ai.AgentSnapshot{p.snap}, nil
}

// TestConnectorForUsesCapturedSnapshotNotLive is the Codex #1 regression guard: when the run's captured connector is
// bound onto the context (the governed execute/verify path), connectorFor resolves against THAT connector, not the live
// agent definition. It simulates an admin editing the connector during the approval wait — pointing it at a different
// host, egress allowlist and approval posture — and asserts the reviewed connector still decides the write (§8.3), and
// that the live definition was never read.
func TestConnectorForUsesCapturedSnapshotNotLive(t *testing.T) {
	// The LIVE definition after the admin's edit: a different target, egress allowlist and approval posture.
	provider := &editableAgentProvider{snap: &ai.AgentSnapshot{
		Name: "triage",
		MCPConnectors: []ai.MCPConnector{{
			Name: "jira_ops", URL: "https://evil.example/mcp", AllowedHosts: []string{"evil.example"},
			Approval: "auto", TrustReadOnlyHint: true,
		}},
	}}
	// rt is unused here: the captured connector declares no secret, so connectorFor never touches the instance.
	r := &mcpToolResolver{provider: provider}

	// The connector the approver reviewed, frozen at run start and bound onto the step context by the workflow.
	captured := ai.MCPConnector{
		Name: "jira_ops", URL: "https://jira.example/mcp", AllowedHosts: []string{"jira.example"},
		Approval: "manual", TrustReadOnlyHint: false,
	}
	ctx := withMCPConnector(context.Background(), captured)

	conn, token, err := r.connectorFor(ctx, "inst-1", "triage", "jira_ops")
	require.NoError(t, err)
	require.Empty(t, token)
	// The write resolves against the REVIEWED connector, not the live-edited one.
	require.Equal(t, "https://jira.example/mcp", conn.URL)
	require.Equal(t, []string{"jira.example"}, conn.AllowedHosts)
	require.Equal(t, "manual", conn.Approval)
	require.False(t, conn.TrustReadOnlyHint)
	require.Zero(t, provider.calls, "the governed path must not re-read the live agent definition")
}

// TestConnectorForFallsBackToLiveWithoutCapture verifies that a path which does not bind a captured connector (outside
// the segmented workflow) keeps the pre-existing behavior: connectorFor resolves the live agent definition.
func TestConnectorForFallsBackToLiveWithoutCapture(t *testing.T) {
	provider := &editableAgentProvider{snap: &ai.AgentSnapshot{
		Name:          "triage",
		MCPConnectors: []ai.MCPConnector{{Name: "jira_ops", URL: "https://jira.example/mcp", AllowedHosts: []string{"jira.example"}}},
	}}
	r := &mcpToolResolver{provider: provider}

	conn, _, err := r.connectorFor(context.Background(), "inst-1", "triage", "jira_ops")
	require.NoError(t, err)
	require.Equal(t, "https://jira.example/mcp", conn.URL)
	require.Equal(t, 1, provider.calls, "with no captured connector the resolver reads the live definition")
}

// TestConnectorForIgnoresCapturedConnectorForDifferentName verifies the capture only satisfies a lookup for its OWN
// connector name: a bound connector for a different tool must not be returned in place of the requested one; the
// resolver falls back to the live definition.
func TestConnectorForIgnoresCapturedConnectorForDifferentName(t *testing.T) {
	provider := &editableAgentProvider{snap: &ai.AgentSnapshot{
		Name:          "triage",
		MCPConnectors: []ai.MCPConnector{{Name: "jira_ops", URL: "https://jira.example/mcp"}},
	}}
	r := &mcpToolResolver{provider: provider}

	ctx := withMCPConnector(context.Background(), ai.MCPConnector{Name: "slack_ops", URL: "https://slack.example/mcp"})
	conn, _, err := r.connectorFor(ctx, "inst-1", "triage", "jira_ops")
	require.NoError(t, err)
	require.Equal(t, "https://jira.example/mcp", conn.URL)
	require.Equal(t, 1, provider.calls)
}

// TestConnectorForResolvesSecretLiveWithFrozenConfig pins the intended split: the captured connector's CONFIG (URL,
// allowed_hosts, approval, secret NAME) is frozen from the snapshot, but its bearer SECRET is resolved LIVE from the
// instance variables, so a rotated token is picked up and the credential value is never frozen. The live agent points
// its connector at a different, unset secret name, so a regression that read the live connector would resolve empty and
// fail closed instead of returning the live token.
func TestConnectorForResolvesSecretLiveWithFrozenConfig(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Variables: map[string]string{"jira_token": "live-bearer-value"},
	})
	r := newMCPToolResolver(rt)
	r.provider = &editableAgentProvider{snap: &ai.AgentSnapshot{
		Name:          "triage",
		MCPConnectors: []ai.MCPConnector{{Name: "jira_ops", URL: "https://evil.example/mcp", AuthSecret: "rotated_away_token"}},
	}}

	captured := ai.MCPConnector{Name: "jira_ops", URL: "https://jira.example/mcp", AuthSecret: "jira_token"}
	conn, token, err := r.connectorFor(withMCPConnector(t.Context(), captured), instanceID, "triage", "jira_ops")
	require.NoError(t, err)
	require.Equal(t, "https://jira.example/mcp", conn.URL, "connector config is frozen from the reviewed snapshot")
	require.Equal(t, "jira_token", conn.AuthSecret, "the frozen secret NAME comes from the snapshot")
	require.Equal(t, "live-bearer-value", token, "the bearer secret is resolved live from the instance variables")
}
