package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestInitializeExchangeRatesDataSource_ArgentinaDatosDataSource(t *testing.T) {
	config := &settings.Config{
		ExchangeRatesDataSource:                 settings.ArgentinaDatosDataSource,
		ExchangeRatesArgentinaDatosExchangeHouse: "blue",
		ExchangeRatesArgentinaDatosRateType:     settings.ArgentinaDatosRateTypeSell,
	}

	err := InitializeExchangeRatesDataSource(config)
	assert.Nil(t, err)
	assert.NotNil(t, Container.current)
}

func TestInitializeExchangeRatesDataSource_ArgentinaDatosEmptyExchangeHouse(t *testing.T) {
	config := &settings.Config{
		ExchangeRatesDataSource:             settings.ArgentinaDatosDataSource,
		ExchangeRatesArgentinaDatosRateType: settings.ArgentinaDatosRateTypeSell,
	}

	err := InitializeExchangeRatesDataSource(config)
	assert.Nil(t, err)
	assert.NotNil(t, Container.current)
}

func TestInitializeExchangeRatesDataSource_ArgentinaDatosInvalidExchangeHouse(t *testing.T) {
	config := &settings.Config{
		ExchangeRatesDataSource:                 settings.ArgentinaDatosDataSource,
		ExchangeRatesArgentinaDatosExchangeHouse: "invalid-house",
	}

	err := InitializeExchangeRatesDataSource(config)
	assert.Equal(t, errs.ErrInvalidExchangeRatesDataSource, err)
}
