package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var en = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verify Email",
		SalutationFormat:          "Hi %s,",
		DescriptionAboveBtn:       "Please click the link below to confirm your email address.",
		VerifyEmail:               "Verify Email",
		DescriptionBelowBtnFormat: "If you did not sign up for %s account, please simply disregard this email. If you cannot click the link above, please copy the above url and paste it into your browser. The verify email link will be expired after %v minutes.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Reset Your Password",
		SalutationFormat:          "Hi %s,",
		DescriptionAboveBtn:       "We recently received a request to reset your password. You can click the link below to reset your password.",
		ResetPassword:             "Reset Password",
		DescriptionBelowBtnFormat: "If you did not request to reset your password, please simply disregard this email. If you cannot click the link above, please copy the above url and paste it into your browser. The password reset link will be expired after %v minutes.",
	},
}
