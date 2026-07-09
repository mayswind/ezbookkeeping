package exchangerates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const (
	argentinaDatosExchangeRateReferenceUrl = "https://argentinadatos.com/docs/operations/get-cotizaciones-dolares"
	argentinaDatosExchangeRateUrlTemplate  = "https://api.argentinadatos.com/v1/cotizaciones/dolares/%s"
	argentinaDatosDataSource               = "ArgentinaDatos"
	argentinaDatosBaseCurrency             = "ARS"

	// ArgentinaDatos API rate type values
	argentinaDatosApiRateTypeBuy  = "compra"
	argentinaDatosApiRateTypeSell = "venta"
)

var argentinaDatosValidExchangeHouses = map[string]bool{
	"oficial":         true,
	"blue":            true,
	"bolsa":           true,
	"contadoconliqui": true,
	"cripto":          true,
	"mayorista":       true,
	"solidario":       true,
	"turista":         true,
}

// ArgentinaDatosDataSource defines the structure of exchange rates data source of ArgentinaDatos API
type ArgentinaDatosDataSource struct {
	HttpExchangeRatesDataSource
	exchangeHouse string
	apiRateType   string
}

// ArgentinaDatosQuote represents a single dollar exchange rate quote from ArgentinaDatos API
type ArgentinaDatosQuote struct {
	ExchangeHouse string  `json:"casa"`
	Buy           float64 `json:"compra"`
	Sell          float64 `json:"venta"`
	Date          string  `json:"fecha"`
}

// IsValidArgentinaDatosExchangeHouse returns whether the given exchange house name is supported by ArgentinaDatos API
func IsValidArgentinaDatosExchangeHouse(exchangeHouse string) bool {
	return argentinaDatosValidExchangeHouses[strings.ToLower(exchangeHouse)]
}

func mapConfigRateTypeToArgentinaDatosApiRateType(rateType string) string {
	if rateType == settings.ArgentinaDatosRateTypeBuy {
		return argentinaDatosApiRateTypeBuy
	}

	return argentinaDatosApiRateTypeSell
}

func newArgentinaDatosDataSource(config *settings.Config) *ArgentinaDatosDataSource {
	exchangeHouse := config.ExchangeRatesArgentinaDatosExchangeHouse

	if exchangeHouse == "" {
		exchangeHouse = "blue"
	}

	rateType := config.ExchangeRatesArgentinaDatosRateType

	if rateType != settings.ArgentinaDatosRateTypeBuy && rateType != settings.ArgentinaDatosRateTypeSell {
		rateType = settings.ArgentinaDatosRateTypeSell
	}

	return &ArgentinaDatosDataSource{
		exchangeHouse: strings.ToLower(exchangeHouse),
		apiRateType:   mapConfigRateTypeToArgentinaDatosApiRateType(rateType),
	}
}

func (q *ArgentinaDatosQuote) getRateByApiRateType(apiRateType string) float64 {
	if q == nil {
		return 0
	}

	if apiRateType == argentinaDatosApiRateTypeBuy {
		return q.Buy
	}

	return q.Sell
}

func (e *ArgentinaDatosDataSource) getArsPerUsd(quote *ArgentinaDatosQuote) float64 {
	return quote.getRateByApiRateType(e.apiRateType)
}

func (e *ArgentinaDatosDataSource) ToLatestExchangeRateResponse(c core.Context, content []byte) *models.LatestExchangeRateResponse {
	var quotes []ArgentinaDatosQuote
	err := json.Unmarshal(content, &quotes)

	if err != nil {
		log.Errorf(c, "[argentina_datos_datasource.ToLatestExchangeRateResponse] failed to parse json data, because %s", err.Error())
		return nil
	}

	if len(quotes) < 1 {
		log.Errorf(c, "[argentina_datos_datasource.ToLatestExchangeRateResponse] quotes is empty")
		return nil
	}

	latestQuote := quotes[len(quotes)-1]
	arsPerUsd := e.getArsPerUsd(&latestQuote)

	if arsPerUsd <= 0 {
		log.Errorf(c, "[argentina_datos_datasource.ToLatestExchangeRateResponse] rate is invalid, ars per usd is %f", arsPerUsd)
		return nil
	}

	updateTime, err := time.Parse("2006-01-02", latestQuote.Date)

	if err != nil {
		log.Errorf(c, "[argentina_datos_datasource.ToLatestExchangeRateResponse] failed to parse update date, date is %s", latestQuote.Date)
		return nil
	}

	return &models.LatestExchangeRateResponse{
		DataSource:   argentinaDatosDataSource,
		ReferenceUrl: argentinaDatosExchangeRateReferenceUrl,
		UpdateTime:   updateTime.Unix(),
		BaseCurrency: argentinaDatosBaseCurrency,
		ExchangeRates: models.LatestExchangeRateSlice{
			{
				Currency: "USD",
				Rate:     utils.Float64ToString(1 / arsPerUsd),
			},
		},
	}
}

// BuildRequests returns the ArgentinaDatos exchange rates http requests
func (e *ArgentinaDatosDataSource) BuildRequests() ([]*http.Request, error) {
	requestUrl := fmt.Sprintf(argentinaDatosExchangeRateUrlTemplate, e.exchangeHouse)
	req, err := http.NewRequest("GET", requestUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the ArgentinaDatos data source raw response
func (e *ArgentinaDatosDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	latestExchangeRateResponse := e.ToLatestExchangeRateResponse(c, content)

	if latestExchangeRateResponse == nil {
		log.Errorf(c, "[argentina_datos_datasource.Parse] failed to parse latest exchange rate data, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return latestExchangeRateResponse, nil
}
