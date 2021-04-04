package exchangerates

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/models"
)

const bankOfCanadaMinimumRequiredContent = "{\n" +
	"    \"observations\": [\n" +
	"        {\n" +
	"            \"d\": \"2019-12-31\",\n" +
	"            \"FXVNDCAD\": {\n" +
	"                \"v\": \"0.000056\"\n" +
	"            }\n" +
	"        },\n" +
	"        {\n" +
	"            \"d\": \"2021-04-01\",\n" +
	"            \"FXCNYCAD\": {\n" +
	"                \"v\": \"0.1913\"\n" +
	"            },\n" +
	"            \"FXUSDCAD\": {\n" +
	"                \"v\": \"1.2565\"\n" +
	"            }\n" +
	"        }\n" +
	"    ]\n" +
	"}"

func TestBankOfCanadaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfCanadaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "CAD", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestBankOfCanadaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfCanadaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.7958615200955034",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "5.2273915316257185",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "VND",
		Rate:     "17857.14285714286",
	})
}

func TestBankOfCanadaDataSource_BlankContent(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestBankOfCanadaDataSource_EmptyJsonObject(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfCanadaDataSource_EmptyObservationsContent(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("{"+
		"    \"observations\": []"+
		"}"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfCanadaDataSource_InvalidObservationFormat(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{"+
		"    \"observations\": [\n"+
		"        {\n"+
		"            \"d\": \"2021-04-01\",\n"+
		"            \"CNYCAD\": {\n"+
		"                \"v\": \"0.1913\"\n"+
		"            }\n"+
		"        }\n"+
		"    ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfCanadaDataSource_InvalidObservationFormat2(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{"+
		"    \"observations\": [\n"+
		"        {\n"+
		"            \"d\": \"2021-04-01\",\n"+
		"            \"FXCADCNY\": {\n"+
		"                \"v\": \"0.1913\"\n"+
		"            }\n"+
		"        }\n"+
		"    ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfCanadaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{"+
		"    \"observations\": [\n"+
		"        {\n"+
		"            \"d\": \"2021-04-01\",\n"+
		"            \"FXXXXCAD\": {\n"+
		"                \"v\": \"1\"\n"+
		"            }\n"+
		"        }\n"+
		"    ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfCanadaDataSource_EmptyRate(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{"+
		"    \"observations\": [\n"+
		"        {\n"+
		"            \"d\": \"2021-04-01\",\n"+
		"            \"FXUSDCAD\": {\n"+
		"                \"v\": \"\"\n"+
		"            }\n"+
		"        }\n"+
		"    ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfCanadaDataSource_InvalidRate(t *testing.T) {
	dataSource := &BankOfCanadaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("{"+
		"    \"observations\": [\n"+
		"        {\n"+
		"            \"d\": \"2021-04-01\",\n"+
		"            \"FXUSDCAD\": {\n"+
		"                \"v\": null\n"+
		"            }\n"+
		"        }\n"+
		"    ]\n"+
		"}"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
