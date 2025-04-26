package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TransactionPicturesApi represents transaction pictures api
type TransactionPicturesApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	users    *services.UserService
	pictures *services.TransactionPictureService
}

// Initialize a transaction api singleton instance
var (
	TransactionPictures = &TransactionPicturesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		users:    services.Users,
		pictures: services.TransactionPictures,
	}
)

// TransactionPictureUploadHandler saves transaction picture by request parameters for current user
func (a *TransactionPicturesApi) TransactionPictureUploadHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	form, err := c.MultipartForm()

	if err != nil {
		log.Errorf(c, "[transaction_pictures.TransactionPictureUploadHandler] failed to get multi-part form data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	pictureFiles := form.File["picture"]

	if len(pictureFiles) < 1 {
		log.Warnf(c, "[transaction_pictures.TransactionPictureUploadHandler] there is no transaction picture in request for user \"uid:%d\"", uid)
		return nil, errs.ErrNoTransactionPicture
	}

	if pictureFiles[0].Size < 1 {
		log.Warnf(c, "[transaction_pictures.TransactionPictureUploadHandler] the size of transaction picture in request is zero for user \"uid:%d\"", uid)
		return nil, errs.ErrTransactionPictureIsEmpty
	}

	if pictureFiles[0].Size > int64(a.CurrentConfig().MaxTransactionPictureFileSize) {
		log.Warnf(c, "[transaction_pictures.TransactionPictureUploadHandler] the upload file size \"%d\" exceeds the maximum size \"%d\" of transaction picture for user \"uid:%d\"", pictureFiles[0].Size, a.CurrentConfig().MaxTransactionPictureFileSize, uid)
		return nil, errs.ErrExceedMaxTransactionPictureFileSize
	}

	fileExtension := utils.GetFileNameExtension(pictureFiles[0].Filename)

	if utils.GetImageContentType(fileExtension) == "" {
		log.Warnf(c, "[transaction_pictures.TransactionPictureUploadHandler] the file extension \"%s\" of transaction picture in request is not supported for user \"uid:%d\"", fileExtension, uid)
		return nil, errs.ErrImageTypeNotSupported
	}

	pictureFile, err := pictureFiles[0].Open()

	if err != nil {
		log.Errorf(c, "[transaction_pictures.TransactionPictureUploadHandler] failed to get transaction picture file from request for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	pictureInfo := a.createNewPictureInfoModel(uid, fileExtension, c.ClientIP())

	clientSessionIds := form.Value["clientSessionId"]
	clientSessionId := ""

	if len(clientSessionIds) > 0 {
		clientSessionId = clientSessionIds[0]
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && clientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_PICTURE, uid, clientSessionId)

		if found {
			log.Infof(c, "[transaction_pictures.TransactionPictureUploadHandler] another transaction picture \"id:%s\" has been uploaded for user \"uid:%d\"", remark, uid)
			pictureId, err := utils.StringToInt64(remark)

			if err == nil {
				pictureInfo, err = a.pictures.GetPictureInfoByPictureId(c, uid, pictureId)

				if err != nil {
					log.Errorf(c, "[transaction_pictures.TransactionPictureUploadHandler] failed to get existed transaction picture \"id:%d\" for user \"uid:%d\", because %s", pictureId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				pictureInfoResp := a.GetTransactionPictureInfoResponse(pictureInfo)

				return pictureInfoResp, nil
			}
		}
	}

	err = a.pictures.UploadPicture(c, pictureInfo, pictureFile)

	if err != nil {
		log.Errorf(c, "[transaction_pictures.TransactionPictureUploadHandler] failed to update transaction picture for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_PICTURE, uid, clientSessionId, utils.Int64ToString(pictureInfo.PictureId))
	pictureInfoResp := a.GetTransactionPictureInfoResponse(pictureInfo)

	return pictureInfoResp, nil
}

// TransactionPictureGetHandler returns transaction picture data for current user
func (a *TransactionPicturesApi) TransactionPictureGetHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	fileName := c.Param("fileName")
	fileExtension := utils.GetFileNameExtension(fileName)
	contentType := utils.GetImageContentType(fileExtension)

	if contentType == "" {
		return nil, "", errs.ErrImageTypeNotSupported
	}

	fileBaseName := utils.GetFileNameWithoutExtension(fileName)
	pictureId, err := utils.StringToInt64(fileBaseName)

	if err != nil {
		return nil, "", errs.ErrTransactionPictureIdInvalid
	}

	uid := c.GetCurrentUid()
	pictureData, err := a.pictures.GetPictureByPictureId(c, uid, pictureId, fileExtension)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transaction_pictures.TransactionPictureUploadHandler] failed to get transaction picture, because %s", err.Error())
		}

		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	return pictureData, contentType, nil
}

// TransactionPictureRemoveUnusedHandler removes unused transaction picture by request parameters for current user
func (a *TransactionPicturesApi) TransactionPictureRemoveUnusedHandler(c *core.WebContext) (any, *errs.Error) {
	var pictureDeleteReq models.TransactionPictureUnusedDeleteRequest
	err := c.ShouldBindJSON(&pictureDeleteReq)

	if err != nil {
		log.Warnf(c, "[transaction_pictures.TransactionPictureRemoveUnusedHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.pictures.RemoveUnusedTransactionPicture(c, uid, pictureDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_pictures.TransactionPictureRemoveUnusedHandler] failed to remove unused transaction picture for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

func (a *TransactionPicturesApi) createNewPictureInfoModel(uid int64, fileExtension string, clientIp string) *models.TransactionPictureInfo {
	return &models.TransactionPictureInfo{
		Uid:              uid,
		TransactionId:    models.TransactionPictureNewPictureTransactionId,
		PictureExtension: fileExtension,
		CreatedIp:        clientIp,
	}
}
