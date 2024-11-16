package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const nationalBankOfPolandMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
	"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n" +
	"  <ExchangeRatesTable>\n" +
	"    <EffectiveDate>2024-02-28</EffectiveDate>\n" +
	"    <Rates>\n" +
	"      <Rate>\n" +
	"        <Code>USD</Code>\n" +
	"        <Mid>3.9922</Mid>\n" +
	"      </Rate>\n" +
	"      <Rate>\n" +
	"        <Code>CNY</Code>\n" +
	"        <Mid>0.5545</Mid>\n" +
	"      </Rate>\n" +
	"    </Rates>\n" +
	"  </ExchangeRatesTable>\n" +
	"</ArrayOfExchangeRatesTable>"

func TestNationalBankOfPolandDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfPolandMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "PLN", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNationalBankOfPolandDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfPolandMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1709118900), actualLatestExchangeRateResponse.UpdateTime)
}

func TestNationalBankOfPolandDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfPolandMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.2504884524823406",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "1.8034265103697025",
	})
}

func TestNationalBankOfPolandDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_EmptyArrayOfExchangeRatesTable(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_EmptyExchangeRatesTable(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"  <ExchangeRatesTable>\n"+
		"  </ExchangeRatesTable>\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_EmptyExchangeRatesContent(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"  <ExchangeRatesTable>\n"+
		"    <EffectiveDate>2024-02-28</EffectiveDate>\n"+
		"    <Rates>\n"+
		"    </Rates>\n"+
		"  </ExchangeRatesTable>\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfPolandDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"  <ExchangeRatesTable>\n"+
		"    <EffectiveDate>2024-02-28</EffectiveDate>\n"+
		"    <Rates>\n"+
		"      <Rate>\n"+
		"        <Code>XXX</Code>\n"+
		"        <Mid>1</Mid>\n"+
		"      </Rate>\n"+
		"    </Rates>\n"+
		"  </ExchangeRatesTable>\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfPolandDataSource_EmptyRate(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"  <ExchangeRatesTable>\n"+
		"    <EffectiveDate>2024-02-28</EffectiveDate>\n"+
		"    <Rates>\n"+
		"      <Rate>\n"+
		"        <Code>USD</Code>\n"+
		"        <Mid></Mid>\n"+
		"      </Rate>\n"+
		"    </Rates>\n"+
		"  </ExchangeRatesTable>\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfPolandDataSource_InvalidRate(t *testing.T) {
	dataSource := &NationalBankOfPolandDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"  <ExchangeRatesTable>\n"+
		"    <EffectiveDate>2024-02-28</EffectiveDate>\n"+
		"    <Rates>\n"+
		"      <Rate>\n"+
		"        <Code>USD</Code>\n"+
		"        <Mid>null</Mid>\n"+
		"      </Rate>\n"+
		"    </Rates>\n"+
		"  </ExchangeRatesTable>\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<ArrayOfExchangeRatesTable xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n"+
		"  <ExchangeRatesTable>\n"+
		"    <EffectiveDate>2024-02-28</EffectiveDate>\n"+
		"    <Rates>\n"+
		"      <Rate>\n"+
		"        <Code>USD</Code>\n"+
		"        <Mid>0</Mid>\n"+
		"      </Rate>\n"+
		"    </Rates>\n"+
		"  </ExchangeRatesTable>\n"+
		"</ArrayOfExchangeRatesTable>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
