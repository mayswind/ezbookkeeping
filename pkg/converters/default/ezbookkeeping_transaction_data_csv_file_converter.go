package _default

// ezBookKeepingTransactionDataCSVFileConverter defines the structure of CSV file converter
type ezBookKeepingTransactionDataCSVFileConverter struct {
	ezBookKeepingTransactionDataPlainTextConverter
}

// Initialize an ezbookkeeping transaction data csv file converter singleton instance
var (
	EzBookKeepingTransactionDataCSVFileConverter = &ezBookKeepingTransactionDataCSVFileConverter{
		ezBookKeepingTransactionDataPlainTextConverter{
			columnSeparator: ",",
		},
	}
)
