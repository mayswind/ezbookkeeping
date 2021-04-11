package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// TransactionTagsApi represents transaction tag api
type TransactionTagsApi struct {
	tags *services.TransactionTagService
}

// Initialize a transaction tag api singleton instance
var (
	TransactionTags = &TransactionTagsApi{
		tags: services.TransactionTags,
	}
)

// TagListHandler returns transaction tag list of current user
func (a *TransactionTagsApi) TagListHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	tags, err := a.tags.GetAllTagsByUid(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagListHandler] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	tagResps := make(models.TransactionTagInfoResponseSlice, len(tags))

	for i := 0; i < len(tags); i++ {
		tagResps[i] = tags[i].ToTransactionTagInfoResponse()
	}

	sort.Sort(tagResps)

	return tagResps, nil
}

// TagGetHandler returns one specific transaction tag of current user
func (a *TransactionTagsApi) TagGetHandler(c *core.Context) (interface{}, *errs.Error) {
	var tagGetReq models.TransactionTagGetRequest
	err := c.ShouldBindQuery(&tagGetReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_tags.TagGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tag, err := a.tags.GetTagByTagId(uid, tagGetReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagGetHandler] failed to get tag \"id:%d\" for user \"uid:%d\", because %s", tagGetReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	tagResp := tag.ToTransactionTagInfoResponse()

	return tagResp, nil
}

// TagCreateHandler saves a new transaction tag by request parameters for current user
func (a *TransactionTagsApi) TagCreateHandler(c *core.Context) (interface{}, *errs.Error) {
	var tagCreateReq models.TransactionTagCreateRequest
	err := c.ShouldBindJSON(&tagCreateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_tags.TagCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	maxOrderId, err := a.tags.GetMaxDisplayOrder(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	tag := a.createNewTagModel(uid, &tagCreateReq, maxOrderId+1)

	err = a.tags.CreateTag(tag)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagCreateHandler] failed to create tag \"id:%d\" for user \"uid:%d\", because %s", tag.TagId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_tags.TagCreateHandler] user \"uid:%d\" has created a new tag \"id:%d\" successfully", uid, tag.TagId)

	tagResp := tag.ToTransactionTagInfoResponse()

	return tagResp, nil
}

// TagModifyHandler saves an existed transaction tag by request parameters for current user
func (a *TransactionTagsApi) TagModifyHandler(c *core.Context) (interface{}, *errs.Error) {
	var tagModifyReq models.TransactionTagModifyRequest
	err := c.ShouldBindJSON(&tagModifyReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_tags.TagModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tag, err := a.tags.GetTagByTagId(uid, tagModifyReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagModifyHandler] failed to get tag \"id:%d\" for user \"uid:%d\", because %s", tagModifyReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	newTag := &models.TransactionTag{
		TagId: tag.TagId,
		Uid:   uid,
		Name:  tagModifyReq.Name,
	}

	if newTag.Name == tag.Name {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.tags.ModifyTag(newTag)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagModifyHandler] failed to update tag \"id:%d\" for user \"uid:%d\", because %s", tagModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_tags.TagModifyHandler] user \"uid:%d\" has updated tag \"id:%d\" successfully", uid, tagModifyReq.Id)

	tag.Name = newTag.Name
	tagResp := tag.ToTransactionTagInfoResponse()

	return tagResp, nil
}

// TagHideHandler hides an transaction tag by request parameters for current user
func (a *TransactionTagsApi) TagHideHandler(c *core.Context) (interface{}, *errs.Error) {
	var tagHideReq models.TransactionTagHideRequest
	err := c.ShouldBindJSON(&tagHideReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_tags.CategoryHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.tags.HideTag(uid, []int64{tagHideReq.Id}, tagHideReq.Hidden)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.CategoryHideHandler] failed to hide tag \"id:%d\" for user \"uid:%d\", because %s", tagHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_tags.CategoryHideHandler] user \"uid:%d\" has hidden category \"id:%d\"", uid, tagHideReq.Id)
	return true, nil
}

// TagMoveHandler moves display order of existed transaction tags by request parameters for current user
func (a *TransactionTagsApi) TagMoveHandler(c *core.Context) (interface{}, *errs.Error) {
	var tagMoveReq models.TransactionTagMoveRequest
	err := c.ShouldBindJSON(&tagMoveReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_tags.CategoryMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tags := make([]*models.TransactionTag, len(tagMoveReq.NewDisplayOrders))

	for i := 0; i < len(tagMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := tagMoveReq.NewDisplayOrders[i]
		tag := &models.TransactionTag{
			Uid:          uid,
			TagId:        newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		tags[i] = tag
	}

	err = a.tags.ModifyTagDisplayOrders(uid, tags)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.CategoryMoveHandler] failed to move tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_tags.CategoryMoveHandler] user \"uid:%d\" has moved categories", uid)
	return true, nil
}

// TagDeleteHandler deletes an existed transaction tag by request parameters for current user
func (a *TransactionTagsApi) TagDeleteHandler(c *core.Context) (interface{}, *errs.Error) {
	var tagDeleteReq models.TransactionTagDeleteRequest
	err := c.ShouldBindJSON(&tagDeleteReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_tags.TagDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.tags.DeleteTag(uid, tagDeleteReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_tags.TagDeleteHandler] failed to delete tag \"id:%d\" for user \"uid:%d\", because %s", tagDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_tags.TagDeleteHandler] user \"uid:%d\" has deleted tag \"id:%d\"", uid, tagDeleteReq.Id)
	return true, nil
}

func (a *TransactionTagsApi) createNewTagModel(uid int64, tagCreateReq *models.TransactionTagCreateRequest, order int) *models.TransactionTag {
	return &models.TransactionTag{
		Uid:          uid,
		Name:         tagCreateReq.Name,
		DisplayOrder: order,
	}
}
