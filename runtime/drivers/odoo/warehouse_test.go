package odoo_test

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"os"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/apache/arrow-go/v18/parquet/pqarrow"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/drivers/odoo/odootest"
	"github.com/stretchr/testify/require"
)

func invoices() odootest.Model {
	return odootest.Model{
		Fields: map[string]string{
			"id": "integer", "name": "char", "invoice_date": "date", "write_date": "datetime",
			"amount_untaxed_signed": "monetary", "is_move_sent": "boolean", "partner_id": "many2one",
			"payment_state": "selection", "line_ids": "one2many",
		},
		Records: []map[string]any{
			{"id": 1, "name": "INV/2026/00001", "invoice_date": "2026-08-03", "write_date": "2026-09-01 02:00:00",
				"amount_untaxed_signed": 3900.5, "is_move_sent": true, "partner_id": []any{7, "Talleres Norte SL"}, "payment_state": "paid"},
			{"id": 2, "name": "INV/2026/00002", "invoice_date": false, "write_date": "2026-09-01 02:00:01",
				"amount_untaxed_signed": 0, "is_move_sent": false, "partner_id": false, "payment_state": "not_paid"},
			{"id": 5, "name": "INV/2026/00003", "invoice_date": "2026-08-04", "write_date": "2026-09-01 02:00:02",
				"amount_untaxed_signed": 120, "is_move_sent": false, "partner_id": []any{8, "Obras Sur SA"}, "payment_state": "not_paid"},
			{"id": 8, "name": "INV/2026/00004", "invoice_date": "2026-08-05", "write_date": "2026-09-01 02:00:03",
				"amount_untaxed_signed": 75.25, "is_move_sent": true, "partner_id": []any{7, "Talleres Norte SL"}, "payment_state": "in_payment"},
			{"id": 9, "name": "INV/2026/00005", "invoice_date": "2026-08-06", "write_date": "2026-09-01 02:00:04",
				"amount_untaxed_signed": 10, "is_move_sent": false, "partner_id": []any{9, "Instalaciones Centro SL"}, "payment_state": "paid"},
		},
	}
}

func openWarehouse(t *testing.T, srv *odootest.Server, extra map[string]any) drivers.Warehouse {
	cfg := map[string]any{"url": srv.URL, "database": "acme", "api_key": srv.APIKey}
	maps.Copy(cfg, extra)
	w, ok := openHandle(t, cfg).AsWarehouse()
	require.True(t, ok)
	return w
}

// readModel runs QueryAsFiles and returns the single parquet file it produces as a table.
func readModel(t *testing.T, w drivers.Warehouse, props map[string]any) arrow.Table {
	it, err := w.QueryAsFiles(t.Context(), props)
	require.NoError(t, err)
	t.Cleanup(func() { _ = it.Close() })
	files, err := it.Next(t.Context())
	require.NoError(t, err)
	require.Len(t, files, 1)
	_, err = it.Next(t.Context())
	require.ErrorIs(t, err, io.EOF)

	fh, err := os.Open(files[0])
	require.NoError(t, err)
	t.Cleanup(func() { _ = fh.Close() })
	tbl, err := pqarrow.ReadTable(t.Context(), fh, nil, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	require.NoError(t, err)
	t.Cleanup(tbl.Release)
	return tbl
}

func chunk(t *testing.T, tbl arrow.Table, col string) arrow.Array {
	idx := tbl.Schema().FieldIndices(col)
	require.Len(t, idx, 1, "column %s", col)
	return tbl.Column(idx[0]).Data().Chunk(0)
}

func TestReadsTypesAndTurnsOdooFalseIntoNull(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{"account.move": invoices()})
	tbl := readModel(t, openWarehouse(t, srv, nil), map[string]any{
		"model":  "account.move",
		"fields": []any{"name", "invoice_date", "write_date", "amount_untaxed_signed", "is_move_sent", "partner_id", "payment_state"},
	})

	require.Equal(t, int64(5), tbl.NumRows())
	var names []string
	for _, f := range tbl.Schema().Fields() {
		names = append(names, f.Name)
	}
	require.Equal(t, []string{"id", "name", "invoice_date", "write_date", "amount_untaxed_signed", "is_move_sent", "partner_id", "partner_id_name", "payment_state"}, names)

	require.Equal(t, int64(7), chunk(t, tbl, "partner_id").(*array.Int64).Value(0))
	require.Equal(t, "Talleres Norte SL", chunk(t, tbl, "partner_id_name").(*array.String).Value(0))
	require.True(t, chunk(t, tbl, "partner_id").IsNull(1))
	require.True(t, chunk(t, tbl, "partner_id_name").IsNull(1))
	require.True(t, chunk(t, tbl, "invoice_date").IsNull(1))
	require.False(t, chunk(t, tbl, "is_move_sent").IsNull(1))
	require.False(t, chunk(t, tbl, "is_move_sent").(*array.Boolean).Value(1))
	require.Equal(t, "2026-08-03", chunk(t, tbl, "invoice_date").(*array.Date32).Value(0).ToTime().Format(time.DateOnly))
	require.Equal(t, 3900.5, chunk(t, tbl, "amount_untaxed_signed").(*array.Float64).Value(0))
}

func TestPagesByIDAndKeepsTheModelDomain(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{"account.move": invoices()})
	tbl := readModel(t, openWarehouse(t, srv, nil), map[string]any{
		"model": "account.move", "fields": []any{"name"}, "page_size": 2,
		"domain": []any{[]any{"state", "=", "posted"}},
	})

	require.Equal(t, int64(5), tbl.NumRows())
	calls := srv.CallsTo("search_read")
	require.Len(t, calls, 3) // ids 1,2 | 5,8 | 9
	for i, floor := range []string{"0", "2", "8"} {
		domain := calls[i].Body["domain"].([]any)
		require.Equal(t, []any{"state", "=", "posted"}, domain[0])
		require.Equal(t, []any{"id", ">", json.Number(floor)}, domain[1])
		require.Equal(t, "id asc", calls[i].Body["order"])
		require.Equal(t, json.Number("2"), calls[i].Body["limit"])
	}
}

func TestEmptyModelYieldsTheSchemaWithNoRows(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{
		"sale.order": {Fields: map[string]string{"id": "integer", "name": "char"}},
	})
	tbl := readModel(t, openWarehouse(t, srv, nil), map[string]any{"model": "sale.order", "fields": []any{"name"}})
	require.Equal(t, int64(0), tbl.NumRows())
	require.Equal(t, "name", tbl.Schema().Field(1).Name)
}

func TestSendsLangAndModelContextOnEveryCall(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{"account.move": invoices()})
	readModel(t, openWarehouse(t, srv, map[string]any{"lang": "es_ES"}), map[string]any{
		"model": "account.move", "fields": []any{"name"}, "context": map[string]any{"active_test": false},
	})
	calls := srv.Calls()
	require.Len(t, calls, 2) // fields_get + one page
	for _, c := range calls {
		require.Equal(t, "bearer k", c.Header.Get("Authorization"))
		require.Equal(t, "acme", c.Header.Get("X-Odoo-Database"))
		require.Equal(t, map[string]any{"lang": "es_ES", "active_test": false}, c.Body["context"])
	}
}

func TestRejectsBadModelProperties(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{"account.move": invoices()})
	w := openWarehouse(t, srv, nil)
	cases := []struct {
		props map[string]any
		want  string
	}{
		{map[string]any{"fields": []any{"name"}}, `"model" is mandatory`},
		{map[string]any{"model": "account.move"}, `"fields" is mandatory`},
		{map[string]any{"model": "account.move", "fields": []any{"line_ids"}}, `type "one2many"`},
		{map[string]any{"model": "account.move", "fields": []any{"nope"}}, `no field "nope"`},
		{map[string]any{"model": "no.such.model", "fields": []any{"name"}}, "HTTP 404"},
	}
	for _, tc := range cases {
		_, err := w.QueryAsFiles(t.Context(), tc.props)
		require.ErrorContains(t, err, tc.want, "%v", tc.props)
	}
}

func TestWrongKeyFailsWithOdoosMessageAndWithoutTheKey(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{"account.move": invoices()})
	w := openWarehouse(t, srv, map[string]any{"api_key": "otra-clave"})
	_, err := w.QueryAsFiles(t.Context(), map[string]any{"model": "account.move", "fields": []any{"name"}})
	require.ErrorContains(t, err, "Invalid apikey")
	require.NotContains(t, err.Error(), "otra-clave")
}

func TestNextStopsOnACancelledContext(t *testing.T) {
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{"account.move": invoices()})
	ctx, cancel := context.WithCancel(t.Context())
	it, err := openWarehouse(t, srv, nil).QueryAsFiles(ctx, map[string]any{"model": "account.move", "fields": []any{"name"}})
	require.NoError(t, err)
	t.Cleanup(func() { _ = it.Close() })
	cancel()
	_, err = it.Next(ctx)
	require.ErrorIs(t, err, context.Canceled)
}
