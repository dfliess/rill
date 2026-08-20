package act

import (
	"context"
	"errors"
	"fmt"
	"slices"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/drivers"
)

// CatalogAgentProvider resolves dynamic agent definitions from the runtime's resource catalog. It is the production
// implementation of ai.AgentDefinitionProvider: where StaticAgentProvider holds a fixed in-memory map, this reads
// the reconciled Agent resource for an instance.
//
// It resolves each agent from AgentState.ValidSpec, the spec the reconciler has proven valid, not the raw Spec. An
// agent that failed validation therefore looks absent rather than runnable. The snapshot is a copy taken at read
// time; the executor captures it once at run start, so later edits to the resource never reach an in-flight run.
type CatalogAgentProvider struct {
	rt *runtime.Runtime
}

var _ ai.AgentDefinitionProvider = (*CatalogAgentProvider)(nil)

// NewCatalogAgentProvider returns a provider backed by rt's catalog.
func NewCatalogAgentProvider(rt *runtime.Runtime) *CatalogAgentProvider {
	return &CatalogAgentProvider{rt: rt}
}

func (p *CatalogAgentProvider) GetAgent(ctx context.Context, instanceID, name string) (*ai.AgentSnapshot, error) {
	ctrl, err := p.rt.Controller(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	res, err := ctrl.Get(ctx, &runtimev1.ResourceName{Kind: runtime.ResourceKindAgent, Name: name}, false)
	if err != nil {
		if errors.Is(err, drivers.ErrResourceNotFound) {
			return nil, fmt.Errorf("%w: %q", ai.ErrAgentNotFound, name)
		}
		return nil, err
	}

	spec := res.GetAgent().State.ValidSpec
	if spec == nil {
		// The resource exists but has not reconciled to a valid spec yet: not runnable, so treat it as absent.
		return nil, fmt.Errorf("%w: %q (no valid spec)", ai.ErrAgentNotFound, name)
	}
	snap := snapshotFromSpec(res.Meta.Name.Name, res.GetAgent().State.SpecHash, spec)
	if err := p.resolveModelSnapshot(ctx, instanceID, snap); err != nil {
		return nil, fmt.Errorf("resolve model for agent %q: %w", name, err)
	}
	return snap, nil
}

func (p *CatalogAgentProvider) ListAgents(ctx context.Context, instanceID string) ([]*ai.AgentSnapshot, error) {
	ctrl, err := p.rt.Controller(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	resources, err := ctrl.List(ctx, runtime.ResourceKindAgent, "", false)
	if err != nil {
		return nil, err
	}

	snapshots := make([]*ai.AgentSnapshot, 0, len(resources))
	for _, res := range resources {
		spec := res.GetAgent().State.ValidSpec
		if spec == nil {
			continue
		}
		snap := snapshotFromSpec(res.Meta.Name.Name, res.GetAgent().State.SpecHash, spec)
		if err := p.resolveModelSnapshot(ctx, instanceID, snap); err != nil {
			return nil, fmt.Errorf("resolve model for agent %q: %w", res.Meta.Name.Name, err)
		}
		snapshots = append(snapshots, snap)
	}
	return snapshots, nil
}

// resolveModelSnapshot freezes the effective non-secret AI connector configuration into snap. The connector name is
// resolved once here (including the project default), so changing either the agent or rill.yaml cannot reroute a run
// after DBOS checkpoints the snapshot. Secrets are excluded and intentionally resolved later by the session factory.
func (p *CatalogAgentProvider) resolveModelSnapshot(ctx context.Context, instanceID string, snap *ai.AgentSnapshot) error {
	connector := snap.ModelConnector
	if connector == "" {
		inst, err := p.rt.Instance(ctx, instanceID)
		if err != nil {
			return err
		}
		connector = inst.ResolveAIConnector()
	}
	// Preserve compatibility for installations with no AI connector configured. A production SessionFactory will
	// return ErrAINotConfigured when execution reaches the model; tests may still inject a scripted LLM.
	if connector == "" {
		return nil
	}

	cfg, err := p.rt.ConnectorConfig(ctx, instanceID, connector)
	if err != nil {
		return err
	}
	properties, err := runtime.NonSecretConnectorProperties(cfg.Driver, cfg.Resolve())
	if err != nil {
		return err
	}
	// The agent-level name is the most specific setting and therefore overrides the connector's default model. Store
	// the effective value in ModelProperties so the exact configuration handed to the driver is checkpointed.
	if snap.ModelName != "" {
		// Rill's managed admin proxy does not expose a model parameter in its completion API. Accepting an override here
		// would record intent that the provider silently ignores, so require an explicit provider connector instead.
		if cfg.Driver == "admin" {
			return fmt.Errorf("model.name is not supported by the managed AI connector; select an explicit AI connector")
		}
		properties["model"] = snap.ModelName
	}

	snap.ModelConnector = connector
	snap.ModelDriver = cfg.Driver
	snap.ModelProperties = properties
	return nil
}

// snapshotFromSpec projects a validated AgentSpec onto the ai.AgentSnapshot the executor runs from. specHash is the
// reconciler's stable hash of the effective spec, carried onto the snapshot so a run can be bound to its version.
func snapshotFromSpec(name, specHash string, spec *runtimev1.AgentSpec) *ai.AgentSnapshot {
	snap := &ai.AgentSnapshot{
		Name:           name,
		SpecHash:       specHash,
		DisplayName:    spec.DisplayName,
		Instructions:   spec.Instructions,
		ModelConnector: spec.ModelConnector,
		ModelName:      spec.ModelName,
		Tools:          slices.Clone(spec.Tools),
	}
	if spec.Limits != nil {
		snap.MaxSteps = int(spec.Limits.MaxSteps)
		snap.TimeoutSeconds = int(spec.Limits.TimeoutSeconds)
	}
	if len(spec.Mcp) > 0 {
		snap.MCPConnectors = make([]ai.MCPConnector, len(spec.Mcp))
		for i, m := range spec.Mcp {
			snap.MCPConnectors[i] = ai.MCPConnector{
				Name:              m.Name,
				URL:               m.Url,
				AuthSecret:        m.AuthSecret,
				AllowedHosts:      slices.Clone(m.AllowedHosts),
				TrustReadOnlyHint: m.TrustReadOnlyHint,
				Approval:          m.Approval,
				RequireApproval:   slices.Clone(m.RequireApproval),
				AutoApprove:       slices.Clone(m.AutoApprove),
			}
		}
	}
	return snap
}
