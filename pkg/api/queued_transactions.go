package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// QueuedTransactionsApi represents the queued transaction (draft queue) api
type QueuedTransactionsApi struct {
	ApiUsingConfig
	queuedTransactions *services.QueuedTransactionService
}

// Initialize a queued transaction api singleton instance
var (
	QueuedTransactions = &QueuedTransactionsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		queuedTransactions: services.QueuedTransactions,
	}
)

// QueuedTransactionListHandler returns the draft queue of the current user
func (a *QueuedTransactionsApi) QueuedTransactionListHandler(c *core.WebContext) (any, *errs.Error) {
	var listReq models.QueuedTransactionListRequest
	err := c.ShouldBindQuery(&listReq)

	if err != nil {
		log.Warnf(c, "[queued_transactions.QueuedTransactionListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	queuedTransactions, err := a.queuedTransactions.GetAllQueuedTransactionsByUid(c, uid, listReq.IncludeArchived)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionListHandler] failed to get queued transactions for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	queuedTransactionResps := make(models.QueuedTransactionInfoResponseSlice, len(queuedTransactions))

	for i := 0; i < len(queuedTransactions); i++ {
		queuedTransactionResps[i] = queuedTransactions[i].ToQueuedTransactionInfoResponse()
	}

	sort.Sort(queuedTransactionResps)

	return queuedTransactionResps, nil
}

// QueuedTransactionCountHandler returns how many items the draft queue of the current user holds
func (a *QueuedTransactionsApi) QueuedTransactionCountHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	pendingCount, archivedCount, err := a.queuedTransactions.GetQueuedTransactionCountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionCountHandler] failed to get queued transaction counts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return &models.QueuedTransactionCountResponse{
		PendingCount:  pendingCount,
		ArchivedCount: archivedCount,
	}, nil
}

// QueuedTransactionGetHandler returns one queued transaction of the current user
func (a *QueuedTransactionsApi) QueuedTransactionGetHandler(c *core.WebContext) (any, *errs.Error) {
	var getReq models.QueuedTransactionGetRequest
	err := c.ShouldBindQuery(&getReq)

	if err != nil {
		log.Warnf(c, "[queued_transactions.QueuedTransactionGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	queuedTransaction, err := a.queuedTransactions.GetQueuedTransactionById(c, uid, getReq.Id)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionGetHandler] failed to get queued transaction \"id:%d\" for user \"uid:%d\", because %s", getReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return queuedTransaction.ToQueuedTransactionInfoResponse(), nil
}

// QueuedTransactionCreateHandler pushes a new draft into the queue of the current user
func (a *QueuedTransactionsApi) QueuedTransactionCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var createReq models.QueuedTransactionCreateRequest
	err := c.ShouldBindJSON(&createReq)

	if err != nil {
		log.Warnf(c, "[queued_transactions.QueuedTransactionCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	queuedTransaction, tagIds, errResp := a.buildQueuedTransactionFromCreateRequest(c, uid, &createReq, models.QUEUED_TRANSACTION_SOURCE_API, 0)

	if errResp != nil {
		return nil, errResp
	}

	err = a.queuedTransactions.CreateQueuedTransaction(c, queuedTransaction, tagIds)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionCreateHandler] failed to create queued transaction for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[queued_transactions.QueuedTransactionCreateHandler] user \"uid:%d\" has queued a new draft transaction \"id:%d\" successfully", uid, queuedTransaction.QueuedTransactionId)

	return queuedTransaction.ToQueuedTransactionInfoResponse(), nil
}

// QueuedTransactionArchiveHandler rejects one queued transaction of the current user
func (a *QueuedTransactionsApi) QueuedTransactionArchiveHandler(c *core.WebContext) (any, *errs.Error) {
	var archiveReq models.QueuedTransactionArchiveRequest
	err := c.ShouldBindJSON(&archiveReq)

	if err != nil {
		log.Warnf(c, "[queued_transactions.QueuedTransactionArchiveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.queuedTransactions.ArchiveQueuedTransaction(c, uid, archiveReq.Id)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionArchiveHandler] failed to archive queued transaction \"id:%d\" for user \"uid:%d\", because %s", archiveReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[queued_transactions.QueuedTransactionArchiveHandler] user \"uid:%d\" has archived queued transaction \"id:%d\"", uid, archiveReq.Id)

	return true, nil
}

// QueuedTransactionArchiveAllHandler rejects every pending queued transaction of the current user
func (a *QueuedTransactionsApi) QueuedTransactionArchiveAllHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	affectedCount, err := a.queuedTransactions.ArchiveAllPendingQueuedTransactions(c, uid)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionArchiveAllHandler] failed to archive all queued transactions for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[queued_transactions.QueuedTransactionArchiveAllHandler] user \"uid:%d\" has archived %d queued transactions", uid, affectedCount)

	return &models.QueuedTransactionBatchOperationResponse{
		AffectedCount: affectedCount,
	}, nil
}

// QueuedTransactionRestoreHandler moves an archived queued transaction of the current user back to the queue
func (a *QueuedTransactionsApi) QueuedTransactionRestoreHandler(c *core.WebContext) (any, *errs.Error) {
	var restoreReq models.QueuedTransactionRestoreRequest
	err := c.ShouldBindJSON(&restoreReq)

	if err != nil {
		log.Warnf(c, "[queued_transactions.QueuedTransactionRestoreHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.queuedTransactions.RestoreQueuedTransaction(c, uid, restoreReq.Id)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionRestoreHandler] failed to restore queued transaction \"id:%d\" for user \"uid:%d\", because %s", restoreReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[queued_transactions.QueuedTransactionRestoreHandler] user \"uid:%d\" has restored queued transaction \"id:%d\"", uid, restoreReq.Id)

	return true, nil
}

// QueuedTransactionDeleteHandler removes one queued transaction of the current user from the queue,
// which is what the client does after the real transaction has been created out of the draft
func (a *QueuedTransactionsApi) QueuedTransactionDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var deleteReq models.QueuedTransactionDeleteRequest
	err := c.ShouldBindJSON(&deleteReq)

	if err != nil {
		log.Warnf(c, "[queued_transactions.QueuedTransactionDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.queuedTransactions.DeleteQueuedTransaction(c, uid, deleteReq.Id)

	if err != nil {
		log.Errorf(c, "[queued_transactions.QueuedTransactionDeleteHandler] failed to delete queued transaction \"id:%d\" for user \"uid:%d\", because %s", deleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[queued_transactions.QueuedTransactionDeleteHandler] user \"uid:%d\" has deleted queued transaction \"id:%d\"", uid, deleteReq.Id)

	return true, nil
}

func (a *QueuedTransactionsApi) buildQueuedTransactionFromCreateRequest(c *core.WebContext, uid int64, createReq *models.QueuedTransactionCreateRequest, source models.QueuedTransactionSource, sourceTemplateId int64) (*models.QueuedTransaction, []int64, *errs.Error) {
	tagIds, err := utils.StringArrayToInt64Array(createReq.TagIds)

	if err != nil {
		log.Warnf(c, "[queued_transactions.buildQueuedTransactionFromCreateRequest] parse tag ids failed, because %s", err.Error())
		return nil, nil, errs.ErrTransactionTagIdInvalid
	}

	if len(tagIds) > models.MaximumTagsCountOfTransaction {
		return nil, nil, errs.ErrQueuedTransactionHasTooManyTags
	}

	if createReq.Type != 0 && (createReq.Type < models.TRANSACTION_TYPE_MODIFY_BALANCE || createReq.Type > models.TRANSACTION_TYPE_TRANSFER) {
		log.Warnf(c, "[queued_transactions.buildQueuedTransactionFromCreateRequest] transaction type is invalid")
		return nil, nil, errs.ErrTransactionTypeInvalid
	}

	queuedTransaction := &models.QueuedTransaction{
		Uid:               uid,
		Source:            source,
		SourceTemplateId:  sourceTemplateId,
		Type:              createReq.Type,
		CategoryId:        createReq.CategoryId,
		AccountId:         createReq.SourceAccountId,
		TimezoneUtcOffset: createReq.UtcOffset,
		Amount:            createReq.SourceAmount,
		HideAmount:        createReq.HideAmount,
		Comment:           createReq.Comment,
		CreatedIp:         c.ClientIP(),
	}

	if createReq.Time > 0 {
		queuedTransaction.TransactionTime = utils.GetMinTransactionTimeFromUnixTime(createReq.Time)
	}

	if createReq.Type == models.TRANSACTION_TYPE_TRANSFER {
		queuedTransaction.RelatedAccountId = createReq.DestinationAccountId
		queuedTransaction.RelatedAccountAmount = createReq.DestinationAmount
	}

	if createReq.GeoLocation != nil {
		queuedTransaction.GeoLongitude = createReq.GeoLocation.Longitude
		queuedTransaction.GeoLatitude = createReq.GeoLocation.Latitude
	}

	return queuedTransaction, tagIds, nil
}
