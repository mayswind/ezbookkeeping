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

const bankOfCanadaExchangeRateUrl = "https://www.bankofcanada.ca/valet/observations/group/FX_RATES_DAILY/json?recent=1"
const bankOfCanadaExchangeRateReferenceUrl = "https://www.bankofcanada.ca/rates/exchange/daily-exchange-rates/"
const bankOfCanadaDataSource = "Bank of Canada"
const bankOfCanadaBaseCurrency = "CAD"

const bankOfCanadaDataUpdateDateFormat = "2006-01-02 15:04"
const bankOfCanadaDataUpdateDateTimezone = "America/Toronto"

// BankOfCanadaDataSource defines the structure of exchange rates data source of bank of Canada
type BankOfCanadaDataSource struct {
	ExchangeRatesDataSource
}

// BankOfCanadaExchangeRateData represents the whole data from bank of Canada
type BankOfCanadaExchangeRateData struct {
	Observations []BankOfCanadaObservationData `json:"observations"`
}

// BankOfCanadaObservationData represents the observation data from bank of Canada
type BankOfCanadaObservationData map[string]any

// ToLatestExchangeRateResponse returns a view-object according to original data from bank of Canada
func (e *BankOfCanadaExchangeRateData) ToLatestExchangeRateResponse(c *core.Context) *models.LatestExchangeRateResponse {
	if len(e.Observations) < 1 {
		log.ErrorfWithRequestId(c, "[bank_of_canada_datasource.ToLatestExchangeRateResponse] observations is empty")
		return nil
	}

	exchangeRateMap := make(map[string]string)
	latestUpdateDate := ""

	for i := 0; i < len(e.Observations); i++ {
		observation := e.Observations[i]
		updateDateData := observation["d"]

		if updateDate, ok := updateDateData.(string); ok {
			if latestUpdateDate == "" || strings.Compare(updateDate, latestUpdateDate) > 0 {
				latestUpdateDate = updateDate
			}
		}

		for typeName, exchangeRateData := range observation {
			if len(typeName) < 8 || !strings.HasPrefix(typeName, "FX") || !strings.HasSuffix(typeName, bankOfCanadaBaseCurrency) {
				continue
			}

			currencyCode := utils.SubString(typeName, 2, 3)

			if data, ok := exchangeRateData.(map[string]any); ok {
				exchangeRate := data["v"]

				if exchangeRateValue, ok2 := exchangeRate.(string); ok2 {
					exchangeRateMap[currencyCode] = exchangeRateValue
				}
			}
		}
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(exchangeRateMap))

	for currencyCode, exchangeRate := range exchangeRateMap {
		if _, exists := validators.AllCurrencyNames[currencyCode]; !exists {
			continue
		}

		rate, err := utils.StringToFloat64(exchangeRate)

		if err != nil {
			log.WarnfWithRequestId(c, "[bank_of_canada_datasource.ToLatestExchangeRateResponse] failed to parse rate, rate is %s", exchangeRate)
			continue
		}

		if rate <= 0 {
			log.WarnfWithRequestId(c, "[bank_of_canada_datasource.ToLatestExchangeRateResponse] rate is invalid, rate is %s", exchangeRate)
			continue
		}

		finalRate := 1 / rate

		if math.IsInf(finalRate, 0) {
			continue
		}

		exchangeRates = append(exchangeRates, &models.LatestExchangeRate{
			Currency: currencyCode,
			Rate:     utils.Float64ToString(finalRate),
		})
	}

	timezone, err := time.LoadLocation(bankOfCanadaDataUpdateDateTimezone)

	if err != nil {
		log.ErrorfWithRequestId(c, "[bank_of_canada_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", bankOfCanadaDataUpdateDateTimezone)
		return nil
	}

	updateDateTime := latestUpdateDate + " 16:30" // Daily average exchange rates - published once each business day by 16:30 ET.
	updateTime, err := time.ParseInLocation(bankOfCanadaDataUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.ErrorfWithRequestId(c, "[bank_of_canada_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    bankOfCanadaDataSource,
		ReferenceUrl:  bankOfCanadaExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  bankOfCanadaBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// GetRequestUrls returns the bank of Canada data source urls
func (e *BankOfCanadaDataSource) GetRequestUrls() []string {
	return []string{bankOfCanadaExchangeRateUrl}
}

// Parse returns the common response entity according to the bank of Canada data source raw response
func (e *BankOfCanadaDataSource) Parse(c *core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	bankOfCanadaData := &BankOfCanadaExchangeRateData{}
	err := json.Unmarshal(content, bankOfCanadaData)

	if err != nil {
		log.ErrorfWithRequestId(c, "[bank_of_canada_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := bankOfCanadaData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.ErrorfWithRequestId(c, "[bank_of_canada_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
