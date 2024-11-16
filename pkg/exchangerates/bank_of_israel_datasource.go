package exchangerates

import (
	"bytes"
	"encoding/xml"
	"math"
	"time"

	"golang.org/x/net/html/charset"

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

// BankOfIsraelExchangeRateData represents the whole data from bank of Israel
type BankOfIsraelExchangeRateData struct {
	XMLName          xml.Name                    `xml:"ExchangeRatesResponseCollectioDTO"`
	AllExchangeRates []*BankOfIsraelExchangeRate `xml:"ExchangeRates>ExchangeRateResponseDTO"`
}

// BankOfIsraelExchangeRate represents the exchange rate data from bank of Israel
type BankOfIsraelExchangeRate struct {
	Currency   string `xml:"Key"`
	Rate       string `xml:"CurrentExchangeRate"`
	LastUpdate string `xml:"LastUpdate"`
	Unit       string `xml:"Unit"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from bank of Israel
func (e *BankOfIsraelExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
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

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
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
func (e *BankOfIsraelExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(e.Rate)

	if err != nil {
		log.Warnf(c, "[bank_of_israel_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[bank_of_israel_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	unit, err := utils.StringToFloat64(e.Unit)

	if err != nil {
		log.Warnf(c, "[bank_of_israel_datasource.ToLatestExchangeRate] failed to parse unit, currency is %s, unit is %s", e.Currency, e.Unit)
		return nil
	}

	finalRate := unit / rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// GetRequestUrls returns the bank of Israel data source urls
func (e *BankOfIsraelDataSource) GetRequestUrls() []string {
	return []string{bankOfIsraelExchangeRateUrl}
}

// Parse returns the common response entity according to the bank of Israel data source raw response
func (e *BankOfIsraelDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	bankOfIsraelData := &BankOfIsraelExchangeRateData{}
	err := xmlDecoder.Decode(bankOfIsraelData)

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
