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

// aiRecognizedTextTransactionDataImporter represents transaction data importer using AI text recognition
type aiRecognizedTextTransactionDataImporter struct{}

// AIRecognizedTextTransactionDataImporter is the singleton instance
var AIRecognizedTextTransactionDataImporter = &aiRecognizedTextTransactionDataImporter{}

// aiRecognizedImageTransactionDataImporter represents transaction data importer using AI image recognition
type aiRecognizedImageTransactionDataImporter struct{}

// AIRecognizedImageTransactionDataImporter is the singleton instance
var AIRecognizedImageTransactionDataImporter = &aiRecognizedImageTransactionDataImporter{}

// ParseImportedData returns the imported transaction data parsed by AI text recognition
func (c *aiRecognizedTextTransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, fileData []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	aiRecognizedTransactionDataParser, err := createNewAITextTransactionDataParser(additionalOptions.GetCurrentConfig())

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	recognizedTransactions, err := aiRecognizedTransactionDataParser.parseText(ctx, user, string(fileData), additionalOptions.GetAIAdditionalPrompt(), defaultTimezone, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return getImportTransactionResponse(ctx, user, recognizedTransactions, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

// ParseImportedData returns the imported transaction data parsed by AI image recognition
func (c *aiRecognizedImageTransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, fileData []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	aiRecognizedTransactionDataParser, err := createNewAIImageTransactionDataParser(additionalOptions.GetCurrentConfig())

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	recognizedTransactions, err := aiRecognizedTransactionDataParser.parseImage(ctx, user, fileData, additionalOptions.GetAIAdditionalPrompt(), additionalOptions, defaultTimezone, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return getImportTransactionResponse(ctx, user, recognizedTransactions, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

// getImportTransactionResponse returns the imported transaction data parsed by AI recognized results
func getImportTransactionResponse(ctx core.Context, user *models.User, results []*models.RecognizedTransactionResult, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	transactionDataTable, err := createNewAIRecognizedTransactionDataTable(results)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(aiRecognizedTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
