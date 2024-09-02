package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestExchangeRateSliceLess(t *testing.T) {
	var latestExchangeRateSlice LatestExchangeRateSlice
	latestExchangeRateSlice = append(latestExchangeRateSlice, &LatestExchangeRate{
		Currency: "USD",
	})
	latestExchangeRateSlice = append(latestExchangeRateSlice, &LatestExchangeRate{
		Currency: "EUR",
	})
	latestExchangeRateSlice = append(latestExchangeRateSlice, &LatestExchangeRate{
		Currency: "CNY",
	})

	sort.Sort(latestExchangeRateSlice)

	assert.Equal(t, "CNY", latestExchangeRateSlice[0].Currency)
	assert.Equal(t, "EUR", latestExchangeRateSlice[1].Currency)
	assert.Equal(t, "USD", latestExchangeRateSlice[2].Currency)
}
