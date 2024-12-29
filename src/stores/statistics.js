import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.js';
import { useUserStore } from './user.js';
import { useAccountsStore } from './account.js';
import { useTransactionCategoriesStore } from './transactionCategory.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import { CategoryType } from '@/core/category.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';
import datetimeConstants from '@/consts/datetime.js';
import statisticsConstants from '@/consts/statistics.js';
import { DEFAULT_ACCOUNT_ICON, DEFAULT_CATEGORY_ICON } from '@/consts/icon.ts';
import { DEFAULT_ACCOUNT_COLOR, DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import {
    isEquals,
    isNumber,
    isString,
    isObject,
    isInteger,
    isYearMonth,
    isYearMonthEquals,
    isObjectEmpty,
    objectFieldToArrayItem
} from '@/lib/common.ts';
import {
    getYearAndMonthFromUnixTime,
    getDateRangeByDateType
} from '@/lib/datetime.js';
import {
    getFinalAccountIdsByFilteredAccountIds
} from '@/lib/account.js';
import {
    getFinalCategoryIdsByFilteredCategoryIds
} from '@/lib/category.js';
import {
    sortStatisticsItems
} from '@/lib/statistics.js';

function assembleAccountAndCategoryInfo(userStore, accountsStore, transactionCategoriesStore, exchangeRatesStore, items) {
    const finalItems = [];
    const defaultCurrency = userStore.currentUserDefaultCurrency;

    for (let i = 0; i < items.length; i++) {
        const dataItem = items[i];
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

        finalItems.push(item);
    }

    return finalItems;
}

function getCategoryTotalAmountItems(items, transactionStatisticsFilter) {
    const allDataItems = {};
    let totalAmount = 0;
    let totalNonNegativeAmount = 0;

    for (let i = 0; i < items.length; i++) {
        const item = items[i];

        if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category) {
            continue;
        }

        if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalExpense.type) {
            if (item.category.type !== CategoryType.Expense) {
                continue;
            }
        } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalIncome.type) {
            if (item.category.type !== CategoryType.Income) {
                continue;
            }
        } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalBalance.type) {
            // Do Nothing
        } else {
            continue;
        }

        if (transactionStatisticsFilter.filterAccountIds && transactionStatisticsFilter.filterAccountIds[item.account.id]) {
            continue;
        }

        if (transactionStatisticsFilter.filterCategoryIds && transactionStatisticsFilter.filterCategoryIds[item.category.id]) {
            continue;
        }

        if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type) {
            if (isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems[item.account.id];

                if (data) {
                    data.totalAmount += item.amountInDefaultCurrency;
                } else {
                    data = {
                        name: item.account.name,
                        type: 'account',
                        id: item.account.id,
                        icon: item.account.icon || DEFAULT_ACCOUNT_ICON.icon,
                        color: item.account.color || DEFAULT_ACCOUNT_COLOR,
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
        } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type) {
            if (isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems[item.primaryCategory.id];

                if (data) {
                    data.totalAmount += item.amountInDefaultCurrency;
                } else {
                    data = {
                        name: item.primaryCategory.name,
                        type: 'category',
                        id: item.primaryCategory.id,
                        icon: item.primaryCategory.icon || DEFAULT_CATEGORY_ICON.icon,
                        color: item.primaryCategory.color || DEFAULT_CATEGORY_COLOR,
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
        } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type) {
            if (isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems[item.category.id];

                if (data) {
                    data.totalAmount += item.amountInDefaultCurrency;
                } else {
                    data = {
                        name: item.category.name,
                        type: 'category',
                        id: item.category.id,
                        icon: item.category.icon || DEFAULT_CATEGORY_ICON.icon,
                        color: item.category.color || DEFAULT_CATEGORY_COLOR,
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
        } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalExpense.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalIncome.type ||
            transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalBalance.type) {
            if (isNumber(item.amountInDefaultCurrency)) {
                let data = allDataItems['total'];
                let amount = item.amountInDefaultCurrency;

                if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalBalance.type &&
                    item.category.type === CategoryType.Expense) {
                    amount = -amount;
                }

                if (data) {
                    data.totalAmount += amount;
                } else {
                    let name = '';

                    if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalExpense.type) {
                        name = statisticsConstants.allChartDataTypes.TotalExpense.name;
                    } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalIncome.type) {
                        name = statisticsConstants.allChartDataTypes.TotalIncome.name;
                    } else if (transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalBalance.type) {
                        name = statisticsConstants.allChartDataTypes.TotalBalance.name;
                    }

                    data = {
                        name: name,
                        type: 'total',
                        id: 'total',
                        icon: '',
                        color: '',
                        hidden: false,
                        displayOrders: [1],
                        totalAmount: amount
                    }
                }

                totalAmount += amount;

                if (item.amountInDefaultCurrency > 0) {
                    totalNonNegativeAmount += amount;
                }

                allDataItems['total'] = data;
            }
        }
    }

    return {
        totalAmount: totalAmount,
        totalNonNegativeAmount: totalNonNegativeAmount,
        items: allDataItems
    };
}

function sortCategoryTotalAmountItems(items, transactionStatisticsFilter) {
    sortStatisticsItems(items, transactionStatisticsFilter.sortingType)
}

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
            filterCategoryIds: {},
            tagIds: '',
            tagFilterType: TransactionTagFilterType.Default.type
        },
        transactionCategoryStatisticsData: {},
        transactionCategoryTrendsData: {},
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

                finalStatistics.items = assembleAccountAndCategoryInfo(userStore, accountsStore, transactionCategoriesStore, exchangeRatesStore, statistics.items);
            }

            return finalStatistics;
        },
        transactionCategoryTotalAmountAnalysisData(state) {
            if (!state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo || !state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo.items) {
                return null;
            }

            return getCategoryTotalAmountItems(state.transactionCategoryStatisticsDataWithCategoryAndAccountInfo.items, state.transactionStatisticsFilter);
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
                    icon: account.icon || DEFAULT_ACCOUNT_ICON.icon,
                    color: account.color || DEFAULT_ACCOUNT_COLOR,
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

            sortCategoryTotalAmountItems(allStatisticsItems, state.transactionStatisticsFilter);

            return {
                totalAmount: combinedData.totalAmount,
                items: allStatisticsItems
            };
        },
        transactionCategoryTrendsDataWithCategoryAndAccountInfo(state) {
            const trendsData = state.transactionCategoryTrendsData;
            const finalTrendsData = [];

            if (trendsData && trendsData.length) {
                const userStore = useUserStore();
                const accountsStore = useAccountsStore();
                const transactionCategoriesStore = useTransactionCategoriesStore();
                const exchangeRatesStore = useExchangeRatesStore();

                for (let i = 0; i < trendsData.length; i++) {
                    const trendItem = trendsData[i];
                    const finalTrendItem = {
                        year: trendItem.year,
                        month: trendItem.month,
                        items: []
                    };

                    if (trendItem && trendItem.items && trendItem.items.length) {
                        finalTrendItem.items = assembleAccountAndCategoryInfo(userStore, accountsStore, transactionCategoriesStore, exchangeRatesStore, trendItem.items);
                    }

                    finalTrendsData.push(finalTrendItem);
                }
            }

            return finalTrendsData;
        },
        trendsAnalysisData(state) {
            if (!state.transactionCategoryTrendsDataWithCategoryAndAccountInfo || !state.transactionCategoryTrendsDataWithCategoryAndAccountInfo.length) {
                return null;
            }

            const combinedDataMap = {};

            for (let i = 0; i < state.transactionCategoryTrendsDataWithCategoryAndAccountInfo.length; i++) {
                const trendItem = state.transactionCategoryTrendsDataWithCategoryAndAccountInfo[i];
                const totalAmountItems = getCategoryTotalAmountItems(trendItem.items, state.transactionStatisticsFilter);

                for (let id in totalAmountItems.items) {
                    if (!Object.prototype.hasOwnProperty.call(totalAmountItems.items, id)) {
                        continue;
                    }

                    const item = totalAmountItems.items[id];
                    let combinedData = combinedDataMap[id];

                    if (!combinedData) {
                        combinedData = {
                            name: item.name,
                            type: item.type,
                            id: item.id,
                            icon: item.icon,
                            color: item.color,
                            hidden: item.hidden,
                            displayOrders: item.displayOrders,
                            totalAmount: 0,
                            items: []
                        };
                    }

                    combinedData.items.push({
                        year: trendItem.year,
                        month: trendItem.month,
                        totalAmount: item.totalAmount
                    });

                    combinedData.totalAmount += item.totalAmount;
                    combinedDataMap[id] = combinedData;
                }
            }

            const totalAmountsTrends = [];

            for (let id in combinedDataMap) {
                if (!Object.prototype.hasOwnProperty.call(combinedDataMap, id)) {
                    continue;
                }

                const trendData = combinedDataMap[id];
                totalAmountsTrends.push(trendData);
            }

            sortCategoryTotalAmountItems(totalAmountsTrends, state.transactionStatisticsFilter);

            return {
                items: totalAmountsTrends
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
            this.transactionStatisticsFilter.tagIds = '';
            this.transactionStatisticsFilter.tagFilterType = TransactionTagFilterType.Default.type;
            this.transactionCategoryStatisticsData = {};
            this.transactionCategoryTrendsData = {};
            this.transactionStatisticsStateInvalid = true;
        },
        initTransactionStatisticsFilter(analysisType, filter) {
            const settingsStore = useSettingsStore();
            const userStore = useUserStore();

            if (filter && isInteger(filter.chartDataType)) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            } else {
                this.transactionStatisticsFilter.chartDataType = settingsStore.appSettings.statistics.defaultChartDataType;
            }

            if (analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis || analysisType === statisticsConstants.allAnalysisTypes.TrendAnalysis) {
                if (!statisticsConstants.allChartDataTypesMap[this.transactionStatisticsFilter.chartDataType] ||
                    !statisticsConstants.allChartDataTypesMap[this.transactionStatisticsFilter.chartDataType].availableAnalysisTypes[analysisType]) {
                    this.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
                }
            }

            if (filter && isInteger(filter.categoricalChartType)) {
                this.transactionStatisticsFilter.categoricalChartType = filter.categoricalChartType;
            } else {
                this.transactionStatisticsFilter.categoricalChartType = settingsStore.appSettings.statistics.defaultCategoricalChartType;
            }

            if (this.transactionStatisticsFilter.categoricalChartType !== statisticsConstants.allCategoricalChartTypes.Pie && this.transactionStatisticsFilter.categoricalChartType !== statisticsConstants.allCategoricalChartTypes.Bar) {
                this.transactionStatisticsFilter.categoricalChartType = statisticsConstants.defaultCategoricalChartType;
            }

            if (filter && isInteger(filter.categoricalChartDateType)) {
                this.transactionStatisticsFilter.categoricalChartDateType = filter.categoricalChartDateType;
            } else {
                this.transactionStatisticsFilter.categoricalChartDateType = settingsStore.appSettings.statistics.defaultCategoricalChartDataRangeType;
            }

            let categoricalChartDateTypeValid = true;

            if (!datetimeConstants.allDateRangesMap[this.transactionStatisticsFilter.categoricalChartDateType] ||
                !datetimeConstants.allDateRangesMap[this.transactionStatisticsFilter.categoricalChartDateType].availableScenes[datetimeConstants.allDateRangeScenes.Normal]) {
                this.transactionStatisticsFilter.categoricalChartDateType = statisticsConstants.defaultCategoricalChartDataRangeType;
                categoricalChartDateTypeValid = false;
            }

            if (categoricalChartDateTypeValid && this.transactionStatisticsFilter.categoricalChartDateType === datetimeConstants.allDateRanges.Custom.type) {
                if (filter && isInteger(filter.categoricalChartStartTime)) {
                    this.transactionStatisticsFilter.categoricalChartStartTime = filter.categoricalChartStartTime;
                } else {
                    this.transactionStatisticsFilter.categoricalChartStartTime = 0;
                }

                if (filter && isInteger(filter.categoricalChartEndTime)) {
                    this.transactionStatisticsFilter.categoricalChartEndTime = filter.categoricalChartEndTime;
                } else {
                    this.transactionStatisticsFilter.categoricalChartEndTime = 0;
                }
            } else {
                const categoricalChartDateRange = getDateRangeByDateType(this.transactionStatisticsFilter.categoricalChartDateType, userStore.currentUserFirstDayOfWeek);
                this.transactionStatisticsFilter.categoricalChartDateType = categoricalChartDateRange.dateType;
                this.transactionStatisticsFilter.categoricalChartStartTime = categoricalChartDateRange.minTime;
                this.transactionStatisticsFilter.categoricalChartEndTime = categoricalChartDateRange.maxTime;
            }

            if (filter && isInteger(filter.trendChartType)) {
                this.transactionStatisticsFilter.trendChartType = filter.trendChartType;
            } else {
                this.transactionStatisticsFilter.trendChartType = settingsStore.appSettings.statistics.defaultTrendChartType;
            }

            if (this.transactionStatisticsFilter.trendChartType !== statisticsConstants.allTrendChartTypes.Area && this.transactionStatisticsFilter.trendChartType !== statisticsConstants.allTrendChartTypes.Column) {
                this.transactionStatisticsFilter.trendChartType = statisticsConstants.defaultTrendChartType;
            }

            if (filter && isInteger(filter.trendChartDateType)) {
                this.transactionStatisticsFilter.trendChartDateType = filter.trendChartDateType;
            } else {
                this.transactionStatisticsFilter.trendChartDateType = settingsStore.appSettings.statistics.defaultTrendChartDataRangeType;
            }

            let trendChartDateTypeValid = true;

            if (!datetimeConstants.allDateRangesMap[this.transactionStatisticsFilter.trendChartDateType] ||
                !datetimeConstants.allDateRangesMap[this.transactionStatisticsFilter.trendChartDateType].availableScenes[datetimeConstants.allDateRangeScenes.TrendAnalysis]) {
                this.transactionStatisticsFilter.trendChartDateType = statisticsConstants.defaultTrendChartDataRangeType;
                trendChartDateTypeValid = false;
            }

            if (trendChartDateTypeValid && this.transactionStatisticsFilter.trendChartDateType === datetimeConstants.allDateRanges.Custom.type) {
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
            } else {
                const trendChartDateRange = getDateRangeByDateType(this.transactionStatisticsFilter.trendChartDateType, userStore.currentUserFirstDayOfWeek);
                this.transactionStatisticsFilter.trendChartDateType = trendChartDateRange.dateType;
                this.transactionStatisticsFilter.trendChartStartYearMonth = getYearAndMonthFromUnixTime(trendChartDateRange.minTime);
                this.transactionStatisticsFilter.trendChartEndYearMonth = getYearAndMonthFromUnixTime(trendChartDateRange.maxTime);
            }

            if (filter && isObject(filter.filterAccountIds)) {
                this.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            } else {
                this.transactionStatisticsFilter.filterAccountIds = settingsStore.appSettings.statistics.defaultAccountFilter || {};
            }

            if (filter && isObject(filter.filterCategoryIds)) {
                this.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            } else {
                this.transactionStatisticsFilter.filterCategoryIds = settingsStore.appSettings.statistics.defaultTransactionCategoryFilter || {};
            }

            if (filter && isString(filter.tagIds)) {
                this.transactionStatisticsFilter.tagIds = filter.tagIds;
            } else {
                this.transactionStatisticsFilter.tagIds = '';
            }

            if (filter && isInteger(filter.tagFilterType)) {
                this.transactionStatisticsFilter.tagFilterType = filter.tagFilterType;
            } else {
                this.transactionStatisticsFilter.tagFilterType = TransactionTagFilterType.Default.type;
            }

            if (filter && isInteger(filter.sortingType)) {
                this.transactionStatisticsFilter.sortingType = filter.sortingType;
            } else {
                this.transactionStatisticsFilter.sortingType = settingsStore.appSettings.statistics.defaultSortingType;
            }

            if (this.transactionStatisticsFilter.sortingType < statisticsConstants.allSortingTypes.Amount.type || this.transactionStatisticsFilter.sortingType > statisticsConstants.allSortingTypes.Name.type) {
                this.transactionStatisticsFilter.sortingType = statisticsConstants.defaultSortingType;
            }
        },
        updateTransactionStatisticsFilter(filter) {
            let changed = false;

            if (filter && isInteger(filter.chartDataType) && this.transactionStatisticsFilter.chartDataType !== filter.chartDataType) {
                this.transactionStatisticsFilter.chartDataType = filter.chartDataType;
                changed = true;
            }

            if (filter && isInteger(filter.categoricalChartType) && this.transactionStatisticsFilter.categoricalChartType !== filter.categoricalChartType) {
                this.transactionStatisticsFilter.categoricalChartType = filter.categoricalChartType;
                changed = true;
            }

            if (filter && isInteger(filter.categoricalChartDateType) && this.transactionStatisticsFilter.categoricalChartDateType !== filter.categoricalChartDateType) {
                this.transactionStatisticsFilter.categoricalChartDateType = filter.categoricalChartDateType;
                changed = true;
            }

            if (filter && isInteger(filter.categoricalChartStartTime) && this.transactionStatisticsFilter.categoricalChartStartTime !== filter.categoricalChartStartTime) {
                this.transactionStatisticsFilter.categoricalChartStartTime = filter.categoricalChartStartTime;
                changed = true;
            }

            if (filter && isInteger(filter.categoricalChartEndTime) && this.transactionStatisticsFilter.categoricalChartEndTime !== filter.categoricalChartEndTime) {
                this.transactionStatisticsFilter.categoricalChartEndTime = filter.categoricalChartEndTime;
                changed = true;
            }

            if (filter && isInteger(filter.trendChartType) && this.transactionStatisticsFilter.trendChartType !== filter.trendChartType) {
                this.transactionStatisticsFilter.trendChartType = filter.trendChartType;
                changed = true;
            }

            if (filter && isInteger(filter.trendChartDateType) && this.transactionStatisticsFilter.trendChartDateType !== filter.trendChartDateType) {
                this.transactionStatisticsFilter.trendChartDateType = filter.trendChartDateType;
                changed = true;
            }

            if (filter && (isYearMonth(filter.trendChartStartYearMonth) || filter.trendChartStartYearMonth === '') && !isYearMonthEquals(this.transactionStatisticsFilter.trendChartStartYearMonth, filter.trendChartStartYearMonth)) {
                this.transactionStatisticsFilter.trendChartStartYearMonth = filter.trendChartStartYearMonth;
                changed = true;
            }

            if (filter && (isYearMonth(filter.trendChartEndYearMonth) || filter.trendChartEndYearMonth === '') && !isYearMonthEquals(this.transactionStatisticsFilter.trendChartEndYearMonth, filter.trendChartEndYearMonth)) {
                this.transactionStatisticsFilter.trendChartEndYearMonth = filter.trendChartEndYearMonth;
                changed = true;
            }

            if (filter && isObject(filter.filterAccountIds) && !isEquals(this.transactionStatisticsFilter.filterAccountIds, filter.filterAccountIds)) {
                this.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
                changed = true;
            }

            if (filter && isObject(filter.filterCategoryIds) && !isEquals(this.transactionStatisticsFilter.filterCategoryIds, filter.filterCategoryIds)) {
                this.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
                changed = true;
            }

            if (filter && isString(filter.tagIds) && this.transactionStatisticsFilter.tagIds !== filter.tagIds) {
                this.transactionStatisticsFilter.tagIds = filter.tagIds;
                changed = true;
            }

            if (filter && isInteger(filter.tagFilterType) && this.transactionStatisticsFilter.tagFilterType !== filter.tagFilterType) {
                this.transactionStatisticsFilter.tagFilterType = filter.tagFilterType;
                changed = true;
            }

            if (filter && isInteger(filter.sortingType) && this.transactionStatisticsFilter.sortingType !== filter.sortingType) {
                this.transactionStatisticsFilter.sortingType = filter.sortingType;
                changed = true;
            }

            return changed;
        },
        getTransactionStatisticsPageParams(analysisType, trendDateAggregationType) {
            const querys = [];

            querys.push('analysisType=' + analysisType);
            querys.push('chartDataType=' + this.transactionStatisticsFilter.chartDataType);

            if (analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                querys.push('chartType=' + this.transactionStatisticsFilter.categoricalChartType);
                querys.push('chartDateType=' + this.transactionStatisticsFilter.categoricalChartDateType);

                if (this.transactionStatisticsFilter.categoricalChartDateType === datetimeConstants.allDateRanges.Custom.type) {
                    querys.push('startTime=' + this.transactionStatisticsFilter.categoricalChartStartTime);
                    querys.push('endTime=' + this.transactionStatisticsFilter.categoricalChartEndTime);
                }
            } else if (analysisType === statisticsConstants.allAnalysisTypes.TrendAnalysis) {
                querys.push('chartType=' + this.transactionStatisticsFilter.trendChartType);
                querys.push('chartDateType=' + this.transactionStatisticsFilter.trendChartDateType);

                if (this.transactionStatisticsFilter.trendChartDateType === datetimeConstants.allDateRanges.Custom.type) {
                    querys.push('startTime=' + this.transactionStatisticsFilter.trendChartStartYearMonth);
                    querys.push('endTime=' + this.transactionStatisticsFilter.trendChartEndYearMonth);
                }

                if (trendDateAggregationType !== statisticsConstants.allDateAggregationTypes.Month.type) {
                    querys.push('trendDateAggregationType=' + trendDateAggregationType);
                }
            }

            if (this.transactionStatisticsFilter.filterAccountIds) {
                const ids = objectFieldToArrayItem(this.transactionStatisticsFilter.filterAccountIds);

                if (ids && ids.length) {
                    querys.push('filterAccountIds=' + ids.join(','));
                }
            }

            if (this.transactionStatisticsFilter.filterCategoryIds) {
                const ids = objectFieldToArrayItem(this.transactionStatisticsFilter.filterCategoryIds);

                if (ids && ids.length) {
                    querys.push('filterCategoryIds=' + ids.join(','));
                }
            }

            if (this.transactionStatisticsFilter.tagIds) {
                querys.push('tagIds=' + this.transactionStatisticsFilter.tagIds);
            }

            if (this.transactionStatisticsFilter.tagFilterType) {
                querys.push('tagFilterType=' + this.transactionStatisticsFilter.tagFilterType);
            }

            querys.push('sortingType=' + this.transactionStatisticsFilter.sortingType);

            return querys.join('&');
        },
        getTransactionListPageParams(analysisType, itemId, dateRange) {
            const accountsStore = useAccountsStore();
            const transactionCategoriesStore = useTransactionCategoriesStore();
            const querys = [];

            if (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalIncome.type) {
                querys.push('type=2');
            } else if (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.TotalExpense.type) {
                querys.push('type=3');
            }

            if (itemId && (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByAccount.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalAssets.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type)) {
                querys.push('accountIds=' + itemId);

                if (!isObjectEmpty(this.transactionStatisticsFilter.filterCategoryIds)) {
                    querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, this.transactionStatisticsFilter.filterCategoryIds));
                }
            } else if (itemId && (this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.IncomeBySecondaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.transactionStatisticsFilter.chartDataType === statisticsConstants.allChartDataTypes.ExpenseBySecondaryCategory.type)) {
                querys.push('categoryIds=' + itemId);

                if (!isObjectEmpty(this.transactionStatisticsFilter.filterAccountIds)) {
                    querys.push('accountIds=' + getFinalAccountIdsByFilteredAccountIds(accountsStore.allAccountsMap, this.transactionStatisticsFilter.filterAccountIds));
                }
            } else if (!itemId) {
                if (!isObjectEmpty(this.transactionStatisticsFilter.filterCategoryIds)) {
                    querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, this.transactionStatisticsFilter.filterCategoryIds));
                }

                if (!isObjectEmpty(this.transactionStatisticsFilter.filterAccountIds)) {
                    querys.push('accountIds=' + getFinalAccountIdsByFilteredAccountIds(accountsStore.allAccountsMap, this.transactionStatisticsFilter.filterAccountIds));
                }
            }

            if (this.transactionStatisticsFilter.tagIds) {
                querys.push('tagIds=' + this.transactionStatisticsFilter.tagIds);
            }

            if (this.transactionStatisticsFilter.tagFilterType) {
                querys.push('tagFilterType=' + this.transactionStatisticsFilter.tagFilterType);
            }

            if (analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis
                && this.transactionStatisticsFilter.chartDataType !== statisticsConstants.allChartDataTypes.AccountTotalAssets.type
                && this.transactionStatisticsFilter.chartDataType !== statisticsConstants.allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('dateType=' + this.transactionStatisticsFilter.categoricalChartDateType);

                if (this.transactionStatisticsFilter.categoricalChartDateType === datetimeConstants.allDateRanges.Custom.type) {
                    querys.push('minTime=' + this.transactionStatisticsFilter.categoricalChartStartTime);
                    querys.push('maxTime=' + this.transactionStatisticsFilter.categoricalChartEndTime);
                }
            } else if (analysisType === statisticsConstants.allAnalysisTypes.TrendAnalysis && dateRange) {
                querys.push('dateType=' + dateRange.type);
                querys.push('minTime=' + dateRange.minTime);
                querys.push('maxTime=' + dateRange.maxTime);
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
                    tagIds: self.transactionStatisticsFilter.tagIds,
                    tagFilterType: self.transactionStatisticsFilter.tagFilterType,
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
        loadTrendAnalysis({ force }) {
            const self = this;
            const settingsStore = useSettingsStore();

            return new Promise((resolve, reject) => {
                services.getTransactionStatisticsTrends({
                    startYearMonth: self.transactionStatisticsFilter.trendChartStartYearMonth,
                    endYearMonth: self.transactionStatisticsFilter.trendChartEndYearMonth,
                    tagIds: self.transactionStatisticsFilter.tagIds,
                    tagFilterType: self.transactionStatisticsFilter.tagFilterType,
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

                    if (force && data.result && isEquals(self.transactionCategoryTrendsData, data.result)) {
                        reject({ message: 'Data is up to date' });
                        return;
                    }

                    self.transactionCategoryTrendsData = data.result;

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
