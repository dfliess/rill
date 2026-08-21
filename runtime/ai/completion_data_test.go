package ai

import (
	"testing"

	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestCompletionDataRoundTripAndToolCallGrouping(t *testing.T) {
	providerData, err := structpb.NewStruct(map[string]any{"reasoning_content": "exact provider reasoning"})
	require.NoError(t, err)

	call1 := &Message{
		ID: "11111111-1111-1111-1111-111111111111", Role: RoleAssistant, Type: MessageTypeCall,
		Tool: "first", ContentType: MessageContentTypeJSON, Content: `{}`,
		ProviderData: providerData, CompletionID: "completion-1", CompletionToolCallID: "provider-call-1",
	}
	result1 := &Message{
		ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", ParentID: call1.ID, Role: RoleAssistant, Type: MessageTypeResult,
		Tool: "first", ContentType: MessageContentTypeJSON, Content: `{"ok":true}`,
	}
	call2 := &Message{
		ID: "22222222-2222-2222-2222-222222222222", Role: RoleAssistant, Type: MessageTypeCall,
		Tool: "second", ContentType: MessageContentTypeJSON, Content: `{}`,
		ProviderData: providerData, CompletionID: "completion-1", CompletionToolCallID: "provider-call-2",
	}
	result2 := &Message{
		ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", ParentID: call2.ID, Role: RoleAssistant, Type: MessageTypeResult,
		Tool: "second", ContentType: MessageContentTypeJSON, Content: `{"ok":true}`,
	}

	// Exercise the same JSON envelope that Flush writes and Session loads after a runtime restart.
	persisted, err := call1.marshalCompletionData()
	require.NoError(t, err)
	reloadedCall1 := *call1
	reloadedCall1.ProviderData = nil
	reloadedCall1.CompletionID = ""
	reloadedCall1.CompletionToolCallID = ""
	require.NoError(t, reloadedCall1.loadCompletionData(persisted))
	require.Equal(t, call1.CompletionID, reloadedCall1.CompletionID)
	require.Equal(t, call1.CompletionToolCallID, reloadedCall1.CompletionToolCallID)
	require.Equal(t, call1.ProviderData.AsMap(), reloadedCall1.ProviderData.AsMap())

	messages := []*Message{&reloadedCall1, result1, call2, result2}
	s := &Session{BaseSession: &BaseSession{messages: messages}}
	completionMessages := s.NewCompletionMessages(messages)
	require.Len(t, completionMessages, 3)

	assistant := completionMessages[0]
	require.Equal(t, "assistant", assistant.Role)
	require.Equal(t, providerData.AsMap(), assistant.ProviderData.AsMap())
	require.Len(t, assistant.Content, 2, "parallel tool calls must remain in one assistant message")
	require.Equal(t, "provider-call-1", assistant.Content[0].GetToolCall().Id)
	require.Equal(t, "provider-call-2", assistant.Content[1].GetToolCall().Id)

	require.Equal(t, "provider-call-1", completionMessages[1].Content[0].GetToolResult().Id)
	require.Equal(t, "provider-call-2", completionMessages[2].Content[0].GetToolResult().Id)
}

func TestMaybeTruncateMessagesDropsPartialParallelToolBatch(t *testing.T) {
	messages := make([]*aiv1.CompletionMessage, 105)
	for i := range messages {
		messages[i] = NewTextCompletionMessage(RoleUser, "filler")
	}
	messages[3] = &aiv1.CompletionMessage{
		Role: "assistant",
		Content: []*aiv1.ContentBlock{
			{BlockType: &aiv1.ContentBlock_ToolCall{ToolCall: &aiv1.ToolCall{Id: "call-1", Name: "first"}}},
			{BlockType: &aiv1.ContentBlock_ToolCall{ToolCall: &aiv1.ToolCall{Id: "call-2", Name: "second"}}},
		},
	}
	// call-1's result falls into the omitted middle while call-2's result remains in the retained suffix. The grouped
	// assistant and the surviving result must both be dropped, otherwise the provider receives a partial tool batch.
	messages[4] = &aiv1.CompletionMessage{Role: "tool", Content: []*aiv1.ContentBlock{{
		BlockType: &aiv1.ContentBlock_ToolResult{ToolResult: &aiv1.ToolResult{Id: "call-1", Content: "one"}},
	}}}
	messages[89] = &aiv1.CompletionMessage{Role: "tool", Content: []*aiv1.ContentBlock{{
		BlockType: &aiv1.ContentBlock_ToolResult{ToolResult: &aiv1.ToolResult{Id: "call-2", Content: "two"}},
	}}}

	truncated := maybeTruncateMessages(messages)
	for _, message := range truncated {
		for _, block := range message.Content {
			require.Nil(t, block.GetToolCall())
			require.Nil(t, block.GetToolResult())
		}
	}
}
