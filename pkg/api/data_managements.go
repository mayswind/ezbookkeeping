package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/converters"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const pageCountForClearTransactions = 1000
const pageCountForDataExport = 1000

// DataManagementsApi represents data management api
type DataManagementsApi struct {
	ApiUsingConfig
	tokens                  *services.TokenService
	users                   *services.UserService
	accounts                *services.AccountService
	transactions            *services.TransactionService
	categories              *services.TransactionCategoryService
	tags                    *services.TransactionTagService
	pictures                *services.TransactionPictureService
	templates               *services.TransactionTemplateService
	userCustomExchangeRates *services.UserCustomExchangeRatesService
}

// Initialize a data management api singleton instance
var (
	DataManagements = &DataManagementsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		tokens:                  services.Tokens,
		users:                   services.Users,
		accounts:                services.Accounts,
		transactions:            services.Transactions,
		categories:              services.TransactionCategories,
		tags:                    services.TransactionTags,
		pictures:                services.TransactionPictures,
		templates:               services.TransactionTemplates,
		userCustomExchangeRates: services.UserCustomExchangeRates,
	}
)

// ExportDataToEzbookkeepingCSVHandler returns exported data in csv format
func (a *DataManagementsApi) ExportDataToEzbookkeepingCSVHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	return a.getExportedFileContent(c, "csv")
}

// ExportDataToEzbookkeepingTSVHandler returns exported data in csv format
func (a *DataManagementsApi) ExportDataToEzbookkeepingTSVHandler(c *core.WebContext) ([]byte, string, *errs.Error) {
	return a.getExportedFileContent(c, "tsv")
}

// DataStatisticsHandler returns user data statistics
func (a *DataManagementsApi) DataStatisticsHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	totalAccountCount, err := a.accounts.GetTotalAccountCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total account count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionCategoryCount, err := a.categories.GetTotalCategoryCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total transaction category count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionTagCount, err := a.tags.GetTotalTagCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total transaction tag count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionCount, err := a.transactions.GetTotalTransactionCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total transaction count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionPictureCount, err := a.pictures.GetTotalTransactionPicturesCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total transaction picture count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionTemplateCount, err := a.templates.GetTotalNormalTemplateCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total transaction template count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalScheduledTransactionCount, err := a.templates.GetTotalScheduledTemplateCountByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.DataStatisticsHandler] failed to get total scheduled transaction count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	dataStatisticsResp := &models.DataStatisticsResponse{
		TotalAccountCount:              totalAccountCount,
		TotalTransactionCategoryCount:  totalTransactionCategoryCount,
		TotalTransactionTagCount:       totalTransactionTagCount,
		TotalTransactionCount:          totalTransactionCount,
		TotalTransactionPictureCount:   totalTransactionPictureCount,
		TotalTransactionTemplateCount:  totalTransactionTemplateCount,
		TotalScheduledTransactionCount: totalScheduledTransactionCount,
	}

	return dataStatisticsResp, nil
}

// ClearAllDataHandler deletes all user data
func (a *DataManagementsApi) ClearAllDataHandler(c *core.WebContext) (any, *errs.Error) {
	var clearDataReq models.ClearDataRequest
	err := c.ShouldBindJSON(&clearDataReq)

	if err != nil {
		log.Warnf(c, "[data_managements.ClearAllDataHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[data_managements.ClearAllDataHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(clearDataReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_CLEAR_ALL_DATA) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	err = a.templates.DeleteAllTemplates(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllDataHandler] failed to delete all transaction templates, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.transactions.DeleteAllTransactions(c, uid, true)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllDataHandler] failed to delete all transactions, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.categories.DeleteAllCategories(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllDataHandler] failed to delete all transaction categories, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.tags.DeleteAllTags(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllDataHandler] failed to delete all transaction tags, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.userCustomExchangeRates.DeleteAllCustomExchangeRates(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllDataHandler] failed to delete all user custom exchange rates, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[data_managements.ClearAllDataHandler] user \"uid:%d\" has cleared all data", uid)
	return true, nil
}

// ClearAllTransactionsHandler deletes all transactions
func (a *DataManagementsApi) ClearAllTransactionsHandler(c *core.WebContext) (any, *errs.Error) {
	var clearDataReq models.ClearDataRequest
	err := c.ShouldBindJSON(&clearDataReq)

	if err != nil {
		log.Warnf(c, "[data_managements.ClearAllTransactionsHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[data_managements.ClearAllTransactionsHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(clearDataReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_CLEAR_ALL_DATA) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	err = a.transactions.DeleteAllTransactions(c, uid, false)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllTransactionsHandler] failed to delete all transactions, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[data_managements.ClearAllTransactionsHandler] user \"uid:%d\" has cleared all transactions", uid)
	return true, nil
}

// ClearAllTransactionsByAccountHandler deletes all transactions of specified account
func (a *DataManagementsApi) ClearAllTransactionsByAccountHandler(c *core.WebContext) (any, *errs.Error) {
	var clearDataReq models.ClearAccountTransactionsRequest
	err := c.ShouldBindJSON(&clearDataReq)

	if err != nil {
		log.Warnf(c, "[data_managements.ClearAllTransactionsByAccountHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[data_managements.ClearAllTransactionsByAccountHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(clearDataReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_CLEAR_ALL_DATA) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	account, err := a.accounts.GetAccountByAccountId(c, uid, clearDataReq.AccountId)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllTransactionsByAccountHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", uid, clearDataReq.AccountId, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if account.Hidden {
		return nil, errs.ErrCannotDeleteTransactionInHiddenAccount
	}

	if account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		return nil, errs.ErrCannotDeleteTransactionInParentAccount
	}

	err = a.transactions.DeleteAllTransactionsOfAccount(c, uid, account.AccountId, pageCountForClearTransactions)

	if err != nil {
		log.Errorf(c, "[data_managements.ClearAllTransactionsByAccountHandler] failed to delete all transactions in account \"id:%d\", because %s", account.AccountId, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[data_managements.ClearAllTransactionsByAccountHandler] user \"uid:%d\" has cleared all transactions in account \"id:%d\"", uid, account.AccountId)
	return true, nil
}

func (a *DataManagementsApi) getExportedFileContent(c *core.WebContext, fileType string) ([]byte, string, *errs.Error) {
	if !a.CurrentConfig().EnableDataExport {
		return nil, "", errs.ErrDataExportNotAllowed
	}

	var exportTransactionDataReq models.ExportTransactionDataRequest
	err := c.ShouldBindQuery(&exportTransactionDataReq)

	if err != nil {
		log.Warnf(c, "[data_managements.getExportedFileContent] parse request failed, because %s", err.Error())
		return nil, "", errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	clientTimezone, _, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[data_managements.getExportedFileContent] cannot get client timezone, because %s", err.Error())
		clientTimezone = time.Local
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Warnf(c, "[data_managements.getExportedFileContent] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, "", errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_EXPORT_TRANSACTION) {
		return nil, "", errs.ErrNotPermittedToPerformThisAction
	}

	accounts, err := a.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.getExportedFileContent] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	categories, err := a.categories.GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.Errorf(c, "[data_managements.getExportedFileContent] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	tags, err := a.tags.GetAllTagsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.getExportedFileContent] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	tagIndexes, err := a.tags.GetAllTagIdsMapOfAllTransactions(c, uid)

	if err != nil {
		log.Errorf(c, "[data_managements.getExportedFileContent] failed to get tag index for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	accountMap := a.accounts.GetAccountMapByList(accounts)
	categoryMap := a.categories.GetCategoryMapByList(categories)
	tagMap := a.tags.GetTagMapByList(tags)

	allAccountIds, err := a.accounts.GetAccountOrSubAccountIds(c, exportTransactionDataReq.AccountIds, uid)

	if err != nil {
		log.Warnf(c, "[data_managements.getExportedFileContent] get account error, because %s", err.Error())
		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	allCategoryIds, err := a.categories.GetCategoryOrSubCategoryIds(c, exportTransactionDataReq.CategoryIds, uid)

	if err != nil {
		log.Warnf(c, "[data_managements.getExportedFileContent] get transaction category error, because %s", err.Error())
		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	noTags := exportTransactionDataReq.TagFilter == models.TransactionNoTagFilterValue
	var tagFilters []*models.TransactionTagFilter

	if !noTags {
		tagFilters, err = models.ParseTransactionTagFilter(exportTransactionDataReq.TagFilter)

		if err != nil {
			log.Warnf(c, "[data_managements.getExportedFileContent] parse transaction tag filters error, because %s", err.Error())
			return nil, "", errs.Or(err, errs.ErrOperationFailed)
		}
	}

	maxTransactionTime := int64(math.MaxInt64)
	minTransactionTime := int64(0)

	if exportTransactionDataReq.MaxTime > 0 {
		maxTransactionTime = utils.GetMaxTransactionTimeFromUnixTime(exportTransactionDataReq.MaxTime)
	}

	if exportTransactionDataReq.MinTime > 0 {
		minTransactionTime = utils.GetMinTransactionTimeFromUnixTime(exportTransactionDataReq.MinTime)
	}

	allTransactions, err := a.transactions.GetAllSpecifiedTransactions(c, uid, maxTransactionTime, minTransactionTime, exportTransactionDataReq.Type, allCategoryIds, allAccountIds, tagFilters, noTags, exportTransactionDataReq.AmountFilter, exportTransactionDataReq.Keyword, pageCountForDataExport, true)

	if err != nil {
		log.Errorf(c, "[data_managements.getExportedFileContent] failed to all transactions user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	dataExporter := converters.GetTransactionDataExporter(fileType)

	if dataExporter == nil {
		return nil, "", errs.ErrNotImplemented
	}

	result, err := dataExporter.ToExportedContent(c, uid, allTransactions, accountMap, categoryMap, tagMap, tagIndexes)

	if err != nil {
		log.Errorf(c, "[data_managements.getExportedFileContent] failed to get exported data for \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	fileName := a.getFileName(user, clientTimezone, fileType)

	return result, fileName, nil
}

func (a *DataManagementsApi) getFileName(user *models.User, clientTimezone *time.Location, fileExtension string) string {
	currentTime := utils.FormatUnixTimeToLongDateTimeWithoutSecond(time.Now().Unix(), clientTimezone)
	currentTime = strings.Replace(currentTime, "-", "_", -1)
	currentTime = strings.Replace(currentTime, " ", "_", -1)
	currentTime = strings.Replace(currentTime, ":", "_", -1)

	return fmt.Sprintf("%s_%s.%s", user.Username, currentTime, fileExtension)
}
