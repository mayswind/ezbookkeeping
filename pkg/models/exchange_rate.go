package models

import "encoding/xml"

const euroCentralBankDataSource = "European Central Bank"
const euroCentralBankBaseCurrency = "EUR"

// LatestExchangeRateResponse returns a view-object which contains latest exchange rate
type LatestExchangeRateResponse struct {
	DataSource    string                `json:"dataSource"`
	Date          string                `json:"date"`
	BaseCurrency  string                `json:"baseCurrency"`
	ExchangeRates []*LatestExchangeRate `json:"exchangeRates"`
}

// LatestExchangeRate represents a data pair of currency and exchange rate
type LatestExchangeRate struct {
	Currency string `json:"currency"`
	Rate     string `json:"rate"`
}

// EuroCentralBankExchangeRateData represents the whole data from euro central bank
type EuroCentralBankExchangeRateData struct {
	XMLName          xml.Name                        `xml:"Envelope"`
	AllExchangeRates []*EuroCentralBankExchangeRates `xml:"Cube>Cube"`
}

// EuroCentralBankExchangeRates represents the exchange rates data from euro central bank
type EuroCentralBankExchangeRates struct {
	Date          string                         `xml:"time,attr"`
	ExchangeRates []*EuroCentralBankExchangeRate `xml:"Cube"`
}

// EuroCentralBankExchangeRate represents the exchange rate data from euro central bank
type EuroCentralBankExchangeRate struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

// ToLatestExchangeRateResponse returns a view-object according to original data from euro central bank
func (e EuroCentralBankExchangeRateData) ToLatestExchangeRateResponse() *LatestExchangeRateResponse {
	if len(e.AllExchangeRates) < 1 {
		return nil
	}

	latestEuroCentralBankExchangeRate := e.AllExchangeRates[0]

	if len(latestEuroCentralBankExchangeRate.ExchangeRates) < 1 {
		return nil
	}

	exchangeRates := make([]*LatestExchangeRate, len(latestEuroCentralBankExchangeRate.ExchangeRates))

	for i := 0; i < len(latestEuroCentralBankExchangeRate.ExchangeRates); i++ {
		exchangeRates[i] = latestEuroCentralBankExchangeRate.ExchangeRates[i].ToLatestExchangeRate()
	}

	latestExchangeRateResp := &LatestExchangeRateResponse{
		DataSource:    euroCentralBankDataSource,
		Date:          latestEuroCentralBankExchangeRate.Date,
		BaseCurrency:  euroCentralBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

// ToLatestExchangeRate returns a data pair according to original data from euro central bank
func (e EuroCentralBankExchangeRate) ToLatestExchangeRate() *LatestExchangeRate {
	return &LatestExchangeRate{
		Currency: e.Currency,
		Rate:     e.Rate,
	}
}
