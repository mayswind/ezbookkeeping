package requestid

import (
	"github.com/mayswind/lab/pkg/settings"
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
func InitializeRequestIdGenerator(config *settings.Config) error {
	generator, err := NewDefaultRequestIdGenerator(config)

	if err != nil {
		return err
	}

	Container.Current = generator
	return nil
}

// GenerateRequestId returns a new request id by the current request id generator
func (u *RequestIdContainer) GenerateRequestId(clientIpAddr string) string {
	return u.Current.GenerateRequestId(clientIpAddr)
}
