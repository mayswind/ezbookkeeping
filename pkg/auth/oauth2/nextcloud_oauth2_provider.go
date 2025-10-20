package oauth2

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
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

// NextcloudOAuth2Provider represents Nextcloud OAuth 2.0 provider
type NextcloudOAuth2Provider struct {
	baseUrl string
}

// NewNextcloudOAuth2Provider creates a new Nextcloud OAuth 2.0 provider instance
func NewNextcloudOAuth2Provider(baseUrl string) OAuth2Provider {
	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl += "/"
	}

	return &NextcloudOAuth2Provider{
		baseUrl: baseUrl,
	}
}

// GetAuthUrl returns the authentication url of the Nextcloud provider
func (p *NextcloudOAuth2Provider) GetAuthUrl() string {
	return p.baseUrl + "apps/oauth2/authorize"
}

// GetTokenUrl returns the token url of the Nextcloud provider
func (p *NextcloudOAuth2Provider) GetTokenUrl() string {
	return p.baseUrl + "apps/oauth2/api/v1/token"
}

// GetUserInfo returns the user info by the Nextcloud provider
func (p *NextcloudOAuth2Provider) GetUserInfo(c core.Context, oauth2Client *http.Client) (*OAuth2UserInfo, error) {
	url := p.baseUrl + "ocs/v2.php/cloud/user?format=json"
	resp, err := oauth2Client.Get(url)

	if err != nil {
		log.Errorf(c, "[nextcloud_oauth2_provider.GetUserInfo] failed to get user info response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	log.Debugf(c, "[nextcloud_oauth2_provider.GetUserInfo] response is %s", body)

	if resp.StatusCode != 200 {
		log.Errorf(c, "[nextcloud_oauth2_provider.GetUserInfo] failed to get user info response, because response code is %d", resp.StatusCode)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return p.parseUserInfo(c, body)
}

// GetScopes returns the scopes required by the Nextcloud provider
func (p *NextcloudOAuth2Provider) GetScopes() []string {
	return []string{"profile", "email"}
}

func (p *NextcloudOAuth2Provider) parseUserInfo(c core.Context, body []byte) (*OAuth2UserInfo, error) {
	userInfoResp := &nextcloudUserInfoResponse{}
	err := json.Unmarshal(body, &userInfoResp)

	if err != nil {
		log.Warnf(c, "[nextcloud_oauth2_provider.parseUserInfo] failed to parse user info response body, because %s", err.Error())
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.OCS == nil || userInfoResp.OCS.Meta == nil || userInfoResp.OCS.Data == nil {
		log.Warnf(c, "[nextcloud_oauth2_provider.parseUserInfo] invalid user info response body")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.OCS.Meta.StatusCode != 200 {
		log.Warnf(c, "[nextcloud_oauth2_provider.parseUserInfo] user info response status code is %d", userInfoResp.OCS.Meta.StatusCode)
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.OCS.Data.ID == "" {
		log.Warnf(c, "[nextcloud_oauth2_provider.parseUserInfo] user info id is empty")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	return &OAuth2UserInfo{
		UserName: userInfoResp.OCS.Data.ID,
		Email:    userInfoResp.OCS.Data.Email,
		NickName: userInfoResp.OCS.Data.DisplayName,
	}, nil
}
