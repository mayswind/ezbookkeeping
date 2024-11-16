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

const swissNationalBankExchangeRateUrl = "https://www.snb.ch/public/en/rss/exchangeRates"
const swissNationalBankExchangeRateReferenceUrl = "https://www.snb.ch/en/the-snb/mandates-goals/statistics/statistics-pub/current_interest_exchange_rates"
const swissNationalBankDataSource = "Swiss National Bank"
const swissNationalBankBaseCurrency = "CHF"

const swissNationalBankDataUpdateDateFormat = "Mon, 02 Jan 2006 15:04:05 MST"
const swissNationalBankExchangeRatePeriodDateFormat = "2006-01-02"

// SwissNationalBankDataSource defines the structure of exchange rates data source of the reserve Swiss National Bank
type SwissNationalBankDataSource struct {
	ExchangeRatesDataSource
}

// SwissNationalBankData represents the whole data from the reserve Swiss National Bank
type SwissNationalBankData struct {
	XMLName xml.Name                     `xml:"rss"`
	Channel *SwissNationalBankRssChannel `xml:"channel"`
}

// SwissNationalBankRssChannel represents the rss channel from the reserve Swiss National Bank
type SwissNationalBankRssChannel struct {
	PublishDate string                          `xml:"pubDate"`
	Items       []*SwissNationalBankChannelItem `xml:"item"`
}

// SwissNationalBankChannelItem represents the channel item from the reserve Swiss National Bank
type SwissNationalBankChannelItem struct {
	Statistics *SwissNationalBankItemStatistics `xml:"statistics"`
}

// SwissNationalBankItemStatistics represents the item statistics from the reserve Swiss National Bank
type SwissNationalBankItemStatistics struct {
	ExchangeRate *SwissNationalBankExchangeRate `xml:"exchangeRate"`
}

// SwissNationalBankExchangeRate represents the exchange rate from the reserve Swiss National Bank
type SwissNationalBankExchangeRate struct {
	BaseCurrency      string                                          `xml:"baseCurrency"`
	TargetCurrency    string                                          `xml:"targetCurrency"`
	Observation       *SwissNationalBankExchangeRateObservation       `xml:"observation"`
	ObservationPeriod *SwissNationalBankExchangeRateObservationPeriod `xml:"observationPeriod"`
}

// SwissNationalBankExchangeRateObservation represents the exchange rate data from the reserve Swiss National Bank
type SwissNationalBankExchangeRateObservation struct {
	Value        string `xml:"value"`
	Unit         string `xml:"unit"`
	UnitExponent string `xml:"unit_mult"`
}

// SwissNationalBankExchangeRateObservationPeriod represents the exchange rate period data from the reserve Swiss National Bank
type SwissNationalBankExchangeRateObservationPeriod struct {
	Period string `xml:"period"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from the reserve Swiss National Bank
func (e *SwissNationalBankData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if e.Channel == nil {
		log.Errorf(c, "[swiss_national_bank_datasource.ToLatestExchangeRateResponse] rss channel does not exist")
		return nil
	}

	if len(e.Channel.Items) < 1 {
		log.Errorf(c, "[swiss_national_bank_datasource.ToLatestExchangeRateResponse] channel items is empty")
		return nil
	}

	latestCurrencyExchangeRateDate := make(map[string]int64)
	latestExchangeRates := make(map[string]*models.LatestExchangeRate)

	for i := 0; i < len(e.Channel.Items); i++ {
		item := e.Channel.Items[i]

		if item.Statistics == nil || item.Statistics.ExchangeRate == nil || item.Statistics.ExchangeRate.Observation == nil || item.Statistics.ExchangeRate.ObservationPeriod == nil {
			continue
		}

		if item.Statistics.ExchangeRate.BaseCurrency != swissNationalBankBaseCurrency || item.Statistics.ExchangeRate.Observation.Unit != swissNationalBankBaseCurrency {
			continue
		}

		if _, exists := validators.AllCurrencyNames[item.Statistics.ExchangeRate.TargetCurrency]; !exists {
			continue
		}

		date, err := time.Parse(swissNationalBankExchangeRatePeriodDateFormat, item.Statistics.ExchangeRate.ObservationPeriod.Period)

		if err != nil {
			log.Warnf(c, "[swiss_national_bank_datasource.ToLatestExchangeRateResponse] failed to parse exchange rate period date, period is %s", item.Statistics.ExchangeRate.ObservationPeriod.Period)
			continue
		}

		currency := item.Statistics.ExchangeRate.TargetCurrency
		latestDate, exists := latestCurrencyExchangeRateDate[currency]

		if !exists || date.Unix() > latestDate {
			finalExchangeRate := item.Statistics.ExchangeRate.ToLatestExchangeRate(c)

			if finalExchangeRate != nil {
				latestCurrencyExchangeRateDate[currency] = date.Unix()
				latestExchangeRates[currency] = finalExchangeRate
			}
		}
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.Channel.Items))

	for _, exchangeRate := range latestExchangeRates {
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	updateDateTime := e.Channel.PublishDate
	updateTime, err := time.Parse(swissNationalBankDataUpdateDateFormat, updateDateTime)

	if err != nil {
		log.Errorf(c, "[swiss_national_bank_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    swissNationalBankDataSource,
		ReferenceUrl:  swissNationalBankExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  swissNationalBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from the reserve Swiss National Bank
func (e *SwissNationalBankExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(e.Observation.Value)

	if err != nil {
		log.Warnf(c, "[swiss_national_bank_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.TargetCurrency, e.Observation.Value)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[swiss_national_bank_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.TargetCurrency, e.Observation.Value)
		return nil
	}

	unitExponent, err := utils.StringToInt(e.Observation.UnitExponent)

	if err != nil {
		log.Warnf(c, "[swiss_national_bank_datasource.ToLatestExchangeRate] failed to parse unit, currency is %s, unit exponent is %s", e.TargetCurrency, e.Observation.UnitExponent)
		return nil
	}

	finalRate := 1 / rate

	if unitExponent > 1 {
		finalRate = finalRate / math.Pow10(unitExponent-1)
	} else if unitExponent < 0 {
		finalRate = finalRate * math.Pow10(-unitExponent)
	}

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.TargetCurrency,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// GetRequestUrls returns the the reserve Swiss National Bank data source urls
func (e *SwissNationalBankDataSource) GetRequestUrls() []string {
	return []string{swissNationalBankExchangeRateUrl}
}

// Parse returns the common response entity according to the the reserve Swiss National Bank data source raw response
func (e *SwissNationalBankDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	xmlDecoder := xml.NewDecoder(bytes.NewReader(content))
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	swissNationalBankData := &SwissNationalBankData{}
	err := xmlDecoder.Decode(swissNationalBankData)

	if err != nil {
		log.Errorf(c, "[swiss_national_bank_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := swissNationalBankData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[swiss_national_bank_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
