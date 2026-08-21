package locales

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

var el = &LocaleTextItems{
	GlobalTextItems: &GlobalTextItems{
		AppName: "ezBookkeeping",
	},
	DefaultTypes: &DefaultTypes{
		DecimalSeparator:    core.DECIMAL_SEPARATOR_COMMA,
		DigitGroupingSymbol: core.DIGIT_GROUPING_SYMBOL_DOT,
	},
	DataConverterTextItems: &DataConverterTextItems{
		Alipay:       "Alipay",
		WeChatWallet: "Πορτοφόλι",
	},
	VerifyEmailTextItems: &VerifyEmailTextItems{
		Title:                     "Επαλήθευση email",
		SalutationFormat:          "Γεια σας %s,",
		DescriptionAboveBtn:       "Κάντε κλικ στον παρακάτω σύνδεσμο για να επιβεβαιώσετε τη διεύθυνση email σας.",
		VerifyEmail:               "Επαλήθευση email",
		DescriptionBelowBtnFormat: "Αν δεν δημιουργήσατε εσείς λογαριασμό στο %s, απλώς αγνοήστε αυτό το email. Αν δεν μπορείτε να κάνετε κλικ στον παραπάνω σύνδεσμο, αντιγράψτε τη διεύθυνση και επικολλήστε την στο πρόγραμμα περιήγησής σας. Ο σύνδεσμος επαλήθευσης email θα λήξει μετά από %v λεπτά.",
	},
	ForgetPasswordMailTextItems: &ForgetPasswordMailTextItems{
		Title:                     "Επαναφορά κωδικού πρόσβασης",
		SalutationFormat:          "Γεια σας %s,",
		DescriptionAboveBtn:       "Λάβαμε πρόσφατα αίτημα επαναφοράς του κωδικού πρόσβασής σας. Μπορείτε να κάνετε κλικ στον παρακάτω σύνδεσμο για να επαναφέρετε τον κωδικό πρόσβασής σας.",
		ResetPassword:             "Επαναφορά κωδικού πρόσβασης",
		DescriptionBelowBtnFormat: "Αν δεν ζητήσατε εσείς επαναφορά του κωδικού πρόσβασής σας, απλώς αγνοήστε αυτό το email. Αν δεν μπορείτε να κάνετε κλικ στον παραπάνω σύνδεσμο, αντιγράψτε τη διεύθυνση και επικολλήστε την στο πρόγραμμα περιήγησής σας. Ο σύνδεσμος επαναφοράς κωδικού πρόσβασης θα λήξει μετά από %v λεπτά.",
	},
}
