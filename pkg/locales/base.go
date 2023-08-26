package locales

// LocaleTextItems represents all text items need to be translated
type LocaleTextItems struct {
	ForgetPasswordMailTextItems *ForgetPasswordMailTextItems
}

// ForgetPasswordMailTextItems represents text items need to be translated in forget password mail
type ForgetPasswordMailTextItems struct {
	Title                     string
	SalutationFormat          string
	DescriptionAboveBtn       string
	ResetPassword             string
	DescriptionBelowBtnFormat string
}
