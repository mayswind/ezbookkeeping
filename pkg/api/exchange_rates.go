package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/exchangerates"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// ExchangeRatesApi represents exchange rate api
type ExchangeRatesApi struct {
	ApiUsingConfig
}

// Initialize a exchange rate api singleton instance
var (
	ExchangeRates = &ExchangeRatesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// LatestExchangeRateHandler returns latest exchange rate data
func (a *ExchangeRatesApi) LatestExchangeRateHandler(c *core.WebContext) (any, *errs.Error) {
	dataSource := exchangerates.Container.Current

	if dataSource == nil {
		return nil, errs.ErrInvalidExchangeRatesDataSource
	}

	exchangeRateResponse, err := dataSource.GetLatestExchangeRates(c, c.GetCurrentUid(), a.container.Current)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return exchangeRateResponse, nil
}
