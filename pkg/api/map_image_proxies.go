package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const openStreetMapTileImageUrlFormat = "https://tile.openstreetmap.org/%s/%s/%s" // https://tile.openstreetmap.org/{z}/{x}/{y}.png

// MapImageProxy represents map image proxy
type MapImageProxy struct {
}

// Initialize a map image proxy singleton instance
var (
	MapImages = &MapImageProxy{}
)

// MapTileImageProxyHandler returns map tile image
func (p *MapImageProxy) MapTileImageProxyHandler(c *core.Context) (*httputil.ReverseProxy, *errs.Error) {
	mapProvider := c.Query("provider")
	targetUrl := ""

	if mapProvider == settings.OpenStreetMapProvider {
		targetUrl = openStreetMapTileImageUrlFormat
	} else {
		return nil, errs.ErrParameterInvalid
	}

	director := func(req *http.Request) {
		zoomLevel := c.Param("zoomLevel")
		coordinateX := c.Param("coordinateX")
		fileName := c.Param("fileName")

		imageRawUrl := fmt.Sprintf(targetUrl, zoomLevel, coordinateX, fileName)
		imageUrl, _ := url.Parse(imageRawUrl)

		req.URL = imageUrl
		req.RequestURI = req.URL.RequestURI()
		req.Host = imageUrl.Host
	}

	return &httputil.ReverseProxy{Director: director}, nil
}
