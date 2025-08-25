<template>
    <v-card :class="{ 'disabled': disabled }">
        <template #title>
            <span>{{ tt('Income and Expense Trends') }}</span>
        </template>

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
                <h2 style="margin-top: -40px">{{ tt('No data') }}</h2>
            </div>
        </v-card-text>

        <v-chart autoresize class="overview-monthly-chart-container"
                 :class="{ 'readonly': !hasAnyData }" :option="chartOptions" @click="clickItem"/>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

import { TextDirection } from '@/core/text.ts';
import type { HiddenAmount } from '@/core/numeral.ts';
import { TransactionType } from '@/core/transaction.ts';
import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import { type TransactionMonthlyIncomeAndExpenseData } from '@/models/transaction.ts';

import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { getExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export interface MonthlyIncomeAndExpenseCardClickEvent {
    transactionType: TransactionType;
    monthStartTime: number;
}

const props = defineProps<{
    loading: boolean;
    data: TransactionMonthlyIncomeAndExpenseData[];
    disabled: boolean;
    isDarkMode?: boolean;
    enableClickItem?: boolean;
}>();

const emit = defineEmits<{
    (e: 'click', event: MonthlyIncomeAndExpenseCardClickEvent): void;
}>();

const { tt, getCurrentLanguageTextDirection, formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const hasAnyData = computed<boolean>(() => {
    if (!props.data || !props.data.length || props.data.length < 1) {
        return false;
    }

    for (let i = 0; i < props.data.length; i++) {
        const item = props.data[i];

        if (item.incomeAmount > 0 || item.incomeAmount < 0 || item.expenseAmount > 0 || item.expenseAmount < 0) {
            return true;
        }
    }

    return false;
});

const chartOptions = computed<object>(() => {
    const monthNames: string[] = [];
    const incomeAmounts: number[] = [];
    const expenseAmounts: number[] = [];
    let minAmount = 0;
    let maxAmount = 0;

    const expenseIncomeAmountColor = getExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor, props.isDarkMode);

    if (props.data) {
        for (let i = 0; i < props.data.length; i++) {
            const item = props.data[i];
            const monthShortName = parseDateTimeFromUnixTime(item.monthStartTime).getGregorianCalendarMonthDisplayShortName();

            monthNames.push(monthShortName);
            incomeAmounts.push(item.incomeAmount);
            expenseAmounts.push(-item.expenseAmount);

            if (item.incomeAmount > maxAmount) {
                maxAmount = item.incomeAmount;
            }

            if (-item.expenseAmount > maxAmount) {
                maxAmount = -item.expenseAmount;
            }

            if (item.incomeAmount < minAmount) {
                minAmount = item.incomeAmount;
            }

            if (-item.expenseAmount < minAmount) {
                minAmount = -item.expenseAmount;
            }
        }
    }

    const amountGap = maxAmount - minAmount;

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: {
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

                for (let i = 0; i < params.length; i++) {
                    const param = params[i];
                    const dataIndex = param.dataIndex;
                    const data = props.data[dataIndex];

                    if (param.seriesId === 'seriesIncome') {
                        incomeAmount = getDisplayIncomeAmount(data);
                    } else if (param.seriesId === 'seriesExpense') {
                        expenseAmount = getDisplayExpenseAmount(data);
                    }
                }

                return `<table>` +
                    `<thead>` +
                    `<tr>` +
                    `<td colspan="2" class="text-start">${params[0].name}</td>` +
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
            bottom: 20,
            itemWidth: 14,
            itemHeight: 14,
            textStyle: {
                color: props.isDarkMode ? '#eee' : '#333'
            },
            icon: 'circle',
            data: [ tt('Income'), tt('Expense') ]
        },
        grid: {
            left: '20px',
            right: '20px',
            top: '10px',
            bottom: '100px'
        },
        xAxis: [
            {
                type: 'category',
                data: monthNames,
                inverse: textDirection.value === TextDirection.RTL,
                axisLine: {
                    show: false
                },
                axisTick: {
                    show: false
                },
                axisLabel: {
                    padding: [ 20, 0, 0, 0 ]
                }
            }
        ],
        yAxis: [
            {
                type: 'value',
                min: minAmount - amountGap / 20,
                max: maxAmount,
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
                min: minAmount,
                max: maxAmount + amountGap / 20,
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
            {
                type: 'bar',
                id: 'seriesIncome',
                name: tt('Income'),
                yAxisIndex: 0,
                stack: 'Total',
                itemStyle: {
                    color: expenseIncomeAmountColor.incomeAmountColor,
                    borderRadius: 16
                },
                emphasis: {
                    focus: 'series',
                    labelLine: {
                        show: false
                    }
                },
                barMaxWidth: 16,
                data: incomeAmounts
            },
            {
                type: 'bar',
                id: 'seriesExpense',
                name: tt('Expense'),
                yAxisIndex: 1,
                stack: 'Total',
                itemStyle: {
                    color: expenseIncomeAmountColor.expenseAmountColor,
                    borderRadius: 16
                },
                emphasis: {
                    focus: 'series',
                    labelLine: {
                        show: false
                    }
                },
                barMaxWidth: 16,
                data: expenseAmounts
            }
        ]
    };
});

function getDisplayCurrency(value: number | HiddenAmount, currencyCode: string): string {
    return formatAmountToLocalizedNumeralsWithCurrency(value, currencyCode);
}

function getDisplayAmount(amount: number, incomplete: boolean): string {
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
.overview-monthly-chart-container {
    width: 100%;
    height: 400px;
}

.overview-monthly-chart-overlay {
    position: absolute !important;
    z-index: 10;
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
