package models

import "encoding/xml"

const EuroCentralBankDataSource = "European Central Bank"
const EuroCentralBankBaseCurrency = "EUR"

type LatestExchangeRateResponse struct {
	DataSource    string                `json:"dataSource"`
	Date          string                `json:"date"`
	BaseCurrency  string                `json:"baseCurrency"`
	ExchangeRates []*LatestExchangeRate `json:"exchangeRates"`
}

type LatestExchangeRate struct {
	Currency string `json:"currency"`
	Rate     string `json:"rate"`
}

type EuroCentralBankExchangeRateData struct {
	XMLName          xml.Name                        `xml:"Envelope"`
	AllExchangeRates []*EuroCentralBankExchangeRates `xml:"Cube>Cube"`
}

type EuroCentralBankExchangeRates struct {
	Date          string                         `xml:"time,attr"`
	ExchangeRates []*EuroCentralBankExchangeRate `xml:"Cube"`
}

type EuroCentralBankExchangeRate struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

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
		DataSource:    EuroCentralBankDataSource,
		Date:          latestEuroCentralBankExchangeRate.Date,
		BaseCurrency:  EuroCentralBankBaseCurrency,
		ExchangeRates: exchangeRates,
	}

	return latestExchangeRateResp
}

func (e EuroCentralBankExchangeRate) ToLatestExchangeRate() *LatestExchangeRate {
	return &LatestExchangeRate{
		Currency: e.Currency,
		Rate:     e.Rate,
	}
}
