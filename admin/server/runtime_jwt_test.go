package server

import (
	"testing"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/runtime"
	"github.com/stretchr/testify/require"
)

// TestSecurityRulesFromMagicAuthToken pins the admin-side link of the confinement that keeps Act agents (and
// any other unlisted resource) closed to public links: a magic auth token's rules are either a blanket deny
// (no resources) or transitive access rules, one per named resource, and nothing else that grants access. The
// runtime expands each transitive rule into an exclusive access rule that denies everything it does not name
// (see expandTransitiveAccessRules in runtime/security.go), so if this function ever emitted a plain allow, or
// stopped emitting rules at all, share links would silently stop being confined.
func TestSecurityRulesFromMagicAuthToken(t *testing.T) {
	// No resources: exactly one rule, an unconditional deny.
	rules, err := securityRulesFromMagicAuthToken(&database.MagicAuthToken{})
	require.NoError(t, err)
	require.Len(t, rules, 1)
	access := rules[0].GetAccess()
	require.NotNil(t, access)
	require.False(t, access.Allow)

	// With resources: one transitive access rule per resource, and no other access-granting rule. The field
	// allow-list that rides along is not an access rule (and stays exclusive).
	mdl := &database.MagicAuthToken{
		Resources: []database.ResourceName{
			{Type: runtime.ResourceKindExplore, Name: "e1"},
			{Type: runtime.ResourceKindCanvas, Name: "c1"},
		},
		Fields: []string{"foo"},
	}
	rules, err = securityRulesFromMagicAuthToken(mdl)
	require.NoError(t, err)
	require.Len(t, rules, 3)
	for i, res := range mdl.Resources {
		transitive := rules[i].GetTransitiveAccess()
		require.NotNil(t, transitive, "resource rules must be transitive, never plain allows")
		require.Equal(t, res.Type, transitive.Resource.Kind)
		require.Equal(t, res.Name, transitive.Resource.Name)
	}
	require.Nil(t, rules[2].GetAccess())
	fieldAccess := rules[2].GetFieldAccess()
	require.NotNil(t, fieldAccess)
	require.True(t, fieldAccess.Exclusive)
}
