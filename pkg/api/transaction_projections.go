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
// transaction templates of the user will produce afterwards, plus the balance every account holds
// just before the period starts.
//
// The monthly amounts have the same shape as the ones of TransactionStatisticsTrendsHandler, so the
// client can read both with the same code. Every item carries its account, which the client needs to
// convert amounts of accounts held in different currencies before adding them up, and so does every
// opening balance.
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

	openingBalances, err := a.transactions.GetProjectionOpeningBalances(c, uid, startYear, startMonth, time.Now().Unix())

	if err != nil {
		log.Errorf(c, "[transaction_projections.TransactionProjectionsHandler] failed to get projection opening balances for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	monthsResp := make(models.TransactionStatisticTrendsResponseItemSlice, 0, len(allMonthlyTotalAmounts))

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

		monthsResp = append(monthsResp, monthlyProjectionResp)
	}

	sort.Sort(monthsResp)

	// An account whose opening balance is zero is left out: it carries nothing into the period, and
	// the client adds up whatever it receives
	openingBalanceResp := make(models.TransactionProjectionOpeningBalanceItemSlice, 0, len(openingBalances))

	for accountId, openingBalance := range openingBalances {
		if openingBalance.Sign() == 0 {
			continue
		}

		openingBalanceResp = append(openingBalanceResp, &models.TransactionProjectionOpeningBalanceItem{
			AccountId: accountId,
			Amount:    openingBalance.String(),
		})
	}

	sort.Sort(openingBalanceResp)

	return &models.TransactionProjectionResponse{
		OpeningBalance: openingBalanceResp,
		Months:         monthsResp,
	}, nil
}
