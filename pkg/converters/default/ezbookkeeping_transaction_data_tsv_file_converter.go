package _default

// ezBookKeepingTransactionDataTSVFileConverter defines the structure of TSV file converter
type ezBookKeepingTransactionDataTSVFileConverter struct {
	ezBookKeepingTransactionDataPlainTextConverter
}

// Initialize an ezbookkeeping transaction data tsv file converter singleton instance
var (
	EzBookKeepingTransactionDataTSVFileConverter = &ezBookKeepingTransactionDataTSVFileConverter{
		ezBookKeepingTransactionDataPlainTextConverter{
			columnSeparator: "\t",
		},
	}
)
