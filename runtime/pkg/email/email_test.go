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
