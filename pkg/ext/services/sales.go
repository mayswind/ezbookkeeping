package extservices

import (
	"fmt"
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// SaleLineInput is one requested line of a sale
type SaleLineInput struct {
	ItemId    int64
	Qty       int64  // scaled by QtyScale
	UnitPrice *int64 // nil means the item's sale price
}

// SaleInput is a request to record a sale
type SaleInput struct {
	LocationId          int64 // 0 means the default location
	CustomerId          int64 // 0 means a walk-in customer, who must pay in full
	Time                int64 // 0 means now
	UtcOffset           int16
	Lines               []SaleLineInput
	Discount            int64
	AmountPaid          int64 // paid now, the rest is put on credit
	PaymentAccountId    int64 // receives AmountPaid
	ReceivableAccountId int64 // receives the credit part, must be a "receivables" account
	CategoryId          int64 // income category used for the sale
	Note                string
}

// SaleDetail is a sale with its lines
type SaleDetail struct {
	Sale  *extmodels.Sale
	Lines []*extmodels.SaleLine
}

// SaleService records sales and keeps stock and the books in step
type SaleService struct{}

// Sales is the sale service singleton
var Sales = &SaleService{}

type plannedLine struct {
	item      *extmodels.Item
	qty       int64
	unitPrice int64
	total     int64
}

// plan validates the request and prices it without writing anything
func (s *SaleService) plan(c core.Context, ownerUid int64, in *SaleInput) (lines []*plannedLine, subtotal int64, err error) {
	if len(in.Lines) == 0 {
		return nil, 0, exterrs.ErrSaleHasNoLines
	}

	items := make(map[int64]*extmodels.Item)

	for _, line := range in.Lines {
		if _, loaded := items[line.ItemId]; loaded {
			continue
		}

		item, err := Items.Get(c, ownerUid, line.ItemId)

		if err != nil {
			return nil, 0, err
		}

		items[line.ItemId] = item
	}

	for _, line := range in.Lines {
		item := items[line.ItemId]
		price := item.SalePrice

		if line.UnitPrice != nil {
			price = *line.UnitPrice
		}

		total, ok := lineTotal(line.Qty, price)

		if !ok {
			return nil, 0, exterrs.ErrQuantityInvalid
		}

		subtotal += total

		if subtotal > maxAmount {
			return nil, 0, exterrs.ErrSaleTotalInvalid
		}

		lines = append(lines, &plannedLine{item: item, qty: line.Qty, unitPrice: price, total: total})
	}

	return lines, subtotal, nil
}

// Create records a sale: stock leaves the location, the paid part is booked as income into the payment account
// and the unpaid part is booked as income into the receivable account and added to what the customer owes.
func (s *SaleService) Create(c core.Context, ownerUid int64, actorUid int64, in SaleInput) (*SaleDetail, error) {
	lines, subtotal, err := s.plan(c, ownerUid, &in)

	if err != nil {
		return nil, err
	}

	if in.Discount < 0 || in.Discount > subtotal {
		return nil, exterrs.ErrSaleDiscountInvalid
	}

	total := subtotal - in.Discount

	if total <= 0 {
		return nil, exterrs.ErrSaleTotalInvalid
	}

	if in.AmountPaid < 0 || in.AmountPaid > total {
		return nil, exterrs.ErrSaleAmountPaidInvalid
	}

	credit := total - in.AmountPaid

	if in.CustomerId > 0 {
		if _, err = Customers.Get(c, ownerUid, in.CustomerId); err != nil {
			return nil, err
		}
	}

	if credit > 0 && in.CustomerId <= 0 {
		return nil, exterrs.ErrCreditSaleRequiresCustomer
	}

	if (in.AmountPaid > 0 && in.PaymentAccountId <= 0) || (credit > 0 && in.ReceivableAccountId <= 0) {
		return nil, exterrs.ErrAccountsRequired
	}

	if err = s.checkAccounts(c, ownerUid, in.AmountPaid > 0, credit > 0, in.PaymentAccountId, in.ReceivableAccountId); err != nil {
		return nil, err
	}

	if err = loadIncomeCategory(c, ownerUid, in.CategoryId); err != nil {
		return nil, err
	}

	location, err := Locations.Resolve(c, ownerUid, in.LocationId)

	if err != nil {
		return nil, err
	}

	if in.Time <= 0 {
		in.Time = nowUnix()
	}

	now := nowUnix()
	sale := &extmodels.Sale{
		OwnerUid: ownerUid, LocationId: location.LocationId, CustomerId: in.CustomerId, SaleTime: in.Time,
		Subtotal: subtotal, Discount: in.Discount, Total: total, Paid: in.AmountPaid,
		PaymentAccountId: in.PaymentAccountId, ReceivableAccountId: in.ReceivableAccountId, CategoryId: in.CategoryId,
		Note: truncate(in.Note, 255), ActorUid: actorUid, CreatedUnixTime: now,
	}
	saleLines := make([]*extmodels.SaleLine, 0, len(lines))

	// (a) stock and sale rows in one database transaction; the stock check is repeated here because it must be race-safe
	err = ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		needed := make(map[int64]int64)

		for _, line := range lines {
			if line.item.TrackStock {
				needed[line.item.ItemId] += line.qty
			}
		}

		// lock the location, then the items in ascending id order (the same order everywhere, so no deadlocks),
		// and only then read stock: concurrent sales of the same item queue up here instead of overselling
		if err := requireLocation(sess, ownerUid, location.LocationId); err != nil {
			return err
		}

		itemIds := make([]int64, 0, len(needed))

		for itemId := range needed {
			itemIds = append(itemIds, itemId)
		}

		sort.Slice(itemIds, func(i, j int) bool { return itemIds[i] < itemIds[j] })

		for _, itemId := range itemIds {
			if _, err := requireStockItem(sess, ownerUid, itemId); err != nil {
				return err
			}

			current, err := onHand(sess, ownerUid, itemId, location.LocationId)

			if err != nil {
				return err
			} else if current < needed[itemId] {
				return exterrs.ErrInsufficientStock
			}
		}

		if _, err := sess.Insert(sale); err != nil {
			return err
		}

		for _, line := range lines {
			saleLine := &extmodels.SaleLine{OwnerUid: ownerUid, SaleId: sale.SaleId, ItemId: line.item.ItemId, Qty: line.qty, UnitPrice: line.unitPrice, LineTotal: line.total}

			if _, err := sess.Insert(saleLine); err != nil {
				return err
			}

			saleLines = append(saleLines, saleLine)

			if !line.item.TrackStock {
				continue
			}

			movement := &extmodels.StockMovement{
				OwnerUid: ownerUid, ItemId: line.item.ItemId, LocationId: location.LocationId, QtyChange: -line.qty,
				Reason: extmodels.StockReasonSale, RefType: "sale", RefId: sale.SaleId, UnitCost: line.item.CostPrice,
				ActorUid: actorUid, MovementTime: in.Time, CreatedUnix: now,
			}

			if _, err := sess.Insert(movement); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// (b) book the money; if this fails the sale rows are removed again
	comment := fmt.Sprintf("Sale #%d", sale.SaleId)

	if sale.Note != "" {
		comment = truncate(comment+": "+sale.Note, 255)
	}

	comment = stamped(c, ownerUid, actorUid, comment)

	creditComment := stamped(c, ownerUid, actorUid, truncate(fmt.Sprintf("Sale #%d (on credit)", sale.SaleId)+noteSuffix(sale.Note), 255))
	var createdIds []int64

	rollback := func(cause error) (*SaleDetail, error) {
		for _, id := range createdIds {
			if delErr := services.Transactions.DeleteTransaction(c, ownerUid, id); delErr != nil {
				log.Errorf(c, "[ext.sales.Create] failed to delete transaction %d while rolling back sale %d, because %s", id, sale.SaleId, delErr.Error())
			}
		}

		if delErr := s.purge(c, ownerUid, sale.SaleId); delErr != nil {
			log.Errorf(c, "[ext.sales.Create] failed to remove sale %d while rolling back, because %s", sale.SaleId, delErr.Error())
		}

		return nil, cause
	}

	if in.AmountPaid > 0 {
		tx := newIncomeTransaction(ownerUid, in.PaymentAccountId, in.CategoryId, in.AmountPaid, in.Time, in.UtcOffset, comment)

		if err = services.Transactions.CreateTransaction(c, tx, nil, nil); err != nil {
			return rollback(err)
		}

		sale.PaidTransactionId = tx.TransactionId
		createdIds = append(createdIds, tx.TransactionId)
	}

	if credit > 0 {
		tx := newIncomeTransaction(ownerUid, in.ReceivableAccountId, in.CategoryId, credit, in.Time, in.UtcOffset, creditComment)

		if err = services.Transactions.CreateTransaction(c, tx, nil, nil); err != nil {
			return rollback(err)
		}

		sale.CreditTransactionId = tx.TransactionId
		createdIds = append(createdIds, tx.TransactionId)
	}

	// (c) remember which transactions belong to the sale; if that cannot be saved, undo everything so the client
	// gets an error for a sale that really did not happen and a retry cannot double-sell
	_, err = ownerDB(ownerUid).NewSession(c).ID(sale.SaleId).Cols("paid_transaction_id", "credit_transaction_id").Update(sale)

	if err != nil {
		log.Errorf(c, "[ext.sales.Create] could not save the transaction ids of sale %d, rolling it back, because %s", sale.SaleId, err.Error())
		return rollback(err)
	}

	return &SaleDetail{Sale: sale, Lines: saleLines}, nil
}

// checkAccounts validates the accounts a sale or repayment will touch
func (s *SaleService) checkAccounts(c core.Context, ownerUid int64, usePayment bool, useReceivable bool, paymentAccountId int64, receivableAccountId int64) error {
	var paymentCurrency, receivableCurrency string

	if usePayment {
		account, err := loadUsableAccount(c, ownerUid, paymentAccountId)

		if err != nil {
			return err
		}

		paymentCurrency = account.Currency
	}

	if useReceivable {
		account, err := loadUsableAccount(c, ownerUid, receivableAccountId)

		if err != nil {
			return err
		}

		if account.Category != models.ACCOUNT_CATEGORY_RECEIVABLES {
			return exterrs.ErrReceivableAccountInvalid
		}

		receivableCurrency = account.Currency
	}

	if usePayment && useReceivable && paymentCurrency != receivableCurrency {
		return exterrs.ErrAccountCurrencyMismatch
	}

	return nil
}

// purge hard-deletes a sale that was never fully booked
func (s *SaleService) purge(c core.Context, ownerUid int64, saleId int64) error {
	return ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		if _, err := sess.Where("owner_uid=? AND ref_type=? AND ref_id=?", ownerUid, "sale", saleId).Delete(&extmodels.StockMovement{}); err != nil {
			return err
		}

		if _, err := sess.Where("owner_uid=? AND sale_id=?", ownerUid, saleId).Delete(&extmodels.SaleLine{}); err != nil {
			return err
		}

		_, err := sess.Where("owner_uid=? AND sale_id=?", ownerUid, saleId).Delete(&extmodels.Sale{})

		return err
	})
}

// Void cancels a sale: its stock returns to a location and its transactions are removed from the books.
// A sale that already received repayments cannot be voided.
//
// The sale is flipped to voided and restocked first, in one database transaction that also checks for repayments.
// Only then are the transactions removed, one at a time, clearing each id on the sale as it goes. If something
// fails halfway the error is returned and calling Void again finishes the job instead of failing forever.
func (s *SaleService) Void(c core.Context, ownerUid int64, actorUid int64, saleId int64) error {
	detail, err := s.Get(c, ownerUid, saleId)

	if err != nil {
		return err
	}

	sale := detail.Sale

	if sale.Voided && sale.PaidTransactionId == 0 && sale.CreditTransactionId == 0 {
		return exterrs.ErrSaleAlreadyVoided
	}

	if !sale.Voided {
		if err = s.voidState(c, ownerUid, actorUid, sale); err != nil {
			return err
		}
	}

	for _, column := range []string{"paid_transaction_id", "credit_transaction_id"} {
		id := sale.PaidTransactionId

		if column == "credit_transaction_id" {
			id = sale.CreditTransactionId
		}

		if id == 0 {
			continue
		}

		// "not found" means an earlier attempt already removed it
		if err = services.Transactions.DeleteTransaction(c, ownerUid, id); err != nil && err != errs.ErrTransactionNotFound {
			log.Errorf(c, "[ext.sales.Void] sale %d is voided but transaction %d could not be removed, call void again to finish, because %s", saleId, id, err.Error())
			return err
		}

		if column == "paid_transaction_id" {
			sale.PaidTransactionId = 0
		} else {
			sale.CreditTransactionId = 0
		}

		if _, err = ownerDB(ownerUid).NewSession(c).ID(saleId).Cols(column).Update(sale); err != nil {
			log.Errorf(c, "[ext.sales.Void] sale %d: transaction %d was removed but its id could not be cleared, because %s", saleId, id, err.Error())
			return err
		}
	}

	return nil
}

// voidState marks the sale voided and returns its stock, atomically
func (s *SaleService) voidState(c core.Context, ownerUid int64, actorUid int64, sale *extmodels.Sale) error {
	// stock goes back to the sale's location, or to the default one if that location was deleted since
	restock, err := Locations.Resolve(c, ownerUid, sale.LocationId)

	if err != nil {
		restock, err = Locations.Resolve(c, ownerUid, 0)

		if err != nil {
			return err
		}
	}

	now := nowUnix()

	return ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		// checked in the same transaction as the update, and a repayment's own update requires voided=false
		hasRepayments, err := sess.Where("owner_uid=? AND sale_id=?", ownerUid, sale.SaleId).Exist(&extmodels.RepaymentAllocation{})

		if err != nil {
			return err
		} else if hasRepayments {
			return exterrs.ErrSaleHasRepayments
		}

		updated, err := sess.Where("owner_uid=? AND sale_id=? AND voided=?", ownerUid, sale.SaleId, false).Cols("voided", "voided_unix_time").Update(&extmodels.Sale{Voided: true, VoidedUnixTime: now})

		if err != nil {
			return err
		} else if updated != 1 {
			return exterrs.ErrConcurrentModification
		}

		original := make([]*extmodels.StockMovement, 0)

		if err = sess.Where("owner_uid=? AND ref_type=? AND ref_id=?", ownerUid, "sale", sale.SaleId).Find(&original); err != nil {
			return err
		}

		if len(original) > 0 {
			if err = requireLocation(sess, ownerUid, restock.LocationId); err != nil {
				return err
			}
		}

		for _, m := range original {
			reverse := &extmodels.StockMovement{
				OwnerUid: ownerUid, ItemId: m.ItemId, LocationId: restock.LocationId, QtyChange: -m.QtyChange, Reason: extmodels.StockReasonSaleVoid,
				RefType: "sale_void", RefId: sale.SaleId, UnitCost: m.UnitCost, ActorUid: actorUid, MovementTime: now, CreatedUnix: now,
			}

			if _, err = sess.Insert(reverse); err != nil {
				return err
			}
		}

		sale.Voided, sale.VoidedUnixTime = true, now

		return nil
	})
}

// Get returns a sale with its lines
func (s *SaleService) Get(c core.Context, ownerUid int64, saleId int64) (*SaleDetail, error) {
	sale := &extmodels.Sale{}
	has, err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND sale_id=?", ownerUid, saleId).Get(sale)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, exterrs.ErrSaleNotFound
	}

	lines := make([]*extmodels.SaleLine, 0)

	if err = ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND sale_id=?", ownerUid, saleId).OrderBy("line_id").Find(&lines); err != nil {
		return nil, err
	}

	return &SaleDetail{Sale: sale, Lines: lines}, nil
}

// SaleFilter filters the sale listing
type SaleFilter struct {
	CustomerId int64
	OnlyOpen   bool // only sales that still have an outstanding balance
	BeforeId   int64
	Limit      int
}

// List returns sales, newest first
func (s *SaleService) List(c core.Context, ownerUid int64, f SaleFilter) ([]*extmodels.Sale, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}

	sales := make([]*extmodels.Sale, 0, f.Limit)
	sess := ownerDB(ownerUid).NewSession(c).Where("owner_uid=?", ownerUid)

	if f.CustomerId > 0 {
		sess = sess.And("customer_id=?", f.CustomerId)
	}

	if f.OnlyOpen {
		sess = sess.And("voided=? AND paid<total", false)
	}

	if f.BeforeId > 0 {
		sess = sess.And("sale_id<?", f.BeforeId)
	}

	err := sess.OrderBy("sale_id desc").Limit(f.Limit).Find(&sales)

	return sales, err
}

// openSalesOldestFirst returns every sale of a customer that still has an outstanding balance, oldest first.
// It is deliberately not paged: repayments must see the customer's whole balance.
func (s *SaleService) openSalesOldestFirst(c core.Context, ownerUid int64, customerId int64) ([]*extmodels.Sale, error) {
	sales := make([]*extmodels.Sale, 0)
	err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND customer_id=? AND voided=? AND paid<total", ownerUid, customerId, false).OrderBy("sale_time, sale_id").Find(&sales)

	return sales, err
}

// noteSuffix is ": note" for a sale note, or nothing
func noteSuffix(note string) string {
	if note == "" {
		return ""
	}

	return ": " + note
}
