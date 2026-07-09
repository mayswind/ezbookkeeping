package exchangerates

import (
	"encoding/json"
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

const (
	centralBankOfTheArgentineRepublicExchangeRateUrl          = "https://api.bcra.gob.ar/estadisticascambiarias/v1.0/Cotizaciones"
	centralBankOfTheArgentineRepublicExchangeRateReferenceUrl = "https://www.bcra.gob.ar/Estadisticas/EstadisticasCambiarias"
	centralBankOfTheArgentineRepublicDataSource               = "Central Bank of the Argentine Republic"
	centralBankOfTheArgentineRepublicBaseCurrency               = "ARS"
	centralBankOfTheArgentineRepublicDataUpdateDateFormat       = "2006-01-02"
	centralBankOfTheArgentineRepublicUsdCurrencyCode            = "USD"
)

// CentralBankOfTheArgentineRepublicDataSource defines the structure of exchange rates data source of central bank of the Argentine Republic
type CentralBankOfTheArgentineRepublicDataSource struct {
	HttpExchangeRatesDataSource
}

// CentralBankOfTheArgentineRepublicExchangeRateData represents the whole data from BCRA Cotizaciones API
type CentralBankOfTheArgentineRepublicExchangeRateData struct {
	Status  int                                                 `json:"status"`
	Results CentralBankOfTheArgentineRepublicResultsData        `json:"results"`
}

// CentralBankOfTheArgentineRepublicResultsData represents the results section from BCRA Cotizaciones API
type CentralBankOfTheArgentineRepublicResultsData struct {
	Fecha   string                                              `json:"fecha"`
	Detalle []CentralBankOfTheArgentineRepublicDetalleItem      `json:"detalle"`
}

// CentralBankOfTheArgentineRepublicDetalleItem represents a single currency quote from BCRA Cotizaciones API
type CentralBankOfTheArgentineRepublicDetalleItem struct {
	CodigoMoneda   string  `json:"codigoMoneda"`
	Descripcion    string  `json:"descripcion"`
	TipoPase       float64 `json:"tipoPase"`
	TipoCotizacion float64 `json:"tipoCotizacion"`
}

func (e *CentralBankOfTheArgentineRepublicExchangeRateData) findUsdArsPerUsd(c core.Context) float64 {
	for i := 0; i < len(e.Results.Detalle); i++ {
		item := e.Results.Detalle[i]

		if item.CodigoMoneda != centralBankOfTheArgentineRepublicUsdCurrencyCode {
			continue
		}

		if item.TipoCotizacion <= 0 {
			log.Warnf(c, "[central_bank_of_the_argentine_republic_datasource.ToLatestExchangeRateResponse] tipo cotizacion is invalid for USD, value is %f", item.TipoCotizacion)
			return 0
		}

		return item.TipoCotizacion
	}

	return 0
}

// ToLatestExchangeRateResponse returns a view-object according to original data from central bank of the Argentine Republic
func (e *CentralBankOfTheArgentineRepublicExchangeRateData) ToLatestExchangeRateResponse(c core.Context) *models.LatestExchangeRateResponse {
	if len(e.Results.Detalle) < 1 {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.ToLatestExchangeRateResponse] detalle is empty")
		return nil
	}

	arsPerUsd := e.findUsdArsPerUsd(c)

	if arsPerUsd <= 0 {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.ToLatestExchangeRateResponse] USD quote is missing or invalid")
		return nil
	}

	if _, exists := validators.AllCurrencyNames[centralBankOfTheArgentineRepublicUsdCurrencyCode]; !exists {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.ToLatestExchangeRateResponse] USD is not a supported currency")
		return nil
	}

	finalRate := 1 / arsPerUsd

	if math.IsInf(finalRate, 0) {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.ToLatestExchangeRateResponse] rate is invalid, ars per usd is %f", arsPerUsd)
		return nil
	}

	updateTime, err := time.Parse(centralBankOfTheArgentineRepublicDataUpdateDateFormat, e.Results.Fecha)

	if err != nil {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.ToLatestExchangeRateResponse] failed to parse update date, date is %s", e.Results.Fecha)
		return nil
	}

	return &models.LatestExchangeRateResponse{
		DataSource:   centralBankOfTheArgentineRepublicDataSource,
		ReferenceUrl: centralBankOfTheArgentineRepublicExchangeRateReferenceUrl,
		UpdateTime:   updateTime.Unix(),
		BaseCurrency: centralBankOfTheArgentineRepublicBaseCurrency,
		ExchangeRates: models.LatestExchangeRateSlice{
			{
				Currency: centralBankOfTheArgentineRepublicUsdCurrencyCode,
				Rate:     utils.Float64ToString(finalRate),
			},
		},
	}
}

// BuildRequests returns the central bank of the Argentine Republic exchange rates http requests
func (e *CentralBankOfTheArgentineRepublicDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", centralBankOfTheArgentineRepublicExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the central bank of the Argentine Republic data source raw response
func (e *CentralBankOfTheArgentineRepublicDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	exchangeRateData := &CentralBankOfTheArgentineRepublicExchangeRateData{}
	err := json.Unmarshal(content, exchangeRateData)

	if err != nil {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResponse := exchangeRateData.ToLatestExchangeRateResponse(c)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[central_bank_of_the_argentine_republic_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
