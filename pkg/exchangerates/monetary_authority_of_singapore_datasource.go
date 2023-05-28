package exchangerates

import (
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const monetaryAuthorityOfSingaporeExchangeRateUrl = "https://eservices.mas.gov.sg/api/action/datastore/search.json?resource_id=95932927-c8bc-4e7a-b484-68a66a24edfe&sort=end_of_day+desc&limit=1"
const monetaryAuthorityOfSingaporeExchangeRateReferenceUrl = "https://eservices.mas.gov.sg/Statistics/msb/ExchangeRates.aspx"
const monetaryAuthorityOfSingaporeDataSource = "Monetary Authority of Singapore"
const monetaryAuthorityOfSingaporeBaseCurrency = "SGD"

const monetaryAuthorityOfSingaporeDataUpdateDateFormat = "2006-01-02 15"
const monetaryAuthorityOfSingaporeDataUpdateDateTimezone = "Asia/Singapore"

// MonetaryAuthorityOfSingaporeDataSource defines the structure of exchange rates data source of Monetary Authority of Singapore
type MonetaryAuthorityOfSingaporeDataSource struct {
	ExchangeRatesDataSource
}

// MonetaryAuthorityOfSingaporeExchangeRateData represents the whole data from Monetary Authority of Singapore
type MonetaryAuthorityOfSingaporeExchangeRateData struct {
	Success bool                                `json:"success"`
	Result  *MonetaryAuthorityOfSingaporeResult `json:"result"`
}

// MonetaryAuthorityOfSingaporeResult represents the actual result from Monetary Authority of Singapore
type MonetaryAuthorityOfSingaporeResult struct {
	Records []MonetaryAuthorityOfSingaporeRecord `json:"records"`
}

// MonetaryAuthorityOfSingaporeRecord represents the record from Monetary Authority of Singapore
type MonetaryAuthorityOfSingaporeRecord map[string]string

// ToLatestExchangeRateResponse returns a view-object according to original data from Monetary Authority of Singapore
func (e *MonetaryAuthorityOfSingaporeExchangeRateData) ToLatestExchangeRateResponse(c *core.Context) *models.LatestExchangeRateResponse {
	if !e.Success {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.ToLatestExchangeRateResponse] response is not success")
		return nil
	}

	if e.Result == nil {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.ToLatestExchangeRateResponse] result is null")
		return nil
	}

	if len(e.Result.Records) < 1 {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.ToLatestExchangeRateResponse] records is empty")
		return nil
	}

	lastDayRecord := e.Result.Records[0]
	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(lastDayRecord))
	latestUpdateDate := ""

	for key, value := range lastDayRecord {
		if key == "end_of_day" {
			latestUpdateDate = value
			continue
		}

		exchangeRate := e.parseExchangeRateResponse(c, key, value)

		if exchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, exchangeRate)
	}

	timezone, err := time.LoadLocation(monetaryAuthorityOfSingaporeDataUpdateDateTimezone)

	if err != nil {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", monetaryAuthorityOfSingaporeDataUpdateDateTimezone)
		return nil
	}

	updateDateTime := latestUpdateDate + " 12" // These rates are the average of buying and selling interbank rates quoted around midday in Singapore
	updateTime, err := time.ParseInLocation(monetaryAuthorityOfSingaporeDataUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    monetaryAuthorityOfSingaporeDataSource,
		ReferenceUrl:  monetaryAuthorityOfSingaporeExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  monetaryAuthorityOfSingaporeBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

func (e *MonetaryAuthorityOfSingaporeExchangeRateData) parseExchangeRateResponse(c *core.Context, key string, value string) *models.LatestExchangeRate {
	if !strings.Contains(key, "_") {
		return nil
	}

	items := strings.Split(key, "_")

	if len(items) < 2 {
		return nil
	}

	fromCurrencyCode := strings.ToUpper(items[0])
	toCurrencyCode := strings.ToUpper(items[1])

	if _, exists := validators.AllCurrencyNames[fromCurrencyCode]; !exists {
		return nil
	}

	if toCurrencyCode != monetaryAuthorityOfSingaporeBaseCurrency {
		return nil
	}

	rate, err := utils.StringToFloat64(value)

	if err != nil {
		log.WarnfWithRequestId(c, "[monetary_authority_of_singapore_datasource.parseExchangeRateResponse] failed to parse rate, rate is %s", value)
		return nil
	}

	if rate <= 0 {
		log.WarnfWithRequestId(c, "[monetary_authority_of_singapore_datasource.parseExchangeRateResponse] rate is invalid, rate is %s", value)
		return nil
	}

	finalRate := 1 / rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	if len(items) == 3 && items[2] == "100" {
		finalRate = finalRate * 100
	}

	return &models.LatestExchangeRate{
		Currency: fromCurrencyCode,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// GetRequestUrls returns the Monetary Authority of Singapore data source urls
func (e *MonetaryAuthorityOfSingaporeDataSource) GetRequestUrls() []string {
	return []string{monetaryAuthorityOfSingaporeExchangeRateUrl}
}

// Parse returns the common response entity according to the Monetary Authority of Singapore data source raw response
func (e *MonetaryAuthorityOfSingaporeDataSource) Parse(c *core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	monetaryAuthorityOfSingaporeData := &MonetaryAuthorityOfSingaporeExchangeRateData{}
	err := json.Unmarshal(content, monetaryAuthorityOfSingaporeData)

	if err != nil {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := monetaryAuthorityOfSingaporeData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.ErrorfWithRequestId(c, "[monetary_authority_of_singapore_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
