package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const SwissNationalBankMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
	"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n" +
	"  <channel>\n" +
	"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n" +
	"    <item>\n" +
	"      <cb:statistics rdf:parseType=\"Resource\">\n" +
	"        <cb:exchangeRate rdf:parseType=\"Resource\">\n" +
	"          <cb:observation rdf:parseType=\"Resource\">\n" +
	"            <cb:value>0.9378</cb:value>\n" +
	"            <cb:unit>CHF</cb:unit>\n" +
	"            <cb:unit_mult>1</cb:unit_mult>\n" +
	"          </cb:observation>\n" +
	"          <cb:baseCurrency>CHF</cb:baseCurrency>\n" +
	"          <cb:targetCurrency>EUR</cb:targetCurrency>\n" +
	"          <cb:observationPeriod rdf:parseType=\"Resource\">\n" +
	"            <cb:period>2024-11-12</cb:period>\n" +
	"          </cb:observationPeriod>\n" +
	"        </cb:exchangeRate>\n" +
	"      </cb:statistics>\n" +
	"    </item>\n" +
	"    <item>\n" +
	"      <cb:statistics rdf:parseType=\"Resource\">\n" +
	"        <cb:exchangeRate rdf:parseType=\"Resource\">\n" +
	"          <cb:observation rdf:parseType=\"Resource\">\n" +
	"            <cb:value>0.5727</cb:value>\n" +
	"            <cb:unit>CHF</cb:unit>\n" +
	"            <cb:unit_mult>-2</cb:unit_mult>\n" +
	"          </cb:observation>\n" +
	"          <cb:baseCurrency>CHF</cb:baseCurrency>\n" +
	"          <cb:targetCurrency>JPY</cb:targetCurrency>\n" +
	"          <cb:observationPeriod rdf:parseType=\"Resource\">\n" +
	"            <cb:period>2024-11-12</cb:period>\n" +
	"          </cb:observationPeriod>\n" +
	"        </cb:exchangeRate>\n" +
	"      </cb:statistics>\n" +
	"    </item>\n" +
	"  </channel>\n" +
	"</rss>"

func TestSwissNationalBankDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(SwissNationalBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "CHF", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestSwissNationalBankDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(SwissNationalBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731409250), actualLatestExchangeRateResponse.UpdateTime)
}

func TestSwissNationalBankDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(SwissNationalBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "EUR",
		Rate:     "1.0663254425250588",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "174.6114894360049",
	})
}

func TestSwissNationalBankDataSource_MultipleDateExchanges(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9378</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9381</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-11</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "EUR",
		Rate:     "1.0663254425250588",
	})
}

func TestSwissNationalBankDataSource_BlankContent(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestSwissNationalBankDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestSwissNationalBankDataSource_EmptyRDFContent(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"</rss>"))
	assert.NotEqual(t, nil, err)
}

func TestSwissNationalBankDataSource_EmptyChannelContent(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.NotEqual(t, nil, err)
}

func TestSwissNationalBankDataSource_NoItem(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.NotEqual(t, nil, err)
}

func TestSwissNationalBankDataSource_BaseCurrencyNotEqualPreset(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9378</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>EUR</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>CHF</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestSwissNationalBankDataSource_UnitCurrencyNotEqualPreset(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9378</cb:value>\n"+
		"            <cb:unit>EUR</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestSwissNationalBankDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9378</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>XXX</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestSwissNationalBankDataSource_EmptyRate(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value></cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestSwissNationalBankDataSource_InvalidRate(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>null</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>1</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestSwissNationalBankDataSource_InvalidUnit(t *testing.T) {
	dataSource := &SwissNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9378</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>null</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
		"<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" xmlns:cb=\"http://www.cbwiki.net/wiki/index.php/Specification_1.2/\" xmlns:dc=\"http://purl.org/dc/elements/1.1/\" xmlns:dcterms=\"http://purl.org/dc/terms/\" xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\" version=\"2.0\">\n"+
		"  <channel>\n"+
		"    <pubDate>Tue, 12 Nov 2024 11:00:50 GMT</pubDate>\n"+
		"    <item>\n"+
		"      <cb:statistics rdf:parseType=\"Resource\">\n"+
		"        <cb:exchangeRate rdf:parseType=\"Resource\">\n"+
		"          <cb:observation rdf:parseType=\"Resource\">\n"+
		"            <cb:value>0.9378</cb:value>\n"+
		"            <cb:unit>CHF</cb:unit>\n"+
		"            <cb:unit_mult>0</cb:unit_mult>\n"+
		"          </cb:observation>\n"+
		"          <cb:baseCurrency>CHF</cb:baseCurrency>\n"+
		"          <cb:targetCurrency>EUR</cb:targetCurrency>\n"+
		"          <cb:observationPeriod rdf:parseType=\"Resource\">\n"+
		"            <cb:period>2024-11-12</cb:period>\n"+
		"          </cb:observationPeriod>\n"+
		"        </cb:exchangeRate>\n"+
		"      </cb:statistics>\n"+
		"    </item>\n"+
		"  </channel>\n"+
		"</rss>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
