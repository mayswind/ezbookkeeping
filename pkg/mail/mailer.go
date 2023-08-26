package mail

// Mailer is email sender interface
type Mailer interface {
	SendMail(message *MailMessage) error
}
