package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const openStreetMapTileImageUrlFormat = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"                       // https://tile.openstreetmap.org/{z}/{x}/{y}.png
const openStreetMapHumanitarianStyleTileImageUrlFormat = "https://a.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png" // https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png
const openTopoMapTileImageUrlFormat = "https://tile.opentopomap.org/{z}/{x}/{y}.png"                           // https://tile.opentopomap.org/{z}/{x}/{y}.png
const opnvKarteMapTileImageUrlFormat = "https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png"                // https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png
const cyclOSMMapTileImageUrlFormat = "https://a.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png"         // https://{s}.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png
const tomtomMapTileImageUrlFormat = "https://api.tomtom.com/map/1/tile/basic/main/{z}/{x}/{y}.png"             // https://api.tomtom.com/map/{versionNumber}/tile/{layer}/{style}/{z}/{x}/{y}.png?key={key}&language={language}

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

	if mapProvider != settings.Container.Current.MapProvider {
		return nil, errs.ErrMapProviderNotCurrent
	}

	zoomLevel := c.Param("zoomLevel")
	coordinateX := c.Param("coordinateX")
	fileName := c.Param("fileName")
	fileNameParts := strings.Split(fileName, ".")
	coordinateY := fileNameParts[0]

	if len(fileNameParts) != 2 || fileNameParts[len(fileNameParts)-1] != "png" {
		return nil, errs.ErrImageExtensionNotSupported
	}

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
	} else if mapProvider == settings.TomTomMapProvider {
		targetUrl = tomtomMapTileImageUrlFormat + "?key=" + settings.Container.Current.TomTomMapAPIKey
		language := c.Query("language")

		if language != "" {
			targetUrl = targetUrl + "&language=" + language
		}
	} else if mapProvider == settings.CustomProvider {
		targetUrl = settings.Container.Current.CustomMapTileServerUrl
	} else {
		return nil, errs.ErrParameterInvalid
	}

	director := func(req *http.Request) {
		imageRawUrl := targetUrl
		imageRawUrl = strings.Replace(imageRawUrl, "{z}", zoomLevel, -1)
		imageRawUrl = strings.Replace(imageRawUrl, "{x}", coordinateX, -1)
		imageRawUrl = strings.Replace(imageRawUrl, "{y}", coordinateY, -1)
		imageUrl, _ := url.Parse(imageRawUrl)

		req.URL = imageUrl
		req.RequestURI = req.URL.RequestURI()
		req.Host = imageUrl.Host
	}

	return &httputil.ReverseProxy{Director: director}, nil
}
