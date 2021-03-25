package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/exporters"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/utils"
)

const pageCountForDataExport = 1000

// DataManagementsApi represents data management api
type DataManagementsApi struct {
	exporter     *exporters.CSVFileExporter
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
		exporter:     &exporters.CSVFileExporter{},
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

	accounts, err := a.accounts.GetAllAccountsByUid(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	categories, err := a.categories.GetAllCategoriesByUid(uid, 0, -1)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	tags, err := a.tags.GetAllTagsByUid(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get tags for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	tagIndexs, err := a.tags.GetAllTagIdsOfAllTransactions(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get tag index for user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	accountMap := a.accounts.GetAccountMapByList(accounts)
	categoryMap := a.categories.GetCategoryMapByList(categories)
	tagMap := a.tags.GetTagMapByList(tags)

	allTransactions, err := a.transactions.GetAllTransactions(uid, pageCountForDataExport, true)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to all transactions user \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.ErrOperationFailed
	}

	result, err := a.exporter.GetOutputContent(uid, timezone, allTransactions, accountMap, categoryMap, tagMap, tagIndexs)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get csv format exported data for \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	fileName := a.getFileName(timezone)

	return result, fileName, nil
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
	user, err := a.users.GetUserById(uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.WarnfWithRequestId(c, "[data_managements.ClearDataHandler] failed to get user for user \"uid:%d\", because %s", uid, err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if !a.users.IsPasswordEqualsUserPassword(clearDataReq.Password, user) {
		return nil, errs.ErrUserPasswordWrong
	}

	err = a.transactions.DeleteAllTransactions(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ClearDataHandler] failed to delete all transactions, because %s", err.Error())
		return nil, errs.ErrOperationFailed
	}

	err = a.categories.DeleteAllCategories(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ClearDataHandler] failed to delete all transaction categories, because %s", err.Error())
		return nil, errs.ErrOperationFailed
	}

	err = a.tags.DeleteAllTags(uid)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ClearDataHandler] failed to delete all transaction tags, because %s", err.Error())
		return nil, errs.ErrOperationFailed
	}

	log.InfofWithRequestId(c, "[data_managements.ClearDataHandler] user \"uid:%d\" has cleared all data", uid)
	return true, nil
}

func (a *DataManagementsApi) getFileName(timezone *time.Location) string {
	currentTime := utils.FormatUnixTimeToLongDateTimeWithoutSecond(time.Now().Unix(), timezone)
	currentTime = strings.Replace(currentTime, "-", "_", -1)
	currentTime = strings.Replace(currentTime, " ", "_", -1)
	currentTime = strings.Replace(currentTime, ":", "_", -1)

	return fmt.Sprintf("%s.csv", currentTime)
}
