package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestInitializeExchangeRatesDataSource_CentralBankOfTheArgentineRepublicDataSource(t *testing.T) {
	config := &settings.Config{
		ExchangeRatesDataSource: settings.CentralBankOfTheArgentineRepublicDataSource,
	}

	err := InitializeExchangeRatesDataSource(config)

	assert.Nil(t, err)
	assert.NotNil(t, Container.current)
}
