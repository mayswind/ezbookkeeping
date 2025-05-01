package api

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/mayswind/ezbookkeeping/pkg/converters"
	"github.com/mayswind/ezbookkeeping/pkg/converters/converter"
	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const maximumTagsCountOfTransaction = 10
const maximumPicturesCountOfTransaction = 10

// TransactionsApi represents transaction api
type TransactionsApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	transactions          *services.TransactionService
	transactionCategories *services.TransactionCategoryService
	transactionTags       *services.TransactionTagService
	transactionPictures   *services.TransactionPictureService
	accounts              *services.AccountService
	users                 *services.UserService
}

// Initialize a transaction api singleton instance
var (
	Transactions = &TransactionsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		transactions:          services.Transactions,
		transactionCategories: services.TransactionCategories,
		transactionTags:       services.TransactionTags,
		transactionPictures:   services.TransactionPictures,
		accounts:              services.Accounts,
		users:                 services.Users,
	}
)

// TransactionCountHandler returns transaction total count of current user
func (a *TransactionsApi) TransactionCountHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionCountReq models.TransactionCountRequest
	err := c.ShouldBindQuery(&transactionCountReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionCountHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	allAccountIds, err := a.getAccountOrSubAccountIds(c, transactionCountReq.AccountIds, uid)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionCountHandler] get account error, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	allCategoryIds, err := a.getCategoryOrSubCategoryIds(c, transactionCountReq.CategoryIds, uid)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionCountHandler] get transaction category error, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var allTagIds []int64
	noTags := transactionCountReq.TagIds == "none"

	if !noTags {
		allTagIds, err = a.getTagIds(transactionCountReq.TagIds)

		if err != nil {
			log.Warnf(c, "[transactions.TransactionCountHandler] get transaction tag ids error, because %s", err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	totalCount, err := a.transactions.GetTransactionCount(c, uid, transactionCountReq.MaxTime, transactionCountReq.MinTime, transactionCountReq.Type, allCategoryIds, allAccountIds, allTagIds, noTags, transactionCountReq.TagFilterType, transactionCountReq.AmountFilter, transactionCountReq.Keyword)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionCountHandler] failed to get transaction count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	countResp := &models.TransactionCountResponse{
		TotalCount: totalCount,
	}

	return countResp, nil
}

// TransactionListHandler returns transaction list of current user
func (a *TransactionsApi) TransactionListHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionListReq models.TransactionListByMaxTimeRequest
	err := c.ShouldBindQuery(&transactionListReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionListHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionListHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	allAccountIds, err := a.getAccountOrSubAccountIds(c, transactionListReq.AccountIds, uid)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionListHandler] get account error, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	allCategoryIds, err := a.getCategoryOrSubCategoryIds(c, transactionListReq.CategoryIds, uid)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionListHandler] get transaction category error, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var allTagIds []int64
	noTags := transactionListReq.TagIds == "none"

	if !noTags {
		allTagIds, err = a.getTagIds(transactionListReq.TagIds)

		if err != nil {
			log.Warnf(c, "[transactions.TransactionListHandler] get transaction tag ids error, because %s", err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	var totalCount int64

	if transactionListReq.WithCount {
		totalCount, err = a.transactions.GetTransactionCount(c, uid, transactionListReq.MaxTime, transactionListReq.MinTime, transactionListReq.Type, allCategoryIds, allAccountIds, allTagIds, noTags, transactionListReq.TagFilterType, transactionListReq.AmountFilter, transactionListReq.Keyword)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionListHandler] failed to get transaction count for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	transactions, err := a.transactions.GetTransactionsByMaxTime(c, uid, transactionListReq.MaxTime, transactionListReq.MinTime, transactionListReq.Type, allCategoryIds, allAccountIds, allTagIds, noTags, transactionListReq.TagFilterType, transactionListReq.AmountFilter, transactionListReq.Keyword, transactionListReq.Page, transactionListReq.Count, true, true)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionListHandler] failed to get transactions earlier than \"%d\" for user \"uid:%d\", because %s", transactionListReq.MaxTime, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	hasMore := false
	var nextTimeSequenceId *int64

	if len(transactions) > int(transactionListReq.Count) {
		hasMore = true
		nextTimeSequenceId = &transactions[transactionListReq.Count].TransactionTime
		transactions = transactions[:transactionListReq.Count]
	}

	transactionResult, err := a.getTransactionResponseListResult(c, user, transactions, utcOffset, transactionListReq.WithPictures, transactionListReq.TrimAccount, transactionListReq.TrimCategory, transactionListReq.TrimTag)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionListHandler] failed to assemble transaction result for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	transactionResps := &models.TransactionInfoPageWrapperResponse{
		Items: transactionResult,
	}

	if hasMore {
		transactionResps.NextTimeSequenceId = nextTimeSequenceId
	}

	if transactionListReq.WithCount {
		transactionResps.TotalCount = &totalCount
	}

	return transactionResps, nil
}

// TransactionMonthListHandler returns all transaction list of current user by month
func (a *TransactionsApi) TransactionMonthListHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionListReq models.TransactionListInMonthByPageRequest
	err := c.ShouldBindQuery(&transactionListReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionMonthListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionMonthListHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionMonthListHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	allAccountIds, err := a.getAccountOrSubAccountIds(c, transactionListReq.AccountIds, uid)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionMonthListHandler] get account error, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	allCategoryIds, err := a.getCategoryOrSubCategoryIds(c, transactionListReq.CategoryIds, uid)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionMonthListHandler] get transaction category error, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var allTagIds []int64
	noTags := transactionListReq.TagIds == "none"

	if !noTags {
		allTagIds, err = a.getTagIds(transactionListReq.TagIds)

		if err != nil {
			log.Warnf(c, "[transactions.TransactionMonthListHandler] get transaction tag ids error, because %s", err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	transactions, err := a.transactions.GetTransactionsInMonthByPage(c, uid, transactionListReq.Year, transactionListReq.Month, transactionListReq.Type, allCategoryIds, allAccountIds, allTagIds, noTags, transactionListReq.TagFilterType, transactionListReq.AmountFilter, transactionListReq.Keyword)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionMonthListHandler] failed to get transactions in month \"%d-%d\" for user \"uid:%d\", because %s", transactionListReq.Year, transactionListReq.Month, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	transactionResult, err := a.getTransactionResponseListResult(c, user, transactions, utcOffset, transactionListReq.WithPictures, transactionListReq.TrimAccount, transactionListReq.TrimCategory, transactionListReq.TrimTag)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionMonthListHandler] failed to assemble transaction result for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	transactionResps := &models.TransactionInfoPageWrapperResponse2{
		Items:      transactionResult,
		TotalCount: int64(transactionResult.Len()),
	}

	return transactionResps, nil
}

// TransactionStatisticsHandler returns transaction statistics of current user
func (a *TransactionsApi) TransactionStatisticsHandler(c *core.WebContext) (any, *errs.Error) {
	var statisticReq models.TransactionStatisticRequest
	err := c.ShouldBindQuery(&statisticReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionStatisticsHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionStatisticsHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	var allTagIds []int64
	noTags := statisticReq.TagIds == "none"

	if !noTags {
		allTagIds, err = a.getTagIds(statisticReq.TagIds)

		if err != nil {
			log.Warnf(c, "[transactions.TransactionStatisticsHandler] get transaction tag ids error, because %s", err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	uid := c.GetCurrentUid()
	totalAmounts, err := a.transactions.GetAccountsAndCategoriesTotalIncomeAndExpense(c, uid, statisticReq.StartTime, statisticReq.EndTime, allTagIds, noTags, statisticReq.TagFilterType, utcOffset, statisticReq.UseTransactionTimezone)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionStatisticsHandler] failed to get accounts and categories total income and expense for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	statisticResp := &models.TransactionStatisticResponse{
		StartTime: statisticReq.StartTime,
		EndTime:   statisticReq.EndTime,
	}

	statisticResp.Items = make([]*models.TransactionStatisticResponseItem, len(totalAmounts))

	for i := 0; i < len(totalAmounts); i++ {
		totalAmountItem := totalAmounts[i]
		statisticResp.Items[i] = &models.TransactionStatisticResponseItem{
			CategoryId:  totalAmountItem.CategoryId,
			AccountId:   totalAmountItem.AccountId,
			TotalAmount: totalAmountItem.Amount,
		}
	}

	return statisticResp, nil
}

// TransactionStatisticsTrendsHandler returns transaction statistics trends of current user
func (a *TransactionsApi) TransactionStatisticsTrendsHandler(c *core.WebContext) (any, *errs.Error) {
	var statisticTrendsReq models.TransactionStatisticTrendsRequest
	err := c.ShouldBindQuery(&statisticTrendsReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionStatisticsTrendsHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionStatisticsTrendsHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	startYear, startMonth, endYear, endMonth, err := statisticTrendsReq.GetNumericYearMonthRange()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionStatisticsTrendsHandler] cannot parse year month, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var allTagIds []int64
	noTags := statisticTrendsReq.TagIds == "none"

	if !noTags {
		allTagIds, err = a.getTagIds(statisticTrendsReq.TagIds)

		if err != nil {
			log.Warnf(c, "[transactions.TransactionStatisticsTrendsHandler] get transaction tag ids error, because %s", err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	uid := c.GetCurrentUid()
	allMonthlyTotalAmounts, err := a.transactions.GetAccountsAndCategoriesMonthlyIncomeAndExpense(c, uid, startYear, startMonth, endYear, endMonth, allTagIds, noTags, statisticTrendsReq.TagFilterType, utcOffset, statisticTrendsReq.UseTransactionTimezone)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionStatisticsTrendsHandler] failed to get accounts and categories total income and expense for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	statisticTrendsResp := make(models.TransactionStatisticTrendsResponseItemSlice, 0, len(allMonthlyTotalAmounts))

	for yearMonth, monthlyTotalAmounts := range allMonthlyTotalAmounts {
		monthlyStatisticResp := &models.TransactionStatisticTrendsResponseItem{
			Year:  yearMonth / 100,
			Month: yearMonth % 100,
			Items: make([]*models.TransactionStatisticResponseItem, len(monthlyTotalAmounts)),
		}

		for i := 0; i < len(monthlyTotalAmounts); i++ {
			totalAmountItem := monthlyTotalAmounts[i]
			monthlyStatisticResp.Items[i] = &models.TransactionStatisticResponseItem{
				CategoryId:  totalAmountItem.CategoryId,
				AccountId:   totalAmountItem.AccountId,
				TotalAmount: totalAmountItem.Amount,
			}
		}

		statisticTrendsResp = append(statisticTrendsResp, monthlyStatisticResp)
	}

	sort.Sort(statisticTrendsResp)

	return statisticTrendsResp, nil
}

// TransactionAmountsHandler returns transaction amounts of current user
func (a *TransactionsApi) TransactionAmountsHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionAmountsReq models.TransactionAmountsRequest
	err := c.ShouldBindQuery(&transactionAmountsReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionAmountsHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	requestItems, err := transactionAmountsReq.GetTransactionAmountsRequestItems()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionAmountsHandler] get request item failed, because %s", err.Error())
		return nil, errs.ErrQueryItemsInvalid
	}

	if len(requestItems) < 1 {
		log.Warnf(c, "[transactions.TransactionAmountsHandler] parse request failed, because there are no valid items")
		return nil, errs.ErrQueryItemsEmpty
	}

	if len(requestItems) > 20 {
		log.Warnf(c, "[transactions.TransactionAmountsHandler] parse request failed, because there are too many items")
		return nil, errs.ErrQueryItemsTooMuch
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionAmountsHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	uid := c.GetCurrentUid()

	accounts, err := a.accounts.GetAllAccountsByUid(c, uid)
	accountMap := a.accounts.GetAccountMapByList(accounts)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionAmountsHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	amountsResp := orderedmap.New[string, *models.TransactionAmountsResponseItem]()

	for i := 0; i < len(requestItems); i++ {
		requestItem := requestItems[i]

		incomeAmounts, expenseAmounts, err := a.transactions.GetAccountsTotalIncomeAndExpense(c, uid, requestItem.StartTime, requestItem.EndTime, utcOffset, transactionAmountsReq.UseTransactionTimezone)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionAmountsHandler] failed to get transaction amounts item for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		amountsMap := make(map[string]*models.TransactionAmountsResponseItemAmountInfo)

		for accountId, incomeAmount := range incomeAmounts {
			account, exists := accountMap[accountId]

			if !exists {
				log.Warnf(c, "[transactions.TransactionAmountsHandler] cannot find account for account \"id:%d\" of user \"uid:%d\"", accountId, uid)
				continue
			}

			totalAmounts, exists := amountsMap[account.Currency]

			if !exists {
				totalAmounts = &models.TransactionAmountsResponseItemAmountInfo{
					Currency:      account.Currency,
					IncomeAmount:  0,
					ExpenseAmount: 0,
				}
			}

			totalAmounts.IncomeAmount += incomeAmount
			amountsMap[account.Currency] = totalAmounts
		}

		for accountId, expenseAmount := range expenseAmounts {
			account, exists := accountMap[accountId]

			if !exists {
				log.Warnf(c, "[transactions.TransactionAmountsHandler] cannot find account for account \"id:%d\" of user \"uid:%d\"", accountId, uid)
				continue
			}

			totalAmounts, exists := amountsMap[account.Currency]

			if !exists {
				totalAmounts = &models.TransactionAmountsResponseItemAmountInfo{
					Currency:      account.Currency,
					IncomeAmount:  0,
					ExpenseAmount: 0,
				}
			}

			totalAmounts.ExpenseAmount += expenseAmount
			amountsMap[account.Currency] = totalAmounts
		}

		allTotalAmounts := make(models.TransactionAmountsResponseItemAmountInfoSlice, 0)

		for _, totalAmounts := range amountsMap {
			allTotalAmounts = append(allTotalAmounts, totalAmounts)
		}

		sort.Sort(allTotalAmounts)

		amountsResp.Set(requestItem.Name, &models.TransactionAmountsResponseItem{
			StartTime: requestItem.StartTime,
			EndTime:   requestItem.EndTime,
			Amounts:   allTotalAmounts,
		})
	}

	return amountsResp, nil
}

// TransactionGetHandler returns one specific transaction of current user
func (a *TransactionsApi) TransactionGetHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionGetReq models.TransactionGetRequest
	err := c.ShouldBindQuery(&transactionGetReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionGetHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionGetHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	transaction, err := a.transactions.GetTransactionByTransactionId(c, uid, transactionGetReq.Id)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionGetHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", transactionGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		transaction = a.transactions.GetRelatedTransferTransaction(transaction)
	}

	accountIds := make([]int64, 0, 2)
	accountIds = append(accountIds, transaction.AccountId)

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		accountIds = append(accountIds, transaction.RelatedAccountId)
		accountIds = utils.ToUniqueInt64Slice(accountIds)
	}

	accountMap, err := a.accounts.GetAccountsByAccountIds(c, uid, accountIds)

	if _, exists := accountMap[transaction.AccountId]; !exists {
		log.Warnf(c, "[transactions.TransactionGetHandler] account of transaction \"id:%d\" does not exist for user \"uid:%d\"", transaction.TransactionId, uid)
		return nil, errs.ErrTransactionNotFound
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if _, exists := accountMap[transaction.RelatedAccountId]; !exists {
			log.Warnf(c, "[transactions.TransactionGetHandler] related account of transaction \"id:%d\" does not exist for user \"uid:%d\"", transaction.TransactionId, uid)
			return nil, errs.ErrTransactionNotFound
		}
	}

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(c, uid, []int64{transaction.TransactionId})

	if err != nil {
		log.Errorf(c, "[transactions.TransactionGetHandler] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var category *models.TransactionCategory
	var tagMap map[int64]*models.TransactionTag
	var pictureInfos []*models.TransactionPictureInfo

	if !transactionGetReq.TrimCategory {
		category, err = a.transactionCategories.GetCategoryByCategoryId(c, uid, transaction.CategoryId)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionGetHandler] failed to get transactions category for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	if !transactionGetReq.TrimTag {
		tagMap, err = a.transactionTags.GetTagsByTagIds(c, uid, utils.ToUniqueInt64Slice(a.getTransactionTagIds(allTransactionTagIds)))

		if err != nil {
			log.Errorf(c, "[transactions.TransactionGetHandler] failed to get transactions tags for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	if transactionGetReq.WithPictures && a.CurrentConfig().EnableTransactionPictures {
		pictureInfos, err = a.transactionPictures.GetPictureInfosByTransactionId(c, uid, transaction.TransactionId)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionGetHandler] failed to get transactions pictures for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}
	}

	transactionEditable := transaction.IsEditable(user, utcOffset, accountMap[transaction.AccountId], accountMap[transaction.RelatedAccountId])
	transactionTagIds := allTransactionTagIds[transaction.TransactionId]
	transactionResp := transaction.ToTransactionInfoResponse(transactionTagIds, transactionEditable)

	if !transactionGetReq.TrimAccount {
		if sourceAccount := accountMap[transaction.AccountId]; sourceAccount != nil {
			transactionResp.SourceAccount = sourceAccount.ToAccountInfoResponse()
		}

		if destinationAccount := accountMap[transaction.RelatedAccountId]; destinationAccount != nil {
			transactionResp.DestinationAccount = destinationAccount.ToAccountInfoResponse()
		}
	}

	if !transactionGetReq.TrimCategory {
		if category != nil {
			transactionResp.Category = category.ToTransactionCategoryInfoResponse()
		}
	}

	if !transactionGetReq.TrimTag {
		transactionResp.Tags = a.getTransactionTagInfoResponses(transactionTagIds, tagMap)
	}

	if transactionGetReq.WithPictures && a.CurrentConfig().EnableTransactionPictures {
		transactionResp.Pictures = a.GetTransactionPictureInfoResponseList(pictureInfos)
	}

	return transactionResp, nil
}

// TransactionCreateHandler saves a new transaction by request parameters for current user
func (a *TransactionsApi) TransactionCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionCreateReq models.TransactionCreateRequest
	err := c.ShouldBindJSON(&transactionCreateReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	tagIds, err := utils.StringArrayToInt64Array(transactionCreateReq.TagIds)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionCreateHandler] parse tag ids failed, because %s", err.Error())
		return nil, errs.ErrTransactionTagIdInvalid
	}

	if len(tagIds) > maximumTagsCountOfTransaction {
		return nil, errs.ErrTransactionHasTooManyTags
	}

	pictureIds, err := utils.StringArrayToInt64Array(transactionCreateReq.PictureIds)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionCreateHandler] parse picture ids failed, because %s", err.Error())
		return nil, errs.ErrTransactionPictureIdInvalid
	}

	if len(pictureIds) > maximumPicturesCountOfTransaction {
		return nil, errs.ErrTransactionHasTooManyPictures
	}

	if transactionCreateReq.Type < models.TRANSACTION_TYPE_MODIFY_BALANCE || transactionCreateReq.Type > models.TRANSACTION_TYPE_TRANSFER {
		log.Warnf(c, "[transactions.TransactionCreateHandler] transaction type is invalid")
		return nil, errs.ErrTransactionTypeInvalid
	}

	if transactionCreateReq.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE && transactionCreateReq.CategoryId > 0 {
		log.Warnf(c, "[transactions.TransactionCreateHandler] balance modification transaction cannot set category id")
		return nil, errs.ErrBalanceModificationTransactionCannotSetCategory
	}

	if transactionCreateReq.Type != models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.DestinationAccountId != 0 {
		log.Warnf(c, "[transactions.TransactionCreateHandler] non-transfer transaction destination account cannot be set")
		return nil, errs.ErrTransactionDestinationAccountCannotBeSet
	} else if transactionCreateReq.Type == models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.SourceAccountId == transactionCreateReq.DestinationAccountId {
		log.Warnf(c, "[transactions.TransactionCreateHandler] transfer transaction source account must not be destination account")
		return nil, errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
	}

	if transactionCreateReq.Type != models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.DestinationAmount != 0 {
		log.Warnf(c, "[transactions.TransactionCreateHandler] non-transfer transaction destination amount cannot be set")
		return nil, errs.ErrTransactionDestinationAmountCannotBeSet
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionCreateHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	transaction := a.createNewTransactionModel(uid, &transactionCreateReq, c.ClientIP())
	transactionEditable := user.CanEditTransactionByTransactionTime(transaction.TransactionTime, transactionCreateReq.UtcOffset)

	if !transactionEditable {
		return nil, errs.ErrCannotCreateTransactionWithThisTransactionTime
	}

	var pictureInfos []*models.TransactionPictureInfo

	if len(pictureIds) > 0 {
		pictureInfos, err = a.transactionPictures.GetNewPictureInfosByPictureIds(c, uid, pictureIds)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionCreateHandler] failed to get transactions pictures for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		notExistsPictureIds := utils.Int64SliceMinus(pictureIds, a.transactionPictures.GetTransactionPictureIds(pictureInfos))

		if len(notExistsPictureIds) > 0 {
			log.Errorf(c, "[transactions.TransactionCreateHandler] some pictures \"ids:%s\" does not exists for user \"uid:%d\"", strings.Join(utils.Int64ArrayToStringArray(notExistsPictureIds), ","), uid)
			return nil, errs.ErrTransactionPictureNotFound
		}
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && transactionCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, transactionCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[transactions.TransactionCreateHandler] another transaction \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			transactionId, err := utils.StringToInt64(remark)

			if err == nil {
				transaction, err = a.transactions.GetTransactionByTransactionId(c, uid, transactionId)

				if err != nil {
					log.Errorf(c, "[transactions.TransactionCreateHandler] failed to get existed transaction \"id:%d\" for user \"uid:%d\", because %s", transactionId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				transactionResp := transaction.ToTransactionInfoResponse(tagIds, transactionEditable)
				transactionResp.Pictures = a.GetTransactionPictureInfoResponseList(pictureInfos)

				return transactionResp, nil
			}
		}
	}

	err = a.transactions.CreateTransaction(c, transaction, tagIds, pictureIds)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionCreateHandler] failed to create transaction \"id:%d\" for user \"uid:%d\", because %s", transaction.TransactionId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transactions.TransactionCreateHandler] user \"uid:%d\" has created a new transaction \"id:%d\" successfully", uid, transaction.TransactionId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, transactionCreateReq.ClientSessionId, utils.Int64ToString(transaction.TransactionId))
	transactionResp := transaction.ToTransactionInfoResponse(tagIds, transactionEditable)
	transactionResp.Pictures = a.GetTransactionPictureInfoResponseList(pictureInfos)

	return transactionResp, nil
}

// TransactionModifyHandler saves an existed transaction by request parameters for current user
func (a *TransactionsApi) TransactionModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionModifyReq models.TransactionModifyRequest
	err := c.ShouldBindJSON(&transactionModifyReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	tagIds, err := utils.StringArrayToInt64Array(transactionModifyReq.TagIds)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionModifyHandler] parse tag ids failed, because %s", err.Error())
		return nil, errs.ErrTransactionTagIdInvalid
	}

	if len(tagIds) > maximumTagsCountOfTransaction {
		return nil, errs.ErrTransactionHasTooManyTags
	}

	pictureIds, err := utils.StringArrayToInt64Array(transactionModifyReq.PictureIds)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionModifyHandler] parse picture ids failed, because %s", err.Error())
		return nil, errs.ErrTransactionPictureIdInvalid
	}

	if len(pictureIds) > maximumPicturesCountOfTransaction {
		return nil, errs.ErrTransactionHasTooManyPictures
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionModifyHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	transaction, err := a.transactions.GetTransactionByTransactionId(c, uid, transactionModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionModifyHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", transactionModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		log.Warnf(c, "[transactions.TransactionModifyHandler] cannot modify transaction \"id:%d\" for user \"uid:%d\", because transaction type is transfer in", transactionModifyReq.Id, uid)
		return nil, errs.ErrTransactionTypeInvalid
	}

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(c, uid, []int64{transaction.TransactionId})

	if err != nil {
		log.Errorf(c, "[transactions.TransactionModifyHandler] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	transactionTagIds := allTransactionTagIds[transaction.TransactionId]

	if transactionTagIds == nil {
		transactionTagIds = make([]int64, 0, 0)
	}

	transactionPictureInfos, err := a.transactionPictures.GetPictureInfosByTransactionId(c, uid, transaction.TransactionId)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionModifyHandler] failed to get transaction picture infos for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	transactionPictureIds := a.transactionPictures.GetTransactionPictureIds(transactionPictureInfos)

	newTransaction := &models.Transaction{
		TransactionId:     transaction.TransactionId,
		Uid:               uid,
		CategoryId:        transactionModifyReq.CategoryId,
		TransactionTime:   utils.GetMinTransactionTimeFromUnixTime(transactionModifyReq.Time),
		TimezoneUtcOffset: transactionModifyReq.UtcOffset,
		AccountId:         transactionModifyReq.SourceAccountId,
		Amount:            transactionModifyReq.SourceAmount,
		HideAmount:        transactionModifyReq.HideAmount,
		Comment:           transactionModifyReq.Comment,
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
		newTransaction.RelatedAccountId = transactionModifyReq.DestinationAccountId
		newTransaction.RelatedAccountAmount = transactionModifyReq.DestinationAmount
	}

	if transactionModifyReq.GeoLocation != nil {
		newTransaction.GeoLongitude = transactionModifyReq.GeoLocation.Longitude
		newTransaction.GeoLatitude = transactionModifyReq.GeoLocation.Latitude
	}

	if newTransaction.CategoryId == transaction.CategoryId &&
		utils.GetUnixTimeFromTransactionTime(newTransaction.TransactionTime) == utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime) &&
		newTransaction.TimezoneUtcOffset == transaction.TimezoneUtcOffset &&
		newTransaction.AccountId == transaction.AccountId &&
		newTransaction.Amount == transaction.Amount &&
		(transaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_OUT || newTransaction.RelatedAccountId == transaction.RelatedAccountId) &&
		(transaction.Type != models.TRANSACTION_DB_TYPE_TRANSFER_OUT || newTransaction.RelatedAccountAmount == transaction.RelatedAccountAmount) &&
		newTransaction.HideAmount == transaction.HideAmount &&
		newTransaction.Comment == transaction.Comment &&
		newTransaction.GeoLongitude == transaction.GeoLongitude &&
		newTransaction.GeoLatitude == transaction.GeoLatitude &&
		utils.Int64SliceEquals(tagIds, transactionTagIds) &&
		utils.Int64SliceEquals(pictureIds, transactionPictureIds) {
		return nil, errs.ErrNothingWillBeUpdated
	}

	transactionEditable := user.CanEditTransactionByTransactionTime(transaction.TransactionTime, transaction.TimezoneUtcOffset)
	newTransactionEditable := user.CanEditTransactionByTransactionTime(newTransaction.TransactionTime, transactionModifyReq.UtcOffset)

	if !transactionEditable || !newTransactionEditable {
		return nil, errs.ErrCannotModifyTransactionWithThisTransactionTime
	}

	var addTransactionTagIds []int64
	var removeTransactionTagIds []int64

	if !utils.Int64SliceEquals(tagIds, transactionTagIds) {
		removeTransactionTagIds = transactionTagIds
		addTransactionTagIds = tagIds
	}

	addTransactionPictureIds := utils.Int64SliceMinus(pictureIds, transactionPictureIds)
	removeTransactionPictureIds := utils.Int64SliceMinus(transactionPictureIds, pictureIds)
	var newPictureInfos []*models.TransactionPictureInfo

	if !utils.Int64SliceEquals(pictureIds, transactionPictureIds) {
		oldAndNewPictureIds := transactionPictureIds
		oldAndNewPictureInfoMap := a.transactionPictures.GetPictureInfoMapByList(transactionPictureInfos)

		if len(addTransactionPictureIds) > 0 {
			addPictureInfos, err := a.transactionPictures.GetNewPictureInfosByPictureIds(c, uid, addTransactionPictureIds)

			if err != nil {
				log.Errorf(c, "[transactions.TransactionModifyHandler] failed to get transactions pictures for user \"uid:%d\", because %s", uid, err.Error())
				return nil, errs.Or(err, errs.ErrOperationFailed)
			}

			oldAndNewPictureIds = append(oldAndNewPictureIds, a.transactionPictures.GetTransactionPictureIds(addPictureInfos)...)
			notExistsPictureIds := utils.Int64SliceMinus(pictureIds, oldAndNewPictureIds)

			if len(notExistsPictureIds) > 0 {
				log.Errorf(c, "[transactions.TransactionModifyHandler] some pictures \"ids:%s\" does not exists for user \"uid:%d\"", strings.Join(utils.Int64ArrayToStringArray(notExistsPictureIds), ","), uid)
				return nil, errs.ErrTransactionPictureNotFound
			}

			for i := 0; i < len(addPictureInfos); i++ {
				oldAndNewPictureInfoMap[addPictureInfos[i].PictureId] = addPictureInfos[i]
			}
		}

		for i := 0; i < len(pictureIds); i++ {
			pictureId := pictureIds[i]
			pictureInfo, exists := oldAndNewPictureInfoMap[pictureId]

			if exists {
				newPictureInfos = append(newPictureInfos, pictureInfo)
			}
		}
	}

	err = a.transactions.ModifyTransaction(c, newTransaction, len(transactionTagIds), addTransactionTagIds, removeTransactionTagIds, addTransactionPictureIds, removeTransactionPictureIds)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionModifyHandler] failed to update transaction \"id:%d\" for user \"uid:%d\", because %s", transactionModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transactions.TransactionModifyHandler] user \"uid:%d\" has updated transaction \"id:%d\" successfully", uid, transactionModifyReq.Id)

	newTransaction.Type = transaction.Type
	newTransactionResp := newTransaction.ToTransactionInfoResponse(tagIds, transactionEditable)
	newTransactionResp.Pictures = a.GetTransactionPictureInfoResponseList(newPictureInfos)

	return newTransactionResp, nil
}

// TransactionDeleteHandler deletes an existed transaction by request parameters for current user
func (a *TransactionsApi) TransactionDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionDeleteReq models.TransactionDeleteRequest
	err := c.ShouldBindJSON(&transactionDeleteReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionDeleteHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionDeleteHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	transaction, err := a.transactions.GetTransactionByTransactionId(c, uid, transactionDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionDeleteHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", transactionDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
		log.Warnf(c, "[transactions.TransactionDeleteHandler] cannot delete transaction \"id:%d\" for user \"uid:%d\", because transaction type is transfer in", transactionDeleteReq.Id, uid)
		return nil, errs.ErrTransactionTypeInvalid
	}

	transactionEditable := user.CanEditTransactionByTransactionTime(transaction.TransactionTime, utcOffset)

	if !transactionEditable {
		return nil, errs.ErrCannotDeleteTransactionWithThisTransactionTime
	}

	err = a.transactions.DeleteTransaction(c, uid, transactionDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionDeleteHandler] failed to delete transaction \"id:%d\" for user \"uid:%d\", because %s", transactionDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transactions.TransactionDeleteHandler] user \"uid:%d\" has deleted transaction \"id:%d\"", uid, transactionDeleteReq.Id)
	return true, nil
}

// TransactionParseImportDsvFileDataHandler returns the parsed file data by request parameters for current user
func (a *TransactionsApi) TransactionParseImportDsvFileDataHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	form, err := c.MultipartForm()

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportDsvFileDataHandler] failed to get multi-part form data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	fileTypes := form.Value["fileType"]

	if len(fileTypes) < 1 || fileTypes[0] == "" {
		return nil, errs.ErrImportFileTypeIsEmpty
	}

	fileType := fileTypes[0]

	if !converters.IsCustomDelimiterSeparatedValuesFileType(fileType) {
		return nil, errs.Or(err, errs.ErrImportFileTypeNotSupported)
	}

	fileEncodings := form.Value["fileEncoding"]

	if len(fileEncodings) < 1 || fileEncodings[0] == "" {
		return nil, errs.ErrImportFileEncodingIsEmpty
	}

	fileEncoding := fileEncodings[0]
	dataParser, err := converters.CreateNewDelimiterSeparatedValuesDataParser(fileType, fileEncoding)

	if err != nil {
		return nil, errs.Or(err, errs.ErrImportFileTypeNotSupported)
	}

	importFiles := form.File["file"]

	if len(importFiles) < 1 {
		log.Warnf(c, "[transactions.TransactionParseImportDsvFileDataHandler] there is no import file in request for user \"uid:%d\"", uid)
		return nil, errs.ErrNoFilesUpload
	}

	if importFiles[0].Size < 1 {
		log.Warnf(c, "[transactions.TransactionParseImportDsvFileDataHandler] the size of import file in request is zero for user \"uid:%d\"", uid)
		return nil, errs.ErrUploadedFileEmpty
	}

	if importFiles[0].Size > int64(a.CurrentConfig().MaxImportFileSize) {
		log.Warnf(c, "[transactions.TransactionParseImportDsvFileDataHandler] the upload file size \"%d\" exceeds the maximum size \"%d\" of import file for user \"uid:%d\"", importFiles[0].Size, a.CurrentConfig().MaxImportFileSize, uid)
		return nil, errs.ErrExceedMaxUploadFileSize
	}

	importFile, err := importFiles[0].Open()

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportDsvFileDataHandler] failed to get import file from request for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	defer importFile.Close()
	fileData, err := io.ReadAll(importFile)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportDsvFileDataHandler] failed to read import file data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	allLines, err := dataParser.ParseDsvFileLines(c, fileData)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportDsvFileDataHandler] failed to parse import file data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return allLines, nil
}

// TransactionParseImportFileHandler returns the parsed transaction data by request parameters for current user
func (a *TransactionsApi) TransactionParseImportFileHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	form, err := c.MultipartForm()

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to get multi-part form data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrParameterInvalid
	}

	utcOffset, err := c.GetClientTimezoneOffset()

	if err != nil {
		log.Warnf(c, "[transactions.TransactionParseImportFileHandler] cannot get client timezone offset, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	fileTypes := form.Value["fileType"]

	if len(fileTypes) < 1 || fileTypes[0] == "" {
		return nil, errs.ErrImportFileTypeIsEmpty
	}

	fileType := fileTypes[0]

	var dataImporter converter.TransactionDataImporter

	if converters.IsCustomDelimiterSeparatedValuesFileType(fileType) {
		fileEncodings := form.Value["fileEncoding"]

		if len(fileEncodings) < 1 || fileEncodings[0] == "" {
			return nil, errs.ErrImportFileEncodingIsEmpty
		}

		fileEncoding := fileEncodings[0]

		columnMappings := form.Value["columnMapping"]

		if len(columnMappings) < 1 || columnMappings[0] == "" {
			return nil, errs.ErrImportFileColumnMappingInvalid
		}

		var columnIndexMapping = map[datatable.TransactionDataTableColumn]int{}
		err = json.Unmarshal([]byte(columnMappings[0]), &columnIndexMapping)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to parse column mapping for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrImportFileColumnMappingInvalid
		}

		transactionTypeMappings := form.Value["transactionTypeMapping"]

		if len(transactionTypeMappings) < 1 || transactionTypeMappings[0] == "" {
			return nil, errs.ErrImportFileTransactionTypeMappingInvalid
		}

		var transactionTypeNameMapping = map[string]models.TransactionType{}
		err = json.Unmarshal([]byte(transactionTypeMappings[0]), &transactionTypeNameMapping)

		if err != nil {
			log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to parse transaction type mapping for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrImportFileTransactionTypeMappingInvalid
		}

		hasHeaderLines := form.Value["hasHeaderLine"]
		hasHeaderLine := false

		if len(hasHeaderLines) > 0 {
			hasHeaderLine = hasHeaderLines[0] == "true"
		}

		timeFormats := form.Value["timeFormat"]

		if len(timeFormats) < 1 || timeFormats[0] == "" {
			return nil, errs.ErrImportFileTransactionTimeFormatInvalid
		}

		timezoneFormats := form.Value["timezoneFormat"]
		timezoneFormat := ""

		if len(timezoneFormats) > 0 {
			timezoneFormat = timezoneFormats[0]
		}

		amountDecimalSeparators := form.Value["amountDecimalSeparator"]
		amountDecimalSeparator := ""

		if len(amountDecimalSeparators) > 0 {
			amountDecimalSeparator = amountDecimalSeparators[0]
		}

		amountDigitGroupingSymbols := form.Value["amountDigitGroupingSymbol"]
		amountDigitGroupingSymbol := ""

		if len(amountDigitGroupingSymbols) > 0 {
			amountDigitGroupingSymbol = amountDigitGroupingSymbols[0]
		}

		geoLocationSeparators := form.Value["geoSeparator"]
		geoLocationSeparator := ""

		if len(geoLocationSeparators) > 0 {
			geoLocationSeparator = geoLocationSeparators[0]
		}

		transactionTagSeparators := form.Value["tagSeparator"]
		transactionTagSeparator := ""

		if len(transactionTagSeparators) > 0 {
			transactionTagSeparator = transactionTagSeparators[0]
		}

		dataImporter, err = converters.CreateNewDelimiterSeparatedValuesDataImporter(fileType, fileEncoding, columnIndexMapping, transactionTypeNameMapping, hasHeaderLine, timeFormats[0], timezoneFormat, amountDecimalSeparator, amountDigitGroupingSymbol, geoLocationSeparator, transactionTagSeparator)
	} else {
		dataImporter, err = converters.GetTransactionDataImporter(fileType)
	}

	if err != nil {
		return nil, errs.Or(err, errs.ErrImportFileTypeNotSupported)
	}

	importFiles := form.File["file"]

	if len(importFiles) < 1 {
		log.Warnf(c, "[transactions.TransactionParseImportFileHandler] there is no import file in request for user \"uid:%d\"", uid)
		return nil, errs.ErrNoFilesUpload
	}

	if importFiles[0].Size < 1 {
		log.Warnf(c, "[transactions.TransactionParseImportFileHandler] the size of import file in request is zero for user \"uid:%d\"", uid)
		return nil, errs.ErrUploadedFileEmpty
	}

	if importFiles[0].Size > int64(a.CurrentConfig().MaxImportFileSize) {
		log.Warnf(c, "[transactions.TransactionParseImportFileHandler] the upload file size \"%d\" exceeds the maximum size \"%d\" of import file for user \"uid:%d\"", importFiles[0].Size, a.CurrentConfig().MaxImportFileSize, uid)
		return nil, errs.ErrExceedMaxUploadFileSize
	}

	importFile, err := importFiles[0].Open()

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to get import file from request for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	defer importFile.Close()
	fileData, err := io.ReadAll(importFile)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to read import file data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_IMPORT_TRANSACTION) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	accounts, err := a.accounts.GetAllAccountsByUid(c, user.Uid)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to get accounts for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	accountMap := a.accounts.GetVisibleAccountNameMapByList(accounts)

	categories, err := a.transactionCategories.GetAllCategoriesByUid(c, user.Uid, 0, -1)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to get categories for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	expenseCategoryMap, incomeCategoryMap, transferCategoryMap := a.transactionCategories.GetVisibleSubCategoryNameMapByList(categories)

	tags, err := a.transactionTags.GetAllTagsByUid(c, user.Uid)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to get tags for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tagMap := a.transactionTags.GetTagNameMapByList(tags)

	parsedTransactions, _, _, _, _, _, err := dataImporter.ParseImportedData(c, user, fileData, utcOffset, accountMap, expenseCategoryMap, incomeCategoryMap, transferCategoryMap, tagMap)

	if err != nil {
		log.Errorf(c, "[transactions.TransactionParseImportFileHandler] failed to parse imported data for user \"uid:%d\", because %s", user.Uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	parsedTransactionRespsList := parsedTransactions.ToImportTransactionResponseList()

	if len(parsedTransactionRespsList) < 1 {
		return nil, errs.ErrNoDataToImport
	}

	parsedTransactionResps := &models.ImportTransactionResponsePageWrapper{
		Items:      parsedTransactionRespsList,
		TotalCount: int64(len(parsedTransactionRespsList)),
	}

	return parsedTransactionResps, nil
}

// TransactionImportHandler imports transactions by request parameters for current user
func (a *TransactionsApi) TransactionImportHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionImportReq models.TransactionImportRequest
	err := c.ShouldBindJSON(&transactionImportReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionImportHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && transactionImportReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_IMPORT_TRANSACTIONS, uid, transactionImportReq.ClientSessionId)

		if found {
			items := strings.Split(remark, ":")

			if len(items) >= 2 {
				if items[0] == "finished" {
					log.Infof(c, "[transactions.TransactionImportHandler] another \"%s\" transactions has been imported for user \"uid:%d\"", items[1], uid)
					count, err := utils.StringToInt(items[1])

					if err == nil {
						return count, nil
					}
				} else if items[0] == "processing" {
					return nil, errs.ErrRepeatedRequest
				}
			} else {
				log.Warnf(c, "[transactions.TransactionImportHandler] another transaction import task may be executing, but remark \"%s\" is invalid", remark)
			}
		}
	}

	newTransactionTagIdsMap := make(map[int][]int64, len(transactionImportReq.Transactions))

	for i := 0; i < len(transactionImportReq.Transactions); i++ {
		transactionCreateReq := transactionImportReq.Transactions[i]
		tagIds, err := utils.StringArrayToInt64Array(transactionCreateReq.TagIds)

		if err != nil {
			log.Warnf(c, "[transactions.TransactionImportHandler] parse tag ids failed of transaction \"index:%d\", because %s", i, err.Error())
			return nil, errs.ErrTransactionTagIdInvalid
		}

		if len(tagIds) > maximumTagsCountOfTransaction {
			return nil, errs.ErrTransactionHasTooManyTags
		}

		if transactionCreateReq.Type < models.TRANSACTION_TYPE_MODIFY_BALANCE || transactionCreateReq.Type > models.TRANSACTION_TYPE_TRANSFER {
			log.Warnf(c, "[transactions.TransactionImportHandler] transaction type of transaction \"index:%d\" is invalid", i)
			return nil, errs.ErrTransactionTypeInvalid
		}

		if transactionCreateReq.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE && transactionCreateReq.CategoryId > 0 {
			log.Warnf(c, "[transactions.TransactionImportHandler] balance modification transaction \"index:%d\" cannot set category id", i)
			return nil, errs.ErrBalanceModificationTransactionCannotSetCategory
		}

		if transactionCreateReq.Type != models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.DestinationAccountId != 0 {
			log.Warnf(c, "[transactions.TransactionImportHandler] non-transfer transaction \"index:%d\" destination account cannot be set", i)
			return nil, errs.ErrTransactionDestinationAccountCannotBeSet
		} else if transactionCreateReq.Type == models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.SourceAccountId == transactionCreateReq.DestinationAccountId {
			log.Warnf(c, "[transactions.TransactionImportHandler] transfer transaction \"index:%d\" source account must not be destination account", i)
			return nil, errs.ErrTransactionSourceAndDestinationIdCannotBeEqual
		}

		if transactionCreateReq.Type != models.TRANSACTION_TYPE_TRANSFER && transactionCreateReq.DestinationAmount != 0 {
			log.Warnf(c, "[transactions.TransactionImportHandler] non-transfer transaction \"index:%d\" destination amount cannot be set", i)
			return nil, errs.ErrTransactionDestinationAmountCannotBeSet
		}

		newTransactionTagIdsMap[i] = tagIds
	}

	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		if !errs.IsCustomError(err) {
			log.Errorf(c, "[transactions.TransactionImportHandler] failed to get user, because %s", err.Error())
		}

		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_IMPORT_TRANSACTION) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	newTransactions := make([]*models.Transaction, len(transactionImportReq.Transactions))

	for i := 0; i < len(transactionImportReq.Transactions); i++ {
		transactionCreateReq := transactionImportReq.Transactions[i]
		transaction := a.createNewTransactionModel(uid, transactionCreateReq, c.ClientIP())
		transactionEditable := user.CanEditTransactionByTransactionTime(transaction.TransactionTime, transactionCreateReq.UtcOffset)

		if !transactionEditable {
			return nil, errs.ErrCannotCreateTransactionWithThisTransactionTime
		}

		newTransactions[i] = transaction
	}

	err = a.transactions.BatchCreateTransactions(c, user.Uid, newTransactions, newTransactionTagIdsMap, func(currentProcess float64) {
		a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_IMPORT_TRANSACTIONS, uid, transactionImportReq.ClientSessionId, fmt.Sprintf("processing:%.2f", currentProcess))
	})
	count := len(newTransactions)

	if err != nil {
		a.RemoveSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_IMPORT_TRANSACTIONS, uid, transactionImportReq.ClientSessionId)
		log.Errorf(c, "[transactions.TransactionImportHandler] failed to import %d transactions for user \"uid:%d\", because %s", count, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transactions.TransactionImportHandler] user \"uid:%d\" has imported %d transactions successfully", uid, count)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_IMPORT_TRANSACTIONS, uid, transactionImportReq.ClientSessionId, fmt.Sprintf("finished:%d", count))

	return count, nil
}

// TransactionImportProcessHandler returns the process of specified transaction import task by request parameters for current user
func (a *TransactionsApi) TransactionImportProcessHandler(c *core.WebContext) (any, *errs.Error) {
	var transactionImportProcessReq models.TransactionImportProcessRequest
	err := c.ShouldBindQuery(&transactionImportProcessReq)

	if err != nil {
		log.Warnf(c, "[transactions.TransactionImportProcessHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if !a.CurrentConfig().EnableDuplicateSubmissionsCheck {
		return nil, nil
	}

	found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_IMPORT_TRANSACTIONS, uid, transactionImportProcessReq.ClientSessionId)

	if !found {
		return nil, nil
	}

	items := strings.Split(remark, ":")

	if len(items) < 2 {
		return nil, nil
	}

	if items[0] == "finished" {
		return 100, nil
	} else if items[0] != "processing" {
		return nil, nil
	}

	process, err := utils.StringToFloat64(items[1])

	if err != nil {
		log.Warnf(c, "[transactions.TransactionImportProcessHandler] parse process failed, because %s", err.Error())
		return nil, nil
	}

	if process < 0 {
		return nil, nil
	} else if process >= 100 {
		process = 100
	}

	return process, nil
}

func (a *TransactionsApi) filterTransactions(c *core.WebContext, uid int64, transactions []*models.Transaction, accountMap map[int64]*models.Account) []*models.Transaction {
	finalTransactions := make([]*models.Transaction, 0, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if _, exists := accountMap[transaction.AccountId]; !exists {
			log.Warnf(c, "[transactions.filterTransactions] account of transaction \"id:%d\" does not exist for user \"uid:%d\"", transaction.TransactionId, uid)
			continue
		}

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN || transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			if _, exists := accountMap[transaction.RelatedAccountId]; !exists {
				log.Warnf(c, "[transactions.filterTransactions] related account of transaction \"id:%d\" does not exist for user \"uid:%d\"", transaction.TransactionId, uid)
				continue
			}
		}

		finalTransactions = append(finalTransactions, transaction)
	}

	return finalTransactions
}

func (a *TransactionsApi) getAccountOrSubAccountIds(c *core.WebContext, accountIds string, uid int64) ([]int64, error) {
	if accountIds == "" || accountIds == "0" {
		return nil, nil
	}

	requestAccountIds, err := utils.StringArrayToInt64Array(strings.Split(accountIds, ","))

	if err != nil {
		return nil, errs.Or(err, errs.ErrAccountIdInvalid)
	}

	var allAccountIds []int64

	if len(requestAccountIds) > 0 {
		allSubAccounts, err := a.accounts.GetSubAccountsByAccountIds(c, uid, requestAccountIds)

		if err != nil {
			return nil, err
		}

		accountIdsMap := make(map[int64]int32, len(requestAccountIds))

		for i := 0; i < len(requestAccountIds); i++ {
			accountIdsMap[requestAccountIds[i]] = 0
		}

		for i := 0; i < len(allSubAccounts); i++ {
			subAccount := allSubAccounts[i]

			if refCount, exists := accountIdsMap[subAccount.ParentAccountId]; exists {
				accountIdsMap[subAccount.ParentAccountId] = refCount + 1
			} else {
				accountIdsMap[subAccount.ParentAccountId] = 1
			}

			if _, exists := accountIdsMap[subAccount.AccountId]; exists {
				delete(accountIdsMap, subAccount.AccountId)
			}

			allAccountIds = append(allAccountIds, subAccount.AccountId)
		}

		for accountId, refCount := range accountIdsMap {
			if refCount < 1 {
				allAccountIds = append(allAccountIds, accountId)
			}
		}
	}

	return allAccountIds, nil
}

func (a *TransactionsApi) getCategoryOrSubCategoryIds(c *core.WebContext, categoryIds string, uid int64) ([]int64, error) {
	if categoryIds == "" || categoryIds == "0" {
		return nil, nil
	}

	requestCategoryIds, err := utils.StringArrayToInt64Array(strings.Split(categoryIds, ","))

	if err != nil {
		return nil, errs.Or(err, errs.ErrTransactionCategoryIdInvalid)
	}

	var allCategoryIds []int64

	if len(requestCategoryIds) > 0 {
		allSubCategories, err := a.transactionCategories.GetSubCategoriesByCategoryIds(c, uid, requestCategoryIds)

		if err != nil {
			return nil, err
		}

		categoryIdsMap := make(map[int64]int32, len(requestCategoryIds))

		for i := 0; i < len(requestCategoryIds); i++ {
			categoryIdsMap[requestCategoryIds[i]] = 0
		}

		for i := 0; i < len(allSubCategories); i++ {
			subCategory := allSubCategories[i]

			if refCount, exists := categoryIdsMap[subCategory.ParentCategoryId]; exists {
				categoryIdsMap[subCategory.ParentCategoryId] = refCount + 1
			} else {
				categoryIdsMap[subCategory.ParentCategoryId] = 1
			}

			if _, exists := categoryIdsMap[subCategory.CategoryId]; exists {
				delete(categoryIdsMap, subCategory.CategoryId)
			}

			allCategoryIds = append(allCategoryIds, subCategory.CategoryId)
		}

		for accountId, refCount := range categoryIdsMap {
			if refCount < 1 {
				allCategoryIds = append(allCategoryIds, accountId)
			}
		}
	}

	return allCategoryIds, nil
}

func (a *TransactionsApi) getTagIds(tagIds string) ([]int64, error) {
	if tagIds == "" || tagIds == "0" {
		return nil, nil
	}

	requestTagIds, err := utils.StringArrayToInt64Array(strings.Split(tagIds, ","))

	if err != nil {
		return nil, errs.Or(err, errs.ErrTransactionTagIdInvalid)
	}

	return requestTagIds, nil
}

func (a *TransactionsApi) getTransactionTagIds(allTransactionTagIds map[int64][]int64) []int64 {
	allTagIds := make([]int64, 0, len(allTransactionTagIds))

	for _, tagIds := range allTransactionTagIds {
		allTagIds = append(allTagIds, tagIds...)
	}

	return allTagIds
}

func (a *TransactionsApi) getTransactionTagInfoResponses(tagIds []int64, allTransactionTags map[int64]*models.TransactionTag) []*models.TransactionTagInfoResponse {
	allTags := make([]*models.TransactionTagInfoResponse, 0, len(tagIds))

	for i := 0; i < len(tagIds); i++ {
		tag := allTransactionTags[tagIds[i]]

		if tag == nil {
			continue
		}

		allTags = append(allTags, tag.ToTransactionTagInfoResponse())
	}

	return allTags
}

func (a *TransactionsApi) getTransactionResponseListResult(c *core.WebContext, user *models.User, transactions []*models.Transaction, utcOffset int16, withPictures bool, trimAccount bool, trimCategory bool, trimTag bool) (models.TransactionInfoResponseSlice, error) {
	uid := user.Uid
	transactionIds := make([]int64, len(transactions))
	accountIds := make([]int64, 0, len(transactions)*2)
	categoryIds := make([]int64, 0, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transactionId := transactions[i].TransactionId

		if transactions[i].Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			transactionId = transactions[i].RelatedId
		}

		transactionIds[i] = transactionId
		accountIds = append(accountIds, transactions[i].AccountId)

		if transactions[i].Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN || transactions[i].Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT {
			accountIds = append(accountIds, transactions[i].RelatedAccountId)
		}

		categoryIds = append(categoryIds, transactions[i].CategoryId)
	}

	allAccounts, err := a.accounts.GetAccountsByAccountIds(c, uid, utils.ToUniqueInt64Slice(accountIds))

	if err != nil {
		log.Errorf(c, "[transactions.getTransactionResponseListResult] failed to get accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, err
	}

	transactions = a.filterTransactions(c, uid, transactions, allAccounts)

	allTransactionTagIds, err := a.transactionTags.GetAllTagIdsOfTransactions(c, uid, transactionIds)

	if err != nil {
		log.Errorf(c, "[transactions.getTransactionResponseListResult] failed to get transactions tag ids for user \"uid:%d\", because %s", uid, err.Error())
		return nil, err
	}

	var categoryMap map[int64]*models.TransactionCategory
	var tagMap map[int64]*models.TransactionTag
	var pictureInfoMap map[int64][]*models.TransactionPictureInfo

	if !trimCategory {
		categoryMap, err = a.transactionCategories.GetCategoriesByCategoryIds(c, uid, utils.ToUniqueInt64Slice(categoryIds))

		if err != nil {
			log.Errorf(c, "[transactions.getTransactionResponseListResult] failed to get transactions categories for user \"uid:%d\", because %s", uid, err.Error())
			return nil, err
		}
	}

	if !trimTag {
		tagMap, err = a.transactionTags.GetTagsByTagIds(c, uid, utils.ToUniqueInt64Slice(a.getTransactionTagIds(allTransactionTagIds)))

		if err != nil {
			log.Errorf(c, "[transactions.getTransactionResponseListResult] failed to get transactions tags for user \"uid:%d\", because %s", uid, err.Error())
			return nil, err
		}
	}

	if withPictures && a.CurrentConfig().EnableTransactionPictures {
		pictureInfoMap, err = a.transactionPictures.GetPictureInfosByTransactionIds(c, uid, utils.ToUniqueInt64Slice(a.transactions.GetTransactionIds(transactions)))

		if err != nil {
			log.Errorf(c, "[transactions.getTransactionResponseListResult] failed to get transactions pictures for user \"uid:%d\", because %s", uid, err.Error())
			return nil, err
		}
	}

	result := make(models.TransactionInfoResponseSlice, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if transaction.Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
			transaction = a.transactions.GetRelatedTransferTransaction(transaction)
		}

		transactionEditable := transaction.IsEditable(user, utcOffset, allAccounts[transaction.AccountId], allAccounts[transaction.RelatedAccountId])
		transactionTagIds := allTransactionTagIds[transaction.TransactionId]
		result[i] = transaction.ToTransactionInfoResponse(transactionTagIds, transactionEditable)

		if !trimAccount {
			if sourceAccount := allAccounts[transaction.AccountId]; sourceAccount != nil {
				result[i].SourceAccount = sourceAccount.ToAccountInfoResponse()
			}

			if destinationAccount := allAccounts[transaction.RelatedAccountId]; destinationAccount != nil {
				result[i].DestinationAccount = destinationAccount.ToAccountInfoResponse()
			}
		}

		if !trimCategory {
			if category := categoryMap[transaction.CategoryId]; category != nil {
				result[i].Category = category.ToTransactionCategoryInfoResponse()
			}
		}

		if !trimTag {
			result[i].Tags = a.getTransactionTagInfoResponses(transactionTagIds, tagMap)
		}

		if withPictures && a.CurrentConfig().EnableTransactionPictures {
			pictureInfos, exists := pictureInfoMap[transaction.TransactionId]

			if exists {
				result[i].Pictures = a.GetTransactionPictureInfoResponseList(pictureInfos)
			}
		}
	}

	sort.Sort(result)

	return result, nil
}

func (a *TransactionsApi) createNewTransactionModel(uid int64, transactionCreateReq *models.TransactionCreateRequest, clientIp string) *models.Transaction {
	var transactionDbType models.TransactionDbType

	if transactionCreateReq.Type == models.TRANSACTION_TYPE_MODIFY_BALANCE {
		transactionDbType = models.TRANSACTION_DB_TYPE_MODIFY_BALANCE
	} else if transactionCreateReq.Type == models.TRANSACTION_TYPE_EXPENSE {
		transactionDbType = models.TRANSACTION_DB_TYPE_EXPENSE
	} else if transactionCreateReq.Type == models.TRANSACTION_TYPE_INCOME {
		transactionDbType = models.TRANSACTION_DB_TYPE_INCOME
	} else if transactionCreateReq.Type == models.TRANSACTION_TYPE_TRANSFER {
		transactionDbType = models.TRANSACTION_DB_TYPE_TRANSFER_OUT
	}

	transaction := &models.Transaction{
		Uid:               uid,
		Type:              transactionDbType,
		CategoryId:        transactionCreateReq.CategoryId,
		TransactionTime:   utils.GetMinTransactionTimeFromUnixTime(transactionCreateReq.Time),
		TimezoneUtcOffset: transactionCreateReq.UtcOffset,
		AccountId:         transactionCreateReq.SourceAccountId,
		Amount:            transactionCreateReq.SourceAmount,
		HideAmount:        transactionCreateReq.HideAmount,
		Comment:           transactionCreateReq.Comment,
		CreatedIp:         clientIp,
	}

	if transactionCreateReq.Type == models.TRANSACTION_TYPE_TRANSFER {
		transaction.RelatedAccountId = transactionCreateReq.DestinationAccountId
		transaction.RelatedAccountAmount = transactionCreateReq.DestinationAmount
	}

	if transactionCreateReq.GeoLocation != nil {
		transaction.GeoLongitude = transactionCreateReq.GeoLocation.Longitude
		transaction.GeoLatitude = transactionCreateReq.GeoLocation.Latitude
	}

	return transaction
}
