package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenInfoResponseSliceLess(t *testing.T) {
	var tokenInfoRespSlice TokenInfoResponseSlice
	tokenInfoRespSlice = append(tokenInfoRespSlice, &TokenInfoResponse{
		TokenId:  "1",
		LastSeen: int64(1),
	})
	tokenInfoRespSlice = append(tokenInfoRespSlice, &TokenInfoResponse{
		TokenId:  "2",
		LastSeen: int64(3),
	})
	tokenInfoRespSlice = append(tokenInfoRespSlice, &TokenInfoResponse{
		TokenId:  "3",
		LastSeen: int64(2),
	})

	sort.Sort(tokenInfoRespSlice)

	assert.Equal(t, "2", tokenInfoRespSlice[0].TokenId)
	assert.Equal(t, "3", tokenInfoRespSlice[1].TokenId)
	assert.Equal(t, "1", tokenInfoRespSlice[2].TokenId)
}
