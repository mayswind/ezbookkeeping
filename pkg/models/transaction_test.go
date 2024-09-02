package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionInfoResponseSliceLess(t *testing.T) {
	var transactionRespSlice TransactionInfoResponseSlice
	transactionRespSlice = append(transactionRespSlice, &TransactionInfoResponse{
		Id:   2,
		Time: 3,
	})
	transactionRespSlice = append(transactionRespSlice, &TransactionInfoResponse{
		Id:   3,
		Time: 2,
	})
	transactionRespSlice = append(transactionRespSlice, &TransactionInfoResponse{
		Id:   5,
		Time: 2,
	})
	transactionRespSlice = append(transactionRespSlice, &TransactionInfoResponse{
		Id:   4,
		Time: 1,
	})
	transactionRespSlice = append(transactionRespSlice, &TransactionInfoResponse{
		Id:   1,
		Time: 3,
	})

	sort.Sort(transactionRespSlice)

	assert.Equal(t, int64(2), transactionRespSlice[0].Id)
	assert.Equal(t, int64(1), transactionRespSlice[1].Id)
	assert.Equal(t, int64(5), transactionRespSlice[2].Id)
	assert.Equal(t, int64(3), transactionRespSlice[3].Id)
	assert.Equal(t, int64(4), transactionRespSlice[4].Id)
}

func TestTransactionStatisticTrendsItemSliceLess(t *testing.T) {
	var transactionTrendsSlice TransactionStatisticTrendsItemSlice
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsItem{
		Year:  2024,
		Month: 9,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsItem{
		Year:  2022,
		Month: 10,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsItem{
		Year:  2023,
		Month: 1,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsItem{
		Year:  2022,
		Month: 2,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsItem{
		Year:  2024,
		Month: 1,
	})

	sort.Sort(transactionTrendsSlice)

	assert.Equal(t, int32(2022), transactionTrendsSlice[0].Year)
	assert.Equal(t, int32(2), transactionTrendsSlice[0].Month)
	assert.Equal(t, int32(2022), transactionTrendsSlice[1].Year)
	assert.Equal(t, int32(10), transactionTrendsSlice[1].Month)
	assert.Equal(t, int32(2023), transactionTrendsSlice[2].Year)
	assert.Equal(t, int32(1), transactionTrendsSlice[2].Month)
	assert.Equal(t, int32(2024), transactionTrendsSlice[3].Year)
	assert.Equal(t, int32(1), transactionTrendsSlice[3].Month)
	assert.Equal(t, int32(2024), transactionTrendsSlice[4].Year)
	assert.Equal(t, int32(9), transactionTrendsSlice[4].Month)
}

func TestTransactionAmountsResponseItemAmountInfoSliceLess(t *testing.T) {
	var amountInfoSlice TransactionAmountsResponseItemAmountInfoSlice
	amountInfoSlice = append(amountInfoSlice, &TransactionAmountsResponseItemAmountInfo{
		Currency: "USD",
	})
	amountInfoSlice = append(amountInfoSlice, &TransactionAmountsResponseItemAmountInfo{
		Currency: "EUR",
	})
	amountInfoSlice = append(amountInfoSlice, &TransactionAmountsResponseItemAmountInfo{
		Currency: "CNY",
	})

	sort.Sort(amountInfoSlice)

	assert.Equal(t, "CNY", amountInfoSlice[0].Currency)
	assert.Equal(t, "EUR", amountInfoSlice[1].Currency)
	assert.Equal(t, "USD", amountInfoSlice[2].Currency)
}
