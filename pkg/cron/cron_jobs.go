package cron

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// RemoveExpiredTokensJob represents the cron job which periodically remove expired user tokens from the database
var RemoveExpiredTokensJob = &CronJob{
	Name:        "RemoveExpiredTokens",
	Description: "Periodically remove expired user tokens from the database.",
	Period: CronJobFixedHourPeriod{
		Hour: 0,
	},
	Run: func(c *core.CronContext) error {
		return services.Tokens.DeleteAllExpiredTokens(c)
	},
}

// CreateScheduledTransactionJob represents the cron job which periodically create transaction by scheduled transaction template
var CreateScheduledTransactionJob = &CronJob{
	Name:        "CreateScheduledTransaction",
	Description: "Periodically create transaction by scheduled transaction template.",
	Period: CronJobEvery15MinutesPeriod{
		Second: 0,
	},
	Run: func(c *core.CronContext) error {
		return services.Transactions.CreateScheduledTransactions(c, time.Now().Unix(), c.GetInterval())
	},
}

// RemoveExpiredArchivedQueuedTransactionsJob represents the cron job which periodically removes the
// rejected draft transactions whose retention period has elapsed
var RemoveExpiredArchivedQueuedTransactionsJob = &CronJob{
	Name:        "RemoveExpiredArchivedQueuedTransactions",
	Description: "Periodically remove rejected draft transactions from the database once their retention period has elapsed.",
	Period: CronJobIntervalPeriod{
		Interval: time.Hour,
	},
	Run: func(c *core.CronContext) error {
		return services.QueuedTransactions.DeleteAllExpiredArchivedQueuedTransactions(c, time.Now().Unix())
	},
}
