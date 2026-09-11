package duckdb_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/drivers/odoo/odootest"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	_ "github.com/rilldata/rill/runtime/drivers/duckdb"
	_ "github.com/rilldata/rill/runtime/drivers/odoo"
)

func TestOdooToDuckDBTransfer(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{
		"account.move": {
			Fields: map[string]string{
				"id": "integer", "name": "char", "invoice_date": "date", "write_date": "datetime",
				"amount_untaxed_signed": "monetary", "partner_id": "many2one", "is_move_sent": "boolean",
			},
			Records: []map[string]any{
				{"id": 1, "name": "INV/1", "invoice_date": "2026-08-03", "write_date": "2026-09-01 02:00:00",
					"amount_untaxed_signed": 3900.5, "partner_id": []any{7, "Talleres Norte SL"}, "is_move_sent": true},
				{"id": 2, "name": "INV/2", "invoice_date": false, "write_date": "2026-09-01 02:00:01",
					"amount_untaxed_signed": 0, "partner_id": false, "is_move_sent": false},
			},
		},
		"sale.order": {Fields: map[string]string{"id": "integer", "name": "char"}},
	})

	facturas := odooToDuckDB(t, srv, map[string]any{
		"model":  "account.move",
		"fields": []any{"name", "invoice_date", "write_date", "amount_untaxed_signed", "partner_id", "is_move_sent"},
	})
	res, err := facturas.Query(t.Context(), &drivers.Statement{Query: `
		SELECT count(*), count(invoice_date), count(partner_id), typeof(any_value(id)), typeof(any_value(name)),
		       typeof(any_value(invoice_date)), typeof(any_value(write_date)), typeof(any_value(amount_untaxed_signed)),
		       typeof(any_value(partner_id)), typeof(any_value(partner_id_name)), typeof(any_value(is_move_sent))
		FROM sink`})
	require.NoError(t, err)
	// Assert the row exists before scanning: inside a bare `for res.Next()` the checks would be
	// skipped on an empty result and the test would pass without asserting anything.
	require.True(t, res.Next())
	var n, withDate, withPartner int
	var types [8]string
	require.NoError(t, res.Rows.Scan(&n, &withDate, &withPartner, &types[0], &types[1], &types[2], &types[3], &types[4], &types[5], &types[6], &types[7]))
	require.Equal(t, 2, n)
	require.Equal(t, 1, withDate)    // Odoo's false became NULL
	require.Equal(t, 1, withPartner) // likewise for a many2one
	require.Equal(t, [8]string{"BIGINT", "VARCHAR", "DATE", "TIMESTAMP WITH TIME ZONE", "DOUBLE", "BIGINT", "VARCHAR", "BOOLEAN"}, types)
	require.NoError(t, res.Err())
	require.NoError(t, res.Close())

	// An empty Odoo model still creates the table, with its columns and no rows.
	pedidos := odooToDuckDB(t, srv, map[string]any{"model": "sale.order", "fields": []any{"name"}})
	res, err = pedidos.Query(t.Context(), &drivers.Statement{Query: "SELECT count(*), count(name) FROM sink"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var rows, names int
	require.NoError(t, res.Rows.Scan(&rows, &names))
	require.Equal(t, 0, rows)
	require.Equal(t, 0, names)
	require.NoError(t, res.Close())
}

// odooToDuckDB runs one odoo model into a fresh DuckDB table called "sink".
func odooToDuckDB(t *testing.T, srv *odootest.Server, props map[string]any) drivers.OLAPStore {
	duckDB, err := drivers.Open("duckdb", "", "default", map[string]any{"data_dir": t.TempDir()}, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	t.Cleanup(func() { _ = duckDB.Close() })
	input, err := drivers.Open("odoo", "odoo", "default", map[string]any{"url": srv.URL, "api_key": srv.APIKey}, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	t.Cleanup(func() { _ = input.Close() })

	opts := &drivers.ModelExecutorOptions{
		InputHandle:                 input,
		InputConnector:              "odoo",
		OutputHandle:                duckDB,
		OutputConnector:             "duckdb",
		Env:                         &drivers.ModelEnv{AllowHostAccess: false, StageChanges: true},
		PreliminaryInputProperties:  props,
		PreliminaryOutputProperties: map[string]any{"table": "sink"},
	}
	me, err := duckDB.AsModelExecutor("default", opts)
	require.NoError(t, err)
	_, err = me.Execute(t.Context(), &drivers.ModelExecuteOptions{
		ModelExecutorOptions: opts,
		InputProperties:      opts.PreliminaryInputProperties,
		OutputProperties:     opts.PreliminaryOutputProperties,
	})
	require.NoError(t, err)

	olap, ok := duckDB.AsOLAP("default")
	require.True(t, ok)
	return olap
}
