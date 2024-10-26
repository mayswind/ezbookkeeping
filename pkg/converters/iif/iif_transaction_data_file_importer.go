package iif

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// iifTransactionDataFileImporter defines the structure of intuit interchange format (iif) for transaction data
type iifTransactionDataFileImporter struct{}

// Initialize an intuit interchange format (iif) file importer singleton instance
var (
	IifTransactionDataFileImporter = &iifTransactionDataFileImporter{}
)

// ParseImportedData returns the imported data by parsing the intuit interchange format (iif) data
func (c *iifTransactionDataFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	iifDataReader := createNewIifDataReader(data)
	accountDatasets, transactionDatasets, err := iifDataReader.read(ctx)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionDataTable, err := createNewIIfTransactionDataTable(ctx, accountDatasets, transactionDatasets)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := datatable.CreateNewSimpleImporter(iifTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
