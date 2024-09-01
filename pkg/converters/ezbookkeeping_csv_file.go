package converters

import (
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// EzBookKeepingCSVFileConverter defines the structure of CSV file converter
type EzBookKeepingCSVFileConverter struct {
	EzBookKeepingPlainFileConverter
}

const csvSeparator = ","

// ToExportedContent returns the exported CSV data
func (e *EzBookKeepingCSVFileConverter) ToExportedContent(uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexes map[int64][]int64) ([]byte, error) {
	return e.toExportedContent(uid, csvSeparator, transactions, accountMap, categoryMap, tagMap, allTagIndexes)
}

// ParseImportedData parses transactions of ezbookkeeping CSV data
func (e *EzBookKeepingCSVFileConverter) ParseImportedData(user *models.User, data []byte, defaultTimezoneOffset int16, accountMap map[string]*models.Account, categoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) ([]*models.Transaction, []*models.Account, []*models.TransactionCategory, []*models.TransactionTag, error) {
	return e.parseImportedData(user, csvSeparator, data, defaultTimezoneOffset, accountMap, categoryMap, tagMap)
}
