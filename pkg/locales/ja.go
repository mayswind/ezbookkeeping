package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var ja = &LocaleTextItems{
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_DOT,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_COMMA,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Wallet",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "メールの確認",
		SalutationFormat:          "こんにちは%s,",
		DescriptionAboveBtn:       "次のリンクをクリックしてメールアドレスを確認してください。",
		VerifyEmail:               "メールを確認",
		DescriptionBelowBtnFormat: "%sアカウントに登録していない場合はこのメールを無視してください。上記のリンクをクリックできない場合は上記のURLをコピーしてブラウザに貼り付けてください。メールの確認リンクは%v分後に期限切れになります。",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "パスワードのリセット",
		SalutationFormat:          "こんにちは%s,",
		DescriptionAboveBtn:       "パスワードリセットのリクエストを受け取りました。次のリンクをクリックしてパスワードをリセットしてください。",
		ResetPassword:             "パスワードをリセット",
		DescriptionBelowBtnFormat: "パスワードのリセットをリクエストしていない場合はこのメールを無視してください。上記のリンクをクリックできない場合は、上記のURLをコピーしてブラウザに貼り付けてください。パスワードリセットのリンクは%v分後に期限切れになります。",
	},
}
