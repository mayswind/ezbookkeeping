package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const nationalBankOfRomaniaMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
	"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">\n" +
	"  <Header>\n" +
	"    <PublishingDate>2024-11-15</PublishingDate>\n" +
	"  </Header>\n" +
	"  <Body>\n" +
	"    <OrigCurrency>RON</OrigCurrency>\n" +
	"    <Cube date=\"2024-11-15\">\n" +
	"      <Rate currency=\"JPY\" multiplier=\"100\">3.0303</Rate>\n" +
	"      <Rate currency=\"USD\">4.7057</Rate>\n" +
	"    </Cube>\n" +
	"  </Body>\n" +
	"</DataSet>"

func TestNationalBankOfRomaniaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfRomaniaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "RON", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNationalBankOfRomaniaDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfRomaniaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731668400), actualLatestExchangeRateResponse.UpdateTime)
}

func TestNationalBankOfRomaniaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfRomaniaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "33.000033000033",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.21250823469409438",
	})
}

func TestNationalBankOfRomaniaDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfRomaniaDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfRomaniaDataSource_EmptyExchangeRatesDataset(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"</DataSet>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfRomaniaDataSource_NoDailyRatesHeader(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"</DataSet>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfRomaniaDataSource_NoDailyRatesBody(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"</DataSet>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfRomaniaDataSource_NoDailyRatesCube(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfRomaniaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"    <Cube date=\"2024-11-15\">\n"+
		"      <Rate currency=\"XXX\">1</Rate>\n"+
		"    </Cube>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfRomaniaDataSource_EmptyRate(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"    <Cube date=\"2024-11-15\">\n"+
		"      <Rate currency=\"USD\"></Rate>\n"+
		"    </Cube>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfRomaniaDataSource_InvalidRate(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"    <Cube date=\"2024-11-15\">\n"+
		"      <Rate currency=\"USD\">null</Rate>\n"+
		"    </Cube>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"    <Cube date=\"2024-11-15\">\n"+
		"      <Rate currency=\"USD\">0</Rate>\n"+
		"    </Cube>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNationalBankOfRomaniaDataSource_InvalidMultiplier(t *testing.T) {
	dataSource := &NationalBankOfRomaniaDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"    <Cube date=\"2024-11-15\">\n"+
		"      <Rate currency=\"JPY\" multiplier=\"null\">3.0303</Rate>\n"+
		"    </Cube>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<DataSet xmlns=\"http://www.bnr.ro/xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.bnr.ro/xsd nbrfxrates.xsd\">"+
		"  <Header>\n"+
		"    <PublishingDate>2024-11-15</PublishingDate>\n"+
		"  </Header>\n"+
		"  <Body>\n"+
		"    <OrigCurrency>RON</OrigCurrency>\n"+
		"    <Cube date=\"2024-11-15\">\n"+
		"      <Rate currency=\"JPY\" multiplier=\"0\">3.0303</Rate>\n"+
		"    </Cube>\n"+
		"  </Body>\n"+
		"</DataSet>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
