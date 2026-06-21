package templates

type KnownTemplate string

// Known templates
const (
	TEMPLATE_VERIFY_EMAIL                      KnownTemplate = "email/verify_email"
	TEMPLATE_PASSWORD_RESET                    KnownTemplate = "email/password_reset"
	SYSTEM_PROMPT_RECEIPT_IMAGE_RECOGNITION    KnownTemplate = "prompt/receipt_image_recognition"
	SYSTEM_PROMPT_TRANSACTION_TEXT_RECOGNITION KnownTemplate = "prompt/transaction_text_recognition"
)
