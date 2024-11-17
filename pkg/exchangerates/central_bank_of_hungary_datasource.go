package exchangerates

import (
	"bytes"
	"encoding/xml"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const centralBankOfHungaryExchangeRateServiceUrl = "http://www.mnb.hu/arfolyamok.asmx"
const centralBankOfHungaryExchangeRateServiceCurrentExchangeRatesSoapAction = "http://www.mnb.hu/webservices/MNBArfolyamServiceSoap/GetCurrentExchangeRates"
const centralBankOfHungaryExchangeRateReferenceUrl = "https://www.mnb.hu/en/arfolyamok"
const centralBankOfHungaryDataSource = "Magyar Nemzeti Bank"
const centralBankOfHungaryBaseCurrency = "HUF"

const centralBankOfHungaryUpdateDateFormat = "2006-01-02 15"
const centralBankOfHungaryUpdateDateTimezone = "Europe/Budapest"

// CentralBankOfHungaryDataSource defines the structure of exchange rates data source of central bank of Hungary
type CentralBankOfHungaryDataSource struct {
	ExchangeRatesDataSource
}

// CentralBankOfHungaryExchangeRateServiceResponse represents the response data of exchange rate service for central bank of Hungary
type CentralBankOfHungaryExchangeRateServiceResponse struct {
	XMLName                       xml.Name `xml:"Envelope"`
	GetCurrentExchangeRatesResult string   `xml:"Body>GetCurrentExchangeRatesResponse>GetCurrentExchangeRatesResult"`
}

// CentralBankOfHungaryCurrentExchangeRatesResult represents the current exchange rate result data from central bank of Hungary
type CentralBankOfHungaryCurrentExchangeRatesResult struct {
	XMLName          xml.Name                             `xml:"MNBCurrentExchangeRates"`
	AllExchangeRates []*CentralBankOfHungaryExchangeRates `xml:"Day"`
}

// CentralBankOfHungaryExchangeRates represents the exchange rates data from Danmarks Nationalbank
type CentralBankOfHungaryExchangeRates struct {
	Date          string                              `xml:"date,attr"`
	ExchangeRates []*CentralBankOfHungaryExchangeRate `xml:"Rate"`
}

// CentralBankOfHungaryExchangeRate represents the exchange rate data from central bank of Hungary
type CentralBankOfHungaryExchangeRate struct {
	Currency string `xml:"curr,attr"`
	Unit     string `xml:"unit,attr"`
	Rate     string `xml:",chardata"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from central bank of Hungary
func (e *CentralBankOfHungaryCurrentExchangeRatesResult) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.AllExchangeRates) < 1 {
		log.Errorf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRateResponse] all exchange rates is empty")
		return nil
	}

	latestCentralBankOfHungaryExchangeRate := e.AllExchangeRates[0]

	if len(latestCentralBankOfHungaryExchangeRate.ExchangeRates) < 1 {
		log.Errorf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRateResponse] exchange rates is empty")
		return nil
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(e.AllExchangeRates))

	for i := 0; i < len(latestCentralBankOfHungaryExchangeRate.ExchangeRates); i++ {
		exchangeRate := latestCentralBankOfHungaryExchangeRate.ExchangeRates[i]

		if _, exists := validators.AllCurrencyNames[exchangeRate.Currency]; !exists {
			continue
		}

		finalExchangeRate := exchangeRate.ToLatestExchangeRate(c)

		if finalExchangeRate == nil {
			continue
		}

		exchangeRates = append(exchangeRates, finalExchangeRate)
	}

	timezone, err := time.LoadLocation(centralBankOfHungaryUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRateResponse] failed to get timezone, timezone name is %s", centralBankOfHungaryUpdateDateTimezone)
		return nil
	}

	updateDateTime := latestCentralBankOfHungaryExchangeRate.Date + " 11" // The exchange rates are fixed at 11 am.
	updateTime, err := time.ParseInLocation(centralBankOfHungaryUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.Errorf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRateResponse] failed to parse update date, datetime is %s", updateDateTime)
		return nil
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    centralBankOfHungaryDataSource,
		ReferenceUrl:  centralBankOfHungaryExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  centralBankOfHungaryBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from central bank of Hungary
func (e *CentralBankOfHungaryExchangeRate) ToLatestExchangeRate(c core.Context) *models.LatestExchangeRate {
	rate, err := utils.StringToFloat64(strings.ReplaceAll(e.Rate, ",", "."))

	if err != nil {
		log.Warnf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRate] failed to parse rate, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRate] rate is invalid, currency is %s, rate is %s", e.Currency, e.Rate)
		return nil
	}

	unit, err := utils.StringToFloat64(e.Unit)

	if err != nil {
		log.Warnf(c, "[central_bank_of_hungary_datasource.ToLatestExchangeRate] failed to parse unit, currency is %s, unit is %s", e.Currency, e.Unit)
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

// BuildRequests returns the central bank of Hungary exchange rates http requests
func (e *CentralBankOfHungaryDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("POST", centralBankOfHungaryExchangeRateServiceUrl, bytes.NewReader([]byte(
		"<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
			"<s:Body>"+
			"<GetCurrentExchangeRates xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\"/>"+
			"</s:Body>"+
			"</s:Envelope>")))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", centralBankOfHungaryExchangeRateServiceCurrentExchangeRatesSoapAction)

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the central bank of Hungary data source raw response
func (e *CentralBankOfHungaryDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	responseXmlDecoder := xml.NewDecoder(bytes.NewReader(content))

	centralBankOfHungaryServiceResponse := &CentralBankOfHungaryExchangeRateServiceResponse{}
	err := responseXmlDecoder.Decode(centralBankOfHungaryServiceResponse)

	if err != nil {
		log.Errorf(c, "[central_bank_of_hungary_datasource.Parse] failed to parse service response xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if len(centralBankOfHungaryServiceResponse.GetCurrentExchangeRatesResult) < 1 {
		log.Errorf(c, "[central_bank_of_hungary_datasource.Parse] exchange rates response is empty")
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	resultXmlDecoder := xml.NewDecoder(strings.NewReader(centralBankOfHungaryServiceResponse.GetCurrentExchangeRatesResult))

	centralBankOfHungaryExchangeRatesResult := &CentralBankOfHungaryCurrentExchangeRatesResult{}
	err = resultXmlDecoder.Decode(centralBankOfHungaryExchangeRatesResult)

	if err != nil {
		log.Errorf(c, "[central_bank_of_hungary_datasource.Parse] failed to parse exchange rates response xml data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := centralBankOfHungaryExchangeRatesResult.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[central_bank_of_hungary_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
