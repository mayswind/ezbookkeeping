package mt

import (
	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var mt940TransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE:  utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
	models.TRANSACTION_TYPE_TRANSFER: utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER)),
}

// mt940TransactionDataFileImporter defines the structure of mt940 file importer for statement data
type mt940TransactionDataFileImporter struct{}

// Initialize a mt940 statement data importer singleton instance
var (
	MT940TransactionDataFileImporter = &mt940TransactionDataFileImporter{}
)

// ParseImportedData returns the imported data by parsing the mt940 file statement data
func (c *mt940TransactionDataFileImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, []string, error) {
	mt940DataReader := createNewMT940FileReader(data)
	mt940Data, err := mt940DataReader.read(ctx)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	transactionDataTable, err := createNewMT940TransactionDataTable(mt940Data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(mt940TransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezoneOffset, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
