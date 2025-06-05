import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { type TimeRangeAndDateType, DateRangeScene, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';
import {
    StatisticsAnalysisType,
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType,
    DEFAULT_CATEGORICAL_CHART_DATA_RANGE,
    DEFAULT_TREND_CHART_DATA_RANGE
} from '@/core/statistics.ts';
import { DEFAULT_ACCOUNT_ICON, DEFAULT_CATEGORY_ICON } from '@/consts/icon.ts';
import { DEFAULT_ACCOUNT_COLOR, DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type {
    TransactionStatisticResponse,
    TransactionStatisticResponseItem,
    TransactionStatisticTrendsResponseItem,
    TransactionStatisticDataItemType,
    TransactionStatisticDataItemBase,
    TransactionCategoricalAnalysisData,
    TransactionCategoricalAnalysisDataItem,
    TransactionTrendsAnalysisData,
    TransactionTrendsAnalysisDataItem,
    TransactionTrendsAnalysisDataAmount
} from '@/models/transaction.ts';

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
import { getYearAndMonthFromUnixTime, getDateRangeByDateType } from '@/lib/datetime.ts';
import { getFinalAccountIdsByFilteredAccountIds } from '@/lib/account.ts';
import { getFinalCategoryIdsByFilteredCategoryIds } from '@/lib/category.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';
import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

interface TransactionStatisticResponseItemWithInfo extends TransactionStatisticResponseItem {
    categoryId: string;
    accountId: string;
    amount: number;
    account?: Account;
    primaryAccount?: Account;
    category?: TransactionCategory;
    primaryCategory?: TransactionCategory;
    amountInDefaultCurrency: number | null;
}

interface TransactionStatisticResponseWithInfo {
    readonly startTime: number;
    readonly endTime: number;
    readonly items: TransactionStatisticResponseItemWithInfo[];
}

interface TransactionStatisticTrendsResponseItemWithInfo {
    readonly year: number;
    readonly month: number;
    readonly items: TransactionStatisticResponseItemWithInfo[];
}

interface WritableTransactionCategoricalAnalysisData {
    totalAmount: number;
    totalNonNegativeAmount: number;
    items: Record<string, WritableTransactionCategoricalAnalysisDataItem>;
}

interface WritableTransactionCategoricalAnalysisDataItem {
    name: string;
    type: TransactionStatisticDataItemType;
    id: string;
    icon: string;
    color: string;
    hidden: boolean;
    displayOrders: number[];
    totalAmount: number;
    percent?: number;
}

interface WritableTransactionTrendsAnalysisDataItem {
    name: string;
    type: TransactionStatisticDataItemType;
    id: string;
    icon: string;
    color: string;
    hidden: boolean;
    displayOrders: number[];
    totalAmount: number;
    items: TransactionTrendsAnalysisDataAmount[];
}

export interface TransactionStatisticsPartialFilter {
    chartDataType?: number;
    categoricalChartType?: number;
    categoricalChartDateType?: number;
    categoricalChartStartTime?: number;
    categoricalChartEndTime?: number;
    trendChartType?: number;
    trendChartDateType?: number;
    trendChartStartYearMonth?: string;
    trendChartEndYearMonth?: string;
    filterAccountIds?: Record<string, boolean>;
    filterCategoryIds?: Record<string, boolean>;
    tagIds?: string;
    tagFilterType?: number;
    sortingType?: number;
}

export interface TransactionStatisticsFilter extends TransactionStatisticsPartialFilter {
    chartDataType: number;
    categoricalChartType: number;
    categoricalChartDateType: number;
    categoricalChartStartTime: number;
    categoricalChartEndTime: number;
    trendChartType: number;
    trendChartDateType: number;
    trendChartStartYearMonth: string;
    trendChartEndYearMonth: string;
    filterAccountIds: Record<string, boolean>;
    filterCategoryIds: Record<string, boolean>;
    tagIds: string;
    tagFilterType: number;
    sortingType: number;
}

export const useStatisticsStore = defineStore('statistics', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const transactionStatisticsFilter = ref<TransactionStatisticsFilter>({
        chartDataType: ChartDataType.Default.type,
        categoricalChartType: CategoricalChartType.Default.type,
        categoricalChartDateType: DEFAULT_CATEGORICAL_CHART_DATA_RANGE.type,
        categoricalChartStartTime: 0,
        categoricalChartEndTime: 0,
        trendChartType: TrendChartType.Default.type,
        trendChartDateType: DEFAULT_TREND_CHART_DATA_RANGE.type,
        trendChartStartYearMonth: '',
        trendChartEndYearMonth: '',
        filterAccountIds: {},
        filterCategoryIds: {},
        tagIds: '',
        tagFilterType: TransactionTagFilterType.Default.type,
        sortingType: ChartSortingType.Default.type
    });

    const transactionCategoryStatisticsData = ref<TransactionStatisticResponse | null>(null);
    const transactionCategoryTrendsData = ref<TransactionStatisticTrendsResponseItem[]>([]);
    const transactionStatisticsStateInvalid = ref<boolean>(true);

    const categoricalAnalysisChartDataCategory = computed<string>(() => {
        if (transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
            return 'account';
        } else if (transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
            return 'category';
        } else {
            return '';
        }
    });

    const transactionCategoryStatisticsDataWithCategoryAndAccountInfo = computed<TransactionStatisticResponseWithInfo | null>(() => {
        const statistics = transactionCategoryStatisticsData.value;

        if (!statistics) {
            return null;
        }

        const finalStatistics: TransactionStatisticResponseWithInfo = {
            startTime: statistics.startTime,
            endTime: statistics.endTime,
            items: []
        };

        if (statistics && statistics.items && statistics.items.length) {
            finalStatistics.items.push(...assembleAccountAndCategoryInfo(statistics.items));
        }

        return finalStatistics;
    });

    const transactionCategoryTotalAmountAnalysisData = computed<WritableTransactionCategoricalAnalysisData | null>(() => {
        if (!transactionCategoryStatisticsDataWithCategoryAndAccountInfo.value || !transactionCategoryStatisticsDataWithCategoryAndAccountInfo.value.items) {
            return null;
        }

        return getCategoryTotalAmountItems(transactionCategoryStatisticsDataWithCategoryAndAccountInfo.value.items, transactionStatisticsFilter.value);
    });

    const accountTotalAmountAnalysisData = computed<WritableTransactionCategoricalAnalysisData | null>(() => {
        if (!accountsStore.allPlainAccounts) {
            return null;
        }

        const allDataItems: Record<string, WritableTransactionCategoricalAnalysisDataItem> = {};
        let totalAmount = 0;
        let totalNonNegativeAmount = 0;

        for (let i = 0; i < accountsStore.allPlainAccounts.length; i++) {
            const account = accountsStore.allPlainAccounts[i];

            if (transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type) {
                if (!account.isAsset) {
                    continue;
                }
            } else if (transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                if (!account.isLiability) {
                    continue;
                }
            }

            if (transactionStatisticsFilter.value.filterAccountIds && transactionStatisticsFilter.value.filterAccountIds[account.id]) {
                continue;
            }

            let primaryAccount = accountsStore.allAccountsMap[account.parentId];

            if (!primaryAccount) {
                primaryAccount = account;
            }

            let amount = account.balance;

            if (account.currency !== userStore.currentUserDefaultCurrency) {
                const finalAmount = exchangeRatesStore.getExchangedAmount(amount, account.currency, userStore.currentUserDefaultCurrency);

                if (!isNumber(finalAmount)) {
                    continue;
                }

                amount = Math.floor(finalAmount);
            }

            if (account.isLiability) {
                amount = -amount;
            }

            const data: WritableTransactionCategoricalAnalysisDataItem = {
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
        };
    });

    const categoricalAnalysisData = computed<TransactionCategoricalAnalysisData>(() => {
        let combinedData: WritableTransactionCategoricalAnalysisData | null = null;

        if (transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
            combinedData = transactionCategoryTotalAmountAnalysisData.value;
        } else if (transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
            combinedData = accountTotalAmountAnalysisData.value;
        }

        const allStatisticsItems: TransactionCategoricalAnalysisDataItem[] = [];

        if (combinedData && combinedData.items) {
            for (const id in combinedData.items) {
                if (!Object.prototype.hasOwnProperty.call(combinedData.items, id)) {
                    continue;
                }

                const dataItem = combinedData.items[id];
                let percent = 0;

                if (dataItem.totalAmount > 0) {
                    percent = dataItem.totalAmount * 100 / combinedData.totalNonNegativeAmount;
                } else {
                    percent = 0;
                }

                if (percent < 0) {
                    percent = 0;
                }

                const statisticDataItem: TransactionCategoricalAnalysisDataItem = {
                    name: dataItem.name,
                    type: dataItem.type,
                    id: dataItem.id,
                    icon: dataItem.icon,
                    color: dataItem.color,
                    hidden: dataItem.hidden,
                    displayOrders: dataItem.displayOrders,
                    totalAmount: dataItem.totalAmount,
                    percent: percent
                };

                allStatisticsItems.push(statisticDataItem);
            }
        }

        sortCategoryTotalAmountItems(allStatisticsItems, transactionStatisticsFilter.value);

        const statisticData: TransactionCategoricalAnalysisData = {
            totalAmount: combinedData?.totalAmount || 0,
            items: allStatisticsItems
        };

        return statisticData;
    });

    const transactionCategoryTrendsDataWithCategoryAndAccountInfo = computed<TransactionStatisticTrendsResponseItemWithInfo[]>(() => {
        const trendsData = transactionCategoryTrendsData.value;
        const finalTrendsData: TransactionStatisticTrendsResponseItemWithInfo[] = [];

        if (trendsData && trendsData.length) {
            for (let i = 0; i < trendsData.length; i++) {
                const trendItem = trendsData[i];
                const finalTrendItem: TransactionStatisticTrendsResponseItemWithInfo = {
                    year: trendItem.year,
                    month: trendItem.month,
                    items: []
                };

                if (trendItem && trendItem.items && trendItem.items.length) {
                    finalTrendItem.items.push(...assembleAccountAndCategoryInfo(trendItem.items));
                }

                finalTrendsData.push(finalTrendItem);
            }
        }

        return finalTrendsData;
    });

    const trendsAnalysisData = computed<TransactionTrendsAnalysisData | null>(() => {
        if (!transactionCategoryTrendsDataWithCategoryAndAccountInfo.value || !transactionCategoryTrendsDataWithCategoryAndAccountInfo.value.length) {
            return null;
        }

        const combinedDataMap: Record<string, WritableTransactionTrendsAnalysisDataItem> = {};

        for (let i = 0; i < transactionCategoryTrendsDataWithCategoryAndAccountInfo.value.length; i++) {
            const trendItem = transactionCategoryTrendsDataWithCategoryAndAccountInfo.value[i];
            const totalAmountItems = getCategoryTotalAmountItems(trendItem.items, transactionStatisticsFilter.value);

            for (const id in totalAmountItems.items) {
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

        const totalAmountsTrends: TransactionTrendsAnalysisDataItem[] = [];

        for (const id in combinedDataMap) {
            if (!Object.prototype.hasOwnProperty.call(combinedDataMap, id)) {
                continue;
            }

            const trendData = combinedDataMap[id];
            totalAmountsTrends.push(trendData);
        }

        sortCategoryTotalAmountItems(totalAmountsTrends, transactionStatisticsFilter.value);

        const trendsData: TransactionTrendsAnalysisData = {
            items: totalAmountsTrends
        };

        return trendsData;
    });

    function assembleAccountAndCategoryInfo(items: TransactionStatisticResponseItem[]): TransactionStatisticResponseItemWithInfo[] {
        const finalItems: TransactionStatisticResponseItemWithInfo[] = [];
        const defaultCurrency = userStore.currentUserDefaultCurrency;

        for (let i = 0; i < items.length; i++) {
            const dataItem = items[i];
            const item: TransactionStatisticResponseItemWithInfo = {
                categoryId: dataItem.categoryId,
                accountId: dataItem.accountId,
                amount: dataItem.amount,
                amountInDefaultCurrency: null
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

    function getCategoryTotalAmountItems(items: TransactionStatisticResponseItemWithInfo[], transactionStatisticsFilter: TransactionStatisticsFilter): WritableTransactionCategoricalAnalysisData {
        const allDataItems: Record<string, WritableTransactionCategoricalAnalysisDataItem> = {};
        let totalAmount = 0;
        let totalNonNegativeAmount = 0;

        for (let i = 0; i < items.length; i++) {
            const item = items[i];

            if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category) {
                continue;
            }

            if (transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalExpense.type) {
                if (item.category.type !== CategoryType.Expense) {
                    continue;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalIncome.type) {
                if (item.category.type !== CategoryType.Income) {
                    continue;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalBalance.type) {
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

            if (transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByAccount.type) {
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
                        };
                    }

                    totalAmount += item.amountInDefaultCurrency;

                    if (item.amountInDefaultCurrency > 0) {
                        totalNonNegativeAmount += item.amountInDefaultCurrency;
                    }

                    allDataItems[item.account.id] = data;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByPrimaryCategory.type) {
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
                        };
                    }

                    totalAmount += item.amountInDefaultCurrency;

                    if (item.amountInDefaultCurrency > 0) {
                        totalNonNegativeAmount += item.amountInDefaultCurrency;
                    }

                    allDataItems[item.primaryCategory.id] = data;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
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
                        };
                    }

                    totalAmount += item.amountInDefaultCurrency;

                    if (item.amountInDefaultCurrency > 0) {
                        totalNonNegativeAmount += item.amountInDefaultCurrency;
                    }

                    allDataItems[item.category.id] = data;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalExpense.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalIncome.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalBalance.type) {
                if (isNumber(item.amountInDefaultCurrency)) {
                    let data = allDataItems['total'];
                    let amount = item.amountInDefaultCurrency;

                    if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalBalance.type &&
                        item.category.type === CategoryType.Expense) {
                        amount = -amount;
                    }

                    if (data) {
                        data.totalAmount += amount;
                    } else {
                        let name = '';

                        if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalExpense.type) {
                            name = ChartDataType.TotalExpense.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalIncome.type) {
                            name = ChartDataType.TotalIncome.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalBalance.type) {
                            name = ChartDataType.TotalBalance.name;
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
                        };
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

    function sortCategoryTotalAmountItems(items: TransactionStatisticDataItemBase[], transactionStatisticsFilter: TransactionStatisticsFilter): void {
        sortStatisticsItems(items, transactionStatisticsFilter.sortingType);
    }

    function updateTransactionStatisticsInvalidState(invalidState: boolean): void {
        transactionStatisticsStateInvalid.value = invalidState;
    }

    function resetTransactionStatistics(): void {
        transactionStatisticsFilter.value.chartDataType = ChartDataType.Default.type;
        transactionStatisticsFilter.value.categoricalChartType = CategoricalChartType.Default.type;
        transactionStatisticsFilter.value.categoricalChartDateType = DEFAULT_CATEGORICAL_CHART_DATA_RANGE.type;
        transactionStatisticsFilter.value.categoricalChartStartTime = 0;
        transactionStatisticsFilter.value.categoricalChartEndTime = 0;
        transactionStatisticsFilter.value.trendChartType = TrendChartType.Default.type;
        transactionStatisticsFilter.value.trendChartDateType = DEFAULT_TREND_CHART_DATA_RANGE.type;
        transactionStatisticsFilter.value.trendChartStartYearMonth = '';
        transactionStatisticsFilter.value.trendChartEndYearMonth = '';
        transactionStatisticsFilter.value.filterAccountIds = {};
        transactionStatisticsFilter.value.filterCategoryIds = {};
        transactionStatisticsFilter.value.tagIds = '';
        transactionStatisticsFilter.value.tagFilterType = TransactionTagFilterType.Default.type;
        transactionCategoryStatisticsData.value = null;
        transactionCategoryTrendsData.value = [];
        transactionStatisticsStateInvalid.value = true;
    }

    function initTransactionStatisticsFilter(analysisType: StatisticsAnalysisType, filter?: TransactionStatisticsPartialFilter): void {
        if (filter && isInteger(filter.chartDataType)) {
            transactionStatisticsFilter.value.chartDataType = filter.chartDataType;
        } else {
            transactionStatisticsFilter.value.chartDataType = settingsStore.appSettings.statistics.defaultChartDataType;
        }

        if (analysisType === StatisticsAnalysisType.CategoricalAnalysis || analysisType === StatisticsAnalysisType.TrendAnalysis) {
            if (!ChartDataType.isAvailableForAnalysisType(transactionStatisticsFilter.value.chartDataType, analysisType)) {
                transactionStatisticsFilter.value.chartDataType = ChartDataType.Default.type;
            }
        }

        if (filter && isInteger(filter.categoricalChartType)) {
            transactionStatisticsFilter.value.categoricalChartType = filter.categoricalChartType;
        } else {
            transactionStatisticsFilter.value.categoricalChartType = settingsStore.appSettings.statistics.defaultCategoricalChartType;
        }

        if (transactionStatisticsFilter.value.categoricalChartType !== CategoricalChartType.Pie.type && transactionStatisticsFilter.value.categoricalChartType !== CategoricalChartType.Bar.type) {
            transactionStatisticsFilter.value.categoricalChartType = CategoricalChartType.Default.type;
        }

        if (filter && isInteger(filter.categoricalChartDateType)) {
            transactionStatisticsFilter.value.categoricalChartDateType = filter.categoricalChartDateType;
        } else {
            transactionStatisticsFilter.value.categoricalChartDateType = settingsStore.appSettings.statistics.defaultCategoricalChartDataRangeType;
        }

        let categoricalChartDateTypeValid = true;

        if (!DateRange.isAvailableForScene(transactionStatisticsFilter.value.categoricalChartDateType, DateRangeScene.Normal)) {
            transactionStatisticsFilter.value.categoricalChartDateType = DEFAULT_CATEGORICAL_CHART_DATA_RANGE.type;
            categoricalChartDateTypeValid = false;
        }

        if (categoricalChartDateTypeValid && transactionStatisticsFilter.value.categoricalChartDateType === DateRange.Custom.type) {
            if (filter && isInteger(filter.categoricalChartStartTime)) {
                transactionStatisticsFilter.value.categoricalChartStartTime = filter.categoricalChartStartTime;
            } else {
                transactionStatisticsFilter.value.categoricalChartStartTime = 0;
            }

            if (filter && isInteger(filter.categoricalChartEndTime)) {
                transactionStatisticsFilter.value.categoricalChartEndTime = filter.categoricalChartEndTime;
            } else {
                transactionStatisticsFilter.value.categoricalChartEndTime = 0;
            }
        } else {
            const categoricalChartDateRange = getDateRangeByDateType(transactionStatisticsFilter.value.categoricalChartDateType, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart);

            if (categoricalChartDateRange) {
                transactionStatisticsFilter.value.categoricalChartDateType = categoricalChartDateRange.dateType;
                transactionStatisticsFilter.value.categoricalChartStartTime = categoricalChartDateRange.minTime;
                transactionStatisticsFilter.value.categoricalChartEndTime = categoricalChartDateRange.maxTime;
            }
        }

        if (filter && isInteger(filter.trendChartType)) {
            transactionStatisticsFilter.value.trendChartType = filter.trendChartType;
        } else {
            transactionStatisticsFilter.value.trendChartType = settingsStore.appSettings.statistics.defaultTrendChartType;
        }

        if (transactionStatisticsFilter.value.trendChartType !== TrendChartType.Area.type && transactionStatisticsFilter.value.trendChartType !== TrendChartType.Column.type) {
            transactionStatisticsFilter.value.trendChartType = TrendChartType.Default.type;
        }

        if (filter && isInteger(filter.trendChartDateType)) {
            transactionStatisticsFilter.value.trendChartDateType = filter.trendChartDateType;
        } else {
            transactionStatisticsFilter.value.trendChartDateType = settingsStore.appSettings.statistics.defaultTrendChartDataRangeType;
        }

        let trendChartDateTypeValid = true;

        if (!DateRange.isAvailableForScene(transactionStatisticsFilter.value.trendChartDateType, DateRangeScene.TrendAnalysis)) {
            transactionStatisticsFilter.value.trendChartDateType = DEFAULT_TREND_CHART_DATA_RANGE.type;
            trendChartDateTypeValid = false;
        }

        if (trendChartDateTypeValid && transactionStatisticsFilter.value.trendChartDateType === DateRange.Custom.type) {
            if (filter && isYearMonth(filter.trendChartStartYearMonth)) {
                transactionStatisticsFilter.value.trendChartStartYearMonth = filter.trendChartStartYearMonth;
            } else {
                transactionStatisticsFilter.value.trendChartStartYearMonth = '';
            }

            if (filter && isYearMonth(filter.trendChartEndYearMonth)) {
                transactionStatisticsFilter.value.trendChartEndYearMonth = filter.trendChartEndYearMonth;
            } else {
                transactionStatisticsFilter.value.trendChartEndYearMonth = '';
            }
        } else {
            const trendChartDateRange = getDateRangeByDateType(transactionStatisticsFilter.value.trendChartDateType, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart);

            if (trendChartDateRange) {
                transactionStatisticsFilter.value.trendChartDateType = trendChartDateRange.dateType;
                transactionStatisticsFilter.value.trendChartStartYearMonth = getYearAndMonthFromUnixTime(trendChartDateRange.minTime);
                transactionStatisticsFilter.value.trendChartEndYearMonth = getYearAndMonthFromUnixTime(trendChartDateRange.maxTime);
            }
        }

        if (filter && isObject(filter.filterAccountIds)) {
            transactionStatisticsFilter.value.filterAccountIds = filter.filterAccountIds;
        } else {
            transactionStatisticsFilter.value.filterAccountIds = settingsStore.appSettings.statistics.defaultAccountFilter || {};
        }

        if (filter && isObject(filter.filterCategoryIds)) {
            transactionStatisticsFilter.value.filterCategoryIds = filter.filterCategoryIds;
        } else {
            transactionStatisticsFilter.value.filterCategoryIds = settingsStore.appSettings.statistics.defaultTransactionCategoryFilter || {};
        }

        if (filter && isString(filter.tagIds)) {
            transactionStatisticsFilter.value.tagIds = filter.tagIds;
        } else {
            transactionStatisticsFilter.value.tagIds = '';
        }

        if (filter && isInteger(filter.tagFilterType)) {
            transactionStatisticsFilter.value.tagFilterType = filter.tagFilterType;
        } else {
            transactionStatisticsFilter.value.tagFilterType = TransactionTagFilterType.Default.type;
        }

        if (filter && isInteger(filter.sortingType)) {
            transactionStatisticsFilter.value.sortingType = filter.sortingType;
        } else {
            transactionStatisticsFilter.value.sortingType = settingsStore.appSettings.statistics.defaultSortingType;
        }

        if (transactionStatisticsFilter.value.sortingType < ChartSortingType.Amount.type || transactionStatisticsFilter.value.sortingType > ChartSortingType.Name.type) {
            transactionStatisticsFilter.value.sortingType = ChartSortingType.Default.type;
        }
    }

    function updateTransactionStatisticsFilter(filter: TransactionStatisticsPartialFilter): boolean {
        let changed = false;

        if (filter && isInteger(filter.chartDataType) && transactionStatisticsFilter.value.chartDataType !== filter.chartDataType) {
            transactionStatisticsFilter.value.chartDataType = filter.chartDataType;
            changed = true;
        }

        if (filter && isInteger(filter.categoricalChartType) && transactionStatisticsFilter.value.categoricalChartType !== filter.categoricalChartType) {
            transactionStatisticsFilter.value.categoricalChartType = filter.categoricalChartType;
            changed = true;
        }

        if (filter && isInteger(filter.categoricalChartDateType) && transactionStatisticsFilter.value.categoricalChartDateType !== filter.categoricalChartDateType) {
            transactionStatisticsFilter.value.categoricalChartDateType = filter.categoricalChartDateType;
            changed = true;
        }

        if (filter && isInteger(filter.categoricalChartStartTime) && transactionStatisticsFilter.value.categoricalChartStartTime !== filter.categoricalChartStartTime) {
            transactionStatisticsFilter.value.categoricalChartStartTime = filter.categoricalChartStartTime;
            changed = true;
        }

        if (filter && isInteger(filter.categoricalChartEndTime) && transactionStatisticsFilter.value.categoricalChartEndTime !== filter.categoricalChartEndTime) {
            transactionStatisticsFilter.value.categoricalChartEndTime = filter.categoricalChartEndTime;
            changed = true;
        }

        if (filter && isInteger(filter.trendChartType) && transactionStatisticsFilter.value.trendChartType !== filter.trendChartType) {
            transactionStatisticsFilter.value.trendChartType = filter.trendChartType;
            changed = true;
        }

        if (filter && isInteger(filter.trendChartDateType) && transactionStatisticsFilter.value.trendChartDateType !== filter.trendChartDateType) {
            transactionStatisticsFilter.value.trendChartDateType = filter.trendChartDateType;
            changed = true;
        }

        if (filter && (isYearMonth(filter.trendChartStartYearMonth) || filter.trendChartStartYearMonth === '') && !isYearMonthEquals(transactionStatisticsFilter.value.trendChartStartYearMonth, filter.trendChartStartYearMonth)) {
            transactionStatisticsFilter.value.trendChartStartYearMonth = filter.trendChartStartYearMonth;
            changed = true;
        }

        if (filter && (isYearMonth(filter.trendChartEndYearMonth) || filter.trendChartEndYearMonth === '') && !isYearMonthEquals(transactionStatisticsFilter.value.trendChartEndYearMonth, filter.trendChartEndYearMonth)) {
            transactionStatisticsFilter.value.trendChartEndYearMonth = filter.trendChartEndYearMonth;
            changed = true;
        }

        if (filter && isObject(filter.filterAccountIds) && !isEquals(transactionStatisticsFilter.value.filterAccountIds, filter.filterAccountIds)) {
            transactionStatisticsFilter.value.filterAccountIds = filter.filterAccountIds;
            changed = true;
        }

        if (filter && isObject(filter.filterCategoryIds) && !isEquals(transactionStatisticsFilter.value.filterCategoryIds, filter.filterCategoryIds)) {
            transactionStatisticsFilter.value.filterCategoryIds = filter.filterCategoryIds;
            changed = true;
        }

        if (filter && isString(filter.tagIds) && transactionStatisticsFilter.value.tagIds !== filter.tagIds) {
            transactionStatisticsFilter.value.tagIds = filter.tagIds;
            changed = true;
        }

        if (filter && isInteger(filter.tagFilterType) && transactionStatisticsFilter.value.tagFilterType !== filter.tagFilterType) {
            transactionStatisticsFilter.value.tagFilterType = filter.tagFilterType;
            changed = true;
        }

        if (filter && isInteger(filter.sortingType) && transactionStatisticsFilter.value.sortingType !== filter.sortingType) {
            transactionStatisticsFilter.value.sortingType = filter.sortingType;
            changed = true;
        }

        return changed;
    }

    function getTransactionStatisticsPageParams(analysisType: StatisticsAnalysisType, trendDateAggregationType: number): string {
        const querys: string[] = [];

        querys.push('analysisType=' + analysisType);
        querys.push('chartDataType=' + transactionStatisticsFilter.value.chartDataType);

        if (analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
            querys.push('chartType=' + transactionStatisticsFilter.value.categoricalChartType);
            querys.push('chartDateType=' + transactionStatisticsFilter.value.categoricalChartDateType);

            if (transactionStatisticsFilter.value.categoricalChartDateType === DateRange.Custom.type) {
                querys.push('startTime=' + transactionStatisticsFilter.value.categoricalChartStartTime);
                querys.push('endTime=' + transactionStatisticsFilter.value.categoricalChartEndTime);
            }
        } else if (analysisType === StatisticsAnalysisType.TrendAnalysis) {
            querys.push('chartType=' + transactionStatisticsFilter.value.trendChartType);
            querys.push('chartDateType=' + transactionStatisticsFilter.value.trendChartDateType);

            if (transactionStatisticsFilter.value.trendChartDateType === DateRange.Custom.type) {
                querys.push('startTime=' + transactionStatisticsFilter.value.trendChartStartYearMonth);
                querys.push('endTime=' + transactionStatisticsFilter.value.trendChartEndYearMonth);
            }

            if (trendDateAggregationType !== ChartDateAggregationType.Month.type) {
                querys.push('trendDateAggregationType=' + trendDateAggregationType);
            }
        }

        if (transactionStatisticsFilter.value.filterAccountIds) {
            const ids = objectFieldToArrayItem(transactionStatisticsFilter.value.filterAccountIds);

            if (ids && ids.length) {
                querys.push('filterAccountIds=' + ids.join(','));
            }
        }

        if (transactionStatisticsFilter.value.filterCategoryIds) {
            const ids = objectFieldToArrayItem(transactionStatisticsFilter.value.filterCategoryIds);

            if (ids && ids.length) {
                querys.push('filterCategoryIds=' + ids.join(','));
            }
        }

        if (transactionStatisticsFilter.value.tagIds) {
            querys.push('tagIds=' + transactionStatisticsFilter.value.tagIds);
        }

        if (transactionStatisticsFilter.value.tagFilterType) {
            querys.push('tagFilterType=' + transactionStatisticsFilter.value.tagFilterType);
        }

        querys.push('sortingType=' + transactionStatisticsFilter.value.sortingType);

        return querys.join('&');
    }

    function getTransactionListPageParams(analysisType: StatisticsAnalysisType, itemId: string, dateRange?: TimeRangeAndDateType): string {
        const querys: string[] = [];

        if (transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByAccount.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.TotalIncome.type) {
            querys.push('type=2');
        } else if (transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.TotalExpense.type) {
            querys.push('type=3');
        }

        if (itemId && (transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByAccount.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type)) {
            querys.push('accountIds=' + itemId);

            if (!isObjectEmpty(transactionStatisticsFilter.value.filterCategoryIds)) {
                querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, transactionStatisticsFilter.value.filterCategoryIds));
            }
        } else if (itemId && (transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type
            || transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type)) {
            querys.push('categoryIds=' + itemId);

            if (!isObjectEmpty(transactionStatisticsFilter.value.filterAccountIds)) {
                querys.push('accountIds=' + getFinalAccountIdsByFilteredAccountIds(accountsStore.allAccountsMap, transactionStatisticsFilter.value.filterAccountIds));
            }
        } else if (!itemId) {
            if (!isObjectEmpty(transactionStatisticsFilter.value.filterCategoryIds)) {
                querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, transactionStatisticsFilter.value.filterCategoryIds));
            }

            if (!isObjectEmpty(transactionStatisticsFilter.value.filterAccountIds)) {
                querys.push('accountIds=' + getFinalAccountIdsByFilteredAccountIds(accountsStore.allAccountsMap, transactionStatisticsFilter.value.filterAccountIds));
            }
        }

        if (transactionStatisticsFilter.value.tagIds) {
            querys.push('tagIds=' + transactionStatisticsFilter.value.tagIds);
        }

        if (transactionStatisticsFilter.value.tagFilterType) {
            querys.push('tagFilterType=' + transactionStatisticsFilter.value.tagFilterType);
        }

        if (analysisType === StatisticsAnalysisType.CategoricalAnalysis
            && transactionStatisticsFilter.value.chartDataType !== ChartDataType.AccountTotalAssets.type
            && transactionStatisticsFilter.value.chartDataType !== ChartDataType.AccountTotalLiabilities.type) {
            querys.push('dateType=' + transactionStatisticsFilter.value.categoricalChartDateType);

            if (transactionStatisticsFilter.value.categoricalChartDateType === DateRange.Custom.type) {
                querys.push('minTime=' + transactionStatisticsFilter.value.categoricalChartStartTime);
                querys.push('maxTime=' + transactionStatisticsFilter.value.categoricalChartEndTime);
            }
        } else if (analysisType === StatisticsAnalysisType.TrendAnalysis && dateRange) {
            querys.push('dateType=' + dateRange.dateType);
            querys.push('minTime=' + dateRange.minTime);
            querys.push('maxTime=' + dateRange.maxTime);
        }

        return querys.join('&');
    }

    function loadCategoricalAnalysis({ force }: { force: boolean }): Promise<TransactionStatisticResponse> {
        return new Promise((resolve, reject) => {
            services.getTransactionStatistics({
                startTime: transactionStatisticsFilter.value.categoricalChartStartTime,
                endTime: transactionStatisticsFilter.value.categoricalChartEndTime,
                tagIds: transactionStatisticsFilter.value.tagIds,
                tagFilterType: transactionStatisticsFilter.value.tagFilterType,
                useTransactionTimezone: settingsStore.appSettings.statistics.defaultTimezoneType === TimezoneTypeForStatistics.TransactionTimezone.type
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve transaction statistics' });
                    return;
                }

                if (transactionStatisticsStateInvalid.value) {
                    updateTransactionStatisticsInvalidState(false);
                }

                if (force && data.result && isEquals(transactionCategoryStatisticsData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionCategoryStatisticsData.value = data.result;

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
    }

    function loadTrendAnalysis({ force }: { force: boolean }): Promise<TransactionStatisticTrendsResponseItem[]> {
        return new Promise((resolve, reject) => {
            services.getTransactionStatisticsTrends({
                startYearMonth: transactionStatisticsFilter.value.trendChartStartYearMonth,
                endYearMonth: transactionStatisticsFilter.value.trendChartEndYearMonth,
                tagIds: transactionStatisticsFilter.value.tagIds,
                tagFilterType: transactionStatisticsFilter.value.tagFilterType,
                useTransactionTimezone: settingsStore.appSettings.statistics.defaultTimezoneType === TimezoneTypeForStatistics.TransactionTimezone.type
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve transaction statistics' });
                    return;
                }

                if (transactionStatisticsStateInvalid.value) {
                    updateTransactionStatisticsInvalidState(false);
                }

                if (force && data.result && isEquals(transactionCategoryTrendsData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionCategoryTrendsData.value = data.result;

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
    }

    return {
        // states
        transactionStatisticsFilter,
        transactionCategoryStatisticsData,
        transactionCategoryTrendsData,
        transactionStatisticsStateInvalid,
        // computed states
        categoricalAnalysisChartDataCategory,
        categoricalAnalysisData,
        trendsAnalysisData,
        // functions
        updateTransactionStatisticsInvalidState,
        resetTransactionStatistics,
        initTransactionStatisticsFilter,
        updateTransactionStatisticsFilter,
        getTransactionStatisticsPageParams,
        getTransactionListPageParams,
        loadCategoricalAnalysis,
        loadTrendAnalysis
    };
});
