package api

import (
	"sort"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/validators"
)

type AccountsApi struct {
	accounts *services.AccountService
}

var (
	Accounts = &AccountsApi{
		accounts: services.Accounts,
	}
)

func (a *AccountsApi) AccountListHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountListReq models.AccountListRequest
	err := c.ShouldBindQuery(&accountListReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accounts, err := a.accounts.GetAllAccountsByUid(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountListHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
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

		if userAccountResp.ParentId <= models.ACCOUNT_PARENT_ID_LEVEL_ONE {
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
		if userAllAccountResps[i].ParentId == models.ACCOUNT_PARENT_ID_LEVEL_ONE && (!accountListReq.VisibleOnly || !userAllAccountResps[i].Hidden) {
			sort.Sort(userAllAccountResps[i].SubAccounts)
			userFinalAccountResps = append(userFinalAccountResps, userAllAccountResps[i])
		}
	}

	sort.Sort(userFinalAccountResps)

	return userFinalAccountResps, nil
}

func (a *AccountsApi) AccountGetHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountGetReq models.AccountGetRequest
	err := c.ShouldBindQuery(&accountGetReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(uid, accountGetReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountGetHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", accountGetReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	accountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(accountAndSubAccounts); i++ {
		acccountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
		accountRespMap[acccountResp.Id] = acccountResp
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

func (a *AccountsApi) AccountCreateHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountCreateReq models.AccountCreateRequest
	err := c.ShouldBindJSON(&accountCreateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if accountCreateReq.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
		if len(accountCreateReq.SubAccounts) > 0 {
			log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] account cannot have any sub accounts")
			return nil, errs.ErrAccountCannotHaveSubAccounts
		}
	} else if accountCreateReq.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		if len(accountCreateReq.SubAccounts) < 1 {
			log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] account does not have any sub accounts")
			return nil, errs.ErrAccountHaveNoSubAccount
		}

		if accountCreateReq.Currency != validators.PARENT_ACCOUNT_CURRENCY_PLACEHODLER {
			log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] parent account cannot set currency")
			return nil, errs.ErrParentAccountCannotSetCurrency
		}

		if accountCreateReq.Balance != 0 {
			log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] parent account cannot set balance")
			return nil, errs.ErrParentAccountCannotSetBalance
		}

		for i := 0; i < len(accountCreateReq.SubAccounts); i++ {
			subAccount := accountCreateReq.SubAccounts[i]

			if subAccount.Category != accountCreateReq.Category {
				log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] category of sub account not equals to parent")
				return nil, errs.ErrSubAccountCategoryNotEqualsToParent
			}

			if subAccount.Type != models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
				log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] sub account type invalid")
				return nil, errs.ErrSubAccountTypeInvalid
			}
		}
	} else {
		log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] account type invalid, type is %d", accountCreateReq.Type)
		return nil, errs.ErrAccountTypeInvalid
	}

	uid := c.GetCurrentUid()
	maxOrderId, err := a.accounts.GetMaxDisplayOrder(uid, accountCreateReq.Category)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	mainAccount := a.createNewAccountModel(uid, &accountCreateReq, maxOrderId+1)
	childrenAccounts := a.createSubAccountModels(uid, &accountCreateReq)

	err = a.accounts.CreateAccounts(mainAccount, childrenAccounts)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountCreateHandler] failed to create account \"id:%d\" for user \"uid:%d\", because %s", mainAccount.AccountId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[accounts.AccountCreateHandler] user \"uid:%d\" has created a new account \"id:%d\" successfully", uid, mainAccount.AccountId)

	accountInfoResp := mainAccount.ToAccountInfoResponse()

	if len(childrenAccounts) > 0 {
		accountInfoResp.SubAccounts = make([]*models.AccountInfoResponse, len(childrenAccounts))

		for i := 0; i < len(childrenAccounts); i++ {
			accountInfoResp.SubAccounts[i] = childrenAccounts[i].ToAccountInfoResponse()
		}
	}

	return accountInfoResp, nil
}

func (a *AccountsApi) AccountModifyHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountModifyReq models.AccountModifyRequest
	err := c.ShouldBindJSON(&accountModifyReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(uid, accountModifyReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountModifyHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", accountModifyReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	accountMap := make(map[int64]*models.Account)

	for i := 0; i < len(accountAndSubAccounts); i++ {
		acccount := accountAndSubAccounts[i]
		accountMap[acccount.AccountId] = acccount
	}

	if _, exists := accountMap[accountModifyReq.Id]; !exists {
		return nil, errs.ErrAccountNotFound
	}

	if len(accountModifyReq.SubAccounts)+1 != len(accountAndSubAccounts) {
		return nil, errs.ErrCannotAddOrDeleteSubAccountsWhenModify
	}

	anythingUpdate := false
	var toUpdateAccounts []*models.Account

	toUpdateAccount := a.getToUpdateAccount(uid, &accountModifyReq, accountMap[accountModifyReq.Id])

	if toUpdateAccount != nil {
		anythingUpdate = true
		toUpdateAccounts = append(toUpdateAccounts, toUpdateAccount)
	}

	for i := 0; i < len(accountModifyReq.SubAccounts); i++ {
		subAcccountReq := accountModifyReq.SubAccounts[i]

		if _, exists := accountMap[subAcccountReq.Id]; !exists {
			return nil, errs.ErrAccountNotFound
		}

		toUpdateSubAccount := a.getToUpdateAccount(uid, subAcccountReq, accountMap[subAcccountReq.Id])

		if toUpdateSubAccount != nil {
			anythingUpdate = true
			toUpdateAccounts = append(toUpdateAccounts, toUpdateSubAccount)
		}
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.accounts.ModifyAccounts(uid, toUpdateAccounts)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountModifyHandler] failed to update account \"id:%d\" for user \"uid:%d\", because %s", accountModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[accounts.AccountModifyHandler] user \"uid:%d\" has updated account \"id:%d\" successfully", uid, accountModifyReq.Id)

	return true, nil
}

func (a *AccountsApi) AccountHideHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountHideReq models.AccountHideRequest
	err := c.ShouldBindJSON(&accountHideReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.HideAccount(uid, []int64{accountHideReq.Id}, accountHideReq.Hidden)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountHideHandler] failed to hide account \"id:%d\" for user \"uid:%d\", because %s", accountHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[accounts.AccountHideHandler] user \"uid:%d\" has hidden account \"id:%d\"", uid, accountHideReq.Id)
	return true, nil
}

func (a *AccountsApi) AccountMoveHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountMoveReq models.AccountMoveRequest
	err := c.ShouldBindJSON(&accountMoveReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountMoveHandler] parse request failed, because %s", err.Error())
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

	err = a.accounts.ModifyAccountDisplayOrders(uid, accounts)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountMoveHandler] failed to move accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[accounts.AccountMoveHandler] user \"uid:%d\" has moved accounts", uid)
	return true, nil
}

func (a *AccountsApi) AccountDeleteHandler(c *core.Context) (interface{}, *errs.Error) {
	var accountDeleteReq models.AccountDeleteRequest
	err := c.ShouldBindJSON(&accountDeleteReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[accounts.AccountDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.DeleteAccount(uid, accountDeleteReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountDeleteHandler] failed to delete account \"id:%d\" for user \"uid:%d\", because %s", accountDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[accounts.AccountDeleteHandler] user \"uid:%d\" has deleted account \"id:%d\"", uid, accountDeleteReq.Id)
	return true, nil
}

func (a *AccountsApi) createNewAccountModel(uid int64, accountCreateReq *models.AccountCreateRequest, order int) *models.Account {
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
	}
}

func (a *AccountsApi) createSubAccountModels(uid int64, accountCreateReq *models.AccountCreateRequest) []*models.Account {
	if len(accountCreateReq.SubAccounts) <= 0 {
		return nil
	}

	childrenAccounts := make([]*models.Account, len(accountCreateReq.SubAccounts))

	for i := 0; i < len(accountCreateReq.SubAccounts); i++ {
		childrenAccounts[i] = a.createNewAccountModel(uid, accountCreateReq.SubAccounts[i], i+1)
	}

	return childrenAccounts
}

func (a *AccountsApi) getToUpdateAccount(uid int64, accountModifyReq *models.AccountModifyRequest, oldAccount *models.Account) *models.Account {
	newAccount := &models.Account{
		AccountId: oldAccount.AccountId,
		Uid:       uid,
		Name:      accountModifyReq.Name,
		Category:  accountModifyReq.Category,
		Icon:      accountModifyReq.Icon,
		Color:     accountModifyReq.Color,
		Comment:   accountModifyReq.Comment,
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

	return nil
}
