<template>
    <v-chart autoresize class="radar-chart-container" :class="{ 'transition-in': skeleton }"
             :option="chartOptions" :update-options="{ notMerge: true }" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { ChartValueType, type CategoricalChartSourceDataItem } from '@/core/chart.ts';

import { isNumber } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, isBigDecimal } from '@/lib/numeral.ts';
import { max } from '@/lib/math.ts';
import { getDisplayColor } from '@/lib/color.ts';

interface RadarChartData {
    totalValidValue: BigDecimal;
    maxValue: BigDecimal;
    indicators: RadarChartDataItem[];
    values: number[]; // only used for echarts rendering
    tooltip: string;
}

interface RadarChartDataItem {
    name: string;
    max: number; // only used for echarts rendering
    color: ColorStyleValue;
}

const props = defineProps<{
    skeleton?: boolean;
    items: CategoricalChartSourceDataItem[];
    valueType: ChartValueType;
    defaultCurrency?: string;
    showValue?: boolean;
    showPercent?: boolean;
    useCustomColor?: boolean;
}>();

const theme = useTheme();

const {
    formatPercentToLocalizedNumerals,
    formatChartValueToLocalizedNumerals
} = useI18n();

const settingsStore = useSettingsStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);

const radarData = computed<RadarChartData>(() => {
    let totalValidValue: BigDecimal = BIG_DECIMAL_ZERO;
    let maxValue: BigDecimal = BIG_DECIMAL_ZERO;
    const indicators: RadarChartDataItem[] = [];
    const values: number[] = [];
    let tooltip = '';

    if (props.items.length) {
        for (const item of props.items) {
            if (isBigDecimal(item.value) && item.value.isPositive() && !item.hidden) {
                totalValidValue = totalValidValue.add(item.value);

                if (item.value.greaterThan(maxValue)) {
                    maxValue = item.value;
                }
            }
        }

        for (const item of props.items) {
            if (isBigDecimal(item.value) && !item.hidden) {
                let percent: number = isNumber(item.percent) ? item.percent : -1;

                if (percent < 0) {
                    if (item.value.isPositive()) {
                        percent = item.value.divide(totalValidValue).multiply(100).toDoubleNumber();
                    } else {
                        percent = 0;
                    }
                }

                const color = getDisplayColor(props.useCustomColor && item.color ? item.color : chartColors.value[indicators.length % chartColors.value.length]);
                const displayValue = formatChartValueToLocalizedNumerals(item.value, props.valueType, props.defaultCurrency);
                const displayPercent = formatPercentToLocalizedNumerals(percent, 2, '<0.01');

                indicators.push({
                    name: item.name,
                    max: maxValue.toDoubleNumber(),
                    color: isDarkMode.value ? '#ccc' : '#333'
                });

                values.push(max(item.value, BIG_DECIMAL_ZERO).toDoubleNumber());

                tooltip += '<div><span class="chart-pointer" style="background-color: ' + color + '"></span>';
                tooltip += `<span>${item.name}</span>`;

                const showValue = props.showValue;
                const showPercent = props.showPercent && item.value.isPositive();

                if (showValue && showPercent) {
                    tooltip += `<span class="ms-1" style="float: inline-end">(${displayPercent})</span><span class="ms-5" style="float: inline-end">${displayValue}</span>`;
                } else if (showValue && !showPercent) {
                    tooltip += `<span class="ms-5" style="float: inline-end">${displayValue}</span>`;
                } else if (!showValue && showPercent) {
                    tooltip += `<span class="ms-5" style="float: inline-end">${displayPercent}</span>`;
                }

                tooltip += '</div>';
            }
        }
    } else {
        for (let i = 0; i < 6; i++) {
            indicators.push({
                name: '',
                max: 0,
                color: isDarkMode.value ? '#ccc' : '#333'
            });
            values.push(0);
        }
    }

    const ret: RadarChartData = {
        totalValidValue: totalValidValue,
        maxValue: maxValue,
        indicators: indicators,
        values: values,
        tooltip: tooltip
    };

    return ret;
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
            formatter: () => radarData.value.tooltip
        },
        radar: {
            radius: '75%',
            splitNumber: (!props.skeleton && props.items.length) ? 5 : 1,
            splitLine: {
                lineStyle: {
                    color: (!props.skeleton && props.items.length) ? '#e8e8e7' : '#d3d3d3'
                }
            },
            splitArea: {
                areaStyle: {
                    color: (!props.skeleton && props.items.length) ? (isDarkMode.value ? ['#363534', '#1a1a1a'] : ['#faf8f4', '#fff']) : ['#d3d3d3', '#d3d3d3']
                }
            },
            indicator: radarData.value.indicators
        },
        series: (!props.skeleton && props.items.length) ? [
            {
                type: 'radar',
                data: [
                    {
                        value: radarData.value.values,
                        itemStyle: {
                            color: '#c07d43'
                        },
                        lineStyle: {
                            color: '#c07d43'
                        },
                        areaStyle: {
                            color: isDarkMode.value ? '#c07d4380' : '#c07d4340'
                        }
                    }
                ],
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
        ] : []
    };
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
