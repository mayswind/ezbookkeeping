package provider

import (
	"golang.org/x/oauth2"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/data"
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// OAuth2Provider defines the structure of OAuth 2.0 provider
type OAuth2Provider interface {
	// GetOAuth2AuthUrl returns the authentication url of the provider
	GetOAuth2AuthUrl(c core.Context, state string, opts ...oauth2.AuthCodeOption) (string, error)

	// GetOAuth2Token returns the OAuth 2.0 token of the provider
	GetOAuth2Token(c core.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)

	// GetUserInfo returns the user info
	GetUserInfo(c core.Context, oauth2Token *oauth2.Token) (*data.OAuth2UserInfo, error)
}
