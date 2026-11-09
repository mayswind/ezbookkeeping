package middlewares

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// AmapApiProxyAuthCookie adds amap api proxy auth cookie to cookies in response
func AmapApiProxyAuthCookie(c *core.WebContext, token string, expiresAt int64) {
	maxAge := int(expiresAt - time.Now().Unix())

	if maxAge <= 0 {
		token = ""
	}

	c.SetTokenStringToCookie(token, maxAge, "/_AMapService")
}
