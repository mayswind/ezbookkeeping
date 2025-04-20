package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var zhHant = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "支付寶",
		WeChatWallet: "零錢",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "驗證郵箱",
		SalutationFormat:          "%s 您好，",
		DescriptionAboveBtn:       "請點擊下方的連結確認您的郵箱地址。",
		VerifyEmail:               "驗證郵箱",
		DescriptionBelowBtnFormat: "如果您沒有註冊 %s 帳戶，請直接忽略本郵件。如果您無法點擊上述連結，請複製下方的地址然後在您的瀏覽器中貼上。郵箱驗證連結將在 %v 分鐘後過期。",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "重設密碼",
		SalutationFormat:          "%s 您好，",
		DescriptionAboveBtn:       "我們剛才收到重設您密碼的請求。您可以點擊下方連結重設您的密碼。",
		ResetPassword:             "重設密碼",
		DescriptionBelowBtnFormat: "如果您沒有請求重設密碼，請直接忽略本郵件。如果您無法點擊上述連結，請複製下方的地址然後在您的瀏覽器中貼上。重設密碼連結將在 %v 分鐘後過期。",
	},
}
