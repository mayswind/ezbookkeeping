package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var uk = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Підтвердіть електронну пошту",
		SalutationFormat:          "Вітаємо, %s!",
		DescriptionAboveBtn:       "Натисніть на посилання нижче, щоб підтвердити вашу електронну адресу.",
		VerifyEmail:               "Підтвердити електронну пошту",
		DescriptionBelowBtnFormat: "Якщо ви не створювали обліковий запис %s, просто проігноруйте цей лист. Якщо ви не можете натиснути на посилання вище, скопіюйте вказану URL-адресу та вставте її у свій браузер. Посилання для підтвердження електронної пошти буде дійсне протягом %v хвилин.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Скидання пароля",
		SalutationFormat:          "Вітаємо, %s!",
		DescriptionAboveBtn:       "Нещодавно ми отримали запит на скидання вашого пароля. Натисніть на посилання нижче, щоб скинути пароль.",
		ResetPassword:             "Скинути пароль",
		DescriptionBelowBtnFormat: "Якщо ви не надсилали запит на скидання пароля, просто проігноруйте цей лист. Якщо ви не можете натиснути на посилання вище, скопіюйте вказану URL-адресу та вставте її у свій браузер. Посилання для скидання пароля буде дійсне протягом %v хвилин.",
	},
}
