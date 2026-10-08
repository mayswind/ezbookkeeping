package extapi

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
)

// Handlers below act on the business named by c.GetCurrentUid(): the caller's own, or the owner's when
// a manager or staff member works through the business header (see pkg/ext/middleware).

func locationView(l *extmodels.Location) *LocationView {
	return &LocationView{Id: l.LocationId, Name: l.Name, IsDefault: l.IsDefault}
}

func itemView(i *extmodels.Item) *ItemView {
	return &ItemView{Id: i.ItemId, Sku: i.Sku, Name: i.Name, Unit: i.Unit, CostPrice: i.CostPrice, SalePrice: i.SalePrice, ReorderLevel: i.ReorderLevel, TrackStock: i.TrackStock}
}

func movementView(m *extmodels.StockMovement) *MovementView {
	return &MovementView{Id: m.MovementId, ItemId: m.ItemId, LocationId: m.LocationId, QtyChange: m.QtyChange, Reason: int(m.Reason),
		RefType: m.RefType, RefId: m.RefId, UnitCost: m.UnitCost, Note: m.Note, ActorUid: m.ActorUid, Time: m.MovementTime}
}

// LocationListHandler lists locations
func (h *Handlers) LocationListHandler(c *core.WebContext) (any, *errs.Error) {
	locations, err := extservices.Locations.List(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "locations.list", err)
	}

	result := make([]*LocationView, 0, len(locations))

	for _, l := range locations {
		result = append(result, locationView(l))
	}

	return result, nil
}

// LocationAddHandler adds a location
func (h *Handlers) LocationAddHandler(c *core.WebContext) (any, *errs.Error) {
	var req LocationRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	location, err := extservices.Locations.Add(c, c.GetCurrentUid(), req.Name)

	if err != nil {
		return nil, fail(c, "locations.add", err)
	}

	return locationView(location), nil
}

// LocationModifyHandler renames a location
func (h *Handlers) LocationModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req LocationRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	location, err := extservices.Locations.Modify(c, c.GetCurrentUid(), req.Id, req.Name)

	if err != nil {
		return nil, fail(c, "locations.modify", err)
	}

	return locationView(location), nil
}

// LocationDeleteHandler deletes an empty location
func (h *Handlers) LocationDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req IdRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Locations.Delete(c, c.GetCurrentUid(), req.Id); err != nil {
		return nil, fail(c, "locations.delete", err)
	}

	return true, nil
}

// ItemListHandler lists items
func (h *Handlers) ItemListHandler(c *core.WebContext) (any, *errs.Error) {
	items, err := extservices.Items.List(c, c.GetCurrentUid())

	if err != nil {
		return nil, fail(c, "items.list", err)
	}

	result := make([]*ItemView, 0, len(items))

	for _, i := range items {
		result = append(result, itemView(i))
	}

	return result, nil
}

func itemInput(req *ItemRequest) extservices.ItemInput {
	return extservices.ItemInput{Sku: req.Sku, Name: req.Name, Unit: req.Unit, CostPrice: req.CostPrice, SalePrice: req.SalePrice, ReorderLevel: req.ReorderLevel, TrackStock: req.TrackStock}
}

// ItemAddHandler adds an item
func (h *Handlers) ItemAddHandler(c *core.WebContext) (any, *errs.Error) {
	var req ItemRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	item, err := extservices.Items.Add(c, c.GetCurrentUid(), itemInput(&req))

	if err != nil {
		return nil, fail(c, "items.add", err)
	}

	return itemView(item), nil
}

// ItemModifyHandler updates an item
func (h *Handlers) ItemModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req ItemRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	item, err := extservices.Items.Modify(c, c.GetCurrentUid(), req.Id, itemInput(&req))

	if err != nil {
		return nil, fail(c, "items.modify", err)
	}

	return itemView(item), nil
}

// ItemDeleteHandler deletes an item
func (h *Handlers) ItemDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req IdRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Items.Delete(c, c.GetCurrentUid(), req.Id); err != nil {
		return nil, fail(c, "items.delete", err)
	}

	return true, nil
}

// StockLevelsHandler returns stock on hand per item and location
func (h *Handlers) StockLevelsHandler(c *core.WebContext) (any, *errs.Error) {
	var req StockLevelRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	levels, err := extservices.Stock.Levels(c, c.GetCurrentUid(), req.ItemId, req.LocationId)

	if err != nil {
		return nil, fail(c, "items.stock", err)
	}

	result := make([]*StockLevelView, 0, len(levels))

	for _, l := range levels {
		result = append(result, &StockLevelView{ItemId: l.ItemId, LocationId: l.LocationId, Qty: l.Qty})
	}

	return result, nil
}

func stockInput(req *StockChangeRequest) extservices.StockInput {
	return extservices.StockInput{ItemId: req.ItemId, LocationId: req.LocationId, Qty: req.Qty, UnitCost: req.UnitCost, Note: req.Note, Time: req.Time, Opening: req.Opening}
}

// StockReceiveHandler adds stock
func (h *Handlers) StockReceiveHandler(c *core.WebContext) (any, *errs.Error) {
	var req StockChangeRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	movement, err := extservices.Stock.Receive(c, c.GetCurrentUid(), c.GetActualUid(), stockInput(&req))

	if err != nil {
		return nil, fail(c, "stock.receive", err)
	}

	return movementView(movement), nil
}

// StockAdjustHandler corrects stock by a positive or negative quantity
func (h *Handlers) StockAdjustHandler(c *core.WebContext) (any, *errs.Error) {
	var req StockChangeRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	movement, err := extservices.Stock.Adjust(c, c.GetCurrentUid(), c.GetActualUid(), stockInput(&req))

	if err != nil {
		return nil, fail(c, "stock.adjust", err)
	}

	return movementView(movement), nil
}

// StockTransferHandler moves stock between locations
func (h *Handlers) StockTransferHandler(c *core.WebContext) (any, *errs.Error) {
	var req StockTransferRequest

	if err := bindBody(c, &req); err != nil {
		return nil, err
	}

	if err := extservices.Stock.Transfer(c, c.GetCurrentUid(), c.GetActualUid(), req.ItemId, req.FromLocationId, req.ToLocationId, req.Qty, req.Note); err != nil {
		return nil, fail(c, "stock.transfer", err)
	}

	return true, nil
}

// StockMovementsHandler lists the stock ledger
func (h *Handlers) StockMovementsHandler(c *core.WebContext) (any, *errs.Error) {
	var req MovementListRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	movements, err := extservices.Stock.Movements(c, c.GetCurrentUid(), extservices.MovementFilter{ItemId: req.ItemId, LocationId: req.LocationId, BeforeId: req.BeforeId, Limit: req.Limit})

	if err != nil {
		return nil, fail(c, "stock.movements", err)
	}

	result := make([]*MovementView, 0, len(movements))

	for _, m := range movements {
		result = append(result, movementView(m))
	}

	return result, nil
}
