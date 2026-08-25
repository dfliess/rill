package server

import (
	"testing"

	"github.com/rilldata/rill/admin/database"
	"github.com/stretchr/testify/require"
)

// The notification history is a window over a record, not a page of an archive (kairos-cloud#148). These cover the
// two decisions that make it one: how wide the window is, and how a reader learns there is something behind it.

func TestNotificationHistoryWindow(t *testing.T) {
	require.Equal(t, notificationHistoryLimit, notificationHistoryWindow(0), "no page size asked for is the deployment's own limit")
	require.Equal(t, 3, notificationHistoryWindow(3), "a caller may ask for fewer")
	require.Equal(t, notificationHistoryLimit, notificationHistoryWindow(notificationHistoryLimit+1),
		"and may not talk the deployment into serving more than it wants to")
	require.Equal(t, notificationHistoryLimit, notificationHistoryWindow(-1), "nor into a nonsensical one")
}

func TestTrimToWindow(t *testing.T) {
	entries := func(n int) []*database.NotificationDelivery {
		out := make([]*database.NotificationDelivery, n)
		for i := range out {
			out[i] = &database.NotificationDelivery{Title: string(rune('a' + i))}
		}
		return out
	}

	// The read asks for one more than the window precisely so this can tell the two apart.
	shown, hasMore := trimToWindow(entries(4), 3)
	require.Len(t, shown, 3)
	require.True(t, hasMore, "the extra entry is what says there is more behind the window")
	require.Equal(t, "a", shown[0].Title, "and trimming keeps the newest, which is where the listing starts")

	shown, hasMore = trimToWindow(entries(3), 3)
	require.Len(t, shown, 3)
	require.False(t, hasMore, "a history that exactly fills the window is not truncated")

	shown, hasMore = trimToWindow(entries(1), 3)
	require.Len(t, shown, 1)
	require.False(t, hasMore)

	shown, hasMore = trimToWindow(nil, 3)
	require.Empty(t, shown)
	require.False(t, hasMore, "an empty history is empty, not truncated")
}
