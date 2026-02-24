<template>
    <v-chart autoresize :class="finalClass" :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { itemAndIndex } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import type { ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';

import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

import type { SortableTransactionStatisticDataItem } from '@/models/transaction.ts';

import { isArray } from '@/lib/common.ts';
import { getDisplayColor } from '@/lib/color.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';

export type AxisChartDisplayType = 'line' | 'area' | 'column' | 'bubble';

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
    data: number[];
}

interface AxisChartTooltipItem extends SortableTransactionStatisticDataItem {
    readonly name: string;
    readonly color: unknown;
    readonly displayOrders: number[];
    readonly totalAmount: number;
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
    items: Record<string, unknown>[];
    idField?: string;
    nameField: string;
    valuesField: string;
    colorField?: string;
    hiddenField?: string;
    displayOrdersField?: string;
    translateName?: boolean;
    amountValue?: boolean;
    defaultCurrency?: string;
    enableClickItem?: boolean;
}>();

const emit = defineEmits<{
    (e: 'click', itemId: string, categoryIndex: number, item: Record<string, unknown>): void;
}>();

const theme = useTheme();

const {
    tt,
    getCurrentLanguageTextDirection,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatNumberToLocalizedNumerals,
    formatNumberToWesternArabicNumerals,
    formatPercentToLocalizedNumerals
} = useI18n();

const selectedLegends = ref<Record<string, boolean>>({});

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
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

const itemsMap = computed<Record<string, Record<string, unknown>>>(() => {
    const map: Record<string, Record<string, unknown>> = {};

    for (const item of props.items) {
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

const allSeries = computed<AxisChartDataItem[]>(() => {
    const allSeries: AxisChartDataItem[] = [];
    const categoryTotalAmount: Record<number, number> = {};
    let maxAmountOfAllData: number = 0;

    for (const item of props.items) {
        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        if (!isArray(item[props.valuesField])) {
            continue;
        }

        const allAmounts: number[] = item[props.valuesField] as number[];

        for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
            let totalAmount: number = categoryTotalAmount[categoryIndex] ?? 0;
            totalAmount += amount;
            categoryTotalAmount[categoryIndex] = totalAmount;

            if (amount > maxAmountOfAllData) {
                maxAmountOfAllData = amount;
            }
        }
    }

    for (const item of props.items) {
        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        if (!isArray(item[props.valuesField])) {
            continue;
        }

        const allAmounts: number[] = item[props.valuesField] as number[];

        if (props.oneHundredPercentStacked) {
            for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
                const totalAmount: number = categoryTotalAmount[categoryIndex] ?? 0;
                allAmounts[categoryIndex] = totalAmount !== 0 ? amount * 100.0 / totalAmount : 0;
            }
        }

        const finalItem: AxisChartDataItem = {
            id: (props.idField && item[props.idField]) ? item[props.idField] as string : getItemName(item[props.nameField] as string),
            name: (props.idField && item[props.idField]) ? item[props.idField] as string : getItemName(item[props.nameField] as string),
            itemStyle: {
                color: getDisplayColor(props.colorField && item[props.colorField] ? item[props.colorField] as string : DEFAULT_CHART_COLORS[allSeries.length % DEFAULT_CHART_COLORS.length]),
            },
            selected: true,
            type: 'line',
            animation: !props.skeleton,
            data: allAmounts
        };

        if (props.stacked) {
            finalItem.stack = 'a';
        } else if (props.idField && item[props.idField]) {
            finalItem.stack = item[props.idField] as string;
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
                return Math.sqrt(data) / Math.sqrt(maxAmountOfAllData) * 80 + 5;
            }
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

    for (const series of allSeries.value) {
        for (const value of series.data) {
            if (value > maxValue) {
                maxValue = value;
            }

            if (value < minValue) {
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
                let totalAmount = 0;
                let actualDisplayItemCount = 0;
                const displayItems: AxisChartTooltipItem[] = [];

                for (const param of params) {
                    const id = param.seriesId as string;
                    const name = itemsMap.value[id] && props.nameField && itemsMap.value[id][props.nameField] ? getItemName(itemsMap.value[id][props.nameField] as string) : id;
                    const color = param.color;
                    const displayOrders = itemsMap.value[id] && props.displayOrdersField && itemsMap.value[id][props.displayOrdersField] ? itemsMap.value[id][props.displayOrdersField] as number[] : [0];
                    const amount = param.data as number;

                    displayItems.push({
                        name: name,
                        color: color,
                        displayOrders: displayOrders,
                        totalAmount: amount
                    });

                    totalAmount += amount;
                }

                sortStatisticsItems(displayItems, props.sortingType);

                for (const item of displayItems) {
                    if (displayItems.length === 1 || item.totalAmount !== 0) {
                        const value = getDisplayValue(item.totalAmount);
                        tooltip += '<div><span class="chart-pointer" style="background-color: ' + item.color + '"></span>';
                        tooltip += `<span>${item.name}</span><span class="ms-5" style="float: inline-end">${value}</span><br/>`;
                        tooltip += '</div>';
                        actualDisplayItemCount++;
                    }
                }

                if (props.showTotalAmountInTooltip && !props.oneHundredPercentStacked) {
                    const displayTotalAmount = getDisplayValue(totalAmount);
                    tooltip = (actualDisplayItemCount > 0 ? '<div style="border-bottom: ' + (isDarkMode.value ? '#eee' : '#333') + ' dashed 1px">' : '<div></div>')
                        + '<span class="chart-pointer" style="background-color: ' + (isDarkMode.value ? '#eee' : '#333') + '"></span>'
                        + `<span>${props.totalNameInTooltip}</span><span class="ms-5" style="float: inline-end">${displayTotalAmount}</span><br/>`
                        + '</div>' + tooltip;
                }

                if (params.length && params[0] && params[0].name) {
                    tooltip = `${params[0].name}<br/>` + tooltip;
                }

                return tooltip;
            }
        },
        legend: {
            orient: 'horizontal',
            type: 'scroll',
            top: 0,
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
                    formatter: (value: string) => {
                        return getDisplayValue(parseInt(value));
                    }
                },
                axisPointer: {
                    label: {
                        formatter: (params: CallbackDataParams) => {
                            return getDisplayValue(Math.trunc(params.value as number));
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
        series: allSeries.value
    };
});

function getItemName(name: string): string {
    return props.translateName ? tt(name) : name;
}

function getDisplayValue(value: number): string {
    if (props.oneHundredPercentStacked) {
        return formatPercentToLocalizedNumerals(value, 2, '<0.01');
    }

    if (props.amountValue) {
        return formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
    }

    return formatNumberToLocalizedNumerals(value);
}

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series') {
        return;
    }

    const id = e.seriesId as string;
    const item = itemsMap.value[id] as Record<string, unknown>;
    const itemId = props.idField ? item[props.idField] as string : '';
    const category = props.allCategoryNames[e.dataIndex];

    if (!category) {
        return;
    }

    emit('click', itemId, e.dataIndex, item);
}

function exportData(): { headers: string[], data: string[][] } {
    const headers: string[] = [];
    const data: string[][] = [];

    headers.push(props.categoryTypeName);

    for (const series of allSeries.value) {
        const id = series.id;
        const name = itemsMap.value[id] && props.nameField && itemsMap.value[id][props.nameField] ? getItemName(itemsMap.value[id][props.nameField] as string) : id;
        headers.push(name);
    }

    for (const [categoryName, index] of itemAndIndex(props.allCategoryNames)) {
        const row: string[] = [];
        row.push(categoryName);
        row.push(...allSeries.value.map(item => {
            if (props.oneHundredPercentStacked) {
                return formatNumberToWesternArabicNumerals(item.data[index] ?? 0);
            } else {
                return formatAmountToWesternArabicNumeralsWithoutDigitGrouping(item.data[index] ?? 0);
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
        height: 630px;
    }
}
</style>
