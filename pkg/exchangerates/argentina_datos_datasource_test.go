package exchangerates

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const argentinaDatosMinimumRequiredContent = "[\n" +
	"  {\n" +
	"    \"casa\": \"blue\",\n" +
	"    \"compra\": 1490,\n" +
	"    \"venta\": 1510,\n" +
	"    \"fecha\": \"2026-07-06\"\n" +
	"  },\n" +
	"  {\n" +
	"    \"casa\": \"blue\",\n" +
	"    \"compra\": 1495,\n" +
	"    \"venta\": 1515,\n" +
	"    \"fecha\": \"2026-07-07\"\n" +
	"  }\n" +
	"]"

func TestMapConfigRateTypeToArgentinaDatosApiRateType(t *testing.T) {
	assert.Equal(t, argentinaDatosApiRateTypeBuy, mapConfigRateTypeToArgentinaDatosApiRateType(settings.ArgentinaDatosRateTypeBuy))
	assert.Equal(t, argentinaDatosApiRateTypeSell, mapConfigRateTypeToArgentinaDatosApiRateType(settings.ArgentinaDatosRateTypeSell))
	assert.Equal(t, argentinaDatosApiRateTypeSell, mapConfigRateTypeToArgentinaDatosApiRateType("invalid"))
}

func TestArgentinaDatosDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(argentinaDatosMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "ARS", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestArgentinaDatosDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(argentinaDatosMinimumRequiredContent))
	assert.Equal(t, nil, err)

	expectedUpdateTime, _ := time.Parse("2006-01-02", "2026-07-07")
	assert.Equal(t, expectedUpdateTime.Unix(), actualLatestExchangeRateResponse.UpdateTime)
}

func TestArgentinaDatosDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(argentinaDatosMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.0006600660066006601",
	})
}

func TestArgentinaDatosDataSource_BuyRateType(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeBuy,
	}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(argentinaDatosMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.0006688963210702341",
	})
}

func TestArgentinaDatosDataSource_BlankContent(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestArgentinaDatosDataSource_EmptyArray(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[]"))
	assert.NotEqual(t, nil, err)
}

func TestIsValidArgentinaDatosExchangeHouse(t *testing.T) {
	assert.True(t, IsValidArgentinaDatosExchangeHouse("blue"))
	assert.True(t, IsValidArgentinaDatosExchangeHouse("BLUE"))
	assert.False(t, IsValidArgentinaDatosExchangeHouse("invalid"))
}

func TestNewArgentinaDatosDataSource_DefaultExchangeHouse(t *testing.T) {
	dataSource := newArgentinaDatosDataSource(&settings.Config{})

	assert.Equal(t, "blue", dataSource.exchangeHouse)
	assert.Equal(t, argentinaDatosApiRateTypeSell, dataSource.apiRateType)
}

func TestNewArgentinaDatosDataSource_BuyRateTypeMapping(t *testing.T) {
	dataSource := newArgentinaDatosDataSource(&settings.Config{
		ExchangeRatesArgentinaDatosRateType: settings.ArgentinaDatosRateTypeBuy,
	})

	assert.Equal(t, argentinaDatosApiRateTypeBuy, dataSource.apiRateType)
}

func TestArgentinaDatosDataSource_BuildRequests_BlueExchangeHouse(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}

	requests, err := dataSource.BuildRequests()
	assert.Equal(t, nil, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, http.MethodGet, requests[0].Method)
	assert.Equal(t, "https://api.argentinadatos.com/v1/cotizaciones/dolares/blue", requests[0].URL.String())
}

func TestArgentinaDatosDataSource_LatestEntrySelection(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	content := "[\n" +
		"  {\n" +
		"    \"casa\": \"blue\",\n" +
		"    \"compra\": 1000,\n" +
		"    \"venta\": 1100,\n" +
		"    \"fecha\": \"2026-01-01\"\n" +
		"  },\n" +
		"  {\n" +
		"    \"casa\": \"blue\",\n" +
		"    \"compra\": 1495,\n" +
		"    \"venta\": 1515,\n" +
		"    \"fecha\": \"2026-07-07\"\n" +
		"  }\n" +
		"]"

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)

	expectedUpdateTime, _ := time.Parse("2006-01-02", "2026-07-07")
	assert.Equal(t, expectedUpdateTime.Unix(), actualLatestExchangeRateResponse.UpdateTime)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.0006600660066006601",
	})
}

func TestArgentinaDatosDataSource_InvalidJson(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("not-json"))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestArgentinaDatosDataSource_ZeroRate(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	content := "[\n" +
		"  {\n" +
		"    \"casa\": \"blue\",\n" +
		"    \"compra\": 0,\n" +
		"    \"venta\": 0,\n" +
		"    \"fecha\": \"2026-07-07\"\n" +
		"  }\n" +
		"]"

	_, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}

func TestArgentinaDatosDataSource_InvalidDate(t *testing.T) {
	dataSource := &ArgentinaDatosDataSource{
		exchangeHouse: "blue",
		apiRateType:   argentinaDatosApiRateTypeSell,
	}
	context := core.NewNullContext()

	content := "[\n" +
		"  {\n" +
		"    \"casa\": \"blue\",\n" +
		"    \"compra\": 1495,\n" +
		"    \"venta\": 1515,\n" +
		"    \"fecha\": \"not-a-date\"\n" +
		"  }\n" +
		"]"

	_, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, errs.ErrFailedToRequestRemoteApi, err)
}
