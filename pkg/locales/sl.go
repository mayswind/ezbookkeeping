package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var sl = &LocaleTextItems{
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
		Title:                     "Potrditev e-pošte",
		SalutationFormat:          "Zdravo %s,",
		DescriptionAboveBtn:       "Za potrditev svojega e-poštnega naslova kliknite spodnjo povezavo.",
		VerifyEmail:               "Potrdi e-pošto",
		DescriptionBelowBtnFormat: "Če se niste registrirali za %s račun, prosimo, da to e-pošto preprosto prezrete. Če ne morete klikniti zgornje povezave, kopirajte zgornji URL in ga prilepite v brskalnik. Povezava za potrditev e-pošte bo potekla po %v minutah.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Ponastavitev gesla",
		SalutationFormat:          "Zdravo %s,",
		DescriptionAboveBtn:       "Pred kratkim smo prejeli zahtevo za ponastavitev gesla. Za ponastavitev gesla lahko kliknete spodnjo povezavo.",
		ResetPassword:             "Ponastavi geslo",
		DescriptionBelowBtnFormat: "Če niste zahtevali ponastavitve gesla, prosimo, da to e-poštno sporočilo preprosto prezrete. Če ne morete klikniti zgornje povezave, kopirajte zgornji URL in ga prilepite v brskalnik. Povezava za ponastavitev gesla bo potekla po %v minutah.",
	},
}
