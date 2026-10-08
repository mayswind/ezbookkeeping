package extservices

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/core"
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

		for itemId, qty := range needed {
			current, err := onHand(sess, ownerUid, itemId, location.LocationId)

			if err != nil {
				return err
			} else if current < qty {
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
		tx := newIncomeTransaction(ownerUid, in.ReceivableAccountId, in.CategoryId, credit, in.Time, in.UtcOffset, comment+" (on credit)")

		if err = services.Transactions.CreateTransaction(c, tx, nil, nil); err != nil {
			return rollback(err)
		}

		sale.CreditTransactionId = tx.TransactionId
		createdIds = append(createdIds, tx.TransactionId)
	}

	// (c) remember which transactions belong to the sale
	_, err = ownerDB(ownerUid).NewSession(c).ID(sale.SaleId).Cols("paid_transaction_id", "credit_transaction_id").Update(sale)

	if err != nil {
		log.Errorf(c, "[ext.sales.Create] sale %d was booked but its transaction ids could not be saved, because %s", sale.SaleId, err.Error())
		return nil, err
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

// Void cancels a sale: its transactions are deleted and its stock returns to the location.
// A sale that already received repayments cannot be voided.
func (s *SaleService) Void(c core.Context, ownerUid int64, actorUid int64, saleId int64) error {
	detail, err := s.Get(c, ownerUid, saleId)

	if err != nil {
		return err
	}

	sale := detail.Sale

	if sale.Voided {
		return exterrs.ErrSaleAlreadyVoided
	}

	hasRepayments, err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND sale_id=?", ownerUid, saleId).Exist(&extmodels.RepaymentAllocation{})

	if err != nil {
		return err
	} else if hasRepayments {
		return exterrs.ErrSaleHasRepayments
	}

	for _, id := range []int64{sale.PaidTransactionId, sale.CreditTransactionId} {
		if id == 0 {
			continue
		}

		if err = services.Transactions.DeleteTransaction(c, ownerUid, id); err != nil {
			return err
		}
	}

	now := nowUnix()

	err = ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		// re-check inside the transaction so two concurrent voids cannot both restock
		updated, err := sess.Where("owner_uid=? AND sale_id=? AND voided=?", ownerUid, saleId, false).Cols("voided", "voided_unix_time").Update(&extmodels.Sale{Voided: true, VoidedUnixTime: now})

		if err != nil {
			return err
		} else if updated != 1 {
			return exterrs.ErrConcurrentModification
		}

		original := make([]*extmodels.StockMovement, 0)

		if err = sess.Where("owner_uid=? AND ref_type=? AND ref_id=?", ownerUid, "sale", saleId).Find(&original); err != nil {
			return err
		}

		for _, m := range original {
			reverse := &extmodels.StockMovement{
				OwnerUid: ownerUid, ItemId: m.ItemId, LocationId: m.LocationId, QtyChange: -m.QtyChange, Reason: extmodels.StockReasonSaleVoid,
				RefType: "sale_void", RefId: saleId, UnitCost: m.UnitCost, ActorUid: actorUid, MovementTime: now, CreatedUnix: now,
			}

			if _, err = sess.Insert(reverse); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf(c, "[ext.sales.Void] sale %d transactions were deleted but its state could not be updated, because %s", saleId, err.Error())
	}

	return err
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
