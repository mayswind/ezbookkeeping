package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/cmd"
)

const LAB_VERSION = "0.1.0"

func main() {
	app := &cli.App{
		Name:    "lab",
		Usage:   "A lightweight account book app hosted by yourself.",
		Version: LAB_VERSION,
		Commands: []*cli.Command{
			cmd.WebServer,
			cmd.Database,
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
