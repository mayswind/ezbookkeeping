package oidc

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/data"
	"github.com/mayswind/ezbookkeeping/pkg/auth/oauth2/provider"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/httpclient"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// OIDCClaims represents OIDC claims
type OIDCClaims struct {
	PreferredUserName string `json:"preferred_username"`
	UserName          string `json:"username"`
	Name              string `json:"name"`
	Email             string `json:"email"`
}

// OIDCProvider represents OIDC provider
type OIDCProvider struct {
	provider.OAuth2Provider
	oidcIssuerURL      string
	oidcCheckIssuerURL bool
	redirectUrl        string
	oauth2ClientID     string
	oauth2ClientSecret string
	oauth2Config       *oauth2.Config
	oidcProvider       *oidc.Provider
	oidcVerifier       *oidc.IDTokenVerifier
}

// GetOAuth2AuthUrl returns the authentication url of the OIDC provider
func (p *OIDCProvider) GetOAuth2AuthUrl(c core.Context, state string, opts ...oauth2.AuthCodeOption) (string, error) {
	oauth2Config, err := p.getOAuth2Config(c)

	if err != nil {
		return "", err
	}

	return oauth2Config.AuthCodeURL(state, opts...), nil
}

// GetOAuth2Token returns the OAuth 2.0 token of the OIDC provider
func (p *OIDCProvider) GetOAuth2Token(c core.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	oauth2Config, err := p.getOAuth2Config(c)

	if err != nil {
		return nil, err
	}

	return oauth2Config.Exchange(c, code, opts...)
}

// GetUserInfo returns the user info by the OIDC provider
func (p *OIDCProvider) GetUserInfo(c core.Context, oauth2Token *oauth2.Token) (*data.OAuth2UserInfo, error) {
	_, err := p.getOAuth2Config(c)

	if err != nil {
		return nil, err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)

	if !ok {
		log.Errorf(c, "[oidc_provider.GetUserInfo] missing \"id_token\" field in oauth 2.0 token")
		return nil, errs.ErrInvalidOAuth2Token
	}

	idToken, err := p.oidcVerifier.Verify(c, rawIDToken)

	if err != nil {
		log.Errorf(c, "[oidc_provider.GetUserInfo] failed to verify \"id_token\" field in oauth 2.0 token, because %s", err.Error())
		return nil, errs.ErrInvalidOAuth2Token
	}

	var claims OIDCClaims
	err = idToken.Claims(&claims)

	if err != nil {
		log.Errorf(c, "[oidc_provider.GetUserInfo] failed to parse claims in oauth 2.0 token, because %s", err.Error())
		return nil, errs.ErrInvalidOAuth2Token
	}

	userName := claims.PreferredUserName
	email := claims.Email
	nickName := claims.Name

	if userName == "" || email == "" || nickName == "" {
		userInfo, err := p.oidcProvider.UserInfo(httpclient.CustomHttpResponseLog(c, func(data []byte) {
			log.Debugf(c, "[oidc_provider.GetUserInfo] response is %s", data)
		}), oauth2.StaticTokenSource(oauth2Token))

		if err != nil {
			log.Errorf(c, "[oidc_provider.GetUserInfo] failed to get user info, because %s", err.Error())
			return nil, errs.ErrCannotRetrieveUserInfo
		}

		err = userInfo.Claims(&claims)

		if err != nil {
			log.Errorf(c, "[oidc_provider.GetUserInfo] failed to parse user info, because %s", err.Error())
			return nil, errs.ErrCannotRetrieveUserInfo
		}

		if userName == "" {
			userName = claims.PreferredUserName
		}

		if userName == "" {
			userName = claims.UserName
		}

		if email == "" {
			email = claims.Email
		}

		if nickName == "" {
			nickName = claims.Name
		}
	}

	return &data.OAuth2UserInfo{
		UserName: userName,
		Email:    email,
		NickName: nickName,
	}, nil
}

func (p *OIDCProvider) getOAuth2Config(c core.Context) (*oauth2.Config, error) {
	if p.oauth2Config != nil {
		return p.oauth2Config, nil
	}

	var ctx context.Context = c

	if !p.oidcCheckIssuerURL {
		ctx = oidc.InsecureIssuerURLContext(c, p.oidcIssuerURL)
	}

	oidcProvider, err := oidc.NewProvider(ctx, p.oidcIssuerURL)

	if err != nil {
		log.Errorf(c, "[oidc_provider.getOAuth2Config] failed to create oidc provider, because %s", err.Error())
		return nil, err
	}

	oidcVerifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:        p.oauth2ClientID,
		SkipIssuerCheck: !p.oidcCheckIssuerURL,
	})

	oauth2Config := &oauth2.Config{
		ClientID:     p.oauth2ClientID,
		ClientSecret: p.oauth2ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  p.redirectUrl,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	p.oauth2Config = oauth2Config
	p.oidcProvider = oidcProvider
	p.oidcVerifier = oidcVerifier
	return oauth2Config, nil
}

// NewOIDCProvider returns a new OIDC provider
func NewOIDCProvider(config *settings.Config, redirectUrl string) (*OIDCProvider, error) {
	if len(config.OAuth2OIDCProviderIssuerURL) < 1 {
		return nil, errs.ErrInvalidOAuth2Config
	}

	return &OIDCProvider{
		oidcIssuerURL:      config.OAuth2OIDCProviderIssuerURL,
		oidcCheckIssuerURL: config.OAuth2OIDCProviderCheckIssuerURL,
		redirectUrl:        redirectUrl,
		oauth2ClientID:     config.OAuth2ClientID,
		oauth2ClientSecret: config.OAuth2ClientSecret,
		oauth2Config:       nil,
	}, nil
}
