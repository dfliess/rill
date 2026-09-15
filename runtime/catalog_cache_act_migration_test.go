package runtime_test

import (
	"testing"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

// TODO(kairos-cloud#158): remove once dev and prod have started with fields 1000/1001.
func TestCatalogDropsActResourceWithStaleFieldNumber(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": "",
			"agents/stale.yaml": `
type: agent
instructions: Investigate using governed metrics only.
`,
			"agents/healthy.yaml": `
type: agent
instructions: Summarize the weekly report.
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, -1, 0, 0)
	ctx := t.Context()

	stale := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "stale")
	healthy := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "healthy")
	require.NotNil(t, stale.GetAgent().Spec)

	// Persist what the controller holds in memory, so both rows exist in the store.
	ctrl, err := rt.Controller(ctx, id)
	require.NoError(t, err)
	require.NoError(t, ctrl.Flush(ctx))

	// Rewrite the stale agent's row as it was stored while Agent was field 18 of Resource.
	data, err := proto.Marshal(&runtimev1.Resource{Meta: stale.Meta})
	require.NoError(t, err)
	agent, err := proto.Marshal(stale.GetAgent())
	require.NoError(t, err)
	data = protowire.AppendTag(data, 18, protowire.BytesType)
	data = protowire.AppendBytes(data, agent)

	decoded := &runtimev1.Resource{}
	require.NoError(t, proto.Unmarshal(data, decoded))
	require.Nil(t, decoded.GetAgent())

	catalog, release, err := rt.Catalog(ctx, id)
	require.NoError(t, err)
	v, err := catalog.NextControllerVersion(ctx)
	require.NoError(t, err)
	err = catalog.UpdateResource(ctx, v, drivers.Resource{Kind: runtime.ResourceKindAgent, Name: "stale", Data: data})
	require.NoError(t, err)
	release()

	// Restart the controller, so it loads the catalog from the store again.
	inst, err := rt.Instance(ctx, id)
	require.NoError(t, err)
	require.NoError(t, rt.EditInstance(ctx, inst, true))
	require.Eventually(t, func() bool {
		c, err := rt.Controller(ctx, id)
		return err == nil && c != ctrl
	}, 10*time.Second, 50*time.Millisecond)

	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, -1, 0, 0)

	// The stale row was dropped and the parser recreated the agent from its file.
	recreated := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "stale")
	require.NotNil(t, recreated.GetAgent())
	require.Equal(t, "Investigate using governed metrics only.", recreated.GetAgent().State.ValidSpec.Instructions)
	require.True(t, recreated.Meta.CreatedOn.AsTime().After(stale.Meta.CreatedOn.AsTime()))

	// The healthy row was loaded as it was.
	kept := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "healthy")
	require.NotNil(t, kept.GetAgent())
	require.True(t, kept.Meta.CreatedOn.AsTime().Equal(healthy.Meta.CreatedOn.AsTime()))
	require.Equal(t, healthy.GetAgent().State.SpecHash, kept.GetAgent().State.SpecHash)
}
