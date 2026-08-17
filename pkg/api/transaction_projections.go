package api

import (
	"sort"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// TransactionProjectionsHandler returns the monthly projection of transaction amounts of the current
// user, combining the transactions already recorded up to now with the ones the scheduled
// transaction templates of the user will produce afterwards.
//
// The response has the same shape as the one of TransactionStatisticsTrendsHandler, so the client can
// read both with the same code. Every item carries its account, which the client needs to convert
// amounts of accounts held in different currencies before adding them up.
func (a *TransactionsApi) TransactionProjectionsHandler(c *core.WebContext) (any, *errs.Error) {
	var projectionReq models.TransactionProjectionRequest
	err := c.ShouldBindQuery(&projectionReq)

	if err != nil {
		log.Warnf(c, "[transaction_projections.TransactionProjectionsHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[transaction_projections.TransactionProjectionsHandler] cannot get client timezone, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	startYear, startMonth, endYear, endMonth, err := projectionReq.GetNumericYearMonthRange()

	if err != nil {
		log.Warnf(c, "[transaction_projections.TransactionProjectionsHandler] cannot parse year month, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	uid := c.GetCurrentUid()
	allMonthlyTotalAmounts, err := a.transactions.GetProjectedCategoryAmountsByMonth(c, uid, startYear, startMonth, endYear, endMonth, clientTimezone, projectionReq.UseTransactionTimezone, time.Now().Unix())

	if err != nil {
		log.Errorf(c, "[transaction_projections.TransactionProjectionsHandler] failed to get projected category amounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	projectionResp := make(models.TransactionStatisticTrendsResponseItemSlice, 0, len(allMonthlyTotalAmounts))

	for yearMonth, monthlyTotalAmounts := range allMonthlyTotalAmounts {
		monthlyProjectionResp := &models.TransactionStatisticTrendsResponseItem{
			Year:  yearMonth / 100,
			Month: yearMonth % 100,
			Items: make([]*models.TransactionStatisticResponseItem, len(monthlyTotalAmounts)),
		}

		for i := 0; i < len(monthlyTotalAmounts); i++ {
			totalAmountItem := monthlyTotalAmounts[i]
			monthlyProjectionResp.Items[i] = &models.TransactionStatisticResponseItem{
				CategoryId:  totalAmountItem.CategoryId,
				AccountId:   totalAmountItem.AccountId,
				TotalAmount: totalAmountItem.Amount.String(),
			}
		}

		projectionResp = append(projectionResp, monthlyProjectionResp)
	}

	sort.Sort(projectionResp)

	return projectionResp, nil
}
