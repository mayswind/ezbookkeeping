package exchangerates

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const euroCentralBankExchangeRateUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
const euroCentralBankExchangeRateReferenceUrl = "https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html"
const euroCentralBankDataSource = "European Central Bank"
const euroCentralBankBaseCurrency = "EUR"

const euroCentralBankDataUpdateDateFormat = "2006-01-02 15"
const euroCentralBankDataUpdateDateTimezone = "Europe/Berlin"

// EuroCentralBankDataSource defines the structure of exchange rates data source of euro central bank
type EuroCentralBankDataSource struct {
	HttpExchangeRatesDataSource
}

// EuroCentralBankExchangeRateData represents the whole data from euro central bank
type EuroCentralBankExchangeRateData struct {
	XMLName          xml.Name                        `xml:"Envelope"`
	AllExchangeRates []*EuroCentralBankExchangeRates `xml:"Cube>Cube"`
}

// EuroCentralBankExchangeRates represents the exchange rates data from euro central bank
type EuroCentralBankExchangeRates struct {
	Date          string                         `xml:"time,attr"`
	ExchangeRates []*EuroCentralBankExchangeRate `xml:"Cube"`
}

// EuroCentralBankExchangeRate represents the exchange rate data from euro central bank
type EuroCentralBankExchangeRate struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from euro central bank
func (e *EuroCentralBankExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.AllExchangeRates) < 1 {
		log.Errorf(c, "[euro_central_bank_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil
	}

	latestEuroCentralBankExchangeRate := e.AllExchangeRates[0]

	if len(latestEuroCentralBankExchangeRate.ExchangeRates) < 1 {
		log.Errorf(c, "[euro_central_bank_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(latestEuroCentralBankExchangeRate.ExchangeRates))

	for i := 0; i < len(latestEuroCentralBankExchangeRate.ExchangeRates); i++ {
		exchangeRate := latestEuroCentralBankExchangeRate.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		if _, err := utils.StringToFloat64(exchangeRate.Rate); err != nil {
			continue
		}

		exchangeRates = append(exchangeRates, exchangeRate.ToLatestExchangeRate())
	}

	timezone, err := time.LoadLocation(euroCentralBankDataUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[euro_central_bank_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", euroCentralBankDataUpdateDateTimezone)
		return nil
	}

	updateDateTime := latestEuroCentralBankExchangeRate.Date + " 16" // The reference rates are usually updated around 16:00 CET on every working day
	updateTime, err := time.ParseInLocation(euroCentralBankDataUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.Errorf(c, "[euro_central_bank_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    euroCentralBankDataSource,
		ReferenceUrl:  euroCentralBankExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  euroCentralBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from euro central bank
func (e *EuroCentralBankExchangeRate) ToLatestExchangeRate() *models.LatestExchangeRate {
	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     e.Rate,
	}
}

// BuildRequests returns the euro central bank exchange rates http requests
func (e *EuroCentralBankDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", euroCentralBankExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the euro central bank data source raw response
func (e *EuroCentralBankDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	euroCentralBankData := &EuroCentralBankExchangeRateData{}
	err := xmlDecoder.Decode(euroCentralBankData)

	if err != nil {
		log.Errorf(c, "[euro_central_bank_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := euroCentralBankData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[euro_central_bank_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
