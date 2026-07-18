package ai

import (
	"context"
	"errors"
	"fmt"
	"slices"
)

// ErrAgentNotFound is returned by an AgentDefinitionProvider when no agent matches the requested name.
var ErrAgentNotFound = errors.New("agent not found")

// AgentSnapshot is an immutable description of a dynamic agent.
//
// Unlike the built-in agents (which are Go types registered in NewRunner), a dynamic agent is defined
// entirely by data: its instructions, the model to run it on, and the set of tools it may call. A snapshot
// is resolved once at the start of a run and is not expected to change while that run executes, which keeps
// the agent's authority fixed and auditable.
type AgentSnapshot struct {
	// Name is the stable identifier of the agent (used to look it up and to label its call in a session).
	Name string
	// SpecHash is the stable hash of the agent's effective spec at the time the snapshot was taken. It binds a run
	// to the exact agent version it executed (§8.3): the executor records it on the run so an audit can tell which
	// definition produced an action even after the agent is edited. The catalog-backed provider fills it from
	// AgentState.SpecHash; the in-memory StaticAgentProvider leaves it empty.
	SpecHash string
	// DisplayName is a human-friendly label for the agent.
	DisplayName string
	// Instructions is the agent-specific system prompt. It is combined with the project's ai_instructions at runtime.
	Instructions string
	// ModelConnector and ModelName identify the LLM to run the agent on.
	//
	// NOTE: RunDynamicAgent currently executes on the Session's configured LLM (mirroring the built-in agents,
	// which don't select a model). Per-agent model routing is left to a later integration; these fields carry the
	// intent so the resolved snapshot is complete.
	ModelConnector string
	ModelName      string
	// Tools is the exact set of built-in (analytical) tool names the agent may call. It is the upper bound on the
	// agent's authority over Rill's own tools: the model never sees a built-in tool outside this list, and a
	// proposed call outside it is rejected, not executed.
	//
	// It deliberately does NOT include the agent's MCP tools. MCP tools are not declared here and are not part of
	// declarableAgentTools: they are discovered live at run start from MCPConnectors (their schemas pinned by that
	// discovery) and added to the callable set only for the duration of the run. Keeping them out of Tools is what
	// lets IsDeclarableAgentTool stay a fixed allowlist of read-only analytical tools.
	Tools []string
	// MCPConnectors are the outbound MCP servers this agent may reach. At run start the executor connects to each,
	// discovers its tools (pinning their schemas), and offers them to the model under namespaced names
	// (mcp.<name>.<tool>). They live only for the run and are never added to Tools.
	MCPConnectors []MCPConnector
	// MaxSteps bounds the number of model/tool iterations. Zero means use the default.
	MaxSteps int
	// TimeoutSeconds bounds the wall-clock duration of a run. Zero means no agent-level timeout is applied here
	// (the Session still applies its per-LLM-request timeout).
	TimeoutSeconds int
}

// MCPConnector is the resolved, immutable description of one outbound MCP server an agent may reach. AuthSecret
// is a reference to the credential (a variable/secret name), never the literal: the executor resolves it
// server-side at run start and hands the token to the mcpconn client below its egress guard.
type MCPConnector struct {
	// Name is the local connector name; its discovered tools are namespaced as mcp.<name>.<tool>.
	Name string
	// URL is the absolute https endpoint of the remote MCP server.
	URL string
	// AuthSecret is the name of the variable/secret holding the bearer token. Empty means no bearer is sent.
	AuthSecret string
	// AllowedHosts is the egress allowlist for this connector.
	AllowedHosts []string
	// TrustReadOnlyHint opts the connector in to treating a tool's server-advertised readOnlyHint as meaningful.
	TrustReadOnlyHint bool
	// Approval is the connector's approval posture for its actions: "auto" auto-approves them (subject to
	// RequireApproval), anything else (the default) requires human approval. It never affects a trusted read-only
	// tool, which runs inline. RequireApproval and AutoApprove are tool-name glob exceptions to the posture.
	Approval        string
	RequireApproval []string
	AutoApprove     []string
}

// declarableAgentTools is the allowlist of tools a dynamic (YAML-defined) agent may declare in v1.
// Per the Act alignment decisions, a dynamic agent may use only Rill's read-only analytical tools.
// Development and administration tools (file access, query_sql, create_chart, project_status, bucket
// tools, navigate) and the built-in agent tools (router/analyst/developer/feedback, whose nested loops
// would escape this agent's tool boundary) are deliberately excluded.
// TODO(act, phase 1): when the outbound MCP client lands, namespaced MCP tool names (mcp.<connector>.<tool>)
// become declarable through the tool gateway, each gated by approval policy.
var declarableAgentTools = map[string]bool{
	ListMetricsViewsName:        true,
	GetMetricsViewName:          true,
	QueryMetricsViewName:        true,
	QueryMetricsViewSummaryName: true,
	GetCanvasName:               true,
}

// IsDeclarableAgentTool reports whether a dynamic agent may declare the named tool.
func IsDeclarableAgentTool(name string) bool { return declarableAgentTools[name] }

// AgentDefinitionProvider resolves dynamic agent definitions for an instance.
//
// It is the seam that lets the runtime discover agents that are not compiled-in built-ins. The in-memory
// StaticAgentProvider below is a stub for tests; a later frontend replaces it with one backed by the project's
// resource catalog without any change to the execution path.
type AgentDefinitionProvider interface {
	// GetAgent returns the snapshot for the named agent, or ErrAgentNotFound if it does not exist.
	GetAgent(ctx context.Context, instanceID, name string) (*AgentSnapshot, error)
	// ListAgents returns all agents available for the instance.
	ListAgents(ctx context.Context, instanceID string) ([]*AgentSnapshot, error)
}

// StaticAgentProvider is an in-memory AgentDefinitionProvider backed by a fixed map of snapshots keyed by name.
//
// It is instance-agnostic: the same set of agents is returned for every instanceID. This is intentional for a
// stub used by tests and local wiring; the real provider keys off the instance's resource catalog.
type StaticAgentProvider struct {
	agents map[string]*AgentSnapshot
}

var _ AgentDefinitionProvider = (*StaticAgentProvider)(nil)

// NewStaticAgentProvider creates a StaticAgentProvider from the given snapshots. Later snapshots with the same
// name override earlier ones. It stores copies so a caller that retains and mutates its argument cannot alter the
// provider's definitions.
func NewStaticAgentProvider(snapshots ...*AgentSnapshot) *StaticAgentProvider {
	agents := make(map[string]*AgentSnapshot, len(snapshots))
	for _, s := range snapshots {
		agents[s.Name] = cloneSnapshot(s)
	}
	return &StaticAgentProvider{agents: agents}
}

func (p *StaticAgentProvider) GetAgent(_ context.Context, _, name string) (*AgentSnapshot, error) {
	s, ok := p.agents[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrAgentNotFound, name)
	}
	// Hand out a copy, not the stored pointer: a snapshot is immutable authority, so a consumer that mutates its
	// Tools must not alter the provider's stored definition (or any other run resolving the same agent).
	return cloneSnapshot(s), nil
}

func (p *StaticAgentProvider) ListAgents(_ context.Context, _ string) ([]*AgentSnapshot, error) {
	res := make([]*AgentSnapshot, 0, len(p.agents))
	for _, s := range p.agents {
		res = append(res, cloneSnapshot(s))
	}
	return res, nil
}

// cloneSnapshot returns a deep-enough copy of a snapshot that mutating the copy (including its Tools slice or its
// MCP connectors) cannot reach the original. All other fields are value types.
func cloneSnapshot(s *AgentSnapshot) *AgentSnapshot {
	snap := *s
	snap.Tools = slices.Clone(s.Tools)
	if s.MCPConnectors != nil {
		snap.MCPConnectors = make([]MCPConnector, len(s.MCPConnectors))
		for i := range s.MCPConnectors {
			c := s.MCPConnectors[i] // explicit copy: the range value is mutated below and must not alias the original
			c.AllowedHosts = slices.Clone(c.AllowedHosts)
			c.RequireApproval = slices.Clone(c.RequireApproval)
			c.AutoApprove = slices.Clone(c.AutoApprove)
			snap.MCPConnectors[i] = c
		}
	}
	return &snap
}
