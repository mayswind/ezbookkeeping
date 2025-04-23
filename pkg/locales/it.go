package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var it = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verifica il tuo indirizzo e-mail",
		SalutationFormat:          "Ciao %s,",
		DescriptionAboveBtn:       "Clicca il link sotto per confermare il tuo indirizzo e-amil",
		VerifyEmail:               "Verifica il tuo indirizzo e-mail",
		DescriptionBelowBtnFormat: "Se non hai creato un account %s, puoi ignorare questa e-mail. Se non riesci a cliccare il link, copia l'indirizzo URL qui sopra e incollalo nel tuo browser preferito. Il link di verifica scadrà tra %v minuti.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Reimposta password",
		SalutationFormat:          "Ciao %s,",
		DescriptionAboveBtn:       "Abbiamo ricevuto la tua richiesta di modifica della tua password. Puoi cliccare sul link qui sotto per impostare nuovamente la tua password.",
		ResetPassword:             "Reimposta password",
		DescriptionBelowBtnFormat: "Se non hai chiesto alcun cambio della password, puoi ignorare questa mail. Se non riesci a cliccare il link, copia l'indirizzo URL qui sopra e incollalo nel tuo browser preferito. Il link di verifica scadrà tra %v minuti.",
	},
}
