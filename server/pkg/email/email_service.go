package email

import (
	"context"

	"github.com/spf13/viper"
)

type EmailService struct {
	mailSender MailSender
}

func NewEmailService(mailSender MailSender) *EmailService {
	if mailSender == nil {
		mailSender = NewSMTPSender(SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		})
	}
	return &EmailService{mailSender: mailSender}
}

func (s *EmailService) SendEmail(ctx context.Context, email string, subject string, body string) error {
	return s.mailSender.Send(ctx, email, subject, body)
}

func (s *EmailService) SendEmailInBackground(email string, subject string, body string) {
	go func() {
		bgCtx := context.Background()
		s.mailSender.Send(bgCtx, email, subject, body)
	}()
}
