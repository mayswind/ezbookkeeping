package extservices

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	"github.com/mayswind/ezbookkeeping/pkg/ext/testdb"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// fixture is a business with the accounts and categories a sale needs
type fixture struct {
	c           core.Context
	ownerUid    int64
	cashId      int64
	bankId      int64
	receivable  int64
	incomeCat   int64
	transferCat int64
}

// newFixture boots an isolated sqlite database with upstream and ext tables
func newFixture(t *testing.T) *fixture {
	t.Helper()

	config := testdb.Config(t)
	require.NoError(t, datastore.InitializeDataStore(config))
	require.NoError(t, uuid.InitializeUuidGenerator(config))

	require.NoError(t, datastore.Container.UserStore.SyncStructs(new(models.User)))
	require.NoError(t, datastore.Container.UserDataStore.SyncStructs(
		new(models.Account), new(models.Transaction), new(models.TransactionCategory),
		new(models.TransactionTagIndex), new(models.TransactionTag), new(models.TransactionPictureInfo)))
	require.NoError(t, datastore.Container.UserStore.SyncStructs(extmodels.GlobalTables()...))
	require.NoError(t, datastore.Container.UserDataStore.SyncStructs(extmodels.OwnerTables()...))

	c := core.NewNullContext()
	f := &fixture{c: c}
	f.ownerUid = f.newUser(t, "owner")

	f.cashId = f.newAccount(t, models.ACCOUNT_CATEGORY_CASH, "Cash", "NGN")
	f.bankId = f.newAccount(t, models.ACCOUNT_CATEGORY_CHECKING_ACCOUNT, "Bank", "NGN")
	f.receivable = f.newAccount(t, models.ACCOUNT_CATEGORY_RECEIVABLES, "Customers owe", "NGN")

	f.incomeCat = f.newCategory(t, models.CATEGORY_TYPE_INCOME, "Sales")
	f.transferCat = f.newCategory(t, models.CATEGORY_TYPE_TRANSFER, "Repayments")

	return f
}

func (f *fixture) newUser(t *testing.T, name string) int64 {
	t.Helper()

	user := &models.User{Username: name, Email: name + "@example.com", Nickname: name, Language: "en", DefaultCurrency: "NGN"}
	require.NoError(t, services.Users.CreateUser(f.c, user, true)) // no password: hashing would dominate the test time

	return user.Uid
}

func (f *fixture) newAccount(t *testing.T, category models.AccountCategory, name string, currency string) int64 {
	t.Helper()

	account := &models.Account{
		AccountId: uuid.Container.GenerateUuid(uuid.UUID_TYPE_ACCOUNT), Uid: f.ownerUid, Category: category,
		Type: models.ACCOUNT_TYPE_SINGLE_ACCOUNT, Name: name, Color: "000000", Currency: currency,
	}
	_, err := datastore.Container.UserDataStore.Choose(f.ownerUid).NewSession(f.c).Insert(account)
	require.NoError(t, err)

	return account.AccountId
}

func (f *fixture) newCategory(t *testing.T, categoryType models.TransactionCategoryType, name string) int64 {
	t.Helper()

	db := datastore.Container.UserDataStore.Choose(f.ownerUid)
	parent := &models.TransactionCategory{CategoryId: uuid.Container.GenerateUuid(uuid.UUID_TYPE_CATEGORY), Uid: f.ownerUid, Type: categoryType, Name: name + " group", Color: "000000"}
	_, err := db.NewSession(f.c).Insert(parent)
	require.NoError(t, err)

	child := &models.TransactionCategory{CategoryId: uuid.Container.GenerateUuid(uuid.UUID_TYPE_CATEGORY), Uid: f.ownerUid, Type: categoryType, ParentCategoryId: parent.CategoryId, Name: name, Color: "000000"}
	_, err = db.NewSession(f.c).Insert(child)
	require.NoError(t, err)

	return child.CategoryId
}

func (f *fixture) accountBalance(t *testing.T, accountId int64) int64 {
	t.Helper()

	account := &models.Account{}
	has, err := datastore.Container.UserDataStore.Choose(f.ownerUid).NewSession(f.c).ID(accountId).Get(account)
	require.NoError(t, err)
	require.True(t, has)

	return account.Balance
}

func (f *fixture) newItem(t *testing.T, sku string, price int64, track bool) *extmodels.Item {
	t.Helper()

	item, err := Items.Add(f.c, f.ownerUid, ItemInput{Sku: sku, Name: "Item " + sku, Unit: "pc", CostPrice: price / 2, SalePrice: price, TrackStock: track})
	require.NoError(t, err)

	return item
}

func (f *fixture) stock(t *testing.T, itemId int64, locationId int64) int64 {
	t.Helper()

	levels, err := Stock.Levels(f.c, f.ownerUid, itemId, locationId)
	require.NoError(t, err)

	var total int64

	for _, level := range levels {
		total += level.Qty
	}

	return total
}
