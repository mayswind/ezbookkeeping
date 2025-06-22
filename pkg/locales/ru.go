package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ru = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_SPACE,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Подтвердите электронную почту",
		SalutationFormat:          "Здравствуйте %s,",
		DescriptionAboveBtn:       "Нажмите на ссылку ниже, чтобы подтвердить свой адрес электронной почты.",
		VerifyEmail:               "Подтвердить электронную почту",
		DescriptionBelowBtnFormat: "Если вы не регистрировали учетную запись %s, просто проигнорируйте это письмо. Если вы не можете нажать на ссылку выше, скопируйте указанный выше URL и вставьте его в свой браузер. Срок действия ссылки для проверки электронной почты истечет через %v минут.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Сброс пароля",
		SalutationFormat:          "Здравствуйте %s,",
		DescriptionAboveBtn:       "Недавно мы получили запрос на сброс вашего пароля. Нажмите на ссылку ниже, чтобы сбросить свой пароль.",
		ResetPassword:             "Сбросить пароль",
		DescriptionBelowBtnFormat: "Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо. Если вы не можете нажать на ссылку выше, скопируйте указанный выше URL и вставьте его в браузер. Ссылка для сброса пароля истечет через %v минут.",
	},
}
