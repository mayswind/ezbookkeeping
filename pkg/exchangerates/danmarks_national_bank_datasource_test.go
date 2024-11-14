package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const danmarksNationalbankMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
	"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">\n" +
	"  <dailyrates id=\"2024-11-14\">\n" +
	"    <currency code=\"CNY\" desc=\"Chinese yuan renminbi\" rate=\"97.81\" />\n" +
	"    <currency code=\"USD\" desc=\"US dollars\" rate=\"708.18\" />\n" +
	"  </dailyrates>\n" +
	"</exchangerates>"

func TestDanmarksNationalbankDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(danmarksNationalbankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "DKK", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestDanmarksNationalbankDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(danmarksNationalbankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731596400), actualLatestExchangeRateResponse.UpdateTime)
}

func TestDanmarksNationalbankDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(danmarksNationalbankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.1412070377587619",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "1.022390348635109",
	})
}

func TestDanmarksNationalbankDataSource_BlankContent(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestDanmarksNationalbankDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestDanmarksNationalbankDataSource_EmptyExchangeRatesContent(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">"+
		"</exchangerates>"))
	assert.NotEqual(t, nil, err)
}

func TestDanmarksNationalbankDataSource_EmptyDailyRatesContent(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">"+
		"<dailyrates id=\"2024-11-14\">"+
		"</dailyrates>"+
		"</exchangerates>"))
	assert.NotEqual(t, nil, err)
}

func TestDanmarksNationalbankDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">"+
		"  <dailyrates id=\"2024-11-14\">\n"+
		"    <currency code=\"XXX\" desc=\"XXX\" rate=\"1\" />\n"+
		"  </dailyrates>\n"+
		"</exchangerates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestDanmarksNationalbankDataSource_EmptyRate(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">"+
		"  <dailyrates id=\"2024-11-14\">\n"+
		"    <currency code=\"USD\" desc=\"US dollars\" rate=\"\" />\n"+
		"  </dailyrates>\n"+
		"</exchangerates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestDanmarksNationalbankDataSource_InvalidRate(t *testing.T) {
	dataSource := &DanmarksNationalbankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">"+
		"  <dailyrates id=\"2024-11-14\">\n"+
		"    <currency code=\"USD\" desc=\"US dollars\" rate=\"null\" />\n"+
		"  </dailyrates>\n"+
		"</exchangerates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<exchangerates type=\"Exchange rates\" author=\"Danmarks Nationalbank\" refcur=\"DKK\" refamt=\"1\">"+
		"  <dailyrates id=\"2024-11-14\">\n"+
		"    <currency code=\"USD\" desc=\"US dollars\" rate=\"0\" />\n"+
		"  </dailyrates>\n"+
		"</exchangerates>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
