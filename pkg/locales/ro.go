package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ro = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Portofel",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verificare Email",
		SalutationFormat:          "Bună ziua, %s,",
		DescriptionAboveBtn:       "Vă rugăm să faceți clic pe linkul de mai jos pentru a confirma adresa de email.",
		VerifyEmail:               "Verificare Email",
		DescriptionBelowBtnFormat: "Dacă nu v-ați înregistrat pentru un cont %s, ignorați acest email. Dacă nu puteți accesa linkul de mai sus, copiați adresa URL și inserați-o în browser. Linkul de verificare va expira după %v minute.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Resetare Parolă",
		SalutationFormat:          "Bună ziua, %s,",
		DescriptionAboveBtn:       "Am primit recent o solicitare de resetare a parolei. Puteți face clic pe linkul de mai jos pentru a vă reseta parola.",
		ResetPassword:             "Resetare Parolă",
		DescriptionBelowBtnFormat: "Dacă nu ați solicitat resetarea parolei, ignorați acest email. Dacă nu puteți accesa linkul de mai sus, copiați adresa URL și inserați-o în browser. Linkul de resetare va expira după %v minute.",
	},
}
