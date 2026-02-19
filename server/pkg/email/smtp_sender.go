package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig holds SMTP connection settings.
type SMTPConfig struct {
	Host     string // e.g. "smtp.gmail.com"
	Port     int    // e.g. 587
	Username string
	Password string
	From     string // sender email address
}

// SMTPSender sends email via SMTP.
type SMTPSender struct {
	config SMTPConfig
}

// NewSMTPSender returns a MailSender that uses SMTP.
func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{config: config}
}

// Send implements MailSender.
func (s *SMTPSender) Send(ctx context.Context, to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"From: " + s.config.From + "\r\n" +
			"Subject: " + encodeSubject(subject) + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	return smtp.SendMail(addr, auth, s.config.From, []string{to}, msg)
}

// encodeSubject encodes subject for MIME (RFC 2047 style optional; simple ASCII pass-through).
func encodeSubject(s string) string {
	if strings.ContainsAny(s, "\r\n") {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\r", ""), "\n", " ")
	}
	return s
}
