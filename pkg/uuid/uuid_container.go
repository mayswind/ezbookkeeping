package uuid

import (
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/settings"
)

// UuidContainer contains the current uuid generator
type UuidContainer struct {
	Current UuidGenerator
}

// Initialize a uuid container singleton instance
var (
	Container = &UuidContainer{}
)

// InitializeUuidGenerator initialized the current uuid generator according to the config
func InitializeUuidGenerator(config *settings.Config) error {
	if config.UuidGeneratorType == settings.UUID_GENERATOR_TYPE_INTERNAL {
		generator, err := NewInternalUuidGenerator(config)
		Container.Current = generator

		return err
	}

	return errs.ErrInvalidUuidMode
}

// GenerateUuid returns a new uuid by the current uuid generator
func (u *UuidContainer) GenerateUuid(uuidType UuidType) int64 {
	return u.Current.GenerateUuid(uuidType)
}
