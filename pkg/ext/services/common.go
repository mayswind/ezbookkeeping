// Package extservices contains the business logic of the ext module (staff, inventory, sales).
package extservices

import (
	"time"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"

	"github.com/mayswind/ezbookkeeping/pkg/datastore"
)

// globalDB returns the database holding data that must be looked up across owners (memberships, audit log)
func globalDB() *datastore.Database {
	return datastore.Container.UserStore.Choose(0)
}

// ownerDB returns the database holding the business data of an owner
func ownerDB(ownerUid int64) *datastore.Database {
	return datastore.Container.UserDataStore.Choose(ownerUid)
}

// nowUnix returns the current unix time, replaceable in tests
var nowUnix = func() int64 {
	return time.Now().Unix()
}

// xormSession is a short alias to keep service signatures readable
type xormSession = xorm.Session

// lockRows makes the next query take row locks (SELECT ... FOR UPDATE) on databases that support it.
// This serializes concurrent stock changes of one item or location on MySQL and PostgreSQL. SQLite has no row locks,
// but it only ever runs one writer at a time, so nothing is needed there.
func lockRows(sess *xormSession) *xormSession {
	switch sess.Engine().Dialect().URI().DBType {
	case schemas.MYSQL, schemas.POSTGRES:
		return sess.ForUpdate()
	default:
		return sess
	}
}
