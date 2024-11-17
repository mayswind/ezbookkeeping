package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const centralBankOfUzbekistanMinimumRequiredContent = "[\n" +
	"  {\n" +
	"    \"Ccy\": \"USD\",\n" +
	"    \"Nominal\": \"1\",\n" +
	"    \"Rate\": \"12800.13\",\n" +
	"    \"Date\": \"15.11.2024\"\n" +
	"  },\n" +
	"  {\n" +
	"    \"Ccy\": \"VND\",\n" +
	"    \"Nominal\": \"10\",\n" +
	"    \"Rate\": \"5.04\",\n" +
	"    \"Date\": \"15.11.2024\"\n" +
	"  }\n" +
	"]"

func TestCentralBankOfUzbekistanDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfUzbekistanMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "UZS", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCentralBankOfUzbekistanDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfUzbekistanMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731610800), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCentralBankOfUzbekistanDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfUzbekistanMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.07812420655102723",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "VND",
		Rate:     "1984.126984126984",
	})
}

func TestCentralBankOfUzbekistanDataSource_BlankContent(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfUzbekistanDataSource_EmptyData(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[]"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfUzbekistanDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Ccy\": \"XXX\",\n"+
		"    \"Nominal\": \"1\",\n"+
		"    \"Rate\": \"1\",\n"+
		"    \"Date\": \"15.11.2024\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfUzbekistanDataSource_InvalidNominal(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Ccy\": \"USD\",\n"+
		"    \"Nominal\": null,\n"+
		"    \"Rate\": \"12800.13\",\n"+
		"    \"Date\": \"15.11.2024\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Ccy\": \"USD\",\n"+
		"    \"Nominal\": \"0\",\n"+
		"    \"Rate\": \"12800.13\",\n"+
		"    \"Date\": \"15.11.2024\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfUzbekistanDataSource_InvalidRate(t *testing.T) {
	dataSource := &CentralBankOfUzbekistanDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Ccy\": \"USD\",\n"+
		"    \"Nominal\": \"1\",\n"+
		"    \"Rate\": null,\n"+
		"    \"Date\": \"15.11.2024\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Ccy\": \"USD\",\n"+
		"    \"Nominal\": \"1\",\n"+
		"    \"Rate\": \"0\",\n"+
		"    \"Date\": \"15.11.2024\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
