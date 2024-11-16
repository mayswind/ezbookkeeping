package exchangerates

import (
	"encoding/json"
	"math"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const nationalBankOfGeorgiaExchangeRateUrl = "https://nbg.gov.ge/gw/api/ct/monetarypolicy/currencies/en/json"
const nationalBankOfGeorgiaExchangeRateReferenceUrl = "https://nbg.gov.ge/en/monetary-policy/currency"
const nationalBankOfGeorgiaDataSource = "National Bank of Georgia"
const nationalBankOfGeorgiaBaseCurrency = "GEL"

const nationalBankOfGeorgiaUpdateDateFormat = "2006-01-02T15:04:05.999Z"

// NationalBankOfGeorgiaDataSource defines the structure of exchange rates data source of national bank of Georgia
type NationalBankOfGeorgiaDataSource struct {
	ExchangeRatesDataSource
}

// NationalBankOfGeorgiaExchangeRates represents the exchange rates data from national bank of Georgia
type NationalBankOfGeorgiaExchangeRates struct {
	Date          string                               `json:"date"`
	ExchangeRates []*NationalBankOfGeorgiaExchangeRate `json:"currencies"`
}

// NationalBankOfGeorgiaExchangeRate represents the exchange rate data from national bank of Georgia
type NationalBankOfGeorgiaExchangeRate struct {
	Currency string  `json:"code"`
	Quantity float64 `json:"quantity"`
	Rate     float64 `json:"rate"`
	Date     string  `json:"date"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from national bank of Georgia
func (e *NationalBankOfGeorgiaExchangeRates) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.ExchangeRates) < 1 {
		log.Errorf(c, "[national_bank_of_georgia_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.ExchangeRates))
	latestUpdateTime := int64(0)

	for i := 0; i < len(e.ExchangeRates); i++ {
		exchangeRate := e.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		updateTime, err := time.Parse(nationalBankOfGeorgiaUpdateDateFormat, exchangeRate.Date)

		if err != nil {
			log.Errorf(c, "[national_bank_of_georgia_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Date)
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
		DataSource:    nationalBankOfGeorgiaDataSource,
		ReferenceUrl:  nationalBankOfGeorgiaExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  nationalBankOfGeorgiaBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from national bank of Georgia
func (e *NationalBankOfGeorgiaExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	if e.Rate <= 0 {
		log.Warnf(c, "[national_bank_of_georgia_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %f", e.Currency, e.Rate)
		return nil
	}

	if e.Quantity <= 0 {
		log.Warnf(c, "[national_bank_of_georgia_datasource.ToLatestExchangeRate] quantity is invalid, currency is %s, quantity is %f", e.Currency, e.Quantity)
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

// GetRequestUrls returns the national bank of Georgia data source urls
func (e *NationalBankOfGeorgiaDataSource) GetRequestUrls() []string {
	return []string{nationalBankOfGeorgiaExchangeRateUrl}
}

// Parse returns the common response entity according to the national bank of Georgia data source raw response
func (e *NationalBankOfGeorgiaDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	nationalBankOfGeorgiaData := &[]*NationalBankOfGeorgiaExchangeRates{}
	err := json.Unmarshal(content, nationalBankOfGeorgiaData)

	if err != nil {
		log.Errorf(c, "[national_bank_of_georgia_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if nationalBankOfGeorgiaData == nil || len(*nationalBankOfGeorgiaData) < 1 {
		log.Errorf(c, "[national_bank_of_georgia_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := (*nationalBankOfGeorgiaData)[0].ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[national_bank_of_georgia_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
