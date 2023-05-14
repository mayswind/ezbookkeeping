package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

const openStreetMapTileImageUrlFormat = "https://tile.openstreetmap.org/%s/%s/%s" // https://tile.openstreetmap.org/{z}/{x}/{y}.png

// MapImageProxy represents map image proxy
type MapImageProxy struct {
}

// Initialize a map image proxy singleton instance
var (
	MapImages = &MapImageProxy{}
)

// OpenStreetMapTileImageProxyHandler returns open street map tile image
func (p *MapImageProxy) OpenStreetMapTileImageProxyHandler(c *core.Context) (*httputil.ReverseProxy, *errs.Error) {
	director := func(req *http.Request) {
		zoomLevel := c.Param("zoomLevel")
		coordinateX := c.Param("coordinateX")
		fileName := c.Param("fileName")

		imageRawUrl := fmt.Sprintf(openStreetMapTileImageUrlFormat, zoomLevel, coordinateX, fileName)
		imageUrl, _ := url.Parse(imageRawUrl)

		req.Header.Del("Authorization")
		req.URL = imageUrl
		req.RequestURI = req.URL.RequestURI()
		req.Host = imageUrl.Host
	}

	return &httputil.ReverseProxy{Director: director}, nil
}
