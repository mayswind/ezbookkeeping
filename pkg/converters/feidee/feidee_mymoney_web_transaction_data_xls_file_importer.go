package feidee

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/converters/excel"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

var feideeMymoneyWebDataColumnNameMapping = map[datatable.TransactionDataTableColumn]string{
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TIME:     "日期",
	datatable.TRANSACTION_DATA_TABLE_TRANSACTION_TYPE:     "交易类型",
	datatable.TRANSACTION_DATA_TABLE_CATEGORY:             "分类",
	datatable.TRANSACTION_DATA_TABLE_SUB_CATEGORY:         "子分类",
	datatable.TRANSACTION_DATA_TABLE_ACCOUNT_NAME:         "账户1",
	datatable.TRANSACTION_DATA_TABLE_AMOUNT:               "金额",
	datatable.TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME: "账户2",
	datatable.TRANSACTION_DATA_TABLE_DESCRIPTION:          "备注",
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
func (c *feideeMymoneyWebTransactionDataXlsFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := excel.CreateNewExcelMSCFBFileImportedDataTable(data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionRowParser := createFeideeMymoneyTransactionDataRowParser()
	transactionDataTable := datatable.CreateNewImportedTransactionDataTableWithRowParser(dataTable, feideeMymoneyWebDataColumnNameMapping, transactionRowParser)
	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(feideeMymoneyTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
