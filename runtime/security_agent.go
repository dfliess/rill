package runtime

import (
	"context"
	"fmt"
	"maps"
	"slices"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/parser"
)

// AgentGates are one caller's resolved permissions on one agent (Kairos Act): whether they can see it (and
// its runs and approvals) and start a run of it. Launch is a subset of access by construction: a caller
// without access holds no gate at all.
//
// There is deliberately no per-agent approve gate: approval authority can discriminate by action
// (.action.tool/.action.connector), so it is a property of one concrete approval and is resolved with
// ResolveAgentApprove instead.
type AgentGates struct {
	Access bool
	Launch bool
}

// AgentActionContext is the context of one proposed action for evaluating an agent's approve expression.
// Beyond `.user`, the expression may reference `.action.tool`, `.action.connector` and `.run.actor`.
// The action's arguments are deliberately NOT part of the context: they are model-chosen, so a rule that
// read them would let the prompt influence who may approve the action.
type AgentActionContext struct {
	// Tool and Connector identify the proposed action.
	Tool      string
	Connector string
	// RunActor is the subject the run acts for (for a manual run, who launched it).
	RunActor string
}

// ResolveAgentGates resolves the caller's gates on an agent resource. Access comes from the security engine
// (the agent's declared security rules plus the built-in admin rule, see resolveRules). Launch is a spec
// expression evaluated against the same template context as security conditions.
func (r *Runtime) ResolveAgentGates(ctx context.Context, instanceID string, claims *SecurityClaims, res *runtimev1.Resource) (AgentGates, error) {
	if claims.SkipChecks {
		return AgentGates{Access: true, Launch: true}, nil
	}

	resolved, err := r.ResolveSecurity(ctx, instanceID, claims, res)
	if err != nil {
		return AgentGates{}, err
	}
	if !resolved.CanAccess() {
		return AgentGates{}, nil
	}

	spec := agentSpecForGates(res)
	gates := AgentGates{Access: true}

	if spec.LaunchExpression == "" {
		// An absent launch inherits access (already granted here): starting a run is inert until an action
		// passes the approval gate, so it opens with visibility by default.
		gates.Launch = true
	} else {
		gates.Launch, err = r.evaluateAgentGate(ctx, instanceID, claims, res, spec.LaunchExpression, nil)
		if err != nil {
			return AgentGates{}, fmt.Errorf("failed to evaluate launch expression of agent %q: %w", res.Meta.Name.Name, err)
		}
	}
	return gates, nil
}

// ResolveAgentApprove reports whether the caller may decide an approval of the given proposed action of the
// agent. It requires access (approve is a subset of access by construction) and then evaluates the agent's
// approve expression with the action bound. It does NOT apply the EditTrigger break-glass: that override is
// the caller's decision, so it can be audited where it is used.
func (r *Runtime) ResolveAgentApprove(ctx context.Context, instanceID string, claims *SecurityClaims, res *runtimev1.Resource, action AgentActionContext) (bool, error) {
	if claims.SkipChecks {
		return true, nil
	}

	resolved, err := r.ResolveSecurity(ctx, instanceID, claims, res)
	if err != nil {
		return false, err
	}
	if !resolved.CanAccess() {
		return false, nil
	}
	return r.resolveAgentApprove(ctx, instanceID, claims, res, agentSpecForGates(res), action)
}

// resolveAgentApprove evaluates the approve gate for a caller that is already known to have access.
func (r *Runtime) resolveAgentApprove(ctx context.Context, instanceID string, claims *SecurityClaims, res *runtimev1.Resource, spec *runtimev1.AgentSpec, action AgentActionContext) (bool, error) {
	if spec.ApproveExpression == "" {
		// An absent approve means only project admins: deciding an action touches the outside world, so unlike
		// launch it does not inherit access. The gate is the EditTrigger permission, matching the pre-policy
		// behavior of the approval endpoints.
		//
		// It is deliberately NOT claims.Admin(). That reads the "admin" user attribute, which a magic-link or
		// embed token copies from whoever created it, so an admin sharing a dashboard would mint a token that
		// signs actions. Access may still be granted by the attribute (the built-in rule in resolveRules, same
		// posture upstream takes for alerts and reports), but seeing an agent must not imply signing for it.
		return claims.Can(EditTrigger), nil
	}

	extra := map[string]any{
		"action": map[string]any{
			"tool":      action.Tool,
			"connector": action.Connector,
		},
		"run": map[string]any{
			"actor": action.RunActor,
		},
	}
	ok, err := r.evaluateAgentGate(ctx, instanceID, claims, res, spec.ApproveExpression, extra)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate approve expression of agent %q: %w", res.Meta.Name.Name, err)
	}
	return ok, nil
}

// ResolveAgentApprovers returns the emails of the project members who may decide one proposed action of an agent.
// It is the inverse of ResolveAgentApprove, which answers the same question for one caller: a notification has to be
// addressed before anyone asks for it (kairos-cloud#143). It grants nothing — the decision itself is always
// re-resolved against the deciding caller's own claims — so a member listed here still has to pass the same gate
// when they act.
//
// Each member is evaluated with the claims their own runtime token would carry (see agentApproverClaims), so the
// policy sees here exactly what it will see when they open the approval. Members without an email are skipped:
// there is nobody to notify. A member whose evaluation errors fails the whole enumeration, because an approve
// expression is a property of the agent, not of the member: it would fail for everyone.
//
// When the policy matches nobody, it falls back to the members holding EditTrigger. An approval nobody is told
// about is a run that waits forever, and those members are exactly who the API's break-glass
// (authorizeApprovalDecision) lets decide it anyway.
func (r *Runtime) ResolveAgentApprovers(ctx context.Context, instanceID string, res *runtimev1.Resource, action AgentActionContext, members []drivers.ProjectMember) ([]string, error) {
	var approvers, operators []string
	for _, m := range members {
		if m.Email == "" {
			continue
		}
		if m.EditTrigger {
			operators = append(operators, m.Email)
		}

		// ResolveAgentApprove resolves access itself and denies without it, so this is the same verdict
		// resolveApprovalDecision reaches for a caller, minus its EditTrigger break-glass.
		ok, err := r.ResolveAgentApprove(ctx, instanceID, agentApproverClaims(m), res, action)
		if err != nil {
			return nil, err
		}
		if ok {
			approvers = append(approvers, m.Email)
		}
	}

	if len(approvers) == 0 {
		return operators, nil
	}
	return approvers, nil
}

// agentApproverPermissions are the instance permissions every member's runtime token carries whatever their role
// is (see issueRuntimeToken in admin/server/runtime_jwt.go), so they never discriminate between members. The
// permissions that do — EditTrigger above all — depend on the member's role in the deployment's environment, and
// the admin service resolves them per member into drivers.ProjectMember.EditTrigger.
var agentApproverPermissions = []Permission{ReadAPI, ReadMetrics, ReadObjects, UseAI}

// agentApproverClaims rebuilds the SecurityClaims a runtime JWT issued for the member would resolve to (the
// jwtClaims provider in runtime/server/auth): their attributes with "id" defaulted to the token's subject, the
// token's instance permissions, and their resource restrictions as additional rules.
func agentApproverClaims(m drivers.ProjectMember) *SecurityClaims {
	attrs := maps.Clone(m.Attributes)
	if attrs == nil {
		attrs = make(map[string]any)
	}
	if _, ok := attrs["id"]; !ok {
		attrs["id"] = m.UserID
	}

	permissions := slices.Clone(agentApproverPermissions)
	if m.EditTrigger {
		permissions = append(permissions, EditTrigger)
	}

	return &SecurityClaims{
		UserID:          m.UserID,
		UserAttributes:  attrs,
		Permissions:     permissions,
		AdditionalRules: m.SecurityRules,
	}
}

// evaluateAgentGate resolves and evaluates one boolean gate expression against the caller's attributes,
// mirroring how the security engine evaluates access condition expressions (see evaluateConditions).
func (r *Runtime) evaluateAgentGate(ctx context.Context, instanceID string, claims *SecurityClaims, res *runtimev1.Resource, expr string, extra map[string]any) (bool, error) {
	inst, err := r.Instance(ctx, instanceID)
	if err != nil {
		return false, err
	}

	attrs := claims.UserAttributes
	if attrs == nil {
		attrs = make(map[string]any)
	}
	templateData := parser.TemplateData{
		Environment: inst.Environment,
		User:        attrs,
		Variables:   inst.ResolveVariables(false),
		ExtraProps:  extra,
		Self:        parser.TemplateResource{Meta: res.Meta},
		Resolve: func(ref parser.ResourceName) (string, error) {
			return ref.Name, nil
		},
	}

	resolved, err := parser.ResolveTemplate(expr, templateData, false)
	if err != nil {
		return false, err
	}
	return parser.EvaluateBoolExpression(resolved)
}

// agentSpecForGates returns the agent's validated spec, falling back to the raw spec when the resource has not
// reconciled yet (same fallback the security engine uses when resolving the agent's access rules).
func agentSpecForGates(res *runtimev1.Resource) *runtimev1.AgentSpec {
	spec := res.GetAgent().State.ValidSpec
	if spec == nil {
		spec = res.GetAgent().Spec
	}
	return spec
}
