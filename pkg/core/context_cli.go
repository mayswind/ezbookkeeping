package core

import (
	"context"

	"github.com/urfave/cli/v3"
)

// CliContext represents the command-line context
type CliContext struct {
	context.Context
	command *cli.Command
}

// GetContextId returns the current context id
func (c *CliContext) GetContextId() string {
	return ""
}

// GetClientLocale returns the client locale name
func (c *CliContext) GetClientLocale() string {
	return ""
}

// Bool returns the boolean value of parameter
func (c *CliContext) Bool(name string) bool {
	return c.command.Bool(name)
}

// Int returns the integer value of parameter
func (c *CliContext) Int(name string) int {
	return c.command.Int(name)
}

// String returns the string value of parameter
func (c *CliContext) String(name string) string {
	return c.command.String(name)
}

// WrapCliContext returns a context wrapped by this file
func WrapCilContext(ctx context.Context, cmd *cli.Command) *CliContext {
	return &CliContext{
		Context: ctx,
		command: cmd,
	}
}
