package server_test

import (
	"context"
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/ai"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/pkg/ratelimit"
	"github.com/rilldata/rill/runtime/server"
	"github.com/rilldata/rill/runtime/server/auth"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestConversationWithReasoning verifies that a stored reasoning message is returned as a reasoning block
// instead of failing the whole conversation.
func TestConversationWithReasoning(t *testing.T) {
	rt, instanceID := testruntime.NewInstance(t)
	srv, err := server.NewServer(context.Background(), &server.Options{}, rt, zap.NewNop(), ratelimit.NewNoop(), activity.NewNoopClient())
	require.NoError(t, err)

	claims := &runtime.SecurityClaims{
		UserID:      "foo",
		Permissions: []runtime.Permission{runtime.ReadObjects, runtime.ReadMetrics, runtime.UseAI},
	}
	sess, err := ai.NewRunner(rt, activity.NewNoopClient()).Session(t.Context(), &ai.SessionOptions{
		InstanceID:        instanceID,
		CreateIfNotExists: true,
		Claims:            claims,
		UserAgent:         "test",
	})
	require.NoError(t, err)
	sess.AddMessage(&ai.AddMessageOptions{
		Role:        ai.RoleAssistant,
		Type:        ai.MessageTypeReasoning,
		ContentType: ai.MessageContentTypeJSON,
		Content:     `{"signature":"sig","model":"claude-test"}`,
	})
	require.NoError(t, sess.Flush(t.Context()))

	res, err := srv.GetConversation(auth.WithClaims(t.Context(), claims), &runtimev1.GetConversationRequest{
		InstanceId:     instanceID,
		ConversationId: sess.ID(),
	})
	require.NoError(t, err)
	require.Len(t, res.Messages, 1)
	require.Equal(t, "sig", res.Messages[0].Content[0].GetReasoning().GetSignature())
}
