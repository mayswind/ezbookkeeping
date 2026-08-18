<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <template #title>
            <div class="d-flex align-center">
                <div class="d-flex align-baseline">
                    <span class="text-headline-small font-weight-bold">{{ displayDateRange?.thisMonth?.displayTime }}</span>
                    <span class="text-title-large">·</span>
                    <span class="text-title-small">{{ tt('Expense') }}</span>
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

        <v-card-text class="mt-4">
            <span class="text-headline-small font-weight-medium text-primary">
                <span v-if="!loading || (transactionOverview && transactionOverview.thisMonth && transactionOverview.thisMonth.valid)">{{ transactionOverview && transactionOverview.thisMonth ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-' }}</span>
                <v-skeleton-loader class="d-inline-block skeleton-no-margin mt-3 pb-1" width="120px" type="text" :loading="true" v-else-if="loading && (!transactionOverview || !transactionOverview.thisMonth || !transactionOverview.thisMonth.valid)"></v-skeleton-loader>
                <v-btn class="ms-1" density="compact" color="default" variant="text"
                       :icon="true" @click="showAmountInHomePage = !showAmountInHomePage">
                    <v-icon :icon="showAmountInHomePage ? mdiEyeOffOutline : mdiEyeOutline" size="20" />
                </v-btn>
            </span>
            <div class="mt-3 mb-2" style="padding-bottom: 1px">
                <span class="me-2">{{ tt('Monthly income') }}</span>
                <span class="text-body-medium" v-if="!loading || (transactionOverview && transactionOverview.thisMonth && transactionOverview.thisMonth.valid)">{{ transactionOverview && transactionOverview.thisMonth ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-' }}</span>
                <v-skeleton-loader class="d-inline-block skeleton-no-margin mt-1 pb-1" width="120px" type="text" :loading="true" v-else-if="loading && (!transactionOverview || !transactionOverview.thisMonth || !transactionOverview.thisMonth.valid)"></v-skeleton-loader>
            </div>
            <v-btn class="mt-4" variant="tonal" :to="detailsUrl">{{ tt('View Details') }}</v-btn>
            <v-img class="overview-card-background img-with-direction" src="img/desktop/card-background.png" />
            <v-img class="overview-card-background-image img-with-direction" width="116" src="img/desktop/document.svg" />
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useHomePageBase } from '@/views/base/HomePageBase.ts';

import { useOverviewStore } from '@/stores/overview.ts';

import { DateRange } from '@/core/datetime.ts';

import {
    mdiRefresh,
    mdiEyeOutline,
    mdiEyeOffOutline
} from '@mdi/js';

defineProps<{
    loading: boolean
}>();

defineEmits<{
    (e: 'refresh'): void
}>();

const { tt } = useI18n();

const overviewStore = useOverviewStore();

const {
    showAmountInHomePage,
    displayDateRange,
    transactionOverview,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount
} = useHomePageBase();

const detailsUrl = computed<string>(() => `/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: DateRange.ThisMonth.type })}`);
</script>
