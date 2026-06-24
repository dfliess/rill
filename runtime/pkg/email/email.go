package email

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rilldata/rill/admin/database"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"golang.org/x/text/language"
)

//go:embed templates/gen/*
var templatesFS embed.FS

//go:embed catalog/*.toml
var catalogFS embed.FS

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

// Client sends transactional emails with i18n support.
type Client struct {
	Sender        Sender
	branding      *Branding
	bundle        *i18n.Bundle
	defaultLocale string
	templates     *template.Template
}

// New creates an email Client. Options configure branding and default locale.
func New(sender Sender, opts ...Option) *Client {
	c := &Client{
		Sender:        sender,
		branding:      DefaultBranding(),
		defaultLocale: "en",
	}
	for _, opt := range opts {
		opt(c)
	}

	// Init i18n bundle
	c.bundle = i18n.NewBundle(language.English)
	c.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// Load all catalog files; ignore errors for missing locales.
	entries, _ := fs.ReadDir(catalogFS, "catalog")
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
			_, _ = c.bundle.LoadMessageFileFS(catalogFS, "catalog/"+e.Name())
		}
	}

	templateFuncs := template.FuncMap{
		"now":   time.Now,
		"brand": func() *Branding { return c.branding },
	}
	c.templates = template.Must(template.New("").Funcs(templateFuncs).ParseFS(templatesFS, "templates/gen/*.html"))
	return c
}

// localizer returns an i18n.Localizer for the given locale with fallback to
// the client's default locale.
func (c *Client) localizer(locale string) *i18n.Localizer {
	if locale == "" {
		locale = c.defaultLocale
	}
	return i18n.NewLocalizer(c.bundle, locale, c.defaultLocale)
}

// T returns a localized, HTML-escaped string. If the message ID is not found,
// the messageID itself is returned as a fallback.
func (c *Client) T(locale, messageID string, data map[string]any) string {
	loc := c.localizer(locale)
	s, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return s
}

// THtml returns a localized string as template.HTML. The caller must audit
// catalog entries used with THtml for XSS safety — only use for messages that
// contain trusted HTML markup.
func (c *Client) THtml(locale, messageID string, data map[string]any) template.HTML {
	return template.HTML(c.T(locale, messageID, data))
}

// ---------------------------------------------------------------------------
// brandData returns a map of common branding fields for use in go-i18n
// template data.
// ---------------------------------------------------------------------------
func (c *Client) brandData() map[string]any {
	return map[string]any{
		"ProductName":      c.branding.ProductName,
		"ProductNameCloud": c.branding.ProductName + " Cloud",
		"SupportEmail":     c.branding.SupportEmail,
		"PrimaryColor":     c.branding.PrimaryColor,
		"ChatURL":          c.branding.ChatURL,
		"DocsURL":          c.branding.DocsURL,
		"ReleaseNotesURL":  c.branding.ReleaseNotesURL,
	}
}

// mergeData merges the branding base map with additional key-value pairs.
func (c *Client) mergeData(extra map[string]any) map[string]any {
	m := c.brandData()
	for k, v := range extra {
		m[k] = v
	}
	return m
}

// ---------------------------------------------------------------------------
// Scheduled reports
// ---------------------------------------------------------------------------

// ScheduledReport holds the parameters for a scheduled-report email.
type ScheduledReport struct {
	ToEmail         string
	ToName          string
	Locale          string
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
	ReportTimeString string
	DownloadFormat   string
	OpenLink         template.URL
	DownloadLink     template.URL
	EditLink         template.URL
	UnsubscribeLink  template.URL
	Summary          string
	// Localized fields
	ReadyText          template.HTML
	OpenButtonText     string
	DownloadButtonText string
	EditLinkText       template.HTML
	UnsubLinkText      template.HTML
}

func (c *Client) SendScheduledReport(opts *ScheduledReport) error {
	reportTimeStr := opts.ReportTime.Format(time.RFC1123)
	locale := opts.Locale

	d := c.mergeData(map[string]any{
		"DisplayName":      opts.DisplayName,
		"ReportTimeString": reportTimeStr,
		"DownloadFormat":   opts.DownloadFormat,
		"EditLink":         opts.EditLink,
		"UnsubscribeLink":  opts.UnsubscribeLink,
	})

	openButton := c.T(locale, "email.button.open_in_browser", nil)
	if opts.Summary != "" {
		openButton = c.T(locale, "email.button.view_full_analysis", nil)
	}

	data := &scheduledReportData{
		DisplayName:        opts.DisplayName,
		ReportTimeString:   reportTimeStr,
		DownloadFormat:     opts.DownloadFormat,
		OpenLink:           template.URL(opts.OpenLink),
		DownloadLink:       template.URL(opts.DownloadLink),
		EditLink:           template.URL(opts.EditLink),
		UnsubscribeLink:    template.URL(opts.UnsubscribeLink),
		Summary:            opts.Summary,
		ReadyText:          c.THtml(locale, "email.body.scheduled_report_ready", d),
		OpenButtonText:     openButton,
		DownloadButtonText: c.T(locale, "email.button.download_format_file", d),
		EditLinkText:       c.THtml(locale, "email.label.edit_report", d),
		UnsubLinkText:      c.THtml(locale, "email.label.unsubscribe_report", d),
	}

	subject := c.T(locale, "email.subject.scheduled_report", d)

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("scheduled_report.html").Execute(buf, data)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	htmlStr := buf.String()

	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

// ---------------------------------------------------------------------------
// Alerts
// ---------------------------------------------------------------------------

func (c *Client) SendAlertStatus(opts *drivers.AlertStatus) error {
	locale := opts.Locale
	if locale == "" {
		locale = c.defaultLocale
	}

	switch opts.Status {
	case runtimev1.AssertionStatus_ASSERTION_STATUS_PASS:
		return c.sendAlertStatus(opts, locale, &alertStatusData{
			DisplayName:         opts.DisplayName,
			ExecutionTimeString: opts.ExecutionTime.Format(time.RFC1123),
			IsPass:              true,
			IsRecover:           opts.IsRecover,
			OpenLink:            template.URL(opts.OpenLink),
			EditLink:            template.URL(opts.EditLink),
			UnsubscribeLink:     template.URL(opts.UnsubscribeLink),
		})
	case runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL:
		return c.sendAlertFail(opts, locale, &alertFailData{
			DisplayName:         opts.DisplayName,
			ExecutionTimeString: opts.ExecutionTime.Format(time.RFC1123),
			FailRow:             opts.FailRow,
			OpenLink:            template.URL(opts.OpenLink),
			EditLink:            template.URL(opts.EditLink),
			UnsubscribeLink:     template.URL(opts.UnsubscribeLink),
		})
	case runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR:
		return c.sendAlertStatus(opts, locale, &alertStatusData{
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
	ExecutionTimeString string
	FailRow             map[string]any
	OpenLink            template.URL
	EditLink            template.URL
	UnsubscribeLink     template.URL
	// Localized fields
	TriggerText    template.HTML
	OpenButtonText string
	EditLinkText   template.HTML
	UnsubLinkText  template.HTML
}

func (c *Client) sendAlertFail(opts *drivers.AlertStatus, locale string, data *alertFailData) error {
	d := c.mergeData(map[string]any{
		"DisplayName":         data.DisplayName,
		"ExecutionTimeString": data.ExecutionTimeString,
		"EditLink":            string(data.EditLink),
		"UnsubscribeLink":     string(data.UnsubscribeLink),
	})

	data.TriggerText = c.THtml(locale, "email.body.alert_fail_trigger", d)
	data.OpenButtonText = c.T(locale, "email.button.open_in_browser", nil)
	data.EditLinkText = c.THtml(locale, "email.label.edit_alert", d)
	data.UnsubLinkText = c.THtml(locale, "email.label.unsubscribe_alert", d)

	subject := c.T(locale, "email.subject.alert_fail", d)

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("alert_fail.html").Execute(buf, data)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	htmlStr := buf.String()

	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

type alertStatusData struct {
	DisplayName         string
	ExecutionTimeString string
	IsPass              bool
	IsRecover           bool
	IsError             bool
	ErrorMessage        string
	OpenLink            template.URL
	EditLink            template.URL
	UnsubscribeLink     template.URL
	// Localized fields
	StatusText     template.HTML
	OpenButtonText string
	EditLinkText   template.HTML
	UnsubLinkText  template.HTML
}

func (c *Client) sendAlertStatus(opts *drivers.AlertStatus, locale string, data *alertStatusData) error {
	d := c.mergeData(map[string]any{
		"DisplayName":         data.DisplayName,
		"ExecutionTimeString": data.ExecutionTimeString,
		"ErrorMessage":        data.ErrorMessage,
		"EditLink":            string(data.EditLink),
		"UnsubscribeLink":     string(data.UnsubscribeLink),
	})

	switch {
	case data.IsError:
		data.StatusText = c.THtml(locale, "email.body.alert_error", d)
	case data.IsRecover:
		data.StatusText = c.THtml(locale, "email.body.alert_recovered", d)
	case data.IsPass:
		data.StatusText = c.THtml(locale, "email.body.alert_passed", d)
	}
	data.OpenButtonText = c.T(locale, "email.button.open_in_browser", nil)
	data.EditLinkText = c.THtml(locale, "email.label.edit_alert", d)
	data.UnsubLinkText = c.THtml(locale, "email.label.unsubscribe_alert", d)

	subject := c.T(locale, "email.subject.alert_status", d)
	if data.IsRecover {
		subject = c.T(locale, "email.subject.alert_recovered", d)
	}

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("alert_status.html").Execute(buf, data)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	htmlStr := buf.String()

	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

// ---------------------------------------------------------------------------
// Call-to-action (generic)
// ---------------------------------------------------------------------------

// CallToAction is a generic CTA email with a single button.
type CallToAction struct {
	ToEmail    string
	ToName     string
	Locale     string
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
	htmlStr := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

// ---------------------------------------------------------------------------
// Informational (generic)
// ---------------------------------------------------------------------------

// Informational is a body-only email with no button.
type Informational struct {
	ToEmail    string
	ToName     string
	Locale     string
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
	htmlStr := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

// ---------------------------------------------------------------------------
// Welcome emails
// ---------------------------------------------------------------------------

// Welcome is used by SendWelcomeToTrial and SendWelcomeToTeam.
type Welcome struct {
	ToEmail     string
	ToName      string
	Locale      string
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
	htmlStr := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

func (c *Client) SendWelcomeToTeam(opts *Welcome) error {
	buf := new(bytes.Buffer)
	err := c.templates.Lookup("welcome_to_team.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	htmlStr := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: opts.Subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

// ---------------------------------------------------------------------------
// Organization invite
// ---------------------------------------------------------------------------

// OrganizationInvite holds the parameters for an organization-invite email.
type OrganizationInvite struct {
	ToEmail       string
	ToName        string
	Locale        string
	AcceptURL     string
	OrgName       string
	RoleName      string
	InvitedByName string
}

func (c *Client) SendOrganizationInvite(opts *OrganizationInvite) error {
	if opts.InvitedByName == "" {
		opts.InvitedByName = c.branding.ProductName
	}
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"InvitedByName": opts.InvitedByName,
		"OrgName":       opts.OrgName,
		"RoleName":      opts.RoleName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.org_invite", d),
		PreButton:  c.THtml(locale, "email.body.org_invite", d),
		ButtonText: c.T(locale, "email.button.accept_invitation", nil),
		ButtonLink: opts.AcceptURL,
	})
}

// ---------------------------------------------------------------------------
// Organization addition
// ---------------------------------------------------------------------------

// OrganizationAddition holds the parameters for an organization-addition email.
type OrganizationAddition struct {
	ToEmail       string
	ToName        string
	Locale        string
	OpenURL       string
	OrgName       string
	RoleName      string
	InvitedByName string
}

func (c *Client) SendOrganizationAddition(opts *OrganizationAddition) error {
	if opts.InvitedByName == "" {
		opts.InvitedByName = c.branding.ProductName
	}
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"InvitedByName": opts.InvitedByName,
		"OrgName":       opts.OrgName,
		"RoleName":      opts.RoleName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.org_addition", d),
		PreButton:  c.THtml(locale, "email.body.org_addition", d),
		ButtonText: c.T(locale, "email.button.view_account", nil),
		ButtonLink: opts.OpenURL,
	})
}

// ---------------------------------------------------------------------------
// Project invite
// ---------------------------------------------------------------------------

// ProjectInvite holds the parameters for a project-invite email.
type ProjectInvite struct {
	ToEmail       string
	ToName        string
	Locale        string
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
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"InvitedByName": opts.InvitedByName,
		"OrgName":       opts.OrgName,
		"ProjectName":   opts.ProjectName,
		"RoleName":      opts.RoleName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.project_invite", d),
		PreButton:  c.THtml(locale, "email.body.project_invite", d),
		ButtonText: c.T(locale, "email.button.accept_invitation", nil),
		ButtonLink: opts.AcceptURL,
	})
}

// ---------------------------------------------------------------------------
// Project addition
// ---------------------------------------------------------------------------

// ProjectAddition holds the parameters for a project-addition email.
type ProjectAddition struct {
	ToEmail       string
	ToName        string
	Locale        string
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
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"InvitedByName": opts.InvitedByName,
		"OrgName":       opts.OrgName,
		"ProjectName":   opts.ProjectName,
		"RoleName":      opts.RoleName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.project_addition", d),
		PreButton:  c.THtml(locale, "email.body.project_addition", d),
		ButtonText: c.T(locale, "email.button.view_account", nil),
		ButtonLink: opts.OpenURL,
	})
}

// ---------------------------------------------------------------------------
// Project access request
// ---------------------------------------------------------------------------

// ProjectAccessRequest holds the parameters for a project-access-request email.
type ProjectAccessRequest struct {
	Title       string
	Body        template.HTML
	ToEmail     string
	ToName      string
	Locale      string
	Email       string
	OrgName     string
	Role        string
	ProjectName string
	ApproveLink string
	DenyLink    string
	// Localized fields set internally
	ApproveButtonText string
	DenyButtonText    string
}

func (c *Client) SendProjectAccessRequest(opts *ProjectAccessRequest) error {
	locale := opts.Locale

	var accessPrefix string
	switch opts.Role {
	case database.ProjectRoleNameAdmin:
		accessPrefix = c.T(locale, "email.label.access_prefix_admin", nil)
	case database.ProjectRoleNameEditor:
		accessPrefix = c.T(locale, "email.label.access_prefix_editor", nil)
	case database.ProjectRoleNameViewer:
		accessPrefix = c.T(locale, "email.label.access_prefix_viewer", nil)
	}

	d := c.mergeData(map[string]any{
		"Email":        opts.Email,
		"OrgName":      opts.OrgName,
		"ProjectName":  opts.ProjectName,
		"AccessPrefix": accessPrefix,
	})

	var subjectKey string
	switch opts.Role {
	case database.ProjectRoleNameAdmin:
		subjectKey = "email.subject.project_access_request_admin"
	case database.ProjectRoleNameEditor:
		subjectKey = "email.subject.project_access_request_editor"
	case database.ProjectRoleNameViewer:
		subjectKey = "email.subject.project_access_request_viewer"
	default:
		subjectKey = "email.subject.project_access_request_viewer"
	}
	subject := c.T(locale, subjectKey, d)

	if opts.Body == "" {
		opts.Body = c.THtml(locale, "email.body.project_access_request", d)
	}
	opts.ApproveButtonText = c.T(locale, "email.button.approve_request", nil)
	opts.DenyButtonText = c.T(locale, "email.button.deny_request", nil)

	buf := new(bytes.Buffer)
	err := c.templates.Lookup("project_access_request.html").Execute(buf, opts)
	if err != nil {
		return fmt.Errorf("email template error: %w", err)
	}
	htmlStr := buf.String()
	return c.Sender.Send(&Message{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Subject: subject,
		HTML:    htmlStr,
		Text:    htmlToText(htmlStr),
	})
}

// ---------------------------------------------------------------------------
// Project access granted
// ---------------------------------------------------------------------------

// ProjectAccessGranted holds the parameters for a project-access-granted email.
type ProjectAccessGranted struct {
	ToEmail     string
	ToName      string
	Locale      string
	OpenURL     string
	OrgName     string
	ProjectName string
}

func (c *Client) SendProjectAccessGranted(opts *ProjectAccessGranted) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":     opts.OrgName,
		"ProjectName": opts.ProjectName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.project_access_granted", d),
		PreButton:  c.THtml(locale, "email.body.project_access_granted", d),
		ButtonText: c.T(locale, "email.button.view_project", d),
		ButtonLink: opts.OpenURL,
	})
}

// ---------------------------------------------------------------------------
// Project access rejected
// ---------------------------------------------------------------------------

// ProjectAccessRejected holds the parameters for a project-access-rejected email.
type ProjectAccessRejected struct {
	ToEmail     string
	ToName      string
	Locale      string
	OrgName     string
	ProjectName string
}

func (c *Client) SendProjectAccessRejected(opts *ProjectAccessRejected) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":     opts.OrgName,
		"ProjectName": opts.ProjectName,
	})

	return c.SendInformational(&Informational{
		ToEmail: opts.ToEmail,
		ToName:  opts.ToName,
		Locale:  locale,
		Subject: c.T(locale, "email.subject.project_access_rejected", d),
		Body:    c.THtml(locale, "email.body.project_access_rejected", d),
	})
}

// ---------------------------------------------------------------------------
// Invoice — payment failed
// ---------------------------------------------------------------------------

// InvoicePaymentFailed holds the parameters for a payment-failed email.
type InvoicePaymentFailed struct {
	ToEmail            string
	ToName             string
	Locale             string
	OrgName            string
	Currency           string
	Amount             string
	PaymentURL         string
	GracePeriodEndDate time.Time
}

func (c *Client) SendInvoicePaymentFailed(opts *InvoicePaymentFailed) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":            opts.OrgName,
		"GracePeriodEndDate": opts.GracePeriodEndDate.Format(dateFormat),
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.invoice_payment_failed", d),
		PreButton:  c.THtml(locale, "email.body.invoice_payment_failed", d),
		ButtonText: c.T(locale, "email.button.update_payment_info", nil),
		ButtonLink: opts.PaymentURL,
		ShowFooter: true,
	})
}

// ---------------------------------------------------------------------------
// Invoice — payment success
// ---------------------------------------------------------------------------

// InvoicePaymentSuccess holds the parameters for a payment-success email.
type InvoicePaymentSuccess struct {
	ToEmail        string
	ToName         string
	Locale         string
	OrgName        string
	PaymentDate    time.Time
	BillingPageURL string
}

// SendInvoicePaymentSuccess is used when a previously failed invoice payment succeeds.
func (c *Client) SendInvoicePaymentSuccess(opts *InvoicePaymentSuccess) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":        opts.OrgName,
		"PaymentDate":    opts.PaymentDate.Format(dateFormat),
		"BillingPageURL": opts.BillingPageURL,
	})

	return c.SendInformational(&Informational{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.invoice_payment_success", d),
		Body:       c.THtml(locale, "email.body.invoice_payment_success", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Invoice — unpaid
// ---------------------------------------------------------------------------

// InvoiceUnpaid holds the parameters for an overdue-invoice email.
type InvoiceUnpaid struct {
	ToEmail    string
	ToName     string
	Locale     string
	OrgName    string
	PaymentURL string
}

// SendInvoiceUnpaid is sent after the payment grace period has ended.
func (c *Client) SendInvoiceUnpaid(opts *InvoiceUnpaid) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName": opts.OrgName,
		"ToName":  opts.ToName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.invoice_unpaid", d),
		PreButton:  c.THtml(locale, "email.body.invoice_unpaid", d),
		ButtonText: c.T(locale, "email.button.update_payment_info", nil),
		ButtonLink: opts.PaymentURL,
		ShowFooter: true,
	})
}

// ---------------------------------------------------------------------------
// Subscription cancelled
// ---------------------------------------------------------------------------

// SubscriptionCancelled holds the parameters for a subscription-cancelled email.
type SubscriptionCancelled struct {
	ToEmail    string
	ToName     string
	Locale     string
	OrgName    string
	PlanName   string
	BillingURL string
	EndDate    time.Time
}

func (c *Client) SendSubscriptionCancelled(opts *SubscriptionCancelled) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":  opts.OrgName,
		"PlanName": opts.PlanName,
		"ToName":   opts.ToName,
		"EndDate":  opts.EndDate.Format(dateFormat),
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.subscription_cancelled", d),
		PreButton:  c.THtml(locale, "email.body.subscription_cancelled", d),
		ButtonText: c.T(locale, "email.button.billing_settings", nil),
		ButtonLink: opts.BillingURL,
		PostButton: c.THtml(locale, "email.postbutton.subscription_cancelled", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Subscription ended
// ---------------------------------------------------------------------------

// SubscriptionEnded holds the parameters for a subscription-ended email.
type SubscriptionEnded struct {
	ToEmail    string
	ToName     string
	Locale     string
	OrgName    string
	BillingURL string
}

func (c *Client) SendSubscriptionEnded(opts *SubscriptionEnded) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName": opts.OrgName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.subscription_ended", d),
		PreButton:  c.THtml(locale, "email.body.subscription_ended", d),
		ButtonText: c.T(locale, "email.button.billing_settings", nil),
		ButtonLink: opts.BillingURL,
		PostButton: c.THtml(locale, "email.postbutton.subscription_ended", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Trial started
// ---------------------------------------------------------------------------

// TrialStarted holds the parameters for a trial-started email.
type TrialStarted struct {
	ToEmail      string
	ToName       string
	Locale       string
	OrgName      string
	FrontendURL  string
	TrialEndDate time.Time
}

func (c *Client) SendTrialStarted(opts *TrialStarted) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":      opts.OrgName,
		"TrialEndDate": opts.TrialEndDate.Format(dateFormat),
	})

	return c.SendWelcomeToTrial(&Welcome{
		ToEmail:     opts.ToEmail,
		ToName:      opts.ToName,
		Locale:      locale,
		Subject:     c.T(locale, "email.subject.trial_started", d),
		FrontendURL: opts.FrontendURL,
		WelcomeText: c.THtml(locale, "email.body.trial_started", d),
	})
}

// ---------------------------------------------------------------------------
// Trial ending soon
// ---------------------------------------------------------------------------

// TrialEndingSoon holds the parameters for a trial-ending-soon email.
type TrialEndingSoon struct {
	ToEmail      string
	ToName       string
	Locale       string
	OrgName      string
	UpgradeURL   string
	TrialEndDate time.Time
}

func (c *Client) SendTrialEndingSoon(opts *TrialEndingSoon) error {
	diff := time.Until(opts.TrialEndDate)
	days := int(math.Round(diff.Hours() / 24))
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":      opts.OrgName,
		"ToName":       opts.ToName,
		"TrialEndDate": opts.TrialEndDate.Format(dateFormat),
		"Days":         days,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.trial_ending_soon", d),
		PreButton:  c.THtml(locale, "email.body.trial_ending_soon", d),
		ButtonText: c.T(locale, "email.button.upgrade_now", nil),
		ButtonLink: opts.UpgradeURL,
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Trial ended
// ---------------------------------------------------------------------------

// TrialEnded holds the parameters for a trial-ended email.
type TrialEnded struct {
	ToEmail            string
	ToName             string
	Locale             string
	OrgName            string
	UpgradeURL         string
	GracePeriodEndDate time.Time
}

func (c *Client) SendTrialEnded(opts *TrialEnded) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":            opts.OrgName,
		"ToName":             opts.ToName,
		"GracePeriodEndDate": opts.GracePeriodEndDate.Format(dateFormat),
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.trial_ended", d),
		PreButton:  c.THtml(locale, "email.body.trial_ended", d),
		ButtonText: c.T(locale, "email.button.upgrade_to_team_plan", nil),
		ButtonLink: opts.UpgradeURL,
		ShowFooter: true,
	})
}

// ---------------------------------------------------------------------------
// Trial grace period ended
// ---------------------------------------------------------------------------

// TrialGracePeriodEnded holds the parameters for a trial-grace-period-ended email.
type TrialGracePeriodEnded struct {
	ToEmail    string
	ToName     string
	Locale     string
	OrgName    string
	UpgradeURL string
}

func (c *Client) SendTrialGracePeriodEnded(opts *TrialGracePeriodEnded) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName": opts.OrgName,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.trial_grace_period_ended", d),
		PreButton:  c.THtml(locale, "email.body.trial_grace_period_ended", d),
		ButtonText: c.T(locale, "email.button.upgrade_to_team_plan", nil),
		ButtonLink: opts.UpgradeURL,
		PostButton: c.THtml(locale, "email.postbutton.trial_grace_period_ended", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Credit trial started
// ---------------------------------------------------------------------------

// CreditTrialStarted holds the parameters for a credit-trial-started email.
type CreditTrialStarted struct {
	ToEmail          string
	ToName           string
	Locale           string
	OrgName          string
	FrontendURL      string
	CreditAllocation int
}

func (c *Client) SendCreditTrialStarted(opts *CreditTrialStarted) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":          opts.OrgName,
		"CreditAllocation": opts.CreditAllocation,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.credit_trial_started", d),
		PreButton:  c.THtml(locale, "email.body.credit_trial_started", d),
		ButtonText: c.T(locale, "email.button.open_product_cloud", d),
		ButtonLink: opts.FrontendURL,
		PostButton: c.THtml(locale, "email.postbutton.credit_trial_started", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Credit trial low
// ---------------------------------------------------------------------------

// CreditTrialLow holds the parameters for a credit-trial-low email.
type CreditTrialLow struct {
	ToEmail          string
	ToName           string
	Locale           string
	OrgName          string
	FrontendURL      string
	UpgradeURL       string
	CreditAllocation int
	RemainingBalance float64
}

func (c *Client) SendCreditTrialLow(opts *CreditTrialLow) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":          opts.OrgName,
		"CreditAllocation": opts.CreditAllocation,
		"RemainingBalance": fmt.Sprintf("%.2f", opts.RemainingBalance),
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.credit_trial_low", d),
		PreButton:  c.THtml(locale, "email.body.credit_trial_low", d),
		ButtonText: c.T(locale, "email.button.upgrade_to_pro", nil),
		ButtonLink: opts.UpgradeURL,
		PostButton: c.THtml(locale, "email.postbutton.credit_trial_low", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Credit trial depleted
// ---------------------------------------------------------------------------

// CreditTrialDepleted holds the parameters for a credit-trial-depleted email.
type CreditTrialDepleted struct {
	ToEmail          string
	ToName           string
	Locale           string
	OrgName          string
	FrontendURL      string
	UpgradeURL       string
	CreditAllocation int
}

func (c *Client) SendCreditTrialDepleted(opts *CreditTrialDepleted) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":          opts.OrgName,
		"CreditAllocation": opts.CreditAllocation,
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.credit_trial_depleted", d),
		PreButton:  c.THtml(locale, "email.body.credit_trial_depleted", d),
		ButtonText: c.T(locale, "email.button.upgrade_to_pro", nil),
		ButtonLink: opts.UpgradeURL,
		PostButton: c.THtml(locale, "email.postbutton.credit_trial_depleted", d),
		ShowFooter: false,
	})
}

// ---------------------------------------------------------------------------
// Plan update
// ---------------------------------------------------------------------------

// PlanUpdate holds the parameters for a plan-update email.
type PlanUpdate struct {
	ToEmail  string
	ToName   string
	Locale   string
	OrgName  string
	PlanName string
}

func (c *Client) SendPlanUpdate(opts *PlanUpdate) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":  opts.OrgName,
		"PlanName": opts.PlanName,
	})

	return c.SendInformational(&Informational{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.plan_update", d),
		Body:       c.THtml(locale, "email.body.plan_update", d),
		ShowFooter: true,
	})
}

// ---------------------------------------------------------------------------
// Subscription renewed
// ---------------------------------------------------------------------------

// SubscriptionRenewed holds the parameters for a subscription-renewed email.
type SubscriptionRenewed struct {
	ToEmail  string
	ToName   string
	Locale   string
	OrgName  string
	PlanName string
}

func (c *Client) SendSubscriptionRenewed(opts *SubscriptionRenewed) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":  opts.OrgName,
		"PlanName": opts.PlanName,
	})

	return c.SendInformational(&Informational{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.subscription_renewed", d),
		Body:       c.THtml(locale, "email.body.subscription_renewed", d),
		ShowFooter: true,
	})
}

// ---------------------------------------------------------------------------
// Paid plan
// ---------------------------------------------------------------------------

// PaidPlan holds the parameters for paid-plan emails (started or renewed).
type PaidPlan struct {
	ToEmail          string
	ToName           string
	Locale           string
	OrgName          string
	FrontendURL      string
	BillingURL       string
	PlanName         string
	BillingStartDate time.Time
}

// SendPaidPlanStarted sends a customised plan-started email for a paid plan (Team or Pro).
func (c *Client) SendPaidPlanStarted(opts *PaidPlan) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":          opts.OrgName,
		"PlanName":         opts.PlanName,
		"BillingStartDate": opts.BillingStartDate.Format(dateFormat),
	})

	return c.SendCallToAction(&CallToAction{
		ToEmail:    opts.ToEmail,
		ToName:     opts.ToName,
		Locale:     locale,
		Subject:    c.T(locale, "email.subject.paid_plan_started", d),
		PreButton:  c.THtml(locale, "email.body.paid_plan_started", d),
		ButtonText: c.T(locale, "email.button.view_billing_dashboard", nil),
		ButtonLink: opts.BillingURL,
		PostButton: c.THtml(locale, "email.postbutton.paid_plan_started", nil),
		ShowFooter: false,
	})
}

// SendPaidPlanRenewal sends a customised plan-renewed email for a paid plan (Team or Pro).
func (c *Client) SendPaidPlanRenewal(opts *PaidPlan) error {
	locale := opts.Locale
	d := c.mergeData(map[string]any{
		"OrgName":          opts.OrgName,
		"PlanName":         opts.PlanName,
		"BillingStartDate": opts.BillingStartDate.Format(dateFormat),
	})

	return c.SendWelcomeToTeam(&Welcome{
		ToEmail:     opts.ToEmail,
		ToName:      opts.ToName,
		Locale:      locale,
		Subject:     c.T(locale, "email.subject.paid_plan_renewal", d),
		FrontendURL: opts.FrontendURL,
		WelcomeText: c.THtml(locale, "email.body.paid_plan_renewal", d),
	})
}
