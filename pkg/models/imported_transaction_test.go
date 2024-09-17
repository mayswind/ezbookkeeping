package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportTransactionSliceLess_NoSameTransactionTime(t *testing.T) {
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

func TestImportTransactionSliceLess_SameTransactionTime(t *testing.T) {
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
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   3,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   4,
			Type:            TRANSACTION_DB_TYPE_TRANSFER_IN,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   5,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   6,
			Type:            TRANSACTION_DB_TYPE_TRANSFER_OUT,
			TransactionTime: 1,
		},
	})

	sort.Sort(transactionSlice)

	assert.Equal(t, int64(3), transactionSlice[0].TransactionId)
	assert.Equal(t, int64(5), transactionSlice[1].TransactionId)
	assert.Equal(t, int64(2), transactionSlice[2].TransactionId)
	assert.Equal(t, int64(1), transactionSlice[3].TransactionId)
	assert.Equal(t, int64(6), transactionSlice[4].TransactionId)
	assert.Equal(t, int64(4), transactionSlice[5].TransactionId)
}

func TestImportTransactionSliceLess_SameTransactionTimeAndSameType(t *testing.T) {
	var transactionSlice ImportedTransactionSlice
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   1,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
		OriginalCategoryName: "3",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   2,
			Type:            TRANSACTION_DB_TYPE_INCOME,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   3,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   4,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
		OriginalCategoryName: "1",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   5,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   6,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
		OriginalCategoryName: "2",
	})

	sort.Sort(transactionSlice)

	assert.Equal(t, int64(3), transactionSlice[0].TransactionId)
	assert.Equal(t, int64(5), transactionSlice[1].TransactionId)
	assert.Equal(t, int64(2), transactionSlice[2].TransactionId)
	assert.Equal(t, int64(4), transactionSlice[3].TransactionId)
	assert.Equal(t, int64(6), transactionSlice[4].TransactionId)
	assert.Equal(t, int64(1), transactionSlice[5].TransactionId)
}

func TestImportTransactionSliceLess_SameTransactionTimeSameTypeAndSameSubCategoryName(t *testing.T) {
	var transactionSlice ImportedTransactionSlice
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   1,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "b",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   2,
			Type:            TRANSACTION_DB_TYPE_INCOME,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   3,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   4,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "c",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   5,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   6,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "a",
	})

	sort.Sort(transactionSlice)

	assert.Equal(t, int64(3), transactionSlice[0].TransactionId)
	assert.Equal(t, int64(5), transactionSlice[1].TransactionId)
	assert.Equal(t, int64(2), transactionSlice[2].TransactionId)
	assert.Equal(t, int64(6), transactionSlice[3].TransactionId)
	assert.Equal(t, int64(1), transactionSlice[4].TransactionId)
	assert.Equal(t, int64(4), transactionSlice[5].TransactionId)
}

func TestImportTransactionSliceLess_SameTransactionTimeSameTypeSameSubCategoryNameAndSameAccountName(t *testing.T) {
	var transactionSlice ImportedTransactionSlice
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   1,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
			Amount:          3,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "a",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   2,
			Type:            TRANSACTION_DB_TYPE_INCOME,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   3,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   4,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
			Amount:          2,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "a",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   5,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 1,
		},
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   6,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
			Amount:          1,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "a",
	})

	sort.Sort(transactionSlice)

	assert.Equal(t, int64(3), transactionSlice[0].TransactionId)
	assert.Equal(t, int64(5), transactionSlice[1].TransactionId)
	assert.Equal(t, int64(2), transactionSlice[2].TransactionId)
	assert.Equal(t, int64(6), transactionSlice[3].TransactionId)
	assert.Equal(t, int64(4), transactionSlice[4].TransactionId)
	assert.Equal(t, int64(1), transactionSlice[5].TransactionId)
}

func TestImportTransactionSliceLess_SameTransactionTimeSameTypeSameSubCategoryNameSameAccountNameAndSameAmount(t *testing.T) {
	var transactionSlice ImportedTransactionSlice
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   1,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 2,
			Amount:          3,
			Comment:         "2",
		},
		OriginalCategoryName:      "2",
		OriginalSourceAccountName: "b",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   2,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 2,
			Amount:          7,
		},
		OriginalCategoryName:      "2",
		OriginalSourceAccountName: "a",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   3,
			Type:            TRANSACTION_DB_TYPE_MODIFY_BALANCE,
			TransactionTime: 8,
		},
		OriginalCategoryName:      "9",
		OriginalSourceAccountName: "x",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   4,
			Type:            TRANSACTION_DB_TYPE_INCOME,
			TransactionTime: 2,
			Amount:          6,
		},
		OriginalCategoryName:      "7",
		OriginalSourceAccountName: "z",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   5,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 2,
			Amount:          3,
			Comment:         "1",
		},
		OriginalCategoryName:      "2",
		OriginalSourceAccountName: "b",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   6,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 2,
			Amount:          1,
		},
		OriginalCategoryName:      "2",
		OriginalSourceAccountName: "b",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   7,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 2,
			Amount:          6,
		},
		OriginalCategoryName:      "1",
		OriginalSourceAccountName: "b",
	})
	transactionSlice = append(transactionSlice, &ImportTransaction{
		Transaction: &Transaction{
			TransactionId:   8,
			Type:            TRANSACTION_DB_TYPE_EXPENSE,
			TransactionTime: 1,
			Amount:          9,
		},
		OriginalCategoryName:      "9",
		OriginalSourceAccountName: "y",
	})

	sort.Sort(transactionSlice)

	assert.Equal(t, int64(3), transactionSlice[0].TransactionId)
	assert.Equal(t, int64(8), transactionSlice[1].TransactionId)
	assert.Equal(t, int64(4), transactionSlice[2].TransactionId)
	assert.Equal(t, int64(7), transactionSlice[3].TransactionId)
	assert.Equal(t, int64(2), transactionSlice[4].TransactionId)
	assert.Equal(t, int64(6), transactionSlice[5].TransactionId)
	assert.Equal(t, int64(5), transactionSlice[6].TransactionId)
	assert.Equal(t, int64(1), transactionSlice[7].TransactionId)
}
