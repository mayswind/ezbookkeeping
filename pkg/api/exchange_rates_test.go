package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/exchangerates"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestExchangeRatesApiLatestExchangeRateHandler_ReserveBankOfAustraliaDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.ReserveBankOfAustraliaDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "AUD", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"CAD", "CHF", "CNY", "EUR", "GBP", "HKD", "IDR", "INR", "JPY", "KRW",
		"MYR", "NZD", "PHP", "SGD", "THB", "TWD", "USD", "VND"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_BankOfCanadaDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.BankOfCanadaDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "CAD", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "BRL", "CHF", "CNY", "EUR", "GBP", "HKD", "IDR", "INR",
		"JPY", "KRW", "MXN", "MYR", "NOK", "NZD", "PEN", "RUB", "SAR", "SEK", "SGD", "THB", "TRY", "TWD",
		"USD", "VND", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_CzechNationalBankDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.CzechNationalBankDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "CZK", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AFN", "ALL", "AMD", "AOA", "ARS", "AUD", "AWG", "AZN",
		"BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BZD",
		"CAD", "CDF", "CHF", "CLP", "CNY", "COP", "CRC", "CUP", "CVE", "DJF", "DKK", "DOP", "DZD",
		"EGP", "ERN", "ETB", "EUR", "FJD", "FKP", "GBP", "GEL", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD",
		"HKD", "HNL", "HTG", "HUF", "IDR", "ILS", "INR", "IQD", "IRR", "ISK", "JMD", "JOD", "JPY",
		"KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LYD",
		"MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRU", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN",
		"NAD", "NGN", "NIO", "NOK", "NPR", "NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG",
		"QAR", "RON", "RSD", "RUB", "RWF",
		"SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLE", "SOS", "SRD", "SSP", "STN", "SVC", "SYP", "SZL",
		"THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD", "UYU", "UZS",
		"VND", "VUV", "WST", "XAF", "XCD", "XOF", "XPF", "YER", "ZAR", "ZMW"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_DanmarksNationalbankDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.DanmarksNationalbankDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "DKK", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "EUR", "GBP", "HKD", "HUF",
		"IDR", "ILS", "INR", "ISK", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD", "PHP", "PLN", "RON", "SEK", "SGD",
		"THB", "TRY", "USD", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_EuroCentralBankDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.EuroCentralBankDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "EUR", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "GBP", "HKD", "HUF",
		"IDR", "ILS", "INR", "ISK", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD", "PHP", "PLN", "RON", "SEK", "SGD",
		"THB", "TRY", "USD", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_NationalBankOfGeorgiaDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.NationalBankOfGeorgiaDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "GEL", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AMD", "AUD", "AZN", "BGN", "BRL", "BYN", "CAD", "CHF", "CNY", "CZK",
		"DKK", "EGP", "EUR", "GBP", "HKD", "HUF", "ILS", "INR", "IRR", "ISK", "JPY", "KGS", "KRW", "KWD", "KZT",
		"MDL", "NOK", "NZD", "PLN", "QAR", "RON", "RSD", "RUB", "SEK", "SGD", "TJS", "TMT", "TRY",
		"UAH", "USD", "UZS", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_CentralBankOfHungaryDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.CentralBankOfHungaryDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "HUF", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "EUR",
		"GBP", "HKD", "IDR", "ILS", "INR", "ISK", "JPY", "KRW", "MXN", "MYR", "NOK", "NZD",
		"PHP", "PLN", "RON", "RSD", "RUB", "SEK", "SGD", "THB", "TRY", "UAH", "USD", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_BankOfIsraelDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.BankOfIsraelDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "ILS", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "CAD", "CHF", "DKK", "EGP", "EUR", "GBP",
		"JOD", "JPY", "LBP", "NOK", "SEK", "USD", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_CentralBankOfMyanmarDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.CentralBankOfMyanmarDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "MMK", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "BDT", "BND", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK",
		"EGP", "EUR", "GBP", "HKD", "IDR", "ILS", "INR", "JPY", "KES", "KHR", "KRW", "KWD", "LAK", "LKR",
		"MYR", "NOK", "NPR", "NZD", "PHP", "PKR", "RSD", "RUB", "SAR", "SEK", "SGD", "THB", "USD", "VND", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_NorgesBankDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.NorgesBankDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "NOK", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AUD", "BDT", "BGN", "BRL", "BYN", "CAD", "CHF", "CNY", "CZK", "DKK",
		"EUR", "GBP", "HKD", "HUF", "IDR", "ILS", "INR", "ISK", "JPY", "KRW", "MMK", "MXN", "MYR", "NZD",
		"PHP", "PKR", "PLN", "RON", "RUB", "SEK", "SGD", "THB", "TRY", "TWD", "USD", "VND", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_NationalBankOfPolandDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.NationalBankOfPolandDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "PLN", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AFN", "ALL", "AMD", "AOA", "ARS", "AUD", "AWG", "AZN",
		"BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BND", "BOB", "BRL", "BSD", "BWP", "BYN", "BZD",
		"CAD", "CDF", "CHF", "CLP", "CNY", "COP", "CRC", "CUP", "CVE", "CZK",
		"DJF", "DKK", "DOP", "DZD", "EGP", "ERN", "ETB", "EUR", "FJD",
		"GBP", "GEL", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD",
		"HKD", "HNL", "HTG", "HUF", "IDR", "ILS", "INR", "IQD", "IRR", "ISK",
		"JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KRW", "KWD", "KZT",
		"LAK", "LBP", "LKR", "LRD", "LSL", "LYD",
		"MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRU", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN",
		"NAD", "NGN", "NIO", "NOK", "NPR", "NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PYG",
		"QAR", "RON", "RSD", "RUB", "RWF",
		"SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SLE", "SOS", "SRD", "SSP", "STN", "SVC", "SYP", "SZL",
		"THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD", "UYU", "UZS",
		"VES", "VND", "VUV", "WST", "XAF", "XCD", "XOF", "XPF", "YER", "ZAR", "ZMW", "ZWG"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_NationalBankOfRomaniaDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.NationalBankOfRomaniaDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "RON", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "EGP",
		"EUR", "GBP", "HKD", "HUF", "IDR", "ILS", "INR", "ISK", "JPY", "KRW", "MDL", "MXN", "MYR",
		"NOK", "NZD", "PHP", "PLN", "RSD", "RUB", "SEK", "SGD", "THB", "TRY", "UAH", "USD", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_BankOfRussiaDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.BankOfRussiaDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "RUB", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AMD", "AUD", "AZN", "BGN", "BRL", "BYN", "CAD", "CHF", "CNY", "CZK",
		"DKK", "EGP", "EUR", "GBP", "GEL", "HKD", "HUF", "IDR", "INR", "JPY", "KGS", "KRW", "KZT", "MDL",
		"NOK", "NZD", "PLN", "QAR", "RON", "RSD", "SEK", "SGD", "THB", "TJS", "TMT", "TRY",
		"UAH", "USD", "UZS", "VND", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_SwissNationalBankDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.SwissNationalBankDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "CHF", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"EUR", "GBP", "JPY", "USD"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_NationalBankOfUkraineDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.NationalBankOfUkraineDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "UAH", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{
		"AED", "AUD", "AZN", "BDT", "BGN", "CAD", "CHF", "CNY", "CZK", "DKK",
		"DZD", "EGP", "EUR", "GBP", "GEL", "HKD", "HUF", "IDR", "ILS", "INR",
		"JPY", "KRW", "KZT", "LBP", "MDL", "MXN", "MYR", "NOK", "NZD", "PLN",
		"RON", "RSD", "SAR", "SEK", "SGD", "THB", "TND", "TRY", "USD", "VND", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_CentralBankOfUzbekistanDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.CentralBankOfUzbekistanDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "UZS", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AFN", "AMD", "ARS", "AUD", "AZN",
		"BDT", "BGN", "BHD", "BND", "BRL", "BYN", "CAD", "CHF", "CNY", "CUP", "CZK",
		"DKK", "DZD", "EGP", "EUR", "GBP", "GEL", "HKD", "HUF", "IDR", "ILS", "INR", "IQD", "IRR", "ISK",
		"JOD", "JPY", "KGS", "KHR", "KRW", "KWD", "KZT", "LAK", "LBP", "LYD",
		"MAD", "MDL", "MMK", "MNT", "MXN", "MYR", "NOK", "NZD", "OMR", "PHP", "PKR", "PLN",
		"QAR", "RON", "RSD", "RUB", "SAR", "SDG", "SEK", "SGD", "SYP",
		"THB", "TJS", "TMT", "TND", "TRY", "UAH", "USD", "UYU", "VES", "VND", "YER", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func TestExchangeRatesApiLatestExchangeRateHandler_InternationalMonetaryFundDataSource(t *testing.T) {
	exchangeRateResponse := executeLatestExchangeRateHandler(t, settings.InternationalMonetaryFundDataSource)

	if exchangeRateResponse == nil {
		return
	}

	assert.Equal(t, "USD", exchangeRateResponse.BaseCurrency)

	supportedCurrencyCodes := []string{"AED", "AUD", "BND", "BRL", "BWP", "CAD", "CHF", "CLP", "CNY", "CZK",
		"DKK", "DZD", "EUR", "GBP", "ILS", "INR", "JPY", "KRW", "KWD", "MUR", "MXN", "MYR", "NOK", "NZD",
		"OMR", "PEN", "PHP", "PLN", "QAR", "RUB", "SAR", "SEK", "SGD", "THB", "TTD", "UYU", "ZAR"}

	checkExchangeRatesHaveSpecifiedCurrencies(t, exchangeRateResponse.BaseCurrency, supportedCurrencyCodes, exchangeRateResponse.ExchangeRates)
}

func executeLatestExchangeRateHandler(t *testing.T, dataSourceType string) *models.LatestExchangeRateResponse {
	config := &settings.Config{
		ExchangeRatesDataSource:     dataSourceType,
		ExchangeRatesRequestTimeout: 10000,
		ExchangeRatesProxy:          "system",
		ExchangeRatesSkipTLSVerify:  true,
	}

	settingsContainer := &settings.ConfigContainer{
		Current: config,
	}

	err := exchangerates.InitializeExchangeRatesDataSource(config)
	assert.Nil(t, err)

	exchangeRatesApi := &ExchangeRatesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settingsContainer,
		},
	}

	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	response, err := exchangeRatesApi.LatestExchangeRateHandler(&core.WebContext{
		Context: ginContext,
	})
	assert.Nil(t, err)

	exchangeRateResponse := response.(*models.LatestExchangeRateResponse)
	assert.NotNil(t, exchangeRateResponse)

	return exchangeRateResponse
}

func checkExchangeRatesHaveSpecifiedCurrencies(t *testing.T, baseCurrency string, currencyCodes []string, exchangeRates []*models.LatestExchangeRate) {
	assert.Equal(t, len(currencyCodes)+1, len(exchangeRates))

	currencyCodesInExchangeRates := make(map[string]*models.LatestExchangeRate, len(exchangeRates))

	for i := 0; i < len(exchangeRates); i++ {
		exchangeRate := exchangeRates[i]
		currencyCodesInExchangeRates[exchangeRate.Currency] = exchangeRate
	}

	allCurrencyCodes := append(currencyCodes, baseCurrency)

	for i := 0; i < len(allCurrencyCodes); i++ {
		exchangeRate, has := currencyCodesInExchangeRates[allCurrencyCodes[i]]
		assert.True(t, has, allCurrencyCodes[i])

		if has {
			rate, err := utils.StringToFloat64(exchangeRate.Rate)
			assert.Nil(t, err)
			assert.Greater(t, rate, float64(0))
		}
	}
}
