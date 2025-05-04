<template>
    <f7-list class="statistics-list-item skeleton-text" v-if="loading">
        <f7-list-item link="#" :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ]">
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

    <f7-list v-else-if="!loading && allDisplayDataItems && allDisplayDataItems.data && allDisplayDataItems.data.length">
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
        <f7-list-item class="statistics-list-item"
                      link="#"
                      :key="idx"
                      v-for="(item, idx) in allDisplayDataItems.data"
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
            </template>

            <template #after>
                <span>{{ formatAmountWithCurrency(item.totalAmount, defaultCurrency) }}</span>
            </template>

            <template #inner-end>
                <div class="statistics-item-end">
                    <div class="statistics-percent-line statistics-multi-percent-line display-flex">
                        <div class="display-inline-flex" :style="{ 'width': (item.percent * data.totalAmount / item.totalPositiveAmount) + '%' }"
                             :key="dataIdx"
                             v-for="(data, dataIdx) in item.items"
                             v-show="data.totalAmount > 0">
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
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonTrendsChartProps, type TrendsBarChartClickEvent, useTrendsChartBase } from '@/components/base/TrendsChartBase.ts'

import { useUserStore } from '@/stores/user.ts';

import { type YearMonth, type UnixTimeRange, DateRangeScene } from '@/core/datetime.ts';
import { ChartDateAggregationType } from '@/core/statistics.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';
import type { YearMonthDataItem, SortableTransactionStatisticDataItem } from '@/models/transaction.ts';

import {
    isNumber
} from '@/lib/common.ts';
import {
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange,
    getFiscalYearFromUnixTime
} from '@/lib/datetime.ts';
import {
    sortStatisticsItems
} from '@/lib/statistics.ts';

interface TrendsBarChartLegend {
    readonly id: string;
    readonly name: string;
    readonly color: string;
    readonly displayOrders: number[];
}

interface TrendsBarChartDataAmount extends SortableTransactionStatisticDataItem, TrendsBarChartLegend {
    totalAmount: number;
}

interface TrendsBarChartDataItem {
    dateRange: UnixTimeRange;
    displayDateRange: string;
    items: TrendsBarChartDataAmount[];
    totalAmount: number;
    totalPositiveAmount: number;
    percent: number;
}

interface TrendsBarChartData {
    readonly data: TrendsBarChartDataItem[];
    readonly legends: TrendsBarChartLegend[];
}

interface MobileTrendsChartProps<T extends YearMonth> extends CommonTrendsChartProps<T> {
    loading?: boolean;
}

const props = defineProps<MobileTrendsChartProps<YearMonthDataItem>>();

const emit = defineEmits<{
    (e: 'click', value: TrendsBarChartClickEvent): void;
}>();

const { tt, formatUnixTimeToShortYear, formatYearQuarter, formatUnixTimeToShortYearMonth, formatUnixTimeToFiscalYear, formatYearToFiscalYear, formatAmountWithCurrency } = useI18n();
const { allDateRanges, getItemName, getColor } = useTrendsChartBase(props);

const userStore = useUserStore();

const unselectedLegends = ref<Record<string, boolean>>({});

const allDisplayDataItems = computed<TrendsBarChartData>(() => {
    const allDateRangeItemsMap: Record<string, TrendsBarChartDataAmount[]> = {};
    const legends: TrendsBarChartLegend[] = [];

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];

        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        const id = (props.idField && item[props.idField]) ? item[props.idField] as string : getItemName(item[props.nameField] as string);

        const legend: TrendsBarChartLegend = {
            id: id,
            name: (props.nameField && item[props.nameField]) ? getItemName(item[props.nameField] as string) : id,
            color: getColor(props.colorField && item[props.colorField] ? item[props.colorField] as string : DEFAULT_CHART_COLORS[i % DEFAULT_CHART_COLORS.length]),
            displayOrders: (props.displayOrdersField && item[props.displayOrdersField]) ? item[props.displayOrdersField] as number[] : [0]
        };

        legends.push(legend);

        if (unselectedLegends.value[id]) {
            continue;
        }

        const dateRangeItemMap: Record<string, TrendsBarChartDataAmount> = {};

        for (let j = 0; j < item.items.length; j++) {
            const dataItem = item.items[j];
            let dateRangeKey = '';

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                dateRangeKey = dataItem.year.toString();
            } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                const fiscalYear = getFiscalYearFromUnixTime(
                    getYearMonthFirstUnixTime({ year: dataItem.year, month: dataItem.month }),
                    props.fiscalYearStart
                );
                dateRangeKey = formatYearToFiscalYear(fiscalYear);
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month - 1) / 3) + 1}`;
            } else { // if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                dateRangeKey = `${dataItem.year}-${dataItem.month}`;
            }

            if (dateRangeItemMap[dateRangeKey]) {
                dateRangeItemMap[dateRangeKey].totalAmount += (props.valueField && isNumber(dataItem[props.valueField])) ? dataItem[props.valueField] as number : 0;
            } else {
                const allDataItems: TrendsBarChartDataAmount[] = allDateRangeItemsMap[dateRangeKey] || [];
                const finalDataItem: TrendsBarChartDataAmount = Object.assign({}, legend, {
                    totalAmount: (props.valueField && isNumber(dataItem[props.valueField])) ? dataItem[props.valueField] as number : 0
                });

                allDataItems.push(finalDataItem);
                dateRangeItemMap[dateRangeKey] = finalDataItem;
                allDateRangeItemsMap[dateRangeKey] = allDataItems;
            }
        }
    }

    const finalDataItems: TrendsBarChartDataItem[] = [];
    let maxTotalAmount = 0;

    for (let i = 0; i < allDateRanges.value.length; i++) {
        const dateRange = allDateRanges.value[i];
        let dateRangeKey = '';

        if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
            dateRangeKey = dateRange.year.toString();
        } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
            dateRangeKey = formatYearToFiscalYear(dateRange.year);
        } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
            dateRangeKey = `${dateRange.year}-${dateRange.quarter}`;
        } else if (props.dateAggregationType === ChartDateAggregationType.Month.type && 'month' in dateRange) {
            dateRangeKey = `${dateRange.year}-${dateRange.month + 1}`;
        }

        let displayDateRange = '';

        if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
            displayDateRange = formatUnixTimeToShortYear(dateRange.minUnixTime);
        } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
            displayDateRange = formatUnixTimeToFiscalYear(dateRange.minUnixTime);
        } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
            displayDateRange = formatYearQuarter(dateRange.year, dateRange.quarter);
        } else { // if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
            displayDateRange = formatUnixTimeToShortYearMonth(dateRange.minUnixTime);
        }

        const dataItems = allDateRangeItemsMap[dateRangeKey] || [];
        let totalAmount = 0;
        let totalPositiveAmount = 0;

        sortStatisticsItems(dataItems, props.sortingType);

        for (let j = 0; j < dataItems.length; j++) {
            if (dataItems[j].totalAmount > 0) {
                totalPositiveAmount += dataItems[j].totalAmount;
            }

            totalAmount += dataItems[j].totalAmount;
        }

        if (totalAmount > maxTotalAmount) {
            maxTotalAmount = totalAmount;
        }

        const finalDataItem: TrendsBarChartDataItem = {
            dateRange: dateRange,
            displayDateRange: displayDateRange,
            items: dataItems,
            totalAmount: totalAmount,
            totalPositiveAmount: totalPositiveAmount,
            percent: 0.0
        };

        finalDataItems.push(finalDataItem);
    }

    for (let i = 0; i < finalDataItems.length; i++) {
        if (maxTotalAmount > 0 && finalDataItems[i].totalAmount > 0) {
            finalDataItems[i].percent = 100.0 * finalDataItems[i].totalAmount / maxTotalAmount;
        } else {
            finalDataItems[i].percent = 0.0;
        }
    }

    return {
        data: finalDataItems,
        legends: legends
    };
});

function clickItem(item: TrendsBarChartDataItem): void {
    let itemId = '';

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];

        if (!props.hiddenField || item[props.hiddenField]) {
            continue;
        }

        const id = (props.idField && item[props.idField]) ? item[props.idField] as string : getItemName(item[props.nameField] as string);

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
</script>

<style>
.trends-bar-chart-legend {
    margin-right: 4px;
    cursor: pointer;
}

.trends-bar-chart-legend-icon.f7-icons {
    font-size: var(--ebk-trends-bar-chart-legend-icon-font-size);
    margin-right: 2px;
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
