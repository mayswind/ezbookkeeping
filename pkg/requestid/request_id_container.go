package requestid

import (
	"github.com/mayswind/lab/pkg/settings"
)

type RequestIdContainer struct {
	Current RequestIdGenerator
}

var (
	Container = &RequestIdContainer{}
)

func InitializeRequestIdGenerator(config *settings.Config) error {
	generator, err := NewDefaultRequestIdGenerator(config)

	if err != nil {
		return err
	}

	Container.Current = generator
	return nil
}

func (u *RequestIdContainer) GenerateRequestId(clientIpAddr string) string {
	return u.Current.GenerateRequestId(clientIpAddr)
}
