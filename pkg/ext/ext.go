// Package ext is the single entry point of the ext module (owner/staff roles, inventory, credit sales).
// Upstream code calls only SyncTables, DelegationMiddleware and RegisterRoutes from this package;
// keep it that way so merging upstream stays cheap (see docs/EXTENSIONS.md).
package ext

import (
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// SyncTables creates or extends the ext tables. It only ever adds, and never touches upstream tables.
func SyncTables() error {
	if err := datastore.Container.UserStore.SyncStructs(extmodels.GlobalTables()...); err != nil {
		return err
	}

	return datastore.Container.UserDataStore.SyncStructs(extmodels.OwnerTables()...)
}
