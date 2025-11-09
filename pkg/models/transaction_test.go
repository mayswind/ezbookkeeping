package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestTransactionAmountsRequestGetTransactionAmountsRequestItems(t *testing.T) {
	transactionAmountsRequest := &TransactionAmountsRequest{
		Query: "name1_1234567890_1234567891|name2_1234567900_1234567901",
	}

	var expectedValue []*TransactionAmountsRequestItem
	expectedValue = append(expectedValue, &TransactionAmountsRequestItem{
		Name:      "name1",
		StartTime: 1234567890,
		EndTime:   1234567891,
	})
	expectedValue = append(expectedValue, &TransactionAmountsRequestItem{
		Name:      "name2",
		StartTime: 1234567900,
		EndTime:   1234567901,
	})

	actualValue, err := transactionAmountsRequest.GetTransactionAmountsRequestItems()

	assert.Nil(t, err)
	assert.EqualValues(t, expectedValue, actualValue)
}

func TestTransactionAmountsRequestGetTransactionAmountsRequestItems_InvalidValue(t *testing.T) {
	transactionAmountsRequest := &TransactionAmountsRequest{
		Query: "name1_1234567890",
	}

	_, err := transactionAmountsRequest.GetTransactionAmountsRequestItems()
	assert.NotNil(t, err)
	assert.EqualError(t, err, errs.ErrQueryItemsInvalid.Message)

	transactionAmountsRequest2 := &TransactionAmountsRequest{
		Query: "name1_123456789f_1234567891",
	}

	_, err = transactionAmountsRequest2.GetTransactionAmountsRequestItems()
	assert.NotNil(t, err)

	transactionAmountsRequest3 := &TransactionAmountsRequest{
		Query: "name1_1234567890_123456789f",
	}

	_, err = transactionAmountsRequest3.GetTransactionAmountsRequestItems()
	assert.NotNil(t, err)
}

func TestYearMonthRangeRequestGetNumericYearMonthRange(t *testing.T) {
	yearMonthRangeRequest := &YearMonthRangeRequest{
		StartYearMonth: "2023-4",
		EndYearMonth:   "2024-05",
	}

	startYear, startMonth, endYear, endMonth, err := yearMonthRangeRequest.GetNumericYearMonthRange()

	assert.Nil(t, err)
	assert.Equal(t, int32(2023), startYear)
	assert.Equal(t, int32(4), startMonth)
	assert.Equal(t, int32(2024), endYear)
	assert.Equal(t, int32(5), endMonth)
}

func TestYearMonthRangeRequestGetNumericYearMonthRange_InvalidValues(t *testing.T) {
	yearMonthRangeRequest := &YearMonthRangeRequest{
		StartYearMonth: "2023/4",
		EndYearMonth:   "2024/05",
	}

	_, _, _, _, err := yearMonthRangeRequest.GetNumericYearMonthRange()
	assert.NotNil(t, err)

	yearMonthRangeRequest2 := &YearMonthRangeRequest{
		StartYearMonth: "2023-April",
	}

	_, _, _, _, err = yearMonthRangeRequest2.GetNumericYearMonthRange()
	assert.NotNil(t, err)

	yearMonthRangeRequest3 := &YearMonthRangeRequest{
		EndYearMonth: "2024-May",
	}

	_, _, _, _, err = yearMonthRangeRequest3.GetNumericYearMonthRange()
	assert.NotNil(t, err)
}

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

func TestTransactionStatisticTrendsResponseItemSliceLess(t *testing.T) {
	var transactionTrendsSlice TransactionStatisticTrendsResponseItemSlice
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsResponseItem{
		Year:  2024,
		Month: 9,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsResponseItem{
		Year:  2022,
		Month: 10,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsResponseItem{
		Year:  2023,
		Month: 1,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsResponseItem{
		Year:  2022,
		Month: 2,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticTrendsResponseItem{
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

func TestTransactionStatisticAssetTrendsResponseItemSliceLess(t *testing.T) {
	var transactionTrendsSlice TransactionStatisticAssetTrendsResponseItemSlice
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticAssetTrendsResponseItem{
		Year:  2024,
		Month: 9,
		Day:   1,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticAssetTrendsResponseItem{
		Year:  2024,
		Month: 9,
		Day:   2,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticAssetTrendsResponseItem{
		Year:  2024,
		Month: 10,
		Day:   1,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticAssetTrendsResponseItem{
		Year:  2022,
		Month: 10,
		Day:   1,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticAssetTrendsResponseItem{
		Year:  2023,
		Month: 1,
		Day:   1,
	})
	transactionTrendsSlice = append(transactionTrendsSlice, &TransactionStatisticAssetTrendsResponseItem{
		Year:  2024,
		Month: 2,
		Day:   2,
	})

	sort.Sort(transactionTrendsSlice)

	assert.Equal(t, int32(2022), transactionTrendsSlice[0].Year)
	assert.Equal(t, int32(10), transactionTrendsSlice[0].Month)
	assert.Equal(t, int32(1), transactionTrendsSlice[0].Day)
	assert.Equal(t, int32(2023), transactionTrendsSlice[1].Year)
	assert.Equal(t, int32(1), transactionTrendsSlice[1].Month)
	assert.Equal(t, int32(1), transactionTrendsSlice[1].Day)
	assert.Equal(t, int32(2024), transactionTrendsSlice[2].Year)
	assert.Equal(t, int32(2), transactionTrendsSlice[2].Month)
	assert.Equal(t, int32(2), transactionTrendsSlice[2].Day)
	assert.Equal(t, int32(2024), transactionTrendsSlice[3].Year)
	assert.Equal(t, int32(9), transactionTrendsSlice[3].Month)
	assert.Equal(t, int32(1), transactionTrendsSlice[3].Day)
	assert.Equal(t, int32(2024), transactionTrendsSlice[4].Year)
	assert.Equal(t, int32(9), transactionTrendsSlice[4].Month)
	assert.Equal(t, int32(2), transactionTrendsSlice[4].Day)
	assert.Equal(t, int32(2024), transactionTrendsSlice[5].Year)
	assert.Equal(t, int32(10), transactionTrendsSlice[5].Month)
	assert.Equal(t, int32(1), transactionTrendsSlice[5].Day)
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
