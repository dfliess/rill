package server

import (
	"testing"

	aiv1 "github.com/rilldata/rill/proto/gen/rill/ai/v1"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/stretchr/testify/require"
)

// TestMessageToPBUserTextTurn is a regression test for the dynamic agent's opening user turn: a RoleUser text
// message (Type == MessageTypeText) must convert to a proto message with role "user" and a plain text content
// block. Before the fix, messageContentToPB's exhaustive type switch rejected the new "text" type with
// "unexpected message type", and the role of a non-call/result user message mapped to "assistant" instead of
// "user". A non-router text message does not consult the session, so a nil session is sufficient here.
func TestMessageToPBUserTextTurn(t *testing.T) {
	const prompt = "What drove the revenue spike?"
	msg := &ai.Message{
		ID:          "msg-1",
		Role:        ai.RoleUser,
		Type:        ai.MessageTypeText,
		ContentType: ai.MessageContentTypeText,
		Content:     prompt,
	}

	pb, err := messageToPB(nil, msg)
	require.NoError(t, err)

	// The role maps to "user", not "assistant" (the second bug).
	require.Equal(t, "user", pb.Role)
	require.Equal(t, string(ai.MessageTypeText), pb.Type)

	// The content is a single plain text block carrying the prompt verbatim (not "unexpected message type").
	require.Len(t, pb.Content, 1)
	textBlock, ok := pb.Content[0].BlockType.(*aiv1.ContentBlock_Text)
	require.True(t, ok, "user text turn must convert to a text content block")
	require.Equal(t, prompt, textBlock.Text)
	require.Equal(t, prompt, pb.Content[0].GetText())
}
