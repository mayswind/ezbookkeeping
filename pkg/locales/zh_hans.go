package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var zhHans = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "验证邮箱",
		SalutationFormat:          "%s 您好，",
		DescriptionAboveBtn:       "请点击下方的链接确认您的邮箱地址。",
		VerifyEmail:               "验证邮箱",
		DescriptionBelowBtnFormat: "如果您没有注册 %s 账户，请直接忽略本邮件。如果您无法点击上述链接，请复制下方的地址然后在您的浏览器中粘贴。邮箱验证链接将在 %v 分钟后过期。",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "重置密码",
		SalutationFormat:          "%s 您好，",
		DescriptionAboveBtn:       "我们刚才收到重置您密码的请求。您可以点击下方链接重置您的密码。",
		ResetPassword:             "重置密码",
		DescriptionBelowBtnFormat: "如果您没有请求重置密码，请直接忽略本邮件。如果您无法点击上述链接，请复制下方的地址然后在您的浏览器中粘贴。重置密码链接将在 %v 分钟后过期。",
	},
}
