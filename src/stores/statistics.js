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
    isYearMonth,
    isObject
} from '@/lib/common.js';
import {
    getDateRangeByDateType
} from '@/lib/datetime.js';

export const useStatisticsStore = defineStore('statistics', {
    state: () => ({
        transactionStatisticsFilter: {
            chartDataType: statisticsConstants.defaultChartDataType,
            categoricalChartType: statisticsConstants.defaultCategoricalChartType,
            categoricalChartDateType: statisticsConstants.defaultCategoricalChartDataRangeType,
            categoricalChartStartTime: 0,
            categoricalChartEndTime: 0,
            trendChartType: statisticsConstants.defaultTrendChartType,
            trendChartDateType: statisticsConstants.defaultTrendChartDataRangeType,
            trendChartStartYearMonth: '',
            trendChartEndYearMonth: '',
            filterAccountIds: {},
            filterCategoryIds: {}
        },
        transactionCategoryStatisticsData: {},
        transactionStatisticsStateInvalid: true
    }),
    getters: {
        categoricalAnalysisChartDataCategory(state) {
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
        transactionCategoryStatisticsDataWithCategoryAndAccountInfo(state) {
            const statistics = state.transactionCategoryStatisticsData;
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
        transactionCategoryTotalAmountAnalysisData(state) {
            if (!state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo || !state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo.items) {
                return null;
            }

            const allDataItems = {};
            let totalAmount = 0;
            let totalNonNegativeAmount = 0;

            for (let i = 0; i < state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo.items.length; i++) {
                const item = state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo.items[i];

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
        accountTotalAmountAnalysisData(state) {
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
        categoricalAnalysisData(state) {
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
                combinedData = state.transactionCategoryTotalAmountAnalysisData;
            } else if (state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type ||
                state.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                combinedData = state.accountTotalAmountAnalysisData;
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
            this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            this.transactionStatisticsFilter.categoricalChartType = statisticsConstants.defaultCategoricalChartType;
            this.transactionStatisticsFilter.categoricalChartDateType = statisticsConstants.defaultCategoricalChartDataRangeType;
            this.transactionStatisticsFilter.categoricalChartStartTime = 0;
            this.transactionStatisticsFilter.categoricalChartEndTime = 0;
            this.transactionStatisticsFilter.trendChartType = statisticsConstants.defaultTrendChartType;
            this.transactionStatisticsFilter.trendChartDateType = statisticsConstants.defaultTrendChartDataRangeType;
            this.transactionStatisticsFilter.trendChartStartYearMonth = '';
            this.transactionStatisticsFilter.trendChartEndYearMonth = '';
            this.transactionStatisticsFilter.filterAccountIds = {};
            this.transactionStatisticsFilter.filterCategoryIds = {};
            this.transactionCategoryStatisticsData = {};
            this.transactionStatisticsStateInvalid = true;
        },
        initTransactionStatisticsFilter(filter) {
            if (!filter) {
                const settingsStore = useSettingsStore();
                const userStore = useUserStore();

                let defaultChartDataType = settingsStore.appSettings.statistics.defaultChartDataType;

                if (defaultChartDataType < statisticsConstants.allChartDataTypes.ExpenseByAccount.type || defaultChartDataType > statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                    defaultChartDataType = statisticsConstants.defaultChartDataType;
                }

                let defaultCategoricalChartType = settingsStore.appSettings.statistics.defaultCategoricalChartType;

                if (defaultCategoricalChartType !== statisticsConstants.allCategoricalChartTypes.Pie && defaultCategoricalChartType !== statisticsConstants.allCategoricalChartTypes.Bar) {
                    defaultCategoricalChartType = statisticsConstants.defaultCategoricalChartType;
                }

                let defaultCategoricalChartDateRange = settingsStore.appSettings.statistics.defaultCategoricalChartDataRangeType;

                if (defaultCategoricalChartDateRange < datetimeConstants.allDateRanges.All.type || defaultCategoricalChartDateRange >= datetimeConstants.allDateRanges.Custom.type) {
                    defaultCategoricalChartDateRange = statisticsConstants.defaultCategoricalChartDataRangeType;
                }

                let defaultTrendChartType = settingsStore.appSettings.statistics.defaultTrendChartType;

                if (defaultTrendChartType !== statisticsConstants.allTrendChartTypes.Area && defaultTrendChartType !== statisticsConstants.allTrendChartTypes.Column) {
                    defaultTrendChartType = statisticsConstants.defaultTrendChartType;
                }

                let defaultTrendChartDateRange = settingsStore.appSettings.statistics.defaultTrendChartDataRangeType;

                if (defaultTrendChartDateRange < datetimeConstants.allDateRanges.All.type || defaultTrendChartDateRange >= datetimeConstants.allDateRanges.Custom.type) {
                    defaultTrendChartDateRange = statisticsConstants.defaultTrendChartDataRangeType;
                }

                let defaultSortType = settingsStore.appSettings.statistics.defaultSortingType;

                if (defaultSortType < statisticsConstants.allSortingTypes.Amount.type || defaultSortType > statisticsConstants.allSortingTypes.Name.type) {
                    defaultSortType = statisticsConstants.defaultSortingType;
                }

                const categoricalChartDateRange = getDateRangeByDateType(defaultCategoricalChartDateRange, userStore.currentUserFirstDayOfWeek);
                const trendChartDateRange = getDateRangeByDateType(defaultTrendChartDateRange, userStore.currentUserFirstDayOfWeek);

                filter = {
                    chartDataType: defaultChartDataType,
                    categoricalChartType: defaultCategoricalChartType,
                    categoricalChartDateType: categoricalChartDateRange ? categoricalChartDateRange.dateType : undefined,
                    categoricalChartStartTime: categoricalChartDateRange ? categoricalChartDateRange.minTime : undefined,
                    categoricalChartEndTime: categoricalChartDateRange ? categoricalChartDateRange.maxTime : undefined,
                    trendChartType: defaultTrendChartType,
                    trendChartDateType: trendChartDateRange ? trendChartDateRange.dateType : undefined,
                    trendChartStartYearMonth: trendChartDateRange ? trendChartDateRange.minTime : undefined,
                    trendChartEndYearMonth: trendChartDateRange ? trendChartDateRange.maxTime : undefined,
                    filterAccountIds: settingsStore.appSettings.statistics.defaultAccountFilter || {},
                    filterCategoryIds: settingsStore.appSettings.statistics.defaultTransactionCategoryFilter || {},
                    sortingType: defaultSortType,
                };
            }

            if (filter && isNumber(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            } else {
                this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            }

            if (filter && isNumber(filter.categoricalChartType)) {
                this.transactionStatisticsFilter.categoricalChartType = filter.categoricalChartType;
            } else {
                this.transactionStatisticsFilter.categoricalChartType = statisticsConstants.defaultCategoricalChartType;
            }

            if (filter && isNumber(filter.categoricalChartDateType)) {
                this.transactionStatisticsFilter.categoricalChartDateType = filter.categoricalChartDateType;
            } else {
                this.transactionStatisticsFilter.categoricalChartDateType = statisticsConstants.defaultCategoricalChartDataRangeType;
            }

            if (filter && isNumber(filter.categoricalChartStartTime)) {
                this.transactionStatisticsFilter.categoricalChartStartTime = filter.categoricalChartStartTime;
            } else {
                this.transactionStatisticsFilter.categoricalChartStartTime = 0;
            }

            if (filter && isNumber(filter.categoricalChartEndTime)) {
                this.transactionStatisticsFilter.categoricalChartEndTime = filter.categoricalChartEndTime;
            } else {
                this.transactionStatisticsFilter.categoricalChartEndTime = 0;
            }

            if (filter && isNumber(filter.trendChartType)) {
                this.transactionStatisticsFilter.trendChartType = filter.trendChartType;
            } else {
                this.transactionStatisticsFilter.trendChartType = statisticsConstants.defaultTrendChartType;
            }

            if (filter && isNumber(filter.trendChartDateType)) {
                this.transactionStatisticsFilter.trendChartDateType = filter.trendChartDateType;
            } else {
                this.transactionStatisticsFilter.trendChartDateType = statisticsConstants.defaultTrendChartDataRangeType;
            }

            if (filter && isYearMonth(filter.trendChartStartYearMonth)) {
                this.transactionStatisticsFilter.trendChartStartYearMonth = filter.trendChartStartYearMonth;
            } else {
                this.transactionStatisticsFilter.trendChartStartYearMonth = '';
            }

            if (filter && isYearMonth(filter.trendChartEndYearMonth)) {
                this.transactionStatisticsFilter.trendChartEndYearMonth = filter.trendChartEndYearMonth;
            } else {
                this.transactionStatisticsFilter.trendChartEndYearMonth = '';
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
            if (filter && isNumber(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            }

            if (filter && isNumber(filter.categoricalChartType)) {
                this.transactionStatisticsFilter.categoricalChartType = filter.categoricalChartType;
            }

            if (filter && isNumber(filter.categoricalChartDateType)) {
                this.transactionStatisticsFilter.categoricalChartDateType = filter.categoricalChartDateType;
            }

            if (filter && isNumber(filter.categoricalChartStartTime)) {
                this.transactionStatisticsFilter.categoricalChartStartTime = filter.categoricalChartStartTime;
            }

            if (filter && isNumber(filter.categoricalChartEndTime)) {
                this.transactionStatisticsFilter.categoricalChartEndTime = filter.categoricalChartEndTime;
            }

            if (filter && isNumber(filter.trendChartType)) {
                this.transactionStatisticsFilter.trendChartType = filter.trendChartType;
            }

            if (filter && isNumber(filter.trendChartDateType)) {
                this.transactionStatisticsFilter.trendChartDateType = filter.trendChartDateType;
            }

            if (filter && isYearMonth(filter.trendChartStartYearMonth)) {
                this.transactionStatisticsFilter.trendChartStartYearMonth = filter.trendChartStartYearMonth;
            }

            if (filter && isYearMonth(filter.trendChartEndYearMonth)) {
                this.transactionStatisticsFilter.trendChartEndYearMonth = filter.trendChartEndYearMonth;
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
                querys.push('dateType=' + this.transactionStatisticsFilter.categoricalChartDateType);

                if (this.transactionStatisticsFilter.dateType === datetimeConstants.allDateRanges.Custom.type) {
                    querys.push('minTime=' + this.transactionStatisticsFilter.categoricalChartStartTime);
                    querys.push('maxTime=' + this.transactionStatisticsFilter.categoricalChartEndTime);
                }
            }

            return querys.join('&');
        },
        loadCategoricalAnalysis({ force }) {
            const self = this;
            const settingsStore = useSettingsStore();

            return new Promise((resolve, reject) => {
                services.getTransactionStatistics({
                    startTime: self.transactionStatisticsFilter.categoricalChartStartTime,
                    endTime: self.transactionStatisticsFilter.categoricalChartEndTime,
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

                    if (force && data.result && isEquals(self.transactionCategoryStatisticsData, data.result)) {
                        reject({ message: 'Data is up to date' });
                        return;
                    }

                    self.transactionCategoryStatisticsData = data.result;

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
        loadTrendAnalysis() {
            return Promise.resolve(true);
        },
    }
});
