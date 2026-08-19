package admin

import (
	"context"
	"errors"

	"github.com/rilldata/rill/admin/pkg/pushnotifications"
	"go.uber.org/zap"
)

// SendPushNotification sends a web push notification to all subscriptions of the given recipients,
// skipping recipients that have opted out of the given category.
// projectID identifies the project the notification originates from:
// preferences are per organization, so they are read for the organization that owns the project.
// It returns the number of notifications successfully sent.
// Individual delivery failures are logged and do not abort the remaining deliveries.
// Subscriptions the push service reports as gone are deleted.
func (s *Service) SendPushNotification(ctx context.Context, projectID, category string, recipientEmails []string, msg *pushnotifications.Message) (int, error) {
	if !s.Push.Enabled() || len(recipientEmails) == 0 {
		return 0, nil
	}

	subs, err := s.DB.FindPushSubscriptionsForRecipients(ctx, projectID, recipientEmails, category)
	if err != nil {
		return 0, err
	}

	sent := 0
	for _, sub := range subs {
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
				s.Logger.Warn("failed to delete gone push subscription", zap.String("subscription_id", sub.ID), zap.String("project_id", projectID), zap.Error(err))
			}
			continue
		}

		s.Logger.Warn("failed to send push notification", zap.String("subscription_id", sub.ID), zap.String("user_id", sub.UserID), zap.String("project_id", projectID), zap.Error(err))
	}
	return sent, nil
}
