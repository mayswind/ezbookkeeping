package duplicatechecker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestSetAndGetSubmissionRemark(t *testing.T) {
	checker, _ := NewInMemoryDuplicateChecker(&settings.Config{
		DuplicateSubmissionsIntervalDuration:            100 * time.Second,
		InMemoryDuplicateCheckerCleanupIntervalDuration: 100 * time.Second,
	})

	uid := int64(1234567890)
	id := "2345678901"
	expectedRemark := "0123456789"

	checker.SetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, id, expectedRemark)
	found, actualRemark := checker.GetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, id)
	assert.Equal(t, true, found)
	assert.Equal(t, expectedRemark, actualRemark)
}

func TestSetAndGetNotExistedSubmissionRemark(t *testing.T) {
	checker, _ := NewInMemoryDuplicateChecker(&settings.Config{
		DuplicateSubmissionsIntervalDuration:            100 * time.Second,
		InMemoryDuplicateCheckerCleanupIntervalDuration: 100 * time.Second,
	})

	uid := int64(1234567890)
	id := "2345678901"
	expectedRemark := "0123456789"

	checker.SetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, id, expectedRemark)
	found, actualRemark := checker.GetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT, uid, id)
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRemark)

	found, actualRemark = checker.GetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, "")
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRemark)

	found, actualRemark = checker.GetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, "2345678900")
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRemark)

	found, actualRemark = checker.GetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, int64(123456791), "")
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRemark)
}

func TestSetAndGetExpiredSubmissionRemark(t *testing.T) {
	checker, _ := NewInMemoryDuplicateChecker(&settings.Config{
		DuplicateSubmissionsIntervalDuration:            time.Second,
		InMemoryDuplicateCheckerCleanupIntervalDuration: time.Second,
	})

	uid := int64(1234567890)
	id := "2345678901"
	expectedRemark := "0123456789"

	checker.SetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, id, expectedRemark)
	time.Sleep(time.Second * 2)

	found, actualRemark := checker.GetSubmissionRemark(DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, id)
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRemark)
}

func TestGetOrSetCronJobRunningInfo(t *testing.T) {
	checker, _ := NewInMemoryDuplicateChecker(&settings.Config{
		DuplicateSubmissionsIntervalDuration:            time.Second,
		InMemoryDuplicateCheckerCleanupIntervalDuration: time.Second,
	})

	jobName := "foo"
	expectedRunningInfo := "bar"

	found, actualRunningInfo := checker.GetOrSetCronJobRunningInfo(jobName, expectedRunningInfo, 100*time.Second)
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRunningInfo)

	found, actualRunningInfo = checker.GetOrSetCronJobRunningInfo(jobName, expectedRunningInfo, 100*time.Second)
	assert.Equal(t, true, found)
	assert.Equal(t, expectedRunningInfo, actualRunningInfo)

	checker.RemoveCronJobRunningInfo(jobName)

	found, actualRunningInfo = checker.GetOrSetCronJobRunningInfo(jobName, expectedRunningInfo, 100*time.Second)
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRunningInfo)
}

func TestGetNotExistedRunningInfo(t *testing.T) {
	checker, _ := NewInMemoryDuplicateChecker(&settings.Config{
		DuplicateSubmissionsIntervalDuration:            time.Second,
		InMemoryDuplicateCheckerCleanupIntervalDuration: time.Second,
	})

	jobName := "foo"
	expectedRunningInfo := "bar"

	found, actualRunningInfo := checker.GetOrSetCronJobRunningInfo(jobName, expectedRunningInfo, 100*time.Second)
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRunningInfo)

	found, actualRunningInfo = checker.GetOrSetCronJobRunningInfo("bar", expectedRunningInfo, 100*time.Second)
	assert.Equal(t, false, found)
	assert.Equal(t, "", actualRunningInfo)
}

func TestGetOrSetRunningInfoConcurrent(t *testing.T) {
	checker, _ := NewInMemoryDuplicateChecker(&settings.Config{
		DuplicateSubmissionsIntervalDuration:            time.Second,
		InMemoryDuplicateCheckerCleanupIntervalDuration: time.Second,
	})

	jobName := "foo"

	concurrentCount := 10
	var setRunningInfoCount atomic.Uint32
	var waitGroup sync.WaitGroup

	for routineIndex := 0; routineIndex < concurrentCount; routineIndex++ {
		waitGroup.Add(1)

		go func(currentRoutineIndex int) {
			randomNumber, _ := utils.GetRandomInteger(10)
			time.Sleep(time.Duration(int64(randomNumber) * int64(time.Millisecond)))

			for cycle := 0; cycle < 100; cycle++ {
				expectedRunningInfo := fmt.Sprintf("%d-%d", currentRoutineIndex, cycle)
				found, _ := checker.GetOrSetCronJobRunningInfo(jobName, expectedRunningInfo, 100*time.Second)

				if found {
					setRunningInfoCount.Add(1)
				} else {
					fmt.Printf("routine#%d set cron job running info %s\n", currentRoutineIndex, expectedRunningInfo)
				}
			}

			waitGroup.Done()
		}(routineIndex)
	}

	waitGroup.Wait()

	assert.Equal(t, uint32(999), setRunningInfoCount.Load())
}
