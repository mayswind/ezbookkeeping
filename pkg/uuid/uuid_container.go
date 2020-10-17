package uuid

import (
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/settings"
)

type UuidContainer struct {
	Current UuidGenerator
}

var (
	Container = &UuidContainer{}
)

func InitializeUuidGenerator(config *settings.Config) error {
	if config.UuidGeneratorType == settings.UUID_GENERATOR_TYPE_INTERNAL {
		generator, err := NewInternalUuidGenerator(config)
		Container.Current = generator

		return err
	}

	return errs.ErrInvalidUuidMode
}

func (u *UuidContainer) GenerateUuid(uuidType UuidType) int64 {
	return u.Current.GenerateUuid(uuidType)
}
