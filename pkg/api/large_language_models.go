package api

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/llm"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/templates"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// LargeLanguageModelsApi represents large language models api
type LargeLanguageModelsApi struct {
	ApiUsingConfig
	transactionCategories *services.TransactionCategoryService
	transactionTags       *services.TransactionTagService
	accounts              *services.AccountService
	users                 *services.UserService
}

// Initialize a large language models api singleton instance
var (
	LargeLanguageModels = &LargeLanguageModelsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		transactionCategories: services.TransactionCategories,
		transactionTags:       services.TransactionTags,
		accounts:              services.Accounts,
		users:                 services.Users,
	}
)

// RecognizeReceiptImageHandler returns the recognized receipt image result
func (a *LargeLanguageModelsApi) RecognizeReceiptImageHandler(c *core.WebContext) (any, *errs.Error) {
	if a.CurrentConfig().ReceiptImageRecognitionLLMConfig == nil || a.CurrentConfig().ReceiptImageRecognitionLLMConfig.LLMProvider == "" || !a.CurrentConfig().TransactionFromAIImageRecognition {
		return nil, errs.ErrLargeLanguageModelProviderNotEnabled
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[large_language_models.RecognizeReceiptImageHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	timezone := time.FixedZone("Client Timezone", int(utcOffset)*60)
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return false, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_CREATE_TRANSACTION_FROM_AI_IMAGE_RECOGNITION) {
		return false, errs.ErrNotPermittedToPerformThisAction
	}

	form, err := c.MultipartForm()

	if err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to get multi-part form data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	imageFiles := form.File["image"]

	if len(imageFiles) < 1 {
		log.Warnf(c, "[large_language_models.RecognizeReceiptImageHandler] there is no image in request for user \"uid:%d\"", uid)
		return nil, errs.ErrNoAIRecognitionImage
	}

	if imageFiles[0].Size < 1 {
		log.Warnf(c, "[large_language_models.RecognizeReceiptImageHandler] the size of image in request is zero for user \"uid:%d\"", uid)
		return nil, errs.ErrAIRecognitionImageIsEmpty
	}

	if imageFiles[0].Size > int64(a.CurrentConfig().MaxAIRecognitionPictureFileSize) {
		log.Warnf(c, "[large_language_models.RecognizeReceiptImageHandler] the upload file size \"%d\" exceeds the maximum size \"%d\" of image for user \"uid:%d\"", imageFiles[0].Size, a.CurrentConfig().MaxAIRecognitionPictureFileSize, uid)
		return nil, errs.ErrExceedMaxAIRecognitionImageFileSize
	}

	fileExtension := utils.GetFileNameExtension(imageFiles[0].Filename)

	if utils.GetImageContentType(fileExtension) == "" {
		log.Warnf(c, "[large_language_models.RecognizeReceiptImageHandler] the file extension \"%s\" of image in request is not supported for user \"uid:%d\"", fileExtension, uid)
		return nil, errs.ErrImageTypeNotSupported
	}

	imageFile, err := imageFiles[0].Open()

	if err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to get image file from request for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	defer imageFile.Close()

	imageData, err := io.ReadAll(imageFile)

	if err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to read image file from request for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	accounts, err := a.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	accountMap := a.accounts.GetVisibleAccountNameMapByList(accounts)
	accountNames := make([]string, 0, len(accounts))

	for i := 0; i < len(accounts); i++ {
		if accounts[i].Hidden || accounts[i].Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
			continue
		}

		accountNames = append(accountNames, accounts[i].Name)
	}

	categories, err := a.transactionCategories.GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	incomeCategoryMap := make(map[string]*models.TransactionCategory)
	incomeCategoryNames := make([]string, 0)

	expenseCategoryMap := make(map[string]*models.TransactionCategory)
	expenseCategoryNames := make([]string, 0)

	transferCategoryMap := make(map[string]*models.TransactionCategory)
	transferCategoryNames := make([]string, 0)

	for i := 0; i < len(categories); i++ {
		category := categories[i]

		if category.Hidden || category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
			continue
		}

		if category.Type == models.CATEGORY_TYPE_INCOME {
			incomeCategoryMap[category.Name] = category
			incomeCategoryNames = append(incomeCategoryNames, category.Name)
		} else if category.Type == models.CATEGORY_TYPE_EXPENSE {
			expenseCategoryMap[category.Name] = category
			expenseCategoryNames = append(expenseCategoryNames, category.Name)
		} else if category.Type == models.CATEGORY_TYPE_TRANSFER {
			transferCategoryMap[category.Name] = category
			transferCategoryNames = append(transferCategoryNames, category.Name)
		}
	}

	tags, err := a.transactionTags.GetAllTagsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tagMap := a.transactionTags.GetVisibleTagNameMapByList(tags)
	tagNames := make([]string, 0, len(tags))

	for i := 0; i < len(tags); i++ {
		if tags[i].Hidden {
			continue
		}

		tagNames = append(tagNames, tags[i].Name)
	}

	systemPrompt, err := templates.GetTemplate(templates.SYSTEM_PROMPT_RECEIPT_IMAGE_RECOGNITION)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	systemPromptParams := map[string]any{
		"CurrentDateTime":          utils.FormatUnixTimeToLongDateTime(time.Now().Unix(), timezone),
		"AllExpenseCategoryNames":  strings.Join(expenseCategoryNames, "\n"),
		"AllIncomeCategoryNames":   strings.Join(incomeCategoryNames, "\n"),
		"AllTransferCategoryNames": strings.Join(transferCategoryNames, "\n"),
		"AllAccountNames":          strings.Join(accountNames, "\n"),
		"AllTagNames":              strings.Join(tagNames, "\n"),
	}

	var bodyBuffer bytes.Buffer
	err = systemPrompt.Execute(&bodyBuffer, systemPromptParams)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	llmRequest := &llm.LargeLanguageModelRequest{
		Stream:         false,
		SystemPrompt:   strings.ReplaceAll(bodyBuffer.String(), "\r\n", "\n"),
		UserPrompt:     imageData,
		UserPromptType: llm.LARGE_LANGUAGE_MODEL_REQUEST_PROMPT_TYPE_IMAGE_URL,
	}

	llmResponse, err := llm.Container.GetJsonResponseByReceiptImageRecognitionModel(c, c.GetCurrentUid(), a.CurrentConfig(), llmRequest)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if llmResponse == nil || len(llmResponse.Content) == 0 || strings.HasPrefix(llmResponse.Content, "{}") {
		return nil, errs.ErrNoTransactionInformationInImage
	}

	var result *models.RecognizedReceiptImageResult

	if err := json.Unmarshal([]byte(llmResponse.Content), &result); err != nil {
		log.Errorf(c, "[large_language_models.RecognizeReceiptImageHandler] failed to unmarshal recognized receipt image result from llm response \"%s\" for user \"uid:%d\", because %s", llmResponse.Content, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return a.parseRecognizedReceiptImageResponse(c, uid, utcOffset, result, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)
}

func (a *LargeLanguageModelsApi) parseRecognizedReceiptImageResponse(c *core.WebContext, uid int64, utcOffset int16, recognizedResult *models.RecognizedReceiptImageResult, accountMap map[string]*models.Account, expenseCategoryMap map[string]*models.TransactionCategory, incomeCategoryMap map[string]*models.TransactionCategory, transferCategoryMap map[string]*models.TransactionCategory, tagMap map[string]*models.TransactionTag) (*models.RecognizedReceiptImageResponse, *errs.Error) {
	recognizedReceiptImageResponse := &models.RecognizedReceiptImageResponse{
		Type: models.TRANSACTION_TYPE_EXPENSE,
	}

	if recognizedResult == nil {
		log.Errorf(c, "[large_language_models.parseRecognizedReceiptImageResponse] recoginzed result is null")
		return nil, errs.ErrNoTransactionInformationInImage
	}

	if recognizedResult.Type == "income" {
		recognizedReceiptImageResponse.Type = models.TRANSACTION_TYPE_INCOME

		if len(recognizedResult.CategoryName) > 0 {
			category, exists := incomeCategoryMap[recognizedResult.CategoryName]

			if exists {
				recognizedReceiptImageResponse.CategoryId = category.CategoryId
			}
		}
	} else if recognizedResult.Type == "expense" {
		recognizedReceiptImageResponse.Type = models.TRANSACTION_TYPE_EXPENSE

		if len(recognizedResult.CategoryName) > 0 {
			category, exists := expenseCategoryMap[recognizedResult.CategoryName]

			if exists {
				recognizedReceiptImageResponse.CategoryId = category.CategoryId
			}
		}
	} else if recognizedResult.Type == "transfer" {
		recognizedReceiptImageResponse.Type = models.TRANSACTION_TYPE_TRANSFER

		if len(recognizedResult.CategoryName) > 0 {
			category, exists := transferCategoryMap[recognizedResult.CategoryName]

			if exists {
				recognizedReceiptImageResponse.CategoryId = category.CategoryId
			}
		}
	} else if len(recognizedResult.Type) == 0 {
		return nil, errs.ErrNoTransactionInformationInImage
	} else {
		log.Errorf(c, "[large_language_models.parseRecognizedReceiptImageResponse] recoginzed transaction type \"%s\" is invalid", recognizedResult.Type)
		return nil, errs.ErrOperationFailed
	}

	if len(recognizedResult.Time) > 0 {
		longDateTime := a.getLongDateTime(recognizedResult.Time)
		timestamp, err := utils.ParseFromLongDateTime(longDateTime, utcOffset)

		if err != nil {
			log.Warnf(c, "[large_language_models.parseRecognizedReceiptImageResponse] recoginzed time \"%s\" is invalid", recognizedResult.Time)
		} else {
			recognizedReceiptImageResponse.Time = timestamp.Unix()
		}
	}

	if len(recognizedResult.Amount) > 0 {
		amount, err := utils.ParseAmount(recognizedResult.Amount)

		if err != nil {
			log.Errorf(c, "[large_language_models.parseRecognizedReceiptImageResponse] recoginzed amount \"%s\" is invalid", recognizedResult.Amount)
			return nil, errs.ErrOperationFailed
		}

		recognizedReceiptImageResponse.SourceAmount = amount

		if recognizedReceiptImageResponse.Type == models.TRANSACTION_TYPE_TRANSFER && len(recognizedResult.DestinationAmount) > 0 {
			destinationAmount, err := utils.ParseAmount(recognizedResult.DestinationAmount)

			if err != nil {
				log.Errorf(c, "[large_language_models.parseRecognizedReceiptImageResponse] recoginzed destination amount \"%s\" is invalid", recognizedResult.DestinationAmount)
				return nil, errs.ErrOperationFailed
			}

			recognizedReceiptImageResponse.DestinationAmount = destinationAmount
		}
	}

	if len(recognizedResult.AccountName) > 0 {
		account, exists := accountMap[recognizedResult.AccountName]

		if exists {
			recognizedReceiptImageResponse.SourceAccountId = account.AccountId
		}
	}

	if len(recognizedResult.DestinationAccountName) > 0 {
		account, exists := accountMap[recognizedResult.DestinationAccountName]

		if exists {
			recognizedReceiptImageResponse.DestinationAccountId = account.AccountId
		}
	}

	if len(recognizedResult.TagNames) > 0 {
		tagIds := make([]string, 0, len(recognizedResult.TagNames))

		for i := 0; i < len(recognizedResult.TagNames); i++ {
			tagName := recognizedResult.TagNames[i]
			tag, exists := tagMap[tagName]

			if exists {
				tagIds = append(tagIds, utils.Int64ToString(tag.TagId))
			}
		}

		recognizedReceiptImageResponse.TagIds = tagIds
	}

	if len(recognizedResult.Description) > 0 {
		recognizedReceiptImageResponse.Comment = recognizedResult.Description
	}

	return recognizedReceiptImageResponse, nil
}

func (a *LargeLanguageModelsApi) getLongDateTime(dateTime string) string {
	if utils.IsValidLongDateTimeFormat(dateTime) {
		return dateTime
	}

	if utils.IsValidLongDateTimeWithoutSecondFormat(dateTime) {
		return dateTime + ":00"
	}

	if utils.IsValidLongDateFormat(dateTime) {
		return dateTime + " 00:00:00"
	}

	return dateTime
}
