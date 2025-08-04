package requestid

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// RequestIdContainer contains the current request id generator
type RequestIdContainer struct {
	current RequestIdGenerator
}

// Initialize a request id container singleton instance
var (
	Container = &RequestIdContainer{}
)

// InitializeRequestIdGenerator initializes the current request id generator according to the config
func InitializeRequestIdGenerator(c core.Context, config *settings.Config) error {
	generator, err := NewDefaultRequestIdGenerator(c, config)

	if err != nil {
		return err
	}

	Container.current = generator
	return nil
}

// GenerateRequestId returns a new request id by the current request id generator
func (r *RequestIdContainer) GenerateRequestId(clientIpAddr string, clientPort uint16) string {
	if r.current == nil {
		return ""
	}

	return r.current.GenerateRequestId(clientIpAddr, clientPort)
}

// GetCurrentServerUniqId returns current server unique id
func (r *RequestIdContainer) GetCurrentServerUniqId() uint16 {
	if r.current == nil {
		return 0
	}

	return r.current.GetCurrentServerUniqId()
}

// GetCurrentInstanceUniqId returns current application instance unique id
func (r *RequestIdContainer) GetCurrentInstanceUniqId() uint16 {
	if r.current == nil {
		return 0
	}

	return r.current.GetCurrentInstanceUniqId()
}
