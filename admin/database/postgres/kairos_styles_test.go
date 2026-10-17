package postgres

import (
	"context"
	"testing"

	"github.com/rilldata/rill/admin/database"
	"github.com/stretchr/testify/require"
)

// testOrganizationStyles runs inside TestPostgres (kairos-cloud#242).
func testOrganizationStyles(t *testing.T, db database.DB) {
	ctx := context.Background()

	org, err := db.InsertOrganization(ctx, &database.InsertOrganizationOptions{Name: randomName()})
	require.NoError(t, err)
	other, err := db.InsertOrganization(ctx, &database.InsertOrganizationOptions{Name: randomName()})
	require.NoError(t, err)

	// No row means the default Style
	style, err := db.FindOrganizationStyle(ctx, org.ID)
	require.NoError(t, err)
	require.Equal(t, "", style)

	require.NoError(t, db.UpdateOrganizationStyle(ctx, org.ID, "aurora"))
	require.NoError(t, db.UpdateOrganizationStyle(ctx, org.ID, "editorial"))
	style, err = db.FindOrganizationStyle(ctx, org.ID)
	require.NoError(t, err)
	require.Equal(t, "editorial", style)

	// Another organization is untouched
	style, err = db.FindOrganizationStyle(ctx, other.ID)
	require.NoError(t, err)
	require.Equal(t, "", style)

	// Upstream's organization update rewrites every orgs column; the Style lives outside it and survives
	_, err = db.UpdateOrganization(ctx, org.ID, &database.UpdateOrganizationOptions{
		Name:                                org.Name,
		DisplayName:                         "Renamed",
		DefaultProjectRoleID:                org.DefaultProjectRoleID,
		QuotaProjects:                       org.QuotaProjects,
		QuotaDeployments:                    org.QuotaDeployments,
		QuotaSlotsTotal:                     org.QuotaSlotsTotal,
		QuotaSlotsPerDeployment:             org.QuotaSlotsPerDeployment,
		QuotaOutstandingInvites:             org.QuotaOutstandingInvites,
		QuotaStorageLimitBytesPerDeployment: org.QuotaStorageLimitBytesPerDeployment,
		QuotaSeats:                          org.QuotaSeats,
	})
	require.NoError(t, err)
	style, err = db.FindOrganizationStyle(ctx, org.ID)
	require.NoError(t, err)
	require.Equal(t, "editorial", style)

	// "" restores the default
	require.NoError(t, db.UpdateOrganizationStyle(ctx, org.ID, ""))
	style, err = db.FindOrganizationStyle(ctx, org.ID)
	require.NoError(t, err)
	require.Equal(t, "", style)

	// Deleting the organization deletes its Style
	require.NoError(t, db.UpdateOrganizationStyle(ctx, other.ID, "aurora"))
	require.NoError(t, db.DeleteOrganization(ctx, other.Name))
	style, err = db.FindOrganizationStyle(ctx, other.ID)
	require.NoError(t, err)
	require.Equal(t, "", style)
}
