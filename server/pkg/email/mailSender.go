package email

import "context"

type MailSender interface {
	Send(ctx context.Context, email string, subject string, body string) error
}
