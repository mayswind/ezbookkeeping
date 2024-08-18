package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const euroCentralBankMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
	"<gesmes:Envelope xmlns:gesmes=\"http://www.gesmes.org/xml/2002-08-01\" xmlns=\"http://www.ecb.int/vocabulary/2002-08-01/eurofxref\">\n" +
	"  <Cube>\n" +
	"    <Cube time=\"2021-04-01\">\n" +
	"      <Cube currency=\"USD\" rate=\"1.1746\" />\n" +
	"      <Cube currency=\"CNY\" rate=\"7.7195\" />\n" +
	"    </Cube>\n" +
	"  </Cube>\n" +
	"</gesmes:Envelope>"

func TestEuroCentralBankDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(euroCentralBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "EUR", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestEuroCentralBankDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(euroCentralBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "1.1746",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "7.7195",
	})
}

func TestEuroCentralBankDataSource_BlankContent(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestEuroCentralBankDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestEuroCentralBankDataSource_EmptyEnvelopeContent(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"+
		"<gesmes:Envelope xmlns:gesmes=\"http://www.gesmes.org/xml/2002-08-01\" xmlns=\"http://www.ecb.int/vocabulary/2002-08-01/eurofxref\">"+
		"</gesmes:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestEuroCentralBankDataSource_EmptyCubeContent(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"+
		"<gesmes:Envelope xmlns:gesmes=\"http://www.gesmes.org/xml/2002-08-01\" xmlns=\"http://www.ecb.int/vocabulary/2002-08-01/eurofxref\">"+
		"<Cube>"+
		"</Cube>"+
		"</gesmes:Envelope>"))
	assert.NotEqual(t, nil, err)
}

func TestEuroCentralBankDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"+
		"<gesmes:Envelope xmlns:gesmes=\"http://www.gesmes.org/xml/2002-08-01\" xmlns=\"http://www.ecb.int/vocabulary/2002-08-01/eurofxref\">"+
		"<Cube>"+
		"<Cube time=\"2021-04-01\">"+
		"<Cube currency=\"XXX\" rate=\"1\" />"+
		"</Cube>"+
		"</Cube>"+
		"</gesmes:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestEuroCentralBankDataSource_EmptyRate(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"+
		"<gesmes:Envelope xmlns:gesmes=\"http://www.gesmes.org/xml/2002-08-01\" xmlns=\"http://www.ecb.int/vocabulary/2002-08-01/eurofxref\">"+
		"<Cube>"+
		"<Cube time=\"2021-04-01\">"+
		"<Cube currency=\"USD\" rate=\"\" />"+
		"</Cube>"+
		"</Cube>"+
		"</gesmes:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestEuroCentralBankDataSource_InvalidRate(t *testing.T) {
	dataSource := &EuroCentralBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"+
		"<gesmes:Envelope xmlns:gesmes=\"http://www.gesmes.org/xml/2002-08-01\" xmlns=\"http://www.ecb.int/vocabulary/2002-08-01/eurofxref\">"+
		"<Cube>"+
		"<Cube time=\"2021-04-01\">"+
		"<Cube currency=\"USD\" rate=\"null\" />"+
		"</Cube>"+
		"</Cube>"+
		"</gesmes:Envelope>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
