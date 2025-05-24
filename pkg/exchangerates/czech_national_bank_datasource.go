package exchangerates

import (
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

const czechNationalBankDailyExchangeRateUrl = "https://www.cnb.cz/en/financial-markets/foreign-exchange-market/central-bank-exchange-rate-fixing/central-bank-exchange-rate-fixing/daily.txt"
const czechNationalBankMonthlyOtherExchangeRateUrl = "https://www.cnb.cz/en/financial-markets/foreign-exchange-market/fx-rates-of-other-currencies/fx-rates-of-other-currencies/fx_rates.txt"
const czechNationalBankExchangeRateReferenceUrl = "https://www.cnb.cz/en/financial-markets/foreign-exchange-market/central-bank-exchange-rate-fixing/central-bank-exchange-rate-fixing/"
const czechNationalBankDataSource = "Česká národní banka"
const czechNationalBankBaseCurrency = "CZK"

const czechNationalBankDataUpdateDateFormat = "02 Jan 2006 15:04"
const czechNationalBankDataUpdateDateTimezone = "Europe/Prague"

// CzechNationalBankDataSource defines the structure of exchange rates data source of Czech National Bank
type CzechNationalBankDataSource struct {
	HttpExchangeRatesDataSource
}

// BuildRequests returns the Czech National Bank exchange rates http requests
func (e *CzechNationalBankDataSource) BuildRequests() ([]*http.Request, error) {
	monthlyReq, err := http.NewRequest("GET", czechNationalBankMonthlyOtherExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	dailyReq, err := http.NewRequest("GET", czechNationalBankDailyExchangeRateUrl, nil)

	if err != nil {
		return nil, err
	}

	return []*http.Request{monthlyReq, dailyReq}, nil
}

// Parse returns the common response entity according to the czech nation bank data source raw response
func (e *CzechNationalBankDataSource) Parse(c core.Context, content []byte) (*models.LatestExchangeRateResponse, error) {
	lines := strings.Split(string(content), "\n")

	if len(lines) < 3 {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] content is invalid, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	headerLineItems := strings.Split(lines[0], "#")

	if len(headerLineItems) != 2 {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] first line of content is invalid, content is %s", lines[0])
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	updateDate := strings.TrimSpace(headerLineItems[0])

	titleLineItems := strings.Split(lines[1], "|")
	titleItemMap := make(map[string]int)

	for i := 0; i < len(titleLineItems); i++ {
		titleItemMap[titleLineItems[i]] = i
	}

	currencyCodeColumnIndex, exists := titleItemMap["Code"]

	if !exists {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] missing currency code column in title line, title line is %s", lines[1])
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	amountColumnIndex, exists := titleItemMap["Amount"]

	if !exists {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] missing amount column in title line, title line is %s", lines[1])
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	rateColumnIndex, exists := titleItemMap["Rate"]

	if !exists {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] missing rate column in title line, title line is %s", lines[1])
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	exchangeRates := make(models.LatestExchangeRateSlice, 0, len(lines)-2)

	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		exchangeRate := e.parseExchangeRate(c, line, currencyCodeColumnIndex, amountColumnIndex, rateColumnIndex)

		if exchangeRate != nil {
			exchangeRates = append(exchangeRates, exchangeRate)
		}
	}

	timezone, err := time.LoadLocation(czechNationalBankDataUpdateDateTimezone)

	if err != nil {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] failed to get timezone, timezone name is %s", czechNationalBankDataUpdateDateTimezone)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	updateDateTime := updateDate + " 14:30" // Exchange rates of commonly traded currencies are declared every working day after 2.30 p.m.
	updateTime, err := time.ParseInLocation(czechNationalBankDataUpdateDateFormat, updateDateTime, timezone)

	if err != nil {
		log.Errorf(c, "[czech_national_bank_datasource.Parse] failed to parse update date, datetime is %s", updateDateTime)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	latestExchangeRateResp := &models.LatestExchangeRateResponse{
		DataSource:    czechNationalBankDataSource,
		ReferenceUrl:  czechNationalBankExchangeRateReferenceUrl,
		UpdateTime:    updateTime.Unix(),
		BaseCurrency:  czechNationalBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp, nil
}

func (e *CzechNationalBankDataSource) parseExchangeRate(c core.Context, line string, currencyCodeColumnIndex int, amountColumnIndex int, rateColumnIndex int) *models.LatestExchangeRate {
	if len(line) < 1 {
		return nil
	}

	items := strings.Split(line, "|")

	if currencyCodeColumnIndex >= len(items) || amountColumnIndex >= len(items) || rateColumnIndex >= len(items) {
		log.Warnf(c, "[czech_national_bank_datasource.parseExchangeRate] missing column in data line, line is %s", line)
		return nil
	}

	currencyCode := items[currencyCodeColumnIndex]

	if _, exists := validators.AllCurrencyNames[currencyCode]; !exists {
		return nil
	}

	amount, err := utils.StringToInt64(items[amountColumnIndex])

	if err != nil {
		log.Warnf(c, "[czech_national_bank_datasource.parseExchangeRate] failed to parse amount, line is %s", line)
		return nil
	}

	if amount <= 0 {
		log.Warnf(c, "[czech_national_bank_datasource.parseExchangeRate] amount is invalid, line is %s", line)
		return nil
	}

	rate, err := utils.StringToFloat64(items[rateColumnIndex])

	if err != nil {
		log.Warnf(c, "[czech_national_bank_datasource.parseExchangeRate] failed to parse rate, line is %s", line)
		return nil
	}

	if rate <= 0 {
		log.Warnf(c, "[czech_national_bank_datasource.parseExchangeRate] rate is invalid, line is %s", line)
		return nil
	}

	finalRate := float64(amount) / rate

	if math.IsInf(finalRate, 0) {
		return nil
	}

	return &models.LatestExchangeRate{
		Currency: currencyCode,
		Rate:     utils.Float64ToString(finalRate),
	}
}
