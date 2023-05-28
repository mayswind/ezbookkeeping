package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const monetaryAuthorityOfSingaporeMinimumRequiredContent = "{\n" +
	"  \"success\": true,\n" +
	"  \"result\": {\n" +
	"    \"records\": [\n" +
	"      {\n" +
	"        \"end_of_day\": \"2023-05-26\",\n" +
	"        \"usd_sgd\": \"1.3528\",\n" +
	"        \"cny_sgd_100\": \"19.16\"\n" +
	"      }\n" +
	"    ]\n" +
	"  }\n" +
	"}"

func TestMonetaryAuthorityOfSingaporeDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(monetaryAuthorityOfSingaporeMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "SGD", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestMonetaryAuthorityOfSingaporeDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(monetaryAuthorityOfSingaporeMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.7392075694855116",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "5.219206680584551",
	})
}

func TestMonetaryAuthorityOfSingaporeDataSource_BlankContent(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestMonetaryAuthorityOfSingaporeDataSource_EmptyJsonObject(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestMonetaryAuthorityOfSingaporeDataSource_ResponseSuccessIsFalseObject(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": false,\n"+
		"  \"result\": {\n"+
		"    \"records\": [\n"+
		"      {\n"+
		"        \"end_of_day\": \"2023-05-26\",\n"+
		"        \"usd_sgd\": \"1.3528\",\n"+
		"        \"cny_sgd_100\": \"19.16\"\n"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestMonetaryAuthorityOfSingaporeDataSource_NoResultContent(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": true\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestMonetaryAuthorityOfSingaporeDataSource_EmptyRecordContent(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": true,\n"+
		"  \"result\": {\n"+
		"    \"records\": [\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestMonetaryAuthorityOfSingaporeDataSource_TargetCurrencyIsNotBaseCurrency(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": true,\n"+
		"  \"result\": {\n"+
		"    \"records\": [\n"+
		"      {\n"+
		"        \"end_of_day\": \"2023-05-26\",\n"+
		"        \"usd_cny\": \"1\""+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestMonetaryAuthorityOfSingaporeDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": true,\n"+
		"  \"result\": {\n"+
		"    \"records\": [\n"+
		"      {\n"+
		"        \"end_of_day\": \"2023-05-26\",\n"+
		"        \"xxx_sgd\": \"1.3528\""+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestMonetaryAuthorityOfSingaporeDataSource_EmptyRate(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": true,\n"+
		"  \"result\": {\n"+
		"    \"records\": [\n"+
		"      {\n"+
		"        \"end_of_day\": \"2023-05-26\",\n"+
		"        \"usd_sgd\": \"\""+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestMonetaryAuthorityOfSingaporeDataSource_InvalidRate(t *testing.T) {
	dataSource := &MonetaryAuthorityOfSingaporeDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{\n"+
		"  \"success\": true,\n"+
		"  \"result\": {\n"+
		"    \"records\": [\n"+
		"      {\n"+
		"        \"end_of_day\": \"2023-05-26\",\n"+
		"        \"usd_sgd\": null"+
		"      }\n"+
		"    ]\n"+
		"  }\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
