package api

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/avatars"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

type apiTestDatabase struct {
	datastore.Database
	handler apiTestModelHandler
}

type apiTestSession struct {
	datastore.DatabaseSession
	database *apiTestDatabase
	query    any
	args     []any
}

type apiTestModelHandler interface {
	Find(s *apiTestSession, rowsSlicePtr any, condiBean ...any) error
	Get(s *apiTestSession, beans ...any) (bool, error)
	Insert(s *apiTestSession, beans ...any) (int64, error)
	Delete(s *apiTestSession, beans ...any) (int64, error)
}

type apiTestModelOperations struct {
}

func (db *apiTestDatabase) NewSession(c core.Context) datastore.DatabaseSession {
	return &apiTestSession{
		database: db,
	}
}

func (db *apiTestDatabase) DoTransaction(c core.Context, fn func(datastore.DatabaseSession) error) error {
	return fn(db.NewSession(c))
}

func (s *apiTestSession) Where(query any, args ...any) datastore.DatabaseSession {
	s.query = query
	s.args = args
	return s
}

func (s *apiTestSession) Find(rowsSlicePtr any, condiBean ...any) error {
	handler := s.database.handler

	if handler == nil {
		return s.unsupportedOperation("Find", rowsSlicePtr)
	}

	return handler.Find(s, rowsSlicePtr, condiBean...)
}

func (s *apiTestSession) Limit(limit int, start ...int) datastore.DatabaseSession {
	return s
}

func (s *apiTestSession) Get(beans ...any) (bool, error) {
	if len(beans) != 1 {
		return false, s.unsupportedOperation("Get", beans)
	}

	handler := s.database.handler

	if handler == nil {
		return false, s.unsupportedOperation("Get", beans[0])
	}

	return handler.Get(s, beans...)
}

func (s *apiTestSession) Insert(beans ...any) (int64, error) {
	if len(beans) != 1 {
		return 0, s.unsupportedOperation("Insert", beans)
	}

	handler := s.database.handler

	if handler == nil {
		return 0, s.unsupportedOperation("Insert", beans[0])
	}

	return handler.Insert(s, beans...)
}

func (s *apiTestSession) Delete(beans ...any) (int64, error) {
	if len(beans) != 1 {
		return 0, s.unsupportedOperation("Delete", beans)
	}

	handler := s.database.handler

	if handler == nil {
		return 0, s.unsupportedOperation("Delete", beans[0])
	}

	return handler.Delete(s, beans...)
}

func (s *apiTestSession) unsupportedOperation(operation string, bean any) error {
	return fmt.Errorf("unsupported mock %s for %T (query: %v, args: %v)", operation, bean, s.query, s.args)
}

func (o *apiTestModelOperations) Find(s *apiTestSession, rowsSlicePtr any, condiBean ...any) error {
	return s.unsupportedOperation("Find", rowsSlicePtr)
}

func (o *apiTestModelOperations) Get(s *apiTestSession, beans ...any) (bool, error) {
	return false, s.unsupportedOperation("Get", beans)
}

func (o *apiTestModelOperations) Insert(s *apiTestSession, beans ...any) (int64, error) {
	return 0, s.unsupportedOperation("Insert", beans)
}

func (o *apiTestModelOperations) Delete(s *apiTestSession, beans ...any) (int64, error) {
	return 0, s.unsupportedOperation("Delete", beans)
}

func initializeAvatarProviderApiTest(t *testing.T, config *settings.Config) *ApiWithUserInfo {
	t.Helper()
	var provider avatars.AvatarProvider

	if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
		provider = avatars.NewInternalStorageAvatarProvider(config)
	} else if config.AvatarProvider == core.USER_AVATAR_PROVIDER_GRAVATAR {
		provider = avatars.NewGravatarAvatarProvider()
	} else if config.AvatarProvider == "" {
		provider = avatars.NewNullAvatarProvider()
	} else {
		t.Fatal(errs.ErrInvalidAvatarProvider)
	}

	configContainer := settings.NewConfigContainer(config)

	return &ApiWithUserInfo{
		ApiUsingConfig: ApiUsingConfig{
			container: configContainer,
		},
		ApiUsingAvatarProvider: ApiUsingAvatarProvider{
			avatarProvider: provider,
			tokens: &services.TokenService{
				ServiceUsingConfig: services.NewServiceUsingConfig(configContainer),
			},
		},
	}
}

func initializeInternalAvatarApiTest(t *testing.T, rootUrl string) (*ApiWithUserInfo, *core.WebContext, *models.User) {
	t.Helper()
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider:                     core.USER_AVATAR_PROVIDER_INTERNAL,
		RootUrl:                            rootUrl,
		SecretKey:                          "avatar-api-test-key",
		TokenExpiredTimeDuration:           2 * time.Hour,
		UserAvatarTokenExpiredTimeDuration: time.Hour,
	})

	tokenStore, err := datastore.NewDataStore(&apiTestDatabase{
		handler: &tokenApiTestData{},
	})

	if err != nil {
		t.Fatal(err)
	}

	api.tokens.ServiceUsingDB = services.NewServiceUsingDB(&datastore.DataStoreContainer{
		TokenStore: tokenStore,
	})

	c := core.WrapWebContext(&gin.Context{}, nil)

	user := &models.User{
		Uid:              123,
		Username:         "test-user",
		CustomAvatarType: "png",
	}

	_, claims, tokenRecord, err := api.tokens.CreateToken(c, user)

	if err != nil {
		t.Fatal(err)
	}

	c.SetTokenClaims(claims)
	c.SetTokenRecord(tokenRecord)

	return api, c, user
}

func TestApiUsingAvatarProvider_GetAvatarUrl_Gravatar(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_GRAVATAR,
	})

	expectedValue := "https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"
	actualValue, err := api.GetAvatarUrl(nil, &models.User{
		Email: "MyEmailAddress@example.com",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestApiUsingAvatarProvider_GetAvatarUrl_EmptyCustomAvatarType(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_INTERNAL,
	})

	expectedValue := ""
	actualValue, err := api.GetAvatarUrl(nil, &models.User{
		Uid: 123,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestApiUsingAvatarProvider_GetAvatarUrl_NoProvider(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{})
	expectedValue := ""

	actualValue, err := api.GetAvatarUrl(nil, &models.User{
		Uid: 123,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestApiUsingAvatarProvider_GetAvatarUrl_Internal(t *testing.T) {
	api, c, user := initializeInternalAvatarApiTest(t, "https://foo.bar/")

	expectedValue := "https://foo.bar/avatar/123.png?token="
	actualValue, err := api.GetAvatarUrl(c, user)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, expectedValue))

	avatarUrl, err := url.Parse(actualValue)
	assert.Equal(t, nil, err)

	claims, _, err := api.tokens.ParseToken(c, avatarUrl.Query().Get("token"), core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Uid, claims.Uid)
	assert.Equal(t, c.GetTokenClaims().UserTokenId, claims.UserTokenId)
}

func TestApiUsingAvatarProvider_GetAvatarUrl_InternalWithQueryString(t *testing.T) {
	api, c, user := initializeInternalAvatarApiTest(t, "https://foo.bar/?path=/")

	expectedValue := "https://foo.bar/?path=/avatar/123.png&token="
	actualValue, err := api.GetAvatarUrl(c, user)
	assert.Equal(t, nil, err)
	assert.True(t, strings.HasPrefix(actualValue, expectedValue))

	avatarUrl, err := url.Parse(actualValue)
	assert.Equal(t, nil, err)
	assert.Equal(t, "/avatar/123.png", avatarUrl.Query().Get("path"))

	claims, _, err := api.tokens.ParseToken(c, avatarUrl.Query().Get("token"), core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Uid, claims.Uid)
}

func TestApiWithUserInfo_GetUserBasicInfo_AvatarTokenError(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_INTERNAL,
	})

	c := core.WrapWebContext(&gin.Context{}, nil)
	c.SetTokenClaims(&core.UserTokenClaims{
		Type:        core.USER_TOKEN_TYPE_NORMAL,
		UserTokenId: "invalid",
	})

	actualValue, err := api.GetUserBasicInfo(c, &models.User{
		Uid:              123,
		CustomAvatarType: "png",
	})
	assert.Equal(t, errs.ErrInvalidToken, err)
	assert.Nil(t, actualValue)
}

func TestApiWithUserInfo_GetUserBasicInfo_InternalAvatar(t *testing.T) {
	api, c, user := initializeInternalAvatarApiTest(t, "https://foo.bar/")
	actualValue, err := api.GetUserBasicInfo(c, user)

	assert.Equal(t, nil, err)
	assert.Equal(t, user.Username, actualValue.Username)
	assert.Equal(t, string(core.USER_AVATAR_PROVIDER_INTERNAL), actualValue.AvatarProvider)
	assert.True(t, strings.HasPrefix(actualValue.AvatarUrl, "https://foo.bar/avatar/123.png?token="))
}

func TestApiUsingAvatarProvider_GetAvatarUrl_APIToken(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_INTERNAL,
		RootUrl:        "https://foo.bar/",
	})
	c := core.WrapWebContext(&gin.Context{}, nil)
	c.SetTokenClaims(&core.UserTokenClaims{Type: core.USER_TOKEN_TYPE_API})

	expectedValue := "https://foo.bar/avatar/123.png"
	actualValue, err := api.GetAvatarUrl(c, &models.User{Uid: 123, CustomAvatarType: "png"})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestApiUsingAvatarProvider_GetAvatarUrl_Require2FAToken(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_INTERNAL,
		RootUrl:        "https://foo.bar/",
	})
	c := core.WrapWebContext(&gin.Context{}, nil)
	c.SetTokenClaims(&core.UserTokenClaims{Type: core.USER_TOKEN_TYPE_REQUIRE_2FA})

	expectedValue := "https://foo.bar/avatar/123.png"
	actualValue, err := api.GetAvatarUrl(c, &models.User{Uid: 123, CustomAvatarType: "png"})
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestApiWithUserInfo_GetUserBasicInfo_APIToken(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_INTERNAL,
		RootUrl:        "https://foo.bar/",
	})
	c := core.WrapWebContext(&gin.Context{}, nil)
	c.SetTokenClaims(&core.UserTokenClaims{Type: core.USER_TOKEN_TYPE_API})

	user := &models.User{
		Uid:              123,
		Username:         "test-user",
		CustomAvatarType: "png",
	}

	expectedValue := user.ToUserBasicInfo(core.USER_AVATAR_PROVIDER_INTERNAL, "https://foo.bar/avatar/123.png")
	actualValue, err := api.GetUserBasicInfo(c, user)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}
