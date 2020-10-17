package api

import (
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
)

type DefaultApi struct {}

var (
	Default = &DefaultApi{}
)

func (a *DefaultApi) ApiNotFound(c *core.Context) (interface{}, *errs.Error) {
	return nil, errs.ErrApiNotFound
}

func (a *DefaultApi) MethodNotAllowed(c *core.Context) (interface{}, *errs.Error) {
	return nil, errs.ErrMethodNotAllowed
}
