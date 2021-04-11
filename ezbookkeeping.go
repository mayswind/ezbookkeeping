package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/ezbookkeeping/cmd"
	"github.com/mayswind/ezbookkeeping/pkg/version"
)

func main() {
	app := &cli.App{
		Name:    "ezBookkeeping",
		Usage:   "A lightweight personal bookkeeping app hosted by yourself.",
		Version: version.GetFullVersion(),
		Commands: []*cli.Command{
			cmd.WebServer,
			cmd.Database,
			cmd.UserData,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "conf-path",
				Usage: "Custom config `FILE` path",
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("Failed to run ezBookkeeping with %s: %v", os.Args, err)
	}
}
