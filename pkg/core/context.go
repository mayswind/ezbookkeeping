package core

import "context"

// Context is the base context of ezBookkeeping
type Context interface {
	context.Context
	GetContextId() string
	GetClientLocale() string
}
