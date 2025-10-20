package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type defaultTransport struct {
	defaultUserAgent string
	baseTransport    http.RoundTripper
}

func (t *defaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Values("User-Agent")) < 1 {
		req.Header.Set("User-Agent", t.defaultUserAgent)
	} else if req.Header.Get("User-Agent") == "" {
		req.Header.Del("User-Agent")
	}

	return t.baseTransport.RoundTrip(req)
}

// NewHttpClient creates and returns a new http client with specified settings
func NewHttpClient(requestTimeout uint32, proxy string, skipTLSVerify bool, defaultUserAgent string) *http.Client {
	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	SetProxyUrl(baseTransport, proxy)

	if skipTLSVerify {
		baseTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Transport: &defaultTransport{
			defaultUserAgent: defaultUserAgent,
			baseTransport:    baseTransport,
		},
		Timeout: time.Duration(requestTimeout) * time.Millisecond,
	}
}

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
