package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/rilldata/rill/runtime/act/mcpconn"
	"github.com/rilldata/rill/runtime/pkg/secretscrub"
)

// newMCPClient constructs the outbound MCP client for a connector. It is a package variable so tests can override
// it, e.g. to set Options.AllowPrivateNetworks for an in-process loopback server. Production never sets that flag:
// a tenant connector is always fail-closed to public https hosts.
var newMCPClient = func(cfg *mcpconn.ConnectorConfig, token string) (mcpconn.MCPClient, error) {
	return mcpconn.NewClient(cfg, mcpconn.Options{BearerToken: token})
}

// dynamicAgentDefaultMaxSteps is the fallback iteration budget when a snapshot leaves MaxSteps unset.
const dynamicAgentDefaultMaxSteps = 10

// InjectedActionResultTool tags the message that carries an executed action's result back into a resumed segment. The
// message feeds the model (so it can react to the outcome and compose a closing) but is hidden in the chat UI: it is an
// internal prompt, not something the user said, so the chat frontend skips messages with this tool tag and the user
// reads only the model's own closing answer.
const InjectedActionResultTool = "act.action_result"

// DynamicAgent runs a model/tool loop for an agent that is defined by data (an AgentSnapshot) rather than by a
// compiled-in Go type. It mirrors the pattern of the built-in AnalystAgent: assemble the messages, then delegate
// the loop to Session.Complete. The difference is that its authority is bounded by the snapshot:
//
//   - The model is only shown the snapshot's tools, intersected with the tools the caller can actually access.
//   - A tool the model proposes outside that set is rejected, not executed (fail-closed, via RestrictToolCalls).
//   - The loop is cut off after MaxSteps iterations.
type DynamicAgent struct {
	Snapshot *AgentSnapshot

	// proposed holds every write (non-read-only) tool call the loop captured instead of executing (see ProposedAction).
	// It is empty for a pure read-only investigation. The capturing handler that compileProposedMCPTool installs appends
	// to it; Run reads it after the loop. The loop is single-goroutine, so no synchronization is needed. The slice
	// preserves the model's emission order: the first element is the first tool call the model proposed, and the executor
	// governs and executes them in that order.
	proposed []*ProposedAction
}

// ProposedAction is a write tool call the agent chose to make, captured instead of executed. In the governed action
// model a tool the connector does not vouch for as read-only never runs inline: the loop records the intended call
// here and ends, and the durable workflow routes it through policy, human approval and an idempotent, exactly-once
// execution. Read-only tools (a connector that trusts the server's readOnlyHint) still run inline, since a read needs
// no approval.
type ProposedAction struct {
	// ToolCallID is the model's identifier for this call: the message-tree ID of the assistant tool-call the loop
	// captured. It is the action's identity within the run (the ledger keys on it and the idempotency key derives from
	// it), and it is the parent the durable loop attaches the executed result to on resume, so the model sees the result
	// paired with its own call. Two calls in one run get distinct IDs, which is what lets a run propose several actions.
	ToolCallID string
	// Connector is the MCP connector the tool belongs to.
	Connector string
	// Tool is the effective, dotted tool name (mcp.<connector>.<rawName>): the identity the action gateway keys on.
	Tool string
	// Args are the arguments the model produced, carrying secret references (never resolved secrets).
	Args map[string]any
	// SchemaHash pins the tool schema discovered at run start, so the post-approval execution re-validates against the
	// exact schema the proposal was built from.
	SchemaHash string
	// Summary is a short human-readable description of the intended effect, for the approval inbox.
	Summary string
}

// DynamicAgentResult is the outcome of a dynamic agent run (one segment of a durable, segmented run).
type DynamicAgentResult struct {
	Response string `json:"response"`
	Agent    string `json:"agent"`
	// Proposed is the list of write actions the agent proposed in this segment (captured, not executed), empty for a pure
	// investigation or a segment that produced a final answer. A non-empty Proposed is the pause point: the durable
	// workflow governs and executes all of them (in the order the model emitted them), then resumes the run so the model
	// sees all the results.
	Proposed []*ProposedAction `json:"-"`
}

// InjectedResult is the outcome of a governed action, fed back into the loop so the model can react to the effect it
// proposed: on resume the model sees the result and either proposes another action or composes a final answer.
type InjectedResult struct {
	// ToolCallID identifies the proposed call this result answers (the captured tool-call message ID).
	ToolCallID string
	// Tool is the effective tool name, for a human-readable framing of the result turn.
	Tool string
	// Message is the tool's text result (already redacted and capped by the gateway); Output is its structured result.
	Message string
	Output  map[string]any
	// IsError reports that the action failed, so the model is told the effect did not happen rather than that it did.
	IsError bool
}

// RunDynamicAgent resolves the named agent from the provider and runs it over an existing Session.
//
// Signature rationale: the built-in agents are invoked as tools on a live Session (s.CallTool(...)), which already
// carries the caller's claims, the LLM handle, the catalog, and the message tree. A dynamic agent is not a
// registered tool, so it cannot be reached through that path by name; a function that takes the *Session mirrors the
// same call site while keeping the provider — an external collaborator that resolves the definition — an explicit
// argument rather than something reached through the session. It returns *DynamicAgentResult to match the shape of
// AnalystAgentResult / RouterAgentResult.
func RunDynamicAgent(ctx context.Context, s *Session, provider AgentDefinitionProvider, agentName, prompt string) (*DynamicAgentResult, error) {
	snapshot, err := provider.GetAgent(ctx, s.InstanceID(), agentName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve agent %q: %w", agentName, err)
	}
	agent := &DynamicAgent{Snapshot: snapshot}
	return agent.Run(ctx, s, prompt, nil)
}

// Run executes one segment of the agent's model/tool loop over the session and returns its response and the actions it
// paused on (if any). resume is nil for the first segment (a fresh run seeded by prompt); for a resumed segment it
// carries the executed actions' results, and the conversation is reconstructed from the session's persisted message
// tree, so the caller must hand Run the same session the prior segment ran on (in production the workflow reopens it
// by ID across the pause). This mirrors how the built-in agents rebuild model context from the tree (router_agent).
func (a *DynamicAgent) Run(ctx context.Context, s *Session, prompt string, resume []*InjectedResult) (*DynamicAgentResult, error) {
	if a.Snapshot == nil {
		return nil, fmt.Errorf("dynamic agent has no snapshot")
	}

	// Attach the session to the context. Tool CheckAccess and Session.Complete both resolve the session via
	// GetSession(ctx); the built-in agents get it for free because they run inside a Call, but a dynamic agent is
	// invoked directly on the session, so we set it here.
	ctx = WithSession(ctx, s)

	// Apply the agent-level wall-clock timeout, if the snapshot sets one. The snapshot carries it as intent; this is
	// where it takes effect, bounding the whole model/tool loop (Session.Complete still applies its per-request timeout).
	if a.Snapshot.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(a.Snapshot.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// Resolve the callable tool set: the snapshot's tools intersected with what the caller can access.
	// Tools the caller cannot access are silently dropped (they are simply not part of this agent's authority);
	// a tool named in the snapshot that does not exist at all is a definition error and fails the run.
	tools, err := a.allowedTools(ctx, s)
	if err != nil {
		return nil, err
	}

	// Discover and register the agent's MCP tools for the duration of this run. They are resolved live from the
	// snapshot's connectors (not from the declarable analytical set), their schemas pinned at discovery, and they
	// exist only while this run executes. A connector that fails setup fails the run (fail-closed).
	mcpNames, closeMCP, err := a.setupMCPTools(ctx, s)
	if err != nil {
		return nil, err
	}
	defer closeMCP()
	tools = append(tools, mcpNames...)

	maxSteps := a.Snapshot.MaxSteps
	if maxSteps <= 0 {
		maxSteps = dynamicAgentDefaultMaxSteps
	}

	// Build the model context for this segment. A fresh run opens with the system prompt and the user's instruction. A
	// resumed segment reconstructs the conversation from the session's persisted message tree (Rill's durable substrate:
	// every call and result is a parented message flushed to the catalog) and injects the executed action's result as the
	// turn the model reacts to. Persisting the injected turn before reconstructing folds it into the model context and
	// keeps it for the UI and any later segment; the system prompt is not persisted (it is regenerated per segment), so
	// it is prepended fresh. This is the pattern the built-in agents use to rebuild context from the tree (router_agent).
	var messages []*aiv1.CompletionMessage
	if len(resume) == 0 {
		messages = []*aiv1.CompletionMessage{
			NewTextCompletionMessage(RoleSystem, a.systemPrompt(s)),
			NewTextCompletionMessage(RoleUser, prompt),
		}
		s.AddMessage(&AddMessageOptions{
			Role:        RoleUser,
			Type:        MessageTypeText,
			ContentType: MessageContentTypeText,
			Content:     prompt,
		})
	} else {
		for _, r := range resume {
			s.AddMessage(&AddMessageOptions{
				Role:        RoleUser,
				Type:        MessageTypeText,
				Tool:        InjectedActionResultTool,
				ContentType: MessageContentTypeText,
				Content:     resumeResultMessage(r),
			})
		}
		messages = []*aiv1.CompletionMessage{NewTextCompletionMessage(RoleSystem, a.systemPrompt(s))}
		messages = append(messages, s.NewCompletionMessages(s.MessagesWithResults(FilterByRoot()))...)
	}

	// Delegate the loop to Session.Complete, which handles truncation, timeouts, telemetry and panic recovery.
	// RestrictToolCalls enforces the snapshot boundary on execution; passing the same tool set as opts.Tools keeps the
	// advertised set and the enforced set identical. UnwrapCall runs the loop in the session's own scope (as the built-in
	// leaf agents do), so every call and result the loop appends lands at the session root where the resume path above
	// reconstructs it from.
	var response string
	err = s.Complete(ctx, a.Snapshot.Name, &response, &CompleteOptions{
		Messages:          messages,
		Tools:             tools,
		MaxIterations:     maxSteps,
		RestrictToolCalls: true,
		UnwrapCall:        true,
		// Pause the loop the instant a governed write is captured, so the model does not close over the pending proposals
		// (no premature answer before approval). The workflow governs and executes the actions, then resumes this loop.
		PauseAfterToolCall: func() bool { return len(a.proposed) > 0 },
	})
	if err != nil {
		return nil, err
	}

	// Persist the segment's final response as an assistant turn. UnwrapCall means Complete returns the final text via out
	// but does not itself write it to the tree; persisting it here gives the next segment's reconstruction (and the UI)
	// what the model said, not just its tool calls.
	if response != "" {
		s.AddMessage(&AddMessageOptions{
			Role:        RoleAssistant,
			Type:        MessageTypeText,
			ContentType: MessageContentTypeText,
			Content:     response,
		})
	}

	return &DynamicAgentResult{Response: response, Agent: a.Snapshot.Name, Proposed: a.proposed}, nil
}

// resumeResultMessage frames an executed action's outcome as the turn the model reacts to on resume: it states whether
// the effect happened and carries the tool's result, so the model can compose a closing answer (or propose a next
// action). The result text is already redacted and capped by the gateway; a structured-only result is rendered as JSON.
func resumeResultMessage(r *InjectedResult) string {
	if r == nil {
		return "The proposed action was processed. Continue."
	}
	outcome := "was approved and executed"
	if r.IsError {
		outcome = "could not be completed"
	}
	result := r.Message
	if result == "" && len(r.Output) > 0 {
		if b, err := json.Marshal(r.Output); err == nil {
			result = string(b)
		}
	}
	if result == "" {
		result = "(no result returned)"
	}
	tool := r.Tool
	if tool == "" {
		tool = "the requested action"
	}
	// This message is both fed to the model and persisted as a visible turn (Option 1), so it reads as a natural update
	// rather than developer-facing steering. "Reply to the user with this outcome" nudges the model to compose a final
	// answer instead of re-proposing the same action, without leaking an engineering instruction into the chat.
	return fmt.Sprintf("The action %q %s. Result: %s. Reply to the user with this outcome.", tool, outcome, result)
}

// allowedTools returns the snapshot's tools that both exist and pass CheckAccess for the current caller.
func (a *DynamicAgent) allowedTools(ctx context.Context, s *Session) ([]string, error) {
	tools := make([]string, 0, len(a.Snapshot.Tools))
	for _, name := range a.Snapshot.Tools {
		// Defense in depth for the StaticAgentProvider path (tests, local wiring), which does not pass through the
		// reconciler that normally rejects non-declarable tools. A snapshot naming a tool a dynamic agent may not
		// use is a definition/security error, not something to silently skip.
		if !IsDeclarableAgentTool(name) {
			return nil, fmt.Errorf("agent %q references tool %q which dynamic agents may not use", a.Snapshot.Name, name)
		}
		tool, ok := s.Tool(name)
		if !ok {
			return nil, fmt.Errorf("agent %q references unknown tool %q", a.Snapshot.Name, name)
		}
		ok, err := tool.CheckAccess(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to check access for tool %q: %w", name, err)
		}
		if !ok {
			continue
		}
		tools = append(tools, name)
	}
	return tools, nil
}

// setupMCPTools connects to each of the snapshot's MCP connectors, discovers its tools (pinning their schemas at
// this point, immutable for the run), and registers each as a per-run tool on the session. It returns the effective
// tool names to add to the callable set and a close function the caller must defer to tear down the connections.
//
// It is fail-closed: any connector that cannot resolve its secret, parse its config, or list its tools fails the
// whole run, and anything already opened is closed before returning. The bearer token is never logged or placed in
// an error message.
func (a *DynamicAgent) setupMCPTools(ctx context.Context, s *Session) (names []string, closeFn func(), err error) {
	if len(a.Snapshot.MCPConnectors) == 0 {
		return nil, func() {}, nil
	}

	// Resolve the project's variables/secrets once. A connector's auth secret is a reference (a variable name); the
	// literal is resolved here, server-side, and handed to the client below its egress guard.
	instance, err := s.runner.Runtime.Instance(ctx, s.InstanceID())
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to load instance for mcp connectors: %w", err)
	}
	vars := instance.ResolveVariables(true)

	var clients []mcpconn.MCPClient
	var registered []string
	closeFn = func() {
		// Scope the discovered tools to this run: remove them from the shared *BaseSession (which sub-calls reach via
		// WithParent) before tearing down connections, so a later sub-agent or run cannot reach a leftover MCP tool.
		for _, n := range registered {
			s.UnregisterRunTool(n)
		}
		for _, c := range clients {
			if closer, ok := c.(interface{ Close() error }); ok {
				_ = closer.Close()
			}
		}
	}
	// On any error below, close what was already opened and return the wrapped error; the call sites hand back a
	// no-op close so the caller's deferred close is harmless.
	fail := func(format string, args ...any) error {
		closeFn()
		return fmt.Errorf(format, args...)
	}

	// seen guards against two effective names sanitizing to the same model-facing name within a run (rare, but it
	// would let one tool shadow another). It spans all connectors, since the connector prefix is part of the name.
	seen := make(map[string]bool)

	for i := range a.Snapshot.MCPConnectors {
		conn := a.Snapshot.MCPConnectors[i]

		var token string
		if conn.AuthSecret != "" {
			// ResolveVariables(true) canonicalizes keys to lower case (Rill variable names are case-insensitive), so
			// look the secret up by its lowercased name rather than the literal as written in the connector.
			token = vars[strings.ToLower(conn.AuthSecret)]
			if token == "" {
				return nil, func() {}, fail("mcp connector %q: secret %q is empty", conn.Name, conn.AuthSecret)
			}
		}

		props := map[string]any{
			"transport": mcpconn.TransportStreamableHTTP,
			"url":       conn.URL,
			"auth": map[string]any{
				"secret": conn.AuthSecret,
			},
			"network": map[string]any{
				"allowed_hosts":  conn.AllowedHosts,
				"private_ranges": mcpconn.PrivateRangesDeny,
			},
			"trust_read_only_hint": conn.TrustReadOnlyHint,
		}
		cfg, err := mcpconn.ParseConnectorConfig(conn.Name, props)
		if err != nil {
			return nil, func() {}, fail("mcp connector %q: %w", conn.Name, err)
		}

		client, err := newMCPClient(cfg, token)
		if err != nil {
			return nil, func() {}, fail("mcp connector %q: %w", conn.Name, err)
		}
		clients = append(clients, client)

		// ListTools is the discovery snapshot: it pins each tool's schema hash for the whole run, so a later schema
		// change on the server is caught by CallTool's per-call re-check rather than silently accepted.
		remoteTools, err := client.ListTools(ctx)
		if err != nil {
			return nil, func() {}, fail("mcp connector %q: discover tools: %w", conn.Name, err)
		}
		for j := range remoteTools {
			// The model-facing name is a sanitized form of the effective name (the raw mcp.<connector>.<tool> has
			// dots, which the LLM providers' function-name constraint rejects). Dispatch still uses the dotted name.
			displayName := sanitizeMCPToolName(remoteTools[j].Name)
			if seen[displayName] {
				return nil, func() {}, fail("mcp connector %q: tool name %q collides with another tool after sanitization", conn.Name, displayName)
			}
			seen[displayName] = true

			// A tool runs inline only when the connector trusts the server's read-only annotation AND the tool declares
			// it: a pure read needs no approval. Every other tool is a governed write, captured as a proposal (never
			// executed in the loop) and routed through human approval by the durable workflow.
			var tool *CompiledTool
			if conn.TrustReadOnlyHint && remoteTools[j].Hints.ReadOnly {
				// Hand the read-only tool the connector's resolved bearer token so its handler can scrub the token from
				// the result before it reaches the model (a read runs inline; its output is not gated by the gateway).
				tool = compileRemoteMCPTool(client, displayName, remoteTools[j], token)
			} else {
				tool = a.compileProposedMCPTool(displayName, remoteTools[j])
			}
			s.RegisterRunTool(tool)
			names = append(names, displayName)
			registered = append(registered, displayName)
		}
	}

	return names, closeFn, nil
}

// compileRemoteMCPTool adapts a discovered remote MCP tool into a CompiledTool the session can advertise and call.
// name is the model-facing (sanitized) tool name; rt.Name is the original dotted effective name used for dispatch.
// The connector's declaration is the authorization to use the tool, so CheckAccess is always true here; egress
// safety and the schema pin are enforced by the mcpconn client. The tool's pinned schema hash is captured in the
// handler so the exact schema discovered at run start is the one re-validated on every call.
//
// token is the connector's resolved bearer credential (empty if the connector has none). A read-only tool runs inline
// and its result flows straight into the model's context (and thus into ai_messages and logs), unlike a governed write
// whose output the gateway redacts. So the handler scrubs the token from the result before returning it: a misbehaving
// or verbose remote server that echoes the credential back (e.g. reflecting the Authorization header in an error) must
// not leak it into the model context in the clear.
func compileRemoteMCPTool(client mcpconn.MCPClient, name string, rt mcpconn.RemoteTool, token string) *CompiledTool {
	spec := &mcp.Tool{
		Name:        name,
		Description: rt.Description,
	}
	// The server's input schema passes through as raw JSON; AsProto marshals it verbatim (it is not a
	// *jsonschema.Schema, so the OpenAI empty-properties shortcut does not apply).
	if len(rt.InputSchema) > 0 {
		spec.InputSchema = rt.InputSchema
	}

	return &CompiledTool{
		Name:        name,
		Spec:        spec,
		CheckAccess: func(context.Context) (bool, error) { return true, nil },
		UnmarshalArgs: func(content string) (any, error) {
			if strings.TrimSpace(content) == "" {
				return map[string]any{}, nil
			}
			var args map[string]any
			if err := json.Unmarshal([]byte(content), &args); err != nil {
				return nil, err
			}
			return args, nil
		},
		UnmarshalResult: func(content string) (any, error) {
			return json.RawMessage(content), nil
		},
		JSONHandler: func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
			res, err := client.CallTool(ctx, rt.Name, rt.SchemaHash, input)
			if err != nil {
				return nil, err
			}
			// Scrub the connector's resolved credential from everything this handler returns to the model. secretscrub
			// skips an empty token, so a connector with no credential incurs no redaction. This mirrors the governed
			// write path, where the gateway redacts the same resolved secret from every tool output (runtime/act).
			secrets := []string{token}
			if res.IsError {
				// A tool-level error (the tool ran and reported failure): surface it so the loop captures it as a tool
				// result the model can react to, rather than crashing the run. res.Text can echo a header, so scrub it.
				return nil, fmt.Errorf("remote tool %q reported an error: %s", rt.Name, secretscrub.String(res.Text, secrets))
			}
			if len(res.StructuredContent) > 0 {
				// Walk the decoded structure so a token is caught wherever the server placed it (a nested value or a
				// key), matching the gateway's RedactMap. If the payload does not decode as JSON (it should, the client
				// marshals it), fall back to scrubbing the raw bytes so nothing passes through unredacted.
				var decoded any
				if err := json.Unmarshal(res.StructuredContent, &decoded); err != nil {
					return json.RawMessage(secretscrub.String(string(res.StructuredContent), secrets)), nil
				}
				out, err := json.Marshal(secretscrub.Value(decoded, secrets))
				if err != nil {
					return nil, err
				}
				return out, nil
			}
			// Wrap plain text output so the result is always a valid JSON object.
			out, err := json.Marshal(map[string]string{"text": secretscrub.String(res.Text, secrets)})
			if err != nil {
				return nil, err
			}
			return out, nil
		},
	}
}

// compileProposedMCPTool adapts a write (non-read-only) MCP tool into a CompiledTool that, when the model calls it,
// captures the intended call as a ProposedAction and returns a placeholder result instead of executing it. Every
// governed call in a turn is captured (appended to the agent's proposed slice), preserving the model's emission order.
// The real, idempotent execution happens post-approval in the action gateway, never here.
func (a *DynamicAgent) compileProposedMCPTool(name string, rt mcpconn.RemoteTool) *CompiledTool {
	spec := &mcp.Tool{
		Name:        name,
		Description: rt.Description,
	}
	if len(rt.InputSchema) > 0 {
		spec.InputSchema = rt.InputSchema
	}

	return &CompiledTool{
		Name:        name,
		Spec:        spec,
		CheckAccess: func(context.Context) (bool, error) { return true, nil },
		UnmarshalArgs: func(content string) (any, error) {
			if strings.TrimSpace(content) == "" {
				return map[string]any{}, nil
			}
			var args map[string]any
			if err := json.Unmarshal([]byte(content), &args); err != nil {
				return nil, err
			}
			return args, nil
		},
		UnmarshalResult: func(content string) (any, error) {
			return json.RawMessage(content), nil
		},
		JSONHandler: func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
			var args map[string]any
			if len(input) > 0 {
				if err := json.Unmarshal(input, &args); err != nil {
					return nil, err
				}
			}
			// The handler runs in the call's own message scope (s.Call sets the session's parent to the call message),
			// so GetSession(ctx).ParentID is this tool-call's message ID: the model-facing call identity and the parent
			// the resume path attaches the executed result under. Fall back to the tool name only if the scope is somehow
			// unset, so the action always has a non-empty identity for the ledger.
			toolCallID := rt.Name
			if sess := GetSession(ctx); sess != nil && sess.ParentID != "" {
				toolCallID = sess.ParentID
			}
			a.proposed = append(a.proposed, &ProposedAction{
				ToolCallID: toolCallID,
				Connector:  rt.Connector,
				Tool:       rt.Name,
				Args:       args,
				SchemaHash: rt.SchemaHash,
				Summary:    proposedSummary(rt, args),
			})
			// The tool did not run: the effect is deferred to human approval, and the loop pauses after all tool calls
			// in this turn are processed (PauseAfterToolCall checks len(a.proposed) > 0), so the model takes no further
			// turn until all actions are governed and executed. This placeholder is the call's result in the message
			// tree; the model only sees it on resume, alongside the injected real result, so it reads neutrally.
			return json.Marshal(map[string]string{
				"status": "This action was proposed and is awaiting human approval; it has not been executed yet.",
			})
		},
	}
}

// proposedSummary builds a one-line, human-readable description of a proposed action for the approval inbox, e.g.
// "mcp.jira.create_issue(project=OPS, summary=…)". It lists the argument names and short values; it never resolves or
// includes a secret, since the args carry only secret references at this point.
func proposedSummary(rt mcpconn.RemoteTool, args map[string]any) string {
	name := rt.RawName
	if name == "" {
		name = rt.Name
	}
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(name)
	b.WriteString("(")
	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(truncateArgValue(args[k]))
	}
	b.WriteString(")")
	return b.String()
}

// truncateArgValue renders one argument value compactly for a proposal summary, capping long values so the summary
// stays a single readable line in the inbox. The cap is per value, not per summary: 120 bytes is enough for a
// realistic title, email subject or URL to survive intact, while a body-sized value is visibly cut (the appended
// ellipsis marks it). The summary is only the inbox one-liner; the exact, complete arguments the approver signs are
// persisted separately as the approval's canonical args.
func truncateArgValue(v any) string {
	s := fmt.Sprintf("%v", v)
	const maxLen = 120
	if len(s) > maxLen {
		return truncateUTF8(s, maxLen) + "…"
	}
	return s
}

// sanitizeMCPToolName maps an effective MCP tool name (mcp.<connector>.<tool>, which mcpconn namespaces with dots)
// to a model-facing name that satisfies the LLM providers' function-name constraint ^[a-zA-Z0-9_-]+$ (OpenAI and
// DeepSeek reject dots, failing the whole request). Every character outside that class becomes '_'. The mcp_ prefix
// survives, so a sanitized name still cannot shadow a built-in tool (none begin with mcp_). Only the announced and
// registered name is sanitized; the original dotted name is retained separately and is what CallTool dispatches on.
func sanitizeMCPToolName(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

// systemPrompt combines the agent's own instructions with the project-wide ai_instructions carried by the session,
// matching how the built-in agents surface project instructions.
func (a *DynamicAgent) systemPrompt(s *Session) string {
	var b strings.Builder
	b.WriteString(a.Snapshot.Instructions)
	if pi := s.ProjectInstructions(); pi != "" {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString("The administrator has provided the following project-wide instructions, which may or may not be relevant to this task:\n")
		b.WriteString(pi)
	}
	return b.String()
}
