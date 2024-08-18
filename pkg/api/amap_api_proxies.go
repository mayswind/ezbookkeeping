package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const amapCustomMapStylesUrl = "https://webapi.amap.com/v4/map/styles"
const amapOverseasMapUrl = "https://fmap01.amap.com/v3/vectormap"
const amapRestApiUrl = "https://restapi.amap.com/"

// AmapApiProxy represents amap api proxy
type AmapApiProxy struct {
	ApiUsingConfig
}

// Initialize a amap api proxy singleton instance
var (
	AmapApis = &AmapApiProxy{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// AmapApiProxyHandler returns amap api response
func (p *AmapApiProxy) AmapApiProxyHandler(c *core.WebContext) (*httputil.ReverseProxy, *errs.Error) {
	var targetUrl string

	if strings.HasPrefix(c.Request.RequestURI, "/_AMapService/v4/map/styles") {
		targetUrl = amapCustomMapStylesUrl + strings.TrimPrefix(c.Request.URL.Path, "/_AMapService/v4/map/styles")
	} else if strings.HasPrefix(c.Request.RequestURI, "/_AMapService/v3/vectormap") {
		targetUrl = amapOverseasMapUrl + strings.TrimPrefix(c.Request.URL.Path, "/_AMapService/v3/vectormap")
	} else {
		targetUrl = amapRestApiUrl + strings.TrimPrefix(c.Request.URL.Path, "/_AMapService/")
	}

	director := func(req *http.Request) {
		targetRawUrl := fmt.Sprintf("%s?%s&jscode=%s", targetUrl, req.URL.RawQuery, p.CurrentConfig().AmapApplicationSecret)
		targetUrl, _ := url.Parse(targetRawUrl)

		oldCookies := req.Cookies()
		req.Header.Del("Cookie")

		for i := 0; i < len(oldCookies); i++ {
			if strings.HasPrefix(oldCookies[i].Name, "ebk_") {
				continue
			}

			req.AddCookie(oldCookies[i])
		}

		req.URL = targetUrl
		req.RequestURI = req.URL.RequestURI()
		req.Host = targetUrl.Host
	}

	return &httputil.ReverseProxy{Director: director}, nil
}
