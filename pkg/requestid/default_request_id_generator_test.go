package requestid

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestNewDefaultRequestIdGenerator_Http(t *testing.T) {
	generator, _ := NewDefaultRequestIdGenerator(&settings.Config{HttpAddr: "123.234.123.234", HttpPort: 8080, SecretKey: "secretkey"})
	requestId := generator.GenerateRequestId("127.0.0.1", 20000)
	requestIdInfo := generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedServerUniqId := uint16(0x2476) // crc32("123.234.123.234" + "_" + "secretkey") & 0xFFFF
	actualServerUniqId := requestIdInfo.ServerUniqId
	assert.Equal(t, expectedServerUniqId, actualServerUniqId)

	expectedInstanceUniqId := uint16(0x0e79) // crc32("8080" + "_" + "secretkey") & 0xFFFF
	actualInstanceUniqId := requestIdInfo.InstanceUniqId
	assert.Equal(t, expectedInstanceUniqId, actualInstanceUniqId)
}

func TestNewDefaultRequestIdGenerator_UnixSocket(t *testing.T) {
	generator, _ := NewDefaultRequestIdGenerator(&settings.Config{HttpAddr: "1.2.3.4", UnixSocketPath: "/var/lib/ezbookkeeping/ezbookkeeping.sock", Protocol: "socket", SecretKey: "secretkey"})
	requestId := generator.GenerateRequestId("127.0.0.1", 20000)
	requestIdInfo := generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedServerUniqId := uint16(0x5bdb) // crc32("1.2.3.4" + "_" + "secretkey") & 0xFFFF
	actualServerUniqId := requestIdInfo.ServerUniqId
	assert.Equal(t, expectedServerUniqId, actualServerUniqId)

	expectedInstanceUniqId := uint16(0x2cc) // crc32("/var/lib/ezbookkeeping/ezbookkeeping.sock" + "_" + "secretkey") & 0xFFFF
	actualInstanceUniqId := requestIdInfo.InstanceUniqId
	assert.Equal(t, expectedInstanceUniqId, actualInstanceUniqId)
}

func TestNewDefaultRequestIdGenerator_ClientIpv4(t *testing.T) {
	generator, _ := NewDefaultRequestIdGenerator(&settings.Config{HttpAddr: "1.2.3.4", UnixSocketPath: "/var/lib/ezbookkeeping/ezbookkeeping.sock", Protocol: "socket", SecretKey: "secretkey"})
	requestId := generator.GenerateRequestId("127.0.0.1", 20000)
	requestIdInfo := generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientIp := uint32(0x7f000001) // 127.0.0.1
	actualClientIp := requestIdInfo.ClientIp
	assert.Equal(t, expectedClientIp, actualClientIp)

	expectedClientIpv6 := false
	actualClientIpv6 := requestIdInfo.IsClientIpv6
	assert.Equal(t, expectedClientIpv6, actualClientIpv6)

	requestId = generator.GenerateRequestId("192.168.1.100", 20000)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientIp = uint32(0xc0a80164) // 192.168.1.100
	actualClientIp = requestIdInfo.ClientIp
	assert.Equal(t, expectedClientIp, actualClientIp)

	expectedClientIpv6 = false
	actualClientIpv6 = requestIdInfo.IsClientIpv6
	assert.Equal(t, expectedClientIpv6, actualClientIpv6)
}

func TestNewDefaultRequestIdGenerator_ClientIpv6(t *testing.T) {
	generator, _ := NewDefaultRequestIdGenerator(&settings.Config{HttpAddr: "1.2.3.4", UnixSocketPath: "/var/lib/ezbookkeeping/ezbookkeeping.sock", Protocol: "socket", SecretKey: "secretkey"})
	requestId := generator.GenerateRequestId("2001:abc:def:1234::1", 20000)
	requestIdInfo := generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientIp := uint32(0x76fe1b98) // crc32("2001:abc:def:1234::1")
	actualClientIp := requestIdInfo.ClientIp
	assert.Equal(t, expectedClientIp, actualClientIp)

	expectedClientIpv6 := true
	actualClientIpv6 := requestIdInfo.IsClientIpv6
	assert.Equal(t, expectedClientIpv6, actualClientIpv6)

	requestId = generator.GenerateRequestId("2400:abcd:1234:1:56ef:ab78:c90d:1e2f", 20000)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientIp = uint32(0xa0a25faa) // crc32("2400:abcd:1234:1:56ef:ab78:c90d:1e2f")
	actualClientIp = requestIdInfo.ClientIp
	assert.Equal(t, expectedClientIp, actualClientIp)

	expectedClientIpv6 = true
	actualClientIpv6 = requestIdInfo.IsClientIpv6
	assert.Equal(t, expectedClientIpv6, actualClientIpv6)
}

func TestNewDefaultRequestIdGenerator_ClientPort(t *testing.T) {
	generator, _ := NewDefaultRequestIdGenerator(&settings.Config{HttpAddr: "1.2.3.4", UnixSocketPath: "/var/lib/ezbookkeeping/ezbookkeeping.sock", Protocol: "socket", SecretKey: "secretkey"})
	requestId := generator.GenerateRequestId("127.0.0.1", 0)
	requestIdInfo := generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientPort := uint16(0)
	actualClientPort := requestIdInfo.ClientPort
	assert.Equal(t, expectedClientPort, actualClientPort)

	requestId = generator.GenerateRequestId("127.0.0.1", 12345)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientPort = uint16(12345)
	actualClientPort = requestIdInfo.ClientPort
	assert.Equal(t, expectedClientPort, actualClientPort)

	requestId = generator.GenerateRequestId("127.0.0.1", 32767)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientPort = uint16(32767)
	actualClientPort = requestIdInfo.ClientPort
	assert.Equal(t, expectedClientPort, actualClientPort)

	requestId = generator.GenerateRequestId("127.0.0.1", 32768)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientPort = uint16(32768)
	actualClientPort = requestIdInfo.ClientPort
	assert.Equal(t, expectedClientPort, actualClientPort)

	requestId = generator.GenerateRequestId("127.0.0.1", 56789)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientPort = uint16(56789)
	actualClientPort = requestIdInfo.ClientPort
	assert.Equal(t, expectedClientPort, actualClientPort)

	requestId = generator.GenerateRequestId("127.0.0.1", 65535)
	requestIdInfo = generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

	expectedClientPort = uint16(65535)
	actualClientPort = requestIdInfo.ClientPort
	assert.Equal(t, expectedClientPort, actualClientPort)
}

func TestGenerateRequestId_100Times(t *testing.T) {
	generator, _ := NewDefaultRequestIdGenerator(&settings.Config{HttpAddr: "1.2.3.4", HttpPort: 1234, SecretKey: "secretkey"})

	for i := 1; i <= 100; i++ {
		requestId := generator.GenerateRequestId("127.0.0.1", 20000)
		requestIdInfo := generator.parseRequestIdInfo(generator.parseRequestIdFromUuid(requestId))

		expectedRequestSeqId := uint32(i)
		actualRequestSeqId := requestIdInfo.RequestSeqId
		assert.Equal(t, expectedRequestSeqId, actualRequestSeqId)
	}
}
