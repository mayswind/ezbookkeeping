package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// DefaultApi represents default api
type DefaultApi struct{}

// Initialize a default api singleton instance
var (
	Default = &DefaultApi{}
)

// ApiNotFound returns api not found error
func (a *DefaultApi) ApiNotFound(c *core.Context) (any, *errs.Error) {
	return nil, errs.ErrApiNotFound
}

// MethodNotAllowed returns method not allowed error
func (a *DefaultApi) MethodNotAllowed(c *core.Context) (any, *errs.Error) {
	return nil, errs.ErrMethodNotAllowed
}
