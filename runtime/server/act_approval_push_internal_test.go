package server

import (
	"context"
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
// notification however many actions it holds, and that nothing about it can affect the run. They also cover the
// number it carries for each recipient's app badge (kairos-cloud#147).

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

// actPushRunStore is the approval store as the badge count reads it. Only ListApprovals is answered, because
// counting is the only thing the notifier asks of the store; the embedded nil interface turns any other call into
// a panic, which is the signal wanted if that ever stops being true.
type actPushRunStore struct {
	act.RunStore
	approvals []*act.Approval
	err       error
	// filter records what the notifier asked for, so a test can pin that the count reads pending approvals of its
	// own instance and nothing else.
	filter act.ListApprovalsFilter
}

func (s *actPushRunStore) ListApprovals(_ context.Context, f act.ListApprovalsFilter) ([]*act.Approval, error) {
	s.filter = f
	if s.err != nil {
		return nil, s.err
	}
	return s.approvals, nil
}

// actStoredApproval is one pending approval as the store holds it, with the fields the approve policy discriminates
// on. runKey is the run's idempotency key: two approvals sharing it belong to the same pause of the same run.
func actStoredApproval(instanceID, agentName, runKey, tool, connector string) *act.Approval {
	return &act.Approval{
		ApprovalID:      runKey + ":" + tool,
		RunID:           act.ComposeRunID(instanceID, agentName, runKey),
		InstanceID:      instanceID,
		AgentName:       agentName,
		RunActorSubject: "usr_marta",
		ToolName:        tool,
		Connector:       connector,
		Status:          act.ApprovalStatusPending,
	}
}

// newActPushNotifier builds the notifier over an instance that runs in Rill Cloud (it has the org and project
// annotations the admin service sets on deployment) and reports the given members.
// Its store holds nothing pending; a test about the badge count replaces it with one that does.
func newActPushNotifier(t *testing.T, members []drivers.ProjectMember) (*actApprovalPushNotifier, string) {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files:       actPushProjectFiles(),
		Annotations: map[string]string{"organization_name": "acme", "project_name": "demo"},
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	testruntime.SetProjectMembers(instanceID, members)
	return &actApprovalPushNotifier{
		runtime: rt,
		runs:    &actPushRunStore{},
		logger:  zap.NewNop(),
	}, instanceID
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
	notifier := &actApprovalPushNotifier{runtime: rt, runs: &actPushRunStore{}, logger: zap.NewNop()}

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
	notifier := &actApprovalPushNotifier{runtime: rt, runs: &actPushRunStore{}, logger: zap.NewNop()}

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

	resolver := notifier.approverResolver(instanceID, actPushMembers())
	res := resolver.agent(t.Context(), "finanzas")
	require.NotNil(t, res)
	require.Equal(t, "finanzas", res.Meta.Name.Name)
	require.Equal(t, runtime.ResourceKindAgent, res.Meta.Name.Kind)
	require.Nil(t, resolver.agent(t.Context(), "no_existe"))

	// The second lookup of a name is answered from the resolver's memory, including the "it is gone" answer.
	require.Same(t, res, resolver.agent(t.Context(), "finanzas"))
	require.Nil(t, resolver.agent(t.Context(), "no_existe"))
}

// TestActApprovalPushBadgeCountsPerRecipient is the core of the app badge (kairos-cloud#147): the number is each
// recipient's own. Authority discriminates by action, so two people told about the same pause can genuinely have
// different amounts of work waiting, and a shared number would be wrong for at least one of them.
func TestActApprovalPushBadgeCountsPerRecipient(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())
	// Two runs of "finanzas" are waiting. Ana signs everything; Anafin signs only refunds, so only the second run
	// is hers.
	notifier.runs = &actPushRunStore{approvals: []*act.Approval{
		actStoredApproval(instanceID, "finanzas", "alerta/2024", "close_account", "core"),
		actStoredApproval(instanceID, "finanzas", "alerta/2025", "issue_refund", "core"),
	}}

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "finanzas",
		act.PendingAction{Tool: "close_account", Connector: "core"},
		act.PendingAction{Tool: "issue_refund", Connector: "core"},
	))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, map[string]int32{
		"ana@example.com":    2,
		"anafin@example.com": 1,
	}, pushes[0].PendingByRecipient)
}

// TestActApprovalPushBadgeCountsRunsNotApprovals pins the unit of the count against the unit of the inbox: a run
// that proposed three actions in one turn is one pause, one notification and one row to look at, so it is one on
// the icon. It mirrors countPendingRuns in web-admin/src/features/notifications/app-badge.ts, which is what the app
// recomputes for itself; if the two disagreed the badge would jump on every launch.
func TestActApprovalPushBadgeCountsRunsNotApprovals(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())
	// One pause of one run that proposed three actions: three approvals, one thing to look at.
	notifier.runs = &actPushRunStore{approvals: []*act.Approval{
		actStoredApproval(instanceID, "cobranza", "alerta/2024", "create_issue", "jira_ops"),
		actStoredApproval(instanceID, "cobranza", "alerta/2024", "close_issue", "jira_ops"),
		actStoredApproval(instanceID, "cobranza", "alerta/2024", "assign_issue", "jira_ops"),
	}}

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, map[string]int32{"ana@example.com": 1}, pushes[0].PendingByRecipient)
}

// TestActApprovalPushBadgeSendsZero covers the number that clears the badge. A recipient with nothing left waiting
// has to be told so explicitly: the app adds up one number per project, so a recipient simply left out of the map
// would keep whatever the icon was showing for this one.
func TestActApprovalPushBadgeSendsZero(t *testing.T) {
	// The store holds nothing pending, which is what a decision taken between the pause and the notification
	// leaves behind.
	notifier, instanceID := newActPushNotifier(t, actPushMembers())

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, map[string]int32{"ana@example.com": 0}, pushes[0].PendingByRecipient)
}

// TestActApprovalPushBadgeReadsOwnInstance pins what the count is read from: the pending approvals of this
// instance. A count that reached across instances would put another tenant's work on someone's icon.
func TestActApprovalPushBadgeReadsOwnInstance(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())
	store := &actPushRunStore{}
	notifier.runs = store

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	require.Equal(t, instanceID, store.filter.InstanceID)
	require.Equal(t, act.ApprovalStatusPending, store.filter.Status)
	require.Equal(t, maxPendingApprovalsCounted, store.filter.Limit)
}

// TestActApprovalPushSurvivesUncountableBadge pins the priority between the two: the notification is what people
// act on, the badge is a convenience. A store that cannot be read costs the number, not the notification.
func TestActApprovalPushSurvivesUncountableBadge(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())
	notifier.runs = &actPushRunStore{err: errors.New("act store unreachable")}

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "cobranza",
		act.PendingAction{Tool: "create_issue", Connector: "jira_ops"}))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, []string{"ana@example.com"}, pushes[0].Recipients)
	// No number rather than a wrong one: the app leaves the icon as it was.
	require.Nil(t, pushes[0].PendingByRecipient)
}

// TestActApprovalPushBadgeNormalisesMCPToolName pins the count against who may actually decide. An MCP action is
// stored under its qualified name ("mcp.<connector>.<tool>") while an approve expression names the bare tool, so
// evaluating the stored name as-is would leave the person the policy names out of their own count: the push would
// carry a zero and the number would appear the moment they opened the app, which resolves the same approval
// through RawToolName. Found by Codex reviewing kairos-cloud#147.
func TestActApprovalPushBadgeNormalisesMCPToolName(t *testing.T) {
	notifier, instanceID := newActPushNotifier(t, actPushMembers())
	// The same refund, reached over MCP rather than the built-in connector.
	notifier.runs = &actPushRunStore{approvals: []*act.Approval{
		actStoredApproval(instanceID, "finanzas", "alerta/2025", "mcp.core.issue_refund", "core"),
	}}

	notifier.NotifyPendingApprovals(t.Context(), actPendingApprovals(instanceID, "finanzas",
		act.PendingAction{Tool: "issue_refund", Connector: "core"},
	))

	pushes := testruntime.PushNotifications(instanceID)
	require.Len(t, pushes, 1)
	require.Equal(t, map[string]int32{
		"ana@example.com":    1,
		"anafin@example.com": 1,
	}, pushes[0].PendingByRecipient)
}
