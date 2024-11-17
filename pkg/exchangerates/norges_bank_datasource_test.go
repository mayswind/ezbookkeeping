package exchangerates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const norgesBankOfRomaniaMinimumRequiredContent = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
	"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">\n" +
	"  <message:DataSet>\n" +
	"    <Series BASE_CUR=\"JPY\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"2\">\n" +
	"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"7.1179\" />\n" +
	"    </Series>\n" +
	"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"0\">\n" +
	"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"11.0545\" />\n" +
	"    </Series>\n" +
	"  </message:DataSet>\n" +
	"</message:StructureSpecificData>"

func TestNorgesBankDataSource_StandardDataExtractBaseCurrency(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(norgesBankOfRomaniaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, "NOK", actualLatestExchangeRateResponse.BaseCurrency)
}

func TestNorgesBankDataSource_StandardDataExtractUpdateTime(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(norgesBankOfRomaniaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1731682800), actualLatestExchangeRateResponse.UpdateTime)
}

func TestNorgesBankDataSource_StandardDataExtractExchangeRates(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte(norgesBankOfRomaniaMinimumRequiredContent))
	assert.Equal(t, nil, err)
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "JPY",
		Rate:     "14.049087511766112",
	})
	assert.Contains(t, actualLatestExchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
		Currency: "USD",
		Rate:     "0.09046089827671988",
	})
}

func TestNorgesBankDataSource_BlankContent(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte(""))
	assert.NotEqual(t, nil, err)
}

func TestNorgesBankDataSource_OnlyXMLHeader(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"))
	assert.NotEqual(t, nil, err)
}

func TestNorgesBankDataSource_MissingExchangeRatesDataset(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"</message:StructureSpecificData>"))
	assert.NotEqual(t, nil, err)
}

func TestNorgesBankDataSource_EmptyExchangeRatesDataset(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	_, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.NotEqual(t, nil, err)
}

func TestNorgesBankDataSource_EmptyExchangeRateObservations(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"0\">\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNorgesBankDataSource_InvalidCurrency(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"XXX\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"0\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"1\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNorgesBankDataSource_InvalidTargetCurrency(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"EUR\" UNIT_MULT=\"0\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"11.0545\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNorgesBankDataSource_EmptyRate(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"0\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNorgesBankDataSource_InvalidRate(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"0\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"null\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"0\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"0\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}

func TestNorgesBankDataSource_InvalidUnit(t *testing.T) {
	dataSource := &NorgesBankDataSource{}
	context := core.NewNullContext()

	actualLatestExchangeRateResponse, err := dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"null\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"11.0545\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"11.0545\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)

	actualLatestExchangeRateResponse, err = dataSource.Parse(context, []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>"+
		"<message:StructureSpecificData xmlns:ss=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/data/structurespecific\" xmlns:footer=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message/footer\" xmlns:ns1=\"urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD\" xmlns:message=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message\" xmlns:common=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/common\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xml=\"http://www.w3.org/XML/1998/namespace\" xsi:schemaLocation=\"http://www.sdmx.org/resources/sdmxml/schemas/v2_1/message https://registry.sdmx.org/schemas/v2_1/SDMXMessage.xsd urn:sdmx:org.sdmx.infomodel.datastructure.Dataflow=NB:EXR(1.0):ObsLevelDim:TIME_PERIOD https://data.norges-bank.no/api/schema/dataflow/NB/EXR/1.0?format=sdmx-2.1\">"+
		"  <message:DataSet>\n"+
		"    <Series BASE_CUR=\"USD\" QUOTE_CUR=\"NOK\" UNIT_MULT=\"-1\">\n"+
		"      <Obs TIME_PERIOD=\"2024-11-15\" OBS_VALUE=\"11.0545\" />\n"+
		"    </Series>\n"+
		"  </message:DataSet>\n"+
		"</message:StructureSpecificData>"))
	assert.Equal(t, nil, err)
	assert.Len(t, actualLatestExchangeRateResponse.ExchangeRates, 0)
}
