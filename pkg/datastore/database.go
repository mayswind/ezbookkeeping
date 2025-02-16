package datastore

import (
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// Database represents a database instance
type Database struct {
	databaseType string
	engineGroup  *xorm.EngineGroup
}

// NewSession starts a new session with the specified context
func (db *Database) NewSession(c core.Context) *xorm.Session {
	return db.engineGroup.Context(NewXOrmContextAdapter(c))
}

// DoTransaction runs a new database transaction
func (db *Database) DoTransaction(c core.Context, fn func(sess *xorm.Session) error) (err error) {
	sess := db.engineGroup.NewSession()

	if c != nil {
		sess.Context(NewXOrmContextAdapter(c))
	}

	defer sess.Close()

	if err = sess.Begin(); err != nil {
		return err
	}

	if err = fn(sess); err != nil {
		_ = sess.Rollback()
		return err
	}

	if err = sess.Commit(); err != nil {
		return err
	}

	return nil
}

// SetSavePoint sets a save point in the current transaction for Postgres
func (db *Database) SetSavePoint(sess *xorm.Session, savePointName string) error {
	if db.databaseType == settings.PostgresDbType {
		_, err := sess.Exec("SAVEPOINT " + savePointName)
		return err
	}

	return nil
}

// RollbackToSavePoint rolls back to the specified save point in the current transaction for Postgres
func (db *Database) RollbackToSavePoint(sess *xorm.Session, savePointName string) error {
	if db.databaseType == settings.PostgresDbType {
		_, err := sess.Exec("ROLLBACK TO SAVEPOINT " + savePointName)
		return err
	}

	return nil
}
