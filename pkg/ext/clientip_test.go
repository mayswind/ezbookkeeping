package ext

import (
	"net"
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

// the private networks, the machine's own IPv6 loopback and Cloudflare's published ranges as of 9 Oct 2026
// (scripts/trusted-proxies.sh prints the current list)
const withCloudflare = "10.0.0.0/8,169.254.0.0/16,127.0.0.0/8,172.16.0.0/12,192.168.0.0/16,::1/128,173.245.48.0/20,103.21.244.0/22,103.22.200.0/22,103.31.4.0/22,141.101.64.0/18,108.162.192.0/18,190.93.240.0/20,188.114.96.0/20,197.234.240.0/22,198.41.128.0/17,162.158.0.0/15,104.16.0.0/13,104.24.0.0/14,172.64.0.0/13,131.0.72.0/22,2400:cb00::/32,2606:4700::/32,2803:f800::/32,2405:b500::/32,2405:8100::/32,2a06:98c0::/29,2c0f:f248::/32"

func clientIpFor(t *testing.T, trusted string, peer string, forwardedFor string) string {
	t.Helper()

	router := gin.New()
	require.NoError(t, router.SetTrustedProxies(strings.Split(trusted, ",")))

	var seen string
	router.GET("/", func(c *gin.Context) { seen = c.ClientIP() })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = net.JoinHostPort(peer, "4000")

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

// What Render shows: every request arrives from ::1 (the machine's own IPv6 loopback address, so a proxy on the same machine
// hands the traffic over). The default trusted list has 127.0.0.0/8, which is the IPv4 loopback only, so ::1 is not trusted,
// the forwarding header is ignored and every visitor is "::1".
func TestClientIP_IPv6LoopbackProxyIsNotTrustedByDefault(t *testing.T) {
	assert.Equal(t, "::1", clientIpFor(t, defaultTrusted, "::1", "203.0.113.5"))
	assert.Equal(t, "::1", clientIpFor(t, defaultTrusted, "::1", "198.51.100.9"), "a different visitor, the same address: they share one failure counter")

	// trusting ::1 as well lets the app read the real visitor from the header
	assert.Equal(t, "203.0.113.5", clientIpFor(t, defaultTrusted+",::1/128", "::1", "203.0.113.5"))
	assert.Equal(t, "198.51.100.9", clientIpFor(t, defaultTrusted+",::1/128", "::1", "198.51.100.9"))
}

func TestClientIP_TrustingLoopbackDoesNotLetOutsidersFakeAnAddress(t *testing.T) {
	// only a connection that really comes from the machine itself is believed; a remote visitor's header is still ignored
	assert.Equal(t, "198.51.100.7", clientIpFor(t, defaultTrusted+",::1/128", "198.51.100.7", "1.2.3.4"))
	// and even from the proxy, extra addresses added on the left by the visitor do not count
	assert.Equal(t, "203.0.113.5", clientIpFor(t, defaultTrusted+",::1/128", "::1", "9.9.9.9, 203.0.113.5"))
}

// What Render shows after ::1 is trusted: the next hop is a Cloudflare server (172.71.150.174 was seen in the log).
func TestClientIP_BehindLoopbackAndCloudflareTheRealVisitorIsFound(t *testing.T) {
	const cloudflareEdge = "172.71.150.174"
	onlyLoopbackTrusted := defaultTrusted + ",::1/128"

	// without Cloudflare's ranges every visitor arriving through that server is "172.71.150.174"
	assert.Equal(t, cloudflareEdge, clientIpFor(t, onlyLoopbackTrusted, "::1", "203.0.113.5, "+cloudflareEdge))
	assert.Equal(t, cloudflareEdge, clientIpFor(t, onlyLoopbackTrusted, "::1", "198.51.100.9, "+cloudflareEdge))

	// with them, each visitor is themselves, IPv4 and IPv6 alike
	assert.Equal(t, "203.0.113.5", clientIpFor(t, withCloudflare, "::1", "203.0.113.5, "+cloudflareEdge))
	assert.Equal(t, "198.51.100.9", clientIpFor(t, withCloudflare, "::1", "198.51.100.9, "+cloudflareEdge))
	assert.Equal(t, "2001:db8::77", clientIpFor(t, withCloudflare, "::1", "2001:db8::77, 2606:4700:10::1"))

	// a visitor still cannot choose their own address by sending a header
	assert.Equal(t, "203.0.113.5", clientIpFor(t, withCloudflare, "::1", "9.9.9.9, 203.0.113.5, "+cloudflareEdge))
}
