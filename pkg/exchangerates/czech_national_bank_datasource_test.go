package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const czechNationalBankMinimumRequiredContent = "01 Apr 2021 #64\n" +
	"Country|Currency|Amount|Code|Rate\n" +
	"China|renminbi|1|CNY|3.379\n" +
	"USA|dollar|1|USD|22.206\n"

func TestCzechNationalBankDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(czechNationalBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "CZK", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestCzechNationalBankDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(czechNationalBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1617280200), actualLatestExchangeRateResponse.UpdateTime)
}

func TestCzechNationalBankDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(czechNationalBankMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.04503287399801856",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "CNY",
		Rate:     "0.2959455460195324",
	})
}

func TestCzechNationalBankDataSource_BlankContent(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_OnlyHeader(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("01 Apr 2021 #64"))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_OnlyHeaderAndTitle(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate"))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_MissingHeader(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("Country|Currency|Amount|Code|Rate\n"+
		"China|renminbi|1|CNY|3.379\n"+
		"USA|dollar|1|USD|22.206\n"))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_TitleMissingCode(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Rate\n"+
		"China|renminbi|1|3.379\n"+
		"USA|dollar|1|22.206\n"))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_TitleMissingAmount(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Code|Rate\n"+
		"China|renminbi|CNY|3.379\n"+
		"USA|dollar|USD|22.206\n"))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_TitleMissingRate(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code\n"+
		"China|renminbi|1|CNY\n"+
		"USA|dollar|1|USD\n"))
	assert.NotEqual(t, nil, err)
}

func TestCzechNationalBankDataSource_MissingDataItem(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"USA|dollar|1|USD\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCzechNationalBankDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"XXX|xxx|1|XXX|1\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCzechNationalBankDataSource_InvalidAmount(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"USA|dollar|null|USD|22.206\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"USA|dollar|0|USD|22.206\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCzechNationalBankDataSource_EmptyRate(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"USA|dollar|1|USD|\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestCzechNationalBankDataSource_InvalidRate(t *testing.T) {
	dataSource := &CzechNationalBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"USA|dollar|1|USD|null\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("01 Apr 2021 #64\n"+
		"Country|Currency|Amount|Code|Rate\n"+
		"USA|dollar|1|USD|0\n"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
