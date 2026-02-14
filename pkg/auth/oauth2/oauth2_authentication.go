package oauth2

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/data"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/gitea"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/github"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/nextcloud"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider/oidc"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/httpclient"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OAuth2Container contains the current OAuth 2.0 authentication provider
type OAuth2Container struct {
	current              provider.OAuth2Provider
	usePKCE              bool
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

	var err error
	var oauth2Provider provider.OAuth2Provider
	var externalUserAuthType core.UserExternalAuthType
	redirectUrl := config.RootUrl + "oauth2/callback"

	if config.OAuth2Provider == settings.OAuth2ProviderOIDC {
		oauth2Provider, err = oidc.NewOIDCProvider(config, redirectUrl)
		externalUserAuthType = core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_OIDC
	} else if config.OAuth2Provider == settings.OAuth2ProviderNextcloud {
		oauth2Provider, err = nextcloud.NewNextcloudOAuth2Provider(config, redirectUrl)
		externalUserAuthType = core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_NEXTCLOUD
	} else if config.OAuth2Provider == settings.OAuth2ProviderGitea {
		oauth2Provider, err = gitea.NewGiteaOAuth2Provider(config, redirectUrl)
		externalUserAuthType = core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITEA
	} else if config.OAuth2Provider == settings.OAuth2ProviderGithub {
		oauth2Provider, err = github.NewGithubOAuth2Provider(config, redirectUrl)
		externalUserAuthType = core.USER_EXTERNAL_AUTH_TYPE_OAUTH2_GITHUB
	} else {
		return errs.ErrInvalidOAuth2Provider
	}

	if err != nil {
		return err
	}

	Container.current = oauth2Provider
	Container.usePKCE = config.OAuth2UsePKCE
	Container.oauth2HttpClient = httpclient.NewHttpClient(config.OAuth2RequestTimeout, config.OAuth2Proxy, config.OAuth2SkipTLSVerify, core.GetOutgoingUserAgent(), config.EnableDebugLog)
	Container.externalUserAuthType = externalUserAuthType

	return nil
}

// GetOAuth2AuthUrl returns the OAuth 2.0 authentication url
func GetOAuth2AuthUrl(c core.Context, state string, verifier string) (string, error) {
	if Container.current == nil {
		return "", errs.ErrOAuth2NotEnabled
	}

	var opts []oauth2.AuthCodeOption

	if Container.usePKCE {
		opts = append(opts, oauth2.S256ChallengeOption(verifier))
	}

	return Container.current.GetOAuth2AuthUrl(wrapOAuth2Context(c, Container.oauth2HttpClient), state, opts...)
}

// GetOAuth2Token exchanges the authorization code for an OAuth 2.0 token
func GetOAuth2Token(c core.Context, code string, verifier string) (*oauth2.Token, error) {
	if Container.current == nil || Container.oauth2HttpClient == nil {
		return nil, errs.ErrOAuth2NotEnabled
	}

	var opts []oauth2.AuthCodeOption

	if Container.usePKCE {
		opts = append(opts, oauth2.VerifierOption(verifier))
	}

	return Container.current.GetOAuth2Token(wrapOAuth2Context(c, Container.oauth2HttpClient), code, opts...)
}

// GetOAuth2UserInfo retrieves the OAuth 2.0 user info using the provided OAuth 2.0 token
func GetOAuth2UserInfo(c core.Context, token *oauth2.Token) (*data.OAuth2UserInfo, error) {
	if Container.current == nil || Container.oauth2HttpClient == nil {
		return nil, errs.ErrOAuth2NotEnabled
	}

	if token == nil {
		return nil, errs.ErrInvalidOAuth2Token
	}

	return Container.current.GetUserInfo(wrapOAuth2Context(c, Container.oauth2HttpClient), token)
}

// GetExternalUserAuthType returns the external user auth type of the current OAuth 2.0 provider
func GetExternalUserAuthType() core.UserExternalAuthType {
	return Container.externalUserAuthType
}
