package odoo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/apache/arrow-go/v18/parquet"
	"github.com/apache/arrow-go/v18/parquet/compress"
	"github.com/apache/arrow-go/v18/parquet/pqarrow"
	"github.com/c2h5oh/datasize"
	"github.com/mitchellh/mapstructure"
	"github.com/rilldata/rill/runtime/drivers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("github.com/rilldata/rill/runtime/drivers/odoo")

const defaultPageSize = 2000

// Same threshold as the mssql driver: flush the buffered row group at 64MB.
const rowGroupBufferSize = int64(datasize.MB) * 64

var _ drivers.Warehouse = &connection{}

// sourceProperties are the properties of a model that reads from Odoo.
type sourceProperties struct {
	Model    string         `mapstructure:"model"`
	Fields   []string       `mapstructure:"fields"`
	Domain   []any          `mapstructure:"domain"`
	Context  map[string]any `mapstructure:"context"`
	PageSize int            `mapstructure:"page_size"`
}

func parseSourceProperties(props map[string]any) (*sourceProperties, error) {
	p := &sourceProperties{}
	if err := mapstructure.WeakDecode(props, p); err != nil {
		return nil, fmt.Errorf("odoo: invalid model properties: %w", err)
	}
	if p.Model == "" {
		return nil, errors.New(`odoo: property "model" is mandatory`)
	}
	if len(p.Fields) == 0 {
		return nil, errors.New(`odoo: property "fields" is mandatory`)
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultPageSize
	}
	if p.Domain == nil {
		p.Domain = []any{}
	}
	// "id" goes first and only once: pagination runs on it.
	fields := []string{"id"}
	for _, f := range p.Fields {
		if !slices.Contains(fields, f) {
			fields = append(fields, f)
		}
	}
	p.Fields = fields
	return p, nil
}

// QueryAsFiles implements drivers.Warehouse. It resolves the column types with fields_get and
// returns an iterator that pages through search_read into a single parquet file.
func (c *connection) QueryAsFiles(ctx context.Context, props map[string]any) (outIt drivers.FileIterator, outErr error) {
	ctx, span := tracer.Start(ctx, "Connection.QueryAsFiles")
	defer func() {
		if outErr != nil {
			span.SetStatus(codes.Error, outErr.Error())
		}
		span.End()
	}()

	src, err := parseSourceProperties(props)
	if err != nil {
		return nil, err
	}
	cl, err := newClient(c.conf)
	if err != nil {
		return nil, err
	}
	odooCtx := map[string]any{}
	if c.conf.Lang != "" {
		odooCtx["lang"] = c.conf.Lang
	}
	maps.Copy(odooCtx, src.Context)

	var info map[string]fieldInfo
	err = cl.call(ctx, src.Model, "fields_get", map[string]any{
		"allfields": src.Fields, "attributes": []string{"type"}, "context": odooCtx,
	}, &info)
	if err != nil {
		return nil, err
	}
	cols, err := columnsFor(src.Fields, info)
	if err != nil {
		return nil, err
	}
	tempDir, err := c.storage.RandomTempDir("odoo")
	if err != nil {
		return nil, err
	}
	return &fileIterator{
		client:  cl,
		src:     src,
		odooCtx: odooCtx,
		cols:    cols,
		schema:  schemaFor(cols),
		tempDir: tempDir,
		logger:  c.logger,
	}, nil
}

type fileIterator struct {
	client  *client
	src     *sourceProperties
	odooCtx map[string]any
	cols    []column
	schema  *arrow.Schema
	tempDir string
	logger  *zap.Logger
	done    bool
}

var _ drivers.FileIterator = &fileIterator{}

// Close implements drivers.FileIterator.
func (f *fileIterator) Close() error {
	return os.RemoveAll(f.tempDir)
}

// Format implements drivers.FileIterator. The file extension already says parquet.
func (f *fileIterator) Format() string {
	return ""
}

// SetKeepFilesUntilClose implements drivers.FileIterator. The only file lives until Close anyway.
func (f *fileIterator) SetKeepFilesUntilClose() {}

// Next pages through search_read by id and writes every page into one parquet file. Paging is
// by keyset — each page asks for id > the last id seen — so records created or deleted while
// reading neither shift nor repeat pages. An empty model yields a file with the schema and no
// rows: a table with no open orders is a valid state, not an error.
func (f *fileIterator) Next(ctx context.Context) ([]string, error) {
	if f.done {
		return nil, io.EOF
	}
	f.done = true

	fw, err := os.CreateTemp(f.tempDir, "odoo*.parquet")
	if err != nil {
		return nil, err
	}
	defer fw.Close()
	writer, err := pqarrow.NewFileWriter(f.schema, fw,
		parquet.NewWriterProperties(
			parquet.WithCompression(compress.Codecs.Snappy),
			parquet.WithRootRepetition(parquet.Repetitions.Required),
			// Same as mssql: DuckDB has issues reading the string statistics this writer produces.
			parquet.WithStats(false),
		),
		pqarrow.NewArrowWriterProperties(pqarrow.WithStoreSchema()))
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	builder := array.NewRecordBuilder(memory.DefaultAllocator, f.schema)
	defer builder.Release()

	var lastID, total int64
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		domain := append(slices.Clone(f.src.Domain), []any{"id", ">", lastID})
		var rows []map[string]any
		err := f.client.call(ctx, f.src.Model, "search_read", map[string]any{
			"domain":  domain,
			"fields":  f.src.Fields,
			"order":   "id asc",
			"limit":   f.src.PageSize,
			"context": f.odooCtx,
		}, &rows)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			for i, col := range f.cols {
				if err := col.appendValue(builder.Field(i), row[col.field]); err != nil {
					return nil, fmt.Errorf("odoo: %s, field %q: %w", f.src.Model, col.name, err)
				}
			}
			if lastID, err = asInt64(row["id"]); err != nil {
				return nil, fmt.Errorf("odoo: %s: id: %w", f.src.Model, err)
			}
		}
		total += int64(len(rows))
		if len(rows) > 0 {
			rec := builder.NewRecordBatch()
			if writer.RowGroupTotalBytesWritten() >= rowGroupBufferSize {
				writer.NewBufferedRowGroup()
			}
			err := writer.WriteBuffered(rec)
			rec.Release()
			if err != nil {
				return nil, err
			}
		}
		if len(rows) < f.src.PageSize {
			break
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	f.logger.Debug("odoo: model read", zap.String("model", f.src.Model), zap.Int64("rows", total))
	return []string{fw.Name()}, nil
}
