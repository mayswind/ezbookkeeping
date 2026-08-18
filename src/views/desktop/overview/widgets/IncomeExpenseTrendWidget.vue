<template>
    <monthly-income-and-expense-chart class="h-100" :data="monthlyIncomeAndExpenseData" :is-dark-mode="isDarkMode"
                                      :loading="loading" :disabled="loading" :enable-click-item="true"
                                      @click="clickMonthlyIncomeOrExpense" />
</template>

<script setup lang="ts">
import { type MonthlyIncomeAndExpenseCardClickEvent } from '@/components/desktop/MonthlyIncomeAndExpenseChart.vue';

import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useHomePageBase } from '@/views/base/HomePageBase.ts';

import { useOverviewStore } from '@/stores/overview.ts';

import { DateRange } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';
import {
    type TransactionMonthlyIncomeAndExpenseData,
    LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES
} from '@/models/transaction.ts';

import { BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { getUnixTimeAfterUnixTime, getUnixTimeBeforeUnixTime } from '@/lib/datetime.ts';

const props = defineProps<{
    loading: boolean;
    months: number
}>();

const router = useRouter();
const theme = useTheme();

const {
    transactionOverview
} = useHomePageBase();

const overviewStore = useOverviewStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

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
