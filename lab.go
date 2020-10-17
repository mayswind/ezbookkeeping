package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/mayswind/lab/cmd"
)


const LAB_VERSION = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "lab"
	app.Usage = "A lightweight account book app hosted by yourself."
	app.Version = LAB_VERSION
	app.Commands = []cli.Command{
		cmd.WebServer,
		cmd.Database,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf-path",
			Usage: "Custom config `FILE` path",
		},
	}
	app.Run(os.Args)
}
