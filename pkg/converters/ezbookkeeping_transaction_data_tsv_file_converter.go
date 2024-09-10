package converters

// ezBookKeepingTransactionDataTSVFileConverter defines the structure of TSV file converter
type ezBookKeepingTransactionDataTSVFileConverter struct {
	ezBookKeepingTransactionDataPlainTextExporter
	ezBookKeepingTransactionDataPlainTextImporter
}

// Initialize an ezbookkeeping transaction data tsv file converter singleton instance
var (
	EzBookKeepingTransactionDataTSVFileConverter = &ezBookKeepingTransactionDataTSVFileConverter{
		ezBookKeepingTransactionDataPlainTextExporter{
			DataTableTransactionDataExporter: DataTableTransactionDataExporter{
				dataColumnMapping:       ezbookkeepingDataColumnNameMapping,
				transactionTypeMapping:  ezbookkeepingTransactionTypeNameMapping,
				columnSeparator:         "\t",
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
				columnSeparator:            "\t",
				lineSeparator:              "\n",
				geoLocationSeparator:       " ",
				transactionTagSeparator:    ";",
			},
		},
	}
)
