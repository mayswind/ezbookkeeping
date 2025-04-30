package core

import (
	"net/http/httputil"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// CliHandlerFunc represents the cli handler function
type CliHandlerFunc func(*CliContext) error

// MiddlewareHandlerFunc represents the middleware handler function
type MiddlewareHandlerFunc func(*WebContext)

// ApiHandlerFunc represents the api handler function
type ApiHandlerFunc func(*WebContext) (any, *errs.Error)

// EventStreamApiHandlerFunc represents the event stream api handler function
type EventStreamApiHandlerFunc func(*WebContext) *errs.Error

// DataHandlerFunc represents the handler function that returns file data byte array and file name
type DataHandlerFunc func(*WebContext) ([]byte, string, *errs.Error)

// ImageHandlerFunc represents the handler function that returns image byte array and content type
type ImageHandlerFunc func(*WebContext) ([]byte, string, *errs.Error)

// ProxyHandlerFunc represents the reverse proxy handler function
type ProxyHandlerFunc func(*WebContext) (*httputil.ReverseProxy, *errs.Error)
