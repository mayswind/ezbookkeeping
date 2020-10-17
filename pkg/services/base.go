package services

import (
	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/uuid"
)

type ServiceUsingDB struct {
	container *datastore.DataStoreContainer
}

func (s *ServiceUsingDB) UserDB() *datastore.Database {
	return s.container.UserStore.Choose(0)
}

func (s *ServiceUsingDB) TokenDB(uid int64) *datastore.Database {
	return s.container.TokenStore.Choose(uid)
}

func (s *ServiceUsingDB) UserDataDB(uid int64) *datastore.Database {
	return s.container.UserDataStore.Choose(uid)
}

type ServiceUsingConfig struct {
	container *settings.ConfigContainer
}

func (s *ServiceUsingConfig) CurrentConfig() *settings.Config {
	return s.container.Current
}

type ServiceUsingUuid struct {
	container *uuid.UuidContainer
}

func (s *ServiceUsingUuid) GenerateUuid(uuidType uuid.UuidType) int64 {
	return s.container.GenerateUuid(uuidType)
}
