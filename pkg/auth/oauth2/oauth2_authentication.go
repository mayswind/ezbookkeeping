package oauth2

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// OAuth2Container contains the current OAuth 2.0 authentication provider
type OAuth2Container struct {
	oauth2Config         *oauth2.Config
	oauth2Provider       OAuth2Provider
	oauth2HttpClient     *http.Client
	externalUserAuthType core.UserExternalAuthType
}

// Initialize a OAuth 2.0 container singleton instance
var (
	Container = &OAuth2Container{}
)

// InitializeOAuth2Provider initializes the current OAuth 2.0 provider according to the config
func InitializeOAuth2Provider(config *settings.Config) error {
	if !config.EnableOAuth2Login {
		return nil
	}

	if config.OAuth2ClientID == "" || config.OAuth2ClientSecret == "" || config.OAuth2UserIdentifier == "" || config.OAuth2Provider == "" {
		return errs.ErrInvalidOAuth2Config
	}

	var oauth2Provider OAuth2Provider
	var externalUserAuthType core.UserExternalAuthType

	if config.OAuth2Provider == settings.OAuth2ProviderNextcloud {
		oauth2Provider = NewNextcloudOAuth2Provider(config.OAuth2NextcloudBaseUrl)
		externalUserAuthType = core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD
	} else if config.OAuth2Provider == settings.OAuth2ProviderGithub {
		oauth2Provider = NewGithubOAuth2Provider()
		externalUserAuthType = core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITHUB
	} else {
		return errs.ErrInvalidOAuth2Provider
	}

	Container.oauth2Config = buildOAuth2Config(config, oauth2Provider)
	Container.oauth2Provider = oauth2Provider
	Container.oauth2HttpClient = utils.NewHttpClient(config.OAuth2RequestTimeout, config.OAuth2Proxy, config.OAuth2SkipTLSVerify, settings.GetUserAgent())
	Container.externalUserAuthType = externalUserAuthType

	return nil
}

// GetOAuth2AuthUrl returns the OAuth 2.0 authentication url
func GetOAuth2AuthUrl(c core.Context, state string) (string, error) {
	if Container.oauth2Config == nil {
		return "", errs.ErrOAuth2NotEnabled
	}

	return Container.oauth2Config.AuthCodeURL(state), nil
}

// GetOAuth2Token exchanges the authorization code for an OAuth 2.0 token
func GetOAuth2Token(c core.Context, code string) (*oauth2.Token, error) {
	if Container.oauth2Config == nil || Container.oauth2HttpClient == nil {
		return nil, errs.ErrOAuth2NotEnabled
	}

	return Container.oauth2Config.Exchange(wrapOAuth2Context(c, Container.oauth2HttpClient), code)
}

// GetOAuth2UserInfo retrieves the OAuth 2.0 user info using the provided OAuth 2.0 token
func GetOAuth2UserInfo(c core.Context, token *oauth2.Token) (*OAuth2UserInfo, error) {
	if Container.oauth2Config == nil || Container.oauth2Provider == nil || Container.oauth2HttpClient == nil {
		return nil, errs.ErrOAuth2NotEnabled
	}

	if token == nil {
		return nil, errs.ErrInvalidOAuth2Token
	}

	oauth2Client := oauth2.NewClient(wrapOAuth2Context(c, Container.oauth2HttpClient), oauth2.StaticTokenSource(token))
	return Container.oauth2Provider.GetUserInfo(c, oauth2Client)
}

// GetExternalUserAuthType returns the external user auth type of the current OAuth 2.0 provider
func GetExternalUserAuthType() core.UserExternalAuthType {
	return Container.externalUserAuthType
}

func buildOAuth2Config(config *settings.Config, oauth2Provider OAuth2Provider) *oauth2.Config {
	redirectURL := config.RootUrl + "oauth2/callback"

	return &oauth2.Config{
		ClientID:     config.OAuth2ClientID,
		ClientSecret: config.OAuth2ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth2Provider.GetAuthUrl(),
			TokenURL: oauth2Provider.GetTokenUrl(),
		},
		RedirectURL: redirectURL,
		Scopes:      oauth2Provider.GetScopes(),
	}
}
