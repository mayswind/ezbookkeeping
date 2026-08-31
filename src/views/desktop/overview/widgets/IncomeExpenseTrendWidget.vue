<template>
    <v-card class="overview-widget d-flex flex-column" :class="{ 'disabled': loading }">
        <template #title>
            <overview-widget-header :title="displayTitle" :icon="mdiChartBar" />
        </template>

        <monthly-income-and-expense-chart :data="monthlyIncomeAndExpenseData" :is-dark-mode="isDarkMode"
                                          :loading="loading" :disabled="loading" :enable-click-item="true"
                                          :chart-type="chartType"
                                          :transaction-types="transactionTypes"
                                          :smooth-curve="smoothCurve"
                                          :hide-x-axis-labels="!showXAxisLabels" :hide-legend="!showLegend"
                                          @click="clickMonthlyIncomeOrExpense" />
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';
import { type MonthlyIncomeAndExpenseCardClickEvent } from '@/components/desktop/MonthlyIncomeAndExpenseChart.vue';

import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useOverviewStore } from '@/stores/overview.ts';

import { DateRange } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    type TransactionOverviewData,
    type TransactionMonthlyIncomeAndExpenseData,
    LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES
} from '@/models/transaction.ts';

import { BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { getUnixTimeAfterUnixTime, getUnixTimeBeforeUnixTime } from '@/lib/datetime.ts';

import {
    mdiChartBar
} from '@mdi/js';

const props = defineProps<{
    loading: boolean;
    title?: string;
    chartType: number;
    transactionTypes: number[];
    months: number;
    smoothCurve: boolean;
    showXAxisLabels: boolean;
    showLegend: boolean;
}>();

const router = useRouter();
const theme = useTheme();

const { tt } = useI18n();

const overviewStore = useOverviewStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const transactionOverview = computed<TransactionOverviewData>(() => overviewStore.transactionOverview);

const displayTitle = computed<string>(() => {
    if (props.title) {
        return props.title;
    }

    const showIncome = props.transactionTypes.includes(TransactionType.Income);
    const showExpense = props.transactionTypes.includes(TransactionType.Expense);

    if (showIncome && !showExpense) {
        return tt('Income Trends');
    } else if (!showIncome && showExpense) {
        return tt('Expense Trends');
    } else {
        return tt('Income and Expense Trends');
    }
});

const monthlyIncomeAndExpenseData = computed<TransactionMonthlyIncomeAndExpenseData[]>(() => {
    const data: TransactionMonthlyIncomeAndExpenseData[] = [];

    if (!transactionOverview.value || !transactionOverview.value.thisMonth || !transactionOverview.value.thisMonth.valid) {
        return data;
    }

    for (const amountRequestType of LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES) {
        const dateRange = overviewStore.transactionDataRange[amountRequestType];

        if (!dateRange) {
            continue;
        }

        const item = transactionOverview.value[amountRequestType];

        data.push({
            monthStartTime: dateRange.startTime,
            incomeAmount: item?.incomeAmount || BIG_DECIMAL_ZERO,
            expenseAmount: item?.expenseAmount || BIG_DECIMAL_ZERO,
            incompleteIncomeAmount: item ? item.incompleteIncomeAmount : true,
            incompleteExpenseAmount: item ? item.incompleteExpenseAmount : true
        });
    }

    if (props.months <= 0) {
        return data;
    }

    return data.slice(-props.months);
});

function clickMonthlyIncomeOrExpense(e: MonthlyIncomeAndExpenseCardClickEvent): void {
    const minTime = e.monthStartTime;
    const maxTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(minTime, 1, 'months'), 1, 'seconds');
    const type = e.transactionType;

    router.push(`/transaction/list?${overviewStore.getTransactionListPageParams({
        type: type,
        dateType: DateRange.Custom.type,
        minTime: minTime,
        maxTime: maxTime
    })}`);
}
</script>
