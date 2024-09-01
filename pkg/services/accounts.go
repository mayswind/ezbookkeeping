package services

import (
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

// GetTotalAccountCountByUid returns total account count of user
func (s *AccountService) GetTotalAccountCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.Account{})

	return count, err
}

// GetAllAccountsByUid returns all account models of user
func (s *AccountService) GetAllAccountsByUid(c core.Context, uid int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).OrderBy("parent_account_id asc, display_order asc").Find(&accounts)

	return accounts, err
}

// GetAccountAndSubAccountsByAccountId returns account model and sub-account models according to account id
func (s *AccountService) GetAccountAndSubAccountsByAccountId(c core.Context, uid int64, accountId int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if accountId <= 0 {
		return nil, errs.ErrAccountIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND (account_id=? OR parent_account_id=?)", uid, false, accountId, accountId).OrderBy("parent_account_id asc, display_order asc").Find(&accounts)

	return accounts, err
}

// GetSubAccountsByAccountId returns sub-account models according to account id
func (s *AccountService) GetSubAccountsByAccountId(c core.Context, uid int64, accountId int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if accountId <= 0 {
		return nil, errs.ErrAccountIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND parent_account_id=?", uid, false, accountId).OrderBy("display_order asc").Find(&accounts)

	return accounts, err
}

// GetSubAccountsByAccountIds returns sub-account models according to account ids
func (s *AccountService) GetSubAccountsByAccountIds(c core.Context, uid int64, accountIds []int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if len(accountIds) <= 0 {
		return nil, errs.ErrAccountIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, len(accountIds)+2)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	var accountIdConditions strings.Builder

	for i := 0; i < len(accountIds); i++ {
		if accountIds[i] <= 0 {
			return nil, errs.ErrAccountIdInvalid
		}

		if accountIdConditions.Len() > 0 {
			accountIdConditions.WriteString(",")
		}

		accountIdConditions.WriteString("?")
		conditionParams = append(conditionParams, accountIds[i])
	}

	if accountIdConditions.Len() > 1 {
		condition = condition + " AND parent_account_id IN (" + accountIdConditions.String() + ")"
	} else {
		condition = condition + " AND parent_account_id = " + accountIdConditions.String()
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("display_order asc").Find(&accounts)

	return accounts, err
}

// GetAccountsByAccountIds returns account models according to account ids
func (s *AccountService) GetAccountsByAccountIds(c core.Context, uid int64, accountIds []int64) (map[int64]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if accountIds == nil {
		return nil, errs.ErrAccountIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("account_id", accountIds).Find(&accounts)

	if err != nil {
		return nil, err
	}

	accountMap := s.GetAccountMapByList(accounts)
	return accountMap, err
}

// GetMaxDisplayOrder returns the max display order according to account category
func (s *AccountService) GetMaxDisplayOrder(c core.Context, uid int64, category models.AccountCategory) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	account := &models.Account{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "parent_account_id", "display_order").Where("uid=? AND deleted=? AND parent_account_id=? AND category=?", uid, false, models.LevelOneAccountParentId, category).OrderBy("display_order desc").Limit(1).Get(account)

	if err != nil {
		return 0, err
	}

	if has {
		return account.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// GetMaxSubAccountDisplayOrder returns the max display order of sub-account according to account category and parent account id
func (s *AccountService) GetMaxSubAccountDisplayOrder(c core.Context, uid int64, category models.AccountCategory, parentAccountId int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	if parentAccountId <= 0 {
		return 0, errs.ErrAccountIdInvalid
	}

	account := &models.Account{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "parent_account_id", "display_order").Where("uid=? AND deleted=? AND parent_account_id=? AND category=?", uid, false, parentAccountId, category).OrderBy("display_order desc").Limit(1).Get(account)

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
func (s *AccountService) CreateAccounts(c core.Context, mainAccount *models.Account, childrenAccounts []*models.Account, utcOffset int16) error {
	if mainAccount.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	allAccounts := make([]*models.Account, len(childrenAccounts)+1)
	var allInitTransactions []*models.Transaction

	mainAccount.AccountId = s.GenerateUuid(uuid.UUID_TYPE_ACCOUNT)

	if mainAccount.AccountId < 1 {
		return errs.ErrSystemIsBusy
	}

	allAccounts[0] = mainAccount

	if mainAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		for i := 0; i < len(childrenAccounts); i++ {
			childAccount := childrenAccounts[i]
			childAccount.AccountId = s.GenerateUuid(uuid.UUID_TYPE_ACCOUNT)

			if childAccount.AccountId < 1 {
				return errs.ErrSystemIsBusy
			}

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
			transactionId := s.GenerateUuid(uuid.UUID_TYPE_TRANSACTION)

			if transactionId < 1 {
				return errs.ErrSystemIsBusy
			}

			newTransaction := &models.Transaction{
				TransactionId:        transactionId,
				Uid:                  allAccounts[i].Uid,
				Deleted:              false,
				Type:                 models.TRANSACTION_DB_TYPE_MODIFY_BALANCE,
				TransactionTime:      transactionTime,
				TimezoneUtcOffset:    utcOffset,
				AccountId:            allAccounts[i].AccountId,
				Amount:               allAccounts[i].Balance,
				RelatedAccountId:     allAccounts[i].AccountId,
				RelatedAccountAmount: allAccounts[i].Balance,
				CreatedUnixTime:      now,
				UpdatedUnixTime:      now,
			}

			transactionTime++
			allInitTransactions = append(allInitTransactions, newTransaction)
		}
	}

	return s.UserDataDB(mainAccount.Uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *AccountService) ModifyAccounts(c core.Context, uid int64, accounts []*models.Account) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	for i := 0; i < len(accounts); i++ {
		accounts[i].UpdatedUnixTime = now
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *AccountService) HideAccount(c core.Context, uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Account{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *AccountService) ModifyAccountDisplayOrders(c core.Context, uid int64, accounts []*models.Account) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(accounts); i++ {
		accounts[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *AccountService) DeleteAccount(c core.Context, uid int64, accountId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Account{
		Balance:         0,
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		var accountAndSubAccounts []*models.Account
		err := sess.Where("uid=? AND deleted=? AND (account_id=? OR parent_account_id=?)", uid, false, accountId, accountId).Find(&accountAndSubAccounts)

		if err != nil {
			return err
		} else if len(accountAndSubAccounts) < 1 {
			return errs.ErrAccountNotFound
		}

		accountAndSubAccountIds := make([]int64, len(accountAndSubAccounts))

		for i := 0; i < len(accountAndSubAccounts); i++ {
			accountAndSubAccountIds[i] = accountAndSubAccounts[i].AccountId
		}

		var relatedTransactionsByAccount []*models.Transaction
		err = sess.Cols("transaction_id", "uid", "deleted", "account_id", "type").Where("uid=? AND deleted=?", uid, false).In("account_id", accountAndSubAccountIds).Limit(len(accountAndSubAccounts) + 1).Find(&relatedTransactionsByAccount)

		if err != nil {
			return err
		} else if len(relatedTransactionsByAccount) > len(accountAndSubAccountIds) {
			return errs.ErrAccountInUseCannotBeDeleted
		} else if len(relatedTransactionsByAccount) > 0 {
			accountTransactionExists := make(map[int64]bool)

			for i := 0; i < len(relatedTransactionsByAccount); i++ {
				transaction := relatedTransactionsByAccount[i]

				if transaction.Type != models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
					return errs.ErrAccountInUseCannotBeDeleted
				} else if _, exists := accountTransactionExists[transaction.AccountId]; exists {
					return errs.ErrAccountInUseCannotBeDeleted
				}

				accountTransactionExists[transaction.AccountId] = true
			}
		}

		deletedRows, err := sess.Cols("balance", "deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).In("account_id", accountAndSubAccountIds).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrAccountNotFound
		}

		if len(relatedTransactionsByAccount) > 0 {
			updateTransaction := &models.Transaction{
				Deleted:         true,
				DeletedUnixTime: now,
			}

			transactionIds := make([]int64, len(relatedTransactionsByAccount))

			for i := 0; i < len(relatedTransactionsByAccount); i++ {
				transactionIds[i] = relatedTransactionsByAccount[i].TransactionId
			}

			deletedTransactionRows, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).In("transaction_id", transactionIds).Update(updateTransaction)

			if err != nil {
				return err
			} else if deletedTransactionRows < int64(len(transactionIds)) {
				return errs.ErrDatabaseOperationFailed
			}
		}

		return err
	})
}

// GetAccountMapByList returns an account map by a list
func (s *AccountService) GetAccountMapByList(accounts []*models.Account) map[int64]*models.Account {
	accountMap := make(map[int64]*models.Account)

	for i := 0; i < len(accounts); i++ {
		account := accounts[i]
		accountMap[account.AccountId] = account
	}
	return accountMap
}

// GetAccountNameMapByList returns an account map by a list
func (s *AccountService) GetAccountNameMapByList(accounts []*models.Account) map[string]*models.Account {
	accountMap := make(map[string]*models.Account)

	for i := 0; i < len(accounts); i++ {
		account := accounts[i]
		accountMap[account.Name] = account
	}
	return accountMap
}
