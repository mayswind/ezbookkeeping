package exchangerates

import (
	"io"
	"net/http"
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/httpclient"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// HttpExchangeRatesDataSource defines the structure of http exchange rates data source
type HttpExchangeRatesDataSource interface {
	// BuildRequests returns the http requests
	BuildRequests() ([]*http.Request, error)

	// Parse returns the common response entity according to the data source raw response
	Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error)
}

// CommonHttpExchangeRatesDataProvider defines the structure of common http exchange rates data provider
type CommonHttpExchangeRatesDataProvider struct {
	ExchangeRatesDataProvider
	dataSource HttpExchangeRatesDataSource
	httpClient *http.Client
}

func (e *CommonHttpExchangeRatesDataProvider) GetLatestExchangeRates(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestExchangeRateResponse, error) {
	requests, err := e.dataSource.BuildRequests()

	if err != nil {
		log.Errorf(c, "[common_http_exchange_rates_data_provider.GetLatestExchangeRates] failed to build requests for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	exchangeRateResps := make([]*models.LatestExchangeRateResponse, 0, len(requests))

	for i := 0; i < len(requests); i++ {
		req := requests[i]
		req = req.WithContext(httpclient.CustomHttpResponseLog(c, func(data []byte) {
			log.Debugf(c, "[common_http_exchange_rates_data_provider.GetLatestExchangeRates] response#%d is %s", i, data)
		}))

		resp, err := e.httpClient.Do(req)

		if err != nil {
			log.Errorf(c, "[common_http_exchange_rates_data_provider.GetLatestExchangeRates] failed to request latest exchange rate data for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			log.Errorf(c, "[common_http_exchange_rates_data_provider.GetLatestExchangeRates] failed to get latest exchange rate data response for user \"uid:%d\", because response code is %d", uid, resp.StatusCode)
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		exchangeRateResp, err := e.dataSource.Parse(c, body)

		if err != nil {
			log.Errorf(c, "[common_http_exchange_rates_data_provider.GetLatestExchangeRates] failed to parse response for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrFailedToRequestRemoteApi)
		}

		exchangeRateResps = append(exchangeRateResps, exchangeRateResp)
	}

	lastExchangeRateResponse := exchangeRateResps[len(exchangeRateResps)-1]
	allExchangeRatesMap := make(map[string]string)

	for i := 0; i < len(exchangeRateResps); i++ {
		exchangeRateResp := exchangeRateResps[i]

		for j := 0; j < len(exchangeRateResp.ExchangeRates); j++ {
			exchangeRate := exchangeRateResp.ExchangeRates[j]
			allExchangeRatesMap[exchangeRate.Currency] = exchangeRate.Rate
		}
	}

	allExchangeRatesMap[lastExchangeRateResponse.BaseCurrency] = "1"
	allExchangeRates := make(models.LatestExchangeRateSlice, 0, len(allExchangeRatesMap))

	for currency, rate := range allExchangeRatesMap {
		allExchangeRates = append(allExchangeRates, &models.LatestExchangeRate{
			Currency: currency,
			Rate:     rate,
		})
	}

	sort.Sort(allExchangeRates)

	finalExchangeRateResponse := &models.LatestExchangeRateResponse{
		DataSource:    lastExchangeRateResponse.DataSource,
		ReferenceUrl:  lastExchangeRateResponse.ReferenceUrl,
		UpdateTime:    lastExchangeRateResponse.UpdateTime,
		BaseCurrency:  lastExchangeRateResponse.BaseCurrency,
		ExchangeRates: allExchangeRates,
	}

	return finalExchangeRateResponse, nil
}

func newCommonHttpExchangeRatesDataProvider(config *settings.Config, dataSource HttpExchangeRatesDataSource) *CommonHttpExchangeRatesDataProvider {
	return &CommonHttpExchangeRatesDataProvider{
		dataSource: dataSource,
		httpClient: httpclient.NewHttpClient(config.ExchangeRatesRequestTimeout, config.ExchangeRatesProxy, config.ExchangeRatesSkipTLSVerify, core.GetOutgoingUserAgent(), config.EnableDebugLog),
	}
}
