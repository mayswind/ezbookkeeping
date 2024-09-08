package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportTransactionSliceLess(t *testing.T) {
	var transactionSlice ImportedTransactionSlice
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   1,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   2,
			Type:            TRANSACTION_DB_TYPE_INCOME,
			TransactionTime: 2,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   3,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 10,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   4,
			Type:            TRANSACTION_DB_TYPE_TRANSFER_IN,
			TransactionTime: 3,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   5,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 11,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   6,
			Type:            TRANSACTION_DB_TYPE_TRANSFER_OUT,
			TransactionTime: 4,
		},
	})

	sort.Sort(transactionSlice)

	assert.Equal(t, int64(3), transactionSlice[0].TransactionId)
	assert.Equal(t, int64(5), transactionSlice[1].TransactionId)
	assert.Equal(t, int64(1), transactionSlice[2].TransactionId)
	assert.Equal(t, int64(2), transactionSlice[3].TransactionId)
	assert.Equal(t, int64(4), transactionSlice[4].TransactionId)
	assert.Equal(t, int64(6), transactionSlice[5].TransactionId)
}
