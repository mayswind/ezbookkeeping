package feidee

import (
	"bytes"
	"encoding/csv"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	csvdatatable "github.com/mayswind/ezbookkeeping/pkg/converters/csv"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const feideeMymoneyAppTransactionDataCsvFileHeader = "随手记导出文件(headers:v5;"

const feideeMymoneyAppTransactionTimeColumnName = "日期"
const feideeMymoneyAppTransactionTypeColumnName = "交易类型"
const feideeMymoneyAppTransactionCategoryColumnName = "类别"
const feideeMymoneyAppTransactionSubCategoryColumnName = "子类别"
const feideeMymoneyAppTransactionAccountNameColumnName = "账户"
const feideeMymoneyAppTransactionAccountCurrencyColumnName = "账户币种"
const feideeMymoneyAppTransactionAmountColumnName = "金额"
const feideeMymoneyAppTransactionDescriptionColumnName = "备注"
const feideeMymoneyAppTransactionRelatedIdColumnName = "关联Id"

const feideeMymoneyAppTransactionTypeModifyBalanceText = "余额变更"
const feideeMymoneyAppTransactionTypeModifyOutstandingBalanceText = "负债变更"
const feideeMymoneyAppTransactionTypeIncomeText = "收入"
const feideeMymoneyAppTransactionTypeExpenseText = "支出"
const feideeMymoneyAppTransactionTypeTransferInText = "转入"
const feideeMymoneyAppTransactionTypeTransferOutText = "转出"

var feideeMymoneyAppDataColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: feideeMymoneyAppTransactionTimeColumnName,
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: feideeMymoneyAppTransactionTypeColumnName,
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:         feideeMymoneyAppTransactionCategoryColumnName,
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:     feideeMymoneyAppTransactionSubCategoryColumnName,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     feideeMymoneyAppTransactionAccountNameColumnName,
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY: feideeMymoneyAppTransactionAccountCurrencyColumnName,
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:           feideeMymoneyAppTransactionAmountColumnName,
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:      feideeMymoneyAppTransactionDescriptionColumnName,
}

// feideeMymoneyAppTransactionDataCsvFileImporter defines the structure of feidee mymoney app csv importer for transaction data
type feideeMymoneyAppTransactionDataCsvFileImporter struct{}

// Initialize a feidee mymoney app transaction data csv file importer singleton instance
var (
	FeideeMymoneyAppTransactionDataCsvFileImporter = &feideeMymoneyAppTransactionDataCsvFileImporter{}
)

// ParseImportedData returns the imported data by parsing the feidee mymoney app transaction csv data
func (c *feideeMymoneyAppTransactionDataCsvFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	fallback := unicode.UTF8.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(data), unicode.BOMOverride(fallback))

	dataTable, err := c.createNewFeideeMymoneyAppImportedDataTable(ctx, reader)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	commonDataTable := datatable.CreateNewImportedCommonDataTable(dataTable)

	if !commonDataTable.HasColumn(feideeMymoneyAppTransactionTimeColumnName) ||
		!commonDataTable.HasColumn(feideeMymoneyAppTransactionTypeColumnName) ||
		!commonDataTable.HasColumn(feideeMymoneyAppTransactionSubCategoryColumnName) ||
		!commonDataTable.HasColumn(feideeMymoneyAppTransactionAccountNameColumnName) ||
		!commonDataTable.HasColumn(feideeMymoneyAppTransactionAmountColumnName) ||
		!commonDataTable.HasColumn(feideeMymoneyAppTransactionRelatedIdColumnName) {
		log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.ParseImportedData] cannot parse import data, because missing essential columns in header row")
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	transactionDataTable, err := c.createNewFeideeMymoneyAppTransactionDataTable(ctx, commonDataTable)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(feideeMymoneyTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *feideeMymoneyAppTransactionDataCsvFileImporter) createNewFeideeMymoneyAppImportedDataTable(ctx core.Context, reader io.Reader) (datatable.ImportedDataTable, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	allOriginalLines := make([][]string, 0)
	hasFileHeader := false

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] cannot parse feidee mymoney csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], feideeMymoneyAppTransactionDataCsvFileHeader) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
				continue
			}
		}

		allOriginalLines = append(allOriginalLines, items)
	}

	if !hasFileHeader {
		return nil, errs.ErrInvalidFileHeader
	}

	if len(allOriginalLines) < 2 {
		log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] cannot parse import data, because data table row count is less 1")
		return nil, errs.ErrNotFoundTransactionDataInFile
	}

	dataTable := csvdatatable.CreateNewCustomCsvImportedDataTable(allOriginalLines)

	return dataTable, nil
}

func (c *feideeMymoneyAppTransactionDataCsvFileImporter) createNewFeideeMymoneyAppTransactionDataTable(ctx core.Context, commonDataTable datatable.CommonDataTable) (datatable.TransactionDataTable, error) {
	newColumns := make([]datatable.TransactionDataTableColumn, 0, 11)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME)

	if commonDataTable.HasColumn(feideeMymoneyAppTransactionCategoryColumnName) {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_CATEGORY)
	}

	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME)

	if commonDataTable.HasColumn(feideeMymoneyAppTransactionAccountCurrencyColumnName) {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY)
	}

	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_AMOUNT)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME)

	if commonDataTable.HasColumn(feideeMymoneyAppTransactionAccountCurrencyColumnName) {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY)
	}

	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT)

	if commonDataTable.HasColumn(feideeMymoneyAppTransactionDescriptionColumnName) {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_DESCRIPTION)
	}

	transactionRowParser := createFeideeMymoneyTransactionDataRowParser()
	transactionDataTable := datatable.CreateNewWritableTransactionDataTableWithRowParser(newColumns, transactionRowParser)
	transferTransactionsMap := make(map[string]map[datatable.TransactionDataTableColumn]string, 0)

	commonDataTableIterator := commonDataTable.DataRowIterator()

	for commonDataTableIterator.HasNext() {
		dataRow := commonDataTableIterator.Next()
		rowId := commonDataTableIterator.CurrentRowId()

		if dataRow.ColumnCount() < commonDataTable.HeaderColumnCount() {
			log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] cannot parse row \"%s\", because may missing some columns (column count %d in data row is less than header column count %d)", rowId, dataRow.ColumnCount(), commonDataTable.HeaderColumnCount())
			return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
		}

		data := make(map[datatable.TransactionDataTableColumn]string, 11)

		for columnType, columnName := range feideeMymoneyAppDataColumnNameMapping {
			if dataRow.HasData(columnName) {
				data[columnType] = dataRow.GetData(columnName)
			}
		}

		transactionType := data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE]

		if transactionType == feideeMymoneyAppTransactionTypeModifyBalanceText || transactionType == feideeMymoneyAppTransactionTypeModifyOutstandingBalanceText ||
			transactionType == feideeMymoneyAppTransactionTypeIncomeText || transactionType == feideeMymoneyAppTransactionTypeExpenseText {
			if transactionType == feideeMymoneyAppTransactionTypeModifyBalanceText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]
			} else if transactionType == feideeMymoneyAppTransactionTypeModifyOutstandingBalanceText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeModifyOutstandingBalanceName
			} else if transactionType == feideeMymoneyAppTransactionTypeIncomeText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
			} else if transactionType == feideeMymoneyAppTransactionTypeExpenseText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
			}

			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = ""
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = ""
			transactionDataTable.Add(data)
		} else if transactionType == feideeMymoneyAppTransactionTypeTransferInText || transactionType == feideeMymoneyAppTransactionTypeTransferOutText {
			relatedId := ""

			if dataRow.HasData(feideeMymoneyAppTransactionRelatedIdColumnName) {
				relatedId = dataRow.GetData(feideeMymoneyAppTransactionRelatedIdColumnName)
			}

			if relatedId == "" {
				log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] transfer transaction has blank related id in row \"%s\"", rowId)
				return nil, errs.ErrRelatedIdCannotBeBlank
			}

			relatedData, exists := transferTransactionsMap[relatedId]

			if !exists {
				transferTransactionsMap[relatedId] = data
				continue
			}

			if transactionType == feideeMymoneyAppTransactionTypeTransferInText && relatedData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == feideeMymoneyAppTransactionTypeTransferOutText {
				relatedData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
				relatedData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				relatedData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY]
				relatedData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = data[datatable.TRANSACTION_DATA_TABLE_AMOUNT]
				transactionDataTable.Add(relatedData)
				delete(transferTransactionsMap, relatedId)
			} else if transactionType == feideeMymoneyAppTransactionTypeTransferOutText && relatedData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == feideeMymoneyAppTransactionTypeTransferInText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = relatedData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = relatedData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY]
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = relatedData[datatable.TRANSACTION_DATA_TABLE_AMOUNT]
				transactionDataTable.Add(data)
				delete(transferTransactionsMap, relatedId)
			} else {
				log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] transfer transaction type \"%s\" is not expected in row \"%s\"", transactionType, rowId)
				return nil, errs.ErrTransactionTypeInvalid
			}
		} else {
			log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] cannot parse transaction type \"%s\" in row \"%s\"", transactionType, rowId)
			return nil, errs.ErrTransactionTypeInvalid
		}
	}

	if len(transferTransactionsMap) > 0 {
		log.Errorf(ctx, "[feidee_mymoney_app_transaction_data_csv_file_importer.createNewFeideeMymoneyAppTransactionDataTable] there are %d transactions (related id is %s) which don't have related records", len(transferTransactionsMap), c.getFeideeMymoneyAppRelatedTransactionIds(transferTransactionsMap))
		return nil, errs.ErrFoundRecordNotHasRelatedRecord
	}

	return transactionDataTable, nil
}

func (c *feideeMymoneyAppTransactionDataCsvFileImporter) getFeideeMymoneyAppRelatedTransactionIds(transferTransactionsMap map[string]map[datatable.TransactionDataTableColumn]string) string {
	builder := strings.Builder{}

	for relatedId := range transferTransactionsMap {
		if builder.Len() > 0 {
			builder.WriteRune(',')
		}

		builder.WriteString(relatedId)
	}

	return builder.String()
}
