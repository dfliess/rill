package parser

import (
	"fmt"
	"math"
	"path"
	"strings"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
)

// AgentYAML is the raw structure of an Agent resource defined in YAML (does not include common fields).
type AgentYAML struct {
	commonYAML  `yaml:",inline"` // Not accessed here, only setting it so we can use KnownFields for YAML parsing
	DisplayName string           `yaml:"display_name"`
	Description string           `yaml:"description"`
	Model       struct {
		Connector string `yaml:"connector"`
		Name      string `yaml:"name"`
	} `yaml:"model"`
	Instructions string   `yaml:"instructions"`
	Tools        []string `yaml:"tools"`
	Limits       struct {
		MaxSteps uint32 `yaml:"max_steps"`
		Timeout  string `yaml:"timeout"`
	} `yaml:"limits"`
	MCP []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
		Auth struct {
			Secret string `yaml:"secret"`
		} `yaml:"auth"`
		Network struct {
			AllowedHosts []string `yaml:"allowed_hosts"`
		} `yaml:"network"`
		TrustReadOnlyHint bool     `yaml:"trust_read_only_hint"`
		Approval          string   `yaml:"approval"`
		RequireApproval   []string `yaml:"require_approval"`
		AutoApprove       []string `yaml:"auto_approve"`
	} `yaml:"mcp"`
	// Triggers declares when this agent runs, inline (ADR-0016). Each entry is an agent_trigger body minus "agent"
	// (which is this agent): the parser desugars it into a standalone AgentTrigger resource. The triggers live
	// outside AgentSpec on purpose, so editing when the agent fires never changes its spec_hash (what it does).
	Triggers []agentTriggerBodyYAML `yaml:"triggers"`
}

// parseAgent parses an agent definition and adds the resulting resource to p.Resources.
// The agent resource is experimental (Kairos Act); the reconciler only validates it and never executes it.
func (p *Parser) parseAgent(node *Node) error {
	// Parse YAML
	tmp := &AgentYAML{}
	err := p.decodeNodeYAML(node, true, tmp)
	if err != nil {
		return err
	}

	// Validate: instructions are the agent's system prompt and are required.
	if strings.TrimSpace(tmp.Instructions) == "" {
		return fmt.Errorf(`agents must set "instructions"`)
	}

	// Parse limits.timeout (supports seconds, Go duration and ISO 8601 durations).
	var timeoutSeconds uint32
	if tmp.Limits.Timeout != "" {
		timeout, err := parseDuration(tmp.Limits.Timeout)
		if err != nil {
			return fmt.Errorf(`invalid value %q for property "limits.timeout"`, tmp.Limits.Timeout)
		}
		// Round up: a sub-second timeout (e.g. "500ms") must not truncate to 0, which would drop the limit entirely.
		timeoutSeconds = uint32(math.Ceil(timeout.Seconds()))
	}

	// Validate the MCP connector blocks structurally. Deep validation (transport, SSRF posture, absolute https
	// URL) belongs to the reconciler via mcpconn.ParseConnectorConfig; the parser only enforces the shape:
	// a name and url are required, and connector names must be unique so tool namespacing is unambiguous.
	seenMCP := make(map[string]bool, len(tmp.MCP))
	for i := range tmp.MCP {
		m := &tmp.MCP[i]
		if strings.TrimSpace(m.Name) == "" {
			return fmt.Errorf(`each "mcp" connector must set "name"`)
		}
		// A dot in the connector name breaks tool namespacing: effective tool names are "mcp.<connector>.<tool>", and the
		// gateway splits on the first dot to recover the connector. A name like "jira.prod" would produce
		// "mcp.jira.prod.create_issue", parsed as connector "jira" — a silent misroute. Reject at definition time.
		if strings.Contains(m.Name, ".") {
			return fmt.Errorf(`mcp connector name %q must not contain dots (used as namespace separator in tool names)`, m.Name)
		}
		if strings.TrimSpace(m.URL) == "" {
			return fmt.Errorf(`mcp connector %q must set "url"`, m.Name)
		}
		if seenMCP[m.Name] {
			return fmt.Errorf(`duplicate mcp connector name %q`, m.Name)
		}
		seenMCP[m.Name] = true

		// Validate the approval posture and the glob patterns so a typo is a definition error, not a silent
		// misconfiguration that changes which actions run unsupervised.
		switch m.Approval {
		case "", "manual", "auto":
		default:
			return fmt.Errorf(`mcp connector %q: "approval" must be "manual" or "auto", got %q`, m.Name, m.Approval)
		}
		for _, pats := range [][]string{m.RequireApproval, m.AutoApprove} {
			for _, pat := range pats {
				if _, err := path.Match(pat, "x"); err != nil {
					return fmt.Errorf(`mcp connector %q: invalid tool pattern %q: %w`, m.Name, pat, err)
				}
			}
		}
	}

	// Desugar the inline triggers (ADR-0016) into standalone AgentTrigger specs. Build and name them, and dry-run their
	// resource names for collisions, all BEFORE inserting anything: the parser must not return an error once a resource
	// is inserted, and emitting the agent plus its triggers must be all-or-nothing rather than a half-built set.
	triggerSpecs := make([]*runtimev1.AgentTriggerSpec, len(tmp.Triggers))
	triggerNames := make([]string, len(tmp.Triggers))
	for i := range tmp.Triggers {
		spec, err := buildAgentTriggerSpec(node.Name, tmp.Triggers[i])
		if err != nil {
			return err
		}
		// Synthetic name: "<agent>__trigger_<i>". The double underscore keeps it distinct from a hand-written trigger
		// while staying a valid resource name (names only need to be collision-free, not to match a charset).
		name := fmt.Sprintf("%s__trigger_%d", node.Name, i)
		if err := p.insertDryRun(ResourceKindAgentTrigger, name); err != nil {
			return fmt.Errorf("inline trigger %d desugars to %q, which collides with an existing resource: %w", i, name, err)
		}
		triggerSpecs[i] = spec
		triggerNames[i] = name
	}

	// Add a ref to the model connector. It's dropped later if it's not a project resource
	// (AI connectors are commonly configured in rill.yaml instead of as a connector resource).
	if tmp.Model.Connector != "" {
		node.Refs = append(node.Refs, ResourceName{Kind: ResourceKindConnector, Name: tmp.Model.Connector})
	}

	// Insert the resource
	r, err := p.insertResource(ResourceKindAgent, node.Name, node.Paths, node.Tags, node.Refs...)
	if err != nil {
		return err
	}
	// NOTE: After calling insertResource, an error must not be returned. Any validation should be done before calling it.

	r.AgentSpec.DisplayName = tmp.DisplayName
	r.AgentSpec.Description = tmp.Description
	r.AgentSpec.ModelConnector = tmp.Model.Connector
	r.AgentSpec.ModelName = tmp.Model.Name
	r.AgentSpec.Instructions = tmp.Instructions
	r.AgentSpec.Tools = tmp.Tools
	if tmp.Limits.MaxSteps != 0 || timeoutSeconds != 0 {
		r.AgentSpec.Limits = &runtimev1.AgentLimits{
			MaxSteps:       tmp.Limits.MaxSteps,
			TimeoutSeconds: timeoutSeconds,
		}
	}
	if len(tmp.MCP) > 0 {
		r.AgentSpec.Mcp = make([]*runtimev1.MCPConnector, len(tmp.MCP))
		for i := range tmp.MCP {
			m := &tmp.MCP[i]
			r.AgentSpec.Mcp[i] = &runtimev1.MCPConnector{
				Name:              m.Name,
				Url:               m.URL,
				AuthSecret:        m.Auth.Secret,
				AllowedHosts:      m.Network.AllowedHosts,
				TrustReadOnlyHint: m.TrustReadOnlyHint,
				Approval:          m.Approval,
				RequireApproval:   m.RequireApproval,
				AutoApprove:       m.AutoApprove,
			}
		}
	}

	// Emit the desugared AgentTrigger resources. They are NOT added to AgentSpec (so the agent's spec_hash is untouched
	// by construction); each is a sibling resource that refs the agent and carries the agent's own paths and tags. The
	// names were dry-run checked above, so insertResource cannot collide here and the "no error after insertResource"
	// rule holds.
	agentRef := ResourceName{Kind: ResourceKindAgent, Name: node.Name}
	for i, spec := range triggerSpecs {
		tr, err := p.insertResource(ResourceKindAgentTrigger, triggerNames[i], node.Paths, node.Tags, agentRef)
		if err != nil {
			return err
		}
		tr.AgentTriggerSpec = spec
	}

	return nil
}
