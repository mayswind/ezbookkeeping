package middlewares

import (
	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/requestid"
	"github.com/Paxtiny/oscar/pkg/settings"
)

const requestIdHeader = "X-Request-ID"

// RequestId generates a new request id and add it to context and response header
func RequestId(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.WebContext) {
		requestId := requestid.Container.GenerateRequestId(c.ClientIP(), c.ClientPort())
		c.SetContextId(requestId)

		if config.EnableRequestIdHeader {
			c.Header(requestIdHeader, requestId)
		}

		c.Next()
	}
}
