package uuid

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func TestGenerateUuid(t *testing.T) {
	expectedUnixTime := time.Now().Unix()
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_TRANSACTION

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

	uuid = generator.GenerateUuid(UUID_TYPE_TRANSACTION)
	uuidInfo = generator.ParseUuidInfo(uuid)
	actualSeqId = uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)
}

func TestGenerateUuid_2000TimesIn2Seconds(t *testing.T) {
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

func TestGenerateUuid_10000TimesConcurrent(t *testing.T) {
	concurrentCount := 10
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 3})
	var mutex sync.Mutex
	var generatedIds sync.Map
	var waitGroup sync.WaitGroup

	for routineIndex := 0; routineIndex < concurrentCount; routineIndex++ {
		go func(currentRoutineIndex int) {
			waitGroup.Add(1)

			for cycle := 0; cycle < 1000; cycle++ {
				expectedUnixTime := time.Now().Unix()
				uuid := generator.GenerateUuid(UUID_TYPE_USER)
				uuidInfo := generator.ParseUuidInfo(uuid)

				if uint32(expectedUnixTime) != uuidInfo.UnixTime {
					mutex.Lock()
					assert.Equal(t, uint32(expectedUnixTime), uuidInfo.UnixTime)
					mutex.Unlock()
				}

				if uuidInfo.SequentialId == 0 {
					if existedRoutineIndex, exists := generatedIds.Load(uuid); exists {
						mutex.Lock()
						assert.Fail(t, fmt.Sprintf("uuid \"%d\" conflicts, unix time is %d, seq id is %d, existed routine index is %d, current routine index is %d", uuid, uuidInfo.UnixTime, uuidInfo.SequentialId, existedRoutineIndex, currentRoutineIndex))
						mutex.Unlock()
					}

					generatedIds.Store(uuid, currentRoutineIndex)
				}
			}

			waitGroup.Done()
		}(routineIndex)
	}

	waitGroup.Wait()
}

func TestGenerateUuids_Count0(t *testing.T) {
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_TRANSACTION

	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: expectedUuidServerId})
	uuids := generator.GenerateUuids(expectedUuidType, 0)

	assert.NotEqual(t, nil, uuids)
	assert.Equal(t, 0, len(uuids))
}

func TestGenerateUuids_Count255(t *testing.T) {
	expectedUnixTime := time.Now().Unix()
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_TRANSACTION
	expectedUuidCount := uint8(255)

	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: expectedUuidServerId})
	uuids := generator.GenerateUuids(expectedUuidType, expectedUuidCount)

	for i := 0; i < int(expectedUuidCount); i++ {
		uuidInfo := generator.ParseUuidInfo(uuids[i])

		actualUnixTime := uuidInfo.UnixTime
		assert.Equal(t, uint32(expectedUnixTime), actualUnixTime)

		actualUuidServerId := uuidInfo.UuidServerId
		assert.Equal(t, expectedUuidServerId, actualUuidServerId)

		actualUuidType := uuidInfo.UuidType
		assert.Equal(t, uint8(expectedUuidType), actualUuidType)

		expectedSeqId := i
		actualSeqId := uuidInfo.SequentialId
		assert.Equal(t, uint32(expectedSeqId), actualSeqId)
	}

	assert.Equal(t, int(expectedUuidCount), len(uuids))
}

func TestGenerateUuids_30TimesIn3Seconds(t *testing.T) {
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_TRANSACTION
	expectedUuidCount := uint8(255)

	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: expectedUuidServerId})

	for cycle := 0; cycle < 30; cycle++ {
		expectedUnixTime := time.Now().Unix()
		uuids := generator.GenerateUuids(expectedUuidType, expectedUuidCount)
		var firstSeqId uint32

		for i := 0; i < int(expectedUuidCount); i++ {
			uuidInfo := generator.ParseUuidInfo(uuids[i])

			actualUnixTime := uuidInfo.UnixTime
			assert.Equal(t, uint32(expectedUnixTime), actualUnixTime)

			actualUuidServerId := uuidInfo.UuidServerId
			assert.Equal(t, expectedUuidServerId, actualUuidServerId)

			actualUuidType := uuidInfo.UuidType
			assert.Equal(t, uint8(expectedUuidType), actualUuidType)

			if i == 0 {
				firstSeqId = uuidInfo.SequentialId
			} else {
				expectedSeqId := firstSeqId + uint32(i)
				actualSeqId := uuidInfo.SequentialId
				assert.Equal(t, expectedSeqId, actualSeqId)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestGenerateUuids_20000TimesConcurrent(t *testing.T) {
	concurrentCount := 10
	expectedUuidCount := uint8(20)
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 3})
	var mutex sync.Mutex
	var generatedIds sync.Map
	var waitGroup sync.WaitGroup

	for routineIndex := 0; routineIndex < concurrentCount; routineIndex++ {
		go func(currentRoutineIndex int) {
			waitGroup.Add(1)

			for cycle := 0; cycle < 100; cycle++ {
				expectedUnixTime := time.Now().Unix()
				uuids := generator.GenerateUuids(UUID_TYPE_USER, expectedUuidCount)

				for i := 0; i < int(expectedUuidCount); i++ {
					uuidInfo := generator.ParseUuidInfo(uuids[i])

					if uint32(expectedUnixTime) != uuidInfo.UnixTime {
						mutex.Lock()
						assert.Equal(t, uint32(expectedUnixTime), uuidInfo.UnixTime)
						mutex.Unlock()
					}

					if existedRoutineIndex, exists := generatedIds.Load(uuids[i]); exists {
						mutex.Lock()
						assert.Fail(t, fmt.Sprintf("uuid \"%d\" conflicts, unix time is %d, seq id is %d, existed routine index is %d, current routine index is %d", uuids[i], uuidInfo.UnixTime, uuidInfo.SequentialId, existedRoutineIndex, currentRoutineIndex))
						mutex.Unlock()
					}

					generatedIds.Store(uuids[i], currentRoutineIndex)
				}
			}

			waitGroup.Done()
		}(routineIndex)
	}

	waitGroup.Wait()
}
