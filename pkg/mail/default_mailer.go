package mail

import (
	"crypto/tls"
	"net"

	"gopkg.in/mail.v2"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// DefaultMailer represents default mailer
type DefaultMailer struct {
	dialer      *mail.Dialer
	fromAddress string
}

// NewDefaultMailer returns a new default mailer
func NewDefaultMailer(smtpConfig *settings.SmtpConfig) (*DefaultMailer, error) {
	host, portStr, err := net.SplitHostPort(smtpConfig.SmtpHost)

	if err != nil {
		return nil, errs.ErrSmtpServerHostInvalid
	}

	port, err := utils.StringToInt(portStr)

	if err != nil {
		return nil, errs.ErrSmtpServerHostInvalid
	}

	dialer := mail.NewDialer(host, port, smtpConfig.SmtpUser, smtpConfig.SmtpPasswd)
	dialer.TLSConfig = &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: smtpConfig.SmtpSkipTLSVerify,
	}

	mailer := &DefaultMailer{
		dialer:      dialer,
		fromAddress: smtpConfig.FromAddress,
	}

	return mailer, nil
}

// SendMail sends an email according to argument
func (m *DefaultMailer) SendMail(message *MailMessage) error {
	if m.dialer == nil {
		return errs.ErrSmtpServerNotEnabled
	}

	mailMessage := mail.NewMessage()
	mailMessage.SetHeader("From", m.fromAddress)
	mailMessage.SetHeader("To", message.To)
	mailMessage.SetHeader("Subject", message.Subject)
	mailMessage.SetBody("text/html", message.Body)

	err := m.dialer.DialAndSend(mailMessage)

	return err
}
