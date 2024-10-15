package _default

// defaultTransactionDataTSVFileConverter defines the structure of ezbookkeeping default tsv file converter
type defaultTransactionDataTSVFileConverter struct {
	defaultTransactionDataPlainTextConverter
}

// Initialize an ezbookkeeping default transaction data tsv file converter singleton instance
var (
	DefaultTransactionDataTSVFileConverter = &defaultTransactionDataTSVFileConverter{
		defaultTransactionDataPlainTextConverter{
			columnSeparator: "\t",
		},
	}
)
