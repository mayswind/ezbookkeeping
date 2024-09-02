package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionTagInfoResponseSliceLess(t *testing.T) {
	var transactionTagRespSlice TransactionTagInfoResponseSlice
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           1,
		DisplayOrder: 3,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           2,
		DisplayOrder: 1,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           3,
		DisplayOrder: 2,
	})

	sort.Sort(transactionTagRespSlice)

	assert.Equal(t, int64(2), transactionTagRespSlice[0].Id)
	assert.Equal(t, int64(3), transactionTagRespSlice[1].Id)
	assert.Equal(t, int64(1), transactionTagRespSlice[2].Id)
}
