package mssql_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestOLAP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	_, olap := acquireTestMSSQL(t)
	t.Run("Test MapScan", func(t *testing.T) {
		testMapScan(t, olap)
	})
	t.Run("Test Scan", func(t *testing.T) {
		testScan(t, olap)
	})
	t.Run("Test Empty Rows", func(t *testing.T) {
		testEmptyRows(t, olap)
	})
	t.Run("Test Exec", func(t *testing.T) {
		testExec(t, olap)
	})
	t.Run("Test Scan Full Table", func(t *testing.T) {
		testFullTableScan(t, olap)
	})
	t.Run("Test LoadDDL", func(t *testing.T) {
		testLoadDDL(t, olap)
	})
	t.Run("Test Information Schema", func(t *testing.T) {
		testInformationSchema(t, olap)
	})
}

func testMapScan(t *testing.T, olap drivers.OLAPStore) {
	tests := []struct {
		query  string
		args   []any
		result map[string]any
	}{
		// Integer types: TINYINT, SMALLINT, INT, BIGINT
		{
			"SELECT tinyint_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"tinyint_col": uint8(255)},
		},
		{
			"SELECT smallint_col FROM all_datatypes WHERE smallint_col = 32767",
			nil,
			map[string]any{"smallint_col": int16(32767)},
		},
		{
			"SELECT int_col FROM all_datatypes WHERE int_col = @p1",
			[]any{2147483647},
			map[string]any{"int_col": int32(2147483647)},
		},
		{
			"SELECT bigint_col FROM all_datatypes WHERE bigint_col = 999999999999999999",
			nil,
			map[string]any{"bigint_col": int64(999999999999999999)},
		},
		// Boolean type
		{
			"SELECT bit_col FROM all_datatypes WHERE bit_col = 1",
			nil,
			map[string]any{"bit_col": true},
		},
		// Floating point types: FLOAT, REAL
		{
			"SELECT float_col FROM all_datatypes WHERE float_col = 2.2",
			nil,
			map[string]any{"float_col": float64(2.2)},
		},
		{
			"SELECT real_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"real_col": float64(1.5)},
		},
		// Decimal types: DECIMAL, NUMERIC, MONEY, SMALLMONEY
		{
			"SELECT decimal_col FROM all_datatypes WHERE decimal_col = 3.3",
			nil,
			map[string]any{"decimal_col": "3.30"},
		},
		{
			"SELECT numeric_col FROM all_datatypes WHERE numeric_col = 4.4",
			nil,
			map[string]any{"numeric_col": "4.4000"},
		},
		{
			"SELECT money_col FROM all_datatypes WHERE money_col = 5.6",
			nil,
			map[string]any{"money_col": "5.6000"},
		},
		{
			"SELECT smallmoney_col FROM all_datatypes WHERE smallmoney_col = 7.8",
			nil,
			map[string]any{"smallmoney_col": "7.8000"},
		},
		// String types: CHAR, VARCHAR, VARCHAR(MAX), NCHAR, NVARCHAR, NVARCHAR(MAX), TEXT, NTEXT
		{
			"SELECT char_col FROM all_datatypes WHERE char_col = 'C'",
			nil,
			map[string]any{"char_col": "C         "}, // CHAR(10) is padded with spaces
		},
		{
			"SELECT varchar_col FROM all_datatypes WHERE varchar_col = 'VarChar'",
			nil,
			map[string]any{"varchar_col": "VarChar"},
		},
		{
			"SELECT varcharmax_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"varcharmax_col": "Long varchar content"},
		},
		{
			"SELECT nchar_col FROM all_datatypes WHERE nchar_col = N'N'",
			nil,
			map[string]any{"nchar_col": "N         "}, // NCHAR(10) is padded with spaces
		},
		{
			"SELECT nvarchar_col FROM all_datatypes WHERE nvarchar_col = N'NVarChar'",
			nil,
			map[string]any{"nvarchar_col": "NVarChar"},
		},
		{
			"SELECT nvarcharmax_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"nvarcharmax_col": "Long nvarchar content"},
		},
		{
			"SELECT text_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"text_col": "Text"},
		},
		{
			"SELECT ntext_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"ntext_col": "NText"},
		},
		// Binary types: BINARY, VARBINARY
		{
			"SELECT binary_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"binary_col": "Binary\x00\x00\x00\x00"},
		},
		{
			"SELECT varbinary_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"varbinary_col": "VarBinary"},
		},
		// Date and Time types: DATE, TIME, DATETIME, DATETIME2, SMALLDATETIME
		{
			"SELECT date_col FROM all_datatypes WHERE date_col = '2024-02-14'",
			nil,
			map[string]any{"date_col": time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)},
		},
		{
			"SELECT time_col FROM all_datatypes WHERE time_col = '12:34:56'",
			nil,
			map[string]any{"time_col": "12:34:56"},
		},
		{
			"SELECT datetime_col FROM all_datatypes WHERE datetime_col = '2025-02-14 12:34:56'",
			nil,
			map[string]any{"datetime_col": time.Date(2025, 2, 14, 12, 34, 56, 0, time.UTC)},
		},
		{
			"SELECT datetime2_col FROM all_datatypes WHERE datetime2_col = '2025-02-14 12:34:56.1234567'",
			nil,
			map[string]any{"datetime2_col": time.Date(2025, 2, 14, 12, 34, 56, 123456700, time.UTC)},
		},
		{
			"SELECT smalldatetime_col FROM all_datatypes WHERE smalldatetime_col = '2025-02-14 12:34:00'",
			nil,
			map[string]any{"smalldatetime_col": time.Date(2025, 2, 14, 12, 34, 0, 0, time.UTC)},
		},
		// UNIQUEIDENTIFIER type
		{
			"SELECT uniqueidentifier_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"uniqueidentifier_col": "AABBCCDD-EEFF-1122-3344-556677889900"},
		},
		// XML type
		{
			"SELECT xml_col FROM all_datatypes WHERE tinyint_col = 255",
			nil,
			map[string]any{"xml_col": "<root><key>value</key></root>"},
		},
		{
			"SELECT NULL AS null_col",
			nil,
			map[string]any{"null_col": nil},
		},
	}
	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			rows, err := olap.Query(t.Context(), &drivers.Statement{Query: test.query, Args: test.args})
			require.NoError(t, err)
			defer rows.Close()
			require.True(t, rows.Next())
			res := make(map[string]any)
			err = rows.MapScan(res)
			require.NoError(t, err)
			require.Equal(t, test.result, res)
			require.False(t, rows.Next())
			require.NoError(t, rows.Err())
			require.NoError(t, rows.Close())
		})
	}
}

func testScan(t *testing.T, olap drivers.OLAPStore) {
	tests := []struct {
		name     string
		query    string
		args     []any
		scanFunc func(rows drivers.Rows) error
	}{
		{
			name:  "Null scan",
			query: "SELECT NULL AS null_col",
			args:  nil,
			scanFunc: func(rows drivers.Rows) error {
				var nullCol any

				require.True(t, rows.Next())
				err := rows.Scan(&nullCol)
				require.NoError(t, err)
				require.Nil(t, nullCol)
				return nil
			},
		},
		{
			name:  "Integer types scan",
			query: "SELECT tinyint_col, smallint_col, int_col, bigint_col FROM all_datatypes WHERE int_col = @p1",
			args:  []any{2147483647},
			scanFunc: func(rows drivers.Rows) error {
				var tinyintCol uint8
				var smallintCol int16
				var intCol int32
				var bigintCol int64

				require.True(t, rows.Next())
				err := rows.Scan(&tinyintCol, &smallintCol, &intCol, &bigintCol)
				require.NoError(t, err)
				require.Equal(t, uint8(255), tinyintCol)
				require.Equal(t, int16(32767), smallintCol)
				require.Equal(t, int32(2147483647), intCol)
				require.Equal(t, int64(999999999999999999), bigintCol)
				return nil
			},
		},
		{
			name:  "Floating point and decimal types scan",
			query: "SELECT float_col, real_col, decimal_col, numeric_col, money_col, smallmoney_col FROM all_datatypes WHERE int_col = @p1",
			args:  []any{2147483647},
			scanFunc: func(rows drivers.Rows) error {
				var floatCol float64
				var realCol float64
				var decimalCol string
				var numericCol string
				var moneyCol string
				var smallmoneyCol string

				require.True(t, rows.Next())
				err := rows.Scan(&floatCol, &realCol, &decimalCol, &numericCol, &moneyCol, &smallmoneyCol)
				require.NoError(t, err)
				require.Equal(t, float64(2.2), floatCol)
				require.Equal(t, float64(1.5), realCol)
				require.Equal(t, "3.30", decimalCol)
				require.Equal(t, "4.4000", numericCol)
				require.Equal(t, "5.6000", moneyCol)
				require.Equal(t, "7.8000", smallmoneyCol)
				return nil
			},
		},
		{
			name:  "String types scan",
			query: "SELECT char_col, varchar_col, varcharmax_col, nchar_col, nvarchar_col, nvarcharmax_col FROM all_datatypes WHERE int_col = @p1",
			args:  []any{2147483647},
			scanFunc: func(rows drivers.Rows) error {
				var charCol string
				var varcharCol string
				var varcharmaxCol string
				var ncharCol string
				var nvarcharCol string
				var nvarcharmaxCol string

				require.True(t, rows.Next())
				err := rows.Scan(&charCol, &varcharCol, &varcharmaxCol, &ncharCol, &nvarcharCol, &nvarcharmaxCol)
				require.NoError(t, err)
				require.Equal(t, "C", strings.TrimRight(charCol, " ")) // char is padded with spaces
				require.Equal(t, "VarChar", varcharCol)
				require.Equal(t, "Long varchar content", varcharmaxCol)
				require.Equal(t, "N", strings.TrimRight(ncharCol, " ")) // nchar is padded with spaces
				require.Equal(t, "NVarChar", nvarcharCol)
				require.Equal(t, "Long nvarchar content", nvarcharmaxCol)
				return nil
			},
		},
		{
			name:  "Binary types scan",
			query: "SELECT binary_col, varbinary_col FROM all_datatypes WHERE int_col = @p1",
			args:  []any{2147483647},
			scanFunc: func(rows drivers.Rows) error {
				var binaryCol []byte
				var varbinaryCol []byte

				require.True(t, rows.Next())
				err := rows.Scan(&binaryCol, &varbinaryCol)
				require.NoError(t, err)
				require.Equal(t, "Binary", strings.TrimRight(string(binaryCol), "\x00")) // binary is padded with null bytes
				require.Equal(t, "VarBinary", string(varbinaryCol))
				return nil
			},
		},
		{
			name:  "Date and Time types scan",
			query: "SELECT date_col, datetime_col, datetime2_col, smalldatetime_col, datetimeoffset_col, time_col FROM all_datatypes WHERE int_col = 2147483647",
			args:  nil,
			scanFunc: func(rows drivers.Rows) error {
				var dateCol time.Time
				var datetimeCol time.Time
				var datetime2Col time.Time
				var smalldatetimeCol time.Time
				var datetimeoffsetCol time.Time
				var timeCol time.Time

				require.True(t, rows.Next())
				err := rows.Scan(&dateCol, &datetimeCol, &datetime2Col, &smalldatetimeCol, &datetimeoffsetCol, &timeCol)
				require.NoError(t, err)
				require.True(t, dateCol.Equal(time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)))
				require.True(t, datetimeCol.Equal(time.Date(2025, 2, 14, 12, 34, 56, 0, time.UTC)))
				require.True(t, datetime2Col.Equal(time.Date(2025, 2, 14, 12, 34, 56, 123456700, time.UTC)))
				require.True(t, smalldatetimeCol.Equal(time.Date(2025, 2, 14, 12, 34, 0, 0, time.UTC)))
				require.True(t, datetimeoffsetCol.Equal(time.Date(2025, 2, 14, 10, 34, 56, 0, time.UTC))) // 12:34:56 +02:00
				require.Equal(t, "12:34:56", timeCol.Format("15:04:05"))
				return nil
			},
		},
		{
			name:  "Boolean type scan",
			query: "SELECT bit_col FROM all_datatypes WHERE int_col = 2147483647",
			args:  nil,
			scanFunc: func(rows drivers.Rows) error {
				var bitCol bool

				require.True(t, rows.Next())
				err := rows.Scan(&bitCol)
				require.NoError(t, err)
				require.True(t, bitCol)
				return nil
			},
		},
		{
			name:  "UNIQUEIDENTIFIER and XML types scan",
			query: "SELECT uniqueidentifier_col, xml_col FROM all_datatypes WHERE int_col = 2147483647",
			args:  nil,
			scanFunc: func(rows drivers.Rows) error {
				var uuidCol []byte
				var xmlCol string

				require.True(t, rows.Next())
				err := rows.Scan(&uuidCol, &xmlCol)
				require.NoError(t, err)
				require.Len(t, uuidCol, 16)
				require.Equal(t, "<root><key>value</key></root>", xmlCol)
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rows, err := olap.Query(t.Context(), &drivers.Statement{Query: test.query, Args: test.args})
			require.NoError(t, err)
			defer rows.Close()

			err = test.scanFunc(rows)
			require.NoError(t, err)
			require.NoError(t, rows.Err())
			require.NoError(t, rows.Close())
		})
	}
}

func testEmptyRows(t *testing.T, olap drivers.OLAPStore) {
	rows, err := olap.Query(t.Context(), &drivers.Statement{Query: "SELECT TOP 0 int_col, float_col FROM all_datatypes"})
	require.NoError(t, err)
	defer rows.Close()

	sc := rows.Schema
	require.Len(t, sc.Fields, 2)
	require.Equal(t, "int_col", sc.Fields[0].Name)
	require.Equal(t, "float_col", sc.Fields[1].Name)
	require.False(t, rows.Next())
	require.Nil(t, rows.Err())
}

func testExec(t *testing.T, olap drivers.OLAPStore) {
	// dry run with SELECT query
	_, err := olap.Query(t.Context(), &drivers.Statement{Query: "SELECT 1", DryRun: true})
	require.NoError(t, err)

	// create table
	err = olap.Exec(t.Context(), &drivers.Statement{Query: "IF OBJECT_ID('exec_test', 'U') IS NULL CREATE TABLE exec_test (id INT, name VARCHAR(255))"})
	require.NoError(t, err)

	// drop table
	err = olap.Exec(t.Context(), &drivers.Statement{Query: "DROP TABLE IF EXISTS exec_test"})
	require.NoError(t, err)
}

func testFullTableScan(t *testing.T, olap drivers.OLAPStore) {
	rows, err := olap.Query(t.Context(), &drivers.Statement{Query: `SELECT
		tinyint_col, smallint_col, int_col, bigint_col, bit_col,
		float_col, real_col, decimal_col, numeric_col, money_col, smallmoney_col,
		char_col, varchar_col, varcharmax_col, nchar_col, nvarchar_col, nvarcharmax_col, text_col, ntext_col,
		binary_col, varbinary_col,
		date_col, time_col, datetime_col, datetime2_col, smalldatetime_col, datetimeoffset_col,
		uniqueidentifier_col, xml_col
		FROM all_datatypes`})
	require.NoError(t, err)
	defer rows.Close()

	count := 0
	for rows.Next() {
		res := make(map[string]any)
		err := rows.MapScan(res)
		require.NoError(t, err)

		// Verify that all expected columns are present
		expectedCols := []string{
			"tinyint_col", "smallint_col", "int_col", "bigint_col", "bit_col",
			"float_col", "real_col", "decimal_col", "numeric_col", "money_col", "smallmoney_col",
			"char_col", "varchar_col", "varcharmax_col", "nchar_col", "nvarchar_col", "nvarcharmax_col", "text_col", "ntext_col",
			"binary_col", "varbinary_col",
			"date_col", "time_col", "datetime_col", "datetime2_col", "smalldatetime_col", "datetimeoffset_col",
			"uniqueidentifier_col", "xml_col",
		}

		for _, col := range expectedCols {
			_, exists := res[col]
			require.True(t, exists, "Column %s should be present in result", col)
		}

		count++
	}
	require.NoError(t, rows.Err())
	require.Equal(t, count, 3)
}

func testLoadDDL(t *testing.T, olap drivers.OLAPStore) {
	// Test DDL for a table
	table, err := olap.InformationSchema().Lookup(t.Context(), "", "dbo", "all_datatypes")
	require.NoError(t, err)
	err = olap.InformationSchema().LoadDDL(t.Context(), table)
	require.NoError(t, err)
	require.Contains(t, table.DDL, "CREATE TABLE")
	require.Contains(t, table.DDL, "all_datatypes")

	// Create a view and test DDL for it
	err = olap.Exec(t.Context(), &drivers.Statement{Query: "CREATE OR ALTER VIEW test_ddl_view AS SELECT int_col, varchar_col FROM all_datatypes"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = olap.Exec(t.Context(), &drivers.Statement{Query: "DROP VIEW IF EXISTS test_ddl_view"})
	})

	view, err := olap.InformationSchema().Lookup(t.Context(), "", "dbo", "test_ddl_view")
	require.NoError(t, err)
	err = olap.InformationSchema().LoadDDL(t.Context(), view)
	require.NoError(t, err)
	require.Contains(t, strings.ToLower(view.DDL), "create")
	require.Contains(t, strings.ToLower(view.DDL), "view")
	require.Contains(t, strings.ToLower(view.DDL), "test_ddl_view")
	require.Contains(t, strings.ToLower(view.DDL), "select")
}

func testInformationSchema(t *testing.T, olap drivers.OLAPStore) {
	is := olap.InformationSchema()

	// create extra tables so pagination is actually exercised
	for _, name := range []string{"pagination_test_a", "pagination_test_b"} {
		err := olap.Exec(t.Context(), &drivers.Statement{Query: fmt.Sprintf("IF OBJECT_ID('%[1]s', 'U') IS NULL CREATE TABLE %[1]s (id INT)", name)})
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		for _, name := range []string{"pagination_test_a", "pagination_test_b"} {
			_ = olap.Exec(context.Background(), &drivers.Statement{Query: fmt.Sprintf("DROP TABLE IF EXISTS %s", name)})
		}
	})

	// list schemas: dbo must be present, system schemas must not
	schemas, next, err := is.ListDatabaseSchemas(t.Context(), 0, "")
	require.NoError(t, err)
	require.Empty(t, next)
	var names []string
	for _, s := range schemas {
		names = append(names, s.DatabaseSchema)
	}
	require.Contains(t, names, "dbo")
	require.NotContains(t, names, "sys")
	require.NotContains(t, names, "INFORMATION_SCHEMA")

	// list tables in dbo
	tables, next, err := is.ListTables(t.Context(), "", "dbo", 0, "")
	require.NoError(t, err)
	require.Empty(t, next)
	var tableNames []string
	for _, tbl := range tables {
		tableNames = append(tableNames, tbl.Name)
	}
	require.Contains(t, tableNames, "all_datatypes")

	// paginate tables one at a time and verify the pages add up to the full listing
	var paged []string
	pageToken := ""
	for {
		page, nextToken, err := is.ListTables(t.Context(), "", "dbo", 1, pageToken)
		require.NoError(t, err)
		require.LessOrEqual(t, len(page), 1)
		for _, tbl := range page {
			paged = append(paged, tbl.Name)
		}
		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}
	require.Equal(t, tableNames, paged)

	// lookup returns the table's schema
	table, err := is.Lookup(t.Context(), "", "dbo", "all_datatypes")
	require.NoError(t, err)
	require.False(t, table.View)
	require.NotEmpty(t, table.Schema.Fields)
	fieldTypes := make(map[string]string)
	for _, f := range table.Schema.Fields {
		fieldTypes[f.Name] = f.Type.Code.String()
	}
	require.Equal(t, "CODE_INT32", fieldTypes["int_col"])
	require.Equal(t, "CODE_DECIMAL", fieldTypes["decimal_col"])
	require.Equal(t, "CODE_TIMESTAMP", fieldTypes["datetime2_col"])
	require.Equal(t, "CODE_DATE", fieldTypes["date_col"])
	require.Equal(t, "CODE_BOOL", fieldTypes["bit_col"])
	require.Equal(t, "CODE_STRING", fieldTypes["nvarchar_col"])
	require.Equal(t, "CODE_UUID", fieldTypes["uniqueidentifier_col"])
}

func acquireTestMSSQL(t *testing.T) (drivers.Handle, drivers.OLAPStore) {
	cfg := testruntime.AcquireConnector(t, "mssql")
	conn, err := drivers.Open("mssql", "", "default", cfg, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	olap, ok := conn.AsOLAP("default")
	require.True(t, ok)

	return conn, olap
}
