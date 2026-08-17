package services

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// GetProjectedCategoryAmountsByMonth itself needs the user databases, and this package has no
// database test harness, so these tests cover the two halves it is made of: the simulation of the
// scheduled templates and the merge of both halves. Everything they assert -- the cutoff at "now",
// the exclusion of transfers, the per-account grouping that keeps currencies apart -- lives in
// those two.
//
// The remaining piece, the cutoff applied to the real half, is the maxTransactionUnixTime parameter
// of GetAccountsAndCategoriesMonthlyInflowAndOutflow, which narrows the SQL range and is verified
// manually in phase 4.3.

const projectionTestAccountId int64 = 3001
const projectionTestOtherAccountId int64 = 3002
const projectionTestCategoryId int64 = 2001
const projectionTestOtherCategoryId int64 = 2002

func newProjectionTemplateForTest(transactionType models.TransactionType, frequencyType models.TransactionScheduleFrequencyType, frequency string, amount int64) *models.TransactionTemplate {
	template := newScheduledTemplateForTest(frequencyType, frequency)
	template.Type = transactionType
	template.Amount = amount
	template.CategoryId = projectionTestCategoryId
	template.AccountId = projectionTestAccountId

	return template
}

func projectForTest(t *testing.T, templates []*models.TransactionTemplate, startYear int32, startMonth int32, endYear int32, endMonth int32, currentUnixTime int64) map[int32][]*models.TransactionTotalAmount {
	t.Helper()

	monthlyAmounts, err := projectScheduledCategoryAmountsByMonth(core.NewNullContext(), templates, startYear, startMonth, endYear, endMonth, time.UTC, false, currentUnixTime)
	assert.Nil(t, err)

	return monthlyAmounts
}

// amountOf returns the total amount of the specified category and account in the specified month, or
// -1 when there is no such entry
func amountOf(monthlyAmounts map[int32][]*models.TransactionTotalAmount, yearMonth int32, categoryId int64, accountId int64) int64 {
	amounts, exists := monthlyAmounts[yearMonth]

	if !exists {
		return -1
	}

	for i := 0; i < len(amounts); i++ {
		if amounts[i].CategoryId == categoryId && amounts[i].AccountId == accountId {
			return amounts[i].Amount.Int64()
		}
	}

	return -1
}

func newTotalAmountForTest(transactionDbType models.TransactionDbType, categoryId int64, accountId int64, amount int64) *models.TransactionTotalAmount {
	return &models.TransactionTotalAmount{
		Type:       transactionDbType,
		CategoryId: categoryId,
		AccountId:  accountId,
		Amount:     big.NewInt(amount),
	}
}

// --- range validation -----------------------------------------------------------------------

func TestValidateTransactionProjectionRange_Valid(t *testing.T) {
	assert.Nil(t, validateTransactionProjectionRange(2026, 8, 2026, 8))
	assert.Nil(t, validateTransactionProjectionRange(2026, 8, 2027, 7))
	assert.Nil(t, validateTransactionProjectionRange(2026, 1, 2030, 12))
}

func TestValidateTransactionProjectionRange_Incomplete(t *testing.T) {
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, validateTransactionProjectionRange(0, 0, 2026, 8))
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, validateTransactionProjectionRange(2026, 8, 0, 0))
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, validateTransactionProjectionRange(2026, 0, 2026, 8))
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, validateTransactionProjectionRange(2026, 13, 2026, 14))
}

func TestValidateTransactionProjectionRange_Inverted(t *testing.T) {
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, validateTransactionProjectionRange(2026, 9, 2026, 8))
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, validateTransactionProjectionRange(2027, 1, 2026, 12))
}

// TestValidateTransactionProjectionRange_TooLong keeps the period bounded: the real half loads every
// transaction of the period into memory.
func TestValidateTransactionProjectionRange_TooLong(t *testing.T) {
	// exactly 60 months
	assert.Nil(t, validateTransactionProjectionRange(2026, 1, 2030, 12))

	// 61 months
	assert.Equal(t, errs.ErrQueryItemsTooMuch, validateTransactionProjectionRange(2026, 1, 2031, 1))
	assert.Equal(t, errs.ErrQueryItemsTooMuch, validateTransactionProjectionRange(2000, 1, 2030, 1))
}

// --- simulated half -------------------------------------------------------------------------

func TestProjectScheduledCategoryAmounts_FutureMonthsOnly(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "15", 20000)
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2026, 8, 2026, 11, now)

	// The occurrence of August is already in the past, so it belongs to the real half
	assert.Equal(t, int64(-1), amountOf(monthlyAmounts, 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(20000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(20000), amountOf(monthlyAmounts, 202610, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(20000), amountOf(monthlyAmounts, 202611, projectionTestCategoryId, projectionTestAccountId))
}

// TestProjectScheduledCategoryAmounts_CurrentMonthIsSplitAtTheCutoff is the invariant the whole
// design rests on: within the month containing "now", only the occurrences still to come are
// simulated, so adding the real half cannot double count.
func TestProjectScheduledCategoryAmounts_CurrentMonthIsSplitAtTheCutoff(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "5,15,25", 10000)
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2026, 8, 2026, 8, now)

	// The 5th and the 15th already happened, only the 25th is still ahead
	assert.Equal(t, int64(10000), amountOf(monthlyAmounts, 202608, projectionTestCategoryId, projectionTestAccountId))
}

func TestProjectScheduledCategoryAmounts_FullyPastPeriodYieldsNothing(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0", 10000)
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2026, 1, 2026, 7, now)

	assert.Equal(t, 0, len(monthlyAmounts))
}

func TestProjectScheduledCategoryAmounts_DailyTemplateCountsEveryRemainingDay(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0", 100)
	now := utcUnixTime(2026, 8, 30, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2026, 8, 2026, 9, now)

	// Only the 31st is left in August, and September has 30 days
	assert.Equal(t, int64(100), amountOf(monthlyAmounts, 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(3000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestAccountId))
}

func TestProjectScheduledCategoryAmounts_IncomeAndExpenseKeepTheirType(t *testing.T) {
	income := newProjectionTemplateForTest(models.TRANSACTION_TYPE_INCOME, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 500000)
	expense := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 20000)
	expense.CategoryId = projectionTestOtherCategoryId
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{income, expense}, 2026, 9, 2026, 9, now)

	assert.Equal(t, 2, len(monthlyAmounts[202609]))

	for _, amount := range monthlyAmounts[202609] {
		if amount.CategoryId == projectionTestCategoryId {
			assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, amount.Type)
			assert.Equal(t, int64(500000), amount.Amount.Int64())
		} else {
			assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, amount.Type)
			assert.Equal(t, int64(20000), amount.Amount.Int64())
		}
	}
}

// TestProjectScheduledCategoryAmounts_TransferTemplatesAreExcluded keeps both halves consistent:
// transfers are dropped from the real half too, so a transfer template must not appear here either.
func TestProjectScheduledCategoryAmounts_TransferTemplatesAreExcluded(t *testing.T) {
	transfer := newProjectionTemplateForTest(models.TRANSACTION_TYPE_TRANSFER, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 100000)
	transfer.RelatedAccountId = projectionTestOtherAccountId
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{transfer}, 2026, 9, 2026, 12, now)

	assert.Equal(t, 0, len(monthlyAmounts))
}

// TestProjectScheduledCategoryAmounts_SameCategoryOnAccountsOfDifferentCurrenciesStaysApart is why
// the account is part of the grouping key: the amount of a template is expressed in the currency of
// its account, so adding both into a single figure here would mix currencies beyond repair. The
// caller converts each one before adding them up.
func TestProjectScheduledCategoryAmounts_SameCategoryOnAccountsOfDifferentCurrenciesStaysApart(t *testing.T) {
	inPesos := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 150000)
	inDollars := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 1000)
	inDollars.AccountId = projectionTestOtherAccountId
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{inPesos, inDollars}, 2026, 9, 2026, 9, now)

	assert.Equal(t, 2, len(monthlyAmounts[202609]))
	assert.Equal(t, int64(150000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(1000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestOtherAccountId))
}

func TestProjectScheduledCategoryAmounts_SameCategoryAndAccountAddsUp(t *testing.T) {
	first := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 30000)
	second := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "20", 12000)
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{first, second}, 2026, 9, 2026, 9, now)

	assert.Equal(t, 1, len(monthlyAmounts[202609]))
	assert.Equal(t, int64(42000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestAccountId))
}

func TestProjectScheduledCategoryAmounts_RespectsTemplateEndTime(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "10", 5000)
	endTime := utcUnixTime(2026, 10, 31, 23, 59)
	template.ScheduledEndTime = &endTime
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2026, 9, 2026, 12, now)

	assert.Equal(t, int64(5000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(5000), amountOf(monthlyAmounts, 202610, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(-1), amountOf(monthlyAmounts, 202611, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(-1), amountOf(monthlyAmounts, 202612, projectionTestCategoryId, projectionTestAccountId))
}

// TestProjectScheduledCategoryAmounts_MalformedTemplateIsSkipped keeps a single broken template from
// failing the whole projection, the way the cron skips it and carries on.
func TestProjectScheduledCategoryAmounts_MalformedTemplateIsSkipped(t *testing.T) {
	malformed := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "monday", 10000)
	healthy := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 7000)
	healthy.CategoryId = projectionTestOtherCategoryId
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{malformed, healthy}, 2026, 9, 2026, 9, now)

	assert.Equal(t, 1, len(monthlyAmounts[202609]))
	assert.Equal(t, int64(7000), amountOf(monthlyAmounts, 202609, projectionTestOtherCategoryId, projectionTestAccountId))
}

// TestProjectScheduledCategoryAmounts_HiddenTemplateIsProjected mirrors the cron, which does not
// filter hidden templates out and keeps creating their transactions.
func TestProjectScheduledCategoryAmounts_HiddenTemplateIsProjected(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 8000)
	template.Hidden = true
	now := utcUnixTime(2026, 8, 17, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2026, 9, 2026, 9, now)

	assert.Equal(t, int64(8000), amountOf(monthlyAmounts, 202609, projectionTestCategoryId, projectionTestAccountId))
}

// TestProjectScheduledCategoryAmounts_NegativeDayIsProjectedPerMonth is the projection level check of
// the bug fixed in phase 1: a template on the last day of the month must land on the last day of
// every projected month, whatever its length.
func TestProjectScheduledCategoryAmounts_NegativeDayIsProjectedPerMonth(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_INCOME, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-1", 250000)
	now := utcUnixTime(2026, 12, 15, 12, 0)

	monthlyAmounts := projectForTest(t, []*models.TransactionTemplate{template}, 2027, 1, 2027, 4, now)

	assert.Equal(t, int64(250000), amountOf(monthlyAmounts, 202701, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(250000), amountOf(monthlyAmounts, 202702, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(250000), amountOf(monthlyAmounts, 202703, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(250000), amountOf(monthlyAmounts, 202704, projectionTestCategoryId, projectionTestAccountId))
}

func TestProjectScheduledCategoryAmounts_NoTemplates(t *testing.T) {
	monthlyAmounts := projectForTest(t, nil, 2026, 9, 2026, 12, utcUnixTime(2026, 8, 17, 12, 0))

	assert.Equal(t, 0, len(monthlyAmounts))
}

func TestProjectScheduledCategoryAmounts_InvalidRange(t *testing.T) {
	template := newProjectionTemplateForTest(models.TRANSACTION_TYPE_EXPENSE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1", 1000)

	_, err := projectScheduledCategoryAmountsByMonth(core.NewNullContext(), []*models.TransactionTemplate{template}, 2026, 13, 2026, 14, time.UTC, false, utcUnixTime(2026, 8, 17, 12, 0))

	assert.NotNil(t, err)
}

// --- transfers on the real half ---------------------------------------------------------------

func TestRemoveTransferTotalAmounts(t *testing.T) {
	monthlyAmounts := map[int32][]*models.TransactionTotalAmount{
		202608: {
			newTotalAmountForTest(models.TRANSACTION_DB_TYPE_INCOME, projectionTestCategoryId, projectionTestAccountId, 500000),
			newTotalAmountForTest(models.TRANSACTION_DB_TYPE_TRANSFER_OUT, projectionTestOtherCategoryId, projectionTestAccountId, 100000),
			newTotalAmountForTest(models.TRANSACTION_DB_TYPE_TRANSFER_IN, projectionTestOtherCategoryId, projectionTestOtherAccountId, 100000),
		},
		202609: {
			newTotalAmountForTest(models.TRANSACTION_DB_TYPE_TRANSFER_OUT, projectionTestOtherCategoryId, projectionTestAccountId, 100000),
		},
	}

	filtered := removeTransferTotalAmounts(monthlyAmounts)

	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, 1, len(filtered[202608]))
	assert.Equal(t, int64(500000), amountOf(filtered, 202608, projectionTestCategoryId, projectionTestAccountId))

	// A month left with nothing but transfers disappears instead of staying as an empty entry
	_, exists := filtered[202609]
	assert.False(t, exists)
}

// --- merge ------------------------------------------------------------------------------------

// TestMergeMonthlyTotalAmounts_PastPresentAndFuture walks the three shapes a projected month can
// have: only real amounts, both halves, and only simulated amounts.
func TestMergeMonthlyTotalAmounts_PastPresentAndFuture(t *testing.T) {
	actual := map[int32][]*models.TransactionTotalAmount{
		202607: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 30000)},
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 20000)},
	}

	projected := map[int32][]*models.TransactionTotalAmount{
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 10000)},
		202609: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 30000)},
	}

	merged := mergeMonthlyTotalAmounts(actual, projected)

	assert.Equal(t, 3, len(merged))
	assert.Equal(t, int64(30000), amountOf(merged, 202607, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(30000), amountOf(merged, 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(30000), amountOf(merged, 202609, projectionTestCategoryId, projectionTestAccountId))

	// The month in the middle is one entry holding the sum, not two competing ones
	assert.Equal(t, 1, len(merged[202608]))
}

func TestMergeMonthlyTotalAmounts_KeepsAccountsApart(t *testing.T) {
	actual := map[int32][]*models.TransactionTotalAmount{
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 150000)},
	}

	projected := map[int32][]*models.TransactionTotalAmount{
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestOtherAccountId, 1000)},
	}

	merged := mergeMonthlyTotalAmounts(actual, projected)

	assert.Equal(t, 2, len(merged[202608]))
	assert.Equal(t, int64(150000), amountOf(merged, 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(1000), amountOf(merged, 202608, projectionTestCategoryId, projectionTestOtherAccountId))
}

func TestMergeMonthlyTotalAmounts_KeepsCategoriesApart(t *testing.T) {
	actual := map[int32][]*models.TransactionTotalAmount{
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 4000)},
	}

	projected := map[int32][]*models.TransactionTotalAmount{
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestOtherCategoryId, projectionTestAccountId, 6000)},
	}

	merged := mergeMonthlyTotalAmounts(actual, projected)

	assert.Equal(t, 2, len(merged[202608]))
	assert.Equal(t, int64(4000), amountOf(merged, 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(6000), amountOf(merged, 202608, projectionTestOtherCategoryId, projectionTestAccountId))
}

// TestMergeMonthlyTotalAmounts_DoesNotMutateItsInputs matters because the real half comes straight
// out of the aggregation of the transactions of the user: adding the simulated amounts on top of
// those big.Int values in place would corrupt them.
func TestMergeMonthlyTotalAmounts_DoesNotMutateItsInputs(t *testing.T) {
	actualAmount := newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 20000)
	projectedAmount := newTotalAmountForTest(models.TRANSACTION_DB_TYPE_EXPENSE, projectionTestCategoryId, projectionTestAccountId, 10000)

	actual := map[int32][]*models.TransactionTotalAmount{202608: {actualAmount}}
	projected := map[int32][]*models.TransactionTotalAmount{202608: {projectedAmount}}

	merged := mergeMonthlyTotalAmounts(actual, projected)

	assert.Equal(t, int64(30000), amountOf(merged, 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(20000), actualAmount.Amount.Int64())
	assert.Equal(t, int64(10000), projectedAmount.Amount.Int64())
}

func TestMergeMonthlyTotalAmounts_EmptySides(t *testing.T) {
	amounts := map[int32][]*models.TransactionTotalAmount{
		202608: {newTotalAmountForTest(models.TRANSACTION_DB_TYPE_INCOME, projectionTestCategoryId, projectionTestAccountId, 500000)},
	}

	assert.Equal(t, 0, len(mergeMonthlyTotalAmounts(nil, nil)))
	assert.Equal(t, int64(500000), amountOf(mergeMonthlyTotalAmounts(amounts, nil), 202608, projectionTestCategoryId, projectionTestAccountId))
	assert.Equal(t, int64(500000), amountOf(mergeMonthlyTotalAmounts(nil, amounts), 202608, projectionTestCategoryId, projectionTestAccountId))
}

// --- template type mapping --------------------------------------------------------------------

func TestGetTransactionDbTypeByTemplateType(t *testing.T) {
	incomeType, err := getTransactionDbTypeByTemplateType(models.TRANSACTION_TYPE_INCOME)
	assert.Nil(t, err)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_INCOME, incomeType)

	expenseType, err := getTransactionDbTypeByTemplateType(models.TRANSACTION_TYPE_EXPENSE)
	assert.Nil(t, err)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_EXPENSE, expenseType)

	transferType, err := getTransactionDbTypeByTemplateType(models.TRANSACTION_TYPE_TRANSFER)
	assert.Nil(t, err)
	assert.Equal(t, models.TRANSACTION_DB_TYPE_TRANSFER_OUT, transferType)

	_, err = getTransactionDbTypeByTemplateType(models.TransactionType(99))
	assert.NotNil(t, err)
}
