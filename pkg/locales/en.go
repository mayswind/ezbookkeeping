package locales

var en = &LocaleTextItems{
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Reset Your Password",
		SalutationFormat:          "Hi %s,",
		DescriptionAboveBtn:       "We recently received a request to reset your password. You can click the below link to reset your password.",
		ResetPassword:             "Reset Password",
		DescriptionBelowBtnFormat: "If you did not request to reset your password, please simply disregard this email. If you cannot click the above link, please copy the above url and paste it into your browser. The password reset link will be expired after %v minutes.",
	},
}
