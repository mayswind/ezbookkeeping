package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var es = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verifique su cuenta de correo",
		SalutationFormat:          "Hola %s,",
		DescriptionAboveBtn:       "Por favor, haga click en el siguiente enlace para confirmar su cuenta de correo.",
		VerifyEmail:               "Verificar correo",
		DescriptionBelowBtnFormat: "Si no registró una cuenta de %s, simplemente haga caso omiso a este correo. Si no puede hacer click en el link de verificación, copie la url arriba mostrada y péguela en su navegador. El enlace de verificación de correo expira pasados %v minutos.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Restablezca su Contraseña",
		SalutationFormat:          "Hola %s,",
		DescriptionAboveBtn:       "Hemos recibido una solicitud para restablecer su contraseña. Puede hacer click en el siguiente link para restablecer su contraseña.",
		ResetPassword:             "Restablecer Contraseña",
		DescriptionBelowBtnFormat: "Si no solicitó un restablecimiento de contraseña, simplemente descarte este correo. Si no puede hacer click en el link anterior, copie la url arriba mostrada y péguela en su navegadror. El enlace de restablecimiento de contraseña expira pasados %v minutos.",
	},
}
