package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var de = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "E-Mail verifizieren",
		SalutationFormat:          "Hallo %s,",
		DescriptionAboveBtn:       "Bitte klicken Sie auf den untenstehenden Link, um Ihre E-Mail-Adresse zu bestätigen.",
		VerifyEmail:               "E-Mail verifizieren",
		DescriptionBelowBtnFormat: "Wenn Sie kein %s Konto erstellt haben, ignorieren Sie bitte diese E-Mail. Wenn Sie den obigen Link nicht anklicken können, kopieren Sie bitte die obige URL und fügen Sie sie in Ihren Browser ein. Der Verifizierungslink wird nach %v Minuten ablaufen.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Passwort zurücksetzen",
		SalutationFormat:          "Hallo %s,",
		DescriptionAboveBtn:       "Wir haben kürzlich eine Anfrage zum Zurücksetzen Ihres Passworts erhalten. Sie können auf den untenstehenden Link klicken, um Ihr Passwort zurückzusetzen.",
		ResetPassword:             "Passwort zurücksetzen",
		DescriptionBelowBtnFormat: "Wenn Sie nicht angefordert haben, Ihr Passwort zurückzusetzen, ignorieren Sie bitte diese E-Mail. Wenn Sie den obigen Link nicht anklicken können, kopieren Sie bitte die obige URL und fügen Sie sie in Ihren Browser ein. Der Link zum Zurücksetzen des Passworts wird nach %v Minuten ablaufen.",
	},
}
