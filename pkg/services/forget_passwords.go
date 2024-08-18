package services

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/locales"
	"github.com/mayswind/ezbookkeeping/pkg/mail"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/templates"
)

const passwordResetUrlFormat = "%sdesktop/#/resetpassword?token=%s"

// ForgetPasswordService represents forget password service
type ForgetPasswordService struct {
	ServiceUsingConfig
	ServiceUsingMailer
}

// Initialize a forget password service singleton instance
var (
	ForgetPasswords = &ForgetPasswordService{
		ServiceUsingConfig: ServiceUsingConfig{
			container: settings.Container,
		},
		ServiceUsingMailer: ServiceUsingMailer{
			container: mail.Container,
		},
	}
)

// SendPasswordResetEmail sends password reset email according to specified parameters
func (s *ForgetPasswordService) SendPasswordResetEmail(c core.Context, user *models.User, passwordResetToken string, backupLocale string) error {
	if !s.CurrentConfig().EnableSMTP {
		return errs.ErrSMTPServerNotEnabled
	}

	locale := user.Language

	if locale == "" {
		locale = backupLocale
	}

	localeTextItems := locales.GetLocaleTextItems(locale)
	forgetPasswordTextItems := localeTextItems.ForgetPasswordMailTextItems

	expireTimeInMinutes := s.CurrentConfig().PasswordResetTokenExpiredTimeDuration.Minutes()
	passwordResetUrl := fmt.Sprintf(passwordResetUrlFormat, s.CurrentConfig().RootUrl, url.QueryEscape(passwordResetToken))

	tmpl, err := templates.GetTemplate(templates.TEMPLATE_PASSWORD_RESET)

	if err != nil {
		return err
	}

	templateParams := map[string]any{
		"AppName": s.CurrentConfig().AppName,
		"ForgetPasswordMail": map[string]any{
			"Title":               forgetPasswordTextItems.Title,
			"Salutation":          fmt.Sprintf(forgetPasswordTextItems.SalutationFormat, user.Nickname),
			"DescriptionAboveBtn": forgetPasswordTextItems.DescriptionAboveBtn,
			"ResetPasswordUrl":    passwordResetUrl,
			"ResetPassword":       forgetPasswordTextItems.ResetPassword,
			"DescriptionBelowBtn": fmt.Sprintf(forgetPasswordTextItems.DescriptionBelowBtnFormat, expireTimeInMinutes),
		},
	}

	var bodyBuffer bytes.Buffer
	err = tmpl.Execute(&bodyBuffer, templateParams)

	if err != nil {
		return err
	}

	message := &mail.MailMessage{
		To:      user.Email,
		Subject: forgetPasswordTextItems.Title,
		Body:    bodyBuffer.String(),
	}

	err = s.SendMail(message)

	return err
}
