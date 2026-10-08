package extservices

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// StockService manages the stock ledger. Stock on hand is always the sum of movements.
type StockService struct{}

// Stock is the stock service singleton
var Stock = &StockService{}

// Levels returns stock on hand grouped by item and location. Pass 0 to not filter by item or location.
func (s *StockService) Levels(c core.Context, ownerUid int64, itemId int64, locationId int64) ([]*extmodels.StockLevel, error) {
	levels := make([]*extmodels.StockLevel, 0)
	sess := ownerDB(ownerUid).NewSession(c).Table("ext_stock_movement").Where("owner_uid=?", ownerUid)

	if itemId > 0 {
		sess = sess.And("item_id=?", itemId)
	}

	if locationId > 0 {
		sess = sess.And("location_id=?", locationId)
	}

	err := sess.Select("item_id, location_id, SUM(qty_change) AS qty").GroupBy("item_id, location_id").OrderBy("item_id, location_id").Find(&levels)

	return levels, err
}

// onHand returns the stock of an item at a location inside an open session
func onHand(sess *xormSession, ownerUid int64, itemId int64, locationId int64) (int64, error) {
	return sess.Where("owner_uid=? AND item_id=? AND location_id=?", ownerUid, itemId, locationId).SumInt(&extmodels.StockMovement{}, "qty_change")
}

// requireStockItem loads a stock-tracked item and locks its row until the transaction ends, so concurrent
// changes to the same item's stock queue up instead of both passing the stock check.
func requireStockItem(sess *xormSession, ownerUid int64, itemId int64) (*extmodels.Item, error) {
	item := &extmodels.Item{}
	has, err := lockRows(sess).Where("owner_uid=? AND item_id=? AND deleted=?", ownerUid, itemId, false).Get(item)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, exterrs.ErrItemNotFound
	}

	if !item.TrackStock {
		return nil, exterrs.ErrItemDoesNotTrackStock
	}

	return item, nil
}

// requireLocation checks that the location is still active and locks its row until the transaction ends,
// so a location cannot be deleted while stock is being added to it. Always lock locations before items.
func requireLocation(sess *xormSession, ownerUid int64, locationId int64) error {
	has, err := lockRows(sess).Where("owner_uid=? AND location_id=? AND deleted=?", ownerUid, locationId, false).Get(&extmodels.Location{})

	if err != nil {
		return err
	} else if !has {
		return exterrs.ErrLocationNotFound
	}

	return nil
}

// StockInput describes a manual stock change
type StockInput struct {
	ItemId     int64
	LocationId int64 // 0 means the default location
	Qty        int64 // scaled by QtyScale
	UnitCost   int64
	Note       string
	Time       int64 // 0 means now
	Opening    bool  // receive only: record as opening stock instead of a purchase
}

func (s *StockService) apply(c core.Context, ownerUid int64, actorUid int64, in StockInput, reason extmodels.StockReason, sign int64) (*extmodels.StockMovement, error) {
	location, err := Locations.Resolve(c, ownerUid, in.LocationId)

	if err != nil {
		return nil, err
	}

	if in.Qty > maxQty || in.Qty < -maxQty || in.UnitCost > maxUnitPrice {
		return nil, exterrs.ErrStockTooLarge
	}

	if in.Time <= 0 {
		in.Time = nowUnix()
	}

	var movement *extmodels.StockMovement

	err = ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		if err := requireLocation(sess, ownerUid, location.LocationId); err != nil {
			return err
		}

		if _, err := requireStockItem(sess, ownerUid, in.ItemId); err != nil {
			return err
		}

		change := in.Qty * sign
		current, err := onHand(sess, ownerUid, in.ItemId, location.LocationId)

		if err != nil {
			return err
		} else if current+change < 0 {
			return exterrs.ErrInsufficientStock
		} else if current+change > maxStockOnHand {
			return exterrs.ErrStockTooLarge
		}

		movement = &extmodels.StockMovement{
			OwnerUid: ownerUid, ItemId: in.ItemId, LocationId: location.LocationId, QtyChange: change, Reason: reason,
			UnitCost: in.UnitCost, Note: in.Note, ActorUid: actorUid, MovementTime: in.Time, CreatedUnix: nowUnix(),
		}
		_, err = sess.Insert(movement)

		return err
	})

	if err != nil {
		return nil, err
	}

	return movement, nil
}

// Receive adds stock, from a purchase or as opening stock
func (s *StockService) Receive(c core.Context, ownerUid int64, actorUid int64, in StockInput) (*extmodels.StockMovement, error) {
	if in.Qty <= 0 || in.UnitCost < 0 {
		return nil, exterrs.ErrQuantityInvalid
	}

	reason := extmodels.StockReasonPurchase

	if in.Opening {
		reason = extmodels.StockReasonOpening
	}

	return s.apply(c, ownerUid, actorUid, in, reason, 1)
}

// Adjust corrects stock by a positive or negative quantity (for example after a stock count)
func (s *StockService) Adjust(c core.Context, ownerUid int64, actorUid int64, in StockInput) (*extmodels.StockMovement, error) {
	if in.Qty == 0 {
		return nil, exterrs.ErrQuantityInvalid
	}

	return s.apply(c, ownerUid, actorUid, in, extmodels.StockReasonAdjustment, 1)
}

// Transfer moves stock between two locations as a pair of movements
func (s *StockService) Transfer(c core.Context, ownerUid int64, actorUid int64, itemId int64, fromLocationId int64, toLocationId int64, qty int64, note string) error {
	if qty <= 0 {
		return exterrs.ErrQuantityInvalid
	}

	from, err := Locations.Resolve(c, ownerUid, fromLocationId)

	if err != nil {
		return err
	}

	to, err := Locations.Resolve(c, ownerUid, toLocationId)

	if err != nil {
		return err
	}

	if from.LocationId == to.LocationId {
		return exterrs.ErrSameLocation
	}

	if qty > maxQty {
		return exterrs.ErrStockTooLarge
	}

	now := nowUnix()

	return ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		first, second := from.LocationId, to.LocationId

		if first > second { // always lock in ascending id order so two opposite transfers cannot deadlock
			first, second = second, first
		}

		for _, locationId := range []int64{first, second} {
			if err := requireLocation(sess, ownerUid, locationId); err != nil {
				return err
			}
		}

		if _, err := requireStockItem(sess, ownerUid, itemId); err != nil {
			return err
		}

		current, err := onHand(sess, ownerUid, itemId, from.LocationId)

		if err != nil {
			return err
		} else if current < qty {
			return exterrs.ErrInsufficientStock
		}

		destination, err := onHand(sess, ownerUid, itemId, to.LocationId)

		if err != nil {
			return err
		} else if destination+qty > maxStockOnHand {
			return exterrs.ErrStockTooLarge
		}

		out := &extmodels.StockMovement{OwnerUid: ownerUid, ItemId: itemId, LocationId: from.LocationId, QtyChange: -qty,
			Reason: extmodels.StockReasonTransferOut, RefType: "transfer", Note: note, ActorUid: actorUid, MovementTime: now, CreatedUnix: now}

		if _, err = sess.Insert(out); err != nil {
			return err
		}

		in := &extmodels.StockMovement{OwnerUid: ownerUid, ItemId: itemId, LocationId: to.LocationId, QtyChange: qty,
			Reason: extmodels.StockReasonTransferIn, RefType: "transfer", RefId: out.MovementId, Note: note, ActorUid: actorUid, MovementTime: now, CreatedUnix: now}
		_, err = sess.Insert(in)

		return err
	})
}

// MovementFilter filters the stock ledger listing
type MovementFilter struct {
	ItemId     int64
	LocationId int64
	BeforeId   int64
	Limit      int
}

// Movements lists the stock ledger, newest first
func (s *StockService) Movements(c core.Context, ownerUid int64, f MovementFilter) ([]*extmodels.StockMovement, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}

	movements := make([]*extmodels.StockMovement, 0, f.Limit)
	sess := ownerDB(ownerUid).NewSession(c).Where("owner_uid=?", ownerUid)

	if f.ItemId > 0 {
		sess = sess.And("item_id=?", f.ItemId)
	}

	if f.LocationId > 0 {
		sess = sess.And("location_id=?", f.LocationId)
	}

	if f.BeforeId > 0 {
		sess = sess.And("movement_id<?", f.BeforeId)
	}

	err := sess.OrderBy("movement_id desc").Limit(f.Limit).Find(&movements)

	return movements, err
}
