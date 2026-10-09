package extservices

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// readZip returns every file of an export, each as parsed CSV rows (header first), and the raw README
func readZip(t *testing.T, data []byte) (tables map[string][][]string, raw map[string]string) {
	t.Helper()

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	tables, raw = map[string][][]string{}, map[string]string{}

	for _, file := range zr.File {
		rc, err := file.Open()
		require.NoError(t, err)
		body, err := io.ReadAll(rc)
		require.NoError(t, err)
		require.NoError(t, rc.Close())
		raw[file.Name] = string(body)

		if strings.HasSuffix(file.Name, ".csv") {
			require.True(t, bytes.HasPrefix(body, []byte("\xEF\xBB\xBF")), file.Name+" should start with a byte order mark so Excel reads it correctly")
			rows, err := csv.NewReader(bytes.NewReader(bytes.TrimPrefix(body, []byte("\xEF\xBB\xBF")))).ReadAll()
			require.NoError(t, err, file.Name)
			tables[file.Name] = rows
		}
	}

	return tables, raw
}

func (f *fixture) export(t *testing.T, ownerUid int64) (map[string][][]string, map[string]string) {
	t.Helper()

	var buffer bytes.Buffer
	require.NoError(t, Exports.WriteZip(f.c, ownerUid, &buffer))

	return readZip(t, buffer.Bytes())
}

func column(t *testing.T, rows [][]string, name string) int {
	t.Helper()

	for i, header := range rows[0] {
		if header == name {
			return i
		}
	}

	require.Failf(t, "column not found", "%q in %v", name, rows[0])

	return -1
}

func TestExport_EmptyBusinessStillHasEveryFileAndHeaders(t *testing.T) {
	f := newFixture(t)

	tables, raw := f.export(t, f.ownerUid)

	for _, name := range []string{"items.csv", "locations.csv", "stock_on_hand.csv", "stock_movements.csv", "customers.csv", "sales.csv", "sale_lines.csv",
		"repayments.csv", "repayment_allocations.csv", "team.csv", "activity_log.csv", "receipt_details.csv"} {
		require.Contains(t, tables, name)
		assert.NotEmpty(t, tables[name][0], name+" has a header")
	}

	assert.Contains(t, raw["README.txt"], "minor units")
	assert.Len(t, tables["items.csv"], 1, "no items: just the header")
	assert.Equal(t, [][]string{{"name", "role", "status"}, {"owner", "owner", "active"}}, tables["team.csv"], "the owner is listed even when nobody else is")
}

func TestExport_ContainsTheBusinessInReadableForm(t *testing.T) {
	f := newFixture(t)
	staff := f.newUser(t, "sam")
	rice := f.newItem(t, "RICE", 150000, true) // cost 75000
	shop, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	gone := f.newItem(t, "GONE", 100, false)
	require.NoError(t, Items.Delete(f.c, f.ownerUid, gone.ItemId))

	f.receive(t, rice.ItemId, main.LocationId, 10500) // 10.5
	f.receive(t, rice.ItemId, shop.LocationId, 2000)
	ada, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada Obi", Phone: "0800", Email: "ada@example.com"})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: rice.ItemId, Qty: 2500}) // 2.5 x 1500.00 = 3750.00
	in.CustomerId = ada.CustomerId
	in.AmountPaid = 100000
	in.Note = "weekend"
	sale, err := Sales.Create(f.c, f.ownerUid, staff, in)
	require.NoError(t, err)
	_, err = Repayments.Add(f.c, f.ownerUid, staff, f.repay(ada.CustomerId, 75000, 0))
	require.NoError(t, err)
	_, err = BusinessProfiles.Save(f.c, f.ownerUid, ProfileInput{ReceiptName: "Ada Stores", Footer: "Thank you"})
	require.NoError(t, err)

	tables, _ := f.export(t, f.ownerUid)

	items := tables["items.csv"]
	require.Len(t, items, 3, "header, the live item and the deleted one")
	row := items[1]
	assert.Equal(t, "RICE", row[column(t, items, "sku")])
	assert.Equal(t, "150000", row[column(t, items, "sale_price_minor_units")])
	assert.Equal(t, "yes", row[column(t, items, "tracks_stock")])
	deleted := items[2]
	assert.Equal(t, "GONE", deleted[column(t, items, "sku")])
	assert.Equal(t, "yes", deleted[column(t, items, "deleted")], "deleted items are kept so old sales still make sense")

	stock := tables["stock_on_hand.csv"]
	require.Len(t, stock, 3)
	byLocation := map[string]string{}
	for _, r := range stock[1:] {
		byLocation[r[column(t, stock, "location")]] = r[column(t, stock, "quantity")]
	}
	assert.Equal(t, map[string]string{"Main": "8", "Shop 2": "2"}, byLocation, "10.5 received, 2.5 sold from the default location")

	movements := tables["stock_movements.csv"]
	var sold []string
	for _, r := range movements[1:] {
		if r[column(t, movements, "reason")] == "sale" {
			sold = r
		}
	}
	require.NotNil(t, sold)
	assert.Equal(t, "-2.5", sold[column(t, movements, "quantity_change")], "decimal quantities, negative when stock leaves")
	assert.Equal(t, "sam", sold[column(t, movements, "recorded_by")], "who did it is named, not just numbered")

	customers := tables["customers.csv"]
	require.Len(t, customers, 2)
	assert.Equal(t, "Ada Obi", customers[1][column(t, customers, "name")])
	assert.Equal(t, "200000", customers[1][column(t, customers, "owes_minor_units")], "3750.00 sale - 1000.00 paid now - 750.00 repaid = 2000.00 owed")

	sales := tables["sales.csv"]
	require.Len(t, sales, 2)
	assert.Equal(t, "375000", sales[1][column(t, sales, "total_minor_units")])
	assert.Equal(t, "175000", sales[1][column(t, sales, "paid_minor_units")], "1000.00 at the till plus the 750.00 repayment")
	assert.Equal(t, "200000", sales[1][column(t, sales, "owed_minor_units")])
	assert.Equal(t, "Ada Obi", sales[1][column(t, sales, "customer")])
	assert.Equal(t, "sam", sales[1][column(t, sales, "recorded_by")])
	assert.Equal(t, "weekend", sales[1][column(t, sales, "note")])
	assert.Equal(t, "no", sales[1][column(t, sales, "voided")])
	assert.Equal(t, itoa(sale.Sale.SaleId), sales[1][column(t, sales, "sale_id")])

	lines := tables["sale_lines.csv"]
	require.Len(t, lines, 2)
	assert.Equal(t, "2.5", lines[1][column(t, lines, "quantity")])
	assert.Equal(t, "Item RICE", lines[1][column(t, lines, "item_name")], "the fixture names items Item <sku>")

	repayments := tables["repayments.csv"]
	require.Len(t, repayments, 2)
	assert.Equal(t, "75000", repayments[1][column(t, repayments, "amount_minor_units")])
	assert.Equal(t, "Ada Obi", repayments[1][column(t, repayments, "customer")])
	assert.Len(t, tables["repayment_allocations.csv"], 2)

	assert.Equal(t, "Ada Stores", tables["receipt_details.csv"][1][0])
	assert.Equal(t, "Thank you", tables["receipt_details.csv"][1][3])
}

func TestExport_TypedTextCannotRunAsASpreadsheetFormula(t *testing.T) {
	f := newFixture(t)

	_, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: `=HYPERLINK("http://evil.example","click")`, Phone: "+234 800", Note: "@SUM(A1)"})
	require.NoError(t, err)
	_, err = Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "-2+3", Email: "a@b.example"})
	require.NoError(t, err)

	tables, _ := f.export(t, f.ownerUid)
	customers := tables["customers.csv"]

	require.Len(t, customers, 3)
	assert.Equal(t, `'=HYPERLINK("http://evil.example","click")`, customers[1][column(t, customers, "name")])
	assert.Equal(t, "'+234 800", customers[1][column(t, customers, "phone")])
	assert.Equal(t, "'@SUM(A1)", customers[1][column(t, customers, "note")])
	assert.Equal(t, "'-2+3", customers[2][column(t, customers, "name")])
	assert.Equal(t, "a@b.example", customers[2][column(t, customers, "email")], "ordinary text is left alone")
}

func TestExport_OnlyContainsThisBusiness(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")

	f.newItem(t, "MINE", 100, true)
	_, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "My customer"})
	require.NoError(t, err)
	_, err = Items.Add(f.c, other, ItemInput{Sku: "THEIRS", Name: "Theirs", SalePrice: 100})
	require.NoError(t, err)
	_, err = Customers.Add(f.c, other, CustomerInput{Name: "Their customer"})
	require.NoError(t, err)

	tables, raw := f.export(t, f.ownerUid)

	for name, body := range raw {
		assert.NotContains(t, body, "THEIRS", name)
		assert.NotContains(t, body, "Their customer", name)
	}

	assert.Len(t, tables["items.csv"], 2)
	assert.Len(t, tables["customers.csv"], 2)
}

func TestExport_ReadsLargeTablesInPages(t *testing.T) {
	f := newFixture(t)
	original := exportChunkSize
	exportChunkSize = 2
	t.Cleanup(func() { exportChunkSize = original })

	for i := 0; i < 7; i++ {
		_, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Customer " + itoa(int64(i))})
		require.NoError(t, err)
	}

	item := f.newItem(t, "A1", 100, true)
	for i := 0; i < 5; i++ {
		f.receive(t, item.ItemId, 0, 1000)
	}

	tables, _ := f.export(t, f.ownerUid)

	assert.Len(t, tables["customers.csv"], 8, "all 7 customers although they are read 2 at a time")
	assert.Len(t, tables["stock_movements.csv"], 6, "all 5 movements")

	seen := map[string]bool{}
	for _, r := range tables["stock_movements.csv"][1:] {
		assert.False(t, seen[r[0]], "no movement may appear twice across pages")
		seen[r[0]] = true
	}
}

func TestExport_NamesWhoWorkedThereIncludingFormerMembers(t *testing.T) {
	f := newFixture(t)
	f.newUser(t, "sam")
	gone := f.newUser(t, "gone")

	for _, email := range []string{"sam@example.com", "gone@example.com"} {
		m, err := Memberships.Invite(f.c, f.ownerUid, email, extmodels.RoleStaff)
		require.NoError(t, err)
		require.NoError(t, Memberships.Respond(f.c, m.StaffUid, f.ownerUid, true))
	}

	require.NoError(t, Memberships.Remove(f.c, f.ownerUid, gone))
	require.NoError(t, Audit.Record(f.c, &extmodels.AuditLog{OwnerUid: f.ownerUid, ActorUid: gone, Role: extmodels.RoleStaff, Method: "POST", Path: "/x", Status: 200, Action: "sales.add", EntityType: "sales", EntityId: 9}))

	tables, _ := f.export(t, f.ownerUid)

	team := map[string]string{}
	for _, r := range tables["team.csv"][1:] {
		team[r[0]] = r[2]
	}
	assert.Equal(t, map[string]string{"owner": "active", "sam": "active", "gone": "removed"}, team)

	activity := tables["activity_log.csv"]
	require.Len(t, activity, 2)
	assert.Equal(t, "gone", activity[1][column(t, activity, "person")], "a record keeps the name of someone who has left")
	assert.Equal(t, "sales.add", activity[1][column(t, activity, "action")])
	assert.Equal(t, "9", activity[1][column(t, activity, "thing_id")])
}

func TestQuantityFormat(t *testing.T) {
	cases := map[int64]string{0: "0", 1000: "1", 2500: "2.5", 125: "0.125", -750: "-0.75", 10: "0.01", 1005: "1.005", -1: "-0.001", 1234567: "1234.567"}

	for in, want := range cases {
		assert.Equal(t, want, quantity(in), in)
	}
}
