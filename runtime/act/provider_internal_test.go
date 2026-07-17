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

	require.Equal(t, "slack", snap.MCPConnectors[1].Name)
	require.Empty(t, snap.MCPConnectors[1].AllowedHosts)
	require.False(t, snap.MCPConnectors[1].TrustReadOnlyHint)

	// The allowed-hosts slice is copied, not aliased: mutating the source spec must not reach the snapshot.
	spec.Mcp[0].AllowedHosts[0] = "evil.example.com"
	require.Equal(t, "jira.example.com", snap.MCPConnectors[0].AllowedHosts[0])
}

// TestSnapshotFromSpecNoMCP verifies the snapshot has no MCP connectors when the spec declares none.
func TestSnapshotFromSpecNoMCP(t *testing.T) {
	snap := snapshotFromSpec("plain", "hash-2", &runtimev1.AgentSpec{Instructions: "Answer."})
	require.Nil(t, snap.MCPConnectors)
}
