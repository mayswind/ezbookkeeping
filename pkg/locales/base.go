package locales

import "github.com/mayswind/ezbookkeeping/pkg/models"

// LocaleTextItems represents all text items need to be translated
type LocaleTextItems struct {
	DefaultTypes                *DefaultTypes
	VerifyEmailTextItems        *VerifyEmailTextItems
	ForgetPasswordMailTextItems *ForgetPasswordMailTextItems
}

type DefaultTypes struct {
	DecimalSeparator    models.DecimalSeparator
	DigitGroupingSymbol models.DigitGroupingSymbol
}

// VerifyEmailTextItems represents text items need to be translated in verify mail
type VerifyEmailTextItems struct {
	Title                     string
	SalutationFormat          string
	DescriptionAboveBtn       string
	VerifyEmail               string
	DescriptionBelowBtnFormat string
}

// ForgetPasswordMailTextItems represents text items need to be translated in forget password mail
type ForgetPasswordMailTextItems struct {
	Title                     string
	SalutationFormat          string
	DescriptionAboveBtn       string
	ResetPassword             string
	DescriptionBelowBtnFormat string
}
