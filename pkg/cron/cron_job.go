package cron

import (
	"fmt"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

type CronJob struct {
	Name        string
	Description string
	Interval    time.Duration
	Run         func() error
}

func (c *CronJob) doRun() {
	start := time.Now()
	localAddr, err := utils.GetLocalIPAddressesString()

	if err != nil {
		log.Warnf("[cron_job.doRun] job \"%s\" cannot get local ipv4 address, because %s", c.Name, err.Error())
		return
	}

	currentInfo := fmt.Sprintf("ip: %s, startTime: %d", localAddr, time.Now().Unix())
	found, runningInfo := duplicatechecker.Container.GetOrSetCronJobRunningInfo(c.Name, currentInfo, c.Interval)

	if found {
		log.Warnf("[cron_job.doRun] job \"%s\" is already running (%s)", c.Name, runningInfo)
		return
	}

	err = c.Run()

	duplicatechecker.Container.Current.RemoveCronJobRunningInfo(c.Name)

	now := time.Now()

	if err != nil {
		log.Errorf("[cron_job.doRun] failed to run job \"%s\", because %s", c.Name, err.Error())
		return
	}

	cost := now.Sub(start).Nanoseconds() / 1e6
	log.Infof("[cron_job.doRun] run job \"%s\" successfully, cost %dms", c.Name, cost)
}
