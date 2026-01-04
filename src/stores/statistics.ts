import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { entries, values } from '@/core/base.ts';
import { type DateTime, type TextualYearMonth, type TimeRangeAndDateType, DateRangeScene, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { CategoryType } from '@/core/category.ts';
import {
    TransactionRelatedAccountType
} from '@/core/transaction.ts';
import {
    StatisticsAnalysisType,
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType,
    DEFAULT_CATEGORICAL_CHART_DATA_RANGE,
    DEFAULT_TREND_CHART_DATA_RANGE,
    DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE
} from '@/core/statistics.ts';
import { DEFAULT_ACCOUNT_ICON, DEFAULT_CATEGORY_ICON } from '@/consts/icon.ts';
import { DEFAULT_ACCOUNT_COLOR, DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';

import {
    type TransactionStatisticResponse,
    type TransactionStatisticResponseItem,
    type TransactionStatisticTrendsResponseItem,
    type TransactionStatisticAssetTrendsResponseItem,
    type TransactionStatisticAssetTrendsResponseDataItem,
    type TransactionStatisticResponseItemWithInfo,
    type TransactionStatisticResponseWithInfo,
    type TransactionStatisticTrendsResponseItemWithInfo,
    type TransactionStatisticAssetTrendsResponseItemWithInfo,
    type TransactionStatisticDataItemType,
    type TransactionStatisticDataItemBase,
    type TransactionCategoricalOverviewAnalysisData,
    type TransactionCategoricalOverviewAnalysisDataItem,
    type TransactionCategoricalAnalysisData,
    type TransactionCategoricalAnalysisDataItem,
    type TransactionTrendsAnalysisData,
    type TransactionTrendsAnalysisDataItem,
    type TransactionTrendsAnalysisDataAmount,
    type TransactionAssetTrendsAnalysisData,
    type TransactionAssetTrendsAnalysisDataItem,
    type TransactionAssetTrendsAnalysisDataAmount,
    TransactionCategoricalOverviewAnalysisDataItemType
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
import {
    getYearMonthDayDateTime,
    getGregorianCalendarYearAndMonthFromUnixTime,
    getDayDifference,
    getDateRangeByDateType
} from '@/lib/datetime.ts';
import { getFinalAccountIdsByFilteredAccountIds } from '@/lib/account.ts';
import { getFinalCategoryIdsByFilteredCategoryIds } from '@/lib/category.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';
import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

interface WritableTransactionCategoricalAnalysisData {
    totalAmount: number;
    totalNonNegativeAmount: number;
    items: Record<string, WritableTransactionCategoricalAnalysisDataItem>;
}

interface WritableTransactionCategoricalAnalysisDataItem extends Record<string, unknown> {
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

interface WritableTransactionTrendsAnalysisDataItem extends Record<string, unknown>, TransactionTrendsAnalysisDataItem {
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

interface WritableTransactionAssetTrendsAnalysisDataItem extends Record<string, unknown>, TransactionAssetTrendsAnalysisDataItem {
    name: string;
    type: TransactionStatisticDataItemType;
    id: string;
    icon: string;
    color: string;
    hidden: boolean;
    displayOrders: number[];
    totalAmount: number;
    items: TransactionAssetTrendsAnalysisDataAmount[];
}

export interface TransactionStatisticsPartialFilter {
    chartDataType?: number;
    categoricalChartType?: number;
    categoricalChartDateType?: number;
    categoricalChartStartTime?: number;
    categoricalChartEndTime?: number;
    trendChartType?: number;
    trendChartDateType?: number;
    trendChartStartYearMonth?: TextualYearMonth | '';
    trendChartEndYearMonth?: TextualYearMonth | '';
    assetTrendsChartType?: number;
    assetTrendsChartDateType?: number;
    assetTrendsChartStartTime?: number;
    assetTrendsChartEndTime?: number;
    filterAccountIds?: Record<string, boolean>;
    filterCategoryIds?: Record<string, boolean>;
    tagFilter?: string;
    keyword?: string;
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
    trendChartStartYearMonth: TextualYearMonth | '';
    trendChartEndYearMonth: TextualYearMonth | '';
    assetTrendsChartType: number;
    assetTrendsChartDateType: number;
    assetTrendsChartStartTime: number;
    assetTrendsChartEndTime: number;
    filterAccountIds: Record<string, boolean>;
    filterCategoryIds: Record<string, boolean>;
    tagFilter: string;
    keyword: string;
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
        assetTrendsChartType: TrendChartType.Default.type,
        assetTrendsChartDateType: DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE.type,
        assetTrendsChartStartTime: 0,
        assetTrendsChartEndTime: 0,
        filterAccountIds: {},
        filterCategoryIds: {},
        tagFilter: '',
        keyword: '',
        sortingType: ChartSortingType.Default.type
    });

    const transactionCategoryStatisticsData = ref<TransactionStatisticResponse | null>(null);
    const transactionCategoryTrendsData = ref<TransactionStatisticTrendsResponseItem[]>([]);
    const transactionAssetTrendsData = ref<TransactionStatisticAssetTrendsResponseItem[]>([]);
    const transactionStatisticsStateInvalid = ref<boolean>(true);

    const categoricalAnalysisChartDataCategory = computed<string>(() => {
        if (transactionStatisticsFilter.value.chartDataType === ChartDataType.OutflowsByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.InflowsByAccount.type ||
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

    const categoricalOverviewAnalysisData = computed<TransactionCategoricalOverviewAnalysisData | null>(() => {
        if (!transactionCategoryStatisticsDataWithCategoryAndAccountInfo.value || !transactionCategoryStatisticsDataWithCategoryAndAccountInfo.value.items) {
            return null;
        }

        const allDataItemsMap: Record<string, TransactionCategoricalOverviewAnalysisDataItem> = {};
        const allIncomeByPrimaryCategoryDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allIncomeBySecondaryCategoryDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allIncomeByAccountDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allExpenseByAccountDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allExpenseBySecondaryCategoryDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allExpenseByPrimaryCategoryDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allOpeningBalanceDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];
        const allNetCashFlowDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [];

        let totalIncome: number = 0;
        let totalExpense: number = 0;

        for (const item of transactionCategoryStatisticsDataWithCategoryAndAccountInfo.value.items) {
            if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category) {
                continue;
            }

            if (item.relatedAccount && item.relatedAccountType === TransactionRelatedAccountType.TransferFrom) {
                continue;
            }

            if (!isNumber(item.amountInDefaultCurrency)) {
                continue;
            }

            if (transactionStatisticsFilter.value.filterAccountIds && transactionStatisticsFilter.value.filterAccountIds[item.account.id]) {
                continue;
            }

            if (transactionStatisticsFilter.value.filterCategoryIds && transactionStatisticsFilter.value.filterCategoryIds[item.category.id]) {
                continue;
            }

            if (item.category.type === CategoryType.Income) {
                totalIncome += item.amountInDefaultCurrency;
            } else if (item.category.type === CategoryType.Expense) {
                totalExpense += item.amountInDefaultCurrency;
            }

            const primaryAccountCategoryDisplayOrder = settingsStore.accountCategoryDisplayOrders[item.primaryAccount.category] || Number.MAX_SAFE_INTEGER;
            const incomeByAccountKey = `${TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount}:${item.account.id}`;
            const expenseByAccountKey = `${TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount}:${item.account.id}`;
            let incomeByAccountItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[incomeByAccountKey];
            let expenseByAccountItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[expenseByAccountKey];

            if (!incomeByAccountItem) {
                incomeByAccountItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                    item.account.id,
                    item.account.name,
                    TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount,
                    [primaryAccountCategoryDisplayOrder, item.primaryAccount.displayOrder, item.account.displayOrder],
                    item.primaryAccount.hidden || item.account.hidden);
                allDataItemsMap[incomeByAccountKey] = incomeByAccountItem;
                allIncomeByAccountDataItems.push(incomeByAccountItem);
            }

            if (!expenseByAccountItem) {
                expenseByAccountItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                    item.account.id,
                    item.account.name,
                    TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount,
                    [primaryAccountCategoryDisplayOrder, item.primaryAccount.displayOrder, item.account.displayOrder],
                    item.primaryAccount.hidden || item.account.hidden);
                allDataItemsMap[expenseByAccountKey] = expenseByAccountItem;
                allExpenseByAccountDataItems.push(expenseByAccountItem);
            }

            if (item.category.type === CategoryType.Income) {
                const primaryCategoryItemKey = `${TransactionCategoricalOverviewAnalysisDataItemType.IncomeByPrimaryCategory}:${item.primaryCategory.id}`;
                const secondaryCategoryItemKey = `${TransactionCategoricalOverviewAnalysisDataItemType.IncomeBySecondaryCategory}:${item.category.id}`;

                let primaryCategoryDataItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[primaryCategoryItemKey];
                let secondaryCategoryDataItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[secondaryCategoryItemKey];

                if (!primaryCategoryDataItem) {
                    primaryCategoryDataItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                        item.primaryCategory.id,
                        item.primaryCategory.name,
                        TransactionCategoricalOverviewAnalysisDataItemType.IncomeByPrimaryCategory,
                        [item.primaryCategory.displayOrder],
                        item.primaryCategory.hidden);
                    allDataItemsMap[primaryCategoryItemKey] = primaryCategoryDataItem;
                    allIncomeByPrimaryCategoryDataItems.push(primaryCategoryDataItem);
                }

                if (!secondaryCategoryDataItem) {
                    secondaryCategoryDataItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                        item.category.id,
                        item.category.name,
                        TransactionCategoricalOverviewAnalysisDataItemType.IncomeBySecondaryCategory,
                        [item.primaryCategory.displayOrder, item.category.displayOrder],
                        item.primaryCategory.hidden || item.category.hidden);
                    allDataItemsMap[secondaryCategoryItemKey] = secondaryCategoryDataItem;
                    allIncomeBySecondaryCategoryDataItems.push(secondaryCategoryDataItem);
                }

                primaryCategoryDataItem.totalAmount += item.amountInDefaultCurrency;
                primaryCategoryDataItem.totalNonNegativeAmount += item.amountInDefaultCurrency > 0 ? item.amountInDefaultCurrency : 0;
                primaryCategoryDataItem.includeInPercent = true;
                primaryCategoryDataItem.outflows.push({ amount: item.amountInDefaultCurrency, relatedItem: secondaryCategoryDataItem });

                secondaryCategoryDataItem.totalAmount += item.amountInDefaultCurrency;
                secondaryCategoryDataItem.totalNonNegativeAmount += item.amountInDefaultCurrency > 0 ? item.amountInDefaultCurrency : 0;
                secondaryCategoryDataItem.includeInPercent = true;
                secondaryCategoryDataItem.inflows.push({ amount: item.amountInDefaultCurrency, relatedItem: primaryCategoryDataItem });
                secondaryCategoryDataItem.outflows.push({ amount: item.amountInDefaultCurrency, relatedItem: incomeByAccountItem });

                incomeByAccountItem.totalAmount += item.amountInDefaultCurrency;
                incomeByAccountItem.totalNonNegativeAmount += item.amountInDefaultCurrency > 0 ? item.amountInDefaultCurrency : 0;
                incomeByAccountItem.includeInPercent = true;
                incomeByAccountItem.inflows.push({ amount: item.amountInDefaultCurrency, relatedItem: secondaryCategoryDataItem });
            } else if (item.category.type === CategoryType.Expense) {
                const primaryCategoryItemKey = `${TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByPrimaryCategory}:${item.primaryCategory.id}`;
                const secondaryCategoryItemKey = `${TransactionCategoricalOverviewAnalysisDataItemType.ExpenseBySecondaryCategory}:${item.category.id}`;

                let primaryCategoryDataItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[primaryCategoryItemKey];
                let secondaryCategoryDataItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[secondaryCategoryItemKey];

                if (!primaryCategoryDataItem) {
                    primaryCategoryDataItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                        item.primaryCategory.id,
                        item.primaryCategory.name,
                        TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByPrimaryCategory,
                        [item.primaryCategory.displayOrder],
                        item.primaryCategory.hidden);
                    allDataItemsMap[primaryCategoryItemKey] = primaryCategoryDataItem;
                    allExpenseByPrimaryCategoryDataItems.push(primaryCategoryDataItem);
                }

                if (!secondaryCategoryDataItem) {
                    secondaryCategoryDataItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                        item.category.id,
                        item.category.name,
                        TransactionCategoricalOverviewAnalysisDataItemType.ExpenseBySecondaryCategory,
                        [item.primaryCategory.displayOrder, item.category.displayOrder],
                        item.primaryCategory.hidden || item.category.hidden);
                    allDataItemsMap[secondaryCategoryItemKey] = secondaryCategoryDataItem;
                    allExpenseBySecondaryCategoryDataItems.push(secondaryCategoryDataItem);
                }

                expenseByAccountItem.totalAmount += item.amountInDefaultCurrency;
                expenseByAccountItem.totalNonNegativeAmount += item.amountInDefaultCurrency > 0 ? item.amountInDefaultCurrency : 0;
                expenseByAccountItem.includeInPercent = true;
                expenseByAccountItem.outflows.push({ amount: item.amountInDefaultCurrency, relatedItem: secondaryCategoryDataItem });

                secondaryCategoryDataItem.totalAmount += item.amountInDefaultCurrency;
                secondaryCategoryDataItem.totalNonNegativeAmount += item.amountInDefaultCurrency > 0 ? item.amountInDefaultCurrency : 0
                secondaryCategoryDataItem.includeInPercent = true;
                secondaryCategoryDataItem.inflows.push({ amount: item.amountInDefaultCurrency, relatedItem: expenseByAccountItem });
                secondaryCategoryDataItem.outflows.push({ amount: item.amountInDefaultCurrency, relatedItem: primaryCategoryDataItem });

                primaryCategoryDataItem.totalAmount += item.amountInDefaultCurrency;
                primaryCategoryDataItem.totalNonNegativeAmount += item.amountInDefaultCurrency > 0 ? item.amountInDefaultCurrency : 0;
                primaryCategoryDataItem.includeInPercent = true;
                primaryCategoryDataItem.inflows.push({ amount: item.amountInDefaultCurrency, relatedItem: secondaryCategoryDataItem });
            } else if (item.category.type === CategoryType.Transfer && item.relatedPrimaryAccount && item.relatedAccount) {
                const transferToAccountKey = `${TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount}:${item.relatedAccount.id}`;
                let transferToAccountItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[transferToAccountKey];

                if (!transferToAccountItem) {
                    const relatedPrimaryAccountCategoryDisplayOrder = settingsStore.accountCategoryDisplayOrders[item.relatedPrimaryAccount.category] || Number.MAX_SAFE_INTEGER;
                    transferToAccountItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                        item.relatedAccount.id,
                        item.relatedAccount.name,
                        TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount,
                        [relatedPrimaryAccountCategoryDisplayOrder, item.relatedPrimaryAccount.displayOrder, item.relatedAccount.displayOrder],
                        item.relatedPrimaryAccount.hidden || item.relatedAccount.hidden);
                    allDataItemsMap[transferToAccountKey] = transferToAccountItem;
                    allExpenseByAccountDataItems.push(transferToAccountItem);
                }

                incomeByAccountItem.outflows.push({ amount: item.amountInDefaultCurrency, relatedItem: transferToAccountItem });
                transferToAccountItem.inflows.push({ amount: item.amountInDefaultCurrency, relatedItem: incomeByAccountItem });
            }
        }

        sortCategoricalOverviewAnalysisDataItems(allIncomeByPrimaryCategoryDataItems, transactionStatisticsFilter.value);
        sortCategoricalOverviewAnalysisDataItems(allIncomeBySecondaryCategoryDataItems, transactionStatisticsFilter.value);
        sortCategoricalOverviewAnalysisDataItems(allIncomeByAccountDataItems, transactionStatisticsFilter.value);
        sortCategoricalOverviewAnalysisDataItems(allExpenseByAccountDataItems, transactionStatisticsFilter.value);
        sortCategoricalOverviewAnalysisDataItems(allExpenseBySecondaryCategoryDataItems, transactionStatisticsFilter.value);
        sortCategoricalOverviewAnalysisDataItems(allExpenseByPrimaryCategoryDataItems, transactionStatisticsFilter.value);

        for (const item of allExpenseByAccountDataItems) {
            const incomeByAccountKey = `${TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount}:${item.id}`;
            const incomeByAccountItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[incomeByAccountKey];

            let accountTotalInflowsAmount: number = 0;
            let accountTotalIncomeAmount: number = 0;
            let accountTotalTransferAmount: number = 0;
            let accountTotalOutflowsAmount: number = 0;

            if (incomeByAccountItem) {
                for (const inflow of incomeByAccountItem.inflows) {
                    accountTotalInflowsAmount += inflow.amount;
                    accountTotalIncomeAmount += inflow.amount;
                }

                for (const outflow of incomeByAccountItem.outflows) {
                    accountTotalTransferAmount += outflow.amount;
                }
            }

            for (const inflow of item.inflows) {
                if (inflow.relatedItem.type === item.type && inflow.relatedItem.id === item.id) {
                    continue;
                }

                accountTotalInflowsAmount += inflow.amount;
            }

            for (const outflow of item.outflows) {
                accountTotalOutflowsAmount += outflow.amount;
            }

            const accountBalance: number = accountTotalIncomeAmount - accountTotalTransferAmount - accountTotalOutflowsAmount;
            const accountNetCashFlow: number = accountTotalInflowsAmount - accountTotalTransferAmount - accountTotalOutflowsAmount;

            if (incomeByAccountItem && accountsStore.allAccountsMap[item.id]?.isAsset) {
                if (accountBalance > 0) { // has positive balance, transfer the amount from income account to expense account
                    incomeByAccountItem.outflows.push({ amount: accountBalance + accountTotalOutflowsAmount, relatedItem: item });
                    item.inflows.push({ amount: accountBalance + accountTotalOutflowsAmount, relatedItem: incomeByAccountItem });
                } else if (accountNetCashFlow < 0) { // has negative net cash flow, add the difference to income account
                    incomeByAccountItem.totalAmount += -accountNetCashFlow;
                    incomeByAccountItem.totalNonNegativeAmount += -accountNetCashFlow > 0 ? -accountNetCashFlow : 0;
                    incomeByAccountItem.outflows.push({ amount: -accountNetCashFlow, relatedItem: item });
                    item.inflows.push({ amount: -accountNetCashFlow, relatedItem: incomeByAccountItem });
                }
            }

            if (accountNetCashFlow > 0) {
                let netCashFlowItem: TransactionCategoricalOverviewAnalysisDataItem | undefined = allDataItemsMap[TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow];

                if (!netCashFlowItem) {
                    netCashFlowItem = createNewTransactionCategoricalOverviewAnalysisDataItem(
                        TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow,
                        'Net Cash Flow',
                        TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow,
                        [Number.MAX_SAFE_INTEGER],
                        false);
                    allDataItemsMap[TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow] = netCashFlowItem;
                    allNetCashFlowDataItems.push(netCashFlowItem);
                }

                item.outflows.push({ amount: accountNetCashFlow, relatedItem: netCashFlowItem });

                netCashFlowItem.totalAmount += accountNetCashFlow;
                netCashFlowItem.totalNonNegativeAmount += accountNetCashFlow > 0 ? accountNetCashFlow : 0;
                netCashFlowItem.inflows.push({ amount: accountNetCashFlow, relatedItem: item });
            }
        }

        const allDataItems: TransactionCategoricalOverviewAnalysisDataItem[] = [
            ...allIncomeByPrimaryCategoryDataItems,
            ...allIncomeBySecondaryCategoryDataItems,
            ...allIncomeByAccountDataItems,
            ...allOpeningBalanceDataItems,
            ...allExpenseByAccountDataItems,
            ...allExpenseBySecondaryCategoryDataItems,
            ...allNetCashFlowDataItems,
            ...allExpenseByPrimaryCategoryDataItems
        ];

        return {
            totalIncome: totalIncome,
            totalExpense: totalExpense,
            items: allDataItems
        };
    });

    const accountTotalAmountAnalysisData = computed<WritableTransactionCategoricalAnalysisData | null>(() => {
        if (!accountsStore.allPlainAccounts) {
            return null;
        }

        const allDataItems: Record<string, WritableTransactionCategoricalAnalysisDataItem> = {};
        let totalAmount = 0;
        let totalNonNegativeAmount = 0;

        for (const account of accountsStore.allPlainAccounts) {
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

            const primaryAccountCategoryDisplayOrder = settingsStore.accountCategoryDisplayOrders[primaryAccount.category] || Number.MAX_SAFE_INTEGER;

            let amount = account.balance;

            if (account.currency !== userStore.currentUserDefaultCurrency) {
                const finalAmount = exchangeRatesStore.getExchangedAmount(amount, account.currency, userStore.currentUserDefaultCurrency);

                if (!isNumber(finalAmount)) {
                    continue;
                }

                amount = Math.trunc(finalAmount);
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
                displayOrders: [primaryAccountCategoryDisplayOrder, primaryAccount.displayOrder, account.displayOrder],
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

        if (transactionStatisticsFilter.value.chartDataType === ChartDataType.OutflowsByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.InflowsByAccount.type ||
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
            let maxTotalAmount = 0;

            for (const dataItem of values(combinedData.items)) {
                if (dataItem.totalAmount > maxTotalAmount) {
                    maxTotalAmount = dataItem.totalAmount;
                }
            }

            for (const dataItem of values(combinedData.items)) {
                let percent = 0;

                if (transactionStatisticsFilter.value.chartDataType === ChartDataType.OutflowsByAccount.type ||
                    transactionStatisticsFilter.value.chartDataType === ChartDataType.InflowsByAccount.type) {
                    if (maxTotalAmount > 0) {
                        percent = dataItem.totalAmount * 100 / maxTotalAmount;
                    } else {
                        percent = 0;
                    }

                    if (percent < 0) {
                        percent = 0;
                    }
                } else {
                    if (dataItem.totalAmount > 0) {
                        percent = dataItem.totalAmount * 100 / combinedData.totalNonNegativeAmount;
                    } else {
                        percent = 0;
                    }

                    if (percent < 0) {
                        percent = 0;
                    }
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
            for (const trendItem of trendsData) {
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

        for (const trendItem of transactionCategoryTrendsDataWithCategoryAndAccountInfo.value) {
            const totalAmountItems = getCategoryTotalAmountItems(trendItem.items, transactionStatisticsFilter.value);

            for (const [id, item] of entries(totalAmountItems.items)) {
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
                    month1base: trendItem.month,
                    totalAmount: item.totalAmount
                });

                combinedData.totalAmount += item.totalAmount;
                combinedDataMap[id] = combinedData;
            }
        }

        const totalAmountsTrends: TransactionTrendsAnalysisDataItem[] = [];

        for (const trendData of values(combinedDataMap)) {
            totalAmountsTrends.push(trendData);
        }

        sortCategoryTotalAmountItems(totalAmountsTrends, transactionStatisticsFilter.value);

        const trendsData: TransactionTrendsAnalysisData = {
            items: totalAmountsTrends
        };

        return trendsData;
    });

    const assetTrendsDataWithAccountInfo = computed<TransactionStatisticAssetTrendsResponseItemWithInfo[]>(() => {
        const assetTrendsData = transactionAssetTrendsData.value;
        const finalAssetTrendsData: TransactionStatisticAssetTrendsResponseItemWithInfo[] = [];

        if (!assetTrendsData || !assetTrendsData.length) {
            return finalAssetTrendsData;
        }

        const firstAssetTrendItem: TransactionStatisticAssetTrendsResponseItem | undefined = assetTrendsData[0];

        if (!firstAssetTrendItem) {
            return finalAssetTrendsData;
        }

        const lastAssetTrendItemMap: Record<string, TransactionStatisticAssetTrendsResponseDataItem> = {};
        let lastAssetTrendItem: TransactionStatisticAssetTrendsResponseItem = firstAssetTrendItem;

        for (const item of firstAssetTrendItem.items) {
            lastAssetTrendItemMap[item.accountId] = item;
        }

        for (const assetTrendItem of assetTrendsData) {
            const statisticResponseItems: TransactionStatisticResponseItem[] = [];
            const existedAccountIds: Record<string, boolean> = {};
            const missingDays: number = getDayDifference(lastAssetTrendItem, assetTrendItem) - 1;
            const lastAssetTrendItemDate: DateTime = getYearMonthDayDateTime(lastAssetTrendItem.year, lastAssetTrendItem.month, lastAssetTrendItem.day);

            // fill in missing days with last known balance
            for (let i = 1; i <= missingDays; i++) {
                const missingStatisticResponseItems: TransactionStatisticResponseItem[] = [];
                const dateTime: DateTime = lastAssetTrendItemDate.add(i, 'days');

                for (const item of values(lastAssetTrendItemMap)) {
                    const statisticResponseItem: TransactionStatisticResponseItem = {
                        categoryId: '',
                        accountId: item.accountId,
                        amount: item.accountClosingBalance
                    };

                    missingStatisticResponseItems.push(statisticResponseItem);
                }

                const finalAssetTrendItem: TransactionStatisticAssetTrendsResponseItemWithInfo = {
                    year: dateTime.getGregorianCalendarYear(),
                    month: dateTime.getGregorianCalendarMonth(),
                    day: dateTime.getGregorianCalendarDay(),
                    items: assembleAccountAndCategoryInfo(missingStatisticResponseItems)
                };

                lastAssetTrendItem = assetTrendItem;
                finalAssetTrendsData.push(finalAssetTrendItem);
            }

            // fill in current day data
            for (const item of assetTrendItem.items) {
                const statisticResponseItem: TransactionStatisticResponseItem = {
                    categoryId: '',
                    accountId: item.accountId,
                    amount: item.accountClosingBalance
                };

                lastAssetTrendItemMap[item.accountId] = item;
                existedAccountIds[item.accountId] = true;
                statisticResponseItems.push(statisticResponseItem);
            }

            // fill in missing accounts with last known balance
            for (const item of values(lastAssetTrendItemMap)) {
                if (existedAccountIds[item.accountId]) {
                    continue;
                }

                const statisticResponseItem: TransactionStatisticResponseItem = {
                    categoryId: '',
                    accountId: item.accountId,
                    amount: item.accountClosingBalance
                };

                existedAccountIds[item.accountId] = true;
                statisticResponseItems.push(statisticResponseItem);
            }

            const finalAssetTrendItem: TransactionStatisticAssetTrendsResponseItemWithInfo = {
                year: assetTrendItem.year,
                month: assetTrendItem.month,
                day: assetTrendItem.day,
                items: assembleAccountAndCategoryInfo(statisticResponseItems)
            };

            lastAssetTrendItem = assetTrendItem;
            finalAssetTrendsData.push(finalAssetTrendItem);
        }

        return finalAssetTrendsData;
    });

    const assetTrendsData = computed<TransactionAssetTrendsAnalysisData | null>(() => {
        if (!assetTrendsDataWithAccountInfo.value || !assetTrendsDataWithAccountInfo.value.length) {
            return null;
        }

        const combinedDataMap: Record<string, WritableTransactionAssetTrendsAnalysisDataItem> = {};

        for (const dailyData of assetTrendsDataWithAccountInfo.value) {
            let dailyTotalAmount: number = 0;

            for (const item of dailyData.items) {
                if (!item.primaryAccount || !item.account) {
                    continue;
                }

                if (transactionStatisticsFilter.value.filterAccountIds && transactionStatisticsFilter.value.filterAccountIds[item.account.id]) {
                    continue;
                }

                if (!isNumber(item.amountInDefaultCurrency)) {
                    continue;
                }

                if (transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type) {
                    if (!item.account.isAsset) {
                        continue;
                    }
                } else if (transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                    if (!item.account.isLiability) {
                        continue;
                    }
                } else if (transactionStatisticsFilter.value.chartDataType === ChartDataType.NetWorth.type) {
                    // Do Nothing
                } else {
                    continue;
                }

                let amount = item.amountInDefaultCurrency;

                if (item.account.isLiability) {
                    amount = -amount;
                }

                if (transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
                    transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                    let data = combinedDataMap[item.account.id];

                    if (data) {
                        data.totalAmount += amount;
                    } else {
                        const primaryAccountCategoryDisplayOrder = settingsStore.accountCategoryDisplayOrders[item.primaryAccount.category] || Number.MAX_SAFE_INTEGER;

                        data = {
                            name: item.account.name,
                            type: 'account',
                            id: item.account.id,
                            icon: item.account.icon || DEFAULT_ACCOUNT_ICON.icon,
                            color: item.account.color || DEFAULT_ACCOUNT_COLOR,
                            hidden: item.primaryAccount.hidden || item.account.hidden,
                            displayOrders: [primaryAccountCategoryDisplayOrder, item.primaryAccount.displayOrder, item.account.displayOrder],
                            totalAmount: amount,
                            items: []
                        };
                    }

                    const amountItem: TransactionAssetTrendsAnalysisDataAmount = {
                        year: dailyData.year,
                        month: dailyData.month,
                        day: dailyData.day,
                        totalAmount: amount
                    };
                    data.items.push(amountItem);
                    combinedDataMap[item.account.id] = data;
                }

                if (item.account.isAsset) {
                    dailyTotalAmount += amount;
                } else if (item.account.isLiability) {
                    dailyTotalAmount -= amount;
                }
            }

            if (transactionStatisticsFilter.value.chartDataType === ChartDataType.NetWorth.type) {
                let data = combinedDataMap['total'];

                if (data) {
                    data.totalAmount += dailyTotalAmount;
                } else {
                    data = {
                        name: ChartDataType.NetWorth.name,
                        type: 'total',
                        id: 'total',
                        icon: '',
                        color: '',
                        hidden: false,
                        displayOrders: [1],
                        totalAmount: dailyTotalAmount,
                        items: []
                    };
                }

                const amountItem: TransactionAssetTrendsAnalysisDataAmount = {
                    year: dailyData.year,
                    month: dailyData.month,
                    day: dailyData.day,
                    totalAmount: dailyTotalAmount
                };
                data.items.push(amountItem);
                combinedDataMap['total'] = data;
            }
        }

        const allAssetTrendsDataItems: TransactionAssetTrendsAnalysisDataItem[] = [];

        for (const assetTrendsDataItem of values(combinedDataMap)) {
            allAssetTrendsDataItems.push(assetTrendsDataItem);
        }

        sortCategoryTotalAmountItems(allAssetTrendsDataItems, transactionStatisticsFilter.value);

        const assetTrendsData: TransactionAssetTrendsAnalysisData = {
            items: allAssetTrendsDataItems
        };

        return assetTrendsData;
    });

    function createNewTransactionCategoricalOverviewAnalysisDataItem(id: string, name: string, type: TransactionCategoricalOverviewAnalysisDataItemType, displayOrders: number[], hidden: boolean): TransactionCategoricalOverviewAnalysisDataItem {
        const dataItem: TransactionCategoricalOverviewAnalysisDataItem = {
            id: id,
            name: name,
            type: type,
            displayOrders: displayOrders,
            hidden: hidden,
            inflows: [],
            outflows: [],
            totalAmount: 0,
            totalNonNegativeAmount: 0
        };

        return dataItem;
    }

    function sortCategoricalOverviewAnalysisDataItems(items: TransactionCategoricalOverviewAnalysisDataItem[], transactionStatisticsFilter: TransactionStatisticsFilter): void {
        let totalNonNegativeAmount: number = 0;

        for (const item of items) {
            totalNonNegativeAmount += item.totalNonNegativeAmount;
        }

        if (totalNonNegativeAmount > 0) {
            for (const item of items) {
                if (!item.includeInPercent) {
                    continue;
                }

                item.percent = item.totalAmount * 100 / totalNonNegativeAmount;
            }
        }

        sortStatisticsItems(items, transactionStatisticsFilter.sortingType);
    }

    function assembleAccountAndCategoryInfo(items: TransactionStatisticResponseItem[]): TransactionStatisticResponseItemWithInfo[] {
        const finalItems: TransactionStatisticResponseItemWithInfo[] = [];
        const defaultCurrency = userStore.currentUserDefaultCurrency;

        for (const dataItem of items) {
            const item: TransactionStatisticResponseItemWithInfo = {
                categoryId: dataItem.categoryId,
                accountId: dataItem.accountId,
                relatedAccountId: dataItem.relatedAccountId,
                relatedAccountType: dataItem.relatedAccountType,
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

            if (item.relatedAccountId) {
                item.relatedAccount = accountsStore.allAccountsMap[item.relatedAccountId];
            }

            if (item.relatedAccount && item.relatedAccount.parentId !== '0') {
                item.relatedPrimaryAccount = accountsStore.allAccountsMap[item.relatedAccount.parentId];
            } else {
                item.relatedPrimaryAccount = item.relatedAccount;
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
                    item.amountInDefaultCurrency = Math.trunc(amount);
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

        for (const item of items) {
            if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category) {
                continue;
            }

            if (transactionStatisticsFilter.chartDataType === ChartDataType.OutflowsByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalOutflows.type) {
                if (item.category.type === CategoryType.Transfer) {
                    if (item.relatedAccountType !== TransactionRelatedAccountType.TransferTo) {
                        continue;
                    }
                } else if (item.category.type !== CategoryType.Expense) {
                    continue;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalExpense.type) {
                if (item.category.type !== CategoryType.Expense) {
                    continue;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.InflowsByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalInflows.type) {
                if (item.category.type === CategoryType.Transfer) {
                    if (item.relatedAccountType !== TransactionRelatedAccountType.TransferFrom) {
                        continue;
                    }
                } else if (item.category.type !== CategoryType.Income) {
                    continue;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalIncome.type) {
                if (item.category.type !== CategoryType.Income) {
                    continue;
                }
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.NetCashFlow.type) {
                // Do Nothing
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.NetIncome.type) {
                if (item.category.type === CategoryType.Transfer) {
                    continue;
                }
            } else {
                continue;
            }

            if (transactionStatisticsFilter.filterAccountIds && transactionStatisticsFilter.filterAccountIds[item.account.id]) {
                continue;
            }

            if (transactionStatisticsFilter.filterCategoryIds && transactionStatisticsFilter.filterCategoryIds[item.category.id]) {
                continue;
            }

            if (transactionStatisticsFilter.chartDataType === ChartDataType.OutflowsByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.ExpenseByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.InflowsByAccount.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.IncomeByAccount.type) {
                if (isNumber(item.amountInDefaultCurrency)) {
                    let data = allDataItems[item.account.id];

                    if (data) {
                        data.totalAmount += item.amountInDefaultCurrency;
                    } else {
                        const primaryAccountCategoryDisplayOrder = settingsStore.accountCategoryDisplayOrders[item.primaryAccount.category] || Number.MAX_SAFE_INTEGER;

                        data = {
                            name: item.account.name,
                            type: 'account',
                            id: item.account.id,
                            icon: item.account.icon || DEFAULT_ACCOUNT_ICON.icon,
                            color: item.account.color || DEFAULT_ACCOUNT_COLOR,
                            hidden: item.primaryAccount.hidden || item.account.hidden,
                            displayOrders: [primaryAccountCategoryDisplayOrder, item.primaryAccount.displayOrder, item.account.displayOrder],
                            totalAmount: item.amountInDefaultCurrency
                        };
                    }

                    let includeInTotal: boolean = true;

                    // total outflows / inflows do not include transfer transactions between unfiltered accounts
                    if (transactionStatisticsFilter.chartDataType === ChartDataType.OutflowsByAccount.type ||
                        transactionStatisticsFilter.chartDataType === ChartDataType.InflowsByAccount.type) {
                        if (item.relatedAccount && (!transactionStatisticsFilter.filterAccountIds || !transactionStatisticsFilter.filterAccountIds[item.relatedAccount.id])) {
                            includeInTotal = false;
                        }
                    }

                    if (includeInTotal) {
                        totalAmount += item.amountInDefaultCurrency;

                        if (item.amountInDefaultCurrency > 0) {
                            totalNonNegativeAmount += item.amountInDefaultCurrency;
                        }
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
            } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalOutflows.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalExpense.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalInflows.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.TotalIncome.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.NetCashFlow.type ||
                transactionStatisticsFilter.chartDataType === ChartDataType.NetIncome.type) {
                if (isNumber(item.amountInDefaultCurrency)) {
                    let data = allDataItems['total'];
                    let amount = item.amountInDefaultCurrency;
                    let includeInTotal: boolean = true;

                    if (transactionStatisticsFilter.chartDataType === ChartDataType.NetCashFlow.type &&
                        (item.category.type === CategoryType.Expense || (item.category.type === CategoryType.Transfer && item.relatedAccountType === TransactionRelatedAccountType.TransferTo))) {
                        amount = -amount;
                    } else if (transactionStatisticsFilter.chartDataType === ChartDataType.NetIncome.type &&
                        item.category.type === CategoryType.Expense) {
                        amount = -amount;
                    }

                    // total outflows / inflows do not include transfer transactions between unfiltered accounts
                    if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalOutflows.type ||
                        transactionStatisticsFilter.chartDataType === ChartDataType.TotalInflows.type ||
                        transactionStatisticsFilter.chartDataType === ChartDataType.NetCashFlow.type) {
                        if (item.relatedAccount && (!transactionStatisticsFilter.filterAccountIds || !transactionStatisticsFilter.filterAccountIds[item.relatedAccount.id])) {
                            includeInTotal = false;
                        }
                    }

                    if (!data) {
                        let name = '';

                        if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalOutflows.type) {
                            name = ChartDataType.TotalOutflows.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalExpense.type) {
                            name = ChartDataType.TotalExpense.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalInflows.type) {
                            name = ChartDataType.TotalInflows.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.TotalIncome.type) {
                            name = ChartDataType.TotalIncome.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.NetCashFlow.type) {
                            name = ChartDataType.NetCashFlow.name;
                        } else if (transactionStatisticsFilter.chartDataType === ChartDataType.NetIncome.type) {
                            name = ChartDataType.NetIncome.name;
                        }

                        data = {
                            name: name,
                            type: 'total',
                            id: 'total',
                            icon: '',
                            color: '',
                            hidden: false,
                            displayOrders: [1],
                            totalAmount: 0
                        };
                    }

                    if (includeInTotal) {
                        data.totalAmount += amount;

                        totalAmount += amount;

                        if (item.amountInDefaultCurrency > 0) {
                            totalNonNegativeAmount += amount;
                        }
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
        transactionStatisticsFilter.value.assetTrendsChartType = TrendChartType.Default.type;
        transactionStatisticsFilter.value.assetTrendsChartDateType = DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE.type;
        transactionStatisticsFilter.value.assetTrendsChartStartTime = 0;
        transactionStatisticsFilter.value.assetTrendsChartEndTime = 0;
        transactionStatisticsFilter.value.filterAccountIds = {};
        transactionStatisticsFilter.value.filterCategoryIds = {};
        transactionStatisticsFilter.value.tagFilter = '';
        transactionStatisticsFilter.value.keyword = '';
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
        } else if (analysisType === StatisticsAnalysisType.AssetTrends) {
            if (!ChartDataType.isAvailableForAnalysisType(transactionStatisticsFilter.value.chartDataType, analysisType)) {
                transactionStatisticsFilter.value.chartDataType = ChartDataType.DefaultForAssetTrends.type;
            }
        }

        // Categorical Analysis filter initialization
        if (filter && isInteger(filter.categoricalChartType)) {
            transactionStatisticsFilter.value.categoricalChartType = filter.categoricalChartType;
        } else {
            transactionStatisticsFilter.value.categoricalChartType = settingsStore.appSettings.statistics.defaultCategoricalChartType;
        }

        if (!CategoricalChartType.isValidType(transactionStatisticsFilter.value.categoricalChartType)) {
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

        // Trend Analysis filter initialization
        if (filter && isInteger(filter.trendChartType)) {
            transactionStatisticsFilter.value.trendChartType = filter.trendChartType;
        } else {
            transactionStatisticsFilter.value.trendChartType = settingsStore.appSettings.statistics.defaultTrendChartType;
        }

        if (!TrendChartType.isValidType(transactionStatisticsFilter.value.trendChartType)) {
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
                transactionStatisticsFilter.value.trendChartStartYearMonth = getGregorianCalendarYearAndMonthFromUnixTime(trendChartDateRange.minTime);
                transactionStatisticsFilter.value.trendChartEndYearMonth = getGregorianCalendarYearAndMonthFromUnixTime(trendChartDateRange.maxTime);
            }
        }

        // Asset Trends filter initialization
        if (filter && isInteger(filter.assetTrendsChartType)) {
            transactionStatisticsFilter.value.assetTrendsChartType = filter.assetTrendsChartType;
        } else {
            transactionStatisticsFilter.value.assetTrendsChartType = settingsStore.appSettings.statistics.defaultAssetTrendsChartType;
        }

        if (!TrendChartType.isValidType(transactionStatisticsFilter.value.assetTrendsChartType)) {
            transactionStatisticsFilter.value.assetTrendsChartType = TrendChartType.Default.type;
        }

        if (filter && isInteger(filter.assetTrendsChartDateType)) {
            transactionStatisticsFilter.value.assetTrendsChartDateType = filter.assetTrendsChartDateType;
        } else {
            transactionStatisticsFilter.value.assetTrendsChartDateType = settingsStore.appSettings.statistics.defaultAssetTrendsChartDataRangeType;
        }

        let assetTrendsChartDateTypeValid = true;

        if (!DateRange.isAvailableForScene(transactionStatisticsFilter.value.assetTrendsChartDateType, DateRangeScene.AssetTrends)) {
            transactionStatisticsFilter.value.assetTrendsChartDateType = DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE.type;
            assetTrendsChartDateTypeValid = false;
        }

        if (assetTrendsChartDateTypeValid && transactionStatisticsFilter.value.assetTrendsChartDateType === DateRange.Custom.type) {
            if (filter && isInteger(filter.assetTrendsChartStartTime)) {
                transactionStatisticsFilter.value.assetTrendsChartStartTime = filter.assetTrendsChartStartTime;
            } else {
                transactionStatisticsFilter.value.assetTrendsChartStartTime = 0;
            }

            if (filter && isInteger(filter.assetTrendsChartEndTime)) {
                transactionStatisticsFilter.value.assetTrendsChartEndTime = filter.assetTrendsChartEndTime;
            } else {
                transactionStatisticsFilter.value.assetTrendsChartEndTime = 0;
            }
        } else {
            const assetTrendsChartDateRange = getDateRangeByDateType(transactionStatisticsFilter.value.assetTrendsChartDateType, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart);

            if (assetTrendsChartDateRange) {
                transactionStatisticsFilter.value.assetTrendsChartDateType = assetTrendsChartDateRange.dateType;
                transactionStatisticsFilter.value.assetTrendsChartStartTime = assetTrendsChartDateRange.minTime;
                transactionStatisticsFilter.value.assetTrendsChartEndTime = assetTrendsChartDateRange.maxTime;
            }
        }

        // Other filter initialization
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

        if (filter && isString(filter.tagFilter)) {
            transactionStatisticsFilter.value.tagFilter = filter.tagFilter;
        } else {
            transactionStatisticsFilter.value.tagFilter = '';
        }

        if (filter && isString(filter.keyword)) {
            transactionStatisticsFilter.value.keyword = filter.keyword;
        } else {
            transactionStatisticsFilter.value.keyword = '';
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

        // Categorical Analysis filter update
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

        // Trend Analysis filter update
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

        // Asset Trends filter update
        if (filter && isInteger(filter.assetTrendsChartType) && transactionStatisticsFilter.value.assetTrendsChartType !== filter.assetTrendsChartType) {
            transactionStatisticsFilter.value.assetTrendsChartType = filter.assetTrendsChartType;
            changed = true;
        }

        if (filter && isInteger(filter.assetTrendsChartDateType) && transactionStatisticsFilter.value.assetTrendsChartDateType !== filter.assetTrendsChartDateType) {
            transactionStatisticsFilter.value.assetTrendsChartDateType = filter.assetTrendsChartDateType;
            changed = true;
        }

        if (filter && isInteger(filter.assetTrendsChartStartTime) && transactionStatisticsFilter.value.assetTrendsChartStartTime !== filter.assetTrendsChartStartTime) {
            transactionStatisticsFilter.value.assetTrendsChartStartTime = filter.assetTrendsChartStartTime;
            changed = true;
        }

        if (filter && isInteger(filter.assetTrendsChartEndTime) && transactionStatisticsFilter.value.assetTrendsChartEndTime !== filter.assetTrendsChartEndTime) {
            transactionStatisticsFilter.value.assetTrendsChartEndTime = filter.assetTrendsChartEndTime;
            changed = true;
        }

        // Other filter update
        if (filter && isObject(filter.filterAccountIds) && !isEquals(transactionStatisticsFilter.value.filterAccountIds, filter.filterAccountIds)) {
            transactionStatisticsFilter.value.filterAccountIds = filter.filterAccountIds;
            changed = true;
        }

        if (filter && isObject(filter.filterCategoryIds) && !isEquals(transactionStatisticsFilter.value.filterCategoryIds, filter.filterCategoryIds)) {
            transactionStatisticsFilter.value.filterCategoryIds = filter.filterCategoryIds;
            changed = true;
        }

        if (filter && isString(filter.tagFilter) && transactionStatisticsFilter.value.tagFilter !== filter.tagFilter) {
            transactionStatisticsFilter.value.tagFilter = filter.tagFilter;
            changed = true;
        }

        if (filter && isString(filter.keyword) && transactionStatisticsFilter.value.keyword !== filter.keyword) {
            transactionStatisticsFilter.value.keyword = filter.keyword;
            changed = true;
        }

        if (filter && isInteger(filter.sortingType) && transactionStatisticsFilter.value.sortingType !== filter.sortingType) {
            transactionStatisticsFilter.value.sortingType = filter.sortingType;
            changed = true;
        }

        return changed;
    }

    function getTransactionStatisticsPageParams(analysisType: StatisticsAnalysisType, trendDateAggregationType: number, assetTrendsDateAggregationType: number): string {
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

            if (trendDateAggregationType !== ChartDateAggregationType.Default.type) {
                querys.push('trendDateAggregationType=' + trendDateAggregationType);
            }
        } else if (analysisType === StatisticsAnalysisType.AssetTrends) {
            querys.push('chartType=' + transactionStatisticsFilter.value.assetTrendsChartType);
            querys.push('chartDateType=' + transactionStatisticsFilter.value.assetTrendsChartDateType);

            if (transactionStatisticsFilter.value.assetTrendsChartDateType === DateRange.Custom.type) {
                querys.push('startTime=' + transactionStatisticsFilter.value.assetTrendsChartStartTime);
                querys.push('endTime=' + transactionStatisticsFilter.value.assetTrendsChartEndTime);
            }

            if (assetTrendsDateAggregationType !== ChartDateAggregationType.Default.type) {
                querys.push('assetTrendsDateAggregationType=' + assetTrendsDateAggregationType);
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

        if (transactionStatisticsFilter.value.tagFilter) {
            querys.push('tagFilter=' + transactionStatisticsFilter.value.tagFilter);
        }

        if (transactionStatisticsFilter.value.keyword) {
            querys.push('keyword=' + encodeURIComponent(transactionStatisticsFilter.value.keyword));
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

        if (itemId && transactionStatisticsFilter.value.chartDataType === ChartDataType.Overview.type) {
            const items = itemId.split('-');
            const sourceItems = (items[0] || '').split(':');
            const queryAccountIds: string[] = [];
            const queryCategoryIds: string[] = [];

            if (sourceItems.length === 2) {
                if (sourceItems[0] === 'account') {
                    queryAccountIds.push(sourceItems[1] as string);
                } else if (sourceItems[0] === 'category') {
                    queryCategoryIds.push(sourceItems[1] as string);
                }
            }

            if (items.length === 2) {
                const targetItems = (items[1] || '').split(':');

                if (targetItems.length === 2) {
                    if (targetItems[0] === 'account') {
                        queryAccountIds.push(targetItems[1] as string);
                    } else if (targetItems[0] === 'category') {
                        queryCategoryIds.push(targetItems[1] as string);
                    }
                }
            }

            if (queryAccountIds.length) {
                if (queryAccountIds.length === 2) {
                    querys.push('type=4');
                }

                querys.push('accountIds=' + queryAccountIds.join(','));
            } else {
                querys.push('accountIds=' + getFinalAccountIdsByFilteredAccountIds(accountsStore.allAccountsMap, transactionStatisticsFilter.value.filterAccountIds));
            }

            if (queryCategoryIds.length) {
                querys.push('categoryIds=' + queryCategoryIds.join(','));
            } else {
                querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, transactionStatisticsFilter.value.filterCategoryIds));
            }
        } else if (itemId && (transactionStatisticsFilter.value.chartDataType === ChartDataType.InflowsByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.OutflowsByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.AccountTotalLiabilities.type)
        ) {
            querys.push('accountIds=' + itemId);

            if ((analysisType === StatisticsAnalysisType.CategoricalAnalysis || analysisType === StatisticsAnalysisType.TrendAnalysis) && !isObjectEmpty(transactionStatisticsFilter.value.filterCategoryIds)) {
                querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, transactionStatisticsFilter.value.filterCategoryIds));
            }
        } else if (itemId && (transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
            transactionStatisticsFilter.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type)
        ) {
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

        if (analysisType === StatisticsAnalysisType.CategoricalAnalysis || analysisType === StatisticsAnalysisType.TrendAnalysis) {
            if (transactionStatisticsFilter.value.tagFilter) {
                querys.push('tagFilter=' + transactionStatisticsFilter.value.tagFilter);
            }

            if (transactionStatisticsFilter.value.keyword) {
                querys.push('keyword=' + encodeURIComponent(transactionStatisticsFilter.value.keyword));
            }
        }

        if (analysisType === StatisticsAnalysisType.CategoricalAnalysis
            && transactionStatisticsFilter.value.chartDataType !== ChartDataType.AccountTotalAssets.type
            && transactionStatisticsFilter.value.chartDataType !== ChartDataType.AccountTotalLiabilities.type) {
            querys.push('dateType=' + transactionStatisticsFilter.value.categoricalChartDateType);

            if (transactionStatisticsFilter.value.categoricalChartDateType === DateRange.Custom.type) {
                querys.push('minTime=' + transactionStatisticsFilter.value.categoricalChartStartTime);
                querys.push('maxTime=' + transactionStatisticsFilter.value.categoricalChartEndTime);
            }
        } else if ((analysisType === StatisticsAnalysisType.TrendAnalysis || analysisType === StatisticsAnalysisType.AssetTrends) && dateRange) {
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
                tagFilter: transactionStatisticsFilter.value.tagFilter,
                keyword: transactionStatisticsFilter.value.keyword,
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
                tagFilter: transactionStatisticsFilter.value.tagFilter,
                keyword: transactionStatisticsFilter.value.keyword,
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

    function loadAssetTrends({ force }: { force: boolean }): Promise<TransactionStatisticAssetTrendsResponseItem[]> {
        return new Promise((resolve, reject) => {
            services.getTransactionStatisticsAssetTrends({
                startTime: transactionStatisticsFilter.value.assetTrendsChartStartTime,
                endTime: transactionStatisticsFilter.value.assetTrendsChartEndTime
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve transaction statistics' });
                    return;
                }

                if (transactionStatisticsStateInvalid.value) {
                    updateTransactionStatisticsInvalidState(false);
                }

                if (force && data.result && isEquals(transactionAssetTrendsData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionAssetTrendsData.value = data.result;

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
        categoricalOverviewAnalysisData,
        categoricalAnalysisData,
        trendsAnalysisData,
        assetTrendsData,
        // functions
        updateTransactionStatisticsInvalidState,
        resetTransactionStatistics,
        initTransactionStatisticsFilter,
        updateTransactionStatisticsFilter,
        getTransactionStatisticsPageParams,
        getTransactionListPageParams,
        loadCategoricalAnalysis,
        loadTrendAnalysis,
        loadAssetTrends
    };
});
