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
                        :disabled="loading"
                        :label="tt('Chart Type')"
                        :items="allTransactionExplorerChartTypes"
                        :model-value="currentChartType"
                        @update:model-value="updateChartType($event as TransactionExplorerChartTypeValue)"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading"
                        :label="tt('Axis / Category')"
                        :items="allTransactionExplorerDataDimensions"
                        :model-value="currentCategoryDimension"
                        @update:model-value="updateCategoryDimension($event as TransactionExplorerDataDimensionType)"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || !TransactionExplorerChartType.valueOf(currentChartType)?.seriesDimensionRequired"
                        :label="tt('Series')"
                        :items="allTransactionExplorerDataDimensions"
                        :model-value="TransactionExplorerChartType.valueOf(currentChartType)?.seriesDimensionRequired ? currentSeriesDimension : TransactionExplorerDataDimension.None.value"
                        @update:model-value="updateSeriesDimension($event as TransactionExplorerDataDimensionType)"
                    >
                        <template #item="{ props, item }">
                            <v-list-item :disabled="item.value === currentCategoryDimension && item.value !== TransactionExplorerDataDimension.SeriesDimensionDefault.value" v-bind="props">
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
                        :disabled="loading"
                        :label="tt('Value Metric')"
                        :items="allTransactionExplorerValueMetrics"
                        :model-value="currentValueMetric"
                        @update:model-value="updateValueMetric($event as TransactionExplorerValueMetricType)"
                    />
                    <v-spacer class="flex-1-1"/>
                </div>
            </v-col>
        </v-row>
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-if="currentChartType === TransactionExplorerChartType.Pie.value">
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
            :min-valid-percent="0.0001"
            :show-value="true"
            :show-percent="true"
            :enable-click-item="true"
            :amount-value="explorersStore.transactionExplorerFilter.valueMetric !== TransactionExplorerValueMetric.TransactionCount.value"
            :default-currency="defaultCurrency"
            id-field="categoryId"
            name-field="categoryDisplayName"
            value-field="value"
            v-else-if="!loading"
            @click="onClickPieChartItem"
        />
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-if="currentChartType === TransactionExplorerChartType.Radar.value">
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
            :amount-value="explorersStore.transactionExplorerFilter.valueMetric !== TransactionExplorerValueMetric.TransactionCount.value"
            :default-currency="defaultCurrency"
            name-field="categoryDisplayName"
            value-field="value"
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

import { type NameValue } from '@/core/base.ts';
import { Month, WeekDay } from '@/core/datetime.ts';
import {
    TransactionExplorerChartTypeValue,
    TransactionExplorerChartType,
    TransactionExplorerDataDimensionType,
    TransactionExplorerDataDimension,
    TransactionExplorerValueMetricType,
    TransactionExplorerValueMetric
} from '@/core/explorer.ts';

import {
    isDefined
} from '@/lib/common.ts';

import {
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

interface InsightsExplorerDataTableTabProps {
    loading?: boolean;
}

interface CategoryDimensionData {
    categoryDisplayName: string;
    categoryId: string;
    categoryIdType: TransactionExplorerDimensionType;
    value: number;
}

defineProps<InsightsExplorerDataTableTabProps>();

const router = useRouter();

const {
    tt,
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
    formatDateTimeToGregorianLikeFiscalYear,
    formatAmountToLocalizedNumerals,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping
} = useI18n();

const userStore = useUserStore();
const explorersStore = useExplorersStore();

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const allTransactionExplorerDataDimensions = computed<NameValue[]>(() => getAllTransactionExplorerDataDimensions());
const allTransactionExplorerValueMetrics = computed<NameValue[]>(() => getAllTransactionExplorerValueMetrics());
const allTransactionExplorerChartTypes = computed<NameValue[]>(() => getAllTransactionExplorerChartTypes());

const currentCategoryDimension = computed<TransactionExplorerDataDimensionType>(() => explorersStore.transactionExplorerFilter.categoryDimension);
const currentSeriesDimension = computed<TransactionExplorerDataDimensionType>(() => explorersStore.transactionExplorerFilter.seriesDimension);
const currentValueMetric = computed<TransactionExplorerValueMetricType>(() => explorersStore.transactionExplorerFilter.valueMetric);
const currentChartType = computed<TransactionExplorerChartTypeValue>(() => explorersStore.transactionExplorerFilter.chartType);

const categoryDimensionTransactionExplorerData = computed<CategoryDimensionData[]>(() => {
    if (currentChartType.value !== TransactionExplorerChartType.Pie.value && currentChartType.value !== TransactionExplorerChartType.Radar.value) {
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
            categoryDisplayName: displayName,
            categoryId: categoriedData.categoryId,
            categoryIdType: categoriedData.categoryIdType,
            value: data.value
        });
    }

    return result;
});

function getCategoriedDataDisplayName(info: CategoriedInfo | SeriesedInfo): string {
    let name: string = '';
    let needI18n: boolean | undefined = false;
    let i18nParameters: Record<string, unknown> | undefined = undefined;
    let dimessionType: TransactionExplorerDataDimensionType = TransactionExplorerDataDimension.None.value;

    if ('categoryName' in info) {
        name = info.categoryName;
        needI18n = info.categoryNameNeedI18n;
        i18nParameters = info.categoryNameI18nParameters;
        dimessionType = explorersStore.transactionExplorerFilter.categoryDimension;
    } else if ('seriesName' in info) {
        name = info.seriesName;
        needI18n = info.seriesNameNeedI18n;
        i18nParameters = info.seriesNameI18nParameters;
        dimessionType = explorersStore.transactionExplorerFilter.seriesDimension;
    }

    let displayName: string = name;

    // convert the name to i18n if needed
    if (needI18n && i18nParameters) {
        displayName = tt(name, i18nParameters);
    } else if (needI18n && !i18nParameters) {
        displayName = tt(name);
    }

    // convert the name to formatted date time if needed
    if (dimessionType === TransactionExplorerDataDimension.DateTime.value) {
        displayName = formatDateTimeToShortDateTime(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByYearMonthDay.value) {
        displayName = formatDateTimeToShortDate(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByYearMonth.value) {
        displayName = formatDateTimeToGregorianLikeShortYearMonth(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByYearQuarter.value) {
        displayName = formatDateTimeToGregorianLikeYearQuarter(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByYear.value) {
        displayName = formatDateTimeToGregorianLikeShortYear(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByFiscalYear.value) {
        displayName = formatDateTimeToGregorianLikeFiscalYear(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByDayOfWeek.value) {
        const weekDay = WeekDay.parse(name);
        displayName = weekDay ? getWeekdayLongName(weekDay) : tt('Unknown');
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByDayOfMonth.value) {
        displayName = getMonthdayShortName(parseInt(name));
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByMonthOfYear.value) {
        const month = Month.valueOf(parseInt(name));
        displayName = month ? getMonthLongName(month.name) : tt('Unknown');
    } else if (dimessionType === TransactionExplorerDataDimension.DateTimeByQuarterOfYear.value) {
        displayName = getQuarterName(parseInt(name));
    }

    if (dimessionType === TransactionExplorerDataDimension.SourceAmount.value
        || dimessionType === TransactionExplorerDataDimension.DestinationAmount.value) {
        if (name !== '' && name !== 'none' && Number.isFinite(parseInt(name))) {
            displayName = formatAmountToLocalizedNumerals(parseInt(name));
        }
    }

    return displayName;
}

function updateCategoryDimension(categoryDimension: TransactionExplorerDataDimensionType): void {
    explorersStore.updateTransactionExplorerFilter({
        categoryDimension: categoryDimension,
    });
}

function updateSeriesDimension(seriesDimension: TransactionExplorerDataDimensionType): void {
    explorersStore.updateTransactionExplorerFilter({
        seriesDimension: seriesDimension,
    });
}

function updateValueMetric(valueMetric: TransactionExplorerValueMetricType): void {
    explorersStore.updateTransactionExplorerFilter({
        valueMetric: valueMetric,
    });
}

function updateChartType(chartType: TransactionExplorerChartTypeValue): void {
    explorersStore.updateTransactionExplorerFilter({
        chartType: chartType,
    });
}

function onClickPieChartItem(item: Record<string, unknown>): void {
    if (!item || !('categoryId' in item) || !('categoryIdType' in item)) {
        return;
    }

    const data = (item as unknown) as CategoryDimensionData;
    const params: string = explorersStore.getTransactionListPageParams(data.categoryIdType, data.categoryId);

    if (params) {
        router.push(`/transaction/list?${params}`);
    }
}

function buildExportResults(): { headers: string[], data: string[][] } | undefined {
    if (currentChartType.value === TransactionExplorerChartType.Pie.value || currentChartType.value === TransactionExplorerChartType.Radar.value) {
        const valueMetric = TransactionExplorerValueMetric.valueOf(explorersStore.transactionExplorerFilter.valueMetric);

        return {
            headers: [
                tt('Name'),
                tt(valueMetric?.name ?? 'Unknown')
            ],
            data: categoryDimensionTransactionExplorerData.value.map(data => [
                data.categoryDisplayName,
                valueMetric?.isAmount ? formatAmountToWesternArabicNumeralsWithoutDigitGrouping(data.value) : data.value.toString(10)
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
