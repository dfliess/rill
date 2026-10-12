package mssql

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	// Load SQL Server driver (registers the "sqlserver" database/sql driver)
	_ "github.com/microsoft/go-mssqldb"
)

func init() {
	drivers.Register("mssql", driver{})
	drivers.RegisterAsConnector("mssql", driver{})
}

var spec = drivers.Spec{
	DisplayName: "SQL Server",
	Description: "Connect to Microsoft SQL Server.",
	ConfigProperties: []*drivers.PropertySpec{
		{
			Key:         "dsn",
			Type:        drivers.StringPropertyType,
			DisplayName: "SQL Server Connection String",
			Required:    true,
			DocsURL:     "https://github.com/microsoft/go-mssqldb#connection-parameters-and-dsn",
			Placeholder: "sqlserver://user:password@host:1433?database=my-db",
			Hint:        "Can be configured here or by setting the 'connector.mssql.dsn' environment variable (using '.env' or '--env')",
			Secret:      true,
		},
		{
			Key:         "user",
			Type:        drivers.StringPropertyType,
			DisplayName: "Username",
			Placeholder: "sa",
			Required:    true,
			Hint:        "SQL Server username for authentication",
		},
		{
			Key:         "password",
			Type:        drivers.StringPropertyType,
			DisplayName: "Password",
			Placeholder: "your_password",
			Hint:        "SQL Server password for authentication",
			Secret:      true,
		},
		{
			Key:         "host",
			Type:        drivers.StringPropertyType,
			DisplayName: "Host",
			Placeholder: "localhost",
			Required:    true,
			Hint:        "SQL Server hostname or IP address",
		},
		{
			Key:         "port",
			Type:        drivers.StringPropertyType,
			DisplayName: "Port",
			Placeholder: "1433",
			Default:     "1433",
			Hint:        "SQL Server port (default is 1433)",
		},
		{
			Key:         "database",
			Type:        drivers.StringPropertyType,
			DisplayName: "Database",
			Placeholder: "your_database",
			Required:    true,
			Hint:        "Name of the SQL Server database to connect to",
		},
		{
			Key:         "instance",
			Type:        drivers.StringPropertyType,
			DisplayName: "Instance",
			Placeholder: "SQLEXPRESS",
			Hint:        "Named instance to connect to (leave empty for the default instance)",
		},
		{
			Key:         "encrypt",
			Type:        drivers.StringPropertyType,
			DisplayName: "Encrypt",
			Placeholder: "mandatory",
			Hint:        "Options include disable, optional, mandatory or strict",
		},
		{
			Key:         "trust_server_certificate",
			Type:        drivers.BooleanPropertyType,
			DisplayName: "Trust Server Certificate",
			Default:     "false",
			Hint:        "Skip verification of the server's TLS certificate (useful for self-signed certificates)",
		},
		{
			Key:         "log_queries",
			Type:        drivers.BooleanPropertyType,
			DisplayName: "Log Queries",
			Default:     "false",
			Hint:        "Enable logging of all SQL queries (useful for debugging)",
		},
	},
	ImplementsSQLStore:  true,
	ImplementsOLAP:      true,
	ImplementsWarehouse: true,
}

type driver struct{}

type ConfigProperties struct {
	DSN                    string `mapstructure:"dsn"`
	Host                   string `mapstructure:"host"`
	Port                   int    `mapstructure:"port"`
	Database               string `mapstructure:"database"`
	Instance               string `mapstructure:"instance"`
	User                   string `mapstructure:"user"`
	Password               string `mapstructure:"password"`
	Encrypt                string `mapstructure:"encrypt"`
	TrustServerCertificate bool   `mapstructure:"trust_server_certificate"`
	LogQueries             bool   `mapstructure:"log_queries"`
}

// ResolveDSN returns a DSN in the URL format understood by go-mssqldb,
// e.g. sqlserver://user:password@host:1433?database=my-db&encrypt=mandatory.
func (c *ConfigProperties) ResolveDSN() (string, error) {
	if c.DSN != "" {
		if c.Host != "" || c.Port != 0 || c.Database != "" || c.Instance != "" || c.User != "" || c.Password != "" || c.Encrypt != "" || c.TrustServerCertificate {
			return "", fmt.Errorf("invalid config: DSN is set but other connection fields are also set")
		}
		return c.DSN, nil
	}

	var userInfo *url.Userinfo
	if c.User != "" {
		if c.Password != "" {
			userInfo = url.UserPassword(c.User, c.Password)
		} else {
			userInfo = url.User(c.User)
		}
	}

	host := c.Host
	if host == "" {
		host = "localhost"
	}
	if c.Port != 0 {
		host = fmt.Sprintf("%s:%d", host, c.Port)
	}

	// A named instance is encoded as the URL path, e.g. sqlserver://host/SQLEXPRESS.
	var path string
	if c.Instance != "" {
		path = "/" + c.Instance
	}

	query := url.Values{}
	if c.Database != "" {
		query.Set("database", c.Database)
	}
	if c.Encrypt != "" {
		switch strings.ToLower(c.Encrypt) {
		case "disable", "optional", "mandatory", "strict", "true", "false":
			query.Set("encrypt", strings.ToLower(c.Encrypt))
		default:
			return "", fmt.Errorf("unsupported encrypt: %s", c.Encrypt)
		}
	}
	if c.TrustServerCertificate {
		query.Set("trustservercertificate", strconv.FormatBool(c.TrustServerCertificate))
	}

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     userInfo,
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	return u.String(), nil
}

func (d driver) Open(_ context.Context, _, instanceID string, config map[string]any, st *storage.Client, ac *activity.Client, logger *zap.Logger) (drivers.Handle, error) {
	if instanceID == "" {
		return nil, errors.New("mssql driver can't be shared")
	}

	conf := &ConfigProperties{}
	if err := mapstructure.WeakDecode(config, conf); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &connection{
		config:     config,
		storage:    st,
		logger:     logger,
		logQueries: conf.LogQueries,
		dbMu:       semaphore.NewWeighted(1),
	}, nil
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
	config     map[string]any
	storage    *storage.Client
	logger     *zap.Logger
	logQueries bool

	db    *sqlx.DB // lazily populated using getDB
	dbErr error
	dbMu  *semaphore.Weighted
}

// Ping implements drivers.Handle.
func (c *connection) Ping(ctx context.Context) error {
	db, err := c.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to open SQL Server connection: %w", err)
	}
	return db.PingContext(ctx)
}

// Migrate implements drivers.Connection.
func (c *connection) Migrate(ctx context.Context) (err error) {
	return nil
}

// MigrationStatus implements drivers.Handle.
func (c *connection) MigrationStatus(ctx context.Context) (current, desired int, err error) {
	return 0, 0, nil
}

// Driver implements drivers.Connection.
func (c *connection) Driver() string {
	return "mssql"
}

// Config implements drivers.Connection.
func (c *connection) Config() map[string]any {
	return maps.Clone(c.config)
}

// Close implements drivers.Connection.
func (c *connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// AsRegistry implements drivers.Connection.
func (c *connection) AsRegistry() (drivers.RegistryStore, bool) {
	return nil, false
}

// AsCatalogStore implements drivers.Connection.
func (c *connection) AsCatalogStore(instanceID string) (drivers.CatalogStore, bool) {
	return nil, false
}

// AsRepoStore implements drivers.Connection.
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

// AsOLAP implements drivers.Connection.
func (c *connection) AsOLAP(instanceID string) (drivers.OLAPStore, bool) {
	return c, true
}

// AsInformationSchema implements drivers.Connection.
func (c *connection) AsInformationSchema() (drivers.InformationSchema, bool) {
	return c, true
}

// AsObjectStore implements drivers.Connection.
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

// AsFileStore implements drivers.Connection.
func (c *connection) AsFileStore() (drivers.FileStore, bool) {
	return nil, false
}

// AsWarehouse implements drivers.Handle.
func (c *connection) AsWarehouse() (drivers.Warehouse, bool) {
	return c, true
}

// AsNotifier implements drivers.Connection.
func (c *connection) AsNotifier(properties map[string]any) (drivers.Notifier, error) {
	return nil, drivers.ErrNotNotifier
}

// getDB lazily initializes and returns a database connection.
func (c *connection) getDB(ctx context.Context) (*sqlx.DB, error) {
	err := c.dbMu.Acquire(ctx, 1)
	if err != nil {
		return nil, err
	}
	defer c.dbMu.Release(1)

	if c.db != nil || c.dbErr != nil {
		return c.db, c.dbErr
	}

	conf := &ConfigProperties{}
	if err := mapstructure.WeakDecode(c.config, conf); err != nil {
		c.dbErr = fmt.Errorf("failed to decode config: %w", err)
		return nil, c.dbErr
	}

	dsn, err := conf.ResolveDSN()
	if err != nil {
		c.dbErr = err
		return nil, c.dbErr
	}

	c.db, c.dbErr = sqlx.Open("sqlserver", dsn)
	if c.dbErr != nil {
		return nil, c.dbErr
	}
	c.db.SetConnMaxIdleTime(time.Minute)

	return c.db, nil
}
