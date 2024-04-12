import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.js';
import { useUserStore } from './user.js';
import { useAccountsStore } from './account.js';
import { useTransactionCategoriesStore } from './transactionCategory.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import datetimeConstants from '@/consts/datetime.js';
import statisticsConstants from '@/consts/statistics.js';
import categoryConstants from '@/consts/category.js';
import iconConstants from '@/consts/icon.js';
import colorConstants from '@/consts/color.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import {
    isEquals,
    isNumber,
    isObject
} from '@/lib/common.js';
import {
    getDateRangeByDateType
} from '@/lib/datetime.js';

export const useStatisticsStore = defineStore('statistics', {
    state: () => ({
        transactionStatisticsFilter: {
            dateType: statisticsConstants.defaultDataRangeType,
            startTime: 0,
            endTime: 0,
            chartType: statisticsConstants.defaultChartType,
            chartDataType: statisticsConstants.defaultChartDataType,
            filterAccountIds: {},
            filterCategoryIds: {}
        },
        transactionStatisticsData: {},
        transactionStatisticsStateInvalid: true
    }),
    getters: {
        transactionStatisticsChartDataCategory(state) {
            if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                return 'account';
            } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
                return 'category';
            } else {
                return '';
            }
        },
        transactionStatistics(state) {
            const statistics = state.transactionStatisticsData;
            const finalStatistics = {
                startTime: statistics.startTime,
                endTime: statistics.endTime,
                items: []
            };

            if (statistics && statistics.items && statistics.items.length) {
                const userStore = useUserStore();
                const accountsStore = useAccountsStore();
                const transactionCategoriesStore = useTransactionCategoriesStore();
                const exchangeRatesStore = useExchangeRatesStore();

                const defaultCurrency = userStore.currentUserDefaultCurrency;

                for (let i = 0; i < statistics.items.length; i++) {
                    const dataItem = statistics.items[i];
                    const item = {
                        categoryId: dataItem.categoryId,
                        accountId: dataItem.accountId,
                        amount: dataItem.amount
                    };

                    if (item.accountId) {
                        item.account = accountsStore.allAccountsMap[item.accountId];
                    }

                    if (item.account && item.account.parentId !== '0') {
                        item.primaryAccount = accountsStore.allAccountsMap[item.account.parentId];
                    } else {
                        item.primaryAccount = item.account;
                    }

                    if (item.categoryId) {
                        item.category = transactionCategoriesStore.allTransactionCategoriesMap[item.categoryId];
                    }

                    if (item.category && item.category.parentId !== '0') {
                        item.primaryCategory = transactionCategoriesStore.allTransactionCategoriesMap[item.category.parentId];
                    } else {
                        item.primaryCategory = item.category;
                    }

                    if (item.account && item.account.currency !== defaultCurrency) {
                        const amount = exchangeRatesStore.getExchangedAmount(item.amount, item.account.currency, defaultCurrency);

                        if (isNumber(amount)) {
                            item.amountInDefaultCurrency = Math.floor(amount);
                        }
                    } else if (item.account && item.account.currency === defaultCurrency) {
                        item.amountInDefaultCurrency = item.amount;
                    } else {
                        item.amountInDefaultCurrency = null;
                    }

                    finalStatistics.items.push(item);
                }
            }

            return finalStatistics;
        },
        statisticsItemsByTransactionStatisticsData(state) {
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
                    if (isNumber(item.amountInDefaultCurrency)) {
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
                    if (isNumber(item.amountInDefaultCurrency)) {
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
                    if (isNumber(item.amountInDefaultCurrency)) {
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
        },
        statisticsItemsByAccountsData(state) {
            const userStore = useUserStore();
            const accountsStore = useAccountsStore();
            const exchangeRatesStore = useExchangeRatesStore();

            if (!accountsStore.allPlainAccounts) {
                return null;
            }

            const allDataItems = {};
            let totalAmount = 0;
            let totalNonNegativeAmount = 0;

            for (let i = 0; i < accountsStore.allPlainAccounts.length; i++) {
                const account = accountsStore.allPlainAccounts[i];

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

                let primaryAccount = accountsStore.allAccountsMap[account.parentId];

                if (!primaryAccount) {
                    primaryAccount = account;
                }

                let amount = account.balance;

                if (account.currency !== userStore.currentUserDefaultCurrency) {
                    amount = Math.floor(exchangeRatesStore.getExchangedAmount(amount, account.currency, userStore.currentUserDefaultCurrency));

                    if (!isNumber(amount)) {
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
        },
        statisticsData(state) {
            let combinedData = {
                items: [],
                totalAmount: 0
            };

            if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
                combinedData = state.statisticsItemsByTransactionStatisticsData;
            } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                combinedData = state.statisticsItemsByAccountsData;
            }

            const allStatisticsItems = [];

            for (let id in combinedData.items) {
                if (!Object.prototype.hasOwnProperty.call(combinedData.items, id)) {
                    continue;
                }

                const data = combinedData.items[id];

                if (data.totalAmount > 0) {
                    data.percent = data.totalAmount * 100 / combinedData.totalNonNegativeAmount;
                } else {
                    data.percent = 0;
                }

                if (data.percent < 0) {
                    data.percent = 0;
                }

                allStatisticsItems.push(data);
            }

            if (state.transactionStatisticsFilter.sortingType === statisticsConstants.allSortingTypes.DisplayOrder.type) {
                allStatisticsItems.sort(function (data1, data2) {
                    for (let i = 0; i < Math.min(data1.displayOrders.length, data2.displayOrders.length); i++) {
                        if (data1.displayOrders[i] !== data2.displayOrders[i]) {
                            return data1.displayOrders[i] - data2.displayOrders[i]; // asc
                        }
                    }

                    return data1.name.localeCompare(data2.name, undefined, { // asc
                        numeric: true,
                        sensitivity: 'base'
                    });
                });
            } else if (state.transactionStatisticsFilter.sortingType === statisticsConstants.allSortingTypes.Name.type) {
                allStatisticsItems.sort(function (data1, data2) {
                    return data1.name.localeCompare(data2.name, undefined, { // asc
                        numeric: true,
                        sensitivity: 'base'
                    });
                });
            } else {
                allStatisticsItems.sort(function (data1, data2) {
                    if (data1.totalAmount !== data2.totalAmount) {
                        return data2.totalAmount - data1.totalAmount; // desc
                    }

                    return data1.name.localeCompare(data2.name, undefined, { // asc
                        numeric: true,
                        sensitivity: 'base'
                    });
                });
            }

            return {
                totalAmount: combinedData.totalAmount,
                items: allStatisticsItems
            };
        }
    },
    actions: {
        updateTransactionStatisticsInvalidState(invalidState) {
            this.transactionStatisticsStateInvalid = invalidState;
        },
        resetTransactionStatistics() {
            this.transactionStatisticsFilter.dateType = statisticsConstants.defaultDataRangeType;
            this.transactionStatisticsFilter.startTime = 0;
            this.transactionStatisticsFilter.endTime = 0;
            this.transactionStatisticsFilter.chartType = statisticsConstants.defaultChartType;
            this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            this.transactionStatisticsFilter.filterAccountIds = {};
            this.transactionStatisticsFilter.filterCategoryIds = {};
            this.transactionStatisticsData = {};
            this.transactionStatisticsStateInvalid = true;
        },
        initTransactionStatisticsFilter(filter) {
            if (!filter) {
                const settingsStore = useSettingsStore();
                const userStore = useUserStore();

                let defaultChartType = settingsStore.appSettings.statistics.defaultChartType;

                if (defaultChartType !== statisticsConstants.allChartTypes.Pie && defaultChartType !== statisticsConstants.allChartTypes.Bar) {
                    defaultChartType = statisticsConstants.defaultChartType;
                }

                let defaultChartDataType = settingsStore.appSettings.statistics.defaultChartDataType;

                if (defaultChartDataType < statisticsConstants.allChartDataTypes.ExpenseByAccount.type || defaultChartDataType > statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                    defaultChartDataType = statisticsConstants.defaultChartDataType;
                }

                let defaultDateRange = settingsStore.appSettings.statistics.defaultDataRangeType;

                if (defaultDateRange < datetimeConstants.allDateRanges.All.type || defaultDateRange >= datetimeConstants.allDateRanges.Custom.type) {
                    defaultDateRange = statisticsConstants.defaultDataRangeType;
                }

                let defaultSortType = settingsStore.appSettings.statistics.defaultSortingType;

                if (defaultSortType < statisticsConstants.allSortingTypes.Amount.type || defaultSortType > statisticsConstants.allSortingTypes.Name.type) {
                    defaultSortType = statisticsConstants.defaultSortingType;
                }

                const dateRange = getDateRangeByDateType(defaultDateRange, userStore.currentUserFirstDayOfWeek);

                filter = {
                    dateType: dateRange ? dateRange.dateType : undefined,
                    startTime: dateRange ? dateRange.minTime : undefined,
                    endTime: dateRange ? dateRange.maxTime : undefined,
                    chartType: defaultChartType,
                    chartDataType: defaultChartDataType,
                    filterAccountIds: settingsStore.appSettings.statistics.defaultAccountFilter || {},
                    filterCategoryIds: settingsStore.appSettings.statistics.defaultTransactionCategoryFilter || {},
                    sortingType: defaultSortType,
                };
            }

            if (filter && isNumber(filter.dateType)) {
                this.transactionStatisticsFilter.dateType = filter.dateType;
            } else {
                this.transactionStatisticsFilter.dateType = statisticsConstants.defaultDataRangeType;
            }

            if (filter && isNumber(filter.startTime)) {
                this.transactionStatisticsFilter.startTime = filter.startTime;
            } else {
                this.transactionStatisticsFilter.startTime = 0;
            }

            if (filter && isNumber(filter.endTime)) {
                this.transactionStatisticsFilter.endTime = filter.endTime;
            } else {
                this.transactionStatisticsFilter.endTime = 0;
            }

            if (filter && isNumber(filter.chartType)) {
                this.transactionStatisticsFilter.chartType = filter.chartType;
            } else {
                this.transactionStatisticsFilter.chartType = statisticsConstants.defaultChartType;
            }

            if (filter && isNumber(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            } else {
                this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            }

            if (filter && isObject(filter.filterAccountIds)) {
                this.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            } else {
                this.transactionStatisticsFilter.filterAccountIds = {};
            }

            if (filter && isObject(filter.filterCategoryIds)) {
                this.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            } else {
                this.transactionStatisticsFilter.filterCategoryIds = {};
            }

            if (filter && isNumber(filter.sortingType)) {
                this.transactionStatisticsFilter.sortingType = filter.sortingType;
            } else {
                this.transactionStatisticsFilter.sortingType = statisticsConstants.defaultSortingType;
            }
        },
        updateTransactionStatisticsFilter(filter) {
            if (filter && isNumber(filter.dateType)) {
                this.transactionStatisticsFilter.dateType = filter.dateType;
            }

            if (filter && isNumber(filter.startTime)) {
                this.transactionStatisticsFilter.startTime = filter.startTime;
            }

            if (filter && isNumber(filter.endTime)) {
                this.transactionStatisticsFilter.endTime = filter.endTime;
            }

            if (filter && isNumber(filter.chartType)) {
                this.transactionStatisticsFilter.chartType = filter.chartType;
            }

            if (filter && isNumber(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            }

            if (filter && isObject(filter.filterAccountIds)) {
                this.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            }

            if (filter && isObject(filter.filterCategoryIds)) {
                this.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            }

            if (filter && isNumber(filter.sortingType)) {
                this.transactionStatisticsFilter.sortingType = filter.sortingType;
            }
        },
        getTransactionListPageParams(item) {
            const querys = [];

            if (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
                querys.push('type=2');
            } else if (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                querys.push('type=3');
            }

            if (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('accountId=' + item.id);
            } else if (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                querys.push('categoryId=' + item.id);
            }

            if (this.transactionStatisticsFilter.chartDataType !== statisticsConstants.allChartDataTypes.AccountTotalAssets.type
                && this.transactionStatisticsFilter.chartDataType !== statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('dateType=' + this.transactionStatisticsFilter.dateType);

                if (this.transactionStatisticsFilter.dateType === datetimeConstants.allDateRanges.Custom.type) {
                    querys.push('minTime=' + this.transactionStatisticsFilter.startTime);
                    querys.push('maxTime=' + this.transactionStatisticsFilter.endTime);
                }
            }

            return querys.join('&');
        },
        loadTransactionStatistics({ force }) {
            const self = this;
            const settingsStore = useSettingsStore();

            return new Promise((resolve, reject) => {
                services.getTransactionStatistics({
                    startTime: self.transactionStatisticsFilter.startTime,
                    endTime: self.transactionStatisticsFilter.endTime,
                    useTransactionTimezone: settingsStore.appSettings.statistics.defaultTimezoneType
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve transaction statistics' });
                        return;
                    }

                    if (self.transactionStatisticsStateInvalid) {
                        self.updateTransactionStatisticsInvalidState(false);
                    }

                    if (force && data.result && isEquals(self.transactionStatisticsData, data.result)) {
                        reject({ message: 'Data is up to date' });
                        return;
                    }

                    self.transactionStatisticsData = data.result;

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to retrieve transaction statistics', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve transaction statistics' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
    }
});
