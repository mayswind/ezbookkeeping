package exchangerates

import (
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

const nationalBankOfKazakhstanMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
	"	<rss version=\"2.0\">\n" +
	"	<channel>\n" +
	"   	<item>\n" +
	"   		<title>USD</title>\n" +
	"			<pubDate>28.04.2026</pubDate>\n" +
	"   		<description>450.50</description>\n" +
	"   		<quant>1</quant>\n" +
	"		</item>\n" +
	"		<item>\n" +
	"   		<title>VND</title>\n" +
	"			<pubDate>28.04.2026</pubDate>\n" +
	"   		<description>0.018</description>\n" +
	"   		<quant>10</quant>\n" +
	"		</item>\n" +
	"	</channel>\n" +
	"</rss>"

func TestNationalBankOfKazakhstanDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	resp, err := dataSource.Parse(context, []byte(nationalBankOfKazakhstanMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "KZT", resp.BaseCurrency)
}

func TestNationalBankOfKazakhstanDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	resp, err := dataSource.Parse(context, []byte(nationalBankOfKazakhstanMinimumRequiredContent))
	assert.Equal(t, nil, err)

	assert.Equal(t, int64(1777316400), resp.UpdateTime)
}

func TestNationalBankOfKazakhstanDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	resp, err := dataSource.Parse(context, []byte(nationalBankOfKazakhstanMinimumRequiredContent))
	assert.Equal(t, nil, err)

	assert.Contains(t, resp.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "450.5",
	})

	assert.Contains(t, resp.ExchangeRates, &models.LatestExchangeRate{
		Currency: "VND",
		Rate:     "0.0018",
	})
}

func TestNationalBankOfKazakhstanDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfKazakhstanDataSource_EmptyData(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	content := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<rss version=\"2.0\">\n" +
		"<channel>\n" +
		"</channel>\n" +
		"</rss>"

	_, err := dataSource.Parse(context, []byte(content))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfKazakhstanDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	content := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"	<rss version=\"2.0\">\n" +
		"	<channel>\n" +
		"   	<item>\n" +
		"   		<title>XXX</title>\n" +
		"			<pubDate>28.04.2026</pubDate>\n" +
		"   		<description>450.50</description>\n" +
		"   		<quant>1</quant>\n" +
		"		</item>\n" +
		"	</channel>\n" +
		"</rss>"

	resp, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)
	assert.Len(t, resp.ExchangeRates, 0)
}

func TestNationalBankOfKazakhstanDataSource_InvalidUnit(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	content := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"	<rss version=\"2.0\">\n" +
		"	<channel>\n" +
		"   	<item>\n" +
		"   		<title>USD</title>\n" +
		"			<pubDate>28.04.2026</pubDate>\n" +
		"   		<description>450.50</description>\n" +
		"   		<quant>null</quant>\n" +
		"		</item>\n" +
		"	</channel>\n" +
		"</rss>"

	resp, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)
	assert.Len(t, resp.ExchangeRates, 0)

	content = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"	<rss version=\"2.0\">\n" +
		"	<channel>\n" +
		"   	<item>\n" +
		"   		<title>USD</title>\n" +
		"			<pubDate>28.04.2026</pubDate>\n" +
		"   		<description>450.50</description>\n" +
		"   		<quant>0</quant>\n" +
		"		</item>\n" +
		"	</channel>\n" +
		"</rss>"

	resp, err = dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)
	assert.Len(t, resp.ExchangeRates, 0)
}

func TestNationalBankOfKazakhstanDataSource_InvalidRate(t *testing.T) {
	dataSource := &NationalBankOfKazakhstanDataSource{}
	context := core.NewNullContext()

	content := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"	<rss version=\"2.0\">\n" +
		"	<channel>\n" +
		"   	<item>\n" +
		"   		<title>USD</title>\n" +
		"			<pubDate>28.04.2026</pubDate>\n" +
		"   		<description>null</description>\n" +
		"   		<quant>1</quant>\n" +
		"		</item>\n" +
		"	</channel>\n" +
		"</rss>"

	resp, err := dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)
	assert.Len(t, resp.ExchangeRates, 0)

	content = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"	<rss version=\"2.0\">\n" +
		"	<channel>\n" +
		"   	<item>\n" +
		"   		<title>USD</title>\n" +
		"			<pubDate>28.04.2026</pubDate>\n" +
		"   		<description>0</description>\n" +
		"   		<quant>1</quant>\n" +
		"		</item>\n" +
		"	</channel>\n" +
		"</rss>"

	resp, err = dataSource.Parse(context, []byte(content))
	assert.Equal(t, nil, err)
	assert.Len(t, resp.ExchangeRates, 0)
}
