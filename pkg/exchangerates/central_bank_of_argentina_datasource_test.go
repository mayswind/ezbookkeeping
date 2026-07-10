package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const centralBankOfArgentinaMinimumRequiredContent = "{\n" +
	"  \"status\": 200,\n" +
	"  \"results\": {\n" +
	"    \"fecha\": \"2026-07-08\",\n" +
	"    \"detalle\": [\n" +
	"      {\n" +
	"        \"codigoMoneda\": \"EUR\",\n" +
	"        \"descripcion\": \"EURO\",\n" +
	"        \"tipoPase\": 1.142,\n" +
	"        \"tipoCotizacion\": 1699.296\n" +
	"      },\n" +
	"      {\n" +
	"        \"codigoMoneda\": \"USD\",\n" +
	"        \"descripcion\": \"DOLAR E.E.U.U.\",\n" +
	"        \"tipoPase\": 0,\n" +
	"        \"tipoCotizacion\": 1488\n" +
	"      }\n" +
	"    ]\n" +
	"  }\n" +
	"}"

func TestCentralBankOfArgentinaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfArgentinaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "ARS", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCentralBankOfArgentinaDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfArgentinaMinimumRequiredContent))
	assert.Equal(t, nil, err)

	assert.Equal(t, int64(1783479600), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCentralBankOfArgentinaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfArgentinaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.0006720430107526882",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "EUR",
		Rate:     "0.0005884789936538425",
	})
}

func TestCentralBankOfArgentinaDataSource_PerUnitQuantity(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"2026-07-08\",\n" +
		"    \"detalle\": [\n" +
		"      {\n" +
		"        \"codigoMoneda\": \"USD\",\n" +
		"        \"descripcion\": \"DOLAR E.E.U.U.\",\n" +
		"        \"tipoPase\": 0,\n" +
		"        \"tipoCotizacion\": 1488\n" +
		"      },\n" +
		"      {\n" +
		"        \"codigoMoneda\": \"VND\",\n" +
		"        \"descripcion\": \"DONG VIETNAM (C/1.000 UNIDADES)\",\n" +
		"        \"tipoPase\": 0.03803,\n" +
		"        \"tipoCotizacion\": 56.588705\n" +
		"      }\n" +
		"    ]\n" +
		"  }\n" +
		"}"
	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)

	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "VND",
		Rate:     "17.671370991790678",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.0006720430107526882",
	})
}

func TestCentralBankOfArgentinaDataSource_BlankContent(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfArgentinaDataSource_EmptyJsonObject(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfArgentinaDataSource_EmptyResultsJsonObject(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {}\n" +
		"}"

	_, err := dataSource.Parse(context, []byte(content))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfArgentinaDataSource_EmptyDetails(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"2026-07-08\",\n" +
		"    \"detalle\": []\n" +
		"  }\n" +
		"}"

	_, err := dataSource.Parse(context, []byte(content))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfArgentinaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"2026-07-08\",\n" +
		"    \"detalle\": [\n" +
		"      {\n" +
		"        \"codigoMoneda\": \"XXX\",\n" +
		"        \"descripcion\": \"XXX\",\n" +
		"        \"tipoPase\": 0,\n" +
		"        \"tipoCotizacion\": 0\n" +
		"      }\n" +
		"    ]\n" +
		"  }\n" +
		"}"

	actualLatestExchangeRateResponse, _ := dataSource.Parse(context, []byte(content))
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfArgentinaDataSource_InvalidRate(t *testing.T) {
	dataSource := &CentralBankOfArgentinaDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"2026-07-08\",\n" +
		"    \"detalle\": [\n" +
		"      {\n" +
		"        \"codigoMoneda\": \"USD\",\n" +
		"        \"descripcion\": \"DOLAR E.E.U.U.\",\n" +
		"        \"tipoPase\": 0,\n" +
		"        \"tipoCotizacion\": 0\n" +
		"      }\n" +
		"    ]\n" +
		"  }\n" +
		"}"

	actualLatestExchangeRateResponse, _ := dataSource.Parse(context, []byte(content))
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
