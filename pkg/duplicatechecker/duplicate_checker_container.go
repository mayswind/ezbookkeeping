package duplicatechecker

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// DuplicateCheckerContainer contains the current duplicate checker
type DuplicateCheckerContainer struct {
	Current DuplicateChecker
}

// Initialize a duplicate checker container singleton instance
var (
	Container = &DuplicateCheckerContainer{}
)

// InitializeDuplicateChecker initializes the current duplicate checker according to the config
func InitializeDuplicateChecker(config *settings.Config) error {
	if config.DuplicateCheckerType == settings.InMemoryDuplicateCheckerType {
		checker, err := NewInMemoryDuplicateChecker(config)
		Container.Current = checker

		return err
	}

	return errs.ErrInvalidDuplicateCheckerType
}

// GetSubmissionRemark returns whether the same submission has been processed and related remark by the current duplicate checker
func (c *DuplicateCheckerContainer) GetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string) (bool, string) {
	return c.Current.GetSubmissionRemark(checkerType, uid, identification)
}

// SetSubmissionRemark saves the identification and remark by the current duplicate checker
func (c *DuplicateCheckerContainer) SetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string, remark string) {
	c.Current.SetSubmissionRemark(checkerType, uid, identification, remark)
}

// RemoveSubmissionRemark removes the identification and remark by the current duplicate checker
func (c *DuplicateCheckerContainer) RemoveSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string) {
	c.Current.RemoveSubmissionRemark(checkerType, uid, identification)
}

// GetOrSetCronJobRunningInfo returns the running info when the cron job is running or saves the running info by the current duplicate checker
func (c *DuplicateCheckerContainer) GetOrSetCronJobRunningInfo(jobName string, runningInfo string, runningInterval time.Duration) (bool, string) {
	return c.Current.GetOrSetCronJobRunningInfo(jobName, runningInfo, runningInterval)
}

// RemoveCronJobRunningInfo removes the running info of the cron job by the current duplicate checker
func (c *DuplicateCheckerContainer) RemoveCronJobRunningInfo(jobName string) {
	c.Current.RemoveCronJobRunningInfo(jobName)
}

// GetFailureCount returns the failure count of the specified failure key
func (c *DuplicateCheckerContainer) GetFailureCount(failureKey string) uint32 {
	return c.Current.GetFailureCount(failureKey)
}

// IncreaseFailureCount increases the failure count of the specified failure key
func (c *DuplicateCheckerContainer) IncreaseFailureCount(failureKey string) uint32 {
	return c.Current.IncreaseFailureCount(failureKey)
}
