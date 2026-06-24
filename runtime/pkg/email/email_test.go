package email

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/stretchr/testify/require"
)

type mockSender struct {
	fromEmail string
	fromName  string
	toEmail   string
	toName    string
	subject   string
	body      string
	text      string
}

func (m *mockSender) Send(msg *Message) error {
	m.toEmail = msg.ToEmail
	m.toName = msg.ToName
	m.subject = msg.Subject
	m.body = msg.HTML
	m.text = msg.Text
	return nil
}

func TestCopyrightYear(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	opts := &CallToAction{
		ToEmail:    uuid.New().String(),
		ToName:     uuid.New().String(),
		Subject:    uuid.New().String(),
		ButtonText: uuid.New().String(),
		ButtonLink: uuid.New().String(),
		ShowFooter: true,
	}
	err := client.SendCallToAction(opts)
	require.NoError(t, err)

	require.Equal(t, opts.ToEmail, mock.toEmail)
	require.Equal(t, opts.ToName, mock.toName)
	require.Equal(t, opts.Subject, mock.subject)
	require.Contains(t, mock.body, opts.ButtonText)
	require.Contains(t, mock.body, opts.ButtonLink)

	year := time.Now().Year()
	require.Contains(t, mock.body, fmt.Sprintf("© %d Rill Data, Inc", year))
}

func TestOrganizationInvite(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	opts := &OrganizationInvite{
		ToEmail:       uuid.New().String(),
		ToName:        uuid.New().String(),
		AcceptURL:     "https://api.example.com",
		OrgName:       uuid.New().String(),
		RoleName:      uuid.New().String(),
		InvitedByName: uuid.New().String(),
	}
	err := client.SendOrganizationInvite(opts)
	require.NoError(t, err)

	require.Equal(t, opts.ToEmail, mock.toEmail)
	require.Equal(t, opts.ToName, mock.toName)
	require.NotEmpty(t, mock.subject)
	require.Contains(t, mock.body, opts.OrgName)
	require.Contains(t, mock.body, opts.RoleName)
	require.Contains(t, mock.body, opts.InvitedByName)
}

func TestAlertFail(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	opts := &drivers.AlertStatus{
		ToEmail:       uuid.New().String(),
		ToName:        uuid.New().String(),
		DisplayName:   "Foobar",
		ExecutionTime: time.Date(2024, 01, 27, 0, 0, 0, 0, time.UTC),
		Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL,
		FailRow:       map[string]any{"hello": "world", "pi": 3.14},
		OpenLink:      "https://example.com",
		EditLink:      "https://example.com",
	}
	err := client.SendAlertStatus(opts)
	require.NoError(t, err)

	require.Equal(t, opts.ToEmail, mock.toEmail)
	require.Equal(t, opts.ToName, mock.toName)
	require.NotEmpty(t, mock.subject)
	require.Contains(t, mock.body, opts.DisplayName)
	require.Contains(t, mock.body, opts.ExecutionTime.Format(time.RFC1123))
	for k, v := range opts.FailRow {
		require.Contains(t, mock.body, k)
		require.Contains(t, mock.body, fmt.Sprintf("%v", v))
	}
}

func TestAlertRecover(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	opts := &drivers.AlertStatus{
		ToEmail:       uuid.New().String(),
		ToName:        uuid.New().String(),
		DisplayName:   "Foobar",
		ExecutionTime: time.Date(2024, 01, 27, 0, 0, 0, 0, time.UTC),
		Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_PASS,
		IsRecover:     true,
		OpenLink:      "https://example.com",
		EditLink:      "https://example.com",
	}
	err := client.SendAlertStatus(opts)
	require.NoError(t, err)

	require.Equal(t, opts.ToEmail, mock.toEmail)
	require.Equal(t, opts.ToName, mock.toName)
	require.NotEmpty(t, mock.subject)
	require.Contains(t, mock.body, opts.DisplayName)
	require.Contains(t, mock.body, opts.ExecutionTime.Format(time.RFC1123))
	require.Contains(t, mock.body, "recovered")
}

func TestAlertError(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	opts := &drivers.AlertStatus{
		ToEmail:        uuid.New().String(),
		ToName:         uuid.New().String(),
		DisplayName:    "Foobar",
		ExecutionTime:  time.Date(2024, 01, 27, 0, 0, 0, 0, time.UTC),
		Status:         runtimev1.AssertionStatus_ASSERTION_STATUS_ERROR,
		ExecutionError: "hello error",
		OpenLink:       "https://example.com",
		EditLink:       "https://example.com",
	}
	err := client.SendAlertStatus(opts)
	require.NoError(t, err)

	require.Equal(t, opts.ToEmail, mock.toEmail)
	require.Equal(t, opts.ToName, mock.toName)
	require.NotEmpty(t, mock.subject)
	require.Contains(t, mock.body, opts.DisplayName)
	require.Contains(t, mock.body, opts.ExecutionTime.Format(time.RFC1123))
	require.Contains(t, mock.body, "hello error")
}

func TestHTMLToText(t *testing.T) {
	html := `<p>Hello <b>World</b></p><ul><li>item one</li><li>item two</li></ul><p>End</p>`
	text := htmlToText(html)
	require.Contains(t, text, "Hello World")
	require.Contains(t, text, "- item one")
	require.Contains(t, text, "- item two")
	require.Contains(t, text, "End")
	require.NotContains(t, text, "<")
}

func TestHTMLToTextEntities(t *testing.T) {
	text := htmlToText("&amp; &lt;b&gt; &#169;")
	require.Contains(t, text, "& <b> ©")
}

func TestMultipartMessageHasText(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	opts := &CallToAction{
		ToEmail:    "test@example.com",
		ToName:     "Test User",
		Subject:    "Test Subject",
		PreButton:  "<p>Hello</p>",
		ButtonText: "Click me",
		ButtonLink: "https://example.com",
	}
	err := client.SendCallToAction(opts)
	require.NoError(t, err)
	require.NotEmpty(t, mock.body)
	require.NotEmpty(t, mock.text)
	require.NotContains(t, mock.text, "<p>")
	require.Contains(t, mock.text, "Click me")
}

func TestTestSenderCapturesText(t *testing.T) {
	ts := NewTestSender().(*TestSender)
	client := New(ts)

	err := client.SendCallToAction(&CallToAction{
		ToEmail:    "to@example.com",
		ToName:     "Recipient",
		Subject:    "Test",
		PreButton:  "<p>Hello</p>",
		ButtonText: "Click",
		ButtonLink: "https://example.com",
	})
	require.NoError(t, err)
	require.Len(t, ts.Emails, 1)
	require.NotEmpty(t, ts.Emails[0].Body)
	require.NotEmpty(t, ts.Emails[0].Text)
	require.NotContains(t, ts.Emails[0].Text, "<p>")
}

func TestDefaultBranding(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendCallToAction(&CallToAction{
		ToEmail:    "test@example.com",
		ToName:     "User",
		Subject:    "Test",
		ButtonText: "Click",
		ButtonLink: "https://example.com",
		ShowFooter: true,
	})
	require.NoError(t, err)
	require.Contains(t, mock.body, "Rill Data, Inc.")
	require.Contains(t, mock.body, "#4736F5")
	require.Contains(t, mock.body, "rill-logo-purple-square.png")
	require.Contains(t, mock.body, "The Rill Team")
}

func TestCustomBranding(t *testing.T) {
	mock := &mockSender{}
	client := New(mock, WithBranding(&Branding{
		LogoURL:         "https://example.com/logo.png",
		CompanyName:     "Acme Corp.",
		ProductName:     "AcmeDash",
		PrimaryColor:    "#FF0000",
		SupportEmail:    "help@acme.com",
		Address:         "1 Main St, Anytown",
		ContactURL:      "https://acme.com/contact",
		CommunityURL:    "https://acme.com/community",
		PrivacyURL:      "https://acme.com/privacy",
		DocsURL:         "https://docs.acme.com",
		ReleaseNotesURL: "https://docs.acme.com/notes",
		ChatURL:         "https://acme.com/chat",
		TeamSignature:   "The Acme Team",
	}))

	err := client.SendCallToAction(&CallToAction{
		ToEmail:    "test@example.com",
		ToName:     "User",
		Subject:    "Test",
		ButtonText: "Click",
		ButtonLink: "https://example.com",
		ShowFooter: true,
	})
	require.NoError(t, err)

	require.Contains(t, mock.body, "Acme Corp.")
	require.Contains(t, mock.body, "https://example.com/logo.png")
	require.Contains(t, mock.body, "#FF0000")
	require.Contains(t, mock.body, "help@acme.com")
	require.Contains(t, mock.body, "1 Main St, Anytown")
	require.Contains(t, mock.body, "The Acme Team")
	require.Contains(t, mock.body, "AcmeDash.")

	require.NotContains(t, mock.body, "Rill Data, Inc.")
	require.NotContains(t, mock.body, "rill-logo-purple")
	require.NotContains(t, mock.body, "The Rill Team")
}

func TestBrandingFromOrg(t *testing.T) {
	mock := &mockSender{}
	// Simulate an org with display name and logo URL.
	b := BrandingFromOrg("Kairos Analytics", "https://cdn.kairos.io/logo.png")
	client := New(mock, WithBranding(b))

	err := client.SendCallToAction(&CallToAction{
		ToEmail:    "test@example.com",
		ToName:     "User",
		Subject:    "Test",
		ButtonText: "Click",
		ButtonLink: "https://example.com",
		ShowFooter: true,
	})
	require.NoError(t, err)

	// Org-derived values appear
	require.Contains(t, mock.body, "Kairos Analytics")
	require.Contains(t, mock.body, "The Kairos Analytics Team")
	require.Contains(t, mock.body, "https://cdn.kairos.io/logo.png")

	// Other defaults are preserved
	require.Contains(t, mock.body, "#4736F5")
	require.Contains(t, mock.body, "support@rilldata.com")
}

func TestBrandingFromOrgEmptyFields(t *testing.T) {
	mock := &mockSender{}
	// No display name or logo — should produce Rill defaults.
	b := BrandingFromOrg("", "")
	client := New(mock, WithBranding(b))

	err := client.SendCallToAction(&CallToAction{
		ToEmail:    "test@example.com",
		ToName:     "User",
		Subject:    "Test",
		ButtonText: "Click",
		ButtonLink: "https://example.com",
	})
	require.NoError(t, err)
	require.Contains(t, mock.body, "Rill Data, Inc.")
	require.Contains(t, mock.body, "rill-logo-purple-square.png")
	require.Contains(t, mock.body, "The Rill Team")
}

func TestCustomBrandingInGoCode(t *testing.T) {
	mock := &mockSender{}
	client := New(mock, WithBranding(&Branding{
		LogoURL:         "https://example.com/logo.png",
		CompanyName:     "Acme Corp.",
		ProductName:     "AcmeDash",
		PrimaryColor:    "#FF0000",
		SupportEmail:    "help@acme.com",
		Address:         "1 Main St, Anytown",
		ContactURL:      "https://acme.com/contact",
		CommunityURL:    "https://acme.com/community",
		PrivacyURL:      "https://acme.com/privacy",
		DocsURL:         "https://docs.acme.com",
		ReleaseNotesURL: "https://docs.acme.com/notes",
		ChatURL:         "https://acme.com/chat",
		TeamSignature:   "The Acme Team",
	}))

	// Test that Go-code brand references use custom branding
	err := client.SendOrganizationInvite(&OrganizationInvite{
		ToEmail:   "test@example.com",
		ToName:    "User",
		AcceptURL: "https://example.com/accept",
		OrgName:   "TestOrg",
		RoleName:  "admin",
	})
	require.NoError(t, err)

	// Subject should contain custom product name, not "Rill"
	require.Contains(t, mock.subject, "AcmeDash")
	require.NotContains(t, mock.subject, "join Rill")

	// InvitedByName defaults to ProductName when empty
	require.Contains(t, mock.body, "AcmeDash has invited you")
	require.Contains(t, mock.body, "AcmeDash account")
}

func TestCustomBrandingProjectAccessGranted(t *testing.T) {
	mock := &mockSender{}
	client := New(mock, WithBranding(&Branding{
		LogoURL:         "https://example.com/logo.png",
		CompanyName:     "Acme Corp.",
		ProductName:     "AcmeDash",
		PrimaryColor:    "#FF0000",
		SupportEmail:    "help@acme.com",
		Address:         "1 Main St, Anytown",
		ContactURL:      "https://acme.com/contact",
		CommunityURL:    "https://acme.com/community",
		PrivacyURL:      "https://acme.com/privacy",
		DocsURL:         "https://docs.acme.com",
		ReleaseNotesURL: "https://docs.acme.com/notes",
		ChatURL:         "https://acme.com/chat",
		TeamSignature:   "The Acme Team",
	}))

	err := client.SendProjectAccessGranted(&ProjectAccessGranted{
		ToEmail:     "test@example.com",
		ToName:      "User",
		OpenURL:     "https://example.com",
		OrgName:     "TestOrg",
		ProjectName: "TestProject",
	})
	require.NoError(t, err)

	// Button text should use custom product name
	require.Contains(t, mock.body, "View project in AcmeDash")
	require.NotContains(t, mock.body, "View project in Rill")
}

// ---------------------------------------------------------------------------
// i18n tests
// ---------------------------------------------------------------------------

func TestDefaultLocaleEnglish(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendOrganizationInvite(&OrganizationInvite{
		ToEmail:       "test@example.com",
		ToName:        "User",
		AcceptURL:     "https://example.com/accept",
		OrgName:       "TestOrg",
		RoleName:      "admin",
		InvitedByName: "Alice",
	})
	require.NoError(t, err)

	// English catalog text
	require.Contains(t, mock.subject, "Alice invited you to join Rill")
	require.Contains(t, mock.body, "Accept invitation")
	require.Contains(t, mock.body, "has invited you to join")
}

func TestLocaleSpanish(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendOrganizationInvite(&OrganizationInvite{
		ToEmail:       "test@example.com",
		ToName:        "User",
		Locale:        "es",
		AcceptURL:     "https://example.com/accept",
		OrgName:       "TestOrg",
		RoleName:      "admin",
		InvitedByName: "Alice",
	})
	require.NoError(t, err)

	// Spanish catalog text
	require.Contains(t, mock.subject, "Alice le invitó a unirse a Rill")
	require.Contains(t, mock.body, "Aceptar invitación")
	require.Contains(t, mock.body, "le ha invitado a unirse a")
}

func TestLocaleSpanishProjectAccessGranted(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendProjectAccessGranted(&ProjectAccessGranted{
		ToEmail:     "test@example.com",
		ToName:      "User",
		Locale:      "es",
		OpenURL:     "https://example.com",
		OrgName:     "TestOrg",
		ProjectName: "TestProject",
	})
	require.NoError(t, err)

	require.Contains(t, mock.subject, "aprobada")
	require.Contains(t, mock.body, "Ver proyecto en Rill")
}

func TestLocaleSpanishAlertFail(t *testing.T) {
	mock := &mockSender{}
	client := New(mock, WithDefaultLocale("es"))

	opts := &drivers.AlertStatus{
		ToEmail:       "test@example.com",
		ToName:        "User",
		DisplayName:   "Foobar",
		ExecutionTime: time.Date(2024, 01, 27, 0, 0, 0, 0, time.UTC),
		Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL,
		FailRow:       map[string]any{"hello": "world"},
		OpenLink:      "https://example.com",
		EditLink:      "https://example.com",
	}
	err := client.SendAlertStatus(opts)
	require.NoError(t, err)

	// Spanish: "Su alerta se activó"
	require.Contains(t, mock.body, "Su alerta se activó")
	require.Contains(t, mock.body, "Abrir en el navegador")
}

func TestLocaleSpanishAlertRecover(t *testing.T) {
	mock := &mockSender{}
	client := New(mock, WithDefaultLocale("es"))

	opts := &drivers.AlertStatus{
		ToEmail:       "test@example.com",
		ToName:        "User",
		DisplayName:   "Foobar",
		ExecutionTime: time.Date(2024, 01, 27, 0, 0, 0, 0, time.UTC),
		Status:        runtimev1.AssertionStatus_ASSERTION_STATUS_PASS,
		IsRecover:     true,
		OpenLink:      "https://example.com",
		EditLink:      "https://example.com",
	}
	err := client.SendAlertStatus(opts)
	require.NoError(t, err)

	// Spanish: "La alerta se ha recuperado"
	require.Contains(t, mock.body, "recuperado")
	require.Contains(t, mock.subject, "Recuperada:")
}

func TestLocaleSpanishScheduledReport(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendScheduledReport(&ScheduledReport{
		ToEmail:        "test@example.com",
		ToName:         "User",
		Locale:         "es",
		DisplayName:    "Daily Report",
		ReportTime:     time.Date(2024, 01, 27, 0, 0, 0, 0, time.UTC),
		DownloadFormat: "CSV",
		OpenLink:       "https://example.com",
		DownloadLink:   "https://example.com/dl",
		EditLink:       "https://example.com/edit",
	})
	require.NoError(t, err)

	// Spanish text
	require.Contains(t, mock.body, "está listo para ver")
	require.Contains(t, mock.body, "Abrir en el navegador")
	require.Contains(t, mock.body, "Descargar archivo CSV")
}

func TestWithDefaultLocaleOption(t *testing.T) {
	mock := &mockSender{}
	client := New(mock, WithDefaultLocale("es"))

	err := client.SendOrganizationInvite(&OrganizationInvite{
		ToEmail:       "test@example.com",
		ToName:        "User",
		AcceptURL:     "https://example.com/accept",
		OrgName:       "TestOrg",
		RoleName:      "admin",
		InvitedByName: "Alice",
	})
	require.NoError(t, err)

	// Default locale is es, so even without Locale on opts, we get Spanish
	require.Contains(t, mock.subject, "le invitó a unirse")
	require.Contains(t, mock.body, "Aceptar invitación")
}

func TestFallbackToEnglishForUnknownLocale(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendOrganizationInvite(&OrganizationInvite{
		ToEmail:       "test@example.com",
		ToName:        "User",
		Locale:        "fr", // no French catalog
		AcceptURL:     "https://example.com/accept",
		OrgName:       "TestOrg",
		RoleName:      "admin",
		InvitedByName: "Alice",
	})
	require.NoError(t, err)

	// Falls back to English
	require.Contains(t, mock.subject, "Alice invited you to join Rill")
	require.Contains(t, mock.body, "Accept invitation")
}

func TestTMethod(t *testing.T) {
	client := New(NewNoopSender())

	// English
	s := client.T("en", "email.button.accept_invitation", nil)
	require.Equal(t, "Accept invitation", s)

	// Spanish
	s = client.T("es", "email.button.accept_invitation", nil)
	require.Equal(t, "Aceptar invitación", s)

	// Unknown message ID returns the ID itself
	s = client.T("en", "email.nonexistent.key", nil)
	require.Equal(t, "email.nonexistent.key", s)
}

func TestTHtmlMethod(t *testing.T) {
	client := New(NewNoopSender())

	h := client.THtml("en", "email.body.org_invite", map[string]any{
		"InvitedByName": "Alice",
		"OrgName":       "TestOrg",
		"RoleName":      "admin",
		"ProductName":   "Rill",
	})
	require.Contains(t, string(h), "Alice has invited you")
	require.Contains(t, string(h), "<b>TestOrg</b>")
}

func TestTHtmlXSSPrevention(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendOrganizationInvite(&OrganizationInvite{
		ToEmail:       "test@example.com",
		ToName:        "User",
		AcceptURL:     "https://example.com/accept",
		OrgName:       "TestOrg",
		RoleName:      "admin",
		InvitedByName: `<script>alert(1)</script>`,
	})
	require.NoError(t, err)

	// The raw script tag must not appear — it should be escaped
	require.NotContains(t, mock.body, "<script>")
	require.Contains(t, mock.body, "&lt;script&gt;")
}

func TestLocaleSpanishInvoicePaymentFailed(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendInvoicePaymentFailed(&InvoicePaymentFailed{
		ToEmail:            "test@example.com",
		ToName:             "User",
		Locale:             "es",
		OrgName:            "TestOrg",
		PaymentURL:         "https://example.com/pay",
		GracePeriodEndDate: time.Date(2024, 06, 15, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	require.Contains(t, mock.subject, "El pago de TestOrg ha fallado")
	require.Contains(t, mock.body, "Actualizar información de pago")
}

func TestLocaleSpanishTrialStarted(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	err := client.SendTrialStarted(&TrialStarted{
		ToEmail:      "test@example.com",
		ToName:       "User",
		Locale:       "es",
		OrgName:      "TestOrg",
		FrontendURL:  "https://example.com",
		TrialEndDate: time.Date(2024, 07, 27, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	require.Contains(t, mock.subject, "prueba gratuita de 30 días")
}

func TestLocalePerRequestOverridesDefault(t *testing.T) {
	mock := &mockSender{}
	// Default locale is English
	client := New(mock)

	// But this specific request uses Spanish
	err := client.SendPlanUpdate(&PlanUpdate{
		ToEmail:  "test@example.com",
		ToName:   "User",
		Locale:   "es",
		OrgName:  "TestOrg",
		PlanName: "Pro",
	})
	require.NoError(t, err)

	require.Contains(t, mock.subject, "ha sido actualizado al plan Pro")
}

func TestTrialEndingSoonPluralDays(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	// days=1 should use singular form "1 day" (not "1 days")
	err := client.SendTrialEndingSoon(&TrialEndingSoon{
		ToEmail:      "test@example.com",
		ToName:       "User",
		OrgName:      "TestOrg",
		UpgradeURL:   "https://example.com/upgrade",
		TrialEndDate: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)
	require.Contains(t, mock.subject, "1 day")
	require.NotContains(t, mock.subject, "1 days")

	// days=5 should use plural form "5 days"
	err = client.SendTrialEndingSoon(&TrialEndingSoon{
		ToEmail:      "test@example.com",
		ToName:       "User",
		OrgName:      "TestOrg",
		UpgradeURL:   "https://example.com/upgrade",
		TrialEndDate: time.Now().Add(5 * 24 * time.Hour),
	})
	require.NoError(t, err)
	require.Contains(t, mock.subject, "5 days")
}

func TestTrialEndingSoonPluralSpanish(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	// days=1 in Spanish should use singular "1 día"
	err := client.SendTrialEndingSoon(&TrialEndingSoon{
		ToEmail:      "test@example.com",
		ToName:       "User",
		Locale:       "es",
		OrgName:      "TestOrg",
		UpgradeURL:   "https://example.com/upgrade",
		TrialEndDate: time.Now().Add(24 * time.Hour),
	})
	require.NoError(t, err)
	require.Contains(t, mock.subject, "1 día")
	require.NotContains(t, mock.subject, "1 días")
}

func TestSanitizeHeader(t *testing.T) {
	result := sanitizeHeader("normal value")
	require.Equal(t, "normal value", result)

	result = sanitizeHeader("injected\r\nBcc: evil@example.com")
	require.Equal(t, "injectedBcc: evil@example.com", result)

	result = sanitizeHeader("line\rfeed\nonly")
	require.Equal(t, "linefeedonly", result)
}
