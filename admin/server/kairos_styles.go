package server

import (
	"context"
	"slices"

	"github.com/rilldata/rill/admin/server/auth"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Kairos Styles (kairos-cloud#242): the app's look an organization's members get unless they pick their own.

// kairosStyles are the Styles an organization can pick, in the order of STYLES in
// web-common/src/features/kairos-styles/styles.ts (scripts/style-guard.js checks it).
// Validated here rather than in the proto, so adding a Style doesn't regenerate api.pb.go.
var kairosStyles = []string{"kairos", "aurora", "nova", "editorial", "lumen", "glass", "apetic", "base"}

// kairosStyleEmbedParam is the embed URL parameter the frontend paints the Style from.
// It is a Kairos name so it can't collide with a parameter upstream adds.
const kairosStyleEmbedParam = "kairos_style"

func (s *Server) GetOrganizationStyle(ctx context.Context, req *adminv1.GetOrganizationStyleRequest) (*adminv1.GetOrganizationStyleResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Org),
	)

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Org)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ReadOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to read org")
	}

	style, err := s.admin.DB.FindOrganizationStyle(ctx, org.ID)
	if err != nil {
		return nil, err
	}

	return &adminv1.GetOrganizationStyleResponse{Style: style}, nil
}

func (s *Server) UpdateOrganizationStyle(ctx context.Context, req *adminv1.UpdateOrganizationStyleRequest) (*adminv1.UpdateOrganizationStyleResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Org),
		attribute.String("args.style", req.Style),
	)

	if req.Style != "" && !slices.Contains(kairosStyles, req.Style) {
		return nil, status.Errorf(codes.InvalidArgument, "unknown style %q", req.Style)
	}

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Org)
	if err != nil {
		return nil, err
	}

	claims := auth.GetClaims(ctx)
	if !claims.OrganizationPermissions(ctx, org.ID).ManageOrg && !claims.Superuser(ctx) {
		return nil, status.Error(codes.PermissionDenied, "not allowed to manage org")
	}

	err = s.admin.DB.UpdateOrganizationStyle(ctx, org.ID, req.Style)
	if err != nil {
		return nil, err
	}

	return &adminv1.UpdateOrganizationStyleResponse{Style: req.Style}, nil
}

// addOrganizationStyle sets the embed's Style parameter to the organization's, unless the caller passed one,
// so an embed looks like the app its organization uses (GetIFrame).
func (s *Server) addOrganizationStyle(ctx context.Context, orgID string, query map[string]string) error {
	if _, ok := query[kairosStyleEmbedParam]; ok {
		return nil
	}
	style, err := s.admin.DB.FindOrganizationStyle(ctx, orgID)
	if err != nil {
		return err
	}
	if style != "" {
		query[kairosStyleEmbedParam] = style
	}
	return nil
}
