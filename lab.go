package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/cmd"
	"github.com/mayswind/lab/pkg/version"
)

func main() {
	app := &cli.App{
		Name:    "lab",
		Usage:   "A lightweight account book app hosted by yourself.",
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
		log.Fatalf("Failed to run lab app with %s: %v", os.Args, err)
	}
}
