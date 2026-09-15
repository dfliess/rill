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
			"triggers/nightly.yaml": `
type: agent_trigger
agent: stale
source:
  kind: schedule
  cron: "0 3 * * *"
input:
  prompt: Run the nightly summary.
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
	trigger := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "nightly")
	healthy := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "healthy")
	require.NotNil(t, stale.GetAgent().State.ValidSpec)
	require.NotNil(t, trigger.GetAgentTrigger().State.ValidSpec)

	// Persist what the controller holds in memory, so the rows exist in the store.
	ctrl, err := rt.Controller(ctx, id)
	require.NoError(t, err)
	require.NoError(t, ctrl.Flush(ctx))

	// Rewrite the agent and trigger rows as they were stored while Agent and AgentTrigger were fields 18 and 19.
	catalog, release, err := rt.Catalog(ctx, id)
	require.NoError(t, err)
	v, err := catalog.NextControllerVersion(ctx)
	require.NoError(t, err)
	err = catalog.UpdateResource(ctx, v, staleActRow(t, stale, 18, stale.GetAgent()))
	require.NoError(t, err)
	err = catalog.UpdateResource(ctx, v, staleActRow(t, trigger, 19, trigger.GetAgentTrigger()))
	require.NoError(t, err)
	release()

	// Restart the controller, so it loads the catalog from the store again.
	inst, err := rt.Instance(ctx, id)
	require.NoError(t, err)
	require.NoError(t, rt.EditInstance(ctx, inst, true))

	// No manual reconcile: starting the controller enqueues every resource, the parser included, and that alone must
	// bring both back. Waiting for idle is retried because the new controller may be ready before it has enqueued
	// anything, and then it is idle right away.
	require.Eventually(t, func() bool {
		c, err := rt.Controller(ctx, id)
		if err != nil || c == ctrl || c.WaitUntilIdle(ctx, false) != nil {
			return false
		}
		a, err := c.Get(ctx, stale.Meta.Name, false)
		if err != nil || a.GetAgent().GetState().GetValidSpec() == nil {
			return false
		}
		tr, err := c.Get(ctx, trigger.Meta.Name, false)
		return err == nil && tr.GetAgentTrigger().GetState().GetValidSpec() != nil
	}, 10*time.Second, 50*time.Millisecond)
	testruntime.RequireReconcileState(t, rt, id, -1, 0, 0)

	// The stale rows were dropped and the parser recreated both resources from their files.
	recreated := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "stale")
	require.Equal(t, "Investigate using governed metrics only.", recreated.GetAgent().State.ValidSpec.Instructions)
	require.True(t, recreated.Meta.CreatedOn.AsTime().After(stale.Meta.CreatedOn.AsTime()))

	recreatedTrigger := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgentTrigger, "nightly")
	require.Empty(t, recreatedTrigger.Meta.ReconcileError)
	require.Equal(t, "stale", recreatedTrigger.GetAgentTrigger().State.ValidSpec.Agent)
	require.True(t, recreatedTrigger.Meta.CreatedOn.AsTime().After(trigger.Meta.CreatedOn.AsTime()))
	require.Len(t, recreatedTrigger.Meta.Refs, 1)
	require.Equal(t, runtime.ResourceKindAgent, recreatedTrigger.Meta.Refs[0].Kind)
	require.Equal(t, "stale", recreatedTrigger.Meta.Refs[0].Name)

	// The healthy row was loaded as it was.
	kept := testruntime.GetResource(t, rt, id, runtime.ResourceKindAgent, "healthy")
	require.NotNil(t, kept.GetAgent())
	require.True(t, kept.Meta.CreatedOn.AsTime().Equal(healthy.Meta.CreatedOn.AsTime()))
	require.Equal(t, healthy.GetAgent().State.SpecHash, kept.GetAgent().State.SpecHash)
}

// staleActRow encodes r as it was stored while its Act payload was the given field of Resource.
func staleActRow(t *testing.T, r *runtimev1.Resource, field protowire.Number, payload proto.Message) drivers.Resource {
	data, err := proto.Marshal(&runtimev1.Resource{Meta: r.Meta})
	require.NoError(t, err)
	b, err := proto.Marshal(payload)
	require.NoError(t, err)
	data = protowire.AppendTag(data, field, protowire.BytesType)
	data = protowire.AppendBytes(data, b)

	// Sanity check: under the current numbering the row decodes with no resource at all.
	decoded := &runtimev1.Resource{}
	require.NoError(t, proto.Unmarshal(data, decoded))
	require.Nil(t, decoded.GetResource())

	return drivers.Resource{Kind: r.Meta.Name.Kind, Name: r.Meta.Name.Name, Data: data}
}
