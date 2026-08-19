package api

import (
	"image"
	_ "image/png"
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

const (
	maximumUserCustomIconPixels      = 256
	allowedUserCustomIconContentType = "image/png"
)

// UserCustomIconsApi represents user custom icons api
type UserCustomIconsApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	icons *services.UserCustomIconService
	users *services.UserService
}

// Initialize a user custom icons api singleton instance
var (
	UserCustomIcons = &UserCustomIconsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		icons: services.UserCustomIcons,
		users: services.Users,
	}
)

// CustomIconListHandler returns all custom icon infos of current user
func (a *UserCustomIconsApi) CustomIconListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	customIcons, err := a.icons.GetAllCustomIconInfosByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconListHandler] failed to get custom icons for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	customIconResps := make(models.UserCustomIconInfoResponseSlice, len(customIcons))

	for i := 0; i < len(customIcons); i++ {
		customIconResps[i] = customIcons[i].ToUserCustomIconInfoResponse()
	}

	sort.Sort(customIconResps)

	return customIconResps, nil
}

// CustomIconUploadHandler saves a new custom icon for current user
func (a *UserCustomIconsApi) CustomIconUploadHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPLOAD_CUSTOM_ICON) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	form, err := c.MultipartForm()

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to get multi-part form data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	customIconFiles := form.File["icon"]

	if len(customIconFiles) < 1 {
		log.Warnf(c, "[user_custom_icons.CustomIconUploadHandler] there is no custom icon in request for user \"uid:%d\"", uid)
		return nil, errs.ErrNoUserCustomIcon
	}

	if customIconFiles[0].Size < 1 {
		log.Warnf(c, "[user_custom_icons.CustomIconUploadHandler] the size of custom icon in request is zero for user \"uid:%d\"", uid)
		return nil, errs.ErrUserCustomIconIsEmpty
	}

	if customIconFiles[0].Size > int64(a.CurrentConfig().MaxUserCustomIconFileSize) {
		log.Warnf(c, "[user_custom_icons.CustomIconUploadHandler] the upload file size \"%d\" exceeds the maximum size \"%d\" of custom icon for user \"uid:%d\"", customIconFiles[0].Size, a.CurrentConfig().MaxUserCustomIconFileSize, uid)
		return nil, errs.ErrExceedMaxUserCustomIconFileSize
	}

	fileExtension := utils.GetFileNameExtension(customIconFiles[0].Filename)

	if utils.GetImageContentType(fileExtension) != allowedUserCustomIconContentType {
		log.Warnf(c, "[user_custom_icons.CustomIconUploadHandler] the file extension \"%s\" of custom icon in request is not supported for user \"uid:%d\"", fileExtension, uid)
		return nil, errs.ErrUserCustomIconExtensionInvalid
	}

	customIconFile, err := customIconFiles[0].Open()

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to get custom icon file from request for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	imageConfig, imageFormat, err := image.DecodeConfig(customIconFile)

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to decode custom icon image config for user \"uid:%d\", because %s", uid, err.Error())
		_ = customIconFile.Close()
		return nil, errs.ErrImageTypeNotSupported
	}

	if imageFormat != models.UserCustomIconFileExtension {
		log.Warnf(c, "[user_custom_icons.CustomIconUploadHandler] unsupported custom icon image format \"%s\" for user \"uid:%d\"", imageFormat, uid)
		_ = customIconFile.Close()
		return nil, errs.ErrUserCustomIconExtensionInvalid
	}

	if imageConfig.Width < 1 || imageConfig.Height < 1 || imageConfig.Width > maximumUserCustomIconPixels || imageConfig.Height > maximumUserCustomIconPixels {
		log.Warnf(c, "[user_custom_icons.CustomIconUploadHandler] invalid custom icon image dimensions \"%dx%d\" for user \"uid:%d\"", imageConfig.Width, imageConfig.Height, uid)
		_ = customIconFile.Close()
		return nil, errs.ErrUserCustomIconDimensionsInvalid
	}

	if _, err = customIconFile.Seek(0, 0); err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to seek to the beginning of custom icon file for user \"uid:%d\", because %s", uid, err.Error())
		_ = customIconFile.Close()
		return nil, errs.ErrOperationFailed
	}

	maxOrderId, err := a.icons.GetMaxDisplayOrder(c, uid)

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	customIconInfo := &models.UserCustomIcon{
		Uid:          uid,
		DisplayOrder: maxOrderId + 1,
	}

	clientSessionIds := form.Value["clientSessionId"]
	clientSessionId := ""

	if len(clientSessionIds) > 0 {
		clientSessionId = clientSessionIds[0]
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && clientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CUSTOM_ICON, uid, clientSessionId)

		if found {
			log.Infof(c, "[user_custom_icons.CustomIconUploadHandler] another custom icon \"id:%s\" has been uploaded for user \"uid:%d\"", remark, uid)
			iconId, err := utils.StringToInt64(remark)

			if err == nil {
				customIconInfo, err = a.icons.GetCustomIconInfoByIconId(c, uid, iconId)

				if err != nil {
					log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to get existed custom icon \"id:%d\" for user \"uid:%d\", because %s", iconId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				customIconResp := customIconInfo.ToUserCustomIconInfoResponse()

				return customIconResp, nil
			}
		}
	}

	err = a.icons.UploadCustomIcon(c, customIconInfo, customIconFile)

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconUploadHandler] failed to upload custom icon for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CUSTOM_ICON, uid, clientSessionId, utils.Int64ToString(customIconInfo.IconId))
	customIconResp := customIconInfo.ToUserCustomIconInfoResponse()

	return customIconResp, nil
}

// CustomIconGetHandler returns custom icon data for current user
func (a *UserCustomIconsApi) CustomIconGetHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	fileName := c.Param("fileName")
	fileExtension := utils.GetFileNameExtension(fileName)
	contentType := utils.GetImageContentType(fileExtension)

	if contentType != allowedUserCustomIconContentType || fileExtension != models.UserCustomIconFileExtension {
		return nil, "", errs.ErrUserCustomIconExtensionInvalid
	}

	fileBaseName := utils.GetFileNameWithoutExtension(fileName)
	iconId, err := utils.StringToInt64(fileBaseName)

	if err != nil {
		return nil, "", errs.ErrUserCustomIconIdInvalid
	}

	uid := c.GetCurrentUid()
	customIconData, err := a.icons.GetCustomIconByIconId(c, uid, iconId)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[user_custom_icons.CustomIconGetHandler] failed to get custom icon \"id:%d\", because %s", iconId, err.Error())
		}

		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	return customIconData, allowedUserCustomIconContentType, nil
}

// CustomIconMoveHandler modifies display order of custom icons for current user
func (a *UserCustomIconsApi) CustomIconMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var iconMoveReq models.UserCustomIconMoveRequest
	err := c.ShouldBindJSON(&iconMoveReq)

	if err != nil {
		log.Warnf(c, "[user_custom_icons.CustomIconMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	customIcons := make([]*models.UserCustomIcon, len(iconMoveReq.NewDisplayOrders))

	for i := 0; i < len(iconMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := iconMoveReq.NewDisplayOrders[i]
		customIcons[i] = &models.UserCustomIcon{
			IconId:       newDisplayOrder.Id,
			Uid:          uid,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}
	}

	err = a.icons.ModifyCustomIconDisplayOrders(c, uid, customIcons)

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconMoveHandler] failed to modify custom icon display orders for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[user_custom_icons.CustomIconMoveHandler] user \"uid:%d\" has moved custom icons", uid)
	return true, nil
}

// CustomIconDeleteHandler deletes specified custom icon for current user
func (a *UserCustomIconsApi) CustomIconDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var iconDeleteReq models.UserCustomIconDeleteRequest
	err := c.ShouldBindJSON(&iconDeleteReq)

	if err != nil {
		log.Warnf(c, "[user_custom_icons.CustomIconDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[user_custom_icons.CustomIconDeleteHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_UPLOAD_CUSTOM_ICON) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	err = a.icons.DeleteCustomIcon(c, uid, iconDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[user_custom_icons.CustomIconDeleteHandler] failed to delete custom icon \"id:%d\" for user \"uid:%d\", because %s", iconDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[user_custom_icons.CustomIconDeleteHandler] user \"uid:%d\" has deleted custom icon \"id:%d\"", uid, iconDeleteReq.Id)
	return true, nil
}
