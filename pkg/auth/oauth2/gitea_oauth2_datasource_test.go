package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestNewGiteaOAuth2Provider(t *testing.T) {
	datasource := NewGiteaOAuth2Provider("https://example.com/")
	assert.Equal(t, "https://example.com/login/oauth/authorize", datasource.GetAuthUrl())
	assert.Equal(t, "https://example.com/login/oauth/access_token", datasource.GetTokenUrl())

	datasource = NewGiteaOAuth2Provider("https://example.com")
	assert.Equal(t, "https://example.com/login/oauth/authorize", datasource.GetAuthUrl())
	assert.Equal(t, "https://example.com/login/oauth/access_token", datasource.GetTokenUrl())
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
