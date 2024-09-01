package uuid

import (
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// UuidContainer contains the current uuid generator
type UuidContainer struct {
	Current UuidGenerator
}

// Initialize a uuid container singleton instance
var (
	Container = &UuidContainer{}
)

// InitializeUuidGenerator initializes the current uuid generator according to the config
func InitializeUuidGenerator(config *settings.Config) error {
	if config.UuidGeneratorType == settings.InternalUuidGeneratorType {
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

// GenerateUuids returns new uuids by the current uuid generator
func (u *UuidContainer) GenerateUuids(uuidType UuidType, count uint16) []int64 {
	return u.Current.GenerateUuids(uuidType, count)
}
