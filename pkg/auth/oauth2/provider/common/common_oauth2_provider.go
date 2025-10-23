package common

import (
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

// CommonOAuth2Provider represents common OAuth 2.0 provider
type CommonOAuth2Provider struct {
	provider.OAuth2Provider
	oauth2Config *oauth2.Config
	dataSource   CommonOAuth2DataSource
}

// CommonOAuth2DataSource defines the structure of OAuth 2.0 data source
type CommonOAuth2DataSource interface {
	// GetAuthUrl returns the authentication url of the data source
	GetAuthUrl() string

	// GetTokenUrl returns the token url of the data source
	GetTokenUrl() string

	// GetUserInfoRequest returns the user info request of the data source
	GetUserInfoRequest() (*http.Request, error)

	// GetScopes returns the scopes required by the data source
	GetScopes() []string

	// ParseUserInfo returns the user info by parsing the response body
	ParseUserInfo(c core.Context, body []byte, oauth2Client *http.Client) (*data.OAuth2UserInfo, error)
}

// GetOAuth2AuthUrl returns the authentication url of the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetOAuth2AuthUrl(c core.Context, state string, challenge string) (string, error) {
	return p.oauth2Config.AuthCodeURL(state), nil
}

// GetOAuth2Token returns the OAuth 2.0 token of the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetOAuth2Token(c core.Context, code string, verifier string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(c, code)
}

// GetUserInfo returns the user info by the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetUserInfo(c core.Context, oauth2Token *oauth2.Token) (*data.OAuth2UserInfo, error) {
	req, err := p.dataSource.GetUserInfoRequest()

	if err != nil {
		log.Errorf(c, "[common_oauth2_provider.GetUserInfo] failed to get user info request, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	oauth2Client := oauth2.NewClient(c, oauth2.StaticTokenSource(oauth2Token))
	resp, err := oauth2Client.Do(req)

	if err != nil {
		log.Errorf(c, "[common_oauth2_provider.GetUserInfo] failed to get user info response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	log.Debugf(c, "[common_oauth2_provider.GetUserInfo] response is %s", body)

	if resp.StatusCode != 200 {
		log.Errorf(c, "[common_oauth2_provider.GetUserInfo] failed to get user info response, because response code is %d", resp.StatusCode)
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return p.dataSource.ParseUserInfo(c, body, oauth2Client)
}

// GetDataSource returns the data source of the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetDataSource() CommonOAuth2DataSource {
	return p.dataSource
}

// NewCommonOAuth2Provider returns a new common OAuth 2.0 provider
func NewCommonOAuth2Provider(config *settings.Config, redirectUrl string, dataSource CommonOAuth2DataSource) *CommonOAuth2Provider {
	oauth2Config := &oauth2.Config{
		ClientID:     config.OAuth2ClientID,
		ClientSecret: config.OAuth2ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  dataSource.GetAuthUrl(),
			TokenURL: dataSource.GetTokenUrl(),
		},
		RedirectURL: redirectUrl,
		Scopes:      dataSource.GetScopes(),
	}

	return &CommonOAuth2Provider{
		oauth2Config: oauth2Config,
		dataSource:   dataSource,
	}
}
