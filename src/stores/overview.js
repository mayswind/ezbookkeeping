import { defineStore } from 'pinia';

import { useExchangeRatesStore } from './exchangeRates.js';

import { isNumber, isEquals } from '@/lib/common.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';

export const useOverviewStore = defineStore('overview', {
    state: () => ({
        transactionOverview: {},
        transactionOverviewStateInvalid: true
    }),
    actions: {
        updateTransactionOverviewInvalidState(invalidState) {
            this.transactionOverviewStateInvalid = invalidState;
        },
        resetTransactionOverview() {
            this.transactionOverview = {};
            this.transactionOverviewStateInvalid = true;
        },
        loadTransactionOverview({ defaultCurrency, dateRange, force }) {
            const self = this;
            const exchangeRatesStore = useExchangeRatesStore();

            if (!force && !self.transactionOverviewStateInvalid) {
                return new Promise((resolve) => {
                    resolve(self.transactionOverview);
                });
            }

            return new Promise((resolve, reject) => {
                services.getTransactionAmounts({
                    today: dateRange.today,
                    thisWeek: dateRange.thisWeek,
                    thisMonth: dateRange.thisMonth,
                    thisYear: dateRange.thisYear
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to get transaction overview' });
                        return;
                    }

                    const overview = data.result;

                    for (let field in overview) {
                        if (!Object.prototype.hasOwnProperty.call(overview, field)) {
                            continue;
                        }

                        const item = overview[field];

                        if (!item.amounts || !item.amounts.length) {
                            item.amounts = [];
                        }

                        let totalIncomeAmount = 0;
                        let totalExpenseAmount = 0;
                        let hasUnCalculatedTotalIncome = false;
                        let hasUnCalculatedTotalExpense = false;

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

                        item.incomeAmount = totalIncomeAmount;
                        item.expenseAmount = totalExpenseAmount;
                        item.incompleteIncomeAmount = hasUnCalculatedTotalIncome;
                        item.incompleteExpenseAmount = hasUnCalculatedTotalExpense;
                    }

                    if (self.transactionOverviewStateInvalid) {
                        self.updateTransactionOverviewInvalidState(false);
                    }

                    if (force && overview && isEquals(self.transactionOverview, overview)) {
                        reject({ message: 'Data is up to date' });
                        return;
                    }

                    self.transactionOverview = overview;

                    resolve(overview);
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
