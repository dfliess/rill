package reconcilers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"github.com/rilldata/rill/runtime/pkg/pbutil"
	"go.uber.org/zap"
)

// This file implements web push mirroring for alerts and reports (kairos-cloud#142).
// Push is an additive channel over the recipients of the "email" notifier:
// wherever a notification is dispatched by email, the same recipients also get one web push through
// the admin service, which resolves their subscriptions and their per-category preferences.
// Because the push is dispatched from the same place as the emails,
// it inherits the alert's renotify semantics without any logic of its own.
// Outside of Rill Cloud there is nothing to push to, so everything degrades to a silent no-op.

// pushBodyMaxLen caps the body of a push notification: payloads should stay small,
// and browsers truncate long bodies anyway.
const pushBodyMaxLen = 200

// pushNotification describes one web push dispatch that mirrors an email notification.
type pushNotification struct {
	category   string
	recipients []string
	title      string
	body       string
	linkPath   string
	tag        string
}

// alertPushNotification builds the push mirror of an alert status notification.
// The wording mirrors the alert emails in runtime/pkg/email/templates.
// It returns nil if there is nothing to push:
// no "email" notifier recipients, no org/project annotations (not running in Rill Cloud),
// or a status that does not notify.
func alertPushNotification(name string, spec *runtimev1.AlertSpec, status *drivers.AlertStatus, org, project string) *pushNotification {
	recipients := emailNotifierRecipients(spec.Notifiers)
	if len(recipients) == 0 || org == "" || project == "" {
		return nil
	}

	// The alert's execution time, like in the emails: it follows the watermark, so it may lag the wall clock.
	executionTime := status.ExecutionTime.Format(time.RFC1123)

	var title, body, linkPath string
	switch {
	case status.IsRecover:
		title = fmt.Sprintf("Recuperada: %s", status.DisplayName)
		body = fmt.Sprintf("La alerta se recuperó el %s de un fallo anterior.", executionTime)
		linkPath = alertOpenPath(org, project, name, status.ExecutionTime)
	case status.Status == runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL:
		title = fmt.Sprintf("Alerta: %s", status.DisplayName)
		body = fmt.Sprintf("Tu alerta se activó el %s.", executionTime)
		linkPath = alertOpenPath(org, project, name, status.ExecutionTime)
	case status.Status == runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR:
		title = fmt.Sprintf("Error en la alerta: %s", status.DisplayName)
		body = truncatePushBody(fmt.Sprintf("La alerta no se pudo evaluar el %s: %s", executionTime, status.ExecutionError))
		// Errors link to the alert itself instead of a dashboard, like the error emails do.
		linkPath = resourcePath(org, project, "alerts", name)
	default:
		return nil
	}

	return &pushNotification{
		category:   drivers.PushCategoryAlerts,
		recipients: recipients,
		title:      title,
		body:       body,
		linkPath:   linkPath,
		tag:        resourcePushTag("alerts", org, project, name),
	}
}

// reportPushNotification builds the push mirror of a "report ready" notification.
// It returns nil if there is nothing to push (see alertPushNotification).
func reportPushNotification(name string, spec *runtimev1.ReportSpec, reportTime time.Time, org, project string) *pushNotification {
	recipients := emailNotifierRecipients(spec.Notifiers)
	if len(recipients) == 0 || org == "" || project == "" {
		return nil
	}

	return &pushNotification{
		category:   drivers.PushCategoryReports,
		recipients: recipients,
		title:      fmt.Sprintf("Tu informe %s está listo", spec.DisplayName),
		body:       fmt.Sprintf("El informe de %s ya está disponible.", reportTime.Format(time.RFC1123)),
		// The emails open a link personalized per recipient (magic tokens, AI sessions),
		// but one push is shared by all recipients, so it links to the report's own page.
		linkPath: resourcePath(org, project, "reports", name),
		tag:      resourcePushTag("reports", org, project, name),
	}
}

// dispatchPushNotification sends n through the admin service.
// It is best-effort: errors are logged as warnings and must never fail the reconcile.
// A nil n is a no-op, as is an admin service that does not implement push (Rill Developer).
func dispatchPushNotification(ctx context.Context, c *runtime.Controller, resourceName string, n *pushNotification) {
	if n == nil {
		return
	}

	admin, release, err := c.Runtime.Admin(ctx, c.InstanceID)
	if err != nil {
		c.Logger.Warn("Failed to acquire admin client for push notification", zap.String("resource", resourceName), zap.Error(err), observability.ZapCtx(ctx))
		return
	}
	defer release()

	sent, err := admin.SendPushNotification(ctx, n.category, n.recipients, n.title, n.body, n.linkPath, n.tag)
	if err != nil {
		if !errors.Is(err, drivers.ErrNotImplemented) && !errors.Is(err, context.Canceled) {
			c.Logger.Warn("Failed to send push notification", zap.String("resource", resourceName), zap.String("category", n.category), zap.Error(err), observability.ZapCtx(ctx))
		}
		return
	}

	c.Logger.Debug("Sent push notifications", zap.String("resource", resourceName), zap.String("category", n.category), zap.Int("sent", sent), observability.ZapCtx(ctx))
}

// instanceOrgProject returns the org and project names from the instance annotations.
// The annotations are set by the admin service on deployment; outside of Rill Cloud they are empty.
func instanceOrgProject(ctx context.Context, c *runtime.Controller) (string, string) {
	inst, err := c.Runtime.Instance(ctx, c.InstanceID)
	if err != nil {
		return "", ""
	}
	return inst.Annotations["organization_name"], inst.Annotations["project_name"]
}

// emailNotifierRecipients returns the combined recipients of the "email" notifiers in the given list.
func emailNotifierRecipients(notifiers []*runtimev1.Notifier) []string {
	var recipients []string
	for _, notifier := range notifiers {
		if notifier.Connector == "email" {
			recipients = append(recipients, pbutil.ToSliceString(notifier.Properties.AsMap()["recipients"])...)
		}
	}
	return recipients
}

// alertOpenPath returns the frontend path that opens the dashboard the alert was checked against,
// with the execution time applied, which is what the "Open" button of the alert emails links to.
func alertOpenPath(org, project, name string, executionTime time.Time) string {
	qry := url.Values{"execution_time": []string{executionTime.UTC().Format(time.RFC3339)}}
	return fmt.Sprintf("%s/open?%s", resourcePath(org, project, "alerts", name), qry.Encode())
}

// resourcePath returns the frontend path of a resource's page, e.g. "/acme/demo/-/alerts/a1".
func resourcePath(org, project, section, name string) string {
	return fmt.Sprintf("/%s/%s/-/%s/%s", url.PathEscape(org), url.PathEscape(project), section, url.PathEscape(name))
}

// resourcePushTag returns a stable tag for a resource's push notifications,
// so newer notifications about the same resource replace older ones in the browser.
func resourcePushTag(section, org, project, name string) string {
	return fmt.Sprintf("%s:%s/%s/%s", section, org, project, name)
}

// truncatePushBody caps s at pushBodyMaxLen characters, appending an ellipsis if it was cut.
func truncatePushBody(s string) string {
	r := []rune(s)
	if len(r) <= pushBodyMaxLen {
		return s
	}
	return string(r[:pushBodyMaxLen]) + "…"
}
