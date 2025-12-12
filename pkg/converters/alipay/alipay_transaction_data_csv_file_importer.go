package alipay

import (
	"bytes"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var alipayTransactionSupportedColumns = map[datatable.TransactionDataTableColumn]bool{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     true,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     true,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         true,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         true,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               true,
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: true,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          true,
}

var alipayTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   "收入",
	models.TRANSACTION_TYPE_EXPENSE:  "支出",
	models.TRANSACTION_TYPE_TRANSFER: "不计收支",
}

// alipayTransactionColumnNames defines the structure of alipay transaction plain text header names
type alipayTransactionColumnNames struct {
	timeColumnName           string
	categoryColumnName       string
	targetNameColumnName     string
	productNameColumnName    string
	amountColumnName         string
	typeColumnName           string
	relatedAccountColumnName string
	statusColumnName         string
	descriptionColumnName    string
}

// alipayTransactionDataCsvFileImporter defines the structure of alipay csv importer for transaction data
type alipayTransactionDataCsvFileImporter struct {
	fileHeaderLine         string
	dataHeaderStartContent []string
	dataBottomEndLineRune  rune
	originalColumnNames    alipayTransactionColumnNames
}

// ParseImportedData returns the imported data by parsing the alipay transaction csv data
func (c *alipayTransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, []string, error) {
	enc := simplifiedchinese.GB18030
	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())

	csvDataTable, err := csv.CreateNewCsvBasicDataTable(ctx, reader, false)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	dataTable, err := createNewAlipayTransactionBasicDataTable(ctx, csvDataTable, c.fileHeaderLine, c.dataHeaderStartContent, c.dataBottomEndLineRune)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewCommonDataTableFromBasicDataTable(dataTable)

	if !commonDataTable.HasColumn(c.originalColumnNames.timeColumnName) ||
		!commonDataTable.HasColumn(c.originalColumnNames.amountColumnName) ||
		!commonDataTable.HasColumn(c.originalColumnNames.typeColumnName) ||
		!commonDataTable.HasColumn(c.originalColumnNames.statusColumnName) {
		log.Errorf(ctx, "[alipay_transaction_data_csv_file_importer.ParseImportedData] cannot parse alipay csv data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionRowParser := createAlipayTransactionDataRowParser(c.originalColumnNames, dataTable.HeaderColumnNames())
	transactionDataTable := datatable.CreateNewTransactionDataTableFromCommonDataTable(commonDataTable, alipayTransactionSupportedColumns, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(alipayTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
