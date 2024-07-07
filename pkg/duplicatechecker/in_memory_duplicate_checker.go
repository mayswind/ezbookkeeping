package duplicatechecker

import (
	"fmt"

	"github.com/patrickmn/go-cache"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// InMemoryDuplicateChecker represents in-memory duplicate checker
type InMemoryDuplicateChecker struct {
	cache *cache.Cache
}

// NewInMemoryDuplicateChecker returns a new in-memory duplicate checker
func NewInMemoryDuplicateChecker(config *settings.Config) (*InMemoryDuplicateChecker, error) {
	checker := &InMemoryDuplicateChecker{
		cache: cache.New(config.DuplicateSubmissionsIntervalDuration, config.InMemoryDuplicateCheckerCleanupIntervalDuration),
	}

	return checker, nil
}

// Get returns whether the same submission has been processed and related remark
func (c *InMemoryDuplicateChecker) Get(checkerType DuplicateCheckerType, uid int64, identification string) (bool, string) {
	existedRemark, found := c.cache.Get(c.getCacheKey(checkerType, uid, identification))

	if found {
		return true, existedRemark.(string)
	}

	return false, ""
}

// Set saves the identification and remark to in-memory cache
func (c *InMemoryDuplicateChecker) Set(checkerType DuplicateCheckerType, uid int64, identification string, remark string) {
	c.cache.Set(c.getCacheKey(checkerType, uid, identification), remark, cache.DefaultExpiration)
}

func (c *InMemoryDuplicateChecker) getCacheKey(checkerType DuplicateCheckerType, uid int64, identification string) string {
	return fmt.Sprintf("%d|%d|%s", checkerType, uid, identification)
}
