package exchangerates

import (
	"encoding/xml"
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

const nationalBankOfKazakhstanExchangeRateUrl = "https://www.nationalbank.kz/rss/rates_all.xml"
const nationalBankOfKazakhstanExchangeRateReferenceUrl = "https://nationalbank.kz/en/exchangerates/ezhednevnye-oficialnye-rynochnye-kursy-valyut"
const nationalBankOfKazakhstanDataSource = "Қазақстан Республикасының Ұлттық Банкі"
const nationalBankOfKazakhstanBaseCurrency = "KZT"

const nationalBankOfKazakhstanUpdateDateFormat = "02.01.2006"
const nationalBankOfKazakhstanUpdateDateTimezone = "Asia/Almaty"

// NationalBankOfKazakhstanDataSource defines the structure of exchange rates data source of the national bank of Kazakhstan
type NationalBankOfKazakhstanDataSource struct {
	HttpExchangeRatesDataSource
}

// NationalBankOfKazakhstanExchangeRates represents the exchange rates data from the national bank of Kazakhstan
type NationalBankOfKazakhstanExchangeRates struct {
	Channel struct {
		Items []*NationalBankOfKazakhstanExchangeRate `xml:"item"`
	} `xml:"channel"`
}

// NationalBankOfKazakhstanExchangeRate represents the exchange rate data from the national bank of Kazakhstan
type NationalBankOfKazakhstanExchangeRate struct {
	Currency string `xml:"title"`
	Rate     string `xml:"description"`
	Unit     string `xml:"quant"`
	Date     string `xml:"pubDate"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from the national bank of Kazakhstan
func (e *NationalBankOfKazakhstanExchangeRates) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if e == nil || len(e.Channel.Items) < 1 {
		log.Errorf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	timezone, err := time.LoadLocation(nationalBankOfKazakhstanUpdateDateTimezone)
	if err != nil {
		log.Errorf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRateResponse] failed to load timezone, timezone name is %s", nationalBankOfKazakhstanUpdateDateTimezone)
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.Channel.Items))
	latestUpdateTime := int64(0)

	for i := 0; i < len(e.Channel.Items); i++ {
		exchangeRate := e.Channel.Items[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		updateTime, err := time.ParseInLocation(nationalBankOfKazakhstanUpdateDateFormat, exchangeRate.Date, timezone)

		if err != nil {
			log.Errorf(c, "[central_bank_of_kazakhstan_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Date)
			return nil
		}

		if updateTime.Unix() > latestUpdateTime {
			latestUpdateTime = updateTime.Unix()
		}

		finalRate := exchangeRate.ToLatestExchangeRate(c)
		if finalRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalRate)
	}

	return &models.LatestExchangeRateResponse{
		DataSource:    nationalBankOfKazakhstanDataSource,
		ReferenceUrl:  nationalBankOfKazakhstanExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  nationalBankOfKazakhstanBaseCurrency,
		ExchangeRates: exchangeRates,
	}
}

// ToLatestExchangeRate returns a data pair according to original data from the national bank of Kazakhstan
func (e *NationalBankOfKazakhstanExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(e.Rate)
	if err != nil {
		log.Warnf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	unit, err := utils.StringToFloat64(e.Unit)
	if err != nil {
		log.Warnf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRate] failed to parse unit, currency=%s, unit=%s", e.Currency, e.Unit)
	}

	if unit <= 0 {
		log.Warnf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRate] unit is less or equal zero, currency is %s, unit is %s", e.Currency, e.Unit)
		return nil
	}

	finalRate := unit / rate
	if math.IsInf(finalRate, 0) {
		log.Warnf(c, "[national_bank_of_kazakhstan_datasource.ToLatestExchangeRate] final exchange rate calculation failed, currency is %s, unit is %s, rate is %s", e.Currency, e.Unit, e.Rate)
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.Currency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the national bank of Kazakhstan exchange rates http requests
func (e *NationalBankOfKazakhstanDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", nationalBankOfKazakhstanExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the national bank of Kazakhstan data source raw response
func (e *NationalBankOfKazakhstanDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	nationalBankOfKazakhstanData := &NationalBankOfKazakhstanExchangeRates{}
	err := xml.Unmarshal(content, nationalBankOfKazakhstanData)

	if err != nil {
		log.Errorf(c, "[national_bank_of_kazakhstan_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := nationalBankOfKazakhstanData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[national_bank_of_kazakhstan_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
