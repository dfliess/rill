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
