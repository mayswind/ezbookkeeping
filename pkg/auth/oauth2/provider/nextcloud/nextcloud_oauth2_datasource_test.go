package nextcloud

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestNewNextcloudOAuth2Provider(t *testing.T) {
	provider, err := NewNextcloudOAuth2Provider(&settings.Config{
		OAuth2NextcloudBaseUrl: "https://example.com/",
	}, "")
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com/apps/oauth2/authorize", provider.(*common.CommonOAuth2Provider).GetDataSource().GetAuthUrl())
	assert.Equal(t, "https://example.com/apps/oauth2/api/v1/token", provider.(*common.CommonOAuth2Provider).GetDataSource().GetTokenUrl())

	provider, err = NewNextcloudOAuth2Provider(&settings.Config{
		OAuth2NextcloudBaseUrl: "https://example.com/index.php",
	}, "")
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com/index.php/apps/oauth2/authorize", provider.(*common.CommonOAuth2Provider).GetDataSource().GetAuthUrl())
	assert.Equal(t, "https://example.com/index.php/apps/oauth2/api/v1/token", provider.(*common.CommonOAuth2Provider).GetDataSource().GetTokenUrl())

	provider, err = NewNextcloudOAuth2Provider(&settings.Config{}, "")
	assert.Equal(t, errs.ErrInvalidOAuth2Config, err)
}

func TestNextcloudOAuth2Datasource_GetUserInfoRequest(t *testing.T) {
	datasource := &NextcloudOAuth2DataSource{baseUrl: "https://example.com/"}
	req, err := datasource.GetUserInfoRequest()

	assert.Nil(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://example.com/ocs/v2.php/cloud/user", req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Accept"))
	assert.Equal(t, "true", req.Header.Get("OCS-APIRequest"))
}

func TestNextcloudOAuth2Datasource_ParseUserInfo_Success(t *testing.T) {
	datasource := &NextcloudOAuth2DataSource{}
	responseContent := `{
		"ocs": {
			"meta": {
				"status": "ok",
				"statuscode": 200
			},
			"data": {
				"id": "user1",
				"email": "user1@example.com",
				"display-name": "User"
			}
		}
	}`
	info, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent), &http.Client{})

	assert.Nil(t, err)
	assert.Equal(t, "user1", info.UserName)
	assert.Equal(t, "user1@example.com", info.Email)
	assert.Equal(t, "User", info.NickName)
}

func TestNextcloudOAuth2Datasource_ParseUserInfo_InvalidJson(t *testing.T) {
	datasource := &NextcloudOAuth2DataSource{}
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte("invalid"), &http.Client{})

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}

func TestNextcloudOAuth2Datasource_ParseUserInfo_MissingFields(t *testing.T) {
	datasource := &NextcloudOAuth2DataSource{}
	responseContent := `{"ocs": {}}`
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent), &http.Client{})

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}

func TestNextcloudOAuth2Datasource_ParseUserInfo_Non200StatusCode(t *testing.T) {
	datasource := &NextcloudOAuth2DataSource{}
	responseContent := `{
		"ocs": {
			"meta": {
				"status": "error",
				"statuscode": 400
			},
			"data": {}
		}
	}`
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent), &http.Client{})

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}

func TestNextcloudOAuth2Datasource_ParseUserInfo_EmptyID(t *testing.T) {
	datasource := &NextcloudOAuth2DataSource{}
	responseContent := `{
		"ocs": {
			"meta": {
				"status": "ok",
				"statuscode": 200
			},
			"data": {
				"id": "",
				"email": "user1@example.com",
				"display-name": "User One"
			}
		}
	}`
	_, err := datasource.ParseUserInfo(core.NewNullContext(), []byte(responseContent), &http.Client{})

	assert.Equal(t, errs.ErrCannotRetrieveUserInfo, err)
}
