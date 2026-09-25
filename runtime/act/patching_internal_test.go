package act

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
	"github.com/google/uuid"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// patchProbeRunner is the minimal Runner NewDBOSExecutor requires; nothing here ever runs it.
type patchProbeRunner struct{}

func (patchProbeRunner) LoadSnapshot(_ context.Context, _, agentName string) (*ai.AgentSnapshot, error) {
	return &ai.AgentSnapshot{Name: agentName, Instructions: "patch probe"}, nil
}

func (patchProbeRunner) RunSegment(context.Context, RunSegmentInput) (RunSegmentResult, error) {
	return RunSegmentResult{Response: "ok", SessionID: "sess-probe"}, nil
}

// TestNewDBOSExecutorArmsPatching pins that the executor's DBOS context has the patching system armed. With the
// pinned ApplicationVersion (see defaultApplicationVersion) a deploy replays checkpoints recorded by older code, and
// a named dbos.Patch is the only supported way to change the workflow's step sequence over runs already in flight —
// so an executor without EnablePatching would leave every future shape change with no safe path, which is how the
// 2026-09-02 deploy killed 19 parked runs. TestPatchGuardsShapeChangeAgainstParkedRuns proves the mechanism works;
// this one proves the production constructor actually arms it.
//
// The probe leans on an ordering dbos.Patch guarantees: it checks the EnablePatching flag before it checks for a
// workflow, so on an ARMED context a call from outside any workflow fails with "patch can only be called within a
// workflow", while a disarmed one fails with ErrorCodePatchingNotEnabled. Remove EnablePatching from NewDBOSExecutor
// and this test fails.
func TestNewDBOSExecutorArmsPatching(t *testing.T) {
	dsn := requirePostgresInternal(t)

	e, err := NewDBOSExecutor(t.Context(), Config{
		DatabaseURL:        dsn,
		DatabaseSchema:     "act_dbos",
		ApplicationVersion: "act-test-" + uuid.NewString(),
		Runner:             patchProbeRunner{},
		Store:              &PostgresRunStore{},
		Gateway:            &Gateway{},
		Proposer:           NewCapturedProposer(),
	})
	require.NoError(t, err)
	t.Cleanup(func() { e.Close(5 * time.Second) })

	_, err = dbos.Patch(e.ctx, "probe")
	require.Error(t, err, "a Patch outside any workflow cannot succeed; the point is WHICH error it fails with")
	var derr *dbos.Error
	if errors.As(err, &derr) {
		require.NotEqual(t, dbos.ErrorCodePatchingNotEnabled, derr.Code,
			"NewDBOSExecutor must set EnablePatching: with the pinned ApplicationVersion, dbos.Patch is the only safe way to change the workflow's shape")
	}
}
