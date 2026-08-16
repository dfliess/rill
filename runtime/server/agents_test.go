package server_test

import (
	"context"
	"net"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/ratelimit"
	"github.com/rilldata/rill/runtime/server"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// These tests exercise the AgentService handlers, their access checks, tenant scoping, approval flow and the SSE
// stream, against the real run store with a fake executor. The real executor <-> store integration (event emission,
// idempotency, recovery) is covered in the runtime/act package; here the executor is a fake so a handler test does
// not need to spin up a DBOS worker, and can drive the store to whatever state a scenario requires.

// requireStorePostgres returns a DSN for the Act store's Postgres, skipping the test if it is unreachable. It reuses
// the DBOS spike container (docker start dbos-spike-pg) unless overridden by env.
func requireStorePostgres(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("ACT_TEST_DBOS_URL")
	if dsn == "" {
		dsn = "postgres://postgres:dbos@localhost:55432/dbos_spike?sslmode=disable"
	}
	u, err := url.Parse(dsn)
	require.NoError(t, err)
	conn, err := net.DialTimeout("tcp", u.Host, 2*time.Second)
	if err != nil {
		t.Skipf("act server tests need Postgres at %s (start it with: docker start dbos-spike-pg): %v", u.Host, err)
	}
	_ = conn.Close()
	return dsn
}

// newActStore returns a migrated store on a schema unique to this test, dropped on cleanup.
func newActStore(t *testing.T) *act.PostgresRunStore {
	t.Helper()
	dsn := requireStorePostgres(t)
	schema := "act_srv_" + uuid.New().String()[:8]
	store, err := act.NewPostgresRunStore(t.Context(), act.StoreConfig{DatabaseURL: dsn, Schema: schema})
	require.NoError(t, err)
	require.NoError(t, store.Migrate(t.Context()))
	t.Cleanup(func() { store.DropSchemaForTest(t.Context()); store.Close() })
	return store
}

// resumeCall records one Resume invocation the handler made on the executor.
type resumeCall struct {
	runID    string
	decision act.ApprovalDecision
}

// fakeAgentExecutor is a stand-in AgentExecutor for handler tests. Start records the run as queued in the store (as
// the real executor does synchronously), so the run is immediately queryable; Resume and Cancel record their calls so
// a test can assert the handler drove the executor, and Cancel reflects the cancellation into the store.
type fakeAgentExecutor struct {
	store act.RunStore

	mu        sync.Mutex
	resumes   []resumeCall
	cancels   []string
	startErr  error
	resumeErr error
}

var _ act.AgentExecutor = (*fakeAgentExecutor)(nil)

func (f *fakeAgentExecutor) Start(ctx context.Context, in act.AgentRunInput) (string, error) {
	if f.startErr != nil {
		return "", f.startErr
	}
	runID := act.ComposeRunID(in.InstanceID, in.AgentName, in.IdempotencyKey)
	trigger := in.Trigger
	if trigger == "" {
		trigger = "manual"
	}
	if err := f.store.CreateRun(ctx, act.NewRun{
		RunID:          runID,
		InstanceID:     in.InstanceID,
		AgentName:      in.AgentName,
		Trigger:        trigger,
		IdempotencyKey: in.IdempotencyKey,
		ConversationID: in.ConversationID,
		Actor:          in.Actor,
	}); err != nil {
		return "", err
	}
	return runID, nil
}

func (f *fakeAgentExecutor) Resume(_ context.Context, runID string, decision act.ApprovalDecision, _ string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.resumeErr != nil {
		return f.resumeErr
	}
	f.resumes = append(f.resumes, resumeCall{runID: runID, decision: decision})
	return nil
}

func (f *fakeAgentExecutor) Cancel(ctx context.Context, instanceID, runID string) error {
	f.mu.Lock()
	f.cancels = append(f.cancels, runID)
	f.mu.Unlock()
	return f.store.RecordTransition(ctx, act.RunTransition{
		InstanceID: instanceID, RunID: runID, Status: act.RunStatusCancelled, EventType: act.EventTypeCancelled,
	})
}

func (f *fakeAgentExecutor) resumeCalls() []resumeCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]resumeCall(nil), f.resumes...)
}

// agentProjectFiles is a minimal project with one agent named "triage". The agents feature flag is on: it is
// the Act kill switch and the handlers enforce it, so a project without it has no Act API at all.
func agentProjectFiles() map[string]string {
	return map[string]string{
		"rill.yaml": `
ai_instructions: "Menciona la residencia de datos cuando sea relevante."
features:
  agents: true
`,
		"triage.yaml": `
type: agent
display_name: Ticket Triage
instructions: "Investiga la alerta y propon un ticket."
tools: []
limits:
  max_steps: 3
`,
	}
}

// newActServer builds a server over an instance that has one reconciled agent, with the Act plane wired to a fresh
// store and a fake executor.
func newActServer(t *testing.T) (*server.Server, *act.PostgresRunStore, *fakeAgentExecutor, string) {
	t.Helper()
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: agentProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)

	store := newActStore(t)
	exec := &fakeAgentExecutor{store: store}
	srv.ConfigureAct(store, exec)
	return srv, store, exec, instanceID
}

// TestAgentServiceDiscovery covers ListAgents and GetAgent over the reconciled catalog.
func TestAgentServiceDiscovery(t *testing.T) {
	srv, _, _, instanceID := newActServer(t)
	ctx := testCtx()

	list, err := srv.ListAgents(ctx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, list.Agents, 1)
	require.Equal(t, "triage", list.Agents[0].Name)
	require.Equal(t, "Ticket Triage", list.Agents[0].DisplayName)

	got, err := srv.GetAgent(ctx, &runtimev1.GetAgentRequest{InstanceId: instanceID, Name: "triage"})
	require.NoError(t, err)
	require.Equal(t, "triage", got.Agent.Name)

	_, err = srv.GetAgent(ctx, &runtimev1.GetAgentRequest{InstanceId: instanceID, Name: "nope"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

// TestAgentServiceStartRunQueryable verifies StartAgentRun returns a queued run that is immediately queryable via
// GetAgentRun, and that a repeated start with the same idempotency key attaches to the same run.
func TestAgentServiceStartRunQueryable(t *testing.T) {
	srv, _, _, instanceID := newActServer(t)
	ctx := testCtx()

	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId:     instanceID,
		Name:           "triage",
		IdempotencyKey: "key-1",
		Prompt:         "El coste subio.",
	})
	require.NoError(t, err)
	require.NotEmpty(t, start.RunId)
	require.Equal(t, string(act.RunStatusQueued), start.Status)
	require.Equal(t, "triage", start.AgentName)

	got, err := srv.GetAgentRun(ctx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.NoError(t, err)
	require.Equal(t, start.RunId, got.Run.RunId)
	require.Equal(t, string(act.RunStatusQueued), got.Run.Status)
	require.Equal(t, "manual", got.Run.Trigger)

	// Same idempotency key: same run id, still one run.
	start2, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId: instanceID, Name: "triage", IdempotencyKey: "key-1", Prompt: "otra vez",
	})
	require.NoError(t, err)
	require.Equal(t, start.RunId, start2.RunId)

	list, err := srv.ListAgentRuns(ctx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, list.Runs, 1)
}

// TestAgentServiceListRunsScoping verifies ListAgentRuns and GetAgentRun never cross tenants: a run created under one
// instance is invisible to another.
func TestAgentServiceListRunsScoping(t *testing.T) {
	srv, store, _, instanceID := newActServer(t)
	ctx := testCtx()

	// One run in the reconciled instance via the handler.
	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k"})
	require.NoError(t, err)

	// One run in a different instance, seeded directly in the store.
	const otherInstance = "inst-other"
	otherRun := act.ComposeRunID(otherInstance, "triage", "k")
	require.NoError(t, store.CreateRun(ctx, act.NewRun{RunID: otherRun, InstanceID: otherInstance, AgentName: "triage"}))

	list, err := srv.ListAgentRuns(ctx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, list.Runs, 1)
	require.Equal(t, start.RunId, list.Runs[0].RunId)

	otherList, err := srv.ListAgentRuns(ctx, &runtimev1.ListAgentRunsRequest{InstanceId: otherInstance})
	require.NoError(t, err)
	require.Len(t, otherList.Runs, 1)
	require.Equal(t, otherRun, otherList.Runs[0].RunId)

	// Cross-tenant fetch by exact id is NotFound.
	_, err = srv.GetAgentRun(ctx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: otherRun})
	require.Equal(t, codes.NotFound, status.Code(err))
}

// TestAgentServiceApproveResumesRun drives a run to the approval gate, then approves it: the handler validates the
// args hash, resolves the approval, and resumes the run through the executor. It also covers the hash mismatch and
// double-submit guards.
func TestAgentServiceApproveResumesRun(t *testing.T) {
	srv, store, exec, instanceID := newActServer(t)
	ctx := testCtx()

	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k"})
	require.NoError(t, err)
	runID := start.RunId

	// Simulate the worker reaching the approval gate: advance the run and persist the approval.
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusWaitingApproval, EventType: act.EventTypeWaitingApproval}))
	proposal := "crear ticket P2"
	approvalID := act.ApprovalIDForRun(runID)
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: approvalID, RunID: runID, InstanceID: instanceID,
		ToolName: "act.propose_action", ArgsHash: act.HashArgs(proposal), CanonicalArgs: proposal, Proposal: proposal, RequestedBy: "user:alice",
	}))

	// It appears in the pending inbox.
	inbox, err := srv.ListAgentApprovals(ctx, &runtimev1.ListAgentApprovalsRequest{InstanceId: instanceID, Status: act.ApprovalStatusPending})
	require.NoError(t, err)
	require.Len(t, inbox.Approvals, 1)
	require.Equal(t, approvalID, inbox.Approvals[0].ApprovalId)

	// A wrong hash is rejected without resuming.
	_, err = srv.ApproveAgentApproval(ctx, &runtimev1.ApproveAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs("otra cosa")})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	require.Empty(t, exec.resumeCalls())

	// The correct hash approves and resumes.
	approve, err := srv.ApproveAgentApproval(ctx, &runtimev1.ApproveAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs(proposal)})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, approve.Approval.Status)

	calls := exec.resumeCalls()
	require.Len(t, calls, 1)
	require.Equal(t, runID, calls[0].runID)
	require.Equal(t, act.ApprovalApproved, calls[0].decision)

	// Re-approving with the same decision is idempotent, not an error: it re-delivers the resume (the run consumes
	// it exactly once), which is what makes a crash between the claim and the resume recoverable by retry.
	reapprove, err := srv.ApproveAgentApproval(ctx, &runtimev1.ApproveAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID, ArgsHash: act.HashArgs(proposal)})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusApproved, reapprove.Approval.Status)
	require.Len(t, exec.resumeCalls(), 2, "the retry re-delivered the resume")

	// A conflicting later decision is rejected: the first recorded decision stands.
	_, err = srv.DenyAgentApproval(ctx, &runtimev1.DenyAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

// TestAgentServiceDenyEndsRun verifies denying an approval resolves it denied and resumes the run with a rejection.
func TestAgentServiceDenyEndsRun(t *testing.T) {
	srv, store, exec, instanceID := newActServer(t)
	ctx := testCtx()

	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k"})
	require.NoError(t, err)
	runID := start.RunId
	approvalID := act.ApprovalIDForRun(runID)
	require.NoError(t, store.CreateApproval(ctx, act.NewApproval{
		ApprovalID: approvalID, RunID: runID, InstanceID: instanceID, ArgsHash: act.HashArgs("x"), CanonicalArgs: "x", RequestedBy: "user:alice",
	}))

	deny, err := srv.DenyAgentApproval(ctx, &runtimev1.DenyAgentApprovalRequest{InstanceId: instanceID, ApprovalId: approvalID})
	require.NoError(t, err)
	require.Equal(t, act.ApprovalStatusDenied, deny.Approval.Status)

	calls := exec.resumeCalls()
	require.Len(t, calls, 1)
	require.Equal(t, act.ApprovalRejected, calls[0].decision)
}

// TestAgentServiceCancelRun verifies CancelAgentRun stops the run and reflects the cancelled state.
func TestAgentServiceCancelRun(t *testing.T) {
	srv, _, exec, instanceID := newActServer(t)
	ctx := testCtx()

	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k"})
	require.NoError(t, err)

	cancel, err := srv.CancelAgentRun(ctx, &runtimev1.CancelAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.NoError(t, err)
	require.Equal(t, string(act.RunStatusCancelled), cancel.Run.Status)
	require.Equal(t, []string{start.RunId}, exec.cancels)

	// Cancelling an already-terminal run is rejected: the run is cancelled, so there is nothing to stop and the
	// executor is not called a second time.
	_, err = srv.CancelAgentRun(ctx, &runtimev1.CancelAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	require.Equal(t, []string{start.RunId}, exec.cancels)

	// Cancelling a run in another instance is NotFound (scoping guard runs before the executor).
	_, err = srv.CancelAgentRun(ctx, &runtimev1.CancelAgentRunRequest{InstanceId: "inst-other", RunId: start.RunId})
	require.Equal(t, codes.NotFound, status.Code(err))
}

// collectStream is a test double for AgentService_StreamAgentRunEventsServer that records the events sent to it.
type collectStream struct {
	ctx    context.Context
	mu     sync.Mutex
	events []*runtimev1.AgentRunEvent
}

func (c *collectStream) Send(r *runtimev1.StreamAgentRunEventsResponse) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, r.Event)
	return nil
}
func (c *collectStream) Context() context.Context     { return c.ctx }
func (c *collectStream) SetHeader(metadata.MD) error  { return nil }
func (c *collectStream) SendHeader(metadata.MD) error { return nil }
func (c *collectStream) SetTrailer(metadata.MD)       {}
func (c *collectStream) SendMsg(any) error            { return nil }
func (c *collectStream) RecvMsg(any) error            { return nil }

// TestAgentServiceStreamEvents verifies the event stream drains a terminal run's events in order and then returns.
func TestAgentServiceStreamEvents(t *testing.T) {
	srv, store, _, instanceID := newActServer(t)
	ctx := testCtx()

	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k"})
	require.NoError(t, err)
	runID := start.RunId
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusRunning, EventType: act.EventTypeRunning}))
	require.NoError(t, store.RecordTransition(ctx, act.RunTransition{InstanceID: instanceID, RunID: runID, Status: act.RunStatusSucceeded, EventType: act.EventTypeSucceeded}))

	stream := &collectStream{ctx: ctx}
	// The run is terminal, so the stream drains the backlog and returns without blocking.
	err = srv.StreamAgentRunEvents(&runtimev1.StreamAgentRunEventsRequest{InstanceId: instanceID, RunId: runID}, stream)
	require.NoError(t, err)

	got := make([]string, len(stream.events))
	for i, e := range stream.events {
		got[i] = e.EventType
	}
	require.Equal(t, []string{act.EventTypeQueued, act.EventTypeRunning, act.EventTypeSucceeded}, got)

	// Resuming after the first event's cursor yields only the later events.
	require.NotEmpty(t, stream.events)
	after := stream.events[0].Id
	stream2 := &collectStream{ctx: ctx}
	err = srv.StreamAgentRunEvents(&runtimev1.StreamAgentRunEventsRequest{InstanceId: instanceID, RunId: runID, AfterId: after}, stream2)
	require.NoError(t, err)
	require.Len(t, stream2.events, 2)
}

// TestAgentServiceRunEndpointsUnimplementedWithoutStore verifies the run endpoints fail cleanly when Act is not
// configured, while discovery keeps working.
func TestAgentServiceRunEndpointsUnimplementedWithoutStore(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: agentProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)
	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)
	ctx := testCtx()

	// Discovery works without the Act plane.
	_, err = srv.ListAgents(ctx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)

	// Run endpoints report Unimplemented.
	_, err = srv.ListAgentRuns(ctx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.Equal(t, codes.Unimplemented, status.Code(err))
	_, err = srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage"})
	require.Equal(t, codes.Unimplemented, status.Code(err))
}
