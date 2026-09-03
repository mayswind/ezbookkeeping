<template>
    <div class="overview-monthly-chart">
        <v-card-text class="overview-monthly-chart-container overview-monthly-chart-overlay" v-if="loading && !hasAnyData">
            <div class="overview-monthly-chart-skeleton-container h-100" style="margin-top: -30px">
                <div class="d-flex w-100 h-100 align-center justify-center"
                     :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ]">
                    <v-skeleton-loader width="16" height="200" :loading="true"></v-skeleton-loader>
                </div>
            </div>
        </v-card-text>

        <v-card-text class="overview-monthly-chart-container overview-monthly-chart-overlay" v-else-if="!loading && !hasAnyData">
            <div class="d-flex flex-column align-center justify-center w-100 h-100">
                <span class="text-title-medium mt-n13">{{ tt('No data') }}</span>
            </div>
        </v-card-text>

        <v-chart autoresize class="overview-monthly-chart-container" :class="{ 'readonly': !hasAnyData }"
                 :option="chartOptions" :update-options="{ notMerge: true }"
                 @click="clickItem"/>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

import { TextDirection } from '@/core/text.ts';
import type { BigDecimal, HiddenAmount } from '@/core/numeral.ts';
import { TrendChartType } from '@/core/statistics.ts';
import { TransactionType } from '@/core/transaction.ts';
import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import { type TransactionMonthlyIncomeAndExpenseData } from '@/models/transaction.ts';

import { BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { getExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export interface MonthlyIncomeAndExpenseCardClickEvent {
    transactionType: TransactionType;
    monthStartTime: number;
}

const props = defineProps<{
    loading: boolean;
    data: TransactionMonthlyIncomeAndExpenseData[];
    chartType: number;
    transactionTypes: number[];
    disabled: boolean;
    isDarkMode?: boolean;
    enableClickItem?: boolean;
    hideLegend?: boolean;
    hideXAxisLabels?: boolean;
    noMargin?: boolean;
    smoothCurve?: boolean;
}>();

const emit = defineEmits<{
    (e: 'click', event: MonthlyIncomeAndExpenseCardClickEvent): void;
}>();

const {
    tt,
    getCurrentLanguageTextDirection,
    formatDateTimeToGregorianLikeShortMonth,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const showIncome = computed<boolean>(() => props.transactionTypes.includes(TransactionType.Income));
const showExpense = computed<boolean>(() => props.transactionTypes.includes(TransactionType.Expense));
const showIncomeAndExpense = computed<boolean>(() => showIncome.value && showExpense.value);
const hasAnyData = computed<boolean>(() => {
    if (!props.data || !props.data.length || props.data.length < 1) {
        return false;
    }

    for (const item of props.data) {
        if ((showIncome.value && !item.incomeAmount.isZero()) || (showExpense.value && !item.expenseAmount.isZero())) {
            return true;
        }
    }

    return false;
});

const chartOptions = computed<object>(() => {
    const monthNames: string[] = [];
    const incomeAmounts: number[] = []; // only used for echarts rendering, the actual value is in props.data
    const expenseAmounts: number[] = []; // only used for echarts rendering, the actual value is in props.data
    let minAmount: BigDecimal = BIG_DECIMAL_ZERO;
    let maxAmount: BigDecimal = BIG_DECIMAL_ZERO;

    const expenseIncomeAmountColor = getExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor, props.isDarkMode);

    if (props.data) {
        for (const item of props.data) {
            const monthStartDateTime = parseDateTimeFromUnixTime(item.monthStartTime);
            const monthShortName = formatDateTimeToGregorianLikeShortMonth(monthStartDateTime);

            monthNames.push(monthShortName);

            if (showIncome.value) {
                incomeAmounts.push(item.incomeAmount.toDoubleNumber());

                if (item.incomeAmount.greaterThan(maxAmount)) {
                    maxAmount = item.incomeAmount;
                }

                if (item.incomeAmount.lessThan(minAmount)) {
                    minAmount = item.incomeAmount;
                }
            }

            if (showExpense.value) {
                const expenseAmount = showIncomeAndExpense.value ? item.expenseAmount.negate() : item.expenseAmount;
                expenseAmounts.push(expenseAmount.toDoubleNumber());

                if (expenseAmount.greaterThan(maxAmount)) {
                    maxAmount = expenseAmount;
                }

                if (expenseAmount.lessThan(minAmount)) {
                    minAmount = expenseAmount;
                }
            }
        }
    }

    const amountGap = maxAmount.subtract(minAmount);

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: props.chartType === TrendChartType.Area.type ? {
                type: 'cross',
                label: {
                    show: showAmountInHomePage.value,
                    backgroundColor: props.isDarkMode ? '#333' : '#fff',
                    color: props.isDarkMode ? '#eee' : '#333'
                }
            } : {
                type: 'shadow',
                shadowStyle: {
                    color: props.isDarkMode ? 'rgba(210, 210, 210, 0.05)' : 'rgba(120, 120, 120, 0.05)'
                }
            },
            backgroundColor: props.isDarkMode ? '#333' : '#fff',
            borderColor: props.isDarkMode ? '#333' : '#fff',
            textStyle: {
                color: props.isDarkMode ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams[]) => {
                let incomeAmount: string | null = null;
                let expenseAmount: string | null = null;

                for (const param of params) {
                    const dataIndex = param.dataIndex;
                    const data = props.data[dataIndex] as TransactionMonthlyIncomeAndExpenseData;

                    if (param.seriesId === 'seriesIncome') {
                        incomeAmount = getDisplayIncomeAmount(data);
                    } else if (param.seriesId === 'seriesExpense') {
                        expenseAmount = getDisplayExpenseAmount(data);
                    }
                }

                return `<table>` +
                    `<thead>` +
                    `<tr>` +
                    `<td colspan="2" class="text-start">${params[0]?.name}</td>` +
                    `</tr>` +
                    `</thead>` +
                    `<tbody>` +
                    (
                        incomeAmount !== null ?
                            `<tr>` +
                            `<td><span class="overview-monthly-chart-tooltip-indicator bg-income me-1"></span><span class="me-4">${tt('Income')}</span></td>` +
                            `<td><strong>${incomeAmount}</strong></td>` +
                            `</tr>` : ''
                    )+
                    (
                        expenseAmount !== null ?
                            `<tr>` +
                            `<td><span class="overview-monthly-chart-tooltip-indicator bg-expense me-1"></span><span class="me-4">${tt('Expense')}</span></td>` +
                            `<td><strong>${expenseAmount}</strong></td>` +
                            `</tr>` : ''
                    ) +
                    `</tbody>` +
                    `</table>`;
            }
        },
        legend: {
            show: !props.hideLegend,
            bottom: 5,
            itemWidth: 14,
            itemHeight: 14,
            textStyle: {
                color: props.isDarkMode ? '#eee' : '#333'
            },
            icon: 'circle',
            data: [
                ...(showIncome.value ? [tt('Income')] : []),
                ...(showExpense.value ? [tt('Expense')] : [])
            ]
        },
        grid: {
            left: props.noMargin ? 0 : 10,
            right: props.noMargin ? 0 : 10,
            top: props.noMargin ? 0 : (showIncomeAndExpense.value ? 10 : 5),
            bottom: props.chartType === TrendChartType.Area.type
                ? props.noMargin ? 0 : (!props.hideXAxisLabels ? 30 : 10) + (!props.hideLegend ? 25 : 0)
                : (!props.hideXAxisLabels ? 20 : 0) + (!props.hideLegend ? 25 : 0) + (showIncomeAndExpense.value ? 20 : 10),
        },
        xAxis: [
            {
                type: 'category',
                data: monthNames,
                boundaryGap: props.chartType !== TrendChartType.Area.type,
                inverse: textDirection.value === TextDirection.RTL,
                axisLine: {
                    show: false
                },
                axisTick: {
                    show: false
                },
                axisLabel: {
                    show: !props.hideXAxisLabels,
                    padding: [ (props.chartType === TrendChartType.Column.type && showIncomeAndExpense.value ? 10 : 0), 0, 0, 0 ]
                }
            }
        ],
        yAxis: props.chartType === TrendChartType.Area.type || !showIncomeAndExpense.value ? [
            {
                type: 'value',
                min: minAmount.toDoubleNumber(),
                max: maxAmount.toDoubleNumber(),
                splitNumber: 10,
                axisLabel: {
                    show: false
                },
                splitLine: {
                    show: false
                }
            }
        ] : [
            {
                type: 'value',
                min: minAmount.subtract(amountGap.divide(20)).toDoubleNumber(),
                max: maxAmount.toDoubleNumber(),
                splitNumber: 10,
                axisLabel: {
                    show: false
                },
                splitLine: {
                    show: false
                }
            },
            {
                type: 'value',
                min: minAmount.toDoubleNumber(),
                max: maxAmount.add(amountGap.divide(20)).toDoubleNumber(),
                splitNumber: 10,
                axisLabel: {
                    show: false
                },
                splitLine: {
                    show: false
                }
            }
        ],
        series: [
            ...(showIncome.value ? [{
                type: props.chartType === TrendChartType.Area.type ? 'line' : 'bar',
                id: 'seriesIncome',
                name: tt('Income'),
                yAxisIndex: 0,
                stack: props.chartType === TrendChartType.Column.type && showIncomeAndExpense.value ? 'Total' : undefined,
                areaStyle: props.chartType === TrendChartType.Area.type ? {} : undefined,
                smooth: props.smoothCurve,
                showSymbol: false,
                itemStyle: {
                    color: expenseIncomeAmountColor.incomeAmountColor,
                    borderRadius: props.chartType === TrendChartType.Area.type ? undefined : 16
                },
                emphasis: {
                    focus: 'series',
                    labelLine: {
                        show: false
                    }
                },
                barMaxWidth: props.chartType === TrendChartType.Area.type ? undefined : 16,
                data: incomeAmounts
            }] : []),
            ...(showExpense.value ? [{
                type: props.chartType === TrendChartType.Area.type ? 'line' : 'bar',
                id: 'seriesExpense',
                name: tt('Expense'),
                yAxisIndex: props.chartType === TrendChartType.Column.type && showIncomeAndExpense.value ? 1 : 0,
                stack: props.chartType === TrendChartType.Column.type && showIncomeAndExpense.value ? 'Total' : undefined,
                areaStyle: props.chartType === TrendChartType.Area.type ? {} : undefined,
                smooth: props.smoothCurve,
                showSymbol: false,
                itemStyle: {
                    color: expenseIncomeAmountColor.expenseAmountColor,
                    borderRadius: props.chartType === TrendChartType.Area.type ? undefined : 16
                },
                emphasis: {
                    focus: 'series',
                    labelLine: {
                        show: false
                    }
                },
                barMaxWidth: props.chartType === TrendChartType.Area.type ? undefined : 16,
                data: expenseAmounts
            }] : [])
        ]
    };
});

function getDisplayCurrency(value: BigDecimal | HiddenAmount, currencyCode: string): string {
    return formatAmountToLocalizedNumeralsWithCurrency(value, currencyCode);
}

function getDisplayAmount(amount: BigDecimal, incomplete: boolean): string {
    if (!showAmountInHomePage.value) {
        return getDisplayCurrency(DISPLAY_HIDDEN_AMOUNT, defaultCurrency.value);
    }

    return getDisplayCurrency(amount, defaultCurrency.value) + (incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
}

function getDisplayIncomeAmount(data: TransactionMonthlyIncomeAndExpenseData): string {
    return getDisplayAmount(data.incomeAmount, data.incompleteIncomeAmount);
}

function getDisplayExpenseAmount(data: TransactionMonthlyIncomeAndExpenseData): string {
    return getDisplayAmount(data.expenseAmount, data.incompleteExpenseAmount);
}

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || !props.data || e.componentType !== 'series') {
        return;
    }

    const clickData = props.data[e.dataIndex];

    if (clickData && e.seriesId === 'seriesIncome') {
        emit('click', {
            transactionType: TransactionType.Income,
            monthStartTime: clickData.monthStartTime
        });
    } else if (clickData && e.seriesId === 'seriesExpense') {
        emit('click', {
            transactionType: TransactionType.Expense,
            monthStartTime: clickData.monthStartTime
        });
    }
}
</script>

<style>
.overview-monthly-chart {
    flex: 1 1 0;
    min-height: 0;
    position: relative;
}

.overview-monthly-chart-overlay {
    position: absolute !important;
    z-index: 10;
    inset: 0;
}

.overview-monthly-chart-skeleton-container {
    display: grid;
    grid-template-columns: repeat(12, minmax(0, 1fr));
}

.overview-monthly-chart-tooltip-indicator {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 10px;
}
</style>
