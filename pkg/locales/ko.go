package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ko = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "이메일 인증",
		SalutationFormat:          "안녕하세요 %s님,",
		DescriptionAboveBtn:       "이메일 주소를 확인하려면 아래 링크를 클릭해주세요.",
		VerifyEmail:               "이메일 인증",
		DescriptionBelowBtnFormat: "%s 계정에 가입하지 않으셨다면 이 이메일을 무시해주세요. 위 링크를 클릭할 수 없는 경우, 위 URL을 복사하여 브라우저에 붙여넣어 주세요. 이메일 인증 링크는 %v분 후에 만료됩니다.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "비밀번호 재설정",
		SalutationFormat:          "안녕하세요 %s님,",
		DescriptionAboveBtn:       "비밀번호 재설정 요청이 있었습니다. 아래 링크를 클릭하시면 비밀번호를 재설정할 수 있습니다.",
		ResetPassword:             "비밀번호 재설정",
		DescriptionBelowBtnFormat: "비밀번호 재설정을 요청하지 않으셨다면 이 이메일을 무시해주세요. 위 링크를 클릭할 수 없는 경우, 위 URL을 복사하여 브라우저에 붙여넣어 주세요. 비밀번호 재설정 링크는 %v분 후에 만료됩니다.",
	},
}
