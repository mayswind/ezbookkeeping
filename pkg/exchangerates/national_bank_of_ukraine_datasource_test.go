package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const nationalBankOfUkraineMinimumRequiredContent = "[\n" +
	"    {\n" +
	"        \"StartDate\": \"21.04.2025\",\n" +
	"        \"TimeSign\": \"0000\",\n" +
	"        \"CurrencyCode\": \"840\",\n" +
	"        \"CurrencyCodeL\": \"USD\",\n" +
	"        \"Units\": 1,\n" +
	"        \"Amount\": 41.3955\n" +
	"    },\n" +
	"    {\n" +
	"        \"StartDate\": \"21.04.2025\",\n" +
	"        \"TimeSign\": \"0000\",\n" +
	"        \"CurrencyCode\": \"392\",\n" +
	"        \"CurrencyCodeL\": \"JPY\",\n" +
	"        \"Units\": 10,\n" +
	"        \"Amount\": 2.907\n" +
	"    }\n" +
	"]"

func TestNationalBankOfUkraineDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfUkraineMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "UAH", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNationalBankOfUkraineDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfUkraineMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1745193600), actualLatestExchangeRateResponse.UpdateTime)
}

func TestNationalBankOfUkraineDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfUkraineMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.02415721515623679",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "3.4399724802201583",
	})
}

func TestNationalBankOfUkraineDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfUkraineDataSource_EmptyData(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[]"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfUkraineDataSource_InvalidDate(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[\n"+
		"    {\n"+
		"        \"StartDate\": \"04.21.2025\",\n"+
		"        \"TimeSign\": \"0000\",\n"+
		"        \"CurrencyCode\": \"840\",\n"+
		"        \"CurrencyCodeL\": \"USD\",\n"+
		"        \"Units\": 1,\n"+
		"        \"Amount\": 41.3955\n"+
		"    }\n"+
		"]"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfUkraineDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"    {\n"+
		"        \"StartDate\": \"21.04.2025\",\n"+
		"        \"TimeSign\": \"0000\",\n"+
		"        \"CurrencyCode\": \"840\",\n"+
		"        \"CurrencyCodeL\": \"XXX\",\n"+
		"        \"Units\": 1,\n"+
		"        \"Amount\": 41.3955\n"+
		"    }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfUkraineDataSource_InvalidUnits(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"    {\n"+
		"        \"StartDate\": \"21.04.2025\",\n"+
		"        \"TimeSign\": \"0000\",\n"+
		"        \"CurrencyCode\": \"840\",\n"+
		"        \"CurrencyCodeL\": \"USD\",\n"+
		"        \"Units\": null,\n"+
		"        \"Amount\": 41.3955\n"+
		"    }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"    {\n"+
		"        \"StartDate\": \"21.04.2025\",\n"+
		"        \"TimeSign\": \"0000\",\n"+
		"        \"CurrencyCode\": \"840\",\n"+
		"        \"CurrencyCodeL\": \"USD\",\n"+
		"        \"Units\": 0,\n"+
		"        \"Amount\": 41.3955\n"+
		"    }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfUkraineDataSource_InvalidAmount(t *testing.T) {
	dataSource := &NationalBankOfUkraineDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"    {\n"+
		"        \"StartDate\": \"21.04.2025\",\n"+
		"        \"TimeSign\": \"0000\",\n"+
		"        \"CurrencyCode\": \"840\",\n"+
		"        \"CurrencyCodeL\": \"USD\",\n"+
		"        \"Units\": 1,\n"+
		"        \"Amount\": null\n"+
		"    }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"    {\n"+
		"        \"StartDate\": \"21.04.2025\",\n"+
		"        \"TimeSign\": \"0000\",\n"+
		"        \"CurrencyCode\": \"840\",\n"+
		"        \"CurrencyCodeL\": \"USD\",\n"+
		"        \"Units\": 1,\n"+
		"        \"Amount\": 0\n"+
		"    }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
