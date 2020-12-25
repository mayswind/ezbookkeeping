package core

import "github.com/mayswind/lab/pkg/errs"

// MiddlewareHandlerFunc represents the middleware handler function
type MiddlewareHandlerFunc func(*Context)

// ApiHandlerFunc represents the api handler function
type ApiHandlerFunc func(*Context) (interface{}, *errs.Error)
