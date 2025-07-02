package exchangerates

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const (
	techcombankExchangeRateUrl          = "https://techcombank.com/content/techcombank/web/vn/vi/cong-cu-tien-ich/ty-gia/_jcr_content.exchange-rates.integration.json"
	techcombankExchangeRateReferenceUrl = "https://techcombank.com/cong-cu-tien-ich/ty-gia"
	techcombankDataSource               = "Techcombank"
	techcombankBaseCurrency             = "VND"
)

type TechcombankExchangeRate struct {
	ExchangeRate struct {
		Data []struct {
			Label          string `json:"label"`
			AskRate        string `json:"askRate"`
			SourceCurrency string `json:"sourceCurrency"`
			TargetCurrency string `json:"targetCurrency"`
			InputDate      string `json:"inputDate"`
		} `json:"data"`
	} `json:"exchangeRate"`
	GoldRate struct {
		Data []struct {
			AskRate string `json:"askRate"`
		} `json:"data"`
	} `json:"goldRate"`
}

type InternationalTechcombankDataSource struct {
	HttpExchangeRatesDataSource
}

func (e *InternationalTechcombankDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", techcombankExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "") // Do not set custom user agent

	return []*http.Request{req}, nil
}

func (e *InternationalTechcombankDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	var response TechcombankExchangeRate

	// Parse JSON response
	err := json.Unmarshal(content, &response)
	if err != nil {
		log.Errorf(c, "[techcombank_datasource.Parse] failed to parse JSON response: %v", err)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(response.ExchangeRate.Data)+1) // +1 for gold

	for _, rate := range response.ExchangeRate.Data {
		if rate.AskRate == "" || rate.SourceCurrency == "" {
			continue
		}

		rateValue, err := utils.StringToFloat64(rate.AskRate)
		if err != nil {
			log.Warnf(c, "[techcombank_datasource.Parse] failed to parse rate %s for currency %s", rate.AskRate, rate.SourceCurrency)
			continue
		}

		// Calculate the inverse rate (1 VND = 1/rate foreign currency)
		inverseRate := 1.0 / rateValue

		exchangeRates = append(exchangeRates, &models.LatestExchangeRate{
			Currency: rate.SourceCurrency,
			Rate:     utils.Float64ToString(inverseRate),
		})
	}

	if len(response.GoldRate.Data) > 0 && response.GoldRate.Data[0].AskRate != "" {
		goldRateValue, err := utils.StringToFloat64(response.GoldRate.Data[0].AskRate)
		if err != nil {
			log.Warnf(c, "[techcombank_datasource.Parse] failed to parse gold rate %s", response.GoldRate.Data[0].AskRate)
			return nil, errs.ErrFailedToRequestRemoteApi
		}
		exchangeRates = append(exchangeRates, &models.LatestExchangeRate{
			Currency: "XAU",
			Rate:    utils.Float64ToString(1.0 / goldRateValue),
		})
	}

	var updateTime time.Time
	if len(response.ExchangeRate.Data) > 0 && response.ExchangeRate.Data[0].InputDate != "" {
		updateTime, err = time.Parse(time.RFC3339, response.ExchangeRate.Data[0].InputDate)
		if err != nil {
			log.Warnf(c, "[techcombank_datasource.Parse] failed to parse update time: %v, using current time", err)
			updateTime = time.Now()
		}
	} else {
		updateTime = time.Now()
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    techcombankDataSource,
		ReferenceUrl:  techcombankExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  techcombankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp, nil
}
