import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { type TransactionStatisticsFilter, useStatisticsStore } from '@/stores/statistics.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type LocalizedDateRange, DateRangeScene, DateRange } from '@/core/datetime.ts';
import { StatisticsAnalysisType, ChartDataType, ChartSortingType, ChartDateAggregationType } from '@/core/statistics.ts';
import type { TransactionCategoricalAnalysisData, TransactionTrendsAnalysisData } from '@/models/transaction.ts';

import { limitText, findNameByType, findDisplayNameByType } from '@/lib/common.ts';
import { getYearMonthFirstUnixTime, getYearMonthLastUnixTime } from '@/lib/datetime.ts';

export function useStatisticsTransactionPageBase() {
    const {
        tt,
        getAllDateRanges,
        getAllStatisticsSortingTypes,
        getAllStatisticsDateAggregationTypes,
        formatUnixTimeToLongDateTime,
        formatUnixTimeToLongYearMonth,
        formatDateRange,
        formatAmountWithCurrency
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const statisticsStore = useStatisticsStore();

    const loading = ref<boolean>(true);
    const analysisType = ref<StatisticsAnalysisType>(StatisticsAnalysisType.CategoricalAnalysis);
    const trendDateAggregationType = ref<number>(ChartDateAggregationType.Default.type);

    const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);

    const allDateRanges = computed<LocalizedDateRange[]>(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return getAllDateRanges(DateRangeScene.Normal, true);
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return getAllDateRanges(DateRangeScene.TrendAnalysis, true);
        } else {
            return [];
        }
    });
    const allSortingTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsSortingTypes());
    const allDateAggregationTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsDateAggregationTypes());

    const query = computed<TransactionStatisticsFilter>(() => statisticsStore.transactionStatisticsFilter);
    const queryChartDataCategory = computed<string>(() => statisticsStore.categoricalAnalysisChartDataCategory);
    const queryDateType = computed<number | null>(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return query.value.categoricalChartDateType;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return query.value.trendChartDateType;
        } else {
            return null;
        }
    });

    const queryStartTime = computed<string>(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return formatUnixTimeToLongDateTime(query.value.categoricalChartStartTime);
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return formatUnixTimeToLongYearMonth(getYearMonthFirstUnixTime(query.value.trendChartStartYearMonth));
        } else {
            return '';
        }
    });

    const queryEndTime = computed<string>(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return formatUnixTimeToLongDateTime(query.value.categoricalChartEndTime);
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return formatUnixTimeToLongYearMonth(getYearMonthLastUnixTime(query.value.trendChartEndYearMonth));
        } else {
            return '';
        }
    });

    const queryDateRangeName = computed<string>(() => {
        if (query.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
            query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
            return tt(DateRange.All.name);
        }

        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return formatDateRange(query.value.categoricalChartDateType, query.value.categoricalChartStartTime, query.value.categoricalChartEndTime);
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return formatDateRange(query.value.trendChartDateType, getYearMonthFirstUnixTime(query.value.trendChartStartYearMonth), getYearMonthLastUnixTime(query.value.trendChartEndYearMonth));
        } else {
            return '';
        }
    });

    const queryChartDataTypeName = computed<string>(() => {
        const queryChartDataTypeName = findNameByType(ChartDataType.values(), query.value.chartDataType) || 'Statistics';
        return tt(queryChartDataTypeName);
    });

    const querySortingTypeName = computed<string>(() => {
        const querySortingTypeName = findNameByType(ChartSortingType.values(), query.value.sortingType) || 'System Default';
        return tt(querySortingTypeName);
    });

    const queryTrendDateAggregationTypeName = computed<string>(() => findDisplayNameByType(allDateAggregationTypes.value, trendDateAggregationType.value) || '');

    const isQueryDateRangeChanged = computed<boolean>(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            if (query.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
                query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                return false;
            }

            if (query.value.categoricalChartDateType === settingsStore.appSettings.statistics.defaultCategoricalChartDataRangeType) {
                return false;
            }

            return !!query.value.categoricalChartStartTime || !!query.value.categoricalChartEndTime;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            if (query.value.trendChartDateType === settingsStore.appSettings.statistics.defaultTrendChartDataRangeType) {
                return false;
            }

            return !!query.value.trendChartStartYearMonth || !!query.value.trendChartEndYearMonth;
        } else {
            return false;
        }
    });

    const canShiftDateRange = computed<boolean>(() => {
        if (query.value.chartDataType === ChartDataType.AccountTotalAssets.type || query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
            return false;
        }

        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return query.value.categoricalChartDateType !== DateRange.All.type;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return query.value.trendChartDateType !== DateRange.All.type;
        } else {
            return false;
        }
    });

    const showCustomDateRange = computed<boolean>(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return query.value.categoricalChartDateType === DateRange.Custom.type && !!query.value.categoricalChartStartTime && !!query.value.categoricalChartEndTime;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return query.value.trendChartDateType === DateRange.Custom.type && !!query.value.trendChartStartYearMonth && !!query.value.trendChartEndYearMonth;
        } else {
            return false;
        }
    });

    const showAmountInChart = computed<boolean>(() => {
        if (!showAccountBalance.value
            && (query.value.chartDataType === ChartDataType.AccountTotalAssets.type || query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type)) {
            return false;
        }

        return true;
    });

    const totalAmountName = computed<string>(() => {
        if (query.value.chartDataType === ChartDataType.IncomeByAccount.type
            || query.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type
            || query.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
            return tt('Total Income');
        } else if (query.value.chartDataType === ChartDataType.ExpenseByAccount.type
            || query.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type
            || query.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type) {
            return tt('Total Expense');
        } else if (query.value.chartDataType === ChartDataType.AccountTotalAssets.type) {
            return tt('Total Assets');
        } else if (query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
            return tt('Total Liabilities');
        }

        return tt('Total Amount');
    });

    const showTotalAmountInTrendsChart = computed<boolean>(() => {
        return query.value.chartDataType !== ChartDataType.TotalExpense.type &&
            query.value.chartDataType !== ChartDataType.TotalIncome.type &&
            query.value.chartDataType !== ChartDataType.TotalBalance.type;
    });

    const translateNameInTrendsChart = computed<boolean>(() => {
        return query.value.chartDataType === ChartDataType.TotalExpense.type ||
            query.value.chartDataType === ChartDataType.TotalIncome.type ||
            query.value.chartDataType === ChartDataType.TotalBalance.type;
    });

    const categoricalAnalysisData = computed<TransactionCategoricalAnalysisData>(() => statisticsStore.categoricalAnalysisData);
    const trendsAnalysisData = computed<TransactionTrendsAnalysisData | null>(() => statisticsStore.trendsAnalysisData);

    function getDisplayAmount(amount: number, currency: string, textLimit?: number): string {
        const finalAmount = formatAmountWithCurrency(amount, currency);

        if (!showAccountBalance.value
            && (query.value.chartDataType === ChartDataType.AccountTotalAssets.type
                || query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type)
        ) {
            return '***';
        }

        if (textLimit) {
            return limitText(finalAmount, textLimit);
        }

        return finalAmount;
    }

    return {
        // states
        loading,
        analysisType,
        trendDateAggregationType,
        // computed states
        showAccountBalance,
        defaultCurrency,
        firstDayOfWeek,
        fiscalYearStart,
        allDateRanges,
        allSortingTypes,
        allDateAggregationTypes,
        query,
        queryChartDataCategory,
        queryDateType,
        queryStartTime,
        queryEndTime,
        queryDateRangeName,
        queryChartDataTypeName,
        querySortingTypeName,
        queryTrendDateAggregationTypeName,
        isQueryDateRangeChanged,
        canShiftDateRange,
        showCustomDateRange,
        showAmountInChart,
        totalAmountName,
        showTotalAmountInTrendsChart,
        translateNameInTrendsChart,
        categoricalAnalysisData,
        trendsAnalysisData,
        // functions
        getDisplayAmount
    };
}
