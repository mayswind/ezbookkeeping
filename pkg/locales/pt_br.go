package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ptBR = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_SPACE,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verificar Email",
		SalutationFormat:          "Olá %s,",
		DescriptionAboveBtn:       "Por favor, clique no link abaixo para confirmar o seu endereço de e-mail.",
		VerifyEmail:               "Verificar Email",
		DescriptionBelowBtnFormat: "Se você não se registrou para uma conta %s, basta ignorar este e-mail. Se não conseguir clicar no link acima, copie a URL acima e cole no seu navegador. O link para verificação de e-mail expirará após %v minutos.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Redefinir Sua Senha",
		SalutationFormat:          "Olá %s,",
		DescriptionAboveBtn:       "Recebemos recentemente uma solicitação para redefinir a sua senha. Você pode clicar no link abaixo para redefinir sua senha.",
		ResetPassword:             "Redefinir Senha",
		DescriptionBelowBtnFormat: "Se você não solicitou a redefinição de senha, basta ignorar este e-mail. Se não conseguir clicar no link acima, copie a URL acima e cole no seu navegador. O link de redefinição de senha expirará após %v minutos.",
	},
}
