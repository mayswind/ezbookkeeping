import { defineStore } from 'pinia';

import { useUserStore } from './user.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import { isNumber, isEquals } from '@/lib/common.js';
import {
    getTodayFirstUnixTime,
    getTodayLastUnixTime,
    getThisWeekFirstUnixTime,
    getThisWeekLastUnixTime,
    getThisMonthFirstUnixTime,
    getThisMonthLastUnixTime,
    getThisYearFirstUnixTime,
    getThisYearLastUnixTime
} from '@/lib/datetime.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';

function updateTransactionDateRange(state) {
    const userStore = useUserStore();

    state.transactionDataRange.today.startTime = getTodayFirstUnixTime();
    state.transactionDataRange.today.endTime = getTodayLastUnixTime();

    state.transactionDataRange.thisWeek.startTime = getThisWeekFirstUnixTime(userStore.currentUserFirstDayOfWeek);
    state.transactionDataRange.thisWeek.endTime = getThisWeekLastUnixTime(userStore.currentUserFirstDayOfWeek);

    state.transactionDataRange.thisMonth.startTime = getThisMonthFirstUnixTime();
    state.transactionDataRange.thisMonth.endTime = getThisMonthLastUnixTime();

    state.transactionDataRange.thisYear.startTime = getThisYearFirstUnixTime();
    state.transactionDataRange.thisYear.endTime = getThisYearLastUnixTime();
}

export const useOverviewStore = defineStore('overview', {
    state: () => ({
        transactionDataRange: {
            today: {
                startTime: getTodayFirstUnixTime(),
                endTime: getTodayLastUnixTime()
            },
            thisWeek: {
                startTime: getThisWeekFirstUnixTime(useUserStore().currentUserFirstDayOfWeek),
                endTime: getThisWeekLastUnixTime(useUserStore().currentUserFirstDayOfWeek)
            },
            thisMonth: {
                startTime: getThisMonthFirstUnixTime(),
                endTime: getThisMonthLastUnixTime()
            },
            thisYear: {
                startTime: getThisYearFirstUnixTime(),
                endTime: getThisYearLastUnixTime()
            }
        },
        transactionOverviewData: {},
        transactionOverviewStateInvalid: true
    }),
    getters: {
        transactionOverview(state) {
            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();

            const overviewData = state.transactionOverviewData;

            if (!overviewData || !overviewData.thisMonth) {
                return {
                    thisMonth: {
                        valid: false,
                        incomeAmount: 0,
                        expenseAmount: 0,
                        incompleteIncomeAmount: false,
                        incompleteExpenseAmount: false
                    }
                };
            }

            const finalOverviewData = {};
            const defaultCurrency = userStore.currentUserDefaultCurrency;

            [ 'today', 'thisWeek', 'thisMonth', 'thisYear' ].forEach(field => {
                if (!Object.prototype.hasOwnProperty.call(overviewData, field)) {
                    return;
                }

                const item = overviewData[field];

                if (!item) {
                    return;
                }

                let totalIncomeAmount = 0;
                let totalExpenseAmount = 0;
                let hasUnCalculatedTotalIncome = false;
                let hasUnCalculatedTotalExpense = false;

                if (item.amounts) {
                    for (let i = 0; i < item.amounts.length; i++) {
                        const amount = item.amounts[i];

                        if (amount.currency !== defaultCurrency) {
                            const incomeAmount = exchangeRatesStore.getExchangedAmount(amount.incomeAmount, amount.currency, defaultCurrency);
                            const expenseAmount = exchangeRatesStore.getExchangedAmount(amount.expenseAmount, amount.currency, defaultCurrency);

                            if (isNumber(incomeAmount)) {
                                totalIncomeAmount += Math.floor(incomeAmount);
                            } else {
                                hasUnCalculatedTotalIncome = true;
                            }

                            if (isNumber(expenseAmount)) {
                                totalExpenseAmount += Math.floor(expenseAmount);
                            } else {
                                hasUnCalculatedTotalExpense = true;
                            }
                        } else {
                            totalIncomeAmount += amount.incomeAmount;
                            totalExpenseAmount += amount.expenseAmount;
                        }
                    }
                }

                finalOverviewData[field] = {
                    valid: true,
                    incomeAmount: totalIncomeAmount,
                    expenseAmount: totalExpenseAmount,
                    incompleteIncomeAmount: hasUnCalculatedTotalIncome,
                    incompleteExpenseAmount: hasUnCalculatedTotalExpense,
                    amounts: item.amounts || []
                };
            });

            return finalOverviewData;
        }
    },
    actions: {
        updateTransactionOverviewInvalidState(invalidState) {
            this.transactionOverviewStateInvalid = invalidState;
        },
        resetTransactionOverview() {
            updateTransactionDateRange(this);
            this.transactionOverviewData = {};
            this.transactionOverviewStateInvalid = true;
        },
        loadTransactionOverview({ force }) {
            const self = this;
            let dateChanged = false;

            if (self.transactionDataRange.today.startTime !== getTodayFirstUnixTime()) {
                dateChanged = true;
                updateTransactionDateRange(self);
            }

            if (!dateChanged && !force && !self.transactionOverviewStateInvalid) {
                return new Promise((resolve) => {
                    resolve(self.transactionOverviewData);
                });
            }

            return new Promise((resolve, reject) => {
                services.getTransactionAmounts({
                    today: self.transactionDataRange.today,
                    thisWeek: self.transactionDataRange.thisWeek,
                    thisMonth: self.transactionDataRange.thisMonth,
                    thisYear: self.transactionDataRange.thisYear
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to get transaction overview' });
                        return;
                    }

                    if (self.transactionOverviewStateInvalid) {
                        self.updateTransactionOverviewInvalidState(false);
                    }

                    if (force && data.result && isEquals(self.transactionOverviewData, data.result)) {
                        reject({ message: 'Data is up to date' });
                        return;
                    }

                    self.transactionOverviewData = data.result;

                    resolve(data.result);
                }).catch(error => {
                    if (force) {
                        logger.error('failed to force load transaction overview', error);
                    } else {
                        logger.error('failed to load transaction overview', error);
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to get transaction overview' });
                    } else {
                        reject(error);
                    }
                });
            });
        }
    }
});
