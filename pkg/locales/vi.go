package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var vi = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA, 
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Ví WeChat",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Xác minh Email",
		SalutationFormat:          "Chào %s,",
		DescriptionAboveBtn:       "Vui lòng nhấp vào liên kết bên dưới để xác nhận địa chỉ email của bạn.",
		VerifyEmail:               "Xác minh Email",
		DescriptionBelowBtnFormat: "Nếu bạn không đăng ký tài khoản %s, vui lòng bỏ qua email này. Nếu bạn không thể nhấp vào liên kết trên, hãy sao chép và dán liên kết vào trình duyệt của bạn. Liên kết xác minh email sẽ hết hạn sau %v phút.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Đặt lại Mật khẩu",
		SalutationFormat:          "Chào %s,",
		DescriptionAboveBtn:       "Chúng tôi vừa nhận được yêu cầu đặt lại mật khẩu của bạn. Bạn có thể nhấp vào liên kết bên dưới để đặt lại mật khẩu.",
		ResetPassword:             "Đặt lại Mật khẩu",
		DescriptionBelowBtnFormat: "Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này. Nếu bạn không thể nhấp vào liên kết trên, hãy sao chép và dán liên kết vào trình duyệt của bạn. Liên kết đặt lại mật khẩu sẽ hết hạn sau %v phút.",
	},
}
