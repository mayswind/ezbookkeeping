package exchangerates

import (
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/models"
)

// ExchangeRatesDataSource defines the structure of exchange rates data source
type ExchangeRatesDataSource interface {
	// GetRequestUrl returns the data source url
	GetRequestUrl() string

	// Parse returns the common response entity according to the data source raw response
	Parse(c *core.Context, content []byte) (*models.LatestExchangeRateResponse, error)
}
