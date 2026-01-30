package camt

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var camtTransactionTypeNameMapping = map[models.TransactionType]string{
	models.TRANSACTION_TYPE_INCOME:   utils.IntToString(int(models.TRANSACTION_TYPE_INCOME)),
	models.TRANSACTION_TYPE_EXPENSE:  utils.IntToString(int(models.TRANSACTION_TYPE_EXPENSE)),
	models.TRANSACTION_TYPE_TRANSFER: utils.IntToString(int(models.TRANSACTION_TYPE_TRANSFER)),
}

// camt052TransactionDataImporter defines the structure of camt.052 file importer for transaction data
type camt052TransactionDataImporter struct {
}

// camt053TransactionDataImporter defines the structure of camt.053 file importer for transaction data
type camt053TransactionDataImporter struct {
}

// Initialize camt.052 and camt.053 transaction data importer singleton instances
var (
	Camt052TransactionDataImporter = &camt052TransactionDataImporter{}
	Camt053TransactionDataImporter = &camt053TransactionDataImporter{}
)

// ParseImportedData returns the imported data by parsing the camt.052 file transaction data
func (c *camt052TransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	camt052DataReader, err := createNewCamt052FileReader(data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	camt052Data, err := camt052DataReader.read(ctx)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if camt052Data.BankToCustomerAccountReport == nil || camt052Data.BankToCustomerAccountReport.Statements == nil {
		return nil, nil, nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	transactionDataTable, err := createNewCamtStatementTransactionDataTable(camt052Data.BankToCustomerAccountReport.Statements)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(camtTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

// ParseImportedData returns the imported data by parsing the camt.053 file transaction data
func (c *camt053TransactionDataImporter) ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezone *time.Location, additionalOptions converter.TransactionDataImporterOptions, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error) {
	camt053DataReader, err := createNewCamt053FileReader(data)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	camt053Data, err := camt053DataReader.read(ctx)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if camt053Data.BankToCustomerStatement == nil || camt053Data.BankToCustomerStatement.Statements == nil {
		return nil, nil, nil, nil, nil, nil, errs.ErrNotFoundTransactionDataInFile
	}

	transactionDataTable, err := createNewCamtStatementTransactionDataTable(camt053Data.BankToCustomerStatement.Statements)

	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	dataTableImporter := converter.CreateNewSimpleImporterWithTypeNameMapping(camtTransactionTypeNameMapping)

	return dataTableImporter.ParseImportedData(ctx, user, transactionDataTable, defaultTimezone, additionalOptions, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}
