package exchangerates

import (
	"bytes"
	"encoding/xml"
	"math"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const bankOfRussiaExchangeRateUrl = "https://cbr.ru/scripts/XML_daily_eng.asp"
const bankOfRussiaExchangeRateReferenceUrl = "https://www.cbr.ru/eng/currency_base/daily/"
const bankOfRussiaDataSource = "Банк России"
const bankOfRussiaBaseCurrency = "RUB"

const bankOfRussiaUpdateDateFormat = "02.01.2006 15:04"
const bankOfRussiaUpdateDateTimezone = "Europe/Moscow"

// BankOfRussiaDataSource defines the structure of exchange rates data source of bank of Russia
type BankOfRussiaDataSource struct {
	HttpExchangeRatesDataSource
}

// BankOfRussiaExchangeRateData represents the whole data from bank of Russia
type BankOfRussiaExchangeRateData struct {
	XMLName       xml.Name                    `xml:"ValCurs"`
	Date          string                      `xml:"Date,attr"`
	ExchangeRates []*BankOfRussiaExchangeRate `xml:"Valute"`
}

// BankOfRussiaExchangeRate represents the exchange rate data from bank of Russia
type BankOfRussiaExchangeRate struct {
	Currency string `xml:"CharCode"`
	Rate     string `xml:"VunitRate"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from bank of Russia
func (e *BankOfRussiaExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.ExchangeRates) < 1 {
		log.Errorf(c, "[bank_of_russia_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.ExchangeRates))

	for i := 0; i < len(e.ExchangeRates); i++ {
		exchangeRate := e.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	timezone, err := time.LoadLocation(bankOfRussiaUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[bank_of_russia_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", bankOfRussiaUpdateDateTimezone)
		return nil
	}

	updateDateTime := e.Date + " 15:30" // the Bank of Russia switches to setting official exchange rates of foreign currencies against the ruble as of 15:30 Moscow time.
	updateTime, err := time.ParseInLocation(bankOfRussiaUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.Errorf(c, "[bank_of_russia_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    bankOfRussiaDataSource,
		ReferenceUrl:  bankOfRussiaExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  bankOfRussiaBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from bank of Russia
func (e *BankOfRussiaExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(strings.ReplaceAll(e.Rate, ",", "."))

	if err != nil {
		log.Warnf(c, "[bank_of_russia_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[bank_of_russia_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	finalRate := 1 / rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the bank of Russia exchange rates http requests
func (e *BankOfRussiaDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", bankOfRussiaExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the bank of Russia data source raw response
func (e *BankOfRussiaDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	bankOfRussiaData := &BankOfRussiaExchangeRateData{}
	err := xmlDecoder.Decode(bankOfRussiaData)

	if err != nil {
		log.Errorf(c, "[bank_of_russia_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := bankOfRussiaData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[bank_of_russia_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
