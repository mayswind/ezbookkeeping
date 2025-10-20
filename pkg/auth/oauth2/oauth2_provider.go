package oauth2

import (
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// OAuth2Provider defines the structure of OAuth 2.0 provider
type OAuth2Provider interface {
	// GetAuthUrl returns the authentication url of the provider
	GetAuthUrl() string

	// GetTokenUrl returns the token url of the provider
	GetTokenUrl() string

	// GetUserInfo returns the user info
	GetUserInfo(c core.Context, oauth2Client *http.Client) (*OAuth2UserInfo, error)

	// GetScopes returns the scopes required by the provider
	GetScopes() []string
}
