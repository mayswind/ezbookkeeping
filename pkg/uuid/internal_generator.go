package uuid

import (
	"sync/atomic"
	"time"

	"github.com/mayswind/lab/pkg/settings"
)

const (
	INTERNAL_UUID_UNIX_TIME_BITS = 32
	INTERNAL_UUID_UNIX_TIME_MASK = (1 << INTERNAL_UUID_UNIX_TIME_BITS) - 1

	INTERNAL_UUID_TYPE_BITS = 4
	INTERNAL_UUID_TYPE_MASK = (1 << INTERNAL_UUID_TYPE_BITS) - 1

	INTERNAL_UUID_SERVER_ID_BITS = 8
	INTERNAL_UUID_SERVER_ID_MASK = (1 << INTERNAL_UUID_SERVER_ID_BITS) - 1

	INTERNAL_UUID_SEQ_ID_BITS = 19
	INTERNAL_UUID_SEQ_ID_MASK = (1 << INTERNAL_UUID_SEQ_ID_BITS) - 1

	SEQ_NUMBER_ID_BITS = 32
	SEQ_NUMBER_ID_MASK = (1 << SEQ_NUMBER_ID_BITS) - 1
)

type InternalUuidInfo struct {
	UnixTime     uint32
	UuidType     uint8
	UuidServerId uint8
	SequentialId uint32
}

type InternalUuidGenerator struct {
	uuidServerId   uint8
	uuidSeqNumbers [1 << INTERNAL_UUID_TYPE_BITS]uint64
}

func NewInternalUuidGenerator(config *settings.Config) (*InternalUuidGenerator, error) {
	generator := &InternalUuidGenerator{
		uuidServerId: config.UuidServerId,
	}

	return generator, nil
}

func (u *InternalUuidGenerator) GenerateUuid(idType UuidType) int64 {
	// 63bits = unixTime(32bits) + uuidType(4bits) + uuidServerId(8bits) + sequentialNumber(19bits)

	var unixTime uint64
	var newSeqId uint64
	uuidType := uint8(idType)

	for {
		unixTime = uint64(time.Now().Unix())
		newSeqId = atomic.AddUint64(&u.uuidSeqNumbers[uuidType], 1)

		if newSeqId>>SEQ_NUMBER_ID_BITS == unixTime {
			break
		}

		currentSeqId := newSeqId
		newSeqId = unixTime << SEQ_NUMBER_ID_BITS

		if atomic.CompareAndSwapUint64(&u.uuidSeqNumbers[uuidType], currentSeqId, newSeqId) {
			break
		}
	}

	seqId := newSeqId & SEQ_NUMBER_ID_MASK

	unixTimePart := (int64(unixTime) & INTERNAL_UUID_UNIX_TIME_MASK) << (INTERNAL_UUID_TYPE_BITS + INTERNAL_UUID_SERVER_ID_BITS + INTERNAL_UUID_SEQ_ID_BITS)
	uuidTypePart := (int64(uuidType) & INTERNAL_UUID_TYPE_MASK) << (INTERNAL_UUID_SERVER_ID_BITS + INTERNAL_UUID_SEQ_ID_BITS)
	uuidServerIdPart := (int64(u.uuidServerId) & INTERNAL_UUID_SERVER_ID_MASK) << INTERNAL_UUID_SEQ_ID_BITS
	seqIdPart := int64(seqId) & INTERNAL_UUID_SEQ_ID_MASK

	uuid := unixTimePart | uuidTypePart | uuidServerIdPart | seqIdPart

	return uuid
}

func (u *InternalUuidGenerator) ParseUuidInfo(uuid int64) *InternalUuidInfo {
	seqId := uint32(uuid & INTERNAL_UUID_SEQ_ID_MASK)
	uuid = uuid >> INTERNAL_UUID_SEQ_ID_BITS

	uuidServerId := uint8(uuid & INTERNAL_UUID_SERVER_ID_MASK)
	uuid = uuid >> INTERNAL_UUID_SERVER_ID_BITS

	uuidType := uint8(uuid & INTERNAL_UUID_TYPE_MASK)
	uuid = uuid >> INTERNAL_UUID_TYPE_BITS

	unixTime := uint32(uuid & INTERNAL_UUID_UNIX_TIME_MASK)

	return &InternalUuidInfo{
		UnixTime:     unixTime,
		UuidType:     uuidType,
		UuidServerId: uuidServerId,
		SequentialId: seqId,
	}
}
