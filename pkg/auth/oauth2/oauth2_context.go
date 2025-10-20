package oauth2

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// OAuth2Context represents the context for OAuth 2.0 operations
type OAuth2Context struct {
	core.Context
	httpClient *http.Client
}

// Value returns the value associated with key
func (o *OAuth2Context) Value(key any) any {
	if key == oauth2.HTTPClient {
		return o.httpClient
	}

	return o.Context.Value(key)
}

func wrapOAuth2Context(ctx core.Context, httpClient *http.Client) core.Context {
	return &OAuth2Context{
		Context:    ctx,
		httpClient: httpClient,
	}
}
