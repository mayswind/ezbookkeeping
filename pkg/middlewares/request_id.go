package middlewares

import (
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/requestid"
	"github.com/mayswind/lab/pkg/settings"
)

const REQUEST_ID_HEADER = "X-Request-ID"

func RequestId(config *settings.Config) core.MiddlewareHandlerFunc {
	return func(c *core.Context) {
		if requestid.Container.Current == nil {
			c.Next()
			return
		}

		requestId := requestid.Container.Current.GenerateRequestId(c.ClientIP())
		c.SetRequestId(requestId)

		if config.EnableRequestIdHeader {
			c.Header(REQUEST_ID_HEADER, requestId)
		}

		c.Next()
	}
}
