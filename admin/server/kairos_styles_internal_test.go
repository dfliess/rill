package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/rilldata/rill/admin"
	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/pkg/pgtestcontainer"
	"github.com/stretchr/testify/require"
)

// TestAddOrganizationStyle covers the Style GetIFrame puts in an embed's URL (kairos-cloud#242).
func TestAddOrganizationStyle(t *testing.T) {
	ctx := context.Background()
	pg := pgtestcontainer.New(t)
	defer pg.Terminate(t)

	keyring, err := database.NewRandomKeyring()
	require.NoError(t, err)
	conf, err := json.Marshal(keyring)
	require.NoError(t, err)
	db, err := database.Open("postgres", pg.DatabaseURL, string(conf))
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()
	require.NoError(t, db.Migrate(ctx))

	s := &Server{admin: &admin.Service{DB: db}}
	org, err := db.InsertOrganization(ctx, &database.InsertOrganizationOptions{Name: "styled"})
	require.NoError(t, err)

	// No Style set: the embed gets no parameter and paints the default
	query := map[string]string{}
	require.NoError(t, s.addOrganizationStyle(ctx, org.ID, query))
	require.NotContains(t, query, kairosStyleEmbedParam)

	require.NoError(t, db.UpdateOrganizationStyle(ctx, org.ID, "aurora"))
	query = map[string]string{}
	require.NoError(t, s.addOrganizationStyle(ctx, org.ID, query))
	require.Equal(t, "aurora", query[kairosStyleEmbedParam])

	// The embedder's own choice wins
	query = map[string]string{kairosStyleEmbedParam: "kairos-v1"}
	require.NoError(t, s.addOrganizationStyle(ctx, org.ID, query))
	require.Equal(t, "kairos-v1", query[kairosStyleEmbedParam])
}
