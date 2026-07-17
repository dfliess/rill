package parser

import (
	"fmt"
	"math"
	"strings"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// agentTriggerBodyYAML is the body an agent trigger shares between its two authoring forms: the standalone
// type: agent_trigger (which adds "agent") and the inline triggers: list on a type: agent, where the agent is
// implicit (ADR-0016). Holding it in one struct is what lets both forms build the identical spec via
// buildAgentTriggerSpec, so they never drift apart.
type agentTriggerBodyYAML struct {
	Source struct {
		Kind   string   `yaml:"kind"`
		Name   string   `yaml:"name"`
		Events []string `yaml:"events"`
		// Cron is the schedule a schedule-kind trigger fires on (validated by the reconciler). Empty for alert/report.
		Cron string `yaml:"cron"`
	} `yaml:"source"`
	Actor struct {
		// Attributes is the identity a started run acts as (its SecurityClaims user attributes), mirroring an alert's
		// query_for_attributes. Explicit attributes take precedence over user_id/user_email. Empty (and no user) means
		// the run carries no claims and fails closed; the reconciler validates the rest.
		Attributes map[string]any `yaml:"attributes"`
		// UserID / UserEmail name the user whose attributes the run acts as; the reconciler resolves them via admin,
		// exactly like an alert's for.user_id / for.user_email. Mutually exclusive.
		UserID    string `yaml:"user_id"`
		UserEmail string `yaml:"user_email"`
	} `yaml:"actor"`
	Input struct {
		Prompt  string            `yaml:"prompt"`
		Context map[string]string `yaml:"context"`
	} `yaml:"input"`
	Deduplication struct {
		Window string `yaml:"window"`
	} `yaml:"deduplication"`
}

// AgentTriggerYAML is the raw structure of an AgentTrigger resource defined in YAML (does not include common fields).
type AgentTriggerYAML struct {
	commonYAML           `yaml:",inline"` // Not accessed here, only setting it so we can use KnownFields for YAML parsing
	Agent                string           `yaml:"agent"`
	agentTriggerBodyYAML `yaml:",inline"`
}

// parseAgentTrigger parses an agent_trigger definition and adds the resulting resource to p.Resources.
// The resource is experimental (Kairos Act); the reconciler only validates it and never starts a run.
func (p *Parser) parseAgentTrigger(node *Node) error {
	// Parse YAML
	tmp := &AgentTriggerYAML{}
	err := p.decodeNodeYAML(node, true, tmp)
	if err != nil {
		return err
	}

	// Validate: a trigger must name the agent it starts.
	if strings.TrimSpace(tmp.Agent) == "" {
		return fmt.Errorf(`agent_trigger must set "agent"`)
	}

	// Build (and fully validate) the spec before inserting: the parser must not return an error after insertResource.
	spec, err := buildAgentTriggerSpec(tmp.Agent, tmp.agentTriggerBodyYAML)
	if err != nil {
		return err
	}

	// Depend on the referenced agent so the trigger reconciles after it and fails closed if it is missing.
	node.Refs = append(node.Refs, ResourceName{Kind: ResourceKindAgent, Name: tmp.Agent})

	// Insert the resource
	r, err := p.insertResource(ResourceKindAgentTrigger, node.Name, node.Paths, node.Tags, node.Refs...)
	if err != nil {
		return err
	}
	// NOTE: After calling insertResource, an error must not be returned. Any validation should be done before calling it.

	r.AgentTriggerSpec = spec

	return nil
}

// buildAgentTriggerSpec maps a validated trigger body to an AgentTriggerSpec for the given agent. It performs every
// validation that can fail (source kind, actor mode, deduplication window) so the caller can insert the resource only
// once it has succeeded: the parser's rule is that no error may be returned after insertResource. Both the standalone
// type: agent_trigger and the inline triggers: list on a type: agent build their spec through it, so the two authoring
// forms stay behaviourally identical (ADR-0016).
func buildAgentTriggerSpec(agent string, body agentTriggerBodyYAML) (*runtimev1.AgentTriggerSpec, error) {
	// Map the source kind. The dispatcher only activates a trigger whose kind matches the emitting resource.
	sourceKind, err := parseAgentTriggerSourceKind(body.Source.Kind)
	if err != nil {
		return nil, err
	}

	// Identity is declared at most one way, like an alert's for.{user_id|user_email|query_for_attributes}: mixing a
	// nominal user with explicit attributes would run under the attributes while the audit subject still reads as the
	// named user. None is allowed too (a service run that fails closed for want of claims).
	identityForms := 0
	if strings.TrimSpace(body.Actor.UserID) != "" {
		identityForms++
	}
	if strings.TrimSpace(body.Actor.UserEmail) != "" {
		identityForms++
	}
	if len(body.Actor.Attributes) > 0 {
		identityForms++
	}
	if identityForms > 1 {
		return nil, fmt.Errorf(`"actor" must set at most one of "user_id", "user_email" or "attributes"`)
	}

	// The actor's run-as attributes become the run's SecurityClaims (mirrors an alert's query_for_attributes). Convert
	// here so a value the Struct cannot represent fails at parse rather than silently dropping the identity.
	var actorAttributes *structpb.Struct
	if len(body.Actor.Attributes) > 0 {
		actorAttributes, err = structpb.NewStruct(body.Actor.Attributes)
		if err != nil {
			return nil, fmt.Errorf(`invalid "actor.attributes": %w`, err)
		}
	}

	// Parse deduplication.window (supports seconds, Go duration and ISO 8601 durations).
	var windowSeconds uint32
	if body.Deduplication.Window != "" {
		window, err := parseDuration(body.Deduplication.Window)
		if err != nil {
			return nil, fmt.Errorf(`invalid value %q for property "deduplication.window"`, body.Deduplication.Window)
		}
		if window < 0 {
			return nil, fmt.Errorf(`"deduplication.window" must not be negative, got %q`, body.Deduplication.Window)
		}
		// Round up: a sub-second window (e.g. "500ms") must not truncate to 0, which would drop deduplication.
		seconds := math.Ceil(window.Seconds())
		if seconds > math.MaxUint32 {
			return nil, fmt.Errorf(`"deduplication.window" is too large: %q`, body.Deduplication.Window)
		}
		windowSeconds = uint32(seconds)
	}

	spec := &runtimev1.AgentTriggerSpec{
		Agent: agent,
		Source: &runtimev1.AgentTriggerSource{
			Kind:   sourceKind,
			Name:   body.Source.Name,
			Events: body.Source.Events,
			Cron:   body.Source.Cron,
		},
		Actor: &runtimev1.AgentTriggerActor{
			Attributes: actorAttributes,
			UserId:     body.Actor.UserID,
			UserEmail:  body.Actor.UserEmail,
		},
		Input: &runtimev1.AgentTriggerInput{
			Prompt:  body.Input.Prompt,
			Context: body.Input.Context,
		},
	}
	if windowSeconds != 0 {
		spec.Deduplication = &runtimev1.AgentTriggerDeduplication{
			WindowSeconds: windowSeconds,
		}
	}
	return spec, nil
}

// parseAgentTriggerSourceKind maps the YAML "source.kind" to its enum. It is required.
func parseAgentTriggerSourceKind(s string) (runtimev1.AgentTriggerSourceKind, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "alert":
		return runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_ALERT, nil
	case "report":
		return runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_REPORT, nil
	case "schedule":
		return runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_SCHEDULE, nil
	case "":
		return runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_UNSPECIFIED, fmt.Errorf(`agent_trigger must set "source.kind"`)
	default:
		return runtimev1.AgentTriggerSourceKind_AGENT_TRIGGER_SOURCE_KIND_UNSPECIFIED, fmt.Errorf(`invalid value %q for property "source.kind"`, s)
	}
}
