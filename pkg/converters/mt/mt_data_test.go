package mt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMtStatementGetInformationToAccountOwnerMap_OneLineMultiTags(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI/test value/ABC/123/FOO/Bar",
		},
	}

	expectedMap := map[string]string{
		"REMI": "test value",
		"ABC":  "123",
		"FOO":  "Bar",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_MultipleLines(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI/test/ABC/123",
			"/FOO/Bar/HELLO/World",
		},
	}

	expectedMap := map[string]string{
		"REMI":  "test",
		"ABC":   "123",
		"FOO":   "Bar",
		"HELLO": "World",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_EmptyInformation(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{},
	}

	expectedMap := map[string]string{}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_InvalidFormat(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/ABCD",
			"EFGH/123",
			"/REMI/123/ABC",
		},
	}

	expectedMap := map[string]string{
		"REMI": "123",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_EmptyKeyValue(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI//ABC/ /DEF/456",
			"/GHI/123/JKL/def",
		},
	}

	expectedMap := map[string]string{
		"REMI": "",
		"ABC":  "",
		"DEF":  "456",
		"GHI":  "123",
		"JKL":  "def",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGStructuredInformation(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/EREF/TEST-REFERENCE-001//CNTP/TESTACCOUNT001/TESTDEFF/Example Supplier/Berlin/REMI/USTD//Service invoice 0042/",
		},
	}

	expectedMap := map[string]string{
		"EREF": "TEST-REFERENCE-001",
		"CNTP": "TESTACCOUNT001/TESTDEFF/Example Supplier/Berlin",
		"REMI": "USTD//Service invoice 0042",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGCounterpartyEmptySubFields(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/IREF/846299123//CNTP/RO64RZBR0000075214917550//BUPA SORTREP ICONIC///REMI/USTD//2025 February ",
			"Payment/",
		},
	}

	expectedMap := map[string]string{
		"IREF": "846299123",
		"CNTP": "RO64RZBR0000075214917550//BUPA SORTREP ICONIC/",
		"REMI": "USTD//2025 February Payment",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGRemittanceTagSplitAcrossLines(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/IREF/846299138//CNTP/RO41INGB0000000008626480//HERTOGH B.V.///R",
			"EMI/USTD//2025 February Payment/",
		},
	}

	expectedMap := map[string]string{
		"IREF": "846299138",
		"CNTP": "RO41INGB0000000008626480//HERTOGH B.V./",
		"REMI": "USTD//2025 February Payment",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGRemittanceWithSlashes(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/IREF/846299125//CNTP/RO34BTRLRONCRT0000052301//FEEL FREE DOT COM///REMI/USTD//INVOICE nr. 15833/ ",
			"01.02.2025/",
		},
	}

	expectedMap := map[string]string{
		"IREF": "846299125",
		"CNTP": "RO34BTRLRONCRT0000052301//FEEL FREE DOT COM/",
		"REMI": "USTD//INVOICE nr. 15833/ 01.02.2025",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGUltimateCreditorAndDebtor(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/ULTC/Jansen/10454//ULTD/Johnson/18888/",
		},
	}

	expectedMap := map[string]string{
		"ULTC": "Jansen/10454",
		"ULTD": "Johnson/18888",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGDutchStructuredRemittance(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI/STRD/CUR/1234567890123456/",
		},
	}

	expectedMap := map[string]string{
		"REMI": "STRD/CUR/1234567890123456",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGISOStructuredRemittance(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI/STRD/ISO/RF18539007547034/",
		},
	}

	expectedMap := map[string]string{
		"REMI": "STRD/ISO/RF18539007547034",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGRemittanceWithPurposeTagInText(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI/USTD//Invoice /PURP/reference//PURP/PENS/",
		},
	}

	expectedMap := map[string]string{
		"REMI": "USTD//Invoice /PURP/reference",
		"PURP": "PENS",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}

func TestMtStatementGetInformationToAccountOwnerMap_INGRemittanceWithTrailingSlash(t *testing.T) {
	statement := &mtStatement{
		InformationToAccountOwner: []string{
			"/REMI/USTD//Invoice 2026/01//",
		},
	}

	expectedMap := map[string]string{
		"REMI": "USTD//Invoice 2026/01/",
	}

	actualMap := statement.GetInformationToAccountOwnerMap()
	assert.Equal(t, expectedMap, actualMap)
}
