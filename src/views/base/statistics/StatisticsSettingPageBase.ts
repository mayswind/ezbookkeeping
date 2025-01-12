import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type LocalizedDateRange, DateRangeScene } from '@/core/datetime.ts';
import { StatisticsAnalysisType } from '@/core/statistics.ts';

export function useStatisticsSettingPageBase() {
    const {
        getAllDateRanges,
        getAllTimezoneTypesUsedForStatistics,
        getAllCategoricalChartTypes,
        getAllTrendChartTypes,
        getAllStatisticsChartDataTypes,
        getAllStatisticsSortingTypes
    } = useI18n();

    const settingsStore = useSettingsStore();

    const allChartDataTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsChartDataTypes(StatisticsAnalysisType.CategoricalAnalysis));
    const allTimezoneTypesUsedForStatistics = computed<TypeAndDisplayName[]>(() => getAllTimezoneTypesUsedForStatistics());
    const allSortingTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsSortingTypes());
    const allCategoricalChartTypes = computed<TypeAndDisplayName[]>(() => getAllCategoricalChartTypes());
    const allCategoricalChartDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.Normal, false));
    const allTrendChartTypes = computed<LocalizedDateRange[]>(() => getAllTrendChartTypes());
    const allTrendChartDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.TrendAnalysis, false));

    const defaultChartDataType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultChartDataType,
        set: (value: number) => settingsStore.setStatisticsDefaultChartDataType(value)
    });

    const defaultTimezoneType = computed<number>({
        get: () => settingsStore.appSettings.statistics.defaultTimezoneType,
        set: (value: number) => settingsStore.setStatisticsDefaultTimezoneType(value)
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

    return {
        // computed states
        allChartDataTypes,
        allTimezoneTypesUsedForStatistics,
        allSortingTypes,
        allCategoricalChartTypes,
        allCategoricalChartDateRanges,
        allTrendChartTypes,
        allTrendChartDateRanges,
        defaultChartDataType,
        defaultTimezoneType,
        defaultSortingType,
        defaultCategoricalChartType,
        defaultCategoricalChartDateRange,
        defaultTrendChartType,
        defaultTrendChartDateRange
    };
}
