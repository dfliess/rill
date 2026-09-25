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
//
// It also carries what each recipient has left waiting in this project, which is what the installed app paints on
// its icon (kairos-cloud#147). Alerts and reports carry no such number: they are read, not resolved.

// actApprovalPushTimeout bounds one notification. It is dispatched between two durable steps of the run, so a hung
// admin call would park the run before it even reaches the approval wait it is announcing.
const actApprovalPushTimeout = 30 * time.Second

// maxPendingApprovalsCounted bounds the listing behind the app badge count. Resolving who may decide an approval
// costs a policy evaluation per member, so the count is a bounded read, not a full scan of a backlog nobody is
// working through. Past this, a recipient's badge understates what is waiting, which is the harmless direction:
// browsers stop rendering the exact number well before it anyway.
const maxPendingApprovalsCounted = 500

// actApprovalPushNotifier implements [act.ApprovalNotifier] by pushing one notification per pause of a run to the
// members that may decide it. Best-effort by contract: every failure is logged and swallowed, because a run must
// never depend on its notification being delivered.
type actApprovalPushNotifier struct {
	runtime *runtime.Runtime
	// runs is the approval store, read to count what is still waiting for each recipient.
	runs   act.RunStore
	logger *zap.Logger
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

	resolver := n.approverResolver(pending.InstanceID, members)
	recipients, err := n.approvalRecipients(ctx, pending, resolver)
	if err != nil {
		n.warn(ctx, "Failed to resolve the approvers of an act approval", pending, err)
		return
	}
	if len(recipients) == 0 {
		return
	}

	displayName := pending.AgentName
	if res := resolver.agent(ctx, pending.AgentName); res != nil && res.GetAgent().State.ValidSpec.DisplayName != "" {
		displayName = res.GetAgent().State.ValidSpec.DisplayName
	}
	body := fmt.Sprintf("La acción %q espera tu aprobación.", pending.Actions[0].Tool)
	if len(pending.Actions) > 1 {
		body = fmt.Sprintf("%d acciones esperan tu aprobación, empezando por %q.", len(pending.Actions), pending.Actions[0].Tool)
	}

	// A count that cannot be read is not worth losing the notification over: the app badge is a convenience, the
	// notification is the thing people act on. So a failure here degrades to a push without a number, which leaves
	// the badge as it was rather than making it wrong.
	pendingByRecipient, err := n.pendingRunsByRecipient(ctx, pending.InstanceID, recipients, resolver)
	if err != nil {
		n.warn(ctx, "Failed to count pending approvals for the act approval push badge", pending, err)
	}

	sent, err := admin.SendPushNotification(ctx, &drivers.PushNotification{
		Category:           drivers.PushCategoryActApprovals,
		Recipients:         recipients,
		Title:              fmt.Sprintf("Aprobación pendiente: %s", displayName),
		Body:               body,
		LinkPath:           agentRunPath(org, project, pending.RunID),
		Tag:                fmt.Sprintf("act:%s/%s/%s", org, project, pending.RunID),
		PendingByRecipient: pendingByRecipient,
	})
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
// least one of its actions. Authority discriminates by action, so a member who may decide only one of the batch's
// actions is still told about the pause, and the inbox shows them which ones are theirs.
func (n *actApprovalPushNotifier) approvalRecipients(ctx context.Context, pending act.PendingApprovals, resolver *approverResolver) ([]string, error) {
	allowed := make(map[string]bool, len(resolver.members))
	for _, a := range pending.Actions {
		emails, err := resolver.approversOf(ctx, pending.AgentName, runtime.AgentActionContext{
			Tool:      a.Tool,
			Connector: a.Connector,
			RunActor:  pending.RunActor,
		})
		if err != nil {
			return nil, err
		}
		for _, email := range emails {
			allowed[email] = true
		}
	}
	return resolver.inMembershipOrder(allowed), nil
}

// pendingRunsByRecipient returns, for each of the given recipients, how many runs of this project are stopped
// waiting for a decision that is theirs to make. Every recipient is present, zero included: the app badge adds up
// one number per project, so "nothing waits for you here any more" is as much a fact it needs as a positive count.
//
// Runs, not approvals, is the unit, and deliberately: a run that proposed three actions in one turn is one pause,
// one notification and one row in the inbox (HomePendingApprovals), so counting its approvals would put a 3 on the
// icon for one thing to look at. It is also the number the app recomputes for itself when it opens, and the two
// have to agree or the badge would jump on every launch.
func (n *actApprovalPushNotifier) pendingRunsByRecipient(ctx context.Context, instanceID string, recipients []string, resolver *approverResolver) (map[string]int32, error) {
	counts := make(map[string]int32, len(recipients))
	for _, email := range recipients {
		counts[email] = 0
	}

	approvals, err := n.runs.ListApprovals(ctx, act.ListApprovalsFilter{
		InstanceID: instanceID,
		Status:     act.ApprovalStatusPending,
		Limit:      maxPendingApprovalsCounted,
	})
	if err != nil {
		return nil, err
	}

	// A run is counted once for a recipient however many of its approvals they may decide.
	counted := make(map[string]map[string]bool, len(counts))
	for _, a := range approvals {
		// The stored tool name is qualified for MCP actions ("mcp.<connector>.<tool>"), while an approve
		// expression is written against the bare name. Normalising here is what makes this count agree with
		// who may actually decide: resolveApprovalDecision and the notification batch both do the same.
		emails, err := resolver.approversOf(ctx, a.AgentName, runtime.AgentActionContext{
			Tool:      act.RawToolName(a.ToolName, a.Connector),
			Connector: a.Connector,
			RunActor:  a.RunActorSubject,
		})
		if err != nil {
			return nil, err
		}
		for _, email := range emails {
			if _, ok := counts[email]; !ok {
				continue // Not a recipient of this notification, so nothing is being sent to their devices.
			}
			if counted[email] == nil {
				counted[email] = make(map[string]bool)
			}
			if counted[email][a.RunID] {
				continue
			}
			counted[email][a.RunID] = true
			counts[email]++
		}
	}
	return counts, nil
}

// approverResolver answers "who may decide this action?" for one instance, remembering what it has already looked
// up. One pause and its badge count ask that question repeatedly about a handful of distinct actions, and each
// fresh answer costs a policy evaluation per project member.
type approverResolver struct {
	notifier   *actApprovalPushNotifier
	instanceID string
	members    []drivers.ProjectMember
	// agents caches resource lookups; an entry present and nil is an agent that is gone, which is a verdict of its
	// own, not a cache miss.
	agents   map[string]*runtimev1.Resource
	verdicts map[approverKey][]string
}

// approverKey is everything an approve policy may discriminate on: the agent whose policy it is, the action it is
// asked about, and the actor whose run proposed it.
type approverKey struct {
	agent     string
	tool      string
	connector string
	runActor  string
}

func (n *actApprovalPushNotifier) approverResolver(instanceID string, members []drivers.ProjectMember) *approverResolver {
	return &approverResolver{
		notifier:   n,
		instanceID: instanceID,
		members:    members,
		agents:     make(map[string]*runtimev1.Resource),
		verdicts:   make(map[approverKey][]string),
	}
}

// approversOf returns the emails of the members who may decide one action of one agent.
//
// An agent deleted (or invalidated) since its run started has no policy left to evaluate, so its approvals go to
// the operators, which is who resolveApprovalDecision falls back to in the same situation.
func (r *approverResolver) approversOf(ctx context.Context, agentName string, action runtime.AgentActionContext) ([]string, error) {
	res := r.agent(ctx, agentName)
	if res == nil {
		return r.operators(), nil
	}

	key := approverKey{agent: agentName, tool: action.Tool, connector: action.Connector, runActor: action.RunActor}
	if emails, ok := r.verdicts[key]; ok {
		return emails, nil
	}
	emails, err := r.notifier.runtime.ResolveAgentApprovers(ctx, r.instanceID, res, action, r.members)
	if err != nil {
		return nil, err
	}
	r.verdicts[key] = emails
	return emails, nil
}

// agent returns the agent's reconciled resource, or nil if it is gone or has no valid spec (the same two cases
// resolveAgentGates reports as not found).
func (r *approverResolver) agent(ctx context.Context, name string) *runtimev1.Resource {
	if res, ok := r.agents[name]; ok {
		return res
	}
	r.agents[name] = nil // Every failure below leaves the agent recorded as gone.

	ctrl, err := r.notifier.runtime.Controller(ctx, r.instanceID)
	if err != nil {
		return nil
	}
	res, err := ctrl.Get(ctx, &runtimev1.ResourceName{Kind: runtime.ResourceKindAgent, Name: name}, false)
	if err != nil || res.GetAgent().State.ValidSpec == nil {
		return nil
	}
	r.agents[name] = res
	return res
}

// operators are the members holding EditTrigger: the break-glass audience for an approval whose policy cannot be
// read any more.
func (r *approverResolver) operators() []string {
	var emails []string
	for _, m := range r.members {
		if m.Email != "" && m.EditTrigger {
			emails = append(emails, m.Email)
		}
	}
	return emails
}

// inMembershipOrder emits the selected emails in membership order (the admin service sorts it by email), so a
// recipient list is deterministic rather than following map iteration.
func (r *approverResolver) inMembershipOrder(selected map[string]bool) []string {
	emails := make([]string, 0, len(selected))
	for _, m := range r.members {
		if selected[m.Email] {
			delete(selected, m.Email)
			emails = append(emails, m.Email)
		}
	}
	return emails
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

// agentRunPath returns the frontend path of a run's detail page, which addresses the run by its ID (ADR-0034). It
// mirrors runDetailPath in web-admin/src/features/agents/run-id.ts.
func agentRunPath(org, project, runID string) string {
	return fmt.Sprintf("/%s/%s/-/agents/runs/%s", url.PathEscape(org), url.PathEscape(project), url.PathEscape(runID))
}
