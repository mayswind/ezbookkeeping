package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const bankOfRussiaDataSourceMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"windows-1251\"?>\n" +
	"<ValCurs Date=\"16.11.2024\">\n" +
	"  <Valute>\n" +
	"    <CharCode>USD</CharCode>\n" +
	"    <VunitRate>99,9971</VunitRate>\n" +
	"  </Valute>\n" +
	"  <Valute>\n" +
	"    <CharCode>CNY</CharCode>\n" +
	"    <VunitRate>13,7992</VunitRate>\n" +
	"  </Valute>\n" +
	"</ValCurs>"

func TestBankOfRussiaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfRussiaDataSourceMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "RUB", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestBankOfRussiaDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfRussiaDataSourceMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731760200), actualLatestExchangeRateResponse.UpdateTime)
}

func TestBankOfRussiaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(bankOfRussiaDataSourceMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.010000290008410243",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "0.07246796915763232",
	})
}

func TestBankOfRussiaDataSource_BlankContent(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestBankOfRussiaDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"windows-1251\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfRussiaDataSource_EmptyExchangeRatesDataset(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"windows-1251\"?>"+
		"<ValCurs Date=\"16.11.2024\">"+
		"</ValCurs>"))
	assert.NotEqual(t, nil, err)
}

func TestBankOfRussiaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"windows-1251\"?>"+
		"<ValCurs Date=\"16.11.2024\">"+
		"  <Valute>\n"+
		"    <CharCode>XXX</CharCode>\n"+
		"    <VunitRate>1</VunitRate>\n"+
		"  </Valute>\n"+
		"</ValCurs>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfRussiaDataSource_EmptyRate(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"windows-1251\"?>"+
		"<ValCurs Date=\"16.11.2024\">"+
		"  <Valute>\n"+
		"    <CharCode>USD</CharCode>\n"+
		"    <VunitRate></VunitRate>\n"+
		"  </Valute>\n"+
		"</ValCurs>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestBankOfRussiaDataSource_InvalidRate(t *testing.T) {
	dataSource := &BankOfRussiaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"windows-1251\"?>"+
		"<ValCurs Date=\"16.11.2024\">"+
		"  <Valute>\n"+
		"    <CharCode>USD</CharCode>\n"+
		"    <VunitRate>null</VunitRate>\n"+
		"  </Valute>\n"+
		"</ValCurs>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"windows-1251\"?>"+
		"<ValCurs Date=\"16.11.2024\">"+
		"  <Valute>\n"+
		"    <CharCode>USD</CharCode>\n"+
		"    <VunitRate>0</VunitRate>\n"+
		"  </Valute>\n"+
		"</ValCurs>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
