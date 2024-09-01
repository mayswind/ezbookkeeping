package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// DataConverter defines the structure of data converter
type DataConverter interface {
	// ToExportedContent returns the exported data
	ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error)

	// ParseImportedData returns the imported data
	ParseImportedData(user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) ([]*models.Transaction, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error)
}
