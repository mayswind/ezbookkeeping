import statisticsConstants from '../consts/statistics.js';
import categoryConstants from '../consts/category.js';
import iconConstants from '../consts/icon.js';
import colorConstants from '../consts/color.js';
import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utilities from '../lib/utilities/index.js';

import { getExchangedAmount } from './exchangeRates.js';

import {
    LOAD_TRANSACTION_STATISTICS,
    INIT_TRANSACTION_STATISTICS_FILTER,
    UPDATE_TRANSACTION_STATISTICS_FILTER,
    UPDATE_TRANSACTION_STATISTICS_INVALID_STATE
} from './mutations.js';

export function loadTransactionStatistics(context, { defaultCurrency }) {
    return new Promise((resolve, reject) => {
        services.getTransactionStatistics({
            startTime: context.state.transactionStatisticsFilter.startTime,
            endTime: context.state.transactionStatisticsFilter.endTime
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get transaction statistics' });
                return;
            }

            context.commit(LOAD_TRANSACTION_STATISTICS, {
                statistics: data.result,
                defaultCurrency: defaultCurrency
            });

            if (context.state.transactionStatisticsStateInvalid) {
                context.commit(UPDATE_TRANSACTION_STATISTICS_INVALID_STATE, false);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to get transaction statistics', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get transaction statistics' });
            } else {
                reject(error);
            }
        });
    });
}

export function initTransactionStatisticsFilter(context, filter) {
    context.commit(INIT_TRANSACTION_STATISTICS_FILTER, filter);
}

export function updateTransactionStatisticsFilter(context, filter) {
    context.commit(UPDATE_TRANSACTION_STATISTICS_FILTER, filter);
}

export function statisticsItemsByTransactionStatisticsData(state) {
    if (!state.transactionStatistics || !state.transactionStatistics.items) {
        return null;
    }

    const allDataItems = {};
    let totalAmount = 0;
    let totalNonNegativeAmount = 0;

    for (let i = 0; i < state.transactionStatistics.items.length; i++) {
        const item = state.transactionStatistics.items[i];

        if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category) {
            continue;
        }

        if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type) {
            if (item.category.type !== categoryConstants.allCategoryTypes.Expense) {
                continue;
            }
        } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
            if (item.category.type !== categoryConstants.allCategoryTypes.Income) {
                continue;
            }
        } else {
            continue;
        }

        if (state.transactionStatisticsFilter.filterAccountIds && state.transactionStatisticsFilter.filterAccountIds[item.account.id]) {
            continue;
        }

        if (state.transactionStatisticsFilter.filterCategoryIds && state.transactionStatisticsFilter.filterCategoryIds[item.category.id]) {
            continue;
        }

        if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type) {
            if (utilities.isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems[item.account.id];

                if (data) {
                    data.totalAmount += item.amountInDefaultCurrency;
                } else {
                    data = {
                        name: item.account.name,
                        type: 'account',
                        id: item.account.id,
                        icon: item.account.icon || iconConstants.defaultAccountIcon.icon,
                        color: item.account.color || colorConstants.defaultAccountColor,
                        hidden: item.primaryAccount.hidden || item.account.hidden,
                        displayOrders: [item.primaryAccount.category, item.primaryAccount.displayOrder, item.account.displayOrder],
                        totalAmount: item.amountInDefaultCurrency
                    }
                }

                totalAmount += item.amountInDefaultCurrency;

                if (item.amountInDefaultCurrency > 0) {
                    totalNonNegativeAmount += item.amountInDefaultCurrency;
                }

                allDataItems[item.account.id] = data;
            }
        } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type) {
            if (utilities.isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems[item.primaryCategory.id];

                if (data) {
                    data.totalAmount += item.amountInDefaultCurrency;
                } else {
                    data = {
                        name: item.primaryCategory.name,
                        type: 'category',
                        id: item.primaryCategory.id,
                        icon: item.primaryCategory.icon || iconConstants.defaultCategoryIcon.icon,
                        color: item.primaryCategory.color || colorConstants.defaultCategoryColor,
                        hidden: item.primaryCategory.hidden,
                        displayOrders: [item.primaryCategory.type, item.primaryCategory.displayOrder],
                        totalAmount: item.amountInDefaultCurrency
                    }
                }

                totalAmount += item.amountInDefaultCurrency;

                if (item.amountInDefaultCurrency > 0) {
                    totalNonNegativeAmount += item.amountInDefaultCurrency;
                }

                allDataItems[item.primaryCategory.id] = data;
            }
        } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type ||
            state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
            if (utilities.isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems[item.category.id];

                if (data) {
                    data.totalAmount += item.amountInDefaultCurrency;
                } else {
                    data = {
                        name: item.category.name,
                        type: 'category',
                        id: item.category.id,
                        icon: item.category.icon || iconConstants.defaultCategoryIcon.icon,
                        color: item.category.color || colorConstants.defaultCategoryColor,
                        hidden: item.primaryCategory.hidden || item.category.hidden,
                        displayOrders: [item.primaryCategory.type, item.primaryCategory.displayOrder, item.category.displayOrder],
                        totalAmount: item.amountInDefaultCurrency
                    }
                }

                totalAmount += item.amountInDefaultCurrency;

                if (item.amountInDefaultCurrency > 0) {
                    totalNonNegativeAmount += item.amountInDefaultCurrency;
                }

                allDataItems[item.category.id] = data;
            }
        }
    }

    return {
        totalAmount: totalAmount,
        totalNonNegativeAmount: totalNonNegativeAmount,
        items: allDataItems
    }
}

export function statisticsItemsByAccountsData(state, getters) {
    if (!getters.allPlainAccounts) {
        return null;
    }

    const allDataItems = {};
    let totalAmount = 0;
    let totalNonNegativeAmount = 0;

    for (let i = 0; i < getters.allPlainAccounts.length; i++) {
        const account = getters.allPlainAccounts[i];

        if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type) {
            if (!account.isAsset) {
                continue;
            }
        } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
            if (!account.isLiability) {
                continue;
            }
        }

        if (state.transactionStatisticsFilter.filterAccountIds && state.transactionStatisticsFilter.filterAccountIds[account.id]) {
            continue;
        }

        let primaryAccount = state.allAccountsMap[account.parentId];

        if (!primaryAccount) {
            primaryAccount = account;
        }

        let amount = account.balance;

        if (account.currency !== getters.currentUserDefaultCurrency) {
            amount = Math.floor(getExchangedAmount(state)(amount, account.currency, getters.currentUserDefaultCurrency));

            if (!utilities.isNumber(amount)) {
                continue;
            }
        }

        if (account.isLiability) {
            amount = -amount;
        }

        const data = {
            name: account.name,
            type: 'account',
            id: account.id,
            icon: account.icon || iconConstants.defaultAccountIcon.icon,
            color: account.color || colorConstants.defaultAccountColor,
            hidden: primaryAccount.hidden || account.hidden,
            displayOrders: [primaryAccount.category, primaryAccount.displayOrder, account.displayOrder],
            totalAmount: amount
        };

        totalAmount += amount;

        if (amount > 0) {
            totalNonNegativeAmount += amount;
        }

        allDataItems[account.id] = data;
    }

    return {
        totalAmount: totalAmount,
        totalNonNegativeAmount: totalNonNegativeAmount,
        items: allDataItems
    }
}
