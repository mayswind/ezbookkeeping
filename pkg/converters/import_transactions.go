package converters

import "github.com/mayswind/ezbookkeeping/pkg/models"

// ImportTransactionSlice represents the slice data structure of import transaction data
type ImportTransactionSlice []*models.Transaction

// Len returns the count of items
func (s ImportTransactionSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ImportTransactionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ImportTransactionSlice) Less(i, j int) bool {
	if s[i].Type != s[j].Type && (s[i].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE || s[j].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE) {
		if s[i].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			return true
		} else if s[j].Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			return false
		}
	}

	return s[i].TransactionTime < s[j].TransactionTime
}
