<template>
    <v-chart autoresize class="trends-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonTrendsChartProps, type TrendsBarChartClickEvent, useTrendsChartBase } from '@/components/base/TrendsChartBase.ts'

import { useUserStore } from '@/stores/user.ts';

import { type YearMonth, DateRangeScene } from '@/core/datetime.ts';
import type { ColorValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { TrendChartType, ChartDateAggregationType } from '@/core/statistics.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';
import type { YearMonthDataItem, SortableTransactionStatisticDataItem } from '@/models/transaction.ts';

import {
    isArray,
    isNumber
} from '@/lib/common.ts';
import {
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange
} from '@/lib/datetime.ts';
import {
    sortStatisticsItems
} from '@/lib/statistics.ts';

interface DesktopTrendsChartProps<T extends YearMonth> extends CommonTrendsChartProps<T> {
    skeleton?: boolean;
    type: number;
    showValue?: boolean;
    showTotalAmountInTooltip?: boolean;
}

interface TrendsChartDataItem {
    id: string;
    name: string;
    itemStyle: {
        color: ColorValue;
    };
    selected: boolean;
    type: string;
    areaStyle?: object;
    stack: string;
    animation: boolean;
    data: number[];
}

interface TrendsChartTooltipItem extends SortableTransactionStatisticDataItem {
    readonly name: string;
    readonly color: unknown;
    readonly displayOrders: number[];
    readonly totalAmount: number;
}

const props = defineProps<DesktopTrendsChartProps<YearMonthDataItem>>();

const emit = defineEmits<{
    (e: 'click', value: TrendsBarChartClickEvent): void;
}>();

const theme = useTheme();
const { tt, formatUnixTimeToShortYear, formatYearQuarter, formatUnixTimeToShortYearMonth, formatAmountWithCurrency } = useI18n();
const { allDateRanges, getItemName, getColor } = useTrendsChartBase(props);

const userStore = useUserStore();

const selectedLegends = ref<Record<string, boolean>>({});

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const itemsMap = computed<Record<string, Record<string, unknown>>>(() => {
    const map: Record<string, Record<string, unknown>> = {};

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];
        let id: string = '';

        if (props.idField && item[props.idField]) {
            id = item[props.idField] as string;
        } else {
            id = getItemName(item[props.nameField] as string);
        }

        const finalItem: Record<string, unknown> = {
            [props.nameField]: item[props.nameField]
        };

        if (props.idField) {
            finalItem[props.idField] = item[props.idField];
        }

        if (props.hiddenField) {
            finalItem[props.hiddenField] = item[props.hiddenField];
        }

        if (props.displayOrdersField) {
            finalItem[props.displayOrdersField] = item[props.displayOrdersField];
        }

        map[id] = finalItem;
    }

    return map;
});

const allDisplayDateRanges = computed<string[]>(() => {
    const allDisplayDateRanges: string[] = [];

    for (let i = 0; i < allDateRanges.value.length; i++) {
        const dateRange = allDateRanges.value[i];

        if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
            allDisplayDateRanges.push(formatUnixTimeToShortYear(dateRange.minUnixTime));
        } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
            allDisplayDateRanges.push(formatYearQuarter(dateRange.year, dateRange.quarter));
        } else { // if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
            allDisplayDateRanges.push(formatUnixTimeToShortYearMonth(dateRange.minUnixTime));
        }
    }

    return allDisplayDateRanges;
});

const allSeries = computed<TrendsChartDataItem[]>(() => {
    const allSeries: TrendsChartDataItem[] = [];

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];

        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        const allAmounts: number[] = [];
        const dateRangeAmountMap: Record<string, YearMonthDataItem[]> = {};

        for (let j = 0; j < item.items.length; j++) {
            const dataItem = item.items[j];
            let dateRangeKey = '';

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                dateRangeKey = dataItem.year.toString();
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month - 1) / 3) + 1}`;
            } else { // if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                dateRangeKey = `${dataItem.year}-${dataItem.month}`;
            }

            const dataItems = dateRangeAmountMap[dateRangeKey] || [];
            dataItems.push(dataItem);

            dateRangeAmountMap[dateRangeKey] = dataItems;
        }

        for (let j = 0; j < allDateRanges.value.length; j++) {
            const dateRange = allDateRanges.value[j];
            let dateRangeKey = '';

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                dateRangeKey = dateRange.year.toString();
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
                dateRangeKey = `${dateRange.year}-${dateRange.quarter}`;
            } else if (props.dateAggregationType === ChartDateAggregationType.Month.type && 'month' in dateRange) {
                dateRangeKey = `${dateRange.year}-${dateRange.month + 1}`;
            }

            let amount = 0;
            const dataItems = dateRangeAmountMap[dateRangeKey];

            if (isArray(dataItems)) {
                for (let i = 0; i < dataItems.length; i++) {
                    const dataItem = dataItems[i];

                    if (isNumber(dataItem[props.valueField])) {
                        amount += dataItem[props.valueField];
                    }
                }
            }

            allAmounts.push(amount);
        }

        const finalItem: TrendsChartDataItem = {
            id: (props.idField && item[props.idField]) ? item[props.idField] as string : getItemName(item[props.nameField] as string),
            name: (props.idField && item[props.idField]) ? item[props.idField] as string : getItemName(item[props.nameField] as string),
            itemStyle: {
                color: getColor(props.colorField && item[props.colorField] ? item[props.colorField] as string : DEFAULT_CHART_COLORS[i % DEFAULT_CHART_COLORS.length]),
            },
            selected: true,
            type: 'line',
            stack: 'a',
            animation: !props.skeleton,
            data: allAmounts
        };

        if (props.type === TrendChartType.Area.type) {
            finalItem.areaStyle = {};
        } else if (props.type === TrendChartType.Column.type) {
            finalItem.type = 'bar';
        }

        allSeries.push(finalItem);
    }

    return allSeries;
});

const yAxisWidth = computed<number>(() => {
    let maxValue = Number.MIN_SAFE_INTEGER;
    let minValue = Number.MAX_SAFE_INTEGER;
    let width = 90;

    if (!allSeries.value || !allSeries.value.length) {
        return width;
    }

    for (let i = 0; i < allSeries.value.length; i++) {
        for (let j = 0; j < allSeries.value[i].data.length; j++) {
            const value = allSeries.value[i].data[j];

            if (value > maxValue) {
                maxValue = value;
            }

            if (value < minValue) {
                minValue = value;
            }
        }
    }

    const maxValueText = formatAmountWithCurrency(maxValue, props.defaultCurrency);
    const minValueText = formatAmountWithCurrency(minValue, props.defaultCurrency);
    const maxLengthText = maxValueText.length > minValueText.length ? maxValueText : minValueText;

    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');

    if (context) {
        context.font = '12px Arial';

        const textMetrics = context.measureText(maxLengthText);
        const actualWidth = Math.round(textMetrics.width) + 20;

        if (actualWidth >= 200) {
            width = 200;
        } if (actualWidth > 90) {
            width = actualWidth;
        }
    }

    return width;
});

const chartOptions = computed<object>(() => {
    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'cross',
                label: {
                    backgroundColor: isDarkMode.value ? '#333' : '#fff',
                    color: isDarkMode.value ? '#eee' : '#333'
                },
            },
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams[]) => {
                let tooltip = '';
                let totalAmount = 0;
                const displayItems: TrendsChartTooltipItem[] = [];

                for (let i = 0; i < params.length; i++) {
                    const id = params[i].seriesId as string;
                    const name = itemsMap.value[id] && props.nameField && itemsMap.value[id][props.nameField] ? getItemName(itemsMap.value[id][props.nameField] as string) : id;
                    const color = params[i].color;
                    const displayOrders = itemsMap.value[id] && props.displayOrdersField && itemsMap.value[id][props.displayOrdersField] ? itemsMap.value[id][props.displayOrdersField] as number[] : [0];
                    const amount = params[i].data as number;

                    displayItems.push({
                        name: name,
                        color: color,
                        displayOrders: displayOrders,
                        totalAmount: amount
                    });

                    totalAmount += amount;
                }

                sortStatisticsItems(displayItems, props.sortingType);

                for (let i = 0; i < displayItems.length; i++) {
                    const item = displayItems[i];

                    if (displayItems.length === 1 || item.totalAmount !== 0) {
                        const value = formatAmountWithCurrency(item.totalAmount, props.defaultCurrency);
                        tooltip += '<div><span class="chart-pointer" style="background-color: ' + item.color + '"></span>';
                        tooltip += `<span>${item.name}</span><span style="margin-left: 20px; float: right">${value}</span><br/>`;
                        tooltip += '</div>';
                    }
                }

                if (props.showTotalAmountInTooltip) {
                    const displayTotalAmount = formatAmountWithCurrency(totalAmount, props.defaultCurrency);
                    tooltip = '<div style="border-bottom: ' + (isDarkMode.value ? '#eee' : '#333') + ' dashed 1px">'
                        + '<span class="chart-pointer" style="background-color: ' + (isDarkMode.value ? '#eee' : '#333') + '"></span>'
                        + `<span>${tt('Total Amount')}</span><span style="margin-left: 20px; float: right">${displayTotalAmount}</span><br/>`
                        + '</div>' + tooltip;
                }

                if (params.length && params[0].name) {
                    tooltip = `${params[0].name}<br/>` + tooltip;
                }

                return tooltip;
            }
        },
        legend: {
            orient: 'horizontal',
            data: allSeries.value.map(item => item.name),
            selected: selectedLegends.value,
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (id: string) => {
                return itemsMap.value[id] && props.nameField && itemsMap.value[id][props.nameField] ? getItemName(itemsMap.value[id][props.nameField] as string) : id;
            }
        },
        grid: {
            left: yAxisWidth.value,
            right: 20
        },
        xAxis: [
            {
                type: 'category',
                data: allDisplayDateRanges.value
            }
        ],
        yAxis: [
            {
                type: 'value',
                axisLabel: {
                    formatter: (value: string) => {
                        return formatAmountWithCurrency(value, props.defaultCurrency);
                    }
                },
                axisPointer: {
                    label: {
                        formatter: (params: CallbackDataParams) => {
                            return formatAmountWithCurrency(Math.floor(params.value as number), props.defaultCurrency);
                        }
                    }
                }
            }
        ],
        series: allSeries.value
    };
});

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series') {
        return;
    }

    const id = e.seriesId as string;
    const item = itemsMap.value[id];
    const itemId = props.idField ? item[props.idField] as string : '';
    const dateRange = allDateRanges.value[e.dataIndex];
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

function onLegendSelectChanged(e: { selected: Record<string, boolean> }): void {
    selectedLegends.value = e.selected;
}
</script>

<style scoped>
.trends-chart-container {
    width: 100%;
    height: 560px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .pie-chart-container {
        height: 500px;
    }
}
</style>
