package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/ezbookkeeping/cmd"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var (
	// Version holds the version of this execution program
	Version string

	// CommitHash holds the git commit hash of this execution program's source code
	CommitHash string

	// BuildUnixTime holds the time when starting building this execution program
	BuildUnixTime string
)

func main() {
	settings.Version = Version
	settings.CommitHash = CommitHash

	app := &cli.App{
		Name:    "ezBookkeeping",
		Usage:   "A lightweight personal bookkeeping app hosted by yourself.",
		Version: GetFullVersion(),
		Commands: []*cli.Command{
			cmd.WebServer,
			cmd.Database,
			cmd.UserData,
			cmd.CronJobs,
			cmd.SecurityUtils,
			cmd.Utilities,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf-path",
				Usage: "Custom config `FILE` path",
			},
			&cli.BoolFlag{
				Name:  "no-boot-log",
				Usage: "Disable boot log",
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("Failed to run ezBookkeeping with %s: %v", os.Args, err)
	}
}

// GetFullVersion returns the full version
func GetFullVersion() string {
	fullVersion := "Local Build"

	if Version != "" {
		fullVersion = Version
	}

	additionalInfos := make([]string, 0, 2)

	if CommitHash != "" {
		additionalInfos = append(additionalInfos, "commit "+CommitHash)
	}

	if BuildUnixTime != "" {
		unixTime, err := utils.StringToInt64(BuildUnixTime)

		if unixTime > 0 && err == nil {
			additionalInfos = append(additionalInfos, "build time "+utils.FormatUnixTimeToLongDateTimeInServerTimezone(unixTime))
		}
	}

	if len(additionalInfos) > 0 {
		fullVersion = fmt.Sprintf("%s (%s)", fullVersion, strings.Join(additionalInfos, ", "))
	}

	return fullVersion
}
