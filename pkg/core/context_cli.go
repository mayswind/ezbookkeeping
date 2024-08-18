package core

import (
	"github.com/urfave/cli/v2"
)

// CliContext represents the command-line context
type CliContext struct {
	*cli.Context
}

// GetContextId returns the current context id
func (c *CliContext) GetContextId() string {
	return ""
}

// WrapCliContext returns a context wrapped by this file
func WrapCilContext(cliCtx *cli.Context) *CliContext {
	return &CliContext{
		Context: cliCtx,
	}
}
