package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportPresetInfoResponseSliceLess(t *testing.T) {
	var insightsExplorerRespSlice ImportPresetInfoResponseSlice
	insightsExplorerRespSlice = append(insightsExplorerRespSlice, &ImportPresetInfoResponse{
		Id:           1,
		DisplayOrder: 3,
	})
	insightsExplorerRespSlice = append(insightsExplorerRespSlice, &ImportPresetInfoResponse{
		Id:           2,
		DisplayOrder: 1,
	})
	insightsExplorerRespSlice = append(insightsExplorerRespSlice, &ImportPresetInfoResponse{
		Id:           3,
		DisplayOrder: 2,
	})

	sort.Sort(insightsExplorerRespSlice)

	assert.Equal(t, int64(2), insightsExplorerRespSlice[0].Id)
	assert.Equal(t, int64(3), insightsExplorerRespSlice[1].Id)
	assert.Equal(t, int64(1), insightsExplorerRespSlice[2].Id)
}
