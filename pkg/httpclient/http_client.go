package httpclient

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"
)

type defaultTransport struct {
	defaultUserAgent      string
	enableHttpResponseLog bool
	baseTransport         http.RoundTripper
}

func (t *defaultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Values("User-Agent")) < 1 {
		req.Header.Set("User-Agent", t.defaultUserAgent)
	} else if req.Header.Get("User-Agent") == "" {
		req.Header.Del("User-Agent")
	}

	resp, err := t.baseTransport.RoundTrip(req)

	if t.enableHttpResponseLog && err == nil {
		ctx := req.Context()

		if handler, ok := ctx.Value(logHandleKey).(HttpResponseLogHandlerFunc); ok {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				return nil, err
			}

			handler(body)

			resp.Body = io.NopCloser(bytes.NewReader(body))
		}
	}

	return resp, err
}

// NewHttpClient creates and returns a new http client with specified settings
func NewHttpClient(requestTimeout uint32, proxy string, skipTLSVerify bool, defaultUserAgent string, enableHttpResponseLog bool, additionTLSCipherSuites []uint16) *http.Client {
	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	SetProxyUrl(baseTransport, proxy)

	if skipTLSVerify || len(additionTLSCipherSuites) > 0 {
		tlsClientConfig := &tls.Config{}

		if skipTLSVerify {
			tlsClientConfig.InsecureSkipVerify = true
		}

		if len(additionTLSCipherSuites) > 0 {
			var tlsCipherSuites []uint16

			for _, cipherSuite := range tls.CipherSuites() {
				tlsCipherSuites = append(tlsCipherSuites, cipherSuite.ID)
			}

			tlsCipherSuites = append(tlsCipherSuites, additionTLSCipherSuites...)
			tlsClientConfig.CipherSuites = tlsCipherSuites
		}

		baseTransport.TLSClientConfig = tlsClientConfig
	}

	return &http.Client{
		Transport: &defaultTransport{
			defaultUserAgent:      defaultUserAgent,
			enableHttpResponseLog: enableHttpResponseLog,
			baseTransport:         baseTransport,
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
