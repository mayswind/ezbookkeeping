package api

import (
	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
)

// StatisticApi represents statistic api
type StatisticApi struct {
	transactions *services.TransactionService
}

// Initialize an statistic api singleton instance
var (
	Statistics = &StatisticApi{
		transactions: services.Transactions,
	}
)

// TransactionStatisticsHandler returns transaction statistics of current user
func (a *StatisticApi) TransactionStatisticsHandler(c *core.Context) (interface{}, *errs.Error) {
	var statisticReq models.TransactionStatisticRequest
	err := c.ShouldBindQuery(&statisticReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[statistics.TransactionOverviewHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	totalAmounts, err := a.transactions.GetAccountsAndCategoriesTotalIncomeAndExpense(uid, statisticReq.StartTime, statisticReq.EndTime)

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
