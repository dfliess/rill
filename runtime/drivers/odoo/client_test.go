package odoo

import (
	"encoding/json"
	"testing"

	"github.com/rilldata/rill/runtime/drivers/odoo/odootest"
	"github.com/stretchr/testify/require"
)

func TestNewClientValidatesConfig(t *testing.T) {
	_, err := newClient(&configProperties{APIKey: "k"})
	require.ErrorContains(t, err, `"url" is required`)
	_, err = newClient(&configProperties{URL: "http://odoo"})
	require.ErrorContains(t, err, `"api_key" is required`)
	_, err = newClient(&configProperties{URL: "odoo.example.com", APIKey: "k"})
	require.ErrorContains(t, err, "invalid url")
}

func TestCallSendsKeyDatabaseAndNamedArguments(t *testing.T) {
	srv := odootest.NewServer(t, "clave", nil)
	c, err := newClient(&configProperties{URL: srv.URL + "/", Database: "acme", APIKey: "clave"})
	require.NoError(t, err)

	var out map[string]any
	require.NoError(t, c.call(t.Context(), "res.users", "context_get", map[string]any{"context": map[string]any{"lang": "es_ES"}}, &out))

	require.Equal(t, "Europe/Madrid", out["tz"])
	call := srv.Calls()[0]
	require.Equal(t, "bearer clave", call.Header.Get("Authorization"))
	require.Equal(t, "acme", call.Header.Get("X-Odoo-Database"))
	require.Equal(t, map[string]any{"lang": "es_ES"}, call.Body["context"])
}

func TestCallKeepsOdoosMessageAndHidesKeyAndTraceback(t *testing.T) {
	srv := odootest.NewServer(t, "la-buena", nil)
	c, err := newClient(&configProperties{URL: srv.URL, APIKey: "la-mala-y-secreta"})
	require.NoError(t, err)

	var out any
	err = c.call(t.Context(), "res.users", "context_get", map[string]any{}, &out)
	require.ErrorContains(t, err, "HTTP 401")
	require.ErrorContains(t, err, "Invalid apikey")
	require.NotContains(t, err.Error(), "la-mala-y-secreta")
	require.NotContains(t, err.Error(), "Traceback")
}

func TestCallDecodesIntegersExactly(t *testing.T) {
	// 2^53 + 1 does not survive a float64: ids have to arrive as json.Number.
	srv := odootest.NewServer(t, "k", map[string]odootest.Model{
		"res.partner": {Fields: map[string]string{"id": "integer"}, Records: []map[string]any{{"id": int64(9007199254740993)}}},
	})
	c, err := newClient(&configProperties{URL: srv.URL, APIKey: "k"})
	require.NoError(t, err)

	var rows []map[string]any
	require.NoError(t, c.call(t.Context(), "res.partner", "search_read", map[string]any{"domain": []any{}, "fields": []string{"id"}}, &rows))
	require.Equal(t, json.Number("9007199254740993"), rows[0]["id"])
}
