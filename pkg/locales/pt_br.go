package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ptBR = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verifique seu e-mail",
		SalutationFormat:          "Olá %s,",
		DescriptionAboveBtn:       "Clique no link abaixo para confirmar seu endereço de e-mail.",
		VerifyEmail:               "Verificar e-mail",
		DescriptionBelowBtnFormat: "Se você não criou uma conta no %s, ignore este e-mail. Se não conseguir clicar no link acima, copie a URL e cole no navegador. O link de verificação de e-mail expira em %v minutos.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Redefina sua senha",
		SalutationFormat:          "Olá %s,",
		DescriptionAboveBtn:       "Recebemos recentemente uma solicitação para redefinir sua senha. Clique no link abaixo para redefini-la.",
		ResetPassword:             "Redefinir senha",
		DescriptionBelowBtnFormat: "Se você não solicitou a redefinição da senha, ignore este e-mail. Se não conseguir clicar no link acima, copie a URL e cole no navegador. O link de redefinição de senha expira em %v minutos.",
	},
}
