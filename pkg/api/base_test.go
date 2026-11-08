package api

import (
	"net/url"
	"path/filepath"
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

	previousStores := *datastore.Container

	t.Cleanup(func() {
		*datastore.Container = previousStores
	})

	err := datastore.InitializeDataStore(&settings.Config{
		DatabaseConfig: &settings.DatabaseConfig{
			DatabaseType: settings.Sqlite3DbType,
			DatabasePath: filepath.Join(t.TempDir(), "tokens.db"),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = datastore.Container.TokenStore.SyncStructs(new(models.TokenRecord))

	if err != nil {
		t.Fatal(err)
	}

	c := core.WrapWebContext(&gin.Context{}, nil)

	user := &models.User{
		Uid:              123,
		Username:         "test-user",
		CustomAvatarType: "png",
	}

	_, claims, tokenRecord, err := services.Tokens.CreateToken(c, user)

	if err != nil {
		t.Fatal(err)
	}

	c.SetTokenClaims(claims)
	c.SetTokenRecord(tokenRecord)
	api.ApiUsingAvatarProvider.tokens = services.Tokens

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

	claims, _, err := services.Tokens.ParseToken(c, avatarUrl.Query().Get("token"), core.NormalTokenScopeUserAvatar)
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

	claims, _, err := services.Tokens.ParseToken(c, avatarUrl.Query().Get("token"), core.NormalTokenScopeUserAvatar)
	assert.Equal(t, nil, err)
	assert.Equal(t, user.Uid, claims.Uid)
}

func TestApiWithUserInfo_GetUserBasicInfo_AvatarTokenError(t *testing.T) {
	api := initializeAvatarProviderApiTest(t, &settings.Config{
		AvatarProvider: core.USER_AVATAR_PROVIDER_INTERNAL,
	})
	api.ApiUsingAvatarProvider.tokens = services.Tokens

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
