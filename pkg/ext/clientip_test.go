package ext

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The login-failure limits (5 per minute per address) count failures per client address, and the address comes from
// gin's ClientIP with the app's "trusted_proxy_ips" setting (cmd/webserver.go). These tests record how that behaves
// behind a hosting proxy, so the setting is not left to guesswork. See docs/OPERATIONS.md.

// the app's default trusted_proxy_ips: private networks only
const defaultTrusted = "10.0.0.0/8,169.254.0.0/16,127.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

// two Cloudflare ranges, for the example only: the real list must be fetched from https://www.cloudflare.com/ips
const withCloudflare = defaultTrusted + ",172.64.0.0/13,104.16.0.0/13"

func clientIpFor(t *testing.T, trusted string, peer string, forwardedFor string) string {
	t.Helper()

	router := gin.New()
	require.NoError(t, router.SetTrustedProxies(strings.Split(trusted, ",")))

	var seen string
	router.GET("/", func(c *gin.Context) { seen = c.ClientIP() })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = peer + ":4000"

	if forwardedFor != "" {
		req.Header.Set("X-Forwarded-For", forwardedFor)
	}

	router.ServeHTTP(httptest.NewRecorder(), req)

	return seen
}

func TestClientIP_DirectConnectionIgnoresForwardedHeaders(t *testing.T) {
	// nobody can pretend to be another address by sending the header themselves
	assert.Equal(t, "198.51.100.7", clientIpFor(t, defaultTrusted, "198.51.100.7", "1.2.3.4"))
}

func TestClientIP_PrivateProxyPassesTheRealAddress(t *testing.T) {
	// proxy on a private address that appends the visitor to X-Forwarded-For: the visitor is found
	assert.Equal(t, "203.0.113.5", clientIpFor(t, defaultTrusted, "10.20.30.40", "203.0.113.5"))
	assert.Equal(t, "203.0.113.5", clientIpFor(t, defaultTrusted, "10.20.30.40", "203.0.113.5, 10.1.1.1"))
}

func TestClientIP_AnAttackerCannotInventAnAddressToDodgeTheLimit(t *testing.T) {
	// the header is read from the right, so extra addresses the visitor adds on the left are ignored
	assert.Equal(t, "203.0.113.5", clientIpFor(t, defaultTrusted, "10.20.30.40", "9.9.9.9, 203.0.113.5"))
}

func TestClientIP_AnotherCompanysProxyInTheChainBecomesEveryonesAddress(t *testing.T) {
	// If the chain is "visitor, cdn" and the cdn's address is not trusted, the cdn's address is taken as the visitor.
	// Every visitor that arrives through that cdn machine then shares one counter, and five wrong passwords from
	// anyone lock them all out for a minute.
	assert.Equal(t, "172.70.1.1", clientIpFor(t, defaultTrusted, "10.20.30.40", "203.0.113.5, 172.70.1.1"))
	assert.Equal(t, "172.70.1.1", clientIpFor(t, defaultTrusted, "10.20.30.40", "198.51.100.9, 172.70.1.1"), "a different visitor, the same address")

	// trusting the cdn's published ranges as well makes the real visitor visible again
	assert.Equal(t, "203.0.113.5", clientIpFor(t, withCloudflare, "10.20.30.40", "203.0.113.5, 172.70.1.1"))
	assert.Equal(t, "198.51.100.9", clientIpFor(t, withCloudflare, "10.20.30.40", "198.51.100.9, 172.70.1.1"))
}
