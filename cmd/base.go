package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func bindAction(fn core.CliHandlerFunc) cli.ActionFunc {
	return func(cliCtx *cli.Context) error {
		c := core.WrapCilContext(cliCtx)
		return fn(c)
	}
}
