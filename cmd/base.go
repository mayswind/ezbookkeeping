package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

func bindAction(fn core.CliHandlerFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		c := core.WrapCilContext(ctx, cmd)
		return fn(c)
	}
}
