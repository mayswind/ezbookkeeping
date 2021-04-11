package exchangerates

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const nationalBankOfPolandMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n" +
	"<exchange_rates table=\"A\" date=\"2021-04-02\" number=\"064/A/NBP/2021\" uid=\"21a064\">\n" +
	"  <mid-rate currency=\"US Dollar\" units=\"1\" code=\"USD\">3.8986</mid-rate>\n" +
	"  <mid-rate currency=\"Yuan Renminbi\" units=\"1\" code=\"CNY\">0.5941</mid-rate>\n" +
	"</exchange_rates>"

func TestNationalBankOfPolandDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfPolandMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "PLN", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNationalBankOfPolandDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfPolandMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.25650233417124096",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "1.68321831341525",
	})
}

func TestNationalBankOfPolandDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_EmptyExchangeRatesContent(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n"+
		"<exchange_rates table=\"A\" date=\"2021-04-02\" number=\"064/A/NBP/2021\" uid=\"21a064\">\n"+
		"</exchange_rates>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n"+
		"<exchange_rates table=\"A\" date=\"2021-04-02\" number=\"064/A/NBP/2021\" uid=\"21a064\">\n"+
		"  <mid-rate currency=\"XXX\" units=\"1\" code=\"XXX\">1</mid-rate>\n"+
		"</exchange_rates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfPolandDataSource_EmptyRate(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n"+
		"<exchange_rates table=\"A\" date=\"2021-04-02\" number=\"064/A/NBP/2021\" uid=\"21a064\">\n"+
		"  <mid-rate currency=\"US Dollar\" units=\"1\" code=\"USD\"></mid-rate>\n"+
		"</exchange_rates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfPolandDataSource_InvalidRate(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\n"+
		"<exchange_rates table=\"A\" date=\"2021-04-02\" number=\"064/A/NBP/2021\" uid=\"21a064\">\n"+
		"  <mid-rate currency=\"US Dollar\" units=\"1\" code=\"USD\">null</mid-rate>\n"+
		"</exchange_rates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
