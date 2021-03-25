package middlewares

import (
	"github.com/mayswind/lab/pkg/core"
)

const utcOffsetQueryStringParam = "utc_offset"

// HeaderInQueryString puts some headers from query string
func HeaderInQueryString(c *core.Context) {
	utcOffset, exists := c.GetQuery(utcOffsetQueryStringParam)

	if exists {
		c.Request.Header.Set(core.ClientTimezoneOffsetHeaderName, utcOffset)
	}
}
