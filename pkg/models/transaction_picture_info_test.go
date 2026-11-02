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

func TestTransactionPictureInfo_ToTransactionPictureInfoBasicResponse_WithToken(t *testing.T) {
	pictureInfo := &TransactionPictureInfo{
		PictureId: 123,
	}

	expectedValue := &TransactionPictureInfoBasicResponse{
		PictureId:   123,
		OriginalUrl: "https://foo.bar/pictures/123.png?token=media%2Btoken",
	}
	actualValue := pictureInfo.ToTransactionPictureInfoBasicResponse("https://foo.bar/pictures/123.png", "media+token")
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionPictureInfo_ToTransactionPictureInfoBasicResponse_WithTokenAndQueryString(t *testing.T) {
	pictureInfo := &TransactionPictureInfo{
		PictureId: 123,
	}

	expectedValue := &TransactionPictureInfoBasicResponse{
		PictureId:   123,
		OriginalUrl: "https://foo.bar/pictures/123.png?size=64&token=media%2Btoken",
	}
	actualValue := pictureInfo.ToTransactionPictureInfoBasicResponse("https://foo.bar/pictures/123.png?size=64", "media+token")
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionPictureInfo_ToTransactionPictureInfoBasicResponse_WithoutTokenAndWithQueryString(t *testing.T) {
	pictureInfo := &TransactionPictureInfo{
		PictureId: 123,
	}

	expectedValue := &TransactionPictureInfoBasicResponse{
		PictureId:   123,
		OriginalUrl: "https://foo.bar/pictures/123.png?size=64",
	}
	actualValue := pictureInfo.ToTransactionPictureInfoBasicResponse("https://foo.bar/pictures/123.png?size=64", "")
	assert.Equal(t, expectedValue, actualValue)
}
