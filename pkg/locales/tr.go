package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var tr = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Cüzdan",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "E-postayı Doğrula",
		SalutationFormat:          "Merhaba %s,",
		DescriptionAboveBtn:       "E-posta adresinizi onaylamak için lütfen aşağıdaki bağlantıya tıklayın.",
		VerifyEmail:               "E-postayı Doğrula",
		DescriptionBelowBtnFormat: "Eğer %s hesabı oluşturmadıysanız, lütfen bu e-postayı dikkate almayın. Eğer yukarıdaki bağlantıya tıklayamıyorsanız, lütfen adresi kopyalayıp tarayıcınıza yapıştırın. Doğrulama bağlantısının süresi %v dakika sonra dolacaktır.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Şifrenizi Sıfırlayın",
		SalutationFormat:          "Merhaba %s,",
		DescriptionAboveBtn:       "Yakın zamanda şifrenizi sıfırlama talebi aldık. Şifrenizi sıfırlamak için aşağıdaki bağlantıya tıklayabilirsiniz.",
		ResetPassword:             "Şifreyi Sıfırla",
		DescriptionBelowBtnFormat: "Eğer şifre sıfırlama talebinde bulunmadıysanız, lütfen bu e-postayı dikkate almayın. Eğer yukarıdaki bağlantıya tıklayamıyorsanız, lütfen adresi kopyalayıp tarayıcınıza yapıştırın. Şifre sıfırlama bağlantısının süresi %v dakika sonra dolacaktır.",
	},
}
