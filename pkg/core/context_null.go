package core

import "context"

const nullContextId = "00000000-0000-0000-0000-00000000"

// NullContext represents the null context
type NullContext struct {
	context.Context
}

// GetContextId returns the current context id
func (c *NullContext) GetContextId() string {
	return nullContextId
}

// NewCronJobContext returns a new null context
func NewNullContext() *NullContext {
	return &NullContext{
		Context: context.Background(),
	}
}
