package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const centralBankOfHungaryDataSourceMinimumRequiredContent = "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
	"<s:Body>" +
	"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">" +
	"<GetCurrentExchangeRatesResult>" +
	"&lt;MNBCurrentExchangeRates&gt;" +
	"&lt;Day date=\"2024-11-15\"&gt;" +
	"&lt;Rate unit=\"100\" curr=\"JPY\"&gt;247,46&lt;/Rate&gt;" +
	"&lt;Rate unit=\"1\" curr=\"USD\"&gt;384,48&lt;/Rate&gt;" +
	"&lt;/Day&gt;" +
	"&lt;/MNBCurrentExchangeRates&gt;" +
	"</GetCurrentExchangeRatesResult>" +
	"</GetCurrentExchangeRatesResponse>" +
	"</s:Body>" +
	"</s:Envelope>"

func TestCentralBankOfHungaryDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfHungaryDataSourceMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "HUF", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCentralBankOfHungaryDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfHungaryDataSourceMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731664800), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCentralBankOfHungaryDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(centralBankOfHungaryDataSourceMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "0.4041057140547967",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.002600915522263837",
	})
}

func TestCentralBankOfHungaryDataSource_BlankContent(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfHungaryDataSource_MissingSoapBody(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"</s:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfHungaryDataSource_MissingGetCurrentExchangeRatesResponse(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfHungaryDataSource_MissingGetCurrentExchangeRatesResult(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfHungaryDataSource_EmptyGetCurrentExchangeRatesResult(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfHungaryDataSource_InvalidGetCurrentExchangeRatesResult(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestCentralBankOfHungaryDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"&lt;Rate unit=\"1\" curr=\"XXX\"&gt;1&lt;/Rate&gt;"+
		"&lt;/Day&gt;"+
		"&lt;/MNBCurrentExchangeRates&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfHungaryDataSource_EmptyRate(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"&lt;Rate unit=\"1\" curr=\"USD\"&gt;&lt;/Rate&gt;"+
		"&lt;/Day&gt;"+
		"&lt;/MNBCurrentExchangeRates&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfHungaryDataSource_InvalidRate(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"&lt;Rate unit=\"1\" curr=\"USD\"&gt;null&lt;/Rate&gt;"+
		"&lt;/Day&gt;"+
		"&lt;/MNBCurrentExchangeRates&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"&lt;Rate unit=\"1\" curr=\"USD\"&gt;0&lt;/Rate&gt;"+
		"&lt;/Day&gt;"+
		"&lt;/MNBCurrentExchangeRates&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCentralBankOfHungaryDataSource_InvalidUnit(t *testing.T) {
	dataSource := &CentralBankOfHungaryDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"&lt;Rate unit=\"null\" curr=\"USD\"&gt;384,48&lt;/Rate&gt;"+
		"&lt;/Day&gt;"+
		"&lt;/MNBCurrentExchangeRates&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">"+
		"<s:Body>"+
		"<GetCurrentExchangeRatesResponse xmlns=\"http://www.mnb.hu/webservices/\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">"+
		"<GetCurrentExchangeRatesResult>"+
		"&lt;MNBCurrentExchangeRates&gt;"+
		"&lt;Day date=\"2024-11-15\"&gt;"+
		"&lt;Rate unit=\"0\" curr=\"USD\"&gt;384,48&lt;/Rate&gt;"+
		"&lt;/Day&gt;"+
		"&lt;/MNBCurrentExchangeRates&gt;"+
		"</GetCurrentExchangeRatesResult>"+
		"</GetCurrentExchangeRatesResponse>"+
		"</s:Body>"+
		"</s:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
