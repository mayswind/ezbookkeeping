package cron

import (
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// CronJobSchedulerContainer contains the current cron job scheduler
type CronJobSchedulerContainer struct {
	scheduler        gocron.Scheduler
	allJobs          []*CronJob
	allJobsMap       map[string]*CronJob
	allGocronJobsMap map[string]gocron.Job
}

// Initialize a cron job scheduler container singleton instance
var (
	Container = &CronJobSchedulerContainer{
		allJobsMap:       make(map[string]*CronJob),
		allGocronJobsMap: make(map[string]gocron.Job),
	}
)

// InitializeCronJobSchedulerContainer initializes the cron job scheduler according to the config
func InitializeCronJobSchedulerContainer(config *settings.Config, startScheduler bool) error {
	var err error

	Container.scheduler, err = gocron.NewScheduler(
		gocron.WithLocation(time.Local),
		gocron.WithLogger(NewGocronLoggerAdapter()),
	)

	if err != nil {
		return err
	}

	Container.registerAllJobs(config)

	if startScheduler {
		Container.scheduler.Start()
	}

	return nil
}

// GetAllJobs returns all the cron jobs
func (c *CronJobSchedulerContainer) GetAllJobs() []*CronJob {
	return c.allJobs
}

// SyncRunJobNow runs the specified cron job synchronously now
func (c *CronJobSchedulerContainer) SyncRunJobNow(jobName string) error {
	if jobName == "" {
		return errs.ErrCronJobNameIsEmpty
	}

	job := c.allJobsMap[jobName]

	if job == nil {
		return errs.ErrCronJobNotExistsOrNotEnabled
	}

	gocronJob := c.allGocronJobsMap[jobName]

	if gocronJob == nil {
		return errs.ErrCronJobNotExistsOrNotEnabled
	}

	job.doRun()
	return nil
}

func (c *CronJobSchedulerContainer) registerAllJobs(config *settings.Config) {
	if config.EnableRemoveExpiredTokens {
		Container.registerIntervalJob(RemoveExpiredTokensJob)
	}
}

func (c *CronJobSchedulerContainer) registerIntervalJob(job *CronJob) {
	gocronJob, err := c.scheduler.NewJob(
		gocron.DurationJob(job.Interval),
		gocron.NewTask(job.doRun),
		gocron.WithName(job.Name),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	if err == nil {
		c.allJobs = append(c.allJobs, job)
		c.allJobsMap[job.Name] = job
		c.allGocronJobsMap[job.Name] = gocronJob
		log.Infof("[cron_container.registerJob] job \"%s\" has been registered", job.Name)
		log.Debugf("[cron_container.registerJob] job \"%s\" gocron id is %s", job.Name, gocronJob.ID())
	} else {
		log.Errorf("[cron_container.registerJob] job \"%s\" cannot be been registered, because %s", job.Name, err.Error())
	}
}
