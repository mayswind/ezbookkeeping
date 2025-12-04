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
		WeChatWallet: "ವಾಲೆಟ್",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "ಇಮೇಲ್ ಪರಿಶೀಲನೆ",
		SalutationFormat:          "ನಮಸ್ಕಾರ %s,",
		DescriptionAboveBtn:       "ನಿಮ್ಮ ಇಮೇಲ್ ವಿಳಾಸ ಪರಿಶೀಲಿಸಲು ಕೆಳಗಿನ ಲಿಂಕ್ ಕ್ಲಿಕ್ ಮಾಡಿ.",
		VerifyEmail:               "ಇಮೇಲ್ ಪರಿಶೀಲನೆ",
		DescriptionBelowBtnFormat: "ನೀವು %s ಖಾತೆಗೆ ಸೈನ್ ಅಪ್ ಮಾಡದಿದ್ದರೆ ಈ ಇಮೇಲ್ ನಿರ್ಲಕ್ಷಿಸಿ. ಮೇಲಿನ ಲಿಂಕ್ ಕ್ಲಿಕ್ ಮಾಡಲು ಸಾಧ್ಯವಾಗದಿದ್ದರೆ, ಮೇಲಿನ ಯುಆರ್ಎಲ್ ಕಾಪಿ ಮಾಡಿ ಮತ್ತು ನಿಮ್ಮ ಬ್ರೌಜರ್‍ನಲ್ಲಿ ಪೇಸ್ಟ್ ಮಾಡಿ. ಇಮೇಲ್ ಪರಿಶೀಲನೆ ಲಿಂಕ್ %v ನಿಮಿಷಗಳ ನಂತರ ಅವಧಿಯನ್ನು ಮೀರುತ್ತದೆ.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "ಪಾಸ್‍ವರ್ಡ್ ಮರುಹೊಂದಾಯಿಸಿ",
		SalutationFormat:          "ನಮಸ್ಕಾರ %s,",
		DescriptionAboveBtn:       "ನಿಮ್ಮ ಪಾಸ್‍ವರ್ಡ್ ಮರುಹೊಂದಾಯಿಸುವ ವಿನಂತಿ ಇತ್ತೀಚೆ ಸ್ವೀಕರಿಸಲಾಯಿತು. ನಿಮ್ಮ ಪಾಸ್‍ವರ್ಡ್ ಮರುಹೊಂದಾಯಿಸಲು ಕೆಳಗಿನ ಲಿಂಕ್ ಕ್ಲಿಕ್ ಮಾಡಿ.",
		ResetPassword:             "ಪಾಸ್‍ವರ್ಡ್ ಮರುಹೊಂದಾಯಿಸಿ",
		DescriptionBelowBtnFormat: "ನೀವು ಪಾಸ್‍ವರ್ಡ್ ಮರುಹೊಂದಾಯಿಸುವ ವಿನಂತಿ ಮಾಡದಿದ್ದರೆ ಈ ಇಮೇಲ್ ನಿರ್ಲಕ್ಷಿಸಿ. ಮೇಲಿನ ಲಿಂಕ್ ಕ್ಲಿಕ್ ಮಾಡಲು ಸಾಧ್ಯವಾಗದಿದ್ದರೆ, ಮೇಲಿನ ಯುಆರ್ಎಲ್ ಕಾಪಿ ಮಾಡಿ ಮತ್ತು ನಿಮ್ಮ ಬ್ರೌಜರ್‍ನಲ್ಲಿ ಪೇಸ್ಟ್ ಮಾಡಿ. ಪಾಸ್‍ವರ್ಡ್ ಮರುಹೊಂದಾಯಿಸುವ ಲಿಂಕ್ %v ನಿಮಿಷಗಳ ನಂತರ ಅವಧಿಯನ್ನು ಮೀರುತ್ತದೆ.",
	},
}
