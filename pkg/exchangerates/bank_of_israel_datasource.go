package exchangerates

import (
	"encoding/xml"
	"math"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const bankOfIsraelExchangeRateUrl = "https://boi.org.il/PublicApi/GetExchangeRates?asXml=true"
const bankOfIsraelExchangeRateReferenceUrl = "https://www.boi.org.il/en/economic-roles/financial-markets/exchange-rates/"
const bankOfIsraelDataSource = "Bank of Israel"
const bankOfIsraelBaseCurrency = "ILS"

const bankOfIsraelDataUpdateDateFormat = "2006-01-02T15:04:05.9999999Z"

// BankOfIsraelDataSource defines the structure of exchange rates data source of bank of Israel
type BankOfIsraelDataSource struct {
	ExchangeRatesDataSource
}

// bankOfIsraelExchangeRateData represents the whole data from bank of Israel
type bankOfIsraelExchangeRateData struct {
	XMLName          xml.Name                    `xml:"ExchangeRatesResponseCollectioDTO"`
	AllExchangeRates []*bankOfIsraelExchangeRate `xml:"ExchangeRates>ExchangeRateResponseDTO"`
}

// bankOfIsraelExchangeRate represents the exchange rate data from bank of Israel
type bankOfIsraelExchangeRate struct {
	Currency   string `xml:"Key"`
	Rate       string `xml:"CurrentExchangeRate"`
	LastUpdate string `xml:"LastUpdate"`
	Unit       string `xml:"Unit"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from bank of Israel
func (e *bankOfIsraelExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.AllExchangeRates) < 1 {
		log.Errorf(c, "[bank_of_israel_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil
	}

	latestUpdateDate := ""
	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.AllExchangeRates))

	for i := 0; i < len(e.AllExchangeRates); i++ {
		exchangeRate := e.AllExchangeRates[i]

		if latestUpdateDate == "" {
			latestUpdateDate = exchangeRate.LastUpdate
		}

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		rate, err := utils.StringToFloat64(exchangeRate.Rate)

		if err != nil {
			log.Warnf(c, "[bank_of_israel_datasource.ToLatestExchangeRateResponse] failed to parse rate, rate is %s", exchangeRate.Rate)
			continue
		}

		if rate <= 0 {
			log.Warnf(c, "[bank_of_israel_datasource.ToLatestExchangeRateResponse] rate is invalid, rate is %s", exchangeRate.Rate)
			continue
		}

		unit, err := utils.StringToFloat64(exchangeRate.Unit)

		if err != nil {
			log.Warnf(c, "[bank_of_israel_datasource.ToLatestExchangeRateResponse] failed to parse unit, unit is %s", exchangeRate.Unit)
			continue
		}

		finalRate := unit / rate

		if math.IsInf(finalRate, 0) {
			continue
		}

		exchangeRates = append(exchangeRates, &models.LatestExchangeRate{
			Currency: exchangeRate.Currency,
			Rate:     utils.Float64ToString(finalRate),
		})
	}

	updateTime, err := time.Parse(bankOfIsraelDataUpdateDateFormat, latestUpdateDate)

	if err != nil {
		log.Errorf(c, "[bank_of_israel_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", latestUpdateDate)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    bankOfIsraelDataSource,
		ReferenceUrl:  bankOfIsraelExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  bankOfIsraelBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from bank of Israel
func (e *bankOfIsraelExchangeRate) ToLatestExchangeRate() *models.LatestExchangeRate {
	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     e.Rate,
	}
}

// GetRequestUrls returns the bank of Israel data source urls
func (e *BankOfIsraelDataSource) GetRequestUrls() []string {
	return []string{bankOfIsraelExchangeRateUrl}
}

// Parse returns the common response entity according to the bank of Israel data source raw response
func (e *BankOfIsraelDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	bankOfIsraelData := &bankOfIsraelExchangeRateData{}
	err := xml.Unmarshal(content, bankOfIsraelData)

	if err != nil {
		log.Errorf(c, "[bank_of_israel_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := bankOfIsraelData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[bank_of_israel_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
