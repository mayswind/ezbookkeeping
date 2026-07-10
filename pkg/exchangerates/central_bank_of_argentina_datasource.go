package exchangerates

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const centralBankOfArgentinaExchangeRateUrl = "https://api.bcra.gob.ar/estadisticascambiarias/v1.0/Cotizaciones"
const centralBankOfArgentinaExchangeRateReferenceUrl = "https://www.bcra.gob.ar/en/central-bank-api-catalog/"
const centralBankOfArgentinaDataSource = "Banco Central de la República Argentina"
const centralBankOfArgentinaBaseCurrency = "ARS"

const centralBankOfArgentinaDataUpdateDateFormat = "2006-01-02"
const centralBankOfArgentinaUpdateDateTimezone = "America/Buenos_Aires"

// centralBankOfArgentinaPerUnitQuantityPattern matches the per-unit quantity declared in the descripcion
var centralBankOfArgentinaPerUnitQuantityPattern = regexp.MustCompile(`C/(\d+(?:\.\d+)*)\s+UNIDADES`)

// CentralBankOfArgentinaDataSource defines the structure of exchange rates data source of the central bank of Argentina
type CentralBankOfArgentinaDataSource struct {
	HttpExchangeRatesDataSource
}

// CentralBankOfArgentinaExchangeRateData represents the whole data from BCRA Cotizaciones API
type CentralBankOfArgentinaExchangeRateData struct {
	Status  int                               `json:"status"`
	Results CentralBankOfArgentinaResultsData `json:"results"`
}

// CentralBankOfArgentinaResultsData represents the results section from BCRA Cotizaciones API
type CentralBankOfArgentinaResultsData struct {
	Date    string                             `json:"fecha"`
	Details []CentralBankOfArgentinaDetailItem `json:"detalle"`
}

// CentralBankOfArgentinaDetailItem represents a single currency quote from BCRA Cotizaciones API
type CentralBankOfArgentinaDetailItem struct {
	CurrencyCode string  `json:"codigoMoneda"`
	Descripcion  string  `json:"descripcion"`
	USDCrossRate float64 `json:"tipoPase"`
	ARSQuoteRate float64 `json:"tipoCotizacion"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from the central bank of Argentina
func (e *CentralBankOfArgentinaExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.Results.Details) < 1 {
		log.Errorf(c, "[central_bank_of_argentina_datasource.ToLatestExchangeRateResponse] detalle is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.Results.Details))

	for i := 0; i < len(e.Results.Details); i++ {
		exchangeRate := e.Results.Details[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.CurrencyCode]; !exists {
			continue
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	timezone, err := time.LoadLocation(centralBankOfArgentinaUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_argentina_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", centralBankOfArgentinaUpdateDateTimezone)
		return nil
	}

	updateTime, err := time.ParseInLocation(centralBankOfArgentinaDataUpdateDateFormat, e.Results.Date, timezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_argentina_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", e.Results.Date)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    centralBankOfArgentinaDataSource,
		ReferenceUrl:  centralBankOfArgentinaExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  centralBankOfArgentinaBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from the central bank of Argentina
func (e *CentralBankOfArgentinaDetailItem) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	if e.ARSQuoteRate <= 0 {
		log.Warnf(c, "[central_bank_of_argentina_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %f", e.CurrencyCode, e.ARSQuoteRate)
		return nil
	}

	unit := 1.0

	if e.Descripcion != "" {
		matches := centralBankOfArgentinaPerUnitQuantityPattern.FindStringSubmatch(e.Descripcion)

		if len(matches) > 1 {
			parsedUnit, err := utils.StringToFloat64(strings.ReplaceAll(matches[1], ".", ""))

			if err != nil {
				log.Warnf(c, "[central_bank_of_argentina_datasource.ToLatestExchangeRate] failed to parse per-unit quantity, currency is %s, descripcion is %s", e.CurrencyCode, e.Descripcion)
				return nil
			}

			if parsedUnit <= 0 {
				log.Warnf(c, "[central_bank_of_argentina_datasource.ToLatestExchangeRate] per-unit quantity is invalid, currency is %s, descripcion is %s", e.CurrencyCode, e.Descripcion)
				return nil
			}

			unit = parsedUnit
		}
	}

	finalRate := unit / e.ARSQuoteRate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: e.CurrencyCode,
		Rate:     utils.Float64ToString(finalRate),
	}
}

// BuildRequests returns the central bank of Argentina exchange rates http requests
func (e *CentralBankOfArgentinaDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", centralBankOfArgentinaExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the central bank of Argentina data source raw response
func (e *CentralBankOfArgentinaDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	exchangeRateData := &CentralBankOfArgentinaExchangeRateData{}
	err := json.Unmarshal(content, exchangeRateData)

	if err != nil {
		log.Errorf(c, "[central_bank_of_argentina_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := exchangeRateData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[central_bank_of_argentina_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
