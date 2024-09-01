package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// EzBookKeepingTSVFileConverter defines the structure of TSV file converter
type EzBookKeepingTSVFileConverter struct {
	EzBookKeepingPlainFileConverter
}

const tsvSeparator = "\t"

// ToExportedContent returns the exported TSV data
func (e *EzBookKeepingTSVFileConverter) ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error) {
	return e.toExportedContent(uid, tsvSeparator, transactions, accountMap, categoryMap, tagMap, allTagIndexes)
}

// ParseImportedData parses transactions of ezbookkeeping TSV data
func (e *EzBookKeepingTSVFileConverter) ParseImportedData(user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) ([]*models.Transaction, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	return e.parseImportedData(user, tsvSeparator, data, defaultTimezoneOffset, accountMap, categoryMap, tagMap)
}
