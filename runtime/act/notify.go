package act

import (
	"context"
	"errors"
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/rilldata/rill/runtime/pkg/email"
)

// notifyActionSink is Act's first real, end-to-end write: on an approved run it emails the agent's response to the
// run's initiator. It replaces the fail-closed denyActionSink so the loop investigate -> propose -> approve actually
// produces an observable effect (in dev the mail lands in Mailpit like any other Rill notification).
//
// It is deliberately minimal — the "rill.notify" internal action of the design, built on the existing email client.
// The full Fase 2 write path (structured tool-call proposals, per-tool policy, the action ledger, and external MCP
// action connectors such as Jira) is separate; this proves the north-star loop with the thinnest possible action.
type notifyActionSink struct {
	email *email.Client
}

var _ ActionSink = (*notifyActionSink)(nil)

// NewNotifyActionSink returns an ActionSink that emails the agent's response to the run's initiator via client. A nil
// client makes every Apply fail closed rather than silently drop an approved notification.
func NewNotifyActionSink(client *email.Client) ActionSink {
	return &notifyActionSink{email: client}
}

func (s *notifyActionSink) Apply(_ context.Context, req ActionRequest) (ActionResult, error) {
	if s.email == nil {
		return ActionResult{}, errors.New("act: notify action sink has no email client configured")
	}
	// The recipient is the run's initiator, resolved from the claims the run executes under. Fail closed if it is
	// absent: an approved action must reach someone, never silently go nowhere.
	if req.Actor.Claims == nil {
		return ActionResult{}, errors.New("act: cannot send notification: run has no initiator claims")
	}
	recipient, _ := req.Actor.Claims.UserAttributes["email"].(string)
	if recipient == "" {
		return ActionResult{}, errors.New("act: cannot send notification: initiator has no email attribute")
	}

	// The model response is plain text; escape it and turn newlines into breaks so it renders safely in the HTML body.
	body := strings.ReplaceAll(html.EscapeString(req.Response), "\n", "<br>")

	err := s.email.SendInformational(&email.Informational{
		ToEmail:    recipient,
		Subject:    fmt.Sprintf("Kairos · %s", req.AgentName),
		Body:       template.HTML(body), //nolint:gosec // body is html.EscapeString'd above.
		ShowFooter: true,
	})
	if err != nil {
		return ActionResult{}, fmt.Errorf("act: send notification email: %w", err)
	}
	return ActionResult{Ref: "email:" + recipient}, nil
}
