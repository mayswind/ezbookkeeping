package converters

import "github.com/mayswind/ezbookkeeping/pkg/models"

// ImportedTransactionSlice represents the slice data structure of import transaction data
type ImportedTransactionSlice []*models.Transaction

// Len returns the count of items
func (s ImportedTransactionSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ImportedTransactionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ImportedTransactionSlice) Less(i, j int) bool {
	if s[i].Type != s[j].Type && (s[i].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE || s[j].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE) {
		if s[i].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			return true
		} else if s[j].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			return false
		}
	}

	return s[i].TransactionTime < s[j].TransactionTime
}
