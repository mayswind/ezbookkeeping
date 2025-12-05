package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var kn = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "ಇಮೇಲ್ ದೃಢೀಕರಿಸಿ",
		SalutationFormat:          "ಹಲೋ %s,",
		DescriptionAboveBtn:       "ದಯವಿಟ್ಟು ನಿಮ್ಮ ಇಮೇಲ್ ವಿಳಾಸವನ್ನು ದೃಢೀಕರಿಸಲು ಕೆಳಗಿನ ಲಿಂಕ್ ಕ್ಲಿಕ್ ಮಾಡಿ.",
		VerifyEmail:               "ಇಮೇಲ್ ದೃಢೀಕರಿಸಿ",
		DescriptionBelowBtnFormat: "%s ಖಾತೆಗೆ ನೀವು ನೋಂದಾಯಿಸದಿದ್ದರೆ, ದಯವಿಟ್ಟು ಈ ಇಮೇಲ್ ಅನ್ನು ನಿರ್ಲಕ್ಷಿಸಿ. ಮೇಲಿನ ಲಿಂಕ್ ಕ್ಲಿಕ್ ಮಾಡಲು ಸಾಧ್ಯವಾಗದಿದ್ದರೆ, ಮೇಲಿನ URL ಅನ್ನು ನಕಲಿಸಿ ಮತ್ತು ನಿಮ್ಮ ಬ್ರೌಸರ್‌ನಲ್ಲಿ ಅಂಟಿಸಿ. ಇಮೇಲ್ ದೃಢೀಕರಣ ಲಿಂಕ್ %v ನಿಮಿಷಗಳ ನಂತರ ಅವಧಿ ಮುಗಿಯುತ್ತದೆ.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Reset Your Password",
		SalutationFormat:          "ಹಲೋ %s,",
		DescriptionAboveBtn:       "We recently received a request to reset your password. You can click the link below to reset your password.",
		ResetPassword:             "Reset Password",
		DescriptionBelowBtnFormat: "If you did not request to reset your password, please simply disregard this email. If you cannot click the link above, please copy the above url and paste it into your browser. The password reset link will be expired after %v minutes.",
	},
}
