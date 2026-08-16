package server_test

import (
	"context"
	"testing"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/ratelimit"
	"github.com/rilldata/rill/runtime/server"
	"github.com/rilldata/rill/runtime/server/auth"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// These tests exercise the Act permission model of issue #135 end to end against the handlers: per-agent
// access resolved by the security engine, the launch/approve gates, the WHERE-clause filtering of runs and
// approvals, the confinement of public-link (exclusive rule) tokens, and the agents kill switch.

// actSecurityProjectFiles is a project with four agents covering the policy shapes:
//   - triage: no security block, so only admins see or launch it.
//   - cobranza: access for operaciones+direccion, launch only operaciones, approve only direccion (disjoint
//     launch/approve: whoever launches cannot sign).
//   - abierta: access for everyone, no launch (inherits access), no approve (admins only).
//   - finanzas: approve discriminates by action (direccion signs everything; Ana additionally signs refunds),
//     the issue #135 example that makes approval authority a property of one approval, not of the agent.
func actSecurityProjectFiles() map[string]string {
	return map[string]string{
		"rill.yaml": `
features:
  agents: true
`,
		"triage.yaml": `
type: agent
display_name: Ticket Triage
instructions: "Investiga la alerta y propon un ticket."
`,
		"cobranza.yaml": `
type: agent
display_name: Cobranza
instructions: "Investiga la morosidad y propone la tarea."
security:
  access: '{{ or (has "operaciones" .user.groups) (has "direccion" .user.groups) }}'
  launch: '{{ has "operaciones" .user.groups }}'
  approve: '{{ has "direccion" .user.groups }}'
`,
		"abierta.yaml": `
type: agent
display_name: Abierta
instructions: "Narra el dashboard."
security:
  access: "true"
`,
		"finanzas.yaml": `
type: agent
display_name: Finanzas
instructions: "Investiga y propone acciones financieras."
security:
  access: "true"
  launch: '{{ has "operaciones" .user.groups }}'
  approve: >-
    {{ or (has "direccion" .user.groups)
          (and (eq .action.tool "issue_refund") (eq .user.email "usr_anafin@example.com")) }}
`,
	}
}

// userCtx returns a context carrying real (non-skip) claims for a project member with the given groups.
func userCtx(userID string, groups []any, admin bool, extraPerms ...runtime.Permission) context.Context {
	perms := append([]runtime.Permission{runtime.ReadObjects, runtime.ReadMetrics, runtime.ReadAPI, runtime.UseAI}, extraPerms...)
	return auth.WithClaims(context.Background(), &runtime.SecurityClaims{
		UserID: userID,
		UserAttributes: map[string]any{
			"id":     userID,
			"email":  userID + "@example.com",
			"domain": "example.com",
			"groups": groups,
			"admin":  admin,
		},
		Permissions: perms,
	})
}

// publicLinkAdminCtx returns claims shaped like a magic auth token created BY AN ADMIN. This is the
// dangerous shape and the reason no authorization shortcut may key on the "admin" attribute: the admin
// service copies the creator's attributes into the token (admin/server/magic_tokens.go, with upstream's own
// warning), so attrs["admin"] is true here — while the token's permissions, derived from
// ProjectPermissionsForMagicAuthToken, still carry no EditTrigger. Everything must deny it exactly as it
// denies the non-admin link.
func publicLinkAdminCtx() context.Context {
	return auth.WithClaims(context.Background(), &runtime.SecurityClaims{
		UserID:         "",
		UserAttributes: map[string]any{"admin": true},
		Permissions:    []runtime.Permission{runtime.ReadObjects, runtime.ReadMetrics, runtime.ReadAPI, runtime.UseAI},
		AdditionalRules: []*runtimev1.SecurityRule{
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
	})
}

// embedAdminCtx returns claims shaped like an EMBED token whose integrator declared admin: true, and which
// carries no exclusive rule to confine it. This is the residual we accept on the access side — the built-in
// rule keys on the attribute, the same posture upstream takes for alerts and reports — so this principal does
// see agents. It exists to pin the line that residual must not cross: seeing an agent is not signing for it,
// so every gate past access has to deny it.
func embedAdminCtx() context.Context {
	return auth.WithClaims(context.Background(), &runtime.SecurityClaims{
		UserID:         "",
		UserAttributes: map[string]any{"admin": true},
		Permissions:    []runtime.Permission{runtime.ReadObjects, runtime.ReadMetrics, runtime.ReadAPI, runtime.UseAI},
	})
}

// publicLinkCtx returns a context carrying claims shaped like a magic auth (public link) token's: full
// instance permissions (they all get UseAI) plus an exclusive access rule for one explore, which the engine
// expands into "deny everything else".
func publicLinkCtx() context.Context {
	return auth.WithClaims(context.Background(), &runtime.SecurityClaims{
		UserID:         "",
		UserAttributes: map[string]any{"admin": false},
		Permissions:    []runtime.Permission{runtime.ReadObjects, runtime.ReadMetrics, runtime.ReadAPI, runtime.UseAI},
		AdditionalRules: []*runtimev1.SecurityRule{
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
	})
}

// newActSecurityServer builds a server over the three-agent project with the Act plane wired.
func newActSecurityServer(t *testing.T) (*server.Server, *act.PostgresRunStore, *fakeAgentExecutor, string) {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: actSecurityProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)

	store := newActStore(t)
	exec := &fakeAgentExecutor{store: store}
	srv.ConfigureAct(store, exec)
	return srv, store, exec, instanceID
}

func agentNames(list *runtimev1.ListAgentsResponse) []string {
	names := make([]string, len(list.Agents))
	for i, a := range list.Agents {
		names[i] = a.Name
	}
	return names
}

func agentByName(t *testing.T, list *runtimev1.ListAgentsResponse, name string) *runtimev1.AgentDefinition {
	t.Helper()
	for _, a := range list.Agents {
		if a.Name == name {
			return a
		}
	}
	t.Fatalf("agent %q not in listing", name)
	return nil
}

// TestAgentSecurityDiscovery covers who sees which agent, with which gates, and that a hidden agent is
// NotFound rather than Forbidden (no existence oracle). A public-link token enumerates nothing.
func TestAgentSecurityDiscovery(t *testing.T) {
	srv, _, _, instanceID := newActSecurityServer(t)

	adminCtx := userCtx("usr_admin", []any{}, true, runtime.EditTrigger)
	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)
	dirCtx := userCtx("usr_ana", []any{"direccion"}, false)
	salesCtx := userCtx("usr_sam", []any{"ventas"}, false)

	// Admins see everything; without a security block the defaults keep every gate theirs.
	list, err := srv.ListAgents(adminCtx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"triage", "cobranza", "abierta", "finanzas"}, agentNames(list))
	triage := agentByName(t, list, "triage")
	require.True(t, triage.CanLaunch)
	// Admins are not in operaciones: the declared launch expression binds them too.
	cobranza := agentByName(t, list, "cobranza")
	require.False(t, cobranza.CanLaunch)

	// Operaciones: cobranza (launch yes) and the open agents; triage is invisible.
	list, err = srv.ListAgents(opsCtx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"cobranza", "abierta", "finanzas"}, agentNames(list))
	cobranza = agentByName(t, list, "cobranza")
	require.True(t, cobranza.CanLaunch)
	// Launch inherits access when absent.
	abierta := agentByName(t, list, "abierta")
	require.True(t, abierta.CanLaunch)

	// Dirección: access yes, launch no. The disjoint launch set is visible in the gate.
	list, err = srv.ListAgents(dirCtx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	cobranza = agentByName(t, list, "cobranza")
	require.False(t, cobranza.CanLaunch)

	// Ventas only sees the open agents; the hidden ones are NotFound by name.
	list, err = srv.ListAgents(salesCtx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"abierta", "finanzas"}, agentNames(list))
	_, err = srv.GetAgent(salesCtx, &runtimev1.GetAgentRequest{InstanceId: instanceID, Name: "cobranza"})
	require.Equal(t, codes.NotFound, status.Code(err))

	// A public-link token enumerates no agents and reads none, even the access-for-everyone one: its exclusive
	// rule denies everything it does not name, and it never names an agent.
	list, err = srv.ListAgents(publicLinkCtx(), &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Empty(t, list.Agents)
	_, err = srv.GetAgent(publicLinkCtx(), &runtimev1.GetAgentRequest{InstanceId: instanceID, Name: "abierta"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

// TestAgentSecurityLaunchGate covers StartAgentRun: launch is granted by the expression (or inherited from
// access), requires no EditTrigger, and access without launch is a plain permission error while no access
// reads as NotFound.
func TestAgentSecurityLaunchGate(t *testing.T) {
	srv, _, _, instanceID := newActSecurityServer(t)

	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)
	dirCtx := userCtx("usr_ana", []any{"direccion"}, false)
	salesCtx := userCtx("usr_sam", []any{"ventas"}, false)

	// A member of operaciones launches without holding EditTrigger (no manage_project involved).
	start, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "cobranza", IdempotencyKey: "k1", Prompt: "Arranca."})
	require.NoError(t, err)
	require.Equal(t, "cobranza", start.AgentName)

	// Access without launch: the agent is visible but launching is denied.
	_, err = srv.StartAgentRun(dirCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "cobranza", IdempotencyKey: "k2", Prompt: "Arranca."})
	require.Equal(t, codes.PermissionDenied, status.Code(err))

	// No access: the agent does not even exist for the caller.
	_, err = srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k3", Prompt: "Arranca."})
	require.Equal(t, codes.NotFound, status.Code(err))

	// An absent launch inherits access: any accessor can launch (this is what keeps the canvas agent note
	// working for a plain viewer with access).
	_, err = srv.StartAgentRun(salesCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "abierta", IdempotencyKey: "k4", Prompt: "Arranca."})
	require.NoError(t, err)

	// A public-link token cannot launch anything.
	_, err = srv.StartAgentRun(publicLinkCtx(), &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "abierta", IdempotencyKey: "k5"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

// TestAgentSecurityRunsFilteredInWhere covers the read side of runs: the access restriction lands in the
// store's WHERE clause (a page with more inaccessible candidates than its size still comes back full), a
// hidden agent's run is NotFound by id, and a public-link token reads nothing.
func TestAgentSecurityRunsFilteredInWhere(t *testing.T) {
	srv, store, _, instanceID := newActSecurityServer(t)

	adminCtx := userCtx("usr_admin", []any{}, true, runtime.EditTrigger)
	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)

	// Two cobranza runs first (older), then three triage runs (newer): a newest-first page of two without the
	// WHERE restriction would contain only triage rows, which operaciones cannot see.
	startRun := func(ctx context.Context, agent, key string) string {
		t.Helper()
		res, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: agent, IdempotencyKey: key, Prompt: "Arranca."})
		require.NoError(t, err)
		return res.RunId
	}
	cob1 := startRun(opsCtx, "cobranza", "c1")
	cob2 := startRun(opsCtx, "cobranza", "c2")
	time.Sleep(20 * time.Millisecond) // separate created_on so newest-first ordering is deterministic
	tri1 := startRun(adminCtx, "triage", "t1")
	startRun(adminCtx, "triage", "t2")
	startRun(adminCtx, "triage", "t3")

	// The page is full of accessible rows: the filter ran in the WHERE, not on the fetched page.
	page, err := srv.ListAgentRuns(opsCtx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID, PageSize: 2})
	require.NoError(t, err)
	require.Len(t, page.Runs, 2)
	require.ElementsMatch(t, []string{cob1, cob2}, []string{page.Runs[0].RunId, page.Runs[1].RunId})

	// Admins see all runs.
	all, err := srv.ListAgentRuns(adminCtx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, all.Runs, 5)

	// A hidden agent's run is NotFound by exact id, indistinguishable from a missing run.
	_, err = srv.GetAgentRun(opsCtx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: tri1})
	require.Equal(t, codes.NotFound, status.Code(err))
	got, err := srv.GetAgentRun(opsCtx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: cob1})
	require.NoError(t, err)
	require.Equal(t, cob1, got.Run.RunId)

	// A public-link token reads no runs at all, by list or by id, and cannot stream events.
	list, err := srv.ListAgentRuns(publicLinkCtx(), &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Empty(t, list.Runs)
	_, err = srv.GetAgentRun(publicLinkCtx(), &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: cob1})
	require.Equal(t, codes.NotFound, status.Code(err))
	stream := &collectStream{ctx: publicLinkCtx()}
	err = srv.StreamAgentRunEvents(&runtimev1.StreamAgentRunEventsRequest{InstanceId: instanceID, RunId: cob1}, stream)
	require.Equal(t, codes.NotFound, status.Code(err))
	require.Empty(t, stream.events)

	// And the same holds for a link created BY AN ADMIN, which inherits attrs["admin"] = true. If any
	// shortcut on the read path keyed on that attribute instead of on a real permission, this link would read
	// every run in the project, prompts and proposed actions included.
	list, err = srv.ListAgentRuns(publicLinkAdminCtx(), &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Empty(t, list.Runs, "a share link created by an admin must not read runs: the admin attribute is inherited, the permission is not")
	_, err = srv.GetAgentRun(publicLinkAdminCtx(), &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: cob1})
	require.Equal(t, codes.NotFound, status.Code(err))

	// Sanity: the runs exist in the store; only the read path hides them.
	raw, err := store.ListRuns(context.Background(), act.ListRunsFilter{InstanceID: instanceID})
	require.NoError(t, err)
	require.Len(t, raw, 5)
}

// TestAgentSecurityApprovalByGroup covers the approve gate: a group member decides without holding
// EditTrigger (no manage_project), the launcher of a disjoint policy cannot sign what they started, an
// operator can still break-glass with EditTrigger, and approvals stay invisible outside access.
func TestAgentSecurityApprovalByGroup(t *testing.T) {
	srv, store, exec, instanceID := newActSecurityServer(t)
	seedApproval := func(runID, proposal string) string {
		t.Helper()
		approvalID := act.ApprovalIDForRun(runID)
		require.NoError(t, store.CreateApproval(context.Background(), act.NewApproval{
			ApprovalID: approvalID, RunID: runID, InstanceID: instanceID,
			ToolName: "create_crm_task", Connector: "crm", ArgsHash: act.HashArgs(proposal), CanonicalArgs: proposal, Proposal: proposal,
			RequestedBy: "usr_marta",
		}))
		return approvalID
	}

	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)
	dirCtx := userCtx("usr_ana", []any{"direccion"}, false)
	salesCtx := userCtx("usr_sam", []any{"ventas"}, false)
	adminCtx := userCtx("usr_admin", []any{}, true, runtime.EditTrigger)

	start, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "cobranza", IdempotencyKey: "k1", Prompt: "Arranca."})
	require.NoError(t, err)
	approvalID := seedApproval(start.RunId, "crear tarea de cobranza")

	// Dirección sees the pending approval and may decide it; operaciones sees it too (access governs reads)
	// but cannot decide, which is what the UI's can_decide bit reflects. Ventas has no access at all.
	inbox, err := srv.ListAgentApprovals(dirCtx, &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Len(t, inbox.Approvals, 1)
	require.True(t, inbox.Approvals[0].CanDecide)
	opsInbox, err := srv.ListAgentApprovals(opsCtx, &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Len(t, opsInbox.Approvals, 1)
	require.False(t, opsInbox.Approvals[0].CanDecide)
	salesInbox, err := srv.ListAgentApprovals(salesCtx, &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Empty(t, salesInbox.Approvals)
	_, err = srv.GetAgentApproval(salesCtx, &runtimev1.GetAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID})
	require.Equal(t, codes.NotFound, status.Code(err))
	publicInbox, err := srv.ListAgentApprovals(publicLinkCtx(), &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Empty(t, publicInbox.Approvals)
	// Same for a link created by an admin, which inherits attrs["admin"] = true but no EditTrigger: it must
	// neither read the pending approvals nor be able to decide one.
	adminLinkInbox, err := srv.ListAgentApprovals(publicLinkAdminCtx(), &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Empty(t, adminLinkInbox.Approvals, "a share link created by an admin must not see pending approvals")
	_, err = srv.ApproveAgentApproval(publicLinkAdminCtx(), &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs("crear tarea de cobranza"),
	})
	require.Equal(t, codes.NotFound, status.Code(err))

	// Whoever launched cannot sign: launch and approve are disjoint groups, so the run's own initiator is
	// denied even though they can see the approval.
	_, err = srv.ApproveAgentApproval(opsCtx, &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs("crear tarea de cobranza"),
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	require.Empty(t, exec.resumeCalls())

	// A member of dirección decides without holding EditTrigger: approval authority comes from the policy,
	// not from manage_project.
	approve, err := srv.ApproveAgentApproval(dirCtx, &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs("crear tarea de cobranza"),
	})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, approve.Approval.Status)
	require.Equal(t, "usr_ana", approve.Approval.DecidedBy)
	require.Len(t, exec.resumeCalls(), 1)

	// EditTrigger break-glass: an operator outside dirección can still deny a second run's approval.
	start2, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "cobranza", IdempotencyKey: "k2", Prompt: "Arranca."})
	require.NoError(t, err)
	approval2 := seedApproval(start2.RunId, "otra tarea")
	deny, err := srv.DenyAgentApproval(adminCtx, &runtimev1.DenyAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approval2})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusDenied, deny.Approval.Status)
}

// TestAgentApprovalDeniedToInheritedAdminAttribute covers the gap between the two halves of the admin
// attribute. An embed or share link copies attrs["admin"] from its creator but never gets EditTrigger, so a
// gate that reads the attribute is forgeable by anyone an admin ever shared with. On the access side that is
// the accepted residual; on the approve side it would let a shared link sign a write to the outside world.
//
// The agent here declares access: "true" and no approve:, which is the default posture and the one most
// projects will have. The token therefore reaches the approval — and must still be refused at it.
func TestAgentApprovalDeniedToInheritedAdminAttribute(t *testing.T) {
	srv, store, exec, instanceID := newActSecurityServer(t)

	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)
	start, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "abierta", IdempotencyKey: "k1", Prompt: "Arranca."})
	require.NoError(t, err)
	approvalID := act.ApprovalIDForRun(start.RunId)
	require.NoError(t, store.CreateApproval(context.Background(), act.NewApproval{
		ApprovalID: approvalID, RunID: start.RunId, InstanceID: instanceID,
		ToolName: "mcp.crm.create_crm_task", Connector: "crm",
		ArgsHash: act.HashArgs("crear tarea"), CanonicalArgs: "crear tarea", Proposal: "crear tarea", RequestedBy: "usr_marta",
	}))

	// The token does reach the approval, because access is open. That is the point: the denial has to come
	// from the approve gate, not from the token failing to see the agent.
	inbox, err := srv.ListAgentApprovals(embedAdminCtx(), &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Len(t, inbox.Approvals, 1)
	require.False(t, inbox.Approvals[0].CanDecide, "an inherited admin attribute must not present itself as able to decide")

	_, err = srv.ApproveAgentApproval(embedAdminCtx(), &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs("crear tarea"),
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = srv.DenyAgentApproval(embedAdminCtx(), &runtimev1.DenyAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	require.Empty(t, exec.resumeCalls(), "no decision may have reached the executor")

	// A real project admin holds the permission, not just the attribute, and decides the same approval.
	approve, err := srv.ApproveAgentApproval(userCtx("usr_admin", []any{}, true, runtime.EditTrigger), &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs("crear tarea"),
	})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, approve.Approval.Status)
}

// TestAgentApprovalRefusedWithoutAVerifiedPreimage pins the guarantee where it has to live: in the server.
// Showing the approver the exact arguments is a client behaviour, and a guarantee that depends on the client
// behaving is not a guarantee — the hash check alone only proves the caller echoed a hash back, never that a
// human saw what it covers. So when there is no verified preimage there is nothing that could have been shown,
// and approving is refused.
//
// Both states are reproduced: a row from before the column existed, and one whose bytes were corrupted after the
// fact. Denying stays open in both, so an approval that cannot be signed can still be closed and re-proposed.
func TestAgentApprovalRefusedWithoutAVerifiedPreimage(t *testing.T) {
	srv, store, exec, instanceID := newActSecurityServer(t)
	ctx := userCtx("usr_admin", []any{}, true, runtime.EditTrigger)

	seed := func(key string) string {
		t.Helper()
		start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "abierta", IdempotencyKey: key, Prompt: "Arranca."})
		require.NoError(t, err)
		approvalID := act.ApprovalIDForRun(start.RunId)
		require.NoError(t, store.CreateApproval(context.Background(), act.NewApproval{
			ApprovalID: approvalID, RunID: start.RunId, InstanceID: instanceID,
			ToolName: "mcp.crm.create_crm_task", Connector: "crm",
			ArgsHash: act.HashArgs("crear tarea"), CanonicalArgs: "crear tarea", Proposal: "crear tarea",
			RequestedBy: "usr_admin",
		}))
		return approvalID
	}

	// A row recorded before the preimage was persisted: its bytes were never stored and cannot be recovered.
	preMigration := seed("k-old")
	require.NoError(t, store.ClobberCanonicalArgsForTest(context.Background(), instanceID, preMigration, nil))
	_, err := srv.ApproveAgentApproval(ctx, &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: preMigration, ArgsHash: act.HashArgs("crear tarea"),
	})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "deny it", "the error has to tell the operator what their way out is")
	require.Empty(t, exec.resumeCalls(), "a refused approval never reaches the executor")

	// Denying it is still allowed, so the run is not stranded: the agent can propose the action again.
	deny, err := srv.DenyAgentApproval(ctx, &runtimev1.DenyAgentApprovalRequest{InstanceId: instanceID, ApprovalId: preMigration})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusDenied, deny.Approval.Status)

	// Bytes that contradict the hash: the same refusal, reached the other way.
	corrupted := seed("k-corrupt")
	other := "crear OTRA tarea"
	require.NoError(t, store.ClobberCanonicalArgsForTest(context.Background(), instanceID, corrupted, &other))
	_, err = srv.ApproveAgentApproval(ctx, &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: corrupted, ArgsHash: act.HashArgs("crear tarea"),
	})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))

	// And the corrupted bytes are never served as the signed material either.
	got, err := srv.GetAgentApproval(ctx, &runtimev1.GetAgentApprovalRequest{InstanceId: instanceID, ApprovalId: corrupted})
	require.NoError(t, err)
	require.Empty(t, got.Approval.CanonicalArgs, "bytes that contradict the hash are omitted, not shown")
}

// TestAgentApprovalCanDecidePerAction covers the issue #135 example that makes approval authority a property
// of one approval, not of the agent: direccion signs everything, and Ana additionally signs refunds. Ana's
// clause only matches a concrete action, so a per-agent bit would hide her buttons; the per-approval
// can_decide must say true for her refund and false for anything else, and the endpoints must agree.
func TestAgentApprovalCanDecidePerAction(t *testing.T) {
	srv, store, _, instanceID := newActSecurityServer(t)

	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)
	dirCtx := userCtx("usr_ana", []any{"direccion"}, false)
	anaCtx := userCtx("usr_anafin", []any{"finanzas-ops"}, false)

	seedApproval := func(runID, tool, proposal string) string {
		t.Helper()
		approvalID := act.ApprovalIDForRun(runID)
		require.NoError(t, store.CreateApproval(context.Background(), act.NewApproval{
			ApprovalID: approvalID, RunID: runID, InstanceID: instanceID,
			ToolName: tool, Connector: "erp", ArgsHash: act.HashArgs(proposal), CanonicalArgs: proposal, Proposal: proposal,
			RequestedBy: "usr_marta",
		}))
		return approvalID
	}
	start1, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "finanzas", IdempotencyKey: "f1", Prompt: "Arranca."})
	require.NoError(t, err)
	// Seeded with the EFFECTIVE tool name, which is what real traffic stores (mcpconn composes
	// "mcp."+connector+"."+tool). Seeding the raw name here would let the expression below pass while failing in
	// production, which is exactly what this test exists to prevent.
	refundApproval := seedApproval(start1.RunId, "mcp.erp.issue_refund", "reembolsar pedido 42")
	start2, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "finanzas", IdempotencyKey: "f2", Prompt: "Arranca."})
	require.NoError(t, err)
	deleteApproval := seedApproval(start2.RunId, "mcp.erp.delete_account", "borrar la cuenta 7")

	// Ana can decide the refund and only the refund; direccion can decide both.
	byID := func(ctx context.Context) map[string]bool {
		t.Helper()
		list, err := srv.ListAgentApprovals(ctx, &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
		require.NoError(t, err)
		out := make(map[string]bool, len(list.Approvals))
		for _, a := range list.Approvals {
			out[a.ApprovalId] = a.CanDecide
		}
		return out
	}
	require.Equal(t, map[string]bool{refundApproval: true, deleteApproval: false}, byID(anaCtx))
	require.Equal(t, map[string]bool{refundApproval: true, deleteApproval: true}, byID(dirCtx))

	// The endpoints enforce the same verdicts the bit promised.
	_, err = srv.ApproveAgentApproval(anaCtx, &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: deleteApproval, ArgsHash: act.HashArgs("borrar la cuenta 7"),
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	approve, err := srv.ApproveAgentApproval(anaCtx, &runtimev1.ApproveAgentApprovalRequest{
		InstanceId: instanceID, ApprovalId: refundApproval, ArgsHash: act.HashArgs("reembolsar pedido 42"),
	})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, approve.Approval.Status)
	require.Equal(t, "usr_anafin", approve.Approval.DecidedBy)
}

// TestAgentSecurityCancelRun covers the cancel gate: the run's actor may stop what they launched without
// holding EditTrigger, other accessors may not (and the can_cancel bit says so), a caller without access gets
// NotFound, and EditTrigger keeps working as the operator's gate.
func TestAgentSecurityCancelRun(t *testing.T) {
	srv, _, _, instanceID := newActSecurityServer(t)

	opsCtx := userCtx("usr_marta", []any{"operaciones"}, false)
	dirCtx := userCtx("usr_ana", []any{"direccion"}, false)
	salesCtx := userCtx("usr_sam", []any{"ventas"}, false)
	adminCtx := userCtx("usr_admin", []any{}, true, runtime.EditTrigger)

	start, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "cobranza", IdempotencyKey: "k1", Prompt: "Arranca."})
	require.NoError(t, err)

	// The can_cancel bit mirrors the gate: true for the actor, false for another accessor.
	got, err := srv.GetAgentRun(opsCtx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.NoError(t, err)
	require.True(t, got.Run.CanCancel)
	got, err = srv.GetAgentRun(dirCtx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.NoError(t, err)
	require.False(t, got.Run.CanCancel)

	// Another accessor cannot cancel; a caller without access does not even learn the run exists.
	_, err = srv.CancelAgentRun(dirCtx, &runtimev1.CancelAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = srv.CancelAgentRun(salesCtx, &runtimev1.CancelAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.Equal(t, codes.NotFound, status.Code(err))

	// The actor cancels their own run without holding EditTrigger.
	cancel, err := srv.CancelAgentRun(opsCtx, &runtimev1.CancelAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.NoError(t, err)
	require.Equal(t, string(act.RunStatusCancelled), cancel.Run.Status)

	// EditTrigger remains the operator's gate for runs they did not start.
	start2, err := srv.StartAgentRun(opsCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "cobranza", IdempotencyKey: "k2", Prompt: "Arranca."})
	require.NoError(t, err)
	cancel, err = srv.CancelAgentRun(adminCtx, &runtimev1.CancelAgentRunRequest{InstanceId: instanceID, RunId: start2.RunId})
	require.NoError(t, err)
	require.Equal(t, string(act.RunStatusCancelled), cancel.Run.Status)
}

// TestAgentKillSwitch covers the flag as enforced in the handlers: a project that does not enable the agents
// feature has no Act API for real (non-skip) claims, admins included.
func TestAgentKillSwitch(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: map[string]string{
		"rill.yaml": ``,
		"triage.yaml": `
type: agent
instructions: "Investiga."
`,
	}})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)

	adminCtx := userCtx("usr_admin", []any{}, true, runtime.EditTrigger)
	_, err = srv.ListAgents(adminCtx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = srv.StartAgentRun(adminCtx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", Prompt: "Arranca."})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = srv.ListAgentRuns(adminCtx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = srv.ListAgentApprovals(adminCtx, &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}
