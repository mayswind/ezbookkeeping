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

// Get returns whether the same submission has been processed and related remark by the current duplicate checker
func (c *DuplicateCheckerContainer) Get(checkerType DuplicateCheckerType, uid int64, identification string) (bool, string) {
	return c.Current.Get(checkerType, uid, identification)
}

// Set saves the identification and remark to in-memory cache by the current duplicate checker
func (c *DuplicateCheckerContainer) Set(checkerType DuplicateCheckerType, uid int64, identification string, remark string) {
	c.Current.Set(checkerType, uid, identification, remark)
}
