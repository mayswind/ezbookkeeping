package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const openStreetMapTileImageUrlFormat = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"                                                                                                                              // https://tile.openstreetmap.org/{z}/{x}/{y}.png
const openStreetMapHumanitarianStyleTileImageUrlFormat = "https://a.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png"                                                                                                        // https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png
const openTopoMapTileImageUrlFormat = "https://tile.opentopomap.org/{z}/{x}/{y}.png"                                                                                                                                  // https://tile.opentopomap.org/{z}/{x}/{y}.png
const opnvKarteMapTileImageUrlFormat = "https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png"                                                                                                                       // https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png
const cyclOSMMapTileImageUrlFormat = "https://a.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png"                                                                                                                // https://{s}.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png
const cartoDBMapTileImageUrlFormat = "https://a.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{scale}.png"                                                                                                     // https://{s}.basemaps.cartocdn.com/{style}/{z}/{x}/{y}{scale}.png
const tomtomMapTileImageUrlFormat = "https://api.tomtom.com/map/1/tile/basic/main/{z}/{x}/{y}.png"                                                                                                                    // https://api.tomtom.com/map/{versionNumber}/tile/{layer}/{style}/{z}/{x}/{y}.png?key={key}&language={language}
const tianDiTuMapTileImageUrlFormat = "https://t0.tianditu.gov.cn/vec_w/wmts?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=vec&STYLE=default&TILEMATRIXSET=w&FORMAT=tiles&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}"  // https://t{s}.tianditu.gov.cn/vec_w/wmts?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=vec&STYLE=default&TILEMATRIXSET=w&FORMAT=tiles&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}&tk={key}
const tianDiTuMapAnnotationUrlFormat = "https://t0.tianditu.gov.cn/cva_w/wmts?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=cva&STYLE=default&TILEMATRIXSET=w&FORMAT=tiles&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}" // https://t{s}.tianditu.gov.cn/cva_w/wmts?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=cva&STYLE=default&TILEMATRIXSET=w&FORMAT=tiles&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}&tk={key}

// MapImageProxy represents map image proxy
type MapImageProxy struct {
	ApiUsingConfig
}

// Initialize a map image proxy singleton instance
var (
	MapImages = &MapImageProxy{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// MapTileImageProxyHandler returns map tile image
func (p *MapImageProxy) MapTileImageProxyHandler(c *core.Context) (*httputil.ReverseProxy, *errs.Error) {
	return p.mapImageProxyHandler(c, func(c *core.Context, mapProvider string) (string, *errs.Error) {
		if mapProvider == settings.OpenStreetMapProvider {
			return openStreetMapTileImageUrlFormat, nil
		} else if mapProvider == settings.OpenStreetMapHumanitarianStyleProvider {
			return openStreetMapHumanitarianStyleTileImageUrlFormat, nil
		} else if mapProvider == settings.OpenTopoMapProvider {
			return openTopoMapTileImageUrlFormat, nil
		} else if mapProvider == settings.OPNVKarteMapProvider {
			return opnvKarteMapTileImageUrlFormat, nil
		} else if mapProvider == settings.CyclOSMMapProvider {
			return cyclOSMMapTileImageUrlFormat, nil
		} else if mapProvider == settings.CartoDBMapProvider {
			return cartoDBMapTileImageUrlFormat, nil
		} else if mapProvider == settings.TomTomMapProvider {
			targetUrl := tomtomMapTileImageUrlFormat + "?key=" + p.CurrentConfig().TomTomMapAPIKey
			language := c.Query("language")

			if language != "" {
				targetUrl = targetUrl + "&language=" + language
			}

			return targetUrl, nil
		} else if mapProvider == settings.TianDiTuProvider {
			return tianDiTuMapTileImageUrlFormat + "&tk=" + p.CurrentConfig().TianDiTuAPIKey, nil
		} else if mapProvider == settings.CustomProvider {
			return p.CurrentConfig().CustomMapTileServerTileLayerUrl, nil
		}

		return "", errs.ErrParameterInvalid
	})
}

// MapAnnotationImageProxyHandler returns map annotation image
func (p *MapImageProxy) MapAnnotationImageProxyHandler(c *core.Context) (*httputil.ReverseProxy, *errs.Error) {
	return p.mapImageProxyHandler(c, func(c *core.Context, mapProvider string) (string, *errs.Error) {
		if mapProvider == settings.TianDiTuProvider {
			return tianDiTuMapAnnotationUrlFormat + "&tk=" + p.CurrentConfig().TianDiTuAPIKey, nil
		} else if mapProvider == settings.CustomProvider {
			return p.CurrentConfig().CustomMapTileServerAnnotationLayerUrl, nil
		}

		return "", errs.ErrParameterInvalid
	})
}

func (p *MapImageProxy) mapImageProxyHandler(c *core.Context, fn func(c *core.Context, mapProvider string) (string, *errs.Error)) (*httputil.ReverseProxy, *errs.Error) {
	mapProvider := strings.Replace(c.Query("provider"), "-", "_", -1)
	targetUrl := ""

	if mapProvider != p.CurrentConfig().MapProvider {
		return nil, errs.ErrMapProviderNotCurrent
	}

	zoomLevel := c.Param("zoomLevel")
	coordinateX := c.Param("coordinateX")
	fileName := c.Param("fileName")
	fileNameParts := strings.Split(fileName, ".")
	coordinateY := fileNameParts[0]
	scale := c.Query("scale")

	if len(fileNameParts) != 2 || fileNameParts[len(fileNameParts)-1] != "png" {
		return nil, errs.ErrImageExtensionNotSupported
	}

	var err *errs.Error
	targetUrl, err = fn(c, mapProvider)

	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	utils.SetProxyUrl(transport, p.CurrentConfig().MapProxy)

	director := func(req *http.Request) {
		imageRawUrl := targetUrl
		imageRawUrl = strings.Replace(imageRawUrl, "{z}", zoomLevel, -1)
		imageRawUrl = strings.Replace(imageRawUrl, "{x}", coordinateX, -1)
		imageRawUrl = strings.Replace(imageRawUrl, "{y}", coordinateY, -1)
		imageRawUrl = strings.Replace(imageRawUrl, "{scale}", scale, -1)
		imageUrl, _ := url.Parse(imageRawUrl)

		req.URL = imageUrl
		req.RequestURI = req.URL.RequestURI()
		req.Host = imageUrl.Host
	}

	return &httputil.ReverseProxy{
		Transport: transport,
		Director:  director,
	}, nil
}
