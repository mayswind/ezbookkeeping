package api

import (
	"strings"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
	"github.com/mayswind/lab/pkg/utils"
)

// OverviewApi represents overview api
type OverviewApi struct {
	transactions *services.TransactionService
}

// Initialize an overview api singleton instance
var (
	Overviews = &OverviewApi{
		transactions: services.Transactions,
	}
)

// TransactionOverviewHandler returns transaction over of current user
func (a *OverviewApi) TransactionOverviewHandler(c *core.Context) (interface{}, *errs.Error) {
	var transactionOverviewReq models.TransactionOverviewRequest
	err := c.ShouldBindQuery(&transactionOverviewReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[overviews.TransactionOverviewHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	items := strings.Split(transactionOverviewReq.Query, "|")
	requestItems := make([]*models.TransactionOverviewRequestItem, 0, len(items))

	for i := 0; i < len(items); i++ {
		itemValues := strings.Split(items[i], "_")

		if len(itemValues) != 3 {
			log.WarnfWithRequestId(c, "[overviews.TransactionOverviewHandler] parse request item failed, because its not valid item, content is \"%s\"", items[i])
			continue
		}

		startTime, err := utils.StringToInt64(itemValues[1])

		if err != nil {
			log.WarnfWithRequestId(c, "[overviews.TransactionOverviewHandler] parse request item start time failed, because %s", err.Error())
			continue
		}

		endTime, err := utils.StringToInt64(itemValues[2])

		if err != nil {
			log.WarnfWithRequestId(c, "[overviews.TransactionOverviewHandler] parse request item end time failed, because %s", err.Error())
			continue
		}

		requestItem := &models.TransactionOverviewRequestItem{
			Name:      itemValues[0],
			StartTime: startTime,
			EndTime:   endTime,
		}

		requestItems = append(requestItems, requestItem)
	}

	if len(requestItems) < 1 {
		log.WarnfWithRequestId(c, "[overviews.TransactionOverviewHandler] parse request failed, because there are no valid items")
		return nil, errs.ErrQueryItemsEmpty
	}

	if len(requestItems) > 5 {
		log.WarnfWithRequestId(c, "[overviews.TransactionOverviewHandler] parse request failed, because there are too many items")
		return nil, errs.ErrQueryItemsTooMuch
	}

	uid := c.GetCurrentUid()

	overviewResp := make(map[string]*models.TransactionOverviewResponseItem)

	for i := 0; i < len(requestItems); i++ {
		requestItem := requestItems[i]

		incomeAmount, expenseAmount, err := a.transactions.GetTotalIncomeAndExpenseByDateRange(uid, requestItem.StartTime, requestItem.EndTime)

		if err != nil {
			log.ErrorfWithRequestId(c, "[overviews.TransactionOverviewHandler] failed to get transaction overview item for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrOperationFailed
		}

		overviewResp[requestItem.Name] = &models.TransactionOverviewResponseItem{
			StartTime:     requestItem.StartTime,
			EndTime:       requestItem.EndTime,
			IncomeAmount:  incomeAmount,
			ExpenseAmount: expenseAmount,
		}
	}

	return overviewResp, nil
}
