package reconcilers

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/rilldata/rill/runtime/ai"
)

func init() {
	runtime.RegisterReconcilerInitializer(runtime.ResourceKindAgent, newAgentReconciler)
}

// AgentReconciler validates declarative agent resources (Kairos Act, experimental).
// It only validates the spec and records a valid snapshot and hash; it never creates
// schedules, starts runs, or makes network calls (beyond resolving connector config).
type AgentReconciler struct {
	C *runtime.Controller
}

func newAgentReconciler(ctx context.Context, c *runtime.Controller) (runtime.Reconciler, error) {
	return &AgentReconciler{C: c}, nil
}

func (r *AgentReconciler) Close(ctx context.Context) error {
	return nil
}

func (r *AgentReconciler) AssignSpec(from, to *runtimev1.Resource) error {
	a := from.GetAgent()
	b := to.GetAgent()
	if a == nil || b == nil {
		return fmt.Errorf("cannot assign spec from %T to %T", from.Resource, to.Resource)
	}
	b.Spec = a.Spec
	return nil
}

func (r *AgentReconciler) AssignState(from, to *runtimev1.Resource) error {
	a := from.GetAgent()
	b := to.GetAgent()
	if a == nil || b == nil {
		return fmt.Errorf("cannot assign state from %T to %T", from.Resource, to.Resource)
	}
	b.State = a.State
	return nil
}

func (r *AgentReconciler) ResetState(res *runtimev1.Resource) error {
	res.GetAgent().State = &runtimev1.AgentState{}
	return nil
}

func (r *AgentReconciler) Reconcile(ctx context.Context, n *runtimev1.ResourceName) runtime.ReconcileResult {
	self, err := r.C.Get(ctx, n, true)
	if err != nil {
		return runtime.ReconcileResult{Err: err}
	}
	a := self.GetAgent()
	if a == nil {
		return runtime.ReconcileResult{Err: errors.New("not an agent")}
	}

	// Exit early for deletion
	if self.Meta.DeletedOn != nil {
		return runtime.ReconcileResult{}
	}

	// Compute a stable hash of the spec.
	specHash, err := r.specHash(a.Spec)
	if err != nil {
		return runtime.ReconcileResult{Err: err}
	}

	// Validate the spec (no execution, no schedules).
	validateErr := r.validateSpec(ctx, a.Spec)
	if validateErr != nil {
		// Clear the previously valid spec and surface the error.
		a.State.ValidSpec = nil
		a.State.SpecHash = specHash
		if err := r.C.UpdateState(ctx, self.Meta.Name, self); err != nil {
			return runtime.ReconcileResult{Err: err}
		}
		return runtime.ReconcileResult{Err: validateErr}
	}

	// Capture the spec, which we now know to be valid.
	a.State.ValidSpec = a.Spec
	a.State.SpecHash = specHash
	if err := r.C.UpdateState(ctx, self.Meta.Name, self); err != nil {
		return runtime.ReconcileResult{Err: err}
	}

	return runtime.ReconcileResult{}
}

func (r *AgentReconciler) ResolveTransitiveAccess(ctx context.Context, claims *runtime.SecurityClaims, res *runtimev1.Resource) ([]*runtimev1.SecurityRule, error) {
	if res.GetAgent() == nil {
		return nil, fmt.Errorf("not an agent resource")
	}
	return []*runtimev1.SecurityRule{{Rule: runtime.SelfAllowRuleAccess(res)}}, nil
}

// validateSpec performs lightweight, execution-free validation of an agent spec.
func (r *AgentReconciler) validateSpec(ctx context.Context, spec *runtimev1.AgentSpec) error {
	// Instructions are the agent's system prompt and are required (also enforced by the parser).
	if strings.TrimSpace(spec.Instructions) == "" {
		return errors.New(`agent must set "instructions"`)
	}

	// Fail-closed on the declared tools: a v1 dynamic agent may only use Rill's read-only analytical tools.
	// A tool outside the allowlist fails validation, so ValidSpec stays nil and the agent is never executable.
	for _, t := range spec.Tools {
		if !ai.IsDeclarableAgentTool(t) {
			return fmt.Errorf("agent tool %q is not allowed: v1 agents may only use read-only analytical tools", t)
		}
	}

	// If a model connector is declared, it must resolve in the project.
	// ConnectorConfig resolves the connector's configuration without opening a connection.
	if spec.ModelConnector != "" {
		_, err := r.C.Runtime.ConnectorConfig(ctx, r.C.InstanceID, spec.ModelConnector)
		if err != nil {
			return fmt.Errorf("invalid model connector %q: %w", spec.ModelConnector, err)
		}
	}

	// Validate each MCP connector fail-closed: an unparseable or unsafe connector (bad transport, non-https URL,
	// a private-range opt-out) fails validation, so ValidSpec stays nil and the agent is never executable. This
	// only parses and validates the definition; it opens no connection and resolves no secret.
	for _, m := range spec.Mcp {
		props := map[string]any{
			"transport": mcpconn.TransportStreamableHTTP,
			"url":       m.Url,
			"auth": map[string]any{
				"secret": m.AuthSecret,
			},
			"network": map[string]any{
				"allowed_hosts":  m.AllowedHosts,
				"private_ranges": mcpconn.PrivateRangesDeny,
			},
			"trust_read_only_hint": m.TrustReadOnlyHint,
		}
		if _, err := mcpconn.ParseConnectorConfig(m.Name, props); err != nil {
			return fmt.Errorf("invalid mcp connector %q: %w", m.Name, err)
		}
	}

	return nil
}

// specHash returns a stable hash of the agent spec, used to detect changes.
func (r *AgentReconciler) specHash(spec *runtimev1.AgentSpec) (string, error) {
	h := sha256.New()
	// writeField length-prefixes each string so field boundaries are unambiguous: without it, concatenating fields
	// collides (e.g. "ab"+"c" hashes the same as "a"+"bc").
	writeField := func(s string) error {
		if err := binary.Write(h, binary.BigEndian, uint32(len(s))); err != nil {
			return err
		}
		_, err := h.Write([]byte(s))
		return err
	}

	for _, s := range []string{
		spec.DisplayName,
		spec.Description,
		spec.ModelConnector,
		spec.ModelName,
		spec.Instructions,
	} {
		if err := writeField(s); err != nil {
			return "", err
		}
	}

	// Prefix the tool list with its length so the boundary between the tools and the fields that follow is fixed.
	if err := binary.Write(h, binary.BigEndian, uint32(len(spec.Tools))); err != nil {
		return "", err
	}
	for _, t := range spec.Tools {
		if err := writeField(t); err != nil {
			return "", err
		}
	}

	if spec.Limits != nil {
		if err := binary.Write(h, binary.BigEndian, spec.Limits.MaxSteps); err != nil {
			return "", err
		}
		if err := binary.Write(h, binary.BigEndian, spec.Limits.TimeoutSeconds); err != nil {
			return "", err
		}
	}

	// Fold in the MCP connectors so editing one changes the hash (the run records the hash as its audit binding).
	// Prefix the connector list with its length so its boundary is fixed, then each connector's fields.
	if err := binary.Write(h, binary.BigEndian, uint32(len(spec.Mcp))); err != nil {
		return "", err
	}
	for _, m := range spec.Mcp {
		for _, s := range []string{m.Name, m.Url, m.AuthSecret} {
			if err := writeField(s); err != nil {
				return "", err
			}
		}
		// Hash the allowed hosts in a stable order so a reordering of an equivalent allowlist does not churn the hash.
		hosts := slices.Clone(m.AllowedHosts)
		slices.Sort(hosts)
		if err := binary.Write(h, binary.BigEndian, uint32(len(hosts))); err != nil {
			return "", err
		}
		for _, host := range hosts {
			if err := writeField(host); err != nil {
				return "", err
			}
		}
		var trust uint8
		if m.TrustReadOnlyHint {
			trust = 1
		}
		if err := binary.Write(h, binary.BigEndian, trust); err != nil {
			return "", err
		}
		// Fold in the approval posture and its glob exceptions so editing which actions run unsupervised changes the
		// hash (a run binds to the hash for audit). Sort the patterns so reordering an equivalent list does not churn.
		if err := writeField(m.Approval); err != nil {
			return "", err
		}
		for _, pats := range [][]string{m.RequireApproval, m.AutoApprove} {
			ps := slices.Clone(pats)
			slices.Sort(ps)
			if err := binary.Write(h, binary.BigEndian, uint32(len(ps))); err != nil {
				return "", err
			}
			for _, p := range ps {
				if err := writeField(p); err != nil {
					return "", err
				}
			}
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
