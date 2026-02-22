package custom

import (
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	csvconverter "github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/converters/excel"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const customOOXMLExcelFileType = "custom_xlsx"
const customMSCFBExcelFileType = "custom_xls"

// customTransactionDataExcelFileImporter defines the structure of custom excel importer for transaction data
type customTransactionDataExcelFileImporter struct {
	fileType                   string
	columnIndexMapping         map[datatable.TransactionDataTableColumn]int
	transactionTypeNameMapping map[string]models.TransactionType
	hasHeaderLine              bool
	timeFormat                 string
	timezoneFormat             string
	amountDecimalSeparator     string
	amountDigitGroupingSymbol  string
	geoLocationSeparator       string
	geoLocationOrder           converter.TransactionGeoLocationOrder
	transactionTagSeparator    string
}

// ParseDataLines returns the parsed file lines for specified the excel file data
func (c *customTransactionDataExcelFileImporter) ParseDataLines(ctx core.Context, data []byte) ([][]string, error) {
	var excelDataTable datatable.BasicDataTable
	var err error

	if c.fileType == customOOXMLExcelFileType {
		excelDataTable, err = excel.CreateNewExcelOOXMLFileBasicDataTable(data, false)
	} else if c.fileType == customMSCFBExcelFileType {
		excelDataTable, err = excel.CreateNewExcelMSCFBFileBasicDataTable(data, false)
	} else {
		return nil, errs.ErrImportFileTypeNotSupported
	}

	if err != nil {
		return nil, err
	}

	iterator := excelDataTable.DataRowIterator()
	allLines := make([][]string, 0)

	for iterator.HasNext() {
		row := iterator.Next()
		items := make([]string, row.ColumnCount())

		for i := 0; i < row.ColumnCount(); i++ {
			items[i] = strings.Trim(row.GetData(i), " ")
		}

		allLines = append(allLines, items)
	}

	return allLines, nil
}

// ParseImportedData returns the imported data by parsing the custom transaction dsv data
func (c *customTransactionDataExcelFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	allLines, err := c.ParseDataLines(ctx, data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTable := csvconverter.CreateNewCustomCsvBasicDataTable(allLines, c.hasHeaderLine)
	transactionDataTable := CreateNewCustomPlainTextDataTable(dataTable, c.columnIndexMapping, c.transactionTypeNameMapping, c.timeFormat, c.timezoneFormat, c.amountDecimalSeparator, c.amountDigitGroupingSymbol)
	dataTableImporter := converter.CreateNewImporterWithTypeNameMapping(customTransactionTypeNameMapping, c.geoLocationSeparator, c.geoLocationOrder, c.transactionTagSeparator)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

// IsCustomExcelFileType returns whether the file type is the custom excel file type
func IsCustomExcelFileType(fileType string) bool {
	return fileType == customOOXMLExcelFileType || fileType == customMSCFBExcelFileType
}

// CreateNewCustomTransactionDataExcelFileParser returns a new custom transaction data parser
func CreateNewCustomTransactionDataExcelFileParser(fileType string) (CustomTransactionDataParser, error) {
	if fileType != customOOXMLExcelFileType && fileType != customMSCFBExcelFileType {
		return nil, errs.ErrImportFileTypeNotSupported
	}

	return &customTransactionDataExcelFileImporter{
		fileType: fileType,
	}, nil
}

// CreateNewCustomTransactionDataExcelFileImporter returns a new custom excel importer for transaction data
func CreateNewCustomTransactionDataExcelFileImporter(fileType string, columnIndexMapping map[datatable.TransactionDataTableColumn]int, transactionTypeNameMapping map[string]models.TransactionType, hasHeaderLine bool, timeFormat string, timezoneFormat string, amountDecimalSeparator string, amountDigitGroupingSymbol string, geoLocationSeparator string, geoLocationOrder string, transactionTagSeparator string) (converter.TransactionDataImporter, error) {
	if fileType != customOOXMLExcelFileType && fileType != customMSCFBExcelFileType {
		return nil, errs.ErrImportFileTypeNotSupported
	}

	if geoLocationOrder == "" {
		geoLocationOrder = string(converter.TRANSACTION_GEO_LOCATION_ORDER_LONGITUDE_LATITUDE)
	} else if geoLocationOrder != string(converter.TRANSACTION_GEO_LOCATION_ORDER_LONGITUDE_LATITUDE) &&
		geoLocationOrder != string(converter.TRANSACTION_GEO_LOCATION_ORDER_LATITUDE_LONGITUDE) {
		return nil, errs.ErrImportFileTypeNotSupported
	}

	if _, exists := columnIndexMapping[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME]; !exists {
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if _, exists := columnIndexMapping[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE]; !exists {
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	if _, exists := columnIndexMapping[datatable.TRANSACTION_DATA_TABLE_AMOUNT]; !exists {
		return nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	return &customTransactionDataExcelFileImporter{
		fileType:                   fileType,
		columnIndexMapping:         columnIndexMapping,
		transactionTypeNameMapping: transactionTypeNameMapping,
		hasHeaderLine:              hasHeaderLine,
		timeFormat:                 timeFormat,
		timezoneFormat:             timezoneFormat,
		amountDecimalSeparator:     amountDecimalSeparator,
		amountDigitGroupingSymbol:  amountDigitGroupingSymbol,
		geoLocationSeparator:       geoLocationSeparator,
		geoLocationOrder:           converter.TransactionGeoLocationOrder(geoLocationOrder),
		transactionTagSeparator:    transactionTagSeparator,
	}, nil
}
