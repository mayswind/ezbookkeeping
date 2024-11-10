package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionTemplateGetTagIds(t *testing.T) {
	template := &TransactionTemplate{
		TagIds: "1,2,3",
	}

	expectedValue := []int64{1, 2, 3}
	assert.EqualValues(t, expectedValue, template.GetTagIds())
}

func TestTransactionTemplateInfoResponseSliceLess(t *testing.T) {
	var transactionTemplateRespSlice TransactionTemplateInfoResponseSlice
	transactionTemplateRespSlice = append(transactionTemplateRespSlice, &TransactionTemplateInfoResponse{
		TransactionInfoResponse: &TransactionInfoResponse{
			Id: 1,
		},
		DisplayOrder: 3,
	})
	transactionTemplateRespSlice = append(transactionTemplateRespSlice, &TransactionTemplateInfoResponse{
		TransactionInfoResponse: &TransactionInfoResponse{
			Id: 2,
		},
		DisplayOrder: 1,
	})
	transactionTemplateRespSlice = append(transactionTemplateRespSlice, &TransactionTemplateInfoResponse{
		TransactionInfoResponse: &TransactionInfoResponse{
			Id: 3,
		},
		DisplayOrder: 2,
	})

	sort.Sort(transactionTemplateRespSlice)

	assert.Equal(t, int64(2), transactionTemplateRespSlice[0].Id)
	assert.Equal(t, int64(3), transactionTemplateRespSlice[1].Id)
	assert.Equal(t, int64(1), transactionTemplateRespSlice[2].Id)
}
