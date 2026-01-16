package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionTagGroupInfoResponseSliceLess(t *testing.T) {
	var transactionTagGroupRespSlice TransactionTagGroupInfoResponseSlice
	transactionTagGroupRespSlice = append(transactionTagGroupRespSlice, &TransactionTagGroupInfoResponse{
		Id:           1,
		DisplayOrder: 3,
	})
	transactionTagGroupRespSlice = append(transactionTagGroupRespSlice, &TransactionTagGroupInfoResponse{
		Id:           2,
		DisplayOrder: 1,
	})
	transactionTagGroupRespSlice = append(transactionTagGroupRespSlice, &TransactionTagGroupInfoResponse{
		Id:           3,
		DisplayOrder: 2,
	})

	sort.Sort(transactionTagGroupRespSlice)

	assert.Equal(t, int64(2), transactionTagGroupRespSlice[0].Id)
	assert.Equal(t, int64(3), transactionTagGroupRespSlice[1].Id)
	assert.Equal(t, int64(1), transactionTagGroupRespSlice[2].Id)
}
