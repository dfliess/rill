package email

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"html/template"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/rilldata/rill/admin/database"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
)

//go:embed templates/gen/*
var templatesFS embed.FS

const dateFormat = "January 2, 2006"

// htmlToText converts HTML to a plain-text approximation suitable for the
// text/plain part of a multipart/alternative email.
func htmlToText(s string) string {
	// Convert block-level closing tags and breaks to newlines
	for _, tag := range []string{"<br>", "<br/>", "<br />", "<BR>", "<BR/>", "<BR />"} {
		s = strings.ReplaceAll(s, tag, "\n")
	}
	for _, tag := range []string{"</p>", "</div>", "</tr>", "</h1>", "</h2>", "</h3>", "</h4>"} {
		s = strings.ReplaceAll(s, tag, "\n")
	}
	s = strings.ReplaceAll(s, "</li>", "\n")
	s = strings.ReplaceAll(s, "<li>", "- ")
	s = strings.ReplaceAll(s, "<LI>", "- ")

	// Strip remaining HTML tags
	s = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(s, "")

	// Decode HTML entities
	s = html.UnescapeString(s)

	// Collapse runs of blank lines
	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")

	// Trim leading/trailing whitespace per line, then overall
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	s = strings.Join(lines, "\n")
	return strings.TrimSpace(s)
}

type Client struct {
	Sender    Sender
	branding  *Branding
	templates *template.Template
}

func New(sender Sender, opts ...Option) *Client {
	c := &Client{
		Sender:   sender,
		branding: DefaultBranding(),
	}
	for _, opt := range opts {
		opt(c)
	}
	templateFuncs := template.FuncMap{
		"now":   time.Now,
		"brand": func() *Branding { return c.branding },
	}
	c.templates = template.Must(template.New("").Funcs(templateFuncs).ParseFS(templatesFS, "templates/gen/*.html"))
	return c
}

type ScheduledReport struct {
	ToEmail         string
	ToName          string
	DisplayName     string
	ReportTime      time.Time
	DownloadFormat  string
	OpenLink        string
	DownloadLink    string
	EditLink        string
	UnsubscribeLink string
	Summary         string // For AI reports
}

type scheduledReportData struct {
	DisplayName      string
	ReportTimeString string // Will be inferred from ReportTime
	DownloadFormat   string
	OpenLink         template.URL
	DownloadLink     template.URL
	EditLink         template.URL
	UnsubscribeLink  template.URL
	Summary          string // For AI reports
}

func (c *Client) SendScheduledReport(opts *ScheduledReport) error {
	// Build template data
	data := &scheduledReportData{
		DisplayName:      opts.DisplayName,
		ReportTimeString: opts.ReportTime.Format(time.RFC1123),
		DownloadFormat:   opts.DownloadFormat,
		OpenLink:         template.URL(opts.OpenLink),
		DownloadLink:     template.URL(opts.DownloadLink),
		EditLink:         template.URL(opts.EditLink),
		UnsubscribeLink:  template.URL(opts.UnsubscribeLink),
		Summary:          opts.Summary,
	}

	// Build subject
	subject := fmt.Sprintf("%s (%s)", opts.DisplayName, data.ReportTimeString)

	var err error
	// Resolve template
	buf := new(bytes.Buffer)
	err = c.templates.Lookup("scheduled_report.html").Execute(buf, data)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()

	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

func (c *Client) SendAlertStatus(opts *drivers.AlertStatus) error {
	switch opts.Status {
	case runtimev1.AssertionStatus_ASSERTION_STATUS_PASS:
		return c.sendAlertStatus(opts, &alertStatusData{
			DisplayName:         opts.DisplayName,
			ExecutionTimeString: opts.ExecutionTime.Format(time.RFC1123),
			IsPass:              true,
			IsRecover:           opts.IsRecover,
			OpenLink:            template.URL(opts.OpenLink),
			EditLink:            template.URL(opts.EditLink),
			UnsubscribeLink:     template.URL(opts.UnsubscribeLink),
		})
	case runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL:
		return c.sendAlertFail(opts, &alertFailData{
			DisplayName:         opts.DisplayName,
			ExecutionTimeString: opts.ExecutionTime.Format(time.RFC1123),
			FailRow:             opts.FailRow,
			OpenLink:            template.URL(opts.OpenLink),
			EditLink:            template.URL(opts.EditLink),
			UnsubscribeLink:     template.URL(opts.UnsubscribeLink),
		})
	case runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR:
		return c.sendAlertStatus(opts, &alertStatusData{
			DisplayName:         opts.DisplayName,
			ExecutionTimeString: opts.ExecutionTime.Format(time.RFC1123),
			IsError:             true,
			ErrorMessage:        opts.ExecutionError,
			OpenLink:            template.URL(opts.EditLink), // NOTE: Using edit link here since for errors, we don't want to open a dashboard, but rather the alert itself
			EditLink:            template.URL(opts.EditLink),
			UnsubscribeLink:     template.URL(opts.UnsubscribeLink),
		})
	default:
		return fmt.Errorf("unknown assertion status: %v", opts.Status)
	}
}

type alertFailData struct {
	DisplayName         string
	ExecutionTimeString string // Will be inferred from ExecutionTime
	FailRow             map[string]any
	OpenLink            template.URL
	EditLink            template.URL
	UnsubscribeLink     template.URL
}

func (c *Client) sendAlertFail(opts *drivers.AlertStatus, data *alertFailData) error {
	subject := fmt.Sprintf("%s (%s)", data.DisplayName, data.ExecutionTimeString)

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("alert_fail.html").Execute(buf, data)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()

	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

type alertStatusData struct {
	DisplayName         string
	ExecutionTimeString string // Will be inferred from ExecutionTime
	IsPass              bool
	IsRecover           bool
	IsError             bool
	ErrorMessage        string
	OpenLink            template.URL
	EditLink            template.URL
	UnsubscribeLink     template.URL
}

func (c *Client) sendAlertStatus(opts *drivers.AlertStatus, data *alertStatusData) error {
	subject := fmt.Sprintf("%s (%s)", data.DisplayName, data.ExecutionTimeString)
	if data.IsRecover {
		subject = fmt.Sprintf("Recovered: %s", subject)
	}

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("alert_status.html").Execute(buf, data)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()

	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

type CallToAction struct {
	ToEmail    string
	ToName     string
	Subject    string
	PreButton  template.HTML
	ButtonText string
	ButtonLink string
	PostButton template.HTML
	ShowFooter bool
}

func (c *Client) SendCallToAction(opts *CallToAction) error {
	buf := new(bytes.Buffer)
	err := c.templates.Lookup("call_to_action.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

type Informational struct {
	ToEmail    string
	ToName     string
	Subject    string
	Body       template.HTML
	ShowFooter bool
}

func (c *Client) SendInformational(opts *Informational) error {
	buf := new(bytes.Buffer)
	err := c.templates.Lookup("informational.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

type Welcome struct {
	ToEmail     string
	ToName      string
	Subject     string
	FrontendURL string
	WelcomeText template.HTML
}

func (c *Client) SendWelcomeToTrial(opts *Welcome) error {
	buf := new(bytes.Buffer)
	err := c.templates.Lookup("welcome_to_trial.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

func (c *Client) SendWelcomeToTeam(opts *Welcome) error {
	buf := new(bytes.Buffer)
	err := c.templates.Lookup("welcome_to_team.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

type OrganizationInvite struct {
	ToEmail       string
	ToName        string
	AcceptURL     string
	OrgName       string
	RoleName      string
	InvitedByName string
}

func (c *Client) SendOrganizationInvite(opts *OrganizationInvite) error {
	if opts.InvitedByName == "" {
		opts.InvitedByName = c.branding.ProductName
	}

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("%s invited you to join %s", opts.InvitedByName, c.branding.ProductName),
		PreButton:  template.HTML(fmt.Sprintf("%s has invited you to join <b>%s</b> as a %s for their %s account. Get started interacting with fast, exploratory dashboards by clicking the button below to sign in and accept your invitation.", opts.InvitedByName, opts.OrgName, opts.RoleName, c.branding.ProductName)),
		ButtonText: "Accept invitation",
		ButtonLink: opts.AcceptURL,
	})
}

type OrganizationAddition struct {
	ToEmail       string
	ToName        string
	OpenURL       string
	OrgName       string
	RoleName      string
	InvitedByName string
}

func (c *Client) SendOrganizationAddition(opts *OrganizationAddition) error {
	if opts.InvitedByName == "" {
		opts.InvitedByName = c.branding.ProductName
	}

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("%s has added you to %s", opts.InvitedByName, opts.OrgName),
		PreButton:  template.HTML(fmt.Sprintf("%s has added you as a %s for <b>%s</b>. Click the button below to view and collaborate on %s dashboard projects for %s.", opts.InvitedByName, opts.RoleName, opts.OrgName, c.branding.ProductName, opts.OrgName)),
		ButtonText: "View account",
		ButtonLink: opts.OpenURL,
	})
}

type ProjectInvite struct {
	ToEmail       string
	ToName        string
	AcceptURL     string
	OrgName       string
	ProjectName   string
	RoleName      string
	InvitedByName string
}

func (c *Client) SendProjectInvite(opts *ProjectInvite) error {
	if opts.InvitedByName == "" {
		opts.InvitedByName = c.branding.ProductName
	}

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("You have been invited to the %s/%s project", opts.OrgName, opts.ProjectName),
		PreButton:  template.HTML(fmt.Sprintf("%s has invited you to collaborate as a %s for the <b>%s/%s</b> project. Click the button below to accept your invitation. ", opts.InvitedByName, opts.RoleName, opts.OrgName, opts.ProjectName)),
		ButtonText: "Accept invitation",
		ButtonLink: opts.AcceptURL,
	})
}

type ProjectAddition struct {
	ToEmail       string
	ToName        string
	OpenURL       string
	OrgName       string
	ProjectName   string
	RoleName      string
	InvitedByName string
}

func (c *Client) SendProjectAddition(opts *ProjectAddition) error {
	if opts.InvitedByName == "" {
		opts.InvitedByName = c.branding.ProductName
	}

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("You have been added to the %s/%s project", opts.OrgName, opts.ProjectName),
		PreButton:  template.HTML(fmt.Sprintf("%s has invited you to collaborate as a %s for the <b>%s</b> project. Click the button below to accept your invitation. ", opts.InvitedByName, opts.RoleName, opts.ProjectName)),
		ButtonText: "View account",
		ButtonLink: opts.OpenURL,
	})
}

type ProjectAccessRequest struct {
	Title       string
	Body        template.HTML
	ToEmail     string
	ToName      string
	Email       string
	OrgName     string
	Role        string
	ProjectName string
	ApproveLink string
	DenyLink    string
}

func (c *Client) SendProjectAccessRequest(opts *ProjectAccessRequest) error {
	var accessPrefix string
	switch opts.Role {
	case database.ProjectRoleNameAdmin:
		accessPrefix = "to be an admin of"
	case database.ProjectRoleNameEditor:
		accessPrefix = "to edit"
	case database.ProjectRoleNameViewer:
		accessPrefix = "to view"
	}

	subject := fmt.Sprintf("%s would like %s %s/%s", opts.Email, accessPrefix, opts.OrgName, opts.ProjectName)
	if opts.Body == "" {
		opts.Body = template.HTML(fmt.Sprintf("<b>%s</b> would like %s <b>%s/%s</b>", opts.Email, accessPrefix, opts.OrgName, opts.ProjectName))
	}

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("project_access_request.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	html := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    html,
		Text:    htmlToText(html),
	})
}

type ProjectAccessGranted struct {
	ToEmail     string
	ToName      string
	OpenURL     string
	OrgName     string
	ProjectName string
}

func (c *Client) SendProjectAccessGranted(opts *ProjectAccessGranted) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("Your request to %s/%s has been approved", opts.OrgName, opts.ProjectName),
		PreButton:  template.HTML(fmt.Sprintf("Your request to <b>%s/%s</b> has been approved", opts.OrgName, opts.ProjectName)),
		ButtonText: fmt.Sprintf("View project in %s", c.branding.ProductName),
		ButtonLink: opts.OpenURL,
	})
}

type ProjectAccessRejected struct {
	ToEmail     string
	ToName      string
	OrgName     string
	ProjectName string
}

func (c *Client) SendProjectAccessRejected(opts *ProjectAccessRejected) error {
	return c.SendInformational(&Informational{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Your request to %s/%s has been denied", opts.OrgName, opts.ProjectName),
		Body:    template.HTML(fmt.Sprintf("Your request to <b>%s/%s</b> has been denied. Contact your project admin for help.", opts.OrgName, opts.ProjectName)),
	})
}

type InvoicePaymentFailed struct {
	ToEmail            string
	ToName             string
	OrgName            string
	Currency           string
	Amount             string
	PaymentURL         string
	GracePeriodEndDate time.Time
}

func (c *Client) SendInvoicePaymentFailed(opts *InvoicePaymentFailed) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Payment failed for %s. Please update your payment method", opts.OrgName),
		PreButton: template.HTML(fmt.Sprintf(`
We couldn’t process your payment for <b>%s</b>. You have until <b>%s</b> to update your payment details before your org is hibernated.
`, opts.OrgName, opts.GracePeriodEndDate.Format(dateFormat))),
		ButtonText: "Update Payment Info",
		ButtonLink: opts.PaymentURL,
		ShowFooter: true,
	})
}

type InvoicePaymentSuccess struct {
	ToEmail        string
	ToName         string
	OrgName        string
	PaymentDate    time.Time
	BillingPageURL string
}

// SendInvoicePaymentSuccess Currently Used only when a previously failed invoice payment succeeds
func (c *Client) SendInvoicePaymentSuccess(opts *InvoicePaymentSuccess) error {
	return c.SendInformational(&Informational{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Successful payment %s", opts.PaymentDate.Format(dateFormat)),
		Body: template.HTML(fmt.Sprintf(`
Thank you for your payment!
<br /><br />
Your payment for <b>%s</b> has been successfully processed.
<br /><br />
If you believe this charge to be in error or have any questions, please email %s.
<br /><br />
You can manage your subscription by visiting the <a href=%q>Billing settings</a>
`, opts.OrgName, c.branding.SupportEmail, opts.BillingPageURL)),
		ShowFooter: false,
	})
}

type InvoiceUnpaid struct {
	ToEmail    string
	ToName     string
	OrgName    string
	PaymentURL string
}

// SendInvoiceUnpaid sent after the payment grace period has ended
func (c *Client) SendInvoiceUnpaid(opts *InvoiceUnpaid) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Invoice for %s is now past due. Org is now hibernated", opts.OrgName),
		PreButton: template.HTML(fmt.Sprintf(`
<b>%s</b> and its projects have been hibernated due to an overdue payment. 
<br /><br />
Restore access by updating your payment information today! 
`, opts.ToName)),
		ButtonText: "Update Payment Info",
		ButtonLink: opts.PaymentURL,
		ShowFooter: true,
	})
}

type SubscriptionCancelled struct {
	ToEmail    string
	ToName     string
	OrgName    string
	PlanName   string
	BillingURL string
	EndDate    time.Time
}

func (c *Client) SendSubscriptionCancelled(opts *SubscriptionCancelled) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("%s for %s is canceled. Access available until %s", opts.PlanName, opts.OrgName, opts.EndDate.Format(dateFormat)),
		PreButton: template.HTML(fmt.Sprintf(`
We’re sorry to see you go!
<br /><br />
You’ve successfully canceled the %s plan for <b>%s</b>. You’ll still have access to %s until <b>%s</b>. After this date, your subscription will expire, and you will no longer have access.
<br /><br />
If you change your mind, you can always reactivate your subscription!
`, opts.PlanName, opts.ToName, c.branding.ProductName+" Cloud", opts.EndDate.Format(dateFormat))),
		ButtonText: "Billing Settings",
		ButtonLink: opts.BillingURL,
		PostButton: template.HTML(fmt.Sprintf(`If you found that our service did not meet your needs, please contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a> and we’ll do our best to address your feedback and concerns.`, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName)),
		ShowFooter: false,
	})
}

type SubscriptionEnded struct {
	ToEmail    string
	ToName     string
	OrgName    string
	BillingURL string
}

func (c *Client) SendSubscriptionEnded(opts *SubscriptionEnded) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Subscription for %s has now ended. Org is hibernated", opts.OrgName),
		PreButton: template.HTML(fmt.Sprintf(`
Your cancelled subscription for <b>%s</b> has ended and its projects are now <a href="%s/developers/other/FAQ#what-is-project-hibernation">hibernating</a>. We hope you enjoyed using %s during your time with us.
<br /><br />
If you’d like to reactivate your subscription and regain access, you can easily do so at any time by renewing your subscription from here:
`, opts.OrgName, c.branding.DocsURL, c.branding.ProductName+" Cloud")),
		ButtonText: "Billing Settings",
		ButtonLink: opts.BillingURL,
		PostButton: template.HTML(fmt.Sprintf(`
If you have any feedback about your experience or how we can improve, please feel free to contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a>
<br /><br />
Thank you for trying %s. We hope to see you again in the future!
`, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName, c.branding.ProductName+" Cloud")),
		ShowFooter: false,
	})
}

type TrialStarted struct {
	ToEmail      string
	ToName       string
	OrgName      string
	FrontendURL  string
	TrialEndDate time.Time
}

func (c *Client) SendTrialStarted(opts *TrialStarted) error {
	return c.SendWelcomeToTrial(&Welcome{
		ToEmail:     opts.ToEmail,
		ToName:      opts.ToName,
		Subject:     fmt.Sprintf("A 30-day free trial for %s has started", opts.OrgName),
		FrontendURL: opts.FrontendURL,
		WelcomeText: template.HTML(fmt.Sprintf(`
You now have access to %s until <b>%s</b> to explore all features including:
<ul>
<li>User management (RBAC)</li>
<li>Embedded dashboards</li>
<li>Alerts and scheduled reports</li>
</ul>
`, c.branding.ProductName+" Cloud", opts.TrialEndDate.Format(dateFormat))),
	})
}

type TrialEndingSoon struct {
	ToEmail      string
	ToName       string
	OrgName      string
	UpgradeURL   string
	TrialEndDate time.Time
}

func (c *Client) SendTrialEndingSoon(opts *TrialEndingSoon) error {
	diff := time.Until(opts.TrialEndDate)
	days := int(math.Round(diff.Hours() / 24))
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Your %s trial for %s is expiring in %d days", c.branding.ProductName+" Cloud", opts.OrgName, days),
		PreButton: template.HTML(fmt.Sprintf(`
Your trial for <b>%s</b> ends on <b>%s</b>.
<br /><br />
How's %s working out for you? Have you checked out our newest features highlighted in our <a href="%s">Release Notes</a>?
<br /><br />
Our team is here to help you in any way we can, so don't hesitate to contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a> if you have a question, encounter an issue, or need guidance.
<br /><br />
If you're ready to upgrade, simply click the button below.
`, opts.ToName, opts.TrialEndDate.Format(dateFormat), c.branding.ProductName, c.branding.ReleaseNotesURL, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName)),
		ButtonText: "Upgrade Now",
		ButtonLink: opts.UpgradeURL,
		ShowFooter: false,
	})
}

type TrialEnded struct {
	ToEmail            string
	ToName             string
	OrgName            string
	UpgradeURL         string
	GracePeriodEndDate time.Time
}

func (c *Client) SendTrialEnded(opts *TrialEnded) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Your %s trial for %s has expired", c.branding.ProductName+" Cloud", opts.OrgName),
		PreButton: template.HTML(fmt.Sprintf(`
Hi %s,
<br /><br />
Your %s trial has now expired. <b>%s</b> will be hibernated on <b>%s</b>. We hope you’ve enjoyed using our software. If you’d like to keep using %s, upgrade to our Team Plan!
`, opts.ToName, c.branding.ProductName+" Cloud", opts.OrgName, opts.GracePeriodEndDate.Format(dateFormat), c.branding.ProductName+" Cloud")),
		ButtonText: "Upgrade to Team Plan",
		ButtonLink: opts.UpgradeURL,
		ShowFooter: true,
	})
}

type TrialGracePeriodEnded struct {
	ToEmail    string
	ToName     string
	OrgName    string
	UpgradeURL string
}

func (c *Client) SendTrialGracePeriodEnded(opts *TrialGracePeriodEnded) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Trial plan grace period for %s has ended. Org is now hibernated", opts.OrgName),
		PreButton: template.HTML(fmt.Sprintf(`
<b>%s</b> and its projects are now <a href="%s/developers/other/FAQ#what-is-project-hibernation">hibernating</a>.
<br /><br />
Reactivate your org by upgrading to the Team Plan today!
`, opts.OrgName, c.branding.DocsURL)),
		ButtonText: "Upgrade to Team Plan",
		ButtonLink: opts.UpgradeURL,
		PostButton: template.HTML(fmt.Sprintf(`
We'd love to hear from you! If you have any feedback about your experience or how we can improve, please feel free to contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a>
<br /><br />
Thank you for trying %s. We hope to see you again in the future!
`, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName, c.branding.ProductName+" Cloud")),
		ShowFooter: false,
	})
}

type CreditTrialStarted struct {
	ToEmail          string
	ToName           string
	OrgName          string
	FrontendURL      string
	CreditAllocation int
}

func (c *Client) SendCreditTrialStarted(opts *CreditTrialStarted) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Welcome to %s — start using your $%d credit", c.branding.ProductName, opts.CreditAllocation),
		PreButton: template.HTML(fmt.Sprintf(`
Hi there,
<br /><br />
Welcome to %s! Your account is ready and loaded with <b>$%d in free credit</b> to explore the full platform — dashboards, metrics, embedded analytics, AI-powered exploration, and more.
<br /><br />
Your $%d credit covers $0.15/compute unit/hr and $1/GB storage/mo for managed data above 1GB. There's no time limit — your credit is only consumed when your dashboards are running. A typical project with 4 compute units gives you a few weeks to build and share real dashboards with your team. Don't forget to hibernate your project when you're not using it!
`, c.branding.ProductName, opts.CreditAllocation, opts.CreditAllocation)),
		ButtonText: fmt.Sprintf("Open %s", c.branding.ProductName+" Cloud"),
		ButtonLink: opts.FrontendURL,
		PostButton: template.HTML(fmt.Sprintf(`
If you have any questions, feel free to contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a> You can also check out <a href="%s" style="color:%s">%s</a> to learn more.
<br /><br />
Happy exploring,
`, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName, c.branding.DocsURL, c.branding.PrimaryColor, c.branding.DocsURL)),
		ShowFooter: false,
	})
}

type CreditTrialLow struct {
	ToEmail          string
	ToName           string
	OrgName          string
	FrontendURL      string
	UpgradeURL       string
	CreditAllocation int
	RemainingBalance float64
}

func (c *Client) SendCreditTrialLow(opts *CreditTrialLow) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Your %s credit is getting low", c.branding.ProductName),
		PreButton: template.HTML(fmt.Sprintf(`
Hi there,
<br /><br />
You have about <b>$%.2f</b> of your original $%d %s credit remaining. Once it's used up, your dashboards will go into hibernation.
<br /><br />
Hibernation means your dashboards pause and go offline — but nothing is deleted. Your models, data connections, and configuration all stay intact. Upgrading to Pro reactivates everything.
<br /><br />
The Pro plan is simple: $0.15/compute unit/hr, $1/GB storage/mo for managed data above 1GB, no monthly minimums, and any remaining free credit carries over so nothing goes to waste.
`, opts.RemainingBalance, opts.CreditAllocation, c.branding.ProductName)),
		ButtonText: "Upgrade to Pro",
		ButtonLink: opts.UpgradeURL,
		PostButton: template.HTML(fmt.Sprintf(`
If you have any questions, feel free to contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a>
<br /><br />
Happy building,
`, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName)),
		ShowFooter: false,
	})
}

type CreditTrialDepleted struct {
	ToEmail          string
	ToName           string
	OrgName          string
	FrontendURL      string
	UpgradeURL       string
	CreditAllocation int
}

func (c *Client) SendCreditTrialDepleted(opts *CreditTrialDepleted) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("Your %s dashboards are now hibernated", c.branding.ProductName),
		PreButton: template.HTML(fmt.Sprintf(`
Hi there,
<br /><br />
Your $%d %s credit has been fully used and your dashboards are now <a href="%s/developers/other/FAQ#what-is-project-hibernation" style="color:%s">hibernated</a>. All deployments are paused — your team and any embedded analytics are currently offline.
<br /><br />
<b>Nothing has been deleted.</b> Your entire project — data connections, models, dashboards, and configuration — is preserved exactly as you left it.
<br /><br />
To bring everything back online, upgrade to the Pro plan. It takes under a minute: add a payment method, confirm the upgrade, and your dashboards reactivate with your existing configuration.
`, opts.CreditAllocation, c.branding.ProductName, c.branding.DocsURL, c.branding.PrimaryColor)),
		ButtonText: "Upgrade to Pro",
		ButtonLink: opts.UpgradeURL,
		PostButton: template.HTML(fmt.Sprintf(`
Pro pricing is straightforward: $0.15/compute unit/hr, $1/GB storage/mo for managed data above 1GB. No contracts, no seat fees. Scale up or down anytime.
<br /><br />
If you'd like to discuss your options or need a custom arrangement, feel free to contact us via <a href="mailto:%s" style="color:%s">email</a>, or via chat on <a href="%s" style="color:%s">%s.</a>
<br /><br />
Hope to see you back soon,
`, c.branding.SupportEmail, c.branding.PrimaryColor, c.branding.ChatURL, c.branding.PrimaryColor, c.branding.ProductName)),
		ShowFooter: false,
	})
}

type PlanUpdate struct {
	ToEmail  string
	ToName   string
	OrgName  string
	PlanName string
}

func (c *Client) SendPlanUpdate(opts *PlanUpdate) error {
	return c.SendInformational(&Informational{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("Your plan for %s has been updated to %s plan", opts.OrgName, opts.PlanName),
		Body:       template.HTML(fmt.Sprintf("<b>%q</b> has been updated to %q plan.", opts.OrgName, opts.PlanName)),
		ShowFooter: true,
	})
}

type SubscriptionRenewed struct {
	ToEmail  string
	ToName   string
	OrgName  string
	PlanName string
}

func (c *Client) SendSubscriptionRenewed(opts *SubscriptionRenewed) error {
	return c.SendInformational(&Informational{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Subject:    fmt.Sprintf("Your %s subscription for %s plan has been renewed", opts.PlanName, opts.OrgName),
		Body:       template.HTML(fmt.Sprintf("Your subscription for <b>%q</b> has been renewed for %q plan.", opts.OrgName, opts.PlanName)),
		ShowFooter: true,
	})
}

type PaidPlan struct {
	ToEmail          string
	ToName           string
	OrgName          string
	FrontendURL      string
	BillingURL       string
	PlanName         string
	BillingStartDate time.Time
}

// SendPaidPlanStarted sends a customised plan-started email for a paid plan (Team or Pro).
func (c *Client) SendPaidPlanStarted(opts *PaidPlan) error {
	return c.SendCallToAction(&CallToAction{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: fmt.Sprintf("You're on the %s plan", opts.PlanName),
		PreButton: template.HTML(fmt.Sprintf(`
Hi there,
<br /><br />
You're all set on the %s %s plan. Your next billing cycle starts on <b>%s</b>.
<br /><br />
Billing is usage-based with no contracts — you'll receive a monthly invoice with a full breakdown of compute hours and data storage.
`, c.branding.ProductName, opts.PlanName, opts.BillingStartDate.Format(dateFormat))),
		ButtonText: "View Your Billing Dashboard",
		ButtonLink: opts.BillingURL,
		PostButton: template.HTML(`
Welcome aboard. Happy building,
`),
		ShowFooter: false,
	})
}

// SendPaidPlanRenewal sends a customised plan-renewed email for a paid plan (Team or Pro).
func (c *Client) SendPaidPlanRenewal(opts *PaidPlan) error {
	return c.SendWelcomeToTeam(&Welcome{
		ToEmail:     opts.ToEmail,
		ToName:      opts.ToName,
		Subject:     fmt.Sprintf("Your %s plan subscription for %s has been renewed", opts.PlanName, opts.OrgName),
		FrontendURL: opts.FrontendURL,
		WelcomeText: template.HTML(fmt.Sprintf(`
Thank you! You’ve successfully renewed %s to the <b>%s</b> plan.
<br /><br />
Your next billing cycle starts on %s.
`, opts.OrgName, opts.PlanName, opts.BillingStartDate.Format(dateFormat))),
	})
}
