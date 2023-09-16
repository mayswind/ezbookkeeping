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
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_TRANSACTION

	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: expectedUuidServerId})

	generationStartUnixTime := time.Now().Unix()
	uuid := generator.GenerateUuid(expectedUuidType)
	generationEndUnixTime := time.Now().Unix()

	uuidInfo := generator.parseInternalUuidInfo(uuid)

	if generationStartUnixTime == generationEndUnixTime {
		actualUnixTime := uuidInfo.UnixTime
		assert.Equal(t, uint32(generationStartUnixTime), actualUnixTime)
	}

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
	uuidInfo := generator.parseInternalUuidInfo(uuid)

	expectedSeqId := 0
	actualSeqId := uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)

	uuid = generator.GenerateUuid(UUID_TYPE_ACCOUNT)
	uuidInfo = generator.parseInternalUuidInfo(uuid)
	actualSeqId = uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)

	uuid = generator.GenerateUuid(UUID_TYPE_TRANSACTION)
	uuidInfo = generator.parseInternalUuidInfo(uuid)
	actualSeqId = uuidInfo.SequentialId
	assert.Equal(t, uint32(expectedSeqId), actualSeqId)
}

func TestGenerateUuid_2000TimesIn2Seconds(t *testing.T) {
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 2})

	for i := 0; i < 1000; i++ {
		generationStartUnixTime := time.Now().Unix()
		uuid := generator.GenerateUuid(UUID_TYPE_USER)
		generationEndUnixTime := time.Now().Unix()

		uuidInfo := generator.parseInternalUuidInfo(uuid)

		if generationStartUnixTime == generationEndUnixTime {
			assert.Equal(t, uint32(generationStartUnixTime), uuidInfo.UnixTime)
		}

		assert.Equal(t, uint32(i), uuidInfo.SequentialId)
	}

	time.Sleep(1 * time.Second)

	for i := 0; i < 1000; i++ {
		generationStartUnixTime := time.Now().Unix()
		uuid := generator.GenerateUuid(UUID_TYPE_USER)
		generationEndUnixTime := time.Now().Unix()

		uuidInfo := generator.parseInternalUuidInfo(uuid)

		if generationStartUnixTime == generationEndUnixTime {
			assert.Equal(t, uint32(generationStartUnixTime), uuidInfo.UnixTime)
		}

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
		waitGroup.Add(1)

		go func(currentRoutineIndex int) {
			for cycle := 0; cycle < 1000; cycle++ {
				generationStartUnixTime := time.Now().Unix()
				uuid := generator.GenerateUuid(UUID_TYPE_USER)
				generationEndUnixTime := time.Now().Unix()

				uuidInfo := generator.parseInternalUuidInfo(uuid)

				if generationStartUnixTime == generationEndUnixTime {
					if uint32(generationStartUnixTime) != uuidInfo.UnixTime {
						mutex.Lock()
						assert.Equal(t, uint32(generationStartUnixTime), uuidInfo.UnixTime)
						mutex.Unlock()
					}
				}

				if uuidInfo.SequentialId == 0 {
					fmt.Printf("routine#%d generate uuid %d\n", currentRoutineIndex, uuid)

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

func TestGenerateUuid_1000000TimesConcurrent(t *testing.T) {
	concurrentCount := 10
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 3})
	var mutex sync.Mutex
	var generatedIds sync.Map
	var waitGroup sync.WaitGroup

	for routineIndex := 0; routineIndex < concurrentCount; routineIndex++ {
		waitGroup.Add(1)

		go func(currentRoutineIndex int) {
			for cycle := 0; cycle < 100000; cycle++ {
				if cycle%50000 == 0 { // each server can only generate 500,000 (10 * 50000) uuids in one second
					time.Sleep(1000 * time.Millisecond)
				}

				generationStartUnixTime := time.Now().Unix()
				uuid := generator.GenerateUuid(UUID_TYPE_USER)
				generationEndUnixTime := time.Now().Unix()

				uuidInfo := generator.parseInternalUuidInfo(uuid)

				if generationStartUnixTime == generationEndUnixTime {
					if uint32(generationStartUnixTime) != uuidInfo.UnixTime {
						mutex.Lock()
						assert.Equal(t, uint32(generationStartUnixTime), uuidInfo.UnixTime)
						mutex.Unlock()
					}
				}

				if uuidInfo.SequentialId == 0 {
					fmt.Printf("routine#%d generate uuid %d\n", currentRoutineIndex, uuid)

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
	expectedUuidServerId := uint8(90)
	expectedUuidType := UUID_TYPE_TRANSACTION
	expectedUuidCount := uint8(255)

	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: expectedUuidServerId})

	generationStartUnixTime := time.Now().Unix()
	uuids := generator.GenerateUuids(expectedUuidType, expectedUuidCount)
	generationEndUnixTime := time.Now().Unix()

	for i := 0; i < int(expectedUuidCount); i++ {
		uuidInfo := generator.parseInternalUuidInfo(uuids[i])

		if generationStartUnixTime == generationEndUnixTime {
			actualUnixTime := uuidInfo.UnixTime
			assert.Equal(t, uint32(generationStartUnixTime), actualUnixTime)
		}

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
		generationStartUnixTime := time.Now().Unix()
		uuids := generator.GenerateUuids(expectedUuidType, expectedUuidCount)
		generationEndUnixTime := time.Now().Unix()
		var firstSeqId uint32

		for i := 0; i < int(expectedUuidCount); i++ {
			uuidInfo := generator.parseInternalUuidInfo(uuids[i])

			if generationStartUnixTime == generationEndUnixTime {
				actualUnixTime := uuidInfo.UnixTime
				assert.Equal(t, uint32(generationStartUnixTime), actualUnixTime)
			}

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
		waitGroup.Add(1)

		go func(currentRoutineIndex int) {
			for cycle := 0; cycle < 100; cycle++ {
				generationStartUnixTime := time.Now().Unix()
				uuids := generator.GenerateUuids(UUID_TYPE_USER, expectedUuidCount)
				generationEndUnixTime := time.Now().Unix()

				for i := 0; i < int(expectedUuidCount); i++ {
					uuidInfo := generator.parseInternalUuidInfo(uuids[i])

					if generationStartUnixTime == generationEndUnixTime {
						if uint32(generationStartUnixTime) != uuidInfo.UnixTime {
							mutex.Lock()
							assert.Equal(t, uint32(generationStartUnixTime), uuidInfo.UnixTime)
							mutex.Unlock()
						}
					}

					if uuidInfo.SequentialId == 0 {
						fmt.Printf("routine#%d generate uuid %d\n", currentRoutineIndex, uuids[i])

						if existedRoutineIndex, exists := generatedIds.Load(uuids[i]); exists {
							mutex.Lock()
							assert.Fail(t, fmt.Sprintf("uuid \"%d\" conflicts, unix time is %d, seq id is %d, existed routine index is %d, current routine index is %d", uuids[i], uuidInfo.UnixTime, uuidInfo.SequentialId, existedRoutineIndex, currentRoutineIndex))
							mutex.Unlock()
						}

						generatedIds.Store(uuids[i], currentRoutineIndex)
					}
				}
			}

			waitGroup.Done()
		}(routineIndex)
	}

	waitGroup.Wait()
}

func TestGenerateUuids_1000000TimesConcurrent(t *testing.T) {
	concurrentCount := 10
	expectedUuidCount := uint8(250)
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 3})
	var mutex sync.Mutex
	var generatedIds sync.Map
	var waitGroup sync.WaitGroup

	for routineIndex := 0; routineIndex < concurrentCount; routineIndex++ {
		waitGroup.Add(1)

		go func(currentRoutineIndex int) {
			for cycle := 0; cycle < 400; cycle++ {
				if cycle%200 == 0 { // each server can only generate 500,000 (10 * 250 * 200) uuids in one second
					time.Sleep(1000 * time.Millisecond)
				}

				generationStartUnixTime := time.Now().Unix()
				uuids := generator.GenerateUuids(UUID_TYPE_USER, expectedUuidCount)
				generationEndUnixTime := time.Now().Unix()

				for i := 0; i < int(expectedUuidCount); i++ {
					uuidInfo := generator.parseInternalUuidInfo(uuids[i])

					if generationStartUnixTime == generationEndUnixTime {
						if uint32(generationStartUnixTime) != uuidInfo.UnixTime {
							mutex.Lock()
							assert.Equal(t, uint32(generationStartUnixTime), uuidInfo.UnixTime)
							mutex.Unlock()
						}
					}

					if uuidInfo.SequentialId == 0 {
						fmt.Printf("routine#%d generate uuid %d\n", currentRoutineIndex, uuids[i])

						if existedRoutineIndex, exists := generatedIds.Load(uuids[i]); exists {
							mutex.Lock()
							assert.Fail(t, fmt.Sprintf("uuid \"%d\" conflicts, unix time is %d, seq id is %d, existed routine index is %d, current routine index is %d", uuids[i], uuidInfo.UnixTime, uuidInfo.SequentialId, existedRoutineIndex, currentRoutineIndex))
							mutex.Unlock()
						}

						generatedIds.Store(uuids[i], currentRoutineIndex)
					}
				}
			}

			waitGroup.Done()
		}(routineIndex)
	}

	waitGroup.Wait()
}

func TestGenerateUuid_Over524287Times(t *testing.T) {
	generator, _ := NewInternalUuidGenerator(&settings.Config{UuidServerId: 1})
	onceGenerateCount := uint8(255)
	generationStartUnixTime := time.Now().Unix()

	for i := 0; i < 2057; i++ { // 2056*255=524280, 2057*255=524,535 (only can generates 524,287 uuids per second)
		uuids := generator.GenerateUuids(UUID_TYPE_USER, onceGenerateCount)

		if i < 2056 {
			if len(uuids) < int(onceGenerateCount) {
				assert.Fail(t, fmt.Sprintf("%d uuids should be generated", onceGenerateCount))
			}
		} else {
			generationEndUnixTime := time.Now().Unix()

			if generationStartUnixTime == generationEndUnixTime && len(uuids) > 0 {
				assert.Fail(t, fmt.Sprintf("uuids should not be generated because there are too many uuids in one second"))
			}
		}
	}
}
