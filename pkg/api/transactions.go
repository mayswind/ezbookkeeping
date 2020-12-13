package api

import (
	"sort"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/utils"
)

type TransactionsApi struct {
	transactions    *services.TransactionService
	transactionTags *services.TransactionTagService
}

var (
	Transactions = &TransactionsApi{
		transactions:    services.Transactions,
		transactionTags: services.TransactionTags,
	}
)

func (a *TransactionsApi) TransactionListHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionListReq models.TransactionListByMaxTimeRequest
	err := c.ShouldBindQuery(&transactionListReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transactions.TransactionListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	transactions, err := a.transactions.GetTransactionsByMaxTime(uid, transactionListReq.MaxTime, transactionListReq.Count + 1)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionListHandler] failed to get transactions earlier than \"%d\" for user \"uid:%d\", because %s", transactionListReq.MaxTime, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	finalCount := transactionListReq.Count

	if len(transactions) < finalCount {
		finalCount = len(transactions)
	}

	transactionIds := make([]int64, finalCount)

	for i := 0; i < finalCount; i++ {
		transactionIds[i] = transactions[i].TransactionId
	}

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(uid, transactionIds)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionListHandler] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	transactionResps := &models.TransactionInfoPageWrapperResponse{}
	transactionResps.Items = make(models.TransactionInfoResponseSlice, finalCount)

	for i := 0; i < finalCount; i++ {
		transactionTagIds := allTransactionTagIds[transactions[i].TransactionId]
		transactionResps.Items[i] = transactions[i].ToTransactionInfoResponse(transactionTagIds)
	}

	sort.Sort(transactionResps.Items)

	if finalCount < len(transactions) {
		transactionResps.NextTime = &transactions[finalCount].TransactionTime
	}

	return transactionResps, nil
}

func (a *TransactionsApi) TransactionMonthListHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionListReq models.TransactionListInMonthByPageRequest
	err := c.ShouldBindQuery(&transactionListReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transactions.TransactionMonthListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	transactions, err := a.transactions.GetTransactionsInMonthByPage(uid, transactionListReq.Year, transactionListReq.Month, transactionListReq.Page, transactionListReq.Count)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionMonthListHandler] failed to get transactions in month \"%d-%d\" for user \"uid:%d\", because %s", transactionListReq.Year, transactionListReq.Month, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	transactionIds := make([]int64, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transactionIds[i] = transactions[i].TransactionId
	}

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(uid, transactionIds)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionMonthListHandler] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	transactionResps := make([]*models.TransactionInfoResponse, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transactionTagIds := allTransactionTagIds[transactions[i].TransactionId]
		transactionResps[i] = transactions[i].ToTransactionInfoResponse(transactionTagIds)
	}

	return transactionResps, nil
}

func (a *TransactionsApi) TransactionGetHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionGetReq models.TransactionGetRequest
	err := c.ShouldBindQuery(&transactionGetReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transactions.TransactionGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	transaction, err := a.transactions.GetTransactionByTransactionId(uid, transactionGetReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionGetHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", transactionGetReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(uid, []int64{transaction.TransactionId})

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionGetHandler] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	transactionTagIds := allTransactionTagIds[transaction.TransactionId]
	transactionResp := transaction.ToTransactionInfoResponse(transactionTagIds)

	return transactionResp, nil
}

func (a *TransactionsApi) TransactionCreateHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionCreateReq models.TransactionCreateRequest
	err := c.ShouldBindJSON(&transactionCreateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transactions.TransactionCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if transactionCreateReq.Type < models.TRANSACTION_TYPE_MODIFY_BALANCE || transactionCreateReq.Type > models.TRANSACTION_TYPE_TRANSFER {
		log.WarnfWithRequestId(c, "[transactions.TransactionCreateHandler] transaction type is invalid")
		return nil, errs.ErrTransactionTypeInvalid
	}

	if transactionCreateReq.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE && transactionCreateReq.CategoryId > 0 {
		log.WarnfWithRequestId(c, "[transactions.TransactionCreateHandler] balance modification transaction cannot set category id")
		return nil, errs.ErrBalanceModificationTransactionCannotSetCategory
	}

	if transactionCreateReq.Type != models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.SourceAccountId != transactionCreateReq.DestinationAccountId {
		log.WarnfWithRequestId(c, "[transactions.TransactionCreateHandler] non-transfer transaction source account is not destination account")
		return nil, errs.ErrTransactionSourceAndDestinationIdNotEqual
	} else if transactionCreateReq.Type == models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.SourceAccountId == transactionCreateReq.DestinationAccountId {
		log.WarnfWithRequestId(c, "[transactions.TransactionCreateHandler] transfer transaction source account must not be destination account")
		return nil, errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
	}

	if transactionCreateReq.Type != models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.SourceAmount != transactionCreateReq.DestinationAmount {
		log.WarnfWithRequestId(c, "[transactions.TransactionCreateHandler] non-transfer transaction source amount is not destination amount")
		return nil, errs.ErrTransactionSourceAndDestinationAmountNotEqual
	}

	uid := c.GetCurrentUid()
	transaction := a.createNewTransactionModel(uid, &transactionCreateReq)

	err = a.transactions.CreateTransaction(transaction, transactionCreateReq.TagIds)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionCreateHandler] failed to create transaction \"id:%d\" for user \"uid:%d\", because %s", transaction.TransactionId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transactions.TransactionCreateHandler] user \"uid:%d\" has created a new transaction \"id:%d\" successfully", uid, transaction.TransactionId)

	transactionResp := transaction.ToTransactionInfoResponse(nil)

	return transactionResp, nil
}

func (a *TransactionsApi) TransactionModifyHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionModifyReq models.TransactionModifyRequest
	err := c.ShouldBindJSON(&transactionModifyReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transactions.TransactionModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	transaction, err := a.transactions.GetTransactionByTransactionId(uid, transactionModifyReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionModifyHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", transactionModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(uid, []int64{transaction.TransactionId})

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionModifyHandler] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	transactionTagIds := allTransactionTagIds[transaction.TransactionId]
	addTransactionTagIds := utils.Int64SliceMinus(transactionModifyReq.TagIds, transactionTagIds)
	removeTransactionTagIds := utils.Int64SliceMinus(transactionTagIds, transactionModifyReq.TagIds)

	newTransaction := &models.Transaction{
		TransactionId:        transaction.TransactionId,
		Uid:                  uid,
		CategoryId:           transactionModifyReq.CategoryId,
		TransactionTime:      transactionModifyReq.Time,
		SourceAccountId:      transactionModifyReq.SourceAccountId,
		DestinationAccountId: transactionModifyReq.DestinationAccountId,
		SourceAmount:         transactionModifyReq.SourceAmount,
		DestinationAmount:    transactionModifyReq.DestinationAmount,
		Comment:              transactionModifyReq.Comment,
	}

	if newTransaction.CategoryId == transaction.CategoryId &&
		newTransaction.TransactionTime / 1000 == transaction.TransactionTime / 1000 &&
		newTransaction.SourceAccountId == transaction.SourceAccountId &&
		newTransaction.DestinationAccountId == transaction.DestinationAccountId &&
		newTransaction.SourceAmount == transaction.SourceAmount &&
		newTransaction.DestinationAmount == transaction.DestinationAmount &&
		newTransaction.Comment == transaction.Comment &&
		len(addTransactionTagIds) < 1 &&
		len(removeTransactionTagIds) < 1 {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.transactions.ModifyTransaction(newTransaction, addTransactionTagIds, removeTransactionTagIds)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionModifyHandler] failed to update transaction \"id:%d\" for user \"uid:%d\", because %s", transactionModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transactions.TransactionModifyHandler] user \"uid:%d\" has updated transaction \"id:%d\" successfully", uid, transactionModifyReq.Id)

	return true, nil
}

func (a *TransactionsApi) TransactionDeleteHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionDeleteReq models.TransactionDeleteRequest
	err := c.ShouldBindJSON(&transactionDeleteReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transactions.TransactionDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.transactions.DeleteTransaction(uid, transactionDeleteReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transactions.TransactionDeleteHandler] failed to delete transaction \"id:%d\" for user \"uid:%d\", because %s", transactionDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transactions.TransactionDeleteHandler] user \"uid:%d\" has deleted transaction \"id:%d\"", uid, transactionDeleteReq.Id)
	return true, nil
}

func (a *TransactionsApi) createNewTransactionModel(uid int64, transactionCreateReq *models.TransactionCreateRequest) *models.Transaction {
	return &models.Transaction{
		Uid:                  uid,
		Type:                 transactionCreateReq.Type,
		CategoryId:           transactionCreateReq.CategoryId,
		TransactionTime:      transactionCreateReq.Time,
		SourceAccountId:      transactionCreateReq.SourceAccountId,
		DestinationAccountId: transactionCreateReq.DestinationAccountId,
		SourceAmount:         transactionCreateReq.SourceAmount,
		DestinationAmount:    transactionCreateReq.DestinationAmount,
		Comment:              transactionCreateReq.Comment,
	}
}
