package services

import (
	"fmt"
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

const pageCountForLoadTransactionAmounts = 1000

// TransactionService represents transaction service
type TransactionService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a transaction service singleton instance
var (
	Transactions = &TransactionService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalTransactionCountByUid returns total transaction count of user
func (s *TransactionService) GetTotalTransactionCountByUid(c *core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.Transaction{})

	return count, err
}

// GetAllTransactions returns all transactions
func (s *TransactionService) GetAllTransactions(c *core.Context, uid int64, pageCount int32, noDuplicated bool) ([]*models.Transaction, error) {
	maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(time.Now().Unix())
	var allTransactions []*models.Transaction

	for maxTransactionTime > 0 {
		transactions, err := s.GetAllTransactionsByMaxTime(c, uid, maxTransactionTime, pageCount, noDuplicated)

		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < int(pageCount) {
			maxTransactionTime = 0
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	return allTransactions, nil
}

// GetAllTransactionsByMaxTime returns all transactions before given time
func (s *TransactionService) GetAllTransactionsByMaxTime(c *core.Context, uid int64, maxTransactionTime int64, count int32, noDuplicated bool) ([]*models.Transaction, error) {
	return s.GetTransactionsByMaxTime(c, uid, maxTransactionTime, 0, 0, nil, nil, "", 1, count, false, noDuplicated)
}

// GetTransactionsByMaxTime returns transactions before given time
func (s *TransactionService) GetTransactionsByMaxTime(c *core.Context, uid int64, maxTransactionTime int64, minTransactionTime int64, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, keyword string, page int32, count int32, needOneMoreItem bool, noDuplicated bool) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if page < 0 {
		return nil, errs.ErrPageIndexInvalid
	} else if page == 0 {
		page = 1
	}

	if count < 1 {
		return nil, errs.ErrPageCountInvalid
	}

	var transactions []*models.Transaction
	var err error

	actualCount := count

	if needOneMoreItem {
		actualCount++
	}

	condition, conditionParams := s.getTransactionQueryCondition(uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, keyword, noDuplicated)
	err = s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).Limit(int(actualCount), int(count*(page-1))).OrderBy("transaction_time desc").Find(&transactions)

	return transactions, err
}

// GetTransactionsInMonthByPage returns all transactions in given year and month
func (s *TransactionService) GetTransactionsInMonthByPage(c *core.Context, uid int64, year int32, month int32, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, keyword string) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	minTransactionTime, maxTransactionTime, err := utils.GetTransactionTimeRangeByYearMonth(year, month)

	if err != nil {
		return nil, errs.ErrSystemError
	}

	var transactions []*models.Transaction

	condition, conditionParams := s.getTransactionQueryCondition(uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, keyword, true)
	err = s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("transaction_time desc").Find(&transactions)

	transactionsInMonth := make([]*models.Transaction, 0, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		transactionUnixTime := utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime)
		transactionTimeZone := time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)

		if utils.IsUnixTimeEqualsYearAndMonth(transactionUnixTime, transactionTimeZone, year, month) {
			transactionsInMonth = append(transactionsInMonth, transaction)
		}
	}

	return transactionsInMonth, err
}

// GetTransactionByTransactionId returns a transaction model according to transaction id
func (s *TransactionService) GetTransactionByTransactionId(c *core.Context, uid int64, transactionId int64) (*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if transactionId <= 0 {
		return nil, errs.ErrTransactionIdInvalid
	}

	transaction := &models.Transaction{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(transaction)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionNotFound
	}

	return transaction, nil
}

// GetAllTransactionCount returns total count of transactions
func (s *TransactionService) GetAllTransactionCount(c *core.Context, uid int64) (int64, error) {
	return s.GetTransactionCount(c, uid, 0, 0, 0, nil, nil, "")
}

// GetMonthTransactionCount returns total count of transactions in given year and month
func (s *TransactionService) GetMonthTransactionCount(c *core.Context, uid int64, year int32, month int32, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, keyword string, utcOffset int16) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	startTime, err := utils.ParseFromLongDateTime(fmt.Sprintf("%d-%02d-01 00:00:00", year, month), utcOffset)

	if err != nil {
		return 0, errs.ErrSystemError
	}

	endTime := startTime.AddDate(0, 1, 0)

	minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startTime.Unix())
	maxTransactionTime := utils.GetMinTransactionTimeFromUnixTime(endTime.Unix()) - 1

	return s.GetTransactionCount(c, uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, keyword)
}

// GetTransactionCount returns count of transactions
func (s *TransactionService) GetTransactionCount(c *core.Context, uid int64, maxTransactionTime int64, minTransactionTime int64, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, keyword string) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	condition, conditionParams := s.getTransactionQueryCondition(uid, maxTransactionTime, minTransactionTime, transactionType, categoryIds, accountIds, keyword, true)
	return s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).Count(&models.Transaction{})
}

// CreateTransaction saves a new transaction to database
func (s *TransactionService) CreateTransaction(c *core.Context, transaction *models.Transaction, tagIds []int64) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	// Check whether account id is valid
	err := s.isAccountIdValid(transaction)

	if err != nil {
		return err
	}

	now := time.Now().Unix()

	needUuidCount := 1

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		needUuidCount = 2
	}

	uuids := s.GenerateUuids(uuid.UUID_TYPE_TRANSACTION, uint8(needUuidCount))

	if len(uuids) < needUuidCount {
		return errs.ErrSystemIsBusy
	}

	transaction.TransactionId = uuids[0]

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		transaction.RelatedId = uuids[1]
	}

	transaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

	transaction.CreatedUnixTime = now
	transaction.UpdatedUnixTime = now

	tagIds = utils.ToUniqueInt64Slice(tagIds)
	transactionTagIndexs := make([]*models.TransactionTagIndex, len(tagIds))

	for i := 0; i < len(tagIds); i++ {
		tagIndexId := s.GenerateUuid(uuid.UUID_TYPE_TAG_INDEX)

		if tagIndexId < 1 {
			return errs.ErrSystemIsBusy
		}

		transactionTagIndexs[i] = &models.TransactionTagIndex{
			TagIndexId:      tagIndexId,
			Uid:             transaction.Uid,
			Deleted:         false,
			TagId:           tagIds[i],
			TransactionId:   transaction.TransactionId,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
	}

	return s.UserDataDB(transaction.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Get and verify source and destination account
		sourceAccount, destinationAccount, err := s.getAccountModels(sess, transaction)

		if err != nil {
			return err
		}

		if sourceAccount.Hidden || (destinationAccount != nil && destinationAccount.Hidden) {
			return errs.ErrCannotAddTransactionToHiddenAccount
		}

		if (transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) &&
			sourceAccount.Currency == destinationAccount.Currency && transaction.Amount != transaction.RelatedAccountAmount {
			return errs.ErrTransactionSourceAndDestinationAmountNotEqual
		}

		// Get and verify category
		err = s.isCategoryValid(sess, transaction)

		if err != nil {
			return err
		}

		// Get and verify tags
		err = s.isTagsValid(sess, transaction, transactionTagIndexs, tagIds)

		if err != nil {
			return err
		}

		// Verify balance modification transaction and calculate real amount
		if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			otherTransactionExists, err := sess.Cols("uid", "deleted", "account_id").Where("uid=? AND deleted=? AND account_id=?", transaction.Uid, false, sourceAccount.AccountId).Limit(1).Exist(&models.Transaction{})

			if err != nil {
				return err
			} else if otherTransactionExists {
				return errs.ErrBalanceModificationTransactionCannotAddWhenNotEmpty
			}

			transaction.RelatedAccountId = transaction.AccountId
			transaction.RelatedAccountAmount = transaction.Amount - sourceAccount.Balance
		}

		// Insert transaction row
		var relatedTransaction *models.Transaction

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			relatedTransaction = s.GetRelatedTransferTransaction(transaction)
		}

		createdRows, err := sess.Insert(transaction)

		if err != nil || createdRows < 1 { // maybe another transaction has same time
			sameSecondLatestTransaction := &models.Transaction{}
			minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
			maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

			has, err := sess.Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<=?", transaction.Uid, false, minTransactionTime, maxTransactionTime).OrderBy("transaction_time desc").Limit(1).Get(sameSecondLatestTransaction)

			if err != nil {
				return err
			} else if !has {
				return errs.ErrDatabaseOperationFailed
			} else if sameSecondLatestTransaction.TransactionTime == maxTransactionTime-1 {
				return errs.ErrTooMuchTransactionInOneSecond
			}

			transaction.TransactionTime = sameSecondLatestTransaction.TransactionTime + 1
			createdRows, err := sess.Insert(transaction)

			if err != nil {
				return err
			} else if createdRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		}

		if relatedTransaction != nil {
			relatedTransaction.TransactionTime = transaction.TransactionTime + 1

			if utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) != utils.GetUnixTimeFromTransactionTime(relatedTransaction.TransactionTime) {
				return errs.ErrTooMuchTransactionInOneSecond
			}

			createdRows, err := sess.Insert(relatedTransaction)

			if err != nil {
				return err
			} else if createdRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		}

		err = nil

		// Insert transaction tag index
		if len(transactionTagIndexs) > 0 {
			for i := 0; i < len(transactionTagIndexs); i++ {
				transactionTagIndex := transactionTagIndexs[i]
				_, err := sess.Insert(transactionTagIndex)

				if err != nil {
					return err
				}
			}
		}

		// Update account table
		if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", transaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedSourceRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", transaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedSourceRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}

			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedDestinationRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedDestinationRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			return errs.ErrTransactionTypeInvalid
		}

		return err
	})
}

// ModifyTransaction saves an existed transaction to database
func (s *TransactionService) ModifyTransaction(c *core.Context, transaction *models.Transaction, addTagIds []int64, removeTagIds []int64) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	updateCols := make([]string, 0, 16)

	now := time.Now().Unix()

	transaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
	transaction.UpdatedUnixTime = now
	updateCols = append(updateCols, "updated_unix_time")

	addTagIds = utils.ToUniqueInt64Slice(addTagIds)
	removeTagIds = utils.ToUniqueInt64Slice(removeTagIds)

	transactionTagIndexs := make([]*models.TransactionTagIndex, len(addTagIds))

	for i := 0; i < len(addTagIds); i++ {
		tagIndexId := s.GenerateUuid(uuid.UUID_TYPE_TAG_INDEX)

		if tagIndexId < 1 {
			return errs.ErrSystemIsBusy
		}

		transactionTagIndexs[i] = &models.TransactionTagIndex{
			TagIndexId:      tagIndexId,
			Uid:             transaction.Uid,
			Deleted:         false,
			TagId:           addTagIds[i],
			TransactionId:   transaction.TransactionId,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
	}

	err := s.UserDataDB(transaction.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Get and verify current transaction
		oldTransaction := &models.Transaction{}
		has, err := sess.ID(transaction.TransactionId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldTransaction)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrTransactionNotFound
		}

		transaction.Type = oldTransaction.Type

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			transaction.RelatedId = oldTransaction.RelatedId
		}

		// Check whether account id is valid
		err = s.isAccountIdValid(transaction)

		if err != nil {
			return err
		}

		// Get and verify source and destination account (if necessary)
		sourceAccount, destinationAccount, err := s.getAccountModels(sess, transaction)

		if err != nil {
			return err
		}

		if sourceAccount.Hidden || (destinationAccount != nil && destinationAccount.Hidden) {
			return errs.ErrCannotModifyTransactionInHiddenAccount
		}

		if (transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) &&
			sourceAccount.Currency == destinationAccount.Currency && transaction.Amount != transaction.RelatedAccountAmount {
			return errs.ErrTransactionSourceAndDestinationAmountNotEqual
		}

		oldSourceAccount, oldDestinationAccount, err := s.getOldAccountModels(sess, transaction, oldTransaction, sourceAccount, destinationAccount)

		if err != nil {
			return err
		}

		if oldSourceAccount.Hidden || (oldDestinationAccount != nil && oldDestinationAccount.Hidden) {
			return errs.ErrCannotAddTransactionToHiddenAccount
		}

		// Append modified columns and verify
		if transaction.CategoryId != oldTransaction.CategoryId {
			// Get and verify category
			err = s.isCategoryValid(sess, transaction)

			if err != nil {
				return err
			}

			updateCols = append(updateCols, "category_id")
		}

		if utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) != utils.GetUnixTimeFromTransactionTime(oldTransaction.TransactionTime) {
			sameSecondLatestTransaction := &models.Transaction{}
			minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
			maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

			has, err = sess.Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<=?", transaction.Uid, false, minTransactionTime, maxTransactionTime).OrderBy("transaction_time desc").Limit(1).Get(sameSecondLatestTransaction)

			if err != nil {
				return err
			}

			if has && sameSecondLatestTransaction.TransactionTime < maxTransactionTime-1 {
				transaction.TransactionTime = sameSecondLatestTransaction.TransactionTime + 1
			} else if has && sameSecondLatestTransaction.TransactionTime == maxTransactionTime-1 {
				return errs.ErrTooMuchTransactionInOneSecond
			}

			updateCols = append(updateCols, "transaction_time")
		}

		if transaction.TimezoneUtcOffset != oldTransaction.TimezoneUtcOffset {
			updateCols = append(updateCols, "timezone_utc_offset")
		}

		if transaction.AccountId != oldTransaction.AccountId {
			updateCols = append(updateCols, "account_id")
		}

		if transaction.Amount != oldTransaction.Amount {
			if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
				originalBalance := sourceAccount.Balance - oldTransaction.RelatedAccountAmount
				transaction.RelatedAccountAmount = transaction.Amount - originalBalance
				updateCols = append(updateCols, "related_account_amount")
			}

			updateCols = append(updateCols, "amount")
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			if transaction.RelatedAccountId != oldTransaction.RelatedAccountId {
				updateCols = append(updateCols, "related_account_id")
			}

			if transaction.RelatedAccountAmount != oldTransaction.RelatedAccountAmount {
				updateCols = append(updateCols, "related_account_amount")
			}
		}

		if transaction.HideAmount != oldTransaction.HideAmount {
			updateCols = append(updateCols, "hide_amount")
		}

		if transaction.Comment != oldTransaction.Comment {
			updateCols = append(updateCols, "comment")
		}

		if transaction.GeoLongitude != oldTransaction.GeoLongitude {
			updateCols = append(updateCols, "geo_longitude")
		}

		if transaction.GeoLatitude != oldTransaction.GeoLatitude {
			updateCols = append(updateCols, "geo_latitude")
		}

		// Get and verify tags
		err = s.isTagsValid(sess, transaction, transactionTagIndexs, addTagIds)

		if err != nil {
			return err
		}

		// Update transaction row
		updatedRows, err := sess.ID(transaction.TransactionId).Cols(updateCols...).Where("uid=? AND deleted=?", transaction.Uid, false).Update(transaction)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionNotFound
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			relatedTransaction := s.GetRelatedTransferTransaction(transaction)

			if utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) != utils.GetUnixTimeFromTransactionTime(relatedTransaction.TransactionTime) {
				return errs.ErrTooMuchTransactionInOneSecond
			}

			relatedUpdateCols := s.getRelatedUpdateColumns(updateCols)
			updatedRows, err := sess.ID(relatedTransaction.TransactionId).Cols(relatedUpdateCols...).Where("uid=? AND deleted=?", relatedTransaction.Uid, false).Update(relatedTransaction)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		}

		// Update transaction tag index
		if len(removeTagIds) > 0 {
			tagIndexUpdateModel := &models.TransactionTagIndex{
				Deleted:         true,
				DeletedUnixTime: now,
			}

			deletedRows, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", transaction.Uid, false, transaction.TransactionId).In("tag_id", removeTagIds).Update(tagIndexUpdateModel)

			if err != nil {
				return err
			} else if deletedRows < 1 {
				return errs.ErrTransactionTagNotFound
			}
		}

		if len(transactionTagIndexs) > 0 {
			for i := 0; i < len(transactionTagIndexs); i++ {
				transactionTagIndex := transactionTagIndexs[i]
				_, err := sess.Insert(transactionTagIndex)

				if err != nil {
					return err
				}
			}
		}

		// Update account table
		if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			if transaction.AccountId != oldTransaction.AccountId {
				return errs.ErrBalanceModificationTransactionCannotChangeAccountId
			}

			if transaction.RelatedAccountAmount != oldTransaction.RelatedAccountAmount {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.RelatedAccountAmount, transaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			var oldAccountNewAmount int64 = 0
			var newAccountNewAmount int64 = 0

			if transaction.AccountId == oldTransaction.AccountId {
				oldAccountNewAmount = transaction.Amount
			} else if transaction.AccountId != oldTransaction.AccountId {
				newAccountNewAmount = transaction.Amount
			}

			if oldAccountNewAmount != oldTransaction.Amount {
				oldSourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldSourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.Amount, oldAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldSourceAccount.Uid, false).Update(oldSourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newAccountNewAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", newAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			var oldAccountNewAmount int64 = 0
			var newAccountNewAmount int64 = 0

			if transaction.AccountId == oldTransaction.AccountId {
				oldAccountNewAmount = transaction.Amount
			} else if transaction.AccountId != oldTransaction.AccountId {
				newAccountNewAmount = transaction.Amount
			}

			if oldAccountNewAmount != oldTransaction.Amount {
				oldSourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldSourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)-(%d)", oldTransaction.Amount, oldAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldSourceAccount.Uid, false).Update(oldSourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newAccountNewAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", newAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			var oldSourceAccountNewAmount int64 = 0
			var newSourceAccountNewAmount int64 = 0

			if transaction.AccountId == oldTransaction.AccountId {
				oldSourceAccountNewAmount = transaction.Amount
			} else if transaction.AccountId != oldTransaction.AccountId {
				newSourceAccountNewAmount = transaction.Amount
			}

			if oldSourceAccountNewAmount != oldTransaction.Amount {
				oldSourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldSourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)-(%d)", oldTransaction.Amount, oldSourceAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldSourceAccount.Uid, false).Update(oldSourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newSourceAccountNewAmount != 0 {
				sourceAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", newSourceAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}

			var oldDestinationAccountNewAmount int64 = 0
			var newDestinationAccountNewAmount int64 = 0

			if transaction.RelatedAccountId == oldTransaction.RelatedAccountId {
				oldDestinationAccountNewAmount = transaction.RelatedAccountAmount
			} else if transaction.RelatedAccountId != oldTransaction.RelatedAccountId {
				newDestinationAccountNewAmount = transaction.RelatedAccountAmount
			}

			if oldDestinationAccountNewAmount != oldTransaction.RelatedAccountAmount {
				oldDestinationAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(oldDestinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.RelatedAccountAmount, oldDestinationAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", oldDestinationAccount.Uid, false).Update(oldDestinationAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}

			if newDestinationAccountNewAmount != 0 {
				destinationAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", newDestinationAccountNewAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			return errs.ErrTransactionTypeInvalid
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// DeleteTransaction deletes an existed transaction from database
func (s *TransactionService) DeleteTransaction(c *core.Context, uid int64, transactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Transaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	tagIndexUpdateModel := &models.TransactionTagIndex{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Get and verify current transaction
		oldTransaction := &models.Transaction{}
		has, err := sess.ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(oldTransaction)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrTransactionNotFound
		}

		// Get and verify source and destination account
		sourceAccount, destinationAccount, err := s.getAccountModels(sess, oldTransaction)

		if err != nil {
			return err
		}

		if sourceAccount.Hidden || (destinationAccount != nil && destinationAccount.Hidden) {
			return errs.ErrCannotDeleteTransactionInHiddenAccount
		}

		// Update transaction row to deleted
		deletedRows, err := sess.ID(oldTransaction.TransactionId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionNotFound
		}

		if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			deletedRows, err = sess.ID(oldTransaction.RelatedId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

			if err != nil {
				return err
			} else if deletedRows < 1 {
				return errs.ErrTransactionNotFound
			}
		}

		// Update transaction tag index
		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", uid, false, oldTransaction.TransactionId).Update(tagIndexUpdateModel)

		if err != nil {
			return err
		}

		// Update account table
		if oldTransaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", oldTransaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedSourceRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", oldTransaction.Amount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedSourceRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}

			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedDestinationRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.RelatedAccountAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedDestinationRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			return errs.ErrTransactionTypeInvalid
		}

		return err
	})
}

// DeleteAllTransactions deletes all existed transactions from database
func (s *TransactionService) DeleteAllTransactions(c *core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Transaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	tagIndexUpdateModel := &models.TransactionTagIndex{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	accountUpdateModel := &models.Account{
		Balance:         0,
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		// Update all transaction to deleted
		_, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		// Update all transaction tag index to deleted
		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(tagIndexUpdateModel)

		if err != nil {
			return err
		}

		// Update all account table to deleted
		_, err = sess.Cols("balance", "deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(accountUpdateModel)

		if err != nil {
			return err
		}

		return nil
	})
}

// GetRelatedTransferTransaction returns the related transaction for transfer transaction
func (s *TransactionService) GetRelatedTransferTransaction(originalTransaction *models.Transaction) *models.Transaction {
	var relatedType models.TransactionDbType
	var relatedTransactionTime int64

	if originalTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		relatedType = models.TRANSACTION_DB_TYPE_TRANSFER_IN
		relatedTransactionTime = originalTransaction.TransactionTime + 1
	} else if originalTransaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		relatedType = models.TRANSACTION_DB_TYPE_TRANSFER_OUT
		relatedTransactionTime = originalTransaction.TransactionTime - 1
	} else {
		return nil
	}

	relatedTransaction := &models.Transaction{
		TransactionId:        originalTransaction.RelatedId,
		Uid:                  originalTransaction.Uid,
		Deleted:              originalTransaction.Deleted,
		Type:                 relatedType,
		CategoryId:           originalTransaction.CategoryId,
		TransactionTime:      relatedTransactionTime,
		TimezoneUtcOffset:    originalTransaction.TimezoneUtcOffset,
		AccountId:            originalTransaction.RelatedAccountId,
		Amount:               originalTransaction.RelatedAccountAmount,
		RelatedId:            originalTransaction.TransactionId,
		RelatedAccountId:     originalTransaction.AccountId,
		RelatedAccountAmount: originalTransaction.Amount,
		Comment:              originalTransaction.Comment,
		GeoLongitude:         originalTransaction.GeoLongitude,
		GeoLatitude:          originalTransaction.GeoLatitude,
		CreatedIp:            originalTransaction.CreatedIp,
		CreatedUnixTime:      originalTransaction.CreatedUnixTime,
		UpdatedUnixTime:      originalTransaction.UpdatedUnixTime,
		DeletedUnixTime:      originalTransaction.DeletedUnixTime,
	}

	return relatedTransaction
}

// GetAccountsTotalIncomeAndExpense returns the every accounts total income and expense amount by specific date range
func (s *TransactionService) GetAccountsTotalIncomeAndExpense(c *core.Context, uid int64, startUnixTime int64, endUnixTime int64, utcOffset int16, useTransactionTimezone bool) (map[int64]int64, map[int64]int64, error) {
	if uid <= 0 {
		return nil, nil, errs.ErrUserIdInvalid
	}

	clientLocation := time.FixedZone("Client Timezone", int(utcOffset)*60)
	startLocalDateTime := utils.FormatUnixTimeToNumericLocalDateTime(startUnixTime, clientLocation)
	endLocalDateTime := utils.FormatUnixTimeToNumericLocalDateTime(endUnixTime, clientLocation)

	startUnixTime = utils.GetMinUnixTimeWithSameLocalDateTime(startUnixTime, utcOffset)
	endUnixTime = utils.GetMaxUnixTimeWithSameLocalDateTime(endUnixTime, utcOffset)

	startTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startUnixTime)
	endTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(endUnixTime)

	condition := "uid=? AND deleted=? AND (type=? OR type=?) AND transaction_time>=? AND transaction_time<=?"
	conditionParams := make([]any, 0, 4)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)

	minTransactionTime := startTransactionTime
	maxTransactionTime := endTransactionTime
	var allTransactions []*models.Transaction

	for maxTransactionTime > 0 {
		var transactions []*models.Transaction

		finalConditionParams := make([]any, 0, 6)
		finalConditionParams = append(finalConditionParams, conditionParams...)
		finalConditionParams = append(finalConditionParams, minTransactionTime)
		finalConditionParams = append(finalConditionParams, maxTransactionTime)

		err := s.UserDataDB(uid).NewSession(c).Select("type, account_id, transaction_time, timezone_utc_offset, amount").Where(condition, finalConditionParams...).Limit(pageCountForLoadTransactionAmounts, 0).OrderBy("transaction_time desc").Find(&transactions)

		if err != nil {
			return nil, nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < pageCountForLoadTransactionAmounts {
			maxTransactionTime = 0
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	incomeAmounts := make(map[int64]int64)
	expenseAmounts := make(map[int64]int64)

	for i := 0; i < len(allTransactions); i++ {
		transaction := allTransactions[i]
		timeZone := clientLocation

		if useTransactionTimezone {
			timeZone = time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		}

		localDateTime := utils.FormatUnixTimeToNumericLocalDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), timeZone)

		if localDateTime < startLocalDateTime || localDateTime > endLocalDateTime {
			continue
		}

		var amountsMap map[int64]int64

		if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME {
			amountsMap = incomeAmounts
		} else if transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			amountsMap = expenseAmounts
		}

		totalAmounts, exists := amountsMap[transaction.AccountId]

		if !exists {
			totalAmounts = 0
		}

		totalAmounts += transaction.Amount
		amountsMap[transaction.AccountId] = totalAmounts
	}

	return incomeAmounts, expenseAmounts, nil
}

// GetAccountsAndCategoriesTotalIncomeAndExpense returns the every accounts and categories total income and expense amount by specific date range
func (s *TransactionService) GetAccountsAndCategoriesTotalIncomeAndExpense(c *core.Context, uid int64, startUnixTime int64, endUnixTime int64, utcOffset int16, useTransactionTimezone bool) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	clientLocation := time.FixedZone("Client Timezone", int(utcOffset)*60)
	startLocalDateTime := utils.FormatUnixTimeToNumericLocalDateTime(startUnixTime, clientLocation)
	endLocalDateTime := utils.FormatUnixTimeToNumericLocalDateTime(endUnixTime, clientLocation)

	startUnixTime = utils.GetMinUnixTimeWithSameLocalDateTime(startUnixTime, utcOffset)
	endUnixTime = utils.GetMaxUnixTimeWithSameLocalDateTime(endUnixTime, utcOffset)

	startTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startUnixTime)
	endTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(endUnixTime)

	condition := "uid=? AND deleted=? AND (type=? OR type=?)"
	conditionParams := make([]any, 0, 4)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
	conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)

	if startUnixTime > 0 {
		condition = condition + " AND transaction_time>=?"
	}

	if endUnixTime > 0 {
		condition = condition + " AND transaction_time<=?"
	}

	minTransactionTime := startTransactionTime
	maxTransactionTime := endTransactionTime
	var allTransactions []*models.Transaction

	for maxTransactionTime > 0 {
		var transactions []*models.Transaction

		finalConditionParams := make([]any, 0, 6)
		finalConditionParams = append(finalConditionParams, conditionParams...)

		if startUnixTime > 0 {
			finalConditionParams = append(finalConditionParams, minTransactionTime)
		}

		if endUnixTime > 0 {
			finalConditionParams = append(finalConditionParams, maxTransactionTime)
		}

		err := s.UserDataDB(uid).NewSession(c).Select("category_id, account_id, transaction_time, timezone_utc_offset, amount").Where(condition, finalConditionParams...).Limit(pageCountForLoadTransactionAmounts, 0).OrderBy("transaction_time desc").Find(&transactions)

		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < pageCountForLoadTransactionAmounts {
			maxTransactionTime = 0
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	transactionTotalAmountsMap := make(map[string]*models.Transaction)

	for i := 0; i < len(allTransactions); i++ {
		transaction := allTransactions[i]
		timeZone := clientLocation

		if useTransactionTimezone {
			timeZone = time.FixedZone("Transaction Timezone", int(transaction.TimezoneUtcOffset)*60)
		}

		localDateTime := utils.FormatUnixTimeToNumericLocalDateTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime), timeZone)

		if localDateTime < startLocalDateTime || localDateTime > endLocalDateTime {
			continue
		}

		groupKey := fmt.Sprintf("%d_%d", transaction.CategoryId, transaction.AccountId)
		totalAmounts, exists := transactionTotalAmountsMap[groupKey]

		if !exists {
			totalAmounts = &models.Transaction{
				CategoryId: transaction.CategoryId,
				AccountId:  transaction.AccountId,
				Amount:     0,
			}

			transactionTotalAmountsMap[groupKey] = totalAmounts
		}

		totalAmounts.Amount += transaction.Amount
	}

	transactionTotalAmounts := make([]*models.Transaction, 0, len(transactionTotalAmountsMap))

	for _, totalAmounts := range transactionTotalAmountsMap {
		transactionTotalAmounts = append(transactionTotalAmounts, totalAmounts)
	}

	return transactionTotalAmounts, nil
}

// GetTransactionMapByList returns a transaction map by a list
func (s *TransactionService) GetTransactionMapByList(transactions []*models.Transaction) map[int64]*models.Transaction {
	transactionMap := make(map[int64]*models.Transaction)

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		transactionMap[transaction.TransactionId] = transaction
	}

	return transactionMap
}

func (s *TransactionService) getTransactionQueryCondition(uid int64, maxTransactionTime int64, minTransactionTime int64, transactionType models.TransactionDbType, categoryIds []int64, accountIds []int64, keyword string, noDuplicated bool) (string, []any) {
	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, 16)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	if maxTransactionTime > 0 {
		condition = condition + " AND transaction_time<=?"
		conditionParams = append(conditionParams, maxTransactionTime)
	}

	if minTransactionTime > 0 {
		condition = condition + " AND transaction_time>=?"
		conditionParams = append(conditionParams, minTransactionTime)
	}

	if models.TRANSACTION_DB_TYPE_MODIFY_BALANCE <= transactionType && transactionType <= models.TRANSACTION_DB_TYPE_EXPENSE {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, transactionType)
	} else if transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		if len(accountIds) == 0 {
			condition = condition + " AND type=?"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		} else {
			condition = condition + " AND (type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_IN)
		}
	} else {
		if noDuplicated && len(accountIds) == 0 {
			condition = condition + " AND (type=? OR type=? OR type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		}
	}

	if len(categoryIds) > 0 {
		var conditions strings.Builder

		for i := 0; i < len(categoryIds); i++ {
			if i > 0 {
				conditions.WriteString(",")
			}

			conditions.WriteString("?")
			conditionParams = append(conditionParams, categoryIds[i])
		}

		condition = condition + " AND category_id IN (" + conditions.String() + ")"
	}

	if len(accountIds) > 0 {
		var conditions strings.Builder

		for i := 0; i < len(accountIds); i++ {
			if i > 0 {
				conditions.WriteString(",")
			}

			conditions.WriteString("?")
			conditionParams = append(conditionParams, accountIds[i])
		}

		condition = condition + " AND account_id IN (" + conditions.String() + ")"
	}

	if keyword != "" {
		condition = condition + " AND comment LIKE ?"
		conditionParams = append(conditionParams, "%%"+keyword+"%%")
	}

	return condition, conditionParams
}

func (s *TransactionService) isAccountIdValid(transaction *models.Transaction) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.RelatedAccountId != 0 && transaction.RelatedAccountId != transaction.AccountId {
			return errs.ErrTransactionDestinationAccountCannotBeSet
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME ||
		transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
		if transaction.RelatedAccountId != 0 {
			return errs.ErrTransactionDestinationAccountCannotBeSet
		} else if transaction.RelatedAccountAmount != 0 {
			return errs.ErrTransactionDestinationAmountCannotBeSet
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if transaction.AccountId == transaction.RelatedAccountId {
			return errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return errs.ErrTransactionTypeInvalid
	} else {
		return errs.ErrTransactionTypeInvalid
	}

	return nil
}

func (s *TransactionService) getAccountModels(sess *xorm.Session, transaction *models.Transaction) (sourceAccount *models.Account, destinationAccount *models.Account, err error) {
	sourceAccount = &models.Account{}
	destinationAccount = &models.Account{}

	has, err := sess.ID(transaction.AccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(sourceAccount)

	if err != nil {
		return nil, nil, err
	} else if !has {
		return nil, nil, errs.ErrSourceAccountNotFound
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.RelatedAccountId != 0 && transaction.RelatedAccountId != transaction.AccountId {
			return nil, nil, errs.ErrAccountIdInvalid
		} else {
			destinationAccount = sourceAccount
		}
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_INCOME || transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
		if transaction.RelatedAccountId != 0 {
			return nil, nil, errs.ErrAccountIdInvalid
		}

		destinationAccount = nil
	} else if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		if transaction.RelatedAccountId <= 0 {
			return nil, nil, errs.ErrAccountIdInvalid
		} else {
			has, err = sess.ID(transaction.RelatedAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(destinationAccount)

			if err != nil {
				return nil, nil, err
			} else if !has {
				return nil, nil, errs.ErrDestinationAccountNotFound
			}
		}
	}

	return sourceAccount, destinationAccount, nil
}

func (s *TransactionService) getOldAccountModels(sess *xorm.Session, transaction *models.Transaction, oldTransaction *models.Transaction, sourceAccount *models.Account, destinationAccount *models.Account) (oldSourceAccount *models.Account, oldDestinationAccount *models.Account, err error) {
	oldSourceAccount = &models.Account{}
	oldDestinationAccount = &models.Account{}

	if transaction.AccountId == oldTransaction.AccountId {
		oldSourceAccount = sourceAccount
	} else {
		has, err := sess.ID(oldTransaction.AccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldSourceAccount)

		if err != nil {
			return nil, nil, err
		} else if !has {
			return nil, nil, errs.ErrSourceAccountNotFound
		}
	}

	if transaction.RelatedAccountId == oldTransaction.RelatedAccountId {
		oldDestinationAccount = destinationAccount
	} else {
		has, err := sess.ID(oldTransaction.RelatedAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldDestinationAccount)

		if err != nil {
			return nil, nil, err
		} else if !has {
			return nil, nil, errs.ErrDestinationAccountNotFound
		}
	}
	return oldSourceAccount, oldDestinationAccount, nil
}

func (s *TransactionService) getRelatedUpdateColumns(updateCols []string) []string {
	relatedUpdateCols := make([]string, len(updateCols))

	for i := 0; i < len(updateCols); i++ {
		if updateCols[i] == "account_id" {
			relatedUpdateCols[i] = "related_account_id"
		} else if updateCols[i] == "related_account_id" {
			relatedUpdateCols[i] = "account_id"
		} else if updateCols[i] == "amount" {
			relatedUpdateCols[i] = "related_account_amount"
		} else if updateCols[i] == "related_account_amount" {
			relatedUpdateCols[i] = "amount"
		} else {
			relatedUpdateCols[i] = updateCols[i]
		}
	}

	return relatedUpdateCols
}

func (s *TransactionService) isCategoryValid(sess *xorm.Session, transaction *models.Transaction) error {
	if transaction.Type == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		if transaction.CategoryId != 0 {
			return errs.ErrBalanceModificationTransactionCannotSetCategory
		}
	} else {
		category := &models.TransactionCategory{}
		has, err := sess.ID(transaction.CategoryId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(category)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrTransactionCategoryNotFound
		}

		if category.ParentCategoryId < 1 {
			return errs.ErrCannotUsePrimaryCategoryForTransaction
		}

		if (transaction.Type == models.TRANSACTION_DB_TYPE_INCOME && category.Type != models.CATEGORY_TYPE_INCOME) ||
			(transaction.Type == models.TRANSACTION_DB_TYPE_EXPENSE && category.Type != models.CATEGORY_TYPE_EXPENSE) ||
			((transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN) && category.Type != models.CATEGORY_TYPE_TRANSFER) {
			return errs.ErrTransactionCategoryTypeInvalid
		}
	}

	return nil
}

func (s *TransactionService) isTagsValid(sess *xorm.Session, transaction *models.Transaction, transactionTagIndexs []*models.TransactionTagIndex, tagIds []int64) error {
	if len(transactionTagIndexs) > 0 {
		var tags []*models.TransactionTag
		err := sess.Where("uid=? AND deleted=?", transaction.Uid, false).In("tag_id", tagIds).Find(&tags)

		if err != nil {
			return err
		}

		tagMap := make(map[int64]*models.TransactionTag)

		for i := 0; i < len(tags); i++ {
			tagMap[tags[i].TagId] = tags[i]
		}

		for i := 0; i < len(transactionTagIndexs); i++ {
			if _, exists := tagMap[transactionTagIndexs[i].TagId]; !exists {
				return errs.ErrTransactionTagNotFound
			}
		}
	}

	return nil
}
