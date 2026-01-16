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
		TagGroupId:   0,
		DisplayOrder: 3,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           2,
		TagGroupId:   1,
		DisplayOrder: 2,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           3,
		TagGroupId:   0,
		DisplayOrder: 1,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           4,
		TagGroupId:   2,
		DisplayOrder: 1,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           5,
		TagGroupId:   1,
		DisplayOrder: 1,
	})
	transactionTagRespSlice = append(transactionTagRespSlice, &TransactionTagInfoResponse{
		Id:           6,
		TagGroupId:   0,
		DisplayOrder: 2,
	})

	sort.Sort(transactionTagRespSlice)

	assert.Equal(t, int64(3), transactionTagRespSlice[0].Id)
	assert.Equal(t, int64(6), transactionTagRespSlice[1].Id)
	assert.Equal(t, int64(1), transactionTagRespSlice[2].Id)
	assert.Equal(t, int64(5), transactionTagRespSlice[3].Id)
	assert.Equal(t, int64(2), transactionTagRespSlice[4].Id)
	assert.Equal(t, int64(4), transactionTagRespSlice[5].Id)
}
