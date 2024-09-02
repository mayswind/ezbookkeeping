package models

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionPictureInfoBasicResponseSliceLess(t *testing.T) {
	var pictureInfoSlice TransactionPictureInfoBasicResponseSlice
	pictureInfoSlice = append(pictureInfoSlice, &TransactionPictureInfoBasicResponse{
		PictureId: 2,
	})
	pictureInfoSlice = append(pictureInfoSlice, &TransactionPictureInfoBasicResponse{
		PictureId: 3,
	})
	pictureInfoSlice = append(pictureInfoSlice, &TransactionPictureInfoBasicResponse{
		PictureId: 1,
	})

	sort.Sort(pictureInfoSlice)

	assert.Equal(t, int64(1), pictureInfoSlice[0].PictureId)
	assert.Equal(t, int64(2), pictureInfoSlice[1].PictureId)
	assert.Equal(t, int64(3), pictureInfoSlice[2].PictureId)
}
