package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ta = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "மின்னஞ்சல் சரிபார்ப்பு",
		SalutationFormat:          "வணக்கம் %s,",
		DescriptionAboveBtn:       "உங்கள் மின்னஞ்சல் முகவரியை உறுதிப்படுத்த கீழே உள்ள இணைப்பைக் கிளிக் செய்யவும்.",
		VerifyEmail:               "மின்னஞ்சலை சரிபார்க்கவும்",
		DescriptionBelowBtnFormat: "நீங்கள் %s கணக்கிற்கு பதிவு செய்யவில்லை என்றால், இந்த மின்னஞ்சலை புறக்கணிக்கவும். மேலே உள்ள இணைப்பைக் கிளிக் செய்ய முடியவில்லை என்றால், மேலே உள்ள URL ஐ நகலெடுத்து உங்கள் உலாவியில் ஒட்டவும். மின்னஞ்சல் சரிபார்ப்பு இணைப்பு %v நிமிடங்களுக்குப் பிறகு காலாவதியாகும்.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "உங்கள் கடவுச்சொல்லை மீட்டமைக்கவும்",
		SalutationFormat:          "வணக்கம் %s,",
		DescriptionAboveBtn:       "உங்கள் கடவுச்சொல்லை மீட்டமைக்க சமீபத்தில் கோரிக்கை பெற்றோம். உங்கள் கடவுச்சொல்லை மீட்டமைக்க கீழே உள்ள இணைப்பைக் கிளிக் செய்யவும்.",
		ResetPassword:             "கடவுச்சொல்லை மீட்டமை",
		DescriptionBelowBtnFormat: "உங்கள் கடவுச்சொல்லை மீட்டமைக்க நீங்கள் கோரவில்லை என்றால், இந்த மின்னஞ்சலை புறக்கணிக்கவும். மேலே உள்ள இணைப்பைக் கிளிக் செய்ய முடியவில்லை என்றால், மேலே உள்ள URL ஐ நகலெடுத்து உங்கள் உலாவியில் ஒட்டவும். கடவுச்சொல் மீட்டமைப்பு இணைப்பு %v நிமிடங்களுக்குப் பிறகு காலாவதியாகும்.",
	},
}
