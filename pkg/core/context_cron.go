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
	contextId string
}

// GetContextId returns the current context id
func (c *CronContext) GetContextId() string {
	return c.contextId
}

// NewCronJobContext returns a new cron job context
func NewCronJobContext(cronJobName string) *CronContext {
	return &CronContext{
		Context:   context.Background(),
		contextId: generateNewRandomCronContextId(cronJobName),
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
