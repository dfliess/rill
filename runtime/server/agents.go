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
	"github.com/rilldata/rill/runtime/ai"
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

// runEventPollInterval is how often the event stream polls the store for new events once it has drained the backlog.
// The store is the source of truth (§7.2), so the stream is a poller over it rather than a pub/sub: simple, and it
// survives a worker on a different process than the API. A run's events are few and terminal-bounded, so this is cheap.
const runEventPollInterval = 500 * time.Millisecond

// ListAgents lists the valid agents defined in an instance's project.
func (s *Server) ListAgents(ctx context.Context, req *runtimev1.ListAgentsRequest) (*runtimev1.ListAgentsResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx, attribute.String("args.instance_id", req.InstanceId))

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return nil, ErrForbidden
	}

	snapshots, err := s.agents.ListAgents(ctx, req.InstanceId)
	if err != nil {
		return nil, err
	}
	agents := make([]*runtimev1.AgentDefinition, len(snapshots))
	for i, snap := range snapshots {
		agents[i] = agentDefinitionToPB(snap)
	}
	return &runtimev1.ListAgentsResponse{Agents: agents}, nil
}

// GetAgent returns a single agent by name.
func (s *Server) GetAgent(ctx context.Context, req *runtimev1.GetAgentRequest) (*runtimev1.GetAgentResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.name", req.Name),
	)

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return nil, ErrForbidden
	}

	snap, err := s.agents.GetAgent(ctx, req.InstanceId, req.Name)
	if err != nil {
		if errors.Is(err, ai.ErrAgentNotFound) {
			return nil, status.Errorf(codes.NotFound, "agent %q not found", req.Name)
		}
		return nil, err
	}
	return &runtimev1.GetAgentResponse{Agent: agentDefinitionToPB(snap)}, nil
}

// StartAgentRun enqueues a durable run of an agent and returns its run id and initial (queued) status. It is the
// write side of manual triggering, so it requires EditTrigger: the same permission that gates a manual refresh or
// reconcile trigger, and the one a project manager holds even on a prod (non-editable) deployment. It stands in for
// the security.execute claim (§8.2) until the policy engine lands.
func (s *Server) StartAgentRun(ctx context.Context, req *runtimev1.StartAgentRunRequest) (*runtimev1.StartAgentRunResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.name", req.Name),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.EditTrigger) {
		return nil, ErrForbidden
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "agent name is required")
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

// ListAgentRuns lists an instance's runs, newest first, optionally filtered by agent and status.
func (s *Server) ListAgentRuns(ctx context.Context, req *runtimev1.ListAgentRunsRequest) (*runtimev1.ListAgentRunsResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx, attribute.String("args.instance_id", req.InstanceId))

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	runs, err := s.agentRuns.ListRuns(ctx, act.ListRunsFilter{
		InstanceID: req.InstanceId,
		AgentName:  req.AgentName,
		Status:     act.RunStatus(req.Status),
		Limit:      int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	pbs := make([]*runtimev1.AgentRun, len(runs))
	for i, r := range runs {
		pbs[i] = runToPB(r)
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

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	run, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
	if err != nil {
		return nil, mapActError(err)
	}
	return &runtimev1.GetAgentRunResponse{Run: runToPB(run)}, nil
}

// CancelAgentRun stops a run and returns its updated state.
func (s *Server) CancelAgentRun(ctx context.Context, req *runtimev1.CancelAgentRunRequest) (*runtimev1.CancelAgentRunResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.run_id", req.RunId),
	)

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.EditTrigger) {
		return nil, ErrForbidden
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	// Validate existence and scoping before cancelling, so a cancel cannot touch another tenant's run.
	run, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId)
	if err != nil {
		return nil, mapActError(err)
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
	return &runtimev1.CancelAgentRunResponse{Run: runToPB(run)}, nil
}

// ListAgentApprovals lists approval requests, newest first, optionally filtered by run and status.
func (s *Server) ListAgentApprovals(ctx context.Context, req *runtimev1.ListAgentApprovalsRequest) (*runtimev1.ListAgentApprovalsResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx, attribute.String("args.instance_id", req.InstanceId))

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approvals, err := s.agentRuns.ListApprovals(ctx, act.ListApprovalsFilter{
		InstanceID: req.InstanceId,
		RunID:      req.RunId,
		Status:     req.Status,
		Limit:      int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	pbs := make([]*runtimev1.AgentApproval, len(approvals))
	for i, a := range approvals {
		pbs[i] = approvalToPB(a)
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

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return nil, ErrForbidden
	}
	if s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approval, err := s.agentRuns.GetApproval(ctx, req.InstanceId, req.ApprovalId)
	if err != nil {
		return nil, mapActError(err)
	}
	return &runtimev1.GetAgentApprovalResponse{Approval: approvalToPB(approval)}, nil
}

// ApproveAgentApproval approves a pending approval and resumes its run. It requires EditTrigger: the write side of
// manual triggering, standing in for the security.execute claim in v1 (§2 alignment). The request must echo the
// args_hash the approver saw, so a decision cannot silently apply to arguments that changed since (§6.4).
func (s *Server) ApproveAgentApproval(ctx context.Context, req *runtimev1.ApproveAgentApprovalRequest) (*runtimev1.ApproveAgentApprovalResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.approval_id", req.ApprovalId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.EditTrigger) {
		return nil, ErrForbidden
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approval, err := s.agentRuns.GetApproval(ctx, req.InstanceId, req.ApprovalId)
	if err != nil {
		return nil, mapActError(err)
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
	return &runtimev1.ApproveAgentApprovalResponse{Approval: approvalToPB(resolved)}, nil
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

// DenyAgentApproval denies a pending approval; its run ends without performing the proposed action.
func (s *Server) DenyAgentApproval(ctx context.Context, req *runtimev1.DenyAgentApprovalRequest) (*runtimev1.DenyAgentApprovalResponse, error) {
	ctx = runtime.WithRequestSource(ctx, runtime.RequestSourceAct)
	s.addInstanceRequestAttributes(ctx, req.InstanceId)
	observability.AddRequestAttributes(ctx,
		attribute.String("args.instance_id", req.InstanceId),
		attribute.String("args.approval_id", req.ApprovalId),
	)

	claims := auth.GetClaims(ctx, req.InstanceId)
	if !claims.Can(runtime.EditTrigger) {
		return nil, ErrForbidden
	}
	if s.agentExecutor == nil || s.agentRuns == nil {
		return nil, errActNotConfigured
	}

	approval, err := s.agentRuns.GetApproval(ctx, req.InstanceId, req.ApprovalId)
	if err != nil {
		return nil, mapActError(err)
	}
	resolved, err := s.claimAndResume(ctx, approval, act.ApprovalRejected, act.ApprovalStatusDenied, claims.UserID)
	if err != nil {
		return nil, err
	}
	return &runtimev1.DenyAgentApprovalResponse{Approval: approvalToPB(resolved)}, nil
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

	if !auth.GetClaims(ctx, req.InstanceId).Can(runtime.UseAI) {
		return ErrForbidden
	}
	if s.agentRuns == nil {
		return errActNotConfigured
	}

	// Validate existence and scoping up front, so a stream on another tenant's (or a missing) run fails cleanly.
	if _, err := s.agentRuns.GetRun(ctx, req.InstanceId, req.RunId); err != nil {
		return mapActError(err)
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

func agentDefinitionToPB(snap *ai.AgentSnapshot) *runtimev1.AgentDefinition {
	return &runtimev1.AgentDefinition{
		Name:           snap.Name,
		DisplayName:    snap.DisplayName,
		Instructions:   snap.Instructions,
		ModelConnector: snap.ModelConnector,
		ModelName:      snap.ModelName,
		Tools:          snap.Tools,
		MaxSteps:       int32(snap.MaxSteps),
		TimeoutSeconds: int32(snap.TimeoutSeconds),
	}
}

func runToPB(r *act.Run) *runtimev1.AgentRun {
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

func approvalToPB(a *act.Approval) *runtimev1.AgentApproval {
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
	}
}

// tsToPB converts a nullable time to a protobuf timestamp, mapping nil to nil so an unset time stays unset on the wire.
func tsToPB(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
