package requestid

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// RequestIdContainer contains the current request id generator
type RequestIdContainer struct {
	Current RequestIdGenerator
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

	Container.Current = generator
	return nil
}

// GenerateRequestId returns a new request id by the current request id generator
func (u *RequestIdContainer) GenerateRequestId(clientIpAddr string, clientPort uint16) string {
	return u.Current.GenerateRequestId(clientIpAddr, clientPort)
}
