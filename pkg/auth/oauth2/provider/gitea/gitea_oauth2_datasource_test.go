package gitea

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestNewGiteaOAuth2Provider(t *testing.T) {
	provider, err := NewGiteaOAuth2Provider(&settings.Config{
		OAuth2GiteaBaseUrl: "https://example.com/",
	}, "")
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com/login/oauth/authorize", provider.(*common.CommonOAuth2Provider).GetDataSource().GetAuthUrl())
	assert.Equal(t, "https://example.com/login/oauth/access_token", provider.(*common.CommonOAuth2Provider).GetDataSource().GetTokenUrl())

	provider, err = NewGiteaOAuth2Provider(&settings.Config{
		OAuth2GiteaBaseUrl: "https://example.com",
	}, "")
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com/login/oauth/authorize", provider.(*common.CommonOAuth2Provider).GetDataSource().GetAuthUrl())
	assert.Equal(t, "https://example.com/login/oauth/access_token", provider.(*common.CommonOAuth2Provider).GetDataSource().GetTokenUrl())

	provider, err = NewGiteaOAuth2Provider(&settings.Config{}, "")
	assert.Equal(t, errs.ErrInvalidOAuth2Config, err)
}

func TestGiteaOAuth2Datasource_GetUserInfoRequest(t *testing.T) {
	datasource := &GiteaOAuth2DataSource{baseUrl: "https://example.com/"}
	req, err := datasource.GetUserInfoRequest()

	assert.Nil(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://example.com/api/v1/user", req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Accept"))
}

func TestGiteaOAuth2Datasource_ParseUserInfo_Success(t *testing.T) {
	datasource := &GiteaOAuth2DataSource{}
	responseContent := `{
		"login": "user1",
		"full_name": "User",
		"email": "user1@example.com"
	}`
	info, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent))

	assert.Nil(t, err)
	assert.Equal(t, "user1", info.UserName)
	assert.Equal(t, "user1@example.com", info.Email)
	assert.Equal(t, "User", info.NickName)
}

func TestGiteaOAuth2Datasource_ParseUserInfo_InvalidJson(t *testing.T) {
	datasource := &GiteaOAuth2DataSource{}
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte("invalid"))

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}

func TestGiteaOAuth2Datasource_ParseUserInfo_EmptyLogin(t *testing.T) {
	datasource := &GiteaOAuth2DataSource{}
	responseContent := `{"login": ""}`
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent))

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}
