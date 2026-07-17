package parser

import (
	"context"
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/stretchr/testify/require"
)

func TestAgentTrigger(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// The agent the trigger starts. It must exist for the trigger's agent ref to resolve.
		`agents/revenue_incident.yaml`: `
type: agent
instructions: Investigate using governed metrics only.
`,
		// A full trigger binding the agent to an alert's transition events.
		`triggers/on_revenue_drop.yaml`: `
type: agent_trigger
agent: revenue_incident
source:
  kind: alert
  name: revenue_drop
  events: [entered_fail, renotify_due]
actor:
  user_email: act@revenue.ops
input:
  prompt: Investiga la alerta y prepara evidencia.
  context:
    alert: revenue_drop
deduplication:
  window: 24h
`,
	})

	resources := []*Resource{
		{
			Name:  ResourceName{Kind: ResourceKindAgent, Name: "revenue_incident"},
			Paths: []string{"/agents/revenue_incident.yaml"},
			AgentSpec: &runtimev1.AgentSpec{
				Instructions: "Investigate using governed metrics only.",
			},
		},
		{
			Name:  ResourceName{Kind: ResourceKindAgentTrigger, Name: "on_revenue_drop"},
			Paths: []string{"/triggers/on_revenue_drop.yaml"},
			Refs:  []ResourceName{{Kind: ResourceKindAgent, Name: "revenue_incident"}},
			AgentTriggerSpec: &runtimev1.AgentTriggerSpec{
				Agent: "revenue_incident",
				Source: &runtimev1.AgentTriggerSource{
					Kind:   runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT,
					Name:   "revenue_drop",
					Events: []string{"entered_fail", "renotify_due"},
				},
				Actor: &runtimev1.AgentTriggerActor{
					UserEmail: "act@revenue.ops",
				},
				Input: &runtimev1.AgentTriggerInput{
					Prompt:  "Investiga la alerta y prepara evidencia.",
					Context: map[string]string{"alert": "revenue_drop"},
				},
				Deduplication: &runtimev1.AgentTriggerDeduplication{WindowSeconds: 86400},
			},
		},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, resources, nil)
}

func TestAgentTriggerSubSecondWindowRoundsUp(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		`agents/a.yaml`: `
type: agent
instructions: Investigate.
`,
		`triggers/t.yaml`: `
type: agent_trigger
agent: a
source:
  kind: alert
  name: x
  events: [entered_fail]
deduplication:
  window: 500ms
`,
	})
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)

	res := p.Resources[ResourceName{Kind: ResourceKindAgentTrigger, Name: "t"}]
	require.NotNil(t, res)
	require.NotNil(t, res.AgentTriggerSpec.Deduplication)
	// A sub-second window rounds up to 1s so it is not truncated to 0 and dropped.
	require.Equal(t, uint32(1), res.AgentTriggerSpec.Deduplication.WindowSeconds)
}

func TestAgentTriggerErrors(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// Missing the required agent field.
		`triggers/no_agent.yaml`: `
type: agent_trigger
source:
  kind: alert
  name: x
  events: [entered_fail]
`,
		// Missing source.kind.
		`triggers/no_kind.yaml`: `
type: agent_trigger
agent: a
source:
  name: x
`,
		// Invalid source.kind.
		`triggers/bad_kind.yaml`: `
type: agent_trigger
agent: a
source:
  kind: sunset
  name: x
`,
		// Invalid deduplication window.
		`triggers/bad_window.yaml`: `
type: agent_trigger
agent: a
source:
  kind: alert
  name: x
  events: [entered_fail]
deduplication:
  window: soon
`,
		// Unknown field (knownFields is enforced).
		`triggers/unknown_field.yaml`: `
type: agent_trigger
agent: a
source:
  kind: alert
  name: x
  events: [entered_fail]
not_a_field: true
`,
		// A named user and explicit attributes at once: identity must be declared at most one way.
		`triggers/mixed_identity.yaml`: `
type: agent_trigger
agent: a
source:
  kind: alert
  name: x
  events: [entered_fail]
actor:
  user_email: alice@example.com
  attributes:
    email: bob@example.com
`,
		// A negative deduplication window would wrap to a huge unsigned value.
		`triggers/negative_window.yaml`: `
type: agent_trigger
agent: a
source:
  kind: alert
  name: x
  events: [entered_fail]
deduplication:
  window: -1s
`,
	})

	wantErrors := []*runtimev1.ParseError{
		{Message: `agent_trigger must set "agent"`, FilePath: "/triggers/no_agent.yaml"},
		{Message: `agent_trigger must set "source.kind"`, FilePath: "/triggers/no_kind.yaml"},
		{Message: `invalid value "sunset" for property "source.kind"`, FilePath: "/triggers/bad_kind.yaml"},
		{Message: `invalid value "soon" for property "deduplication.window"`, FilePath: "/triggers/bad_window.yaml"},
		{Message: "not_a_field", FilePath: "/triggers/unknown_field.yaml"},
		{Message: `at most one of "user_id", "user_email" or "attributes"`, FilePath: "/triggers/mixed_identity.yaml"},
		{Message: `"deduplication.window" must not be negative`, FilePath: "/triggers/negative_window.yaml"},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, nil, wantErrors)
}
