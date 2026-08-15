package runtime

import (
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/parser"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// agentResource builds a reconciled Agent resource with the given security fields on its valid spec.
func agentResource(name string, rules []*runtimev1.SecurityRule) *runtimev1.Resource {
	spec := &runtimev1.AgentSpec{
		Instructions:  "Investigate using governed metrics only.",
		SecurityRules: rules,
	}
	return &runtimev1.Resource{
		Meta: &runtimev1.ResourceMeta{
			Name:           &runtimev1.ResourceName{Kind: ResourceKindAgent, Name: name},
			StateUpdatedOn: timestamppb.Now(),
		},
		Resource: &runtimev1.Resource_Agent{
			Agent: &runtimev1.Agent{
				Spec:  spec,
				State: &runtimev1.AgentState{ValidSpec: spec},
			},
		},
	}
}

func accessRule(condition string) *runtimev1.SecurityRule {
	return &runtimev1.SecurityRule{
		Rule: &runtimev1.SecurityRule_Access{
			Access: &runtimev1.SecurityRuleAccess{ConditionExpression: condition, Allow: true},
		},
	}
}

// TestResolveSecurityAgent covers the built-in rule for the Agent kind: admins always have access, declared
// rules govern everyone else, and an agent that declares none is admin-only.
func TestResolveSecurityAgent(t *testing.T) {
	adminAttrs := map[string]any{"name": "ana", "email": "ana@example.com", "groups": []any{"direccion"}, "admin": true}
	opsAttrs := map[string]any{"name": "marta", "email": "marta@example.com", "groups": []any{"operaciones"}, "admin": false}
	salesAttrs := map[string]any{"name": "sam", "email": "sam@example.com", "groups": []any{"ventas"}, "admin": false}

	groupRule := accessRule(`{{ or (has "operaciones" .user.groups) (has "direccion" .user.groups) }}`)

	tests := []struct {
		name       string
		attrs      map[string]any
		rules      []*runtimev1.SecurityRule
		wantAccess bool
	}{
		{"no_rules_admin", adminAttrs, nil, true},
		{"no_rules_non_admin", opsAttrs, nil, false},
		{"group_rule_member", opsAttrs, []*runtimev1.SecurityRule{groupRule}, true},
		{"group_rule_non_member", salesAttrs, []*runtimev1.SecurityRule{groupRule}, false},
		// The admin allow replaces the declared rules, so even a deny-all policy keeps the agent visible to
		// the people who administer (and can read) the project that defines it.
		{"deny_all_admin", adminAttrs, []*runtimev1.SecurityRule{
			{Rule: &runtimev1.SecurityRule_Access{Access: &runtimev1.SecurityRuleAccess{Allow: false}}},
		}, true},
		{"deny_all_non_admin", opsAttrs, []*runtimev1.SecurityRule{
			{Rule: &runtimev1.SecurityRule_Access{Access: &runtimev1.SecurityRuleAccess{Allow: false}}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newSecurityEngine(1, zap.NewNop(), nil)
			claims := &SecurityClaims{UserAttributes: tt.attrs}
			got, err := p.resolveSecurity(t.Context(), "", "test", nil, claims, agentResource("triage", tt.rules))
			require.NoError(t, err)
			require.Equal(t, tt.wantAccess, got.CanAccess())
		})
	}
}

// exclusiveClaims returns claims shaped like a magic auth (public link) token's: an exclusive access rule that
// allows only the named resources and denies everything else, which is how the admin server confines them.
func exclusiveClaims(attrs map[string]any, resources ...*runtimev1.ResourceName) *SecurityClaims {
	return &SecurityClaims{
		UserAttributes: attrs,
		Permissions:    []Permission{ReadAPI, ReadMetrics, ReadObjects, UseAI},
		AdditionalRules: []*runtimev1.SecurityRule{
			{
				Rule: &runtimev1.SecurityRule_Access{
					Access: &runtimev1.SecurityRuleAccess{
						ConditionResources: resources,
						Allow:              true,
						Exclusive:          true,
					},
				},
			},
		},
	}
}

// TestAgentDeniedToExclusiveTokenClaims verifies that a public-link/embed token's exclusive rule denies agents:
// the token's allow-list names an explore, so the agent falls in the "everything else" the rule denies. This
// holds even when the agent's own rules would allow everyone, and even when the token carries an admin
// attribute, because the built-in agent rule never merges into the exclusive rule (unlike alert/report).
func TestAgentDeniedToExclusiveTokenClaims(t *testing.T) {
	allowEveryone := []*runtimev1.SecurityRule{accessRule("true")}
	explore := &runtimev1.ResourceName{Kind: ResourceKindExplore, Name: "shared_dashboard"}

	p := newSecurityEngine(4, zap.NewNop(), nil)

	claims := exclusiveClaims(map[string]any{"admin": false}, explore)
	got, err := p.resolveSecurity(t.Context(), "", "test", nil, claims, agentResource("triage", allowEveryone))
	require.NoError(t, err)
	require.False(t, got.CanAccess(), "an exclusive token rule must deny an agent it does not name")

	adminClaims := exclusiveClaims(map[string]any{"admin": true}, explore)
	got, err = p.resolveSecurity(t.Context(), "", "test", nil, adminClaims, agentResource("triage", allowEveryone))
	require.NoError(t, err)
	require.False(t, got.CanAccess(), "the built-in admin allow must not merge into (or override) the exclusive rule")
}

// TestExpandTransitiveAccessRulesMergesExclusiveRules pins the merge step of the expansion: multiple exclusive
// access rules collapse into exactly one (so evaluation order cannot produce false rejections) and the merged
// rule keeps the exclusive "denies access to everything else" semantics. The other link the confinement
// depends on, that TRANSITIVE rules expand into exclusive ones in the first place, is pinned by
// TestSecurityRulesFromMagicAuthToken (admin/server) and TestAgentDeniedThroughRealTransitiveAccess
// (runtime/reconcilers), which exercise real magic-auth-shaped inputs end to end.
func TestExpandTransitiveAccessRulesMergesExclusiveRules(t *testing.T) {
	explore := &runtimev1.ResourceName{Kind: ResourceKindExplore, Name: "shared_dashboard"}
	mv := &runtimev1.ResourceName{Kind: ResourceKindMetricsView, Name: "orders"}
	claims := &SecurityClaims{
		AdditionalRules: []*runtimev1.SecurityRule{
			{Rule: &runtimev1.SecurityRule_Access{Access: &runtimev1.SecurityRuleAccess{
				ConditionResources: []*runtimev1.ResourceName{explore}, Allow: true, Exclusive: true,
			}}},
			{Rule: &runtimev1.SecurityRule_Access{Access: &runtimev1.SecurityRuleAccess{
				ConditionResources: []*runtimev1.ResourceName{mv}, Allow: true, Exclusive: true,
			}}},
		},
	}

	p := newSecurityEngine(1, zap.NewNop(), nil)
	rules, err := p.expandTransitiveAccessRules(t.Context(), "", claims)
	require.NoError(t, err)

	// The exclusive rules merge into exactly one, so evaluation order cannot produce false rejections.
	require.Len(t, rules, 1)
	access := rules[0].GetAccess()
	require.NotNil(t, access)
	require.True(t, access.Allow)
	require.True(t, access.Exclusive, "the merged rule must stay exclusive: it is what denies everything else")
	require.ElementsMatch(t, []*runtimev1.ResourceName{explore, mv}, access.ConditionResources)

	// And the exclusive semantics deny a resource outside the list: this is the "denies access to everything
	// else" behavior (see the expansion in expandTransitiveAccessRules) that confines share links and embeds.
	res := &ResolvedSecurity{}
	err = p.applySecurityRuleAccess(res, agentResource("triage", nil), access, parser.TemplateData{})
	require.NoError(t, err)
	require.False(t, res.CanAccess())
}
