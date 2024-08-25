package core

import (
	"context"
	"strconv"
	"strings"
	"time"
)

// CronContext represents the cron job context
type CronContext struct {
	context.Context
	contextId       string
	cronJobInterval time.Duration
}

// GetContextId returns the current context id
func (c *CronContext) GetContextId() string {
	return c.contextId
}

// GetInterval returns the current cron job interval
func (c *CronContext) GetInterval() time.Duration {
	return c.cronJobInterval
}

// NewCronJobContext returns a new cron job context
func NewCronJobContext(cronJobName string, cronJobInterval time.Duration) *CronContext {
	return &CronContext{
		Context:         context.Background(),
		contextId:       generateNewRandomCronContextId(cronJobName),
		cronJobInterval: cronJobInterval,
	}
}

func generateNewRandomCronContextId(cronJobName string) string {
	var ret strings.Builder
	ret.WriteString("cron-job-")
	ret.WriteString(strings.ToLower(cronJobName))
	ret.WriteRune('-')
	ret.WriteString(strconv.FormatInt(time.Now().Unix(), 10))

	return ret.String()
}
