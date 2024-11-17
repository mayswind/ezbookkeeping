package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const bankOfIsraelMinimumRequiredContent = "" +
	"<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n" +
	"  <ExchangeRates>\n" +
	"    <ExchangeRateResponseDTO>\n" +
	"      <CurrentExchangeRate>3.733</CurrentExchangeRate>\n" +
	"      <Key>USD</Key>\n" +
	"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n" +
	"      <Unit>1</Unit>\n" +
	"    </ExchangeRateResponseDTO>\n" +
	"    <ExchangeRateResponseDTO>\n" +
	"      <CurrentExchangeRate>2.4287</CurrentExchangeRate>\n" +
	"      <Key>JPY</Key>\n" +
	"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n" +
	"      <Unit>100</Unit>\n" +
	"    </ExchangeRateResponseDTO>\n" +
	"  </ExchangeRates>\n" +
	"</ExchangeRatesResponseCollectioDTO>"

func TestBankOfIsraelDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfIsraelMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "ILS", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestBankOfIsraelDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfIsraelMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731331565), actualLatestExchangeRateResponse.UpdateTime)
}

func TestBankOfIsraelDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfIsraelMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.2678810608090008",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "41.17429077284144",
	})
}

func TestBankOfIsraelDataSource_BlankContent(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestBankOfIsraelDataSource_EmptyExchangeRatesResponseCollectioDTO(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfIsraelDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <CurrentExchangeRate>1</CurrentExchangeRate>\n"+
		"      <Key>XXX</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"      <Unit>1</Unit>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfIsraelDataSource_EmptyRate(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <Key>USD</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"      <Unit>1</Unit>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfIsraelDataSource_InvalidRate(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <CurrentExchangeRate>null</CurrentExchangeRate>\n"+
		"      <Key>USD</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"      <Unit>1</Unit>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <CurrentExchangeRate>0</CurrentExchangeRate>\n"+
		"      <Key>USD</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"      <Unit>1</Unit>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfIsraelDataSource_EmptyUnit(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <CurrentExchangeRate>1</CurrentExchangeRate>\n"+
		"      <Key>USD</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfIsraelDataSource_InvalidUnit(t *testing.T) {
	dataSource := &BankOfIsraelDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <CurrentExchangeRate>1</CurrentExchangeRate>\n"+
		"      <Key>USD</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"      <Unit>null</Unit>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<ExchangeRatesResponseCollectioDTO xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://schemas.datacontract.org/2004/07/BOI.Core.Models.HotData\">\n"+
		"  <ExchangeRates>\n"+
		"    <ExchangeRateResponseDTO>\n"+
		"      <CurrentExchangeRate>1</CurrentExchangeRate>\n"+
		"      <Key>USD</Key>\n"+
		"      <LastUpdate>2024-11-11T13:26:05.6590204Z</LastUpdate>\n"+
		"      <Unit>0</Unit>\n"+
		"    </ExchangeRateResponseDTO>\n"+
		"  </ExchangeRates>\n"+
		"</ExchangeRatesResponseCollectioDTO>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
