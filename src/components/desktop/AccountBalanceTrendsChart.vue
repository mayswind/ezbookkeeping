<template>
    <v-chart autoresize class="account-balance-trends-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"/>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonAccountBalanceTrendsChartProps, useAccountBalanceTrendsChartBase } from '@/components/base/AccountBalanceTrendsChartBase.ts'

import type { ColorValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { TrendChartType } from '@/core/statistics.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

interface DesktopAccountBalanceTrendsChartProps extends CommonAccountBalanceTrendsChartProps {
    legendName: string;
    skeleton?: boolean;
    type?: number;
}

interface AccountBalanceTrendsChartDataItem {
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

const props = defineProps<DesktopAccountBalanceTrendsChartProps>();

const theme = useTheme();
const { formatAmountWithCurrency } = useI18n();
const { allDataItems, allDisplayDateRanges } = useAccountBalanceTrendsChartBase(props);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const allSeries = computed<AccountBalanceTrendsChartDataItem[]>(() => {
    const series: AccountBalanceTrendsChartDataItem = {
        id: 'accountBalance',
        name: props.legendName,
        itemStyle: {
            color: `#${DEFAULT_CHART_COLORS[0]}`
        },
        selected: true,
        type: 'line',
        stack: 'a',
        animation: !props.skeleton,
        data: []
    };

    if (props.type === TrendChartType.Area.type) {
        series.areaStyle = {};
    } else if (props.type === TrendChartType.Column.type) {
        series.type = 'bar';
    }

    for (let i = 0; i < allDataItems.value.length; i++) {
        const item = allDataItems.value[i];
        series.data.push(item.amount);
    }

    return [series];
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

    const maxValueText = formatAmountWithCurrency(maxValue, props.account.currency);
    const minValueText = formatAmountWithCurrency(minValue, props.account.currency);
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
                const amount = params[0].data as number;
                const value = formatAmountWithCurrency(amount, props.account.currency);

                return `${params[0].name}<br/>`
                    + '<div><span class="chart-pointer" style="background-color: #' + DEFAULT_CHART_COLORS[0] + '"></span>'
                    + `<span>${props.legendName}</span><span style="margin-left: 20px; float: right">${value}</span><br/>`
                    + '</div>';
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
                        return formatAmountWithCurrency(value, props.account.currency);
                    }
                },
                axisPointer: {
                    label: {
                        formatter: (params: CallbackDataParams) => {
                            return formatAmountWithCurrency(Math.floor(params.value as number), props.account.currency);
                        }
                    }
                }
            }
        ],
        series: allSeries.value
    };
});
</script>

<style scoped>
.account-balance-trends-chart-container {
    width: 100%;
    height: 400px;
    margin-top: 10px;
}
</style>
