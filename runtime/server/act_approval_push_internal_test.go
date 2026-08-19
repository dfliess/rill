package server

import (
	"errors"
	"testing"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// These tests cover the push that announces a paused run (kairos-cloud#143): who receives it, that one pause is one
// notification however many actions it holds, and that nothing about it can affect the run.

// actPushProjectFiles is a project with two agents: one whose approve policy names a group, one that discriminates
// by tool (so a batch mixing tools has to reach the union of both audiences).
func actPushProjectFiles() map[string]string {
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
	}
}

// actPushMembers is the project's membership as the admin service reports it: an operator, the group the policies
// name, and the member that may only sign refunds.
func actPushMembers() []drivers.ProjectMember {
	member := func(id, email string, groups []any, admin, editTrigger bool) drivers.ProjectMember {
		return drivers.ProjectMember{
			UserID: id,
			Email:  email,
			Attributes: map[string]any{
				"name": id, "email": email, "domain": "example.com", "groups": groups, "admin": admin,
			},
			EditTrigger: editTrigger,
		}
	}
	return []drivers.ProjectMember{
		member("usr_admin", "admin@example.com", []any{}, true, true),
		member("usr_ana", "ana@example.com", []any{"direccion"}, false, false),
		member("usr_anafin", "anafin@example.com", []any{"finanzas"}, false, false),
		member("usr_marta", "marta@example.com", []any{"operaciones"}, false, false),
	}
}

// newActPushNotifier builds the notifier over an instance that runs in Rill Cloud (it has the org and project
// annotations the admin service sets on deployment) and reports the given members.
func newActPushNotifier(t *testing.T, members []drivers.ProjectMember) (*actApprovalPushNotifier, string) {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files:       actPushProjectFiles(),
		Annotations: map[string]string{"organization_name": "acme", "project_name": "demo"},
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	testruntime.SetProjectMembers(instanceID, members)
	return &actApprovalPushNotifier{runtime: rt, logger: zap.NewNop()}, instanceID
}

func actPendingApprovals(instanceID string, agentName string, actions ...act.PendingAction) act.PendingApprovals {
	return act.PendingApprovals{
		InstanceID:     instanceID,
		RunID:          act.ComposeRunID(instanceID, agentName, "alerta/2024"),
		AgentName:      agentName,
		IdempotencyKey: "alerta/2024",
		RunActor:       "usr_marta",
		Actions:        actions,
	}
}

// TestActApprovalPushNotifiesPolicyApprovers is the core case: the push goes to the members the agent's approve
// policy authorizes, and it opens the run's own page.
func TestActApprovalPushNotifiesPolicyApprovers(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	pending := actPendingApprovals(instanceID, "cobranza", act.PendingAction{Tool: "create_issue", Connector: "jira_ops"})
	notifier.NotifyPendingApprovals(t.Context(), pending)

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, "act_approvals", pushes[0].Category)
	// Only the group the policy names: not the operator (who could still break the glass), not the member whose
	// only claim to the agent is access, and not the actor's group.
	require.Equal(t, []string{"ana@example.com"}, pushes[0].Recipients)
	require.Equal(t, "Aprobación pendiente: Cobranza", pushes[0].Title)
	require.Equal(t, `La acción "create_issue" espera tu aprobación.`, pushes[0].Body)
	// The run segment is the idempotency key, escaped: the key of a trigger-fired run may contain slashes.
	require.Equal(t, "/acme/demo/-/agents/cobranza/runs/alerta%2F2024", pushes[0].LinkPath)
	require.Equal(t, "act:acme/demo/"+pending.RunID, pushes[0].Tag)
}

// TestActApprovalPushBatchIsOneNotification pins the decision of the grill: a batch is one moment of signing, so N
// actions produce ONE push, addressed to everyone who may decide any of them.
func TestActApprovalPushBatchIsOneNotification(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "finanzas",
		act.PendingAction{Tool: "close_account", Connector: "core"},
		act.PendingAction{Tool: "issue_refund", Connector: "core"},
		act.PendingAction{Tool: "issue_refund", Connector: "core"},
	))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1, "one pause of the run is one notification")
	// The union: Ana signs everything, Anafin only the refunds, and she is told about the pause that holds one.
	require.Equal(t, []string{"ana@example.com", "anafin@example.com"}, pushes[0].Recipients)
	require.Equal(t, `3 acciones esperan tu aprobación, empezando por "close_account".`, pushes[0].Body)
}

// TestActApprovalPushOnePerPause checks that a run that pauses twice notifies twice, under one stable tag so the
// browser shows the newer pause instead of stacking them.
func TestActApprovalPushOnePerPause(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	pending := actPendingApprovals(instanceID, "cobranza", act.PendingAction{Tool: "create_issue", Connector: "jira_ops"})
	notifier.NotifyPendingApprovals(t.Context(), pending)
	pending.Actions = []act.PendingAction{{Tool: "close_issue", Connector: "jira_ops"}}
	notifier.NotifyPendingApprovals(t.Context(), pending)

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 2)
	require.Equal(t, `La acción "close_issue" espera tu aprobación.`, pushes[1].Body)
	require.Equal(t, pushes[0].Tag, pushes[1].Tag, "both pauses of a run share a tag")
}

// TestActApprovalPushDeletedAgent covers the agent that is gone (or lost its valid spec) since the run started: its
// policy went with it, so the operators are notified, which is who the API would let decide it.
func TestActApprovalPushDeletedAgent(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "borrado",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, []string{"admin@example.com"}, pushes[0].Recipients)
	require.Equal(t, "Aprobación pendiente: borrado", pushes[0].Title)
}

// TestActApprovalPushWithoutRecipients checks that a batch nobody can decide (and no operator to fall back to)
// sends nothing, instead of a push with an empty recipient list.
func TestActApprovalPushWithoutRecipients(t *testing.T) {
	members := actPushMembers()
	for i := range members {
		members[i].EditTrigger = false
		members[i].Attributes["groups"] = []any{"ventas"}
		members[i].Attributes["admin"] = false
	}
	notifier, instanceID := newActPushNotifier(t, members)

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	require.Empty(t, testruntime.PushNotifications(instanceID))
}

// TestActApprovalPushOutsideRillCloud checks that an instance without the org and project annotations (Rill
// Developer) notifies nobody: there is no frontend to link to and no subscription to reach.
func TestActApprovalPushOutsideRillCloud(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: actPushProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	testruntime.SetProjectMembers(instanceID, actPushMembers())
	notifier := &actApprovalPushNotifier{runtime: rt, logger: zap.NewNop()}

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	require.Empty(t, testruntime.PushNotifications(instanceID))
}

// TestActApprovalPushWithoutMemberAttributes covers an admin service that cannot enumerate the project's members
// (an older admin, or none at all): no membership, no notification, and no failure.
func TestActApprovalPushWithoutMemberAttributes(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files:       actPushProjectFiles(),
		Annotations: map[string]string{"organization_name": "acme", "project_name": "demo"},
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	notifier := &actApprovalPushNotifier{runtime: rt, logger: zap.NewNop()}

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	require.Empty(t, testruntime.PushNotifications(instanceID))
}

// TestActApprovalPushFailureIsSwallowed pins the contract the run depends on: a failing push is logged and dropped,
// never returned, so the approvals are created and the run waits for them either way.
func TestActApprovalPushFailureIsSwallowed(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())
	testruntime.FailPushNotifications(instanceID, errors.New("push service down"))

	require.NotPanics(t, func() {
		notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
			act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))
	})
	require.Empty(t, testruntime.PushNotifications(instanceID))
}

// TestActApprovalPushWithoutActions checks the degenerate payload: a pause with no manual action notifies nobody
// (the executor already filters those out, so this only pins that the notifier does not invent a notification).
func TestActApprovalPushWithoutActions(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza"))

	require.Empty(t, testruntime.PushNotifications(instanceID))
}

// TestAgentRunPathMirrorsFrontend pins the deep link against the frontend's runDetailPath: the run segment is the
// idempotency key (not the composite run ID), and every segment is escaped.
func TestAgentRunPathMirrorsFrontend(t *testing.T) {
	require.Equal(t, "/acme/demo/-/agents/cobranza/runs/run-1", agentRunPath("acme", "demo", "cobranza", "run-1"))
	require.Equal(t, "/acme/demo/-/agents/mi%20agente/runs/alert%2Fa1%2F2024", agentRunPath("acme", "demo", "mi agente", "alert/a1/2024"))
}

// TestActApprovalPushResolvesAgentByName is a guard on the resource lookup: the notifier must read the agent the
// run names, not any agent, or a project's policies would leak into each other's notifications.
func TestActApprovalPushResolvesAgentByName(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	res := notifier.agentResource(t.Context(), instanceID, "finanzas")
	require.NotNil(t, res)
	require.Equal(t, "finanzas", res.Meta.Name.Name)
	require.Equal(t, runtime.ResourceKindAgent, res.Meta.Name.Kind)
	require.Nil(t, notifier.agentResource(t.Context(), instanceID, "no_existe"))
}
