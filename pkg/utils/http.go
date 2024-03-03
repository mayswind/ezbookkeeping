package utils

import (
	"net/http"
	"net/url"
)

// SetProxyUrl sets proxy url to http transport according to specified proxy setting
func SetProxyUrl(transport *http.Transport, proxy string) {
	if proxy == "none" {
		transport.Proxy = nil
	} else if proxy != "system" {
		proxy, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxy)
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}
}
