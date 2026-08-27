<template>
    <v-card class="monthly-expense-widget h-100" :class="{ disabled: loading }">
        <template #title>
            <div class="d-flex align-end">
                <div class="d-flex align-baseline">
                    <span class="text-headline-small font-weight-bold">{{ displayDateRange?.thisMonth?.displayTime }}</span>
                    <span class="mx-1 text-title-large text-medium-emphasis">·</span>
                    <span class="text-medium-emphasis font-weight-bold text-title-small">{{ tt('Expense') }}</span>
                </div>
                <v-btn class="ms-2" density="compact" color="default" variant="text"
                       :icon="true" :loading="loading" @click="$emit('refresh')">
                    <template #loader>
                        <v-progress-circular indeterminate size="20" />
                    </template>
                    <v-icon :icon="mdiRefresh" size="24" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
            </div>
        </template>

        <v-card-text class="pt-2">
            <div class="monthly-expense-widget__content">
                <div class="mt-2 d-flex align-end">
                    <span class="text-headline-small font-weight-medium text-primary text-no-wrap">
                        <span v-if="!loading || (transactionOverview && transactionOverview.thisMonth && transactionOverview.thisMonth.valid)">{{ transactionOverview && transactionOverview.thisMonth ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-' }}</span>
                        <v-skeleton-loader class="d-inline-block skeleton-no-margin mt-3 pb-1" width="120px" type="text" :loading="true" v-else-if="loading && (!transactionOverview || !transactionOverview.thisMonth || !transactionOverview.thisMonth.valid)"></v-skeleton-loader>
                    </span>
                    <v-btn class="ms-1" density="compact" color="primary" variant="text"
                           :icon="true" @click="showAmountInHomePage = !showAmountInHomePage">
                        <v-icon :icon="showAmountInHomePage ? mdiEyeOffOutline : mdiEyeOutline" size="20" />
                    </v-btn>
                </div>

                <div class="monthly-expense-widget__income mt-2">
                    <span class="monthly-expense-widget__income-dot"></span>
                    <span class="text-medium-emphasis text-truncate">{{ tt('Monthly income') }}</span>
                    <span class="font-weight-bold text-body-medium text-truncate" v-if="!loading || (transactionOverview && transactionOverview.thisMonth && transactionOverview.thisMonth.valid)">{{ transactionOverview && transactionOverview.thisMonth ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-' }}</span>
                    <v-skeleton-loader class="d-inline-block skeleton-no-margin mt-1 pb-1" width="100px" type="text" :loading="true" v-else-if="loading && (!transactionOverview || !transactionOverview.thisMonth || !transactionOverview.thisMonth.valid)"></v-skeleton-loader>
                </div>

                <v-btn class="monthly-expense-widget__details mt-4" variant="tonal" :to="currentDetailsUrl">{{ tt('View Details') }}</v-btn>
            </div>

            <div class="monthly-expense-widget__illustration img-with-direction" aria-hidden="true">
                <svg viewBox="0 0 168 132" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <circle class="monthly-expense-widget__illustration-halo" cx="112" cy="73" r="52" />
                    <circle class="monthly-expense-widget__illustration-dot" cx="151" cy="25" r="4" />
                    <circle class="monthly-expense-widget__illustration-dot monthly-expense-widget__illustration-dot--muted" cx="18" cy="108" r="3" />

                    <g class="monthly-expense-widget__illustration-paper">
                        <rect x="28" y="16" width="78" height="101" rx="12" />
                        <path d="M85 16H94C100.627 16 106 21.3726 106 28V38L85 16Z" />
                        <rect class="monthly-expense-widget__illustration-line" x="43" y="39" width="34" height="6" rx="3" />
                        <rect class="monthly-expense-widget__illustration-line monthly-expense-widget__illustration-line--short" x="43" y="53" width="49" height="4" rx="2" />
                        <rect class="monthly-expense-widget__illustration-line monthly-expense-widget__illustration-line--short" x="43" y="64" width="39" height="4" rx="2" />
                        <rect class="monthly-expense-widget__illustration-total" x="42" y="84" width="49" height="18" rx="6" />
                        <path class="monthly-expense-widget__illustration-total-line" d="M50 93H68" />
                    </g>

                    <g class="monthly-expense-widget__illustration-chart">
                        <circle cx="117" cy="82" r="29" />
                        <circle class="monthly-expense-widget__illustration-chart-track" cx="117" cy="82" r="18" />
                        <circle class="monthly-expense-widget__illustration-chart-value" cx="117" cy="82" r="18" pathLength="100" />
                        <circle class="monthly-expense-widget__illustration-chart-center" cx="117" cy="82" r="8" />
                    </g>

                    <g class="monthly-expense-widget__illustration-coins">
                        <ellipse cx="139" cy="111" rx="15" ry="5" />
                        <path d="M124 104V111C124 113.761 130.716 116 139 116C147.284 116 154 113.761 154 111V104" />
                        <ellipse cx="139" cy="104" rx="15" ry="5" />
                    </g>
                </svg>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

import { DateRange } from '@/core/datetime.ts';

import {
    mdiRefresh,
    mdiEyeOutline,
    mdiEyeOffOutline
} from '@mdi/js';

defineProps<{
    loading: boolean;
}>();

defineEmits<{
    (e: 'refresh'): void
}>();

const { tt } = useI18n();

const {
    showAmountInHomePage,
    displayDateRange,
    transactionOverview,
    currentDetailsUrl,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount
} = usePeriodStatisticsWidgetBase({
    dateType: DateRange.ThisMonth.type
});
</script>
