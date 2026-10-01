package act_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/stretchr/testify/require"
)

// panickingRunner panics inside RunSegment when the prompt asks for it, and otherwise delegates to the real runner.
// RunSegment runs inside a checkpointed step, so this is a panic in a step of the durable workflow.
type panickingRunner struct {
	act.Runner
}

func (r *panickingRunner) RunSegment(ctx context.Context, in act.RunSegmentInput) (act.RunSegmentResult, error) {
	if strings.Contains(in.Prompt, "PANIC") {
		panic("boom in a step")
	}
	return r.Runner.RunSegment(ctx, in)
}

// TestExecutorSurvivesPanicInStep verifies kairos-cloud#215: a panic inside a step fails that run and leaves the
// process, and every other run in it, alive. Up to DBOS v1.3.0 the panic escaped the workflow goroutine and killed the
// whole runtime, runs parked on human approval included. Since v1.5.0 DBOS recovers it and returns ErrWorkflowPanic
// as the run's result. Verified by mutation: with v1.3.0 this test binary crashes on the first run.
func TestExecutorSurvivesPanicInStep(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	base := &act.SessionRunner{
		Provider: act.NewCatalogAgentProvider(rt),
		Sessions: scriptedSessionFactory(rt, &scriptedAI{response: "todo en orden"}),
	}
	e := newExecutorWithRunner(t, &panickingRunner{Runner: base})
	ctx := t.Context()

	panicID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "PANIC",
		IdempotencyKey: "run-panic-" + uuid.NewString(),
	})
	require.NoError(t, err)
	_, err = e.Result(panicID)
	require.Error(t, err)
	require.True(t, errors.Is(err, dbos.ErrWorkflowPanic), "the panic comes back as the run's error, got: %v", err)

	okID, err := e.Start(ctx, act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "alerta normal",
		IdempotencyKey: "run-after-panic-" + uuid.NewString(),
	})
	require.NoError(t, err)
	res, err := e.Result(okID)
	require.NoError(t, err)
	require.Equal(t, act.RunStatusSucceeded, res.Status, "the executor keeps serving runs after a panic")
}
