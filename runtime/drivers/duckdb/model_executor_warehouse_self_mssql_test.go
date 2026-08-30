package duckdb_test

import (
	"context"
	"testing"

	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"github.com/rilldata/rill/runtime/testruntime"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	_ "github.com/rilldata/rill/runtime/drivers/duckdb"
	_ "github.com/rilldata/rill/runtime/drivers/mssql"
)

func TestMSSQLToDuckDBTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cfg := testruntime.AcquireConnector(t, "mssql")

	duckDB, err := drivers.Open("duckdb", "", "default", map[string]any{"data_dir": t.TempDir()}, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)

	inputHandle, err := drivers.Open("mssql", "", "default", cfg, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)

	opts := &drivers.ModelExecutorOptions{
		InputHandle:     inputHandle,
		InputConnector:  "mssql",
		OutputHandle:    duckDB,
		OutputConnector: "duckdb",
		Env: &drivers.ModelEnv{
			AllowHostAccess: false,
			StageChanges:    true,
		},
		PreliminaryInputProperties: map[string]any{
			"sql": "SELECT int_col, bigint_col, bit_col, decimal_col, nvarchar_col, date_col, datetime2_col FROM all_datatypes",
		},
		PreliminaryOutputProperties: map[string]any{
			"table": "sink",
		},
	}

	me, err := duckDB.AsModelExecutor("default", opts)
	require.NoError(t, err)

	execOpts := &drivers.ModelExecuteOptions{
		ModelExecutorOptions: opts,
		InputProperties:      opts.PreliminaryInputProperties,
		OutputProperties:     opts.PreliminaryOutputProperties,
	}
	_, err = me.Execute(context.Background(), execOpts)
	require.NoError(t, err)

	olap, ok := duckDB.AsOLAP("default")
	require.True(t, ok)

	res, err := olap.Query(context.Background(), &drivers.Statement{Query: "SELECT COUNT(*) FROM sink"})
	require.NoError(t, err)
	// Assert the row exists before scanning: inside a bare `for res.Next()` these checks
	// are skipped on an empty result and the test passes without asserting anything.
	require.True(t, res.Next())
	var count int
	require.NoError(t, res.Rows.Scan(&count))
	require.Equal(t, 3, count)
	require.NoError(t, res.Err())
	require.NoError(t, res.Close())

	// decimal survives as a decimal with its precision and scale, and its value is exact
	res, err = olap.Query(context.Background(), &drivers.Statement{Query: "SELECT typeof(decimal_col), decimal_col::VARCHAR FROM sink WHERE int_col = 2147483647"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var typ, val string
	require.NoError(t, res.Rows.Scan(&typ, &val))
	require.Equal(t, "DECIMAL(10,2)", typ)
	require.Equal(t, "3.30", val)
	require.NoError(t, res.Err())
	require.NoError(t, res.Close())

	require.NoError(t, duckDB.Close())
}
