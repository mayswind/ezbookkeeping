package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var hi = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "वॉलेट",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "ईमेल सत्यापित करें",
		SalutationFormat:          "नमस्ते %s,",
		DescriptionAboveBtn:       "अपना ईमेल पता सत्यापित करने के लिए कृपया नीचे दिए गए लिंक पर क्लिक करें।",
		VerifyEmail:               "ईमेल सत्यापित करें",
		DescriptionBelowBtnFormat: "अगर आपने %s खाता नहीं बनाया है, तो कृपया इस ईमेल को नज़रअंदाज़ कर दें। अगर आप ऊपर दिए गए लिंक पर क्लिक नहीं कर पा रहे हैं, तो उस URL को कॉपी करके अपने ब्राउज़र में पेस्ट करें। ईमेल सत्यापन लिंक %v मिनट बाद समाप्त हो जाएगा।",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "अपना पासवर्ड रीसेट करें",
		SalutationFormat:          "नमस्ते %s,",
		DescriptionAboveBtn:       "हमें हाल ही में आपका पासवर्ड रीसेट करने का अनुरोध मिला है। आप नीचे दिए गए लिंक पर क्लिक करके अपना पासवर्ड रीसेट कर सकते हैं।",
		ResetPassword:             "पासवर्ड रीसेट करें",
		DescriptionBelowBtnFormat: "अगर आपने अपना पासवर्ड रीसेट करने का अनुरोध नहीं किया है, तो कृपया इस ईमेल को नज़रअंदाज़ कर दें। अगर आप ऊपर दिए गए लिंक पर क्लिक नहीं कर पा रहे हैं, तो उस URL को कॉपी करके अपने ब्राउज़र में पेस्ट करें। पासवर्ड रीसेट लिंक %v मिनट बाद समाप्त हो जाएगा।",
	},
}
