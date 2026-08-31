<template>
    <f7-list strong inset dividers class="overview-transaction-list no-margin-top margin-bottom">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || tt('Monthly Expense Progress') }}</small>
        </f7-list-item>
        <f7-list-item class="item-title-full-line">
            <template #title>
                <div class="overview-widget-primary-amount skeleton-text" v-if="loading">0.00 USD</div>
                <div class="overview-widget-primary-amount" :class="{ 'text-expense': !!currentDisplayExpenseAmount }" v-else-if="!loading">
                    {{ currentDisplayExpenseAmount || tt('No data') }}
                </div>
                <div class="display-flex justify-content-space-between">
                    <small>{{ tt('Month elapsed') }}</small>
                    <small class="margin-left-half text-truncate skeleton-text" v-if="loading">100.0</small>
                    <small v-else-if="!loading">{{ displayElapsedPercent }}</small>
                </div>
                <f7-progressbar class="margin-top-half" :progress="currentMonthElapsedPercent * 100"
                                :aria-label="tt('Month elapsed')"></f7-progressbar>
            </template>
        </f7-list-item>
        <f7-list-item class="item-title-full-line">
            <template #footer>
                <div class="overview-transaction-footer default-text-color display-flex justify-content-space-between">
                    <span>{{ tt('Estimated month-end expense') }}</span>
                    <div class="margin-left-half text-truncate skeleton-text" v-if="loading">0.00 USD</div>
                    <span class="margin-left-half text-truncate" v-else-if="!loading">{{ displayEstimatedExpense }}</span>
                </div>
                <div class="overview-transaction-footer display-flex justify-content-space-between">
                    <span>{{ tt('Last month total') }}</span>
                    <div class="margin-left-half text-truncate skeleton-text" v-if="loading">0.00 USD</div>
                    <span class="margin-left-half text-truncate" v-else-if="!loading">{{ displayLastMonthExpense }}</span>
                </div>
            </template>
        </f7-list-item>
    </f7-list>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { type DateTime, DateRange } from '@/core/datetime.ts';

import { BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { getCurrentDateTime } from '@/lib/datetime.ts';

defineProps<{
    loading: boolean;
    title?: string;
    showTitle: boolean;
}>();

const { tt, formatPercentToLocalizedNumerals } = useI18n();
const {
    transactionOverview,
    currentExpenseAmount,
    currentDisplayExpenseAmount,
    getDisplayAmount
} = usePeriodStatisticsWidgetBase({
    dateType: DateRange.ThisMonth.type
});

const currentDateTime = ref<DateTime>(getCurrentDateTime());

const currentMonthElapsedPercent = computed<number>(() => currentDateTime.value.getGregorianCalendarDay() / currentDateTime.value.getMaxDayOfGregorianCalendarMonth());
const estimatedExpense = computed<BigDecimal>(() => currentExpenseAmount.value.divide(currentMonthElapsedPercent.value));
const lastMonthExpense = computed<BigDecimal>(() => transactionOverview.value.lastMonth?.expenseAmount ?? BIG_DECIMAL_ZERO);

const displayElapsedPercent = formatPercentToLocalizedNumerals(currentMonthElapsedPercent.value * 100.0, 0, '<1');
const displayEstimatedExpense = computed<string>(() => getDisplayAmount(estimatedExpense.value, !!transactionOverview.value.thisMonth?.incompleteExpenseAmount));
const displayLastMonthExpense = computed<string>(() => getDisplayAmount(lastMonthExpense.value, !!transactionOverview.value.lastMonth?.incompleteExpenseAmount));
</script>
