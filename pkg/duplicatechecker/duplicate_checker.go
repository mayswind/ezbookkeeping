package duplicatechecker

import "time"

// DuplicateChecker is common duplicate checker interface
type DuplicateChecker interface {
	GetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string) (bool, string)
	SetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string, remark string)
	RemoveSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string)
	GetOrSetCronJobRunningInfo(jobName string, runningInfo string, runningInterval time.Duration) (bool, string)
	RemoveCronJobRunningInfo(jobName string)
	GetFailureCount(failureKey string) uint32
	IncreaseFailureCount(failureKey string) uint32
}
