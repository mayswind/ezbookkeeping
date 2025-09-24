package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var fr = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_SPACE,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Vérifier l'e-mail",
		SalutationFormat:          "Bonjour %s,",
		DescriptionAboveBtn:       "Cliquez sur le lien ci-dessous pour confirmer votre adresse e-mail.",
		VerifyEmail:               "Vérifier l'e-mail",
		DescriptionBelowBtnFormat: "Si vous n'avez pas créé de compte %s, vous pouvez ignorer cet e-mail. Si vous ne pouvez pas cliquer sur le lien ci-dessus, copiez l'URL ci-dessus et collez-la dans votre navigateur. Le lien de vérification expire après %v minutes.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Réinitialiser le mot de passe",
		SalutationFormat:          "Bonjour %s,",
		DescriptionAboveBtn:       "Nous avons récemment reçu une demande de réinitialisation de votre mot de passe. Cliquez sur le lien ci-dessous pour réinitialiser votre mot de passe.",
		ResetPassword:             "Réinitialiser le mot de passe",
		DescriptionBelowBtnFormat: "Si vous n'avez pas demandé la réinitialisation de votre mot de passe, vous pouvez ignorer cet e-mail. Si vous ne pouvez pas cliquer sur le lien ci-dessus, copiez l'URL ci-dessus et collez-la dans votre navigateur. Le lien de réinitialisation du mot de passe expire après %v minutes.",
	},
}
