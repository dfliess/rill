package drivers

import (
	"context"
	"errors"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
)

var ErrNotAuthenticated = errors.New("not authenticated")

type AdminService interface {
	GetReportMetadata(ctx context.Context, reportName, ownerID, webOpenMode string, emailRecipients []string, anonRecipients bool, executionTime time.Time) (*ReportMetadata, error)
	GetAlertMetadata(ctx context.Context, alertName, ownerID string, emailRecipients []string, anonRecipients bool, annotations map[string]string, queryForUserID, queryForUserEmail string) (*AlertMetadata, error)
	// SendPushNotification sends a web push notification to the subscriptions of its recipients.
	// Recipients that opted out of its category are skipped.
	// It returns the number of notifications sent, which may be zero if push is disabled or nobody is subscribed.
	SendPushNotification(ctx context.Context, n *PushNotification) (int, error)
	// ListProjectMemberAttributes returns the project's members with the identity a runtime JWT issued for them would carry.
	// It lets the runtime evaluate a security policy for members that are not making a request,
	// which is what turns "may this user decide?" into "who may decide?" (see Runtime.ResolveAgentApprovers).
	ListProjectMemberAttributes(ctx context.Context) ([]ProjectMember, error)
	ProvisionConnector(ctx context.Context, name, driver string, args map[string]any) (map[string]any, error)
	GetConfig(ctx context.Context) (*Config, error)
	ListDeployments(ctx context.Context) ([]*Deployment, error)
	UpdateProjectVariables(ctx context.Context, environment string, variables map[string]string) error
}

// PushNotification is one web push dispatch: what it says, who it goes to, and what it leaves behind on their
// app icon.
type PushNotification struct {
	// Category must be one of the PushCategory constants. It is what recipients opt out of.
	Category   string
	Recipients []string
	Title      string
	Body       string
	// LinkPath is the frontend path to open when the notification is clicked. It must start with "/".
	LinkPath string
	// Tag makes notifications replace each other in the browser: a second one with the same tag takes the place
	// of the first instead of stacking under it.
	Tag string
	// PendingByRecipient, when non-nil, is how many things are left waiting for each recipient in this project
	// once this notification lands, keyed by email. It is what the installed app paints on its icon
	// (kairos-cloud#147), so it counts work to resolve, not notifications delivered: a category with nothing to
	// resolve leaves it nil and the badge untouched. A recipient absent from a non-nil map has nothing pending,
	// which the badge needs to hear just as much as a positive count.
	PendingByRecipient map[string]int32
}

// Categories of web push notifications, matching the notification preferences users can toggle in the admin.
const (
	PushCategoryAlerts       = "alerts"
	PushCategoryReports      = "reports"
	PushCategoryActApprovals = "act_approvals"
)

// ProjectMember is one member of the project as the runtime sees them:
// the attributes, permissions and security rules a runtime JWT issued for them would carry.
type ProjectMember struct {
	UserID string
	Email  string
	// Attributes as they would be embedded in the member's runtime JWT, i.e. the `.user` a security policy evaluates against.
	Attributes map[string]any
	// EditTrigger reports whether the member holds the runtime's EditTrigger permission on this deployment.
	// It is the only instance permission that differs between members, see Runtime.ResolveAgentApprovers.
	EditTrigger bool
	// SecurityRules are the resource restrictions that apply to the member; empty for an unrestricted one.
	SecurityRules []*runtimev1.SecurityRule
}

type ReportMetadata struct {
	ReportDelivery map[string]ReportDelivery
}

type ReportDelivery struct {
	OpenURL        string
	ExportURL      string
	EditURL        string
	UnsubscribeURL string
	UserID         string         // user ID of the intended recipient, will be empty for non-Rill users and users not having project access. In creator mode this will be the user ID of the creator.
	UserAttrs      map[string]any // user attrs of the intended recipient, will be empty for non-Rill users and users not having project access. In creator mode this will be the user attrs of the creator.
}

type AlertURLs struct {
	OpenURL        string
	EditURL        string
	UnsubscribeURL string
}

type AlertMetadata struct {
	RecipientURLs      map[string]AlertURLs
	QueryForAttributes map[string]any
}

// Config holds configuration returned by the admin service for the runtime instance.
// For local runtimes only applicable fields like Variables will be populated.
type Config struct {
	Variables             map[string]map[string]string
	SystemVariables       map[string]string
	Annotations           map[string]string
	FrontendURL           string
	UpdatedOn             time.Time
	UsesArchive           bool
	DuckdbConnectorConfig map[string]any
	Editable              bool
}

type Deployment struct {
	Branch   string
	Editable bool
}
