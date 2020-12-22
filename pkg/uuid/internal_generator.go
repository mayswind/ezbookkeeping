package uuid

import (
	"sync/atomic"
	"time"

	"github.com/mayswind/lab/pkg/settings"
)

// Length and mask of all information in uuid
const (
	internalUuidUnixTimeBits = 32
	internalUuidUnixTimeMask = (1 << internalUuidUnixTimeBits) - 1

	internalUuidTypeBits = 4
	internalUuidTypeMask = (1 << internalUuidTypeBits) - 1

	internalUuidServerIdBits = 8
	internalUuidServerIdMask = (1 << internalUuidServerIdBits) - 1

	internalUuidSeqIdBits = 19
	internalUuidSeqIdMask = (1 << internalUuidSeqIdBits) - 1

	seqNumberIdBits = 32
	seqNumberIdMask = (1 << seqNumberIdBits) - 1
)

// InternalUuidInfo represents a struct which has all information in uuid
type InternalUuidInfo struct {
	UnixTime     uint32
	UuidType     uint8
	UuidServerId uint8
	SequentialId uint32
}

// InternalUuidGenerator represents internal bundled uuid generator
type InternalUuidGenerator struct {
	uuidServerId   uint8
	uuidSeqNumbers [1 << internalUuidTypeBits]uint64
}

// NewInternalUuidGenerator returns a new internal uuid generator
func NewInternalUuidGenerator(config *settings.Config) (*InternalUuidGenerator, error) {
	generator := &InternalUuidGenerator{
		uuidServerId: config.UuidServerId,
	}

	return generator, nil
}

// GenerateUuid returns a new uuid
func (u *InternalUuidGenerator) GenerateUuid(idType UuidType) int64 {
	// 63bits = unixTime(32bits) + uuidType(4bits) + uuidServerId(8bits) + sequentialNumber(19bits)

	var unixTime uint64
	var newSeqId uint64
	uuidType := uint8(idType)

	for {
		unixTime = uint64(time.Now().Unix())
		newSeqId = atomic.AddUint64(&u.uuidSeqNumbers[uuidType], 1)

		if newSeqId>>seqNumberIdBits == unixTime {
			break
		}

		currentSeqId := newSeqId
		newSeqId = unixTime << seqNumberIdBits

		if atomic.CompareAndSwapUint64(&u.uuidSeqNumbers[uuidType], currentSeqId, newSeqId) {
			break
		}
	}

	seqId := newSeqId & seqNumberIdMask

	unixTimePart := (int64(unixTime) & internalUuidUnixTimeMask) << (internalUuidTypeBits + internalUuidServerIdBits + internalUuidSeqIdBits)
	uuidTypePart := (int64(uuidType) & internalUuidTypeMask) << (internalUuidServerIdBits + internalUuidSeqIdBits)
	uuidServerIdPart := (int64(u.uuidServerId) & internalUuidServerIdMask) << internalUuidSeqIdBits
	seqIdPart := int64(seqId) & internalUuidSeqIdMask

	uuid := unixTimePart | uuidTypePart | uuidServerIdPart | seqIdPart

	return uuid
}

// ParseUuidInfo returns a info struct which contains all information in uuid
func (u *InternalUuidGenerator) ParseUuidInfo(uuid int64) *InternalUuidInfo {
	seqId := uint32(uuid & internalUuidSeqIdMask)
	uuid = uuid >> internalUuidSeqIdBits

	uuidServerId := uint8(uuid & internalUuidServerIdMask)
	uuid = uuid >> internalUuidServerIdBits

	uuidType := uint8(uuid & internalUuidTypeMask)
	uuid = uuid >> internalUuidTypeBits

	unixTime := uint32(uuid & internalUuidUnixTimeMask)

	return &InternalUuidInfo{
		UnixTime:     unixTime,
		UuidType:     uuidType,
		UuidServerId: uuidServerId,
		SequentialId: seqId,
	}
}
