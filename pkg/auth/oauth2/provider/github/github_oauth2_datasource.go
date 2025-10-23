package github

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/data"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const githubOAuth2AuthUrl = "https://github.com/login/oauth/authorize"     // Reference: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
const githubOAuth2TokenUrl = "https://github.com/login/oauth/access_token" // Reference: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
const githubUserProfileApiUrl = "https://api.github.com/user"              // Reference: https://docs.github.com/en/rest/users/users
const githubUserEmailApiUrl = "https://api.github.com/user/emails"         // Reference: https://docs.github.com/en/rest/users/emails

var githubOAuth2Scopes = []string{"user:email"}

type githubUserProfileResponse struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type githubUserEmailsResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// GithubOAuth2Provider represents Github OAuth 2.0 provider
type GithubOAuth2Provider struct {
	provider.OAuth2Provider
	oauth2Config *oauth2.Config
}

// GetOAuth2AuthUrl returns the authentication url of the GitHub OAuth 2.0 provider
func (p *GithubOAuth2Provider) GetOAuth2AuthUrl(c core.Context, state string, challenge string) (string, error) {
	return p.oauth2Config.AuthCodeURL(state), nil
}

// GetOAuth2Token returns the OAuth 2.0 token of the GitHub OAuth 2.0 provider
func (p *GithubOAuth2Provider) GetOAuth2Token(c core.Context, code string, verifier string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(c, code)
}

// GetUserInfo returns the user info by the Github OAuth 2.0 provider
func (p *GithubOAuth2Provider) GetUserInfo(c core.Context, oauth2Token *oauth2.Token) (*data.OAuth2UserInfo, error) {
	// first get user name and nick name from user profile
	req, err := p.buildAPIRequest(githubUserProfileApiUrl)

	if err != nil {
		log.Errorf(c, "[github_oauth2_datasource_test.GetUserInfo] failed to get user info request, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	oauth2Client := oauth2.NewClient(c, oauth2.StaticTokenSource(oauth2Token))
	resp, err := oauth2Client.Do(req)

	if err != nil {
		log.Errorf(c, "[github_oauth2_datasource_test.GetUserInfo] failed to get user info response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	log.Debugf(c, "[github_oauth2_datasource_test.GetUserInfo] user profile response is %s", body)

	if resp.StatusCode != 200 {
		log.Errorf(c, "[github_oauth2_datasource_test.GetUserInfo] failed to get user info response, because response code is %d", resp.StatusCode)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	userProfileResp, err := p.parseUserProfile(c, body)

	if err != nil {
		return nil, err
	}

	// then get user primary email
	req, err = p.buildAPIRequest(githubUserEmailApiUrl)

	if err != nil {
		log.Errorf(c, "[github_oauth2_datasource_test.GetUserInfo] failed to get user emails request, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	resp, err = oauth2Client.Do(req)

	if err != nil {
		log.Errorf(c, "[github_oauth2_datasource_test.GetUserInfo] failed to get user emails response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	log.Debugf(c, "[github_oauth2_datasource_test.GetUserInfo] user emails response is %s", body)

	if resp.StatusCode != 200 {
		log.Errorf(c, "[github_oauth2_datasource_test.GetUserInfo] failed to get user emails response, because response code is %d", resp.StatusCode)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	email, err := p.parsePrimaryEmail(c, body)

	if err != nil {
		return nil, err
	}

	return &data.OAuth2UserInfo{
		UserName: userProfileResp.Login,
		Email:    email,
		NickName: userProfileResp.Name,
	}, nil
}

func (p *GithubOAuth2Provider) parseUserProfile(c core.Context, body []byte) (*githubUserProfileResponse, error) {
	userProfileResp := &githubUserProfileResponse{}
	err := json.Unmarshal(body, &userProfileResp)

	if err != nil {
		log.Warnf(c, "[github_oauth2_datasource.parseUserProfile] failed to parse user profile response body, because %s", err.Error())
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	if userProfileResp.Login == "" {
		log.Warnf(c, "[github_oauth2_datasource.parseUserProfile] invalid user profile response body")
		return nil, errs.ErrCannotRetrieveUserInfo
	}

	return userProfileResp, nil
}

func (p *GithubOAuth2Provider) parsePrimaryEmail(c core.Context, body []byte) (string, error) {
	emailsResp := make([]githubUserEmailsResponse, 0)
	err := json.Unmarshal(body, &emailsResp)

	if err != nil {
		log.Warnf(c, "[github_oauth2_datasource.parsePrimaryEmail] failed to parse user emails response body, because %s", err.Error())
		return "", errs.ErrCannotRetrieveUserInfo
	}

	for _, emailEntry := range emailsResp {
		if emailEntry.Primary && emailEntry.Verified {
			return emailEntry.Email, nil
		}
	}

	return "", nil
}

func (p *GithubOAuth2Provider) buildAPIRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	return req, nil
}

// NewGithubOAuth2Provider creates a new Github OAuth 2.0 provider instance
func NewGithubOAuth2Provider(config *settings.Config, redirectUrl string) (provider.OAuth2Provider, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     config.OAuth2ClientID,
		ClientSecret: config.OAuth2ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubOAuth2AuthUrl,
			TokenURL: githubOAuth2TokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      githubOAuth2Scopes,
	}

	return &GithubOAuth2Provider{
		oauth2Config: oauth2Config,
	}, nil
}
