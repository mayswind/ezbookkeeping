package uuid

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/lab/pkg/settings"
)

func TestGenerateUuid(t *testing.T) {
	expectedUnixTime := time.Now().Unix()
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_JOURNAL

	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: expectedUuidServerId})
	uuid := generator.GenerateUuid(expectedUuidType)
	uuidInfo := generator.ParseUuidInfo(uuid)

	actualUnixTime := uuidInfo.UnixTime
	assert.Equal(t, uint32(expectedUnixTime), actualUnixTime)

	actualUuidServerId := uuidInfo.UuidServerId
	assert.Equal(t, expectedUuidServerId, actualUuidServerId)

	actualUuidType := uuidInfo.UuidType
	assert.Equal(t, uint8(expectedUuidType), actualUuidType)

	expectedSeqId := 0
	actualSeqId := uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)
}

func TestGenerateUuid_MultiType(t *testing.T) {
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 1})
	uuid := generator.GenerateUuid(UUID_TYPE_USER)
	uuidInfo := generator.ParseUuidInfo(uuid)

	expectedSeqId := 0
	actualSeqId := uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)

	uuid = generator.GenerateUuid(UUID_TYPE_ACCOUNT)
	uuidInfo = generator.ParseUuidInfo(uuid)
	actualSeqId = uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)

	uuid = generator.GenerateUuid(UUID_TYPE_JOURNAL)
	uuidInfo = generator.ParseUuidInfo(uuid)
	actualSeqId = uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)
}

func TestGenerateUuid_1000Times(t *testing.T) {
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 2})
	expectedUnixTime := time.Now().Unix()

	for i := 0; i < 1000; i++ {
		uuid := generator.GenerateUuid(UUID_TYPE_USER)
		uuidInfo := generator.ParseUuidInfo(uuid)

		assert.Equal(t, uint32(expectedUnixTime), uuidInfo.UnixTime)
		assert.Equal(t, uint32(i), uuidInfo.SequentialId)
	}

	time.Sleep(1 * time.Second)
	expectedUnixTime = time.Now().Unix()

	for i := 0; i < 1000; i++ {
		uuid := generator.GenerateUuid(UUID_TYPE_USER)
		uuidInfo := generator.ParseUuidInfo(uuid)

		assert.Equal(t, uint32(expectedUnixTime), uuidInfo.UnixTime)
		assert.Equal(t, uint32(i), uuidInfo.SequentialId)
	}
}

func TestGenerateUuid_1000000TimesConcurrent(t *testing.T) {
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 3})
	var mutex sync.Mutex
	var generatedIds sync.Map
	var waitGroup sync.WaitGroup

	for i := 0; i < 50; i++ {
		go func() {
			waitGroup.Add(1)

			for j := 0; j < 40000; j++ {
				if j%10000 == 0 { // echo server can only generate 500,000 (50 * 10000) uuids in one second
					time.Sleep(1000 * time.Millisecond)
				}

				expectedUnixTime := time.Now().Unix()
				uuid := generator.GenerateUuid(UUID_TYPE_USER)
				uuidInfo := generator.ParseUuidInfo(uuid)

				if uint32(expectedUnixTime) != uuidInfo.UnixTime {
					mutex.Lock()
					assert.Equal(t, uint32(expectedUnixTime), uuidInfo.UnixTime)
					mutex.Unlock()
				}

				if uuidInfo.SequentialId == 0 {
					if existedUnixTime, exists := generatedIds.Load(uuid); exists {
						mutex.Lock()
						assert.Fail(t, fmt.Sprintf("uuid \"%d\" conflicts, seq id is %d, existed unixtime is %d, current unix time is %d", uuid, uuidInfo.SequentialId, existedUnixTime, uuidInfo.UnixTime))
						mutex.Unlock()
					}

					generatedIds.Store(uuid, uuidInfo.UnixTime)
				}
			}

			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
}
