package extservices

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// ClearedCounts says how many rows of each kind ClearBusinessRecords removed
type ClearedCounts struct {
	Items, Locations, StockMovements, Customers, Sales, SaleLines, Repayments, RepaymentAllocations int64
}

// ClearBusinessRecords removes a business's items, locations, stock history, customers, sales, sale lines, repayments and
// repayment allocations, in one database transaction. It runs after the app's own "Clear All Data" has succeeded, so
// that button leaves nothing of the business behind.
//
// Kept on purpose, because they are settings and accountability rather than records of trade: the team and their
// roles, the receipt details, the Terms acceptances, the person's own settings, and the activity log (which also gets a
// line saying the data was cleared). Clearing again is harmless and removes nothing more.
func (s *ClearService) ClearBusinessRecords(c core.Context, ownerUid int64, actorUid int64, clientIp string) (*ClearedCounts, error) {
	if ownerUid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	counts := &ClearedCounts{}

	// children before parents
	steps := []struct {
		bean  any
		count *int64
	}{
		{&extmodels.RepaymentAllocation{}, &counts.RepaymentAllocations},
		{&extmodels.Repayment{}, &counts.Repayments},
		{&extmodels.SaleLine{}, &counts.SaleLines},
		{&extmodels.Sale{}, &counts.Sales},
		{&extmodels.StockMovement{}, &counts.StockMovements},
		{&extmodels.Customer{}, &counts.Customers},
		{&extmodels.Item{}, &counts.Items},
		{&extmodels.Location{}, &counts.Locations},
	}

	err := ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		for _, step := range steps {
			removed, err := sess.Where("owner_uid=?", ownerUid).Delete(step.bean)

			if err != nil {
				return err
			}

			*step.count = removed
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = Audit.Record(c, &extmodels.AuditLog{
		OwnerUid: ownerUid, ActorUid: actorUid, Role: extmodels.RoleOwner, Method: "POST", Path: "/api/v1/data/clear/all.json",
		Status: 200, ClientIp: clientIp, Action: "data.clear_all", EntityType: "data",
	})

	return counts, err
}

// ClearService removes a business's records
type ClearService struct{}

// Clear is the clear service singleton
var Clear = &ClearService{}
