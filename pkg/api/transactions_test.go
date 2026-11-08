package api

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestTransactionsApiGetTransactionResponseListResult_NormalToken(t *testing.T) {
	baseAPI, c, user := initializeInternalAvatarApiTest(t, "https://foo.bar/")
	config := baseAPI.CurrentConfig()
	config.EnableTransactionPictures = true
	config.TransactionPictureTokenExpiredTimeDuration = time.Hour
	api := &TransactionsApi{
		ApiUsingConfig: baseAPI.ApiUsingConfig,
		tokens:         baseAPI.tokens,
	}

	transactions := []*models.Transaction{{
		TransactionId: 1,
		Type:          models.TRANSACTION_DB_TYPE_EXPENSE,
	}}
	pictureInfoMap := map[int64][]*models.TransactionPictureInfo{
		1: {{
			PictureId:        456,
			PictureExtension: "png",
		}},
	}

	expectedValue := "https://foo.bar/pictures/456.png?token="
	actualValue, err := api.getTransactionResponseListResult(c, user, transactions, nil, nil, nil, nil, pictureInfoMap, time.UTC, true, true, true, true)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(actualValue))
	assert.Equal(t, 1, len(actualValue[0].Pictures))
	assert.True(t, strings.HasPrefix(actualValue[0].Pictures[0].OriginalUrl, expectedValue))

	pictureUrl, err := url.Parse(actualValue[0].Pictures[0].OriginalUrl)
	assert.Equal(t, nil, err)

	claims, _, err := api.tokens.ParseToken(c, pictureUrl.Query().Get("token"), core.NormalTokenScopeTransactionPicture)
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Uid, claims.Uid)
	assert.Equal(t, c.GetTokenClaims().UserTokenId, claims.UserTokenId)
}

func TestTransactionsApiGetTransactionResponseListResult_APIToken(t *testing.T) {
	baseAPI := initializeAvatarProviderApiTest(t, &settings.Config{
		RootUrl:                   "https://foo.bar/",
		EnableTransactionPictures: true,
	})
	api := &TransactionsApi{ApiUsingConfig: baseAPI.ApiUsingConfig}
	c := core.WrapWebContext(&gin.Context{}, nil)
	c.SetTokenClaims(&core.UserTokenClaims{Type: core.USER_TOKEN_TYPE_API, Uid: 123})

	transactions := []*models.Transaction{{
		TransactionId: 1,
		Type:          models.TRANSACTION_DB_TYPE_EXPENSE,
	}}
	pictureInfoMap := map[int64][]*models.TransactionPictureInfo{
		1: {{
			PictureId:        456,
			PictureExtension: "png",
		}},
	}

	expectedValue := models.TransactionPictureInfoBasicResponseSlice{{
		PictureId:   456,
		OriginalUrl: "https://foo.bar/pictures/456.png",
	}}
	actualValue, err := api.getTransactionResponseListResult(c, &models.User{Uid: 123}, transactions, nil, nil, nil, nil, pictureInfoMap, time.UTC, true, true, true, true)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(actualValue))
	assert.Equal(t, expectedValue, actualValue[0].Pictures)
}
