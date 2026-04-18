<template>
    <v-card-text class="px-5 py-0 mb-4">
        <v-row>
            <v-col cols="12">
                <div class="d-flex overflow-x-auto align-center gap-2 pt-2">
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || disabled"
                        :label="tt('Chart Type')"
                        :items="allTransactionExplorerChartTypes"
                        v-model="currentExplorer.chartType"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || disabled"
                        :label="tt('Axis / Category')"
                        :items="allTransactionExplorerDataDimensions"
                        :model-value="currentExplorer.categoryDimension"
                        @update:model-value="updateCategoryDimensionType"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || disabled || !TransactionExplorerChartType.valueOf(currentExplorer.chartType)?.seriesDimensionRequired"
                        :label="tt('Series')"
                        :items="allTransactionExplorerDataDimensions"
                        :model-value="TransactionExplorerChartType.valueOf(currentExplorer.chartType)?.seriesDimensionRequired ? currentExplorer.seriesDimension : TransactionExplorerDataDimension.None.value"
                        @update:model-value="currentExplorer.seriesDimension = $event as TransactionExplorerDataDimensionType"
                    >
                        <template #item="{ props, item }">
                            <v-list-item :disabled="item.value === currentExplorer.categoryDimension && item.value !== TransactionExplorerDataDimension.SeriesDimensionDefault.value" v-bind="props">
                                <template #title>
                                    <div class="text-truncate">{{ item.raw.name }}</div>
                                </template>
                            </v-list-item>
                        </template>
                    </v-select>
                    <v-select
                        class="flex-0-0"
                        min-width="220"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || disabled"
                        :label="tt('Number of Amount Ranges')"
                        :items="allAmountRangeCounts"
                        v-model="currentExplorer.amountRangeCount"
                        v-if="isUsingAmountRange"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || disabled"
                        :label="tt('Value Metric')"
                        :items="allTransactionExplorerValueMetrics"
                        v-model="currentExplorer.valueMetric"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="displayName"
                        item-value="type"
                        density="compact"
                        :disabled="loading || disabled"
                        :label="tt('Sort Order')"
                        :items="allTransactionExplorerChartSortingTypes"
                        v-model="currentExplorer.chartSortingType"
                    />
                    <v-spacer class="flex-1-1"/>
                </div>
            </v-col>
        </v-row>
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-if="currentExplorer.chartType === TransactionExplorerChartType.Pie.value">
        <pie-chart
            :items="[
                {id: '1', name: '---', value: 60, color: '7c7c7f'},
                {id: '2', name: '---', value: 20, color: 'a5a5aa'},
                {id: '3', name: '---', value: 20, color: 'c5c5c9'}
            ]"
            :skeleton="true"
            id-field="id"
            name-field="name"
            value-field="value"
            color-field="color"
            v-if="loading"
        />
        <pie-chart
            :items="categoryDimensionTransactionExplorerData && categoryDimensionTransactionExplorerData.length ? categoryDimensionTransactionExplorerData : []"
            :show-value="true"
            :show-percent="true"
            :enable-click-item="true"
            :amount-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isAmount"
            :percent-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isPercent"
            :default-currency="defaultCurrency"
            id-field="id"
            name-field="name"
            value-field="totalAmount"
            v-else-if="!loading"
            @click="onClickPieChartItem"
        />
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-else-if="currentExplorer.chartType === TransactionExplorerChartType.Radar.value">
        <radar-chart
            :items="[
                {name: '---', value: 10},
                {name: '---', value: 10},
                {name: '---', value: 10},
                {name: '---', value: 10},
                {name: '---', value: 10},
                {name: '---', value: 10}
            ]"
            :skeleton="true"
            name-field="name"
            value-field="value"
            v-if="loading"
        />
        <radar-chart
            :items="categoryDimensionTransactionExplorerData && categoryDimensionTransactionExplorerData.length ? categoryDimensionTransactionExplorerData : []"
            :show-value="true"
            :show-percent="true"
            :amount-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isAmount"
            :percent-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isPercent"
            :default-currency="defaultCurrency"
            name-field="name"
            value-field="totalAmount"
            v-else-if="!loading"
        />
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-else-if="TransactionExplorerChartType.valueOf(currentExplorer.chartType)?.seriesDimensionRequired && axisChartDisplayType">
        <axis-chart
            :skeleton="true"
            :type="axisChartDisplayType"
            :sorting-type="currentExplorer.chartSortingType"
            :all-category-names="[]"
            :items="[]"
            category-type-name=""
            name-field="name"
            values-field="values"
            v-if="loading"
        />
        <axis-chart
            ref="axisChart"
            :type="axisChartDisplayType"
            :stacked="axisChartStacked"
            :one-hundred-percent-stacked="axisChart100PercentStacked"
            :sorting-type="currentExplorer.chartSortingType"
            :show-value="true"
            :show-total-amount-in-tooltip="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.supportSum"
            :total-name-in-tooltip="currentExplorer.valueMetric === TransactionExplorerValueMetric.TransactionCount.value ? tt('Total Transactions') : tt('Total Amount')"
            :category-type-name="currentTransactionExplorerCategoryDimensionName"
            :all-category-names="categoriedNamesSortedByDisplayOrder"
            :items="seriesDimensionTransactionExplorerData"
            :amount-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isAmount"
            :percent-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isPercent"
            :default-currency="defaultCurrency"
            :enable-click-item="true"
            :tooltip-extra-column-names="axisChartTooltipExtraColumnNames"
            :tooltip-extra-column-total-values="axisChartShowYearOverYear || axisChartShowPeriodOverPeriod ? getAxisChartTooltipExtraColumnTotalValues : undefined"
            :tooltip-extra-column-values="axisChartShowYearOverYear || axisChartShowPeriodOverPeriod ? getAxisChartTooltipExtraColumnValues : undefined"
            id-field="id"
            name-field="name"
            values-field="categoryValues"
            display-orders-field="displayOrders"
            @click="onClickTrendChartItem"
            v-else-if="!loading"
        />
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-else-if="currentExplorer.chartType === TransactionExplorerChartType.Heatmap.value">
        <heat-map-chart
            :skeleton="true"
            :all-category-names="[]"
            :items="[]"
            :value-type-name="tt(TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.name ?? 'Value')"
            name-field="name"
            values-field="values"
            v-if="loading"
        />
        <heat-map-chart
            :show-value="true"
            :all-category-names="categoriedNamesSortedByDisplayOrder"
            :items="seriesDimensionTransactionExplorerData"
            :value-type-name="tt(TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.name ?? 'Value')"
            :amount-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isAmount"
            :percent-value="TransactionExplorerValueMetric.valueOf(currentExplorer.valueMetric)?.isPercent"
            :default-currency="defaultCurrency"
            name-field="name"
            values-field="categoryValues"
            v-else-if="!loading"
        />
    </v-card-text>
</template>

<script setup lang="ts">
import AxisChart, { type AxisChartDisplayType } from '@/components/desktop/AxisChart.vue';

import { computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import {
    type CategoriedInfo,
    type CategoriedTransactionExplorerData,
    type SeriesInfo,
    type CategoriedTransactionExplorerDataItem,
    TransactionExplorerDimensionType,
    useExplorersStore
} from '@/stores/explorer.ts';

import { type NameValue, type NameNumeralValue, type TypeAndDisplayName, itemAndIndex, entries } from '@/core/base.ts';
import { NumeralSystem } from '@/core/numeral.ts';
import { Month, WeekDay } from '@/core/datetime.ts';
import { ChartSortingType, ExportMermaidChartType } from '@/core/statistics.ts';
import {
    TransactionExplorerChartType,
    TransactionExplorerDataDimensionType,
    TransactionExplorerDataDimension,
    TransactionExplorerValueMetric
} from '@/core/explorer.ts';

import { type SortableTransactionStatisticDataItem } from '@/models/transaction.ts';
import type { InsightsExplorer } from '@/models/explorer.ts';

import { isDefined, isNumber, findNameByValue } from '@/lib/common.ts';
import { getCurrentDateTime, parseDateTimeFromString } from '@/lib/datetime.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';

type AxisChartType = InstanceType<typeof AxisChart>;

interface InsightsExplorerDataTableTabProps {
    loading?: boolean;
    disabled?: boolean;
}

interface CategoryDimensionData extends SortableTransactionStatisticDataItem {
    id: string;
    dimension: TransactionExplorerDimensionType;
    name: string;
    displayOrders: number[];
    totalAmount: number;
}

interface SortableCategoriedTransactionExplorerDataItem extends SortableTransactionStatisticDataItem {
    name: string;
    displayOrders: number[];
    totalAmount: number;
    originalItem: CategoriedTransactionExplorerData;
}

interface SeriesDimensionData extends SortableTransactionStatisticDataItem, Record<string, unknown> {
    id: string;
    dimension: TransactionExplorerDimensionType;
    name: string;
    displayOrders: number[];
    categoryValues: number[];
    totalAmount: number;
}

defineProps<InsightsExplorerDataTableTabProps>();

const router = useRouter();

const {
    tt,
    getAllStatisticsSortingTypes,
    getAllTransactionExplorerDataDimensions,
    getAllTransactionExplorerValueMetrics,
    getAllTransactionExplorerChartTypes,
    getMonthLongName,
    getMonthdayShortName,
    getWeekdayLongName,
    getQuarterName,
    getCurrentNumeralSystemType,
    getCurrencyName,
    formatDateTimeToShortDateTime,
    formatDateTimeToShortDate,
    formatDateTimeToShortTime,
    formatDateTimeToGregorianLikeShortYear,
    formatDateTimeToGregorianLikeShortYearMonth,
    formatDateTimeToGregorianLikeYearQuarter,
    formatGregorianYearToGregorianLikeFiscalYear,
    formatAmountToLocalizedNumerals,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatPercentToLocalizedNumerals
} = useI18n();

const userStore = useUserStore();
const explorersStore = useExplorersStore();

const axisChart = useTemplateRef<AxisChartType>('axisChart');

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const allTransactionExplorerDataDimensions = computed<NameValue[]>(() => getAllTransactionExplorerDataDimensions());
const allTransactionExplorerValueMetrics = computed<NameValue[]>(() => getAllTransactionExplorerValueMetrics());
const allTransactionExplorerChartTypes = computed<NameValue[]>(() => getAllTransactionExplorerChartTypes());
const allTransactionExplorerChartSortingTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsSortingTypes(true));
const currentTransactionExplorerCategoryDimensionName = computed<string>(() => findNameByValue(allTransactionExplorerDataDimensions.value, currentExplorer.value.categoryDimension) ?? tt('Unknown'));

const currentExplorer = computed<InsightsExplorer>(() => explorersStore.currentInsightsExplorer);
const isUsingAmountRange = computed<boolean>(() => explorersStore.isUsingAmountRange);

const allAmountRangeCounts = computed<NameNumeralValue[]>(() => {
    const pageCounts: NameNumeralValue[] = [];

    for (let i = 3; i <= 20; i++) {
        pageCounts.push({ value: i, name: numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(i.toString()) });
    }

    return pageCounts;
});

const categoryDimensionTransactionExplorerData = computed<CategoryDimensionData[]>(() => {
    if (currentExplorer.value.chartType !== TransactionExplorerChartType.Pie.value && currentExplorer.value.chartType !== TransactionExplorerChartType.Radar.value) {
        return [];
    }

    if (!explorersStore.categoriedTransactionExplorerData || !explorersStore.categoriedTransactionExplorerData.length) {
        return [];
    }

    const result: CategoryDimensionData[] = [];

    for (const categoriedData of explorersStore.categoriedTransactionExplorerData) {
        const data = categoriedData.data[0];

        if (!isDefined(data)) {
            continue;
        }

        const displayName = getCategoriedDataDisplayName(categoriedData);

        result.push({
            id: categoriedData.categoryId,
            dimension: categoriedData.categoryIdType,
            name: displayName,
            displayOrders: categoriedData.categoryDisplayOrders,
            totalAmount: data.value
        });
    }

    sortStatisticsItems(result, currentExplorer.value.chartSortingType);

    return result;
});

const categoriedDataSortedByDisplayOrder = computed<SortableCategoriedTransactionExplorerDataItem[]>(() => {
    if (!explorersStore.categoriedTransactionExplorerData || !explorersStore.categoriedTransactionExplorerData.length) {
        return [];
    }

    const result: SortableCategoriedTransactionExplorerDataItem[] = [];

    for (const categoriedData of explorersStore.categoriedTransactionExplorerData) {
        result.push({
            name: getCategoriedDataDisplayName(categoriedData),
            displayOrders: categoriedData.categoryDisplayOrders,
            totalAmount: 0,
            originalItem: categoriedData
        });
    }

    sortStatisticsItems(result, ChartSortingType.DisplayOrder.type);

    return result;
});

const categoriedNamesSortedByDisplayOrder = computed<string[]>(() => {
    if (!categoriedDataSortedByDisplayOrder.value || !categoriedDataSortedByDisplayOrder.value.length) {
        return [];
    }

    const result: string[] = [];

    for (const categoriedData of categoriedDataSortedByDisplayOrder.value) {
        result.push(categoriedData.name);
    }

    return result;
});

const seriesDimensionTransactionExplorerData = computed<SeriesDimensionData[]>(() => {
    if (!TransactionExplorerChartType.valueOf(currentExplorer.value.chartType)?.seriesDimensionRequired) {
        return [];
    }

    if (!explorersStore.categoriedTransactionExplorerData || !explorersStore.categoriedTransactionExplorerData.length) {
        return [];
    }

    const result: SeriesDimensionData[] = [];
    const seriesDimensionDataMap: Record<string, SeriesDimensionData> = {};

    for (const categoriedData of explorersStore.categoriedTransactionExplorerData) {
        for (const seriesData of categoriedData.data) {
            const displayName = getCategoriedDataDisplayName(seriesData);
            let seriesDimensionData: SeriesDimensionData | undefined = seriesDimensionDataMap[seriesData.seriesId];

            if (!seriesDimensionData) {
                seriesDimensionData = {
                    id: seriesData.seriesId,
                    dimension: seriesData.seriesIdType,
                    name: displayName,
                    displayOrders: seriesData.seriesDisplayOrders,
                    categoryValues: [],
                    totalAmount: 0
                };

                seriesDimensionDataMap[seriesData.seriesId] = seriesDimensionData;
                result.push(seriesDimensionData);
            }
        }
    }

    for (const categoriedData of categoriedDataSortedByDisplayOrder.value) {
        const seriesDataMap: Record<string, CategoriedTransactionExplorerDataItem> = {};

        for (const seriesData of categoriedData.originalItem.data) {
            seriesDataMap[seriesData.seriesId] = seriesData;
        }

        for (const seriesDimensionData of result) {
            const seriesData = seriesDataMap[seriesDimensionData.id];

            if (isDefined(seriesData)) {
                seriesDimensionData.categoryValues.push(seriesData.value);
                seriesDimensionData.totalAmount += seriesData.value;
            } else {
                seriesDimensionData.categoryValues.push(0);
            }
        }
    }

    sortStatisticsItems(result, currentExplorer.value.chartSortingType);

    return result;
});

const seriesDimensionTransactionExplorerDataMap = computed<Record<string, SeriesDimensionData>>(() => {
    const result: Record<string, SeriesDimensionData> = {};

    for (const seriesDimensionData of seriesDimensionTransactionExplorerData.value) {
        result[seriesDimensionData.id] = seriesDimensionData;
    }

    return result;
});

const axisChartCategoryIndexYoYMap = computed<Record<number, number>>(() => {
    const result: Record<number, number> = {};

    if (!axisChartShowYearOverYear.value) {
        return result;
    }

    const dateKeyToIndex: Record<string, number> = {};
    const dateKeyToPreviousYearDateKey: Record<string, string> = {};

    for (const [item, categoryIndex] of itemAndIndex(categoriedDataSortedByDisplayOrder.value)) {
        const categoriedData = item.originalItem;
        const name = categoriedData.categoryId;
        const dimessionType = categoriedData.categoryIdType;
        const dimension = currentExplorer.value.categoryDimension;

        if (dimension === TransactionExplorerDataDimension.DateTimeByYearMonthDay.value) {
            const dateTime = parseDateTimeFromString(name, dimessionType);

            if (dateTime) {
                dateKeyToIndex[dateTime.getGregorianCalendarYearDashMonthDashDay()] = categoryIndex;
                dateKeyToPreviousYearDateKey[dateTime.getGregorianCalendarYearDashMonthDashDay()] = dateTime.add(-1, 'years').getGregorianCalendarYearDashMonthDashDay();
            }
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYearMonth.value) {
            const dateTime = parseDateTimeFromString(name, dimessionType);

            if (dateTime) {
                dateKeyToIndex[dateTime.getGregorianCalendarYearDashMonth()] = categoryIndex;
                dateKeyToPreviousYearDateKey[dateTime.getGregorianCalendarYearDashMonth()] = dateTime.add(-1, 'years').getGregorianCalendarYearDashMonth();
            }
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYearQuarter.value) {
            const parts = name.split('-');
            const year = parts.length === 2 ? parseInt(parts[0] as string) : 0;
            const quarter = parts.length === 2 ? parseInt(parts[1] as string) : 0;

            dateKeyToIndex[`${year}-Q${quarter}`] = categoryIndex;
            dateKeyToPreviousYearDateKey[`${year}-Q${quarter}`] = `${year - 1}-Q${quarter}`;
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYear.value) {
            const year = parseInt(name);
            dateKeyToIndex[name] = categoryIndex;
            dateKeyToPreviousYearDateKey[name] = (year - 1).toString(10);
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByFiscalYear.value) {
            const year = parseInt(name);
            dateKeyToIndex[name] = categoryIndex;
            dateKeyToPreviousYearDateKey[name] = (year - 1).toString(10);
        }
    }

    for (const [dateKey, previousYearDateKey] of entries(dateKeyToPreviousYearDateKey)) {
        const categoryIndex = dateKeyToIndex[dateKey];
        const previousYearCategoryIndex = dateKeyToIndex[previousYearDateKey];

        if (isNumber(categoryIndex) && isNumber(previousYearCategoryIndex)) {
            result[categoryIndex] = previousYearCategoryIndex;
        }
    }

    return result;
});

const axisChartDisplayType = computed<AxisChartDisplayType | undefined>(() => {
    if (currentExplorer.value.chartType === TransactionExplorerChartType.ColumnStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.Column100PercentStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.ColumnGrouped.value) {
        return 'column';
    } else if (currentExplorer.value.chartType === TransactionExplorerChartType.LineGrouped.value) {
        return 'line';
    } else if (currentExplorer.value.chartType === TransactionExplorerChartType.AreaStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.Area100PercentStacked.value) {
        return 'area';
    } else if (currentExplorer.value.chartType === TransactionExplorerChartType.BubbleGrouped.value) {
        return 'bubble';
    } else {
        return undefined;
    }
});

const axisChartStacked = computed<boolean>(() => {
    return (currentExplorer.value.chartType === TransactionExplorerChartType.ColumnStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.Column100PercentStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.AreaStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.Area100PercentStacked.value);
});

const axisChart100PercentStacked = computed<boolean>(() => {
    return (currentExplorer.value.chartType === TransactionExplorerChartType.Column100PercentStacked.value
        || currentExplorer.value.chartType === TransactionExplorerChartType.Area100PercentStacked.value);
});

const axisChartShowYearOverYear = computed<boolean>(() => {
    const dimession = currentExplorer.value.categoryDimension;

    return dimession === TransactionExplorerDataDimension.DateTimeByYearMonthDay.value
        || dimession === TransactionExplorerDataDimension.DateTimeByYearMonth.value
        || dimession === TransactionExplorerDataDimension.DateTimeByYearQuarter.value
        || dimession === TransactionExplorerDataDimension.DateTimeByYear.value
        || dimession === TransactionExplorerDataDimension.DateTimeByFiscalYear.value;
});

const axisChartShowPeriodOverPeriod = computed<boolean>(() => {
    const dimession = currentExplorer.value.categoryDimension;

    return dimession === TransactionExplorerDataDimension.DateTimeByYearMonthDay.value
        || dimession === TransactionExplorerDataDimension.DateTimeByYearMonth.value
        || dimession === TransactionExplorerDataDimension.DateTimeByYearQuarter.value;
});

const axisChartTooltipExtraColumnNames = computed<string[]>(() => {
    const extraColumnNames: string[] = [];
    const dimession = currentExplorer.value.categoryDimension;

    if (axisChartShowYearOverYear.value) {
        extraColumnNames.push(tt('Year-over-Year'));
    }

    if (axisChartShowPeriodOverPeriod.value) {
        if (dimession === TransactionExplorerDataDimension.DateTimeByYearQuarter.value) {
            extraColumnNames.push(tt('Quarter-over-Quarter'));
        } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearMonth.value) {
            extraColumnNames.push(tt('Month-over-Month'));
        } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearMonthDay.value) {
            extraColumnNames.push(tt('Day-over-Day'));
        } else {
            extraColumnNames.push(tt('Period-over-Period'));
        }
    }

    return extraColumnNames;
});

function getCategoriedDataDisplayName(info: CategoriedInfo | SeriesInfo): string {
    let name: string = '';
    let needI18n: boolean | undefined = false;
    let i18nParameters: Record<string, unknown> | undefined = undefined;
    let dimessionType: TransactionExplorerDimensionType = TransactionExplorerDimensionType.Other;
    let dimessionValue: TransactionExplorerDataDimensionType = TransactionExplorerDataDimension.None.value;

    if ('categoryName' in info) {
        name = info.categoryName;
        needI18n = info.categoryNameNeedI18n;
        i18nParameters = info.categoryNameI18nParameters;
        dimessionType = info.categoryIdType;
        dimessionValue = currentExplorer.value.categoryDimension;
    } else if ('seriesName' in info) {
        name = info.seriesName;
        needI18n = info.seriesNameNeedI18n;
        i18nParameters = info.seriesNameI18nParameters;
        dimessionType = info.seriesIdType;
        dimessionValue = currentExplorer.value.seriesDimension;
    }

    const dimession = TransactionExplorerDataDimension.valueOf(dimessionValue);
    let displayName: string = name;

    // convert the name to i18n if needed
    if (needI18n && i18nParameters) {
        displayName = tt(name, i18nParameters);
    } else if (needI18n && !i18nParameters) {
        displayName = tt(name);
    }

    // convert the name to formatted date time if needed
    if (dimession === TransactionExplorerDataDimension.DateTime) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToShortDateTime(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearMonthDay) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToShortDate(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearMonth) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToGregorianLikeShortYearMonth(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearQuarter) {
        const parts = name.split('-');
        const year = parts.length === 2 ? parts[0] : '';
        const quarter = parts.length === 2 ? parseInt(parts[1] as string) : 0;
        const quarterLastMonth = quarter * 3;
        const dateTime = year && quarterLastMonth ? parseDateTimeFromString(`${year}-${quarterLastMonth.toString(10).padStart(2, NumeralSystem.WesternArabicNumerals.digitZero)}`, TransactionExplorerDimensionType.YearMonth) : undefined;
        displayName = dateTime ? formatDateTimeToGregorianLikeYearQuarter(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYear) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToGregorianLikeShortYear(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByFiscalYear) {
        displayName = formatGregorianYearToGregorianLikeFiscalYear(parseInt(name));
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByDayOfWeek) {
        const weekDay = WeekDay.parse(name);
        displayName = weekDay ? getWeekdayLongName(weekDay) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByDayOfMonth) {
        displayName = getMonthdayShortName(parseInt(name));
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByMonthOfYear) {
        const month = Month.valueOf(parseInt(name));
        displayName = month ? getMonthLongName(month.name) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByQuarterOfYear) {
        displayName = getQuarterName(parseInt(name));
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByHourOfDay) {
        const dateTime = getCurrentDateTime().set({
            hour: parseInt(name),
            minute: 0,
            second: 0,
            millisecond: 0
        });
        displayName = formatDateTimeToShortTime(dateTime);
    } else if (dimession === TransactionExplorerDataDimension.SourceAccountCurrency || dimession === TransactionExplorerDataDimension.DestinationAccountCurrency) {
        if (!needI18n) {
            displayName = getCurrencyName(name);
        }
    }

    if (dimession === TransactionExplorerDataDimension.SourceAmount
        || dimession === TransactionExplorerDataDimension.DestinationAmount) {
        if (name !== '' && name !== 'none' && Number.isFinite(parseInt(name))) {
            displayName = formatAmountToLocalizedNumerals(parseInt(name), defaultCurrency.value);
        }
    }

    if (dimession?.isSourceAmountRange || dimession?.isDestinationAmountRange) {
        const rangeParts = name.split('|');

        if (rangeParts && rangeParts.length === 2 && Number.isFinite(parseInt(rangeParts[0] as string)) && Number.isFinite(parseInt(rangeParts[1] as string))) {
            const from = formatAmountToLocalizedNumerals(parseInt(rangeParts[0] as string), defaultCurrency.value);
            const to = formatAmountToLocalizedNumerals(parseInt(rangeParts[1] as string), defaultCurrency.value);
            displayName = `${from} ~ ${to}`;
        }
    }

    return displayName;
}

function formatDisplayChangeRate(current: number, reference: number): string {
    if (reference === 0 && current === 0) {
        return formatPercentToLocalizedNumerals(0, 2, '<0.01');
    }

    if (reference === 0) {
        return '-';
    }

    const rate = (current - reference) / reference * 100;
    return formatPercentToLocalizedNumerals(rate, 2, '<0.01');
}

function getAxisChartTooltipExtraColumnTotalValues(categoryIndex: number, totalValue: number, visibleSeriesIds: string[]): string[] {
    const extraColumnValues: string[] = [];

    if (!axisChartShowYearOverYear.value && !axisChartShowPeriodOverPeriod.value) {
        return extraColumnValues;
    }

    if (axisChartShowYearOverYear.value) {
        const yoyReferenceIndex = axisChartCategoryIndexYoYMap.value[categoryIndex];
        let displayChangeRate = '-';

        if (isNumber(yoyReferenceIndex)) {
            let referenceTotalValue = 0;

            for (const seriesId of visibleSeriesIds) {
                const seriesDimensionData = seriesDimensionTransactionExplorerDataMap.value[seriesId];

                if (seriesDimensionData && seriesDimensionData.categoryValues) {
                    referenceTotalValue += seriesDimensionData.categoryValues[yoyReferenceIndex] ?? 0;
                }
            }

            displayChangeRate = formatDisplayChangeRate(totalValue, referenceTotalValue);
        }

        extraColumnValues.push(displayChangeRate);
    }

    if (axisChartShowPeriodOverPeriod.value) {
        const popReferenceIndex = categoryIndex - 1;
        let displayChangeRate = '-';

        if (popReferenceIndex >= 0) {
            let referenceTotalValue = 0;

            for (const seriesId of visibleSeriesIds) {
                const seriesDimensionData = seriesDimensionTransactionExplorerDataMap.value[seriesId];

                if (seriesDimensionData && seriesDimensionData.categoryValues) {
                    referenceTotalValue += seriesDimensionData.categoryValues[popReferenceIndex] ?? 0;
                }
            }

            displayChangeRate = formatDisplayChangeRate(totalValue, referenceTotalValue);
        }

        extraColumnValues.push(displayChangeRate);
    }

    return extraColumnValues;
}

function getAxisChartTooltipExtraColumnValues(seriesId: string, categoryIndex: number, currentValue: number): string[] {
    const extraColumnValues: string[] = [];

    if (!axisChartShowYearOverYear.value && !axisChartShowPeriodOverPeriod.value) {
        return extraColumnValues;
    }

    const seriesDimensionData = seriesDimensionTransactionExplorerDataMap.value[seriesId];

    if (!seriesDimensionData || !seriesDimensionData.categoryValues) {
        return extraColumnValues;
    }

    const values = seriesDimensionData.categoryValues;

    if (axisChartShowYearOverYear.value) {
        const yoyReferenceIndex = axisChartCategoryIndexYoYMap.value[categoryIndex];
        let displayChangeRate = '-';

        if (isNumber(yoyReferenceIndex) && yoyReferenceIndex >= 0 && yoyReferenceIndex < values.length) {
            displayChangeRate = formatDisplayChangeRate(currentValue, values[yoyReferenceIndex] ?? 0);
        }

        extraColumnValues.push(displayChangeRate);
    }

    if (axisChartShowPeriodOverPeriod.value) {
        const popReferenceIndex = categoryIndex - 1;
        let displayChangeRate = '-';

        if (popReferenceIndex >= 0 && popReferenceIndex < values.length) {
            displayChangeRate = formatDisplayChangeRate(currentValue, values[popReferenceIndex] ?? 0);
        }

        extraColumnValues.push(displayChangeRate);
    }

    return extraColumnValues;
}

function updateCategoryDimensionType(dimensionType: TransactionExplorerDataDimensionType): void {
    if (currentExplorer.value.categoryDimension !== dimensionType) {
        currentExplorer.value.categoryDimension = dimensionType;
    }

    if (currentExplorer.value.seriesDimension === currentExplorer.value.categoryDimension) {
        currentExplorer.value.seriesDimension = TransactionExplorerDataDimension.None.value;
    }
}

function onClickPieChartItem(item: Record<string, unknown>): void {
    if (!item || !('id' in item) || !('dimension' in item)) {
        return;
    }

    const data = (item as unknown) as CategoryDimensionData;
    const params: string = explorersStore.getTransactionListPageParams(data.dimension, data.id);

    if (params) {
        router.push(`/transaction/list?${params}`);
    }
}

function onClickTrendChartItem(itemId: string, categoryIndex: number): void {
    const categoryData = categoriedDataSortedByDisplayOrder.value[categoryIndex];

    if (categoryData) {
        const params: string = explorersStore.getTransactionListPageParams(categoryData.originalItem.categoryIdType, categoryData.originalItem.categoryId);

        if (params && categoryData.originalItem.categoryId !== 'none' && categoryData.originalItem.categoryId !== 'unknown') {
            router.push(`/transaction/list?${params}`);
        }
    }
}

function buildExportResults(): { headers: string[], data: string[][], supportedMermaidCharts?: ExportMermaidChartType[] } | undefined {
    if (currentExplorer.value.chartType === TransactionExplorerChartType.Pie.value || currentExplorer.value.chartType === TransactionExplorerChartType.Radar.value) {
        const valueMetric = TransactionExplorerValueMetric.valueOf(currentExplorer.value.valueMetric);
        let supportedMermaidCharts: ExportMermaidChartType[] | undefined = undefined;

        if (currentExplorer.value.chartType === TransactionExplorerChartType.Pie.value) {
            supportedMermaidCharts = [ ExportMermaidChartType.PieChart ];
        }

        return {
            headers: [
                tt('Name'),
                tt(valueMetric?.name ?? 'Unknown')
            ],
            data: categoryDimensionTransactionExplorerData.value.map(data => [
                data.name,
                valueMetric?.isAmount ? formatAmountToWesternArabicNumeralsWithoutDigitGrouping(data.totalAmount, defaultCurrency.value) : data.totalAmount.toString(10)
            ]),
            supportedMermaidCharts: supportedMermaidCharts
        };
    } else if (TransactionExplorerChartType.valueOf(currentExplorer.value.chartType)?.seriesDimensionRequired && axisChartDisplayType.value) {
        const results = axisChart.value?.exportData();

        if (!results) {
            return undefined;
        }

        let supportedMermaidCharts: ExportMermaidChartType[] | undefined = undefined;

        if (results.headers.length === 2 &&
            (
                currentExplorer.value.chartType === TransactionExplorerChartType.ColumnStacked.value ||
                currentExplorer.value.chartType === TransactionExplorerChartType.Column100PercentStacked.value ||
                currentExplorer.value.chartType === TransactionExplorerChartType.ColumnGrouped.value
            )
        ) {
            supportedMermaidCharts = [ ExportMermaidChartType.XYChartBar ];
        } else if (results.headers.length === 2 &&
            (
                currentExplorer.value.chartType === TransactionExplorerChartType.AreaStacked.value ||
                currentExplorer.value.chartType === TransactionExplorerChartType.Area100PercentStacked.value
            )
        ) {
            supportedMermaidCharts = [ ExportMermaidChartType.XYChartLine ];
        } else if (currentExplorer.value.chartType === TransactionExplorerChartType.LineGrouped.value) {
            supportedMermaidCharts = [ ExportMermaidChartType.XYChartLine ];
        }

        return {
            headers: results.headers,
            data: results.data,
            supportedMermaidCharts: supportedMermaidCharts
        };
    } else {
        return undefined;
    }
}

defineExpose({
    buildExportResults
});
</script>
