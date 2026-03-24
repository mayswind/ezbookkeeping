package middlewares

import (
	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/settings"
)

// AmapApiProxyAuthCookie adds amap api proxy auth cookie to cookies in response
func AmapApiProxyAuthCookie(c *core.WebContext, config *settings.Config) {
	token := c.GetTextualToken()
	c.SetTokenStringToCookie(token, int(config.TokenExpiredTime), "/_AMapService")
}
