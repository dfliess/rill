package testruntime

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"go.uber.org/zap"
)

// noopAdminService is the admin service used in tests.
// It is registered as a driver and doubles as the handle it opens,
// so there is exactly one of them and it can hold the state tests observe or configure, keyed by instance.
type noopAdminService struct {
	mu                sync.Mutex
	pushNotifications map[string][]PushNotification
	pushErrs          map[string]error
	reportDelivery    map[string]bool
}

var (
	_ drivers.AdminService = &noopAdminClient{}
	_ drivers.Handle       = &noopAdminService{}
	_ drivers.Driver       = &noopAdminService{}

	noopAdmin = &noopAdminService{
		pushNotifications: make(map[string][]PushNotification),
		pushErrs:          make(map[string]error),
		reportDelivery:    make(map[string]bool),
	}
)

func init() {
	drivers.Register("noop_admin", noopAdmin)
}

func (n *noopAdminService) GetAlertMetadata(ctx context.Context, alertName, ownerID string, emailRecipients []string, anonRecipients bool, annotations map[string]string, queryForUserID, queryForUserEmail string) (*drivers.AlertMetadata, error) {
	return nil, drivers.ErrNotImplemented
}

func (n *noopAdminService) GetConfig(ctx context.Context) (*drivers.Config, error) {
	return &drivers.Config{}, nil
}

func (n *noopAdminService) ProvisionConnector(ctx context.Context, name, driver string, args map[string]any) (map[string]any, error) {
	return nil, drivers.ErrNotImplemented
}

func (n *noopAdminService) ListDeployments(ctx context.Context) ([]*drivers.Deployment, error) {
	return nil, drivers.ErrNotImplemented
}

func (n *noopAdminService) UpdateProjectVariables(ctx context.Context, environment string, variables map[string]string) error {
	return drivers.ErrNotImplemented
}

// noopAdminClient scopes the noop admin service to an instance,
// so push notifications can be recorded per instance for test assertions.
type noopAdminClient struct {
	*noopAdminService
	instanceID string
}

// GetReportMetadata implements [drivers.AdminService] by delivering to every recipient,
// but only for instances opted in with EnableReportDelivery:
// reports are skipped by default, like on a deployment whose admin service does not support them.
func (c *noopAdminClient) GetReportMetadata(ctx context.Context, reportName, ownerID, webOpenMode string, emailRecipients []string, anonRecipients bool, executionTime time.Time) (*drivers.ReportMetadata, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.reportDelivery[c.instanceID] {
		return nil, drivers.ErrNotImplemented
	}

	meta := &drivers.ReportMetadata{ReportDelivery: make(map[string]drivers.ReportDelivery, len(emailRecipients)+1)}
	delivery := drivers.ReportDelivery{
		OpenURL:        "https://example.com/open",
		ExportURL:      "https://example.com/export",
		EditURL:        "https://example.com/edit",
		UnsubscribeURL: "https://example.com/unsubscribe",
	}
	for _, recipient := range emailRecipients {
		meta.ReportDelivery[recipient] = delivery
	}
	if anonRecipients {
		meta.ReportDelivery[""] = delivery
	}
	return meta, nil
}

// SendPushNotification implements [drivers.AdminService] by recording the notification,
// so tests can assert on push dispatches the way they assert on emails through [email.TestSender].
// It fails instead of recording if the instance was set up with FailPushNotifications.
func (c *noopAdminClient) SendPushNotification(ctx context.Context, category string, emails []string, title, body, linkPath, tag string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.pushErrs[c.instanceID]; err != nil {
		return 0, err
	}

	c.pushNotifications[c.instanceID] = append(c.pushNotifications[c.instanceID], PushNotification{
		Category:   category,
		Recipients: emails,
		Title:      title,
		Body:       body,
		LinkPath:   linkPath,
		Tag:        tag,
	})
	return len(emails), nil
}

// PushNotification records a push notification sent through the noop admin service.
type PushNotification struct {
	Category   string
	Recipients []string
	Title      string
	Body       string
	LinkPath   string
	Tag        string
}

// PushNotifications returns the push notifications recorded for the given instance, in dispatch order.
func PushNotifications(instanceID string) []PushNotification {
	noopAdmin.mu.Lock()
	defer noopAdmin.mu.Unlock()
	return slices.Clone(noopAdmin.pushNotifications[instanceID])
}

// FailPushNotifications makes the admin service of the given instance fail every push notification,
// so tests can assert that a failing push does not affect the resource it mirrors.
func FailPushNotifications(instanceID string, err error) {
	noopAdmin.mu.Lock()
	defer noopAdmin.mu.Unlock()
	noopAdmin.pushErrs[instanceID] = err
}

// EnableReportDelivery makes the admin service of the given instance return delivery metadata for reports,
// so tests can exercise the report notifications of an instance that runs in Rill Cloud.
func EnableReportDelivery(instanceID string) {
	noopAdmin.mu.Lock()
	defer noopAdmin.mu.Unlock()
	noopAdmin.reportDelivery[instanceID] = true
}

// HasAnonymousSourceAccess implements [drivers.Driver].
func (n *noopAdminService) HasAnonymousSourceAccess(ctx context.Context, srcProps map[string]any, logger *zap.Logger) (bool, error) {
	return true, nil
}

// Open implements [drivers.Driver].
func (n *noopAdminService) Open(connectorName, instanceID string, config map[string]any, st *storage.Client, ac *activity.Client, logger *zap.Logger) (drivers.Handle, error) {
	return n, nil
}

// Spec implements [drivers.Driver].
func (n *noopAdminService) Spec() drivers.Spec {
	return drivers.Spec{
		ImplementsAdmin: true,
	}
}

// TertiarySourceConnectors implements [drivers.Driver].
func (n *noopAdminService) TertiarySourceConnectors(ctx context.Context, srcProps map[string]any, logger *zap.Logger) ([]string, error) {
	return nil, nil
}

// AsAI implements [drivers.Handle].
func (n *noopAdminService) AsAI(instanceID string) (drivers.AIService, bool) {
	return nil, false
}

// AsAdmin implements [drivers.Handle].
func (n *noopAdminService) AsAdmin(instanceID string) (drivers.AdminService, bool) {
	return &noopAdminClient{noopAdminService: n, instanceID: instanceID}, true
}

// AsCatalogStore implements [drivers.Handle].
func (n *noopAdminService) AsCatalogStore(instanceID string) (drivers.CatalogStore, bool) {
	return nil, false
}

// AsFileStore implements [drivers.Handle].
func (n *noopAdminService) AsFileStore() (drivers.FileStore, bool) {
	return nil, false
}

// AsInformationSchema implements [drivers.Handle].
func (n *noopAdminService) AsInformationSchema() (drivers.InformationSchema, bool) {
	return nil, false
}

// AsModelExecutor implements [drivers.Handle].
func (n *noopAdminService) AsModelExecutor(instanceID string, opts *drivers.ModelExecutorOptions) (drivers.ModelExecutor, error) {
	return nil, drivers.ErrNotImplemented
}

// AsModelManager implements [drivers.Handle].
func (n *noopAdminService) AsModelManager(instanceID string) (drivers.ModelManager, error) {
	return nil, drivers.ErrNotImplemented
}

// AsNotifier implements [drivers.Handle].
func (n *noopAdminService) AsNotifier(properties map[string]any) (drivers.Notifier, error) {
	return nil, drivers.ErrNotImplemented
}

// AsOLAP implements [drivers.Handle].
func (n *noopAdminService) AsOLAP(instanceID string) (drivers.OLAPStore, bool) {
	return nil, false
}

// AsObjectStore implements [drivers.Handle].
func (n *noopAdminService) AsObjectStore() (drivers.ObjectStore, bool) {
	return nil, false
}

// AsRegistry implements [drivers.Handle].
func (n *noopAdminService) AsRegistry() (drivers.RegistryStore, bool) {
	return nil, false
}

// AsRepoStore implements [drivers.Handle].
func (n *noopAdminService) AsRepoStore(instanceID string) (drivers.RepoStore, bool) {
	return nil, false
}

// AsWarehouse implements [drivers.Handle].
func (n *noopAdminService) AsWarehouse() (drivers.Warehouse, bool) {
	return nil, false
}

// Close implements [drivers.Handle].
func (n *noopAdminService) Close() error {
	return nil
}

// Config implements [drivers.Handle].
func (n *noopAdminService) Config() map[string]any {
	return map[string]any{}
}

// Driver implements [drivers.Handle].
func (n *noopAdminService) Driver() string {
	return "noop_admin"
}

// Migrate implements [drivers.Handle].
func (n *noopAdminService) Migrate(ctx context.Context) error {
	return nil
}

// MigrationStatus implements [drivers.Handle].
func (n *noopAdminService) MigrationStatus(ctx context.Context) (current, desired int, err error) {
	return 0, 0, nil
}

// Ping implements [drivers.Handle].
func (n *noopAdminService) Ping(ctx context.Context) error {
	return nil
}
