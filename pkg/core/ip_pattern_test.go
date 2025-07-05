package core

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestIPPattern_GobEncode(t *testing.T) {
	pattern, err := ParseIPPattern("192.168.1.*")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(pattern)
	assert.Nil(t, err)

	newPattern := &IPPattern{}
	err = gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(newPattern)
	assert.Nil(t, err)
	assert.NotNil(t, newPattern)

	assert.Equal(t, pattern.Pattern, newPattern.Pattern)
	assert.Equal(t, pattern.regex.String(), newPattern.regex.String())

	assert.True(t, newPattern.Match("192.168.1.1"))
	assert.True(t, newPattern.Match("192.168.1.255"))
}

func TestParseIPPattern(t *testing.T) {
	pattern, err := ParseIPPattern("")
	assert.Nil(t, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPPattern("invalid")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPPattern("192.1:2:3.4")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPPattern("0:0:0:0:0:0:1.2.3.4") // not support IPv6 with embedded IPv4
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPPattern("192.168.1.*")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("192.168.1.1"))
	assert.True(t, pattern.Match("192.168.1.255"))
	assert.False(t, pattern.Match("192.168.2.1"))

	pattern, err = ParseIPPattern("2001:db8::*")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("2001:db8::1"))
	assert.True(t, pattern.Match("2001:db8::ffff"))
	assert.False(t, pattern.Match("2001:db9::1"))
}

func TestParseIPv4Pattern(t *testing.T) {
	pattern, err := ParseIPv4Pattern("192.168.1.1")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("192.168.1.1"))
	assert.False(t, pattern.Match("192.168.1.2"))

	pattern, err = ParseIPv4Pattern("192.168.*.1")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("192.168.1.1"))
	assert.True(t, pattern.Match("192.168.255.1"))
	assert.False(t, pattern.Match("192.168.1.2"))

	pattern, err = ParseIPv4Pattern("*.*.*.*")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("0.0.0.0"))
	assert.True(t, pattern.Match("255.255.255.255"))

	pattern, err = ParseIPv4Pattern("256.256.256.256")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPv4Pattern("1.2.3")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPv4Pattern("1.2.3.4.5")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPv4Pattern("a.b.c.d")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)
}

func TestParseIPv6Pattern(t *testing.T) {
	pattern, err := ParseIPv6Pattern("2001:db8:85a3:8d3:1319:8a2e:370:7348")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
	assert.False(t, pattern.Match("2001:db8:85a3:8d3:1319:8a2e:370:7349"))

	pattern, err = ParseIPv6Pattern("2001:db8::*")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("2001:db8::0"))
	assert.True(t, pattern.Match("2001:db8::ffff"))
	assert.False(t, pattern.Match("2001:db9::0"))

	pattern, err = ParseIPv6Pattern("::*")
	assert.Nil(t, err)
	assert.NotNil(t, pattern)
	assert.True(t, pattern.Match("::1"))
	assert.True(t, pattern.Match("::2"))
	assert.False(t, pattern.Match(":1:1"))

	pattern, err = ParseIPv6Pattern("2001:db8:85a3:8d3:1319:8a2e:370:7348:extra")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPv6Pattern("g001:db8:85a3:8d3")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)

	pattern, err = ParseIPv6Pattern("2001:db8:")
	assert.Equal(t, errs.ErrInvalidIpAddressPattern, err)
	assert.Nil(t, pattern)
}
