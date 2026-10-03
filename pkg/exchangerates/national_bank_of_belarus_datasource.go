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

// API Documentation:
// https://www.nbrb.by/apihelp/exrates

const nationalBankOfBelarusExchangeRateUrl = "https://api.nbrb.by/exrates/rates?periodicity=0"
const nationalBankOfBelarusExchangeRateReferenceUrl = "https://www.nbrb.by/engl/statistics/rates/ratesdaily"
const nationalBankOfBelarusDataSource = "Нацыянальны банк Рэспублiкi Беларусь"
const nationalBankOfBelarusBaseCurrency = "BYN"

const nationalBankOfBelarusDataUpdateDateFormat = "2006-01-02T15:04:05"
const nationalBankOfBelarusDataUpdateDateTimezone = "Europe/Minsk"

// NationalBankOfBelarusDataSource defines the structure of exchange rates data source of national bank of Belarus
type NationalBankOfBelarusDataSource struct {
	HttpExchangeRatesDataSource
}

// NationalBankOfBelarusExchangeRateData represents the whole data from national bank of Belarus
type NationalBankOfBelarusExchangeRateData []*NationalBankOfBelarusExchangeRate

// NationalBankOfBelarusExchangeRate represents a currency quote from national bank of Belarus
type NationalBankOfBelarusExchangeRate struct {
	CurrencyCode string  `json:"Cur_Abbreviation"`
	Unit         float64 `json:"Cur_Scale"`
	Rate         float64 `json:"Cur_OfficialRate"`
	Date         string  `json:"Date"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from national bank of Belarus
func (e *NationalBankOfBelarusExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(*e) < 1 {
		log.Errorf(c, "[national_bank_of_belarus_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	timezone, err := time.LoadLocation(nationalBankOfBelarusDataUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[national_bank_of_belarus_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", nationalBankOfBelarusDataUpdateDateTimezone)
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(*e))
	latestUpdateTime := int64(0)

	for _, exchangeRate := range *e {
		if _, exists := validators.AllCurrencyNames[exchangeRate.CurrencyCode]; !exists {
			continue
		}

		updateTime, err := time.ParseInLocation(nationalBankOfBelarusDataUpdateDateFormat, exchangeRate.Date, timezone)

		if err != nil {
			log.Errorf(c, "[national_bank_of_belarus_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Date)
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
		DataSource:    nationalBankOfBelarusDataSource,
		ReferenceUrl:  nationalBankOfBelarusExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  nationalBankOfBelarusBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from national bank of Belarus
func (e *NationalBankOfBelarusExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	if e.Unit <= 0 {
		log.Warnf(c, "[national_bank_of_belarus_datasource.ToLatestExchangeRate] unit is invalid, currency is %s, unit is %f", e.CurrencyCode, e.Unit)
		return nil
	}

	if e.Rate <= 0 {
		log.Warnf(c, "[national_bank_of_belarus_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %f", e.CurrencyCode, e.Rate)
		return nil
	}

	finalRate := e.Unit / e.Rate

	if math.IsInf(finalRate, 0) {
		log.Warnf(c, "[national_bank_of_belarus_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %f", e.CurrencyCode, finalRate)
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.CurrencyCode,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the national bank of Belarus exchange rates http requests
func (e *NationalBankOfBelarusDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", nationalBankOfBelarusExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the national bank of Belarus data source raw response
func (e *NationalBankOfBelarusDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	nationalBankOfBelarusData := &NationalBankOfBelarusExchangeRateData{}
	err := json.Unmarshal(content, nationalBankOfBelarusData)

	if err != nil {
		log.Errorf(c, "[national_bank_of_belarus_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := nationalBankOfBelarusData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[national_bank_of_belarus_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
