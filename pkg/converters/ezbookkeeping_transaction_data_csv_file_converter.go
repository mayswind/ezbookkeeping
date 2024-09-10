package converters

// ezBookKeepingTransactionDataCSVFileConverter defines the structure of CSV file converter
type ezBookKeepingTransactionDataCSVFileConverter struct {
	ezBookKeepingTransactionDataPlainTextExporter
	ezBookKeepingTransactionDataPlainTextImporter
}

// Initialize an ezbookkeeping transaction data csv file converter singleton instance
var (
	EzBookKeepingTransactionDataCSVFileConverter = &ezBookKeepingTransactionDataCSVFileConverter{
		ezBookKeepingTransactionDataPlainTextExporter{
			DataTableTransactionDataExporter: DataTableTransactionDataExporter{
				dataColumnMapping:       ezbookkeepingDataColumnNameMapping,
				transactionTypeMapping:  ezbookkeepingTransactionTypeNameMapping,
				columnSeparator:         ",",
				lineSeparator:           "\n",
				geoLocationSeparator:    " ",
				transactionTagSeparator: ";",
			},
			columns: ezbookkeepingDataColumns,
		},
		ezBookKeepingTransactionDataPlainTextImporter{
			DataTableTransactionDataImporter: DataTableTransactionDataImporter{
				dataColumnMapping:          ezbookkeepingDataColumnNameMapping,
				transactionTypeNameMapping: ezbookkeepingNameTransactionTypeMapping,
				columnSeparator:            ",",
				lineSeparator:              "\n",
				geoLocationSeparator:       " ",
				transactionTagSeparator:    ";",
			},
		},
	}
)
