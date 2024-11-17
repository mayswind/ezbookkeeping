package exchangerates

import (
	"net/http"
	"strings"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const internationalMonetaryFundExchangeRateUrl = "https://www.imf.org/external/np/fin/data/rms_five.aspx?tsvflag=Y"
const internationalMonetaryFundExchangeRateReferenceUrl = "https://www.imf.org/external/np/fin/data/param_rms_mth.aspx"
const internationalMonetaryFundDataSource = "International Monetary Fund"
const internationalMonetaryFundBaseCurrency = "USD"

const internationalMonetaryFundDataUpdateDateFormat = "January 02, 2006 15:04"
const internationalMonetaryFundDataUpdateDateTimezone = "America/New_York"

var internationalMonetaryFundCurrencyNameCodeMap map[string]string

// InternationalMonetaryFundDataSource defines the structure of exchange rates data source of international monetary fund
type InternationalMonetaryFundDataSource struct {
	ExchangeRatesDataSource
}

func init() {
	internationalMonetaryFundCurrencyNameCodeMap = make(map[string]string, 38)
	internationalMonetaryFundCurrencyNameCodeMap["Chinese yuan"] = "CNY"
	internationalMonetaryFundCurrencyNameCodeMap["Euro"] = "EUR"
	internationalMonetaryFundCurrencyNameCodeMap["Japanese yen"] = "JPY"
	internationalMonetaryFundCurrencyNameCodeMap["U.K. pound"] = "GBP"
	internationalMonetaryFundCurrencyNameCodeMap["U.S. dollar"] = "USD"
	internationalMonetaryFundCurrencyNameCodeMap["Algerian dinar"] = "DZD"
	internationalMonetaryFundCurrencyNameCodeMap["Australian dollar"] = "AUD"
	internationalMonetaryFundCurrencyNameCodeMap["Botswana pula"] = "BWP"
	internationalMonetaryFundCurrencyNameCodeMap["Brazilian real"] = "BRL"
	internationalMonetaryFundCurrencyNameCodeMap["Brunei dollar"] = "BND"
	internationalMonetaryFundCurrencyNameCodeMap["Canadian dollar"] = "CAD"
	internationalMonetaryFundCurrencyNameCodeMap["Chilean peso"] = "CLP"
	internationalMonetaryFundCurrencyNameCodeMap["Czech koruna"] = "CZK"
	internationalMonetaryFundCurrencyNameCodeMap["Danish krone"] = "DKK"
	internationalMonetaryFundCurrencyNameCodeMap["Indian rupee"] = "INR"
	internationalMonetaryFundCurrencyNameCodeMap["Israeli New Shekel"] = "ILS"
	internationalMonetaryFundCurrencyNameCodeMap["Korean won"] = "KRW"
	internationalMonetaryFundCurrencyNameCodeMap["Kuwaiti dinar"] = "KWD"
	internationalMonetaryFundCurrencyNameCodeMap["Malaysian ringgit"] = "MYR"
	internationalMonetaryFundCurrencyNameCodeMap["Mauritian rupee"] = "MUR"
	internationalMonetaryFundCurrencyNameCodeMap["Mexican peso"] = "MXN"
	internationalMonetaryFundCurrencyNameCodeMap["New Zealand dollar"] = "NZD"
	internationalMonetaryFundCurrencyNameCodeMap["Norwegian krone"] = "NOK"
	internationalMonetaryFundCurrencyNameCodeMap["Omani rial"] = "OMR"
	internationalMonetaryFundCurrencyNameCodeMap["Peruvian sol"] = "PEN"
	internationalMonetaryFundCurrencyNameCodeMap["Philippine peso"] = "PHP"
	internationalMonetaryFundCurrencyNameCodeMap["Polish zloty"] = "PLN"
	internationalMonetaryFundCurrencyNameCodeMap["Qatari riyal"] = "QAR"
	internationalMonetaryFundCurrencyNameCodeMap["Russian ruble"] = "RUB"
	internationalMonetaryFundCurrencyNameCodeMap["Saudi Arabian riyal"] = "SAR"
	internationalMonetaryFundCurrencyNameCodeMap["Singapore dollar"] = "SGD"
	internationalMonetaryFundCurrencyNameCodeMap["South African rand"] = "ZAR"
	internationalMonetaryFundCurrencyNameCodeMap["Swedish krona"] = "SEK"
	internationalMonetaryFundCurrencyNameCodeMap["Swiss franc"] = "CHF"
	internationalMonetaryFundCurrencyNameCodeMap["Thai baht"] = "THB"
	internationalMonetaryFundCurrencyNameCodeMap["Trinidadian dollar"] = "TTD"
	internationalMonetaryFundCurrencyNameCodeMap["U.A.E. dirham"] = "AED"
	internationalMonetaryFundCurrencyNameCodeMap["Uruguayan peso"] = "UYU"
}

// BuildRequests returns the international monetary fund exchange rates http requests
func (e *InternationalMonetaryFundDataSource) BuildRequests() ([]*http.Request, error) {
	req, err := http.NewRequest("GET", internationalMonetaryFundExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the international monetary fund data source raw response
func (e *InternationalMonetaryFundDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	lines := strings.Split(string(content), "\n")

	if len(lines) < 1 {
		log.Errorf(c, "[international_monetary_fund_datasource.Parse] content is invalid, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	exchangeRatesToSDR := orderedmap.New[string, float64]()
	latestUpdateDate := ""

	findSDRsPerCurrencyUnitLine := false
	findExchangeRateDataHeader := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if line == "" {
			continue
		}

		line = strings.ReplaceAll(line, "\r", "")

		if strings.Index(line, "Currency units per SDR") == 0 {
			break
		}

		if strings.Index(line, "SDRs per Currency unit") == 0 {
			findSDRsPerCurrencyUnitLine = true
			continue
		}

		if findExchangeRateDataHeader {
			items := strings.Split(line, "\t")

			if len(items) != 6 {
				continue
			}

			currencyCode, exchangeRate := e.parseExchangeRate(c, line, items)

			if currencyCode != nil && exchangeRate != nil {
				exchangeRatesToSDR.Set(*currencyCode, *exchangeRate)
			}

			continue
		}

		if findSDRsPerCurrencyUnitLine {
			items := strings.Split(line, "\t")

			if len(items) != 6 {
				continue
			}

			if items[0] == "Currency" {
				findExchangeRateDataHeader = true
				latestUpdateDate = items[1]
				continue
			}
		}
	}

	if latestUpdateDate == "" {
		log.Errorf(c, "[international_monetary_fund_datasource.Parse] latest update date is empty")
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if exchangeRatesToSDR.Len() < 1 {
		log.Errorf(c, "[international_monetary_fund_datasource.Parse] exchange rates date is empty")
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defaultCurrencyExchangeRateToSDR, exists := exchangeRatesToSDR.Get(internationalMonetaryFundBaseCurrency)

	if !exists {
		log.Errorf(c, "[international_monetary_fund_datasource.Parse] exchange rates date does not have default currency \"%s\"", internationalMonetaryFundBaseCurrency)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, exchangeRatesToSDR.Len())

	for pair := exchangeRatesToSDR.Oldest(); pair != nil; pair = pair.Next() {
		exchangeRates = append(exchangeRates, &models.LatestExchangeRate{
			Currency: pair.Key,
			Rate:     utils.Float64ToString(defaultCurrencyExchangeRateToSDR / pair.Value),
		})
	}

	timezone, err := time.LoadLocation(internationalMonetaryFundDataUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[international_monetary_fund_datasource.Parse] failed to get timezone, timezone name is %s", internationalMonetaryFundDataUpdateDateTimezone)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	updateDateTime := latestUpdateDate + " 11:00" // The IMF posts Representative and SDR exchange rates every 20 minutes from 11:00 AM to 6:00 PM U.S. EST Monday to Friday except for these holidays
	updateTime, err := time.ParseInLocation(internationalMonetaryFundDataUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.Errorf(c, "[international_monetary_fund_datasource.Parse] failed to parse update date, datetime is %s", updateDateTime)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    internationalMonetaryFundDataSource,
		ReferenceUrl:  internationalMonetaryFundExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  internationalMonetaryFundBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp, nil
}

func (e *InternationalMonetaryFundDataSource) parseExchangeRate(c core.Context, line string, lineItems []string) (*string, *float64) {
	currencyCode, exists := internationalMonetaryFundCurrencyNameCodeMap[lineItems[0]]

	if !exists {
		log.Warnf(c, "[international_monetary_fund_datasource.parseExchangeRate] unknown currency name %s, line is %s", lineItems[0], line)
		return nil, nil
	}

	if _, exists := validators.AllCurrencyNames[currencyCode]; !exists {
		return nil, nil
	}

	for i := 1; i < 6; i++ {
		item := lineItems[i]

		if item == "" {
			continue
		}

		rate, err := utils.StringToFloat64(item)

		if err != nil {
			log.Warnf(c, "[international_monetary_fund_datasource.parseExchangeRate] failed to parse rate, line is %s", line)
			return nil, nil
		}

		if rate <= 0 {
			log.Warnf(c, "[international_monetary_fund_datasource.parseExchangeRate] rate is invalid, line is %s", line)
			return nil, nil
		}

		return &currencyCode, &rate
	}

	log.Warnf(c, "[international_monetary_fund_datasource.parseExchangeRate] no exchange rate data exists for currency \"%s\", line is %s", currencyCode, line)
	return nil, nil
}
