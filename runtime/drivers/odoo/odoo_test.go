package odoo_test

import (
	"testing"

	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/drivers/odoo/odootest"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/rilldata/rill/runtime/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	_ "github.com/rilldata/rill/runtime/drivers/odoo"
)

func openHandle(t *testing.T, cfg map[string]any) drivers.Handle {
	h, err := drivers.Open("odoo", "odoo", "default", cfg, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	t.Cleanup(func() { _ = h.Close() })
	return h
}

func TestOpenGivesAWarehouseAndNothingElse(t *testing.T) {
	h := openHandle(t, map[string]any{"url": "http://localhost:1", "api_key": "k"})
	require.Equal(t, "odoo", h.Driver())
	_, ok := h.AsWarehouse()
	require.True(t, ok)
	_, ok = h.AsOLAP("default")
	require.False(t, ok)

	_, err := drivers.Open("odoo", "odoo", "", map[string]any{}, storage.MustNew(t.TempDir(), nil), activity.NewNoopClient(), zap.NewNop())
	require.ErrorContains(t, err, "can't be shared")
}

func TestSpecMarksTheAPIKeyAsSecret(t *testing.T) {
	secret := map[string]bool{}
	for _, p := range drivers.Drivers["odoo"].Spec().ConfigProperties {
		secret[p.Key] = p.Secret
	}
	require.Equal(t, map[string]bool{"url": false, "database": false, "api_key": true, "lang": false}, secret)
}

func TestPing(t *testing.T) {
	srv := odootest.NewServer(t, "k", nil)
	require.NoError(t, openHandle(t, map[string]any{"url": srv.URL, "api_key": "k"}).Ping(t.Context()))
	require.ErrorContains(t, openHandle(t, map[string]any{"url": srv.URL, "api_key": "otra"}).Ping(t.Context()), "Invalid apikey")
	require.ErrorContains(t, openHandle(t, map[string]any{"api_key": "k"}).Ping(t.Context()), `"url" is required`)
}
