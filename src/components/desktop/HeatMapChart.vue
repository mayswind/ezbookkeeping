<template>
    <v-chart autoresize :class="finalClass" :style="finalStyle"
             :option="chartOptions" :update-options="{ notMerge: true }"
             @click="clickItem" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { itemAndIndex } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import { ThemeType } from '@/core/theme.ts';
import { type AxisChartSourceDataItem, ChartValueType } from '@/core/chart.ts';

import { isArray } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, BIG_DECIMAL_POSITIVE_INFINITY, parseBigDecimal } from '@/lib/numeral.ts';

interface HeatMapData {
    allSeriesNames: string[];
    allOriginalDataMap: Record<string, BigDecimal>;
    data: [number, number, number][]; // third value only used for echarts rendering, the actual value is in allOriginalDataMap
    minValue: BigDecimal;
    maxValue: BigDecimal;
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    showValue?: boolean;
    enableClickItem?: boolean;
    categoryTypeName: string;
    allCategoryNames: string[];
    items: AxisChartSourceDataItem[];
    valueType: ChartValueType;
    valueTypeName: string;
    translateName?: boolean;
    defaultCurrency?: string;
}>();

const emit = defineEmits<{
    (e: 'click', categoryIndex: number, seriesIndex: number): void;
}>();

const theme = useTheme();

const {
    tt,
    getCurrentLanguageTextDirection,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping,
    formatChartValueToLocalizedNumerals
} = useI18n();

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
        finalClass += ' heatmap-chart-container';
    }

    return finalClass;
});
const finalStyle = computed<Record<string, string>>(() => {
    const style: Record<string, string> = {};

    if (heatMapData.value.allSeriesNames && heatMapData.value.allSeriesNames.length > 15) {
        style['height'] = `${heatMapData.value.allSeriesNames.length * 40}px`;
    }

    return style;
});

const heatMapData = computed<HeatMapData>(() => {
    const allData: [number, number, number][] = [];
    const allOriginalDataMap: Record<string, BigDecimal> = {};
    const allSeriesNames: string[] = [];
    let minValue: BigDecimal = BIG_DECIMAL_POSITIVE_INFINITY;
    let maxValue: BigDecimal = BIG_DECIMAL_ZERO;

    for (const [item, seriesIndex] of itemAndIndex(props.items)) {
        if (item.hidden) {
            continue;
        }

        if (!isArray(item.values)) {
            continue;
        }

        allSeriesNames.push(props.translateName ? tt(item.name) : item.name);

        for (const [amount, categoryIndex] of itemAndIndex(item.values)) {
            if (amount.greaterThan(maxValue)) {
                maxValue = amount;
            }

            if (amount.lessThan(minValue)) {
                minValue = amount;
            }

            allOriginalDataMap[`${categoryIndex}-${seriesIndex}`] = amount;
            allData.push([categoryIndex, seriesIndex, amount.toDoubleNumber()]);
        }
    }

    const ret: HeatMapData = {
        allSeriesNames: allSeriesNames,
        allOriginalDataMap: allOriginalDataMap,
        data: allData,
        minValue: minValue.isPositiveInfinity() ? BIG_DECIMAL_ZERO : minValue,
        maxValue: maxValue
    };

    return ret;
});

const yAxisWidth = computed<number>(() => {
    let width: number = 60;

    if (!heatMapData.value || !heatMapData.value.allSeriesNames) {
        return width;
    }

    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');

    if (context) {
        context.font = '12px Arial';

        for (const seriesName of heatMapData.value.allSeriesNames) {
            const textMetrics = context.measureText(seriesName);
            const actualWidth = Math.round(textMetrics.width) + 20;

            if (actualWidth > width) {
                width = actualWidth;
            }
        }
    }

    if (width >= 200) {
        width = 200;
    }

    return width;
});

const chartOptions = computed<object>(() => {
    return {
        tooltip: {
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams) => {
                if (!props.showValue) {
                    return '';
                }

                const dataItem = params.data as [number, number, number];
                const name = props.valueTypeName;
                const displayValue: string = formatDataItemDisplayValue(dataItem);

                return `<div class="d-inline-flex">${params.name}</div><br/>`
                    + `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`
                    + `<span>${name}</span>`
                    + `<span class="ms-5">${displayValue}</span>`
                    + '</div>';
            }
        },
        visualMap: [
            {
                type: 'continuous',
                orient: 'horizontal',
                top: 0,
                left: 'center',
                itemHeight: 320,
                min: heatMapData.value.minValue.toDoubleNumber(),
                max: heatMapData.value.maxValue.toDoubleNumber(),
                calculable: true,
                inRange: {
                    color: isDarkMode.value ? [ '#1a1a1a', '#c67e48' ] : [ '#faf8f4', '#c67e48' ]
                },
                textStyle: {
                    color: isDarkMode.value ? '#888' : '#666'
                },
                formatter: (value: number) => {
                    if (!props.showValue) {
                        return '';
                    }

                    let actualValue: BigDecimal = parseBigDecimal(value);

                    if (value === heatMapData.value.minValue.toDoubleNumber()) {
                        actualValue = heatMapData.value.minValue;
                    } else if (value === heatMapData.value.maxValue.toDoubleNumber()) {
                        actualValue = heatMapData.value.maxValue;
                    }

                    return formatChartValueToLocalizedNumerals(actualValue, props.valueType, props.defaultCurrency);
                }
            }
        ],
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
                type: 'category',
                data: heatMapData.value.allSeriesNames,
                inverse: true,
                axisLabel: {
                    color: isDarkMode.value ? '#888' : '#666'
                }
            }
        ],
        series: [
            {
                type: 'heatmap',
                animation: !props.skeleton,
                data: heatMapData.value.data,
                label: {
                    show: props.showValue ?? false,
                    color: isDarkMode.value ? '#eee' : '#333',
                    formatter: (params: CallbackDataParams) => {
                        if (!props.showValue) {
                            return '';
                        }

                        const data: [number, number, number] = params.data as [number, number, number];
                        return formatDataItemDisplayValue(data);
                    }
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 6,
                        shadowColor: isDarkMode.value ? 'rgba(255, 255, 255, 0.5)' : 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    };
});

function formatDataItemDisplayValue(dataItem: [number, number, number]): string {
    const categoryIndex = dataItem[0];
    const seriesIndex = dataItem[1];
    const value: BigDecimal | undefined = heatMapData.value.allOriginalDataMap[`${categoryIndex}-${seriesIndex}`];
    return value ? formatChartValueToLocalizedNumerals(value, props.valueType, props.defaultCurrency) : '0';
}

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series') {
        return;
    }

    const dataItem = e.data as [number, number, number];

    if (!dataItem) {
        return;
    }

    const categoryIndex = dataItem[0];
    const seriesIndex = dataItem[1];
    emit('click', categoryIndex, seriesIndex);
}

function exportData(): { headers: string[], data: string[][] } {
    const headers: string[] = [];
    const data: string[][] = [];

    headers.push(props.categoryTypeName);

    for (const categoryName of props.allCategoryNames) {
        headers.push(categoryName);
    }

    for (const [seriesName, seriesIndex] of itemAndIndex(heatMapData.value.allSeriesNames)) {
        const row: string[] = [];
        row.push(seriesName);
        for (let categoryIndex = 0; categoryIndex < props.allCategoryNames.length; categoryIndex++) {
            const value: BigDecimal = heatMapData.value.allOriginalDataMap[`${categoryIndex}-${seriesIndex}`] ?? BIG_DECIMAL_ZERO;

            if (props.valueType === ChartValueType.Amount) {
                row.push(formatAmountToWesternArabicNumeralsWithoutDigitGrouping(value, props.defaultCurrency));
            } else {
                row.push(formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping(value));
            }
        }
        data.push(row);
    }

    return {
        headers: headers,
        data: data
    };
}

defineExpose({
    exportData
});
</script>

<style scoped>
.heatmap-chart-container {
    width: 100%;
    height: 560px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .heatmap-chart-container {
        height: 630px;
    }
}
</style>
