package api

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

type userCustomIconApiTestData struct {
	apiTestModelOperations
	records []*models.UserCustomIcon
}

func (d *userCustomIconApiTestData) Find(s *apiTestSession, rowsSlicePtr any, condiBean ...any) error {
	rows, ok := rowsSlicePtr.(*[]*models.UserCustomIcon)

	if !ok || rows == nil || s.query != "uid=? AND deleted=?" || len(s.args) != 2 || len(condiBean) != 0 {
		return s.unsupportedOperation("Find", rowsSlicePtr)
	}

	uid, uidValid := s.args[0].(int64)
	deleted, deletedValid := s.args[1].(bool)

	if !uidValid || !deletedValid {
		return s.unsupportedOperation("Find", rowsSlicePtr)
	}

	*rows = nil

	for _, icon := range d.records {
		if icon != nil && icon.Uid == uid && icon.Deleted == deleted {
			*rows = append(*rows, icon)
		}
	}

	return nil
}

func TestUserCustomIconsApi_CustomIconListHandler(t *testing.T) {
	configContainer := settings.NewConfigContainer(&settings.Config{
		RootUrl:                                "https://foo.bar/",
		SecretKey:                              "custom-icons-api-test-key",
		UserCustomIconTokenExpiredTimeDuration: time.Hour,
	})

	now := time.Now()
	expiresAt := now.Add(2 * time.Hour)

	claims := &core.UserTokenClaims{
		Uid:         123,
		UserTokenId: "789",
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		IssuedAt:    now.Unix(),
		ExpiresAt:   expiresAt.Unix(),
	}
	record := &models.TokenRecord{
		Uid:             claims.Uid,
		UserTokenId:     789,
		TokenType:       claims.Type,
		Secret:          "0123456789",
		CreatedUnixTime: claims.IssuedAt,
		ExpiredUnixTime: claims.ExpiresAt,
	}

	userDataStore, err := datastore.NewDataStore(&apiTestDatabase{
		handler: &userCustomIconApiTestData{
			records: []*models.UserCustomIcon{
				{
					IconId:       456,
					Uid:          claims.Uid,
					DisplayOrder: 2,
				},
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	tokenStore, err := datastore.NewDataStore(&apiTestDatabase{
		handler: &tokenApiTestData{
			records: []*models.TokenRecord{record},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	container := &datastore.DataStoreContainer{
		UserDataStore: userDataStore,
		TokenStore:    tokenStore,
	}

	api := &UserCustomIconsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: configContainer,
		},
		icons: &services.UserCustomIconService{
			ServiceUsingDB: services.NewServiceUsingDB(container),
		},
		tokens: &services.TokenService{
			ServiceUsingDB:     services.NewServiceUsingDB(container),
			ServiceUsingConfig: services.NewServiceUsingConfig(configContainer),
		},
	}
	c := core.WrapWebContext(&gin.Context{}, nil)
	c.SetTokenClaims(claims)
	c.SetTokenRecord(record)

	result, err := api.CustomIconListHandler(c)
	assert.Nil(t, err)

	iconList := result.(*models.UserCustomIconListResponse)
	icons := iconList.Items
	assert.Len(t, icons, 1)
	assert.True(t, strings.HasPrefix(icons[0].OriginalUrl, "https://foo.bar/icons/456.png?token=ebk-tk.i-"))

	iconUrl, err := url.Parse(icons[0].OriginalUrl)
	assert.Nil(t, err)

	actualClaims, _, err := api.tokens.ParseToken(c, iconUrl.Query().Get("token"), core.NormalTokenScopeUserCustomIcon)
	assert.Nil(t, err)
	assert.Equal(t, claims.Uid, actualClaims.Uid)
	assert.Equal(t, claims.UserTokenId, actualClaims.UserTokenId)
	assert.Equal(t, actualClaims.ExpiresAt, iconList.ExpiresAt)
	assert.Equal(t, now.Add(1*time.Hour).Unix(), iconList.ExpiresAt)

	_, _, err = api.tokens.ParseToken(c, iconUrl.Query().Get("token"))
	assert.Equal(t, errs.ErrCurrentInvalidToken, err)

	c.SetTokenClaims(&core.UserTokenClaims{
		Uid:  claims.Uid,
		Type: core.USER_TOKEN_TYPE_API,
	})
	result, err = api.CustomIconListHandler(c)
	assert.Nil(t, err)

	iconList = result.(*models.UserCustomIconListResponse)
	icons = iconList.Items
	assert.Equal(t, "https://foo.bar/icons/456.png", icons[0].OriginalUrl)
	assert.Equal(t, int64(0), iconList.ExpiresAt)
}
