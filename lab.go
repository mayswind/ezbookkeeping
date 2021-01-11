package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/mayswind/lab/cmd"
)

var version string
var commitHash string

func main() {
	app := &cli.App{
		Name:    "lab",
		Usage:   "A lightweight account book app hosted by yourself.",
		Version: getVersion(),
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

func getVersion() string {
	fullVersion := "Local Build"

	if version != "" {
		fullVersion = version
	}

	if commitHash != "" {
		fullVersion = fmt.Sprintf("%s (%s)", fullVersion, commitHash)
	}

	return fullVersion
}
