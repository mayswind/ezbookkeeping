package exchangerates

import (
	"encoding/xml"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const reserveBankOfAustraliaExchangeRateUrl = "https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml"
const reserveBankOfAustraliaExchangeRateReferenceUrl = "https://www.rba.gov.au/statistics/frequency/exchange-rates.html"
const reserveBankOfAustraliaDataSource = "Reserve Bank of Australia"
const reserveBankOfAustraliaBaseCurrency = "AUD"

const reserveBankOfAustraliaDataUpdateDateFormat = "2006-01-02T15:04:05Z07:00"

// ReserveBankOfAustraliaDataSource defines the structure of exchange rates data source of the reserve bank of Australia
type ReserveBankOfAustraliaDataSource struct {
	ExchangeRatesDataSource
}

// ReserveBankOfAustraliaData represents the whole data from the reserve bank of Australia
type ReserveBankOfAustraliaData struct {
	XMLName xml.Name                          `xml:"RDF"`
	Channel *ReserveBankOfAustraliaRssChannel `xml:"channel"`
	Items   []*ReserveBankOfAustraliaRssItem  `xml:"item"`
}

// ReserveBankOfAustraliaRssChannel represents the rss channel from the reserve bank of Australia
type ReserveBankOfAustraliaRssChannel struct {
	Date string `xml:"date"`
}

// ReserveBankOfAustraliaRssItem represents the rss item from the reserve bank of Australia
type ReserveBankOfAustraliaRssItem struct {
	Statistics *ReserveBankOfAustraliaItemStatistics `xml:"statistics"`
}

// ReserveBankOfAustraliaItemStatistics represents the item statistics from the reserve bank of Australia
type ReserveBankOfAustraliaItemStatistics struct {
	ExchangeRate *ReserveBankOfAustraliaExchangeRate `xml:"exchangeRate"`
}

// ReserveBankOfAustraliaExchangeRate represents the exchange rate from the reserve bank of Australia
type ReserveBankOfAustraliaExchangeRate struct {
	BaseCurrency   string                                         `xml:"baseCurrency"`
	TargetCurrency string                                         `xml:"targetCurrency"`
	Observation    *ReserveBankOfAustraliaExchangeRateObservation `xml:"observation"`
}

// ReserveBankOfAustraliaExchangeRateObservation represents the exchange rate data from the reserve bank of Australia
type ReserveBankOfAustraliaExchangeRateObservation struct {
	Value string `xml:"value"`
	Unit  string `xml:"unit"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from the reserve bank of Australia
func (e *ReserveBankOfAustraliaData) ToLatestExchangeRateResponse(c *core.Context) *models.LatestExchangeRateResponse {
	if e.Channel == nil {
		log.ErrorfWithRequestId(c, "[reserve_bank_of_australia_datasource.ToLatestExchangeRateResponse] rss channel does not exist")
		return nil
	}

	if len(e.Items) < 1 {
		log.ErrorfWithRequestId(c, "[reserve_bank_of_australia_datasource.ToLatestExchangeRateResponse] rss items is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.Items))

	for i := 0; i < len(e.Items); i++ {
		item := e.Items[i]

		if item.Statistics == nil || item.Statistics.ExchangeRate == nil || item.Statistics.ExchangeRate.Observation == nil {
			continue
		}

		if item.Statistics.ExchangeRate.BaseCurrency != reserveBankOfAustraliaBaseCurrency || item.Statistics.ExchangeRate.Observation.Unit != reserveBankOfAustraliaBaseCurrency {
			continue
		}

		if _, exists := validators.AllCurrencyNames[item.Statistics.ExchangeRate.TargetCurrency]; !exists {
			continue
		}

		if _, err := utils.StringToFloat64(item.Statistics.ExchangeRate.Observation.Value); err != nil {
			continue
		}

		exchangeRates = append(exchangeRates, item.Statistics.ExchangeRate.ToLatestExchangeRate())
	}

	updateDateTime := e.Channel.Date
	updateTime, err := time.Parse(reserveBankOfAustraliaDataUpdateDateFormat, updateDateTime)

	if err != nil {
		log.ErrorfWithRequestId(c, "[reserve_bank_of_australia_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    reserveBankOfAustraliaDataSource,
		ReferenceUrl:  reserveBankOfAustraliaExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  reserveBankOfAustraliaBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from the reserve bank of Australia
func (e *ReserveBankOfAustraliaExchangeRate) ToLatestExchangeRate() *models.LatestExchangeRate {
	return &models.LatestExchangeRate{
		Currency: e.TargetCurrency,
		Rate:     e.Observation.Value,
	}
}

// GetRequestUrls returns the the reserve bank of Australia data source urls
func (e *ReserveBankOfAustraliaDataSource) GetRequestUrls() []string {
	return []string{reserveBankOfAustraliaExchangeRateUrl}
}

// Parse returns the common response entity according to the the reserve bank of Australia data source raw response
func (e *ReserveBankOfAustraliaDataSource) Parse(c *core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	reserveBankOfAustraliaData := &ReserveBankOfAustraliaData{}
	err := xml.Unmarshal(content, reserveBankOfAustraliaData)

	if err != nil {
		log.ErrorfWithRequestId(c, "[reserve_bank_of_australia_datasource.Parse] failed to parse xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := reserveBankOfAustraliaData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.ErrorfWithRequestId(c, "[reserve_bank_of_australia_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
