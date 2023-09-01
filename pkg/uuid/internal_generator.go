package uuid

import (
	"sync/atomic"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
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
	uuidSeqNumbers [1 << internalUuidTypeBits]atomic.Uint64
	uuidServerId   uint8
}

// NewInternalUuidGenerator returns a new internal uuid generator
func NewInternalUuidGenerator(config *settings.Config) (*InternalUuidGenerator, error) {
	generator := &InternalUuidGenerator{
		uuidServerId: config.UuidServerId,
	}

	return generator, nil
}

// GenerateUuid generates a new uuid
func (u *InternalUuidGenerator) GenerateUuid(idType UuidType) int64 {
	uuids := u.GenerateUuids(idType, 1)
	return uuids[0]
}

// GenerateUuids generates new uuids
func (u *InternalUuidGenerator) GenerateUuids(idType UuidType, count uint8) []int64 {
	// 63bits = unixTime(32bits) + uuidType(4bits) + uuidServerId(8bits) + sequentialNumber(19bits)

	uuids := make([]int64, count)

	if count < 1 {
		return uuids
	}

	var unixTime uint64
	var newFirstSeqId uint64
	var newLastSeqId uint64
	uuidType := uint8(idType)

	for {
		unixTime = uint64(time.Now().Unix())
		newLastSeqId = u.uuidSeqNumbers[uuidType].Add(uint64(count))

		if newLastSeqId>>seqNumberIdBits == unixTime {
			newFirstSeqId = newLastSeqId - uint64(count-1)
			break
		}

		currentSeqId := newLastSeqId
		newFirstSeqId = unixTime << seqNumberIdBits
		newLastSeqId = newFirstSeqId + uint64(count-1)

		if u.uuidSeqNumbers[uuidType].CompareAndSwap(currentSeqId, newLastSeqId) {
			break
		}
	}

	for i := 0; i < int(count); i++ {
		seqId := (newFirstSeqId + uint64(i)) & seqNumberIdMask
		uuids[i] = u.assembleUuid(unixTime, uuidType, seqId)
	}

	return uuids
}

func (u *InternalUuidGenerator) parseInternalUuidInfo(uuid int64) *InternalUuidInfo {
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

func (u *InternalUuidGenerator) assembleUuid(unixTime uint64, uuidType uint8, seqId uint64) int64 {
	unixTimePart := (int64(unixTime) & internalUuidUnixTimeMask) << (internalUuidTypeBits + internalUuidServerIdBits + internalUuidSeqIdBits)
	uuidTypePart := (int64(uuidType) & internalUuidTypeMask) << (internalUuidServerIdBits + internalUuidSeqIdBits)
	uuidServerIdPart := (int64(u.uuidServerId) & internalUuidServerIdMask) << internalUuidSeqIdBits
	seqIdPart := int64(seqId) & internalUuidSeqIdMask

	return unixTimePart | uuidTypePart | uuidServerIdPart | seqIdPart
}
