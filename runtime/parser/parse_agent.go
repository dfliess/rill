package parser

import (
	"fmt"
	"math"
	"path"
	"strings"
	"text/template"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"gopkg.in/yaml.v3"
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
	Instructions string `yaml:"instructions"`
	// Prompt is the user turn sent when a run is started without one, for an agent launched by hand rather than
	// by a trigger. Optional: an agent that declares triggers already has their prompts to fall back on.
	Prompt string   `yaml:"prompt"`
	Tools  []string `yaml:"tools"`
	Limits struct {
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
	// Security declares who can see the agent (access), start a run of it (launch) and decide the approvals its
	// actions request (approve). Absent means only admins. It is deliberately NOT folded into spec_hash-relevant
	// behavior snapshots: policy is read live, so tightening it reaches approvals that are already waiting.
	Security *AgentSecurityYAML `yaml:"security"`
}

// AgentSecurityYAML is the security block of an agent. It embeds SecurityPolicyYAML so `access:` shares the
// grammar, validation and errors of a dashboard's security policy, and adds the two agent verbs on top:
//   - launch: who may start a run. Absent inherits access (a run is inert until an action passes the approval gate).
//   - approve: who may decide an approval the connector's posture marked as needed. Absent means only admins:
//     the asymmetry with launch is deliberate, approving touches the outside world.
//
// Both are templated boolean expressions over `.user`; approve may additionally reference `.action.tool`,
// `.action.connector` and `.run.actor`.
type AgentSecurityYAML struct {
	SecurityPolicyYAML `yaml:",inline"`
	Launch             string `yaml:"launch"`
	Approve            string `yaml:"approve"`
	// Execute is reserved for a future per-action execution gate. Setting it is an error rather than an unknown
	// key, so the reserved meaning cannot be squatted by accident.
	Execute yaml.Node `yaml:"execute"`
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

	// Parse and validate the security block (also before any insert: the parser must not error afterwards).
	var securityRules []*runtimev1.SecurityRule
	var launchExpr, approveExpr string
	if tmp.Security != nil {
		securityRules, launchExpr, approveExpr, err = parseAgentSecurity(tmp.Security)
		if err != nil {
			return err
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
	r.AgentSpec.Prompt = tmp.Prompt
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
	r.AgentSpec.SecurityRules = securityRules
	r.AgentSpec.LaunchExpression = launchExpr
	r.AgentSpec.ApproveExpression = approveExpr

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

// parseAgentSecurity validates an agent's security block and maps it to the spec's fields: the access rules
// (through SecurityPolicyYAML.Proto, so `access:` gets the same validation and deny-all-when-absent default as
// a dashboard policy) and the launch/approve gate expressions.
func parseAgentSecurity(sec *AgentSecurityYAML) (rules []*runtimev1.SecurityRule, launch, approve string, err error) {
	// The reserved key fails loudly instead of decoding as unknown, so its future meaning cannot be squatted.
	if !sec.Execute.IsZero() {
		return nil, "", "", fmt.Errorf(`invalid 'security': 'execute' is reserved and not implemented; use 'launch' to control who may start a run and the connector's approval posture to control which actions need a decision`)
	}

	// The embedded policy members that only make sense on a queryable resource are rejected rather than silently
	// accepted: an agent has no fields to include/exclude and no rows to filter, so their presence is a mistake.
	if sec.RowFilter != "" {
		return nil, "", "", fmt.Errorf(`invalid 'security': 'row_filter' is not supported on agents (an agent has no rows); only 'access', 'launch' and 'approve' apply`)
	}
	if len(sec.Include) > 0 || len(sec.Exclude) > 0 {
		return nil, "", "", fmt.Errorf(`invalid 'security': 'include'/'exclude' are not supported on agents (an agent has no fields); only 'access', 'launch' and 'approve' apply`)
	}
	if len(sec.Rules) > 0 {
		return nil, "", "", fmt.Errorf(`invalid 'security': 'rules' is not supported on agents; only 'access', 'launch' and 'approve' apply`)
	}

	// The approval-time context is only in scope while deciding an approval: in access/launch it would resolve
	// to nothing at request time, silently turning the policy into a constant, so referencing it is an error.
	for key, expr := range map[string]string{"access": sec.Access, "launch": sec.Launch} {
		refs, err := templateContextRefs(expr, "action", "run")
		if err != nil {
			return nil, "", "", fmt.Errorf(`invalid 'security': %q templating is not valid: %w`, key, err)
		}
		if len(refs) > 0 {
			return nil, "", "", fmt.Errorf(`invalid 'security': %q cannot reference %q: the action context exists only while deciding an approval, so it is available only in 'approve'`, key, "."+refs[0])
		}
	}

	// In approve, the action context is exactly the identity of the proposed action plus the run's actor. The
	// action's arguments are model-chosen, so exposing them would let the prompt influence who may approve; an
	// unknown field is rejected instead of resolving empty.
	approveRefs, err := templateContextRefs(sec.Approve, "action", "run")
	if err != nil {
		return nil, "", "", fmt.Errorf(`invalid 'security': 'approve' templating is not valid: %w`, err)
	}
	for _, ref := range approveRefs {
		switch ref {
		case "action.tool", "action.connector", "run.actor":
		default:
			return nil, "", "", fmt.Errorf(`invalid 'security': 'approve' cannot reference %q: beyond '.user', only '.action.tool', '.action.connector' and '.run.actor' are available (the action's arguments are model-chosen and must never select the approver)`, "."+ref)
		}
	}

	// Validate that launch and approve render and evaluate as booleans, exactly like `access:` (approve gets
	// sample values for its extra context, so a valid reference does not fail the dry run).
	if sec.Launch != "" {
		if err := validateAgentGateExpression("launch", sec.Launch, nil); err != nil {
			return nil, "", "", err
		}
	}
	if sec.Approve != "" {
		extra := map[string]any{
			"action": map[string]any{"tool": "dummy_tool", "connector": "dummy_connector"},
			"run":    map[string]any{"actor": "dummy_actor"},
		}
		if err := validateAgentGateExpression("approve", sec.Approve, extra); err != nil {
			return nil, "", "", err
		}
	}

	// Delegate access to the shared policy parser: it validates the template, validates that it evaluates as a
	// boolean, and emits the deny-all rule when `security:` is present without `access:`.
	rules, err = sec.SecurityPolicyYAML.Proto()
	if err != nil {
		return nil, "", "", err
	}
	return rules, sec.Launch, sec.Approve, nil
}

// validateAgentGateExpression dry-runs a gate expression with sample template data, mirroring how
// SecurityPolicyYAML.Proto validates `access:`.
func validateAgentGateExpression(key, expr string, extra map[string]any) error {
	data := validationTemplateData
	data.ExtraProps = extra
	tmp, err := ResolveTemplate(expr, data, false)
	if err != nil {
		return fmt.Errorf(`invalid 'security': %q templating is not valid: %w`, key, err)
	}
	_, err = EvaluateBoolExpression(tmp)
	if err != nil {
		return fmt.Errorf(`invalid 'security': %q expression error: %w`, key, err)
	}
	return nil
}

// templateContextRefs returns the field references in expr that address one of the given context roots (e.g.
// "action" matches both ".action" and ".action.tool"). It only parses the template, never executes it, so a
// reference outside the sample data is reported by name instead of failing as a template execution error.
func templateContextRefs(expr string, roots ...string) ([]string, error) {
	if expr == "" {
		return nil, nil
	}
	// The func map only needs the names to exist for parsing; none of the functions run here.
	funcMap := newFuncMap("", nil)
	for _, name := range []string{"configure", "dependency", "ref", "lookup", "env"} {
		funcMap[name] = func(...string) (string, error) { return "", nil }
	}
	t, err := template.New("").Funcs(funcMap).Option("missingkey=default").Parse(expr)
	if err != nil {
		return nil, err
	}
	var refs []string
	for _, v := range extractVariablesFromTemplate(t.Tree) {
		for _, root := range roots {
			if v == root || strings.HasPrefix(v, root+".") {
				refs = append(refs, v)
			}
		}
	}
	return refs, nil
}
