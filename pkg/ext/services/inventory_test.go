package extservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
)

func TestLocations_DefaultIsCreatedOnFirstUse(t *testing.T) {
	f := newFixture(t)

	locations, err := Locations.List(f.c, f.ownerUid)
	require.NoError(t, err)
	require.Len(t, locations, 1)
	assert.Equal(t, DefaultLocationName, locations[0].Name)
	assert.True(t, locations[0].IsDefault)

	again, err := Locations.List(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Len(t, again, 1, "the default location must not be created twice")
}

func TestLocations_AddRenameDelete(t *testing.T) {
	f := newFixture(t)

	second, err := Locations.Add(f.c, f.ownerUid, "Warehouse")
	require.NoError(t, err)

	_, err = Locations.Add(f.c, f.ownerUid, "warehouse")
	assert.Equal(t, exterrs.ErrLocationNameExists, err)

	_, err = Locations.Add(f.c, f.ownerUid, "  ")
	assert.Equal(t, exterrs.ErrLocationNameIsEmpty, err)

	renamed, err := Locations.Modify(f.c, f.ownerUid, second.LocationId, "Back store")
	require.NoError(t, err)
	assert.Equal(t, "Back store", renamed.Name)

	require.NoError(t, Locations.Delete(f.c, f.ownerUid, second.LocationId))

	remaining, err := Locations.List(f.c, f.ownerUid)
	require.NoError(t, err)
	require.Len(t, remaining, 1)
	assert.Equal(t, exterrs.ErrCannotDeleteLastLocation, Locations.Delete(f.c, f.ownerUid, remaining[0].LocationId))
}

func TestLocations_DeletingTheDefaultPromotesAnother(t *testing.T) {
	f := newFixture(t)

	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	other, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)

	require.NoError(t, Locations.Delete(f.c, f.ownerUid, main.LocationId))

	resolved, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	assert.Equal(t, other.LocationId, resolved.LocationId)
	assert.True(t, resolved.IsDefault)
}

func TestLocations_CannotDeleteLocationWithStock(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	other, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: other.LocationId, Qty: 5 * 1000})
	require.NoError(t, err)

	assert.Equal(t, exterrs.ErrLocationHasStock, Locations.Delete(f.c, f.ownerUid, other.LocationId))
}

func TestItems_SkuMustBeUniquePerBusinessAndIsReusableAfterDelete(t *testing.T) {
	f := newFixture(t)

	item := f.newItem(t, "RICE", 5000, true)

	_, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: "RICE", Name: "Another rice"})
	assert.Equal(t, exterrs.ErrItemSkuExists, err)

	require.NoError(t, Items.Delete(f.c, f.ownerUid, item.ItemId))

	again, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: "RICE", Name: "New rice"})
	require.NoError(t, err)
	assert.NotEqual(t, item.ItemId, again.ItemId)

	_, err = Items.Get(f.c, f.ownerUid, item.ItemId)
	assert.Equal(t, exterrs.ErrItemNotFound, err)
}

func TestItems_Validation(t *testing.T) {
	f := newFixture(t)

	_, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: "", Name: "x"})
	assert.Equal(t, exterrs.ErrItemSkuIsEmpty, err)
	_, err = Items.Add(f.c, f.ownerUid, ItemInput{Sku: "x", Name: ""})
	assert.Equal(t, exterrs.ErrItemNameIsEmpty, err)
	_, err = Items.Add(f.c, f.ownerUid, ItemInput{Sku: "x", Name: "x", SalePrice: -1})
	assert.Equal(t, exterrs.ErrItemPriceInvalid, err)
}

func TestItems_ModifyKeepsSkuUniqueness(t *testing.T) {
	f := newFixture(t)
	a := f.newItem(t, "A", 100, true)
	f.newItem(t, "B", 100, true)

	_, err := Items.Modify(f.c, f.ownerUid, a.ItemId, ItemInput{Sku: "B", Name: "A renamed"})
	assert.Equal(t, exterrs.ErrItemSkuExists, err)

	updated, err := Items.Modify(f.c, f.ownerUid, a.ItemId, ItemInput{Sku: "A", Name: "A renamed", SalePrice: 250, TrackStock: true})
	require.NoError(t, err)
	assert.Equal(t, "A renamed", updated.Name)
	assert.Equal(t, int64(250), updated.SalePrice)
}

func TestStock_ReceiveAdjustAndNeverGoNegative(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)

	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000, UnitCost: 400})
	require.NoError(t, err)
	assert.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0))

	_, err = Stock.Adjust(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: -3 * 1000, Note: "damaged"})
	require.NoError(t, err)
	assert.Equal(t, int64(7*1000), f.stock(t, item.ItemId, 0))

	_, err = Stock.Adjust(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: -8 * 1000})
	assert.Equal(t, exterrs.ErrInsufficientStock, err)
	assert.Equal(t, int64(7*1000), f.stock(t, item.ItemId, 0), "a rejected adjustment must not change stock")

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 0})
	assert.Equal(t, exterrs.ErrQuantityInvalid, err)
	_, err = Stock.Adjust(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 0})
	assert.Equal(t, exterrs.ErrQuantityInvalid, err)
}

func TestStock_FractionalQuantities(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "KG", 1000, true)

	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 2500}) // 2.5 kg
	require.NoError(t, err)
	_, err = Stock.Adjust(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: -750}) // 0.75 kg
	require.NoError(t, err)

	assert.Equal(t, int64(1750), f.stock(t, item.ItemId, 0))
}

func TestStock_UntrackedItemsCannotHaveStock(t *testing.T) {
	f := newFixture(t)
	service := f.newItem(t, "SVC", 5000, false)

	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: service.ItemId, Qty: 1000})
	assert.Equal(t, exterrs.ErrItemDoesNotTrackStock, err)
}

func TestStock_TransferBetweenLocations(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	shop2, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: main.LocationId, Qty: 10 * 1000})
	require.NoError(t, err)

	require.NoError(t, Stock.Transfer(f.c, f.ownerUid, f.ownerUid, item.ItemId, main.LocationId, shop2.LocationId, 4*1000, "restock"))
	assert.Equal(t, int64(6*1000), f.stock(t, item.ItemId, main.LocationId))
	assert.Equal(t, int64(4*1000), f.stock(t, item.ItemId, shop2.LocationId))
	assert.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0), "a transfer must not change the total")

	assert.Equal(t, exterrs.ErrInsufficientStock, Stock.Transfer(f.c, f.ownerUid, f.ownerUid, item.ItemId, shop2.LocationId, main.LocationId, 5*1000, ""))
	assert.Equal(t, exterrs.ErrSameLocation, Stock.Transfer(f.c, f.ownerUid, f.ownerUid, item.ItemId, main.LocationId, main.LocationId, 1000, ""))
	assert.Equal(t, exterrs.ErrQuantityInvalid, Stock.Transfer(f.c, f.ownerUid, f.ownerUid, item.ItemId, main.LocationId, shop2.LocationId, 0, ""))
}

func TestStock_BusinessesAreIsolated(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")
	item := f.newItem(t, "A1", 1000, true)

	_, err := Items.Get(f.c, other, item.ItemId)
	assert.Equal(t, exterrs.ErrItemNotFound, err, "another business must not see this item")

	_, err = Stock.Receive(f.c, other, other, StockInput{ItemId: item.ItemId, Qty: 1000})
	assert.Equal(t, exterrs.ErrItemNotFound, err, "another business must not move this item's stock")
}

func TestStock_MovementsAreListedNewestFirst(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)

	for i := 0; i < 3; i++ {
		_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 1000})
		require.NoError(t, err)
	}

	movements, err := Stock.Movements(f.c, f.ownerUid, MovementFilter{ItemId: item.ItemId, Limit: 2})
	require.NoError(t, err)
	require.Len(t, movements, 2)
	assert.Greater(t, movements[0].MovementId, movements[1].MovementId)

	older, err := Stock.Movements(f.c, f.ownerUid, MovementFilter{ItemId: item.ItemId, BeforeId: movements[1].MovementId})
	require.NoError(t, err)
	assert.Len(t, older, 1)
}
