package core

import "github.com/mayswind/lab/pkg/errs"

type MiddlewareHandlerFunc func(*Context)

type ApiHandlerFunc func(*Context) (interface{}, *errs.Error)
