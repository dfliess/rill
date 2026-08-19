package runtime_test

import (
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// These tests cover the inverse of the approve gate (kairos-cloud#143): given a project's membership, which of
// those members may decide one proposed action. It is what addresses an approval notification, so it has to agree
// with the gate the deciding caller will hit later (runtime/server/agents.go resolveApprovalDecision).

// approversProjectFiles is a project whose agents cover the shapes an approve policy takes.
func approversProjectFiles() map[string]string {
	return map[string]string{
		"rill.yaml": `
features:
  agents: true
`,
		"cobranza.yaml": `
type: agent
display_name: Cobranza
instructions: "Investiga la morosidad y propone la tarea."
security:
  access: '{{ or (has "operaciones" .user.groups) (has "direccion" .user.groups) }}'
  approve: '{{ has "direccion" .user.groups }}'
`,
		"finanzas.yaml": `
type: agent
display_name: Finanzas
instructions: "Investiga y propone acciones financieras."
security:
  access: "true"
  approve: >-
    {{ or (has "direccion" .user.groups)
          (and (eq .action.tool "issue_refund") (eq .user.email "anafin@example.com")) }}
`,
		"cuatroojos.yaml": `
type: agent
display_name: Cuatro Ojos
instructions: "Propone acciones que otro debe firmar."
security:
  access: "true"
  approve: '{{ ne .user.id .run.actor }}'
`,
		"nadie.yaml": `
type: agent
display_name: Nadie
instructions: "Propone acciones que ningun grupo puede firmar."
security:
  access: "true"
  approve: '{{ has "inexistente" .user.groups }}'
`,
		"abierta.yaml": `
type: agent
display_name: Abierta
instructions: "Todos la ven, solo los operadores la firman."
security:
  access: "true"
`,
		"sinpolitica.yaml": `
type: agent
display_name: Sin Politica
instructions: "No declara seguridad."
`,
	}
}

// approversMembers is the project's membership as the admin service reports it (sorted by email, as
// ListProjectMemberAttributes returns it), with the attributes a runtime JWT would carry for each of them.
func approversMembers() []drivers.ProjectMember {
	member := func(id, email string, groups []any, admin, editTrigger bool) drivers.ProjectMember {
		return drivers.ProjectMember{
			UserID: id,
			Email:  email,
			// NOTE: no "id" attribute. The admin service does not send one; the JWT claims provider defaults it to
			// the token's subject, and so must the claims rebuilt here, or `.user.id` policies would never match.
			Attributes: map[string]any{
				"name":   id,
				"email":  email,
				"domain": "example.com",
				"groups": groups,
				"admin":  admin,
			},
			EditTrigger: editTrigger,
		}
	}

	return []drivers.ProjectMember{
		member("usr_admin", "admin@example.com", []any{}, true, true),
		member("usr_ana", "ana@example.com", []any{"direccion"}, false, false),
		member("usr_anafin", "anafin@example.com", []any{"finanzas"}, false, false),
		// A member reachable through a public link: same group as Ana, but confined by an exclusive rule that names
		// one explore, which is how the admin service restricts a magic auth token.
		{
			UserID: "usr_link",
			Email:  "link@example.com",
			Attributes: map[string]any{
				"name": "link", "email": "link@example.com", "domain": "example.com",
				"groups": []any{"direccion"}, "admin": false,
			},
			SecurityRules: []*runtimev1.SecurityRule{
				{
					Rule: &runtimev1.SecurityRule_Access{
						Access: &runtimev1.SecurityRuleAccess{
							ConditionResources: []*runtimev1.ResourceName{{Kind: runtime.ResourceKindExplore, Name: "shared_dashboard"}},
							Allow:              true,
							Exclusive:          true,
						},
					},
				},
			},
		},
		member("usr_marta", "marta@example.com", []any{"operaciones"}, false, false),
		// A member with no email: there is nothing to notify, so they are skipped even where the policy allows them.
		member("usr_sinmail", "", []any{"direccion"}, false, false),
	}
}

func approversInstance(t *testing.T) (*runtime.Runtime, string) {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: approversProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	return rt, instanceID
}

func approversAgent(t *testing.T, rt *runtime.Runtime, instanceID, name string) *runtimev1.Resource {
	t.Helper()
	ctrl, err := rt.Controller(t.Context(), instanceID)
	require.NoError(t, err)
	res, err := ctrl.Get(t.Context(), &runtimev1.ResourceName{Kind: runtime.ResourceKindAgent, Name: name}, false)
	require.NoError(t, err)
	require.NotNil(t, res.GetAgent().State.ValidSpec, "agent %q must have reconciled", name)
	return res
}

// TestResolveAgentApprovers covers who is notified for each policy shape: a group policy, a policy that
// discriminates by tool, the four-eyes shape that reads the run's actor, and the two shapes that resolve to the
// operators (an absent approve expression, and one that matches nobody).
func TestResolveAgentApprovers(t *testing.T) {
	rt, instanceID := approversInstance(t)
	members := approversMembers()

	tests := []struct {
		name   string
		agent  string
		action runtime.AgentActionContext
		want   []string
	}{
		{
			// The policy names a group, so authority is exactly that group: neither the operator (who could still
			// break the glass) nor the member the policy does not name is notified.
			name:  "group_policy",
			agent: "cobranza",
			want:  []string{"ana@example.com"},
		},
		{
			// The run's actor is not excluded: if whoever launched the run may sign, they are told about it.
			name:   "actor_not_excluded",
			agent:  "cobranza",
			action: runtime.AgentActionContext{RunActor: "usr_ana"},
			want:   []string{"ana@example.com"},
		},
		{
			// Same agent, same members, different tool: authority discriminates by action, so the sets differ.
			name:   "tool_discriminates_refund",
			agent:  "finanzas",
			action: runtime.AgentActionContext{Tool: "issue_refund", Connector: "jira_ops"},
			want:   []string{"ana@example.com", "anafin@example.com"},
		},
		{
			name:   "tool_discriminates_other",
			agent:  "finanzas",
			action: runtime.AgentActionContext{Tool: "close_account", Connector: "jira_ops"},
			want:   []string{"ana@example.com"},
		},
		{
			// `.run.actor` is bound, and `.user.id` resolves even though no member carries an "id" attribute.
			name:   "four_eyes_excludes_the_actor",
			agent:  "cuatroojos",
			action: runtime.AgentActionContext{RunActor: "usr_ana"},
			want:   []string{"admin@example.com", "anafin@example.com", "marta@example.com"},
		},
		{
			// Nobody matches: fall back to the operators rather than notify nobody and strand the run.
			name:  "fallback_to_operators",
			agent: "nadie",
			want:  []string{"admin@example.com"},
		},
		{
			// An agent everyone sees but only the operators sign: an absent approve expression resolves to the
			// EditTrigger permission, which is the one thing that differs between members' tokens.
			name:  "absent_policy_is_operators",
			agent: "abierta",
			want:  []string{"admin@example.com"},
		},
		{
			// No security block at all: no access for anyone but an admin, and the operators sign.
			name:  "no_security_block",
			agent: "sinpolitica",
			want:  []string{"admin@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rt.ResolveAgentApprovers(t.Context(), instanceID, approversAgent(t, rt, instanceID, tt.agent), tt.action, members)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestResolveAgentApproversHonorsMemberRestrictions pins that a member's own security rules are applied: the
// public-link member is in the group the policy names, but their token's exclusive rule denies the agent, and
// deciding is a subset of access. Without this, an approval notification would reach a share link's audience.
func TestResolveAgentApproversHonorsMemberRestrictions(t *testing.T) {
	rt, instanceID := approversInstance(t)
	res := approversAgent(t, rt, instanceID, "cobranza")

	got, err := rt.ResolveAgentApprovers(t.Context(), instanceID, res, runtime.AgentActionContext{}, approversMembers())
	require.NoError(t, err)
	require.NotContains(t, got, "link@example.com")

	// The same member without the restriction is an approver, so it is the rule that excludes them, not their attributes.
	unrestricted := approversMembers()
	for i := range unrestricted {
		if unrestricted[i].Email == "link@example.com" {
			unrestricted[i].SecurityRules = nil
		}
	}
	got, err = rt.ResolveAgentApprovers(t.Context(), instanceID, res, runtime.AgentActionContext{}, unrestricted)
	require.NoError(t, err)
	require.Contains(t, got, "link@example.com")
}

// TestResolveAgentApproversNoMembers checks the degenerate input: a project the admin service reports no members
// for resolves to nobody, not to an error.
func TestResolveAgentApproversNoMembers(t *testing.T) {
	rt, instanceID := approversInstance(t)

	got, err := rt.ResolveAgentApprovers(t.Context(), instanceID, approversAgent(t, rt, instanceID, "cobranza"), runtime.AgentActionContext{}, nil)
	require.NoError(t, err)
	require.Empty(t, got)
}
