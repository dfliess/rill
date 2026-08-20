package pushnotifications

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/stretchr/testify/require"
)

func TestEnabled(t *testing.T) {
	c := New("", "", "")
	require.False(t, c.Enabled())
	require.Empty(t, c.VAPIDPublicKey())

	err := c.Send(context.Background(), &Subscription{Endpoint: "https://example.com"}, &Message{})
	require.ErrorContains(t, err, "not configured")

	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	require.NoError(t, err)
	c = New(publicKey, privateKey, "mailto:test@example.com")
	require.True(t, c.Enabled())
	require.Equal(t, publicKey, c.VAPIDPublicKey())
}

func TestSend(t *testing.T) {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	require.NoError(t, err)
	client := New(publicKey, privateKey, "mailto:test@example.com")

	// Test subscription keys generated in a browser (any valid P-256 point and 16-byte auth secret works).
	sub := &Subscription{
		P256dh: "BNcRdreALRFXTkOOUHK1EtK2wtaz5Ry4YfYCA_0QTpQtUbVlUls0VJXg7A8u-Ts1XbjhazAkj7I99e8QcYP7DkM=",
		Auth:   "zqbxT6JKstKSY9JKibZLSQ==",
	}

	t.Run("Delivered", func(t *testing.T) {
		var gotTTL string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotTTL = r.Header.Get("TTL")
			w.WriteHeader(http.StatusCreated)
		}))
		defer srv.Close()
		sub := &Subscription{Endpoint: srv.URL, P256dh: sub.P256dh, Auth: sub.Auth}

		err := client.Send(context.Background(), sub, &Message{Title: "Hello", Body: "World", Link: "https://example.com/x", Category: "alerts", Tag: "t1"})
		require.NoError(t, err)
		require.Equal(t, "86400", gotTTL)
	})

	t.Run("SubscriptionGone", func(t *testing.T) {
		for _, code := range []int{http.StatusNotFound, http.StatusGone} {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			sub := &Subscription{Endpoint: srv.URL, P256dh: sub.P256dh, Auth: sub.Auth}

			err := client.Send(context.Background(), sub, &Message{Title: "Hello"})
			require.ErrorIs(t, err, ErrSubscriptionGone)
			srv.Close()
		}
	})

	t.Run("PushServiceError", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()
		sub := &Subscription{Endpoint: srv.URL, P256dh: sub.P256dh, Auth: sub.Auth}

		err := client.Send(context.Background(), sub, &Message{Title: "Hello"})
		require.ErrorContains(t, err, "status 429")
	})
}

// TestSubjectReachesTheJWTOnce guards a trap in webpush-go: it prefixes "mailto:" to any subject that is
// not an https URL, so configuring the documented "mailto:someone@example.com" produced a JWT claiming
// "sub":"mailto:mailto:someone@example.com". Google accepts that and Apple answers 403 BadJwtToken, so the
// bug hides until an iPhone subscribes. Asserting on the header the push service actually receives is the
// only thing that catches it: a mock that returns 201 never will.
func TestSubjectReachesTheJWTOnce(t *testing.T) {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	require.NoError(t, err)

	for _, configured := range []string{"mailto:alerts@example.com", "alerts@example.com"} {
		t.Run(configured, func(t *testing.T) {
			var auth string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				auth = r.Header.Get("Authorization")
				w.WriteHeader(http.StatusCreated)
			}))
			defer srv.Close()

			client := New(publicKey, privateKey, configured)
			err := client.Send(context.Background(), &Subscription{
				Endpoint: srv.URL,
				P256dh:   "BNcRdreALRFXTkOOUHK1EtK2wtaz5Ry4YfYCA_0QTpQtUbVlUls0VJXg7A8u-Ts1XbjhazAkj7I99e8QcYP7DkM=",
				Auth:     "zqbxT6JKstKSY9JKibZLSQ==",
			}, &Message{Title: "Hello"})
			require.NoError(t, err)

			token := strings.TrimPrefix(strings.Split(auth, ",")[0], "vapid t=")
			parts := strings.Split(token, ".")
			require.Len(t, parts, 3, "unexpected JWT %q", token)
			claims, err := base64.RawURLEncoding.DecodeString(parts[1])
			require.NoError(t, err)
			var parsed struct {
				Sub string `json:"sub"`
			}
			require.NoError(t, json.Unmarshal(claims, &parsed))
			require.Equal(t, "mailto:alerts@example.com", parsed.Sub)
		})
	}
}
