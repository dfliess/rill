package reconcilers

import (
	"strings"
	"testing"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

var pushTestTime = time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)

func TestAlertPushNotificationMessage(t *testing.T) {
	spec := &runtimev1.AlertSpec{Notifiers: []*runtimev1.Notifier{emailNotifier(t, "somebody@example.com")}}

	t.Run("fail", func(t *testing.T) {
		n := alertPushNotification("a1", spec, &drivers.AlertStatus{
			DisplayName:   "Ventas diarias",
			ExecutionTime: pushTestTime,
			Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL,
		}, "acme", "demo")
		require.NotNil(t, n)
		require.Equal(t, drivers.PushCategoryAlerts, n.Category)
		require.Equal(t, []string{"somebody@example.com"}, n.Recipients)
		require.Equal(t, "Alerta: Ventas diarias", n.Title)
		require.Equal(t, "Tu alerta se activó el Thu, 04 Jan 2024 00:00:00 UTC.", n.Body)
		// Same destination as the "Open" button of the alert email: the dashboard at the execution time.
		require.Equal(t, "/acme/demo/-/alerts/a1/open?execution_time=2024-01-04T00%3A00%3A00Z", n.LinkPath)
		require.Equal(t, "alerts:acme/demo/a1", n.Tag)
	})

	t.Run("recover", func(t *testing.T) {
		n := alertPushNotification("a1", spec, &drivers.AlertStatus{
			DisplayName:   "Ventas diarias",
			ExecutionTime: pushTestTime,
			Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_PASS,
			IsRecover:     true,
		}, "acme", "demo")
		require.NotNil(t, n)
		require.Equal(t, "Recuperada: Ventas diarias", n.Title)
		require.Equal(t, "La alerta se recuperó el Thu, 04 Jan 2024 00:00:00 UTC de un fallo anterior.", n.Body)
		require.Equal(t, "/acme/demo/-/alerts/a1/open?execution_time=2024-01-04T00%3A00%3A00Z", n.LinkPath)
		// A recovery replaces the notification of the failure it recovers from.
		require.Equal(t, "alerts:acme/demo/a1", n.Tag)
	})

	t.Run("error", func(t *testing.T) {
		n := alertPushNotification("a1", spec, &drivers.AlertStatus{
			DisplayName:    "Ventas diarias",
			ExecutionTime:  pushTestTime,
			Status:         runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR,
			ExecutionError: "table not found",
		}, "acme", "demo")
		require.NotNil(t, n)
		require.Equal(t, "Error en la alerta: Ventas diarias", n.Title)
		require.Equal(t, "La alerta no se pudo evaluar el Thu, 04 Jan 2024 00:00:00 UTC: table not found", n.Body)
		// Errors open the alert itself instead of a dashboard, like the alert error email.
		require.Equal(t, "/acme/demo/-/alerts/a1", n.LinkPath)
	})

	t.Run("error with a long message", func(t *testing.T) {
		n := alertPushNotification("a1", spec, &drivers.AlertStatus{
			DisplayName:    "Ventas diarias",
			ExecutionTime:  pushTestTime,
			Status:         runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR,
			ExecutionError: strings.Repeat("é", 500),
		}, "acme", "demo")
		require.NotNil(t, n)
		require.Equal(t, pushBodyMaxLen+1, len([]rune(n.Body)))
		require.True(t, strings.HasSuffix(n.Body, "…"))
	})

	t.Run("a pass that is not a recovery does not notify", func(t *testing.T) {
		require.Nil(t, alertPushNotification("a1", spec, &drivers.AlertStatus{
			DisplayName:   "Ventas diarias",
			ExecutionTime: pushTestTime,
			Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_PASS,
		}, "acme", "demo"))
	})

	t.Run("without email recipients", func(t *testing.T) {
		status := &drivers.AlertStatus{DisplayName: "Ventas diarias", ExecutionTime: pushTestTime, Status: runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL}

		// A notifier that is not email: push mirrors email recipients only.
		props, err := structpb.NewStruct(map[string]any{"channels": []any{"#general"}})
		require.NoError(t, err)
		slack := &runtimev1.AlertSpec{Notifiers: []*runtimev1.Notifier{{Connector: "slack", Properties: props}}}
		require.Nil(t, alertPushNotification("a1", slack, status, "acme", "demo"))

		// No notifiers at all.
		require.Nil(t, alertPushNotification("a1", &runtimev1.AlertSpec{}, status, "acme", "demo"))

		// An email notifier with an empty recipient list.
		empty := &runtimev1.AlertSpec{Notifiers: []*runtimev1.Notifier{emailNotifier(t)}}
		require.Nil(t, alertPushNotification("a1", empty, status, "acme", "demo"))
	})

	t.Run("outside Rill Cloud", func(t *testing.T) {
		status := &drivers.AlertStatus{DisplayName: "Ventas diarias", ExecutionTime: pushTestTime, Status: runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL}
		require.Nil(t, alertPushNotification("a1", spec, status, "", ""))
		require.Nil(t, alertPushNotification("a1", spec, status, "acme", ""))
		require.Nil(t, alertPushNotification("a1", spec, status, "", "demo"))
	})

	t.Run("recipients of several email notifiers", func(t *testing.T) {
		two := &runtimev1.AlertSpec{Notifiers: []*runtimev1.Notifier{
			emailNotifier(t, "a@example.com", "b@example.com"),
			emailNotifier(t, "c@example.com"),
		}}
		n := alertPushNotification("a1", two, &drivers.AlertStatus{
			DisplayName:   "Ventas diarias",
			ExecutionTime: pushTestTime,
			Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL,
		}, "acme", "demo")
		require.NotNil(t, n)
		require.Equal(t, []string{"a@example.com", "b@example.com", "c@example.com"}, n.Recipients)
	})

	t.Run("names that need escaping", func(t *testing.T) {
		n := alertPushNotification("mi alerta", spec, &drivers.AlertStatus{
			DisplayName:   "Ventas diarias",
			ExecutionTime: pushTestTime,
			Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR,
		}, "acme corp", "demo/1")
		require.NotNil(t, n)
		require.Equal(t, "/acme%20corp/demo%2F1/-/alerts/mi%20alerta", n.LinkPath)
	})
}

func TestReportPushNotificationMessage(t *testing.T) {
	spec := &runtimev1.ReportSpec{
		DisplayName: "Ventas semanales",
		Notifiers:   []*runtimev1.Notifier{emailNotifier(t, "somebody@example.com")},
	}

	t.Run("ready", func(t *testing.T) {
		n := reportPushNotification("r1", spec, pushTestTime, "acme", "demo", "")
		require.NotNil(t, n)
		require.Equal(t, drivers.PushCategoryReports, n.Category)
		require.Equal(t, []string{"somebody@example.com"}, n.Recipients)
		require.Equal(t, "Tu informe Ventas semanales está listo", n.Title)
		require.Equal(t, "El informe de Thu, 04 Jan 2024 00:00:00 UTC ya está disponible.", n.Body)
		require.Equal(t, "/acme/demo/-/reports/r1/open?execution_time=2024-01-04T00%3A00%3A00Z", n.LinkPath)
		require.Equal(t, "reports:acme/demo/r1", n.Tag)
	})

	// An AI report's analysis IS a conversation, so the push has to land on it the way the email does.
	// Without the session the /open route has nothing to resolve and drops the reader on the report's own
	// page, which is the bug this pins: the reader tapped a notification about one run and got a listing.
	t.Run("AI report deep-links to its conversation", func(t *testing.T) {
		n := reportPushNotification("r1", spec, pushTestTime, "acme", "demo", "fe6da861-0a68-4cdf-b822-d6ea333c7e9f")
		require.NotNil(t, n)
		require.Equal(t,
			"/acme/demo/-/reports/r1/open?execution_time=2024-01-04T00%3A00%3A00Z&session_id=fe6da861-0a68-4cdf-b822-d6ea333c7e9f",
			n.LinkPath)
	})


	t.Run("without email recipients", func(t *testing.T) {
		require.Nil(t, reportPushNotification("r1", &runtimev1.ReportSpec{DisplayName: "Ventas semanales"}, pushTestTime, "acme", "demo", ""))
	})

	t.Run("outside Rill Cloud", func(t *testing.T) {
		require.Nil(t, reportPushNotification("r1", spec, pushTestTime, "", "", ""))
	})
}

func emailNotifier(t *testing.T, recipients ...string) *runtimev1.Notifier {
	t.Helper()
	props := make([]any, 0, len(recipients))
	for _, r := range recipients {
		props = append(props, r)
	}
	s, err := structpb.NewStruct(map[string]any{"recipients": props})
	require.NoError(t, err)
	return &runtimev1.Notifier{Connector: "email", Properties: s}
}
