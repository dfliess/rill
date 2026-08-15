package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"github.com/rilldata/rill/runtime/server/auth"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// errActNotConfigured is returned by the run/approval endpoints when the Act plane (store + executor) has not been
// wired via ConfigureAct. Discovery (ListAgents/GetAgent) does not need it, so those endpoints stay available.
var errActNotConfigured = status.Error(codes.Unimplemented, "act run store is not configured on this runtime")

// errActDisabled is returned by every Act endpoint when the "agents" feature flag is off. The flag is the kill
// switch (issue #135): most feature flags only gate the UI, but this one is enforced in the handlers so turning
// it off also closes the API, per project and per environment, without a redeploy.
var errActDisabled = status.Error(codes.FailedPrecondition, `the "agents" feature is disabled for this project`)

// runEventPollInterval is how often the event stream polls the store for new events once it has drained the backlog.
// The store is the source of truth (§7.2), so the stream is a poller over it rather than a pub/sub: simple, and it
// survives a worker on a different process than the API. A run's events are few and terminal-bounded, so this is cheap.
const runEventPollInterval = 500 * time.Millisecond

// ListAgents lists the valid agents defined in an instance's project that the caller has access to. Access is
// resolved per agent by the security engine (the agent's security rules plus the built-in admin rule), which is
// what keeps an agent's definition (its instructions) out of reach of callers it was not opened to, including
// public-link and embed tokens whose exclusive rules never include an agent.
func (s *Server) ListAgents(ctx context.Context, req *runtimev1.ListAgentsRequest) (*runtimev1.ListAgentsResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx, attribute.String("args.instance_id", req.InstanceId))

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}

	resources, err := s.listAgentResources(ctx, req.InstanceId)
	if err != nil {
		return nil, err
	}
	agents := make([]*runtimev1.AgentDefinition, 0, len(resources))
	for _, res := range resources {
		if res.GetAgent().State.ValidSpec == nil {
			// Not reconciled to a valid spec: not runnable, so treat it as absent (same as the executor's provider).
			continue
		}
		gates, err := s.runtime.ResolveAgentGates(ctx, req.InstanceId, claims, res)
		if err != nil {
			return nil, err
		}
		if !gates.Access {
			continue
		}
		agents = append(agents, agentToPB(res, gates))
	}
	return &runtimev1.ListAgentsResponse{Agents: agents}, nil
}

// GetAgent returns a single agent by name. An agent the caller has no access to is reported as not found, so
// the endpoint is not an existence oracle for hidden agents.
func (s *Server) GetAgent(ctx context.Context, req *runtimev1.GetAgentRequest) (*runtimev1.GetAgentResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.name", req.Name),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}

	res, gates, err := s.resolveAgentGates(ctx, req.InstanceId, req.Name, claims)
	if err != nil {
		return nil, err
	}
	if !gates.Access {
		return nil, status.Errorf(codes.NotFound, "agent %q not found", req.Name)
	}
	return &runtimev1.GetAgentResponse{Agent: agentToPB(res, gates)}, nil
}

// checkActEnabled enforces the "agents" feature flag as the Act kill switch. It is checked in every Act
// handler, not only in the UI: the flag's value can point at a project variable, so Act can be shut off per
// project and per environment without git or a redeploy, and shutting it off also closes the API.
func (s *Server) checkActEnabled(ctx context.Context, instanceID string) error {
	claims := auth.GetClaims(ctx, instanceID)
	if claims.SkipChecks {
		// Local development skips all access checks, and the kill switch with it.
		return nil
	}
	ff, err := s.runtime.FeatureFlags(ctx, instanceID, claims)
	if err != nil {
		return err
	}
	if !ff["agents"] {
		return errActDisabled
	}
	return nil
}

// listAgentResources returns the instance's Agent resources from the catalog.
func (s *Server) listAgentResources(ctx context.Context, instanceID string) ([]*runtimev1.Resource, error) {
	ctrl, err := s.runtime.Controller(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	return ctrl.List(ctx, runtime.ResourceKindAgent, "", false)
}

// resolveAgentGates fetches the named agent from the catalog and resolves the caller's gates on it. A missing
// agent and one that has not reconciled to a valid spec are both NotFound, mirroring the executor's provider.
func (s *Server) resolveAgentGates(ctx context.Context, instanceID, name string, claims *runtime.SecurityClaims) (*runtimev1.Resource, runtime.AgentGates, error) {
	ctrl, err := s.runtime.Controller(ctx, instanceID)
	if err != nil {
		return nil, runtime.AgentGates{}, err
	}
	res, err := ctrl.Get(ctx, &runtimev1.ResourceName{Kind: runtime.ResourceKindAgent, Name: name}, false)
	if err != nil {
		if errors.Is(err, drivers.ErrResourceNotFound) {
			return nil, runtime.AgentGates{}, status.Errorf(codes.NotFound, "agent %q not found", name)
		}
		return nil, runtime.AgentGates{}, err
	}
	if res.GetAgent().State.ValidSpec == nil {
		return nil, runtime.AgentGates{}, status.Errorf(codes.NotFound, "agent %q not found", name)
	}
	gates, err := s.runtime.ResolveAgentGates(ctx, instanceID, claims, res)
	if err != nil {
		return nil, runtime.AgentGates{}, err
	}
	return res, gates, nil
}

// accessibleAgentNames resolves the read-side allow-list for run and approval listings: the agents the caller
// has access to, evaluated once per agent (a project has few). Nil means unrestricted and is reserved for
// admins (and skipped checks), which also keeps the runs of a since-deleted agent visible to operators: they
// are audit trail. For everyone else the list is exact, and empty (non-nil) matches nothing.
func (s *Server) accessibleAgentNames(ctx context.Context, instanceID string, claims *runtime.SecurityClaims) ([]string, error) {
	// NOTE: The shortcut keys on EditTrigger, NOT on claims.Admin(). The "admin" attribute is copied from the
	// creator into a magic auth token's attributes (admin/server/magic_tokens.go, and upstream's own note
	// there), so a share link created by an admin reports Admin() == true. Permissions cannot be forged that
	// way: they are derived by the admin service from the token's own project permissions, and a magic auth
	// token never gets EditTrigger (ProjectPermissionsForMagicAuthToken sets ManageProd false). Keying on the
	// attribute here would hand such a link every run in the project, prompts and proposed actions included —
	// bypassing the security engine, which confines it correctly on the discovery path.
	if claims.SkipChecks || claims.Can(runtime.EditTrigger) {
		return nil, nil
	}
	resources, err := s.listAgentResources(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(resources))
	for _, res := range resources {
		gates, err := s.runtime.ResolveAgentGates(ctx, instanceID, claims, res)
		if err != nil {
			return nil, err
		}
		if gates.Access {
			names = append(names, res.Meta.Name.Name)
		}
	}
	return names, nil
}

// checkRunAccessible enforces access to the agent a stored run (or an approval addressed through one) belongs
// to. A run outside the caller's access is reported as not found, indistinguishable from a missing run. The
// run of a since-deleted agent stays reachable for admins only (see accessibleAgentNames).
func (s *Server) checkRunAccessible(ctx context.Context, instanceID, agentName string, claims *runtime.SecurityClaims) error {
	// EditTrigger, not claims.Admin(): see the note in accessibleAgentNames. A magic auth token inherits the
	// creator's "admin" attribute but never the permission, and this is the by-id read path, so keying on the
	// attribute would let a share link fetch any run it can name.
	if claims.SkipChecks || claims.Can(runtime.EditTrigger) {
		return nil
	}
	_, gates, err := s.resolveAgentGates(ctx, instanceID, agentName, claims)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return status.Error(codes.NotFound, act.ErrRunNotFound.Error())
		}
		return err
	}
	if !gates.Access {
		return status.Error(codes.NotFound, act.ErrRunNotFound.Error())
	}
	return nil
}

// StartAgentRun enqueues a durable run of an agent and returns its run id and initial (queued) status. It is
// gated by the agent's launch policy (issue #135): the `launch:` expression when the agent declares one, and
// otherwise the agent's resolved access, since starting a run is inert until an action passes the approval
// gate. It no longer requires EditTrigger, so launching stops implying administering the project.
func (s *Server) StartAgentRun(ctx context.Context, req *runtimev1.StartAgentRunRequest) (*runtimev1.StartAgentRunResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.name", req.Name),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "agent name is required")
	}

	// Resolve the launch gate. No access reads as not found (the agent must stay invisible); access without
	// launch is a plain permission error.
	_, gates, err := s.resolveAgentGates(ctx, req.InstanceId, req.Name, claims)
	if err != nil {
		return nil, err
	}
	if !gates.Access {
		return nil, status.Errorf(codes.NotFound, "agent %q not found", req.Name)
	}
	if !gates.Launch {
		return nil, ErrForbidden
	}

	// A manual run without a caller-supplied key gets a fresh one: the run still starts, it just is not deduplicated.
	// A caller that wants OAOO (e.g. a form submit guard) supplies a stable idempotency_key, which is scoped to the
	// user it acts for before it can deduplicate anything.
	idempotencyKey := req.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	} else {
		idempotencyKey = scopeKeyToActor(claims.UserID, idempotencyKey)
	}

	// This is the public manual-trigger API, so the run is always recorded as "manual": we ignore any caller-
	// supplied req.Trigger. Letting a caller set the trigger would let them forge provenance — labelling a run they
	// initiated as an automatic "alert"/"report", or the reverse — which the audit trail (§18.3) must not permit.
	// An automatic run's truthful trigger is set by the dispatcher, not here.
	runID, err := s.agentExecutor.Start(ctx, act.AgentRunInput{
		InstanceID:     req.InstanceId,
		AgentName:      req.Name,
		Prompt:         promptWithDashboardContext(req.Prompt, req.DashboardContext),
		IdempotencyKey: idempotencyKey,
		Trigger:        "manual",
		ConversationID: req.ConversationId,
		// The run executes as the calling user: it carries their claims, so it can read and act on exactly what
		// they could directly, never more (§17.3).
		Actor: act.Actor{Subject: claims.UserID, Claims: claims},
	})
	if err != nil {
		return nil, err
	}

	// Read the run back so the response reflects the stored state (status, spec_hash) rather than assuming queued.
	run, err := s.agentRuns.GetRun(ctx, req.InstanceId, runID)
	if err != nil {
		return nil, mapActError(err)
	}
	return &runtimev1.StartAgentRunResponse{
		RunId:     run.RunID,
		Status:    string(run.Status),
		AgentName: run.AgentName,
		SpecHash:  run.SpecHash,
	}, nil
}

// scopeKeyToActor namespaces a caller-supplied idempotency key by the user the run acts for.
//
// A run executes with its caller's claims (§17.3), so the same agent, prompt and dashboard state can yield
// materially different answers per user: an access policy narrows what each of them may read. Deduplicating across
// users would hand the second caller an answer computed under the first one's access, which the session owner check
// in runtime/ai then refuses to serve them, so they would get a permission error where an answer should be.
//
// An empty user id, an anonymous reader of a public project, namespaces to nothing, and those readers go on sharing
// one run. That is the same trade-off runtime/ai makes when it skips the owner check on a session nobody owns.
//
// Length-prefixed like ComposeRunID, and for the same reason: without the prefix ("a","b/c") and ("a/b","c") would
// render alike and collide onto one run.
func scopeKeyToActor(userID, idempotencyKey string) string {
	return fmt.Sprintf("%d:%s/%s", len(userID), userID, idempotencyKey)
}

// promptWithDashboardContext folds the calling surface's state into the run's prompt.
//
// A dynamic agent's whole input contract is its prompt, so the structured context has to become text somewhere; doing
// it here, once, keeps the API typed for callers and leaves the stored run self-contained, which is what lets a durable
// run resume without re-resolving a dashboard state that has since moved on.
//
// The rendering is deterministic (filters are emitted in sorted order): callers derive their idempotency key from the
// same context, so an unchanged dashboard must produce a byte-identical prompt or the key would promise reuse that the
// stored run does not deliver.
func promptWithDashboardContext(prompt string, dc *runtimev1.AnalystAgentContext) string {
	if dc == nil {
		return prompt
	}

	var b strings.Builder
	if dc.Canvas != "" {
		fmt.Fprintf(&b, "dashboard: %s\n", dc.Canvas)
	}
	if dc.Explore != "" {
		fmt.Fprintf(&b, "dashboard: %s\n", dc.Explore)
	}
	if dc.TimeStart != nil && dc.TimeEnd != nil {
		fmt.Fprintf(&b, "time range: %s to %s\n",
			dc.TimeStart.AsTime().Format(time.RFC3339),
			dc.TimeEnd.AsTime().Format(time.RFC3339))
	}
	if len(dc.Dimensions) > 0 {
		fmt.Fprintf(&b, "dimensions: %s\n", strings.Join(dc.Dimensions, ", "))
	}
	if len(dc.Measures) > 0 {
		fmt.Fprintf(&b, "measures: %s\n", strings.Join(dc.Measures, ", "))
	}
	if len(dc.WherePerMetricsView) > 0 {
		metricsViews := make([]string, 0, len(dc.WherePerMetricsView))
		for mv := range dc.WherePerMetricsView {
			metricsViews = append(metricsViews, mv)
		}
		sort.Strings(metricsViews)

		b.WriteString("filters in force (JSON, one per metrics view):\n")
		for _, mv := range metricsViews {
			expr, err := protojson.Marshal(dc.WherePerMetricsView[mv])
			if err != nil {
				continue // A filter we cannot render is dropped rather than failing the run: the note degrades, it does not break.
			}
			fmt.Fprintf(&b, "  %s: %s\n", mv, expr)
		}
	}

	if b.Len() == 0 {
		return prompt
	}
	return fmt.Sprintf(
		"The user is looking at this dashboard state. Scope your answer to it and say so if a filter changes what you would otherwise report.\n\n%s\n%s",
		b.String(), prompt,
	)
}

// ListAgentRuns lists an instance's runs, newest first, optionally filtered by agent and status. The caller
// only sees runs of agents they have access to; the restriction is pushed into the store's WHERE clause so a
// page is full of visible rows rather than filtered after the fact (which would make pagination lie).
func (s *Server) ListAgentRuns(ctx context.Context, req *runtimev1.ListAgentRunsRequest) (*runtimev1.ListAgentRunsResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx, attribute.String("args.instance_id", req.InstanceId))

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	accessible, err := s.accessibleAgentNames(ctx, req.InstanceId, claims)
	if err != nil {
		return nil, err
	}
	runs, err := s.agentRuns.ListRuns(ctx, act.ListRunsFilter{
		InstanceID:       req.InstanceId,
		AgentName:        req.AgentName,
		Status:           act.RunStatus(req.Status),
		AccessibleAgents: accessible,
		Limit:            int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	pbs := make([]*runtimev1.AgentRun, len(runs))
	for i, r := range runs {
		pbs[i] = runToPB(r, canCancelRun(claims, r))
	}
	return &runtimev1.ListAgentRunsResponse{Runs: pbs}, nil
}

// GetAgentRun returns a single run by id.
func (s *Server) GetAgentRun(ctx context.Context, req *runtimev1.GetAgentRunRequest) (*runtimev1.GetAgentRunResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.run_id", req.RunId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	run, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
	if err != nil {
		return nil, mapActError(err)
	}
	// A run of an agent outside the caller's access reads as not found, same as a missing run.
	if err := s.checkRunAccessible(ctx, req.InstanceId, run.AgentName, claims); err != nil {
		return nil, err
	}
	return &runtimev1.GetAgentRunResponse{Run: runToPB(run, canCancelRun(claims, run))}, nil
}

// CancelAgentRun stops a run and returns its updated state. It is allowed to the run's actor (whoever may
// launch may stop what they launched) and to EditTrigger holders (the operator's gate); see canCancelRun.
func (s *Server) CancelAgentRun(ctx context.Context, req *runtimev1.CancelAgentRunRequest) (*runtimev1.CancelAgentRunResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.run_id", req.RunId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	// Validate existence and scoping before cancelling, so a cancel cannot touch another tenant's run.
	run, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
	if err != nil {
		return nil, mapActError(err)
	}
	// Cancelling is gated on being the run's actor (whoever may launch may stop what they launched) or on
	// EditTrigger (the operator's gate). A run outside the caller's access stays not-found rather than
	// forbidden, so cancel is not an existence oracle either.
	if !canCancelRun(claims, run) {
		if err := s.checkRunAccessible(ctx, req.InstanceId, run.AgentName, claims); err != nil {
			return nil, err
		}
		return nil, ErrForbidden
	}
	// Reject cancelling a run that has already finished: a terminal run has no work to stop, and overwriting its
	// outcome (e.g. turning succeeded into cancelled) would falsify the record. This matches DBOS, which treats
	// cancelling a completed workflow as a no-op.
	if run.Status.IsTerminal() {
		return nil, status.Errorf(codes.FailedPrecondition, "run %q is already %s and cannot be cancelled", req.RunId, run.Status)
	}
	if err := s.agentExecutor.Cancel(ctx, req.InstanceId, req.RunId); err != nil {
		return nil, err
	}
	run, err = s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
	if err != nil {
		return nil, mapActError(err)
	}
	return &runtimev1.CancelAgentRunResponse{Run: runToPB(run, canCancelRun(claims, run))}, nil
}

// ListAgentApprovals lists approval requests, newest first, optionally filtered by run and status. Like runs,
// the caller only sees approvals of agents they have access to, restricted in the store's WHERE clause.
func (s *Server) ListAgentApprovals(ctx context.Context, req *runtimev1.ListAgentApprovalsRequest) (*runtimev1.ListAgentApprovalsResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx, attribute.String("args.instance_id", req.InstanceId))

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	accessible, err := s.accessibleAgentNames(ctx, req.InstanceId, claims)
	if err != nil {
		return nil, err
	}
	approvals, err := s.agentRuns.ListApprovals(ctx, act.ListApprovalsFilter{
		InstanceID:       req.InstanceId,
		RunID:            req.RunId,
		Status:           req.Status,
		AccessibleAgents: accessible,
		Limit:            int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	pbs := make([]*runtimev1.AgentApproval, len(approvals))
	for i, a := range approvals {
		// can_decide is resolved per approval, with its concrete action bound: approval authority can
		// discriminate by tool, so one caller may decide some of an agent's approvals and not others.
		canDecide, _, _, err := s.resolveApprovalDecision(ctx, req.InstanceId, a, claims)
		if err != nil {
			return nil, err
		}
		pbs[i] = approvalToPB(a, canDecide)
	}
	return &runtimev1.ListAgentApprovalsResponse{Approvals: pbs}, nil
}

// GetAgentApproval returns a single approval by id.
func (s *Server) GetAgentApproval(ctx context.Context, req *runtimev1.GetAgentApprovalRequest) (*runtimev1.GetAgentApprovalResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.approval_id", req.ApprovalId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approval, err := s.agentRuns.GetApproval(ctx, req.InstanceId, req.ApprovalId)
	if err != nil {
		return nil, mapActError(err)
	}
	// An approval belongs to its run's agent; outside the caller's access it reads as not found.
	if err := s.checkRunAccessible(ctx, req.InstanceId, approval.AgentName, claims); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, act.ErrApprovalNotFound.Error())
		}
		return nil, err
	}
	canDecide, _, _, err := s.resolveApprovalDecision(ctx, req.InstanceId, approval, claims)
	if err != nil {
		return nil, err
	}
	return &runtimev1.GetAgentApprovalResponse{Approval: approvalToPB(approval, canDecide)}, nil
}

// resolveApprovalDecision evaluates whether the caller may decide the given approval, with its concrete
// action bound (`.action.tool`, `.action.connector`, `.run.actor` all come from the approval's stored row).
// It reports three things: whether deciding is allowed, whether that verdict rests on the EditTrigger
// break-glass (allowed despite a failing approve policy), and whether the approval is visible to the caller
// at all (deciding is a subset of access, so an approval outside access must read as not found).
//
// It is both the enforcement input (authorizeApprovalDecision) and the source of the API's can_decide bit:
// authority can discriminate by action, so it is a property of one approval, never of the agent.
func (s *Server) resolveApprovalDecision(ctx context.Context, instanceID string, approval *act.Approval, claims *runtime.SecurityClaims) (allowed, breakGlass, visible bool, err error) {
	if claims.SkipChecks {
		return true, false, true, nil
	}

	res, gates, err := s.resolveAgentGates(ctx, instanceID, approval.AgentName, claims)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// The agent was deleted (or lost its valid spec) after the run started, and its policy went with
			// it: fall back to the operator. For anyone else the approval reads as not found.
			// Keyed on EditTrigger and not on claims.Admin() for the same reason as accessibleAgentNames: the
			// "admin" attribute is inherited by a share link, the permission is not. Without the agent there
			// is no policy left to confine the decision, so the weaker check must not be the one that runs.
			ok := claims.Can(runtime.EditTrigger)
			return ok, false, ok, nil
		}
		return false, false, false, err
	}
	if !gates.Access {
		return false, false, false, nil
	}

	allowed, err = s.runtime.ResolveAgentApprove(ctx, instanceID, claims, res, runtime.AgentActionContext{
		// The RAW tool name, not the mcp.<connector>.<tool> form the ledger stores. An approve expression is written
		// by the same author who wrote the connector's auto_approve globs a few lines above it in the same file, and
		// those match the raw name; binding the prefixed form here would mean two naming conventions in one YAML.
		// It also fails in the dangerous direction: `eq .action.tool "delete_account"` would silently never match, so
		// a rule meant to reserve deletions for one group would hand them to everyone the fallback clause allows.
		Tool:      act.RawToolName(approval.ToolName, approval.Connector),
		Connector: approval.Connector,
		RunActor:  approval.RunActorSubject,
	})
	if err != nil {
		return false, false, true, err
	}
	if allowed {
		return true, false, true, nil
	}
	if claims.Can(runtime.EditTrigger) {
		return true, true, true, nil
	}
	return false, false, true, nil
}

// authorizeApprovalDecision gates an approval decision on the agent's approve policy, with EditTrigger as the
// break-glass: an operator can always decide, but overriding a declared approve policy is exceptional, so
// that path is logged for audit.
func (s *Server) authorizeApprovalDecision(ctx context.Context, instanceID string, approval *act.Approval, claims *runtime.SecurityClaims) error {
	allowed, breakGlass, visible, err := s.resolveApprovalDecision(ctx, instanceID, approval, claims)
	if err != nil {
		return err
	}
	if !visible {
		return status.Error(codes.NotFound, act.ErrApprovalNotFound.Error())
	}
	if !allowed {
		return ErrForbidden
	}
	if breakGlass {
		s.logger.Warn("agent approval decided via EditTrigger break-glass, overriding the agent's approve policy",
			zap.String("instance_id", instanceID),
			zap.String("agent", approval.AgentName),
			zap.String("approval_id", approval.ApprovalID),
			zap.String("user_id", claims.UserID),
		)
		observability.AddRequestAttributes(ctx, attribute.Bool("act.approval_break_glass", true))
	}
	return nil
}

// ApproveAgentApproval approves a pending approval and resumes its run. The decision is gated by the agent's
// approve policy (issue #135), evaluated with the concrete action bound, with EditTrigger retained as the
// audited break-glass (see authorizeApprovalDecision). The request must echo the args_hash the approver saw,
// so a decision cannot silently apply to arguments that changed since (§6.4).
func (s *Server) ApproveAgentApproval(ctx context.Context, req *runtimev1.ApproveAgentApprovalRequest) (*runtimev1.ApproveAgentApprovalResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.approval_id", req.ApprovalId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approval, err := s.agentRuns.GetApproval(ctx, req.InstanceId, req.ApprovalId)
	if err != nil {
		return nil, mapActError(err)
	}
	if err := s.authorizeApprovalDecision(ctx, req.InstanceId, approval, claims); err != nil {
		return nil, err
	}
	// Bind the decision to the exact proposed arguments: the approver must confirm the hash they were shown.
	if req.ArgsHash == "" {
		return nil, status.Error(codes.InvalidArgument, "args_hash is required to approve")
	}
	if req.ArgsHash != approval.ArgsHash {
		return nil, status.Error(codes.FailedPrecondition, "args_hash does not match the proposed action; the proposal changed and must be re-reviewed")
	}

	resolved, err := s.claimAndResume(ctx, approval, act.ApprovalApproved, act.ApprovalStatusApproved, claims.UserID)
	if err != nil {
		return nil, err
	}
	// The caller just decided it, so their can_decide is true by construction.
	return &runtimev1.ApproveAgentApprovalResponse{Approval: approvalToPB(resolved, true)}, nil
}

// claimAndResume delivers an approval decision to its run durably, in an order that is both double-submit safe and
// crash-recoverable.
//
// The claim (ResolveApproval, a CAS on the pending row) is the single-winner gate: of two concurrent decisions only
// one moves the approval out of pending, so only one is ever delivered — two conflicting decisions can never both
// resume the run. The DBOS Send (Resume) is the durable source of truth for the run continuing: the workflow
// consumes exactly one approval message (a single Recv), so any Send after the first is ignored.
//
// Ordering: claim, then Send. Recording the decision before the Send is what keeps two different concurrent
// decisions from both being delivered — the loser's claim fails, so it never Sends. The one hazard that ordering
// creates is a crash after the claim but before the Send: the approval is no longer pending, so a naive retry
// could never deliver the resume and the run would hang forever. We close that gap by making a retry that finds
// the approval already resolved to the SAME decision re-deliver the Send (idempotent, thanks to the single Recv)
// instead of failing. A retry carrying a DIFFERENT decision is rejected: the first recorded decision stands.
// TODO(act, phase 1): recovery from a crash-between-claim-and-Send relies on an EXTERNAL retry arriving before the
// run's approval timeout; no reconciler completes a claimed-but-unsent decision on its own. A durable resume outbox
// (claim and enqueue the Send in one transaction, drained by a background worker) would make it self-healing. Same
// class as known limit #2 (the trigger outbox).
func (s *Server) claimAndResume(ctx context.Context, approval *act.Approval, decision act.ApprovalDecision, storeStatus, decidedBy string) (*act.Approval, error) {
	var resolved *act.Approval
	switch approval.Status {
	case act.ApprovalStatusPending:
		r, err := s.agentRuns.ResolveApproval(ctx, approval.InstanceID, approval.ApprovalID, storeStatus, decidedBy)
		if err != nil {
			return nil, mapActError(err)
		}
		resolved = r
	case storeStatus:
		// Recovery: a prior attempt already claimed this approval for the same decision but may have crashed before
		// delivering the resume. Fall through to re-deliver; the run's single Recv makes the repeat harmless.
		resolved = approval
	default:
		return nil, status.Errorf(codes.FailedPrecondition, "approval %q is already resolved and cannot be set to %s", approval.ApprovalID, storeStatus)
	}
	if err := s.agentExecutor.Resume(ctx, approval.RunID, decision, approval.ToolCallID); err != nil {
		return nil, err
	}
	return resolved, nil
}

// DenyAgentApproval denies a pending approval; its run ends without performing the proposed action. It is
// gated exactly like ApproveAgentApproval: denying is the other half of the same decision.
func (s *Server) DenyAgentApproval(ctx context.Context, req *runtimev1.DenyAgentApprovalRequest) (*runtimev1.DenyAgentApprovalResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.approval_id", req.ApprovalId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return nil, err
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approval, err := s.agentRuns.GetApproval(ctx, req.InstanceId, req.ApprovalId)
	if err != nil {
		return nil, mapActError(err)
	}
	if err := s.authorizeApprovalDecision(ctx, req.InstanceId, approval, claims); err != nil {
		return nil, err
	}
	resolved, err := s.claimAndResume(ctx, approval, act.ApprovalRejected, act.ApprovalStatusDenied, claims.UserID)
	if err != nil {
		return nil, err
	}
	// The caller just decided it, so their can_decide is true by construction.
	return &runtimev1.DenyAgentApprovalResponse{Approval: approvalToPB(resolved, true)}, nil
}

// StreamAgentRunEvents streams a run's lifecycle events in order, resuming after req.AfterId, until the run reaches a
// terminal state and the backlog is drained (or the client disconnects). It polls the store, which is the source of
// truth, so it works whether the worker runs in this process or another.
func (s *Server) StreamAgentRunEvents(req *runtimev1.StreamAgentRunEventsRequest, stream runtimev1.AgentService_StreamAgentRunEventsServer) error {
	ctx := runtime.WithRequestSource(stream.Context(), runtime.RequestSourceAct)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.run_id", req.RunId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.UseAI) {
		return ErrForbidden
	}
	if err := s.checkActEnabled(ctx, req.InstanceId); err != nil {
		return err
	}
	if s.agentRuns == nil {
		return errActNotConfigured
	}

	// Validate existence, scoping and access up front, so a stream on another tenant's, a missing, or a hidden
	// run fails cleanly before any event is sent.
	run, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
	if err != nil {
		return mapActError(err)
	}
	if err := s.checkRunAccessible(ctx, req.InstanceId, run.AgentName, claims); err != nil {
		return err
	}

	const pageLimit = 200
	cursor := req.AfterId
	for {
		events, err := s.agentRuns.ListRunEvents(ctx, req.InstanceId, req.RunId, cursor, pageLimit)
		if err != nil {
			return err
		}
		for _, e := range events {
			if err := stream.Send(&runtimev1.StreamAgentRunEventsResponse{Event: runEventToPB(e)}); err != nil {
				return err
			}
			cursor = e.ID
		}
		// A full page may have more behind it: page again immediately before deciding we are caught up.
		if len(events) == pageLimit {
			continue
		}

		// Caught up. If the run is terminal, no more events will ever arrive, so end the stream.
		run, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
		if err != nil {
			return mapActError(err)
		}
		if run.Status.IsTerminal() {
			// The terminal event may have been appended between the last page and this status read (the status
			// update and its event are one transaction, so the event is committed once the status is terminal).
			// Drain from the cursor one last time before closing, or that final event would be lost to the race.
			for {
				tail, err := s.agentRuns.ListRunEvents(ctx, req.InstanceId, req.RunId, cursor, pageLimit)
				if err != nil {
					return err
				}
				for _, e := range tail {
					if err := stream.Send(&runtimev1.StreamAgentRunEventsResponse{Event: runEventToPB(e)}); err != nil {
						return err
					}
					cursor = e.ID
				}
				if len(tail) < pageLimit {
					return nil
				}
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(runEventPollInterval):
		}
	}
}

// StreamAgentRunEventsHandler is the SSE HTTP handler for StreamAgentRunEvents. Vanguard does not map streaming RPCs
// to SSE, so (as with chat) this shim adapts the gRPC streaming method to an SSE response (§chat.go pattern).
func (s *Server) StreamAgentRunEventsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	instanceID := req.PathValue("instance_id")
	runID := req.PathValue("run_id")
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", instanceID),
		attribute.String("args.run_id", runID),
	)

	if !auth.GetClaims(ctx, instanceID).Can(runtime.UseAI) {
		http.Error(w, "action not allowed", http.StatusUnauthorized)
		return
	}

	var afterID int64
	if v := req.URL.Query().Get("after_id"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			http.Error(w, "invalid 'after_id' parameter", http.StatusBadRequest)
			return
		}
		afterID = parsed
	}

	streamReq := &runtimev1.StreamAgentRunEventsRequest{InstanceId: instanceID, RunId: runID, AfterId: afterID}

	events := make(chan *sseEvent)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("panic in StreamAgentRunEventsHandler goroutine", zap.Any("recover", r), zap.Stack("stack"))
			}
		}()
		defer close(events)

		shim := &grpcStreamingShim[*runtimev1.StreamAgentRunEventsResponse]{
			ctx: ctx,
			fn: func(data []byte) error {
				events <- &sseEvent{Event: "agent_run_event", Data: data}
				return nil
			},
		}
		err := s.StreamAgentRunEvents(streamReq, shim)
		if err != nil && !errors.Is(err, context.Canceled) {
			code := codes.Unknown
			msg := err.Error()
			if st, ok := status.FromError(err); ok {
				code = st.Code()
				msg = st.Message()
			}
			errJSON, mErr := json.Marshal(map[string]string{"code": code.String(), "error": msg})
			if mErr != nil {
				s.logger.Error("failed to marshal error as json", zap.Error(mErr))
			}
			events <- &sseEvent{Event: "error", Data: errJSON}
		}
	}()

	serveSSEUntilClose(w, events)
}

// mapActError maps the act store's sentinel errors to gRPC status codes. Other errors pass through to the generic
// error-mapping interceptor.
func mapActError(err error) error {
	switch {
	case errors.Is(err, act.ErrRunNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, act.ErrApprovalNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, act.ErrApprovalNotResolvable):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return err
	}
}

// agentToPB projects a reconciled agent resource (its valid spec) plus the caller's resolved gates onto the
// API's AgentDefinition. The launch gate rides along so the UI can hide the launch affordances the caller
// cannot use; enforcement stays in the handlers. Approval authority is per approval (see AgentApproval's
// can_decide), so it deliberately does not appear here.
func agentToPB(res *runtimev1.Resource, gates runtime.AgentGates) *runtimev1.AgentDefinition {
	spec := res.GetAgent().State.ValidSpec
	pb := &runtimev1.AgentDefinition{
		Name:           res.Meta.Name.Name,
		DisplayName:    spec.DisplayName,
		Instructions:   spec.Instructions,
		ModelConnector: spec.ModelConnector,
		ModelName:      spec.ModelName,
		Tools:          spec.Tools,
		CanLaunch:      gates.Launch,
	}
	if spec.Limits != nil {
		pb.MaxSteps = int32(spec.Limits.MaxSteps)
		pb.TimeoutSeconds = int32(spec.Limits.TimeoutSeconds)
	}
	return pb
}

// canCancelRun reports whether the caller may cancel the run: its actor (whoever may launch may stop what
// they launched) or an EditTrigger holder. The empty-subject guard matters: an anonymous caller must not
// match a service run's empty actor.
func canCancelRun(claims *runtime.SecurityClaims, r *act.Run) bool {
	if claims.SkipChecks || claims.Can(runtime.EditTrigger) {
		return true
	}
	return claims.UserID != "" && claims.UserID == r.Actor.Subject
}

func runToPB(r *act.Run, canCancel bool) *runtimev1.AgentRun {
	return &runtimev1.AgentRun{
		RunId:                 r.RunID,
		InstanceId:            r.InstanceID,
		OrganizationId:        r.OrganizationID,
		ProjectId:             r.ProjectID,
		AgentName:             r.AgentName,
		SpecHash:              r.SpecHash,
		Trigger:               r.Trigger,
		TriggerRef:            r.TriggerRef,
		IdempotencyKey:        r.IdempotencyKey,
		ConversationId:        r.ConversationID,
		ActorSubject:          r.Actor.Subject,
		ActorServicePrincipal: r.Actor.ServicePrincipal,
		Status:                string(r.Status),
		Error:                 r.Error,
		CreatedOn:             timestamppb.New(r.CreatedOn),
		UpdatedOn:             timestamppb.New(r.UpdatedOn),
		StartedOn:             tsToPB(r.StartedOn),
		FinishedOn:            tsToPB(r.FinishedOn),
		CanCancel:             canCancel,
	}
}

func runEventToPB(e *act.RunEvent) *runtimev1.AgentRunEvent {
	pb := &runtimev1.AgentRunEvent{
		Id:         e.ID,
		RunId:      e.RunID,
		Seq:        e.Seq,
		EventType:  e.EventType,
		Status:     string(e.Status),
		Visibility: e.Visibility,
		CreatedOn:  timestamppb.New(e.CreatedOn),
	}
	if e.Payload != nil {
		if payload, err := structpb.NewStruct(e.Payload); err == nil {
			pb.Payload = payload
		}
	}
	return pb
}

func approvalToPB(a *act.Approval, canDecide bool) *runtimev1.AgentApproval {
	return &runtimev1.AgentApproval{
		ApprovalId:  a.ApprovalID,
		RunId:       a.RunID,
		InstanceId:  a.InstanceID,
		ToolName:    a.ToolName,
		Connector:   a.Connector,
		ToolCallId:  a.ToolCallID,
		ArgsHash:    a.ArgsHash,
		Proposal:    a.Proposal,
		Policy:      a.Policy,
		Status:      a.Status,
		RequestedBy: a.RequestedBy,
		DecidedBy:   a.DecidedBy,
		CreatedOn:   timestamppb.New(a.CreatedOn),
		DecidedOn:   tsToPB(a.DecidedOn),
		Position:    int32(a.Position),
		Total:       int32(a.Total),
		CanDecide:   canDecide,
	}
}

// tsToPB converts a nullable time to a protobuf timestamp, mapping nil to nil so an unset time stays unset on the wire.
func tsToPB(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
