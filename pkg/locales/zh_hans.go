package locales

var zhHans = &LocaleTextItems{
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "重置密码",
		SalutationFormat:          "%s 你好，",
		DescriptionAboveBtn:       "我们刚才收到重置您密码的请求。您可以点击下方链接重置您的密码。",
		ResetPassword:             "重置密码",
		DescriptionBelowBtnFormat: "如果您没有请求重置密码，请直接忽略本邮件。如果您无法点击上述链接，请复制下方的地址然后在您的浏览器中粘贴。重置密码链接将在 %v 分钟后过期。",
	},
}
