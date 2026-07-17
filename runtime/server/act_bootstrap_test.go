package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/ratelimit"
	"github.com/rilldata/rill/runtime/server"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestBootstrapActEndToEnd drives the real bootstrap (not the fake executor of the other agent tests): it stands up a
// runtime with one reconciled agent, runs BootstrapAct against the local Postgres, and verifies the plane serves the
// AgentService end to end — discovery lists the agent, and a manual StartAgentRun enqueues a durable run that is
// immediately queryable through the real store. It does not assert the run succeeds: without a model connector the
// worker will fail the run asynchronously, but the run is created and queryable regardless, which is what wiring must
// guarantee. A long dispatch interval keeps the trigger loop out of this test's way; the loop is exercised by the
// runtime/act/trigger package.
func TestBootstrapActEndToEnd(t *testing.T) {
	dsn := requireStorePostgres(t)

	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: agentProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)

	// Unique product schema + application version isolate this test's runs on the shared Postgres/DBOS instance.
	unique := uuid.New().String()[:8]
	closer, err := srv.BootstrapAct(t.Context(), server.ActConfig{
		PostgresDSN:        dsn,
		ProductSchema:      "act_prod_it_" + unique,
		ApplicationVersion: "act-it-" + unique,
		DispatchInterval:   time.Hour,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = closer.Close() })

	ctx := testCtx()

	// Discovery works over the real plane.
	list, err := srv.ListAgents(ctx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, list.Agents, 1)
	require.Equal(t, "triage", list.Agents[0].Name)

	// StartAgentRun enqueues a durable run through the real executor and returns it read back from the real store.
	start, err := srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{
		InstanceId:     instanceID,
		Name:           "triage",
		IdempotencyKey: "it-key-1",
		Prompt:         "El coste subio.",
	})
	require.NoError(t, err)
	require.NotEmpty(t, start.RunId)
	require.Equal(t, "triage", start.AgentName)

	// The run is immediately queryable, scoped to its instance and recorded as a manual trigger.
	got, err := srv.GetAgentRun(ctx, &runtimev1.GetAgentRunRequest{InstanceId: instanceID, RunId: start.RunId})
	require.NoError(t, err)
	require.Equal(t, start.RunId, got.Run.RunId)
	require.Equal(t, "manual", got.Run.Trigger)

	// It shows up in the instance's run list on the isolated product schema.
	runs, err := srv.ListAgentRuns(ctx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, runs.Runs, 1)
	require.Equal(t, start.RunId, runs.Runs[0].RunId)
}

// TestActDisabledServesWithoutRunPlane is the zero-regression case: without BootstrapAct (no Act Postgres DSN) the
// server still stands up and agent discovery works, but the run/approval surface reports Unimplemented rather than
// erroring hard or panicking. This is exactly how a deployment that does not provision Act behaves. It needs no
// Postgres, so it always runs.
func TestActDisabledServesWithoutRunPlane(t *testing.T) {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{Files: agentProjectFiles()})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)
	// BootstrapAct is intentionally NOT called: Act is disabled.

	ctx := testCtx()

	// Discovery works without Act (it only needs the resource catalog).
	list, err := srv.ListAgents(ctx, &runtimev1.ListAgentsRequest{InstanceId: instanceID})
	require.NoError(t, err)
	require.Len(t, list.Agents, 1)

	// The run/approval endpoints report Unimplemented until Act is configured.
	_, err = srv.StartAgentRun(ctx, &runtimev1.StartAgentRunRequest{InstanceId: instanceID, Name: "triage", IdempotencyKey: "k"})
	require.Equal(t, codes.Unimplemented, status.Code(err))

	_, err = srv.ListAgentRuns(ctx, &runtimev1.ListAgentRunsRequest{InstanceId: instanceID})
	require.Equal(t, codes.Unimplemented, status.Code(err))
}
