package exchangerates

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const centralBankOfMyanmarExchangeRateUrl = "https://forex.cbm.gov.mm/api/latest"
const centralBankOfMyanmarExchangeRateReferenceUrl = "https://forex.cbm.gov.mm/index.php/fxrate"
const centralBankOfMyanmarDataSource = "မြန်မာနိုင်ငံတော်ဗဟိုဘဏ်"
const centralBankOfMyanmarBaseCurrency = "MMK"

var centralBankOfMyanmarSpecialCurrencyUnits = map[string]int32{
	"JPY": 100,
	"KHR": 100,
	"IDR": 100,
	"KRW": 100,
	"LAK": 100,
	"VND": 100,
}

// CentralBankOfMyanmarDataSource defines the structure of exchange rates data source of central bank of Myanmar
type CentralBankOfMyanmarDataSource struct {
	ExchangeRatesDataSource
}

// CentralBankOfMyanmarExchangeRate represents the exchange rate data from central bank of Myanmar
type CentralBankOfMyanmarExchangeRate struct {
	Timestamp     string            `json:"timestamp"`
	ExchangeRates map[string]string `json:"rates"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from central bank of Myanmar
func (e *CentralBankOfMyanmarExchangeRate) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.ExchangeRates))

	for currencyCode, exchangeRate := range e.ExchangeRates {
		if _, exists := validators.AllCurrencyNames[currencyCode]; !exists {
			continue
		}

		finalExchangeRate := e.BuildLatestExchangeRate(c, currencyCode, exchangeRate)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	updateTime, err := utils.StringToInt64(e.Timestamp)

	if err != nil {
		log.Errorf(c, "[central_bank_of_myanmar_datasource.ToLatestExchangeRateResponse] failed to parse timestamp, timestamp is %s", e.Timestamp)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    centralBankOfMyanmarDataSource,
		ReferenceUrl:  centralBankOfMyanmarExchangeRateReferenceUrl,
		UpdateTime:    updateTime,
		BaseCurrency:  centralBankOfMyanmarBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// BuildLatestExchangeRate returns a data pair according to original data from central bank of Myanmar
func (e *CentralBankOfMyanmarExchangeRate) BuildLatestExchangeRate(c core.Context, currencyCode string, exchangeRate string) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(strings.ReplaceAll(exchangeRate, ",", ""))

	if err != nil {
		log.Warnf(c, "[central_bank_of_myanmar_datasource.BuildLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", currencyCode, exchangeRate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[central_bank_of_myanmar_datasource.BuildLatestExchangeRate] rate is invalid, currency is %s, rate is %s", currencyCode, exchangeRate)
		return nil
	}

	unit, has := centralBankOfMyanmarSpecialCurrencyUnits[currencyCode]

	if !has {
		unit = 1
	}

	finalRate := float64(unit) / rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: currencyCode,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the central bank of Myanmar exchange rates http requests
func (e *CentralBankOfMyanmarDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", centralBankOfMyanmarExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the central bank of Myanmar data source raw response
func (e *CentralBankOfMyanmarDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	centralBankOfMyanmarData := &CentralBankOfMyanmarExchangeRate{}
	err := json.Unmarshal(content, centralBankOfMyanmarData)

	if err != nil {
		log.Errorf(c, "[central_bank_of_myanmar_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := centralBankOfMyanmarData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[central_bank_of_myanmar_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
