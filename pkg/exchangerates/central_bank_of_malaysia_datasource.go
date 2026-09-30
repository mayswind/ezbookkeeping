package exchangerates

import (
	"encoding/json/v2"
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

const centralBankOfMalaysiaExchangeRateUrl = "https://api.bnm.gov.my/public/exchange-rate?quote=rm"
const centralBankOfMalaysiaExchangeRateReferenceUrl = "https://www.bnm.gov.my/rates-statistics"
const centralBankOfMalaysiaDataSource = "Bank Negara Malaysia"
const centralBankOfMalaysiaBaseCurrency = "MYR"

const centralBankOfMalaysiaDataUpdateDateFormat = "2006-01-02 15:04:05"
const centralBankOfMalaysiaDataUpdateDateTimezone = "Asia/Kuala_Lumpur"
const centralBankOfMalaysiaAcceptHeader = "application/vnd.BNM.API.v1+json"

// CentralBankOfMalaysiaDataSource defines the structure of exchange rates data source of central bank of Malaysia
type CentralBankOfMalaysiaDataSource struct {
	HttpExchangeRatesDataSource
}

// CentralBankOfMalaysiaExchangeRateData represents the whole data from central bank of Malaysia
type CentralBankOfMalaysiaExchangeRateData struct {
	ExchangeRates []*CentralBankOfMalaysiaExchangeRate `json:"data"`
	MetaInfo      *CentralBankOfMalaysiaMetaData       `json:"meta"`
}

// CentralBankOfMalaysiaMetaData represents the quotation metadata from central bank of Malaysia
type CentralBankOfMalaysiaMetaData struct {
	LastUpdated string `json:"last_updated"`
}

// CentralBankOfMalaysiaExchangeRate represents a currency quote from central bank of Malaysia
type CentralBankOfMalaysiaExchangeRate struct {
	CurrencyCode string                         `json:"currency_code"`
	Unit         float64                        `json:"unit"`
	Rate         *CentralBankOfMalaysiaRateData `json:"rate"`
}

// CentralBankOfMalaysiaRateData represents the rate and its date from central bank of Malaysia
type CentralBankOfMalaysiaRateData struct {
	Date        string  `json:"date"`
	BuyingRate  float64 `json:"buying_rate"`
	SellingRate float64 `json:"selling_rate"`
	MiddleRate  float64 `json:"middle_rate"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from central bank of Malaysia
func (e *CentralBankOfMalaysiaExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.ExchangeRates) < 1 {
		log.Errorf(c, "[central_bank_of_malaysia_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.ExchangeRates))

	for i := 0; i < len(e.ExchangeRates); i++ {
		exchangeRate := e.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.CurrencyCode]; !exists {
			continue
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	timezone, err := time.LoadLocation(centralBankOfMalaysiaDataUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_malaysia_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", centralBankOfMalaysiaDataUpdateDateTimezone)
		return nil
	}

	updateTime, err := time.ParseInLocation(centralBankOfMalaysiaDataUpdateDateFormat, e.MetaInfo.LastUpdated, timezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_malaysia_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", e.MetaInfo.LastUpdated)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    centralBankOfMalaysiaDataSource,
		ReferenceUrl:  centralBankOfMalaysiaExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  centralBankOfMalaysiaBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from central bank of Malaysia
func (e *CentralBankOfMalaysiaExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	if e.Unit <= 0 {
		log.Warnf(c, "[central_bank_of_malaysia_datasource.ToLatestExchangeRate] unit is invalid, currency is %s, unit is %f", e.CurrencyCode, e.Unit)
		return nil
	}

	if e.Rate == nil {
		log.Warnf(c, "[central_bank_of_malaysia_datasource.ToLatestExchangeRate] rate data is invalid, currency is %s", e.CurrencyCode)
		return nil
	}

	if e.Rate.MiddleRate <= 0 {
		log.Warnf(c, "[central_bank_of_malaysia_datasource.ToLatestExchangeRate] middle rate is invalid, currency is %s, middle rate is %f", e.CurrencyCode, e.Rate.MiddleRate)
		return nil
	}

	finalRate := e.Unit / e.Rate.MiddleRate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.CurrencyCode,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the central bank of Malaysia exchange rates http requests
func (e *CentralBankOfMalaysiaDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", centralBankOfMalaysiaExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", centralBankOfMalaysiaAcceptHeader)
	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the central bank of Malaysia data source raw response
func (e *CentralBankOfMalaysiaDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	centralBankOfMalaysiaData := &CentralBankOfMalaysiaExchangeRateData{}
	err := json.Unmarshal(content, centralBankOfMalaysiaData)

	if err != nil {
		log.Errorf(c, "[central_bank_of_malaysia_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if centralBankOfMalaysiaData == nil || centralBankOfMalaysiaData.MetaInfo == nil {
		log.Errorf(c, "[central_bank_of_malaysia_datasource.Parse] data is invalid, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := centralBankOfMalaysiaData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[central_bank_of_malaysia_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
