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

const norgesBankExchangeRateUrl = "https://data.norges-bank.no/api/data/EXR/B..NOK.SP?format=sdmx-compact-2.1&lastNObservations=1"
const norgesBankExchangeRateReferenceUrl = "https://www.norges-bank.no/en/topics/Statistics/exchange_rates/"
const norgesBankDataSource = "Norges Bank"
const norgesBankBaseCurrency = "NOK"

const norgesBankUpdateDateFormat = "2006-01-02 15"
const norgesBankUpdateDateTimezone = "Europe/Oslo"

// NorgesBankDataSource defines the structure of exchange rates data source of Norges Bank
type NorgesBankDataSource struct {
	ExchangeRatesDataSource
}

// NorgesBankExchangeRateData represents the whole data from Norges Bank
type NorgesBankExchangeRateData struct {
	XMLName xml.Name                       `xml:"StructureSpecificData"`
	DataSet *NorgesBankExchangeRateDataSet `xml:"DataSet"`
}

// NorgesBankExchangeRateDataSet represents the dataset for exchange rates data of Norges Bank
type NorgesBankExchangeRateDataSet struct {
	ExchangeRates []*NorgesBankExchangeRate `xml:"Series"`
}

// NorgesBankExchangeRate represents the exchange rate data from Norges Bank
type NorgesBankExchangeRate struct {
	BaseCurrency   string                               `xml:"BASE_CUR,attr"`
	TargetCurrency string                               `xml:"QUOTE_CUR,attr"`
	UnitExponent   string                               `xml:"UNIT_MULT,attr"`
	Observations   []*NorgesBankExchangeRateObservation `xml:"Obs"`
}

// NorgesBankExchangeRateObservation represents the observation data of exchange rate data from Norges Bank
type NorgesBankExchangeRateObservation struct {
	Date string `xml:"TIME_PERIOD,attr"`
	Rate string `xml:"OBS_VALUE,attr"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from Norges Bank
func (e *NorgesBankExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if e.DataSet == nil || len(e.DataSet.ExchangeRates) < 1 {
		log.Errorf(c, "[norges_bank_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil
	}

	timezone, err := time.LoadLocation(norgesBankUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[norges_bank_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", norgesBankUpdateDateTimezone)
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.DataSet.ExchangeRates))
	latestUpdateTime := int64(0)

	for i := 0; i < len(e.DataSet.ExchangeRates); i++ {
		exchangeRate := e.DataSet.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.BaseCurrency]; !exists {
			continue
		}

		if exchangeRate.TargetCurrency != norgesBankBaseCurrency {
			continue
		}

		if len(exchangeRate.Observations) < 1 {
			continue
		}

		updateDateTime := exchangeRate.Observations[0].Date + " 16" // Publication time of daily exchange rates is approximately 16:00 CET.
		updateTime, err := time.ParseInLocation(norgesBankUpdateDateFormat, updateDateTime, timezone)

		if err != nil {
			log.Errorf(c, "[norges_bank_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", exchangeRate.Observations[0].Date)
			return nil
		}

		if updateTime.Unix() > latestUpdateTime {
			latestUpdateTime = updateTime.Unix()
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c, exchangeRate.Observations[0].Rate)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    norgesBankDataSource,
		ReferenceUrl:  norgesBankExchangeRateReferenceUrl,
		UpdateTime:    latestUpdateTime,
		BaseCurrency:  norgesBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from Norges Bank
func (e *NorgesBankExchangeRate) ToLatestExchangeRate(c core.Context, exchangeRate string) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(exchangeRate)

	if err != nil {
		log.Warnf(c, "[norges_bank_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.BaseCurrency, exchangeRate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[norges_bank_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.BaseCurrency, exchangeRate)
		return nil
	}

	unitExponent, err := utils.StringToInt(e.UnitExponent)

	if err != nil {
		log.Warnf(c, "[norges_bank_datasource.ToLatestExchangeRate] failed to parse unit, currency is %s, unit exponent is %s", e.BaseCurrency, e.UnitExponent)
		return nil
	}

	finalRate := 1 / rate

	if unitExponent > 0 {
		finalRate = finalRate / math.Pow10(-unitExponent)
	}

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.BaseCurrency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// GetRequestUrls returns the Norges Bank data source urls
func (e *NorgesBankDataSource) GetRequestUrls() []string {
	return []string{norgesBankExchangeRateUrl}
}

// Parse returns the common response entity according to the Norges Bank data source raw response
func (e *NorgesBankDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	norgesBankData := &NorgesBankExchangeRateData{}
	err := xmlDecoder.Decode(norgesBankData)

	if err != nil {
		log.Errorf(c, "[norges_bank_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := norgesBankData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[norges_bank_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
