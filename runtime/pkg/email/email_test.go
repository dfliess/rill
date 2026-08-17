package email

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rilldata/rill/admin/database"
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
}

func (m *mockSender) Send(toEmail, toName, subject, body string) error {
	m.toEmail = toEmail
	m.toName = toName
	m.subject = subject
	m.body = body
	return nil
}

func TestCallToActionBranding(t *testing.T) {
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
	require.Contains(t, mock.body, "Inteligencia para empresas")
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
	require.Contains(t, mock.body, "recuperó")
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

// TestKairosBranding renders a representative email of each layout and asserts
// it carries the Kairos brand (logo + footer) and none of Rill's default marks.
// It guards the templates after regenerating templates/gen/*.html.
func TestKairosBranding(t *testing.T) {
	mock := &mockSender{}
	client := New(mock)

	cases := map[string]func() error{
		"scheduled_report": func() error {
			return client.SendScheduledReport(&ScheduledReport{
				ToEmail: "a@b.com", DisplayName: "Ventas", ReportTime: time.Now(),
				OpenLink: "https://example.com", UnsubscribeLink: "https://example.com/u",
			})
		},
		"alert_fail": func() error {
			return client.SendAlertStatus(&drivers.AlertStatus{
				ToEmail: "a@b.com", DisplayName: "Alerta", ExecutionTime: time.Now(),
				Status:  runtimev1.AssertionStatus_ASSERTION_STATUS_FAIL,
				FailRow: map[string]any{"x": 1}, OpenLink: "https://example.com", EditLink: "https://example.com",
			})
		},
		"organization_invite": func() error {
			return client.SendOrganizationInvite(&OrganizationInvite{
				ToEmail: "a@b.com", AcceptURL: "https://example.com", OrgName: "acme", RoleName: "editor",
			})
		},
		"project_access_request": func() error {
			return client.SendProjectAccessRequest(&ProjectAccessRequest{
				ToEmail: "a@b.com", Email: "user@b.com", OrgName: "acme", ProjectName: "sales",
				Role: database.ProjectRoleNameViewer, ApproveLink: "https://example.com/a", DenyLink: "https://example.com/d",
			})
		},
		"project_access_rejected": func() error {
			return client.SendProjectAccessRejected(&ProjectAccessRejected{
				ToEmail: "a@b.com", OrgName: "acme", ProjectName: "sales",
			})
		},
	}

	rillMarks := []string{"rilldata.com", "Rill Data", "rill-logo", "#4736F5", "#3524C7", "Bartol"}

	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, fn())
			require.Contains(t, mock.body, "kairosagentica.com/img/logo/kairos-lockup-horizontal.png", "falta el logo de Kairos")
			require.Contains(t, mock.body, "Inteligencia para empresas", "falta el footer de Kairos")
			for _, mark := range rillMarks {
				require.NotContains(t, mock.body, mark, "marca Rill residual: %s", mark)
			}
		})
	}
}
