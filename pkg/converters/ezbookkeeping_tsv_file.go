package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// EzBookKeepingTSVFileExporter defines the structure of TSV file exporter
type EzBookKeepingTSVFileExporter struct {
	EzBookKeepingPlainFileExporter
}

// ToExportedContent returns the exported TSV data
func (e *EzBookKeepingTSVFileExporter) ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexs map[int64][]int64) ([]byte, error) {
	return e.toExportedContent(uid, "\t", transactions, accountMap, categoryMap, tagMap, allTagIndexs)
}
