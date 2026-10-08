package extservices

import (
	"fmt"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// ItemInput is the editable part of an item
type ItemInput struct {
	Sku          string
	Name         string
	Unit         string
	CostPrice    int64
	SalePrice    int64
	ReorderLevel int64
	TrackStock   bool
}

// ItemService manages the product catalog
type ItemService struct{}

// Items is the item service singleton
var Items = &ItemService{}

func normalizeItemInput(in *ItemInput) error {
	in.Sku = strings.TrimSpace(in.Sku)
	in.Name = strings.TrimSpace(in.Name)
	in.Unit = strings.TrimSpace(in.Unit)

	if in.Sku == "" {
		return exterrs.ErrItemSkuIsEmpty
	}

	if in.Name == "" {
		return exterrs.ErrItemNameIsEmpty
	}

	if in.CostPrice < 0 || in.SalePrice < 0 || in.ReorderLevel < 0 {
		return exterrs.ErrItemPriceInvalid
	}

	if len(in.Sku) > 32 {
		in.Sku = in.Sku[:32]
	}

	if len(in.Name) > 128 {
		in.Name = in.Name[:128]
	}

	if len(in.Unit) > 16 {
		in.Unit = in.Unit[:16]
	}

	return nil
}

// Add creates an item
func (s *ItemService) Add(c core.Context, ownerUid int64, in ItemInput) (*extmodels.Item, error) {
	if err := normalizeItemInput(&in); err != nil {
		return nil, err
	}

	now := nowUnix()
	item := &extmodels.Item{
		OwnerUid: ownerUid, Sku: in.Sku, Name: in.Name, Unit: in.Unit, CostPrice: in.CostPrice,
		SalePrice: in.SalePrice, ReorderLevel: in.ReorderLevel, TrackStock: in.TrackStock,
		CreatedUnixTime: now, UpdatedUnixTime: now,
	}

	err := ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		exists, err := sess.Where("owner_uid=? AND sku=?", ownerUid, in.Sku).Exist(&extmodels.Item{})

		if err != nil {
			return err
		} else if exists {
			return exterrs.ErrItemSkuExists
		}

		_, err = sess.Insert(item)

		return err
	})

	if err != nil {
		return nil, err
	}

	return item, nil
}

// Modify updates an item
func (s *ItemService) Modify(c core.Context, ownerUid int64, itemId int64, in ItemInput) (*extmodels.Item, error) {
	if err := normalizeItemInput(&in); err != nil {
		return nil, err
	}

	var item *extmodels.Item

	err := ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		existing := &extmodels.Item{}
		has, err := sess.Where("owner_uid=? AND item_id=? AND deleted=?", ownerUid, itemId, false).Get(existing)

		if err != nil {
			return err
		} else if !has {
			return exterrs.ErrItemNotFound
		}

		clash, err := sess.Where("owner_uid=? AND sku=? AND item_id<>?", ownerUid, in.Sku, itemId).Exist(&extmodels.Item{})

		if err != nil {
			return err
		} else if clash {
			return exterrs.ErrItemSkuExists
		}

		existing.Sku, existing.Name, existing.Unit = in.Sku, in.Name, in.Unit
		existing.CostPrice, existing.SalePrice, existing.ReorderLevel = in.CostPrice, in.SalePrice, in.ReorderLevel
		existing.TrackStock = in.TrackStock
		existing.UpdatedUnixTime = nowUnix()
		_, err = sess.ID(itemId).Cols("sku", "name", "unit", "cost_price", "sale_price", "reorder_level", "track_stock", "updated_unix_time").Update(existing)
		item = existing

		return err
	})

	if err != nil {
		return nil, err
	}

	return item, nil
}

// Delete soft-deletes an item. Its stock history is kept and the sku is freed for reuse.
func (s *ItemService) Delete(c core.Context, ownerUid int64, itemId int64) error {
	return ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		item := &extmodels.Item{}
		has, err := sess.Where("owner_uid=? AND item_id=? AND deleted=?", ownerUid, itemId, false).Get(item)

		if err != nil {
			return err
		} else if !has {
			return exterrs.ErrItemNotFound
		}

		item.Deleted = true
		item.DeletedUnixTime = nowUnix()
		item.Sku = fmt.Sprintf("%s~%d", item.Sku, item.ItemId)
		_, err = sess.ID(itemId).Cols("deleted", "deleted_unix_time", "sku").Update(item)

		return err
	})
}

// List returns the active items of a business
func (s *ItemService) List(c core.Context, ownerUid int64) ([]*extmodels.Item, error) {
	items := make([]*extmodels.Item, 0)
	err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND deleted=?", ownerUid, false).OrderBy("name, item_id").Find(&items)

	return items, err
}

// Get returns one active item
func (s *ItemService) Get(c core.Context, ownerUid int64, itemId int64) (*extmodels.Item, error) {
	item := &extmodels.Item{}
	has, err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND item_id=? AND deleted=?", ownerUid, itemId, false).Get(item)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, exterrs.ErrItemNotFound
	}

	return item, nil
}
