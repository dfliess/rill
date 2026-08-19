package pushnotifications

import (
	"context"
	"net/http"
	"net/http/httptest"
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
