package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var id = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Dompet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Verifikasi Email",
		SalutationFormat:          "Halo %s,",
		DescriptionAboveBtn:       "Silakan klik tautan di bawah ini untuk mengonfirmasi alamat email Anda.",
		VerifyEmail:               "Verifikasi Email",
		DescriptionBelowBtnFormat: "Jika Anda tidak mendaftar akun %s, silakan abaikan email ini. Jika Anda tidak dapat mengklik tautan di atas, silakan salin URL di atas dan tempelkan ke peramban Anda. Tautan verifikasi email akan kedaluwarsa setelah %v menit.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Atur Ulang Kata Sandi Anda",
		SalutationFormat:          "Halo %s,",
		DescriptionAboveBtn:       "Kami baru saja menerima permintaan untuk mengatur ulang kata sandi Anda. Anda dapat mengklik tautan di bawah ini untuk mengatur ulang kata sandi Anda.",
		ResetPassword:             "Atur Ulang Kata Sandi",
		DescriptionBelowBtnFormat: "Jika Anda tidak meminta untuk mengatur ulang kata sandi Anda, silakan abaikan email ini. Jika Anda tidak dapat mengklik tautan di atas, silakan salin URL di atas dan tempelkan ke peramban Anda. Tautan pengaturan ulang kata sandi akan kedaluwarsa setelah %v menit.",
	},
}
