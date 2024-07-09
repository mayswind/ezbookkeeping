package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// DataConverter defines the structure of data exporter
type DataConverter interface {
	// ToExportedContent returns the exported data
	ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error)
}
