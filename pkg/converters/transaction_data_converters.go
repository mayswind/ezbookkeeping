package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/alipay"
	"github.com/mayswind/ezbookkeeping/pkg/converters/base"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/converters/default"
	"github.com/mayswind/ezbookkeeping/pkg/converters/dsv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/feidee"
	"github.com/mayswind/ezbookkeeping/pkg/converters/fireflyIII"
	"github.com/mayswind/ezbookkeeping/pkg/converters/gnucash"
	"github.com/mayswind/ezbookkeeping/pkg/converters/iif"
	"github.com/mayswind/ezbookkeeping/pkg/converters/ofx"
	"github.com/mayswind/ezbookkeeping/pkg/converters/qif"
	"github.com/mayswind/ezbookkeeping/pkg/converters/wechat"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// GetTransactionDataExporter returns the transaction data exporter according to the file type
func GetTransactionDataExporter(fileType string) base.TransactionDataExporter {
	if fileType == "csv" {
		return _default.DefaultTransactionDataCSVFileConverter
	} else if fileType == "tsv" {
		return _default.DefaultTransactionDataTSVFileConverter
	} else {
		return nil
	}
}

// GetTransactionDataImporter returns the transaction data importer according to the file type
func GetTransactionDataImporter(fileType string) (base.TransactionDataImporter, error) {
	if fileType == "ezbookkeeping_csv" {
		return _default.DefaultTransactionDataCSVFileConverter, nil
	} else if fileType == "ezbookkeeping_tsv" {
		return _default.DefaultTransactionDataTSVFileConverter, nil
	} else if fileType == "ofx" {
		return ofx.OFXTransactionDataImporter, nil
	} else if fileType == "qfx" {
		return ofx.OFXTransactionDataImporter, nil
	} else if fileType == "qif_ymd" {
		return qif.QifYearMonthDayTransactionDataImporter, nil
	} else if fileType == "qif_mdy" {
		return qif.QifMonthDayYearTransactionDataImporter, nil
	} else if fileType == "qif_dmy" {
		return qif.QifDayMonthYearTransactionDataImporter, nil
	} else if fileType == "iif" {
		return iif.IifTransactionDataFileImporter, nil
	} else if fileType == "gnucash" {
		return gnucash.GnuCashTransactionDataImporter, nil
	} else if fileType == "firefly_iii_csv" {
		return fireflyIII.FireflyIIITransactionDataCsvFileImporter, nil
	} else if fileType == "feidee_mymoney_csv" {
		return feidee.FeideeMymoneyAppTransactionDataCsvFileImporter, nil
	} else if fileType == "feidee_mymoney_xls" {
		return feidee.FeideeMymoneyWebTransactionDataXlsFileImporter, nil
	} else if fileType == "alipay_app_csv" {
		return alipay.AlipayAppTransactionDataCsvFileImporter, nil
	} else if fileType == "alipay_web_csv" {
		return alipay.AlipayWebTransactionDataCsvFileImporter, nil
	} else if fileType == "wechat_pay_app_csv" {
		return wechat.WeChatPayTransactionDataCsvFileImporter, nil
	} else {
		return nil, errs.ErrImportFileTypeNotSupported
	}
}

// IsCustomDelimiterSeparatedValuesFileType returns whether the file type is the delimiter-separated values file type
func IsCustomDelimiterSeparatedValuesFileType(fileType string) bool {
	return dsv.IsDelimiterSeparatedValuesFileType(fileType)
}

// CreateNewDelimiterSeparatedValuesDataParser returns a new delimiter-separated values data parser according to the file type and encoding
func CreateNewDelimiterSeparatedValuesDataParser(fileType string, fileEncoding string) (dsv.CustomTransactionDataDsvFileParser, error) {
	return dsv.CreateNewCustomTransactionDataDsvFileParser(fileType, fileEncoding)
}

// CreateNewDelimiterSeparatedValuesDataImporter returns a new delimiter-separated values data importer according to the file type and encoding
func CreateNewDelimiterSeparatedValuesDataImporter(fileType string, fileEncoding string, columnIndexMapping map[datatable.TransactionDataTableColumn]int, transactionTypeNameMapping map[string]models.TransactionType, hasHeaderLine bool, timeFormat string, timezoneFormat string, geoLocationSeparator string, transactionTagSeparator string) (base.TransactionDataImporter, error) {
	return dsv.CreateNewCustomTransactionDataDsvFileImporter(fileType, fileEncoding, columnIndexMapping, transactionTypeNameMapping, hasHeaderLine, timeFormat, timezoneFormat, geoLocationSeparator, transactionTagSeparator)
}
