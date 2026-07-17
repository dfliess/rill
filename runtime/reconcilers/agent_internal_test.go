package reconcilers

import (
	"context"
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/stretchr/testify/require"
)

// specWithMCP builds a minimal agent spec carrying a single MCP connector, so the tests below can vary one field
// at a time. It declares no model connector or tools, so validateSpec and specHash never touch the controller.
func specWithMCP(m *runtimev1.MCPConnector) *runtimev1.AgentSpec {
	return &runtimev1.AgentSpec{
		Instructions: "Investigate.",
		Mcp:          []*runtimev1.MCPConnector{m},
	}
}

// TestSpecHashSensitiveToMCP verifies the hash folds in each MCP connector field: changing name, url, auth secret,
// allowed hosts or the trust flag changes the hash (so editing a connector rebinds the run's audit hash), while an
// equivalent spec hashes the same and reordering an allowlist does not churn it.
func TestSpecHashSensitiveToMCP(t *testing.T) {
	r := &AgentReconciler{}

	base := specWithMCP(&runtimev1.MCPConnector{
		Name:              "jira",
		Url:               "https://jira.example.com/mcp",
		AuthSecret:        "JIRA_TOKEN",
		AllowedHosts:      []string{"jira.example.com", "cdn.example.com"},
		TrustReadOnlyHint: false,
	})

	baseHash, err := r.specHash(base)
	require.NoError(t, err)

	// An equivalent spec (fresh objects, same values) hashes identically.
	sameHash, err := r.specHash(specWithMCP(&runtimev1.MCPConnector{
		Name:         "jira",
		Url:          "https://jira.example.com/mcp",
		AuthSecret:   "JIRA_TOKEN",
		AllowedHosts: []string{"jira.example.com", "cdn.example.com"},
	}))
	require.NoError(t, err)
	require.Equal(t, baseHash, sameHash)

	// Reordering an equivalent allowlist does not change the hash (hosts are hashed in a stable order).
	reorderedHash, err := r.specHash(specWithMCP(&runtimev1.MCPConnector{
		Name:         "jira",
		Url:          "https://jira.example.com/mcp",
		AuthSecret:   "JIRA_TOKEN",
		AllowedHosts: []string{"cdn.example.com", "jira.example.com"},
	}))
	require.NoError(t, err)
	require.Equal(t, baseHash, reorderedHash)

	// Each of these single-field edits must change the hash.
	cases := map[string]*runtimev1.MCPConnector{
		"name":     {Name: "confluence", Url: "https://jira.example.com/mcp", AuthSecret: "JIRA_TOKEN", AllowedHosts: []string{"jira.example.com", "cdn.example.com"}},
		"url":      {Name: "jira", Url: "https://other.example.com/mcp", AuthSecret: "JIRA_TOKEN", AllowedHosts: []string{"jira.example.com", "cdn.example.com"}},
		"secret":   {Name: "jira", Url: "https://jira.example.com/mcp", AuthSecret: "OTHER_TOKEN", AllowedHosts: []string{"jira.example.com", "cdn.example.com"}},
		"hosts":    {Name: "jira", Url: "https://jira.example.com/mcp", AuthSecret: "JIRA_TOKEN", AllowedHosts: []string{"jira.example.com"}},
		"trustFlag": {Name: "jira", Url: "https://jira.example.com/mcp", AuthSecret: "JIRA_TOKEN", AllowedHosts: []string{"jira.example.com", "cdn.example.com"}, TrustReadOnlyHint: true},
	}
	for name, m := range cases {
		h, err := r.specHash(specWithMCP(m))
		require.NoError(t, err)
		require.NotEqual(t, baseHash, h, "hash should change when %s changes", name)
	}
}

// TestValidateSpecRejectsNonHTTPSMCP verifies fail-closed validation: an MCP connector with a plaintext http URL is
// rejected (a tenant connector must reach a public https host), while the same connector over https validates.
func TestValidateSpecRejectsNonHTTPSMCP(t *testing.T) {
	r := &AgentReconciler{}
	ctx := context.Background()

	err := r.validateSpec(ctx, specWithMCP(&runtimev1.MCPConnector{
		Name: "jira",
		Url:  "http://jira.example.com/mcp",
	}))
	require.Error(t, err)
	require.Contains(t, err.Error(), "jira")

	err = r.validateSpec(ctx, specWithMCP(&runtimev1.MCPConnector{
		Name: "jira",
		Url:  "https://jira.example.com/mcp",
	}))
	require.NoError(t, err)
}
