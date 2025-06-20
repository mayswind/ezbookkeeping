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
