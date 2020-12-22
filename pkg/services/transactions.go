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

type TransactionService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

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

func (s *TransactionService) GetTransactionsByMaxTime(uid int64, maxTime int64, count int) ([]*models.Transaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if count < 1 {
		return nil, errs.ErrPageCountInvalid
	}

	var transactions []*models.Transaction
	err := s.UserDataDB(uid).Where("uid=? AND deleted=? AND transaction_time<=?", uid, false, maxTime).Limit(count, 0).OrderBy("transaction_time desc").Find(&transactions)

	return transactions, err
}

func (s *TransactionService) GetTransactionsInMonthByPage(uid int64, year int, month int, page int, count int) ([]*models.Transaction, error) {
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

	startUnixTime := startTime.Unix()
	endUnixTime := endTime.Unix()

	var transactions []*models.Transaction
	err = s.UserDataDB(uid).Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<?", uid, false, startUnixTime, endUnixTime).Limit(count, count*(page-1)).OrderBy("transaction_time desc").Find(&transactions)

	return transactions, err
}

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

func (s *TransactionService) GetAllTransactionCount(uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	return s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).Count(&models.Transaction{})
}

func (s *TransactionService) GetMonthTransactionCount(uid int64, year int64, month int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	startTime, err := utils.ParseFromLongDateTime(fmt.Sprintf("%d-%d-01 00:00:00", year, month))

	if err != nil {
		return 0, errs.ErrSystemError
	}

	endTime := startTime.AddDate(0, 1, 0)

	startUnixTime := startTime.Unix()
	endUnixTime := endTime.Unix()

	return s.UserDataDB(uid).Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<?", uid, false, startUnixTime, endUnixTime).Count(&models.Transaction{})
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction, tagIds []int64) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if transaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE ||
		transaction.Type == models.TRANSACTION_TYPE_INCOME ||
		transaction.Type == models.TRANSACTION_TYPE_EXPENSE {
		if transaction.SourceAccountId != transaction.DestinationAccountId {
			return errs.ErrTransactionSourceAndDestinationIdNotEqual
		} else if transaction.SourceAmount != transaction.DestinationAmount {
			return errs.ErrTransactionSourceAndDestinationAmountNotEqual
		}
	} else if transaction.Type == models.TRANSACTION_TYPE_TRANSFER {
		if transaction.SourceAccountId == transaction.DestinationAccountId {
			return errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
		}
	} else {
		return errs.ErrTransactionTypeInvalid
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
		sourceAccount := &models.Account{}
		destinationAccount := &models.Account{}
		has, err := sess.ID(transaction.SourceAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(sourceAccount)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrSourceAccountNotFound
		} else if sourceAccount.Hidden {
			return errs.ErrCannotAddTransactionToHiddenAccount
		}

		if transaction.DestinationAccountId == transaction.SourceAccountId {
			destinationAccount = sourceAccount
		} else {
			has, err = sess.ID(transaction.DestinationAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(destinationAccount)

			if err != nil {
				return err
			} else if !has {
				return errs.ErrDestinationAccountNotFound
			} else if destinationAccount.Hidden {
				return errs.ErrCannotAddTransactionToHiddenAccount
			}
		}

		if sourceAccount.Currency == destinationAccount.Currency && transaction.SourceAmount != transaction.DestinationAmount {
			return errs.ErrTransactionSourceAndDestinationAmountNotEqual
		}

		// Get and verify category
		category := &models.TransactionCategory{}

		if transaction.Type != models.TRANSACTION_TYPE_MODIFY_BALANCE {
			has, err = sess.ID(transaction.CategoryId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(category)

			if err != nil {
				return err
			} else if !has {
				return errs.ErrTransactionCategoryNotFound
			}

			if category.ParentCategoryId < 1 {
				return errs.ErrCannotUsePrimaryCategoryForTransaction
			}
		}

		if (transaction.Type == models.TRANSACTION_TYPE_INCOME && category.Type != models.CATEGORY_TYPE_INCOME) ||
			(transaction.Type == models.TRANSACTION_TYPE_EXPENSE && category.Type != models.CATEGORY_TYPE_EXPENSE) ||
			(transaction.Type == models.TRANSACTION_TYPE_TRANSFER && category.Type != models.CATEGORY_TYPE_TRANSFER) {
			return errs.ErrTransactionCategoryTypeInvalid
		}

		// Get and verify tags
		if len(transactionTagIndexs) > 0 {
			var tags []*models.TransactionTag
			err := sess.Where("uid=? AND deleted=?", transaction.Uid, false).In("tag_ids", tagIds).Find(&tags)

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

		// Verify balance modification transaction and calculate real amount
		if transaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
			otherTransactionExists, err := sess.Cols("uid", "deleted", "destination_account_id").Where("uid=? AND deleted=? AND destination_account_id=?", transaction.Uid, false, destinationAccount.AccountId).Limit(1).Exist(&models.Transaction{})

			if err != nil {
				return err
			} else if otherTransactionExists {
				return errs.ErrBalanceModificationTransactionCannotAddWhenNotEmpty
			}

			transaction.DestinationAmount = transaction.SourceAmount - destinationAccount.Balance
		}

		// Insert transaction row
		createdRows, err := sess.Insert(transaction)

		if err != nil || createdRows < 1 { // maybe another transaction has same time
			sameSecondLatestTransaction := &models.Transaction{}
			minTransactionTime := utils.GetMinTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))
			maxTransactionTime := utils.GetMaxTransactionTimeFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime))

			has, err = sess.Where("uid=? AND deleted=? AND transaction_time>=? AND transaction_time<=?", transaction.Uid, false, minTransactionTime, maxTransactionTime).OrderBy("transaction_time desc").Limit(1).Get(sameSecondLatestTransaction)

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
		if transaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_TYPE_INCOME {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_TYPE_EXPENSE {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", transaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if transaction.Type == models.TRANSACTION_TYPE_TRANSFER {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedSourceRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", transaction.SourceAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedSourceRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}

			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedDestinationRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", transaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedDestinationRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		}

		return err
	})
}

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

		if oldTransaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE ||
			oldTransaction.Type == models.TRANSACTION_TYPE_INCOME ||
			oldTransaction.Type == models.TRANSACTION_TYPE_EXPENSE {
			if transaction.SourceAccountId != transaction.DestinationAccountId {
				return errs.ErrTransactionSourceAndDestinationIdNotEqual
			} else if transaction.SourceAmount != transaction.DestinationAmount {
				return errs.ErrTransactionSourceAndDestinationAmountNotEqual
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_TRANSFER {
			if transaction.SourceAccountId == transaction.DestinationAccountId {
				return errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
			}
		} else {
			return errs.ErrTransactionTypeInvalid
		}

		// Get and verify source and destination account (if necessary)
		sourceAccount := &models.Account{}
		destinationAccount := &models.Account{}
		has, err = sess.ID(transaction.SourceAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(sourceAccount)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrSourceAccountNotFound
		} else if sourceAccount.Hidden {
			return errs.ErrCannotModifyTransactionInHiddenAccount
		}

		if transaction.DestinationAccountId == transaction.SourceAccountId {
			destinationAccount = sourceAccount
		} else {
			has, err = sess.ID(transaction.DestinationAccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(destinationAccount)

			if err != nil {
				return err
			} else if !has {
				return errs.ErrDestinationAccountNotFound
			} else if destinationAccount.Hidden {
				return errs.ErrCannotModifyTransactionInHiddenAccount
			}
		}

		// Append modified columns and verify
		if transaction.CategoryId != oldTransaction.CategoryId {
			if oldTransaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
				if transaction.CategoryId > 0 {
					return errs.ErrBalanceModificationTransactionCannotSetCategory
				}
			} else {
				category := &models.TransactionCategory{}
				has, err = sess.ID(transaction.CategoryId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(category)

				if err != nil {
					return err
				} else if !has {
					return errs.ErrTransactionCategoryNotFound
				}

				if category.ParentCategoryId < 1 {
					return errs.ErrCannotUsePrimaryCategoryForTransaction
				}

				if (oldTransaction.Type == models.TRANSACTION_TYPE_INCOME && category.Type != models.CATEGORY_TYPE_INCOME) ||
					(oldTransaction.Type == models.TRANSACTION_TYPE_EXPENSE && category.Type != models.CATEGORY_TYPE_EXPENSE) ||
					(oldTransaction.Type == models.TRANSACTION_TYPE_TRANSFER && category.Type != models.CATEGORY_TYPE_TRANSFER) {
					return errs.ErrTransactionCategoryTypeInvalid
				}
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

		if transaction.SourceAccountId != oldTransaction.SourceAccountId {
			updateCols = append(updateCols, "source_account_id")
		}

		if transaction.DestinationAccountId != oldTransaction.DestinationAccountId {
			updateCols = append(updateCols, "destination_account_id")
		}

		if transaction.SourceAmount != oldTransaction.SourceAmount {
			if oldTransaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
				originalBalance := sourceAccount.Balance - oldTransaction.DestinationAmount
				transaction.DestinationAmount = transaction.SourceAmount - originalBalance
			}

			updateCols = append(updateCols, "source_amount")
		}

		if transaction.DestinationAmount != oldTransaction.DestinationAmount {
			updateCols = append(updateCols, "destination_amount")
		}

		if transaction.Comment != oldTransaction.Comment {
			updateCols = append(updateCols, "comment")
		}

		// Get and verify tags
		if len(transactionTagIndexs) > 0 {
			var tags []*models.TransactionTag
			err := sess.Where("uid=? AND deleted=?", transaction.Uid, false).In("tag_ids", addTagIds).Find(&tags)

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

		// Update transaction row
		updatedRows, err := sess.ID(transaction.TransactionId).Cols(updateCols...).Where("uid=? AND deleted=?", transaction.Uid, false).Update(transaction)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionNotFound
		}

		// Update transaction tag index
		if len(removeTagIds) > 0 {
			deletedRows, err := sess.Where("uid=?", transaction.Uid).In("tag_id", removeTagIds).Delete(&models.TransactionTagIndex{})

			if err != nil {
				return err
			} else if deletedRows < 1 {
				return errs.ErrTransactionTagNotFound
			}

			return err
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
		if oldTransaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
			if transaction.SourceAccountId != oldTransaction.SourceAccountId {
				return errs.ErrBalanceModificationTransactionCannotChangeAccountId
			}

			if transaction.SourceAmount != oldTransaction.SourceAmount {
				destinationAccount.UpdatedUnixTime = time.Now().Unix()
				updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)+(%d)", oldTransaction.DestinationAmount, transaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

				if err != nil {
					return err
				} else if updatedRows < 1 {
					return errs.ErrDatabaseOperationFailed
				}
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_INCOME {
			if transaction.SourceAccountId != oldTransaction.SourceAccountId && transaction.DestinationAmount != oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.SourceAccountId != oldTransaction.SourceAccountId && transaction.DestinationAmount == oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.SourceAccountId == oldTransaction.SourceAccountId && transaction.DestinationAmount != oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_EXPENSE {
			if transaction.SourceAccountId != oldTransaction.SourceAccountId && transaction.DestinationAmount != oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.SourceAccountId != oldTransaction.SourceAccountId && transaction.DestinationAmount == oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.SourceAccountId == oldTransaction.SourceAccountId && transaction.DestinationAmount != oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_TRANSFER {
			if transaction.SourceAccountId != oldTransaction.SourceAccountId && transaction.SourceAmount != oldTransaction.SourceAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.SourceAccountId != oldTransaction.SourceAccountId && transaction.SourceAmount == oldTransaction.SourceAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.SourceAccountId == oldTransaction.SourceAccountId && transaction.SourceAmount != oldTransaction.SourceAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			}

			if transaction.DestinationAccountId != oldTransaction.DestinationAccountId && transaction.DestinationAmount != oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.DestinationAccountId != oldTransaction.DestinationAccountId && transaction.DestinationAmount == oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			} else if transaction.DestinationAccountId == oldTransaction.DestinationAccountId && transaction.DestinationAmount != oldTransaction.DestinationAmount {
				// TODO: implement
				return errs.ErrNotImplemented
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

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
		sourceAccount := &models.Account{}
		destinationAccount := &models.Account{}
		has, err = sess.ID(oldTransaction.SourceAccountId).Where("uid=? AND deleted=?", oldTransaction.Uid, false).Get(sourceAccount)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrSourceAccountNotFound
		} else if sourceAccount.Hidden {
			return errs.ErrCannotDeleteTransactionInHiddenAccount
		}

		if oldTransaction.DestinationAccountId == oldTransaction.SourceAccountId {
			destinationAccount = sourceAccount
		} else {
			has, err = sess.ID(oldTransaction.DestinationAccountId).Where("uid=? AND deleted=?", oldTransaction.Uid, false).Get(destinationAccount)

			if err != nil {
				return err
			} else if !has {
				return errs.ErrDestinationAccountNotFound
			} else if destinationAccount.Hidden {
				return errs.ErrCannotDeleteTransactionInHiddenAccount
			}
		}

		// Update transaction row to deleted
		deletedRows, err := sess.ID(transactionId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionNotFound
		}

		// Update account table
		if oldTransaction.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_INCOME {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_EXPENSE {
			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", oldTransaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		} else if oldTransaction.Type == models.TRANSACTION_TYPE_TRANSFER {
			sourceAccount.UpdatedUnixTime = time.Now().Unix()
			updatedSourceRows, err := sess.ID(sourceAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", oldTransaction.SourceAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", sourceAccount.Uid, false).Update(sourceAccount)

			if err != nil {
				return err
			} else if updatedSourceRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}

			destinationAccount.UpdatedUnixTime = time.Now().Unix()
			updatedDestinationRows, err := sess.ID(destinationAccount.AccountId).SetExpr("balance", fmt.Sprintf("balance-(%d)", oldTransaction.DestinationAmount)).Cols("updated_unix_time").Where("uid=? AND deleted=?", destinationAccount.Uid, false).Update(destinationAccount)

			if err != nil {
				return err
			} else if updatedDestinationRows < 1 {
				return errs.ErrDatabaseOperationFailed
			}
		}

		return err
	})
}
