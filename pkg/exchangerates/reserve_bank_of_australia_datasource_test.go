package exchangerates

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const reserveBankOfAustraliaMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
	"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n" +
	"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n" +
	"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n" +
	"  </channel>\n" +
	"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n" +
	"    <cb:statistics rdf:parseType=\"Resource\">\n" +
	"      <cb:exchangeRate rdf:parseType=\"Resource\">\n" +
	"        <cb:observation rdf:parseType=\"Resource\">\n" +
	"          <cb:value>0.7543</cb:value>\n" +
	"          <cb:unit>AUD</cb:unit>\n" +
	"        </cb:observation>\n" +
	"        <cb:baseCurrency>AUD</cb:baseCurrency>\n" +
	"        <cb:targetCurrency>USD</cb:targetCurrency>\n" +
	"      </cb:exchangeRate>\n" +
	"    </cb:statistics>\n" +
	"  </item>\n" +
	"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#CNY\">\n" +
	"    <cb:statistics rdf:parseType=\"Resource\">\n" +
	"      <cb:exchangeRate rdf:parseType=\"Resource\">\n" +
	"        <cb:observation rdf:parseType=\"Resource\">\n" +
	"          <cb:value>4.9577</cb:value>\n" +
	"          <cb:unit>AUD</cb:unit>\n" +
	"        </cb:observation>\n" +
	"        <cb:baseCurrency>AUD</cb:baseCurrency>\n" +
	"        <cb:targetCurrency>CNY</cb:targetCurrency>\n" +
	"      </cb:exchangeRate>\n" +
	"    </cb:statistics>\n" +
	"  </item>\n" +
	"</rdf:RDF>"

func TestReserveBankOfAustraliaDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(reserveBankOfAustraliaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "AUD", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestReserveBankOfAustraliaDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(reserveBankOfAustraliaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.7543",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "4.9577",
	})
}

func TestReserveBankOfAustraliaDataSource_BlankContent(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestReserveBankOfAustraliaDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestReserveBankOfAustraliaDataSource_EmptyRDFContent(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"</rdf:RDF>"))
	assert.NotEqual(t, nil, err)
}

func TestReserveBankOfAustraliaDataSource_EmptyChannelContent(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"  </channel>"+
		"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n"+
		"    <cb:statistics rdf:parseType=\"Resource\">\n"+
		"      <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"        <cb:observation rdf:parseType=\"Resource\">\n"+
		"          <cb:value>0.7543</cb:value>\n"+
		"          <cb:unit>AUD</cb:unit>\n"+
		"        </cb:observation>\n"+
		"        <cb:baseCurrency>AUD</cb:baseCurrency>\n"+
		"        <cb:targetCurrency>USD</cb:targetCurrency>\n"+
		"      </cb:exchangeRate>\n"+
		"    </cb:statistics>\n"+
		"  </item>\n"+
		"</rdf:RDF>"))
	assert.NotEqual(t, nil, err)
}

func TestReserveBankOfAustraliaDataSource_NoItem(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n"+
		"  </channel>\n"+
		"</rdf:RDF>"))
	assert.NotEqual(t, nil, err)
}

func TestReserveBankOfAustraliaDataSource_BaseCurrencyNotEqualPreset(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n"+
		"  </channel>\n"+
		"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n"+
		"    <cb:statistics rdf:parseType=\"Resource\">\n"+
		"      <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"        <cb:observation rdf:parseType=\"Resource\">\n"+
		"          <cb:value>0.7543</cb:value>\n"+
		"          <cb:unit>AUD</cb:unit>\n"+
		"        </cb:observation>\n"+
		"        <cb:baseCurrency>USD</cb:baseCurrency>\n"+
		"        <cb:targetCurrency>USD</cb:targetCurrency>\n"+
		"      </cb:exchangeRate>\n"+
		"    </cb:statistics>\n"+
		"  </item>\n"+
		"</rdf:RDF>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestReserveBankOfAustraliaDataSource_UnitCurrencyNotEqualPreset(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n"+
		"  </channel>\n"+
		"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n"+
		"    <cb:statistics rdf:parseType=\"Resource\">\n"+
		"      <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"        <cb:observation rdf:parseType=\"Resource\">\n"+
		"          <cb:value>0.7543</cb:value>\n"+
		"          <cb:unit>USD</cb:unit>\n"+
		"        </cb:observation>\n"+
		"        <cb:baseCurrency>AUD</cb:baseCurrency>\n"+
		"        <cb:targetCurrency>USD</cb:targetCurrency>\n"+
		"      </cb:exchangeRate>\n"+
		"    </cb:statistics>\n"+
		"  </item>\n"+
		"</rdf:RDF>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestReserveBankOfAustraliaDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n"+
		"  </channel>\n"+
		"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n"+
		"    <cb:statistics rdf:parseType=\"Resource\">\n"+
		"      <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"        <cb:observation rdf:parseType=\"Resource\">\n"+
		"          <cb:value>1</cb:value>\n"+
		"          <cb:unit>AUD</cb:unit>\n"+
		"        </cb:observation>\n"+
		"        <cb:baseCurrency>AUD</cb:baseCurrency>\n"+
		"        <cb:targetCurrency>XXX</cb:targetCurrency>\n"+
		"      </cb:exchangeRate>\n"+
		"    </cb:statistics>\n"+
		"  </item>\n"+
		"</rdf:RDF>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestReserveBankOfAustraliaDataSource_EmptyRate(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n"+
		"  </channel>\n"+
		"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n"+
		"    <cb:statistics rdf:parseType=\"Resource\">\n"+
		"      <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"        <cb:observation rdf:parseType=\"Resource\">\n"+
		"          <cb:value></cb:value>\n"+
		"          <cb:unit>AUD</cb:unit>\n"+
		"        </cb:observation>\n"+
		"        <cb:baseCurrency>AUD</cb:baseCurrency>\n"+
		"        <cb:targetCurrency>USD</cb:targetCurrency>\n"+
		"      </cb:exchangeRate>\n"+
		"    </cb:statistics>\n"+
		"  </item>\n"+
		"</rdf:RDF>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestReserveBankOfAustraliaDataSource_InvalidRate(t *testing.T) {
	dataSource := &ReserveBankOfAustraliaDataSource{}
	context := &core.Context{
		Context: &gin.Context{},
	}

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" xmlns:rba=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns=\"http://purl.org/rss/1.0/\" xsi:schemaLocation=\"http://www.w3.org/1999/02/22-rdf-syntax-ns# rdf.xsd\">\n"+
		"  <channel rdf:about=\"https://www.rba.gov.au/rss/rss-cb-exchange-rates.xml\">\n"+
		"    <dc:date>2021-04-01T16:45:00+11:00</dc:date>\n"+
		"  </channel>\n"+
		"  <item rdf:about=\"https://www.rba.gov.au/statistics/frequency/exchange-rates.html#USD\">\n"+
		"    <cb:statistics rdf:parseType=\"Resource\">\n"+
		"      <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"        <cb:observation rdf:parseType=\"Resource\">\n"+
		"          <cb:value>null</cb:value>\n"+
		"          <cb:unit>AUD</cb:unit>\n"+
		"        </cb:observation>\n"+
		"        <cb:baseCurrency>AUD</cb:baseCurrency>\n"+
		"        <cb:targetCurrency>USD</cb:targetCurrency>\n"+
		"      </cb:exchangeRate>\n"+
		"    </cb:statistics>\n"+
		"  </item>\n"+
		"</rdf:RDF>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
