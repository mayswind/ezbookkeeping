package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// AccountsApi represents account api
type AccountsApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	accounts *services.AccountService
}

// Initialize an account api singleton instance
var (
	Accounts = &AccountsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			container: duplicatechecker.Container,
		},
		accounts: services.Accounts,
	}
)

// AccountListHandler returns accounts list of current user
func (a *AccountsApi) AccountListHandler(c *core.WebContext) (any, *errs.Error) {
	var accountListReq models.AccountListRequest
	err := c.ShouldBindQuery(&accountListReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accounts, err := a.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[accounts.AccountListHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	userAllAccountResps := make([]*models.AccountInfoResponse, len(accounts))
	userAllAccountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(accounts); i++ {
		userAllAccountResps[i] = accounts[i].ToAccountInfoResponse()
		userAllAccountRespMap[userAllAccountResps[i].Id] = userAllAccountResps[i]
	}

	for i := 0; i < len(userAllAccountResps); i++ {
		userAccountResp := userAllAccountResps[i]

		if accountListReq.VisibleOnly && userAccountResp.Hidden {
			continue
		}

		if userAccountResp.ParentId <= models.LevelOneAccountParentId {
			continue
		}

		parentAccount, parentExists := userAllAccountRespMap[userAccountResp.ParentId]

		if !parentExists || parentAccount == nil {
			continue
		}

		parentAccount.SubAccounts = append(parentAccount.SubAccounts, userAccountResp)
	}

	userFinalAccountResps := make(models.AccountInfoResponseSlice, 0, len(userAllAccountResps))

	for i := 0; i < len(userAllAccountResps); i++ {
		if userAllAccountResps[i].ParentId == models.LevelOneAccountParentId && (!accountListReq.VisibleOnly || !userAllAccountResps[i].Hidden) {
			sort.Sort(userAllAccountResps[i].SubAccounts)
			userFinalAccountResps = append(userFinalAccountResps, userAllAccountResps[i])
		}
	}

	sort.Sort(userFinalAccountResps)

	return userFinalAccountResps, nil
}

// AccountGetHandler returns one specific account of current user
func (a *AccountsApi) AccountGetHandler(c *core.WebContext) (any, *errs.Error) {
	var accountGetReq models.AccountGetRequest
	err := c.ShouldBindQuery(&accountGetReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountGetReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.AccountGetHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", accountGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	accountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(accountAndSubAccounts); i++ {
		accountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
		accountRespMap[accountResp.Id] = accountResp
	}

	accountResp, exists := accountRespMap[accountGetReq.Id]

	if !exists {
		return nil, errs.ErrAccountNotFound
	}

	for i := 0; i < len(accountAndSubAccounts); i++ {
		if accountAndSubAccounts[i].ParentAccountId == accountResp.Id {
			subAccountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
			accountResp.SubAccounts = append(accountResp.SubAccounts, subAccountResp)
		}
	}

	sort.Sort(accountResp.SubAccounts)

	return accountResp, nil
}

// AccountCreateHandler saves a new account by request parameters for current user
func (a *AccountsApi) AccountCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var accountCreateReq models.AccountCreateRequest
	err := c.ShouldBindJSON(&accountCreateReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[accounts.AccountCreateHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	if accountCreateReq.Category < models.ACCOUNT_CATEGORY_CASH || accountCreateReq.Category > models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT {
		log.Warnf(c, "[accounts.AccountCreateHandler] account category invalid, category is %d", accountCreateReq.Category)
		return nil, errs.ErrAccountCategoryInvalid
	}

	if accountCreateReq.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
		if len(accountCreateReq.SubAccounts) > 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] account cannot have any sub-accounts")
			return nil, errs.ErrAccountCannotHaveSubAccounts
		}

		if accountCreateReq.Currency == validators.ParentAccountCurrencyPlaceholder {
			log.Warnf(c, "[accounts.AccountCreateHandler] account cannot set currency placeholder")
			return nil, errs.ErrAccountCurrencyInvalid
		}

		if accountCreateReq.Balance != 0 && accountCreateReq.BalanceTime <= 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] account balance time is not set")
			return nil, errs.ErrAccountBalanceTimeNotSet
		}

		if accountCreateReq.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD && accountCreateReq.CreditCardStatementDate != 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] cannot set statement date with category \"%d\"", accountCreateReq.Category)
			return nil, errs.ErrCannotSetStatementDateForNonCreditCard
		}
	} else if accountCreateReq.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		if len(accountCreateReq.SubAccounts) < 1 {
			log.Warnf(c, "[accounts.AccountCreateHandler] account does not have any sub-accounts")
			return nil, errs.ErrAccountHaveNoSubAccount
		}

		if accountCreateReq.Currency != validators.ParentAccountCurrencyPlaceholder {
			log.Warnf(c, "[accounts.AccountCreateHandler] parent account cannot set currency")
			return nil, errs.ErrParentAccountCannotSetCurrency
		}

		if accountCreateReq.Balance != 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] parent account cannot set balance")
			return nil, errs.ErrParentAccountCannotSetBalance
		}

		for i := 0; i < len(accountCreateReq.SubAccounts); i++ {
			subAccount := accountCreateReq.SubAccounts[i]

			if subAccount.Category != accountCreateReq.Category {
				log.Warnf(c, "[accounts.AccountCreateHandler] category of sub-account#%d not equals to parent", i)
				return nil, errs.ErrSubAccountCategoryNotEqualsToParent
			}

			if subAccount.Type != models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d type invalid", i)
				return nil, errs.ErrSubAccountTypeInvalid
			}

			if subAccount.Currency == validators.ParentAccountCurrencyPlaceholder {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d cannot set currency placeholder", i)
				return nil, errs.ErrAccountCurrencyInvalid
			}

			if subAccount.Balance != 0 && subAccount.BalanceTime <= 0 {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d balance time is not set", i)
				return nil, errs.ErrAccountBalanceTimeNotSet
			}

			if subAccount.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD && subAccount.CreditCardStatementDate != 0 {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d cannot set statement date with category \"%d\"", i, subAccount.Category)
				return nil, errs.ErrCannotSetStatementDateForNonCreditCard
			}
		}
	} else {
		log.Warnf(c, "[accounts.AccountCreateHandler] account type invalid, type is %d", accountCreateReq.Type)
		return nil, errs.ErrAccountTypeInvalid
	}

	uid := c.GetCurrentUid()
	maxOrderId, err := a.accounts.GetMaxDisplayOrder(c, uid, accountCreateReq.Category)

	if err != nil {
		log.Errorf(c, "[accounts.AccountCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	mainAccount := a.createNewAccountModel(uid, &accountCreateReq, maxOrderId+1)
	childrenAccounts, childrenAccountBalanceTimes := a.createSubAccountModels(uid, &accountCreateReq)

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && accountCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT, uid, accountCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[accounts.AccountCreateHandler] another account \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			accountId, err := utils.StringToInt64(remark)

			if err == nil {
				accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountId)

				if err != nil {
					log.Errorf(c, "[accounts.AccountCreateHandler] failed to get existed account \"id:%d\" for user \"uid:%d\", because %s", accountId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				accountMap := a.accounts.GetAccountMapByList(accountAndSubAccounts)
				mainAccount, exists := accountMap[accountId]

				if !exists {
					return nil, errs.ErrOperationFailed
				}

				accountInfoResp := mainAccount.ToAccountInfoResponse()

				for i := 0; i < len(accountAndSubAccounts); i++ {
					if accountAndSubAccounts[i].ParentAccountId == mainAccount.AccountId {
						subAccountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
						accountInfoResp.SubAccounts = append(accountInfoResp.SubAccounts, subAccountResp)
					}
				}

				return accountInfoResp, nil
			}
		}
	}

	err = a.accounts.CreateAccounts(c, mainAccount, accountCreateReq.BalanceTime, childrenAccounts, childrenAccountBalanceTimes, utcOffset)

	if err != nil {
		log.Errorf(c, "[accounts.AccountCreateHandler] failed to create account \"id:%d\" for user \"uid:%d\", because %s", mainAccount.AccountId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountCreateHandler] user \"uid:%d\" has created a new account \"id:%d\" successfully", uid, mainAccount.AccountId)

	a.SetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT, uid, accountCreateReq.ClientSessionId, utils.Int64ToString(mainAccount.AccountId))
	accountInfoResp := mainAccount.ToAccountInfoResponse()

	if len(childrenAccounts) > 0 {
		accountInfoResp.SubAccounts = make([]*models.AccountInfoResponse, len(childrenAccounts))

		for i := 0; i < len(childrenAccounts); i++ {
			accountInfoResp.SubAccounts[i] = childrenAccounts[i].ToAccountInfoResponse()
		}
	}

	return accountInfoResp, nil
}

// AccountModifyHandler saves an existed account by request parameters for current user
func (a *AccountsApi) AccountModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var accountModifyReq models.AccountModifyRequest
	err := c.ShouldBindJSON(&accountModifyReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if accountModifyReq.Category < models.ACCOUNT_CATEGORY_CASH || accountModifyReq.Category > models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT {
		log.Warnf(c, "[accounts.AccountModifyHandler] account category invalid, category is %d", accountModifyReq.Category)
		return nil, errs.ErrAccountCategoryInvalid
	}

	uid := c.GetCurrentUid()
	accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.AccountModifyHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", accountModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	accountMap := a.accounts.GetAccountMapByList(accountAndSubAccounts)
	mainAccount, exists := accountMap[accountModifyReq.Id]

	if !exists {
		return nil, errs.ErrAccountNotFound
	}

	if len(accountModifyReq.SubAccounts)+1 != len(accountAndSubAccounts) {
		return nil, errs.ErrCannotAddOrDeleteSubAccountsWhenModify
	}

	if mainAccount.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
		if accountModifyReq.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD && accountModifyReq.CreditCardStatementDate != 0 {
			log.Warnf(c, "[accounts.AccountModifyHandler] cannot set statement date with category \"%d\"", accountModifyReq.Category)
			return nil, errs.ErrCannotSetStatementDateForNonCreditCard
		}
	} else if mainAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		for i := 0; i < len(accountModifyReq.SubAccounts); i++ {
			subAccount := accountModifyReq.SubAccounts[i]

			if subAccount.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD && subAccount.CreditCardStatementDate != 0 {
				log.Warnf(c, "[accounts.AccountModifyHandler] sub-account#%d cannot set statement date with category \"%d\"", i, subAccount.Category)
				return nil, errs.ErrCannotSetStatementDateForNonCreditCard
			}
		}
	}

	anythingUpdate := false
	var toUpdateAccounts []*models.Account

	toUpdateAccount := a.getToUpdateAccount(uid, &accountModifyReq, mainAccount)

	if toUpdateAccount != nil {
		anythingUpdate = true
		toUpdateAccounts = append(toUpdateAccounts, toUpdateAccount)
	}

	for i := 0; i < len(accountModifyReq.SubAccounts); i++ {
		subAccountReq := accountModifyReq.SubAccounts[i]

		if _, exists := accountMap[subAccountReq.Id]; !exists {
			return nil, errs.ErrAccountNotFound
		}

		toUpdateSubAccount := a.getToUpdateAccount(uid, subAccountReq, accountMap[subAccountReq.Id])

		if toUpdateSubAccount != nil {
			anythingUpdate = true
			toUpdateAccounts = append(toUpdateAccounts, toUpdateSubAccount)
		}
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.accounts.ModifyAccounts(c, uid, toUpdateAccounts)

	if err != nil {
		log.Errorf(c, "[accounts.AccountModifyHandler] failed to update account \"id:%d\" for user \"uid:%d\", because %s", accountModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountModifyHandler] user \"uid:%d\" has updated account \"id:%d\" successfully", uid, accountModifyReq.Id)

	accountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(toUpdateAccounts); i++ {
		account := toUpdateAccounts[i]
		oldAccount := accountMap[account.AccountId]

		account.Type = oldAccount.Type
		account.ParentAccountId = oldAccount.ParentAccountId
		account.DisplayOrder = oldAccount.DisplayOrder
		account.Currency = oldAccount.Currency
		account.Balance = oldAccount.Balance

		accountResp := account.ToAccountInfoResponse()
		accountRespMap[accountResp.Id] = accountResp
	}

	for i := 0; i < len(accountAndSubAccounts); i++ {
		oldAccount := accountAndSubAccounts[i]
		_, exists := accountRespMap[oldAccount.AccountId]

		if !exists {
			oldAccountResp := oldAccount.ToAccountInfoResponse()
			accountRespMap[oldAccountResp.Id] = oldAccountResp
		}
	}

	accountResp := accountRespMap[accountModifyReq.Id]

	for i := 0; i < len(accountAndSubAccounts); i++ {
		if accountAndSubAccounts[i].ParentAccountId == accountResp.Id {
			subAccountResp := accountRespMap[accountAndSubAccounts[i].AccountId]
			accountResp.SubAccounts = append(accountResp.SubAccounts, subAccountResp)
		}
	}

	sort.Sort(accountResp.SubAccounts)

	return accountResp, nil
}

// AccountHideHandler hides an existed account by request parameters for current user
func (a *AccountsApi) AccountHideHandler(c *core.WebContext) (any, *errs.Error) {
	var accountHideReq models.AccountHideRequest
	err := c.ShouldBindJSON(&accountHideReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.HideAccount(c, uid, []int64{accountHideReq.Id}, accountHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[accounts.AccountHideHandler] failed to hide account \"id:%d\" for user \"uid:%d\", because %s", accountHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountHideHandler] user \"uid:%d\" has hidden account \"id:%d\"", uid, accountHideReq.Id)
	return true, nil
}

// AccountMoveHandler moves display order of existed accounts by request parameters for current user
func (a *AccountsApi) AccountMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var accountMoveReq models.AccountMoveRequest
	err := c.ShouldBindJSON(&accountMoveReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accounts := make([]*models.Account, len(accountMoveReq.NewDisplayOrders))

	for i := 0; i < len(accountMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := accountMoveReq.NewDisplayOrders[i]
		account := &models.Account{
			Uid:          uid,
			AccountId:    newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		accounts[i] = account
	}

	err = a.accounts.ModifyAccountDisplayOrders(c, uid, accounts)

	if err != nil {
		log.Errorf(c, "[accounts.AccountMoveHandler] failed to move accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountMoveHandler] user \"uid:%d\" has moved accounts", uid)
	return true, nil
}

// AccountDeleteHandler deletes an existed account by request parameters for current user
func (a *AccountsApi) AccountDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var accountDeleteReq models.AccountDeleteRequest
	err := c.ShouldBindJSON(&accountDeleteReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.DeleteAccount(c, uid, accountDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.AccountDeleteHandler] failed to delete account \"id:%d\" for user \"uid:%d\", because %s", accountDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountDeleteHandler] user \"uid:%d\" has deleted account \"id:%d\"", uid, accountDeleteReq.Id)
	return true, nil
}

func (a *AccountsApi) createNewAccountModel(uid int64, accountCreateReq *models.AccountCreateRequest, order int32) *models.Account {
	accountExtend := &models.AccountExtend{}

	if accountCreateReq.Category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
		accountExtend.CreditCardStatementDate = &accountCreateReq.CreditCardStatementDate
	}

	return &models.Account{
		Uid:          uid,
		Name:         accountCreateReq.Name,
		DisplayOrder: order,
		Category:     accountCreateReq.Category,
		Type:         accountCreateReq.Type,
		Icon:         accountCreateReq.Icon,
		Color:        accountCreateReq.Color,
		Currency:     accountCreateReq.Currency,
		Balance:      accountCreateReq.Balance,
		Comment:      accountCreateReq.Comment,
		Extend:       accountExtend,
	}
}

func (a *AccountsApi) createSubAccountModels(uid int64, accountCreateReq *models.AccountCreateRequest) ([]*models.Account, []int64) {
	if len(accountCreateReq.SubAccounts) <= 0 {
		return nil, nil
	}

	childrenAccounts := make([]*models.Account, len(accountCreateReq.SubAccounts))
	childrenAccountBalanceTimes := make([]int64, len(accountCreateReq.SubAccounts))

	for i := int32(0); i < int32(len(accountCreateReq.SubAccounts)); i++ {
		childrenAccounts[i] = a.createNewAccountModel(uid, accountCreateReq.SubAccounts[i], i+1)
		childrenAccountBalanceTimes[i] = accountCreateReq.SubAccounts[i].BalanceTime
	}

	return childrenAccounts, childrenAccountBalanceTimes
}

func (a *AccountsApi) getToUpdateAccount(uid int64, accountModifyReq *models.AccountModifyRequest, oldAccount *models.Account) *models.Account {
	newAccountExtend := &models.AccountExtend{}

	if accountModifyReq.Category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
		newAccountExtend.CreditCardStatementDate = &accountModifyReq.CreditCardStatementDate
	}

	newAccount := &models.Account{
		AccountId: oldAccount.AccountId,
		Uid:       uid,
		Name:      accountModifyReq.Name,
		Category:  accountModifyReq.Category,
		Icon:      accountModifyReq.Icon,
		Color:     accountModifyReq.Color,
		Comment:   accountModifyReq.Comment,
		Extend:    newAccountExtend,
		Hidden:    accountModifyReq.Hidden,
	}

	if newAccount.Name != oldAccount.Name ||
		newAccount.Category != oldAccount.Category ||
		newAccount.Icon != oldAccount.Icon ||
		newAccount.Color != oldAccount.Color ||
		newAccount.Comment != oldAccount.Comment ||
		newAccount.Hidden != oldAccount.Hidden {
		return newAccount
	}

	if (newAccount.Extend != nil && oldAccount.Extend == nil) ||
		(newAccount.Extend == nil && oldAccount.Extend != nil) {
		return newAccount
	}

	oldAccountExtend := oldAccount.Extend

	if newAccountExtend.CreditCardStatementDate != oldAccountExtend.CreditCardStatementDate {
		return newAccount
	}

	return nil
}
