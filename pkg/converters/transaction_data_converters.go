package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/alipay"
	"github.com/mayswind/ezbookkeeping/pkg/converters/base"
	"github.com/mayswind/ezbookkeeping/pkg/converters/default"
	"github.com/mayswind/ezbookkeeping/pkg/converters/feidee"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// GetTransactionDataExporter returns the transaction data exporter according to the file type
func GetTransactionDataExporter(fileType string) base.TransactionDataExporter {
	if fileType == "csv" {
		return _default.EzBookKeepingTransactionDataCSVFileConverter
	} else if fileType == "tsv" {
		return _default.EzBookKeepingTransactionDataTSVFileConverter
	} else {
		return nil
	}
}

// GetTransactionDataImporter returns the transaction data importer according to the file type
func GetTransactionDataImporter(fileType string) (base.TransactionDataImporter, error) {
	if fileType == "ezbookkeeping_csv" {
		return _default.EzBookKeepingTransactionDataCSVFileConverter, nil
	} else if fileType == "ezbookkeeping_tsv" {
		return _default.EzBookKeepingTransactionDataTSVFileConverter, nil
	} else if fileType == "feidee_mymoney_csv" {
		return feidee.FeideeMymoneyTransactionDataCsvImporter, nil
	} else if fileType == "feidee_mymoney_xls" {
		return feidee.FeideeMymoneyTransactionDataXlsImporter, nil
	} else if fileType == "alipay_web_csv" {
		return alipay.AlipayWebTransactionDataCsvImporter, nil
	} else {
		return nil, errs.ErrImportFileTypeNotSupported
	}
}
