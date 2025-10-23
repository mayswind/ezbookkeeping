package nextcloud

import (
	"encoding/json"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/data"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/common"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

type nextcloudUserInfoResponse struct {
	OCS *struct {
		Meta *struct {
			Status     string `json:"status"`
			StatusCode int    `json:"statuscode"`
		} `json:"meta"`
		Data *struct {
			ID          string `json:"id"`
			Email       string `json:"email"`
			DisplayName string `json:"display-name"`
		} `json:"data"`
	} `json:"ocs"`
}

// NextcloudOAuth2DataSource represents Nextcloud OAuth 2.0 data source
type NextcloudOAuth2DataSource struct {
	common.CommonOAuth2DataSource
	baseUrl string
}

// GetAuthUrl returns the authentication url of the Nextcloud data source
func (s *NextcloudOAuth2DataSource) GetAuthUrl() string {
	// Reference: https://docs.nextcloud.com/server/stable/developer_manual/_static/openapi.html#/operations/oauth2-login_redirector-authorize
	return s.baseUrl + "apps/oauth2/authorize"
}

// GetTokenUrl returns the token url of the Nextcloud data source
func (s *NextcloudOAuth2DataSource) GetTokenUrl() string {
	// Reference: https://docs.nextcloud.com/server/stable/developer_manual/_static/openapi.html#/operations/oauth2-oauth_api-get-token
	return s.baseUrl + "apps/oauth2/api/v1/token"
}

// GetUserInfoRequest returns the user info request of the Nextcloud data source
func (s *NextcloudOAuth2DataSource) GetUserInfoRequest() (*http.Request, error) {
	// Reference: https://docs.nextcloud.com/server/stable/developer_manual/_static/openapi.html#/operations/provisioning_api-users-get-current-user
	req, err := http.NewRequest("GET", s.baseUrl+"ocs/v2.php/cloud/user", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("OCS-APIRequest", "true")
	return req, nil
}

// GetScopes returns the scopes required by the Nextcloud provider
func (s *NextcloudOAuth2DataSource) GetScopes() []string {
	return []string{}
}

// ParseUserInfo returns the user info by parsing the response body
func (s *NextcloudOAuth2DataSource) ParseUserInfo(c core.Context, body []byte, oauth2Client *http.Client) (*data.OAuth2UserInfo, error) {
	userInfoResp := &nextcloudUserInfoResponse{}
	err := json.Unmarshal(body, &userInfoResp)

	if err != nil {
		log.Warnf(c, "[nextcloud_oauth2_datasource.ParseUserInfo] failed to parse user info response body, because %s", err.Error())
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.OCS == nil || userInfoResp.OCS.Meta == nil || userInfoResp.OCS.Data == nil {
		log.Warnf(c, "[nextcloud_oauth2_datasource.ParseUserInfo] invalid user info response body")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.OCS.Meta.StatusCode != 200 {
		log.Warnf(c, "[nextcloud_oauth2_datasource.ParseUserInfo] user info response status code is %d", userInfoResp.OCS.Meta.StatusCode)
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.OCS.Data.ID == "" {
		log.Warnf(c, "[nextcloud_oauth2_datasource.ParseUserInfo] user info id is empty")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	return &data.OAuth2UserInfo{
		UserName: userInfoResp.OCS.Data.ID,
		Email:    userInfoResp.OCS.Data.Email,
		NickName: userInfoResp.OCS.Data.DisplayName,
	}, nil
}

// NewNextcloudOAuth2Provider creates a new Nextcloud OAuth 2.0 provider instance
func NewNextcloudOAuth2Provider(config *settings.Config, redirectUrl string) (provider.OAuth2Provider, error) {
	if len(config.OAuth2NextcloudBaseUrl) < 1 {
		return nil, errs.ErrInvalidOAuth2Config
	}

	baseUrl := config.OAuth2NextcloudBaseUrl

	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl += "/"
	}

	return common.NewCommonOAuth2Provider(config, redirectUrl, &NextcloudOAuth2DataSource{
		baseUrl: baseUrl,
	}), nil
}
