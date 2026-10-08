package extservices

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// RepaymentInput is a request to record money received from a customer
type RepaymentInput struct {
	CustomerId          int64
	Amount              int64
	SaleId              int64 // 0 means pay off the oldest open sales first
	Time                int64 // 0 means now
	UtcOffset           int16
	PaymentAccountId    int64 // receives the money
	ReceivableAccountId int64 // the receivables account that was debited by the credit sales
	CategoryId          int64 // a transfer category, required by the books for transfers
	Note                string
}

// RepaymentResult is a recorded repayment and how it was spread over sales
type RepaymentResult struct {
	Repayment   *extmodels.Repayment
	Allocations []*extmodels.RepaymentAllocation
}

// RepaymentService records credit repayments
type RepaymentService struct{}

// Repayments is the repayment service singleton
var Repayments = &RepaymentService{}

// Add records a repayment as a transfer from the receivables account to the payment account
// and reduces what the customer owes, oldest sale first unless a sale is given.
func (s *RepaymentService) Add(c core.Context, ownerUid int64, actorUid int64, in RepaymentInput) (*RepaymentResult, error) {
	if in.Amount <= 0 || in.Amount > maxAmount {
		return nil, exterrs.ErrRepaymentAmountInvalid
	}

	if in.PaymentAccountId <= 0 || in.ReceivableAccountId <= 0 {
		return nil, exterrs.ErrAccountsRequired
	}

	if _, err := Customers.Get(c, ownerUid, in.CustomerId); err != nil {
		return nil, err
	}

	if err := Sales.checkAccounts(c, ownerUid, true, true, in.PaymentAccountId, in.ReceivableAccountId); err != nil {
		return nil, err
	}

	openSales, err := Sales.openSalesOldestFirst(c, ownerUid, in.CustomerId)

	if err != nil {
		return nil, err
	}

	if in.SaleId > 0 {
		var chosen *extmodels.Sale

		for _, sale := range openSales {
			if sale.SaleId == in.SaleId {
				chosen = sale
			}
		}

		if chosen == nil {
			return nil, exterrs.ErrSaleNotFound
		}

		openSales = []*extmodels.Sale{chosen}
	}

	var outstanding int64

	for _, sale := range openSales {
		outstanding += sale.Outstanding()
	}

	if in.Amount > outstanding {
		return nil, exterrs.ErrRepaymentExceedsBalance
	}

	allocations := make([]*extmodels.RepaymentAllocation, 0)
	remaining := in.Amount

	for _, sale := range openSales {
		if remaining == 0 {
			break
		}

		portion := sale.Outstanding()

		if portion > remaining {
			portion = remaining
		}

		allocations = append(allocations, &extmodels.RepaymentAllocation{OwnerUid: ownerUid, SaleId: sale.SaleId, Amount: portion})
		remaining -= portion
	}

	if in.Time <= 0 {
		in.Time = nowUnix()
	}

	comment := truncate(fmt.Sprintf("Repayment from customer #%d", in.CustomerId), 255)

	if note := truncate(in.Note, 200); note != "" {
		comment = truncate(comment+": "+note, 255)
	}

	tx := newTransferTransaction(ownerUid, in.ReceivableAccountId, in.PaymentAccountId, in.CategoryId, in.Amount, in.Time, in.UtcOffset, comment)

	if err = services.Transactions.CreateTransaction(c, tx, nil, nil); err != nil {
		return nil, err
	}

	repayment := &extmodels.Repayment{
		OwnerUid: ownerUid, CustomerId: in.CustomerId, Amount: in.Amount, PaymentAccountId: in.PaymentAccountId,
		ReceivableAccountId: in.ReceivableAccountId, TransactionId: tx.TransactionId, RepaymentTime: in.Time,
		Note: truncate(in.Note, 255), ActorUid: actorUid, CreatedUnixTime: nowUnix(),
	}

	err = ownerDB(ownerUid).DoTransaction(c, func(sess *xormSession) error {
		if _, err := sess.Insert(repayment); err != nil {
			return err
		}

		for _, allocation := range allocations {
			allocation.RepaymentId = repayment.RepaymentId

			// optimistic guard: never push a sale past its total, even with concurrent repayments
			result, err := sess.Exec("UPDATE ext_sale SET paid = paid + ? WHERE owner_uid = ? AND sale_id = ? AND voided = ? AND paid + ? <= total",
				allocation.Amount, ownerUid, allocation.SaleId, false, allocation.Amount)

			if err != nil {
				return err
			}

			affected, err := result.RowsAffected()

			if err != nil {
				return err
			} else if affected != 1 {
				return exterrs.ErrConcurrentModification
			}

			if _, err = sess.Insert(allocation); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if delErr := services.Transactions.DeleteTransaction(c, ownerUid, tx.TransactionId); delErr != nil {
			log.Errorf(c, "[ext.repayments.Add] failed to delete transaction %d while rolling back repayment, because %s", tx.TransactionId, delErr.Error())
		}

		return nil, err
	}

	return &RepaymentResult{Repayment: repayment, Allocations: allocations}, nil
}

// List returns the repayments of a business, newest first, optionally for one customer
func (s *RepaymentService) List(c core.Context, ownerUid int64, customerId int64, beforeId int64, limit int) ([]*extmodels.Repayment, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	repayments := make([]*extmodels.Repayment, 0, limit)
	sess := ownerDB(ownerUid).NewSession(c).Where("owner_uid=?", ownerUid)

	if customerId > 0 {
		sess = sess.And("customer_id=?", customerId)
	}

	if beforeId > 0 {
		sess = sess.And("repayment_id<?", beforeId)
	}

	err := sess.OrderBy("repayment_id desc").Limit(limit).Find(&repayments)

	return repayments, err
}
