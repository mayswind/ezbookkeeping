package datastore

import (
	"xorm.io/xorm"

	"github.com/mayswind/lab/pkg/errs"
)

type DataStore struct {
	databases []*Database
}

func (s *DataStore) Choose(key int64) *Database {
	return s.databases[0]
}

func (s *DataStore) Query(key int64) *xorm.Session {
	return s.Choose(key).NewSession()
}

func (s *DataStore) DoTransaction(key int64, fn func(sess *xorm.Session) error) (err error) {
	return s.Choose(key).DoTransaction(fn)
}

func (s *DataStore) SyncStructs(beans ...interface{}) error {
	var err error

	for i := 0; i < len(s.databases); i++ {
		err = s.databases[i].Sync2(beans...)

		if err != nil {
			return err
		}
	}

	return err
}

func NewDataStore(databases ...*Database) (*DataStore, error) {
	if len(databases) < 1 {
		return nil, errs.ErrDatabaseIsNull
	}

	return &DataStore{
		databases: databases,
	}, nil
}
