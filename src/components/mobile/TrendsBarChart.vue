<template>
    <f7-list class="skeleton-text" v-if="loading">
        <f7-list-item class="statistics-list-item" link="#" :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ]">
            <template #media>
                <div class="display-flex no-padding-horizontal">
                    <div class="display-flex align-items-center statistics-icon">
                        <f7-icon f7="app_fill"></f7-icon>
                    </div>
                </div>
            </template>
            <template #title>
                <div class="statistics-list-item-text">
                    <span>Date Range</span>
                </div>
            </template>
            <template #after>
                <span>0.00 USD</span>
            </template>
            <template #inner-end>
                <div class="statistics-item-end">
                    <div class="statistics-percent-line">
                        <f7-progressbar></f7-progressbar>
                    </div>
                </div>
            </template>
        </f7-list-item>
    </f7-list>

    <f7-list v-else-if="!loading && (!allDisplayDataItems || !allDisplayDataItems.data || !allDisplayDataItems.data.length)">
        <f7-list-item :title="tt('No transaction data')"></f7-list-item>
    </f7-list>

    <f7-list v-if="!loading && allDisplayDataItems && allDisplayDataItems.data && allDisplayDataItems.data.length">
        <f7-list-item v-if="allDisplayDataItems.legends && allDisplayDataItems.legends.length > 1">
            <div class="display-flex" style="flex-wrap: wrap">
                <div class="trends-bar-chart-legend display-flex align-items-center"
                     :class="{ 'trends-bar-chart-legend-unselected': !!unselectedLegends[legend.id] }"
                     :key="idx"
                     v-for="(legend, idx) in allDisplayDataItems.legends"
                     @click="toggleLegend(legend)">
                    <f7-icon f7="app_fill" class="trends-bar-chart-legend-icon" :style="{ 'color': unselectedLegends[legend.id] ? '' : legend.color }"></f7-icon>
                    <span class="trends-bar-chart-legend-text">{{ legend.name }}</span>
                </div>
            </div>
        </f7-list-item>
    </f7-list>

    <f7-list :key="`trends-bar-chart-${allDisplayDataItemsVersion}`"
             :virtual-list="useVirtualList"
             :virtual-list-params="useVirtualList ? { items: allDisplayDataItems.data, renderExternal, height: 'auto' } : undefined"
             v-if="!loading && allDisplayDataItems && allDisplayDataItems.data && allDisplayDataItems.data.length">
        <f7-list-item link="#"
                      :key="item.index"
                      :class="{ 'statistics-list-item': true, 'statistics-list-item-stacked': stacked, 'statistics-list-item-non-stacked': !stacked }"
                      :style="useVirtualList ? `top: ${virtualDataItems.topPosition}px` : undefined"
                      :virtual-list-index="item.index"
                      v-for="item in (useVirtualList ? virtualDataItems.items : allDisplayDataItems.data)"
                      @click="clickItem(item)"
        >
            <template #media>
                <div class="display-flex no-padding-horizontal">
                    <div class="display-flex align-items-center statistics-icon">
                        <f7-icon f7="calendar"></f7-icon>
                    </div>
                </div>
            </template>

            <template #title>
                <div class="statistics-list-item-text">
                    <span>{{ item.displayDateRange }}</span>
                </div>
                <div class="full-line statistics-percent-line statistics-multi-percent-line display-flex flex-direction-column" v-if="!stacked && item.items.length > 1">
                    <div class="display-flex flex-direction-column"
                         style="margin-top: 4px"
                         :key="dataIdx"
                         v-for="(data, dataIdx) in item.items"
                         v-show="data.value.isPositive()">
                        <div class="full-line display-flex flex-direction-row">
                            <div class="display-inline-flex" :style="{ 'width': data.value.divide(item.maxAmount).multiply(item.percent).toDoubleNumber() + '%' }">
                                <f7-progressbar :progress="100" :style="{ '--f7-progressbar-progress-color': (data.color ? data.color : '') } "></f7-progressbar>
                            </div>
                            <div class="display-inline-flex" :style="{ 'width': (100.0 - data.value.divide(item.maxAmount).multiply(item.percent).toDoubleNumber()) + '%' }"
                                 v-if="data.value.divide(item.maxAmount).multiply(item.percent).lessThan(100)">
                                <f7-progressbar :progress="0"></f7-progressbar>
                            </div>
                        </div>
                    </div>
                </div>
            </template>

            <template #after>
                <span v-if="stacked || item.items.length <= 1">{{ formatAmountToLocalizedNumeralsWithCurrency(item.totalAmount, defaultCurrency) }}</span>
            </template>

            <template #inner-end>
                <div class="statistics-item-end" v-if="stacked || item.items.length <= 1">
                    <div class="statistics-percent-line statistics-multi-percent-line display-flex">
                        <div class="display-inline-flex" :style="{ 'width': (data.value.divide(item.totalPositiveAmount).multiply(item.percent).toDoubleNumber()) + '%' }"
                             :key="dataIdx"
                             v-for="(data, dataIdx) in item.items"
                             v-show="data.value.isPositive()">
                            <f7-progressbar :progress="100" :style="{ '--f7-progressbar-progress-color': (data.color ? data.color : '') } "></f7-progressbar>
                        </div>
                        <div class="display-inline-flex" :style="{ 'width': (100.0 - item.percent) + '%' }"
                             v-if="item.percent < 100.0">
                            <f7-progressbar :progress="0"></f7-progressbar>
                        </div>
                    </div>
                </div>
            </template>
        </f7-list-item>
    </f7-list>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type TrendsChartDateType,
    type CommonTrendsChartProps,
    type TrendsBarChartClickEvent,
    useTrendsChartBase
} from '@/components/base/TrendsChartBase.ts'

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

import { itemAndIndex } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import {
    type UnixTimeRange,
    DateRangeScene
} from '@/core/datetime.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import {
    ChartDataAggregationType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import type { SortableTransactionStatisticDataItem } from '@/models/transaction.ts';

import {
    BIG_DECIMAL_ZERO,
    isBigDecimal
} from '@/lib/numeral.ts';
import {
    parseDateTimeFromUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange,
    getFiscalYearFromUnixTime
} from '@/lib/datetime.ts';
import {
    getDisplayColor
} from '@/lib/color.ts';
import {
    sortStatisticsItems
} from '@/lib/statistics.ts';

interface TrendsBarChartLegend {
    readonly id: string;
    readonly name: string;
    readonly color: ColorStyleValue;
    readonly displayOrders: number[];
}

interface TrendsBarChartDataAmount extends SortableTransactionStatisticDataItem, TrendsBarChartLegend {
    value: BigDecimal;
}

interface TrendsBarChartDataItem {
    index: number;
    dateRange: UnixTimeRange;
    displayDateRange: string;
    items: TrendsBarChartDataAmount[];
    totalAmount: BigDecimal;
    totalPositiveAmount: BigDecimal;
    maxAmount: BigDecimal;
    percent: number;
}

interface TrendsBarChartVirtualListData {
    items: TrendsBarChartDataItem[],
    topPosition: number
}

interface TrendsBarChartData {
    readonly data: TrendsBarChartDataItem[];
    readonly legends: TrendsBarChartLegend[];
}

interface MobileTrendsChartProps<T extends TrendsChartDateType> extends CommonTrendsChartProps<T> {
    loading?: boolean;
}

const props = defineProps<MobileTrendsChartProps<TrendsChartDateType>>();

const emit = defineEmits<{
    (e: 'click', value: TrendsBarChartClickEvent): void;
}>();

const {
    tt,
    formatDateTimeToShortDate,
    formatDateTimeToGregorianLikeShortYear,
    formatDateTimeToGregorianLikeShortYearMonth,
    formatYearQuarterToGregorianLikeYearQuarter,
    formatDateTimeToGregorianLikeFiscalYear,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const { allDateRanges, getItemName } = useTrendsChartBase(props);

const settingsStore = useSettingsStore();
const userStore = useUserStore();

const allDisplayDataItemsVersion = ref<number>(0);
const unselectedLegends = ref<Record<string, boolean>>({});

const virtualDataItems = ref<TrendsBarChartVirtualListData>({
    items: [],
    topPosition: 0
});

const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);
const useVirtualList = computed<boolean>(() => allDisplayDataItems.value.legends.length <= 1 || props.stacked);

const allDisplayDataItems = computed<TrendsBarChartData>(() => {
    const allDateRangeItemsMap: Record<string, TrendsBarChartDataAmount[]> = {};
    const legends: TrendsBarChartLegend[] = [];

    for (const [item, index] of itemAndIndex(props.items)) {
        if (item.hidden) {
            continue;
        }

        const id = item.id ?? getItemName(item.name);

        const legend: TrendsBarChartLegend = {
            id: id,
            name: item.name ? getItemName(item.name) : id,
            color: getDisplayColor(props.useCustomColor && item.color ? item.color : chartColors.value[index % chartColors.value.length]),
            displayOrders: item.displayOrders ? item.displayOrders : [0]
        };

        legends.push(legend);

        if (unselectedLegends.value[id]) {
            continue;
        }

        const dateRangeItemMap: Record<string, TrendsBarChartDataAmount> = {};

        for (const dataItem of item.items) {
            let dateRangeKey = '';

            if (props.chartMode === 'daily' && 'month' in dataItem) {
                if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                    dateRangeKey = dataItem.year.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                    const fiscalYear = getFiscalYearFromUnixTime(
                        getYearMonthFirstUnixTime({ year: dataItem.year, month1base: dataItem.month }),
                        props.fiscalYearStart
                    );
                    dateRangeKey = fiscalYear.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                    dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month - 1) / 3) + 1}`;
                } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                    dateRangeKey = `${dataItem.year}-${dataItem.month}`;
                } else { // if (props.dateAggregationType === ChartDateAggregationType.Day.type) {
                    dateRangeKey = `${dataItem.year}-${dataItem.month}-${dataItem.day}`;
                }
            } else if (props.chartMode === 'monthly' && 'month1base' in dataItem) {
                if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                    dateRangeKey = dataItem.year.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                    const fiscalYear = getFiscalYearFromUnixTime(
                        getYearMonthFirstUnixTime({ year: dataItem.year, month1base: dataItem.month1base }),
                        props.fiscalYearStart
                    );
                    dateRangeKey = fiscalYear.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                    dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month1base - 1) / 3) + 1}`;
                } else { // if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                    dateRangeKey = `${dataItem.year}-${dataItem.month1base}`;
                }
            }

            const value = dataItem.value;

            if (dateRangeItemMap[dateRangeKey]) {
                if (isBigDecimal(value)) {
                    if (props.dataAggregationType === ChartDataAggregationType.Sum) {
                        dateRangeItemMap[dateRangeKey]!.value = dateRangeItemMap[dateRangeKey]!.value.add(value);
                    } else if (props.dataAggregationType === ChartDataAggregationType.Last) {
                        dateRangeItemMap[dateRangeKey]!.value = value;
                    }
                }
            } else {
                const allDataItems: TrendsBarChartDataAmount[] = allDateRangeItemsMap[dateRangeKey] || [];
                const finalDataItem: TrendsBarChartDataAmount = Object.assign({}, legend, {
                    value: isBigDecimal(value) ? value : BIG_DECIMAL_ZERO
                });

                allDataItems.push(finalDataItem);
                dateRangeItemMap[dateRangeKey] = finalDataItem;
                allDateRangeItemsMap[dateRangeKey] = allDataItems;
            }
        }
    }

    const finalDataItems: TrendsBarChartDataItem[] = [];
    let maxTotalAmount: BigDecimal = BIG_DECIMAL_ZERO;

    for (const dateRange of allDateRanges.value) {
        let dateRangeKey = '';

        if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
            dateRangeKey = dateRange.year.toString();
        } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
            dateRangeKey = dateRange.year.toString();
        } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
            dateRangeKey = `${dateRange.year}-${dateRange.quarter}`;
        } else if (props.dateAggregationType === ChartDateAggregationType.Month.type && 'month0base' in dateRange) {
            dateRangeKey = `${dateRange.year}-${dateRange.month0base + 1}`;
        } else if (props.dateAggregationType === ChartDateAggregationType.Day.type && 'day' in dateRange && props.chartMode === 'daily') {
            dateRangeKey = `${dateRange.year}-${dateRange.month}-${dateRange.day}`;
        }

        const minDateTime = parseDateTimeFromUnixTime(dateRange.minUnixTime);
        let displayDateRange = '';

        if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
            displayDateRange = formatDateTimeToGregorianLikeShortYear(minDateTime);
        } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
            displayDateRange = formatDateTimeToGregorianLikeFiscalYear(minDateTime);
        } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
            displayDateRange = formatYearQuarterToGregorianLikeYearQuarter(dateRange.year, dateRange.quarter);
        } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
            displayDateRange = formatDateTimeToGregorianLikeShortYearMonth(minDateTime);
        } else if (props.dateAggregationType === ChartDateAggregationType.Day.type && props.chartMode === 'daily') {
            displayDateRange = formatDateTimeToShortDate(minDateTime);
        }

        const dataItems = allDateRangeItemsMap[dateRangeKey] || [];
        let totalAmount: BigDecimal = BIG_DECIMAL_ZERO;
        let totalPositiveAmount = BIG_DECIMAL_ZERO;
        let maxAmount = BIG_DECIMAL_ZERO;

        sortStatisticsItems(dataItems, props.sortingType);

        for (const dataItem of dataItems) {
            if (dataItem.value.isPositive()) {
                totalPositiveAmount = totalPositiveAmount.add(dataItem.value);
            }

            totalAmount = totalAmount.add(dataItem.value);

            if (dataItem.value.greaterThan(maxAmount)) {
                maxAmount = dataItem.value;
            }
        }

        if (totalAmount.greaterThan(maxTotalAmount)) {
            maxTotalAmount = totalAmount;
        }

        const finalDataItem: TrendsBarChartDataItem = {
            index: finalDataItems.length,
            dateRange: dateRange,
            displayDateRange: displayDateRange,
            items: dataItems,
            totalAmount: totalAmount,
            totalPositiveAmount: totalPositiveAmount,
            maxAmount: maxAmount,
            percent: 0.0
        };

        finalDataItems.push(finalDataItem);
    }

    for (const finalDataItem of finalDataItems) {
        if (maxTotalAmount.isPositive() && finalDataItem.totalAmount.isPositive()) {
            finalDataItem.percent = finalDataItem.totalAmount.divide(maxTotalAmount).multiply(100).toDoubleNumber();
        } else {
            finalDataItem.percent = 0.0;
        }
    }

    return {
        data: finalDataItems,
        legends: legends
    };
});

function clickItem(item: TrendsBarChartDataItem): void {
    let itemId = '';

    for (const item of props.items) {
        if (item.hidden) {
            continue;
        }

        const id = item.id ?? getItemName(item.name);

        if (unselectedLegends.value[id]) {
            continue;
        }

        if (itemId.length) {
            itemId += ',';
        }

        itemId += id;
    }

    const dateRange = item.dateRange;
    let minUnixTime = dateRange.minUnixTime;
    let maxUnixTime = dateRange.maxUnixTime;

    if (props.chartMode === 'daily') {
        if (props.startTime) {
            if (props.startTime > minUnixTime) {
                minUnixTime = props.startTime;
            }
        }

        if (props.endTime) {
            if (props.endTime < maxUnixTime) {
                maxUnixTime = props.endTime;
            }
        }
    } else if (props.chartMode === 'monthly') {
        if (props.startYearMonth) {
            const startMinUnixTime = getYearMonthFirstUnixTime(props.startYearMonth);

            if (startMinUnixTime > minUnixTime) {
                minUnixTime = startMinUnixTime;
            }
        }

        if (props.endYearMonth) {
            const endMaxUnixTime = getYearMonthLastUnixTime(props.endYearMonth);

            if (endMaxUnixTime < maxUnixTime) {
                maxUnixTime = endMaxUnixTime;
            }
        }
    }

    const dateRangeType = getDateTypeByDateRange(minUnixTime, maxUnixTime, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart, DateRangeScene.Normal);

    emit('click', {
        itemId: itemId,
        dateRange: {
            minTime: minUnixTime,
            maxTime: maxUnixTime,
            dateType: dateRangeType
        }
    });
}

function toggleLegend(legend: TrendsBarChartLegend): void {
    if (unselectedLegends.value[legend.id]) {
        delete unselectedLegends.value[legend.id];
    } else {
        unselectedLegends.value[legend.id] = true;
    }
}

function renderExternal(vl: unknown, vlData: TrendsBarChartVirtualListData): void {
    virtualDataItems.value = vlData;
}

watch(allDisplayDataItems, () => {
    allDisplayDataItemsVersion.value++;
}, {
    deep: true
});
</script>

<style>
.trends-bar-chart-legend {
    margin-inline-end: 4px;
    cursor: pointer;
}

.trends-bar-chart-legend-icon.f7-icons {
    font-size: var(--ebk-trends-bar-chart-legend-icon-font-size);
    margin-inline-end: 2px;
}

.trends-bar-chart-legend-unselected .trends-bar-chart-legend-icon.f7-icons {
    color: #cccccc;
}

.trends-bar-chart-legend-text {
    font-size: var(--ebk-trends-bar-chart-legend-text-font-size);
}

.trends-bar-chart-legend-unselected .trends-bar-chart-legend-text {
    color: #cccccc;
}
</style>
