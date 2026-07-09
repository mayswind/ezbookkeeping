package exchangerates

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const centralBankOfTheArgentineRepublicMinimumRequiredContent = "{\n" +
	"  \"status\": 200,\n" +
	"  \"results\": {\n" +
	"    \"fecha\": \"2026-07-08\",\n" +
	"    \"detalle\": [\n" +
	"      {\n" +
	"        \"codigoMoneda\": \"EUR\",\n" +
	"        \"descripcion\": \"EURO\",\n" +
	"        \"tipoPase\": 0.672,\n" +
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

func TestCentralBankOfTheArgentineRepublicDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfTheArgentineRepublicMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "ARS", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCentralBankOfTheArgentineRepublicDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfTheArgentineRepublicMinimumRequiredContent))
	assert.Equal(t, nil, err)

	expectedUpdateTime, _ := time.Parse("2006-01-02", "2026-07-08")
	assert.Equal(t, expectedUpdateTime.Unix(), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCentralBankOfTheArgentineRepublicDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfTheArgentineRepublicMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     utils.Float64ToString(1 / 1488.0),
	})
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
}

func TestCentralBankOfTheArgentineRepublicDataSource_BuildRequests(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}

	requests, err := dataSource.BuildRequests()
	assert.Equal(t, nil, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, http.MethodGet, requests[0].Method)
	assert.Equal(t, centralBankOfTheArgentineRepublicExchangeRateUrl, requests[0].URL.String())
}

func TestCentralBankOfTheArgentineRepublicDataSource_BlankContent(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestCentralBankOfTheArgentineRepublicDataSource_EmptyJsonObject(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestCentralBankOfTheArgentineRepublicDataSource_EmptyDetalle(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"2026-07-08\",\n" +
		"    \"detalle\": []\n" +
		"  }\n" +
		"}"

	_, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestCentralBankOfTheArgentineRepublicDataSource_MissingUSD(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"2026-07-08\",\n" +
		"    \"detalle\": [\n" +
		"      {\n" +
		"        \"codigoMoneda\": \"EUR\",\n" +
		"        \"descripcion\": \"EURO\",\n" +
		"        \"tipoPase\": 0.672,\n" +
		"        \"tipoCotizacion\": 1699.296\n" +
		"      }\n" +
		"    ]\n" +
		"  }\n" +
		"}"

	_, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestCentralBankOfTheArgentineRepublicDataSource_InvalidTipoCotizacion(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
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

	_, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestCentralBankOfTheArgentineRepublicDataSource_InvalidFecha(t *testing.T) {
	dataSource := &CentralBankOfTheArgentineRepublicDataSource{}
	context := core.NewNullContext()

	content := "{\n" +
		"  \"status\": 200,\n" +
		"  \"results\": {\n" +
		"    \"fecha\": \"not-a-date\",\n" +
		"    \"detalle\": [\n" +
		"      {\n" +
		"        \"codigoMoneda\": \"USD\",\n" +
		"        \"descripcion\": \"DOLAR E.E.U.U.\",\n" +
		"        \"tipoPase\": 0,\n" +
		"        \"tipoCotizacion\": 1488\n" +
		"      }\n" +
		"    ]\n" +
		"  }\n" +
		"}"

	_, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}
