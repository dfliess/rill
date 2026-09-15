package runtime

import (
	"context"
	"fmt"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// dropStaleActResource deletes an Agent or AgentTrigger row that was stored while those resources used fields 18/19
// of Resource, and reports whether it did. Such a row no longer decodes as its kind: field 18/19 is unknown now, and
// decodes as a Skill once upstream's Skill = 18 lands. It is safe to drop because its state only holds a valid spec
// and a spec hash, which the reconciler recomputes after the project parser recreates the resource.
//
// TODO(kairos-cloud#158): remove once dev and prod have started with fields 1000/1001.
func (c *catalogCache) dropStaleActResource(ctx context.Context, r drivers.Resource) (bool, error) {
	if r.Kind != ResourceKindAgent && r.Kind != ResourceKindAgentTrigger {
		return false, nil
	}

	res := &runtimev1.Resource{}
	if err := proto.Unmarshal(r.Data, res); err == nil {
		if r.Kind == ResourceKindAgent && res.GetAgent() != nil {
			return false, nil
		}
		if r.Kind == ResourceKindAgentTrigger && res.GetAgentTrigger() != nil {
			return false, nil
		}
	}

	err := c.store.DeleteResource(ctx, c.version, r.Kind, r.Name)
	if err != nil {
		return false, fmt.Errorf("failed to drop stale %s %q: %w", r.Kind, r.Name, err)
	}
	c.ctrl.Logger.Warn("dropped catalog resource stored with a stale field number; the project parser will recreate it", zap.String("kind", r.Kind), zap.String("name", r.Name))
	return true, nil
}
