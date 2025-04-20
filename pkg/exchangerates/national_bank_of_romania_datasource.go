package exchangerates

import (
	"bytes"
	"encoding/xml"
	"math"
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

const nationalBankOfRomaniaExchangeRateUrl = "https://www.bnr.ro/nbrfxrates.xml"
const nationalBankOfRomaniaExchangeRateReferenceUrl = "https://bnro.ro/Exchange-rates-1224.aspx"
const nationalBankOfRomaniaDataSource = "Banca Naţională a României"

const nationalBankOfRomaniaUpdateDateFormat = "2006-01-02 15"
const nationalBankOfRomaniaUpdateDateTimezone = "Europe/Bucharest"

// NationalBankOfRomaniaDataSource defines the structure of exchange rates data source of national bank of Romania
type NationalBankOfRomaniaDataSource struct {
	ExchangeRatesDataSource
}

// NationalBankOfRomaniaExchangeRateData represents the whole data from national bank of Romania
type NationalBankOfRomaniaExchangeRateData struct {
	XMLName xml.Name                                     `xml:"DataSet"`
	Header  *NationalBankOfRomaniaExchangeRateDataHeader `xml:"Header"`
	Body    *NationalBankOfRomaniaExchangeRateDataBody   `xml:"Body"`
}

// NationalBankOfRomaniaExchangeRateDataHeader represents the header for exchange rates data of national bank of Romania
type NationalBankOfRomaniaExchangeRateDataHeader struct {
	PublishingDate string `xml:"PublishingDate"`
}

// NationalBankOfRomaniaExchangeRateDataBody represents the body for exchange rates data of national bank of Romania
type NationalBankOfRomaniaExchangeRateDataBody struct {
	OrigCurrency     string                                `xml:"OrigCurrency"`
	AllExchangeRates []*NationalBankOfRomaniaExchangeRates `xml:"Cube"`
}

// NationalBankOfRomaniaExchangeRates represents the exchange rates data from national bank of Romania
type NationalBankOfRomaniaExchangeRates struct {
	Date          string                               `xml:"date,attr"`
	ExchangeRates []*NationalBankOfRomaniaExchangeRate `xml:"Rate"`
}

// NationalBankOfRomaniaExchangeRate represents the exchange rate data from national bank of Romania
type NationalBankOfRomaniaExchangeRate struct {
	Currency   string `xml:"currency,attr"`
	Multiplier string `xml:"multiplier,attr"`
	Rate       string `xml:",chardata"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from national bank of Romania
func (e *NationalBankOfRomaniaExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if e.Header == nil || e.Body == nil {
		log.Errorf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRateResponse] header or body is empty")
		return nil
	}

	if len(e.Body.AllExchangeRates) < 1 {
		log.Errorf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil
	}

	latestNationalBankOfRomaniaExchangeRate := e.Body.AllExchangeRates[0]

	if len(latestNationalBankOfRomaniaExchangeRate.ExchangeRates) < 1 {
		log.Errorf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(latestNationalBankOfRomaniaExchangeRate.ExchangeRates))

	for i := 0; i < len(latestNationalBankOfRomaniaExchangeRate.ExchangeRates); i++ {
		exchangeRate := latestNationalBankOfRomaniaExchangeRate.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	timezone, err := time.LoadLocation(nationalBankOfRomaniaUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", nationalBankOfRomaniaUpdateDateTimezone)
		return nil
	}

	updateDateTime := e.Header.PublishingDate + " 13" // The data are updated in real time, shortly after 13:00, every banking day.
	updateTime, err := time.ParseInLocation(nationalBankOfRomaniaUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.Errorf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    nationalBankOfRomaniaDataSource,
		ReferenceUrl:  nationalBankOfRomaniaExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  e.Body.OrigCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from national bank of Romania
func (e *NationalBankOfRomaniaExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(e.Rate)

	if err != nil {
		log.Warnf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	unit := float64(1)

	if e.Multiplier != "" {
		unit, err = utils.StringToFloat64(e.Multiplier)

		if err != nil || unit <= 0 {
			log.Warnf(c, "[national_bank_of_romania_datasource.ToLatestExchangeRate] failed to parse unit, currency is %s, unit is %s", e.Currency, e.Multiplier)
			return nil
		}
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

// BuildRequests returns the national bank of Romania exchange rates http requests
func (e *NationalBankOfRomaniaDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", nationalBankOfRomaniaExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the national bank of Romania data source raw response
func (e *NationalBankOfRomaniaDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	nationalBankOfRomaniaData := &NationalBankOfRomaniaExchangeRateData{}
	err := xmlDecoder.Decode(nationalBankOfRomaniaData)

	if err != nil {
		log.Errorf(c, "[national_bank_of_romania_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := nationalBankOfRomaniaData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[national_bank_of_romania_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
