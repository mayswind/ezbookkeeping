package core

import (
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const webContextRequestIdFieldKey = "REQUEST_ID"
const webContextTextualTokenFieldKey = "TOKEN_STRING"
const webContextTokenClaimsFieldKey = "TOKEN_CLAIMS"
const webContextResponseErrorFieldKey = "RESPONSE_ERROR"

// AcceptLanguageHeaderName represents the header name of accept language
const AcceptLanguageHeaderName = "Accept-Language"

// RemoteClientPortHeader represents the header name of remote client source port
const RemoteClientPortHeader = "X-Real-Port"

// ClientTimezoneOffsetHeaderName represents the header name of client timezone offset
const ClientTimezoneOffsetHeaderName = "X-Timezone-Offset"

const tokenHeaderName = "Authorization"
const tokenHeaderValuePrefix = "bearer "
const tokenQueryStringParam = "token"
const tokenCookieParam = "ebk_auth_token"

// WebContext represents the request and response context
type WebContext struct {
	*gin.Context
	// DO NOT ADD ANY FIELD IN THIS CONTEXT, THIS CONTEXT IS JUST A WRAPPER
}

func (c *WebContext) ClientPort() uint16 {
	remotePort := c.GetHeader(RemoteClientPortHeader)

	if remotePort != "" {
		remotePortNum, err := strconv.ParseInt(remotePort, 10, 32)

		if err == nil {
			return uint16(remotePortNum)
		}
	}

	if c.Request == nil {
		return 0
	}

	_, remotePort, err := net.SplitHostPort(c.Request.RemoteAddr)

	if err != nil {
		return 0
	}

	remotePortNum, err := strconv.ParseInt(remotePort, 10, 32)

	if err != nil {
		return 0
	}

	return uint16(remotePortNum)
}

// SetContextId sets the given request id to context
func (c *WebContext) SetContextId(requestId string) {
	c.Set(webContextRequestIdFieldKey, requestId)
}

// GetContextId returns the current request id
func (c *WebContext) GetContextId() string {
	requestId, exists := c.Get(webContextRequestIdFieldKey)

	if !exists {
		return ""
	}

	return requestId.(string)
}

// SetTextualToken sets the given user token to context
func (c *WebContext) SetTextualToken(token string) {
	c.Set(webContextTextualTokenFieldKey, token)
}

// GetTextualToken returns the current user textual token
func (c *WebContext) GetTextualToken() string {
	token, exists := c.Get(webContextTextualTokenFieldKey)

	if !exists {
		return ""
	}

	return token.(string)
}

// SetTokenClaims sets the given user token to context
func (c *WebContext) SetTokenClaims(claims *UserTokenClaims) {
	c.Set(webContextTokenClaimsFieldKey, claims)
}

// GetTokenClaims returns the current user token
func (c *WebContext) GetTokenClaims() *UserTokenClaims {
	claims, exists := c.Get(webContextTokenClaimsFieldKey)

	if !exists {
		return nil
	}

	return claims.(*UserTokenClaims)
}

// GetCurrentUid returns the current user uid by the current user token
func (c *WebContext) GetCurrentUid() int64 {
	claims := c.GetTokenClaims()

	if claims == nil {
		return 0
	}

	return claims.Uid
}

// GetTokenStringFromHeader returns the token string from the request header
func (c *WebContext) GetTokenStringFromHeader() string {
	tokenHeader := c.GetHeader(tokenHeaderName)

	if len(tokenHeader) < 7 || !strings.EqualFold(tokenHeader[:7], tokenHeaderValuePrefix) {
		return ""
	}

	return tokenHeader[7:]
}

// GetTokenStringFromQueryString returns the token string from the request query string
func (c *WebContext) GetTokenStringFromQueryString() string {
	return c.Query(tokenQueryStringParam)
}

// GetTokenStringFromCookie returns the token string from the request cookie
func (c *WebContext) GetTokenStringFromCookie() string {
	tokenCookie, err := c.Cookie(tokenCookieParam)

	if err != nil {
		return ""
	}

	return tokenCookie
}

func (c *WebContext) SetTokenStringToCookie(token string, tokenExpiredTime int, path string) {
	if token != "" {
		c.SetCookie(tokenCookieParam, token, tokenExpiredTime, path, "", false, true)
	} else {
		c.SetCookie(tokenCookieParam, "", -1, path, "", false, true)
	}
}

// GetClientLocale returns the client locale name
func (c *WebContext) GetClientLocale() string {
	value := c.GetHeader(AcceptLanguageHeaderName)

	return value
}

// GetClientTimezoneOffset returns the client timezone offset
func (c *WebContext) GetClientTimezoneOffset() (int16, error) {
	value := c.GetHeader(ClientTimezoneOffsetHeaderName)
	offset, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return int16(offset), nil
}

// SetResponseError sets the response error
func (c *WebContext) SetResponseError(error *errs.Error) {
	c.Set(webContextResponseErrorFieldKey, error)
}

// GetResponseError returns the response error
func (c *WebContext) GetResponseError() *errs.Error {
	err, exists := c.Get(webContextResponseErrorFieldKey)

	if !exists {
		return nil
	}

	return err.(*errs.Error)
}

// WrapWebContext returns a context wrapped by this file
func WrapWebContext(ginCtx *gin.Context) *WebContext {
	return &WebContext{
		Context: ginCtx,
	}
}
