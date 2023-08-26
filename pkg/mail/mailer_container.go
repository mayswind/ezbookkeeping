package mail

import (
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MailerContainer contains the current mailer
type MailerContainer struct {
	Current Mailer
}

// Initialize a mailer container singleton instance
var (
	Container = &MailerContainer{}
)

// InitializeMailer initializes the current mailer according to the config
func InitializeMailer(config *settings.Config) error {
	if !config.EnableSmtp {
		Container.Current = nil
		return nil
	}

	mailer, err := NewDefaultMailer(config.SmtpConfig)

	if err != nil {
		return err
	}

	Container.Current = mailer
	return nil
}

// SendMail sends an email according to argument
func (u *MailerContainer) SendMail(message *MailMessage) error {
	return u.Current.SendMail(message)
}
