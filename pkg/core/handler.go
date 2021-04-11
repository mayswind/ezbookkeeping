package core

import "github.com/mayswind/ezbookkeeping/pkg/errs"

// MiddlewareHandlerFunc represents the middleware handler function
type MiddlewareHandlerFunc func(*Context)

// ApiHandlerFunc represents the api handler function
type ApiHandlerFunc func(*Context) (interface{}, *errs.Error)

// DataHandlerFunc represents the handler function that returns byte array
type DataHandlerFunc func(*Context) ([]byte, string, *errs.Error)
