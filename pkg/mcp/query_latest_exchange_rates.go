package mcp

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/exchangerates"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// MCPQueryExchangeRatesRequest represents all parameters of the query exchange rates request
type MCPQueryExchangeRatesRequest struct {
	Currencies string `json:"currencies" jsonschema_description:"Comma-separated list of currencies to query exchange rates for (e.g. USD,CNY,EUR)"`
}

// MCPQueryExchangeRatesResponse represents the response structure for querying exchange rates
type MCPQueryExchangeRatesResponse struct {
	BaseCurrency string                      `json:"base_currency" jsonschema_description:"Base currency code (e.g. USD)"`
	UpdateTime   string                      `json:"update_time" jsonschema_description:"Last update time of the exchange rates in RFC 3339 format (e.g. '2023-01-01T12:00:00Z')"`
	Rates        []*MCPQueryExchangeRateInfo `json:"rates" jsonschema_description:"Exchange rates for the specified currencies"`
}

// MCPQueryExchangeRateInfo defines the structure of exchange rate information for a specific currency
type MCPQueryExchangeRateInfo struct {
	Currency string `json:"currency" jsonschema_description:"Currency code (e.g. USD)"`
	Rate     string `json:"rate" jsonschema_description:"The amount of the base currency that can be exchanged for 1 of this currency"`
}

type mcpQueryLatestExchangeRatesToolHandler struct{}

var MCPQueryLatestExchangeRatesToolHandler = &mcpQueryLatestExchangeRatesToolHandler{}

// Name returns the name of the MCP tool
func (h *mcpQueryLatestExchangeRatesToolHandler) Name() string {
	return "query_latest_exchange_rates"
}

// Description returns the description of the MCP tool
func (h *mcpQueryLatestExchangeRatesToolHandler) Description() string {
	return "Query latest exchange rates with specified currencies."
}

// InputType returns the input type for the MCP tool request
func (h *mcpQueryLatestExchangeRatesToolHandler) InputType() reflect.Type {
	return reflect.TypeOf(&MCPQueryExchangeRatesRequest{})
}

// OutputType returns the output type for the MCP tool response
func (h *mcpQueryLatestExchangeRatesToolHandler) OutputType() reflect.Type {
	return reflect.TypeOf(&MCPQueryExchangeRatesResponse{})
}

// Handle processes the MCP call tool request and returns the response
func (h *mcpQueryLatestExchangeRatesToolHandler) Handle(c *core.WebContext, callToolReq *MCPCallToolRequest, currentConfig *settings.Config, services MCPAvailableServices) ([]*MCPTextContent, *errs.Error) {
	var exchangeRatesRequest MCPQueryExchangeRatesRequest

	if callToolReq.Arguments != nil {
		if err := json.Unmarshal(callToolReq.Arguments, &exchangeRatesRequest); err != nil {
			return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	dataSource := exchangerates.Container.Current

	if dataSource == nil {
		return nil, errs.ErrInvalidExchangeRatesDataSource
	}

	exchangeRateResponse, err := dataSource.GetLatestExchangeRates(c, c.GetCurrentUid(), currentConfig)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	response, err := h.createNewMCPQueryExchangeRatesResponse(exchangeRatesRequest.Currencies, exchangeRateResponse)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return response, nil
}

func (h *mcpQueryLatestExchangeRatesToolHandler) createNewMCPQueryExchangeRatesResponse(currencies string, exchangeRatesResp *models.LatestExchangeRateResponse) ([]*MCPTextContent, error) {
	queryCurrencies := make(map[string]bool)

	for _, currency := range strings.Split(currencies, ",") {
		currency = strings.TrimSpace(currency)

		if currency != "" {
			queryCurrencies[currency] = true
		}
	}

	response := &MCPQueryExchangeRatesResponse{
		BaseCurrency: exchangeRatesResp.BaseCurrency,
		UpdateTime:   utils.FormatUnixTimeToLongDateTimeWithTimezoneRFC3389Format(exchangeRatesResp.UpdateTime, time.UTC),
		Rates:        make([]*MCPQueryExchangeRateInfo, 0, len(exchangeRatesResp.ExchangeRates)),
	}

	for _, rate := range exchangeRatesResp.ExchangeRates {
		if _, exists := queryCurrencies[rate.Currency]; rate.Currency != exchangeRatesResp.BaseCurrency && !exists {
			continue
		}

		response.Rates = append(response.Rates, &MCPQueryExchangeRateInfo{
			Currency: rate.Currency,
			Rate:     rate.Rate,
		})
	}

	content, err := json.Marshal(response)

	if err != nil {
		return nil, err
	}

	return []*MCPTextContent{
		NewMCPTextContent(string(content)),
	}, nil
}
