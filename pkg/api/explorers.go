package api

import (
	"encoding/json"
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// InsightsExplorersApi represents insights explorers api
type InsightsExplorersApi struct {
	insightsExploreres *services.InsightsExplorerService
}

// Initialize a insights explorers api singleton instance
var (
	InsightsExplorers = &InsightsExplorersApi{
		insightsExploreres: services.InsightsExplorers,
	}
)

// InsightsExplorerListHandler returns insights explorer list of current user
func (a *InsightsExplorersApi) InsightsExplorerListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	explorers, err := a.insightsExploreres.GetAllInsightsExplorerNamesByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerListHandler] failed to get insights explorers for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	explorerResps := make(models.InsightsExplorerInfoResponseSlice, len(explorers))

	for i := 0; i < len(explorers); i++ {
		explorerResps[i], err = explorers[i].ToInsightsExplorerInfoResponse()

		if err != nil {
			log.Errorf(c, "[explorers.InsightsExplorerListHandler] failed to get insights explorer response for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrInsightsExplorerDataInvalid
		}
	}

	sort.Sort(explorerResps)

	return explorerResps, nil
}

// InsightsExplorerGetHandler returns one specific insights explorer of current user
func (a *InsightsExplorersApi) InsightsExplorerGetHandler(c *core.WebContext) (any, *errs.Error) {
	var explorerGetReq models.InsightsExplorerGetRequest
	err := c.ShouldBindQuery(&explorerGetReq)

	if err != nil {
		log.Warnf(c, "[explorers.InsightsExplorerGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	explorer, err := a.insightsExploreres.GetInsightsExplorerByExplorerId(c, uid, explorerGetReq.Id)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerGetHandler] failed to get insights explorer \"id:%d\" for user \"uid:%d\", because %s", explorerGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	explorerResp, err := explorer.ToInsightsExplorerInfoResponse()

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerGetHandler] failed to get insights explorer response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrInsightsExplorerDataInvalid
	}

	return explorerResp, nil
}

// InsightsExplorerCreateHandler saves a new insights explorer by request parameters for current user
func (a *InsightsExplorersApi) InsightsExplorerCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var explorerCreateReq models.InsightsExplorerCreateRequest
	err := c.ShouldBindJSON(&explorerCreateReq)

	if err != nil {
		log.Warnf(c, "[explorers.InsightsExplorerCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	maxOrderId, err := a.insightsExploreres.GetMaxDisplayOrder(c, uid)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	explorer, err := a.createNewInsightsExplorerModel(uid, &explorerCreateReq, maxOrderId+1)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerCreateHandler] failed to parse insights explorer data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrInsightsExplorerDataInvalid
	}

	err = a.insightsExploreres.CreateInsightsExplorer(c, explorer)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerCreateHandler] failed to create insights explorer \"id:%d\" for user \"uid:%d\", because %s", explorer.ExplorerId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[explorers.InsightsExplorerCreateHandler] user \"uid:%d\" has created a new insights explorer \"id:%d\" successfully", uid, explorer.ExplorerId)

	explorerResp, err := explorer.ToInsightsExplorerInfoResponse()

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerCreateHandler] failed to get insights explorer response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrInsightsExplorerDataInvalid
	}

	return explorerResp, nil
}

// InsightsExplorerModifyHandler saves an existed insights explorer by request parameters for current user
func (a *InsightsExplorersApi) InsightsExplorerModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var explorerModifyReq models.InsightsExplorerModifyRequest
	err := c.ShouldBindJSON(&explorerModifyReq)

	if err != nil {
		log.Warnf(c, "[explorers.InsightsExplorerModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	explorer, err := a.insightsExploreres.GetInsightsExplorerByExplorerId(c, uid, explorerModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerModifyHandler] failed to get insights explorer \"id:%d\" for user \"uid:%d\", because %s", explorerModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	newData, err := json.Marshal(explorerModifyReq.Data)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerModifyHandler] failed to parse insights explorer data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrInsightsExplorerDataInvalid
	}

	newExplorer := &models.InsightsExplorer{
		ExplorerId: explorer.ExplorerId,
		Uid:        uid,
		Name:       explorerModifyReq.Name,
		Data:       string(newData),
	}

	if newExplorer.Name == explorer.Name && newExplorer.Data == explorer.Data {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.insightsExploreres.ModifyInsightsExplorer(c, newExplorer)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerModifyHandler] failed to update insights explorer \"id:%d\" for user \"uid:%d\", because %s", explorerModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[explorers.InsightsExplorerModifyHandler] user \"uid:%d\" has updated insights explorer \"id:%d\" successfully", uid, explorerModifyReq.Id)

	explorer.Name = newExplorer.Name
	explorer.Data = newExplorer.Data
	explorerResp, err := explorer.ToInsightsExplorerInfoResponse()

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerModifyHandler] failed to get insights explorer response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrInsightsExplorerDataInvalid
	}

	return explorerResp, nil
}

// InsightsExplorerHideHandler hides a insights explorer by request parameters for current user
func (a *InsightsExplorersApi) InsightsExplorerHideHandler(c *core.WebContext) (any, *errs.Error) {
	var explorerHideReq models.InsightsExplorerHideRequest
	err := c.ShouldBindJSON(&explorerHideReq)

	if err != nil {
		log.Warnf(c, "[explorers.InsightsExplorerHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.insightsExploreres.HideInsightsExplorer(c, uid, []int64{explorerHideReq.Id}, explorerHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerHideHandler] failed to hide insights explorer \"id:%d\" for user \"uid:%d\", because %s", explorerHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[explorers.InsightsExplorerHideHandler] user \"uid:%d\" has hidden insights explorer \"id:%d\"", uid, explorerHideReq.Id)
	return true, nil
}

// InsightsExplorerMoveHandler moves display order of existed insights explorers by request parameters for current user
func (a *InsightsExplorersApi) InsightsExplorerMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var explorerMoveReq models.InsightsExplorerMoveRequest
	err := c.ShouldBindJSON(&explorerMoveReq)

	if err != nil {
		log.Warnf(c, "[explorers.InsightsExplorerMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	explorers := make([]*models.InsightsExplorer, len(explorerMoveReq.NewDisplayOrders))

	for i := 0; i < len(explorerMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := explorerMoveReq.NewDisplayOrders[i]
		explorer := &models.InsightsExplorer{
			Uid:          uid,
			ExplorerId:   newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		explorers[i] = explorer
	}

	err = a.insightsExploreres.ModifyInsightsExplorerDisplayOrders(c, uid, explorers)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerMoveHandler] failed to move insights explorers for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[explorers.InsightsExplorerMoveHandler] user \"uid:%d\" has moved insights explorers", uid)
	return true, nil
}

// InsightsExplorerDeleteHandler deletes an existed insights explorer by request parameters for current user
func (a *InsightsExplorersApi) InsightsExplorerDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var explorerDeleteReq models.InsightsExplorerDeleteRequest
	err := c.ShouldBindJSON(&explorerDeleteReq)

	if err != nil {
		log.Warnf(c, "[explorers.InsightsExplorerDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.insightsExploreres.DeleteInsightsExplorer(c, uid, explorerDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[explorers.InsightsExplorerDeleteHandler] failed to delete insights explorer \"id:%d\" for user \"uid:%d\", because %s", explorerDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[explorers.InsightsExplorerDeleteHandler] user \"uid:%d\" has deleted insights explorer \"id:%d\"", uid, explorerDeleteReq.Id)
	return true, nil
}

func (a *InsightsExplorersApi) createNewInsightsExplorerModel(uid int64, explorerCreateReq *models.InsightsExplorerCreateRequest, order int32) (*models.InsightsExplorer, error) {
	data, err := json.Marshal(explorerCreateReq.Data)

	if err != nil {
		return nil, err
	}

	return &models.InsightsExplorer{
		Uid:          uid,
		Name:         explorerCreateReq.Name,
		Data:         string(data),
		DisplayOrder: order,
	}, nil
}
