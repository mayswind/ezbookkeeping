package extservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

func (f *fixture) saleInput(lines ...SaleLineInput) SaleInput {
	return SaleInput{
		Lines: lines, PaymentAccountId: f.cashId, ReceivableAccountId: f.receivable, CategoryId: f.incomeCat,
	}
}

func TestSale_CashSaleBooksIncomeAndReducesStock(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1500, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 2 * 1000})
	in.AmountPaid = 3000
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	assert.Equal(t, int64(3000), detail.Sale.Total)
	assert.Equal(t, int64(3000), detail.Sale.Paid)
	assert.Equal(t, int64(0), detail.Sale.Outstanding())
	assert.NotZero(t, detail.Sale.PaidTransactionId)
	assert.Zero(t, detail.Sale.CreditTransactionId)
	require.Len(t, detail.Lines, 1)

	assert.Equal(t, int64(3000), f.accountBalance(t, f.cashId), "cash account receives the money")
	assert.Equal(t, int64(0), f.accountBalance(t, f.receivable))
	assert.Equal(t, int64(8*1000), f.stock(t, item.ItemId, 0))
}

func TestSale_FractionalQuantityRoundsHalfUp(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "KG", 999, true) // 9.99 per kg
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 500}) // 0.5 kg = 499.5 -> 500
	in.AmountPaid = 500
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
	assert.Equal(t, int64(500), detail.Sale.Total)
}

func TestSale_CreditSaleSplitsPaidAndOwedAndTracksCustomer(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 5 * 1000})
	in.CustomerId = customer.CustomerId
	in.AmountPaid = 2000
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	assert.Equal(t, int64(5000), detail.Sale.Total)
	assert.Equal(t, int64(3000), detail.Sale.Outstanding())
	assert.NotZero(t, detail.Sale.PaidTransactionId)
	assert.NotZero(t, detail.Sale.CreditTransactionId)

	assert.Equal(t, int64(2000), f.accountBalance(t, f.cashId))
	assert.Equal(t, int64(3000), f.accountBalance(t, f.receivable), "the unpaid part sits in the receivables account")

	owed, err := Customers.Outstanding(f.c, f.ownerUid, customer.CustomerId)
	require.NoError(t, err)
	assert.Equal(t, int64(3000), owed)

	balances, err := Customers.Balances(f.c, f.ownerUid)
	require.NoError(t, err)
	assert.Equal(t, int64(3000), balances[customer.CustomerId])
}

func TestSale_FullyOnCreditNeedsNoPaymentAccount(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 1000})
	in.PaymentAccountId = 0
	in.CustomerId = customer.CustomerId
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	assert.Zero(t, detail.Sale.PaidTransactionId)
	assert.NotZero(t, detail.Sale.CreditTransactionId)
	assert.Equal(t, int64(1000), f.accountBalance(t, f.receivable))
}

func TestSale_Validation(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	line := SaleLineInput{ItemId: item.ItemId, Qty: 1000}

	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, f.saleInput())
	assert.Equal(t, exterrs.ErrSaleHasNoLines, err)

	in := f.saleInput(line) // nothing paid and no customer: a walk-in cannot buy on credit
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrCreditSaleRequiresCustomer, err)

	in = f.saleInput(line)
	in.AmountPaid = 2000
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrSaleAmountPaidInvalid, err)

	in = f.saleInput(line)
	in.AmountPaid = 1000
	in.Discount = 1001
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrSaleDiscountInvalid, err)

	in = f.saleInput(line)
	in.AmountPaid = 1000
	in.Discount = 1000
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrSaleTotalInvalid, err)

	in = f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 0})
	in.AmountPaid = 1000
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrQuantityInvalid, err)

	in = f.saleInput(line)
	in.AmountPaid = 500
	in.CustomerId = customer.CustomerId
	in.ReceivableAccountId = f.cashId // not a receivables account
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrReceivableAccountInvalid, err)

	in = f.saleInput(line)
	in.AmountPaid = 1000
	in.CategoryId = f.transferCat // not an income category
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrSaleCategoryInvalid, err)

	in = f.saleInput(line)
	in.AmountPaid = 1000
	in.PaymentAccountId = 99999
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrAccountNotUsable, err)

	assert.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0), "no rejected sale may touch stock")
	assert.Equal(t, int64(0), f.accountBalance(t, f.cashId), "no rejected sale may touch the books")
}

func TestSale_InsufficientStockIsRejectedWithoutSideEffects(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 2 * 1000})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 3 * 1000})
	in.AmountPaid = 3000
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrInsufficientStock, err)

	assert.Equal(t, int64(2*1000), f.stock(t, item.ItemId, 0))
	assert.Equal(t, int64(0), f.accountBalance(t, f.cashId))

	sales, err := Sales.List(f.c, f.ownerUid, SaleFilter{})
	require.NoError(t, err)
	assert.Empty(t, sales, "a rejected sale must not leave a row behind")
}

func TestSale_SameItemOnTwoLinesCountsTowardsStockTogether(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 3 * 1000})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 2 * 1000}, SaleLineInput{ItemId: item.ItemId, Qty: 2 * 1000})
	in.AmountPaid = 4000
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrInsufficientStock, err)
}

func TestSale_UnitPriceOverrideAndDiscount(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)

	price := int64(800)
	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 3 * 1000, UnitPrice: &price})
	in.Discount = 400
	in.AmountPaid = 2000
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	assert.Equal(t, int64(2400), detail.Sale.Subtotal)
	assert.Equal(t, int64(2000), detail.Sale.Total)
	assert.Equal(t, int64(800), detail.Lines[0].UnitPrice)
}

func TestSale_UntrackedServiceNeedsNoStock(t *testing.T) {
	f := newFixture(t)
	service := f.newItem(t, "SVC", 5000, false)

	in := f.saleInput(SaleLineInput{ItemId: service.ItemId, Qty: 1000})
	in.AmountPaid = 5000
	_, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
}

func TestSale_SellsFromTheChosenLocation(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	main, err := Locations.Resolve(f.c, f.ownerUid, 0)
	require.NoError(t, err)
	shop2, err := Locations.Add(f.c, f.ownerUid, "Shop 2")
	require.NoError(t, err)

	_, err = Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, LocationId: shop2.LocationId, Qty: 5 * 1000})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 1000})
	in.AmountPaid = 1000
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in) // default location has no stock
	assert.Equal(t, exterrs.ErrInsufficientStock, err)

	in.LocationId = shop2.LocationId
	_, err = Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
	assert.Equal(t, int64(4*1000), f.stock(t, item.ItemId, shop2.LocationId))
	assert.Equal(t, int64(0), f.stock(t, item.ItemId, main.LocationId))
}

func TestSale_VoidRestoresStockAndRemovesMoney(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, true)
	_, err := Stock.Receive(f.c, f.ownerUid, f.ownerUid, StockInput{ItemId: item.ItemId, Qty: 10 * 1000})
	require.NoError(t, err)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 4 * 1000})
	in.CustomerId = customer.CustomerId
	in.AmountPaid = 1000
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
	require.Equal(t, int64(6*1000), f.stock(t, item.ItemId, 0))

	require.NoError(t, Sales.Void(f.c, f.ownerUid, f.ownerUid, detail.Sale.SaleId))

	assert.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0))
	assert.Equal(t, int64(0), f.accountBalance(t, f.cashId))
	assert.Equal(t, int64(0), f.accountBalance(t, f.receivable))

	owed, err := Customers.Outstanding(f.c, f.ownerUid, customer.CustomerId)
	require.NoError(t, err)
	assert.Equal(t, int64(0), owed, "a voided sale is not owed")

	assert.Equal(t, exterrs.ErrSaleAlreadyVoided, Sales.Void(f.c, f.ownerUid, f.ownerUid, detail.Sale.SaleId))
	assert.Equal(t, int64(10*1000), f.stock(t, item.ItemId, 0), "voiding twice must not restock twice")
}

func TestSale_Get_ListAndIsolation(t *testing.T) {
	f := newFixture(t)
	other := f.newUser(t, "other")
	item := f.newItem(t, "A1", 1000, false)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 1000})
	in.AmountPaid = 1000
	first, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)
	second, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	list, err := Sales.List(f.c, f.ownerUid, SaleFilter{})
	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.Equal(t, second.Sale.SaleId, list[0].SaleId, "newest first")

	got, err := Sales.Get(f.c, f.ownerUid, first.Sale.SaleId)
	require.NoError(t, err)
	assert.Len(t, got.Lines, 1)

	_, err = Sales.Get(f.c, other, first.Sale.SaleId)
	assert.Equal(t, exterrs.ErrSaleNotFound, err, "another business must not read this sale")
	assert.Equal(t, exterrs.ErrSaleNotFound, Sales.Void(f.c, other, other, first.Sale.SaleId))
}

func TestSale_StaffActorIsRecorded(t *testing.T) {
	f := newFixture(t)
	staff := f.newUser(t, "staff")
	item := f.newItem(t, "A1", 1000, false)

	in := f.saleInput(SaleLineInput{ItemId: item.ItemId, Qty: 1000})
	in.AmountPaid = 1000
	detail, err := Sales.Create(f.c, f.ownerUid, staff, in)
	require.NoError(t, err)

	assert.Equal(t, f.ownerUid, detail.Sale.OwnerUid)
	assert.Equal(t, staff, detail.Sale.ActorUid)
}

// ---- repayments

func (f *fixture) creditSale(t *testing.T, customerId int64, itemId int64, qty int64, paid int64) *SaleDetail {
	t.Helper()

	in := f.saleInput(SaleLineInput{ItemId: itemId, Qty: qty})
	in.CustomerId = customerId
	in.AmountPaid = paid
	detail, err := Sales.Create(f.c, f.ownerUid, f.ownerUid, in)
	require.NoError(t, err)

	return detail
}

func (f *fixture) repay(customerId int64, amount int64, saleId int64) RepaymentInput {
	return RepaymentInput{
		CustomerId: customerId, Amount: amount, SaleId: saleId, PaymentAccountId: f.bankId,
		ReceivableAccountId: f.receivable, CategoryId: f.transferCat,
	}
}

func TestRepayment_PaysOldestSaleFirstAndMovesMoney(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	older := f.creditSale(t, customer.CustomerId, item.ItemId, 3*1000, 0) // owes 3000
	newer := f.creditSale(t, customer.CustomerId, item.ItemId, 2*1000, 0) // owes 2000
	require.Equal(t, int64(5000), f.accountBalance(t, f.receivable))

	result, err := Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 3500, 0))
	require.NoError(t, err)

	require.Len(t, result.Allocations, 2)
	assert.Equal(t, older.Sale.SaleId, result.Allocations[0].SaleId)
	assert.Equal(t, int64(3000), result.Allocations[0].Amount)
	assert.Equal(t, newer.Sale.SaleId, result.Allocations[1].SaleId)
	assert.Equal(t, int64(500), result.Allocations[1].Amount)
	assert.NotZero(t, result.Repayment.TransactionId)

	assert.Equal(t, int64(1500), f.accountBalance(t, f.receivable), "receivables shrink by the repayment")
	assert.Equal(t, int64(3500), f.accountBalance(t, f.bankId), "the bank account receives the money")

	owed, err := Customers.Outstanding(f.c, f.ownerUid, customer.CustomerId)
	require.NoError(t, err)
	assert.Equal(t, int64(1500), owed)

	gotOlder, err := Sales.Get(f.c, f.ownerUid, older.Sale.SaleId)
	require.NoError(t, err)
	assert.Equal(t, int64(0), gotOlder.Sale.Outstanding())
	gotNewer, err := Sales.Get(f.c, f.ownerUid, newer.Sale.SaleId)
	require.NoError(t, err)
	assert.Equal(t, int64(1500), gotNewer.Sale.Outstanding())
}

func TestRepayment_CanTargetOneSale(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)

	f.creditSale(t, customer.CustomerId, item.ItemId, 3*1000, 0)
	newer := f.creditSale(t, customer.CustomerId, item.ItemId, 2*1000, 0)

	result, err := Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 2000, newer.Sale.SaleId))
	require.NoError(t, err)
	require.Len(t, result.Allocations, 1)
	assert.Equal(t, newer.Sale.SaleId, result.Allocations[0].SaleId)

	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 1, newer.Sale.SaleId))
	assert.Equal(t, exterrs.ErrSaleNotFound, err, "a settled sale is no longer open")
}

func TestRepayment_CannotExceedWhatIsOwed(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	f.creditSale(t, customer.CustomerId, item.ItemId, 1000, 0)

	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 1001, 0))
	assert.Equal(t, exterrs.ErrRepaymentExceedsBalance, err)
	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 0, 0))
	assert.Equal(t, exterrs.ErrRepaymentAmountInvalid, err)

	assert.Equal(t, int64(1000), f.accountBalance(t, f.receivable), "rejected repayments must not touch the books")
	assert.Equal(t, int64(0), f.accountBalance(t, f.bankId))
}

func TestRepayment_FullSettlementClearsTheBalance(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	f.creditSale(t, customer.CustomerId, item.ItemId, 4*1000, 1000)

	for _, amount := range []int64{1000, 2000} {
		_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, amount, 0))
		require.NoError(t, err)
	}

	owed, err := Customers.Outstanding(f.c, f.ownerUid, customer.CustomerId)
	require.NoError(t, err)
	assert.Equal(t, int64(0), owed)
	assert.Equal(t, int64(0), f.accountBalance(t, f.receivable))
	assert.Equal(t, int64(1000), f.accountBalance(t, f.cashId))
	assert.Equal(t, int64(3000), f.accountBalance(t, f.bankId))

	open, err := Sales.List(f.c, f.ownerUid, SaleFilter{OnlyOpen: true})
	require.NoError(t, err)
	assert.Empty(t, open)

	repayments, err := Repayments.List(f.c, f.ownerUid, customer.CustomerId, 0, 0)
	require.NoError(t, err)
	assert.Len(t, repayments, 2)
}

func TestRepayment_SaleWithRepaymentsCannotBeVoided(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	sale := f.creditSale(t, customer.CustomerId, item.ItemId, 2*1000, 0)

	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, f.repay(customer.CustomerId, 500, 0))
	require.NoError(t, err)

	assert.Equal(t, exterrs.ErrSaleHasRepayments, Sales.Void(f.c, f.ownerUid, f.ownerUid, sale.Sale.SaleId))
	assert.Equal(t, int64(1500), f.accountBalance(t, f.receivable), "a refused void must not touch the books")
}

func TestRepayment_CurrencyMismatchIsRejected(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	customer, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Ada"})
	require.NoError(t, err)
	f.creditSale(t, customer.CustomerId, item.ItemId, 1000, 0)

	usd := f.newAccount(t, 2, "USD bank", "USD")
	in := f.repay(customer.CustomerId, 500, 0)
	in.PaymentAccountId = usd
	_, err = Repayments.Add(f.c, f.ownerUid, f.ownerUid, in)
	assert.Equal(t, exterrs.ErrAccountCurrencyMismatch, err)
}

func TestCustomers_DeleteRules(t *testing.T) {
	f := newFixture(t)
	item := f.newItem(t, "A1", 1000, false)
	debtor, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Debtor"})
	require.NoError(t, err)
	clean, err := Customers.Add(f.c, f.ownerUid, CustomerInput{Name: "Clean"})
	require.NoError(t, err)
	f.creditSale(t, debtor.CustomerId, item.ItemId, 1000, 0)

	assert.Equal(t, exterrs.ErrCustomerHasBalance, Customers.Delete(f.c, f.ownerUid, debtor.CustomerId))
	require.NoError(t, Customers.Delete(f.c, f.ownerUid, clean.CustomerId))

	list, err := Customers.List(f.c, f.ownerUid)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "Debtor", list[0].Name)

	_, err = Customers.Add(f.c, f.ownerUid, CustomerInput{Name: " "})
	assert.Equal(t, exterrs.ErrCustomerNameIsEmpty, err)
}

var _ = extmodels.QtyScale
