package mssql_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/apache/arrow-go/v18/parquet/pqarrow"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/stretchr/testify/require"
)

func TestQueryAsFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	conn, olap := acquireTestMSSQL(t)
	w, ok := conn.AsWarehouse()
	require.True(t, ok)

	t.Run("Test Type Round Trip", func(t *testing.T) {
		testTypeRoundTrip(t, olap, w)
	})
	t.Run("Test Batches", func(t *testing.T) {
		testBatches(t, w)
	})
	t.Run("Test Empty Result", func(t *testing.T) {
		testEmptyResult(t, w)
	})
	t.Run("Test Cancellation", func(t *testing.T) {
		testCancellation(t, w)
	})
	t.Run("Test Unsupported Type", func(t *testing.T) {
		testUnsupportedType(t, w)
	})
}

func testTypeRoundTrip(t *testing.T, olap drivers.OLAPStore, w drivers.Warehouse) {
	err := olap.Exec(t.Context(), &drivers.Statement{Query: `
		IF OBJECT_ID('warehouse_types', 'U') IS NOT NULL DROP TABLE warehouse_types;
		CREATE TABLE warehouse_types (
			id INT NOT NULL,
			dec38_col DECIMAL(38,10),
			dec18_col DECIMAL(18,2),
			money_col MONEY,
			bit_col BIT,
			int_col INT,
			bigint_col BIGINT,
			nvarchar_col NVARCHAR(100),
			date_col DATE,
			datetime_col DATETIME,
			datetime2_col DATETIME2
		);
		INSERT INTO warehouse_types VALUES
		(1, 1234567890123456789012345678.0123456789, 9999999999999999.99, 922337203685477.5807, 1, 2147483647, 9223372036854775807, N'áñ 中文 🚀', '0001-01-01', '1753-01-01 00:00:00', '0001-01-01 00:00:00'),
		(2, -0.0000000001, -12345678901234.56, -5.6, 0, -2147483648, -9223372036854775808, N'', '9999-12-31', '2025-02-14 12:34:56.997', '2025-02-14 12:34:56.1234567'),
		(3, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
	`})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = olap.Exec(context.Background(), &drivers.Statement{Query: "DROP TABLE IF EXISTS warehouse_types"})
	})

	iter, err := w.QueryAsFiles(t.Context(), map[string]any{"sql": "SELECT * FROM warehouse_types ORDER BY id"})
	require.NoError(t, err)
	defer iter.Close()

	files, err := iter.Next(t.Context())
	require.NoError(t, err)
	require.Len(t, files, 1)

	tbl := readParquet(t, files[0])
	defer tbl.Release()
	require.Equal(t, int64(3), tbl.NumRows())

	// schema fidelity: decimals keep their precision and scale, timestamps are microsecond, date is date32
	sc := tbl.Schema()
	requireField := func(name string, typ arrow.DataType) {
		f, ok := sc.FieldsByName(name)
		require.True(t, ok, "field %s missing", name)
		require.Truef(t, arrow.TypeEqual(f[0].Type, typ), "field %s has type %s, expected %s", name, f[0].Type, typ)
	}
	requireField("dec38_col", &arrow.Decimal128Type{Precision: 38, Scale: 10})
	requireField("dec18_col", &arrow.Decimal128Type{Precision: 18, Scale: 2})
	requireField("money_col", &arrow.Decimal128Type{Precision: 19, Scale: 4})
	requireField("bit_col", arrow.FixedWidthTypes.Boolean)
	requireField("int_col", arrow.PrimitiveTypes.Int32)
	requireField("bigint_col", arrow.PrimitiveTypes.Int64)
	requireField("nvarchar_col", arrow.BinaryTypes.String)
	requireField("date_col", arrow.FixedWidthTypes.Date32)
	requireField("datetime_col", &arrow.TimestampType{Unit: arrow.Microsecond})
	requireField("datetime2_col", &arrow.TimestampType{Unit: arrow.Microsecond})

	// value fidelity, ordered by id; decimals are compared through their exact unscaled
	// representation (ValueStr renders small values in scientific notation)
	require.Equal(t, "1234567890123456789012345678.0123456789", columnDecimalStr(t, tbl, "dec38_col", 0))
	require.Equal(t, "-0.0000000001", columnDecimalStr(t, tbl, "dec38_col", 1))
	require.Equal(t, "9999999999999999.99", columnDecimalStr(t, tbl, "dec18_col", 0))
	require.Equal(t, "-12345678901234.56", columnDecimalStr(t, tbl, "dec18_col", 1))
	require.Equal(t, "922337203685477.5807", columnDecimalStr(t, tbl, "money_col", 0))
	require.Equal(t, "-5.6000", columnDecimalStr(t, tbl, "money_col", 1))
	require.Equal(t, "true", columnStr(t, tbl, "bit_col", 0))
	require.Equal(t, "false", columnStr(t, tbl, "bit_col", 1))
	require.Equal(t, "2147483647", columnStr(t, tbl, "int_col", 0))
	require.Equal(t, "-2147483648", columnStr(t, tbl, "int_col", 1))
	require.Equal(t, "9223372036854775807", columnStr(t, tbl, "bigint_col", 0))
	require.Equal(t, "-9223372036854775808", columnStr(t, tbl, "bigint_col", 1))
	require.Equal(t, "áñ 中文 🚀", columnStr(t, tbl, "nvarchar_col", 0))
	require.Equal(t, "", columnStr(t, tbl, "nvarchar_col", 1))

	// dates and timestamps, including dates far before the nanosecond-timestamp range
	dateCol := columnOf(t, tbl, "date_col").(*array.Date32)
	require.Equal(t, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), dateCol.Value(0).ToTime())
	require.Equal(t, time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC), dateCol.Value(1).ToTime())
	dtCol := columnOf(t, tbl, "datetime_col").(*array.Timestamp)
	require.Equal(t, time.Date(1753, 1, 1, 0, 0, 0, 0, time.UTC), dtCol.Value(0).ToTime(arrow.Microsecond))
	dt2Col := columnOf(t, tbl, "datetime2_col").(*array.Timestamp)
	require.Equal(t, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), dt2Col.Value(0).ToTime(arrow.Microsecond))
	// datetime2(7) is truncated from 100ns to microsecond precision
	require.Equal(t, time.Date(2025, 2, 14, 12, 34, 56, 123456000, time.UTC), dt2Col.Value(1).ToTime(arrow.Microsecond))

	// the all-NULLs row survives as NULLs
	for _, name := range []string{"dec38_col", "dec18_col", "money_col", "bit_col", "int_col", "bigint_col", "nvarchar_col", "date_col", "datetime_col", "datetime2_col"} {
		require.Truef(t, columnOf(t, tbl, name).IsNull(2), "column %s row 2 should be NULL", name)
	}
}

func testBatches(t *testing.T, w drivers.Warehouse) {
	// more rows than the flush batch size, and not a multiple of it, so both the intermediate
	// and the final partial flush are exercised
	const n = 20001
	iter, err := w.QueryAsFiles(t.Context(), map[string]any{"sql": "SELECT CAST(value AS BIGINT) AS v FROM GENERATE_SERIES(1, 20001)"})
	require.NoError(t, err)
	defer iter.Close()

	files, err := iter.Next(t.Context())
	require.NoError(t, err)
	require.Len(t, files, 1)

	tbl := readParquet(t, files[0])
	defer tbl.Release()
	require.Equal(t, int64(n), tbl.NumRows())

	// sum detects lost or duplicated rows, not just a wrong count
	var sum int64
	col := tbl.Column(0)
	for _, chunk := range col.Data().Chunks() {
		vs := chunk.(*array.Int64)
		for i := 0; i < vs.Len(); i++ {
			sum += vs.Value(i)
		}
	}
	require.Equal(t, int64(n)*(n+1)/2, sum)

	_, err = iter.Next(t.Context())
	require.ErrorIs(t, err, io.EOF)
}

func testEmptyResult(t *testing.T, w drivers.Warehouse) {
	iter, err := w.QueryAsFiles(t.Context(), map[string]any{"sql": "SELECT 1 AS x WHERE 1 = 0"})
	require.NoError(t, err)
	defer iter.Close()

	_, err = iter.Next(t.Context())
	require.ErrorIs(t, err, drivers.ErrNoRows)
}

func testCancellation(t *testing.T, w drivers.Warehouse) {
	// cancelled before iterating
	ctx, cancel := context.WithCancel(t.Context())
	iter, err := w.QueryAsFiles(ctx, map[string]any{"sql": "SELECT value FROM GENERATE_SERIES(1, 100000)"})
	require.NoError(t, err)
	defer iter.Close()
	cancel()
	_, err = iter.Next(ctx)
	require.ErrorIs(t, err, context.Canceled)

	// cancelled while streaming a large result
	ctx, cancel = context.WithCancel(t.Context())
	defer cancel()
	iter2, err := w.QueryAsFiles(ctx, map[string]any{"sql": "SELECT a.value AS v FROM GENERATE_SERIES(1, 10000) a CROSS JOIN GENERATE_SERIES(1, 1000) b"})
	require.NoError(t, err)
	defer iter2.Close()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	_, err = iter2.Next(ctx)
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled))
}

func testUnsupportedType(t *testing.T, w drivers.Warehouse) {
	_, err := w.QueryAsFiles(t.Context(), map[string]any{"sql": "SELECT CAST('x' AS sql_variant) AS v"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not supported")
}

func readParquet(t *testing.T, path string) arrow.Table {
	fw, err := os.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = fw.Close() })

	tbl, err := pqarrow.ReadTable(t.Context(), fw, nil, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	require.NoError(t, err)
	return tbl
}

func columnOf(t *testing.T, tbl arrow.Table, name string) arrow.Array {
	idx := tbl.Schema().FieldIndices(name)
	require.Len(t, idx, 1, "column %s", name)
	col := tbl.Column(idx[0])
	require.Len(t, col.Data().Chunks(), 1, "expected a single chunk for column %s", name)
	return col.Data().Chunk(0)
}

func columnStr(t *testing.T, tbl arrow.Table, name string, row int) string {
	arr := columnOf(t, tbl, name)
	require.Falsef(t, arr.IsNull(row), "column %s row %d is NULL", name, row)
	return arr.ValueStr(row)
}

func columnDecimalStr(t *testing.T, tbl arrow.Table, name string, row int) string {
	arr := columnOf(t, tbl, name)
	require.Falsef(t, arr.IsNull(row), "column %s row %d is NULL", name, row)
	dec, ok := arr.(*array.Decimal128)
	require.Truef(t, ok, "column %s is not decimal128", name)
	return dec.Value(row).ToString(dec.DataType().(*arrow.Decimal128Type).Scale)
}
