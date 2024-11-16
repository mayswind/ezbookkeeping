package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const nationalBankOfGeorgiaMinimumRequiredContent = "[\n" +
	"  {\n" +
	"    \"date\": \"2024-11-16T00:00:00.000Z\",\n" +
	"    \"currencies\": [\n" +
	"      {\n" +
	"        \"code\": \"JPY\",\n" +
	"        \"quantity\": 100,\n" +
	"        \"rate\": 1.7589,\n" +
	"        \"date\": \"2024-11-15T17:01:11.702Z\"\n" +
	"      },\n" +
	"      {\n" +
	"        \"code\": \"USD\",\n" +
	"        \"quantity\": 1,\n" +
	"        \"rate\": 2.7311,\n" +
	"        \"date\": \"2024-11-15T17:01:11.702Z\"\n" +
	"      }\n" +
	"    ]\n" +
	"  }\n" +
	"]"

func TestNationalBankOfGeorgiaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfGeorgiaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "GEL", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNationalBankOfGeorgiaDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfGeorgiaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731690071), actualLatestExchangeRateResponse.UpdateTime)
}

func TestNationalBankOfGeorgiaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfGeorgiaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "56.853715390300756",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.366152832192157",
	})
}

func TestNationalBankOfGeorgiaDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfGeorgiaDataSource_EmptyData(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[]"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfGeorgiaDataSource_EmptyExchangeRatesData(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[{}]"))
	assert.NotEqual(t, nil, err)

	_, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"date\": \"2024-11-16T00:00:00.000Z\",\n"+
		"    \"currencies\": [\n"+
		"    ]\n"+
		"  }\n"+
		"]"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfGeorgiaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"date\": \"2024-11-16T00:00:00.000Z\",\n"+
		"    \"currencies\": [\n"+
		"      {\n"+
		"        \"code\": \"XXX\",\n"+
		"        \"quantity\": 1,\n"+
		"        \"rate\": 1,\n"+
		"        \"date\": \"2024-11-15T17:01:11.702Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfGeorgiaDataSource_InvalidQuantity(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"date\": \"2024-11-16T00:00:00.000Z\",\n"+
		"    \"currencies\": [\n"+
		"      {\n"+
		"        \"code\": \"USD\",\n"+
		"        \"quantity\": null,\n"+
		"        \"rate\": 2.7311,\n"+
		"        \"date\": \"2024-11-15T17:01:11.702Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"date\": \"2024-11-16T00:00:00.000Z\",\n"+
		"    \"currencies\": [\n"+
		"      {\n"+
		"        \"code\": \"USD\",\n"+
		"        \"quantity\": 0,\n"+
		"        \"rate\": 2.7311,\n"+
		"        \"date\": \"2024-11-15T17:01:11.702Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfGeorgiaDataSource_InvalidRate(t *testing.T) {
	dataSource := &NationalBankOfGeorgiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"date\": \"2024-11-16T00:00:00.000Z\",\n"+
		"    \"currencies\": [\n"+
		"      {\n"+
		"        \"code\": \"USD\",\n"+
		"        \"quantity\": 1,\n"+
		"        \"rate\": null,\n"+
		"        \"date\": \"2024-11-15T17:01:11.702Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"date\": \"2024-11-16T00:00:00.000Z\",\n"+
		"    \"currencies\": [\n"+
		"      {\n"+
		"        \"code\": \"USD\",\n"+
		"        \"quantity\": 1,\n"+
		"        \"rate\": 0,\n"+
		"        \"date\": \"2024-11-15T17:01:11.702Z\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
