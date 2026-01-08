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
                        v-model="currentExplorer.categoryDimension"
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
            :amount-value="currentExplorer.valueMetric !== TransactionExplorerValueMetric.TransactionCount.value"
            :default-currency="defaultCurrency"
            id-field="id"
            name-field="name"
            value-field="totalAmount"
            v-else-if="!loading"
            @click="onClickPieChartItem"
        />
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-if="currentExplorer.chartType === TransactionExplorerChartType.Radar.value">
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
            :min-valid-percent="0.0001"
            :show-value="true"
            :show-percent="true"
            :amount-value="currentExplorer.valueMetric !== TransactionExplorerValueMetric.TransactionCount.value"
            :default-currency="defaultCurrency"
            name-field="name"
            value-field="totalAmount"
            v-else-if="!loading"
        />
    </v-card-text>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import {
    type CategoriedInfo,
    type SeriesedInfo,
    TransactionExplorerDimensionType,
    useExplorersStore
} from '@/stores/explorer.ts';

import { type NameValue, type TypeAndDisplayName } from '@/core/base.ts';
import { Month, WeekDay } from '@/core/datetime.ts';
import {
    TransactionExplorerChartType,
    TransactionExplorerDataDimensionType,
    TransactionExplorerDataDimension,
    TransactionExplorerValueMetric
} from '@/core/explorer.ts';

import { type SortableTransactionStatisticDataItem } from '@/models/transaction.ts';
import type { InsightsExplorer } from '@/models/explorer.ts';

import { isDefined } from '@/lib/common.ts';
import { parseDateTimeFromString } from '@/lib/datetime.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';

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
    formatDateTimeToShortDateTime,
    formatDateTimeToShortDate,
    formatDateTimeToGregorianLikeShortYear,
    formatDateTimeToGregorianLikeShortYearMonth,
    formatDateTimeToGregorianLikeYearQuarter,
    formatGregorianYearToGregorianLikeFiscalYear,
    formatAmountToLocalizedNumerals,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping
} = useI18n();

const userStore = useUserStore();
const explorersStore = useExplorersStore();

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const allTransactionExplorerDataDimensions = computed<NameValue[]>(() => getAllTransactionExplorerDataDimensions());
const allTransactionExplorerValueMetrics = computed<NameValue[]>(() => getAllTransactionExplorerValueMetrics());
const allTransactionExplorerChartTypes = computed<NameValue[]>(() => getAllTransactionExplorerChartTypes());
const allTransactionExplorerChartSortingTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsSortingTypes());

const currentExplorer = computed<InsightsExplorer>(() => explorersStore.currentInsightsExplorer);

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

function getCategoriedDataDisplayName(info: CategoriedInfo | SeriesedInfo): string {
    let name: string = '';
    let needI18n: boolean | undefined = false;
    let i18nParameters: Record<string, unknown> | undefined = undefined;
    let dimessionType: TransactionExplorerDimensionType = TransactionExplorerDimensionType.Other;
    let dimession: TransactionExplorerDataDimensionType = TransactionExplorerDataDimension.None.value;

    if ('categoryName' in info) {
        name = info.categoryName;
        needI18n = info.categoryNameNeedI18n;
        i18nParameters = info.categoryNameI18nParameters;
        dimessionType = info.categoryIdType;
        dimession = currentExplorer.value.categoryDimension;
    } else if ('seriesName' in info) {
        name = info.seriesName;
        needI18n = info.seriesNameNeedI18n;
        i18nParameters = info.seriesNameI18nParameters;
        dimessionType = info.seriesIdType;
        dimession = currentExplorer.value.seriesDimension;
    }

    let displayName: string = name;

    // convert the name to i18n if needed
    if (needI18n && i18nParameters) {
        displayName = tt(name, i18nParameters);
    } else if (needI18n && !i18nParameters) {
        displayName = tt(name);
    }

    // convert the name to formatted date time if needed
    if (dimession === TransactionExplorerDataDimension.DateTime.value) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToShortDateTime(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearMonthDay.value) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToShortDate(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearMonth.value) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToGregorianLikeShortYearMonth(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYearQuarter.value) {
        const parts = name.split('-');
        const year = parts.length === 2 ? parts[0] : '';
        const quarter = parts.length === 2 ? parseInt(parts[1] as string) : 0;
        const dateTime = year && quarter ? parseDateTimeFromString(`${year}-${quarter * 3}`, TransactionExplorerDimensionType.YearMonth) : undefined;
        displayName = dateTime ? formatDateTimeToGregorianLikeYearQuarter(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByYear.value) {
        const dateTime = parseDateTimeFromString(name, dimessionType);
        displayName = dateTime ? formatDateTimeToGregorianLikeShortYear(dateTime) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByFiscalYear.value) {
        displayName = formatGregorianYearToGregorianLikeFiscalYear(parseInt(name));
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByDayOfWeek.value) {
        const weekDay = WeekDay.parse(name);
        displayName = weekDay ? getWeekdayLongName(weekDay) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByDayOfMonth.value) {
        displayName = getMonthdayShortName(parseInt(name));
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByMonthOfYear.value) {
        const month = Month.valueOf(parseInt(name));
        displayName = month ? getMonthLongName(month.name) : tt('Unknown');
    } else if (dimession === TransactionExplorerDataDimension.DateTimeByQuarterOfYear.value) {
        displayName = getQuarterName(parseInt(name));
    }

    if (dimession === TransactionExplorerDataDimension.SourceAmount.value
        || dimession === TransactionExplorerDataDimension.DestinationAmount.value) {
        if (name !== '' && name !== 'none' && Number.isFinite(parseInt(name))) {
            displayName = formatAmountToLocalizedNumerals(parseInt(name));
        }
    }

    return displayName;
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

function buildExportResults(): { headers: string[], data: string[][] } | undefined {
    if (currentExplorer.value.chartType === TransactionExplorerChartType.Pie.value || currentExplorer.value.chartType === TransactionExplorerChartType.Radar.value) {
        const valueMetric = TransactionExplorerValueMetric.valueOf(currentExplorer.value.valueMetric);

        return {
            headers: [
                tt('Name'),
                tt(valueMetric?.name ?? 'Unknown')
            ],
            data: categoryDimensionTransactionExplorerData.value.map(data => [
                data.name,
                valueMetric?.isAmount ? formatAmountToWesternArabicNumeralsWithoutDigitGrouping(data.totalAmount) : data.totalAmount.toString(10)
            ])
        };
    } else {
        return undefined;
    }
}

defineExpose({
    buildExportResults
});
</script>
