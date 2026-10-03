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
// https://tassidicambio.bancaditalia.it/terzevalute-wf-ui-web/
// https://tassidicambio.bancaditalia.it/terzevalute-wf-ui-web/assets/files/Operating_Instructions.pdf

const bankOfItalyExchangeRateUrl = "https://tassidicambio.bancaditalia.it/terzevalute-wf-web/rest/v1.0/latestRates?lang=en"
const bankOfItalyExchangeRateReferenceUrl = "https://tassidicambio.bancaditalia.it/terzevalute-wf-ui-web/latestRates"
const bankOfItalyDataSource = "Banca d'Italia"
const bankOfItalyBaseCurrency = "EUR"

const bankOfItalyDataUpdateDateFormat = "2006-01-02"
const bankOfItalyDataUpdateDateTimezone = "Europe/Rome"
const bankOfItalyAcceptHeader = "application/json"

// BankOfItalyDataSource defines the structure of exchange rates data source of bank of Italy
type BankOfItalyDataSource struct {
	HttpExchangeRatesDataSource
}

// BankOfItalyExchangeRateData represents the whole data from bank of Italy
type BankOfItalyExchangeRateData struct {
	ExchangeRates []*BankOfItalyExchangeRate `json:"latestRates"`
}

// BankOfItalyExchangeRate represents a currency quote from bank of Italy
type BankOfItalyExchangeRate struct {
	CurrencyCode string `json:"isoCode"`
	Rate         string `json:"eurRate"`
	Date         string `json:"referenceDate"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from bank of Italy
func (e *BankOfItalyExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.ExchangeRates) < 1 {
		log.Errorf(c, "[bank_of_italy_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	timezone, err := time.LoadLocation(bankOfItalyDataUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[bank_of_italy_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", bankOfItalyDataUpdateDateTimezone)
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.ExchangeRates))
	exchangeRateCurrencyIndexMap := make(map[string]int)
	exchangeRateLatestUpdateUnixTimeMap := make(map[string]int64)
	latestUpdateTime := int64(0)

	for i := 0; i < len(e.ExchangeRates); i++ {
		exchangeRate := e.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.CurrencyCode]; !exists {
			continue
		}

		updateTime, err := time.ParseInLocation(bankOfItalyDataUpdateDateFormat, exchangeRate.Date, timezone)

		if err != nil {
			log.Errorf(c, "[bank_of_italy_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Date)
			return nil
		}

		if updateTime.Unix() > latestUpdateTime {
			latestUpdateTime = updateTime.Unix()
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		if index, exists := exchangeRateCurrencyIndexMap[exchangeRate.CurrencyCode]; exists {
			if updateTime.Unix() > exchangeRateLatestUpdateUnixTimeMap[exchangeRate.CurrencyCode] {
				exchangeRates[index] = finalExchangeRate
				exchangeRateLatestUpdateUnixTimeMap[exchangeRate.CurrencyCode] = updateTime.Unix()
			}
		} else {
			exchangeRateCurrencyIndexMap[exchangeRate.CurrencyCode] = len(exchangeRates)
			exchangeRateLatestUpdateUnixTimeMap[exchangeRate.CurrencyCode] = updateTime.Unix()
			exchangeRates = append(exchangeRates, finalExchangeRate)
		}
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    bankOfItalyDataSource,
		ReferenceUrl:  bankOfItalyExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  bankOfItalyBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from bank of Italy
func (e *BankOfItalyExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(e.Rate)

	if err != nil {
		log.Warnf(c, "[bank_of_italy_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.CurrencyCode, e.Rate)
		return nil
	}

	if rate <= 0 || math.IsInf(rate, 0) {
		log.Warnf(c, "[bank_of_italy_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.CurrencyCode, e.Rate)
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.CurrencyCode,
		Rate:     utils.Float64ToString(rate),
	}
}

// BuildRequests returns the bank of Italy exchange rates http requests
func (e *BankOfItalyDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", bankOfItalyExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", bankOfItalyAcceptHeader)
	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the bank of Italy data source raw response
func (e *BankOfItalyDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	bankOfItalyData := &BankOfItalyExchangeRateData{}
	err := json.Unmarshal(content, bankOfItalyData)

	if err != nil {
		log.Errorf(c, "[bank_of_italy_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := bankOfItalyData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[bank_of_italy_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
