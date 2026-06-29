package ai

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var aiRecognizedTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE:  utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
	models.TRANSACTION_TYPE_TRANSFER: utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER)),
}

// aiRecognizedTransactionDataImporter represents transaction data importer using AI
type aiRecognizedTransactionDataImporter struct{}

// AIRecognizedTransactionDataImporter is the singleton instance
var AIRecognizedTransactionDataImporter = &aiRecognizedTransactionDataImporter{}

// ParseImportedData returns the imported transaction data parsed by AI
func (c *aiRecognizedTransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, fileData []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	aiRecognizedTransactionDataParser, err := createNewAITransactionDataParser(additionalOptions.GetCurrentConfig())

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	recognizedTransactions, err := aiRecognizedTransactionDataParser.parse(ctx, user, string(fileData), additionalOptions.GetAIAdditionalPrompt(), defaultTimezone, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	transactionDataTable, err := createNewAIRecognizedTransactionDataTable(recognizedTransactions)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(aiRecognizedTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
