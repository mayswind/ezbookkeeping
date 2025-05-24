package exchangerates

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// HttpExchangeRatesDataSource defines the structure of http exchange rates data source
type HttpExchangeRatesDataSource interface {
	// BuildRequests returns the http requests
	BuildRequests() ([]*http.Request, error)

	// Parse returns the common response entity according to the data source raw response
	Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error)
}

// CommonHttpExchangeRatesDataSource defines the structure of common http exchange rates data source
type CommonHttpExchangeRatesDataSource struct {
	ExchangeRatesDataSource
	dataSource HttpExchangeRatesDataSource
}

func (e *CommonHttpExchangeRatesDataSource) GetLatestExchangeRates(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestExchangeRateResponse, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	utils.SetProxyUrl(transport, currentConfig.ExchangeRatesProxy)

	if currentConfig.ExchangeRatesSkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(currentConfig.ExchangeRatesRequestTimeout) * time.Millisecond,
	}

	requests, err := e.dataSource.BuildRequests()

	if err != nil {
		log.Errorf(c, "[http_exchange_rates_datasource.GetLatestExchangeRates] failed to build requests for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	exchangeRateResps := make([]*models.LatestExchangeRateResponse, 0, len(requests))

	for i := 0; i < len(requests); i++ {
		req := requests[i]

		if len(req.Header.Values("User-Agent")) < 1 {
			req.Header.Set("User-Agent", fmt.Sprintf("ezBookkeeping/%s", settings.Version))
		} else if req.Header.Get("User-Agent") == "" {
			req.Header.Del("User-Agent")
		}

		resp, err := client.Do(req)

		if err != nil {
			log.Errorf(c, "[http_exchange_rates_datasource.GetLatestExchangeRates] failed to request latest exchange rate data for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		if resp.StatusCode != 200 {
			log.Errorf(c, "[http_exchange_rates_datasource.GetLatestExchangeRates] failed to get latest exchange rate data response for user \"uid:%d\", because response code is not 200", uid)
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		log.Debugf(c, "[http_exchange_rates_datasource.GetLatestExchangeRates] response#%d is %s", i, body)

		exchangeRateResp, err := e.dataSource.Parse(c, body)

		if err != nil {
			log.Errorf(c, "[http_exchange_rates_datasource.GetLatestExchangeRates] failed to parse response for user \"uid:%d\", because %s", uid, err.Error())
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

func newCommonHttpExchangeRatesDataSource(dataSource HttpExchangeRatesDataSource) *CommonHttpExchangeRatesDataSource {
	return &CommonHttpExchangeRatesDataSource{
		dataSource: dataSource,
	}
}
