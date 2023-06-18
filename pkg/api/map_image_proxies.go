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

const openStreetMapTileImageUrlFormat = "https://tile.openstreetmap.org/%s/%s/%s"                       // https://tile.openstreetmap.org/{z}/{x}/{y}.png
const openStreetMapHumanitarianStyleTileImageUrlFormat = "https://a.tile.openstreetmap.fr/hot/%s/%s/%s" // https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png
const openTopoMapTileImageUrlFormat = "https://tile.opentopomap.org/%s/%s/%s"                           // https://tile.opentopomap.org/{z}/{x}/{y}.png
const opnvKarteMapTileImageUrlFormat = "https://tileserver.memomaps.de/tilegen/%s/%s/%s"                // https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png
const cyclOSMMapTileImageUrlFormat = "https://a.tile-cyclosm.openstreetmap.fr/cyclosm/%s/%s/%s"         // https://{s}.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png

// MapImageProxy represents map image proxy
type MapImageProxy struct {
}

// Initialize a map image proxy singleton instance
var (
	MapImages = &MapImageProxy{}
)

// MapTileImageProxyHandler returns map tile image
func (p *MapImageProxy) MapTileImageProxyHandler(c *core.Context) (*httputil.ReverseProxy, *errs.Error) {
	mapProvider := strings.Replace(c.Query("provider"), "-", "_", -1)
	targetUrl := ""

	if mapProvider == settings.OpenStreetMapProvider {
		targetUrl = openStreetMapTileImageUrlFormat
	} else if mapProvider == settings.OpenStreetMapHumanitarianStyleProvider {
		targetUrl = openStreetMapHumanitarianStyleTileImageUrlFormat
	} else if mapProvider == settings.OpenTopoMapProvider {
		targetUrl = openTopoMapTileImageUrlFormat
	} else if mapProvider == settings.OPNVKarteMapProvider {
		targetUrl = opnvKarteMapTileImageUrlFormat
	} else if mapProvider == settings.CyclOSMMapProvider {
		targetUrl = cyclOSMMapTileImageUrlFormat
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
