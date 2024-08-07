package datastore

import (
	"fmt"
	"time"

	"xorm.io/xorm/log"

	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// XOrmContextAdapter represents the context adapter for xorm
type XOrmContextAdapter struct {
	requestId string
}

// Deadline does nothing
func (c *XOrmContextAdapter) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil
func (c *XOrmContextAdapter) Done() <-chan struct{} {
	return nil
}

// Err always returns nil
func (c *XOrmContextAdapter) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.
func (c *XOrmContextAdapter) Value(key any) any {
	if key == log.SessionIDKey && c.requestId != "" {
		return fmt.Sprintf("%s", c.requestId)
	}

	return nil
}

func NewXOrmContextAdapter(c *core.Context) *XOrmContextAdapter {
	if c != nil {
		return &XOrmContextAdapter{
			requestId: c.GetRequestId(),
		}
	}

	return &XOrmContextAdapter{}
}
