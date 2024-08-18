package requestid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net"
	"sync/atomic"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// Length and mask of all information in request id
const (
	requestIdLength               = 36
	secondsTodayBits              = 17
	secondsTodayBitsMask          = (1 << secondsTodayBits) - 1
	clientPortNumberAllBits       = 16
	clientPortNumberHigh1Bit      = 1
	clientPortNumberLow15Bits     = clientPortNumberAllBits - clientPortNumberHigh1Bit
	clientPortNumberHigh1BitMask  = 1 << clientPortNumberLow15Bits
	clientPortNumberLow15BitsMask = clientPortNumberHigh1BitMask - 1
	reqSeqNumberBits              = 30
	reqSeqNumberBitsMask          = (1 << reqSeqNumberBits) - 1
	clientIpv6Bit                 = 1
	clientIpv6BitMask             = 1
)

// RequestIdInfo represents a struct which has all information in request id
type RequestIdInfo struct {
	ServerUniqId        uint16
	InstanceUniqId      uint16
	SecondsElapsedToday uint32
	RequestSeqId        uint32
	IsClientIpv6        bool
	ClientIp            uint32
	ClientPort          uint16
}

// DefaultRequestIdGenerator represents default request id generator
type DefaultRequestIdGenerator struct {
	serverUniqId   uint16
	instanceUniqId uint16
	requestSeqId   atomic.Uint32
}

// NewDefaultRequestIdGenerator returns a new default request id generator
func NewDefaultRequestIdGenerator(c core.Context, config *settings.Config) (*DefaultRequestIdGenerator, error) {
	serverUniqId, err := getServerUniqId(c, config)

	if err != nil {
		return nil, err
	}

	instanceUniqId := getInstanceUniqId(config)

	generator := &DefaultRequestIdGenerator{
		serverUniqId:   serverUniqId,
		instanceUniqId: instanceUniqId,
	}

	return generator, nil
}

func getServerUniqId(c core.Context, config *settings.Config) (uint16, error) {
	localAddr := ""
	settingAddr := net.ParseIP(config.HttpAddr)

	if settingAddr != nil && !settingAddr.IsUnspecified() {
		localAddr = settingAddr.String()
	} else {
		var err error
		localAddr, err = utils.GetLocalIPAddressesString()

		if err != nil {
			log.Warnf(c, "[default_request_id_generator.getServerUniqId] failed to get local ipv4 address, because %s", err.Error())
			return 0, err
		}
	}

	serverUniqFlag := fmt.Sprintf("%s_%s", localAddr, config.SecretKey)

	return uint16(crc32.ChecksumIEEE([]byte(serverUniqFlag))), nil
}

func getInstanceUniqId(config *settings.Config) uint16 {
	var instanceUniqFlag string

	if config.Protocol == settings.SCHEME_SOCKET {
		instanceUniqFlag = fmt.Sprintf("%s_%s", config.UnixSocketPath, config.SecretKey)
	} else {
		instanceUniqFlag = fmt.Sprintf("%d_%s", config.HttpPort, config.SecretKey)
	}

	return uint16(crc32.ChecksumIEEE([]byte(instanceUniqFlag)))

}

// ParseRequestIdInfo returns a info struct which contains all information in request id
func (r *DefaultRequestIdGenerator) ParseRequestIdInfo(requestId string) (*RequestIdInfo, error) {
	if requestId == "" || len(requestId) != requestIdLength {
		return nil, errs.ErrRequestIdInvalid
	}

	requestIdData := r.parseRequestIdFromUuid(requestId)
	return r.parseRequestIdInfo(requestIdData), nil
}

// GetCurrentServerUniqId returns current server unique id
func (r *DefaultRequestIdGenerator) GetCurrentServerUniqId() uint16 {
	return r.serverUniqId
}

// GetCurrentInstanceUniqId returns current application instance unique id
func (r *DefaultRequestIdGenerator) GetCurrentInstanceUniqId() uint16 {
	return r.instanceUniqId
}

// GenerateRequestId returns a new request id
func (r *DefaultRequestIdGenerator) GenerateRequestId(clientIpAddr string, clientPort uint16) string {
	ip := net.ParseIP(clientIpAddr)
	isClientIpv6 := ip.To4() == nil
	var clientIp uint32

	if isClientIpv6 {
		clientIp = crc32.ChecksumIEEE([]byte(ip.String()))
	} else {
		clientIp = binary.BigEndian.Uint32(ip.To4())
	}

	requestId := r.getRequestId(r.serverUniqId, r.instanceUniqId, isClientIpv6, clientIp, clientPort)

	return requestId
}

func (r *DefaultRequestIdGenerator) getRequestId(serverUniqId uint16, instanceUniqId uint16, clientIpV6 bool, clientIp uint32, clientPort uint16) string {
	clientIpv6Flag := uint32(0)

	if clientIpV6 {
		clientIpv6Flag = uint32(1)
	}

	// 128bits = serverUniqId(16bits) + instanceUniqId(16bits) + secondsElapsedToday(17bits) + clientPortLow15Bits(15bits) + sequentialNumber(30bits) + clientPortHigh1Bit(1bit) + clientIpv6Flag(1bit) + clientIp(32bits)

	secondsElapsedToday := r.getSecondsElapsedToday()
	secondsLow17bits := uint32(secondsElapsedToday & secondsTodayBitsMask)

	clientPortHigh1bit := uint32((clientPort & clientPortNumberHigh1BitMask) >> clientPortNumberLow15Bits)
	clientPortLow15bits := uint32(clientPort & clientPortNumberLow15BitsMask)

	secondsAndClientPortLowBits := (secondsLow17bits << clientPortNumberLow15Bits) | clientPortLow15bits

	seqId := r.requestSeqId.Add(1)
	seqIdLow30bits := seqId & reqSeqNumberBitsMask

	seqIdAndClientPortHighBitAndClientIpv6Flag := (seqIdLow30bits << (clientPortNumberHigh1Bit + clientIpv6Bit)) | (clientPortHigh1bit << clientPortNumberHigh1Bit) | (clientIpv6Flag & clientIpv6BitMask)

	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, serverUniqId)
	_ = binary.Write(buf, binary.BigEndian, instanceUniqId)
	_ = binary.Write(buf, binary.BigEndian, secondsAndClientPortLowBits)
	_ = binary.Write(buf, binary.BigEndian, seqIdAndClientPortHighBitAndClientIpv6Flag)
	_ = binary.Write(buf, binary.BigEndian, clientIp)

	return r.getUuidFromRequestId(buf)
}

func (r *DefaultRequestIdGenerator) getSecondsElapsedToday() int {
	now := time.Now()
	seconds := now.Hour()*60*60 + now.Minute()*60 + now.Second()

	return seconds
}

func (r *DefaultRequestIdGenerator) getUuidFromRequestId(buffer *bytes.Buffer) string {
	data := buffer.Bytes()
	result := make([]byte, 36)

	hex.Encode(result[0:8], data[0:4])
	result[8] = '-'
	hex.Encode(result[9:13], data[4:6])
	result[13] = '-'
	hex.Encode(result[14:18], data[6:8])
	result[18] = '-'
	hex.Encode(result[19:23], data[8:10])
	result[23] = '-'
	hex.Encode(result[24:], data[10:])

	return string(result)
}

func (r *DefaultRequestIdGenerator) parseRequestIdInfo(data []byte) *RequestIdInfo {
	buf := bytes.NewBuffer(data)

	var serverUniqId uint16
	var instanceUniqId uint16
	var secondsAndClientPortLowBits uint32
	var seqIdAndClientPortHighBitAndClientIpv6Flag uint32
	var clientIp uint32

	_ = binary.Read(buf, binary.BigEndian, &serverUniqId)
	_ = binary.Read(buf, binary.BigEndian, &instanceUniqId)
	_ = binary.Read(buf, binary.BigEndian, &secondsAndClientPortLowBits)
	_ = binary.Read(buf, binary.BigEndian, &seqIdAndClientPortHighBitAndClientIpv6Flag)
	_ = binary.Read(buf, binary.BigEndian, &clientIp)

	secondsElapsedToday := (secondsAndClientPortLowBits >> clientPortNumberLow15Bits) & secondsTodayBitsMask

	seqId := (seqIdAndClientPortHighBitAndClientIpv6Flag >> (clientPortNumberHigh1Bit + clientIpv6Bit)) & reqSeqNumberBitsMask
	clientPortHigh1Bit := ((seqIdAndClientPortHighBitAndClientIpv6Flag >> clientIpv6Bit) << clientPortNumberLow15Bits) & clientPortNumberHigh1BitMask
	isClientIpv6Flag := seqIdAndClientPortHighBitAndClientIpv6Flag & clientIpv6BitMask
	isClientIpv6 := false

	clientPort := uint16(clientPortHigh1Bit | (secondsAndClientPortLowBits & clientPortNumberLow15BitsMask))

	if isClientIpv6Flag == 1 {
		isClientIpv6 = true
	}

	return &RequestIdInfo{
		ServerUniqId:        serverUniqId,
		InstanceUniqId:      instanceUniqId,
		SecondsElapsedToday: secondsElapsedToday,
		RequestSeqId:        seqId,
		IsClientIpv6:        isClientIpv6,
		ClientIp:            clientIp,
		ClientPort:          clientPort,
	}
}

func (r *DefaultRequestIdGenerator) parseRequestIdFromUuid(uuid string) []byte {
	data := []byte(uuid)
	result := make([]byte, 16)

	_, _ = hex.Decode(result[0:4], data[0:8])
	_, _ = hex.Decode(result[4:6], data[9:13])
	_, _ = hex.Decode(result[6:8], data[14:18])
	_, _ = hex.Decode(result[8:10], data[19:23])
	_, _ = hex.Decode(result[10:], data[24:])

	return result
}
