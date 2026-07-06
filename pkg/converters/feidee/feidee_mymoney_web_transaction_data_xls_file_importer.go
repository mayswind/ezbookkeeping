package feidee

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/converters/excel"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// feideeMymoneyWebDataLegacyColumnNameMapping is the column name mapping for the old format
var feideeMymoneyWebDataLegacyColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     "交易类型",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:             "分类",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         "子分类",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         "账户1",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               "金额",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: "账户2",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          "备注",
	datatable.TRANSACTION_DATA_TABLE_MEMBER:               "成员",
	datatable.TRANSACTION_DATA_TABLE_PROJECT:              "项目",
	datatable.TRANSACTION_DATA_TABLE_MERCHANT:             "商家",
}

// feideeMymoneyWebDataExpenseColumnNameMapping is the column name mapping for the expense transactions sheet
var feideeMymoneyWebDataExpenseColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "交易类型",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:         "一级分类",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:     "二级分类",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     "支出账户",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:           "金额",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:      "备注",
	datatable.TRANSACTION_DATA_TABLE_MEMBER:           "成员",
	datatable.TRANSACTION_DATA_TABLE_MERCHANT:         "商家",
	datatable.TRANSACTION_DATA_TABLE_PROJECT:          "项目",
}

// feideeMymoneyWebDataIncomeColumnNameMapping is the column name mapping for the income transactions sheet
var feideeMymoneyWebDataIncomeColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME: "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE: "交易类型",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:         "一级分类",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:     "二级分类",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:     "收入账户",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:           "金额",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:      "备注",
	datatable.TRANSACTION_DATA_TABLE_MEMBER:           "成员",
	datatable.TRANSACTION_DATA_TABLE_MERCHANT:         "商家",
	datatable.TRANSACTION_DATA_TABLE_PROJECT:          "项目",
}

// feideeMymoneyWebDataTransferColumnNameMapping is the column name mapping for the transfer transactions sheet
var feideeMymoneyWebDataTransferColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     "交易类型",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         "转出账户",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: "转入账户",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               "金额",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          "备注",
	datatable.TRANSACTION_DATA_TABLE_MEMBER:               "成员",
	datatable.TRANSACTION_DATA_TABLE_MERCHANT:             "商家",
	datatable.TRANSACTION_DATA_TABLE_PROJECT:              "项目",
}

// feideeMymoneyWebDataBalanceModificationColumnNameMapping is the column name mapping for the balance modification transactions sheet
var feideeMymoneyWebDataBalanceModificationColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     "交易类型",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:             "一级分类",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         "二级分类",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         "账户1",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               "金额",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: "账户2",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          "备注",
	datatable.TRANSACTION_DATA_TABLE_MEMBER:               "成员",
	datatable.TRANSACTION_DATA_TABLE_MERCHANT:             "商家",
	datatable.TRANSACTION_DATA_TABLE_PROJECT:              "项目",
}

// feideeMymoneyWebTransactionTypeNames maps transaction type names to transaction types for the web importer
var feideeMymoneyWebTransactionTypeNames = map[string]models.TransactionType{
	feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_MODIFY_BALANCE]: models.TRANSACTION_TYPE_MODIFY_BALANCE,
	feideeMymoneyTransactionTypeModifyOutstandingBalanceName:                        models.TRANSACTION_TYPE_MODIFY_BALANCE,
	feideeMymoneyTransactionTypeDebtModificationName:                                models.TRANSACTION_TYPE_MODIFY_BALANCE,
	feideeMymoneyTransactionTypeReceivableModificationName:                          models.TRANSACTION_TYPE_MODIFY_BALANCE,
	feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_INCOME]:         models.TRANSACTION_TYPE_INCOME,
	feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_EXPENSE]:        models.TRANSACTION_TYPE_EXPENSE,
	feideeMymoneyTransactionTypeNameMapping[models.TRANSACTION_TYPE_TRANSFER]:       models.TRANSACTION_TYPE_TRANSFER,
}

// feideeMymoneyWebTransactionDataXlsFileImporter defines the structure of feidee mymoney (web) xls importer for transaction data
type feideeMymoneyWebTransactionDataXlsFileImporter struct {
	converter.DataTableTransactionDataImporter
}

// Initialize a feidee mymoney (web) transaction data xls file importer singleton instance
var (
	FeideeMymoneyWebTransactionDataXlsFileImporter = &feideeMymoneyWebTransactionDataXlsFileImporter{}
)

// ParseImportedData returns the imported data by parsing the feidee mymoney (web) transaction xls data
func (c *feideeMymoneyWebTransactionDataXlsFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTables, err := excel.CreateNewExcelMSCFBFileBasicDataTables(data, true)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionRowParser := createFeideeMymoneyTransactionDataRowParser()
	transactionDataTables := make([]datatable.TransactionDataTable, 0, len(dataTables))

	for _, dataTable := range dataTables {
		columnMapping := c.getFeideeMymoneyWebColumnNameMapping(dataTable.HeaderColumnNames())
		transactionDataTable := datatable.CreateNewTransactionDataTableFromBasicDataTableWithRowParser(dataTable, columnMapping, transactionRowParser)
		transactionDataTables = append(transactionDataTables, transactionDataTable)
	}

	mergedTransactionDataTable := datatable.CreateNewMergedTransactionDataTable(transactionDataTables)
	dataTableImporter := converter.CreateNewSimpleImporter(feideeMymoneyWebTransactionTypeNames)

	return dataTableImporter.ParseImportedData(ctx, user, mergedTransactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

// getFeideeMymoneyWebColumnNameMapping returns the appropriate column name mapping based on the header column names
func (c *feideeMymoneyWebTransactionDataXlsFileImporter) getFeideeMymoneyWebColumnNameMapping(headerColumns []string) map[datatable.TransactionDataTableColumn]string {
	headerSet := make(map[string]bool, len(headerColumns))

	for _, col := range headerColumns {
		headerSet[col] = true
	}

	if headerSet["转出账户"] {
		return feideeMymoneyWebDataTransferColumnNameMapping
	} else if headerSet["一级分类"] {
		if headerSet["支出账户"] {
			return feideeMymoneyWebDataExpenseColumnNameMapping
		}

		if headerSet["收入账户"] {
			return feideeMymoneyWebDataIncomeColumnNameMapping
		}

		return feideeMymoneyWebDataBalanceModificationColumnNameMapping
	}

	return feideeMymoneyWebDataLegacyColumnNameMapping
}
