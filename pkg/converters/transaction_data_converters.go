package converters

import "github.com/mayswind/ezbookkeeping/pkg/errs"

// GetTransactionDataExporter returns the transaction data exporter according to the file type
func GetTransactionDataExporter(fileType string) TransactionDataExporter {
	if fileType == "csv" {
		return EzBookKeepingTransactionDataCSVFileConverter
	} else if fileType == "tsv" {
		return EzBookKeepingTransactionDataTSVFileConverter
	} else {
		return nil
	}
}

// GetTransactionDataImporter returns the transaction data importer according to the file type
func GetTransactionDataImporter(fileType string) (TransactionDataImporter, error) {
	if fileType == "ezbookkeeping_csv" {
		return EzBookKeepingTransactionDataCSVFileConverter, nil
	} else if fileType == "ezbookkeeping_tsv" {
		return EzBookKeepingTransactionDataTSVFileConverter, nil
	} else if fileType == "feidee_mymoney_csv" {
		return FeideeMymoneyTransactionDataCsvImporter, nil
	} else if fileType == "feidee_mymoney_xls" {
		return FeideeMymoneyTransactionDataXlsImporter, nil
	} else {
		return nil, errs.ErrImportFileTypeNotSupported
	}
}
