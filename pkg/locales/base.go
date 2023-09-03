package locales

// LocaleTextItems represents all text items need to be translated
type LocaleTextItems struct {
	VerifyEmailTextItems        *VerifyEmailTextItems
	ForgetPasswordMailTextItems *ForgetPasswordMailTextItems
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
