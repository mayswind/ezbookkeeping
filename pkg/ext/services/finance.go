package extservices

import (
	"math/big"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// Limits that keep integer arithmetic far away from overflow
const (
	maxQty       int64 = 1_000_000_000_000     // one billion units
	maxUnitPrice int64 = 1_000_000_000_000_0   // generous upper bound in minor units
	maxAmount    int64 = 1_000_000_000_000_000 // 10^15 minor units

	// maxStockOnHand caps stock of one item at one location, so sums of movements can never overflow an int64
	maxStockOnHand int64 = 1_000_000_000_000_000
)

// lineTotal returns qty (scaled by QtyScale) times unit price per whole unit, rounded half up
func lineTotal(qty int64, unitPrice int64) (int64, bool) {
	if qty <= 0 || qty > maxQty || unitPrice < 0 || unitPrice > maxUnitPrice {
		return 0, false
	}

	product := new(big.Int).Mul(big.NewInt(qty), big.NewInt(unitPrice))
	product.Add(product, big.NewInt(500))
	product.Div(product, big.NewInt(1000))

	if !product.IsInt64() || product.Int64() > maxAmount {
		return 0, false
	}

	return product.Int64(), true
}

// loadUsableAccount returns an account that can receive transactions
func loadUsableAccount(c core.Context, ownerUid int64, accountId int64) (*models.Account, error) {
	account := &models.Account{}
	has, err := ownerDB(ownerUid).NewSession(c).Where("uid=? AND account_id=? AND deleted=?", ownerUid, accountId, false).Get(account)

	if err != nil {
		return nil, err
	} else if !has || account.Hidden || account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		return nil, exterrs.ErrAccountNotUsable
	}

	return account, nil
}

// loadIncomeCategory checks that the category exists and is an income category
func loadIncomeCategory(c core.Context, ownerUid int64, categoryId int64) error {
	category := &models.TransactionCategory{}
	has, err := ownerDB(ownerUid).NewSession(c).Where("uid=? AND category_id=? AND deleted=?", ownerUid, categoryId, false).Get(category)

	if err != nil {
		return err
	} else if !has || category.Type != models.CATEGORY_TYPE_INCOME {
		return exterrs.ErrSaleCategoryInvalid
	}

	return nil
}

func newIncomeTransaction(ownerUid int64, accountId int64, categoryId int64, amount int64, unixTime int64, utcOffset int16, comment string) *models.Transaction {
	return &models.Transaction{
		Uid:               ownerUid,
		Type:              models.TRANSACTION_DB_TYPE_INCOME,
		CategoryId:        categoryId,
		AccountId:         accountId,
		Amount:            amount,
		TransactionTime:   utils.GetMinTransactionTimeFromUnixTime(unixTime),
		TimezoneUtcOffset: utcOffset,
		Comment:           comment,
	}
}

func newTransferTransaction(ownerUid int64, fromAccountId int64, toAccountId int64, categoryId int64, amount int64, unixTime int64, utcOffset int16, comment string) *models.Transaction {
	return &models.Transaction{
		Uid:                  ownerUid,
		Type:                 models.TRANSACTION_DB_TYPE_TRANSFER_OUT,
		CategoryId:           categoryId,
		AccountId:            fromAccountId,
		RelatedAccountId:     toAccountId,
		Amount:               amount,
		RelatedAccountAmount: amount,
		TransactionTime:      utils.GetMinTransactionTimeFromUnixTime(unixTime),
		TimezoneUtcOffset:    utcOffset,
		Comment:              comment,
	}
}
