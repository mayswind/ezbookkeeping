package datastore

import (
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// Database represents a database instance
type Database interface {
	// NewSession starts a new session with the specified context
	NewSession(c core.Context) DatabaseSession

	// DoTransaction runs a new database transaction
	DoTransaction(c core.Context, fn func(sess DatabaseSession) error) error

	// SetSavePoint sets a save point in the current transaction for Postgres
	SetSavePoint(sess DatabaseSession, savePointName string) error

	// RollbackToSavePoint rolls back to the specified save point in the current transaction for Postgres
	RollbackToSavePoint(sess DatabaseSession, savePointName string) error

	// SyncStructs updates database structs by database models
	SyncStructs(beans ...any) error
}

// XOrmDatabase represents a database instance using xorm
type XOrmDatabase struct {
	databaseType string
	engineGroup  *xorm.EngineGroup
}

// NewSession starts a new session with the specified context
func (db *XOrmDatabase) NewSession(c core.Context) DatabaseSession {
	return &XOrmSession{
		Session: db.engineGroup.Context(NewXOrmContextAdapter(c)),
	}
}

// DoTransaction runs a new database transaction
func (db *XOrmDatabase) DoTransaction(c core.Context, fn func(sess DatabaseSession) error) (err error) {
	sess := db.engineGroup.NewSession()

	if c != nil {
		sess.Context(NewXOrmContextAdapter(c))
	}

	defer sess.Close()

	if err = sess.Begin(); err != nil {
		return err
	}

	xormSession := &XOrmSession{
		Session: sess,
	}

	if err = fn(xormSession); err != nil {
		_ = sess.Rollback()
		return err
	}

	if err = sess.Commit(); err != nil {
		return err
	}

	return nil
}

// SetSavePoint sets a save point in the current transaction for Postgres
func (db *XOrmDatabase) SetSavePoint(sess DatabaseSession, savePointName string) error {
	if db.databaseType == settings.PostgresDbType {
		_, err := sess.Exec("SAVEPOINT " + savePointName)
		return err
	}

	return nil
}

// RollbackToSavePoint rolls back to the specified save point in the current transaction for Postgres
func (db *XOrmDatabase) RollbackToSavePoint(sess DatabaseSession, savePointName string) error {
	if db.databaseType == settings.PostgresDbType {
		_, err := sess.Exec("ROLLBACK TO SAVEPOINT " + savePointName)
		return err
	}

	return nil
}

// SyncStructs updates database structs by database models
func (db *XOrmDatabase) SyncStructs(beans ...any) error {
	return db.engineGroup.Sync2(beans...)
}
