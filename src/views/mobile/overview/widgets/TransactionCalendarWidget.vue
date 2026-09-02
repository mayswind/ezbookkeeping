<template>
    <f7-card class="no-margin-top margin-bottom" :class="{ disabled: loading }">
        <f7-card-content class="transaction-calendar-container no-padding">
            <transaction-calendar calendar-class="justify-content-center" week-day-name-type="short"
                                  :readonly="loading || editing" :is-dark-mode="isDarkMode"
                                  :default-currency="false"
                                  :min-date="transactionCalendarMinDate" :max-date="transactionCalendarMaxDate"
                                  :daily-total-amounts="dailyTotalAmounts"
                                  :show-amount="showAmountInCalendar"
                                  :show-income-amount="showIncome"
                                  :show-expense-amount="showExpense"
                                  :show-alternate-date="showAlternateDate"
                                  :model-value="currentCalendarDate"
                                  @update:model-value="selectDate" />
        </f7-card-content>
    </f7-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useTransactionCalendarWidgetBase } from '@/views/base/overview/TransactionCalendarWidgetBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

import type { TextualYearMonthDay } from '@/core/datetime.ts';

const props = defineProps<{
    loading: boolean;
    editing?: boolean;
    title?: string;
    transactionTypes: number[];
    showAlternateDate: boolean;
    showAmount: boolean;
}>();

const environmentsStore = useEnvironmentsStore();

const {
    currentCalendarDate,
    showAmountInCalendar,
    showIncome,
    showExpense,
    dailyTotalAmounts,
    transactionCalendarMinDate,
    transactionCalendarMaxDate,
    getTransactionListUrl
} = useTransactionCalendarWidgetBase(props);

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

function selectDate(date: TextualYearMonthDay): void {
    currentCalendarDate.value = date;
    const url = getTransactionListUrl(date);

    if (url) {

    }
}
</script>
