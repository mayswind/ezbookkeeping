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
                        :items="allTransactionExploreChartTypes"
                        :model-value="currentChartType"
                        @update:model-value="updateChartType($event as TransactionExploreChartTypeValue)"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading"
                        :label="tt('Axis / Category')"
                        :items="allTransactionExploreDataDimensions"
                        :model-value="currentCategoryDimension"
                        @update:model-value="updateCategoryDimension($event as TransactionExploreDataDimensionType)"
                    />
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || !TransactionExploreChartType.valueOf(currentChartType)?.seriesDimensionRequired"
                        :label="tt('Series')"
                        :items="allTransactionExploreDataDimensions"
                        :model-value="TransactionExploreChartType.valueOf(currentChartType)?.seriesDimensionRequired ? currentSeriesDimension : TransactionExploreDataDimension.None.value"
                        @update:model-value="updateSeriesDimension($event as TransactionExploreDataDimensionType)"
                    >
                        <template #item="{ props, item }">
                            <v-list-item :disabled="item.value === currentCategoryDimension && item.value !== TransactionExploreDataDimension.SeriesDimensionDefault.value" v-bind="props">
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
                        :items="allTransactionExploreValueMetrics"
                        :model-value="currentValueMetric"
                        @update:model-value="updateValueMetric($event as TransactionExploreValueMetricType)"
                    />
                    <v-spacer class="flex-1-1"/>
                </div>
            </v-col>
        </v-row>
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-if="currentChartType === TransactionExploreChartType.Pie.value">
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
            :items="categoryDimensionTransactionExploreData && categoryDimensionTransactionExploreData.length ? categoryDimensionTransactionExploreData : []"
            :min-valid-percent="0.0001"
            :show-value="true"
            :show-percent="true"
            :enable-click-item="true"
            :amount-value="exploresStore.transactionExploreFilter.valueMetric !== TransactionExploreValueMetric.TransactionCount.value"
            :default-currency="defaultCurrency"
            id-field="categoryId"
            name-field="categoryDisplayName"
            value-field="value"
            v-else-if="!loading"
            @click="onClickPieChartItem"
        />
    </v-card-text>
    <v-card-text :class="{ 'readonly': loading }" v-if="currentChartType === TransactionExploreChartType.Radar.value">
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
            :items="categoryDimensionTransactionExploreData && categoryDimensionTransactionExploreData.length ? categoryDimensionTransactionExploreData : []"
            :min-valid-percent="0.0001"
            :show-value="true"
            :show-percent="true"
            :amount-value="exploresStore.transactionExploreFilter.valueMetric !== TransactionExploreValueMetric.TransactionCount.value"
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
    TransactionExploreDimensionType,
    useExploresStore
} from '@/stores/explore.ts';

import { type NameValue } from '@/core/base.ts';
import {
    TransactionExploreChartTypeValue,
    TransactionExploreChartType,
    TransactionExploreDataDimensionType,
    TransactionExploreDataDimension,
    TransactionExploreValueMetricType,
    TransactionExploreValueMetric
} from '@/core/explore.ts';

import {
    isDefined
} from '@/lib/common.ts';

import {
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

interface InsightsExploreDataTableTabProps {
    loading?: boolean;
}

interface CategoryDimensionData {
    categoryDisplayName: string;
    categoryId: string;
    categoryIdType: TransactionExploreDimensionType;
    value: number;
}

defineProps<InsightsExploreDataTableTabProps>();

const router = useRouter();

const {
    tt,
    getAllTransactionExploreDataDimensions,
    getAllTransactionExploreValueMetrics,
    getAllTransactionExploreChartTypes,
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
const exploresStore = useExploresStore();

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const allTransactionExploreDataDimensions = computed<NameValue[]>(() => getAllTransactionExploreDataDimensions());
const allTransactionExploreValueMetrics = computed<NameValue[]>(() => getAllTransactionExploreValueMetrics());
const allTransactionExploreChartTypes = computed<NameValue[]>(() => getAllTransactionExploreChartTypes());

const currentCategoryDimension = computed<TransactionExploreDataDimensionType>(() => exploresStore.transactionExploreFilter.categoryDimension);
const currentSeriesDimension = computed<TransactionExploreDataDimensionType>(() => exploresStore.transactionExploreFilter.seriesDimension);
const currentValueMetric = computed<TransactionExploreValueMetricType>(() => exploresStore.transactionExploreFilter.valueMetric);
const currentChartType = computed<TransactionExploreChartTypeValue>(() => exploresStore.transactionExploreFilter.chartType);

const categoryDimensionTransactionExploreData = computed<CategoryDimensionData[]>(() => {
    if (currentChartType.value !== TransactionExploreChartType.Pie.value && currentChartType.value !== TransactionExploreChartType.Radar.value) {
        return [];
    }

    if (!exploresStore.categoriedTransactionExploreData || !exploresStore.categoriedTransactionExploreData.length) {
        return [];
    }

    const result: CategoryDimensionData[] = [];

    for (const categoriedData of exploresStore.categoriedTransactionExploreData) {
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
    let dimessionType: TransactionExploreDataDimensionType = TransactionExploreDataDimension.None.value;

    if ('categoryName' in info) {
        name = info.categoryName;
        needI18n = info.categoryNameNeedI18n;
        i18nParameters = info.categoryNameI18nParameters;
        dimessionType = exploresStore.transactionExploreFilter.categoryDimension;
    } else if ('seriesName' in info) {
        name = info.seriesName;
        needI18n = info.seriesNameNeedI18n;
        i18nParameters = info.seriesNameI18nParameters;
        dimessionType = exploresStore.transactionExploreFilter.seriesDimension;
    }

    let displayName: string = name;

    // convert the name to i18n if needed
    if (needI18n && i18nParameters) {
        displayName = tt(name, i18nParameters);
    } else if (needI18n && !i18nParameters) {
        displayName = tt(name);
    }

    // convert the name to formatted date time if needed
    if (dimessionType === TransactionExploreDataDimension.DateTime.value) {
        displayName = formatDateTimeToShortDateTime(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExploreDataDimension.DateTimeByYearMonthDay.value) {
        displayName = formatDateTimeToShortDate(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExploreDataDimension.DateTimeByYearMonth.value) {
        displayName = formatDateTimeToGregorianLikeShortYearMonth(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExploreDataDimension.DateTimeByYearQuarter.value) {
        displayName = formatDateTimeToGregorianLikeYearQuarter(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExploreDataDimension.DateTimeByYear.value) {
        displayName = formatDateTimeToGregorianLikeShortYear(parseDateTimeFromUnixTime(parseInt(name)));
    } else if (dimessionType === TransactionExploreDataDimension.DateTimeByFiscalYear.value) {
        displayName = formatDateTimeToGregorianLikeFiscalYear(parseDateTimeFromUnixTime(parseInt(name)));
    }

    if (dimessionType === TransactionExploreDataDimension.SourceAmount.value
        || dimessionType === TransactionExploreDataDimension.DestinationAmount.value) {
        if (name !== '' && name !== 'none' && Number.isFinite(parseInt(name))) {
            displayName = formatAmountToLocalizedNumerals(parseInt(name));
        }
    }

    return displayName;
}

function updateCategoryDimension(categoryDimension: TransactionExploreDataDimensionType): void {
    exploresStore.updateTransactionExploreFilter({
        categoryDimension: categoryDimension,
    });
}

function updateSeriesDimension(seriesDimension: TransactionExploreDataDimensionType): void {
    exploresStore.updateTransactionExploreFilter({
        seriesDimension: seriesDimension,
    });
}

function updateValueMetric(valueMetric: TransactionExploreValueMetricType): void {
    exploresStore.updateTransactionExploreFilter({
        valueMetric: valueMetric,
    });
}

function updateChartType(chartType: TransactionExploreChartTypeValue): void {
    exploresStore.updateTransactionExploreFilter({
        chartType: chartType,
    });
}

function onClickPieChartItem(item: Record<string, unknown>): void {
    if (!item || !('categoryId' in item) || !('categoryIdType' in item)) {
        return;
    }

    const data = (item as unknown) as CategoryDimensionData;
    const params: string = exploresStore.getTransactionListPageParams(data.categoryIdType, data.categoryId);

    if (params) {
        router.push(`/transaction/list?${params}`);
    }
}

function buildExportResults(): { headers: string[], data: string[][] } | undefined {
    if (currentChartType.value === TransactionExploreChartType.Pie.value || currentChartType.value === TransactionExploreChartType.Radar.value) {
        const valueMetric = TransactionExploreValueMetric.valueOf(exploresStore.transactionExploreFilter.valueMetric);

        return {
            headers: [
                tt('Name'),
                tt(valueMetric?.name ?? 'Unknown')
            ],
            data: categoryDimensionTransactionExploreData.value.map(data => [
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
