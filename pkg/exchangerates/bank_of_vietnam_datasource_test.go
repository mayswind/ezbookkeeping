package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const techcombankMinimumRequiredContent = "{\n" +
	"  \"exchangeRate\": {\n" +
	"    \"data\": [\n" +
	"      {\n" +
	"        \"label\": \"USD (1,2)\",\n" +
	"        \"askRate\": \"25900.0\",\n" +
	"        \"sourceCurrency\": \"USD\",\n" +
	"        \"targetCurrency\": \"VND\",\n" +
	"        \"inputDate\": \"2025-07-02T17:02:00.638Z\"\n" +
	"      },\n" +
	"      {\n" +
	"        \"label\": \"EUR\",\n" +
	"        \"askRate\": \"30441.0\",\n" +
	"        \"sourceCurrency\": \"EUR\",\n" +
	"        \"targetCurrency\": \"VND\",\n" +
	"        \"inputDate\": \"2025-07-02T17:02:00.638Z\"\n" +
	"      }\n" +
	"    ]\n" +
	"  },\n" +
	"  \"goldRate\": {\n" +
	"    \"data\": [\n" +
	"      {\n" +
	"        \"askRate\": \"12100000.0\"\n" +
	"      }\n" +
	"    ]\n" +
	"  }\n" +
	"}"

func TestInternationalTechcombankDataSource_BuildRequests(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}

	requests, err := dataSource.BuildRequests()
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "GET", requests[0].Method)
	assert.Equal(t, techcombankExchangeRateUrl, requests[0].URL.String())
	assert.Equal(t, "", requests[0].Header.Get("User-Agent"))
}

func TestInternationalTechcombankDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(techcombankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "VND", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestInternationalTechcombankDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(techcombankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.LessOrEqual(t, int64(1705312800), actualLatestExchangeRateResponse.UpdateTime)
}

func TestInternationalTechcombankDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(techcombankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.00003861003861003861",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "EUR",
		Rate:     "0.000032850431983180576",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "XAU",
		Rate:     "0.00000008264462809917355",
	})
}

func TestInternationalTechcombankDataSource_BlankContent(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestInternationalTechcombankDataSource_EmptyJsonObject(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestInternationalTechcombankDataSource_EmptyExchangeRateData(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": []\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": []\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestInternationalTechcombankDataSource_EmptyAskRate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"\",\n"+
		"        \"sourceCurrency\": \"USD\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"2024-01-15T10:00:00Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": []\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestInternationalTechcombankDataSource_EmptySourceCurrency(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"25900.0\",\n"+
		"        \"sourceCurrency\": \"\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"2025-07-02T17:02:00.638Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": []\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestInternationalTechcombankDataSource_InvalidAskRate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"invalid\",\n"+
		"        \"sourceCurrency\": \"USD\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"2024-01-15T10:00:00Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": []\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestInternationalTechcombankDataSource_EmptyGoldRate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"25900.0\",\n"+
		"        \"sourceCurrency\": \"USD\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"2025-07-02T17:02:00.638Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"askRate\": \"\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1) // Only USD, no gold
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.00003861003861003861",
	})
}

func TestInternationalTechcombankDataSource_InvalidGoldRate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"24500.0\",\n"+
		"        \"sourceCurrency\": \"USD\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"2024-01-15T10:00:00Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"askRate\": \"invalid\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestInternationalTechcombankDataSource_OnlyGoldRate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": []\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"askRate\": \"75000000.0\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "XAU",
		Rate:     "0.000000013333333333333334",
	})
}

func TestInternationalTechcombankDataSource_InvalidInputDate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"25900.0\",\n"+
		"        \"sourceCurrency\": \"USD\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"invalid-date\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": []\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
	// Should still parse exchange rate but fall back to current time for update time
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.00003861003861003861",
	})
}

func TestInternationalTechcombankDataSource_EmptyInputDate(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"exchangeRate\": {\n"+
		"    \"data\": [\n"+
		"      {\n"+
		"        \"label\": \"USD\",\n"+
		"        \"askRate\": \"25900.0\",\n"+
		"        \"sourceCurrency\": \"USD\",\n"+
		"        \"targetCurrency\": \"VND\",\n"+
		"        \"inputDate\": \"\"\n"+
		"      }\n"+
		"    ]\n"+
		"  },\n"+
		"  \"goldRate\": {\n"+
		"    \"data\": []\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
	// Should still parse exchange rate but fall back to current time for update time
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.00003861003861003861",
	})
}

func TestInternationalTechcombankDataSource_InvalidJsonFormat(t *testing.T) {
	dataSource := &InternationalTechcombankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("invalid json"))
	assert.NotEqual(t, nil, err)
}
