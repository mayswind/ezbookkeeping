import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type LocalizedDateRange, DateRangeScene } from '@/core/datetime.ts';
import { StatisticsAnalysisType } from '@/core/statistics.ts';

import { getIncludedAccountsDisplayContent } from '@/lib/account.ts';
import { getIncludedTransactionCategoriesDisplayContent } from '@/lib/category.ts';

export function useStatisticsSettingPageBase() {
    const {
        tt,
        getAllDateRanges,
        getAllTimezoneTypesUsedForStatistics,
        getAllKeywordMatchModes,
        getAllCategoricalChartTypes,
        getAllTrendChartTypes,
        getAllStatisticsChartDataTypes,
        getAllStatisticsSortingTypes
    } = useI18n();

    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();

    const loadingAccounts = ref<boolean>(false);
    const loadingTransactionCategories = ref<boolean>(false);

    const allChartDataTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsChartDataTypes(StatisticsAnalysisType.CategoricalAnalysis));
    const allTimezoneTypesUsedForStatistics = computed<TypeAndDisplayName[]>(() => getAllTimezoneTypesUsedForStatistics());
    const allKeywordMatchModes = computed<TypeAndDisplayName[]>(() => getAllKeywordMatchModes());
    const allSortingTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsSortingTypes());
    const allCategoricalChartTypes = computed<TypeAndDisplayName[]>(() => getAllCategoricalChartTypes());
    const allCategoricalChartDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.Normal, {}));
    const allTrendChartTypes = computed<TypeAndDisplayName[]>(() => getAllTrendChartTypes());
    const allTrendChartDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.TrendAnalysis, {}));
    const allAssetTrendsChartDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.AssetTrends, {}));

    const defaultChartDataType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultChartDataType,
        set: (value: number) => settingsStore.setStatisticsDefaultChartDataType(value)
    });

    const defaultTimezoneType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultTimezoneType,
        set: (value: number) => settingsStore.setStatisticsDefaultTimezoneType(value)
    });

    const defaultKeywordMatchMode = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultKeywordMatchMode,
        set: (value: number) => settingsStore.setStatisticsDefaultKeywordMatchMode(value)
    });

    const defaultAccountFilterDisplayContent = computed<string>(() => {
        if (loadingAccounts.value) {
            return '';
        }

        const excludeAccountIds = settingsStore.appSettings.statistics.defaultAccountFilter;
        const displayContent = getIncludedAccountsDisplayContent(excludeAccountIds, accountsStore.allPlainAccounts, accountsStore.allAccountsMap);
        return displayContent ? tt(displayContent) : displayContent;
    });

    const defaultTransactionCategoryFilterDisplayContent = computed<string>(() => {
        if (loadingTransactionCategories.value) {
            return '';
        }

        const excludeAccountIds = settingsStore.appSettings.statistics.defaultTransactionCategoryFilter;
        const displayContent = getIncludedTransactionCategoriesDisplayContent(excludeAccountIds, transactionCategoriesStore.allTransactionCategoriesMap);
        return displayContent ? tt(displayContent) : displayContent;
    });

    const defaultSortingType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultSortingType,
        set: (value: number) => settingsStore.setStatisticsSortingType(value)
    });

    const defaultCategoricalChartType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultCategoricalChartType,
        set: (value: number) => settingsStore.setStatisticsDefaultCategoricalChartType(value)
    });

    const defaultCategoricalChartDateRange = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultCategoricalChartDataRangeType,
        set: (value: number) => settingsStore.setStatisticsDefaultCategoricalChartDateRange(value)
    });

    const defaultTrendChartType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultTrendChartType,
        set: (value: number) => settingsStore.setStatisticsDefaultTrendChartType(value)
    });

    const defaultTrendChartDateRange = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultTrendChartDataRangeType,
        set: (value: number) => settingsStore.setStatisticsDefaultTrendChartDateRange(value)
    });

    const defaultAssetTrendsChartType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultAssetTrendsChartType,
        set: (value: number) => settingsStore.setStatisticsDefaultAssetTrendsChartType(value)
    });

    const defaultAssetTrendsChartDateRange = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultAssetTrendsChartDataRangeType,
        set: (value: number) => settingsStore.setStatisticsDefaultAssetTrendsChartDateRange(value)
    });

    return {
        // states,
        loadingAccounts,
        loadingTransactionCategories,
        // computed states
        allChartDataTypes,
        allTimezoneTypesUsedForStatistics,
        allKeywordMatchModes,
        allSortingTypes,
        allCategoricalChartTypes,
        allCategoricalChartDateRanges,
        allTrendChartTypes,
        allTrendChartDateRanges,
        allAssetTrendsChartDateRanges,
        defaultChartDataType,
        defaultTimezoneType,
        defaultKeywordMatchMode,
        defaultAccountFilterDisplayContent,
        defaultTransactionCategoryFilterDisplayContent,
        defaultSortingType,
        defaultCategoricalChartType,
        defaultCategoricalChartDateRange,
        defaultTrendChartType,
        defaultTrendChartDateRange,
        defaultAssetTrendsChartType,
        defaultAssetTrendsChartDateRange
    };
}
