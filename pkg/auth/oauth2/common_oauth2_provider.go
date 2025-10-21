package oauth2

import (
	"io"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// CommonOAuth2Provider represents common OAuth 2.0 provider
type CommonOAuth2Provider struct {
	OAuth2Provider
	dataSource CommonOAuth2DataSource
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
	ParseUserInfo(c core.Context, body []byte) (*OAuth2UserInfo, error)
}

// GetAuthUrl returns the authentication url of the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetAuthUrl() string {
	return p.dataSource.GetAuthUrl()
}

// GetTokenUrl returns the token url of the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetTokenUrl() string {
	return p.dataSource.GetTokenUrl()
}

// GetScopes returns the scopes required by the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetScopes() []string {
	return p.dataSource.GetScopes()
}

// GetUserInfo returns the user info by the common OAuth 2.0 provider
func (p *CommonOAuth2Provider) GetUserInfo(c core.Context, oauth2Client *http.Client) (*OAuth2UserInfo, error) {
	req, err := p.dataSource.GetUserInfoRequest()

	if err != nil {
		log.Errorf(c, "[common_oauth2_provider.GetUserInfo] failed to get user info request, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

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

	return p.dataSource.ParseUserInfo(c, body)
}
