package httpclient

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

const (
	logHandleKey = "log_handler"
)

// HttpResponseLogHandlerFunc represents the http response log handler function
type HttpResponseLogHandlerFunc func([]byte)

// httpRequestContext represents the context for http request
type httpRequestContext struct {
	core.Context
	logHandler HttpResponseLogHandlerFunc
}

// Value returns the value associated with key
func (c *httpRequestContext) Value(key any) any {
	if key == logHandleKey {
		return c.logHandler
	}

	return c.Context.Value(key)
}

// CustomHttpResponseLog returns a context with http response log handler
func CustomHttpResponseLog(c core.Context, responseLogHandler HttpResponseLogHandlerFunc) core.Context {
	return &httpRequestContext{
		Context:    c,
		logHandler: responseLogHandler,
	}
}
