<template>
    <v-card class="overview-widget expense-progress-widget h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Monthly Expense Progress')" :icon="mdiCalendarClockOutline" />
        </template>

        <v-card-text class="overview-widget__body">
            <div class="expense-progress-widget__amount text-truncate">
                <span class="overview-widget__amount text-headline-small" :class="{ 'text-expense': !!currentDisplayExpenseAmount }" v-if="!loading || currentDisplayExpenseAmount">{{ currentDisplayExpenseAmount !== '' ? currentDisplayExpenseAmount : tt('No data') }}</span>
                <v-skeleton-loader class="skeleton-no-margin mt-2 mb-4" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
            </div>
            <div class="overview-widget__caption d-flex justify-space-between mt-2">
                <span>{{ tt('Month elapsed') }}</span>
                <span class="font-weight-medium">{{ displayElapsedPercent }}</span>
            </div>
            <v-progress-linear class="mt-2" color="primary" rounded height="6" :model-value="currentMonthElapsedPercent * 100" :aria-label="tt('Month elapsed')" />
            <div class="expense-progress-widget__projection mt-3 pt-3">
                <div class="overview-widget__detail-row">
                    <span class="text-truncate">{{ tt('Estimated month-end expense') }}</span>
                    <span class="overview-widget__amount font-weight-medium text-truncate" v-if="!loading || currentDisplayExpenseAmount">{{ displayEstimatedExpense }}</span>
                    <v-skeleton-loader class="skeleton-no-margin" type="text" width="100px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
                </div>
                <div class="overview-widget__detail-row overview-widget__caption mt-2">
                    <span class="text-truncate">{{ tt('Last month total') }}</span>
                    <span class="overview-widget__amount text-truncate" v-if="!loading || currentDisplayExpenseAmount">{{ displayLastMonthExpense }}</span>
                    <v-skeleton-loader class="skeleton-no-margin" type="text" width="100px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
                </div>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { type DateTime, DateRange } from '@/core/datetime.ts';

import { BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { getCurrentDateTime } from '@/lib/datetime.ts';

import {
    mdiCalendarClockOutline
} from '@mdi/js';

defineProps<{
    loading?: boolean;
    editing?: boolean;
    title?: string;
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
