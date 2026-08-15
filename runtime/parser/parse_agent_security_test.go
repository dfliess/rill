package parser

import (
	"context"
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/stretchr/testify/require"
)

// TestAgentSecurityBlock covers the happy path of the security block: access becomes security rules on the
// spec (same grammar as a dashboard policy), and launch/approve ride as validated gate expressions.
func TestAgentSecurityBlock(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		`agents/cobranza.yaml`: `
type: agent
instructions: Investiga con las métricas gobernadas.
security:
  access: '{{ or (has "operaciones" .user.groups) (has "direccion" .user.groups) }}'
  launch: '{{ has "operaciones" .user.groups }}'
  approve: '{{ or (has "direccion" .user.groups) (eq .action.tool "noop") }}'
`,
	})
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	require.Empty(t, p.Errors)

	spec := p.Resources[ResourceName{Kind: ResourceKindAgent, Name: "cobranza"}].AgentSpec
	require.Len(t, spec.SecurityRules, 1)
	access := spec.SecurityRules[0].GetAccess()
	require.NotNil(t, access)
	require.True(t, access.Allow)
	require.Equal(t, `{{ or (has "operaciones" .user.groups) (has "direccion" .user.groups) }}`, access.ConditionExpression)
	require.Equal(t, `{{ has "operaciones" .user.groups }}`, spec.LaunchExpression)
	require.Equal(t, `{{ or (has "direccion" .user.groups) (eq .action.tool "noop") }}`, spec.ApproveExpression)
}

// TestAgentSecurityWithoutAccessDeniesAll verifies the shared policy default: a security block without access
// emits the unconditional deny rule (only admins keep access, via the engine's built-in rule).
func TestAgentSecurityWithoutAccessDeniesAll(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		`agents/a.yaml`: `
type: agent
instructions: Investiga.
security:
  approve: '{{ has "direccion" .user.groups }}'
`,
	})
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	require.Empty(t, p.Errors)

	spec := p.Resources[ResourceName{Kind: ResourceKindAgent, Name: "a"}].AgentSpec
	require.Len(t, spec.SecurityRules, 1)
	access := spec.SecurityRules[0].GetAccess()
	require.NotNil(t, access)
	require.False(t, access.Allow)
	require.Empty(t, access.ConditionExpression)
}

// TestAgentSecurityErrors covers the rejections: unknown keys, the reserved execute key, malformed gate
// expressions, and the approval-time context escaping into keys where it does not exist.
func TestAgentSecurityErrors(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// Unknown key inside security (knownFields is enforced recursively).
		`agents/unknown_key.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  aprove: 'true'
`,
		// The reserved execute key fails with an explanatory error, not as an unknown field.
		`agents/execute_reserved.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  execute: 'true'
`,
		// launch must be a valid template.
		`agents/launch_malformed.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  launch: '{{ has "x" .user.groups'
`,
		// launch must evaluate as a boolean expression.
		`agents/launch_not_bool.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  launch: 'not a boolean at all'
`,
		// The approval-time context is not available in launch.
		`agents/launch_action.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  launch: '{{ eq .action.tool "x" }}'
`,
		// Nor in access.
		`agents/access_action.yaml`: `
type: agent
instructions: X.
security:
  access: '{{ eq .action.tool "x" }}'
`,
		// approve may only see the action's identity, never its (model-chosen) arguments.
		`agents/approve_args.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  approve: '{{ eq .action.args.amount "1" }}'
`,
		// The field-level policy members do not apply to agents.
		`agents/include.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  include:
    - if: 'true'
      names: [x]
`,
		`agents/row_filter.yaml`: `
type: agent
instructions: X.
security:
  access: 'true'
  row_filter: "x = 1"
`,
	})

	wantErrors := []*runtimev1.ParseError{
		{Message: `field aprove not found`, FilePath: "/agents/unknown_key.yaml"},
		{Message: `'execute' is reserved`, FilePath: "/agents/execute_reserved.yaml"},
		{Message: `"launch" templating is not valid`, FilePath: "/agents/launch_malformed.yaml"},
		{Message: `"launch" expression error`, FilePath: "/agents/launch_not_bool.yaml"},
		{Message: `"launch" cannot reference ".action.tool"`, FilePath: "/agents/launch_action.yaml"},
		{Message: `"access" cannot reference ".action.tool"`, FilePath: "/agents/access_action.yaml"},
		{Message: `'approve' cannot reference ".action.args.amount"`, FilePath: "/agents/approve_args.yaml"},
		{Message: `'include'/'exclude' are not supported on agents`, FilePath: "/agents/include.yaml"},
		{Message: `'row_filter' is not supported on agents`, FilePath: "/agents/row_filter.yaml"},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, nil, wantErrors)
}
