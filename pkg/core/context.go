package core

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mayswind/lab/pkg/errs"
)

const requestIdFieldKey = "REQUEST_ID"
const tokenClaimsFieldKey = "TOKEN_CLAIMS"
const responseErrorFieldKey = "RESPONSE_ERROR"

type Context struct {
	*gin.Context
	// DO NOT ADD ANY FIELD IN THIS CONTEXT, THIS CONTEXT IS JUST A WRAPPER
}

func (c *Context) SetRequestId(requestId string) {
	c.Set(requestIdFieldKey, requestId)
}

func (c *Context) GetRequestId() string {
	requestId, exists := c.Get(requestIdFieldKey)

	if !exists {
		return ""
	}

	return requestId.(string)
}

func (c *Context) SetTokenClaims(claims *UserTokenClaims) {
	c.Set(tokenClaimsFieldKey, claims)
}

func (c *Context) GetTokenClaims() *UserTokenClaims {
	claims, exists := c.Get(tokenClaimsFieldKey)

	if !exists {
		return nil
	}

	return claims.(*UserTokenClaims)
}

func (c *Context) GetCurrentUid() int64 {
	claims := c.GetTokenClaims()

	if claims == nil {
		return 0
	}

	uid, err := strconv.ParseInt(claims.Id, 10, 64)

	if err != nil {
		return 0
	}

	return uid
}

func (c *Context) SetResponseError(error *errs.Error) {
	c.Set(responseErrorFieldKey, error)
}

func (c *Context) GetResponseError() *errs.Error {
	err, exists := c.Get(responseErrorFieldKey)

	if !exists {
		return nil
	}

	return err.(*errs.Error)
}

func WrapContext(ginCtx *gin.Context) *Context {
	return &Context{
		Context: ginCtx,
	}
}
