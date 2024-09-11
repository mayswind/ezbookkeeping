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
				geoLocationSeparator:    " ",
				transactionTagSeparator: ";",
			},
			columns:         ezbookkeepingDataColumns,
			columnSeparator: ",",
			lineSeparator:   "\n",
		},
		ezBookKeepingTransactionDataPlainTextImporter{
			DataTableTransactionDataImporter: DataTableTransactionDataImporter{
				dataColumnMapping:       ezbookkeepingDataColumnNameMapping,
				transactionTypeMapping:  ezbookkeepingTransactionTypeNameMapping,
				geoLocationSeparator:    " ",
				transactionTagSeparator: ";",
			},
			columnSeparator: ",",
			lineSeparator:   "\n",
		},
	}
)
