package reconcilers_test

import (
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// TestAgentDeniedThroughRealTransitiveAccess pins the runtime-side link of the confinement that keeps Act
// agents closed to public links, using the REAL expansion path rather than hand-built exclusive rules: the
// claims carry a transitive access rule on an explore (exactly what a magic auth token encodes, see
// securityRulesFromMagicAuthToken in admin/server/runtime_jwt.go), the explore reconciler's
// ResolveTransitiveAccess expands it, and the engine merges the result into an exclusive rule that denies
// everything else. The agent declares access for EVERYONE, so it is only denied if that whole chain holds:
// if an upstream rebase made a reconciler return non-exclusive rules, the agent's own allow rule would win
// and this test would go red.
func TestAgentDeniedThroughRealTransitiveAccess(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"m1.sql":    `SELECT 'foo' as foo, 1 as x`,
			"mv1.yaml": `
version: 1
type: metrics_view
model: m1
dimensions:
- column: foo
measures:
- name: x
  expression: sum(x)
`,
			"e1.yaml": `
type: explore
metrics_view: mv1
`,
			"agents/abierta.yaml": `
type: agent
instructions: Narra el dashboard.
security:
  access: "true"
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	// Claims shaped like a magic auth token for the explore: a transitive access rule, nothing else.
	ctx := t.Context()
	claims := &runtime.SecurityClaims{
		UserAttributes: map[string]any{"admin": false},
		AdditionalRules: []*runtimev1.SecurityRule{
			{
				Rule: &runtimev1.SecurityRule_TransitiveAccess{
					TransitiveAccess: &runtimev1.SecurityRuleTransitiveAccess{
						Resource: &runtimev1.ResourceName{Kind: runtime.ResourceKindExplore, Name: "e1"},
					},
				},
			},
		},
	}

	// Sanity: the expansion really ran and opened the explore and its metrics view.
	e1 := testruntime.GetResource(t, rt, id, runtime.ResourceKindExplore, "e1")
	sec, err := rt.ResolveSecurity(ctx, id, claims, e1)
	require.NoError(t, err)
	require.True(t, sec.CanAccess())
	mv1 := testruntime.GetResource(t, rt, id, runtime.ResourceKindMetricsView, "mv1")
	sec, err = rt.ResolveSecurity(ctx, id, claims, mv1)
	require.NoError(t, err)
	require.True(t, sec.CanAccess())

	// The agent is open to everyone, yet the token cannot see it: the expanded rule is exclusive and denies
	// everything it does not name. This is the property that keeps agents out of share links.
	agent := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "abierta")
	sec, err = rt.ResolveSecurity(ctx, id, claims, agent)
	require.NoError(t, err)
	require.False(t, sec.CanAccess(), "a transitive token rule must confine: if this is accessible, the expansion no longer produces an exclusive deny-rest rule")
}
