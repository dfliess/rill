package act

import (
	"context"
	"testing"

	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/pkg/email"
	"github.com/stretchr/testify/require"
)

// recordingSender captures the last email handed to the email client, so the test can assert what the notify sink sent
// without a real SMTP server.
type recordingSender struct {
	called                 bool
	toEmail, subject, body string
}

func (r *recordingSender) Send(toEmail, toName, subject, body string) error {
	r.called = true
	r.toEmail = toEmail
	r.subject = subject
	r.body = body
	return nil
}

func TestNotifyActionSinkEmailsInitiator(t *testing.T) {
	rec := &recordingSender{}
	sink := NewNotifyActionSink(email.New(rec))

	claims := &runtime.SecurityClaims{UserAttributes: map[string]any{"email": "diego@example.com"}}
	res, err := sink.Apply(context.Background(), ActionRequest{
		AgentName: "welcome-mcp",
		Actor:     Actor{Claims: claims},
		Response:  "Todo bien. Diego saludado.",
	})
	require.NoError(t, err)
	require.True(t, rec.called, "an approved action must send the notification")
	require.Equal(t, "diego@example.com", rec.toEmail, "the notification goes to the run's initiator")
	require.Contains(t, rec.body, "Diego saludado", "the body carries the agent's response")
	require.Equal(t, "email:diego@example.com", res.Ref, "the run records an observable reference")
}

func TestNotifyActionSinkFailsClosed(t *testing.T) {
	withEmail := ActionRequest{Actor: Actor{Claims: &runtime.SecurityClaims{UserAttributes: map[string]any{"email": "x@y.com"}}}}

	// No email client: an approved notification must not be silently dropped.
	_, err := NewNotifyActionSink(nil).Apply(context.Background(), withEmail)
	require.Error(t, err)

	sink := NewNotifyActionSink(email.New(&recordingSender{}))

	// No initiator claims: cannot resolve a recipient.
	_, err = sink.Apply(context.Background(), ActionRequest{Actor: Actor{}})
	require.Error(t, err)

	// Claims without an email attribute: cannot resolve a recipient.
	_, err = sink.Apply(context.Background(), ActionRequest{Actor: Actor{Claims: &runtime.SecurityClaims{UserAttributes: map[string]any{}}}})
	require.Error(t, err)
}
