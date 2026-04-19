<template>
    <v-chart autoresize :class="finalClass" :style="finalStyle" :option="chartOptions" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { itemAndIndex } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import { ThemeType } from '@/core/theme.ts';

import { isArray, isNumber } from '@/lib/common.ts';

interface HeatMapData {
    allSeriesNames: string[];
    data: [number, number, number][];
    minValue: number;
    maxValue: number;
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    showValue?: boolean;
    allCategoryNames: string[];
    items: Record<string, unknown>[];
    nameField: string;
    valuesField: string;
    hiddenField?: string;
    translateName?: boolean;
    valueTypeName: string;
    amountValue?: boolean;
    percentValue?: boolean;
    defaultCurrency?: string;
}>();

const theme = useTheme();

const {
    tt,
    getCurrentLanguageTextDirection,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatNumberToLocalizedNumerals,
    formatPercentToLocalizedNumerals
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
    const allSeriesNames: string[] = [];
    let minValue: number = Number.POSITIVE_INFINITY;
    let maxValue: number = 0;

    for (const [item, seriesIndex] of itemAndIndex(props.items)) {
        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        if (!isArray(item[props.valuesField])) {
            continue;
        }

        allSeriesNames.push(getItemName(item[props.nameField] as string));

        const allAmounts: number[] = item[props.valuesField] as number[];

        for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
            if (amount > maxValue) {
                maxValue = amount;
            }

            if (amount < minValue) {
                minValue = amount;
            }

            allData.push([categoryIndex, seriesIndex, amount]);
        }
    }

    const ret: HeatMapData = {
        allSeriesNames: allSeriesNames,
        data: allData,
        minValue: minValue === Number.POSITIVE_INFINITY ? 0 : minValue,
        maxValue: maxValue
    };

    return ret;
});

const yAxisWidth = computed<number>(() => {
    let width: number = 90;

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
                const value = dataItem && isNumber(dataItem[2]) ? getDisplayValue(dataItem[2]) : '';

                return `<div class="d-inline-flex">${params.name}</div><br/>`
                    + `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`
                    + `<span>${name}</span>`
                    + `<span class="ms-5">${value}</span>`
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
                min: heatMapData.value.minValue,
                max: heatMapData.value.maxValue,
                calculable: true,
                inRange: {
                    color: isDarkMode.value ? [ '#1a1a1a', '#c67e48' ] : [ '#faf8f4', '#c67e48' ]
                },
                textStyle: {
                    color: isDarkMode.value ? '#888' : '#666'
                },
                formatter: (value: string) => {
                    if (!props.showValue) {
                        return '';
                    }

                    return getDisplayValue(parseInt(value));
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
                        const value: number = data && isNumber(data[2]) ? data[2] : 0;
                        return getDisplayValue(value);
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

function getItemName(name: string): string {
    return props.translateName ? tt(name) : name;
}

function getDisplayValue(value: number): string {
    if (props.percentValue) {
        return formatPercentToLocalizedNumerals(value, 2, '<0.01');
    }

    if (props.amountValue) {
        return formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
    }

    return formatNumberToLocalizedNumerals(value, 2);
}
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
