package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const bankOfItalyMinimumRequiredContent = "{\n" +
	"  \"latestRates\": [\n" +
	"    {\n" +
	"      \"isoCode\": \"USD\",\n" +
	"      \"eurRate\": \"1.25\",\n" +
	"      \"referenceDate\": \"2026-10-02\"\n" +
	"    },\n" +
	"    {\n" +
	"      \"isoCode\": \"JPY\",\n" +
	"      \"eurRate\": \"160\",\n" +
	"      \"referenceDate\": \"2026-10-02\"\n" +
	"    }\n" +
	"  ]\n" +
	"}"

func TestBankOfItalyDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfItalyMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "EUR", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestBankOfItalyDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfItalyMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1790892000), actualLatestExchangeRateResponse.UpdateTime)
}

func TestBankOfItalyDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfItalyMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "1.25",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "160",
	})
}

func TestBankOfItalyDataSource_BlankContent(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestBankOfItalyDataSource_EmptyData(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfItalyDataSource_EmptyExchangeRatesData(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"resultsInfo\": {}\n"+
		"}"))
	assert.NotEqual(t, nil, err)

	_, err = dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": []\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfItalyDataSource_InvalidUpdateTime(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": [\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"1.25\",\n"+
		"      \"referenceDate\": \"invalid\"\n"+
		"    }\n"+
		"  ]\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfItalyDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": [\n"+
		"    {\n"+
		"      \"isoCode\": \"XXX\",\n"+
		"      \"eurRate\": \"1.25\",\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    }\n"+
		"  ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfItalyDataSource_InvalidRate(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": [\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": null,\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    }\n"+
		"  ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": [\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"0\",\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    }\n"+
		"  ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": [\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"N.A.\",\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    }\n"+
		"  ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfItalyDataSource_DuplicateCurrencies(t *testing.T) {
	dataSource := &BankOfItalyDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"latestRates\": [\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"2\",\n"+
		"      \"referenceDate\": \"2026-10-01\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"1.25\",\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"2\",\n"+
		"      \"referenceDate\": \"2026-10-01\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"isoCode\": \"USD\",\n"+
		"      \"eurRate\": \"1.25\",\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"isoCode\": \"JPY\",\n"+
		"      \"eurRate\": \"160\",\n"+
		"      \"referenceDate\": \"2026-10-02\"\n"+
		"    }\n"+
		"  ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 2)
	assert.Equal(t, int64(1790892000), actualLatestExchangeRateResponse.UpdateTime)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "1.25",
	})
}
