package api

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/exchangerates"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/settings"
)

// ExchangeRatesApi represents exchange rate api
type ExchangeRatesApi struct{}

// Initialize a exchange rate api singleton instance
var (
	ExchangeRates = &ExchangeRatesApi{}
)

// LatestExchangeRateHandler returns latest exchange rate data
func (a *ExchangeRatesApi) LatestExchangeRateHandler(c *core.Context) (interface{}, *errs.Error) {
	dataSource := exchangerates.Container.Current

	if dataSource == nil {
		return nil, errs.ErrInvalidExchangeRatesDataSource
	}

	uid := c.GetCurrentUid()

	client := &http.Client{
		Timeout: time.Duration(settings.Container.Current.ExchangeRatesRequestTimeout) * time.Millisecond,
	}

	resp, err := client.Get(dataSource.GetRequestUrl())

	if err != nil {
		log.ErrorfWithRequestId(c, "[exchange_rates.LatestExchangeRateHandler] failed to request latest exchange rate data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if resp.StatusCode != 200 {
		log.ErrorfWithRequestId(c, "[exchange_rates.LatestExchangeRateHandler] failed to get latest exchange rate data response for user \"uid:%d\", because response code is not 200", uid)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	latestExchangeRateResponse, err := dataSource.Parse(c, body)

	if err != nil {
		log.ErrorfWithRequestId(c, "[exchange_rates.LatestExchangeRateHandler] failed to parse response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrFailedToRequestRemoteApi)
	}

	return latestExchangeRateResponse, nil
}
