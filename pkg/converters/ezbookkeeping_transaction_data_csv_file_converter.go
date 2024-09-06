package converters

// ezBookKeepingTransactionDataCSVFileConverter defines the structure of CSV file converter
type ezBookKeepingTransactionDataCSVFileConverter struct {
	ezBookKeepingTransactionDataPlainTextConverter
}

// Initialize an ezbookkeeping transaction data csv file converter singleton instance
var (
	EzBookKeepingTransactionDataCSVFileConverter = &ezBookKeepingTransactionDataCSVFileConverter{
		ezBookKeepingTransactionDataPlainTextConverter{
			DataTableTransactionDataConverter: DataTableTransactionDataConverter{
				dataColumnMapping:          ezbookkeepingDataColumnNameMapping,
				transactionTypeMapping:     ezbookkeepingTransactionTypeNameMapping,
				transactionTypeNameMapping: ezbookkeepingNameTransactionTypeMapping,
				columnSeparator:            ",",
				lineSeparator:              "\n",
				geoLocationSeparator:       " ",
				transactionTagSeparator:    ";",
			},
			columns: ezbookkeepingDataColumns,
		},
	}
)
