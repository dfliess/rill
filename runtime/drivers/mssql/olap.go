package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	mssqldb "github.com/microsoft/go-mssqldb"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.uber.org/zap"
)

var _ drivers.OLAPStore = (*connection)(nil)

// Dialect implements drivers.OLAPStore.
func (c *connection) Dialect() drivers.Dialect {
	return DialectMSSQL
}

// Exec implements drivers.OLAPStore.
func (c *connection) Exec(ctx context.Context, stmt *drivers.Statement) error {
	if stmt.DryRun {
		return nil
	}
	res, err := c.Query(ctx, stmt)
	if err != nil {
		return err
	}
	if res != nil {
		return res.Close()
	}
	return nil
}

// InformationSchema implements drivers.OLAPStore.
func (c *connection) InformationSchema() drivers.InformationSchema {
	return c
}

// EstimateSize implements drivers.OLAPStore.
func (c *connection) EstimateSize(ctx context.Context) (int64, error) {
	return -1, nil
}

// MayBeScaledToZero implements drivers.OLAPStore.
func (c *connection) MayBeScaledToZero(ctx context.Context) bool {
	return false // SQL Server instances are typically always running
}

// Query implements drivers.OLAPStore.
func (c *connection) Query(ctx context.Context, stmt *drivers.Statement) (res *drivers.Result, resErr error) {
	if c.logQueries {
		fields := []zap.Field{
			zap.String("sql", c.Dialect().SanitizeQueryForLogging(stmt.Query)),
			zap.Any("args", stmt.Args),
			observability.ZapCtx(ctx),
		}
		if len(stmt.QueryAttributes) > 0 {
			fields = append(fields, zap.Any("query_attributes", stmt.QueryAttributes))
		}
		c.logger.Info("SQL Server query", fields...)
	}

	db, err := c.getDB(ctx)
	if err != nil {
		return nil, err
	}

	if stmt.DryRun {
		// SQL Server has no EXPLAIN; sp_describe_first_result_set compiles the query and returns
		// its result schema without executing it. Note: it can't validate queries that reference
		// undeclared parameters, so stmt.Args are ignored here.
		_, err = db.ExecContext(ctx, "EXEC sp_describe_first_result_set @tsql = @p1", stmt.Query)
		return nil, err
	}

	rows, err := db.QueryxContext(ctx, stmt.Query, stmt.Args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resErr != nil {
			_ = rows.Close()
		}
	}()

	schema, err := rowsToSchema(rows)
	if err != nil {
		return nil, err
	}

	cts, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	msSQLRows := &mssqlRows{
		Rows:     rows,
		scanDest: prepareScanDest(schema),
		colTypes: cts,
	}
	res = &drivers.Result{Rows: msSQLRows, Schema: schema}
	return res, nil
}

func (c *connection) Head(ctx context.Context, db, schema, table string, limit int64) (*drivers.Result, error) {
	tbl, err := c.InformationSchema().Lookup(ctx, db, schema, table)
	if err != nil {
		return nil, err
	}

	var columns []string
	for _, field := range tbl.Schema.Fields {
		columns = append(columns, c.Dialect().EscapeIdentifier(field.Name))
	}

	// SQL Server has no LIMIT; use TOP.
	topClause := ""
	if limit > 0 {
		topClause = fmt.Sprintf("TOP %d ", limit)
	}

	return c.Query(ctx, &drivers.Statement{
		Query: fmt.Sprintf("SELECT %s%s FROM %s", topClause, strings.Join(columns, ", "), c.Dialect().EscapeTable(db, schema, table)),
	})
}

// QuerySchema implements drivers.OLAPStore.
func (c *connection) QuerySchema(ctx context.Context, query string, args []any) (*runtimev1.StructType, error) {
	return nil, drivers.ErrNotImplemented
}

// WithConnection implements drivers.OLAPStore.
func (c *connection) WithConnection(ctx context.Context, priority int, fn drivers.WithConnectionFunc) error {
	return drivers.ErrNotImplemented
}

func rowsToSchema(r *sqlx.Rows) (*runtimev1.StructType, error) {
	if r == nil {
		return nil, nil
	}

	cts, err := r.ColumnTypes()
	if err != nil {
		return nil, err
	}

	fields := make([]*runtimev1.StructType_Field, len(cts))
	for i, ct := range cts {
		nullable, ok := ct.Nullable()
		if !ok {
			nullable = true
		}

		fields[i] = &runtimev1.StructType_Field{
			Name: ct.Name(),
			Type: databaseTypeToPB(ct.DatabaseTypeName(), nullable),
		}
	}

	return &runtimev1.StructType{Fields: fields}, nil
}

func databaseTypeToPB(dbt string, nullable bool) *runtimev1.Type {
	t := &runtimev1.Type{Nullable: nullable, RawType: dbt}
	switch strings.ToUpper(dbt) {
	case "TINYINT": // unsigned (0-255) in SQL Server, unlike MySQL
		t.Code = runtimev1.Type_CODE_UINT8
	case "SMALLINT":
		t.Code = runtimev1.Type_CODE_INT16
	case "INT", "INTEGER":
		t.Code = runtimev1.Type_CODE_INT32
	case "BIGINT":
		t.Code = runtimev1.Type_CODE_INT64
	case "BIT":
		t.Code = runtimev1.Type_CODE_BOOL
	case "FLOAT", "REAL":
		t.Code = runtimev1.Type_CODE_FLOAT64
	case "DECIMAL", "NUMERIC", "MONEY", "SMALLMONEY":
		t.Code = runtimev1.Type_CODE_DECIMAL
	case "CHAR", "VARCHAR", "TEXT", "NCHAR", "NVARCHAR", "NTEXT", "XML":
		t.Code = runtimev1.Type_CODE_STRING
	case "BINARY", "VARBINARY", "IMAGE":
		t.Code = runtimev1.Type_CODE_STRING
	case "DATE":
		t.Code = runtimev1.Type_CODE_DATE
	case "TIME":
		t.Code = runtimev1.Type_CODE_TIME
	case "DATETIME", "DATETIME2", "SMALLDATETIME", "DATETIMEOFFSET":
		t.Code = runtimev1.Type_CODE_TIMESTAMP
	case "UNIQUEIDENTIFIER":
		t.Code = runtimev1.Type_CODE_UUID
	default:
		t.Code = runtimev1.Type_CODE_UNSPECIFIED
	}
	return t
}

// mssqlRows wraps sqlx.Rows to provide a MapScan method.
// This is required because scanning into any leaves DECIMAL/MONEY values as raw byte slices
// and UNIQUEIDENTIFIER values as raw 16-byte arrays (in mixed byte order).
// Scanning into typed destinations derived from the schema yields clean Go values.
type mssqlRows struct {
	*sqlx.Rows
	scanDest []any
	colTypes []*sql.ColumnType
}

func (r *mssqlRows) MapScan(dest map[string]any) error {
	err := r.Rows.Scan(r.scanDest...)
	if err != nil {
		return err
	}
	for i, ct := range r.colTypes {
		fieldName := ct.Name()
		valPtr := r.scanDest[i]
		if valPtr == nil {
			dest[fieldName] = nil
			continue
		}
		switch valPtr := valPtr.(type) {
		case *sql.NullBool:
			if valPtr.Valid {
				dest[fieldName] = valPtr.Bool
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullByte:
			if valPtr.Valid {
				dest[fieldName] = valPtr.Byte
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullInt16:
			if valPtr.Valid {
				dest[fieldName] = valPtr.Int16
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullInt32:
			if valPtr.Valid {
				dest[fieldName] = valPtr.Int32
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullInt64:
			if valPtr.Valid {
				dest[fieldName] = valPtr.Int64
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullFloat64:
			if valPtr.Valid {
				dest[fieldName] = valPtr.Float64
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullString:
			if valPtr.Valid {
				dest[fieldName] = valPtr.String
			} else {
				dest[fieldName] = nil
			}
		case *sql.NullTime:
			if valPtr.Valid {
				if strings.ToUpper(r.colTypes[i].DatabaseTypeName()) == "TIME" {
					// The driver returns TIME as a time.Time on a zero date; format just the time part.
					dest[fieldName] = valPtr.Time.Format("15:04:05.9999999")
				} else {
					dest[fieldName] = valPtr.Time
				}
			} else {
				dest[fieldName] = nil
			}
		case *mssqldb.NullUniqueIdentifier:
			if valPtr.Valid {
				dest[fieldName] = valPtr.UUID.String()
			} else {
				dest[fieldName] = nil
			}
		default:
			dest[fieldName] = *(valPtr.(*any))
		}
	}
	return nil
}

func prepareScanDest(schema *runtimev1.StructType) []any {
	scanList := make([]any, len(schema.Fields))
	for i, field := range schema.Fields {
		var dest any
		switch field.Type.Code {
		case runtimev1.Type_CODE_BOOL:
			dest = &sql.NullBool{}
		case runtimev1.Type_CODE_UINT8:
			dest = &sql.NullByte{}
		case runtimev1.Type_CODE_INT16:
			dest = &sql.NullInt16{}
		case runtimev1.Type_CODE_INT32:
			dest = &sql.NullInt32{}
		case runtimev1.Type_CODE_INT64:
			dest = &sql.NullInt64{}
		case runtimev1.Type_CODE_FLOAT64:
			dest = &sql.NullFloat64{}
		case runtimev1.Type_CODE_DECIMAL:
			// scan decimals into strings to avoid precision loss
			dest = &sql.NullString{}
		case runtimev1.Type_CODE_STRING:
			dest = &sql.NullString{}
		case runtimev1.Type_CODE_DATE, runtimev1.Type_CODE_TIME, runtimev1.Type_CODE_TIMESTAMP:
			dest = &sql.NullTime{}
		case runtimev1.Type_CODE_UUID:
			dest = &mssqldb.NullUniqueIdentifier{}
		default:
			dest = new(any)
		}
		scanList[i] = dest
	}
	return scanList
}
