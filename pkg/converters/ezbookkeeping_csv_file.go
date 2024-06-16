package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// EzBookKeepingCSVFileExporter defines the structure of CSV file exporter
type EzBookKeepingCSVFileExporter struct {
	EzBookKeepingPlainFileExporter
}

const csvSeparator = ","

// ToExportedContent returns the exported CSV data
func (e *EzBookKeepingCSVFileExporter) ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexs map[int64][]int64) ([]byte, error) {
	return e.toExportedContent(uid, csvSeparator, transactions, accountMap, categoryMap, tagMap, allTagIndexs)
}
