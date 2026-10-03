package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const nationalBankOfBelarusMinimumRequiredContent = "[\n" +
	"  {\n" +
	"    \"Cur_Abbreviation\": \"USD\",\n" +
	"    \"Cur_Scale\": 1,\n" +
	"    \"Cur_OfficialRate\": 4,\n" +
	"    \"Date\": \"2026-10-02T00:00:00\"\n" +
	"  },\n" +
	"  {\n" +
	"    \"Cur_Abbreviation\": \"JPY\",\n" +
	"    \"Cur_Scale\": 100,\n" +
	"    \"Cur_OfficialRate\": 2,\n" +
	"    \"Date\": \"2026-10-02T00:00:00\"\n" +
	"  }\n" +
	"]"

func TestNationalBankOfBelarusDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfBelarusMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "BYN", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNationalBankOfBelarusDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfBelarusMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1790888400), actualLatestExchangeRateResponse.UpdateTime)
}

func TestNationalBankOfBelarusDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(nationalBankOfBelarusMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.25",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "50",
	})
}

func TestNationalBankOfBelarusDataSource_BlankContent(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfBelarusDataSource_EmptyData(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("{}"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfBelarusDataSource_EmptyExchangeRatesData(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[]"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfBelarusDataSource_InvalidUpdateTime(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"USD\",\n"+
		"    \"Cur_Scale\": 1,\n"+
		"    \"Cur_OfficialRate\": 4,\n"+
		"    \"Date\": \"invalid\"\n"+
		"  }\n"+
		"]"))
	assert.NotEqual(t, nil, err)
}

func TestNationalBankOfBelarusDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"XXX\",\n"+
		"    \"Cur_Scale\": 1,\n"+
		"    \"Cur_OfficialRate\": 4,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  },\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"JPY\",\n"+
		"    \"Cur_Scale\": 100,\n"+
		"    \"Cur_OfficialRate\": 2,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
}

func TestNationalBankOfBelarusDataSource_InvalidUnit(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"USD\",\n"+
		"    \"Cur_Scale\": null,\n"+
		"    \"Cur_OfficialRate\": 4,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  },\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"JPY\",\n"+
		"    \"Cur_Scale\": 100,\n"+
		"    \"Cur_OfficialRate\": 2,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"USD\",\n"+
		"    \"Cur_Scale\": 0,\n"+
		"    \"Cur_OfficialRate\": 4,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  },\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"JPY\",\n"+
		"    \"Cur_Scale\": 100,\n"+
		"    \"Cur_OfficialRate\": 2,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
}

func TestNationalBankOfBelarusDataSource_InvalidRate(t *testing.T) {
	dataSource := &NationalBankOfBelarusDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"USD\",\n"+
		"    \"Cur_Scale\": 1,\n"+
		"    \"Cur_OfficialRate\": null,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  },\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"JPY\",\n"+
		"    \"Cur_Scale\": 100,\n"+
		"    \"Cur_OfficialRate\": 2,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("[\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"USD\",\n"+
		"    \"Cur_Scale\": 1,\n"+
		"    \"Cur_OfficialRate\": 0,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  },\n"+
		"  {\n"+
		"    \"Cur_Abbreviation\": \"JPY\",\n"+
		"    \"Cur_Scale\": 100,\n"+
		"    \"Cur_OfficialRate\": 2,\n"+
		"    \"Date\": \"2026-10-02T00:00:00\"\n"+
		"  }\n"+
		"]"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 1)
}
