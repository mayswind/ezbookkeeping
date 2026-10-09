package extservices

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// ReportService answers the "how much stock do I have, what is running out, who owes me" questions.
// Everything is read-only and computed from the stock ledger and the sales, so it can never disagree with them.
type ReportService struct{}

// Reports is the report service singleton
var Reports = &ReportService{}

// LocationQty is stock of one item at one location
type LocationQty struct {
	LocationId int64
	Qty        int64
}

// StockValueRow is the stock of one item and what it is worth
type StockValueRow struct {
	Item        *extmodels.Item
	Qty         int64         // scaled by QtyScale
	CostValue   int64         // quantity x the item's cost price, in minor currency units
	RetailValue int64         // quantity x the item's sale price
	Locations   []LocationQty // only when no location filter is applied
}

// StockValueReport is the value of the stock on hand
type StockValueReport struct {
	Rows             []*StockValueRow
	TotalCostValue   int64
	TotalRetailValue int64
}

// StockValue values the stock of every stock-tracked item at its current cost and sale price.
// Pass a location id to value one location only, or 0 for all of them. Items with no stock are left out.
func (s *ReportService) StockValue(c core.Context, ownerUid int64, locationId int64) (*StockValueReport, error) {
	items, levels, err := s.itemsAndLevels(c, ownerUid, locationId)

	if err != nil {
		return nil, err
	}

	report := &StockValueReport{Rows: make([]*StockValueRow, 0)}

	for _, item := range items {
		qty := int64(0)
		var byLocation []LocationQty

		for _, level := range levels[item.ItemId] {
			qty += level.Qty

			if level.Qty != 0 {
				byLocation = append(byLocation, level)
			}
		}

		if qty == 0 {
			continue
		}

		row := &StockValueRow{Item: item, Qty: qty}
		row.CostValue, _ = lineTotal(qty, item.CostPrice)
		row.RetailValue, _ = lineTotal(qty, item.SalePrice)

		if locationId == 0 {
			row.Locations = byLocation
		}

		report.Rows = append(report.Rows, row)
		report.TotalCostValue += row.CostValue
		report.TotalRetailValue += row.RetailValue
	}

	sort.SliceStable(report.Rows, func(i, j int) bool { return report.Rows[i].CostValue > report.Rows[j].CostValue })

	return report, nil
}

// LowStockRow is an item that has reached its reorder level
type LowStockRow struct {
	Item      *extmodels.Item
	Qty       int64
	Shortfall int64 // reorder level minus stock, never negative
}

// LowStock lists stock-tracked items whose stock is at or below their reorder level (items without a reorder level are
// never listed), the most urgent first. Pass a location id to look at one location only, or 0 for all of them.
func (s *ReportService) LowStock(c core.Context, ownerUid int64, locationId int64) ([]*LowStockRow, error) {
	items, levels, err := s.itemsAndLevels(c, ownerUid, locationId)

	if err != nil {
		return nil, err
	}

	rows := make([]*LowStockRow, 0)

	for _, item := range items {
		if item.ReorderLevel <= 0 {
			continue
		}

		qty := int64(0)

		for _, level := range levels[item.ItemId] {
			qty += level.Qty
		}

		if qty <= item.ReorderLevel {
			rows = append(rows, &LowStockRow{Item: item, Qty: qty, Shortfall: item.ReorderLevel - qty})
		}
	}

	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Shortfall > rows[j].Shortfall })

	return rows, nil
}

// itemsAndLevels returns the stock-tracked items and their stock per location
func (s *ReportService) itemsAndLevels(c core.Context, ownerUid int64, locationId int64) ([]*extmodels.Item, map[int64][]LocationQty, error) {
	allItems, err := Items.List(c, ownerUid)

	if err != nil {
		return nil, nil, err
	}

	items := make([]*extmodels.Item, 0, len(allItems))

	for _, item := range allItems {
		if item.TrackStock {
			items = append(items, item)
		}
	}

	rawLevels, err := Stock.Levels(c, ownerUid, 0, locationId)

	if err != nil {
		return nil, nil, err
	}

	levels := make(map[int64][]LocationQty)

	for _, level := range rawLevels {
		levels[level.ItemId] = append(levels[level.ItemId], LocationQty{LocationId: level.LocationId, Qty: level.Qty})
	}

	return items, levels, nil
}

// ReceivableRow is what one customer owes, split by how long ago the sales were made
type ReceivableRow struct {
	Customer       *extmodels.Customer
	Outstanding    int64
	Current        int64 // sales up to 30 days old
	Days31To60     int64
	Days61To90     int64
	Over90         int64
	OldestSaleTime int64
	OpenSales      int
}

// ReceivablesReport is who owes what
type ReceivablesReport struct {
	Rows             []*ReceivableRow
	TotalOutstanding int64
	Current          int64
	Days31To60       int64
	Days61To90       int64
	Over90           int64
}

const secondsPerDay = 24 * 60 * 60

// Receivables lists every customer with unpaid sales, the largest debt first, with the debt aged from the sale date.
// now is the reference time (unix seconds).
func (s *ReportService) Receivables(c core.Context, ownerUid int64, now int64) (*ReceivablesReport, error) {
	sales := make([]*extmodels.Sale, 0)
	err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND voided=? AND customer_id>0 AND paid<total", ownerUid, false).OrderBy("sale_time, sale_id").Find(&sales)

	if err != nil {
		return nil, err
	}

	customers, err := Customers.List(c, ownerUid)

	if err != nil {
		return nil, err
	}

	byId := make(map[int64]*extmodels.Customer, len(customers))

	for _, customer := range customers {
		byId[customer.CustomerId] = customer
	}

	report := &ReceivablesReport{Rows: make([]*ReceivableRow, 0)}
	rows := make(map[int64]*ReceivableRow)

	for _, sale := range sales {
		customer := byId[sale.CustomerId]

		if customer == nil {
			continue
		}

		row := rows[sale.CustomerId]

		if row == nil {
			row = &ReceivableRow{Customer: customer, OldestSaleTime: sale.SaleTime}
			rows[sale.CustomerId] = row
			report.Rows = append(report.Rows, row)
		}

		owed := sale.Outstanding()
		days := (now - sale.SaleTime) / secondsPerDay

		switch {
		case days <= 30:
			row.Current += owed
			report.Current += owed
		case days <= 60:
			row.Days31To60 += owed
			report.Days31To60 += owed
		case days <= 90:
			row.Days61To90 += owed
			report.Days61To90 += owed
		default:
			row.Over90 += owed
			report.Over90 += owed
		}

		row.Outstanding += owed
		row.OpenSales++
		report.TotalOutstanding += owed
	}

	sort.SliceStable(report.Rows, func(i, j int) bool { return report.Rows[i].Outstanding > report.Rows[j].Outstanding })

	return report, nil
}
