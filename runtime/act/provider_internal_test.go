package act

import (
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/stretchr/testify/require"
)

// TestSnapshotFromSpecMapsMCP verifies snapshotFromSpec projects an AgentSpec's MCP connectors onto the snapshot
// the executor runs from, carrying every field through and deep-copying the allowed-hosts slice.
func TestSnapshotFromSpecMapsMCP(t *testing.T) {
	spec := &runtimev1.AgentSpec{
		Instructions: "Investigate.",
		Tools:        []string{"query_metrics_view"},
		Mcp: []*runtimev1.MCPConnector{
			{
				Name:              "jira",
				Url:               "https://jira.example.com/mcp",
				AuthSecret:        "JIRA_TOKEN",
				AllowedHosts:      []string{"jira.example.com"},
				TrustReadOnlyHint: true,
				Approval:          "auto",
				RequireApproval:   []string{"delete_*"},
				AutoApprove:       []string{"search_*"},
			},
			{
				Name: "slack",
				Url:  "https://slack.example.com/mcp",
			},
		},
	}

	snap := snapshotFromSpec("incident", "hash-1", spec)

	require.Equal(t, "incident", snap.Name)
	require.Equal(t, "hash-1", snap.SpecHash)
	require.Len(t, snap.MCPConnectors, 2)

	require.Equal(t, "jira", snap.MCPConnectors[0].Name)
	require.Equal(t, "https://jira.example.com/mcp", snap.MCPConnectors[0].URL)
	require.Equal(t, "JIRA_TOKEN", snap.MCPConnectors[0].AuthSecret)
	require.Equal(t, []string{"jira.example.com"}, snap.MCPConnectors[0].AllowedHosts)
	require.True(t, snap.MCPConnectors[0].TrustReadOnlyHint)
	// The approval posture and its glob exceptions must survive the projection, or the executor would fall back to
	// manual-approve-everything (a silent no-op of the author's auto_approve/require_approval declarations).
	require.Equal(t, "auto", snap.MCPConnectors[0].Approval)
	require.Equal(t, []string{"delete_*"}, snap.MCPConnectors[0].RequireApproval)
	require.Equal(t, []string{"search_*"}, snap.MCPConnectors[0].AutoApprove)

	require.Equal(t, "slack", snap.MCPConnectors[1].Name)
	require.Empty(t, snap.MCPConnectors[1].AllowedHosts)
	require.False(t, snap.MCPConnectors[1].TrustReadOnlyHint)
	require.Empty(t, snap.MCPConnectors[1].Approval)

	// The slices are copied, not aliased: mutating the source spec must not reach the snapshot.
	spec.Mcp[0].AllowedHosts[0] = "evil.example.com"
	require.Equal(t, "jira.example.com", snap.MCPConnectors[0].AllowedHosts[0])
	spec.Mcp[0].RequireApproval[0] = "*"
	require.Equal(t, "delete_*", snap.MCPConnectors[0].RequireApproval[0])
	spec.Mcp[0].AutoApprove[0] = "*"
	require.Equal(t, "search_*", snap.MCPConnectors[0].AutoApprove[0])
}

// TestSnapshotFromSpecNoMCP verifies the snapshot has no MCP connectors when the spec declares none.
func TestSnapshotFromSpecNoMCP(t *testing.T) {
	snap := snapshotFromSpec("plain", "hash-2", &runtimev1.AgentSpec{Instructions: "Answer."})
	require.Nil(t, snap.MCPConnectors)
}
