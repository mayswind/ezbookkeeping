import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utilities from '../lib/utilities/index.js';

import { getExchangedAmount } from './exchangeRates.js';

import {
    LOAD_TRANSACTION_OVERVIEW,
    UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE,
} from './mutations.js';

export function loadTransactionOverview(context, { defaultCurrency, dateRange, force }) {
    if (!force && !context.state.transactionOverviewStateInvalid) {
        return new Promise((resolve) => {
            resolve(context.state.transactionOverview);
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
                        const incomeAmount = getExchangedAmount(context.state)(amount.incomeAmount, amount.currency, defaultCurrency);
                        const expenseAmount = getExchangedAmount(context.state)(amount.expenseAmount, amount.currency, defaultCurrency);

                        if (utilities.isNumber(incomeAmount)) {
                            totalIncomeAmount += Math.floor(incomeAmount);
                        } else {
                            hasUnCalculatedTotalIncome = true;
                        }

                        if (utilities.isNumber(expenseAmount)) {
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

            context.commit(LOAD_TRANSACTION_OVERVIEW, overview);

            if (context.state.transactionOverviewStateInvalid) {
                context.commit(UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE, false);
            }

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
