package converter

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// TransactionDataExporter defines the structure of transaction data exporter
type TransactionDataExporter interface {
	// ToExportedContent returns the exported data
	ToExportedContent(ctx core.Context, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error)
}

// TransactionDataImporter defines the structure of transaction data importer
type TransactionDataImporter interface {
	// ParseImportedData returns the imported data
	ParseImportedData(ctx core.Context, user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (models.ImportedTransactionSlice, []*models.Account, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionCategory, []*models.TransactionTag, error)
}

// TransactionDataConverter defines the structure of transaction data converter
type TransactionDataConverter interface {
	TransactionDataExporter
	TransactionDataImporter
}
