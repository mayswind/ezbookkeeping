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
				geoLocationSeparator:    " ",
				transactionTagSeparator: ";",
			},
			columns:         ezbookkeepingDataColumns,
			columnSeparator: "\t",
			lineSeparator:   "\n",
		},
		ezBookKeepingTransactionDataPlainTextImporter{
			DataTableTransactionDataImporter: DataTableTransactionDataImporter{
				dataColumnMapping:          ezbookkeepingDataColumnNameMapping,
				transactionTypeNameMapping: ezbookkeepingNameTransactionTypeMapping,
				geoLocationSeparator:       " ",
				transactionTagSeparator:    ";",
			},
			columnSeparator: "\t",
			lineSeparator:   "\n",
		},
	}
)
