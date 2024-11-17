package exchangerates

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const centralBankOfUzbekistanExchangeRateUrl = "https://cbu.uz/ru/arkhiv-kursov-valyut/json/"
const centralBankOfUzbekistanExchangeRateReferenceUrl = "https://cbu.uz/en/arkhiv-kursov-valyut/"
const centralBankOfUzbekistanDataSource = "O‘zbekiston Respublikasi Markaziy banki"
const centralBankOfUzbekistanBaseCurrency = "UZS"

const centralBankOfUzbekistanUpdateDateFormat = "02.01.2006"
const centralBankOfUzbekistanUpdateDateTimezone = "Asia/Samarkand"

// CentralBankOfUzbekistanDataSource defines the structure of exchange rates data source of the central bank of the Republic of Uzbekistan
type CentralBankOfUzbekistanDataSource struct {
	ExchangeRatesDataSource
}

// CentralBankOfUzbekistanExchangeRates represents the exchange rates data from the central bank of the Republic of Uzbekistan
type CentralBankOfUzbekistanExchangeRates []*CentralBankOfUzbekistanExchangeRate

// CentralBankOfUzbekistanExchangeRate represents the exchange rate data from the central bank of the Republic of Uzbekistan
type CentralBankOfUzbekistanExchangeRate struct {
	Currency string `json:"Ccy"`
	Unit     string `json:"Nominal"`
	Rate     string `json:"Rate"`
	Date     string `json:"Date"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from the central bank of the Republic of Uzbekistan
func (e CentralBankOfUzbekistanExchangeRates) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e) < 1 {
		log.Errorf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	timezone, err := time.LoadLocation(centralBankOfUzbekistanUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", danmarksNationalbankDataUpdateDateTimezone)
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e))
	latestUpdateTime := int64(0)

	for i := 0; i < len(e); i++ {
		exchangeRate := e[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		updateTime, err := time.ParseInLocation(centralBankOfUzbekistanUpdateDateFormat, exchangeRate.Date, timezone)

		if err != nil {
			log.Errorf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Date)
			return nil
		}

		if updateTime.Unix() > latestUpdateTime {
			latestUpdateTime = updateTime.Unix()
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    centralBankOfUzbekistanDataSource,
		ReferenceUrl:  centralBankOfUzbekistanExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  centralBankOfUzbekistanBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from the central bank of the Republic of Uzbekistan
func (e *CentralBankOfUzbekistanExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(e.Rate)

	if err != nil {
		log.Warnf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	unit, err := utils.StringToFloat64(e.Unit)

	if err != nil {
		log.Warnf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRate] failed to parse unit, currency is %s, unit is %s", e.Currency, e.Unit)
		return nil
	}

	if unit <= 0 {
		log.Warnf(c, "[central_bank_of_uzbekistan_datasource.ToLatestExchangeRate] unit is less or equal zero, currency is %s, unit is %s", e.Currency, e.Unit)
		return nil
	}

	finalRate := 1000 * unit / rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the the central bank of the Republic of Uzbekistan exchange rates http requests
func (e *CentralBankOfUzbekistanDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", centralBankOfUzbekistanExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the the central bank of the Republic of Uzbekistan data source raw response
func (e *CentralBankOfUzbekistanDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	centralBankOfUzbekistanData := &CentralBankOfUzbekistanExchangeRates{}
	err := json.Unmarshal(content, centralBankOfUzbekistanData)

	if err != nil {
		log.Errorf(c, "[central_bank_of_uzbekistan_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := centralBankOfUzbekistanData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[central_bank_of_uzbekistan_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
