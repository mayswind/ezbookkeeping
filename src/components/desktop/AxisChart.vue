<template>
    <v-chart autoresize :class="finalClass"
             :option="chartOptions" :update-options="{ notMerge: true }"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import { itemAndIndex } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { type AxisChartSourceDataItem, ChartValueType } from '@/core/chart.ts';

import type { SortableTransactionStatisticDataItem } from '@/models/transaction.ts';

import { isArray, isNumber } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, BIG_DECIMAL_NEGATIVE_INFINITY, BIG_DECIMAL_POSITIVE_INFINITY, parseBigDecimal } from '@/lib/numeral.ts';
import { getDisplayColor } from '@/lib/color.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';

export type AxisChartDisplayType = 'line' | 'area' | 'column' | 'bubble';

interface AxisChartData {
    allSeries: AxisChartDataItem[];
    allOriginalData: BigDecimal[][];
}

interface AxisChartDataItem {
    id: string;
    name: string;
    itemStyle: {
        color: ColorStyleValue;
    };
    selected: boolean;
    type: string;
    areaStyle?: object;
    stack?: string;
    symbolSize?: (data: number) => number;
    animation: boolean;
    data: number[];  // only used for echarts rendering, the actual value is in allOriginalData
}

interface AxisChartTooltipItem extends SortableTransactionStatisticDataItem {
    readonly id: string;
    readonly name: string;
    readonly color: unknown;
    readonly displayOrders: number[];
    readonly value: BigDecimal;
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    type: AxisChartDisplayType;
    stacked?: boolean;
    oneHundredPercentStacked?: boolean;
    sortingType: number;
    showValue?: boolean;
    showTotalAmountInTooltip?: boolean;
    totalNameInTooltip?: string;
    categoryTypeName: string;
    allCategoryNames: string[];
    items: AxisChartSourceDataItem[];
    valueType: ChartValueType;
    translateName?: boolean;
    defaultCurrency?: string;
    useCustomColor?: boolean;
    enableClickItem?: boolean;
    tooltipExtraColumnNames?: string[];
    tooltipExtraColumnTotalValues?: (categoryIndex: number, totalValue: BigDecimal, visibleSeriesIds: string[]) => string[];
    tooltipExtraColumnValues?: (seriesId: string, categoryIndex: number, currentValue: BigDecimal) => string[];
}>();

const emit = defineEmits<{
    (e: 'click', itemId: string, categoryIndex: number, item: AxisChartSourceDataItem): void;
}>();

const theme = useTheme();

const {
    tt,
    getCurrentLanguageTextDirection,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping,
    formatChartValueToLocalizedNumerals
} = useI18n();

const settingsStore = useSettingsStore();

const selectedLegends = ref<Record<string, boolean>>({});

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);

const finalClass = computed<string>(() => {
    let finalClass = '';

    if (props.skeleton) {
        finalClass += 'transition-in';
    }

    if (props.class) {
        finalClass += ` ${props.class}`;
    } else {
        finalClass += ' axis-chart-container';
    }

    return finalClass;
});

const allItemsMap = computed<Record<string, AxisChartSourceDataItem>>(() => {
    const map: Record<string, AxisChartSourceDataItem> = {};

    for (const item of props.items) {
        map[item.id ?? getItemName(item.name)] = item;
    }

    return map;
});

const axisChartData = computed<AxisChartData>(() => {
    const allSeries: AxisChartDataItem[] = [];
    const allOriginalData: BigDecimal[][] = [];
    const categoryTotalAmount: Record<number, BigDecimal> = {};
    let maxAmountOfAllData: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of props.items) {
        if (item.hidden) {
            continue;
        }

        if (!isArray(item.values)) {
            continue;
        }

        const allAmounts: BigDecimal[] = item.values;

        for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
            let totalAmount: BigDecimal = categoryTotalAmount[categoryIndex] ?? BIG_DECIMAL_ZERO;
            totalAmount = totalAmount.add(amount);
            categoryTotalAmount[categoryIndex] = totalAmount;

            if (amount.greaterThan(maxAmountOfAllData)) {
                maxAmountOfAllData = amount;
            }
        }
    }

    for (const item of props.items) {
        if (item.hidden) {
            continue;
        }

        if (!isArray(item.values)) {
            continue;
        }

        const allAmounts: BigDecimal[] = item.values;

        if (props.oneHundredPercentStacked) {
            for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
                const totalAmount: BigDecimal = categoryTotalAmount[categoryIndex] ?? BIG_DECIMAL_ZERO;
                allAmounts[categoryIndex] = !totalAmount.isZero() ? amount.divide(totalAmount).multiply(100) : BIG_DECIMAL_ZERO;
            }
        }

        const finalItem: AxisChartDataItem = {
            id: item.id ?? getItemName(item.name),
            name: item.id ?? getItemName(item.name),
            itemStyle: {
                color: getDisplayColor(props.useCustomColor && item.color ? item.color : chartColors.value[allSeries.length % chartColors.value.length]),
            },
            selected: true,
            type: 'line',
            animation: !props.skeleton,
            data: allAmounts.map(amount => amount.toDoubleNumber())
        };

        if (props.stacked) {
            finalItem.stack = 'a';
        } else if (item.id) {
            finalItem.stack = item.id;
        }

        if (props.type === 'line') {
            finalItem.areaStyle = undefined;
        } else if (props.type === 'area') {
            finalItem.areaStyle = {};
        } else if (props.type === 'column') {
            finalItem.type = 'bar';
        } else if (props.type === 'bubble') {
            finalItem.type = 'scatter';
            finalItem.symbolSize = (data: number): number => {
                return Math.sqrt(data) / Math.sqrt(maxAmountOfAllData.toDoubleNumber()) * 80 + 5;
            }
        }

        allSeries.push(finalItem);
        allOriginalData.push(allAmounts);
    }

    return {
        allSeries: allSeries,
        allOriginalData: allOriginalData
    };
});

const yAxisWidth = computed<number>(() => {
    let maxValue: BigDecimal = BIG_DECIMAL_NEGATIVE_INFINITY;
    let minValue: BigDecimal = BIG_DECIMAL_POSITIVE_INFINITY;
    let width = 90;

    if (!axisChartData.value || !axisChartData.value.allOriginalData || !axisChartData.value.allOriginalData.length) {
        return width;
    }

    for (const seriesData of axisChartData.value.allOriginalData) {
        for (const value of seriesData) {
            if (value.greaterThan(maxValue)) {
                maxValue = value;
            }

            if (value.lessThan(minValue)) {
                minValue = value;
            }
        }
    }

    const maxValueText = getDisplayValue(maxValue);
    const minValueText = getDisplayValue(minValue);
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
                let totalAmount: BigDecimal = BIG_DECIMAL_ZERO;
                let actualDisplayItemCount = 0;
                const displayItems: AxisChartTooltipItem[] = [];
                const categoryIndex = params.length > 0 && params[0] ? (params[0].dataIndex ?? 0) : 0;
                const visibleSeriesIds: string[] = [];

                for (const param of params) {
                    if (!isNumber(param.seriesIndex)) {
                        continue;
                    }

                    const seriesIndex = param.seriesIndex;
                    const seriesData = axisChartData.value.allSeries[seriesIndex];

                    if (!seriesData) {
                        continue;
                    }

                    const seriesId = seriesData.id;
                    const name = allItemsMap.value[seriesId] ? getItemName(allItemsMap.value[seriesId].name) : seriesId;
                    const color = seriesData.itemStyle.color;
                    const displayOrders = allItemsMap.value[seriesId]?.displayOrders ?? [0];
                    const amount: BigDecimal = axisChartData.value.allOriginalData[seriesIndex]?.[categoryIndex] ?? BIG_DECIMAL_ZERO;

                    displayItems.push({
                        id: seriesId,
                        name: name,
                        color: color,
                        displayOrders: displayOrders,
                        value: amount
                    });

                    visibleSeriesIds.push(seriesId);
                    totalAmount = totalAmount.add(amount);
                }

                sortStatisticsItems(displayItems, props.sortingType);

                const extraColumnValuesMap: Record<number, string[]> = {};
                const extraColumnTotalValues: string[] = [];
                const hasExtraColumnIndexes: Record<number, boolean> = {};

                if (props.tooltipExtraColumnNames) {
                    if (props.tooltipExtraColumnValues) {
                        for (const [item, index] of itemAndIndex(displayItems)) {
                            const values = props.tooltipExtraColumnValues(item.id, categoryIndex, item.value);
                            extraColumnValuesMap[index] = values;

                            for (const [value, columnIndex] of itemAndIndex(values)) {
                                if (value && value !== '-') {
                                    hasExtraColumnIndexes[columnIndex] = true;
                                }
                            }
                        }
                    }

                    if (props.tooltipExtraColumnTotalValues) {
                        const values = props.tooltipExtraColumnTotalValues(categoryIndex, totalAmount, visibleSeriesIds);
                        extraColumnTotalValues.push(...values);

                        for (const [value, columnIndex] of itemAndIndex(values)) {
                            if (value && value !== '-') {
                                hasExtraColumnIndexes[columnIndex] = true;
                            }
                        }
                    }
                }

                for (const [item, index] of itemAndIndex(displayItems)) {
                    if (displayItems.length === 1 || !item.value.isZero()) {
                        const value = getDisplayValue(item.value);
                        tooltip += '<tr><td><span class="chart-pointer" style="background-color: ' + item.color + '"></span>';
                        tooltip += `<span>${item.name}</span></td><td><span class="ms-5" style="float: inline-end">${value}</span></td>`;

                        if (props.tooltipExtraColumnNames) {
                            const values = extraColumnValuesMap[index] ?? [];

                            for (let i = 0; i < props.tooltipExtraColumnNames.length; i++) {
                                if (!hasExtraColumnIndexes[i]) {
                                    continue;
                                }

                                const value = values[i] ?? '-';
                                tooltip += `<td><span class="ms-5" style="float: inline-end">${value}</span></td>`;
                            }
                        }

                        tooltip += '</tr>';
                        actualDisplayItemCount++;
                    }
                }

                if (props.showTotalAmountInTooltip && !props.oneHundredPercentStacked) {
                    const displayTotalAmount = getDisplayValue(totalAmount);
                    let totalColumnCount = 2;

                    let totalTooltip = `<tr><td><span class="chart-pointer" style="background-color: ${isDarkMode.value ? '#eee' : '#333'}"></span>`
                        + `<span>${props.totalNameInTooltip}</span></td><td><span class="ms-5" style="float: inline-end">${displayTotalAmount}</span></td>`;

                    if (props.tooltipExtraColumnNames) {
                        for (let i = 0; i < props.tooltipExtraColumnNames.length; i++) {
                            if (!hasExtraColumnIndexes[i]) {
                                continue;
                            }

                            const value = extraColumnTotalValues[i] ?? '-';
                            totalTooltip += `<td><span class="ms-5" style="float: inline-end">${value}</span></td>`;
                            totalColumnCount++;
                        }
                    }

                    totalTooltip += '</tr>';
                    totalTooltip += `<tr><td colspan="${totalColumnCount}" ${actualDisplayItemCount > 0 ? 'style="border-bottom: ' + (isDarkMode.value ? '#eee' : '#333') + ' dashed 1px"' : ''}></td></tr>`;
                    tooltip = totalTooltip + tooltip;
                }

                if (params.length && params[0] && params[0].name) {
                    let tooltipHeader = `<td>${params[0].name}</td><td></td>`;

                    if (props.tooltipExtraColumnNames) {
                        for (const [columnName, columnIndex] of itemAndIndex(props.tooltipExtraColumnNames)) {
                            if (!hasExtraColumnIndexes[columnIndex]) {
                                continue;
                            }

                            tooltipHeader += `<td><span class="ms-5" style="float: inline-end">${columnName}</span></td>`;
                        }
                    }

                    tooltip = `<table class="chart-tooltip-table"><tbody><tr>${tooltipHeader}</tr>${tooltip}</tbody></table>`
                }

                return tooltip;
            }
        },
        legend: {
            orient: 'horizontal',
            type: 'scroll',
            top: 0,
            data: axisChartData.value.allSeries.map(item => item.name),
            selected: selectedLegends.value,
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (id: string) => allItemsMap.value[id] ? getItemName(allItemsMap.value[id].name) : id
        },
        grid: {
            left: yAxisWidth.value,
            right: 20,
            bottom: 40
        },
        xAxis: [
            {
                type: 'category',
                data: props.allCategoryNames,
                inverse: textDirection.value === TextDirection.RTL,
                axisLabel: {
                    color: isDarkMode.value ? '#888' : '#666'
                }
            }
        ],
        yAxis: [
            {
                type: 'value',
                min: props.oneHundredPercentStacked ? 0 : undefined,
                max: props.oneHundredPercentStacked ? 100 : undefined,
                axisLabel: {
                    color: isDarkMode.value ? '#888' : '#666',
                    formatter: (value: number) => {
                        return getDisplayValue(parseBigDecimal(value));
                    }
                },
                axisPointer: {
                    label: {
                        formatter: (params: CallbackDataParams) => {
                            return getDisplayValue(parseBigDecimal(params.value as number).truncate());
                        }
                    }
                },
                splitLine: {
                    lineStyle: {
                        color: isDarkMode.value ? '#4f4f4f' : '#e1e6f2',
                    }
                }
            }
        ],
        series: axisChartData.value.allSeries
    };
});

function getItemName(name: string): string {
    return props.translateName ? tt(name) : name;
}

function getDisplayValue(value: BigDecimal): string {
    if (props.oneHundredPercentStacked) {
        return formatChartValueToLocalizedNumerals(value, ChartValueType.Percent, props.defaultCurrency);
    } else {
        return formatChartValueToLocalizedNumerals(value, props.valueType, props.defaultCurrency);
    }
}

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series') {
        return;
    }

    const id = e.seriesId as string;
    const item = allItemsMap.value[id] as AxisChartSourceDataItem;
    const itemId = item?.id ?? '';
    const category = props.allCategoryNames[e.dataIndex];

    if (!item || !category) {
        return;
    }

    emit('click', itemId, e.dataIndex, item);
}

function exportData(): { headers: string[], data: string[][] } {
    const headers: string[] = [];
    const data: string[][] = [];

    headers.push(props.categoryTypeName);

    for (const series of axisChartData.value.allSeries) {
        const id = series.id;
        const name = allItemsMap.value[id] ? getItemName(allItemsMap.value[id].name) : id;
        headers.push(name);
    }

    for (const [categoryName, index] of itemAndIndex(props.allCategoryNames)) {
        const row: string[] = [];
        row.push(categoryName);
        row.push(...axisChartData.value.allOriginalData.map(item => {
            if (props.oneHundredPercentStacked) {
                return formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping(item[index] ?? BIG_DECIMAL_ZERO);
            } else if (props.valueType === ChartValueType.Amount) {
                return formatAmountToWesternArabicNumeralsWithoutDigitGrouping(item[index] ?? BIG_DECIMAL_ZERO, props.defaultCurrency);
            } else {
                return formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping(item[index] ?? BIG_DECIMAL_ZERO);
            }
        }));
        data.push(row);
    }

    return {
        headers: headers,
        data: data
    };
}

function onLegendSelectChanged(e: { selected: Record<string, boolean> }): void {
    selectedLegends.value = e.selected;
}

defineExpose({
    exportData
});
</script>

<style scoped>
.axis-chart-container {
    width: 100%;
    height: 560px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .axis-chart-container {
        height: 650px;
    }
}
</style>
