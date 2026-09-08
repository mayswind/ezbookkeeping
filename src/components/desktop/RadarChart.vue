<template>
    <v-chart autoresize class="radar-chart-container" :class="{ 'transition-in': skeleton }"
             :option="chartOptions" :update-options="{ notMerge: true }" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import { itemAndIndex } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { type AxisChartSourceDataItem, ChartValueType } from '@/core/chart.ts';

import { BIG_DECIMAL_ZERO, isBigDecimal } from '@/lib/numeral.ts';
import { max } from '@/lib/math.ts';
import { getDisplayColor } from '@/lib/color.ts';

interface RadarChartData {
    indicators: RadarChartDataItem[];
    allSeries: RadarChartSeriesDataItem[];
    seriesIdMap: Record<string, RadarChartSeriesDataItem>;
}

interface RadarChartDataItem {
    name: string;
    max: number; // only used for echarts rendering
    color: ColorStyleValue;
}

interface RadarChartSeriesDataItem {
    id: string;
    name: string;
    values: number[]; // only used for echarts rendering
    tooltip: string;
    color: ColorStyleValue;
}

const props = defineProps<{
    skeleton?: boolean;
    categoryTypeName: string;
    items: AxisChartSourceDataItem[];
    allCategoryNames?: string[];
    hideLegend?: boolean;
    valueType: ChartValueType;
    defaultCurrency?: string;
    showValue?: boolean;
    showPercent?: boolean;
    useCustomColor?: boolean;
}>();

const theme = useTheme();

const {
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping,
    formatPercentToLocalizedNumerals,
    formatChartValueToLocalizedNumerals
} = useI18n();

const settingsStore = useSettingsStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);

const radarData = computed<RadarChartData>(() => {
    let maxValue: BigDecimal = BIG_DECIMAL_ZERO;
    const indicators: RadarChartDataItem[] = [];
    const allSeries: RadarChartSeriesDataItem[] = [];
    const seriesIdMap: Record<string, RadarChartSeriesDataItem> = {};

    if (props.items.length) {
        for (const item of props.items) {
            if (!item.values || item.hidden) {
                continue;
            }

            for (const value of item.values) {
                if (isBigDecimal(value) && value.greaterThan(maxValue)) {
                    maxValue = value;
                }
            }
        }

        for (const name of props.allCategoryNames ?? []) {
            indicators.push({
                name: name,
                max: maxValue.toDoubleNumber(),
                color: isDarkMode.value ? '#ccc' : '#333'
            });
        }

        for (const item of props.items) {
            if (!item.values || item.hidden) {
                continue;
            }

            let totalValidValue: BigDecimal = BIG_DECIMAL_ZERO;
            const color = props.hideLegend ? '#c07d43' : getDisplayColor(props.useCustomColor && item.color ? item.color : chartColors.value[allSeries.length % chartColors.value.length]);
            let tooltip = props.hideLegend ? '' : `<div><span class="chart-pointer" style="background-color: ${color}"></span><span>${item.name}</span></div>`;

            for (const value of item.values) {
                if (isBigDecimal(value) && value.isPositive()) {
                    totalValidValue = totalValidValue.add(value);
                }
            }

            for (let i = 0; i < indicators.length; i++) {
                const value = item.values[i] ?? BIG_DECIMAL_ZERO;
                const percent = value.isPositive() && !totalValidValue.isZero() ? value.divide(totalValidValue).multiply(100).toDoubleNumber() : 0;
                const displayValue = formatChartValueToLocalizedNumerals(value, props.valueType, props.defaultCurrency);
                const displayPercent = formatPercentToLocalizedNumerals(percent, 2, '<0.01');

                const categoryColor = getDisplayColor(chartColors.value[i % chartColors.value.length]);
                tooltip += `<div>${props.hideLegend ? `<span class="chart-pointer" style="background-color: ${categoryColor}"></span>` : ''}<span>${indicators[i]?.name ?? ''}</span>`;

                const showValue = props.showValue;
                const showPercent = props.showPercent && value.isPositive();

                if (showValue && showPercent) {
                    tooltip += `<span class="ms-1" style="float: inline-end">(${displayPercent})</span><span class="ms-5" style="float: inline-end">${displayValue}</span>`;
                } else if (showValue && !showPercent) {
                    tooltip += `<span class="ms-5" style="float: inline-end">${displayValue}</span>`;
                } else if (!showValue && showPercent) {
                    tooltip += `<span class="ms-5" style="float: inline-end">${displayPercent}</span>`;
                }

                tooltip += '</div>';
            }

            const seriesItem: RadarChartSeriesDataItem = {
                id: item.id ?? item.name,
                name: item.name,
                values: item.values.map(value => max(value, BIG_DECIMAL_ZERO).toDoubleNumber()),
                tooltip: tooltip,
                color: color
            };

            allSeries.push(seriesItem);
            seriesIdMap[item.id ?? item.name] = seriesItem;
        }
    } else {
        for (let i = 0; i < 6; i++) {
            indicators.push({
                name: '',
                max: 0,
                color: isDarkMode.value ? '#ccc' : '#333'
            });
        }
    }

    return {
        indicators: indicators,
        allSeries: allSeries,
        seriesIdMap: seriesIdMap
    };
});

const chartOptions = computed<object>(() => {
    return {
        tooltip: {
            trigger: 'item',
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams) => {
                return radarData.value.allSeries[params.dataIndex ?? 0]?.tooltip ?? '';
            }
        },
        legend: {
            show: !props.hideLegend,
            orient: 'horizontal',
            type: 'scroll',
            top: 0,
            data: radarData.value.allSeries.map(item => item.id),
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (id: string) => radarData.value.seriesIdMap[id]?.name ?? id
        },
        radar: {
            radius: '75%',
            splitNumber: (!props.skeleton && radarData.value.allSeries.length) ? 5 : 1,
            splitLine: {
                lineStyle: {
                    color: (!props.skeleton && radarData.value.allSeries.length) ? '#e8e8e7' : '#d3d3d3'
                }
            },
            splitArea: {
                areaStyle: {
                    color: (!props.skeleton && radarData.value.allSeries.length) ? (isDarkMode.value ? ['#363534', '#1a1a1a'] : ['#faf8f4', '#fff']) : ['#d3d3d3', '#d3d3d3']
                }
            },
            indicator: radarData.value.indicators
        },
        series: (!props.skeleton && radarData.value.allSeries.length) ? [
            {
                type: 'radar',
                data: radarData.value.allSeries.map(item => ({
                    name: item.id,
                    value: item.values,
                    itemStyle: {
                        color: item.color
                    },
                    lineStyle: {
                        color: item.color
                    },
                    areaStyle: radarData.value.allSeries.length > 1 ? {
                        color: item.color,
                        opacity: isDarkMode.value ? 0.5 : 0.2
                    } : {
                        color: isDarkMode.value ? '#c07d4380' : '#c07d4340'
                    }
                })),
                top: 0,
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)',
                    }
                },
                animation: !props.skeleton
            }
        ] : [],
        media: [
            {
                query: {
                    minWidth: 600,
                },
                option: {
                    legend: {
                        orient: 'vertical',
                        left: 'left'
                    }
                }
            }
        ]
    };
});

function exportData(): { headers: string[], data: string[][] } {
    const headers: string[] = [];
    const data: string[][] = [];

    headers.push(props.categoryTypeName);

    for (const series of radarData.value.allSeries) {
        headers.push(series.name);
    }

    for (const [categoryName, index] of itemAndIndex(props.allCategoryNames ?? [''])) {
        const row: string[] = [];
        row.push(categoryName);
        row.push(...props.items.map(item => {
            if (props.valueType === ChartValueType.Amount) {
                return formatAmountToWesternArabicNumeralsWithoutDigitGrouping(item.values[index] ?? BIG_DECIMAL_ZERO, props.defaultCurrency);
            } else {
                return formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping(item.values[index] ?? BIG_DECIMAL_ZERO);
            }
        }));
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
.radar-chart-container {
    width: 100%;
    height: 460px;
}

@media (min-width: 600px) {
    .radar-chart-container {
        height: 660px;
    }
}

.radar-chart-container.transition-in {
    animation: radar-chart-skeleton-fade-in 2s 1;
}

@keyframes radar-chart-skeleton-fade-in {
    0% {
        opacity: 0;
    }
    20% {
        opacity: 0;
    }
    100% {
        opacity: 1;
    }
}
</style>
