package api

import (
	"fmt"
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

const pageCountForDataExport = 1000

// DataManagementsApi represents data management api
type DataManagementsApi struct {
	exporter     *converters.EzBookKeepingCSVFileExporter
	tokens       *services.TokenService
	users        *services.UserService
	accounts     *services.AccountService
	transactions *services.TransactionService
	categories   *services.TransactionCategoryService
	tags         *services.TransactionTagService
}

// Initialize a data management api singleton instance
var (
	DataManagements = &DataManagementsApi{
		exporter:     &converters.EzBookKeepingCSVFileExporter{},
		tokens:       services.Tokens,
		users:        services.Users,
		accounts:     services.Accounts,
		transactions: services.Transactions,
		categories:   services.TransactionCategories,
		tags:         services.TransactionTags,
	}
)

// ExportDataHandler returns exported data in csv format
func (a *DataManagementsApi) ExportDataHandler(c *core.Context) ([]byte, string, *errs.Error) {
	if !settings.Container.Current.EnableDataExport {
		return nil, "", errs.ErrDataExportNotAllowed
	}

	timezone := time.Local
	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.WarnfWithRequestId(c, "[data_managements.ExportDataHandler] cannot get client timezone offset, because %s", err.Error())
	} else {
		timezone = time.FixedZone("Client Timezone", int(utcOffset)*60)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.WarnfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, "", errs.ErrUserNotFound
	}

	accounts, err := a.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	categories, err := a.categories.GetAllCategoriesByUid(c, uid, 0, -1)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	tags, err := a.tags.GetAllTagsByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	tagIndexs, err := a.tags.GetAllTagIdsOfAllTransactions(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get tag index for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	accountMap := a.accounts.GetAccountMapByList(accounts)
	categoryMap := a.categories.GetCategoryMapByList(categories)
	tagMap := a.tags.GetTagMapByList(tags)

	allTransactions, err := a.transactions.GetAllTransactions(c, uid, pageCountForDataExport, true)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to all transactions user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	result, err := a.exporter.ToExportedContent(uid, timezone, allTransactions, accountMap, categoryMap, tagMap, tagIndexs)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get csv format exported data for \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	fileName := a.getFileName(user, timezone)

	return result, fileName, nil
}

// DataStatisticsHandler returns user data statistics
func (a *DataManagementsApi) DataStatisticsHandler(c *core.Context) (interface{}, *errs.Error) {
	uid := c.GetCurrentUid()
	totalAccountCount, err := a.accounts.GetTotalAccountCountByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.DataStatisticsHandler] failed to get total account count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionCategoryCount, err := a.categories.GetTotalCategoryCountByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.DataStatisticsHandler] failed to get total transaction category count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionTagCount, err := a.tags.GetTotalTagCountByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.DataStatisticsHandler] failed to get total transaction tag count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	totalTransactionCount, err := a.transactions.GetTotalTransactionCountByUid(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.DataStatisticsHandler] failed to get total transaction count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	dataStatisticsResp := &models.DataStatisticsResponse{
		TotalAccountCount:             totalAccountCount,
		TotalTransactionCategoryCount: totalTransactionCategoryCount,
		TotalTransactionTagCount:      totalTransactionTagCount,
		TotalTransactionCount:         totalTransactionCount,
	}

	return dataStatisticsResp, nil
}

// ClearDataHandler deletes all user data
func (a *DataManagementsApi) ClearDataHandler(c *core.Context) (interface{}, *errs.Error) {
	var clearDataReq models.ClearDataRequest
	err := c.ShouldBindJSON(&clearDataReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[data_managements.ClearDataHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.WarnfWithRequestId(c, "[data_managements.ClearDataHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(clearDataReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	err = a.transactions.DeleteAllTransactions(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ClearDataHandler] failed to delete all transactions, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.categories.DeleteAllCategories(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ClearDataHandler] failed to delete all transaction categories, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	err = a.tags.DeleteAllTags(c, uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ClearDataHandler] failed to delete all transaction tags, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[data_managements.ClearDataHandler] user \"uid:%d\" has cleared all data", uid)
	return true, nil
}

func (a *DataManagementsApi) getFileName(user *models.User, timezone *time.Location) string {
	currentTime := utils.FormatUnixTimeToLongDateTimeWithoutSecond(time.Now().Unix(), timezone)
	currentTime = strings.Replace(currentTime, "-", "_", -1)
	currentTime = strings.Replace(currentTime, " ", "_", -1)
	currentTime = strings.Replace(currentTime, ":", "_", -1)

	return fmt.Sprintf("%s_%s.csv", user.Username, currentTime)
}
