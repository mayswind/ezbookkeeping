package core

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mayswind/lab/pkg/errs"
)

const FIELD_REQUEST_ID_KEY = "REQUEST_ID"
const FIELD_TOKEN_CLAIMS_KEY = "TOKEN_CLAIMS"
const FIELD_RESPONSE_ERROR = "RESPONSE_ERROR"

type Context struct {
	*gin.Context
	// DO NOT ADD ANY FIELD IN THIS CONTEXT, THIS CONTEXT IS JUST A WRAPPER
}

func (c *Context) SetRequestId(requestId string) {
	c.Set(FIELD_REQUEST_ID_KEY, requestId)
}

func (c *Context) GetRequestId() string {
	requestId, exists := c.Get(FIELD_REQUEST_ID_KEY)

	if !exists {
		return ""
	}

	return requestId.(string)
}

func (c *Context) SetTokenClaims(claims *UserTokenClaims) {
	c.Set(FIELD_TOKEN_CLAIMS_KEY, claims)
}

func (c *Context) GetTokenClaims() *UserTokenClaims {
	claims, exists := c.Get(FIELD_TOKEN_CLAIMS_KEY)

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
	c.Set(FIELD_RESPONSE_ERROR, error)
}

func (c *Context) GetResponseError() *errs.Error {
	err, exists := c.Get(FIELD_RESPONSE_ERROR)

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
