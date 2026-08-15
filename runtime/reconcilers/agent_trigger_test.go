package reconcilers_test

import (
	"testing"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// minimalAlertYAML is an alert that parses and exists as a resource. It references a missing metrics view, so it
// reconciles with its own error, but the trigger reconciler only needs the alert to EXIST, not to be valid.
const minimalAlertYAML = `
type: alert
display_name: Revenue Drop
refs:
- type: MetricsView
  name: mv_missing
watermark: inherit
intervals:
  duration: P1D
query:
  name: MetricsViewAggregation
  args:
    metrics_view: mv_missing
    measures:
    - name: m
email:
  recipients:
  - somebody@example.com
`

func TestAgentTriggerValid(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate using governed metrics only.
`,
			"alerts/revenue_drop.yaml": minimalAlertYAML,
			"triggers/on_fail.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  name: revenue_drop
  events: [entered_fail, recovered]
input:
  prompt: Investiga la caida.
deduplication:
  window: 24h
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	// The trigger itself validates: its agent and its source alert both exist. (The alert has its own reconcile
	// error from the missing metrics view; that does not concern the trigger.)
	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "on_fail")
	trigger := res.GetAgentTrigger()
	require.NotNil(t, trigger)
	require.Empty(t, res.Meta.ReconcileError)
	require.NotNil(t, trigger.State.ValidSpec)
	require.NotEmpty(t, trigger.State.SpecHash)
	require.Equal(t, "incident", trigger.State.ValidSpec.Agent)
	require.Equal(t, "revenue_drop", trigger.State.ValidSpec.Source.Name)
	require.Equal(t, uint32(86400), trigger.State.ValidSpec.Deduplication.WindowSeconds)
}

func TestAgentTriggerWildcardSource(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			// No source.name: a wildcard that fires for any alert of the subscribed events (ADR-0016). There is no
			// source resource to resolve, so the trigger validates with only its agent present.
			"triggers/any_alert.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  events: [entered_fail]
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "any_alert")
	trigger := res.GetAgentTrigger()
	require.NotNil(t, trigger)
	require.Empty(t, res.Meta.ReconcileError)
	require.NotNil(t, trigger.State.ValidSpec, "a wildcard trigger with an existing agent must validate")
	require.Equal(t, "", trigger.State.ValidSpec.Source.Name)
}

// TestAgentInlineTriggersReconcile is the end-to-end ADR-0016 path: an inline triggers: list desugars to standalone
// AgentTrigger resources that reconcile like hand-written ones, without touching the agent's spec hash.
func TestAgentInlineTriggersReconcile(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml":                "",
			"alerts/revenue_drop.yaml": minimalAlertYAML,
			"agents/incident.yaml": `
type: agent
instructions: Investigate using governed metrics only.
triggers:
  - source:
      kind: alert
      name: revenue_drop
      events: [entered_fail]
    input:
      prompt: Investiga la caida.
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	// The synthetic trigger reconciles valid and points back at the agent.
	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "incident__trigger_0")
	trigger := res.GetAgentTrigger()
	require.NotNil(t, trigger)
	require.Empty(t, res.Meta.ReconcileError)
	require.NotNil(t, trigger.State.ValidSpec)
	require.Equal(t, "incident", trigger.State.ValidSpec.Agent)
	require.Equal(t, "revenue_drop", trigger.State.ValidSpec.Source.Name)
	require.Equal(t, "Investiga la caida.", trigger.State.ValidSpec.Input.Prompt)

	// The agent is a separate, valid resource; the triggers are not part of it.
	agent := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "incident")
	require.NotNil(t, agent.GetAgent().State.ValidSpec)
}

func TestAgentTriggerUnknownAgent(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"triggers/orphan.yaml": `
type: agent_trigger
agent: does_not_exist
source:
  kind: alert
  name: some_alert
  events: [entered_fail]
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgentTrigger, "orphan", "unknown agent")

	// Fail-closed: no valid spec, so the dispatcher never sees this trigger.
	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "orphan")
	require.Nil(t, res.GetAgentTrigger().State.ValidSpec)
}

func TestAgentTriggerUnknownSource(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			"triggers/on_ghost.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  name: ghost_alert
  events: [entered_fail]
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// Exactly one reconcile error (the trigger), no parse errors: the agent is valid, only the source is missing.
	testruntime.RequireReconcileState(t, rt, id, -1, 1, 0)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgentTrigger, "on_ghost", "unknown alert")

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "on_ghost")
	require.Nil(t, res.GetAgentTrigger().State.ValidSpec)
}

func TestAgentTriggerInvalidEvent(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			"alerts/revenue_drop.yaml": minimalAlertYAML,
			"triggers/bad_event.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  name: revenue_drop
  events: [completed]
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// "completed" is a report event, not an alert event: the trigger fails closed.
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgentTrigger, "bad_event", "not valid for a")

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "bad_event")
	require.Nil(t, res.GetAgentTrigger().State.ValidSpec)
}

func TestAgentTriggerScheduleValid(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			// A schedule trigger owns its own clock: no source resource, no events, just a cron.
			"triggers/nightly.yaml": `
type: agent_trigger
agent: incident
source:
  kind: schedule
  cron: "0 3 * * *"
input:
  prompt: Corre el resumen diario.
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "nightly")
	trigger := res.GetAgentTrigger()
	require.NotNil(t, trigger)
	require.Empty(t, res.Meta.ReconcileError)
	require.NotNil(t, trigger.State.ValidSpec, "a schedule trigger with a valid cron must validate")
	require.Equal(t, "0 3 * * *", trigger.State.ValidSpec.Source.Cron)
	require.Empty(t, trigger.State.ValidSpec.Source.Name)
}

func TestAgentTriggerScheduleRequiresCron(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			"triggers/no_cron.yaml": `
type: agent_trigger
agent: incident
source:
  kind: schedule
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgentTrigger, "no_cron", `must set "source.cron"`)
	require.Nil(t, testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "no_cron").GetAgentTrigger().State.ValidSpec)
}

func TestAgentTriggerScheduleInvalidCron(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			"triggers/bad_cron.yaml": `
type: agent_trigger
agent: incident
source:
  kind: schedule
  cron: "not a cron"
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgentTrigger, "bad_cron", "invalid")
}

func TestAgentTriggerActorAttributes(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			// The actor's run-as attributes become the started run's SecurityClaims (like an alert's query_for).
			"triggers/on_fail.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  events: [entered_fail]
actor:
  attributes:
    email: act@kairosagentica.com
    admin: false
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "on_fail")
	trigger := res.GetAgentTrigger()
	require.NotNil(t, trigger)
	require.Empty(t, res.Meta.ReconcileError)
	require.NotNil(t, trigger.State.ValidSpec.Actor.Attributes, "run-as attributes must reach the valid spec")
	attrs := trigger.State.ValidSpec.Actor.Attributes.AsMap()
	require.Equal(t, "act@kairosagentica.com", attrs["email"])
	require.Equal(t, false, attrs["admin"])
}

func TestAgentTriggerActorUserMutuallyExclusive(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			"triggers/on_fail.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  events: [entered_fail]
actor:
  user_id: usr_1
  attributes:
    email: explicit@kairosagentica.com
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// Mixing a named user with explicit attributes is rejected at parse (the resource is never created): it
	// would audit the run as one user while it executes with a different identity's attributes.
	testruntime.RequireParseErrors(t, rt, id, map[string]string{
		"/triggers/on_fail.yaml": "at most one",
	})
}

func TestAgentTriggerActorUserIDFallsBackWithoutAdmin(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			// user_id names the run-as user, resolved via admin. With no admin service (the test's noop admin), it
			// falls back to no attributes — the trigger still validates, and the run fails closed until an identity
			// resolves. It must NOT fail the reconcile.
			"triggers/on_fail.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  events: [entered_fail]
actor:
  user_id: usr_act
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	res := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "on_fail")
	require.Empty(t, res.Meta.ReconcileError)
	require.NotNil(t, res.GetAgentTrigger().State.ValidSpec)
	require.Equal(t, "usr_act", res.GetAgentTrigger().State.ValidSpec.Actor.UserId)
	require.Nil(t, res.GetAgentTrigger().State.ValidSpec.Actor.Attributes, "no admin resolution => nil attributes (fail closed)")
}

func TestAgentTriggerActorRejectsUserEmail(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			// Naming the actor by email is rejected: the address can change or leave the org, at which point the
			// trigger silently stops resolving an identity and every run fails closed with nobody noticing.
			"triggers/on_fail.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  events: [entered_fail]
actor:
  user_email: other@kairosagentica.com
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	// A rejected identity is a parse error (the resource never reaches the reconciler), like an invalid alert.
	testruntime.RequireParseErrors(t, rt, id, map[string]string{
		"/triggers/on_fail.yaml": `"actor.user_email" is not supported`,
	})
}

func TestAgentTriggerAlertCronForbidden(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/incident.yaml": `
type: agent
instructions: Investigate.
`,
			// Cron belongs to a schedule source; setting it on an alert source is a misconfiguration that fails closed.
			"triggers/alert_cron.yaml": `
type: agent_trigger
agent: incident
source:
  kind: alert
  events: [entered_fail]
  cron: "0 3 * * *"
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileErrorContains(t, rt, id, runtime.ResourceKindAgentTrigger, "alert_cron", `only a schedule`)
}
