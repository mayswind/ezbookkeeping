package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var th = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "ยืนยันอีเมล",
		SalutationFormat:          "สวัสดี %s,",
		DescriptionAboveBtn:       "โปรดคลิกที่ลิงค์ด้านล่างเพื่อยืนยันที่อยู่อีเมลของคุณ",
		VerifyEmail:               "ยืนยันอีเมล",
		DescriptionBelowBtnFormat: "หากคุณไม่ได้ลงทะเบียนสำหรับบัญชี %s โปรดละเว้นอีเมลนี้ หากคุณไม่สามารถคลิกลิงก์ด้านบน โปรดคัดลอก URL ด้านบนและวางลงในเบราว์เซอร์ของคุณ ลิงก์ยืนยันอีเมลจะหมดอายุหลังจาก %v นาที",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "รีเซ็ตรหัสผ่านใหม่",
		SalutationFormat:          "สวัสดี %s,",
		DescriptionAboveBtn:       "เมื่อเร็ว ๆ นี้เราได้รับการร้องขอให้รีเซ็ตรหัสผ่านของคุณ คุณสามารถคลิกลิงก์ด้านล่างเพื่อรีเซ็ตรหัสผ่านของคุณ",
		ResetPassword:             "ตั้งรหัสผ่านใหม่",
		DescriptionBelowBtnFormat: "หากคุณไม่ได้ร้องขอให้รีเซ็ตรหัสผ่าน โปรดละเว้นอีเมลนี้ หากคุณไม่สามารถคลิกลิงก์ด้านบน โปรดคัดลอก URL ด้านบนและวางลงในเบราว์เซอร์ของคุณ ลิงก์รีเซ็ตรหัสผ่านจะหมดอายุหลังจาก %v นาที",
	},
}
