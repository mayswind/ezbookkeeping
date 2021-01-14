package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/settings"
	"github.com/mayswind/lab/pkg/utils"
)

const pageCountForDataExport = 1000
const csvHeaderLine = "Time,Type,Category,Sub Category,Account,Amount,Account2,Account2 Amount,Tags,Comment\n"
const csvDataLineFormat = "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n"

// DataManagementsApi represents data management api
type DataManagementsApi struct {
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

	maxTime := utils.GetMaxTransactionTimeFromUnixTime(time.Now().Unix())
	var allTransactions []*models.Transaction

	for maxTime > 0 {
		transactions, err := a.transactions.GetAllTransactionsByMaxTime(uid, maxTime, pageCountForDataExport)

		if err != nil {
			log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get transactions earlier than \"%d\" for user \"uid:%d\", because %s", maxTime, uid, err.Error())
			return nil, "", errs.ErrOperationFailed
		}

		allTransactions = append(allTransactions, transactions...)

		if len(transactions) < pageCountForDataExport {
			maxTime = 0
			break
		}

		maxTime = transactions[len(transactions)-1].TransactionTime - 1
	}

	result, err := a.getCSVFormatData(c, allTransactions, accountMap, categoryMap, tagMap, tagIndexs)

	if err != nil {
		log.ErrorfWithRequestId(c, "[data_managements.ExportDataHandler] failed to get csv format exported data for \"uid:%d\", because %s", uid, err.Error())
		return nil, "", errs.Or(err, errs.ErrOperationFailed)
	}

	fileName := a.getFileName(c)

	return []byte(result), fileName, nil
}

func (a *DataManagementsApi) getCSVFormatData(c *core.Context, transactions []*models.Transaction, accountMap map[int64]*models.Account, categoryMap map[int64]*models.TransactionCategory, tagMap map[int64]*models.TransactionTag, allTagIndexs map[int64][]int64) (string, error) {
	var ret strings.Builder

	ret.Grow(len(transactions) * 100)
	ret.WriteString(csvHeaderLine)

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			continue
		}

		transactionTime := utils.FormatToLongDateTimeWithoutSecond(utils.ParseFromUnixTime(utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime)))
		transactionType := a.getTransactionTypeName(c, transaction.Type)
		category := a.getTransactionCategoryName(c, transaction.CategoryId, categoryMap)
		subCategory := a.getTransactionSubCategoryName(c, transaction.CategoryId, categoryMap)
		account := a.getAccountName(c, transaction.AccountId, accountMap)
		amount := a.getDisplayAmount(c, transaction.Amount)
		account2 := ""
		account2Amount := ""

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			account2 = a.getAccountName(c, transaction.RelatedAccountId, accountMap)
			account2Amount = a.getDisplayAmount(c, transaction.RelatedAccountAmount)
		}

		tags := a.getTags(c, transaction.TransactionId, allTagIndexs, tagMap)
		comment := a.getComment(c, transaction.Comment)

		ret.WriteString(fmt.Sprintf(csvDataLineFormat, transactionTime, transactionType, category, subCategory, account, amount, account2, account2Amount, tags, comment))
	}

	return ret.String(), nil
}

func (a *DataManagementsApi) getTransactionTypeName(c *core.Context, transactionDbType models.TransactionDbType) string {
	if transactionDbType == models.TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		return "Balance Modification"
	} else if transactionDbType == models.TRANSACTION_DB_TYPE_INCOME {
		return "Income"
	} else if transactionDbType == models.TRANSACTION_DB_TYPE_EXPENSE {
		return "Expense"
	} else if transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || transactionDbType == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		return "Transfer"
	} else {
		return ""
	}
}

func (a *DataManagementsApi) getTransactionCategoryName(c *core.Context, categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if !exists {
		return ""
	}

	if category.ParentCategoryId == 0 {
		return category.Name
	}

	parentCategory, exists := categoryMap[category.ParentCategoryId]

	if !exists {
		return ""
	}

	return parentCategory.Name
}

func (a *DataManagementsApi) getTransactionSubCategoryName(c *core.Context, categoryId int64, categoryMap map[int64]*models.TransactionCategory) string {
	category, exists := categoryMap[categoryId]

	if exists {
		return category.Name
	} else {
		return ""
	}
}

func (a *DataManagementsApi) getAccountName(c *core.Context, accountId int64, accountMap map[int64]*models.Account) string {
	account, exists := accountMap[accountId]

	if exists {
		return account.Name
	} else {
		return ""
	}
}

func (a *DataManagementsApi) getDisplayAmount(c *core.Context, amount int64) string {
	displayAmount := utils.Int64ToString(amount)
	integer := utils.SubString(displayAmount, 0, len(displayAmount)-2)
	decimals := utils.SubString(displayAmount, -2, 2)

	if integer == "" {
		integer = "0"
	} else if integer == "-" {
		integer = "-0"
	}

	if len(decimals) == 0 {
		decimals = "00"
	} else if len(decimals) == 1 {
		decimals = "0" + decimals
	}

	return integer + "." + decimals
}

func (a *DataManagementsApi) getTags(c *core.Context, transactionId int64, allTagIndexs map[int64][]int64, tagMap map[int64]*models.TransactionTag) string {
	tagIndexs, exists := allTagIndexs[transactionId]

	if !exists {
		return ""
	}

	var ret strings.Builder

	for i := 0; i < len(tagIndexs); i++ {
		if i > 0 {
			ret.WriteString(";")
		}

		tagIndex := tagIndexs[i]
		tag, exists := tagMap[tagIndex]

		if !exists {
			continue
		}

		ret.WriteString(tag.Name)
	}

	return ret.String()
}

func (a *DataManagementsApi) getComment(c *core.Context, comment string) string {
	comment = strings.Replace(comment, ",", " ", -1)
	comment = strings.Replace(comment, "\r\n", " ", -1)
	comment = strings.Replace(comment, "\n", " ", -1)

	return comment
}

func (a *DataManagementsApi) getFileName(c *core.Context) string {
	currentTime := utils.FormatToLongDateTimeWithoutSecond(time.Now())
	currentTime = strings.Replace(currentTime, "-", "_", -1)
	currentTime = strings.Replace(currentTime, " ", "_", -1)
	currentTime = strings.Replace(currentTime, ":", "_", -1)

	return fmt.Sprintf("%s.csv", currentTime)
}
