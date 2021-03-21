package exporters

import "github.com/mayswind/lab/pkg/models"

// DataExporter defines the structure of data exporter
type DataExporter interface {
	// GetOutputContent returns the exported data
	GetOutputContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexs map[int64][]int64) ([]byte, error)
}
