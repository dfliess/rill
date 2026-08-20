package reconcilers_test

import (
	"errors"
	"strings"
	"testing"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/pkg/email"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
)

// TestAlertPushNotification checks that a dispatched alert notification is mirrored as one web push
// to the recipients of the alert's email notifier.
func TestAlertPushNotification(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Annotations: map[string]string{"organization_name": "acme", "project_name": "demo"},
	})
	putPushAlert(t, rt, id)

	// The alert passes, so nothing is dispatched.
	require.Empty(t, testruntime.PushNotifications(id))

	// Add data until the alert fails.
	failPushAlert(t, rt, id)

	pushes := testruntime.PushNotifications(id)
	require.Len(t, pushes, 1)
	require.Equal(t, "alerts", pushes[0].Category)
	require.Equal(t, []string{"somebody@example.com"}, pushes[0].Recipients)
	require.Equal(t, "Alerta: Test Alert", pushes[0].Title)
	require.Equal(t, "Tu alerta se activó el Thu, 04 Jan 2024 00:00:00 UTC.", pushes[0].Body)
	require.Equal(t, "/acme/demo/-/alerts/a1/open?execution_time=2024-01-04T00%3A00%3A00Z", pushes[0].LinkPath)
	require.Equal(t, "alerts:acme/demo/a1", pushes[0].Tag)

	// The push mirrors the email, it does not replace it.
	emails := rt.Email.Sender.(*email.TestSender).Emails
	require.Len(t, emails, 1)
	require.Equal(t, "somebody@example.com", emails[0].ToEmail)

	// Move the data a week ahead so the failing rows fall out of the alert's time range and it recovers.
	testruntime.PutFiles(t, rt, id, map[string]string{
		"/models/bar.sql": pushAlertModel + `
UNION ALL
SELECT '2024-01-12T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
`,
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	pushes = testruntime.PushNotifications(id)
	require.Len(t, pushes, 2)
	require.Equal(t, "alerts", pushes[1].Category)
	require.Equal(t, "Recuperada: Test Alert", pushes[1].Title)
	require.Contains(t, pushes[1].Body, "La alerta se recuperó el ")
	require.True(t, strings.HasPrefix(pushes[1].LinkPath, "/acme/demo/-/alerts/a1/open?execution_time="), "unexpected link path %q", pushes[1].LinkPath)
	// The recovery replaces the failure in the browser.
	require.Equal(t, pushes[0].Tag, pushes[1].Tag)

	require.Len(t, rt.Email.Sender.(*email.TestSender).Emails, 2)
}

// TestAlertPushNotificationOutsideRillCloud checks that alerts of an instance without the org and project
// annotations (i.e. Rill Developer) dispatch their emails without attempting a push.
func TestAlertPushNotificationOutsideRillCloud(t *testing.T) {
	rt, id := testruntime.NewInstance(t)
	putPushAlert(t, rt, id)
	failPushAlert(t, rt, id)

	require.Len(t, rt.Email.Sender.(*email.TestSender).Emails, 1)
	require.Empty(t, testruntime.PushNotifications(id))
}

// TestAlertPushNotificationFailure checks that a failing push notification does not fail the alert:
// the alert still reconciles cleanly and its email is still sent.
func TestAlertPushNotificationFailure(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Annotations: map[string]string{"organization_name": "acme", "project_name": "demo"},
	})
	testruntime.FailPushNotifications(id, errors.New("push service is down"))
	putPushAlert(t, rt, id)
	failPushAlert(t, rt, id)

	require.Len(t, rt.Email.Sender.(*email.TestSender).Emails, 1)
	require.Empty(t, testruntime.PushNotifications(id))

	// The alert reconciled without errors and recorded the notification as sent.
	testruntime.RequireReconcileState(t, rt, id, 4, 0, 0)
	a1 := testruntime.GetResource(t, rt, id, runtime.ResourceKindAlert, "a1")
	exec := a1.GetAlert().State.ExecutionHistory[0]
	require.True(t, exec.SentNotifications)
	require.Empty(t, exec.Result.ErrorMessage)
}

// TestReportPushNotification checks that a report notification is mirrored as one web push
// to the recipients of the report's email notifier.
func TestReportPushNotification(t *testing.T) {
	rt, id := testruntime.NewInstanceWithOptions(t, testruntime.InstanceOptions{
		Annotations: map[string]string{"organization_name": "acme", "project_name": "demo"},
	})
	testruntime.EnableReportDelivery(id)

	testruntime.PutFiles(t, rt, id, map[string]string{
		"/models/bar.sql": `
SELECT '2024-01-01T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
`,
		"/metrics/mv1.yaml": `
version: 1
type: metrics_view
model: bar
timeseries: __time
dimensions:
- column: country
measures:
- expression: count(*)
`,
		"/reports/r1.yaml": `
type: report
display_name: Ventas semanales
refs:
- type: MetricsView
  name: mv1
watermark: inherit
intervals:
  duration: P1D
query:
  name: MetricsViewAggregation
  args:
    metrics_view: mv1
    dimensions:
    - name: country
    measures:
    - name: measure_0
export:
  format: csv
notify:
  email:
    recipients:
      - somebody@example.com
`,
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, 4, 0, 0)

	// A report only runs on its schedule or on an ad-hoc trigger, which is how the admin service
	// runs one on demand. Nothing is dispatched until it does.
	require.Empty(t, testruntime.PushNotifications(id))
	testruntime.RefreshAndWait(t, rt, id, &runtimev1.ResourceName{Kind: runtime.ResourceKindReport, Name: "r1"})

	pushes := testruntime.PushNotifications(id)
	require.Len(t, pushes, 1)
	require.Equal(t, "reports", pushes[0].Category)
	require.Equal(t, []string{"somebody@example.com"}, pushes[0].Recipients)
	require.Equal(t, "Tu informe Ventas semanales está listo", pushes[0].Title)
	require.Equal(t, "El informe de Mon, 01 Jan 2024 00:00:00 UTC ya está disponible.", pushes[0].Body)
	require.Equal(t, "/acme/demo/-/reports/r1/open?execution_time=2024-01-01T00%3A00%3A00Z", pushes[0].LinkPath)
	require.Equal(t, "reports:acme/demo/r1", pushes[0].Tag)

	// The push mirrors the email, it does not replace it.
	emails := rt.Email.Sender.(*email.TestSender).Emails
	require.Len(t, emails, 1)
	require.Equal(t, "somebody@example.com", emails[0].ToEmail)
}

// TestReportPushNotificationOutsideRillCloud checks that a report of an instance whose admin service
// does not deliver reports is skipped without attempting a push.
func TestReportPushNotificationOutsideRillCloud(t *testing.T) {
	rt, id := testruntime.NewInstance(t)
	testruntime.PutFiles(t, rt, id, map[string]string{
		"/models/bar.sql": `
SELECT '2024-01-01T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
`,
		"/metrics/mv1.yaml": `
version: 1
type: metrics_view
model: bar
timeseries: __time
dimensions:
- column: country
measures:
- expression: count(*)
`,
		"/reports/r1.yaml": `
type: report
display_name: Ventas semanales
refs:
- type: MetricsView
  name: mv1
watermark: inherit
intervals:
  duration: P1D
query:
  name: MetricsViewAggregation
  args:
    metrics_view: mv1
    dimensions:
    - name: country
    measures:
    - name: measure_0
export:
  format: csv
notify:
  email:
    recipients:
      - somebody@example.com
`,
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, 4, 0, 0)
	testruntime.RefreshAndWait(t, rt, id, &runtimev1.ResourceName{Kind: runtime.ResourceKindReport, Name: "r1"})

	require.Empty(t, rt.Email.Sender.(*email.TestSender).Emails)
	require.Empty(t, testruntime.PushNotifications(id))
}

const pushAlertModel = `
SELECT '2024-01-01T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
UNION ALL
SELECT '2024-01-02T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
UNION ALL
SELECT '2024-01-03T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
UNION ALL
SELECT '2024-01-03T12:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
UNION ALL
SELECT '2024-01-04T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
`

// putPushAlert creates an alert that fails when its metrics view has four or more rows in the past week,
// against a model whose data still passes the assertion.
func putPushAlert(t *testing.T, rt *runtime.Runtime, id string) {
	testruntime.PutFiles(t, rt, id, map[string]string{
		"/models/bar.sql": `
SELECT '2024-01-01T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
`,
		"/metrics/mv1.yaml": `
version: 1
type: metrics_view
model: bar
timeseries: __time
dimensions:
- column: country
measures:
- expression: count(*)
`,
		"/alerts/a1.yaml": `
type: alert
display_name: Test Alert
refs:
- type: MetricsView
  name: mv1
watermark: inherit
intervals:
  duration: P1D
query:
  name: MetricsViewAggregation
  args:
    metrics_view: mv1
    dimensions:
    - name: country
    measures:
    - name: measure_0
    time_range:
      iso_duration: P1W
    having:
      cond:
        op: OPERATION_GTE
        exprs:
        - ident: measure_0
        - val: 4
on_recover: true
notify:
  email:
    recipients:
      - somebody@example.com
`,
	})
	testruntime.ReconcileParserAndWait(t, rt, id)
	testruntime.RequireReconcileState(t, rt, id, 4, 0, 0)
}

// failPushAlert adds data in two steps, so that the alert's watermark only advances to the day
// on which the assertion starts failing, and dispatches exactly one notification.
func failPushAlert(t *testing.T, rt *runtime.Runtime, id string) {
	testruntime.PutFiles(t, rt, id, map[string]string{
		"/models/bar.sql": `
SELECT '2024-01-01T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
UNION ALL
SELECT '2024-01-02T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
UNION ALL
SELECT '2024-01-03T00:00:00Z'::TIMESTAMP as __time, 'Denmark' as country
`,
	})
	testruntime.ReconcileParserAndWait(t, rt, id)

	testruntime.PutFiles(t, rt, id, map[string]string{"/models/bar.sql": pushAlertModel})
	testruntime.ReconcileParserAndWait(t, rt, id)
}
