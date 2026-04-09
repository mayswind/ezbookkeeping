<template>
    <v-chart autoresize class="account-balance-trends-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"/>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';
import {
    type AccountBalanceTrendsChartItem,
    type CommonAccountBalanceTrendsChartProps,
    useAccountBalanceTrendsChartBase
} from '@/components/base/AccountBalanceTrendsChartBase.ts'

import { useUserStore } from '@/stores/user.ts';

import { type NameNumeralValue, itemAndIndex } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import type { ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { AccountBalanceTrendChartType } from '@/core/statistics.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

import { isArray } from '@/lib/common.ts';
import { getExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

interface DesktopAccountBalanceTrendsChartProps extends CommonAccountBalanceTrendsChartProps {
    legendName: string;
    skeleton?: boolean;
    type?: number;
}

interface AccountBalanceTrendsChartDataItem {
    id: string;
    name: string;
    itemStyle: {
        color: ColorStyleValue;
        color0?: ColorStyleValue;
        borderColor?: ColorStyleValue;
        borderColor0?: ColorStyleValue;
    };
    selected: boolean;
    type: string;
    areaStyle?: object;
    stack: string;
    animation: boolean;
    data: (number | number[])[];
}

const props = defineProps<DesktopAccountBalanceTrendsChartProps>();

const theme = useTheme();
const {
    tt,
    getCurrentLanguageTextDirection,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatPercentToLocalizedNumerals
} = useI18n();
const {
    showYearOverYearOnTooltip,
    showPeriodOverPeriodOnTooltip,
    allDataItems,
    allDataItemsMap,
    allDisplayDateRanges
} = useAccountBalanceTrendsChartBase(props);

const userStore = useUserStore();

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
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

    if (props.type === AccountBalanceTrendChartType.Area.type) {
        series.areaStyle = {};
    } else if (props.type === AccountBalanceTrendChartType.Column.type) {
        series.type = 'bar';
    } else if (props.type === AccountBalanceTrendChartType.Boxplot.type) {
        series.type = 'boxplot';
        series.itemStyle.borderColor = series.itemStyle.color;
    } else if (props.type === AccountBalanceTrendChartType.Candlestick.type) {
        const expenseIncomeAmountColor = getExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor, isDarkMode.value);
        series.type = 'candlestick';
        series.itemStyle.color = expenseIncomeAmountColor.incomeAmountColor;
        series.itemStyle.color0 = expenseIncomeAmountColor.expenseAmountColor;
        series.itemStyle.borderColor = expenseIncomeAmountColor.incomeAmountColor;
        series.itemStyle.borderColor0 = expenseIncomeAmountColor.expenseAmountColor;
    }

    for (const item of allDataItems.value) {
        if (props.type === AccountBalanceTrendChartType.Boxplot.type) {
            series.data.push([
                item.minimumBalance,
                item.q1Balance,
                item.medianBalance,
                item.q3Balance,
                item.maximumBalance
            ]);
        } else if (props.type === AccountBalanceTrendChartType.Candlestick.type) {
            series.data.push([
                item.openingBalance,
                item.closingBalance,
                item.minimumBalance,
                item.maximumBalance
            ]);
        } else {
            series.data.push(item.closingBalance);
        }
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

    for (const series of allSeries.value) {
        for (const data of series.data) {
            let currentMinValue: number;
            let currentMaxValue: number;

            if (isArray(data) && props.type === AccountBalanceTrendChartType.Boxplot.type) {
                currentMinValue = data[0] as number;
                currentMaxValue = data[4] as number;
            } else if (isArray(data) && props.type === AccountBalanceTrendChartType.Candlestick.type) {
                currentMinValue = data[2] as number;
                currentMaxValue = data[3] as number;
            } else {
                currentMinValue = data as number;
                currentMaxValue = data as number;
            }

            if (currentMaxValue > maxValue) {
                maxValue = currentMaxValue;
            }

            if (currentMinValue < minValue) {
                minValue = currentMinValue;
            }
        }
    }

    const maxValueText = formatAmountToLocalizedNumeralsWithCurrency(maxValue, props.account.currency);
    const minValueText = formatAmountToLocalizedNumeralsWithCurrency(minValue, props.account.currency);
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
                const dataIndex = params[0]!.dataIndex;
                const dataItem: AccountBalanceTrendsChartItem = allDataItems.value[dataIndex] as AccountBalanceTrendsChartItem;
                const yearOverYearDataItem: AccountBalanceTrendsChartItem | undefined = showYearOverYearOnTooltip.value ? allDataItemsMap.value[dataItem.lastYearDateRangeKey] : undefined;
                const periodOverPeriodDataItem: AccountBalanceTrendsChartItem | undefined = showPeriodOverPeriodOnTooltip.value ? allDataItems.value[dataIndex - 1] : undefined;

                let header: string = params[0]!.name;
                let displayItems: NameNumeralValue[] = [];
                let yearOverYearDataItemDisplayItems: NameNumeralValue[] | undefined = undefined;
                let periodOverPeriodDataItemDisplayItems: NameNumeralValue[] | undefined = undefined;
                let separatorLineIndex: number | undefined = undefined;

                if (dataItem.alternativeDisplayDate) {
                    header = dataItem.alternativeDisplayDate;
                }

                if (props.type === AccountBalanceTrendChartType.Boxplot.type) {
                    header += ` ${props.legendName}`;
                    displayItems = getBoxplotChartTooltip(dataItem);
                    yearOverYearDataItemDisplayItems = yearOverYearDataItem ? getBoxplotChartTooltip(yearOverYearDataItem) : undefined;
                    periodOverPeriodDataItemDisplayItems = periodOverPeriodDataItem ? getBoxplotChartTooltip(periodOverPeriodDataItem) : undefined;
                    separatorLineIndex = 5;
                } else if (props.type === AccountBalanceTrendChartType.Candlestick.type) {
                    header += ` ${props.legendName}`;
                    displayItems = getCandlestickChartTooltip(dataItem);
                    yearOverYearDataItemDisplayItems = yearOverYearDataItem ? getCandlestickChartTooltip(yearOverYearDataItem) : undefined;
                    periodOverPeriodDataItemDisplayItems = periodOverPeriodDataItem ? getCandlestickChartTooltip(periodOverPeriodDataItem) : undefined;
                    separatorLineIndex = 4;
                } else {
                    displayItems = getDefaultChartTooltip(dataItem);
                    yearOverYearDataItemDisplayItems = yearOverYearDataItem ? getDefaultChartTooltip(yearOverYearDataItem) : undefined;
                    periodOverPeriodDataItemDisplayItems = periodOverPeriodDataItem ? getDefaultChartTooltip(periodOverPeriodDataItem) : undefined;
                }

                const totalColumnCount = 2 + (yearOverYearDataItemDisplayItems && yearOverYearDataItemDisplayItems.length ? 1 : 0) + (periodOverPeriodDataItemDisplayItems && periodOverPeriodDataItemDisplayItems.length ? 1 : 0);
                let tooltip = `<table class="chart-tooltip-table"><tbody><tr><td colspan="2">${header}</td>`;

                if (yearOverYearDataItemDisplayItems && yearOverYearDataItemDisplayItems.length) {
                    tooltip += `<td><span class="ms-5" style="float: inline-end">${tt('Year-over-Year')}</span></td>`;
                }

                if (periodOverPeriodDataItemDisplayItems && periodOverPeriodDataItemDisplayItems.length) {
                    tooltip += `<td><span class="ms-5" style="float: inline-end">${tt('Period-over-Period')}</span></td>`;
                }

                tooltip += '</tr>';

                for (const [displayItem, index] of itemAndIndex(displayItems)) {
                    const displayValue = formatAmountToLocalizedNumeralsWithCurrency(displayItem.value, props.account.currency);
                    tooltip += `<tr><td><span class="chart-pointer" style="background-color: #${DEFAULT_CHART_COLORS[index]}"></span>`
                        + `<span>${displayItem.name}</span></td><td><span class="ms-5" style="float: inline-end">${displayValue}</span></td>`;

                    if (yearOverYearDataItemDisplayItems && yearOverYearDataItemDisplayItems.length && yearOverYearDataItemDisplayItems[index]) {
                        const yearOverYearDisplayItem = yearOverYearDataItemDisplayItems[index];
                        const displayGrowthRate = formatDisplayChangeRate(displayItem.value, yearOverYearDisplayItem.value);
                        tooltip += `<td><span class="ms-5" style="float: inline-end">${displayGrowthRate}</span></td>`;
                    }

                    if (periodOverPeriodDataItemDisplayItems && periodOverPeriodDataItemDisplayItems.length && periodOverPeriodDataItemDisplayItems[index]) {
                        const periodOverPeriodDisplayItem = periodOverPeriodDataItemDisplayItems[index];
                        const displayGrowthRate = formatDisplayChangeRate(displayItem.value, periodOverPeriodDisplayItem.value);
                        tooltip += `<td><span class="ms-5" style="float: inline-end">${displayGrowthRate}</span></td>`;
                    }

                    tooltip += '</tr>';

                    if (separatorLineIndex !== undefined && index === separatorLineIndex - 1) {
                        tooltip += `<tr><td colspan="${totalColumnCount}" style="border-bottom: ${(isDarkMode.value ? '#eee' : '#333')} dashed 1px"></td></tr>`;
                    }
                }

                tooltip += `</tbody></table>`;
                return tooltip;
            }
        },
        grid: {
            left: yAxisWidth.value,
            right: 20,
            top: 0,
            bottom: 20
        },
        xAxis: [
            {
                type: 'category',
                data: allDisplayDateRanges.value,
                inverse: textDirection.value === TextDirection.RTL,
                axisLabel: {
                    color: isDarkMode.value ? '#888' : '#666'
                }
            }
        ],
        yAxis: [
            {
                type: 'value',
                axisLabel: {
                    color: isDarkMode.value ? '#888' : '#666',
                    formatter: (value: string) => {
                        return formatAmountToLocalizedNumeralsWithCurrency(parseInt(value), props.account.currency);
                    }
                },
                axisPointer: {
                    label: {
                        formatter: (params: CallbackDataParams) => {
                            return formatAmountToLocalizedNumeralsWithCurrency(Math.trunc(params.value as number), props.account.currency);
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

function getBoxplotChartTooltip(dataItem: AccountBalanceTrendsChartItem): NameNumeralValue[] {
    return [
        {
            name: tt('Minimum Balance'),
            value: dataItem.minimumBalance
        },
        {
            name: tt('Q1 Balance (First Quartile)'),
            value: dataItem.q1Balance
        },
        {
            name: tt('Median Balance'),
            value: dataItem.medianBalance
        },
        {
            name: tt('Q3 Balance (Third Quartile)'),
            value: dataItem.q3Balance
        },
        {
            name: tt('Maximum Balance'),
            value: dataItem.maximumBalance
        },
        {
            name: tt('Opening Balance'),
            value: dataItem.openingBalance
        },
        {
            name: tt('Closing Balance'),
            value: dataItem.closingBalance
        }
    ];
}

function getCandlestickChartTooltip(dataItem: AccountBalanceTrendsChartItem): NameNumeralValue[] {
    return [
        {
            name: tt('Opening Balance'),
            value: dataItem.openingBalance
        },
        {
            name: tt('Closing Balance'),
            value: dataItem.closingBalance
        },
        {
            name: tt('Minimum Balance'),
            value: dataItem.minimumBalance
        },
        {
            name: tt('Maximum Balance'),
            value: dataItem.maximumBalance
        },
        {
            name: tt('Median Balance'),
            value: dataItem.medianBalance
        },
        {
            name: tt('Average Balance'),
            value: dataItem.averageBalance
        }
    ];
}

function getDefaultChartTooltip(dataItem: AccountBalanceTrendsChartItem): NameNumeralValue[] {
    return [
        {
            name: props.legendName,
            value: dataItem.closingBalance
        }
    ];
}


function formatDisplayChangeRate(current: number, reference: number): string {
    if (reference === 0 && current === 0) {
        return formatPercentToLocalizedNumerals(0, 2, '<0.01');
    }

    if (reference === 0) {
        return '-';
    }

    const rate = (current - reference) / reference * 100;
    return formatPercentToLocalizedNumerals(rate, 2, '<0.01');
}
</script>

<style scoped>
.account-balance-trends-chart-container {
    width: 100%;
    height: 418px;
    margin-top: 10px;
}
</style>
