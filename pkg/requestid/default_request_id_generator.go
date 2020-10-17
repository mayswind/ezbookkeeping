package requestid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/utils"
)

const REQUEST_ID_LENGTH = 36
const SECONDS_TODAY_BITS = 17
const SECONDS_TODAY_BITS_MASK = (1 << SECONDS_TODAY_BITS) - 1
const RANDOM_NUMBER_BITS = 15
const RANDOM_NUMBER_BITS_MASK = (1 << RANDOM_NUMBER_BITS) - 1
const REQ_SEQ_NUMBER_BITS = 31
const REQ_SEQ_NUMBER_BITS_MASK = (1 << REQ_SEQ_NUMBER_BITS) - 1
const CLIENT_IPV6_BIT = 1
const CLIENT_IPV6_BIT_MASK = 1

type RequestIdInfo struct {
	ServerUniqId        uint16
	InstanceUniqId      uint16
	SecondsElapsedToday uint32
	RandomNumber        uint32
	RequestSeqId        uint32
	IsClientIpv6        bool
	ClientIp            uint32
}

type DefaultRequestIdGenerator struct {
	serverUniqId   uint16
	instanceUniqId uint16
	requestSeqId   uint32
}

func NewDefaultRequestIdGenerator(config *settings.Config) (*DefaultRequestIdGenerator, error) {
	serverUniqId, err := getServerUniqId(config)

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

func getServerUniqId(config *settings.Config) (uint16, error) {
	localAddr := ""
	settingAddr := net.ParseIP(config.HttpAddr)

	if settingAddr != nil && !settingAddr.IsUnspecified() {
		localAddr = settingAddr.String()
	} else {
		var err error
		localAddr, err = utils.GetLocalIPAddressesString()

		if err != nil {
			log.Warnf("[default_request_id_generator.getServerUniqId] failed to get local ipv4 address, because %s", err.Error())
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

func (r *DefaultRequestIdGenerator) ParseRequestIdInfo(requestId string) (*RequestIdInfo, error) {
	if requestId == "" || len(requestId) != REQUEST_ID_LENGTH {
		return nil, errs.ErrRequestIdInvalid
	}

	requestIdData := r.parseRequestIdFromUuid(requestId)
	return r.parseRequestIdInfo(requestIdData), nil
}

func (r *DefaultRequestIdGenerator) GetCurrentServerUniqId() uint16 {
	return r.serverUniqId
}

func (r *DefaultRequestIdGenerator) GetCurrentInstanceUniqId() uint16 {
	return r.instanceUniqId
}

func (r *DefaultRequestIdGenerator) GenerateRequestId(clientIpAddr string) string {
	ip := net.ParseIP(clientIpAddr)
	isClientIpv6 := ip.To4() == nil
	var clientIp uint32

	if isClientIpv6 {
		clientIp = crc32.ChecksumIEEE([]byte(ip.String()))
	} else {
		clientIp = binary.BigEndian.Uint32(ip.To4())
	}

	requestId := r.getRequestId(r.serverUniqId, r.instanceUniqId, isClientIpv6, clientIp)

	return requestId
}

func (r *DefaultRequestIdGenerator) getRequestId(serverUniqId uint16, instanceUniqId uint16, clientIpV6 bool, clientIp uint32) string {
	clientIpv6Flag := uint32(0)

	if clientIpV6 {
		clientIpv6Flag = uint32(1)
	}

	// 128bits = serverUniqId(16bits) + instanceUniqId(16bits) + secondsElapsedToday(17bits) + randomNumber(15bits) + sequentialNumber(31bits) + clientIpv6Flag(1bit) + clientIp(32bits)

	secondsElapsedToday := r.getSecondsElapsedToday()
	secondsLow17bits := uint32(secondsElapsedToday & SECONDS_TODAY_BITS_MASK)

	randomNumber, _ := utils.GetRandomInteger(math.MaxInt16)
	randomNumberLow15bits := uint32(randomNumber & RANDOM_NUMBER_BITS_MASK)

	secondsAndRandomNumber := (secondsLow17bits << RANDOM_NUMBER_BITS) | randomNumberLow15bits

	seqId := atomic.AddUint32(&r.requestSeqId, 1)
	seqIdLow31bits := seqId & REQ_SEQ_NUMBER_BITS_MASK

	seqIdAndClientIpv6Flag := (seqIdLow31bits << CLIENT_IPV6_BIT) | (clientIpv6Flag & CLIENT_IPV6_BIT_MASK)

	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.BigEndian, serverUniqId)
	_ = binary.Write(buf, binary.BigEndian, instanceUniqId)
	_ = binary.Write(buf, binary.BigEndian, secondsAndRandomNumber)
	_ = binary.Write(buf, binary.BigEndian, seqIdAndClientIpv6Flag)
	_ = binary.Write(buf, binary.BigEndian, clientIp)

	return r.getUuidFromRequestId(buf)
}

func (r *DefaultRequestIdGenerator) getSecondsElapsedToday() int {
	now := time.Now()
	seconds := now.Hour()*24*60 + now.Minute()*60 + now.Second()

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
	var secondsAndRandomNumber uint32
	var seqIdAndClientIpv6Flag uint32
	var clientIp uint32

	_ = binary.Read(buf, binary.BigEndian, &serverUniqId)
	_ = binary.Read(buf, binary.BigEndian, &instanceUniqId)
	_ = binary.Read(buf, binary.BigEndian, &secondsAndRandomNumber)
	_ = binary.Read(buf, binary.BigEndian, &seqIdAndClientIpv6Flag)
	_ = binary.Read(buf, binary.BigEndian, &clientIp)

	secondsElapsedToday := (secondsAndRandomNumber >> RANDOM_NUMBER_BITS) & SECONDS_TODAY_BITS_MASK
	randomNumber := (secondsAndRandomNumber & RANDOM_NUMBER_BITS_MASK)

	seqId := (seqIdAndClientIpv6Flag >> CLIENT_IPV6_BIT) & REQ_SEQ_NUMBER_BITS_MASK
	isClientIpv6Flag := (seqIdAndClientIpv6Flag & CLIENT_IPV6_BIT_MASK)
	isClientIpv6 := false

	if isClientIpv6Flag == 1 {
		isClientIpv6 = true
	}

	return &RequestIdInfo{
		ServerUniqId:        serverUniqId,
		InstanceUniqId:      instanceUniqId,
		SecondsElapsedToday: secondsElapsedToday,
		RequestSeqId:        seqId,
		RandomNumber:        randomNumber,
		IsClientIpv6:        isClientIpv6,
		ClientIp:            clientIp,
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
