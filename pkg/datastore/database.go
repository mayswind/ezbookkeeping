package datastore

import "xorm.io/xorm"

// Database represents a database instance
type Database struct {
	*xorm.EngineGroup
}

// DoTransaction runs a new database transaction
func (db *Database) DoTransaction(fn func(sess *xorm.Session) error) (err error) {
	sess := db.NewSession()
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
