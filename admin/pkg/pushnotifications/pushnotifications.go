// Package pushnotifications sends Web Push notifications to browser push subscriptions.
// It authenticates against push services using VAPID (RFC 8292).
package pushnotifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// ErrSubscriptionGone is returned by Send when the push service reports that the subscription
// no longer exists (HTTP 404 or 410). Callers should delete the subscription.
var ErrSubscriptionGone = errors.New("push subscription gone")

// defaultTTL is how long (in seconds) the push service should retain an undelivered notification.
const defaultTTL = 24 * 60 * 60

// Message is the notification payload.
// It is serialized to JSON and parsed by the frontend's service worker, so the field names are part of the contract.
type Message struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Link     string `json:"link"` // Absolute URL to open when the notification is clicked
	Category string `json:"category"`
	Tag      string `json:"tag"` // Notifications with the same tag replace each other
}

// Subscription identifies a browser push subscription as returned by the Push API.
type Subscription struct {
	Endpoint string
	P256dh   string
	Auth     string
}

// Client sends Web Push notifications signed with a VAPID key pair.
// A client created without keys is disabled: Enabled returns false and Send fails.
type Client struct {
	vapidPublicKey  string
	vapidPrivateKey string
	subject         string // VAPID subject, a "mailto:" or "https:" URI identifying the sender
}

// New creates a Client. If the keys are empty, the client is disabled.
// The subject is accepted with or without its "mailto:" prefix; see normalizeSubject.
func New(vapidPublicKey, vapidPrivateKey, subject string) *Client {
	return &Client{
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
		subject:         normalizeSubject(subject),
	}
}

// normalizeSubject strips a "mailto:" prefix from an email subject, because webpush-go adds one of its
// own to anything that is not an https URL. Configured as "mailto:someone@example.com", the JWT ends up
// claiming "sub":"mailto:mailto:someone@example.com": Google's push service accepts that, Apple's rejects
// every message with 403 BadJwtToken, so it breaks iOS alone and only once a real device is subscribed.
// Both spellings of the setting are valid per RFC 8292, so this normalizes rather than rejects.
func normalizeSubject(subject string) string {
	return strings.TrimPrefix(subject, "mailto:")
}

// Enabled returns true if the client is configured with a VAPID key pair.
func (c *Client) Enabled() bool {
	return c.vapidPublicKey != "" && c.vapidPrivateKey != ""
}

// VAPIDPublicKey returns the public VAPID key that browsers need to create subscriptions.
// It returns an empty string if the client is disabled.
func (c *Client) VAPIDPublicKey() string {
	return c.vapidPublicKey
}

// Send delivers msg to a single subscription.
// It returns ErrSubscriptionGone if the push service reports the subscription no longer exists.
func (c *Client) Send(ctx context.Context, sub *Subscription, msg *Message) error {
	if !c.Enabled() {
		return errors.New("push notifications are not configured")
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	res, err := webpush.SendNotificationWithContext(ctx, payload, &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}, &webpush.Options{
		Subscriber:      c.subject,
		TTL:             defaultTTL,
		VAPIDPublicKey:  c.vapidPublicKey,
		VAPIDPrivateKey: c.vapidPrivateKey,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusGone {
		return ErrSubscriptionGone
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("push service returned status %d", res.StatusCode)
	}
	return nil
}
