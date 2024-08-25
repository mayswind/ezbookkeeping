package cron

import (
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func TestCronJobNextRunTimeWithIntervalPeriod(t *testing.T) {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	assert.Nil(t, err)

	job := CronJob{
		Name:        "TestCronJobWithIntervalPeriod",
		Description: "The test cron job",
		Period: CronJobIntervalPeriod{
			Interval: 2*time.Hour + 34*time.Minute + 56*time.Second,
		},
		Run: func(c *core.CronContext) error {
			return nil
		},
	}

	gocronJob, err := scheduler.NewJob(
		job.Period.ToJobDefinition(),
		gocron.NewTask(job.doRun),
		gocron.WithName(job.Name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	assert.Nil(t, err)

	scheduler.Start()

	currentTime := time.Now()
	nextRunTime, err := gocronJob.NextRun()
	assert.Nil(t, err)

	expectedNextTime := currentTime.Add(2 * time.Hour).Add(34 * time.Minute).Add(56 * time.Second)

	assert.Equal(t, expectedNextTime.Year(), nextRunTime.Year())
	assert.Equal(t, expectedNextTime.Month(), nextRunTime.Month())
	assert.Equal(t, expectedNextTime.Day(), nextRunTime.Day())
	assert.Equal(t, expectedNextTime.Hour(), nextRunTime.Hour())
	assert.Equal(t, expectedNextTime.Minute(), nextRunTime.Minute())
	assert.Equal(t, expectedNextTime.Second(), nextRunTime.Second())

	err = scheduler.Shutdown()
	assert.Nil(t, err)
}

func TestCronJobNextRunTimeWithFixedHourPeriod(t *testing.T) {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	assert.Nil(t, err)

	job := CronJob{
		Name:        "TestCronJobWithFixedHourPeriod",
		Description: "The test cron job",
		Period: CronJobFixedHourPeriod{
			Hour: 0,
		},
		Run: func(c *core.CronContext) error {
			return nil
		},
	}

	gocronJob, err := scheduler.NewJob(
		job.Period.ToJobDefinition(),
		gocron.NewTask(job.doRun),
		gocron.WithName(job.Name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	assert.Nil(t, err)

	scheduler.Start()

	nextRunTime, err := gocronJob.NextRun()
	assert.Nil(t, err)

	tomorrow := time.Now().AddDate(0, 0, 1)

	assert.Equal(t, tomorrow.Year(), nextRunTime.Year())
	assert.Equal(t, tomorrow.Month(), nextRunTime.Month())
	assert.Equal(t, tomorrow.Day(), nextRunTime.Day())
	assert.Equal(t, 0, nextRunTime.Hour())
	assert.Equal(t, 0, nextRunTime.Minute())
	assert.Equal(t, 0, nextRunTime.Second())

	err = scheduler.Shutdown()
	assert.Nil(t, err)
}

func TestCronJobNextRunTimeWithEvery15MinutesPeriod(t *testing.T) {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	assert.Nil(t, err)

	expectedSecond := uint32(23)

	job := CronJob{
		Name:        "TestCronJobWithEvery15MinutesPeriod",
		Description: "The test cron job",
		Period: CronJobEvery15MinutesPeriod{
			Second: expectedSecond,
		},
		Run: func(c *core.CronContext) error {
			return nil
		},
	}

	assert.Equal(t, 15*time.Minute, job.Period.GetInterval())

	gocronJob, err := scheduler.NewJob(
		job.Period.ToJobDefinition(),
		gocron.NewTask(job.doRun),
		gocron.WithName(job.Name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	assert.Nil(t, err)

	scheduler.Start()

	nextRunTime, err := gocronJob.NextRun()
	assert.Nil(t, err)

	nextMinuteTime := time.Now()

	if (nextMinuteTime.Minute() == 0 || nextMinuteTime.Minute() == 15 || nextMinuteTime.Minute() == 30 || nextMinuteTime.Minute() == 45) && nextMinuteTime.Second() < int(expectedSecond) {
		// Do Nothing
	} else {
		nextMinute := ((nextMinuteTime.Minute() / 15) + 1) * 15
		minuteDiff := nextMinute - nextMinuteTime.Minute()
		nextMinuteTime = nextMinuteTime.Add(time.Duration(int64(minuteDiff) * int64(time.Minute)))
	}

	assert.Equal(t, nextMinuteTime.Year(), nextRunTime.Year())
	assert.Equal(t, nextMinuteTime.Month(), nextRunTime.Month())
	assert.Equal(t, nextMinuteTime.Day(), nextRunTime.Day())
	assert.Equal(t, nextMinuteTime.Hour(), nextRunTime.Hour())
	assert.Equal(t, nextMinuteTime.Minute(), nextRunTime.Minute())
	assert.Equal(t, int(expectedSecond), nextRunTime.Second())

	err = scheduler.Shutdown()
	assert.Nil(t, err)
}

func TestCronJobNextRunTimeWithFixedTimePeriod(t *testing.T) {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	assert.Nil(t, err)

	expectedTime := time.Now().Add(123456789 * time.Second)

	job := CronJob{
		Name:        "TestCronJobWithFixedTimePeriod",
		Description: "The test cron job",
		Period: CronJobFixedTimePeriod{
			Time: expectedTime,
		},
		Run: func(c *core.CronContext) error {
			return nil
		},
	}

	gocronJob, err := scheduler.NewJob(
		job.Period.ToJobDefinition(),
		gocron.NewTask(job.doRun),
		gocron.WithName(job.Name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	assert.Nil(t, err)

	scheduler.Start()

	nextRunTime, err := gocronJob.NextRun()
	assert.Nil(t, err)

	assert.Equal(t, expectedTime.Year(), nextRunTime.Year())
	assert.Equal(t, expectedTime.Month(), nextRunTime.Month())
	assert.Equal(t, expectedTime.Day(), nextRunTime.Day())
	assert.Equal(t, expectedTime.Hour(), nextRunTime.Hour())
	assert.Equal(t, expectedTime.Minute(), nextRunTime.Minute())
	assert.Equal(t, expectedTime.Second(), nextRunTime.Second())

	err = scheduler.Shutdown()
	assert.Nil(t, err)
}
