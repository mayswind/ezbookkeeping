package services

import (
	"math/big"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// maxTransactionProjectionMonths is the longest period a projection can be requested for. The real
// half of a projection loads every transaction of the period into memory, so the range cannot be
// left open ended.
const maxTransactionProjectionMonths = 60

// projectionOccurrenceMarginSeconds widens the simulated window by one day at each end. Occurrences
// are bucketed into months using the timezone of the client or of the template, which can be up to
// 14 hours away from UTC, so an occurrence just outside the UTC bounds of the period may still
// belong to one of the requested months.
const projectionOccurrenceMarginSeconds int64 = 24 * 60 * 60

// GetProjectedCategoryAmountsByMonth returns the monthly amount of every category of the specified
// user within the specified period, combining what has already happened with what the scheduled
// transaction templates of the user will produce.
//
// Everything up to currentUnixTime comes from transactions already recorded in the database, and
// everything after it is simulated from the scheduled templates using their current amount. The two
// halves are split exactly at currentUnixTime, so they never overlap and the month containing that
// instant is simply the sum of both.
//
// Transfers are excluded from both halves: they move money between accounts of the user instead of
// adding to income or expenses, and a projection that counted them on the real half only would show
// them appearing in past months and vanishing from future ones.
//
// The result is keyed by numeric year month and shaped exactly like
// GetAccountsAndCategoriesMonthlyInflowAndOutflow, including the account of every amount, which the
// caller needs to convert amounts of accounts held in different currencies.
func (s *TransactionService) GetProjectedCategoryAmountsByMonth(c core.Context, uid int64, startYear int32, startMonth int32, endYear int32, endMonth int32, clientTimezone *time.Location, useTransactionTimezone bool, currentUnixTime int64) (map[int32][]*models.TransactionTotalAmount, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	// A cutoff of zero would disable it in GetAccountsAndCategoriesMonthlyInflowAndOutflow and make
	// both halves overlap silently, so it is rejected instead of defaulted
	if currentUnixTime <= 0 {
		return nil, errs.ErrSystemError
	}

	err := validateTransactionProjectionRange(startYear, startMonth, endYear, endMonth)

	if err != nil {
		return nil, err
	}

	actualAmounts, err := s.GetAccountsAndCategoriesMonthlyInflowAndOutflow(c, uid, startYear, startMonth, endYear, endMonth, nil, false, "", core.MATCH_MODE_DEFAULT, clientTimezone, useTransactionTimezone, currentUnixTime)

	if err != nil {
		return nil, err
	}

	projectedAmounts, err := s.getScheduledCategoryAmountsByMonth(c, uid, startYear, startMonth, endYear, endMonth, clientTimezone, useTransactionTimezone, currentUnixTime)

	if err != nil {
		return nil, err
	}

	return mergeMonthlyTotalAmounts(removeTransferTotalAmounts(actualAmounts), projectedAmounts), nil
}

// getScheduledCategoryAmountsByMonth simulates every occurrence the scheduled transaction templates
// of the specified user would produce after currentUnixTime, and returns their monthly amount per
// category and account.
func (s *TransactionService) getScheduledCategoryAmountsByMonth(c core.Context, uid int64, startYear int32, startMonth int32, endYear int32, endMonth int32, clientTimezone *time.Location, useTransactionTimezone bool, currentUnixTime int64) (map[int32][]*models.TransactionTotalAmount, error) {
	templates, err := s.getActiveScheduledTemplates(c, uid)

	if err != nil {
		return nil, err
	}

	return projectScheduledCategoryAmountsByMonth(c, templates, startYear, startMonth, endYear, endMonth, clientTimezone, useTransactionTimezone, currentUnixTime)
}

// projectScheduledCategoryAmountsByMonth is the simulated half of a projection: it accumulates the
// amount of every occurrence the specified templates would produce after currentUnixTime, per month,
// category and account.
//
// Templates of type transfer are skipped, and a template whose scheduled frequency is malformed is
// skipped as well rather than failing the whole projection, in the same way the cron logs it and
// carries on.
func projectScheduledCategoryAmountsByMonth(c core.Context, templates []*models.TransactionTemplate, startYear int32, startMonth int32, endYear int32, endMonth int32, clientTimezone *time.Location, useTransactionTimezone bool, currentUnixTime int64) (map[int32][]*models.TransactionTotalAmount, error) {
	monthlyAmounts := make(map[int32][]*models.TransactionTotalAmount)

	if len(templates) < 1 {
		return monthlyAmounts, nil
	}

	periodStartUnixTime, periodEndUnixTime, err := getTransactionProjectionPeriodUnixTimeRange(startYear, startMonth, endYear, endMonth)

	if err != nil {
		return nil, err
	}

	// The simulated half starts at the cutoff, so a period entirely in the past yields nothing
	from := time.Unix(max(currentUnixTime, periodStartUnixTime-projectionOccurrenceMarginSeconds), 0)
	to := time.Unix(periodEndUnixTime+projectionOccurrenceMarginSeconds, 0)

	startYearMonth := startYear*100 + startMonth
	endYearMonth := endYear*100 + endMonth
	amountsMap := make(map[monthlyTotalAmountKey]*models.TransactionTotalAmount)

	for i := 0; i < len(templates); i++ {
		template := templates[i]

		if template.Type == models.TRANSACTION_TYPE_TRANSFER {
			continue
		}

		transactionDbType, err := getTransactionDbTypeByTemplateType(template.Type)

		if err != nil {
			log.Warnf(c, "[transactions.projectScheduledCategoryAmountsByMonth] transaction template \"id:%d\" has invalid transaction type", template.TemplateId)
			continue
		}

		occurrences, err := GetScheduledOccurrences(template, from, to)

		if err != nil {
			log.Warnf(c, "[transactions.projectScheduledCategoryAmountsByMonth] transaction template \"id:%d\" has invalid scheduled transaction frequency, because %s", template.TemplateId, err.Error())
			continue
		}

		timeZone := clientTimezone

		if useTransactionTimezone {
			timeZone = getScheduledTemplateTimeZone(template)
		}

		for j := 0; j < len(occurrences); j++ {
			yearMonth := utils.FormatUnixTimeToNumericYearMonth(occurrences[j].Unix(), timeZone)

			if (startYearMonth > 0 && yearMonth < startYearMonth) || (endYearMonth > 0 && yearMonth > endYearMonth) {
				continue
			}

			addToMonthlyTotalAmounts(amountsMap, yearMonth, &models.TransactionTotalAmount{
				Type:       transactionDbType,
				CategoryId: template.CategoryId,
				AccountId:  template.AccountId,
				Amount:     big.NewInt(template.Amount),
			})
		}
	}

	return groupTotalAmountsByYearMonth(amountsMap), nil
}

// GetProjectionOpeningBalances returns the balance every account of the specified user holds just
// before the specified period starts, keyed by account id.
//
// The projection table needs it to turn its accumulated row from a running sum of the flows of the
// period into the net worth reached at the end of every month. The instant it stops at is the one
// GetAccountsAndCategoriesMonthlyInflowAndOutflow starts the real half at, so between the opening
// balance and the first month of the projection no transaction is counted twice nor left out.
//
// It is derived from the current balance of every account instead of being replayed from the first
// transaction of the user: walking back from a known balance only reads the transactions of the
// period, while replaying forward would read the whole history of the user. Every transaction type
// takes part, balance modifications and both legs of every transfer included, because the result has
// to be the balance of an account rather than the flow of its categories -- which is also why the
// amounts are keyed by account and not by category.
//
// When the period starts in the future, the occurrences the scheduled templates would produce
// between now and that instant are added as well, so the period opens at the balance the user is
// expected to hold by then and the accumulated row stays continuous whichever period is asked for.
func (s *TransactionService) GetProjectionOpeningBalances(c core.Context, uid int64, startYear int32, startMonth int32, currentUnixTime int64) (map[int64]*big.Int, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if currentUnixTime <= 0 {
		return nil, errs.ErrSystemError
	}

	err := validateTransactionProjectionRange(startYear, startMonth, startYear, startMonth)

	if err != nil {
		return nil, err
	}

	periodStartTransactionTime, _, err := utils.GetTransactionTimeRangeByYearMonth(startYear, startMonth)

	if err != nil {
		return nil, errs.ErrSystemError
	}

	openingBalances, err := s.getAccountBalancesBeforeTransactionTime(c, uid, periodStartTransactionTime)

	if err != nil {
		return nil, err
	}

	periodStartUnixTime := utils.GetUnixTimeFromTransactionTime(periodStartTransactionTime)

	if periodStartUnixTime <= currentUnixTime {
		return openingBalances, nil
	}

	// The period opens after the cutoff, so the stretch between them is simulated exactly like the
	// future half of the projection is
	templates, err := s.getActiveScheduledTemplates(c, uid)

	if err != nil {
		return nil, err
	}

	projectedAmounts := projectScheduledAccountAmounts(c, templates, time.Unix(currentUnixTime, 0), time.Unix(periodStartUnixTime-1, 0))

	for accountId, projectedAmount := range projectedAmounts {
		openingBalance, exists := openingBalances[accountId]

		// getAccountBalancesBeforeTransactionTime returns an entry for every account of the user, so
		// a missing one is a template left pointing at a deleted account. A deleted account holds no
		// balance to open the period with, and the client could not convert its currency either.
		if !exists {
			continue
		}

		openingBalance.Add(openingBalance, projectedAmount)
	}

	return openingBalances, nil
}

// getAccountBalancesBeforeTransactionTime returns the balance every account of the specified user
// held before the specified transaction time, by taking its current balance and undoing every
// transaction recorded at or after that instant.
func (s *TransactionService) getAccountBalancesBeforeTransactionTime(c core.Context, uid int64, minTransactionTime int64) (map[int64]*big.Int, error) {
	var accounts []*models.Account
	err := s.UserDataDB(uid).NewSession(c).Select("account_id, balance").Where("uid=? AND deleted=?", uid, false).Find(&accounts)

	if err != nil {
		return nil, err
	}

	balances := make(map[int64]*big.Int, len(accounts))

	for i := 0; i < len(accounts); i++ {
		balances[accounts[i].AccountId] = big.NewInt(accounts[i].Balance)
	}

	// Paginated the same way the other amount aggregations of this service are, so a long period
	// cannot load an unbounded number of rows at once
	maxTransactionTime := int64(0)

	for {
		var transactions []*models.Transaction
		condition := "uid=? AND deleted=? AND transaction_time>=?"
		conditionParams := []any{uid, false, minTransactionTime}

		if maxTransactionTime > 0 {
			condition = condition + " AND transaction_time<=?"
			conditionParams = append(conditionParams, maxTransactionTime)
		}

		err = s.UserDataDB(uid).NewSession(c).Select("type, account_id, transaction_time, amount, related_account_amount").
			Where(condition, conditionParams...).
			Limit(pageCountForLoadTransactionAmounts, 0).
			OrderBy("transaction_time desc").
			Find(&transactions)

		if err != nil {
			return nil, err
		}

		for i := 0; i < len(transactions); i++ {
			transaction := transactions[i]
			balance, exists := balances[transaction.AccountId]

			if !exists {
				// a transaction of an account that no longer exists cannot be undone against any
				// balance, and its account is not shown by the table either
				continue
			}

			delta, err := getAccountBalanceDelta(transaction)

			if err != nil {
				log.Errorf(c, "[transactions.getAccountBalancesBeforeTransactionTime] transaction \"id:%d\" of user \"uid:%d\" has invalid transaction type %d", transaction.TransactionId, uid, transaction.Type)
				return nil, err
			}

			balance.Sub(balance, delta)
		}

		if len(transactions) < pageCountForLoadTransactionAmounts {
			break
		}

		maxTransactionTime = transactions[len(transactions)-1].TransactionTime - 1

		if maxTransactionTime < minTransactionTime {
			break
		}
	}

	return balances, nil
}

// getAccountBalanceDelta returns how much the specified transaction moved the balance of its own
// account, which is how CreateTransaction keeps the account table up to date.
func getAccountBalanceDelta(transaction *models.Transaction) (*big.Int, error) {
	switch transaction.Type {
	case models.TRANSACTION_DB_TYPE_MODIFY_BALANCE:
		return big.NewInt(transaction.RelatedAccountAmount), nil
	case models.TRANSACTION_DB_TYPE_INCOME, models.TRANSACTION_DB_TYPE_TRANSFER_IN:
		return big.NewInt(transaction.Amount), nil
	case models.TRANSACTION_DB_TYPE_EXPENSE, models.TRANSACTION_DB_TYPE_TRANSFER_OUT:
		return new(big.Int).Neg(big.NewInt(transaction.Amount)), nil
	default:
		return nil, errs.ErrTransactionTypeInvalid
	}
}

// projectScheduledAccountAmounts accumulates, per account, the net amount the specified templates
// would produce between from (exclusive) and to (inclusive).
//
// Unlike the monthly projection this one does not bucket anything by month, so it needs no timezone:
// an occurrence either falls inside the stretch or it does not. Templates of type transfer are
// skipped, as they are everywhere else in a projection, and a template whose frequency is malformed
// is skipped rather than failing the whole projection.
func projectScheduledAccountAmounts(c core.Context, templates []*models.TransactionTemplate, from time.Time, to time.Time) map[int64]*big.Int {
	amounts := make(map[int64]*big.Int)

	for i := 0; i < len(templates); i++ {
		template := templates[i]

		if template.Type != models.TRANSACTION_TYPE_INCOME && template.Type != models.TRANSACTION_TYPE_EXPENSE {
			continue
		}

		occurrences, err := GetScheduledOccurrences(template, from, to)

		if err != nil {
			log.Warnf(c, "[transactions.projectScheduledAccountAmounts] transaction template \"id:%d\" has invalid scheduled transaction frequency, because %s", template.TemplateId, err.Error())
			continue
		}

		if len(occurrences) < 1 {
			continue
		}

		amount, exists := amounts[template.AccountId]

		if !exists {
			amount = big.NewInt(0)
			amounts[template.AccountId] = amount
		}

		occurrencesAmount := new(big.Int).Mul(big.NewInt(int64(len(occurrences))), big.NewInt(template.Amount))

		if template.Type == models.TRANSACTION_TYPE_INCOME {
			amount.Add(amount, occurrencesAmount)
		} else {
			amount.Sub(amount, occurrencesAmount)
		}
	}

	return amounts
}

// getActiveScheduledTemplates returns the scheduled transaction templates of the specified user that
// the cron would consider for creating transactions.
//
// Hidden templates are included on purpose: CreateScheduledTransactions does not filter them out
// either, so hiding a template does not stop it from creating transactions and must not stop it from
// being projected.
func (s *TransactionService) getActiveScheduledTemplates(c core.Context, uid int64) ([]*models.TransactionTemplate, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var templates []*models.TransactionTemplate
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND template_type=? AND scheduled_frequency_type<>?",
		uid,
		false,
		models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED).Find(&templates)

	return templates, err
}

// validateTransactionProjectionRange checks that the specified period is complete, ordered and not
// longer than maxTransactionProjectionMonths
func validateTransactionProjectionRange(startYear int32, startMonth int32, endYear int32, endMonth int32) error {
	if startYear <= 0 || startMonth <= 0 || endYear <= 0 || endMonth <= 0 {
		return errs.ErrIncompleteOrIncorrectSubmission
	}

	if startMonth > 12 || endMonth > 12 {
		return errs.ErrIncompleteOrIncorrectSubmission
	}

	if startYear*100+startMonth > endYear*100+endMonth {
		return errs.ErrIncompleteOrIncorrectSubmission
	}

	monthCount := (endYear-startYear)*12 + (endMonth - startMonth) + 1

	if monthCount > maxTransactionProjectionMonths {
		return errs.ErrQueryItemsTooMuch
	}

	return nil
}

// getTransactionProjectionPeriodUnixTimeRange returns the unix time of the first and the last instant
// of the specified period
func getTransactionProjectionPeriodUnixTimeRange(startYear int32, startMonth int32, endYear int32, endMonth int32) (int64, int64, error) {
	startTransactionTime, _, err := utils.GetTransactionTimeRangeByYearMonth(startYear, startMonth)

	if err != nil {
		return 0, 0, errs.ErrSystemError
	}

	_, endTransactionTime, err := utils.GetTransactionTimeRangeByYearMonth(endYear, endMonth)

	if err != nil {
		return 0, 0, errs.ErrSystemError
	}

	return utils.GetUnixTimeFromTransactionTime(startTransactionTime), utils.GetUnixTimeFromTransactionTime(endTransactionTime), nil
}

// getTransactionDbTypeByTemplateType returns the database transaction type a template of the
// specified transaction type would create
func getTransactionDbTypeByTemplateType(transactionType models.TransactionType) (models.TransactionDbType, error) {
	switch transactionType {
	case models.TRANSACTION_TYPE_INCOME:
		return models.TRANSACTION_DB_TYPE_INCOME, nil
	case models.TRANSACTION_TYPE_EXPENSE:
		return models.TRANSACTION_DB_TYPE_EXPENSE, nil
	case models.TRANSACTION_TYPE_TRANSFER:
		return models.TRANSACTION_DB_TYPE_TRANSFER_OUT, nil
	default:
		return 0, errs.ErrTransactionTypeInvalid
	}
}

// removeTransferTotalAmounts returns the specified monthly amounts without the transfer ones
func removeTransferTotalAmounts(monthlyAmounts map[int32][]*models.TransactionTotalAmount) map[int32][]*models.TransactionTotalAmount {
	filteredMonthlyAmounts := make(map[int32][]*models.TransactionTotalAmount, len(monthlyAmounts))

	for yearMonth, amounts := range monthlyAmounts {
		filteredAmounts := make([]*models.TransactionTotalAmount, 0, len(amounts))

		for i := 0; i < len(amounts); i++ {
			if amounts[i].Type == models.TRANSACTION_DB_TYPE_TRANSFER_OUT || amounts[i].Type == models.TRANSACTION_DB_TYPE_TRANSFER_IN {
				continue
			}

			filteredAmounts = append(filteredAmounts, amounts[i])
		}

		if len(filteredAmounts) > 0 {
			filteredMonthlyAmounts[yearMonth] = filteredAmounts
		}
	}

	return filteredMonthlyAmounts
}

// monthlyTotalAmountKey identifies the amounts that add up together within a projection.
//
// The account is part of the key because the amount of a category is expressed in the currency of the
// account it belongs to: amounts of two accounts held in different currencies must stay apart so that
// the caller can convert each one before adding them up.
type monthlyTotalAmountKey struct {
	yearMonth  int32
	categoryId int64
	accountId  int64
}

// mergeMonthlyTotalAmounts adds up two sets of monthly amounts by year month, category and account
func mergeMonthlyTotalAmounts(left map[int32][]*models.TransactionTotalAmount, right map[int32][]*models.TransactionTotalAmount) map[int32][]*models.TransactionTotalAmount {
	amountsMap := make(map[monthlyTotalAmountKey]*models.TransactionTotalAmount)

	for _, monthlyAmounts := range []map[int32][]*models.TransactionTotalAmount{left, right} {
		for yearMonth, amounts := range monthlyAmounts {
			for i := 0; i < len(amounts); i++ {
				addToMonthlyTotalAmounts(amountsMap, yearMonth, amounts[i])
			}
		}
	}

	return groupTotalAmountsByYearMonth(amountsMap)
}

// addToMonthlyTotalAmounts adds the specified amount to the entry of its year month, category and
// account, creating the entry when it is the first amount of that group
func addToMonthlyTotalAmounts(amountsMap map[monthlyTotalAmountKey]*models.TransactionTotalAmount, yearMonth int32, amount *models.TransactionTotalAmount) {
	groupKey := monthlyTotalAmountKey{
		yearMonth:  yearMonth,
		categoryId: amount.CategoryId,
		accountId:  amount.AccountId,
	}

	totalAmount, exists := amountsMap[groupKey]

	if !exists {
		totalAmount = &models.TransactionTotalAmount{
			Type:       amount.Type,
			CategoryId: amount.CategoryId,
			AccountId:  amount.AccountId,
			Amount:     big.NewInt(0),
		}

		amountsMap[groupKey] = totalAmount
	}

	totalAmount.Amount.Add(totalAmount.Amount, amount.Amount)
}

// groupTotalAmountsByYearMonth turns a map keyed by year month, category and account into one keyed
// by year month only
func groupTotalAmountsByYearMonth(amountsMap map[monthlyTotalAmountKey]*models.TransactionTotalAmount) map[int32][]*models.TransactionTotalAmount {
	monthlyAmounts := make(map[int32][]*models.TransactionTotalAmount)

	for groupKey, totalAmount := range amountsMap {
		monthlyAmounts[groupKey.yearMonth] = append(monthlyAmounts[groupKey.yearMonth], totalAmount)
	}

	return monthlyAmounts
}
