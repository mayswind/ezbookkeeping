package converters

// ezBookKeepingTransactionDataTSVFileConverter defines the structure of TSV file converter
type ezBookKeepingTransactionDataTSVFileConverter struct {
	ezBookKeepingTransactionDataPlainTextConverter
}

// Initialize an ezbookkeeping transaction data tsv file converter singleton instance
var (
	EzBookKeepingTransactionDataTSVFileConverter = &ezBookKeepingTransactionDataTSVFileConverter{
		ezBookKeepingTransactionDataPlainTextConverter{
			DataTableTransactionDataConverter: DataTableTransactionDataConverter{
				dataColumnMapping:          ezbookkeepingDataColumnNameMapping,
				transactionTypeMapping:     ezbookkeepingTransactionTypeNameMapping,
				transactionTypeNameMapping: ezbookkeepingNameTransactionTypeMapping,
				columnSeparator:            "\t",
				lineSeparator:              "\n",
				geoLocationSeparator:       " ",
				transactionTagSeparator:    ";",
			},
			columns: ezbookkeepingDataColumns,
		},
	}
)
