package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/uuid"
)

type AccountService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

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

func (s *AccountService) GetAllAccountsByUid(uid int64) ([]*models.Account, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var accounts []*models.Account
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).OrderBy("parent_account_id asc, display_order asc").Find(&accounts)

	return accounts, err
}

func (s *AccountService) GetAccountByAccountId(uid int64, accountId int64) ([]*models.Account, error) {
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

func (s *AccountService) GetMaxDisplayOrder(uid int64, category models.AccountCategory) (int, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	account := &models.Account{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "parent_account_id", "display_order").Where("uid=? AND deleted=? AND parent_account_id=? AND category=?", uid, false, models.ACCOUNT_PARENT_ID_LEVEL_ONE, category).OrderBy("display_order desc").Limit(1).Get(account)

	if err != nil {
		return 0, err
	}

	if has {
		return account.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

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

func (s *AccountService) CreateAccounts(mainAccount *models.Account, childrenAccounts []*models.Account) error {
	if mainAccount.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	allAccounts := make([]*models.Account, len(childrenAccounts)+1)

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

	for i := 0; i < len(allAccounts); i++ {
		allAccounts[i].Deleted = false
		allAccounts[i].CreatedUnixTime = time.Now().Unix()
		allAccounts[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(mainAccount.Uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(allAccounts); i++ {
			account := allAccounts[i]
			_, err := sess.Insert(account)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

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
			_, err := sess.Cols("name", "category", "icon", "color", "comment", "hidden", "updated_unix_time").Where("account_id=? AND uid=? AND deleted=?", account.AccountId, uid, false).Update(account)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *AccountService) HideAccount(uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Account{
		Hidden: hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		deletedRows, err := sess.Cols("hidden", "updated_unix_time").In("account_id", ids).Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if deletedRows < 1 {
			return errs.ErrAccountNotFound
		}

		return err
	})
}

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
			_, err := sess.Cols("display_order", "updated_unix_time").Where("account_id=? AND uid=? AND deleted=?", account.AccountId, uid, false).Update(account)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *AccountService) DeleteAccounts(uid int64, ids []int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.Account{
		Deleted: true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		deletedRows, err := sess.Cols("deleted", "deleted_unix_time").In("account_id", ids).Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if deletedRows < 1 {
			return errs.ErrAccountNotFound
		}

		return err
	})
}
