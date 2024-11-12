package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const internationalMonetaryFundMinimumRequiredContent = "SDRs per Currency unit and Currency units per SDR (1)\n" +
	"last five days\n" +
	"SDRs per Currency unit (2)\n" +
	"\n" +
	"Currency\tAugust 28, 2024\tAugust 27, 2024\tAugust 26, 2024\tAugust 23, 2024\tAugust 22, 2024\n" +
	"Chinese yuan\t0.1040520000\t0.1039250000\t0.1040370000\t0.1040850000\t0.1040570000\n" +
	"U.S. dollar\t0.7417320000\t0.7410250000\t0.7408270000\t0.7429280000\t0.7423020000\n"

func TestInternationalMonetaryFundDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(internationalMonetaryFundMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "USD", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestInternationalMonetaryFundDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(internationalMonetaryFundMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1724857200), actualLatestExchangeRateResponse.UpdateTime)
}

func TestInternationalMonetaryFundDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(internationalMonetaryFundMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "1",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "7.128474224426247",
	})
}

func TestInternationalMonetaryFundDataSource_BlankContent(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestInternationalMonetaryFundDataSource_OnlyHeader(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("SDRs per Currency unit and Currency units per SDR (1)\n"+
		"last five days\n"+
		"SDRs per Currency unit (2)"))
	assert.NotEqual(t, nil, err)
}

func TestInternationalMonetaryFundDataSource_OnlyHeaderAndTitle(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("SDRs per Currency unit and Currency units per SDR (1)\n"+
		"last five days\n"+
		"SDRs per Currency unit (2)\n"+
		"\n"+
		"Currency\tAugust 28, 2024\tAugust 27, 2024\tAugust 26, 2024\tAugust 23, 2024\tAugust 22, 2024\n"))
	assert.NotEqual(t, nil, err)
}

func TestInternationalMonetaryFundDataSource_MissingHeader(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("Currency\tAugust 28, 2024\tAugust 27, 2024\tAugust 26, 2024\tAugust 23, 2024\tAugust 22, 2024\n"+
		"Chinese yuan\t0.1040520000\t0.1039250000\t0.1040370000\t0.1040850000\t0.1040570000\n"+
		"U.S. dollar\t0.7417320000\t0.7410250000\t0.7408270000\t0.7429280000\t0.7423020000\n"))
	assert.NotEqual(t, nil, err)
}

func TestInternationalMonetaryFundDataSource_MissingDefaultCurrencyData(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("SDRs per Currency unit and Currency units per SDR (1)\n"+
		"last five days\n"+
		"SDRs per Currency unit (2)\n"+
		"\n"+
		"Currency\tAugust 28, 2024\tAugust 27, 2024\tAugust 26, 2024\tAugust 23, 2024\tAugust 22, 2024\n"+
		"Chinese yuan\t0.1040520000\t0.1039250000\t0.1040370000\t0.1040850000\t0.1040570000\n"))
	assert.NotEqual(t, nil, err)
}

func TestInternationalMonetaryFundDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("SDRs per Currency unit and Currency units per SDR (1)\n"+
		"last five days\n"+
		"SDRs per Currency unit (2)\n"+
		"\n"+
		"Currency\tAugust 28, 2024\tAugust 27, 2024\tAugust 26, 2024\tAugust 23, 2024\tAugust 22, 2024\n"+
		"Foo bar\t0.1040520000\t0.1039250000\t0.1040370000\t0.1040850000\t0.1040570000\n"+
		"U.S. dollar\t0.7417320000\t0.7410250000\t0.7408270000\t0.7429280000\t0.7423020000\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
}

func TestInternationalMonetaryFundDataSource_LatestDateNotHasRate(t *testing.T) {
	dataSource := &InternationalMonetaryFundDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("SDRs per Currency unit and Currency units per SDR (1)\n"+
		"last five days\n"+
		"SDRs per Currency unit (2)\n"+
		"\n"+
		"Currency\tAugust 28, 2024\tAugust 27, 2024\tAugust 26, 2024\tAugust 23, 2024\tAugust 22, 2024\n"+
		"U.S. dollar\t0.7417320000\t0.7410250000\t0.7408270000\t0.7429280000\t0.7423020000\n"+
		"U.A.E. dirham\t\t0.2017770000\t0.2017230000\t\t0.2021240000\n"))

	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "AED",
		Rate:     "3.675998751096507",
	})
}
