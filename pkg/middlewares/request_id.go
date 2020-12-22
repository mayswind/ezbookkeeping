package middlewares

import (
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/requestid"
	"github.com/mayswind/lab/pkg/settings"
)

const requestIdHeader = "X-Request-ID"

func RequestId(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.Context) {
		if requestid.Container.Current == nil {
			c.Next()
			return
		}

		requestId := requestid.Container.Current.GenerateRequestId(c.ClientIP())
		c.SetRequestId(requestId)

		if config.EnableRequestIdHeader {
			c.Header(requestIdHeader, requestId)
		}

		c.Next()
	}
}
