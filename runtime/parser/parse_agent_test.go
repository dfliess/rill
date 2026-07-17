package parser

import (
	"context"
	"reflect"
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/stretchr/testify/require"
)

func TestAgent(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// A connector resource that the agent below references.
		`connectors/my_ai.yaml`: `
type: connector
driver: duckdb
`,
		// agent a1: full spec, references a declared connector.
		`agents/a1.yaml`: `
type: agent
display_name: Revenue Incident Agent
description: Investigates revenue anomalies and proposes a handoff.
model:
  connector: my_ai
  name: gpt-4o
instructions: Investigate using governed metrics only.
tools:
  - rill.query_metrics_view
  - jira.create_issue
limits:
  max_steps: 10
  timeout: 15m
`,
		// agent a2: minimal spec, references an undeclared connector (ref is dropped).
		`agents/a2.yaml`: `
type: agent
model:
  connector: openai
instructions: Only answer questions about governed metrics.
`,
	})

	resources := []*Resource{
		{
			Name:          ResourceName{Kind: ResourceKindConnector, Name: "my_ai"},
			Paths:         []string{"/connectors/my_ai.yaml"},
			ConnectorSpec: &runtimev1.ConnectorSpec{Driver: "duckdb"},
		},
		{
			Name:  ResourceName{Kind: ResourceKindAgent, Name: "a1"},
			Paths: []string{"/agents/a1.yaml"},
			Refs:  []ResourceName{{Kind: ResourceKindConnector, Name: "my_ai"}},
			AgentSpec: &runtimev1.AgentSpec{
				DisplayName:    "Revenue Incident Agent",
				Description:    "Investigates revenue anomalies and proposes a handoff.",
				ModelConnector: "my_ai",
				ModelName:      "gpt-4o",
				Instructions:   "Investigate using governed metrics only.",
				Tools:          []string{"rill.query_metrics_view", "jira.create_issue"},
				Limits: &runtimev1.AgentLimits{
					MaxSteps:       10,
					TimeoutSeconds: 900,
				},
			},
		},
		{
			Name:  ResourceName{Kind: ResourceKindAgent, Name: "a2"},
			Paths: []string{"/agents/a2.yaml"},
			AgentSpec: &runtimev1.AgentSpec{
				ModelConnector: "openai",
				Instructions:   "Only answer questions about governed metrics.",
			},
		},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, resources, nil)
}

func TestAgentMCP(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// An agent declaring an outbound MCP connector: the parser maps it structurally onto AgentSpec.Mcp; deep
		// validation (transport, https, SSRF posture) is the reconciler's job.
		`agents/support.yaml`: `
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
	})

	resources := []*Resource{
		{
			Name:  ResourceName{Kind: ResourceKindAgent, Name: "support"},
			Paths: []string{"/agents/support.yaml"},
			AgentSpec: &runtimev1.AgentSpec{
				Instructions: "Triage support tickets.",
				Mcp: []*runtimev1.MCPConnector{
					{
						Name:              "jira",
						Url:               "https://jira.example.com/mcp",
						AuthSecret:        "JIRA_TOKEN",
						AllowedHosts:      []string{"jira.example.com"},
						TrustReadOnlyHint: true,
					},
				},
			},
		},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, resources, nil)
}

func TestAgentMCPErrors(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// A connector with no name.
		`agents/no_name.yaml`: `
type: agent
instructions: Do something.
mcp:
  - url: https://a.example.com/mcp
`,
		// A connector with no url.
		`agents/no_url.yaml`: `
type: agent
instructions: Do something.
mcp:
  - name: jira
`,
		// A connector whose name contains a dot (breaks mcp.<connector>.<tool> namespacing).
		`agents/dotname.yaml`: `
type: agent
instructions: Do something.
mcp:
  - name: jira.prod
    url: https://jira.example.com/mcp
`,
		// Two connectors sharing a name (tool namespacing would be ambiguous).
		`agents/dup.yaml`: `
type: agent
instructions: Do something.
mcp:
  - name: jira
    url: https://a.example.com/mcp
  - name: jira
    url: https://b.example.com/mcp
`,
	})

	wantErrors := []*runtimev1.ParseError{
		{
			Message:  `each "mcp" connector must set "name"`,
			FilePath: "/agents/no_name.yaml",
		},
		{
			Message:  `mcp connector name "jira.prod" must not contain dots (used as namespace separator in tool names)`,
			FilePath: "/agents/dotname.yaml",
		},
		{
			Message:  `mcp connector "jira" must set "url"`,
			FilePath: "/agents/no_url.yaml",
		},
		{
			Message:  `duplicate mcp connector name "jira"`,
			FilePath: "/agents/dup.yaml",
		},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, nil, wantErrors)
}

// TestAgentInlineTriggers checks the ADR-0016 desugaring: a type: agent with a triggers: list produces the Agent
// resource plus one synthetic AgentTrigger per entry, each named "<agent>__trigger_<i>", carrying the agent's paths
// and tags and a ref back to the agent. The two entries exercise a named source and a wildcard (no source.name).
func TestAgentInlineTriggers(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		`agents/revenue_incident.yaml`: `
type: agent
instructions: Investigate anomalies and prepare evidence.
triggers:
  - source:
      kind: alert
      name: revenue_drop
      events: [entered_fail, renotify_due]
    input:
      prompt: Investiga la alerta y prepara evidencia.
      context:
        alert: revenue_drop
    deduplication:
      window: 24h
  - source:
      kind: alert
      events: [entered_fail]
`,
	})

	resources := []*Resource{
		{
			Name:  ResourceName{Kind: ResourceKindAgent, Name: "revenue_incident"},
			Paths: []string{"/agents/revenue_incident.yaml"},
			AgentSpec: &runtimev1.AgentSpec{
				Instructions: "Investigate anomalies and prepare evidence.",
			},
		},
		{
			Name:  ResourceName{Kind: ResourceKindAgentTrigger, Name: "revenue_incident__trigger_0"},
			Paths: []string{"/agents/revenue_incident.yaml"},
			Refs:  []ResourceName{{Kind: ResourceKindAgent, Name: "revenue_incident"}},
			AgentTriggerSpec: &runtimev1.AgentTriggerSpec{
				Agent: "revenue_incident",
				Source: &runtimev1.AgentTriggerSource{
					Kind:   runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT,
					Name:   "revenue_drop",
					Events: []string{"entered_fail", "renotify_due"},
				},
				Actor: &runtimev1.AgentTriggerActor{},
				Input: &runtimev1.AgentTriggerInput{
					Prompt:  "Investiga la alerta y prepara evidencia.",
					Context: map[string]string{"alert": "revenue_drop"},
				},
				Deduplication: &runtimev1.AgentTriggerDeduplication{WindowSeconds: 86400},
			},
		},
		{
			// Wildcard trigger: no source.name, so it fires for any alert of the subscribed events.
			Name:  ResourceName{Kind: ResourceKindAgentTrigger, Name: "revenue_incident__trigger_1"},
			Paths: []string{"/agents/revenue_incident.yaml"},
			Refs:  []ResourceName{{Kind: ResourceKindAgent, Name: "revenue_incident"}},
			AgentTriggerSpec: &runtimev1.AgentTriggerSpec{
				Agent: "revenue_incident",
				Source: &runtimev1.AgentTriggerSource{
					Kind:   runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT,
					Events: []string{"entered_fail"},
				},
				Actor: &runtimev1.AgentTriggerActor{},
				Input: &runtimev1.AgentTriggerInput{},
			},
		},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, resources, nil)
}

// TestAgentInlineTriggersLeaveSpecHashUnchanged is the load-bearing ADR-0016 invariant: adding a triggers: list must
// not change the agent's AgentSpec. The reconciler hashes only AgentSpec, so an identical AgentSpec means an identical
// spec_hash: editing when the agent fires never looks like editing what it does.
func TestAgentInlineTriggersLeaveSpecHashUnchanged(t *testing.T) {
	ctx := context.Background()

	parseAgentSpec := func(t *testing.T, yaml string) *runtimev1.AgentSpec {
		t.Helper()
		repo := makeRepo(t, map[string]string{"rill.yaml": ``, "agents/a.yaml": yaml})
		p, err := Parse(ctx, repo, "", "", "duckdb", true)
		require.NoError(t, err)
		require.Empty(t, p.Errors)
		res := p.Resources[ResourceName{Kind: ResourceKindAgent, Name: "a"}]
		require.NotNil(t, res)
		return res.AgentSpec
	}

	withoutTriggers := parseAgentSpec(t, `
type: agent
instructions: Investigate anomalies and prepare evidence.
`)
	withTriggers := parseAgentSpec(t, `
type: agent
instructions: Investigate anomalies and prepare evidence.
triggers:
  - source:
      kind: alert
      name: revenue_drop
      events: [entered_fail]
`)
	require.True(t, reflect.DeepEqual(withoutTriggers, withTriggers),
		"the triggers: list must not leak into AgentSpec, or the agent's spec_hash would change")
}

// TestAgentInlineTriggerNameCollision checks that a synthetic trigger name colliding with a hand-written
// agent_trigger is reported as a parse error rather than silently overwriting a resource.
func TestAgentInlineTriggerNameCollision(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		`agents/a.yaml`: `
type: agent
instructions: Investigate.
triggers:
  - source:
      kind: alert
      name: revenue_drop
      events: [entered_fail]
`,
		// A standalone trigger whose name is exactly the synthetic name the agent's first inline trigger desugars to.
		`triggers/a__trigger_0.yaml`: `
type: agent_trigger
agent: a
source:
  kind: alert
  name: revenue_drop
  events: [entered_fail]
`,
	})
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	// One of the two files carries the collision error; the exact file depends on parse order, so assert on the message.
	var found bool
	for _, e := range p.Errors {
		if e.FilePath == "/agents/a.yaml" || e.FilePath == "/triggers/a__trigger_0.yaml" {
			require.Contains(t, e.Message, "name collision")
			found = true
		}
	}
	require.True(t, found, "expected a name collision error, got: %v", p.Errors)
}

func TestAgentSubSecondTimeoutRoundsUp(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		`agents/quick.yaml`: `
type: agent
instructions: Answer quickly.
limits:
  timeout: 500ms
`,
	})
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)

	res := p.Resources[ResourceName{Kind: ResourceKindAgent, Name: "quick"}]
	require.NotNil(t, res)
	require.NotNil(t, res.AgentSpec.Limits)
	// A sub-second timeout rounds up to 1s so the limit is not truncated to 0 and dropped.
	require.Equal(t, uint32(1), res.AgentSpec.Limits.TimeoutSeconds)
}

func TestAgentErrors(t *testing.T) {
	ctx := context.Background()
	repo := makeRepo(t, map[string]string{
		`rill.yaml`: ``,
		// Missing the required instructions field.
		`agents/no_instructions.yaml`: `
type: agent
display_name: No Instructions
`,
		// Invalid timeout value.
		`agents/bad_timeout.yaml`: `
type: agent
instructions: Do something.
limits:
  timeout: not-a-duration
`,
		// Unknown field (knownFields is enforced for agents).
		`agents/unknown_field.yaml`: `
type: agent
instructions: Do something.
not_a_field: true
`,
	})

	wantErrors := []*runtimev1.ParseError{
		{
			Message:  `agents must set "instructions"`,
			FilePath: "/agents/no_instructions.yaml",
		},
		{
			Message:  `invalid value "not-a-duration" for property "limits.timeout"`,
			FilePath: "/agents/bad_timeout.yaml",
		},
		{
			Message:  "not_a_field",
			FilePath: "/agents/unknown_field.yaml",
		},
	}
	p, err := Parse(ctx, repo, "", "", "duckdb", true)
	require.NoError(t, err)
	requireResourcesAndErrors(t, p, nil, wantErrors)
}
