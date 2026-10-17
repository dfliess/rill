package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/pkg/pgtestcontainer"
	"github.com/stretchr/testify/require"
)

// TestKairosMigrations checks that Kairos migrations end up in their own sequence whatever state the database was in
// before kairos-cloud#158, and that upstream migrations published after 0096 still run. Each case uses its own
// database in a single Postgres container.
func TestKairosMigrations(t *testing.T) {
	pg := pgtestcontainer.New(t)
	defer pg.Terminate(t)

	t.Run("NewDatabase", func(t *testing.T) {
		ctx := context.Background()
		c := openKairosTestDatabase(t, pg.DatabaseURL, "new_database")

		require.NoError(t, c.Migrate(ctx))
		requireKairosMigrated(t, c, latestUpstreamVersion(t))

		// Migrating again runs nothing: a preference stored in between survives 0002's DROP TABLE
		userID, orgID := insertNotificationPreference(t, c)
		require.NoError(t, c.Migrate(ctx))
		requireKairosMigrated(t, c, latestUpstreamVersion(t))
		requireNotificationPreference(t, c, userID, orgID)
	})

	t.Run("UpstreamOnly", func(t *testing.T) {
		ctx := context.Background()
		c := openKairosTestDatabase(t, pg.DatabaseURL, "upstream_only")

		// A database migrated to 0096 by a binary that had no Kairos migrations
		applyMigrationsBefore158(t, c, 96)
		requireTableExists(t, c, "push_subscriptions", false)

		require.NoError(t, c.Migrate(ctx))
		requireKairosMigrated(t, c, latestUpstreamVersion(t))

		userID, orgID := insertNotificationPreference(t, c)
		require.NoError(t, c.Migrate(ctx))
		requireKairosMigrated(t, c, latestUpstreamVersion(t))
		requireNotificationPreference(t, c, userID, orgID)
	})

	t.Run("KairosNumberedAsUpstream", func(t *testing.T) {
		ctx := context.Background()
		c := openKairosTestDatabase(t, pg.DatabaseURL, "kairos_numbered_as_upstream")

		// A database migrated to 0099 before kairos-cloud#158, when Kairos migrations were upstream 0097 to 0099
		applyMigrationsBefore158(t, c, 99)
		userID, orgID := insertNotificationPreferenceBefore158(t, c)

		// Bootstrapping takes upstream back to 0096, so upstream's own migrations from 0097 on run instead of being skipped
		require.NoError(t, c.Migrate(ctx))
		requireKairosMigrated(t, c, latestUpstreamVersion(t))
		requireNotificationPreference(t, c, userID, orgID)

		require.NoError(t, c.Migrate(ctx))
		requireKairosMigrated(t, c, latestUpstreamVersion(t))
		requireNotificationPreference(t, c, userID, orgID)

		t.Run("NewUpstreamMigration", func(t *testing.T) {
			// Upstream publishes a migration after the last one this binary embeds
			next := latestUpstreamVersion(t) + 1
			upstream := fstest.MapFS{}
			files, err := embeddedMigrationsFS.ReadDir("migrations")
			require.NoError(t, err)
			for _, file := range files {
				data, err := embeddedMigrationsFS.ReadFile(path.Join("migrations", file.Name()))
				require.NoError(t, err)
				upstream[path.Join("migrations", file.Name())] = &fstest.MapFile{Data: data}
			}
			upstream[fmt.Sprintf("migrations/%04d.sql", next)] = &fstest.MapFile{Data: []byte("CREATE TABLE upstream_next (id INT);")}
			migrationsFS = upstream
			t.Cleanup(func() { migrationsFS = embeddedMigrationsFS })

			require.NoError(t, c.Migrate(ctx))
			requireKairosMigrated(t, c, next)
			requireTableExists(t, c, "upstream_next", true)
			requireNotificationPreference(t, c, userID, orgID)

			// The next start does not take upstream back to 0096
			require.NoError(t, c.Migrate(ctx))
			requireKairosMigrated(t, c, next)
		})
	})

	t.Run("UnknownVersionAboveKairos", func(t *testing.T) {
		ctx := context.Background()
		c := openKairosTestDatabase(t, pg.DatabaseURL, "unknown_version_above_kairos")

		// A database at 0100 before kairos-cloud#158: no binary had a migration numbered 0100
		applyMigrationsBefore158(t, c, 99)
		_, err := c.db.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET version=100", migrationVersionTable))
		require.NoError(t, err)

		err = c.Migrate(ctx)
		require.ErrorContains(t, err, "is at version 100")
		require.ErrorContains(t, err, "requires manual intervention")

		// Nothing changed
		version, err := c.FindMigrationVersion(ctx)
		require.NoError(t, err)
		require.Equal(t, 100, version)
		requireTableExists(t, c, kairosMigrationVersionTable, false)
	})
}

// openKairosTestDatabase creates a database named name in the Postgres at databaseURL and connects to it.
func openKairosTestDatabase(t *testing.T, databaseURL, name string) *connection {
	admin, err := sql.Open("pgx", databaseURL)
	require.NoError(t, err)
	defer admin.Close()
	_, err = admin.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", name))
	require.NoError(t, err)

	db, err := driver{}.Open(strings.TrimSuffix(databaseURL, "/postgres")+"/"+name, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db.(*connection)
}

// applyMigrationsBefore158 migrates like Migrate did before kairos-cloud#158, when Kairos migrations followed upstream's
// in admin_migration_version: upstream migrations, then Kairos migrations as 0097, 0098 and 0099, up to version upTo.
func applyMigrationsBefore158(t *testing.T, c *connection, upTo int) {
	ctx := context.Background()

	_, err := c.db.ExecContext(ctx, fmt.Sprintf("create table %s(version int4 not null)", migrationVersionTable))
	require.NoError(t, err)
	_, err = c.db.ExecContext(ctx, fmt.Sprintf("insert into %s(version) values (0)", migrationVersionTable))
	require.NoError(t, err)

	apply := func(fsys migrationFS, dir string, offset, upTo int) {
		files, err := fsys.ReadDir(dir)
		require.NoError(t, err)
		for _, file := range files {
			var version int
			_, err := fmt.Sscanf(file.Name(), "%04d.sql", &version)
			require.NoError(t, err)
			version += offset
			if version > upTo {
				return
			}
			data, err := fs.ReadFile(fsys, path.Join(dir, file.Name()))
			require.NoError(t, err)
			require.NoError(t, c.migrateSingle(ctx, file.Name(), data, version))
		}
	}
	// Binaries from before kairos-cloud#158 had upstream migrations only up to 0096: later ones did not exist yet
	apply(embeddedMigrationsFS, "migrations", 0, min(upTo, kairosLastSharedUpstreamVersion))
	apply(kairosMigrationsFS, "migrations_kairos", kairosLastSharedUpstreamVersion, upTo)

	version, err := c.FindMigrationVersion(ctx)
	require.NoError(t, err)
	require.Equal(t, upTo, version)
}

// requireKairosMigrated checks upstream is at upstreamVersion and every Kairos migration has run.
func requireKairosMigrated(t *testing.T, c *connection, upstreamVersion int) {
	ctx := context.Background()

	version, err := c.FindMigrationVersion(ctx)
	require.NoError(t, err)
	require.Equal(t, upstreamVersion, version)

	var kairosVersion int
	require.NoError(t, c.db.QueryRowContext(ctx, fmt.Sprintf("select version from %s", kairosMigrationVersionTable)).Scan(&kairosVersion))
	require.Equal(t, 4, kairosVersion)

	for _, table := range []string{"push_subscriptions", "notification_preferences", "notification_deliveries", "org_styles"} {
		requireTableExists(t, c, table, true)
	}
}

// latestUpstreamVersion is the number of the last upstream migration this binary embeds.
// It moves every time upstream publishes a migration, so tests must not hardcode it.
func latestUpstreamVersion(t *testing.T) int {
	files, err := embeddedMigrationsFS.ReadDir("migrations")
	require.NoError(t, err)
	var version int
	_, err = fmt.Sscanf(files[len(files)-1].Name(), "%04d.sql", &version)
	require.NoError(t, err)
	return version
}

func requireTableExists(t *testing.T, c *connection, table string, exists bool) {
	var found bool
	require.NoError(t, c.db.QueryRowContext(context.Background(), "select to_regclass($1) is not null", table).Scan(&found))
	require.Equal(t, exists, found, table)
}

// insertNotificationPreference stores a notification preference that differs from the defaults.
func insertNotificationPreference(t *testing.T, c *connection) (userID, orgID string) {
	ctx := context.Background()

	u, err := c.InsertUser(ctx, &database.InsertUserOptions{Email: randomName() + "@rilldata.com"})
	require.NoError(t, err)
	org, err := c.InsertOrganization(ctx, &database.InsertOrganizationOptions{Name: randomName()})
	require.NoError(t, err)
	_, err = c.UpsertNotificationPreferences(ctx, u.ID, org.ID, &database.UpsertNotificationPreferencesOptions{PushAlerts: false, PushReports: true, PushActApprovals: false})
	require.NoError(t, err)
	return u.ID, org.ID
}

// insertNotificationPreferenceBefore158 is insertNotificationPreference with plain SQL, for a database still at its schema
// from before kairos-cloud#158: the current queries write columns that upstream added since, which it does not have yet.
func insertNotificationPreferenceBefore158(t *testing.T, c *connection) (userID, orgID string) {
	ctx := context.Background()

	require.NoError(t, c.db.QueryRowContext(ctx, "INSERT INTO users (email, display_name) VALUES ($1, '') RETURNING id", randomName()+"@rilldata.com").Scan(&userID))
	require.NoError(t, c.db.QueryRowContext(ctx, "INSERT INTO orgs (name, description) VALUES ($1, '') RETURNING id", randomName()).Scan(&orgID))
	_, err := c.db.ExecContext(ctx, "INSERT INTO notification_preferences (user_id, org_id, push_alerts, push_reports, push_act_approvals) VALUES ($1, $2, false, true, false)", userID, orgID)
	require.NoError(t, err)
	return userID, orgID
}

func requireNotificationPreference(t *testing.T, c *connection, userID, orgID string) {
	var count int
	require.NoError(t, c.db.QueryRowContext(context.Background(), "select count(*) from notification_preferences where user_id=$1 and org_id=$2 and not push_alerts", userID, orgID).Scan(&count))
	require.Equal(t, 1, count)
}
