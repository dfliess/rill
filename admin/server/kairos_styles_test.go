package server_test

import (
	"context"
	"testing"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/testadmin"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestOrganizationStyles covers who reads and who sets an organization's Kairos Style (kairos-cloud#242).
func TestOrganizationStyles(t *testing.T) {
	ctx := context.Background()
	fix := testadmin.New(t)

	// The fixture's first user is a superuser; skip it so permissions are the org's own.
	fix.NewUser(t)
	_, admin := fix.NewUser(t)
	viewerUser, viewer := fix.NewUser(t)
	_, outsider := fix.NewUser(t)

	org, err := admin.CreateOrganization(ctx, &adminv1.CreateOrganizationRequest{Name: randomName()})
	require.NoError(t, err)
	name := org.Organization.Name
	_, err = admin.AddOrganizationMemberUser(ctx, &adminv1.AddOrganizationMemberUserRequest{Org: name, Email: viewerUser.Email, Role: database.OrganizationRoleNameViewer})
	require.NoError(t, err)

	// Default until an admin sets one
	got, err := viewer.GetOrganizationStyle(ctx, &adminv1.GetOrganizationStyleRequest{Org: name})
	require.NoError(t, err)
	require.Equal(t, "", got.Style)

	_, err = admin.UpdateOrganizationStyle(ctx, &adminv1.UpdateOrganizationStyleRequest{Org: name, Style: "aurora"})
	require.NoError(t, err)
	got, err = viewer.GetOrganizationStyle(ctx, &adminv1.GetOrganizationStyleRequest{Org: name})
	require.NoError(t, err)
	require.Equal(t, "aurora", got.Style)

	// A member who can't manage the org can't change it
	_, err = viewer.UpdateOrganizationStyle(ctx, &adminv1.UpdateOrganizationStyleRequest{Org: name, Style: "kairos-v1"})
	require.Equal(t, codes.PermissionDenied, status.Code(err))

	// Someone outside the org can't read it
	_, err = outsider.GetOrganizationStyle(ctx, &adminv1.GetOrganizationStyleRequest{Org: name})
	require.Equal(t, codes.PermissionDenied, status.Code(err))

	// Only known Styles
	_, err = admin.UpdateOrganizationStyle(ctx, &adminv1.UpdateOrganizationStyleRequest{Org: name, Style: "retired"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// Empty restores the default
	_, err = admin.UpdateOrganizationStyle(ctx, &adminv1.UpdateOrganizationStyleRequest{Org: name, Style: ""})
	require.NoError(t, err)
	got, err = viewer.GetOrganizationStyle(ctx, &adminv1.GetOrganizationStyleRequest{Org: name})
	require.NoError(t, err)
	require.Equal(t, "", got.Style)
}
