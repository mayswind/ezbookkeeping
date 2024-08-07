package core

import (
	"net/http/httputil"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// MiddlewareHandlerFunc represents the middleware handler function
type MiddlewareHandlerFunc func(*Context)

// ApiHandlerFunc represents the api handler function
type ApiHandlerFunc func(*Context) (any, *errs.Error)

// DataHandlerFunc represents the handler function that returns file data byte array and file name
type DataHandlerFunc func(*Context) ([]byte, string, *errs.Error)

// ImageHandlerFunc represents the handler function that returns image byte array and content type
type ImageHandlerFunc func(*Context) ([]byte, string, *errs.Error)

// ProxyHandlerFunc represents the reverse proxy handler function
type ProxyHandlerFunc func(*Context) (*httputil.ReverseProxy, *errs.Error)
