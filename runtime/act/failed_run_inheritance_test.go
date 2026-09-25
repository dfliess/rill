package act_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/act"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/stretchr/testify/require"
)

// failingAI reproduces a provider outage: every completion is rejected before any tokens are spent. The message is
// the one that actually took a dashboard note down (kairos-cloud#154).
type failingAI struct {
	mu    sync.Mutex
	calls int
}

var _ drivers.AIService = (*failingAI)(nil)

func (f *failingAI) Complete(_ context.Context, _ *drivers.CompleteOptions) (*drivers.CompleteResult, error) {
	f.mu.Lock()
	f.calls++
	f.mu.Unlock()
	return nil, errors.New("Error 403, Message: Your project has been denied access. Please contact support., Status: PERMISSION_DENIED")
}

func (f *failingAI) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// TestExecutorInheritsFailedRun pins the property the TRIGGER DISPATCHER depends on, and that issue #154 must not
// take away while fixing the manual path.
//
// The dispatcher delivers events at-least-once, so the same event can reach Start twice. Inheriting the dead run is
// what stops the second delivery from re-running an agent that may already have performed an external effect before
// it died. The fix for the note therefore lives in the API handler (see resolveKeyFamily in runtime/server), NOT
// here: if this test ever goes red, a redelivered alert has become able to fire its action twice.
func TestExecutorInheritsFailedRun(t *testing.T) {
	rt, instanceID := newInstanceWithAgent(t, "Investiga y propon.")
	llm := &failingAI{}
	e := newExecutor(t, rt, llm)

	in := act.AgentRunInput{
		InstanceID:     instanceID,
		AgentName:      "triage",
		Prompt:         "alerta que muere en el proveedor",
		IdempotencyKey: "run-failed-inherit-" + uuid.NewString(),
	}

	id1, err := e.Start(t.Context(), in)
	require.NoError(t, err)

	// The run reaches a terminal failure: Result surfaces the provider's error rather than a result.
	_, err = e.Result(id1)
	require.Error(t, err)
	callsAfterFailure := llm.callCount()
	require.Positive(t, callsAfterFailure)

	// The redelivery. It must attach to the dead run, not start a second one.
	id2, err := e.Start(t.Context(), in)
	require.NoError(t, err)
	require.Equal(t, id1, id2)
	require.Equal(t, callsAfterFailure, llm.callCount(),
		"a redelivered event re-executed a failed run: the dispatcher's at-least-once guard is gone")
}
