package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const centralBankOfMyanmarMinimumRequiredContent = "{\n" +
	"  \"timestamp\": \"1731571200\",\n" +
	"  \"rates\": {\n" +
	"    \"USD\": \"2,100.0\",\n" +
	"    \"JPY\": \"1,347.6\"\n" +
	"  }\n" +
	"}"

func TestCentralBankOfMyanmarDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfMyanmarMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "MMK", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCentralBankOfMyanmarDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfMyanmarMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731571200), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCentralBankOfMyanmarDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfMyanmarMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "0.07420599584446423",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.0004761904761904762",
	})
}

func TestCentralBankOfMyanmarDataSource_BlankContent(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfMyanmarDataSource_EmptyData(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfMyanmarDataSource_EmptyExchangeRatesData(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"timestamp\": \"1731571200\"\n"+
		"}"))

	_, err = dataSource.Parse(context, []byte("{\n"+
		"  \"timestamp\": \"1731571200\",\n"+
		"  \"rates\": {\n"+
		"  }\n"+
		"}"))
	assert.Nil(t, nil, err)
}

func TestCentralBankOfMyanmarDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"timestamp\": \"1731571200\",\n"+
		"  \"rates\": {\n"+
		"    \"XXX\": \"1\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfMyanmarDataSource_InvalidRate(t *testing.T) {
	dataSource := &CentralBankOfMyanmarDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"timestamp\": \"1731571200\",\n"+
		"  \"rates\": {\n"+
		"    \"USD\": null\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("{\n"+
		"  \"timestamp\": \"1731571200\",\n"+
		"  \"rates\": {\n"+
		"    \"USD\": \"0\"\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
