package mssql

import (
	"context"
	"database/sql"
	"fmt"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/pagination"
)

func (c *connection) ListDatabaseSchemas(ctx context.Context, pageSize uint32, pageToken string) ([]*drivers.DatabaseSchemaInfo, string, error) {
	limit := pagination.ValidPageSize(pageSize, drivers.DefaultPageSize)

	q := `
	SELECT
		schema_name
	FROM information_schema.schemata
	WHERE schema_name NOT IN (
		'sys', 'guest', 'INFORMATION_SCHEMA',
		'db_owner', 'db_accessadmin', 'db_securityadmin', 'db_ddladmin', 'db_backupoperator',
		'db_datareader', 'db_datawriter', 'db_denydatareader', 'db_denydatawriter'
	)
	`
	args := []any{}
	if pageToken != "" {
		var startAfter string
		if err := pagination.UnmarshalPageToken(pageToken, &startAfter); err != nil {
			return nil, "", fmt.Errorf("invalid page token: %w", err)
		}
		q += fmt.Sprintf("	AND schema_name > @p%d", len(args)+1)
		args = append(args, startAfter)
	}
	q += fmt.Sprintf(`
	ORDER BY schema_name
	OFFSET 0 ROWS FETCH NEXT @p%d ROWS ONLY
	`, len(args)+1)
	args = append(args, limit+1)

	db, err := c.getDB(ctx)
	if err != nil {
		return nil, "", err
	}

	rows, err := db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var res []*drivers.DatabaseSchemaInfo
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, "", err
		}
		res = append(res, &drivers.DatabaseSchemaInfo{
			Database:       "",
			DatabaseSchema: schema,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	next := ""
	if len(res) > limit {
		res = res[:limit]
		next = pagination.MarshalPageToken(res[len(res)-1].DatabaseSchema)
	}
	return res, next, nil
}

func (c *connection) ListTables(ctx context.Context, database, databaseSchema string, pageSize uint32, pageToken string) ([]*drivers.TableInfo, string, error) {
	limit := pagination.ValidPageSize(pageSize, drivers.DefaultPageSize)

	q := `
	SELECT
		table_name,
		CAST(CASE WHEN table_type = 'VIEW' THEN 1 ELSE 0 END AS BIT) AS mssql_is_view
	FROM information_schema.tables
	WHERE table_schema = @p1
	`
	args := []any{databaseSchema}
	if pageToken != "" {
		var startAfter string
		if err := pagination.UnmarshalPageToken(pageToken, &startAfter); err != nil {
			return nil, "", fmt.Errorf("invalid page token: %w", err)
		}
		q += fmt.Sprintf("	AND table_name > @p%d", len(args)+1)
		args = append(args, startAfter)
	}
	q += fmt.Sprintf(`
	ORDER BY table_name
	OFFSET 0 ROWS FETCH NEXT @p%d ROWS ONLY
	`, len(args)+1)
	args = append(args, limit+1)

	db, err := c.getDB(ctx)
	if err != nil {
		return nil, "", err
	}

	rows, err := db.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var res []*drivers.TableInfo
	for rows.Next() {
		var name string
		var typ bool
		if err := rows.Scan(&name, &typ); err != nil {
			return nil, "", err
		}
		res = append(res, &drivers.TableInfo{
			Name: name,
			View: typ,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	next := ""
	if len(res) > limit {
		res = res[:limit]
		next = pagination.MarshalPageToken(res[len(res)-1].Name)
	}
	return res, next, nil
}

func (c *connection) Lookup(ctx context.Context, database, databaseSchema, name string) (*drivers.OlapTable, error) {
	q := `
	SELECT
		CAST(CASE WHEN t.table_type = 'VIEW' THEN 1 ELSE 0 END AS BIT) AS mssql_is_view,
		c.column_name,
		c.data_type
	FROM information_schema.tables t
	JOIN information_schema.columns c
	ON t.table_schema = c.table_schema AND t.table_name = c.table_name
	WHERE c.table_schema = @p1 AND c.table_name = @p2
	ORDER BY c.ordinal_position
	`

	db, err := c.getDB(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryxContext(ctx, q, databaseSchema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var view bool
	var col, typ string
	fields := make([]*runtimev1.StructType_Field, 0)
	for rows.Next() {
		if err := rows.Scan(&view, &col, &typ); err != nil {
			return nil, err
		}
		fields = append(fields, &runtimev1.StructType_Field{
			Name: col,
			Type: databaseTypeToPB(typ, true),
		})
	}
	return &drivers.OlapTable{
		Database:       database,
		DatabaseSchema: databaseSchema,
		Name:           name,
		View:           view,
		Schema: &runtimev1.StructType{
			Fields: fields,
		},
		UnsupportedCols:   nil,
		PhysicalSizeBytes: 0,
	}, rows.Err()
}

// All implements drivers.OLAPInformationSchema.
func (c *connection) All(ctx context.Context, like string, pageSize uint32, pageToken string) ([]*drivers.OlapTable, string, error) {
	return drivers.AllFromInformationSchema(ctx, like, pageSize, pageToken, c)
}

// LoadPhysicalSize implements drivers.OLAPInformationSchema.
func (c *connection) LoadPhysicalSize(ctx context.Context, tables []*drivers.OlapTable) error {
	return nil
}

// LoadDDL implements drivers.OLAPInformationSchema.
// Note: table.Database is not used; in SQL Server, the database is determined by the connection.
func (c *connection) LoadDDL(ctx context.Context, table *drivers.OlapTable) error {
	db, err := c.getDB(ctx)
	if err != nil {
		return err
	}

	schema := table.DatabaseSchema
	if schema == "" {
		if err := db.QueryRowContext(ctx, "SELECT SCHEMA_NAME()").Scan(&schema); err != nil {
			return err
		}
	}

	if table.View {
		// For views: OBJECT_DEFINITION returns the original CREATE VIEW statement.
		var ddl sql.NullString
		q := "SELECT OBJECT_DEFINITION(OBJECT_ID(QUOTENAME(@p1) + '.' + QUOTENAME(@p2), 'V'))"
		err = db.QueryRowContext(ctx, q, schema, table.Name).Scan(&ddl)
		if err != nil {
			return err
		}
		table.DDL = ddl.String
		return nil
	}

	// SQL Server does not have a built-in way to get the DDL for a table, so we reconstruct a basic
	// CREATE TABLE statement from the available metadata (won't include indexes, constraints, etc.).
	q := `
	SELECT
		'CREATE TABLE ' + QUOTENAME(c.table_schema) + '.' + QUOTENAME(c.table_name) + ' (' +
		STRING_AGG(
			CAST(
				QUOTENAME(c.column_name) + ' ' + c.data_type +
				CASE
					WHEN c.data_type IN ('char', 'varchar', 'nchar', 'nvarchar', 'binary', 'varbinary') THEN
						'(' + CASE WHEN c.character_maximum_length = -1 THEN 'MAX' ELSE CAST(c.character_maximum_length AS varchar(10)) END + ')'
					WHEN c.data_type IN ('decimal', 'numeric') THEN
						'(' + CAST(c.numeric_precision AS varchar(10)) + ',' + CAST(c.numeric_scale AS varchar(10)) + ')'
					WHEN c.data_type IN ('datetime2', 'datetimeoffset', 'time') THEN
						'(' + CAST(c.datetime_precision AS varchar(10)) + ')'
					ELSE ''
				END +
				CASE WHEN c.is_nullable = 'NO' THEN ' NOT NULL' ELSE '' END
			AS nvarchar(MAX)),
			', '
		) WITHIN GROUP (ORDER BY c.ordinal_position) + ')'
	FROM information_schema.columns c
	WHERE c.table_schema = @p1 AND c.table_name = @p2
	GROUP BY c.table_schema, c.table_name
	`
	var ddl string
	err = db.QueryRowContext(ctx, q, schema, table.Name).Scan(&ddl)
	if err != nil {
		return err
	}
	table.DDL = ddl
	return nil
}
