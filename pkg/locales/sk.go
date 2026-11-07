package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var sk = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_SPACE,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Overenie e-mailu",
		SalutationFormat:          "Dobrý deň, %s,",
		DescriptionAboveBtn:       "Kliknutím na odkaz nižšie potvrďte svoju e-mailovú adresu.",
		VerifyEmail:               "Overiť e-mail",
		DescriptionBelowBtnFormat: "Ak ste sa nezaregistrovali do %s, tento e-mail ignorujte. Ak na odkaz vyššie nemôžete kliknúť, skopírujte adresu a vložte ju do prehliadača. Platnosť odkazu na overenie vyprší o %v minút.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Obnova hesla",
		SalutationFormat:          "Dobrý deň, %s,",
		DescriptionAboveBtn:       "Nedávno sme dostali žiadosť o obnovu vášho hesla. Heslo môžete obnoviť kliknutím na odkaz nižšie.",
		ResetPassword:             "Obnoviť heslo",
		DescriptionBelowBtnFormat: "Ak ste o obnovu hesla nežiadali, tento e-mail ignorujte. Ak na odkaz vyššie nemôžete kliknúť, skopírujte adresu a vložte ju do prehliadača. Platnosť odkazu na obnovu hesla vyprší o %v minút.",
	},
}
