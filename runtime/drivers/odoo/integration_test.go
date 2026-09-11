package odoo_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAgainstLocalOdoo reads the seeded instance zero of Kairos para Odoo. It is skipped unless
// ODOO_TEST_URL is set; see kairos-cloud/docs/kairos-para-odoo/plan-1-instancia-cero-local.md.
func TestAgainstLocalOdoo(t *testing.T) {
	url := os.Getenv("ODOO_TEST_URL")
	if url == "" {
		t.Skip("ODOO_TEST_URL is not set: this test needs the local Odoo of plan 1")
	}
	h := openHandle(t, map[string]any{
		"url": url, "database": os.Getenv("ODOO_TEST_DB"), "api_key": os.Getenv("ODOO_TEST_API_KEY"), "lang": "es_ES",
	})
	require.NoError(t, h.Ping(t.Context()))
	w, ok := h.AsWarehouse()
	require.True(t, ok)

	clientes := readModel(t, w, map[string]any{
		"model":     "res.partner",
		"fields":    []any{"name", "vat", "user_id", "write_date"},
		"domain":    []any{[]any{"customer_rank", ">", 0}, []any{"is_company", "=", true}},
		"page_size": 50,
	})
	require.Equal(t, int64(140), clientes.NumRows())

	facturas := readModel(t, w, map[string]any{
		"model": "account.move",
		"fields": []any{"name", "invoice_date", "invoice_date_due", "amount_untaxed_signed",
			"amount_residual_signed", "payment_state", "partner_id", "invoice_user_id", "company_id"},
		"domain": []any{[]any{"move_type", "=", "out_invoice"}, []any{"state", "=", "posted"}},
	})
	require.Greater(t, facturas.NumRows(), int64(2500))
}
