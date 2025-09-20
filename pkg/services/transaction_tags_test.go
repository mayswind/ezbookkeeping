package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestGetTagMapByList_EmptyList(t *testing.T) {
	tags := make([]*models.TransactionTag, 0)
	actualTagMap := TransactionTags.GetTagMapByList(tags)

	assert.NotNil(t, actualTagMap)
	assert.Equal(t, 0, len(actualTagMap))
}

func TestGetTagMapByList_SingleTag(t *testing.T) {
	tags := []*models.TransactionTag{
		{
			TagId:  1001,
			Name:   "Tag Name",
			Hidden: false,
		},
	}
	actualTagMap := TransactionTags.GetTagMapByList(tags)

	assert.Equal(t, 1, len(actualTagMap))
	assert.Contains(t, actualTagMap, int64(1001))
	assert.Equal(t, "Tag Name", actualTagMap[1001].Name)
	assert.Equal(t, false, actualTagMap[1001].Hidden)
}

func TestGetTagMapByList_MultipleTags(t *testing.T) {
	tags := []*models.TransactionTag{
		{
			TagId:  1001,
			Name:   "Tag Name",
			Hidden: false,
		},
		{
			TagId:  1002,
			Name:   "Tag Name2",
			Hidden: true,
		},
		{
			TagId:  1003,
			Name:   "Tag Name3",
			Hidden: false,
		},
	}
	actualTagMap := TransactionTags.GetTagMapByList(tags)

	assert.Equal(t, 3, len(actualTagMap))
	assert.Contains(t, actualTagMap, int64(1001))
	assert.Contains(t, actualTagMap, int64(1002))
	assert.Contains(t, actualTagMap, int64(1003))
	assert.Equal(t, "Tag Name", actualTagMap[1001].Name)
	assert.Equal(t, "Tag Name2", actualTagMap[1002].Name)
	assert.Equal(t, "Tag Name3", actualTagMap[1003].Name)
	assert.Equal(t, false, actualTagMap[1001].Hidden)
	assert.Equal(t, true, actualTagMap[1002].Hidden)
	assert.Equal(t, false, actualTagMap[1003].Hidden)
}

func TestGetVisibleTagNameMapByList_EmptyList(t *testing.T) {
	tags := make([]*models.TransactionTag, 0)
	actualTagMap := TransactionTags.GetVisibleTagNameMapByList(tags)

	assert.NotNil(t, actualTagMap)
	assert.Equal(t, 0, len(actualTagMap))
}

func TestGetVisibleTagNameMapByList_MixedVisibilityTags(t *testing.T) {
	tags := []*models.TransactionTag{
		{
			TagId:  1001,
			Name:   "Visible Tag",
			Hidden: false,
		},
		{
			TagId:  1002,
			Name:   "Hidden Tag",
			Hidden: true,
		},
		{
			TagId:  1003,
			Name:   "Visible Tag2",
			Hidden: false,
		},
	}
	actualTagMap := TransactionTags.GetVisibleTagNameMapByList(tags)

	assert.Equal(t, 2, len(actualTagMap))
	assert.Contains(t, actualTagMap, "Visible Tag")
	assert.Contains(t, actualTagMap, "Visible Tag2")
	assert.NotContains(t, actualTagMap, "Hidden Tag")
	assert.Equal(t, int64(1001), actualTagMap["Visible Tag"].TagId)
	assert.Equal(t, int64(1003), actualTagMap["Visible Tag2"].TagId)
}

func TestGetTagNames_EmptyList(t *testing.T) {
	tags := make([]*models.TransactionTag, 0)
	actualNames := TransactionTags.GetTagNames(tags)

	assert.NotNil(t, actualNames)
	assert.Equal(t, 0, len(actualNames))
}

func TestGetTagNames_MultipleTags(t *testing.T) {
	tags := []*models.TransactionTag{
		{
			TagId: 1001,
			Name:  "Tag Name",
		},
		{
			TagId: 1002,
			Name:  "Tag Name2",
		},
		{
			TagId: 1003,
			Name:  "Tag Name3",
		},
	}
	actualNames := TransactionTags.GetTagNames(tags)

	assert.Equal(t, 3, len(actualNames))
	assert.Equal(t, "Tag Name", actualNames[0])
	assert.Equal(t, "Tag Name2", actualNames[1])
	assert.Equal(t, "Tag Name3", actualNames[2])
}

func TestGetTagIds_EmptyString(t *testing.T) {
	tagIds, err := TransactionTags.GetTagIds("")

	assert.Nil(t, err)
	assert.Nil(t, tagIds)
}

func TestGetTagIds_ZeroString(t *testing.T) {
	tagIds, err := TransactionTags.GetTagIds("0")

	assert.Nil(t, err)
	assert.Nil(t, tagIds)
}

func TestGetTagIds_SingleId(t *testing.T) {
	tagIds, err := TransactionTags.GetTagIds("1001")

	assert.Nil(t, err)
	assert.Equal(t, 1, len(tagIds))
	assert.Equal(t, int64(1001), tagIds[0])
}

func TestGetTagIds_MultipleIds(t *testing.T) {
	tagIds, err := TransactionTags.GetTagIds("1001,1002,1003")

	assert.Nil(t, err)
	assert.Equal(t, 3, len(tagIds))
	assert.Equal(t, int64(1001), tagIds[0])
	assert.Equal(t, int64(1002), tagIds[1])
	assert.Equal(t, int64(1003), tagIds[2])
}

func TestGetTagIds_InvalidId(t *testing.T) {
	tagIds, err := TransactionTags.GetTagIds("1001,invalid,1003")

	assert.NotNil(t, err)
	assert.Nil(t, tagIds)
}

func TestGetTransactionTagIds_EmptyMap(t *testing.T) {
	allTransactionTagIds := make(map[int64][]int64)
	actualTagIds := TransactionTags.GetTransactionTagIds(allTransactionTagIds)

	assert.NotNil(t, actualTagIds)
	assert.Equal(t, 0, len(actualTagIds))
}

func TestGetTransactionTagIds_MultipleTransactions(t *testing.T) {
	allTransactionTagIds := map[int64][]int64{
		1001: {2001, 2002},
		1002: {2003},
		1003: {2001, 2004},
	}
	actualTagIds := TransactionTags.GetTransactionTagIds(allTransactionTagIds)

	assert.Equal(t, 5, len(actualTagIds))
	assert.Contains(t, actualTagIds, int64(2001))
	assert.Contains(t, actualTagIds, int64(2002))
	assert.Contains(t, actualTagIds, int64(2003))
	assert.Contains(t, actualTagIds, int64(2001))
	assert.Contains(t, actualTagIds, int64(2004))
}

func TestGetTransactionTagIds_EmptyTransactionTagSlices(t *testing.T) {
	allTransactionTagIds := map[int64][]int64{
		1001: {},
		1002: {},
	}
	actualTagIds := TransactionTags.GetTransactionTagIds(allTransactionTagIds)

	assert.NotNil(t, actualTagIds)
	assert.Equal(t, 0, len(actualTagIds))
}

func TestGetTransactionTagIds_MixedTransactionTagSlices(t *testing.T) {
	allTransactionTagIds := map[int64][]int64{
		1001: {2001, 2002},
		1002: {},
		1003: {2003},
	}
	actualTagIds := TransactionTags.GetTransactionTagIds(allTransactionTagIds)

	assert.Equal(t, 3, len(actualTagIds))
	assert.Contains(t, actualTagIds, int64(2001))
	assert.Contains(t, actualTagIds, int64(2002))
	assert.Contains(t, actualTagIds, int64(2003))
}
