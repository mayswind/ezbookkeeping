package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountInfoResponseSliceLess(t *testing.T) {
	var accountRespSlice AccountInfoResponseSlice
	accountRespSlice = append(accountRespSlice, &AccountInfoResponse{
		Id:           1,
		Category:     ACCOUNT_CATEGORY_CHECKING_ACCOUNT,
		DisplayOrder: int32(1),
	})
	accountRespSlice = append(accountRespSlice, &AccountInfoResponse{
		Id:           2,
		Category:     ACCOUNT_CATEGORY_CASH,
		DisplayOrder: int32(3),
	})
	accountRespSlice = append(accountRespSlice, &AccountInfoResponse{
		Id:           3,
		Category:     ACCOUNT_CATEGORY_CREDIT_CARD,
		DisplayOrder: int32(2),
	})
	accountRespSlice = append(accountRespSlice, &AccountInfoResponse{
		Id:           4,
		Category:     ACCOUNT_CATEGORY_CASH,
		DisplayOrder: int32(2),
	})
	accountRespSlice = append(accountRespSlice, &AccountInfoResponse{
		Id:           5,
		Category:     ACCOUNT_CATEGORY_CREDIT_CARD,
		DisplayOrder: int32(1),
	})
	accountRespSlice = append(accountRespSlice, &AccountInfoResponse{
		Id:           6,
		Category:     ACCOUNT_CATEGORY_CASH,
		DisplayOrder: int32(1),
	})

	sort.Sort(accountRespSlice)

	assert.Equal(t, int64(6), accountRespSlice[0].Id)
	assert.Equal(t, int64(4), accountRespSlice[1].Id)
	assert.Equal(t, int64(2), accountRespSlice[2].Id)
	assert.Equal(t, int64(1), accountRespSlice[3].Id)
	assert.Equal(t, int64(5), accountRespSlice[4].Id)
	assert.Equal(t, int64(3), accountRespSlice[5].Id)
}
