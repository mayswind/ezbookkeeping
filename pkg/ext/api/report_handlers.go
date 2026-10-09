package extapi

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extservices "github.com/mayswind/ezbookkeeping/pkg/ext/services"
)

// PeopleListHandler lists the owner and everyone who has worked in the business, so screens can show who did what.
// Removed members are included (marked inactive) because their records keep their name.
func (h *Handlers) PeopleListHandler(c *core.WebContext) (any, *errs.Error) {
	return h.peopleOf(c, c.GetCurrentUid())
}

// MyPeopleListHandler is the same list for the caller's own business, whichever business they are working in.
// The Team page uses it because it is always about the caller's own team.
func (h *Handlers) MyPeopleListHandler(c *core.WebContext) (any, *errs.Error) {
	return h.peopleOf(c, c.GetActualUid())
}

func (h *Handlers) peopleOf(c *core.WebContext, ownerUid int64) (any, *errs.Error) {
	people, err := extservices.People(c, ownerUid)

	if err != nil {
		return nil, fail(c, "people.list", err)
	}

	result := make([]*PersonView, 0, len(people))

	for _, p := range people {
		result = append(result, &PersonView{Uid: p.Uid, Name: p.Name, Role: p.Role.String(), Active: p.Active})
	}

	return result, nil
}

// StockValueReportHandler values the stock on hand
func (h *Handlers) StockValueReportHandler(c *core.WebContext) (any, *errs.Error) {
	var req ReportLocationRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	report, err := extservices.Reports.StockValue(c, c.GetCurrentUid(), req.LocationId)

	if err != nil {
		return nil, fail(c, "reports.stock_value", err)
	}

	view := &StockValueReportView{Rows: make([]StockValueRowView, 0, len(report.Rows)), TotalCostValue: report.TotalCostValue, TotalRetailValue: report.TotalRetailValue}

	for _, row := range report.Rows {
		item := StockValueRowView{Item: itemView(row.Item), Qty: row.Qty, CostValue: row.CostValue, RetailValue: row.RetailValue}

		for _, l := range row.Locations {
			item.Locations = append(item.Locations, LocationQtyView{LocationId: l.LocationId, Qty: l.Qty})
		}

		view.Rows = append(view.Rows, item)
	}

	return view, nil
}

// LowStockReportHandler lists items that reached their reorder level
func (h *Handlers) LowStockReportHandler(c *core.WebContext) (any, *errs.Error) {
	var req ReportLocationRequest

	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}

	rows, err := extservices.Reports.LowStock(c, c.GetCurrentUid(), req.LocationId)

	if err != nil {
		return nil, fail(c, "reports.low_stock", err)
	}

	result := make([]LowStockRowView, 0, len(rows))

	for _, row := range rows {
		result = append(result, LowStockRowView{Item: itemView(row.Item), Qty: row.Qty, Shortfall: row.Shortfall})
	}

	return result, nil
}

// ReceivablesReportHandler lists who owes what, aged from the sale date
func (h *Handlers) ReceivablesReportHandler(c *core.WebContext) (any, *errs.Error) {
	report, err := extservices.Reports.Receivables(c, c.GetCurrentUid(), time.Now().Unix())

	if err != nil {
		return nil, fail(c, "reports.receivables", err)
	}

	view := &ReceivablesReportView{
		Rows: make([]ReceivableRowView, 0, len(report.Rows)), TotalOutstanding: report.TotalOutstanding,
		Current: report.Current, Days31To60: report.Days31To60, Days61To90: report.Days61To90, Over90: report.Over90,
	}

	for _, row := range report.Rows {
		view.Rows = append(view.Rows, ReceivableRowView{
			Customer: customerView(row.Customer, row.Outstanding), Outstanding: row.Outstanding, Current: row.Current,
			Days31To60: row.Days31To60, Days61To90: row.Days61To90, Over90: row.Over90,
			OldestSaleTime: row.OldestSaleTime, OpenSales: row.OpenSales,
		})
	}

	return view, nil
}
