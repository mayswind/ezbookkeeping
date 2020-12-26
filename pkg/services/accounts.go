package services

import (
	"github.com/mayswind/lab/pkg/utils"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/uuid"
)

// AccountService represents account service
type AccountService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a account service singleton instance
var (
	Accounts = &AccountService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllAccountsByUid returns all account models of user
func (s *AccountService) GetAllAccountsByUid(uid int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).OrderBy("parent_account_id asc, display_order asc").Find(&accounts)

	return accounts, err
}

// GetAccountAndSubAccountsByAccountId returns account model and sub account models according to account id
func (s *AccountService) GetAccountAndSubAccountsByAccountId(uid int64, accountId int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if accountId <= 0 {
		return nil, errs.ErrAccountIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).Where("uid=? AND deleted=? AND (account_id=? OR parent_account_id=?)", uid, false, accountId, accountId).OrderBy("parent_account_id asc, display_order asc").Find(&accounts)

	return accounts, err
}

// GetMaxDisplayOrder returns the max display order according to account category
func (s *AccountService) GetMaxDisplayOrder(uid int64, category models.AccountCategory) (int, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	account := &models.Account{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "parent_account_id", "display_order").Where("uid=? AND deleted=? AND parent_account_id=? AND category=?", uid, false, models.LevelOneAccountParentId, category).OrderBy("display_order desc").Limit(1).Get(account)

	if err != nil {
		return 0, err
	}

	if has {
		return account.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// GetMaxSubAccountDisplayOrder returns the max display order of sub account according to account category and parent account id
func (s *AccountService) GetMaxSubAccountDisplayOrder(uid int64, category models.AccountCategory, parentAccountId int64) (int, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	if parentAccountId <= 0 {
		return 0, errs.ErrAccountIdInvalid
	}

	account := &models.Account{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "parent_account_id", "display_order").Where("uid=? AND deleted=? AND parent_account_id=? AND category=?", uid, false, parentAccountId, category).OrderBy("display_order desc").Limit(1).Get(account)

	if err != nil {
		return 0, err
	}

	if has {
		return account.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateAccounts saves a new account model to database
func (s *AccountService) CreateAccounts(mainAccount *models.Account, childrenAccounts []*models.Account) error {
	if mainAccount.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	allAccounts := make([]*models.Account, len(childrenAccounts)+1)
	var allInitTransactions []*models.Transaction

	mainAccount.AccountId = s.GenerateUuid(uuid.UUID_TYPE_ACCOUNT)
	allAccounts[0] = mainAccount

	if mainAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		for i := 0; i < len(childrenAccounts); i++ {
			childAccount := childrenAccounts[i]
			childAccount.AccountId = s.GenerateUuid(uuid.UUID_TYPE_ACCOUNT)
			childAccount.ParentAccountId = mainAccount.AccountId
			childAccount.Uid = mainAccount.Uid
			childAccount.Type = models.ACCOUNT_TYPE_SINGLE_ACCOUNT

			allAccounts[i+1] = childrenAccounts[i]
		}
	}

	transactionTime := utils.GetMinTransactionTimeFromUnixTime(now)

	for i := 0; i < len(allAccounts); i++ {
		allAccounts[i].Deleted = false
		allAccounts[i].CreatedUnixTime = now
		allAccounts[i].UpdatedUnixTime = now

		if allAccounts[i].Balance != 0 {
			newTransaction := &models.Transaction{
				TransactionId:        s.GenerateUuid(uuid.UUID_TYPE_TRANSACTION),
				Uid:                  allAccounts[i].Uid,
				Deleted:              false,
				Type:                 models.TRANSACTION_TYPE_MODIFY_BALANCE,
				TransactionTime:      transactionTime,
				SourceAccountId:      allAccounts[i].AccountId,
				DestinationAccountId: allAccounts[i].AccountId,
				SourceAmount:         allAccounts[i].Balance,
				DestinationAmount:    allAccounts[i].Balance,
				CreatedUnixTime:      now,
				UpdatedUnixTime:      now,
			}

			transactionTime++
			allInitTransactions = append(allInitTransactions, newTransaction)
		}
	}

	return s.UserDataDB(mainAccount.Uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(allAccounts); i++ {
			account := allAccounts[i]
			_, err := sess.Insert(account)

			if err != nil {
				return err
			}
		}

		for i := 0; i < len(allInitTransactions); i++ {
			transaction := allInitTransactions[i]
			_, err := sess.Insert(transaction)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

// ModifyAccounts saves an existed account model to database
func (s *AccountService) ModifyAccounts(uid int64, accounts []*models.Account) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	for i := 0; i < len(accounts); i++ {
		accounts[i].UpdatedUnixTime = now
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(accounts); i++ {
			account := accounts[i]
			updatedRows, err := sess.ID(account.AccountId).Cols("name", "category", "icon", "color", "comment", "hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(account)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrAccountNotFound
			}
		}

		return nil
	})
}

// HideAccount updates hidden field of given accounts
func (s *AccountService) HideAccount(uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Account{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).In("account_id", ids).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrAccountNotFound
		}

		return nil
	})
}

// ModifyAccountDisplayOrders updates display order of given accounts
func (s *AccountService) ModifyAccountDisplayOrders(uid int64, accounts []*models.Account) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(accounts); i++ {
		accounts[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(accounts); i++ {
			account := accounts[i]
			updatedRows, err := sess.ID(account.AccountId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(account)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrAccountNotFound
			}
		}

		return nil
	})
}

// DeleteAccount deletes an existed account from database
func (s *AccountService) DeleteAccount(uid int64, accountId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Account{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		var accountAndSubAccounts []*models.Account
		err := s.UserDataDB(uid).Where("uid=? AND deleted=? AND (account_id=? OR parent_account_id=?)", uid, false, accountId, accountId).Find(&accountAndSubAccounts)

		if err != nil {
			return err
		} else if len(accountAndSubAccounts) < 1 {
			return errs.ErrAccountNotFound
		}

		accountAndSubAccountIds := make([]int64, len(accountAndSubAccounts))

		for i := 0; i < len(accountAndSubAccounts); i++ {
			accountAndSubAccountIds[i] = accountAndSubAccounts[i].AccountId
		}

		exists, err := sess.Cols("uid", "deleted", "source_account_id").Where("uid=? AND deleted=?", uid, false).In("source_account_id", accountAndSubAccountIds).Limit(1).Exist(&models.Transaction{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrAccountInUseCannotBeDeleted
		}

		exists, err = sess.Cols("uid", "deleted", "destination_account_id").Where("uid=? AND deleted=?", uid, false).In("destination_account_id", accountAndSubAccountIds).Limit(1).Exist(&models.Transaction{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrAccountInUseCannotBeDeleted
		}

		deletedRows, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).In("account_id", accountAndSubAccountIds).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrAccountNotFound
		}

		return err
	})
}
