package admin

import (
	"context"
	"errors"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/pkg/pushnotifications"
	"go.uber.org/zap"
)

// SendPushNotificationOptions is one notification and who it goes to.
type SendPushNotificationOptions struct {
	// ProjectID identifies the project the notification originates from.
	// Preferences are per organization, so they are read for the organization that owns it.
	ProjectID  string
	Category   string
	Recipients []string
	Message    *pushnotifications.Message
	// LinkPath is Message.Link before the frontend URL was prepended. The history stores the path rather than the
	// absolute link, so a deployment that changes domain keeps its history usable.
	LinkPath string
	// Pending, when non-nil, makes the notification carry each recipient's own count of work still waiting for
	// them in this project, which is what the installed app puts on its icon. Left nil, the notification says
	// nothing about it and the badge stays as it was.
	Pending *PendingCounts
}

// PendingCounts is how much work is waiting for each recipient in one project.
// ByRecipient is keyed by lowercased email; a recipient missing from it has nothing pending, which is a number the
// app badge needs, not a reason to stay silent (see pushnotifications.Pending).
type PendingCounts struct {
	Scope       string
	ByRecipient map[string]int32
}

// SendPushNotification sends a web push notification to all subscriptions of the given recipients,
// skipping recipients that have opted out of the given category.
// It returns the number of notifications successfully sent.
// Individual delivery failures are logged and do not abort the remaining deliveries.
// Subscriptions the push service reports as gone are deleted.
func (s *Service) SendPushNotification(ctx context.Context, opts *SendPushNotificationOptions) (int, error) {
	if len(opts.Recipients) == 0 {
		return 0, nil
	}

	s.recordNotificationHistory(ctx, opts)

	if !s.Push.Enabled() {
		return 0, nil
	}

	subs, err := s.DB.FindPushSubscriptionsForRecipients(ctx, opts.ProjectID, opts.Recipients, opts.Category)
	if err != nil {
		return 0, err
	}

	sent := 0
	for _, sub := range subs {
		msg := opts.Message
		if opts.Pending != nil {
			// The count differs per recipient, so each subscription gets its own copy of the message rather
			// than everyone sharing one. Authority to decide an approval discriminates by action, so two people
			// told about the same pause can genuinely have different amounts of work waiting.
			perRecipient := *msg
			perRecipient.Pending = &pushnotifications.Pending{
				Scope: opts.Pending.Scope,
				Count: opts.Pending.ByRecipient[sub.Email],
				User:  sub.UserID,
			}
			msg = &perRecipient
		}

		err := s.Push.Send(ctx, &pushnotifications.Subscription{
			Endpoint: sub.Endpoint,
			P256dh:   sub.P256dh,
			Auth:     sub.Auth,
		}, msg)
		if err == nil {
			sent++
			continue
		}

		if errors.Is(err, pushnotifications.ErrSubscriptionGone) {
			err = s.DB.DeletePushSubscriptionByEndpoint(ctx, sub.Endpoint)
			if err != nil {
				s.Logger.Warn("failed to delete gone push subscription", zap.String("subscription_id", sub.ID), zap.String("project_id", opts.ProjectID), zap.Error(err))
			}
			continue
		}

		s.Logger.Warn("failed to send push notification", zap.String("subscription_id", sub.ID), zap.String("user_id", sub.UserID), zap.String("project_id", opts.ProjectID), zap.Error(err))
	}
	return sent, nil
}

// recordNotificationHistory writes the notification into the history of every recipient that is a user of this
// deployment (kairos-cloud#148).
//
// It runs before the notification is sent, and whether or not push is configured at all. What is recorded is that
// this concerned that person, which is true of a deployment with no VAPID keys, of a recipient with no browser
// subscribed, and of one who turned the category off: they turned off the interruption, not the record.
//
// It is best-effort, like the delivery it accompanies: a history that cannot be written is logged and swallowed,
// never returned. Nobody should miss a notification because we failed to remember it.
func (s *Service) recordNotificationHistory(ctx context.Context, opts *SendPushNotificationOptions) {
	_, err := s.DB.InsertNotificationDeliveries(ctx, &database.InsertNotificationDeliveriesOptions{
		ProjectID:  opts.ProjectID,
		Recipients: opts.Recipients,
		Category:   opts.Category,
		Title:      opts.Message.Title,
		Body:       opts.Message.Body,
		LinkPath:   opts.LinkPath,
	})
	if err != nil {
		s.Logger.Warn("failed to record notification history", zap.String("project_id", opts.ProjectID), zap.String("category", opts.Category), zap.Error(err))
	}
}
