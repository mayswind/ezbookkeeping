package duplicatechecker

import (
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

// SetSubmissionRemark saves the identification and remark to in-memory cache by the current duplicate checker
func (c *DuplicateCheckerContainer) SetSubmissionRemark(checkerType DuplicateCheckerType, uid int64, identification string, remark string) {
	c.Current.SetSubmissionRemark(checkerType, uid, identification, remark)
}
