package core

import (
	"net"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const requestIdFieldKey = "REQUEST_ID"
const textualTokenFieldKey = "TOKEN_STRING"
const tokenClaimsFieldKey = "TOKEN_CLAIMS"
const responseErrorFieldKey = "RESPONSE_ERROR"

// AcceptLanguageHeaderName represents the header name of accept language
const AcceptLanguageHeaderName = "Accept-Language"

// RemoteClientPortHeader represents the header name of remote client source port
const RemoteClientPortHeader = "X-Real-Port"

// ClientTimezoneOffsetHeaderName represents the header name of client timezone offset
const ClientTimezoneOffsetHeaderName = "X-Timezone-Offset"

// Context represents the request and response context
type Context struct {
	*gin.Context
	// DO NOT ADD ANY FIELD IN THIS CONTEXT, THIS CONTEXT IS JUST A WRAPPER
}

func (c *Context) ClientPort() uint16 {
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

// SetRequestId sets the given request id to context
func (c *Context) SetRequestId(requestId string) {
	c.Set(requestIdFieldKey, requestId)
}

// GetRequestId returns the current request id
func (c *Context) GetRequestId() string {
	requestId, exists := c.Get(requestIdFieldKey)

	if !exists {
		return ""
	}

	return requestId.(string)
}

// SetTextualToken sets the given user token to context
func (c *Context) SetTextualToken(token string) {
	c.Set(textualTokenFieldKey, token)
}

// GetTextualToken returns the current user textual token
func (c *Context) GetTextualToken() string {
	token, exists := c.Get(textualTokenFieldKey)

	if !exists {
		return ""
	}

	return token.(string)
}

// SetTokenClaims sets the given user token to context
func (c *Context) SetTokenClaims(claims *UserTokenClaims) {
	c.Set(tokenClaimsFieldKey, claims)
}

// GetTokenClaims returns the current user token
func (c *Context) GetTokenClaims() *UserTokenClaims {
	claims, exists := c.Get(tokenClaimsFieldKey)

	if !exists {
		return nil
	}

	return claims.(*UserTokenClaims)
}

// GetCurrentUid returns the current user uid by the current user token
func (c *Context) GetCurrentUid() int64 {
	claims := c.GetTokenClaims()

	if claims == nil {
		return 0
	}

	return claims.Uid
}

// GetClientLocale returns the client locale name
func (c *Context) GetClientLocale() string {
	value := c.GetHeader(AcceptLanguageHeaderName)

	return value
}

// GetClientTimezoneOffset returns the client timezone offset
func (c *Context) GetClientTimezoneOffset() (int16, error) {
	value := c.GetHeader(ClientTimezoneOffsetHeaderName)
	offset, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return int16(offset), nil
}

// SetResponseError sets the response error
func (c *Context) SetResponseError(error *errs.Error) {
	c.Set(responseErrorFieldKey, error)
}

// GetResponseError returns the response error
func (c *Context) GetResponseError() *errs.Error {
	err, exists := c.Get(responseErrorFieldKey)

	if !exists {
		return nil
	}

	return err.(*errs.Error)
}

// WrapContext returns a context wrapped by this file
func WrapContext(ginCtx *gin.Context) *Context {
	return &Context{
		Context: ginCtx,
	}
}
