package mail

import (
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MailerContainer contains the current mailer
type MailerContainer struct {
	current Mailer
}

// Initialize a mailer container singleton instance
var (
	Container = &MailerContainer{}
)

// InitializeMailer initializes the current mailer according to the config
func InitializeMailer(config *settings.Config) error {
	if !config.EnableSMTP {
		Container.current = nil
		return nil
	}

	mailer, err := NewDefaultMailer(config.SMTPConfig)

	if err != nil {
		return err
	}

	Container.current = mailer
	return nil
}

// SendMail sends an email according to argument
func (m *MailerContainer) SendMail(message *MailMessage) error {
	if m.current == nil {
		return errs.ErrSMTPServerNotEnabled
	}

	return m.current.SendMail(message)
}
