package reconcilers

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
)

func init() {
	runtime.RegisterReconcilerInitializer(runtime.ResourceKindAgentTrigger, newAgentTriggerReconciler)
}

// validAgentTriggerEvents lists the event short-names each source kind can emit. It is validation data only:
// the transition semantics and the fully-qualified event types ("alert.entered_fail") live behind the execution
// observer in runtime/act/trigger, which the reconcilers must not depend on. A schedule source has no execution
// events (it fires on its own clock), so it declares none.
var validAgentTriggerEvents = map[runtimev1.AgentTriggerSourceKind]map[string]bool{
	runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT: {
		"entered_fail":  true,
		"recovered":     true,
		"entered_error": true,
		"renotify_due":  true,
		"evaluated":     true,
	},
	runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_REPORT: {
		"completed": true,
		"failed":    true,
	},
	runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_SCHEDULE: {},
}

// AgentTriggerReconciler validates declarative agent_trigger resources (Kairos Act, experimental).
// It only validates the spec and records a valid snapshot and hash; it never starts runs, creates schedules or
// makes network calls. Activation happens out of band in the trigger dispatcher (runtime/act/trigger), so no
// LLM call ever runs inside a reconciler.
type AgentTriggerReconciler struct {
	C *runtime.Controller
}

func newAgentTriggerReconciler(ctx context.Context, c *runtime.Controller) (runtime.Reconciler, error) {
	return &AgentTriggerReconciler{C: c}, nil
}

func (r *AgentTriggerReconciler) Close(ctx context.Context) error {
	return nil
}

func (r *AgentTriggerReconciler) AssignSpec(from, to *runtimev1.Resource) error {
	a := from.GetAgentTrigger()
	b := to.GetAgentTrigger()
	if a == nil || b == nil {
		return fmt.Errorf("cannot assign spec from %T to %T", from.Resource, to.Resource)
	}
	b.Spec = a.Spec
	return nil
}

func (r *AgentTriggerReconciler) AssignState(from, to *runtimev1.Resource) error {
	a := from.GetAgentTrigger()
	b := to.GetAgentTrigger()
	if a == nil || b == nil {
		return fmt.Errorf("cannot assign state from %T to %T", from.Resource, to.Resource)
	}
	b.State = a.State
	return nil
}

func (r *AgentTriggerReconciler) ResetState(res *runtimev1.Resource) error {
	res.GetAgentTrigger().State = &runtimev1.AgentTriggerState{}
	return nil
}

func (r *AgentTriggerReconciler) Reconcile(ctx context.Context, n *runtimev1.ResourceName) runtime.ReconcileResult {
	self, err := r.C.Get(ctx, n, true)
	if err != nil {
		return runtime.ReconcileResult{Err: err}
	}
	t := self.GetAgentTrigger()
	if t == nil {
		return runtime.ReconcileResult{Err: errors.New("not an agent trigger")}
	}

	// Exit early for deletion
	if self.Meta.DeletedOn != nil {
		return runtime.ReconcileResult{}
	}

	// Compute a stable hash of the spec.
	specHash, err := r.specHash(t.Spec)
	if err != nil {
		return runtime.ReconcileResult{Err: err}
	}

	// Validate the spec (no execution, no schedules).
	validateErr := r.validateSpec(ctx, t.Spec)
	if validateErr != nil {
		// Clear the previously valid spec and surface the error.
		t.State.ValidSpec = nil
		t.State.SpecHash = specHash
		if err := r.C.UpdateState(ctx, self.Meta.Name, self); err != nil {
			return runtime.ReconcileResult{Err: err}
		}
		return runtime.ReconcileResult{Err: validateErr}
	}

	// Capture the spec, which we now know to be valid, resolving the actor's named user to attributes (like an alert).
	// A resolution failure is treated like a validation error: clear the valid spec and surface it.
	validSpec, err := r.resolveValidSpec(ctx, t.Spec, self.Meta.Name.Name)
	if err != nil {
		t.State.ValidSpec = nil
		t.State.SpecHash = specHash
		if uerr := r.C.UpdateState(ctx, self.Meta.Name, self); uerr != nil {
			return runtime.ReconcileResult{Err: uerr}
		}
		return runtime.ReconcileResult{Err: err}
	}
	t.State.ValidSpec = validSpec
	t.State.SpecHash = specHash
	if err := r.C.UpdateState(ctx, self.Meta.Name, self); err != nil {
		return runtime.ReconcileResult{Err: err}
	}

	return runtime.ReconcileResult{}
}

// resolveValidSpec returns the spec to store as ValidSpec: a clone of the authored spec with the actor's identity
// resolved to attributes. When the actor names a user (user_id/user_email) and sets no explicit attributes, it
// resolves that user's attributes via the admin service — reusing the alert-metadata path, exactly like an alert's
// query_for — and stores them, so the dispatcher hands the started run governed SecurityClaims. Explicit attributes
// take precedence and skip resolution; an admin service that cannot resolve (e.g. Rill Developer) leaves attributes
// empty, so the run fails closed rather than failing the reconcile. Resolving here (not at dispatch) means a change to
// the user's attributes is only picked up on the trigger's next reconcile — a v1 simplification.
func (r *AgentTriggerReconciler) resolveValidSpec(ctx context.Context, spec *runtimev1.AgentTriggerSpec, name string) (*runtimev1.AgentTriggerSpec, error) {
	valid, ok := proto.Clone(spec).(*runtimev1.AgentTriggerSpec)
	if !ok {
		return nil, errors.New("agent_trigger: clone spec")
	}
	actor := valid.Actor
	// Nothing to resolve: no actor, explicit attributes already set (they win), or no user named.
	if actor == nil || actor.Attributes != nil {
		return valid, nil
	}
	userID, userEmail := strings.TrimSpace(actor.UserId), strings.TrimSpace(actor.UserEmail)
	if userID == "" && userEmail == "" {
		return valid, nil
	}

	admin, release, err := r.C.Runtime.Admin(ctx, r.C.InstanceID)
	if err != nil {
		if errors.Is(err, drivers.ErrNotImplemented) {
			return valid, nil // no admin service (Rill Developer): fail closed, don't fail the reconcile
		}
		return nil, fmt.Errorf("agent_trigger: get admin client: %w", err)
	}
	defer release()

	// GetAlertMetadata resolves query_for to the user's attributes independent of any alert existing; with no
	// recipients it creates no tokens, so it is a clean user-attribute lookup here.
	meta, err := admin.GetAlertMetadata(ctx, name, "", nil, false, nil, userID, userEmail)
	if err != nil {
		if errors.Is(err, drivers.ErrNotImplemented) {
			return valid, nil
		}
		return nil, fmt.Errorf("agent_trigger: resolve actor attributes: %w", err)
	}
	if len(meta.QueryForAttributes) == 0 {
		return valid, nil
	}
	attrs, err := structpb.NewStruct(meta.QueryForAttributes)
	if err != nil {
		return nil, fmt.Errorf("agent_trigger: encode actor attributes: %w", err)
	}
	actor.Attributes = attrs
	return valid, nil
}

func (r *AgentTriggerReconciler) ResolveTransitiveAccess(ctx context.Context, claims *runtime.SecurityClaims, res *runtimev1.Resource) ([]*runtimev1.SecurityRule, error) {
	if res.GetAgentTrigger() == nil {
		return nil, fmt.Errorf("not an agent trigger resource")
	}
	return []*runtimev1.SecurityRule{{Rule: runtime.SelfAllowRuleAccess(res)}}, nil
}

// validateSpec performs lightweight, execution-free validation of an agent trigger spec.
func (r *AgentTriggerReconciler) validateSpec(ctx context.Context, spec *runtimev1.AgentTriggerSpec) error {
	// The trigger must name the agent it starts, and that agent must have reconciled to a valid, runnable spec.
	if strings.TrimSpace(spec.Agent) == "" {
		return errors.New(`agent_trigger must set "agent"`)
	}
	agentRes, err := r.C.Get(ctx, &runtimev1.ResourceName{Kind: runtime.ResourceKindAgent, Name: spec.Agent}, false)
	if err != nil {
		if errors.Is(err, drivers.ErrResourceNotFound) {
			return fmt.Errorf("agent_trigger references unknown agent %q", spec.Agent)
		}
		return err
	}
	if agentRes.GetAgent().State.ValidSpec == nil {
		return fmt.Errorf("agent_trigger references agent %q which has no valid spec", spec.Agent)
	}

	// The source must be coherent: a known kind, a source resource that exists (for alert/report), and an event
	// allowlist drawn from the events that kind can emit.
	if spec.Source == nil {
		return errors.New(`agent_trigger must set "source"`)
	}
	allowed, ok := validAgentTriggerEvents[spec.Source.Kind]
	if !ok {
		return errors.New(`agent_trigger must set a valid "source.kind"`)
	}

	switch spec.Source.Kind {
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT, runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_REPORT:
		// An empty source.name is a wildcard (ADR-0016): the trigger fires for any resource of this kind, so there is
		// nothing to resolve. Only when a name is given do we require it to name a resource that exists, so a typo fails
		// closed rather than silently never firing.
		if strings.TrimSpace(spec.Source.Name) != "" {
			sourceKind := runtime.ResourceKindAlert
			if spec.Source.Kind == runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_REPORT {
				sourceKind = runtime.ResourceKindReport
			}
			if _, err := r.C.Get(ctx, &runtimev1.ResourceName{Kind: sourceKind, Name: spec.Source.Name}, false); err != nil {
				if errors.Is(err, drivers.ErrResourceNotFound) {
					return fmt.Errorf("agent_trigger references unknown %s %q", strings.ToLower(runtime.PrettifyResourceKind(sourceKind)), spec.Source.Name)
				}
				return err
			}
		}
		if len(spec.Source.Events) == 0 {
			return errors.New(`agent_trigger must subscribe to at least one event in "source.events"`)
		}
		if strings.TrimSpace(spec.Source.Cron) != "" {
			return errors.New(`only a schedule agent_trigger may set "source.cron"`)
		}
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_SCHEDULE:
		if len(spec.Source.Events) != 0 {
			return errors.New(`a schedule agent_trigger must not set "source.events": it fires on its own clock`)
		}
		if strings.TrimSpace(spec.Source.Name) != "" {
			return errors.New(`a schedule agent_trigger must not set "source.name": it is its own source`)
		}
		// A schedule fires on its own clock, so it must carry a valid cron. Parse it here (the same parser the
		// dispatcher's ticker uses) so a malformed schedule fails closed at reconcile instead of silently never firing.
		if strings.TrimSpace(spec.Source.Cron) == "" {
			return errors.New(`a schedule agent_trigger must set "source.cron"`)
		}
		if _, err := cron.ParseStandard(spec.Source.Cron); err != nil {
			return fmt.Errorf(`invalid "source.cron" %q: %w`, spec.Source.Cron, err)
		}
	}

	// Every subscribed event must be one the source kind can emit, so a typo fails closed rather than never firing.
	for _, e := range spec.Source.Events {
		if !allowed[strings.TrimSpace(e)] {
			return fmt.Errorf("event %q is not valid for a %s source", e, strings.ToLower(runtime.PrettifyResourceKind(sourceKindResourceName(spec.Source.Kind))))
		}
	}

	// Identity mutual-exclusion (at most one of user_id, user_email, attributes) is enforced at parse, like an alert's
	// for.{...}: a spec only reaches here after that check, so there is nothing to re-validate.

	return nil
}

// sourceKindResourceName maps a trigger source kind to its runtime resource kind for user-facing messages.
func sourceKindResourceName(k runtimev1.AgentTriggerSourceKind) string {
	switch k {
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT:
		return runtime.ResourceKindAlert
	case runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_REPORT:
		return runtime.ResourceKindReport
	default:
		return "Schedule"
	}
}

// specHash returns a stable hash of the agent trigger spec, used to detect changes.
func (r *AgentTriggerReconciler) specHash(spec *runtimev1.AgentTriggerSpec) (string, error) {
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
	writeList := func(items []string) error {
		if err := binary.Write(h, binary.BigEndian, uint32(len(items))); err != nil {
			return err
		}
		for _, it := range items {
			if err := writeField(it); err != nil {
				return err
			}
		}
		return nil
	}

	if err := writeField(spec.Agent); err != nil {
		return "", err
	}

	if spec.Source != nil {
		if err := binary.Write(h, binary.BigEndian, int32(spec.Source.Kind)); err != nil {
			return "", err
		}
		if err := writeField(spec.Source.Name); err != nil {
			return "", err
		}
		if err := writeList(spec.Source.Events); err != nil {
			return "", err
		}
		if err := writeField(spec.Source.Cron); err != nil {
			return "", err
		}
	}

	if spec.Actor != nil {
		// Hash the AUTHORED identity (user_id/user_email), not the admin-resolved attributes: the hash must be stable
		// across environments, so it tracks who the run acts as, not the attributes that resolution happens to yield.
		if err := writeField(spec.Actor.UserId); err != nil {
			return "", err
		}
		if err := writeField(spec.Actor.UserEmail); err != nil {
			return "", err
		}
		// Hash the run-as attributes: json.Marshal of the map sorts keys, so the encoding is stable for an unchanged
		// spec (unlike map iteration order).
		if spec.Actor.Attributes != nil {
			attrs, err := json.Marshal(spec.Actor.Attributes.AsMap())
			if err != nil {
				return "", err
			}
			if err := writeField(string(attrs)); err != nil {
				return "", err
			}
		}
	}

	if spec.Input != nil {
		if err := writeField(spec.Input.Prompt); err != nil {
			return "", err
		}
		// Hash the context map in sorted key order: map iteration order is not deterministic, so an unordered walk
		// would make the hash unstable across reconciles for an unchanged spec.
		keys := make([]string, 0, len(spec.Input.Context))
		for k := range spec.Input.Context {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if err := binary.Write(h, binary.BigEndian, uint32(len(keys))); err != nil {
			return "", err
		}
		for _, k := range keys {
			if err := writeField(k); err != nil {
				return "", err
			}
			if err := writeField(spec.Input.Context[k]); err != nil {
				return "", err
			}
		}
	}

	if spec.Deduplication != nil {
		if err := binary.Write(h, binary.BigEndian, spec.Deduplication.WindowSeconds); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
