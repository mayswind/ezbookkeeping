package api

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
)

const EuroCentralBankExchangeRateUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ExchangeRatesApi struct {}

var (
	ExchangeRates = &ExchangeRatesApi{}
)

func (a *ExchangeRatesApi) LatestExchangeRateHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	resp, err := http.Get(EuroCentralBankExchangeRateUrl)

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
	euroCentralBankData := &models.EuroCentralBankExchangeRateData{}
	err = xml.Unmarshal(body, euroCentralBankData)

	if err != nil {
		log.ErrorfWithRequestId(c, "[exchange_rates.LatestExchangeRateHandler] failed to parse xml data for user \"uid:%d\", response is %s, because %s", uid, string(body), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := euroCentralBankData.ToLatestExchangeRateResponse()

	if latestExchangeRateResponse == nil {
		log.ErrorfWithRequestId(c, "[exchange_rates.LatestExchangeRateHandler] failed to parse latest exchange rate data for user \"uid:%d\", response is %s,", uid, string(body))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
