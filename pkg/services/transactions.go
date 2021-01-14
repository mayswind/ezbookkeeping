package services

import (
	"fmt"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/utils"
	"github.com/mayswind/lab/pkg/uuid"
)

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

// GetAllTransactionsByMaxTime returns all transactions before given time
func (s *TransactionService) GetAllTransactionsByMaxTime(uid int64, maxTime int64, count int, noDuplicated bool) ([]*models.Transaction, error) {
	return s.GetTransactionsByMaxTime(uid, maxTime, 0, 0, 0, 0, "", count, noDuplicated)
}

// GetTransactionsByMaxTime returns transactions before given time
func (s *TransactionService) GetTransactionsByMaxTime(uid int64, maxTime int64, minTime int64, transactionType models.TransactionDbType, categoryId int64, accountId int64, keyword string, count int, noDuplicated bool) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if count < 1 {
		return nil, errs.ErrPageCountInvalid
	}

	var transactions []*models.Transaction
	var err error

	condition := "uid=? AND deleted=?"
	conditionParams := make([]interface{}, 0, 16)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	if models.TRANSACTION_DB_TYPE_MODIFY_BALANCE <= transactionType && transactionType <= models.TRANSACTION_DB_TYPE_EXPENSE {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, transactionType)
	} else if transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		if accountId == 0 {
			condition = condition + " AND type=?"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		} else {
			condition = condition + " AND (type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_IN)
		}
	} else {
		if noDuplicated && accountId == 0 {
			condition = condition + " AND (type=? OR type=? OR type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		}
	}

	if categoryId > 0 {
		condition = condition + " AND category_id=?"
		conditionParams = append(conditionParams, categoryId)
	}

	if accountId > 0 {
		condition = condition + " AND account_id=?"
		conditionParams = append(conditionParams, accountId)
	}

	if keyword != "" {
		condition = condition + " AND comment LIKE ?"
		conditionParams = append(conditionParams, "%%"+keyword+"%%")
	}

	if maxTime > 0 {
		condition = condition + " AND transaction_time<=?"
		conditionParams = append(conditionParams, maxTime)
	}

	if minTime > 0 {
		condition = condition + " AND transaction_time>=?"
		conditionParams = append(conditionParams, minTime)
	}

	err = s.UserDataDB(uid).Where(condition, conditionParams...).Limit(count, 0).OrderBy("transaction_time desc").Find(&transactions)

	return transactions, err
}

// GetTransactionsInMonthByPage returns transactions in given year and month
func (s *TransactionService) GetTransactionsInMonthByPage(uid int64, year int, month int, transactionType models.TransactionDbType, categoryId int64, accountId int64, keyword string, page int, count int) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if page < 1 {
		return nil, errs.ErrPageIndexInvalid
	}

	if count < 1 {
		return nil, errs.ErrPageCountInvalid
	}

	startTime, err := utils.ParseFromLongDateTime(fmt.Sprintf("%d-%d-01 00:00:00", year, month))

	if err != nil {
		return nil, errs.ErrSystemError
	}

	endTime := startTime.AddDate(0, 1, 0)

	startTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startTime.Unix())
	endTransactionTime := utils.GetMinTransactionTimeFromUnixTime(endTime.Unix())

	var transactions []*models.Transaction

	condition := "uid=? AND deleted=? AND transaction_time>=? AND transaction_time<?"
	conditionParams := make([]interface{}, 0, 16)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)
	conditionParams = append(conditionParams, startTransactionTime)
	conditionParams = append(conditionParams, endTransactionTime)

	if models.TRANSACTION_DB_TYPE_MODIFY_BALANCE <= transactionType && transactionType <= models.TRANSACTION_DB_TYPE_EXPENSE {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, transactionType)
	} else if transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transactionType == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		if accountId == 0 {
			condition = condition + " AND type=?"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		} else {
			condition = condition + " AND (type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_IN)
		}
	} else {
		if accountId == 0 {
			condition = condition + " AND (type=? OR type=? OR type=? OR type=?)"
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_MODIFY_BALANCE)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_INCOME)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_EXPENSE)
			conditionParams = append(conditionParams, models.TRANSACTION_DB_TYPE_TRANSFER_OUT)
		}
	}

	if categoryId > 0 {
		condition = condition + " AND category_id=?"
		conditionParams = append(conditionParams, categoryId)
	}

	if accountId > 0 {
		condition = condition + " AND account_id=?"
		conditionParams = append(conditionParams, accountId)
	}

	if keyword != "" {
		condition = condition + " AND comment LIKE ?"
		conditionParams = append(conditionParams, "%%"+keyword+"%%")
	}

	err = s.UserDataDB(uid).Where(condition, conditionParams...).Limit(count, count*(page-1)).OrderBy("transaction_time desc").Find(&transactions)

	return transactions, err
}

// GetTransactionByTransactionId returns a transaction model according to transaction id
func (s *TransactionService) GetTransactionByTransactionId(uid int64, transactionId int64) (*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if transactionId <= 0 {
		return nil, errs.ErrTransactionIdInvalid
	}

	transaction := &models.Transaction{}
	has, err := s.UserDataDB(uid).ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(transaction)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionNotFound
	}

	return transaction, nil
}

// GetAllTransactionCount returns total count of transactions
func (s *TransactionService) GetAllTransactionCount(uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	return s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).Count(&models.Transaction{})
}

// GetMonthTransactionCount returns total count of transactions in given year and month
func (s *TransactionService) GetMonthTransactionCount(uid int64, year int64, month int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	startTime, err := utils.ParseFromLongDateTime(fmt.Sprintf("%d-%d-01 00:00:00", year, month))

	if err != nil {
		return 0, errs.ErrSystemError
	}

	endTime := startTime.AddDate(0, 1, 0)

	startTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startTime.Unix())
	endTransactionTime := utils.GetMinTransactionTimeFromUnixTime(endTime.Unix())

	return s.UserDataDB(uid).Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<?", uid, false, startTransactionTime, endTransactionTime).Count(&models.Transaction{})
}

// CreateTransaction saves a new transaction to database
func (s *TransactionService) CreateTransaction(transaction *models.Transaction, tagIds []int64) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	// Check whether account id is valid
	err := s.isAccountIdValid(transaction)

	if err != nil {
		return err
	}

	now := time.Now().Unix()

	transaction.TransactionId = s.GenerateUuid(uuid.UUID_TYPE_TRANSACTION)
	transaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

	transaction.CreatedUnixTime = now
	transaction.UpdatedUnixTime = now

	tagIds = utils.ToUniqueInt64Slice(tagIds)
	transactionTagIndexs := make([]*models.TransactionTagIndex, len(tagIds))

	for i := 0; i < len(tagIds); i++ {
		transactionTagIndexs[i] = &models.TransactionTagIndex{
			Uid:             transaction.Uid,
			TagId:           tagIds[i],
			TransactionId:   transaction.TransactionId,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
	}

	return s.UserDataDB(transaction.Uid).DoTransaction(func(sess *xorm.Session) error {
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
			relatedTransaction = s.GetRelatedTransferTransaction(transaction, s.GenerateUuid(uuid.UUID_TYPE_TRANSACTION))
			transaction.RelatedId = relatedTransaction.TransactionId
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
func (s *TransactionService) ModifyTransaction(transaction *models.Transaction, addTagIds []int64, removeTagIds []int64) error {
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
		transactionTagIndexs[i] = &models.TransactionTagIndex{
			Uid:             transaction.Uid,
			TagId:           addTagIds[i],
			TransactionId:   transaction.TransactionId,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		}
	}

	err := s.UserDB().DoTransaction(func(sess *xorm.Session) error {
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

		if transaction.Comment != oldTransaction.Comment {
			updateCols = append(updateCols, "comment")
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
			relatedTransaction := s.GetRelatedTransferTransaction(transaction, transaction.RelatedId)

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
			deletedRows, err := sess.Where("uid=?", transaction.Uid).In("tag_id", removeTagIds).Delete(&models.TransactionTagIndex{})

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
func (s *TransactionService) DeleteTransaction(uid int64, transactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Transaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
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

// GetRelatedTransferTransaction returns the related transaction for transfer transaction
func (s *TransactionService) GetRelatedTransferTransaction(originalTransaction *models.Transaction, relatedTransactionId int64) *models.Transaction {
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
		TransactionId:        relatedTransactionId,
		Uid:                  originalTransaction.Uid,
		Deleted:              originalTransaction.Deleted,
		Type:                 relatedType,
		CategoryId:           originalTransaction.CategoryId,
		TransactionTime:      relatedTransactionTime,
		AccountId:            originalTransaction.RelatedAccountId,
		Amount:               originalTransaction.RelatedAccountAmount,
		RelatedId:            originalTransaction.TransactionId,
		RelatedAccountId:     originalTransaction.AccountId,
		RelatedAccountAmount: originalTransaction.Amount,
		Comment:              originalTransaction.Comment,
		CreatedUnixTime:      originalTransaction.CreatedUnixTime,
		UpdatedUnixTime:      originalTransaction.UpdatedUnixTime,
		DeletedUnixTime:      originalTransaction.DeletedUnixTime,
	}

	return relatedTransaction
}

// GetAccountsTotalIncomeAndExpense returns the every accounts total income and expense amount by specific date range
func (s *TransactionService) GetAccountsTotalIncomeAndExpense(uid int64, startUnixTime int64, endUnixTime int64) (map[int64]int64, map[int64]int64, error) {
	if uid <= 0 {
		return nil, nil, errs.ErrUserIdInvalid
	}

	startTransactionTime := utils.GetMinTransactionTimeFromUnixTime(startUnixTime)
	endTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(endUnixTime)

	var transactionTotalAmounts []*models.Transaction
	err := s.UserDataDB(uid).Select("uid, type, account_id, SUM(amount) as amount").Where("uid=? AND deleted=? AND (type=? OR type=?) AND transaction_time>=? AND transaction_time<=?", uid, false, models.TRANSACTION_DB_TYPE_INCOME, models.TRANSACTION_DB_TYPE_EXPENSE, startTransactionTime, endTransactionTime).GroupBy("type, account_id").Find(&transactionTotalAmounts)

	if err != nil {
		return nil, nil, err
	}

	incomeAmounts := make(map[int64]int64)
	expenseAmounts := make(map[int64]int64)

	for i := 0; i < len(transactionTotalAmounts); i++ {
		transactionTotalAmount := transactionTotalAmounts[i]

		if transactionTotalAmount.Type == models.TRANSACTION_DB_TYPE_INCOME {
			incomeAmounts[transactionTotalAmount.AccountId] = transactionTotalAmount.Amount
		} else if transactionTotalAmount.Type == models.TRANSACTION_DB_TYPE_EXPENSE {
			expenseAmounts[transactionTotalAmount.AccountId] = transactionTotalAmount.Amount
		}
	}

	return incomeAmounts, expenseAmounts, nil
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
		err := sess.Where("uid=?", transaction.Uid).In("tag_id", tagIds).Find(&tags)

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
