package server_test

import (
	"context"
	"testing"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/testadmin"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestListProjectMemberAttributes(t *testing.T) {
	ctx := context.Background()
	fix := testadmin.New(t)

	// u1 creates the org and the project, which makes them an org admin and a project admin.
	u1, c1 := fix.NewUser(t)
	orgRes, err := c1.CreateOrganization(ctx, &adminv1.CreateOrganizationRequest{Name: randomName()})
	require.NoError(t, err)
	org := orgRes.Organization.Name

	// Clear the org's default project role, so that org members only reach the project through roles granted explicitly below.
	// Otherwise the autogroup:members usergroup gets a role on every new project and every org member is also a project member.
	_, err = c1.UpdateOrganization(ctx, &adminv1.UpdateOrganizationRequest{Org: org, DefaultProjectRole: toPtr("")})
	require.NoError(t, err)

	projRes, err := c1.CreateProject(ctx, &adminv1.CreateProjectRequest{Org: org, Project: "proj1", ProdSlots: 1, SkipDeploy: true})
	require.NoError(t, err)

	// u2 is a project admin and u3 a project editor. Neither is an org member, so they are added to the org as guests.
	u2, _ := fix.NewUser(t)
	_, err = c1.AddProjectMemberUser(ctx, &adminv1.AddProjectMemberUserRequest{Org: org, Project: "proj1", Email: u2.Email, Role: database.ProjectRoleNameAdmin})
	require.NoError(t, err)
	u3, _ := fix.NewUser(t)
	_, err = c1.AddProjectMemberUser(ctx, &adminv1.AddProjectMemberUserRequest{Org: org, Project: "proj1", Email: u3.Email, Role: database.ProjectRoleNameEditor})
	require.NoError(t, err)

	// u4 is a project viewer restricted to a single explore.
	u4, _ := fix.NewUser(t)
	_, err = c1.AddProjectMemberUser(ctx, &adminv1.AddProjectMemberUserRequest{
		Org:               org,
		Project:           "proj1",
		Email:             u4.Email,
		Role:              database.ProjectRoleNameViewer,
		RestrictResources: toPtr(true),
		Resources:         []*adminv1.ResourceName{{Type: runtime.ResourceKindExplore, Name: "e1"}},
	})
	require.NoError(t, err)

	// u5 is an org admin without any role on the project, so they have implicit access to it.
	u5, _ := fix.NewUser(t)
	_, err = c1.AddOrganizationMemberUser(ctx, &adminv1.AddOrganizationMemberUserRequest{Org: org, Email: u5.Email, Role: database.OrganizationRoleNameAdmin})
	require.NoError(t, err)

	// u6 only reaches the project through a usergroup, and carries a custom org attribute.
	u6, _ := fix.NewUser(t)
	attrs, err := structpb.NewStruct(map[string]any{"region": "eu"})
	require.NoError(t, err)
	_, err = c1.AddOrganizationMemberUser(ctx, &adminv1.AddOrganizationMemberUserRequest{Org: org, Email: u6.Email, Role: database.OrganizationRoleNameViewer, Attributes: attrs})
	require.NoError(t, err)
	_, err = c1.CreateUsergroup(ctx, &adminv1.CreateUsergroupRequest{Org: org, Name: "group1"})
	require.NoError(t, err)
	_, err = c1.AddUsergroupMemberUser(ctx, &adminv1.AddUsergroupMemberUserRequest{Org: org, Usergroup: "group1", Email: u6.Email})
	require.NoError(t, err)
	_, err = c1.AddProjectMemberUsergroup(ctx, &adminv1.AddProjectMemberUsergroupRequest{Org: org, Project: "proj1", Usergroup: "group1", Role: database.ProjectRoleNameViewer})
	require.NoError(t, err)

	// u7 has no relation to the org at all.
	u7, c7 := fix.NewUser(t)

	res, err := c1.ListProjectMemberAttributes(ctx, &adminv1.ListProjectMemberAttributesRequest{ProjectId: projRes.Project.Id})
	require.NoError(t, err)

	members := make(map[string]*adminv1.ProjectMemberAttributes, len(res.Members))
	emails := make([]string, 0, len(res.Members))
	for _, m := range res.Members {
		members[m.Email] = m
		emails = append(emails, m.Email)
	}
	require.ElementsMatch(t, []string{u1.Email, u2.Email, u3.Email, u4.Email, u5.Email, u6.Email}, emails)
	require.NotContains(t, members, u7.Email)

	t.Run("Edit trigger follows manage permissions on the prod deployment", func(t *testing.T) {
		// Project admins and org admins can manage prod; editors and viewers can't.
		require.True(t, members[u1.Email].EditTrigger)
		require.True(t, members[u2.Email].EditTrigger)
		require.False(t, members[u3.Email].EditTrigger)
		require.False(t, members[u4.Email].EditTrigger)
		require.True(t, members[u5.Email].EditTrigger)
		require.False(t, members[u6.Email].EditTrigger)
	})

	t.Run("Attributes match the ones issued in a runtime JWT", func(t *testing.T) {
		attr := members[u6.Email].Attributes.AsMap()
		require.Equal(t, u6.Email, attr["email"])
		require.Equal(t, "test-user.com", attr["domain"])
		require.Equal(t, false, attr["admin"])
		require.Equal(t, "eu", attr["region"], "custom org attributes must be included")
		require.Contains(t, attr["groups"], "group1")

		// "admin" reports project management permission, which the project admin has and the editor doesn't.
		require.Equal(t, true, members[u2.Email].Attributes.AsMap()["admin"])
		require.Equal(t, false, members[u3.Email].Attributes.AsMap()["admin"])
	})

	t.Run("Security rules carry the member's resource restrictions", func(t *testing.T) {
		require.Empty(t, members[u2.Email].SecurityRules)
		require.Empty(t, members[u5.Email].SecurityRules)

		rules := members[u4.Email].SecurityRules
		require.Len(t, rules, 1)
		transitive := rules[0].GetTransitiveAccess()
		require.NotNil(t, transitive)
		require.Equal(t, runtime.ResourceKindExplore, transitive.Resource.Kind)
		require.Equal(t, "e1", transitive.Resource.Name)
	})

	t.Run("A user without access to the project can't list its members", func(t *testing.T) {
		_, err := c7.ListProjectMemberAttributes(ctx, &adminv1.ListProjectMemberAttributesRequest{ProjectId: projRes.Project.Id})
		require.Error(t, err)
		require.Equal(t, codes.PermissionDenied, status.Code(err))
	})
}
