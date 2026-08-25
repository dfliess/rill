package river

import (
	"context"
	"time"

	"github.com/rilldata/rill/admin"
	"github.com/riverqueue/river"
)

// notificationHistoryRetention is how long a notification stays in a user's history (kairos-cloud#148).
// The history answers "what was I told about lately", so it is measured in months, not years: a record nobody will
// ever look at is data kept for its own sake.
const notificationHistoryRetention = 90 * 24 * time.Hour

type DeleteExpiredNotificationsArgs struct{}

func (DeleteExpiredNotificationsArgs) Kind() string { return "delete_expired_notifications" }

type DeleteExpiredNotificationsWorker struct {
	river.WorkerDefaults[DeleteExpiredNotificationsArgs]
	admin *admin.Service
}

func (w *DeleteExpiredNotificationsWorker) Work(ctx context.Context, job *river.Job[DeleteExpiredNotificationsArgs]) error {
	return w.admin.DB.DeleteExpiredNotificationDeliveries(ctx, notificationHistoryRetention)
}
