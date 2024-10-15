package feidee

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/converters/excel"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// feideeMymoneyWebTransactionDataXlsFileImporter defines the structure of feidee mymoney (web) xls importer for transaction data
type feideeMymoneyWebTransactionDataXlsFileImporter struct {
	datatable.DataTableTransactionDataImporter
}

// Initialize a feidee mymoney (web) transaction data xls file importer singleton instance
var (
	FeideeMymoneyWebTransactionDataXlsFileImporter = &feideeMymoneyWebTransactionDataXlsFileImporter{}
)

// ParseImportedData returns the imported data by parsing the feidee mymoney (web) transaction xls data
func (c *feideeMymoneyWebTransactionDataXlsFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	dataTable, err := excel.CreateNewExcelFileImportedDataTable(data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionRowParser := createFeideeMymoneyTransactionDataRowParser()
	transactionDataTable := datatable.CreateImportedTransactionDataTableWithRowParser(dataTable, feideeMymoneyDataColumnNameMapping, transactionRowParser)
	dataTableImporter := datatable.CreateNewSimpleImporter(feideeMymoneyTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
