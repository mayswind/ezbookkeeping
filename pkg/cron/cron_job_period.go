package cron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

// CronJobPeriod represents the cron job period
type CronJobPeriod interface {
	GetInterval() time.Duration
	ToJobDefinition() gocron.JobDefinition
}

// CronJobIntervalPeriod represents the period of execution at intervals
type CronJobIntervalPeriod struct {
	Interval time.Duration
}

// CronJobFixedHourPeriod represents the period of execution at fixed hour
type CronJobFixedHourPeriod struct {
	Hour uint32
}

// CronJobFixedTimePeriod represents the period of execution at fixed time
type CronJobFixedTimePeriod struct {
	Time time.Time
}

// GetInterval returns the interval time of the period of CronJobIntervalPeriod
func (p CronJobIntervalPeriod) GetInterval() time.Duration {
	return p.Interval
}

// ToJobDefinition returns the gocron job definition of the period of CronJobIntervalPeriod
func (p CronJobIntervalPeriod) ToJobDefinition() gocron.JobDefinition {
	return gocron.DurationJob(p.Interval)
}

// GetInterval returns the interval time of the period of CronJobFixedHourPeriod
func (p CronJobFixedHourPeriod) GetInterval() time.Duration {
	return 24 * time.Hour
}

// ToJobDefinition returns the gocron job definition of the period of CronJobFixedHourPeriod
func (p CronJobFixedHourPeriod) ToJobDefinition() gocron.JobDefinition {
	return gocron.DailyJob(
		1,
		gocron.NewAtTimes(
			gocron.NewAtTime(uint(p.Hour), 0, 0),
		),
	)
}

// GetInterval returns the interval time of the period of CronJobFixedTimePeriod
func (p CronJobFixedTimePeriod) GetInterval() time.Duration {
	return 0
}

// ToJobDefinition returns the gocron job definition of the period of CronJobFixedTimePeriod
func (p CronJobFixedTimePeriod) ToJobDefinition() gocron.JobDefinition {
	return gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(p.Time))
}
