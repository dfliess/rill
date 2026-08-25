package server

import (
	"context"
	"net/url"
	"strings"

	"github.com/rilldata/rill/admin"
	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/pkg/pushnotifications"
	"github.com/rilldata/rill/admin/server/auth"
	adminv1 "github.com/rilldata/rill/proto/gen/rill/admin/v1"
	"github.com/rilldata/rill/runtime/pkg/observability"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetPushNotificationConfig(ctx context.Context, req *adminv1.GetPushNotificationConfigRequest) (*adminv1.GetPushNotificationConfigResponse, error) {
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	return &adminv1.GetPushNotificationConfigResponse{
		VapidPublicKey: s.admin.Push.VAPIDPublicKey(),
	}, nil
}

func (s *Server) CreatePushSubscription(ctx context.Context, req *adminv1.CreatePushSubscriptionRequest) (*adminv1.CreatePushSubscriptionResponse, error) {
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	if !s.admin.Push.Enabled() {
		return nil, status.Error(codes.FailedPrecondition, "push notifications are not enabled in this deployment")
	}

	_, err := s.admin.DB.InsertPushSubscription(ctx, &database.InsertPushSubscriptionOptions{
		UserID:    claims.OwnerID(),
		Endpoint:  req.Endpoint,
		P256dh:    req.P256Dh,
		Auth:      req.Auth,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	return &adminv1.CreatePushSubscriptionResponse{}, nil
}

func (s *Server) ListPushSubscriptions(ctx context.Context, req *adminv1.ListPushSubscriptionsRequest) (*adminv1.ListPushSubscriptionsResponse, error) {
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	subs, err := s.admin.DB.FindPushSubscriptionsForUser(ctx, claims.OwnerID())
	if err != nil {
		return nil, err
	}

	dtos := make([]*adminv1.PushSubscription, len(subs))
	for i, sub := range subs {
		dtos[i] = &adminv1.PushSubscription{
			Id:        sub.ID,
			Endpoint:  sub.Endpoint,
			UserAgent: sub.UserAgent,
			CreatedOn: timestamppb.New(sub.CreatedOn),
		}
	}

	return &adminv1.ListPushSubscriptionsResponse{
		Subscriptions: dtos,
	}, nil
}

func (s *Server) DeletePushSubscription(ctx context.Context, req *adminv1.DeletePushSubscriptionRequest) (*adminv1.DeletePushSubscriptionResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.id", req.Id),
	)

	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	err := s.admin.DB.DeletePushSubscription(ctx, req.Id, claims.OwnerID())
	if err != nil {
		return nil, err
	}

	return &adminv1.DeletePushSubscriptionResponse{}, nil
}

// notificationHistoryLimit bounds one read of the history. It is a block on a settings page, meant to answer "what
// have I been told lately", so it is a fortnight's worth at the rate a real deployment notifies, not an archive.
// Anything older is still stored until the retention sweep takes it; it is simply not what this list is for.
const notificationHistoryLimit = 20

// notificationHistoryWindow resolves how many entries one read shows. A caller may ask for fewer than the limit,
// never for more: the limit is what the deployment is willing to serve, not a default the client can talk it out of.
func notificationHistoryWindow(pageSize int32) int {
	if pageSize > 0 && int(pageSize) < notificationHistoryLimit {
		return int(pageSize)
	}
	return notificationHistoryLimit
}

// trimToWindow cuts a listing read one over the window down to it, and reports whether that extra entry was there.
// A list that quietly stops reads as "this is everything", which is a different thing to say.
func trimToWindow(deliveries []*database.NotificationDelivery, limit int) ([]*database.NotificationDelivery, bool) {
	if len(deliveries) > limit {
		return deliveries[:limit], true
	}
	return deliveries, false
}

func (s *Server) ListNotifications(ctx context.Context, req *adminv1.ListNotificationsRequest) (*adminv1.ListNotificationsResponse, error) {
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	// The history is read for the caller and only for the caller: there is no user ID in the request to get wrong.
	limit := notificationHistoryWindow(req.PageSize)

	// One more than the window, to learn whether there is anything behind it without counting the table.
	deliveries, err := s.admin.DB.FindNotificationDeliveriesForUser(ctx, claims.OwnerID(), limit+1)
	if err != nil {
		return nil, err
	}
	deliveries, hasMore := trimToWindow(deliveries, limit)

	dtos := make([]*adminv1.Notification, len(deliveries))
	for i, d := range deliveries {
		dtos[i] = &adminv1.Notification{
			Id:       d.ID,
			Category: d.Category,
			Org:      d.OrgName,
			Project:  d.ProjectName,
			Title:    d.Title,
			Body:     d.Body,
			LinkPath: d.LinkPath,
			SentOn:   timestamppb.New(d.SentOn),
		}
	}
	return &adminv1.ListNotificationsResponse{Notifications: dtos, HasMore: hasMore}, nil
}

func (s *Server) ListNotificationPreferences(ctx context.Context, req *adminv1.ListNotificationPreferencesRequest) (*adminv1.ListNotificationPreferencesResponse, error) {
	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	prefs, err := s.admin.DB.FindNotificationPreferencesForUser(ctx, claims.OwnerID())
	if err != nil {
		return nil, err
	}

	dtos := make([]*adminv1.OrganizationNotificationPreferences, len(prefs))
	for i, pref := range prefs {
		dtos[i] = &adminv1.OrganizationNotificationPreferences{
			Org:            pref.OrgName,
			OrgDisplayName: pref.OrgDisplayName,
			Preferences: &adminv1.NotificationPreferences{
				PushAlerts:       pref.PushAlerts,
				PushReports:      pref.PushReports,
				PushActApprovals: pref.PushActApprovals,
			},
		}
	}

	return &adminv1.ListNotificationPreferencesResponse{
		Organizations: dtos,
	}, nil
}

func (s *Server) GetNotificationPreferences(ctx context.Context, req *adminv1.GetNotificationPreferencesRequest) (*adminv1.GetNotificationPreferencesResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Org),
	)

	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Org)
	if err != nil {
		return nil, err
	}
	if !claims.OrganizationPermissions(ctx, org.ID).ReadOrg {
		return nil, status.Error(codes.PermissionDenied, "not allowed to read org")
	}

	prefs, err := s.admin.DB.FindNotificationPreferences(ctx, claims.OwnerID(), org.ID)
	if err != nil {
		return nil, err
	}

	return &adminv1.GetNotificationPreferencesResponse{
		Preferences: notificationPreferencesToPB(prefs),
	}, nil
}

func (s *Server) UpdateNotificationPreferences(ctx context.Context, req *adminv1.UpdateNotificationPreferencesRequest) (*adminv1.UpdateNotificationPreferencesResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.org", req.Org),
	)

	claims := auth.GetClaims(ctx)
	if claims.OwnerType() != auth.OwnerTypeUser {
		return nil, status.Error(codes.Unauthenticated, "not authenticated as a user")
	}

	org, err := s.admin.DB.FindOrganizationByName(ctx, req.Org)
	if err != nil {
		return nil, err
	}
	if !claims.OrganizationPermissions(ctx, org.ID).ReadOrg {
		return nil, status.Error(codes.PermissionDenied, "not allowed to read org")
	}

	prefs, err := s.admin.DB.UpsertNotificationPreferences(ctx, claims.OwnerID(), org.ID, &database.UpsertNotificationPreferencesOptions{
		PushAlerts:       req.Preferences.PushAlerts,
		PushReports:      req.Preferences.PushReports,
		PushActApprovals: req.Preferences.PushActApprovals,
	})
	if err != nil {
		return nil, err
	}

	return &adminv1.UpdateNotificationPreferencesResponse{
		Preferences: notificationPreferencesToPB(prefs),
	}, nil
}

func (s *Server) SendPushNotification(ctx context.Context, req *adminv1.SendPushNotificationRequest) (*adminv1.SendPushNotificationResponse, error) {
	observability.AddRequestAttributes(ctx,
		attribute.String("args.project_id", req.ProjectId),
		attribute.String("args.category", req.Category),
		attribute.StringSlice("args.recipient_emails", req.RecipientEmails),
	)

	proj, err := s.admin.DB.FindProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}

	permissions := auth.GetClaims(ctx).ProjectPermissions(ctx, proj.OrganizationID, proj.ID)
	if !permissions.ReadProdStatus {
		return nil, status.Error(codes.PermissionDenied, "does not have permission to send push notifications")
	}

	org, err := s.admin.DB.FindOrganization(ctx, proj.OrganizationID)
	if err != nil {
		return nil, err
	}

	// The link path must be relative to the frontend URL: reject absolute URLs and protocol-relative ("//host") paths.
	if req.LinkPath != "" {
		u, err := url.Parse(req.LinkPath)
		if err != nil || !strings.HasPrefix(req.LinkPath, "/") || u.Scheme != "" || u.Host != "" {
			return nil, status.Error(codes.InvalidArgument, "link_path must be a relative path starting with \"/\"")
		}
	}
	link := strings.TrimSuffix(s.admin.URLs.WithCustomDomain(org.CustomDomain).Frontend(), "/") + req.LinkPath

	opts := &admin.SendPushNotificationOptions{
		ProjectID:  proj.ID,
		Category:   req.Category,
		Recipients: req.RecipientEmails,
		LinkPath:   req.LinkPath,
		Message: &pushnotifications.Message{
			Title:    req.Title,
			Body:     req.Body,
			Link:     link,
			Category: req.Category,
			Tag:      req.Tag,
		},
	}
	// A caller that counted what is waiting says so with the map; one that has nothing to count (a report is read,
	// not resolved) omits it, and the notification then leaves the app badge alone. The scope is composed here
	// rather than sent: the caller identifies the project by ID, and the names the browser keys on are the ones
	// this service resolved for the link.
	if req.PendingByRecipient != nil {
		// Keyed like the subscriptions are matched, so a caller that spells an address in mixed case still
		// reaches the right person's devices.
		byRecipient := make(map[string]int32, len(req.PendingByRecipient))
		for email, count := range req.PendingByRecipient {
			byRecipient[strings.ToLower(email)] = count
		}
		opts.Pending = &admin.PendingCounts{
			Scope:       strings.ToLower(org.Name + "/" + proj.Name),
			ByRecipient: byRecipient,
		}
	}

	sent, err := s.admin.SendPushNotification(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &adminv1.SendPushNotificationResponse{
		Sent: int32(sent),
	}, nil
}

func notificationPreferencesToPB(prefs *database.NotificationPreferences) *adminv1.NotificationPreferences {
	return &adminv1.NotificationPreferences{
		PushAlerts:       prefs.PushAlerts,
		PushReports:      prefs.PushReports,
		PushActApprovals: prefs.PushActApprovals,
	}
}
