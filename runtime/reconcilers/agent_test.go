package reconcilers_test

import (
	"testing"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

func TestAgent(t *testing.T) {
	// The model connector references "duckdb", which is a configured connector in the test instance.
	// The reconciler only validates that the connector exists; it does not execute the agent.
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
display_name: Revenue Incident Agent
description: Investigates revenue anomalies.
model:
  connector: duckdb
  name: some-model
instructions: Investigate using governed metrics only.
tools:
  - query_metrics_view
limits:
  max_steps: 5
  timeout: 10m
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// No reconcile errors and no parse errors.
	testruntime.RequireReconcileState(t, rt, id, -1, 0, 0)

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "incident")
	agent := res.GetAgent()
	require.NotNil(t, agent)
	require.Empty(t, res.Meta.ReconcileError)

	// The reconciler records a valid snapshot and a stable hash.
	require.NotNil(t, agent.State.ValidSpec)
	require.NotEmpty(t, agent.State.SpecHash)
	require.Equal(t, "Investigate using governed metrics only.", agent.State.ValidSpec.Instructions)
	require.Equal(t, "duckdb", agent.State.ValidSpec.ModelConnector)
	require.Equal(t, uint32(5), agent.State.ValidSpec.Limits.MaxSteps)
	require.Equal(t, uint32(600), agent.State.ValidSpec.Limits.TimeoutSeconds)
	require.Equal(t, []string{"query_metrics_view"}, agent.State.ValidSpec.Tools)
}

func TestAgentDisallowedTool(t *testing.T) {
	// A v1 dynamic agent may only declare read-only analytical tools. Declaring a write/development tool
	// (here write_file) must fail validation fail-closed, leaving no valid spec behind.
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/writer.yaml": `
type: agent
instructions: Investigate and then fix the file.
tools:
  - write_file
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// Exactly one reconcile error (the agent), no parse errors.
	testruntime.RequireReconcileState(t, rt, id, -1, 1, 0)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgent, "writer", "is not allowed")

	// On validation failure, no valid spec is recorded: the agent is not executable.
	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "writer")
	require.Nil(t, res.GetAgent().State.ValidSpec)
}

func TestAgentDisallowedMCPNamespacedTool(t *testing.T) {
	// The tools: list is a CLOSED allowlist of built-in read-only analytical tools. An MCP tool must never be
	// declarable there: MCP tools are discovered live from the mcp: connectors (gated by approval), so accepting a
	// namespaced name in tools: would smuggle an action tool past the connector's approval posture.
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/smuggler.yaml": `
type: agent
instructions: Investigate and file a ticket.
tools:
  - mcp.jira.create_issue
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, -1, 1, 0)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgent, "smuggler", "is not allowed")

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "smuggler")
	require.Nil(t, res.GetAgent().State.ValidSpec)
}

func TestAgentMCPConnector(t *testing.T) {
	// A valid MCP connector survives the full pipeline (parser to reconciler) and lands in the valid spec.
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/support.yaml": `
type: agent
instructions: Triage support tickets.
mcp:
  - name: jira
    url: https://jira.example.com/mcp
    auth:
      secret: JIRA_TOKEN
    network:
      allowed_hosts:
        - jira.example.com
    trust_read_only_hint: true
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, -1, 0, 0)

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "support")
	agent := res.GetAgent()
	require.NotNil(t, agent.State.ValidSpec)
	require.Len(t, agent.State.ValidSpec.Mcp, 1)
	require.Equal(t, "jira", agent.State.ValidSpec.Mcp[0].Name)
	require.Equal(t, "https://jira.example.com/mcp", agent.State.ValidSpec.Mcp[0].Url)
	require.Equal(t, "JIRA_TOKEN", agent.State.ValidSpec.Mcp[0].AuthSecret)
	require.Equal(t, []string{"jira.example.com"}, agent.State.ValidSpec.Mcp[0].AllowedHosts)
	require.True(t, agent.State.ValidSpec.Mcp[0].TrustReadOnlyHint)
}

func TestAgentMCPConnectorInvalid(t *testing.T) {
	// A plaintext http MCP URL fails validation fail-closed, so no valid spec is recorded.
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/bad.yaml": `
type: agent
instructions: Do something.
mcp:
  - name: jira
    url: http://jira.example.com/mcp
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, -1, 1, 0)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgent, "bad", "invalid mcp connector")

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "bad")
	require.Nil(t, res.GetAgent().State.ValidSpec)
}

func TestAgentInvalidModelConnector(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/bad.yaml": `
type: agent
model:
  connector: does_not_exist
instructions: Do something useful.
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// Exactly one reconcile error (the agent), no parse errors.
	testruntime.RequireReconcileState(t, rt, id, -1, 1, 0)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgent, "bad", "invalid model connector")

	// On validation failure, no valid spec is recorded.
	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "bad")
	require.Nil(t, res.GetAgent().State.ValidSpec)
}
