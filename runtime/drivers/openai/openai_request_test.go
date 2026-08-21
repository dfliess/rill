package openai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestCompleteAppliesConnectorRequestBehavior(t *testing.T) {
	fake := newFakeChatCompletionsServer(t, []string{completionResponse(`{"answer":"ok"}`)})
	defer fake.Close()

	ai := openTestAI(t, fake.URL, map[string]any{
		"structured_output_mode": structuredOutputModeJSONObject,
		"extra_body": map[string]any{
			"thinking": map[string]any{"type": "disabled"},
		},
	})

	_, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{textMessage("user", "answer as JSON")},
		OutputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"answer": {Type: "string"},
			},
			Required: []string{"answer"},
		},
	})
	require.NoError(t, err)

	body := fake.request(t, 0)
	require.Equal(t, map[string]any{"type": "disabled"}, body["thinking"])
	require.Equal(t, map[string]any{"type": "json_object"}, body["response_format"])

	messages := requireJSONArray(t, body["messages"])
	require.Len(t, messages, 2)
	schemaInstruction := requireJSONObject(t, messages[1])
	require.Equal(t, "system", schemaInstruction["role"])
	require.Contains(t, schemaInstruction["content"], "Return ONLY a single valid JSON object")
	require.Contains(t, schemaInstruction["content"], `"required":["answer"]`)
	require.False(t, containsJSONKey(body, "reasoning_content"))
}

func TestCompletePassesQwenVLLMExtraBody(t *testing.T) {
	fake := newFakeChatCompletionsServer(t, []string{completionResponse("ok")})
	defer fake.Close()

	ai := openTestAI(t, fake.URL, map[string]any{
		"extra_body": map[string]any{
			"chat_template_kwargs": map[string]any{"enable_thinking": false},
		},
	})

	_, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{textMessage("user", "hello")},
	})
	require.NoError(t, err)

	body := fake.request(t, 0)
	require.Equal(t, map[string]any{"enable_thinking": false}, body["chat_template_kwargs"])
	require.NotContains(t, body, "response_format")
}

func TestCompleteReportsReasoningTokensSeparately(t *testing.T) {
	fake := newFakeChatCompletionsServer(t, []string{`{
		"id":"chatcmpl-reasoning","object":"chat.completion","created":1,"model":"test-model",
		"choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}],
		"usage":{"prompt_tokens":3,"completion_tokens":7,"total_tokens":10,
			"completion_tokens_details":{"reasoning_tokens":5}}
	}`})
	defer fake.Close()

	ai := openTestAI(t, fake.URL, nil)
	res, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{textMessage("user", "think")},
	})
	require.NoError(t, err)
	require.Equal(t, 7, res.OutputTokens)
	require.Equal(t, 5, res.ReasoningTokens)
}

func TestCompleteReplaysReasoningContentFromMessageState(t *testing.T) {
	t.Setenv("RILL_OPENAI_STRUCTURED_OUTPUT_FALLBACK", structuredOutputModeJSONObject)
	reasoningContent := "provider-private-thinking\nwith \"quotes\" and unicode →"
	reasoningJSON, err := json.Marshal(reasoningContent)
	require.NoError(t, err)
	fake := newFakeChatCompletionsServer(t, []string{
		`{
			"id":"chatcmpl-tool","object":"chat.completion","created":1,"model":"test-model",
			"choices":[{"index":0,"message":{"role":"assistant","content":"","reasoning_content":` + string(reasoningJSON) + `,"tool_calls":[{"id":"call_1","type":"function","function":{"name":"do_work","arguments":"{}"}}]},"finish_reason":"tool_calls"}],
			"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}
		}`,
		completionResponse("done"),
	})
	defer fake.Close()

	ai := openTestAI(t, fake.URL, nil)
	schema := &jsonschema.Schema{Type: "object"}
	first, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages:     []*aiv1.CompletionMessage{textMessage("user", "use the tool")},
		Tools:        []*aiv1.Tool{{Name: "do_work", InputSchema: `{"type":"object"}`}},
		OutputSchema: schema,
	})
	require.NoError(t, err)
	require.NotNil(t, first.Message.ProviderData)
	require.Equal(t, reasoningContent, first.Message.ProviderData.Fields[providerDataReasoningContent].GetStringValue())

	_, err = ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{
			first.Message,
			{
				Role: "tool",
				Content: []*aiv1.ContentBlock{{
					BlockType: &aiv1.ContentBlock_ToolResult{ToolResult: &aiv1.ToolResult{Id: "call_1", Content: "ok"}},
				}},
			},
		},
	})
	require.NoError(t, err)

	firstBody := fake.request(t, 0)
	require.Equal(t, "json_schema", requireJSONObject(t, firstBody["response_format"])["type"])

	secondBody := fake.request(t, 1)
	messages := requireJSONArray(t, secondBody["messages"])
	require.Len(t, messages, 2)
	assistant := requireJSONObject(t, messages[0])
	require.Equal(t, "assistant", assistant["role"])
	require.Equal(t, "", assistant["content"], "tool-calling assistant content must be non-null")
	require.Equal(t, reasoningContent, assistant[providerDataReasoningContent])
}

func TestCompleteOmitsReasoningContentWhenItDoesNotApply(t *testing.T) {
	fake := newFakeChatCompletionsServer(t, []string{
		`{
			"id":"chatcmpl-no-tool","object":"chat.completion","created":1,"model":"test-model",
			"choices":[{"index":0,"message":{"role":"assistant","content":"done","reasoning_content":"not-needed-without-tool-calls"},"finish_reason":"stop"}],
			"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}
		}`,
		completionResponse("done again"),
	})
	defer fake.Close()

	ai := openTestAI(t, fake.URL, nil)
	first, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{textMessage("user", "hello")},
	})
	require.NoError(t, err)
	require.Nil(t, first.Message.ProviderData)

	_, err = ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{first.Message},
	})
	require.NoError(t, err)
	require.False(t, containsJSONKey(fake.request(t, 1), providerDataReasoningContent))
}

func TestCompleteReplaysReasoningWithMultipleToolCalls(t *testing.T) {
	fake := newFakeChatCompletionsServer(t, []string{
		`{
			"id":"chatcmpl-tools","object":"chat.completion","created":1,"model":"test-model",
			"choices":[{"index":0,"message":{"role":"assistant","content":null,"reasoning_content":"plan both calls","tool_calls":[
				{"id":"call_1","type":"function","function":{"name":"first","arguments":"{\"value\":1}"}},
				{"id":"call_2","type":"function","function":{"name":"second","arguments":"{\"value\":2}"}}
			]},"finish_reason":"tool_calls"}],
			"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}
		}`,
		completionResponse("done"),
	})
	defer fake.Close()

	ai := openTestAI(t, fake.URL, nil)
	first, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{textMessage("user", "use both tools")},
		Tools: []*aiv1.Tool{
			{Name: "first", InputSchema: `{"type":"object"}`},
			{Name: "second", InputSchema: `{"type":"object"}`},
		},
	})
	require.NoError(t, err)
	require.Len(t, first.Message.Content, 2)

	_, err = ai.Complete(t.Context(), &drivers.CompleteOptions{
		Messages: []*aiv1.CompletionMessage{
			first.Message,
			{
				Role: "tool",
				Content: []*aiv1.ContentBlock{
					{BlockType: &aiv1.ContentBlock_ToolResult{ToolResult: &aiv1.ToolResult{Id: "call_1", Content: "one"}}},
					{BlockType: &aiv1.ContentBlock_ToolResult{ToolResult: &aiv1.ToolResult{Id: "call_2", Content: "two"}}},
				},
			},
		},
	})
	require.NoError(t, err)

	messages := requireJSONArray(t, fake.request(t, 1)["messages"])
	require.Len(t, messages, 3)
	assistant := requireJSONObject(t, messages[0])
	require.Equal(t, "", assistant["content"], "tool-calling assistant content must be non-null")
	require.Equal(t, "plan both calls", assistant[providerDataReasoningContent])
	require.Len(t, requireJSONArray(t, assistant["tool_calls"]), 2)
	require.Equal(t, "call_1", requireJSONObject(t, messages[1])["tool_call_id"])
	require.Equal(t, "call_2", requireJSONObject(t, messages[2])["tool_call_id"])
}

func TestCompleteAcceptsStringAndJSONObjectResponses(t *testing.T) {
	for _, content := range []string{"plain response", `{"answer":"ok"}`} {
		t.Run(content, func(t *testing.T) {
			fake := newFakeChatCompletionsServer(t, []string{completionResponse(content)})
			defer fake.Close()

			ai := openTestAI(t, fake.URL, nil)
			res, err := ai.Complete(t.Context(), &drivers.CompleteOptions{
				Messages: []*aiv1.CompletionMessage{textMessage("user", "answer")},
			})
			require.NoError(t, err)
			require.Equal(t, content, res.Message.Content[0].GetText())
			require.Nil(t, res.Message.ProviderData)
		})
	}
}

func TestMessageToOpenAIRejectsNonStringReasoningContent(t *testing.T) {
	_, err := messageToOpenAI(&aiv1.CompletionMessage{
		Role:    "assistant",
		Content: textMessage("assistant", "tool call").Content,
		ProviderData: &structpb.Struct{Fields: map[string]*structpb.Value{
			providerDataReasoningContent: structpb.NewBoolValue(true),
		}},
	})
	require.EqualError(t, err, "provider_data.reasoning_content must be a string")
}

func TestOpenValidatesProviderRequestBehavior(t *testing.T) {
	t.Run("structured output mode", func(t *testing.T) {
		_, err := (driver{}).Open("", "", map[string]any{
			"api_key":                "test-key",
			"structured_output_mode": "xml",
		}, nil, nil, nil)
		require.ErrorContains(t, err, `invalid structured_output_mode "xml"`)
	})

	t.Run("reserved extra body fields", func(t *testing.T) {
		_, err := (driver{}).Open("", "", map[string]any{
			"api_key": "test-key",
			"extra_body": map[string]any{
				"audio":      map[string]any{"format": "wav"},
				"modalities": []any{"text", "audio"},
				"Tools":      []any{},
				"model":      "other-model",
				"stream":     true,
			},
		}, nil, nil, nil)
		require.EqualError(t, err, "extra_body cannot override fields controlled by Rill: Tools, audio, modalities, model, stream")
	})

	t.Run("non JSON extra body", func(t *testing.T) {
		_, err := (driver{}).Open("", "", map[string]any{
			"api_key": "test-key",
			"extra_body": map[string]any{
				"extension": make(chan int),
			},
		}, nil, nil, nil)
		require.ErrorContains(t, err, "extra_body must contain JSON-serializable values")
	})
}

func textMessage(role, text string) *aiv1.CompletionMessage {
	return &aiv1.CompletionMessage{
		Role: role,
		Content: []*aiv1.ContentBlock{{
			BlockType: &aiv1.ContentBlock_Text{Text: text},
		}},
	}
}

func openTestAI(t *testing.T, serverURL string, config map[string]any) drivers.AIService {
	t.Helper()
	if config == nil {
		config = make(map[string]any)
	}
	config["api_key"] = "test-key"
	config["base_url"] = serverURL + "/v1"
	config["model"] = "test-model"

	handle, err := (driver{}).Open("", "", config, nil, nil, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, handle.Close()) })
	ai, ok := handle.AsAI("")
	require.True(t, ok)
	return ai
}

type fakeChatCompletionsServer struct {
	*httptest.Server
	t         *testing.T
	mu        sync.Mutex
	requests  []map[string]any
	responses []string
}

func newFakeChatCompletionsServer(t *testing.T, responses []string) *fakeChatCompletionsServer {
	t.Helper()
	fake := &fakeChatCompletionsServer{t: t, responses: responses}
	fake.Server = httptest.NewServer(http.HandlerFunc(fake.serveHTTP))
	return fake
}

func (f *fakeChatCompletionsServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/chat/completions" {
		http.Error(w, "unexpected path "+r.URL.Path, http.StatusNotFound)
		return
	}

	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}

	f.mu.Lock()
	idx := len(f.requests)
	f.requests = append(f.requests, body)
	if idx >= len(f.responses) {
		f.mu.Unlock()
		http.Error(w, "unexpected request", http.StatusInternalServerError)
		return
	}
	response := f.responses[idx]
	f.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(response))
}

func (f *fakeChatCompletionsServer) request(t *testing.T, idx int) map[string]any {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	require.Greater(t, len(f.requests), idx)
	return f.requests[idx]
}

func completionResponse(content string) string {
	encoded, _ := json.Marshal(content)
	return `{"id":"chatcmpl-1","object":"chat.completion","created":1,"model":"test-model","choices":[{"index":0,"message":{"role":"assistant","content":` + string(encoded) + `},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`
}

func requireJSONObject(t *testing.T, value any) map[string]any {
	t.Helper()
	result, ok := value.(map[string]any)
	require.True(t, ok, "expected JSON object, got %T", value)
	return result
}

func requireJSONArray(t *testing.T, value any) []any {
	t.Helper()
	result, ok := value.([]any)
	require.True(t, ok, "expected JSON array, got %T", value)
	return result
}

func containsJSONKey(value any, target string) bool {
	switch value := value.(type) {
	case map[string]any:
		for key, nested := range value {
			if strings.EqualFold(key, target) || containsJSONKey(nested, target) {
				return true
			}
		}
	case []any:
		for _, nested := range value {
			if containsJSONKey(nested, target) {
				return true
			}
		}
	}
	return false
}
