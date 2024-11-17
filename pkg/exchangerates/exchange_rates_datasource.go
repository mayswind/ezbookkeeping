package exchangerates

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ExchangeRatesDataSource defines the structure of exchange rates data source
type ExchangeRatesDataSource interface {
	// BuildRequests returns the http requests
	BuildRequests() ([]*http.Request, error)

	// Parse returns the common response entity according to the data source raw response
	Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error)
}
