package cron

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/services"
)

var RemoveExpiredTokensJob = &CronJob{
	Name:        "RemoveExpiredTokens",
	Description: "Periodically remove expired user tokens from the database.",
	Interval:    24 * time.Hour,
	Run: func() error {
		return services.Tokens.DeleteAllExpiredTokens(nil)
	},
}
