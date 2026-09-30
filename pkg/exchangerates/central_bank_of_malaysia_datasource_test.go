package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const bankNegaraMalaysiaMinimumRequiredContent = "{\n" +
	"  \"data\": [\n" +
	"    {\n" +
	"      \"currency_code\": \"USD\",\n" +
	"      \"unit\": 1,\n" +
	"      \"rate\": {\n" +
	"        \"date\": \"2026-09-29\",\n" +
	"        \"middle_rate\": 4.08\n" +
	"      }\n" +
	"    },\n" +
	"    {\n" +
	"      \"currency_code\": \"JPY\",\n" +
	"      \"unit\": 100,\n" +
	"      \"rate\": {\n" +
	"        \"date\": \"2026-09-29\",\n" +
	"        \"middle_rate\": 2.5926\n" +
	"      }\n" +
	"    }\n" +
	"  ],\n" +
	"  \"meta\": {\n" +
	"    \"last_updated\": \"2026-09-29 23:01:20\"\n" +
	"  }\n" +
	"}"

func TestCentralBankOfMalaysiaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankNegaraMalaysiaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "MYR", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCentralBankOfMalaysiaDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankNegaraMalaysiaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1790694080), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCentralBankOfMalaysiaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankNegaraMalaysiaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.24509803921568626",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "38.5713183676618",
	})
}

func TestCentralBankOfMalaysiaDataSource_BlankContent(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfMalaysiaDataSource_EmptyData(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfMalaysiaDataSource_EmptyExchangeRatesData(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.NotEqual(t, nil, err)

	_, err = dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfMalaysiaDataSource_InvalidUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"USD\",\n"+
		"      \"unit\": 1,\n"+
		"      \"rate\": {\n"+
		"        \"date\": \"2026-09-29\",\n"+
		"        \"middle_rate\": 4.08\n"+
		"      }\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"invalid\"\n"+
		"  }\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfMalaysiaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"XXX\",\n"+
		"      \"unit\": 1,\n"+
		"      \"rate\": {\n"+
		"        \"date\": \"2026-09-29\",\n"+
		"        \"middle_rate\": 4.08\n"+
		"      }\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfMalaysiaDataSource_InvalidUnit(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"USD\",\n"+
		"      \"unit\": null,\n"+
		"      \"rate\": {\n"+
		"        \"date\": \"2026-09-29\",\n"+
		"        \"middle_rate\": 4.08\n"+
		"      }\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"USD\",\n"+
		"      \"unit\": 0,\n"+
		"      \"rate\": {\n"+
		"        \"date\": \"2026-09-29\",\n"+
		"        \"middle_rate\": 4.08\n"+
		"      }\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfMalaysiaDataSource_InvalidRate(t *testing.T) {
	dataSource := &CentralBankOfMalaysiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"USD\",\n"+
		"      \"unit\": 1,\n"+
		"      \"rate\": null\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"USD\",\n"+
		"      \"unit\": 1,\n"+
		"      \"rate\": {\n"+
		"        \"date\": \"2026-09-29\",\n"+
		"        \"middle_rate\": null\n"+
		"      }\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("{\n"+
		"  \"data\": [\n"+
		"    {\n"+
		"      \"currency_code\": \"USD\",\n"+
		"      \"unit\": 1,\n"+
		"      \"rate\": {\n"+
		"        \"date\": \"2026-09-29\",\n"+
		"        \"middle_rate\": 0\n"+
		"      }\n"+
		"    }\n"+
		"  ],\n"+
		"  \"meta\": {\n"+
		"    \"last_updated\": \"2026-09-29 23:01:20\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
