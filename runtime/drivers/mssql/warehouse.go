package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/decimal128"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/apache/arrow-go/v18/parquet"
	"github.com/apache/arrow-go/v18/parquet/compress"
	"github.com/apache/arrow-go/v18/parquet/pqarrow"
	"github.com/c2h5oh/datasize"
	mssqldb "github.com/microsoft/go-mssqldb"
	"github.com/mitchellh/mapstructure"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("github.com/rilldata/rill/runtime/drivers/mssql")

// entire data is buffered in memory before its written to disk so keeping it small reduces memory usage
// but keeping it too small can lead to bad ingestion performance
// 64MB seems to be a good balance
const rowGroupBufferSize = int64(datasize.MB) * 64

// number of rows accumulated in an arrow record before it is flushed to the parquet writer
const batchSizeRows = 8192

// QueryAsFiles implements drivers.Warehouse.
// It executes the query against SQL Server and streams the result into a parquet file in batches.
func (c *connection) QueryAsFiles(ctx context.Context, props map[string]any) (outIt drivers.FileIterator, outErr error) {
	ctx, span := tracer.Start(ctx, "Connection.QueryAsFiles")
	defer func() {
		if outErr != nil {
			span.SetStatus(codes.Error, outErr.Error())
		}
		span.End()
	}()

	srcProps, err := parseSourceProperties(props)
	if err != nil {
		return nil, err
	}

	dsn := srcProps.DSN
	if dsn == "" {
		conf := &ConfigProperties{}
		if err := mapstructure.WeakDecode(c.config, conf); err != nil {
			return nil, fmt.Errorf("failed to decode config: %w", err)
		}
		dsn, err = conf.ResolveDSN()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	defer func() {
		if outErr != nil {
			db.Close()
		}
	}()

	rows, err := db.QueryContext(ctx, srcProps.SQL) //nolint:rowserrcheck // rows.Err is checked in fileIterator.Next
	if err != nil {
		return nil, err
	}
	defer func() {
		if outErr != nil {
			rows.Close()
		}
	}()

	cts, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	schema, appenders, err := arrowSchemaFor(cts)
	if err != nil {
		return nil, err
	}

	tempDir, err := c.storage.RandomTempDir("mssql")
	if err != nil {
		return nil, err
	}
	return &fileIterator{
		db:        db,
		rows:      rows,
		schema:    schema,
		appenders: appenders,
		logger:    c.logger,
		tempDir:   tempDir,
	}, nil
}

type fileIterator struct {
	db        *sql.DB
	rows      *sql.Rows
	schema    *arrow.Schema
	appenders []appendFunc
	logger    *zap.Logger
	tempDir   string
	// Computed while iterating
	totalRecords int64
	downloaded   bool
}

var _ drivers.FileIterator = &fileIterator{}

// Close implements drivers.FileIterator.
func (f *fileIterator) Close() error {
	if f.rows != nil {
		f.rows.Close()
		f.db.Close()
		f.rows = nil
	}
	return os.RemoveAll(f.tempDir)
}

// Format implements drivers.FileIterator.
func (f *fileIterator) Format() string {
	return ""
}

// SetKeepFilesUntilClose implements drivers.FileIterator.
func (f *fileIterator) SetKeepFilesUntilClose() {
	// No-op because it already does this.
}

// Next implements drivers.FileIterator.
// Query result is written to a single parquet file in batches of rows.
func (f *fileIterator) Next(ctx context.Context) ([]string, error) {
	if f.downloaded {
		return nil, io.EOF
	}

	ctx, span := tracer.Start(ctx, "fileIterator.Next")
	defer span.End()

	// close db resources early
	defer func() {
		f.rows.Close()
		f.db.Close()
		// mark rows as nil to prevent double close
		f.rows = nil
	}()

	f.logger.Debug("downloading results in parquet file", observability.ZapCtx(ctx))

	// create a temp file
	fw, err := os.CreateTemp(f.tempDir, "temp*.parquet")
	if err != nil {
		return nil, err
	}
	defer fw.Close()
	f.downloaded = true

	tf := time.Now()
	defer func() {
		f.logger.Debug("time taken to write arrow records in parquet file", zap.Duration("duration", time.Since(tf)), observability.ZapCtx(ctx))
	}()

	writer, err := pqarrow.NewFileWriter(f.schema, fw,
		parquet.NewWriterProperties(
			parquet.WithCompression(compress.Codecs.Snappy),
			parquet.WithRootRepetition(parquet.Repetitions.Required),
			// duckdb has issues reading statistics of string type generated with this write
			// column statistics may not be useful if full file need to be ingested so better to disable to save computations
			parquet.WithStats(false),
		),
		pqarrow.NewArrowWriterProperties(pqarrow.WithStoreSchema()))
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	builder := array.NewRecordBuilder(memory.DefaultAllocator, f.schema)
	defer builder.Release()

	scanDest := make([]any, len(f.appenders))
	for i, a := range f.appenders {
		scanDest[i] = a.newScanDest()
	}

	flush := func() error {
		rec := builder.NewRecordBatch()
		defer rec.Release()
		if writer.RowGroupTotalBytesWritten() >= rowGroupBufferSize {
			writer.NewBufferedRowGroup()
		}
		return writer.WriteBuffered(rec)
	}

	rowsInBatch := 0
	for f.rows.Next() {
		// stop promptly if the context is cancelled; the query itself is also bound to the
		// context passed to QueryAsFiles, which cancels it server-side
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := f.rows.Scan(scanDest...); err != nil {
			return nil, err
		}
		for i, a := range f.appenders {
			if err := a.append(builder.Field(i), scanDest[i]); err != nil {
				return nil, fmt.Errorf("column %q: %w", f.schema.Field(i).Name, err)
			}
		}
		rowsInBatch++
		f.totalRecords++
		if rowsInBatch >= batchSizeRows {
			if err := flush(); err != nil {
				return nil, err
			}
			rowsInBatch = 0
		}
	}
	if err := f.rows.Err(); err != nil {
		return nil, err
	}
	if rowsInBatch > 0 {
		if err := flush(); err != nil {
			return nil, err
		}
	}

	if f.totalRecords == 0 {
		// empty result
		return nil, drivers.ErrNoRows
	}

	writer.Close()
	fw.Close()

	fileInfo, err := os.Stat(fw.Name())
	if err != nil {
		return nil, err
	}
	f.logger.Debug("size of file", zap.String("size", datasize.ByteSize(fileInfo.Size()).HumanReadable()), observability.ZapCtx(ctx))
	return []string{fw.Name()}, nil
}

type sourceProperties struct {
	SQL string `mapstructure:"sql"`
	DSN string `mapstructure:"dsn"`
}

func parseSourceProperties(props map[string]any) (*sourceProperties, error) {
	conf := &sourceProperties{}
	err := mapstructure.Decode(props, conf)
	if err != nil {
		return nil, err
	}
	if conf.SQL == "" {
		return nil, fmt.Errorf("property 'sql' is mandatory for connector \"mssql\"")
	}
	return conf, err
}

// appendFunc converts one scanned SQL Server value into an arrow builder cell.
type appendFunc struct {
	newScanDest func() any
	append      func(b array.Builder, src any) error
}

// arrowSchemaFor maps the result set's SQL Server column types to an arrow schema and matching
// per-column append functions. Types without a clean, lossless mapping return an error rather
// than silently degrading (e.g. sql_variant, spatial and other UDTs).
func arrowSchemaFor(cts []*sql.ColumnType) (*arrow.Schema, []appendFunc, error) {
	fields := make([]arrow.Field, len(cts))
	appenders := make([]appendFunc, len(cts))
	for i, ct := range cts {
		nullable, ok := ct.Nullable()
		if !ok {
			nullable = true
		}

		var typ arrow.DataType
		var app appendFunc
		dbt := strings.ToUpper(ct.DatabaseTypeName())
		switch dbt {
		case "TINYINT": // unsigned (0-255) in SQL Server
			typ = arrow.PrimitiveTypes.Uint8
			app = appendFunc{
				newScanDest: func() any { return &sql.NullByte{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullByte)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Uint8Builder).Append(v.Byte)
					return nil
				},
			}
		case "SMALLINT":
			typ = arrow.PrimitiveTypes.Int16
			app = appendFunc{
				newScanDest: func() any { return &sql.NullInt16{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullInt16)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Int16Builder).Append(v.Int16)
					return nil
				},
			}
		case "INT", "INTEGER":
			typ = arrow.PrimitiveTypes.Int32
			app = appendFunc{
				newScanDest: func() any { return &sql.NullInt32{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullInt32)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Int32Builder).Append(v.Int32)
					return nil
				},
			}
		case "BIGINT":
			typ = arrow.PrimitiveTypes.Int64
			app = appendFunc{
				newScanDest: func() any { return &sql.NullInt64{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullInt64)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Int64Builder).Append(v.Int64)
					return nil
				},
			}
		case "BIT":
			typ = arrow.FixedWidthTypes.Boolean
			app = appendFunc{
				newScanDest: func() any { return &sql.NullBool{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullBool)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.BooleanBuilder).Append(v.Bool)
					return nil
				},
			}
		case "FLOAT":
			typ = arrow.PrimitiveTypes.Float64
			app = appendFunc{
				newScanDest: func() any { return &sql.NullFloat64{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullFloat64)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Float64Builder).Append(v.Float64)
					return nil
				},
			}
		case "REAL": // float32 in SQL Server
			typ = arrow.PrimitiveTypes.Float32
			app = appendFunc{
				newScanDest: func() any { return &sql.NullFloat64{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullFloat64)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Float32Builder).Append(float32(v.Float64))
					return nil
				},
			}
		case "DECIMAL", "NUMERIC", "MONEY", "SMALLMONEY":
			var prec, scale int64
			switch dbt {
			case "MONEY": // fixed decimal(19,4) in SQL Server; the driver does not report its precision
				prec, scale = 19, 4
			case "SMALLMONEY": // fixed decimal(10,4)
				prec, scale = 10, 4
			default:
				var okPS bool
				prec, scale, okPS = ct.DecimalSize()
				if !okPS {
					return nil, nil, fmt.Errorf("column %q: driver did not report precision and scale for type %s", ct.Name(), dbt)
				}
			}
			typ = &arrow.Decimal128Type{Precision: int32(prec), Scale: int32(scale)}
			app = appendFunc{
				newScanDest: func() any { return &sql.NullString{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullString)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					num, err := decimalToNum(v.String, int32(prec), int32(scale))
					if err != nil {
						return err
					}
					b.(*array.Decimal128Builder).Append(num)
					return nil
				},
			}
		case "CHAR", "VARCHAR", "TEXT", "NCHAR", "NVARCHAR", "NTEXT", "XML":
			typ = arrow.BinaryTypes.String
			app = appendFunc{
				newScanDest: func() any { return &sql.NullString{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullString)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.StringBuilder).Append(v.String)
					return nil
				},
			}
		case "BINARY", "VARBINARY", "IMAGE":
			typ = arrow.BinaryTypes.Binary
			app = appendFunc{
				newScanDest: func() any { return new([]byte) },
				append: func(b array.Builder, src any) error {
					v := src.(*[]byte)
					if *v == nil {
						b.AppendNull()
						return nil
					}
					b.(*array.BinaryBuilder).Append(*v)
					return nil
				},
			}
		case "DATE":
			typ = arrow.FixedWidthTypes.Date32
			app = appendFunc{
				newScanDest: func() any { return &sql.NullTime{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullTime)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.Date32Builder).Append(arrow.Date32FromTime(v.Time))
					return nil
				},
			}
		case "TIME":
			// written as a string; DuckDB has issues reading TIME columns from parquet files
			// generated with this writer (see the equivalent restriction in the snowflake driver)
			typ = arrow.BinaryTypes.String
			app = appendFunc{
				newScanDest: func() any { return &sql.NullTime{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullTime)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.StringBuilder).Append(v.Time.Format("15:04:05.9999999"))
					return nil
				},
			}
		case "DATETIME", "DATETIME2", "SMALLDATETIME", "DATETIMEOFFSET":
			// microsecond precision: covers dates far outside the nanosecond timestamp range
			// (old ERP dates) at the cost of truncating datetime2(7)'s sub-microsecond digit
			tsType := &arrow.TimestampType{Unit: arrow.Microsecond}
			if dbt == "DATETIMEOFFSET" {
				tsType.TimeZone = "UTC"
			}
			typ = tsType
			app = appendFunc{
				newScanDest: func() any { return &sql.NullTime{} },
				append: func(b array.Builder, src any) error {
					v := src.(*sql.NullTime)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					ts, err := arrow.TimestampFromTime(v.Time.UTC(), arrow.Microsecond)
					if err != nil {
						return err
					}
					b.(*array.TimestampBuilder).Append(ts)
					return nil
				},
			}
		case "UNIQUEIDENTIFIER":
			typ = arrow.BinaryTypes.String
			app = appendFunc{
				newScanDest: func() any { return &mssqldb.NullUniqueIdentifier{} },
				append: func(b array.Builder, src any) error {
					v := src.(*mssqldb.NullUniqueIdentifier)
					if !v.Valid {
						b.AppendNull()
						return nil
					}
					b.(*array.StringBuilder).Append(v.UUID.String())
					return nil
				},
			}
		default:
			return nil, nil, fmt.Errorf("column %q has type %s which is not supported for ingestion, "+
				"consider excluding it or casting it to another data type", ct.Name(), dbt)
		}

		fields[i] = arrow.Field{Name: ct.Name(), Type: typ, Nullable: nullable}
		appenders[i] = app
	}
	return arrow.NewSchema(fields, nil), appenders, nil
}

// decimalToNum converts the decimal string emitted by the driver (e.g. "-1234.5600") into a
// decimal128 with the column's precision and scale. The conversion is exact integer arithmetic:
// it errors if the value has more fractional digits than the scale or does not fit the precision,
// rather than rounding silently.
func decimalToNum(s string, prec, scale int32) (decimal128.Num, error) {
	neg := false
	switch {
	case strings.HasPrefix(s, "-"):
		neg = true
		s = s[1:]
	case strings.HasPrefix(s, "+"):
		s = s[1:]
	}

	intPart, fracPart, _ := strings.Cut(s, ".")
	if intPart == "" {
		intPart = "0"
	}
	if int32(len(fracPart)) > scale {
		return decimal128.Num{}, fmt.Errorf("value %q has more fractional digits than the column scale %d", s, scale)
	}
	fracPart += strings.Repeat("0", int(scale)-len(fracPart))

	unscaled, ok := new(big.Int).SetString(intPart+fracPart, 10)
	if !ok {
		return decimal128.Num{}, fmt.Errorf("invalid decimal value %q", s)
	}
	if neg {
		unscaled.Neg(unscaled)
	}

	num := decimal128.FromBigInt(unscaled)
	if !num.FitsInPrecision(prec) {
		return decimal128.Num{}, fmt.Errorf("value %q does not fit in decimal(%d,%d)", s, prec, scale)
	}
	return num, nil
}
