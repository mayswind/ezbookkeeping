package extservices

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// ExportService gives a business owner a copy of the business records the app's own export does not cover:
// items, stock, locations, customers, sales, repayments, the team, receipt details and the activity log.
// Accounts, categories, tags and transactions are exported by the app's own "Data management" page.
type ExportService struct{}

// Exports is the export service singleton
var Exports = &ExportService{}

// exportChunkSize is how many rows are read at a time, so a large business never has to fit in memory. Tests lower it.
var exportChunkSize = 1000

// eachChunk reads a table in pages ordered by id and hands every row to fn
func eachChunk[T any](fetch func(afterId int64, limit int) ([]*T, error), idOf func(*T) int64, fn func(*T) error) error {
	var after int64

	for {
		rows, err := fetch(after, exportChunkSize)

		if err != nil {
			return err
		}

		for _, row := range rows {
			if err = fn(row); err != nil {
				return err
			}
		}

		if len(rows) < exportChunkSize {
			return nil
		}

		after = idOf(rows[len(rows)-1])
	}
}

// ---- cell formatting

// text makes free text safe to open in a spreadsheet: a cell that starts with = + - @ (or a tab or return) would be run
// as a formula, so it gets a quote in front. Names and notes are typed by users, so this matters.
func text(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}

	return s
}

func number(n int64) string { return strconv.FormatInt(n, 10) }

func yesNo(b bool) string {
	if b {
		return "yes"
	}

	return "no"
}

// utc formats a unix time as "2026-10-09 07:30:00" (UTC), or nothing for 0
func utc(unix int64) string {
	if unix <= 0 {
		return ""
	}

	return time.Unix(unix, 0).UTC().Format("2006-01-02 15:04:05")
}

// quantity turns a stored quantity (scaled by 1000) into a plain decimal: 2500 -> "2.5", 3000 -> "3", -750 -> "-0.75"
func quantity(scaled int64) string {
	sign := ""
	abs := scaled

	if scaled < 0 {
		sign, abs = "-", -scaled
	}

	whole, fraction := abs/extmodels.QtyScale, abs%extmodels.QtyScale

	if fraction == 0 {
		return sign + strconv.FormatInt(whole, 10)
	}

	return sign + strconv.FormatInt(whole, 10) + "." + strings.TrimRight(fmt.Sprintf("%03d", fraction), "0")
}

func reasonName(reason extmodels.StockReason) string {
	switch reason {
	case extmodels.StockReasonOpening:
		return "opening stock"
	case extmodels.StockReasonPurchase:
		return "purchase"
	case extmodels.StockReasonSale:
		return "sale"
	case extmodels.StockReasonAdjustment:
		return "adjustment"
	case extmodels.StockReasonTransferOut:
		return "transfer out"
	case extmodels.StockReasonTransferIn:
		return "transfer in"
	case extmodels.StockReasonSaleVoid:
		return "sale voided"
	default:
		return strconv.Itoa(int(reason))
	}
}

// ---- writing

type csvTable struct {
	writer *csv.Writer
}

func newCsvTable(zw *zip.Writer, name string, header ...string) (*csvTable, error) {
	entry, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate, Modified: time.Now()})

	if err != nil {
		return nil, err
	}

	// a byte order mark makes Excel read accented letters and non-latin names correctly
	if _, err = entry.Write([]byte("\xEF\xBB\xBF")); err != nil {
		return nil, err
	}

	table := &csvTable{writer: csv.NewWriter(entry)}

	return table, table.writer.Write(header)
}

func (t *csvTable) row(cells ...string) error { return t.writer.Write(cells) }

func (t *csvTable) close() error {
	t.writer.Flush()
	return t.writer.Error()
}

// WriteZip writes the whole export as a ZIP archive of CSV files to w
func (s *ExportService) WriteZip(c core.Context, ownerUid int64, w io.Writer) error {
	if ownerUid <= 0 {
		return errs.ErrUserIdInvalid
	}

	zw := zip.NewWriter(w)
	db := ownerDB(ownerUid)

	// names for the ids used in the other tables: every item, location and customer, including deleted ones
	items := map[int64]*extmodels.Item{}
	locations := map[int64]*extmodels.Location{}
	customers := map[int64]*extmodels.Customer{}

	if err := eachChunk(func(after int64, limit int) ([]*extmodels.Item, error) {
		rows := make([]*extmodels.Item, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND item_id>?", ownerUid, after).OrderBy("item_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.Item) int64 { return r.ItemId }, func(r *extmodels.Item) error { items[r.ItemId] = r; return nil }); err != nil {
		return err
	}

	if err := eachChunk(func(after int64, limit int) ([]*extmodels.Location, error) {
		rows := make([]*extmodels.Location, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND location_id>?", ownerUid, after).OrderBy("location_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.Location) int64 { return r.LocationId }, func(r *extmodels.Location) error { locations[r.LocationId] = r; return nil }); err != nil {
		return err
	}

	if err := eachChunk(func(after int64, limit int) ([]*extmodels.Customer, error) {
		rows := make([]*extmodels.Customer, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND customer_id>?", ownerUid, after).OrderBy("customer_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.Customer) int64 { return r.CustomerId }, func(r *extmodels.Customer) error { customers[r.CustomerId] = r; return nil }); err != nil {
		return err
	}

	people, err := People(c, ownerUid)

	if err != nil {
		return err
	}

	personName := map[int64]string{}

	for _, p := range people {
		personName[p.Uid] = p.Name
	}

	// anybody who left a record is named, even if they are not on the team list (looked up once)
	by := func(uid int64) string {
		if uid <= 0 {
			return ""
		}

		if _, known := personName[uid]; !known {
			personName[uid] = DisplayName(c, uid)
		}

		return text(personName[uid])
	}
	itemSku := func(id int64) string {
		if item := items[id]; item != nil {
			return text(displaySku(item))
		}

		return ""
	}
	itemName := func(id int64) string {
		if item := items[id]; item != nil {
			return text(item.Name)
		}

		return ""
	}
	locationName := func(id int64) string {
		if l := locations[id]; l != nil {
			return text(l.Name)
		}

		return ""
	}
	customerName := func(id int64) string {
		if cu := customers[id]; cu != nil {
			return text(cu.Name)
		}

		return ""
	}

	owed := map[int64]int64{}
	balances, err := Customers.Balances(c, ownerUid)

	if err != nil {
		return err
	}

	for id, amount := range balances {
		owed[id] = amount
	}

	if err = writeReadme(zw, personName[ownerUid]); err != nil {
		return err
	}

	// ---- items
	t, err := newCsvTable(zw, "items.csv", "item_id", "sku", "name", "unit", "cost_price_minor_units", "sale_price_minor_units", "reorder_level", "tracks_stock", "deleted", "created_utc")

	if err != nil {
		return err
	}

	for _, id := range sortedIds(items) {
		i := items[id]

		if err = t.row(number(i.ItemId), text(displaySku(i)), text(i.Name), text(i.Unit), number(i.CostPrice), number(i.SalePrice), quantity(i.ReorderLevel), yesNo(i.TrackStock), yesNo(i.Deleted), utc(i.CreatedUnixTime)); err != nil {
			return err
		}
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- locations
	if t, err = newCsvTable(zw, "locations.csv", "location_id", "name", "is_default", "deleted"); err != nil {
		return err
	}

	for _, id := range sortedIds(locations) {
		l := locations[id]

		if err = t.row(number(l.LocationId), text(l.Name), yesNo(l.IsDefault), yesNo(l.Deleted)); err != nil {
			return err
		}
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- stock on hand now
	levels, err := Stock.Levels(c, ownerUid, 0, 0)

	if err != nil {
		return err
	}

	if t, err = newCsvTable(zw, "stock_on_hand.csv", "item_id", "sku", "item_name", "location_id", "location", "quantity"); err != nil {
		return err
	}

	for _, level := range levels {
		if err = t.row(number(level.ItemId), itemSku(level.ItemId), itemName(level.ItemId), number(level.LocationId), locationName(level.LocationId), quantity(level.Qty)); err != nil {
			return err
		}
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- stock history
	if t, err = newCsvTable(zw, "stock_movements.csv", "movement_id", "time_utc", "item_id", "sku", "item_name", "location_id", "location", "quantity_change", "reason", "unit_cost_minor_units", "note", "recorded_by"); err != nil {
		return err
	}

	if err = eachChunk(func(after int64, limit int) ([]*extmodels.StockMovement, error) {
		rows := make([]*extmodels.StockMovement, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND movement_id>?", ownerUid, after).OrderBy("movement_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.StockMovement) int64 { return r.MovementId }, func(m *extmodels.StockMovement) error {
		return t.row(number(m.MovementId), utc(m.MovementTime), number(m.ItemId), itemSku(m.ItemId), itemName(m.ItemId), number(m.LocationId), locationName(m.LocationId),
			quantity(m.QtyChange), reasonName(m.Reason), number(m.UnitCost), text(m.Note), by(m.ActorUid))
	}); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- customers
	if t, err = newCsvTable(zw, "customers.csv", "customer_id", "name", "phone", "email", "note", "owes_minor_units", "deleted", "created_utc"); err != nil {
		return err
	}

	for _, id := range sortedIds(customers) {
		cu := customers[id]

		if err = t.row(number(cu.CustomerId), text(cu.Name), text(cu.Phone), text(cu.Email), text(cu.Note), number(owed[cu.CustomerId]), yesNo(cu.Deleted), utc(cu.CreatedUnixTime)); err != nil {
			return err
		}
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- sales
	if t, err = newCsvTable(zw, "sales.csv", "sale_id", "time_utc", "customer_id", "customer", "location_id", "location", "subtotal_minor_units", "discount_minor_units", "total_minor_units",
		"paid_minor_units", "owed_minor_units", "voided", "payment_account_id", "receivable_account_id", "category_id", "note", "recorded_by"); err != nil {
		return err
	}

	if err = eachChunk(func(after int64, limit int) ([]*extmodels.Sale, error) {
		rows := make([]*extmodels.Sale, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND sale_id>?", ownerUid, after).OrderBy("sale_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.Sale) int64 { return r.SaleId }, func(sale *extmodels.Sale) error {
		return t.row(number(sale.SaleId), utc(sale.SaleTime), number(sale.CustomerId), customerName(sale.CustomerId), number(sale.LocationId), locationName(sale.LocationId),
			number(sale.Subtotal), number(sale.Discount), number(sale.Total), number(sale.Paid), number(sale.Outstanding()), yesNo(sale.Voided),
			number(sale.PaymentAccountId), number(sale.ReceivableAccountId), number(sale.CategoryId), text(sale.Note), by(sale.ActorUid))
	}); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	if t, err = newCsvTable(zw, "sale_lines.csv", "sale_id", "item_id", "sku", "item_name", "quantity", "unit_price_minor_units", "line_total_minor_units"); err != nil {
		return err
	}

	if err = eachChunk(func(after int64, limit int) ([]*extmodels.SaleLine, error) {
		rows := make([]*extmodels.SaleLine, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND line_id>?", ownerUid, after).OrderBy("line_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.SaleLine) int64 { return r.LineId }, func(l *extmodels.SaleLine) error {
		return t.row(number(l.SaleId), number(l.ItemId), itemSku(l.ItemId), itemName(l.ItemId), quantity(l.Qty), number(l.UnitPrice), number(l.LineTotal))
	}); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- repayments
	if t, err = newCsvTable(zw, "repayments.csv", "repayment_id", "time_utc", "customer_id", "customer", "amount_minor_units", "payment_account_id", "receivable_account_id", "note", "recorded_by"); err != nil {
		return err
	}

	if err = eachChunk(func(after int64, limit int) ([]*extmodels.Repayment, error) {
		rows := make([]*extmodels.Repayment, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND repayment_id>?", ownerUid, after).OrderBy("repayment_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.Repayment) int64 { return r.RepaymentId }, func(r *extmodels.Repayment) error {
		return t.row(number(r.RepaymentId), utc(r.RepaymentTime), number(r.CustomerId), customerName(r.CustomerId), number(r.Amount), number(r.PaymentAccountId), number(r.ReceivableAccountId), text(r.Note), by(r.ActorUid))
	}); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	if t, err = newCsvTable(zw, "repayment_allocations.csv", "repayment_id", "sale_id", "amount_minor_units"); err != nil {
		return err
	}

	if err = eachChunk(func(after int64, limit int) ([]*extmodels.RepaymentAllocation, error) {
		rows := make([]*extmodels.RepaymentAllocation, 0)
		return rows, db.NewSession(c).Where("owner_uid=? AND allocation_id>?", ownerUid, after).OrderBy("allocation_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.RepaymentAllocation) int64 { return r.AllocationId }, func(a *extmodels.RepaymentAllocation) error {
		return t.row(number(a.RepaymentId), number(a.SaleId), number(a.Amount))
	}); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- team and activity
	memberships, err := Memberships.ListAllByOwner(c, ownerUid)

	if err != nil {
		return err
	}

	if t, err = newCsvTable(zw, "team.csv", "name", "role", "status"); err != nil {
		return err
	}

	if err = t.row(text(personName[ownerUid]), "owner", "active"); err != nil {
		return err
	}

	for _, m := range memberships {
		status := "removed"

		if m.Status == extmodels.MembershipStatusActive {
			status = "active"
		}

		if err = t.row(text(personName[m.StaffUid]), m.Role.String(), status); err != nil {
			return err
		}
	}

	if err = t.close(); err != nil {
		return err
	}

	if t, err = newCsvTable(zw, "activity_log.csv", "time_utc", "person", "role", "action", "thing", "thing_id", "result_code"); err != nil {
		return err
	}

	if err = eachChunk(func(after int64, limit int) ([]*extmodels.AuditLog, error) {
		rows := make([]*extmodels.AuditLog, 0)
		return rows, globalDB().NewSession(c).Where("owner_uid=? AND audit_id>?", ownerUid, after).OrderBy("audit_id").Limit(limit).Find(&rows)
	}, func(r *extmodels.AuditLog) int64 { return r.AuditId }, func(a *extmodels.AuditLog) error {
		return t.row(utc(a.CreatedUnix), by(a.ActorUid), a.Role.String(), text(a.Action), text(a.EntityType), number(a.EntityId), strconv.Itoa(a.Status))
	}); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	// ---- receipt details
	profile, err := BusinessProfiles.Get(c, ownerUid)

	if err != nil {
		return err
	}

	if t, err = newCsvTable(zw, "receipt_details.csv", "business_name_on_receipts", "address", "phone", "message_at_the_bottom"); err != nil {
		return err
	}

	if err = t.row(text(profile.ReceiptName), text(profile.Address), text(profile.Phone), text(profile.Footer)); err != nil {
		return err
	}

	if err = t.close(); err != nil {
		return err
	}

	return zw.Close()
}

// displaySku is the code the person typed. Deleting an item renames its code to "CODE~id" so the code can be reused;
// that internal marker is taken off again here.
func displaySku(item *extmodels.Item) string {
	if item.Deleted {
		return strings.TrimSuffix(item.Sku, "~"+number(item.ItemId))
	}

	return item.Sku
}

func sortedIds[T any](m map[int64]*T) []int64 {
	ids := make([]int64, 0, len(m))

	for id := range m {
		ids = append(ids, id)
	}

	for i := 1; i < len(ids); i++ { // small maps; insertion sort keeps the export in id order without another import
		for j := i; j > 0 && ids[j-1] > ids[j]; j-- {
			ids[j-1], ids[j] = ids[j], ids[j-1]
		}
	}

	return ids
}

func writeReadme(zw *zip.Writer, owner string) error {
	entry, err := zw.CreateHeader(&zip.FileHeader{Name: "README.txt", Method: zip.Deflate, Modified: time.Now()})

	if err != nil {
		return err
	}

	_, err = io.WriteString(entry, strings.ReplaceAll(`Business data export
====================
Business: `+owner+`
Created: `+time.Now().UTC().Format("2006-01-02 15:04:05")+` UTC

Each .csv file is a table you can open in Excel, Google Sheets or LibreOffice.

How to read the numbers
- Amounts are whole numbers in the smallest unit of your currency, called "minor units": 1500 means 15.00 for a currency with two decimals
  (cents, kobo, pence). Columns ending in _minor_units hold such amounts.
- Quantities are ordinary decimals: 2.5 means two and a half units.
- Times are UTC (universal time), written year-month-day hour:minute:second.
- yes / no columns answer a yes-or-no question.
- A text cell that starts with a quote mark has had it added to stop spreadsheets running typed text as a formula; ignore it.

Files
- items.csv                 Your items, including deleted ones (the sales that used them still refer to them).
- locations.csv             Your stores or warehouses.
- stock_on_hand.csv         Stock per item and location right now.
- stock_movements.csv       Every change to stock: what, where, why and who. The stock on hand is the sum of these.
- customers.csv             Your customers and how much each owes now.
- sales.csv, sale_lines.csv Every sale and the items on it. Voided sales are included and marked.
- repayments.csv            Money received from customers against what they owed.
- repayment_allocations.csv Which sales each repayment paid off.
- team.csv                  You and the people who work in your business.
- activity_log.csv          What your managers and staff did in the business, and when.
- receipt_details.csv       The details printed on your receipts.

Not in this file
Your accounts, categories, tags and transactions (including the money side of every sale) are exported from the app itself:
Settings > Data management. Account ids in sales.csv (payment_account_id and so on) match those accounts.
`, "\n", "\r\n"))

	return err
}
