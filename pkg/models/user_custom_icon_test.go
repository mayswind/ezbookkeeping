package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCustomIconInfoResponseSliceLess(t *testing.T) {
	var userCustomIconInfoRespSlice UserCustomIconInfoResponseSlice
	userCustomIconInfoRespSlice = append(userCustomIconInfoRespSlice, &UserCustomIconInfoResponse{
		Id:           1,
		DisplayOrder: 3,
	})
	userCustomIconInfoRespSlice = append(userCustomIconInfoRespSlice, &UserCustomIconInfoResponse{
		Id:           2,
		DisplayOrder: 1,
	})
	userCustomIconInfoRespSlice = append(userCustomIconInfoRespSlice, &UserCustomIconInfoResponse{
		Id:           3,
		DisplayOrder: 2,
	})

	sort.Sort(userCustomIconInfoRespSlice)

	assert.Equal(t, int64(2), userCustomIconInfoRespSlice[0].Id)
	assert.Equal(t, int64(3), userCustomIconInfoRespSlice[1].Id)
	assert.Equal(t, int64(1), userCustomIconInfoRespSlice[2].Id)
}
