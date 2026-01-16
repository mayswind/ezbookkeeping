package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// TransactionTagGroupsApi represents transaction tag group api
type TransactionTagGroupsApi struct {
	tagGroups *services.TransactionTagGroupService
}

// Initialize a transaction tag group api singleton instance
var (
	TransactionTagGroups = &TransactionTagGroupsApi{
		tagGroups: services.TransactionTagGroups,
	}
)

// TagGroupListHandler returns transaction tag group list of current user
func (a *TransactionTagGroupsApi) TagGroupListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	tagGroups, err := a.tagGroups.GetAllTagGroupsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupListHandler] failed to get tag groups for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tagGroupResps := make(models.TransactionTagGroupInfoResponseSlice, len(tagGroups))

	for i := 0; i < len(tagGroups); i++ {
		tagGroupResps[i] = tagGroups[i].ToTransactionTagGroupInfoResponse()
	}

	sort.Sort(tagGroupResps)

	return tagGroupResps, nil
}

// TagGroupGetHandler returns one specific transaction tag group of current user
func (a *TransactionTagGroupsApi) TagGroupGetHandler(c *core.WebContext) (any, *errs.Error) {
	var tagGroupGetReq models.TransactionTagGroupGetRequest
	err := c.ShouldBindQuery(&tagGroupGetReq)

	if err != nil {
		log.Warnf(c, "[transaction_tag_groups.TagGroupGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tagGroup, err := a.tagGroups.GetTagGroupByTagGroupId(c, uid, tagGroupGetReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupGetHandler] failed to get tag group \"id:%d\" for user \"uid:%d\", because %s", tagGroupGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tagGroupResp := tagGroup.ToTransactionTagGroupInfoResponse()

	return tagGroupResp, nil
}

// TagGroupCreateHandler saves a new transaction tag group by request parameters for current user
func (a *TransactionTagGroupsApi) TagGroupCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var tagGroupCreateReq models.TransactionTagGroupCreateRequest
	err := c.ShouldBindJSON(&tagGroupCreateReq)

	if err != nil {
		log.Warnf(c, "[transaction_tag_groups.TagGroupCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	maxOrderId, err := a.tagGroups.GetMaxDisplayOrder(c, uid)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tagGroup := a.createNewTagGroupModel(uid, &tagGroupCreateReq, maxOrderId+1)

	err = a.tagGroups.CreateTagGroup(c, tagGroup)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupCreateHandler] failed to create tag group \"id:%d\" for user \"uid:%d\", because %s", tagGroup.TagGroupId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_tag_groups.TagGroupCreateHandler] user \"uid:%d\" has created a new tag group \"id:%d\" successfully", uid, tagGroup.TagGroupId)

	tagGroupResp := tagGroup.ToTransactionTagGroupInfoResponse()

	return tagGroupResp, nil
}

// TagGroupModifyHandler saves an existed transaction tag group by request parameters for current user
func (a *TransactionTagGroupsApi) TagGroupModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var tagGroupModifyReq models.TransactionTagGroupModifyRequest
	err := c.ShouldBindJSON(&tagGroupModifyReq)

	if err != nil {
		log.Warnf(c, "[transaction_tag_groups.TagGroupModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tagGroup, err := a.tagGroups.GetTagGroupByTagGroupId(c, uid, tagGroupModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupModifyHandler] failed to get tag group \"id:%d\" for user \"uid:%d\", because %s", tagGroupModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	newTagGroup := &models.TransactionTagGroup{
		TagGroupId: tagGroup.TagGroupId,
		Uid:        uid,
		Name:       tagGroupModifyReq.Name,
	}

	if newTagGroup.Name == tagGroup.Name {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.tagGroups.ModifyTagGroup(c, newTagGroup)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupModifyHandler] failed to update tag group \"id:%d\" for user \"uid:%d\", because %s", tagGroupModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_tag_groups.TagGroupModifyHandler] user \"uid:%d\" has updated tag group \"id:%d\" successfully", uid, tagGroupModifyReq.Id)

	tagGroup.Name = newTagGroup.Name
	tagGroupResp := tagGroup.ToTransactionTagGroupInfoResponse()

	return tagGroupResp, nil
}

// TagGroupMoveHandler moves display order of existed transaction tag groups by request parameters for current user
func (a *TransactionTagGroupsApi) TagGroupMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var tagGroupMoveReq models.TransactionTagGroupMoveRequest
	err := c.ShouldBindJSON(&tagGroupMoveReq)

	if err != nil {
		log.Warnf(c, "[transaction_tag_groups.TagGroupMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tagGroups := make([]*models.TransactionTagGroup, len(tagGroupMoveReq.NewDisplayOrders))

	for i := 0; i < len(tagGroupMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := tagGroupMoveReq.NewDisplayOrders[i]
		tagGroup := &models.TransactionTagGroup{
			Uid:          uid,
			TagGroupId:   newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		tagGroups[i] = tagGroup
	}

	err = a.tagGroups.ModifyTagGroupDisplayOrders(c, uid, tagGroups)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupMoveHandler] failed to move tag groups for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_tag_groups.TagGroupMoveHandler] user \"uid:%d\" has moved tag groups", uid)
	return true, nil
}

// TagGroupDeleteHandler deletes an existed transaction tag group by request parameters for current user
func (a *TransactionTagGroupsApi) TagGroupDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var tagGroupDeleteReq models.TransactionTagGroupDeleteRequest
	err := c.ShouldBindJSON(&tagGroupDeleteReq)

	if err != nil {
		log.Warnf(c, "[transaction_tag_groups.TagGroupDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.tagGroups.DeleteTagGroup(c, uid, tagGroupDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_tag_groups.TagGroupDeleteHandler] failed to delete tag group \"id:%d\" for user \"uid:%d\", because %s", tagGroupDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_tag_groups.TagGroupDeleteHandler] user \"uid:%d\" has deleted tag group \"id:%d\"", uid, tagGroupDeleteReq.Id)
	return true, nil
}

func (a *TransactionTagGroupsApi) createNewTagGroupModel(uid int64, tagGroupCreateReq *models.TransactionTagGroupCreateRequest, order int32) *models.TransactionTagGroup {
	return &models.TransactionTagGroup{
		Uid:          uid,
		Name:         tagGroupCreateReq.Name,
		DisplayOrder: order,
	}
}
