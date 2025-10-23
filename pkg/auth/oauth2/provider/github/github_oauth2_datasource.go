package github

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

type githubUserProfileResponse struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GithubOAuth2DataSource represents Github OAuth 2.0 data source
type GithubOAuth2DataSource struct {
	common.CommonOAuth2DataSource
}

// GetAuthUrl returns the authentication url of the Github data source
func (s *GithubOAuth2DataSource) GetAuthUrl() string {
	// Reference: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
	return "https://github.com/login/oauth/authorize"
}

// GetTokenUrl returns the token url of the Github data source
func (s *GithubOAuth2DataSource) GetTokenUrl() string {
	// Reference: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
	return "https://github.com/login/oauth/access_token"
}

// GetUserInfoRequest returns the user info request of the Github data source
func (s *GithubOAuth2DataSource) GetUserInfoRequest() (*http.Request, error) {
	// Reference: https://docs.github.com/en/rest/users/users
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	return req, nil
}

// GetScopes returns the scopes required by the Github provider
func (p *GithubOAuth2DataSource) GetScopes() []string {
	return []string{"read:user"}
}

// ParseUserInfo returns the user info by parsing the response body
func (p *GithubOAuth2DataSource) ParseUserInfo(c core.Context, body []byte) (*data.OAuth2UserInfo, error) {
	userInfoResp := &githubUserProfileResponse{}
	err := json.Unmarshal(body, &userInfoResp)

	if err != nil {
		log.Warnf(c, "[github_oauth2_datasource.ParseUserInfo] failed to parse user profile response body, because %s", err.Error())
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userInfoResp.Login == "" {
		log.Warnf(c, "[github_oauth2_datasource.ParseUserInfo] invalid user profile response body")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	return &data.OAuth2UserInfo{
		UserName: userInfoResp.Login,
		Email:    userInfoResp.Email,
		NickName: userInfoResp.Name,
	}, nil
}

// NewGithubOAuth2Provider creates a new Github OAuth 2.0 provider instance
func NewGithubOAuth2Provider(config *settings.Config, redirectUrl string) (provider.OAuth2Provider, error) {
	return common.NewCommonOAuth2Provider(config, redirectUrl, &GithubOAuth2DataSource{}), nil
}
