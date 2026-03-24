package exchangerates

import (
	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/models"
	"github.com/Paxtiny/oscar/pkg/settings"
)

// ExchangeRatesDataProvider defines the structure of exchange rates data provider
type ExchangeRatesDataProvider interface {
	// GetLatestExchangeRates returns the common response entities
	GetLatestExchangeRates(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestExchangeRateResponse, error)
}
