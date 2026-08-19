package server

import (
	"context"
	"slices"
	"strings"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/server/auth"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// projectMemberPageSize is the page size used when enumerating a project's members.
// All members are returned regardless, so it only controls the number of database round trips.
const projectMemberPageSize = 1000

// ListProjectMemberAttributes returns the project's members with the identity a runtime JWT issued for them would carry.
// The runtime uses it to evaluate security policies for users that are not currently making a request,
// such as finding the members authorized to approve an agent's proposed action.
func (s *Server) ListProjectMemberAttributes(ctx context.Context, req *adminv1.ListProjectMemberAttributesRequest) (*adminv1.ListProjectMemberAttributesResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.project_id", req.ProjectId),
	)

	proj, err := s.admin.DB.FindProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	if !claims.ProjectPermissions(ctx, proj.OrganizationID, proj.ID).ReadProdStatus {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to read the project's members")
	}

	// Which permission grants EditTrigger depends on the environment of the deployment the attributes are for; see issueRuntimeToken.
	// The runtime calls this RPC with its deployment's access token; for other callers we describe the prod deployment.
	environment := "prod"
	if claims.OwnerType() == auth.OwnerTypeDeployment {
		depl, err := s.admin.DB.FindDeployment(ctx, claims.OwnerID())
		if err != nil {
			return nil, err
		}
		environment = depl.Environment
	}

	users, err := s.projectMemberUsers(ctx, proj)
	if err != nil {
		return nil, err
	}

	members := make([]*adminv1.ProjectMemberAttributes, len(users))
	for i, u := range users {
		orgPerms, err := s.admin.OrganizationPermissionsForUser(ctx, proj.OrganizationID, u.id)
		if err != nil {
			return nil, err
		}
		projPerms, err := s.admin.ProjectPermissionsForUser(ctx, proj.ID, u.id, orgPerms)
		if err != nil {
			return nil, err
		}

		attr, err := s.jwtAttributesForUser(ctx, u.id, proj.OrganizationID, projPerms)
		if err != nil {
			return nil, err
		}
		attrPB, err := structpb.NewStruct(attr)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to encode attributes for user %q: %s", u.id, err.Error())
		}

		restrictResources, resources, err := s.getResourceRestrictionsForUser(ctx, proj.ID, u.id)
		if err != nil {
			return nil, err
		}

		_, canManage := elevatedPermissionsForEnvironment(environment, projPerms)

		members[i] = &adminv1.ProjectMemberAttributes{
			UserId:        u.id,
			Email:         u.email,
			Attributes:    attrPB,
			EditTrigger:   canManage,
			SecurityRules: securityRulesFromResources(restrictResources, resources),
		}
	}

	return &adminv1.ListProjectMemberAttributesResponse{Members: members}, nil
}

// projectMemberUser identifies a user with access to a project.
type projectMemberUser struct {
	id    string
	email string
}

// projectMemberUsers returns every user that project permissions can resolve for, i.e. everyone with a direct role on the project,
// the members of the usergroups that hold a role on the project, and the org's admins (who have ManageProjects,
// which grants implicit access to every project in the org; see ProjectPermissionsForUser).
// The result is deduplicated by user ID and sorted by email.
func (s *Server) projectMemberUsers(ctx context.Context, proj *database.Project) ([]projectMemberUser, error) {
	directMembers, err := collectPages(
		func(after string, limit int) ([]*database.ProjectMemberUser, error) {
			return s.admin.DB.FindProjectMemberUsers(ctx, proj.OrganizationID, proj.ID, "", after, limit)
		},
		func(m *database.ProjectMemberUser) string { return m.Email },
	)
	if err != nil {
		return nil, err
	}

	adminRole, err := s.admin.DB.FindOrganizationRole(ctx, database.OrganizationRoleNameAdmin)
	if err != nil {
		return nil, err
	}

	orgAdmins, err := s.admin.DB.FindOrganizationMemberUsersByRole(ctx, proj.OrganizationID, adminRole.ID)
	if err != nil {
		return nil, err
	}

	// Roles are also granted through usergroups, both on the project and on the org, and a group's members do not show up in the queries above.
	projectGroups, err := collectPages(
		func(after string, limit int) ([]*database.MemberUsergroup, error) {
			return s.admin.DB.FindProjectMemberUsergroups(ctx, proj.ID, "", false, after, limit)
		},
		func(g *database.MemberUsergroup) string { return g.Name },
	)
	if err != nil {
		return nil, err
	}
	orgAdminGroups, err := collectPages(
		func(after string, limit int) ([]*database.MemberUsergroup, error) {
			return s.admin.DB.FindOrganizationMemberUsergroups(ctx, proj.OrganizationID, adminRole.ID, false, after, limit)
		},
		func(g *database.MemberUsergroup) string { return g.Name },
	)
	if err != nil {
		return nil, err
	}

	var res []projectMemberUser
	seenUsers := make(map[string]bool)
	add := func(id, email string) {
		if seenUsers[id] {
			return
		}
		seenUsers[id] = true
		res = append(res, projectMemberUser{id: id, email: email})
	}
	for _, m := range directMembers {
		add(m.ID, m.Email)
	}
	for _, u := range orgAdmins {
		add(u.ID, u.Email)
	}

	seenGroups := make(map[string]bool)
	for _, g := range slices.Concat(projectGroups, orgAdminGroups) {
		if seenGroups[g.ID] {
			continue
		}
		seenGroups[g.ID] = true

		groupMembers, err := collectPages(
			func(after string, limit int) ([]*database.UsergroupMemberUser, error) {
				return s.admin.DB.FindUsergroupMemberUsers(ctx, g.ID, after, limit)
			},
			func(u *database.UsergroupMemberUser) string { return u.Email },
		)
		if err != nil {
			return nil, err
		}
		for _, u := range groupMembers {
			add(u.ID, u.Email)
		}
	}

	slices.SortFunc(res, func(a, b projectMemberUser) int {
		return strings.Compare(a.email, b.email)
	})
	return res, nil
}

// collectPages returns all items of a paginated lookup.
// It calls fetch until it returns a partial page, passing the cursor of the previous page's last item.
func collectPages[T any](fetch func(after string, limit int) ([]T, error), cursor func(T) string) ([]T, error) {
	var res []T
	var after string
	for {
		page, err := fetch(after, projectMemberPageSize)
		if err != nil {
			return nil, err
		}
		res = append(res, page...)
		if len(page) < projectMemberPageSize {
			return res, nil
		}
		after = cursor(page[len(page)-1])
	}
}
