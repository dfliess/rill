// Package odoo reads any model of an Odoo 19+ database through the JSON-2 API.
//
// It is a warehouse driver: a model names an Odoo model, the fields to read and an optional
// domain; the driver pages through search_read by id and writes one parquet file for the OLAP
// connector to ingest. It is generic on purpose — it knows nothing about invoices or stock —
// so a product built on it changes YAML, not Go (kairos-cloud ADR-0024).
package odoo

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/mitchellh/mapstructure"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"go.uber.org/zap"
)

func init() {
	drivers.Register("odoo", driver{})
	drivers.RegisterAsConnector("odoo", driver{})
}

var spec = drivers.Spec{
	DisplayName: "Odoo",
	Description: "Read Odoo 19+ models through the JSON-2 API.",
	ConfigProperties: []*drivers.PropertySpec{
		{Key: "url", Type: drivers.StringPropertyType, Required: true, DisplayName: "URL", Placeholder: "https://mycompany.odoo.com"},
		{Key: "database", Type: drivers.StringPropertyType, DisplayName: "Database", Hint: "Sent as X-Odoo-Database; needed when the server hosts more than one database."},
		{Key: "api_key", Type: drivers.StringPropertyType, Required: true, Secret: true, DisplayName: "API key"},
		{Key: "lang", Type: drivers.StringPropertyType, DisplayName: "Language", Placeholder: "es_ES", Hint: "Language of translated values, sent as context.lang."},
	},
	ImplementsWarehouse: true,
}

type driver struct{}

var _ drivers.Driver = driver{}

func (d driver) Open(_ context.Context, _, instanceID string, config map[string]any, st *storage.Client, _ *activity.Client, logger *zap.Logger) (drivers.Handle, error) {
	if instanceID == "" {
		return nil, errors.New("odoo driver can't be shared")
	}
	conf := &configProperties{}
	if err := mapstructure.WeakDecode(config, conf); err != nil {
		return nil, fmt.Errorf("odoo: invalid config: %w", err)
	}
	return &connection{config: config, conf: conf, storage: st, logger: logger}, nil
}

func (d driver) Spec() drivers.Spec {
	return spec
}

func (d driver) HasAnonymousSourceAccess(ctx context.Context, src map[string]any, logger *zap.Logger) (bool, error) {
	return false, nil
}

func (d driver) TertiarySourceConnectors(ctx context.Context, src map[string]any, logger *zap.Logger) ([]string, error) {
	return nil, nil
}

type connection struct {
	config  map[string]any
	conf    *configProperties
	storage *storage.Client
	logger  *zap.Logger
}

var _ drivers.Handle = &connection{}

// Ping checks the URL and the API key with a call any user can make.
func (c *connection) Ping(ctx context.Context) error {
	cl, err := newClient(c.conf)
	if err != nil {
		return err
	}
	var out map[string]any
	return cl.call(ctx, "res.users", "context_get", map[string]any{}, &out)
}

// Migrate implements drivers.Handle.
func (c *connection) Migrate(ctx context.Context) error {
	return nil
}

// MigrationStatus implements drivers.Handle.
func (c *connection) MigrationStatus(ctx context.Context) (current, desired int, err error) {
	return 0, 0, nil
}

// Driver implements drivers.Handle.
func (c *connection) Driver() string {
	return "odoo"
}

// Config implements drivers.Handle.
func (c *connection) Config() map[string]any {
	return maps.Clone(c.config)
}

// Close implements drivers.Handle.
func (c *connection) Close() error {
	return nil
}

// AsRegistry implements drivers.Handle.
func (c *connection) AsRegistry() (drivers.RegistryStore, bool) {
	return nil, false
}

// AsCatalogStore implements drivers.Handle.
func (c *connection) AsCatalogStore(instanceID string) (drivers.CatalogStore, bool) {
	return nil, false
}

// AsRepoStore implements drivers.Handle.
func (c *connection) AsRepoStore(instanceID string) (drivers.RepoStore, bool) {
	return nil, false
}

// AsAdmin implements drivers.Handle.
func (c *connection) AsAdmin(instanceID string) (drivers.AdminService, bool) {
	return nil, false
}

// AsAI implements drivers.Handle.
func (c *connection) AsAI(instanceID string) (drivers.AIService, bool) {
	return nil, false
}

// AsOLAP implements drivers.Handle.
func (c *connection) AsOLAP(instanceID string) (drivers.OLAPStore, bool) {
	return nil, false
}

// AsInformationSchema implements drivers.Handle.
func (c *connection) AsInformationSchema() (drivers.InformationSchema, bool) {
	return nil, false
}

// AsObjectStore implements drivers.Handle.
func (c *connection) AsObjectStore() (drivers.ObjectStore, bool) {
	return nil, false
}

// AsModelExecutor implements drivers.Handle.
func (c *connection) AsModelExecutor(instanceID string, opts *drivers.ModelExecutorOptions) (drivers.ModelExecutor, error) {
	return nil, drivers.ErrNotImplemented
}

// AsModelManager implements drivers.Handle.
func (c *connection) AsModelManager(instanceID string) (drivers.ModelManager, error) {
	return nil, drivers.ErrNotImplemented
}

// AsFileStore implements drivers.Handle.
func (c *connection) AsFileStore() (drivers.FileStore, bool) {
	return nil, false
}

// AsWarehouse implements drivers.Handle.
func (c *connection) AsWarehouse() (drivers.Warehouse, bool) {
	return c, true
}

// AsNotifier implements drivers.Handle.
func (c *connection) AsNotifier(properties map[string]any) (drivers.Notifier, error) {
	return nil, drivers.ErrNotNotifier
}
