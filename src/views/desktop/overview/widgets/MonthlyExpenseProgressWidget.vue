<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ title || tt('Monthly Expense Progress') }}</span>
        </template>

        <v-card-text class="py-0">
            <div class="d-flex align-baseline mt-1">
                <span class="text-headline-small" :class="{ 'text-expense': !!currentDisplayExpenseAmount }" v-if="!loading || currentDisplayExpenseAmount">{{ currentDisplayExpenseAmount !== '' ? currentDisplayExpenseAmount : tt('No data') }}</span>
                <v-skeleton-loader class="skeleton-no-margin mt-3 mb-2" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
                <v-spacer />
                <span class="text-body-medium">{{ displayElapsedPercent }}</span>
            </div>
            <v-progress-linear class="mt-3" color="expense" rounded height="10" :model-value="currentMonthElapsedPercent * 100" />
            <div class="text-body-small text-medium-emphasis mt-2">{{ tt('Month elapsed') }}</div>
            <v-divider class="my-3" />
            <div class="d-flex align-center">
                <span>{{ tt('Estimated month-end expense') }}</span>
                <v-spacer />
                <span class="font-weight-medium" v-if="!loading || currentDisplayExpenseAmount">{{ displayEstimatedExpense }}</span>
                <v-skeleton-loader class="skeleton-no-margin" type="text" width="100px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
            </div>
            <div class="d-flex align-center mt-2 text-medium-emphasis">
                <span>{{ tt('Last month total') }}</span>
                <v-spacer />
                <span v-if="!loading || currentDisplayExpenseAmount">{{ displayLastMonthExpense }}</span>
                <v-skeleton-loader class="skeleton-no-margin" type="text" width="100px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
            </div>
        </v-card-text>
    </v-card>
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
