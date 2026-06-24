package email

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Message holds the content of a single outgoing email.
type Message struct {
	ToEmail string
	ToName  string
	Subject string
	HTML    string
	Text    string // plain-text alternative; derived from HTML if empty
}

type Sender interface {
	Send(msg *Message) error
}

type SMTPOptions struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	BCC          string
}

func sanitizeHeader(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

type smtpSender struct {
	opts *SMTPOptions
}

func NewSMTPSender(opts *SMTPOptions) (Sender, error) {
	if opts.SMTPPassword == "" {
		return nil, fmt.Errorf("SMTP server password is required")
	}

	_, err := mail.ParseAddress(opts.FromEmail)
	if err != nil {
		return nil, fmt.Errorf("invalid sender email address %q", opts.FromEmail)
	}

	if opts.BCC != "" {
		_, err := mail.ParseAddress(opts.BCC)
		if err != nil {
			return nil, fmt.Errorf("invalid bcc email address %q", opts.BCC)
		}
	}

	return &smtpSender{opts: opts}, nil
}

func (s *smtpSender) Send(msg *Message) error {
	if _, err := mail.ParseAddress(msg.ToEmail); err != nil {
		return fmt.Errorf("invalid recipient address %q: %w", msg.ToEmail, err)
	}

	from := mail.Address{Name: s.opts.FromName, Address: s.opts.FromEmail}
	to := mail.Address{Name: sanitizeHeader(msg.ToName), Address: msg.ToEmail}

	var buf bytes.Buffer
	// Common headers
	fmt.Fprintf(&buf, "From: %s\r\n", from.String())
	fmt.Fprintf(&buf, "To: %s\r\n", to.String())
	fmt.Fprintf(&buf, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", msg.Subject))
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")

	if msg.Text != "" {
		// multipart/alternative: text/plain before text/html (RFC 2046)
		mw := multipart.NewWriter(&buf)
		fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%q\r\n", mw.Boundary())
		fmt.Fprintf(&buf, "\r\n")

		// text/plain part (first = least preferred per RFC 2046)
		plainHeader := textproto.MIMEHeader{}
		plainHeader.Set("Content-Type", "text/plain; charset=utf-8")
		plainHeader.Set("Content-Transfer-Encoding", "quoted-printable")
		pw, err := mw.CreatePart(plainHeader)
		if err != nil {
			return fmt.Errorf("creating text/plain part: %w", err)
		}
		qpw := quotedprintable.NewWriter(pw)
		if _, err := qpw.Write([]byte(msg.Text)); err != nil {
			return fmt.Errorf("writing text/plain body: %w", err)
		}
		if err := qpw.Close(); err != nil {
			return fmt.Errorf("closing text/plain writer: %w", err)
		}

		// text/html part (last = preferred per RFC 2046)
		htmlHeader := textproto.MIMEHeader{}
		htmlHeader.Set("Content-Type", "text/html; charset=utf-8")
		htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")
		hw, err := mw.CreatePart(htmlHeader)
		if err != nil {
			return fmt.Errorf("creating text/html part: %w", err)
		}
		qpw = quotedprintable.NewWriter(hw)
		if _, err := qpw.Write([]byte(msg.HTML)); err != nil {
			return fmt.Errorf("writing text/html body: %w", err)
		}
		if err := qpw.Close(); err != nil {
			return fmt.Errorf("closing text/html writer: %w", err)
		}

		if err := mw.Close(); err != nil {
			return fmt.Errorf("closing multipart writer: %w", err)
		}
	} else {
		// Single-part HTML (backward compat)
		fmt.Fprintf(&buf, "Content-Type: text/html; charset=utf-8\r\n")
		fmt.Fprintf(&buf, "\r\n")
		buf.WriteString(msg.HTML)
		buf.WriteString("\r\n")
	}

	// Build recipients list
	recipients := []string{msg.ToEmail}
	if s.opts.BCC != "" {
		recipients = append(recipients, s.opts.BCC)
	}

	auth := smtp.PlainAuth("", s.opts.SMTPUsername, s.opts.SMTPPassword, s.opts.SMTPHost)
	return smtp.SendMail(
		s.opts.SMTPHost+":"+strconv.Itoa(s.opts.SMTPPort),
		auth, from.Address, recipients, buf.Bytes(),
	)
}

type consoleSender struct {
	logger    *zap.Logger
	fromEmail string
	fromName  string
}

func NewConsoleSender(logger *zap.Logger, fromEmail, fromName string) (Sender, error) {
	return &consoleSender{logger: logger, fromEmail: fromEmail, fromName: fromName}, nil
}

func (s *consoleSender) Send(msg *Message) error {
	s.logger.Info("email sent",
		zap.String("from_email", s.fromEmail),
		zap.String("from_name", s.fromName),
		zap.String("to_email", msg.ToEmail),
		zap.String("to_name", msg.ToName),
		zap.String("subject", msg.Subject),
		zap.String("body", msg.HTML),
	)
	return nil
}

type noopSender struct{}

func NewNoopSender() Sender {
	return &noopSender{}
}

func (s *noopSender) Send(msg *Message) error {
	return nil
}

type TestSender struct {
	Emails []struct {
		ToEmail string
		ToName  string
		Subject string
		Body    string
		Text    string
	}
}

func NewTestSender() Sender {
	return &TestSender{}
}

func (s *TestSender) Send(msg *Message) error {
	s.Emails = append(s.Emails, struct {
		ToEmail string
		ToName  string
		Subject string
		Body    string
		Text    string
	}{
		ToEmail: msg.ToEmail,
		ToName:  msg.ToName,
		Subject: msg.Subject,
		Body:    msg.HTML,
		Text:    msg.Text,
	})
	return nil
}
