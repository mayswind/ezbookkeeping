package api

import (
	"sort"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const maximumTagsCountOfTemplate = 10

// TransactionTemplatesApi represents transaction template api
type TransactionTemplatesApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	templates *services.TransactionTemplateService
}

// Initialize a transaction template api singleton instance
var (
	TransactionTemplates = &TransactionTemplatesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		templates: services.TransactionTemplates,
	}
)

// TemplateListHandler returns transaction template list of current user
func (a *TransactionTemplatesApi) TemplateListHandler(c *core.WebContext) (any, *errs.Error) {
	var templateListReq models.TransactionTemplateListRequest
	err := c.ShouldBindQuery(&templateListReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if templateListReq.TemplateType < models.TRANSACTION_TEMPLATE_TYPE_NORMAL || templateListReq.TemplateType > models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		log.Warnf(c, "[transaction_templates.TemplateListHandler] template type invalid, type is %d", templateListReq.TemplateType)
		return nil, errs.ErrTransactionTemplateTypeInvalid
	}

	if templateListReq.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
		return nil, errs.ErrScheduledTransactionNotEnabled
	}

	uid := c.GetCurrentUid()
	templates, err := a.templates.GetAllTemplatesByUid(c, uid, templateListReq.TemplateType)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateListHandler] failed to get templates for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	templateResps := make(models.TransactionTemplateInfoResponseSlice, len(templates))
	serverUtcOffset := utils.GetServerTimezoneOffsetMinutes()

	for i := 0; i < len(templates); i++ {
		templateResps[i] = templates[i].ToTransactionTemplateInfoResponse(serverUtcOffset)
	}

	sort.Sort(templateResps)

	return templateResps, nil
}

// TemplateGetHandler returns one specific transaction template of current user
func (a *TransactionTemplatesApi) TemplateGetHandler(c *core.WebContext) (any, *errs.Error) {
	var templateGetReq models.TransactionTemplateGetRequest
	err := c.ShouldBindQuery(&templateGetReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	template, err := a.templates.GetTemplateByTemplateId(c, uid, templateGetReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateGetHandler] failed to get template \"id:%d\" for user \"uid:%d\", because %s", templateGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
		return nil, errs.ErrScheduledTransactionNotEnabled
	}

	serverUtcOffset := utils.GetServerTimezoneOffsetMinutes()
	templateResp := template.ToTransactionTemplateInfoResponse(serverUtcOffset)

	return templateResp, nil
}

// TemplateCreateHandler saves a new transaction template by request parameters for current user
func (a *TransactionTemplatesApi) TemplateCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var templateCreateReq models.TransactionTemplateCreateRequest
	err := c.ShouldBindJSON(&templateCreateReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if templateCreateReq.TemplateType < models.TRANSACTION_TEMPLATE_TYPE_NORMAL || templateCreateReq.TemplateType > models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		log.Warnf(c, "[transaction_templates.TemplateCreateHandler] template type invalid, type is %d", templateCreateReq.TemplateType)
		return nil, errs.ErrTransactionTemplateTypeInvalid
	}

	if templateCreateReq.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
		return nil, errs.ErrScheduledTransactionNotEnabled
	}

	if templateCreateReq.Type <= models.TRANSACTION_TYPE_MODIFY_BALANCE || templateCreateReq.Type > models.TRANSACTION_TYPE_TRANSFER {
		log.Warnf(c, "[transaction_templates.TemplateCreateHandler] transaction type invalid, type is %d", templateCreateReq.Type)
		return nil, errs.ErrTransactionTypeInvalid
	}

	if templateCreateReq.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		if templateCreateReq.ScheduledFrequencyType == nil ||
			templateCreateReq.ScheduledFrequency == nil ||
			templateCreateReq.ScheduledTimezoneUtcOffset == nil {
			return nil, errs.ErrScheduledTransactionFrequencyInvalid
		}

		if *templateCreateReq.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED && *templateCreateReq.ScheduledFrequency != "" {
			return nil, errs.ErrScheduledTransactionFrequencyInvalid
		} else if *templateCreateReq.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED && *templateCreateReq.ScheduledFrequency == "" {
			return nil, errs.ErrScheduledTransactionFrequencyInvalid
		}
	}

	if len(templateCreateReq.TagIds) > maximumTagsCountOfTemplate {
		return nil, errs.ErrTransactionTemplateHasTooManyTags
	}

	uid := c.GetCurrentUid()

	maxOrderId, err := a.templates.GetMaxDisplayOrder(c, uid, templateCreateReq.TemplateType)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	serverUtcOffset := utils.GetServerTimezoneOffsetMinutes()
	template, err := a.createNewTemplateModel(uid, &templateCreateReq, maxOrderId+1)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateCreateHandler] failed to create new template for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && templateCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_TEMPLATE, uid, templateCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[transaction_templates.TemplateCreateHandler] another template \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			templateId, err := utils.StringToInt64(remark)

			if err == nil {
				template, err = a.templates.GetTemplateByTemplateId(c, uid, templateId)

				if err != nil {
					log.Errorf(c, "[transaction_templates.TemplateCreateHandler] failed to get existed template \"id:%d\" for user \"uid:%d\", because %s", templateId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				templateResp := template.ToTransactionTemplateInfoResponse(serverUtcOffset)

				return templateResp, nil
			}
		}
	}

	err = a.templates.CreateTemplate(c, template)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateCreateHandler] failed to create template \"id:%d\" for user \"uid:%d\", because %s", template.TemplateId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_templates.TemplateCreateHandler] user \"uid:%d\" has created a new template \"id:%d\" successfully", uid, template.TemplateId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_TEMPLATE, uid, templateCreateReq.ClientSessionId, utils.Int64ToString(template.TemplateId))
	templateResp := template.ToTransactionTemplateInfoResponse(serverUtcOffset)

	return templateResp, nil
}

// TemplateModifyHandler saves an existed transaction template by request parameters for current user
func (a *TransactionTemplatesApi) TemplateModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var templateModifyReq models.TransactionTemplateModifyRequest
	err := c.ShouldBindJSON(&templateModifyReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if templateModifyReq.Type <= models.TRANSACTION_TYPE_MODIFY_BALANCE || templateModifyReq.Type > models.TRANSACTION_TYPE_TRANSFER {
		log.Warnf(c, "[transaction_templates.TemplateModifyHandler] transaction type invalid, type is %d", templateModifyReq.Type)
		return nil, errs.ErrTransactionTypeInvalid
	}

	uid := c.GetCurrentUid()
	template, err := a.templates.GetTemplateByTemplateId(c, uid, templateModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateModifyHandler] failed to get template \"id:%d\" for user \"uid:%d\", because %s", templateModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
		return nil, errs.ErrScheduledTransactionNotEnabled
	}

	if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		if templateModifyReq.ScheduledFrequencyType == nil ||
			templateModifyReq.ScheduledFrequency == nil ||
			templateModifyReq.ScheduledTimezoneUtcOffset == nil {
			return nil, errs.ErrScheduledTransactionFrequencyInvalid
		}

		if *templateModifyReq.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED && *templateModifyReq.ScheduledFrequency != "" {
			return nil, errs.ErrScheduledTransactionFrequencyInvalid
		} else if *templateModifyReq.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED && *templateModifyReq.ScheduledFrequency == "" {
			return nil, errs.ErrScheduledTransactionFrequencyInvalid
		}
	}

	if len(templateModifyReq.TagIds) > maximumTagsCountOfTemplate {
		return nil, errs.ErrTransactionTemplateHasTooManyTags
	}

	newTemplate := &models.TransactionTemplate{
		TemplateId:           template.TemplateId,
		Uid:                  uid,
		Name:                 templateModifyReq.Name,
		Type:                 templateModifyReq.Type,
		CategoryId:           templateModifyReq.CategoryId,
		AccountId:            templateModifyReq.SourceAccountId,
		TagIds:               strings.Join(templateModifyReq.TagIds, ","),
		Amount:               templateModifyReq.SourceAmount,
		RelatedAccountId:     templateModifyReq.DestinationAccountId,
		RelatedAccountAmount: templateModifyReq.DestinationAmount,
		HideAmount:           templateModifyReq.HideAmount,
		Comment:              templateModifyReq.Comment,
	}

	if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		newTemplate.ScheduledFrequencyType = *templateModifyReq.ScheduledFrequencyType
		newTemplate.ScheduledFrequency = a.getOrderedFrequencyValues(*templateModifyReq.ScheduledFrequency)
		newTemplate.ScheduledAt = a.getUTCScheduledAt(*templateModifyReq.ScheduledTimezoneUtcOffset)
		newTemplate.ScheduledTimezoneUtcOffset = *templateModifyReq.ScheduledTimezoneUtcOffset

		if templateModifyReq.ScheduledStartDate != nil {
			startTime, err := utils.ParseFromLongDateFirstTime(*templateModifyReq.ScheduledStartDate, *templateModifyReq.ScheduledTimezoneUtcOffset)

			if err != nil {
				log.Errorf(c, "[transaction_templates.TemplateModifyHandler] failed to parse scheduled start date for user \"uid:%d\", because %s", uid, err.Error())
				return nil, errs.Or(err, errs.ErrOperationFailed)
			}

			startUnixTime := startTime.Unix()
			newTemplate.ScheduledStartTime = &startUnixTime
		}

		if templateModifyReq.ScheduledEndDate != nil {
			endTime, err := utils.ParseFromLongDateLastTime(*templateModifyReq.ScheduledEndDate, *templateModifyReq.ScheduledTimezoneUtcOffset)

			if err != nil {
				log.Errorf(c, "[transaction_templates.TemplateModifyHandler] failed to parse scheduled end date for user \"uid:%d\", because %s", uid, err.Error())
				return nil, errs.Or(err, errs.ErrOperationFailed)
			}

			endUnixTime := endTime.Unix()
			newTemplate.ScheduledEndTime = &endUnixTime
		}

		if newTemplate.ScheduledStartTime != nil && newTemplate.ScheduledEndTime != nil && *newTemplate.ScheduledStartTime > *newTemplate.ScheduledEndTime {
			return nil, errs.ErrScheduledTransactionTemplateStartDataLaterThanEndDate
		}
	}

	if newTemplate.Name == template.Name &&
		newTemplate.Type == template.Type &&
		newTemplate.CategoryId == template.CategoryId &&
		newTemplate.AccountId == template.AccountId &&
		newTemplate.TagIds == template.TagIds &&
		newTemplate.Amount == template.Amount &&
		newTemplate.RelatedAccountId == template.RelatedAccountId &&
		newTemplate.RelatedAccountAmount == template.RelatedAccountAmount &&
		newTemplate.HideAmount == template.HideAmount &&
		newTemplate.Comment == template.Comment {
		if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_NORMAL {
			return nil, errs.ErrNothingWillBeUpdated
		} else if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
			if newTemplate.ScheduledFrequencyType == template.ScheduledFrequencyType &&
				newTemplate.ScheduledFrequency == template.ScheduledFrequency &&
				newTemplate.ScheduledStartTime == template.ScheduledStartTime &&
				newTemplate.ScheduledEndTime == template.ScheduledEndTime &&
				newTemplate.ScheduledAt == template.ScheduledAt &&
				newTemplate.ScheduledTimezoneUtcOffset == template.ScheduledTimezoneUtcOffset {
				return nil, errs.ErrNothingWillBeUpdated
			}
		}
	}

	err = a.templates.ModifyTemplate(c, newTemplate)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateModifyHandler] failed to update template \"id:%d\" for user \"uid:%d\", because %s", templateModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_templates.TemplateModifyHandler] user \"uid:%d\" has updated template \"id:%d\" successfully", uid, templateModifyReq.Id)

	serverUtcOffset := utils.GetServerTimezoneOffsetMinutes()
	newTemplate.TemplateType = template.TemplateType
	newTemplate.DisplayOrder = template.DisplayOrder
	newTemplate.Hidden = template.Hidden
	templateResp := newTemplate.ToTransactionTemplateInfoResponse(serverUtcOffset)

	return templateResp, nil
}

// TemplateHideHandler hides a transaction template by request parameters for current user
func (a *TransactionTemplatesApi) TemplateHideHandler(c *core.WebContext) (any, *errs.Error) {
	var templateHideReq models.TransactionTemplateHideRequest
	err := c.ShouldBindJSON(&templateHideReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	template, err := a.templates.GetTemplateByTemplateId(c, uid, templateHideReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateHideHandler] failed to get template \"id:%d\" for user \"uid:%d\", because %s", templateHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
		return nil, errs.ErrScheduledTransactionNotEnabled
	}

	err = a.templates.HideTemplate(c, uid, []int64{templateHideReq.Id}, templateHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateHideHandler] failed to hide template \"id:%d\" for user \"uid:%d\", because %s", templateHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_templates.TemplateHideHandler] user \"uid:%d\" has hidden template \"id:%d\"", uid, templateHideReq.Id)
	return true, nil
}

// TemplateMoveHandler moves display order of existed transaction templates by request parameters for current user
func (a *TransactionTemplatesApi) TemplateMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var templateMoveReq models.TransactionTemplateMoveRequest
	err := c.ShouldBindJSON(&templateMoveReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if len(templateMoveReq.NewDisplayOrders) > 0 {
		template, err := a.templates.GetTemplateByTemplateId(c, uid, templateMoveReq.NewDisplayOrders[0].Id)

		if err != nil {
			log.Errorf(c, "[transaction_templates.TemplateMoveHandler] failed to get template \"id:%d\" for user \"uid:%d\", because %s", templateMoveReq.NewDisplayOrders[0].Id, uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
			return nil, errs.ErrScheduledTransactionNotEnabled
		}
	}

	templates := make([]*models.TransactionTemplate, len(templateMoveReq.NewDisplayOrders))

	for i := 0; i < len(templateMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := templateMoveReq.NewDisplayOrders[i]
		template := &models.TransactionTemplate{
			Uid:          uid,
			TemplateId:   newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		templates[i] = template
	}

	err = a.templates.ModifyTemplateDisplayOrders(c, uid, templates)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateMoveHandler] failed to move templates for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_templates.TemplateMoveHandler] user \"uid:%d\" has moved templates", uid)
	return true, nil
}

// TemplateDeleteHandler deletes an existed transaction template by request parameters for current user
func (a *TransactionTemplatesApi) TemplateDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var templateDeleteReq models.TransactionTemplateDeleteRequest
	err := c.ShouldBindJSON(&templateDeleteReq)

	if err != nil {
		log.Warnf(c, "[transaction_templates.TemplateDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	template, err := a.templates.GetTemplateByTemplateId(c, uid, templateDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateDeleteHandler] failed to get template \"id:%d\" for user \"uid:%d\", because %s", templateDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if template.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE && !a.CurrentConfig().EnableScheduledTransaction {
		return nil, errs.ErrScheduledTransactionNotEnabled
	}

	err = a.templates.DeleteTemplate(c, uid, templateDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_templates.TemplateDeleteHandler] failed to delete template \"id:%d\" for user \"uid:%d\", because %s", templateDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_templates.TemplateDeleteHandler] user \"uid:%d\" has deleted template \"id:%d\"", uid, templateDeleteReq.Id)
	return true, nil
}

func (a *TransactionTemplatesApi) createNewTemplateModel(uid int64, templateCreateReq *models.TransactionTemplateCreateRequest, order int32) (*models.TransactionTemplate, error) {
	template := &models.TransactionTemplate{
		Uid:                  uid,
		TemplateType:         templateCreateReq.TemplateType,
		Name:                 templateCreateReq.Name,
		Type:                 templateCreateReq.Type,
		CategoryId:           templateCreateReq.CategoryId,
		AccountId:            templateCreateReq.SourceAccountId,
		TagIds:               strings.Join(templateCreateReq.TagIds, ","),
		Amount:               templateCreateReq.SourceAmount,
		RelatedAccountId:     templateCreateReq.DestinationAccountId,
		RelatedAccountAmount: templateCreateReq.DestinationAmount,
		HideAmount:           templateCreateReq.HideAmount,
		Comment:              templateCreateReq.Comment,
		DisplayOrder:         order,
	}

	if templateCreateReq.TemplateType == models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		template.ScheduledFrequencyType = *templateCreateReq.ScheduledFrequencyType
		template.ScheduledFrequency = a.getOrderedFrequencyValues(*templateCreateReq.ScheduledFrequency)
		template.ScheduledAt = a.getUTCScheduledAt(*templateCreateReq.ScheduledTimezoneUtcOffset)
		template.ScheduledTimezoneUtcOffset = *templateCreateReq.ScheduledTimezoneUtcOffset

		if templateCreateReq.ScheduledStartDate != nil {
			startTime, err := utils.ParseFromLongDateFirstTime(*templateCreateReq.ScheduledStartDate, *templateCreateReq.ScheduledTimezoneUtcOffset)

			if err != nil {
				return nil, err
			}

			startUnixTime := startTime.Unix()
			template.ScheduledStartTime = &startUnixTime
		}

		if templateCreateReq.ScheduledEndDate != nil {
			endTime, err := utils.ParseFromLongDateLastTime(*templateCreateReq.ScheduledEndDate, *templateCreateReq.ScheduledTimezoneUtcOffset)

			if err != nil {
				return nil, err
			}

			endUnixTime := endTime.Unix()
			template.ScheduledEndTime = &endUnixTime
		}

		if template.ScheduledStartTime != nil && template.ScheduledEndTime != nil && *template.ScheduledStartTime > *template.ScheduledEndTime {
			return nil, errs.ErrScheduledTransactionTemplateStartDataLaterThanEndDate
		}
	}

	return template, nil
}

func (a *TransactionTemplatesApi) getUTCScheduledAt(scheduledTimezoneUtcOffset int16) int16 {
	templateTimeZone := time.FixedZone("Template Timezone", int(scheduledTimezoneUtcOffset)*60)
	transactionTime := time.Date(2020, 1, 1, 0, 0, 0, 0, templateTimeZone)
	transactionTimeInUTC := transactionTime.In(time.UTC)

	minutesElapsedOfDayInUtc := transactionTimeInUTC.Hour()*60 + transactionTimeInUTC.Minute()

	return int16(minutesElapsedOfDayInUtc)
}

func (a *TransactionTemplatesApi) getOrderedFrequencyValues(frequencyValue string) string {
	if frequencyValue == "" {
		return ""
	}

	items := strings.Split(frequencyValue, ",")
	values := make([]int, 0, len(items))
	valueExistMap := make(map[int]bool)

	for i := 0; i < len(items); i++ {
		value, err := utils.StringToInt(items[i])

		if err != nil {
			continue
		}

		if _, exists := valueExistMap[value]; !exists {
			values = append(values, value)
			valueExistMap[value] = true
		}
	}

	sort.Ints(values)

	var sortedFrequencyValueBuilder strings.Builder

	for i := 0; i < len(values); i++ {
		if sortedFrequencyValueBuilder.Len() > 0 {
			sortedFrequencyValueBuilder.WriteRune(',')
		}

		sortedFrequencyValueBuilder.WriteString(utils.IntToString(values[i]))
	}

	return sortedFrequencyValueBuilder.String()
}
