package middlewares

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const tokenCookieParam = "ebk_auth_token"

// AmapApiProxyAuthCookie adds amap api proxy auth cookie to cookies in response
func AmapApiProxyAuthCookie(c *core.Context, config *settings.Config) {
	token := c.GetTextualToken()

	if token != "" {
		c.SetCookie(tokenCookieParam, token, int(config.TokenExpiredTime), "/_AMapService", "", false, true)
	} else {
		c.SetCookie(tokenCookieParam, "", -1, "/_AMapService", "", false, true)
	}
}
