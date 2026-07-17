package ai_test

import (
	"context"
	"sync"
	"testing"

	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

// turnFunc produces the assistant message the simulated model returns for one completion call.
type turnFunc func(opts *drivers.CompleteOptions) *aiv1.CompletionMessage

// scriptedAIService is a deterministic drivers.AIService for tests. Each Complete call consumes the next scripted
// turn (falling back to defaultTurn once turns are exhausted) and records the messages and tools it was given, so
// tests can assert what the model saw and how many times it was invoked.
type scriptedAIService struct {
	turns       []turnFunc
	defaultTurn turnFunc

	mu       sync.Mutex
	calls    int
	inputs   [][]*aiv1.CompletionMessage
	toolSets [][]*aiv1.Tool
}

var _ drivers.AIService = (*scriptedAIService)(nil)

func (s *scriptedAIService) Complete(_ context.Context, opts *drivers.CompleteOptions) (*drivers.CompleteResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.inputs = append(s.inputs, opts.Messages)
	s.toolSets = append(s.toolSets, opts.Tools)

	turn := s.defaultTurn
	if s.calls < len(s.turns) {
		turn = s.turns[s.calls]
	}
	s.calls++
	if turn == nil {
		turn = textTurn("done")
	}

	return &drivers.CompleteResult{
		Message:      turn(opts),
		Provider:     "scripted",
		InputTokens:  1,
		OutputTokens: 1,
	}, nil
}

// textTurn makes the model reply with a plain text message (ending the tool loop).
func textTurn(text string) turnFunc {
	return func(_ *drivers.CompleteOptions) *aiv1.CompletionMessage {
		return &aiv1.CompletionMessage{
			Role:    "assistant",
			Content: []*aiv1.ContentBlock{{BlockType: &aiv1.ContentBlock_Text{Text: text}}},
		}
	}
}

// toolCallTurn makes the model request a tool call.
func toolCallTurn(name string, input map[string]any) turnFunc {
	return func(_ *drivers.CompleteOptions) *aiv1.CompletionMessage {
		in, err := structpb.NewStruct(input)
		if err != nil {
			panic(err)
		}
		return &aiv1.CompletionMessage{
			Role: "assistant",
			Content: []*aiv1.ContentBlock{{
				BlockType: &aiv1.ContentBlock_ToolCall{
					ToolCall: &aiv1.ToolCall{Id: "call_" + name, Name: name, Input: in},
				},
			}},
		}
	}
}

// newScriptedSession creates a session backed by the given simulated model.
func newScriptedSession(t *testing.T, script *scriptedAIService) *ai.Session {
	rt, instanceID := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Files: map[string]string{
			"rill.yaml": `ai_instructions: "Project rule: mention data residency when relevant."`,
			"models/orders.yaml": `
type: model
materialize: true
sql: |
  SELECT '2025-01-01T00:00:00Z'::TIMESTAMP AS event_time, 'United States' AS country, 100 AS revenue
  UNION ALL
  SELECT '2025-01-02T00:00:00Z'::TIMESTAMP AS event_time, 'Denmark' AS country, 10 AS revenue
`,
			"metrics/orders.yaml": `
type: metrics_view
model: orders
timeseries: event_time
dimensions:
- column: country
measures:
- name: count
  expression: COUNT(*)
`,
		},
	})
	testruntime.ReconcileParserAndWait(t, rt, instanceID)

	s := newSession(t, rt, instanceID)
	s.SetLLM(func(_ context.Context) (drivers.AIService, func(), error) {
		return script, func() {}, nil
	})
	return s
}

// firstSystemPrompt returns the text of the first system message from the first recorded completion.
func firstSystemPrompt(t *testing.T, script *scriptedAIService) string {
	t.Helper()
	require.NotEmpty(t, script.inputs, "model was never called")
	for _, m := range script.inputs[0] {
		if m.Role == string(ai.RoleSystem) {
			return m.Content[0].GetText()
		}
	}
	t.Fatal("no system message in first completion")
	return ""
}

// advertisedToolNames returns the names of the tools shown to the model on the first completion.
func advertisedToolNames(script *scriptedAIService) []string {
	var names []string
	if len(script.toolSets) > 0 {
		for _, tool := range script.toolSets[0] {
			names = append(names, tool.Name)
		}
	}
	return names
}

// TestDynamicAgentPersistsUserPromptTurn verifies the run persists the prompt as the conversation's opening user turn:
// a single RoleUser text message whose content is the prompt, ordered ahead of the agent's own turns, so the session
// rendered as a chat starts with the instruction. The agent's final response is persisted too, as a trailing assistant
// text turn, so a resumed segment (and the UI) sees what the model said and not just its tool calls.
func TestDynamicAgentPersistsUserPromptTurn(t *testing.T) {
	script := &scriptedAIService{
		turns: []turnFunc{
			textTurn("Here is my analysis."),
		},
	}
	s := newScriptedSession(t, script)

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "insights_agent",
		Instructions: "You are the Kairos insights agent.",
		Tools:        []string{ai.ListMetricsViewsName},
		MaxSteps:     5,
	})

	const prompt = "What drove the revenue spike?"
	_, err := ai.RunDynamicAgent(t.Context(), s, provider, "insights_agent", prompt)
	require.NoError(t, err)

	// The conversation is the user prompt followed by the agent's response, each a text turn carrying its content verbatim.
	textTurns := s.Messages(ai.FilterByType(ai.MessageTypeText))
	require.Len(t, textTurns, 2)
	require.Equal(t, ai.RoleUser, textTurns[0].Role)
	require.Equal(t, prompt, textTurns[0].Content)
	require.Equal(t, ai.MessageContentTypeText, textTurns[0].ContentType)
	require.Equal(t, ai.RoleAssistant, textTurns[1].Role)
	require.Equal(t, "Here is my analysis.", textTurns[1].Content)

	// The prompt is the first message in the session, so the conversation opens with it.
	require.Equal(t, textTurns[0].ID, s.Messages()[0].ID)
}

// TestDynamicAgentRespondsWithInstructions verifies that a dynamic agent runs from a snapshot, surfaces its own
// instructions and the project instructions to the model, can call a metrics-view tool in its allowed set, and
// returns the model's final response.
func TestDynamicAgentRespondsWithInstructions(t *testing.T) {
	script := &scriptedAIService{
		turns: []turnFunc{
			toolCallTurn(ai.ListMetricsViewsName, map[string]any{}),
			textTurn("Available metrics view: orders"),
		},
	}
	s := newScriptedSession(t, script)

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "insights_agent",
		DisplayName:  "Insights Agent",
		Instructions: "You are the Kairos insights agent. Always answer concisely.",
		Tools:        []string{ai.ListMetricsViewsName},
		MaxSteps:     5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "insights_agent", "What metrics views exist?")
	require.NoError(t, err)
	require.Equal(t, "insights_agent", res.Agent)
	require.Equal(t, "Available metrics view: orders", res.Response)

	// The agent's own instructions and the project ai_instructions both reached the model.
	systemPrompt := firstSystemPrompt(t, script)
	require.Contains(t, systemPrompt, "You are the Kairos insights agent")
	require.Contains(t, systemPrompt, "mention data residency when relevant")

	// The listed tool was advertised and actually executed as a sub-call of the agent.
	require.Equal(t, []string{ai.ListMetricsViewsName}, advertisedToolNames(script))
	toolCalls := s.Messages(ai.FilterByType(ai.MessageTypeCall), ai.FilterByTool(ai.ListMetricsViewsName))
	require.Len(t, toolCalls, 1)
}

// TestDynamicAgentRejectsUnlistedTool verifies fail-closed enforcement: a tool that is otherwise accessible but is
// not in the snapshot is neither advertised to the model nor executed when the model forcibly proposes it.
func TestDynamicAgentRejectsUnlistedTool(t *testing.T) {
	script := &scriptedAIService{
		turns: []turnFunc{
			// The model forces a call to `navigate`, which the caller can access but the snapshot does not list.
			toolCallTurn(ai.NavigateName, map[string]any{"kind": "metrics_view", "name": "orders"}),
		},
	}
	s := newScriptedSession(t, script)

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "restricted_agent",
		Instructions: "You may only list metrics views.",
		Tools:        []string{ai.ListMetricsViewsName},
		MaxSteps:     5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "restricted_agent", "Take me to the orders dashboard.")
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), ai.NavigateName)
	require.Contains(t, err.Error(), "not in the allowed set")

	// Invisible: navigate was never advertised to the model.
	require.Equal(t, []string{ai.ListMetricsViewsName}, advertisedToolNames(script))
	require.NotContains(t, advertisedToolNames(script), ai.NavigateName)

	// Not executed: no navigate tool call was ever recorded on the session.
	require.Empty(t, s.Messages(ai.FilterByTool(ai.NavigateName)))
}

// TestDynamicAgentRejectsNonDeclarableTool verifies the defense-in-depth check on the StaticAgentProvider path: a
// snapshot that declares a write/development tool (here write_file) is rejected before the model is ever invoked,
// even though that path does not pass through the reconciler that normally rejects such tools.
func TestDynamicAgentRejectsNonDeclarableTool(t *testing.T) {
	script := &scriptedAIService{}
	s := newScriptedSession(t, script)

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "writer_agent",
		Instructions: "You may edit files.",
		Tools:        []string{ai.WriteFileName},
		MaxSteps:     5,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "writer_agent", "Fix the config file.")
	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), ai.WriteFileName)
	require.Contains(t, err.Error(), "may not use")

	// The model was never invoked: the snapshot was rejected before any completion.
	script.mu.Lock()
	calls := script.calls
	script.mu.Unlock()
	require.Equal(t, 0, calls)
}

// TestDynamicAgentMaxStepsCutsLoop verifies that MaxSteps bounds the loop: a model that never stops calling tools
// is invoked exactly MaxSteps times and then the run is failed rather than looping unbounded.
func TestDynamicAgentMaxStepsCutsLoop(t *testing.T) {
	const maxSteps = 3
	script := &scriptedAIService{
		// No scripted turns: every completion falls back to this, so the model always asks to call the tool again.
		defaultTurn: toolCallTurn(ai.ListMetricsViewsName, map[string]any{}),
	}
	s := newScriptedSession(t, script)

	provider := ai.NewStaticAgentProvider(&ai.AgentSnapshot{
		Name:         "looping_agent",
		Instructions: "You keep working.",
		Tools:        []string{ai.ListMetricsViewsName},
		MaxSteps:     maxSteps,
	})

	res, err := ai.RunDynamicAgent(t.Context(), s, provider, "looping_agent", "Keep going forever.")
	require.Error(t, err)
	require.Nil(t, res)

	// The model was invoked exactly MaxSteps times: the loop was cut, not run unbounded.
	script.mu.Lock()
	calls := script.calls
	script.mu.Unlock()
	require.Equal(t, maxSteps, calls)
}
