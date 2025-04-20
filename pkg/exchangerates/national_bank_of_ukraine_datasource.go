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

const nationalBankOfUkraineExchangeRateUrl = "https://bank.gov.ua/NBU_Exchange/exchange?json"
const nationalBankOfUkraineExchangeRateReferenceUrl = "https://bank.gov.ua/en/markets/exchangerates"
const nationalBankOfUkraineDataSource = "Національний банк України"
const nationalBankOfUkraineBaseCurrency = "UAH"

const nationalBankOfUkraineUpdateDateFormat = "02.01.2006"

// NationalBankOfUkraineDataSource defines the structure of exchange rates data source of National Bank of Ukraine
type NationalBankOfUkraineDataSource struct {
	ExchangeRatesDataSource
}

// NationalBankOfUkraineExchangeRates  represents the exchange rates data from National Bank of Ukraine
type NationalBankOfUkraineExchangeRates []NaionalBankOfUkraineExchangeRate

// NaionalBankOfUkraineExchangeRate represents the exchange rate data from National Bank of Ukraine
type NaionalBankOfUkraineExchangeRate struct {
	Currency string  `json:"CurrencyCodeL"`
	Quantity float64 `json:"Units"`
	Rate     float64 `json:"Amount"`
	Date     string  `json:"StartDate"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from National Bank of Ukraine
func (e *NationalBankOfUkraineExchangeRates) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(*e))
	latestUpdateTime := int64(0)

	for _, exchangeRate := range *e {
		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		updateTime, err := time.Parse(nationalBankOfUkraineUpdateDateFormat, exchangeRate.Date)

		if err != nil {
			log.Errorf(c, "[national_bank_of_ukraine_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Date)
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
		DataSource:    nationalBankOfUkraineDataSource,
		ReferenceUrl:  nationalBankOfUkraineExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  nationalBankOfUkraineBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from National Bank of Ukraine
func (e *NaionalBankOfUkraineExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	if e.Rate <= 0 {
		log.Warnf(c, "[national_bank_of_ukraine_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %f", e.Currency, e.Rate)
		return nil
	}

	if e.Quantity <= 0 {
		log.Warnf(c, "[national_bank_of_ukraine_datasource.ToLatestExchangeRate] quantity is invalid, currency is %s, quantity is %f", e.Currency, e.Quantity)
		return nil
	}

	finalRate := e.Quantity / e.Rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the National Bank of Ukraine exchange rates http requests
func (e *NationalBankOfUkraineDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", nationalBankOfUkraineExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the National Bank of Ukraine data source raw response
func (e *NationalBankOfUkraineDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	var nationalBankOfUkraineData NationalBankOfUkraineExchangeRates
	err := json.Unmarshal(content, &nationalBankOfUkraineData)

	if err != nil {
		log.Errorf(c, "[national_bank_of_ukraine_datasource.Parse] failed to parse JSON data, content: %s, error: %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if len(nationalBankOfUkraineData) == 0 {
		log.Errorf(c, "[national_bank_of_ukraine_datasource.Parse] exchange rate list is empty")
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := nationalBankOfUkraineData.ToLatestExchangeRateResponse(c)
	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[national_bank_of_ukraine_datasource.Parse] failed to parse latest exchange rate data, content: %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
