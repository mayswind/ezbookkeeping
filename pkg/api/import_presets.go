package api

import (
	"encoding/json/v2"
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ImportPresetsApi represents import preset api
type ImportPresetsApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	importPresets *services.ImportPresetsService
}

// Initialize a import preset api singleton instance
var (
	ImportPresets = &ImportPresetsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		importPresets: services.ImportPresets,
	}
)

// ImportPresetListHandler returns preset list of current user
func (a *ImportPresetsApi) ImportPresetListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	presets, err := a.importPresets.GetAllPresetNamesByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetListHandler] failed to get import preset for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	presetResps := make(models.ImportPresetInfoResponseSlice, len(presets))

	for i := 0; i < len(presets); i++ {
		presetResps[i], err = presets[i].ToImportPresetInfoResponse()

		if err != nil {
			log.Errorf(c, "[import_presets.ImportPresetListHandler] failed to get import preset response for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrImportPresetDataInvalid
		}
	}

	sort.Sort(presetResps)

	return presetResps, nil
}

// ImportPresetGetHandler returns the full preset configuration of current user
func (a *ImportPresetsApi) ImportPresetGetHandler(c *core.WebContext) (any, *errs.Error) {
	var presetGetReq models.ImportPresetGetRequest
	err := c.ShouldBindQuery(&presetGetReq)

	if err != nil {
		log.Warnf(c, "[import_presets.ImportPresetGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	preset, err := a.importPresets.GetPresetByPresetId(c, uid, presetGetReq.Id)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetGetHandler] failed to get import preset \"id:%d\" for user \"uid:%d\", because %s", presetGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	presetResp, err := preset.ToImportPresetInfoResponse()

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetGetHandler] failed to get import preset response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrImportPresetDataInvalid
	}

	return presetResp, nil
}

// ImportPresetCreateHandler creates a new preset by request parameters for current user
func (a *ImportPresetsApi) ImportPresetCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var presetCreateReq models.ImportPresetCreateRequest
	err := c.ShouldBindJSON(&presetCreateReq)

	if err != nil {
		log.Warnf(c, "[import_presets.ImportPresetCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	maxOrderId, err := a.importPresets.GetMaxDisplayOrder(c, uid)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	preset, err := a.createNewImportPresetModel(uid, &presetCreateReq, maxOrderId+1)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetCreateHandler] failed to parse import preset data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrImportPresetDataInvalid)
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && presetCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_IMPORT_PRESET, uid, presetCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[import_presets.ImportPresetCreateHandler] another import preset \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			presetId, err := utils.StringToInt64(remark)

			if err == nil {
				preset, err = a.importPresets.GetPresetByPresetId(c, uid, presetId)

				if err != nil {
					log.Errorf(c, "[import_presets.ImportPresetCreateHandler] failed to get existed import preset \"id:%d\" for user \"uid:%d\", because %s", presetId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				presetResp, err := preset.ToImportPresetInfoResponse()

				if err != nil {
					log.Errorf(c, "[import_presets.ImportPresetCreateHandler] failed to get import preset response for user \"uid:%d\", because %s", uid, err.Error())
					return nil, errs.ErrImportPresetDataInvalid
				}

				return presetResp, nil
			}
		}
	}

	err = a.importPresets.CreatePreset(c, preset)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetCreateHandler] failed to create import preset \"id:%d\" for user \"uid:%d\", because %s", preset.PresetId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[import_presets.ImportPresetCreateHandler] user \"uid:%d\" has created a new import preset \"id:%d\" successfully", uid, preset.PresetId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_IMPORT_PRESET, uid, presetCreateReq.ClientSessionId, utils.Int64ToString(preset.PresetId))
	presetResp, err := preset.ToImportPresetInfoResponse()

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetCreateHandler] failed to get import preset response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrImportPresetDataInvalid
	}

	return presetResp, nil
}

// ImportPresetModifyHandler saves an existed import preset by request parameters for current user
func (a *ImportPresetsApi) ImportPresetModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var presetModifyReq models.ImportPresetModifyRequest
	err := c.ShouldBindJSON(&presetModifyReq)

	if err != nil {
		log.Warnf(c, "[import_presets.ImportPresetModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	preset, err := a.importPresets.GetPresetByPresetId(c, uid, presetModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetModifyHandler] failed to get import preset \"id:%d\" for user \"uid:%d\", because %s", presetModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	newData, err := a.serializePresetData(&presetModifyReq.ImportPresetData)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetModifyHandler] failed to parse import preset data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrImportPresetDataInvalid)
	}

	newPreset := &models.ImportPreset{
		PresetId: preset.PresetId,
		Uid:      uid,
		Name:     presetModifyReq.Name,
		Data:     newData,
	}

	if newPreset.Name == preset.Name && newPreset.Data == preset.Data {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.importPresets.ModifyPreset(c, newPreset)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetModifyHandler] failed to update import preset \"id:%d\" for user \"uid:%d\", because %s", preset.PresetId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[import_presets.ImportPresetModifyHandler] user \"uid:%d\" has updated import preset \"id:%d\" successfully", uid, presetModifyReq.Id)

	preset.Name = newPreset.Name
	preset.Data = newPreset.Data
	presetResp, err := preset.ToImportPresetInfoResponse()

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetModifyHandler] failed to get import preset response for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrImportPresetDataInvalid
	}

	return presetResp, nil
}

// ImportPresetMoveHandler moves display order of existed import preset by request parameters for current user
func (a *ImportPresetsApi) ImportPresetMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var presetMoveReq models.ImportPresetMoveRequest
	err := c.ShouldBindJSON(&presetMoveReq)

	if err != nil {
		log.Warnf(c, "[import_presets.ImportPresetMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	presets := make([]*models.ImportPreset, len(presetMoveReq.NewDisplayOrders))

	for i := 0; i < len(presetMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := presetMoveReq.NewDisplayOrders[i]
		preset := &models.ImportPreset{
			Uid:          uid,
			PresetId:     newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		presets[i] = preset
	}

	err = a.importPresets.ModifyPresetDisplayOrders(c, uid, presets)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetMoveHandler] failed to move import presets for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[import_presets.ImportPresetMoveHandler] user \"uid:%d\" has moved import presets", uid)
	return true, nil
}

// ImportPresetDeleteHandler deletes an existed import preset by request parameters for current user
func (a *ImportPresetsApi) ImportPresetDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var presetDeleteReq models.ImportPresetDeleteRequest
	err := c.ShouldBindJSON(&presetDeleteReq)

	if err != nil {
		log.Warnf(c, "[import_presets.ImportPresetDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.importPresets.DeletePreset(c, uid, presetDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[import_presets.ImportPresetDeleteHandler] failed to delete import preset \"id:%d\" for user \"uid:%d\", because %s", presetDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[import_presets.ImportPresetDeleteHandler] user \"uid:%d\" has deleted import preset \"id:%d\"", uid, presetDeleteReq.Id)
	return true, nil
}

func (a *ImportPresetsApi) createNewImportPresetModel(uid int64, presetCreateReq *models.ImportPresetCreateRequest, order int32) (*models.ImportPreset, error) {
	data, err := a.serializePresetData(&presetCreateReq.ImportPresetData)

	if err != nil {
		return nil, err
	}

	return &models.ImportPreset{
		Uid:          uid,
		Name:         presetCreateReq.Name,
		Data:         data,
		DisplayOrder: order,
	}, nil
}

func (a *ImportPresetsApi) serializePresetData(preset *models.ImportPresetData) (string, error) {
	data, err := json.Marshal(preset)

	if err != nil {
		return "", errs.ErrImportPresetDataInvalid
	}

	if len(data) > models.MaxImportPresetDataSize {
		return "", errs.ErrImportPresetDataTooLarge
	}

	return string(data), nil
}
