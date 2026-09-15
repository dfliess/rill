package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"strings"
)

// Kairos keeps its own migration sequence, tracked apart from upstream's, so that upstream's next migration number
// can never collide with one Kairos already applied (kairos-cloud#158).

// Embed Kairos migrations directory in the binary
//
//go:embed migrations_kairos/*.sql
var kairosMigrationsFS embed.FS

// Name of the table that tracks Kairos migrations.
var kairosMigrationVersionTable = "kairos_migration_version"

// kairosLastSharedUpstreamVersion is the last upstream migration when Kairos started its own sequence
// (kairos-cloud#158). Until then Kairos numbered its migrations after upstream's: 0097, 0098 and 0099 were Kairos,
// and are now 0001, 0002 and 0003 in migrations_kairos. Never change it: the bootstrap relies on it to tell which
// versions recorded in admin_migration_version were Kairos, including in databases migrated long after upstream
// moves past 0096.
const kairosLastSharedUpstreamVersion = 96

// kairosMigrationsNumberedAsUpstream is how many Kairos migrations were numbered after upstream's (0097 to 0099).
// Never change it, for the same reason as kairosLastSharedUpstreamVersion.
const kairosMigrationsNumberedAsUpstream = 3

// migrationFS is what the migration loops need from a migrations directory.
type migrationFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

// bootstrapKairosMigrations creates kairosMigrationVersionTable the first time it runs. In a database that applied
// Kairos migrations as upstream 0097 to 0099, it moves them to kairosMigrationVersionTable and sets
// migrationVersionTable back to the last upstream version, so upstream's own 0097 is not skipped.
// It must run under the migration lock, before upstream migrations.
func (c *connection) bootstrapKairosMigrations(ctx context.Context) (err error) {
	// A single transaction: DDL is transactional in Postgres, so either both version tables change or neither does
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var exists int
	err = tx.QueryRowContext(ctx, "select count(*) from pg_catalog.pg_class where relname=$1 and relkind='r' and pg_table_is_visible(oid)", kairosMigrationVersionTable).Scan(&exists)
	if err != nil {
		return err
	}
	if exists != 0 {
		return nil
	}

	var upstreamExists int
	err = tx.QueryRowContext(ctx, "select count(*) from pg_catalog.pg_class where relname=$1 and relkind='r' and pg_table_is_visible(oid)", migrationVersionTable).Scan(&upstreamExists)
	if err != nil {
		return err
	}
	var upstreamVersion int
	if upstreamExists != 0 {
		err = tx.QueryRowContext(ctx, fmt.Sprintf("select version from %s", migrationVersionTable)).Scan(&upstreamVersion)
		if err != nil {
			return err
		}
	}

	kairosVersion := min(max(upstreamVersion-kairosLastSharedUpstreamVersion, 0), kairosMigrationsNumberedAsUpstream)

	_, err = tx.ExecContext(ctx, fmt.Sprintf("create table %s(version int4 not null)", kairosMigrationVersionTable))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf("insert into %s(version) values ($1)", kairosMigrationVersionTable), kairosVersion)
	if err != nil {
		return err
	}

	if upstreamVersion > kairosLastSharedUpstreamVersion {
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET version=$1", migrationVersionTable), kairosLastSharedUpstreamVersion)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// migrateKairos runs Kairos migrations. It must run under the migration lock, after bootstrapKairosMigrations.
func (c *connection) migrateKairos(ctx context.Context) error {
	// Get version of latest Kairos migration
	var currentVersion int
	err := c.db.QueryRowContext(ctx, fmt.Sprintf("select version from %s", kairosMigrationVersionTable)).Scan(&currentVersion)
	if err != nil {
		return err
	}

	// Iterate over migrations (sorted by filename)
	files, err := kairosMigrationsFS.ReadDir("migrations_kairos")
	if err != nil {
		return err
	}

	for _, file := range files {
		// Extract version number from filename
		version, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".sql"))
		if err != nil {
			return fmt.Errorf("unexpected Kairos migration filename: %s", file.Name())
		}

		// Skip migrations below current version
		if version <= currentVersion {
			continue
		}

		// Read SQL
		sql, err := kairosMigrationsFS.ReadFile(path.Join("migrations_kairos", file.Name()))
		if err != nil {
			return err
		}

		err = c.migrateSingleKairos(ctx, file.Name(), sql, version)
		if err != nil {
			return err
		}
	}

	return nil
}

// migrateSingleKairos is migrateSingle for kairosMigrationVersionTable.
func (c *connection) migrateSingleKairos(ctx context.Context, name string, sql []byte, version int) (err error) {
	// Start a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Run migration
	_, err = tx.ExecContext(ctx, string(sql))
	if err != nil {
		return fmt.Errorf("failed to run Kairos migration '%s': %w", name, err)
	}

	// Update migration version
	_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET version=$1", kairosMigrationVersionTable), version)
	if err != nil {
		return err
	}

	// Commit migration
	return tx.Commit()
}
