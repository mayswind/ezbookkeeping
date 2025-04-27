package cmd

import (
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cron"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// CronJobs represents the cron command
var CronJobs = &cli.Command{
	Name:  "cron",
	Usage: "ezBookkeeping cron job utilities",
	Commands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "List all enabled cron jobs",
			Action: bindAction(listAllCronJobs),
		},
		{
			Name:   "run",
			Usage:  "Run specified cron job",
			Action: bindAction(runCronJob),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Cron job name",
				},
			},
		},
	},
}

func listAllCronJobs(c *core.CliContext) error {
	config, err := initializeSystem(c)

	if err != nil {
		return err
	}

	err = cron.InitializeCronJobSchedulerContainer(c, config, false)

	if err != nil {
		log.CliErrorf(c, "[cron_jobs.listAllCronJobs] initializes cron job scheduler failed, because %s", err.Error())
		return err
	}

	cronJobs := cron.Container.GetAllJobs()

	if len(cronJobs) < 1 {
		log.CliErrorf(c, "[cron_jobs.listAllCronJobs] there are no enabled cron jobs")
		return err
	}

	for i := 0; i < len(cronJobs); i++ {
		if i > 0 {
			fmt.Printf("---\n")
		}

		cronJob := cronJobs[i]

		fmt.Printf("[Name] %s\n", cronJob.Name)
		fmt.Printf("[Description] %s\n", cronJob.Description)
		fmt.Printf("[Interval] Every %s\n", cronJob.Period.GetInterval())
	}

	return nil
}

func runCronJob(c *core.CliContext) error {
	config, err := initializeSystem(c)

	if err != nil {
		return err
	}

	err = cron.InitializeCronJobSchedulerContainer(c, config, false)

	if err != nil {
		log.CliErrorf(c, "[cron_jobs.runCronJob] initializes cron job scheduler failed, because %s", err.Error())
		return err
	}

	jobName := c.String("name")
	err = cron.Container.SyncRunJobNow(jobName)

	if err != nil {
		log.CliErrorf(c, "[cron_jobs.runCronJob] failed to run cron job \"%s\", because %s", jobName, err.Error())
		return err
	}

	log.CliInfof(c, "[cron_jobs.runCronJob] run cron job \"%s\" successfully", jobName)

	return nil
}
