package datastore

import (
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// DataStore represents a data storage containing a series of database shards
type DataStore struct {
	databases []*Database
}

// Choose returns a database instance by sharding key
func (s *DataStore) Choose(key int64) *Database {
	return s.databases[0]
}

// Query returns a new database session in a specific database by sharding key
func (s *DataStore) Query(c *core.Context, key int64) *xorm.Session {
	return s.Choose(key).NewSession(c)
}

// DoTransaction runs a new database transaction in a specific database by sharding key
func (s *DataStore) DoTransaction(key int64, c *core.Context, fn func(sess *xorm.Session) error) (err error) {
	return s.Choose(key).DoTransaction(c, fn)
}

// SyncStructs updates database structs by database models
func (s *DataStore) SyncStructs(beans ...any) error {
	var err error

	for i := 0; i < len(s.databases); i++ {
		err = s.databases[i].engineGroup.Sync2(beans...)

		if err != nil {
			return err
		}
	}

	return err
}

// NewDataStore returns a new data storage by a series of database
func NewDataStore(databases ...*Database) (*DataStore, error) {
	if len(databases) < 1 {
		return nil, errs.ErrDatabaseIsNull
	}

	return &DataStore{
		databases: databases,
	}, nil
}
