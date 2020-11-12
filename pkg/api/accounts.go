package api

import (
	"sort"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
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

		if userAccountResp.ParentId <= models.ACCOUNT_PARENT_ID_LEVEL_ONE {
			continue
		}

		parentAccount, parentExists := userAllAccountRespMap[userAccountResp.ParentId]

		if !parentExists || parentAccount == nil {
			continue
		}

		parentAccount.SubAccounts = append(parentAccount.SubAccounts, userAccountResp)
	}

	userFinalAccountResps := make(models.AccountInfoResponseSlice, 0)

	for i := 0; i < len(userAllAccountResps); i++ {
		if userAllAccountResps[i].ParentId == models.ACCOUNT_PARENT_ID_LEVEL_ONE {
			sort.Sort(userAllAccountResps[i].SubAccounts)
			userFinalAccountResps = append(userFinalAccountResps, userAllAccountResps[i])
		}
	}

	sort.Sort(userFinalAccountResps)

	return userFinalAccountResps, nil
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
	} else {
		log.WarnfWithRequestId(c, "[accounts.AccountCreateHandler] account type invalid, type is %d", accountCreateReq.Type)
		return nil, errs.ErrAccountTypeInvalid
	}

	uid := c.GetCurrentUid()
	maxOrderId, err := a.accounts.GetMaxDisplayOrder(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	mainAccount := a.createNewAccount(uid, &accountCreateReq, maxOrderId+1)
	childrenAccounts := a.createSubAccounts(uid, &accountCreateReq)

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
	err = a.accounts.DeleteAccounts(uid, []int64{accountDeleteReq.Id})

	if err != nil {
		log.ErrorfWithRequestId(c, "[accounts.AccountDeleteHandler] failed to delete account \"id:%s\" for user \"uid:%d\", because %s", accountDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[accounts.AccountDeleteHandler] user \"uid:%d\" has deleted account \"id:%s\"", uid, accountDeleteReq.Id)
	return true, nil
}

func (a *AccountsApi) createNewAccount(uid int64, accountCreateReq *models.AccountCreateRequest, order int) *models.Account {
	return &models.Account{
		Uid:          uid,
		Name:         accountCreateReq.Name,
		DisplayOrder: order,
		Category:     accountCreateReq.Category,
		Icon:         accountCreateReq.Icon,
		Currency:     accountCreateReq.Currency,
		Comment:      accountCreateReq.Comment,
	}
}

func (a *AccountsApi) createSubAccounts(uid int64, accountCreateReq *models.AccountCreateRequest) []*models.Account {
	if len(accountCreateReq.SubAccounts) <= 0 {
		return nil
	}

	childrenAccounts := make([]*models.Account, len(accountCreateReq.SubAccounts))

	for i := 0; i < len(accountCreateReq.SubAccounts); i++ {
		childrenAccounts[i] = a.createNewAccount(uid, accountCreateReq.SubAccounts[i], i+1)
	}

	return childrenAccounts
}
