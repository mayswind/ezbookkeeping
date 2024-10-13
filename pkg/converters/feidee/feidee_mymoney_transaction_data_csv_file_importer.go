package feidee

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const feideeMymoneyTransactionDataCsvFileHeader = "随手记导出文件(headers:v5;"
const feideeMymoneyTransactionDataCsvFileHeaderWithUtf8Bom = "\xEF\xBB\xBF" + feideeMymoneyTransactionDataCsvFileHeader

const feideeMymoneyCsvFileTransactionTypeModifyBalanceText = "余额变更"
const feideeMymoneyCsvFileTransactionTypeIncomeText = "收入"
const feideeMymoneyCsvFileTransactionTypeExpenseText = "支出"
const feideeMymoneyCsvFileTransactionTypeTransferInText = "转入"
const feideeMymoneyCsvFileTransactionTypeTransferOutText = "转出"

// feideeMymoneyTransactionDataCsvImporter defines the structure of feidee mymoney csv importer for transaction data
type feideeMymoneyTransactionDataCsvImporter struct{}

// Initialize a feidee mymoney transaction data csv file importer singleton instance
var (
	FeideeMymoneyTransactionDataCsvImporter = &feideeMymoneyTransactionDataCsvImporter{}
)

// ParseImportedData returns the imported data by parsing the feidee mymoney transaction csv data
func (c *feideeMymoneyTransactionDataCsvImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	content := string(data)

	if strings.Index(content, feideeMymoneyTransactionDataCsvFileHeader) != 0 && strings.Index(content, feideeMymoneyTransactionDataCsvFileHeaderWithUtf8Bom) != 0 {
		return nil, nil, nil, nil, nil, nil, errs.ErrInvalidFileHeader
	}

	allLines, err := c.parseAllLinesFromCsvData(ctx, content)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if len(allLines) < 2 {
		log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] cannot parse import data for user \"uid:%d\", because data table row count is less 1", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	headerLineItems := allLines[0]
	headerItemMap := make(map[string]int)

	for i := 0; i < len(headerLineItems); i++ {
		headerItemMap[headerLineItems[i]] = i
	}

	timeColumnIdx, timeColumnExists := headerItemMap["日期"]
	typeColumnIdx, typeColumnExists := headerItemMap["交易类型"]
	categoryColumnIdx, categoryColumnExists := headerItemMap["类别"]
	subCategoryColumnIdx, subCategoryColumnExists := headerItemMap["子类别"]
	accountColumnIdx, accountColumnExists := headerItemMap["账户"]
	accountCurrencyColumnIdx, accountCurrencyColumnExists := headerItemMap["账户币种"]
	amountColumnIdx, amountColumnExists := headerItemMap["金额"]
	descriptionColumnIdx, descriptionColumnExists := headerItemMap["备注"]
	relatedIdColumnIdx, relatedIdColumnExists := headerItemMap["关联Id"]

	if !timeColumnExists || !typeColumnExists || !subCategoryColumnExists ||
		!accountColumnExists || !amountColumnExists || !relatedIdColumnExists {
		log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] cannot parse import data for user \"uid:%d\", because missing essential columns in header row", user.Uid)
		return nil, nil, nil, nil, nil, nil, errs.ErrMissingRequiredFieldInHeaderRow
	}

	newColumns := make([]datatable.TransactionDataTableColumn, 0, 11)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME)

	if categoryColumnExists {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_CATEGORY)
	}

	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME)

	if accountCurrencyColumnExists {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY)
	}

	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_AMOUNT)
	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME)

	if accountCurrencyColumnExists {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY)
	}

	newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT)

	if descriptionColumnExists {
		newColumns = append(newColumns, datatable.TRANSACTION_DATA_TABLE_DESCRIPTION)
	}

	transactionRowParser := createFeideeMymoneyTransactionDataRowParser()
	dataTable := datatable.CreateNewWritableTransactionDataTableWithRowParser(newColumns, transactionRowParser)
	transferTransactionsMap := make(map[string]map[datatable.TransactionDataTableColumn]string, 0)

	for i := 1; i < len(allLines); i++ {
		items := allLines[i]

		if len(items) < len(headerLineItems) {
			log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] cannot parse row \"index:%d\" for user \"uid:%d\", because may missing some columns (column count %d in data row is less than header column count %d)", i, user.Uid, len(items), len(headerLineItems))
			return nil, nil, nil, nil, nil, nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
		}

		data, relatedId := c.parseTransactionData(items,
			timeColumnIdx,
			timeColumnExists,
			typeColumnIdx,
			typeColumnExists,
			categoryColumnIdx,
			categoryColumnExists,
			subCategoryColumnIdx,
			subCategoryColumnExists,
			accountColumnIdx,
			accountColumnExists,
			accountCurrencyColumnIdx,
			accountCurrencyColumnExists,
			amountColumnIdx,
			amountColumnExists,
			descriptionColumnIdx,
			descriptionColumnExists,
			relatedIdColumnIdx,
			relatedIdColumnExists,
		)

		transactionType := data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE]

		if transactionType == feideeMymoneyCsvFileTransactionTypeModifyBalanceText || transactionType == feideeMymoneyCsvFileTransactionTypeIncomeText || transactionType == feideeMymoneyCsvFileTransactionTypeExpenseText {
			if transactionType == feideeMymoneyCsvFileTransactionTypeModifyBalanceText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]
			} else if transactionType == feideeMymoneyCsvFileTransactionTypeIncomeText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]
			} else if transactionType == feideeMymoneyCsvFileTransactionTypeExpenseText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]
			}

			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = ""
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = ""
			data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = ""
			dataTable.Add(data)
		} else if transactionType == feideeMymoneyCsvFileTransactionTypeTransferInText || transactionType == feideeMymoneyCsvFileTransactionTypeTransferOutText {
			if relatedId == "" {
				log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] transfer transaction has blank related id in row \"index:%d\" for user \"uid:%d\"", i, user.Uid)
				return nil, nil, nil, nil, nil, nil, errs.ErrRelatedIdCannotBeBlank
			}

			relatedData, exists := transferTransactionsMap[relatedId]

			if !exists {
				transferTransactionsMap[relatedId] = data
				continue
			}

			if transactionType == feideeMymoneyCsvFileTransactionTypeTransferInText && relatedData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == feideeMymoneyCsvFileTransactionTypeTransferOutText {
				relatedData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
				relatedData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				relatedData[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY]
				relatedData[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = data[datatable.TRANSACTION_DATA_TABLE_AMOUNT]
				dataTable.Add(relatedData)
				delete(transferTransactionsMap, relatedId)
			} else if transactionType == feideeMymoneyCsvFileTransactionTypeTransferOutText && relatedData[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] == feideeMymoneyCsvFileTransactionTypeTransferInText {
				data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME] = relatedData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME]
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY] = relatedData[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY]
				data[datatable.TRANSACTION_DATA_TABLE_RELATED_AMOUNT] = relatedData[datatable.TRANSACTION_DATA_TABLE_AMOUNT]
				dataTable.Add(data)
				delete(transferTransactionsMap, relatedId)
			} else {
				log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] transfer transaction type \"%s\" is not expected in row \"index:%d\" for user \"uid:%d\"", transactionType, i, user.Uid)
				return nil, nil, nil, nil, nil, nil, errs.ErrTransactionTypeInvalid
			}
		} else {
			log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] cannot parse transaction type \"%s\" in row \"index:%d\" for user \"uid:%d\"", transactionType, i, user.Uid)
			return nil, nil, nil, nil, nil, nil, errs.ErrTransactionTypeInvalid
		}
	}

	if len(transferTransactionsMap) > 0 {
		log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.ParseImportedData] there are %d transactions (related id is %s) which don't have related records", len(transferTransactionsMap), c.getRelatedIds(transferTransactionsMap))
		return nil, nil, nil, nil, nil, nil, errs.ErrFoundRecordNotHasRelatedRecord
	}

	dataTableImporter := datatable.CreateNewSimpleImporter(feideeMymoneyTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, dataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (c *feideeMymoneyTransactionDataCsvImporter) parseAllLinesFromCsvData(ctx core.Context, content string) ([][]string, error) {
	csvReader := csv.NewReader(strings.NewReader(content))
	csvReader.FieldsPerRecord = -1

	allLines := make([][]string, 0)
	hasFileHeader := false

	for {
		items, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Errorf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.parseAllLinesFromCsvData] cannot parse feidee mymoney csv data, because %s", err.Error())
			return nil, errs.ErrInvalidCSVFile
		}

		if !hasFileHeader {
			if len(items) <= 0 {
				continue
			} else if strings.Index(items[0], feideeMymoneyTransactionDataCsvFileHeader) == 0 || strings.Index(items[0], feideeMymoneyTransactionDataCsvFileHeaderWithUtf8Bom) == 0 {
				hasFileHeader = true
				continue
			} else {
				log.Warnf(ctx, "[feidee_mymoney_transaction_data_csv_file_importer.parseAllLinesFromCsvData] read unexpected line before read file header, line content is %s", strings.Join(items, ","))
			}
		}

		allLines = append(allLines, items)
	}

	return allLines, nil
}

func (c *feideeMymoneyTransactionDataCsvImporter) parseTransactionData(
	items []string,
	timeColumnIdx int,
	timeColumnExists bool,
	typeColumnIdx int,
	typeColumnExists bool,
	categoryColumnIdx int,
	categoryColumnExists bool,
	subCategoryColumnIdx int,
	subCategoryColumnExists bool,
	accountColumnIdx int,
	accountColumnExists bool,
	accountCurrencyColumnIdx int,
	accountCurrencyColumnExists bool,
	amountColumnIdx int,
	amountColumnExists bool,
	descriptionColumnIdx int,
	descriptionColumnExists bool,
	relatedIdColumnIdx int,
	relatedIdColumnExists bool,
) (map[datatable.TransactionDataTableColumn]string, string) {
	data := make(map[datatable.TransactionDataTableColumn]string, 11)
	relatedId := ""

	if timeColumnExists && timeColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME] = items[timeColumnIdx]
	}

	if typeColumnExists && typeColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE] = items[typeColumnIdx]
	}

	if categoryColumnExists && categoryColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_CATEGORY] = items[categoryColumnIdx]
	}

	if subCategoryColumnExists && subCategoryColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY] = items[subCategoryColumnIdx]
	}

	if accountColumnExists && accountColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME] = items[accountColumnIdx]
	}

	if accountCurrencyColumnExists && accountCurrencyColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY] = items[accountCurrencyColumnIdx]
	}

	if amountColumnExists && amountColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_AMOUNT] = items[amountColumnIdx]
	}

	if descriptionColumnExists && descriptionColumnIdx < len(items) {
		data[datatable.TRANSACTION_DATA_TABLE_DESCRIPTION] = items[descriptionColumnIdx]
	}

	if relatedIdColumnExists && relatedIdColumnIdx < len(items) {
		relatedId = items[relatedIdColumnIdx]
	}

	return data, relatedId
}

func (c *feideeMymoneyTransactionDataCsvImporter) getRelatedIds(transferTransactionsMap map[string]map[datatable.TransactionDataTableColumn]string) string {
	builder := strings.Builder{}

	for relatedId := range transferTransactionsMap {
		if builder.Len() > 0 {
			builder.WriteRune(',')
		}

		builder.WriteString(relatedId)
	}

	return builder.String()
}
