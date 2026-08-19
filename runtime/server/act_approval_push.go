package server

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.uber.org/zap"
)

// This file mirrors an Act run pausing on approvals as a web push (kairos-cloud#143), the third push category after
// the alerts and reports of runtime/reconcilers/push_notifications.go. Unlike those, the recipients are not written
// anywhere: an approval's audience is whoever the agent's approve policy authorizes, so it is enumerated by
// evaluating that policy against every project member (Runtime.ResolveAgentApprovers).
//
// The push is never authority: it only tells people to come and look. The decision itself is re-checked against the
// deciding caller's own claims by authorizeApprovalDecision, exactly as it is for someone who found the approval by
// browsing the inbox.

// actApprovalPushTimeout bounds one notification. It is dispatched between two durable steps of the run, so a hung
// admin call would park the run before it even reaches the approval wait it is announcing.
const actApprovalPushTimeout = 30 * time.Second

// actApprovalPushNotifier implements [act.ApprovalNotifier] by pushing one notification per pause of a run to the
// members that may decide it. Best-effort by contract: every failure is logged and swallowed, because a run must
// never depend on its notification being delivered.
type actApprovalPushNotifier struct {
	runtime *runtime.Runtime
	logger  *zap.Logger
}

var _ act.ApprovalNotifier = (*actApprovalPushNotifier)(nil)

// NotifyPendingApprovals sends one web push for the batch to everyone the agent's approve policy lets decide any of
// its actions. One batch is one moment of signing, so it is one notification however many actions it holds; the
// body carries the count and the tag is stable per run, so a second pause replaces the first in the browser.
func (n *actApprovalPushNotifier) NotifyPendingApprovals(ctx context.Context, pending act.PendingApprovals) {
	if len(pending.Actions) == 0 {
		return
	}

	// Outside Rill Cloud there is no frontend to link to and no subscription to push to (Rill Developer).
	org, project := n.runtime.InstanceOrgProject(ctx, pending.InstanceID)
	if org == "" || project == "" {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, actApprovalPushTimeout)
	defer cancel()

	admin, release, err := n.runtime.Admin(ctx, pending.InstanceID)
	if err != nil {
		n.warn(ctx, "Failed to acquire admin client for act approval push", pending, err)
		return
	}
	defer release()

	members, err := admin.ListProjectMemberAttributes(ctx)
	if err != nil {
		n.warn(ctx, "Failed to list project members for act approval push", pending, err)
		return
	}

	res := n.agentResource(ctx, pending.InstanceID, pending.AgentName)
	recipients, err := n.approvalRecipients(ctx, pending, res, members)
	if err != nil {
		n.warn(ctx, "Failed to resolve the approvers of an act approval", pending, err)
		return
	}
	if len(recipients) == 0 {
		return
	}

	displayName := pending.AgentName
	if res != nil && res.GetAgent().State.ValidSpec.DisplayName != "" {
		displayName = res.GetAgent().State.ValidSpec.DisplayName
	}
	body := fmt.Sprintf("La acción %q espera tu aprobación.", pending.Actions[0].Tool)
	if len(pending.Actions) > 1 {
		body = fmt.Sprintf("%d acciones esperan tu aprobación, empezando por %q.", len(pending.Actions), pending.Actions[0].Tool)
	}

	sent, err := admin.SendPushNotification(ctx,
		drivers.PushCategoryActApprovals,
		recipients,
		fmt.Sprintf("Aprobación pendiente: %s", displayName),
		body,
		agentRunPath(org, project, pending.AgentName, pending.IdempotencyKey),
		fmt.Sprintf("act:%s/%s/%s", org, project, pending.RunID),
	)
	if err != nil {
		n.warn(ctx, "Failed to send act approval push", pending, err)
		return
	}

	n.logger.Debug("Sent act approval push notifications",
		zap.String("instance_id", pending.InstanceID),
		zap.String("run_id", pending.RunID),
		zap.Int("recipients", len(recipients)),
		zap.Int("sent", sent),
		observability.ZapCtx(ctx),
	)
}

// approvalRecipients returns the emails to notify about the batch: everyone the approve policy authorizes for at
// least one of its actions. Authority discriminates by action, so an action is resolved once per distinct
// (tool, connector) pair and the results are unioned; a member who may decide only one of the batch's actions is
// still told about the pause, and the inbox shows them which ones are theirs.
//
// A nil res is an agent deleted (or invalidated) since the run started: its policy is gone with it, so the
// notification goes to the operators, which is who resolveApprovalDecision falls back to in the same situation.
func (n *actApprovalPushNotifier) approvalRecipients(ctx context.Context, pending act.PendingApprovals, res *runtimev1.Resource, members []drivers.ProjectMember) ([]string, error) {
	allowed := make(map[string]bool, len(members))
	if res == nil {
		for _, m := range members {
			if m.Email != "" && m.EditTrigger {
				allowed[m.Email] = true
			}
		}
	} else {
		resolved := make(map[act.PendingAction]bool, len(pending.Actions))
		for _, a := range pending.Actions {
			if resolved[a] {
				continue
			}
			resolved[a] = true

			emails, err := n.runtime.ResolveAgentApprovers(ctx, pending.InstanceID, res, runtime.AgentActionContext{
				Tool:      a.Tool,
				Connector: a.Connector,
				RunActor:  pending.RunActor,
			}, members)
			if err != nil {
				return nil, err
			}
			for _, email := range emails {
				allowed[email] = true
			}
		}
	}

	// Emit in membership order (the admin service sorts it by email) so the recipient list is deterministic.
	recipients := make([]string, 0, len(allowed))
	for _, m := range members {
		if allowed[m.Email] {
			delete(allowed, m.Email)
			recipients = append(recipients, m.Email)
		}
	}
	return recipients, nil
}

// agentResource returns the agent's reconciled resource, or nil if it is gone or has no valid spec (the same two
// cases resolveAgentGates reports as not found).
func (n *actApprovalPushNotifier) agentResource(ctx context.Context, instanceID, name string) *runtimev1.Resource {
	ctrl, err := n.runtime.Controller(ctx, instanceID)
	if err != nil {
		return nil
	}
	res, err := ctrl.Get(ctx, &runtimev1.ResourceName{Kind: runtime.ResourceKindAgent, Name: name}, false)
	if err != nil || res.GetAgent().State.ValidSpec == nil {
		return nil
	}
	return res
}

// warn logs a failed notification. A cancelled context (shutdown) and an admin service without push (Rill
// Developer, self-hosted without the admin plane) are expected outcomes, not incidents.
func (n *actApprovalPushNotifier) warn(ctx context.Context, msg string, pending act.PendingApprovals, err error) {
	if errors.Is(err, drivers.ErrNotImplemented) || errors.Is(err, context.Canceled) {
		return
	}
	n.logger.Warn(msg,
		zap.String("instance_id", pending.InstanceID),
		zap.String("run_id", pending.RunID),
		zap.String("agent", pending.AgentName),
		zap.Error(err),
		observability.ZapCtx(ctx),
	)
}

// agentRunPath returns the frontend path of a run's detail page. It mirrors runDetailPath in
// web-admin/src/features/agents/run-id.ts: the URL carries the agent name and the run's idempotency key, NOT the
// composite run ID, which the frontend recomposes from the path.
func agentRunPath(org, project, agentName, idempotencyKey string) string {
	return fmt.Sprintf("/%s/%s/-/agents/%s/runs/%s",
		url.PathEscape(org), url.PathEscape(project), url.PathEscape(agentName), url.PathEscape(idempotencyKey))
}
