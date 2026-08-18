<template>
    <v-card class="h-100" :class="{ 'disabled': loading }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="grey" size="32">
                <v-icon size="22" :icon="currentPeriodConfig.icon" />
            </v-avatar>
            <span class="text-title-small font-weight-bold ms-2">{{ tt(currentPeriodConfig.title) }}</span>
            <v-spacer/>
            <v-btn density="comfortable" color="default" variant="text" class="ms-2" :icon="true" v-if="!editing">
                <v-icon :icon="mdiDotsVertical" />
                <v-menu activator="parent">
                    <v-list>
                        <v-list-item :prepend-icon="mdiListBoxOutline" :to="detailsUrl">
                            <v-list-item-title>{{ tt('View Details') }}</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-btn>
            <div class="v-btn--size-default" style="height: var(--v-btn-height)" v-else-if="editing"></div><!-- add a button placeholder to avoid height change -->
        </v-card-text>
        <v-card-text class="py-3">
            <div class="text-truncate text-headline-small text-income me-2 mb-2" v-if="!loading || displayIncomeAmount">{{ displayIncomeAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-3 mb-6" type="text" width="120px" :loading="true" v-else-if="loading && !displayIncomeAmount"></v-skeleton-loader>
            <div class="text-truncate text-title-medium text-expense pt-1" v-if="!loading || displayExpenseAmount">{{ displayExpenseAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin pt-1" style="padding-bottom: 7px" type="text" width="120px" :loading="true" v-else-if="loading && !displayExpenseAmount"></v-skeleton-loader>
            <div class="text-truncate text-title-medium mt-3 mb-5" v-if="!loading && !displayIncomeAmount && !displayExpenseAmount">{{ tt('No data') }}</div>
        </v-card-text>
        <v-card-text>
            <span class="text-body-medium">{{ displayDatetime }}</span>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useHomePageBase } from '@/views/base/HomePageBase.ts';

import { useOverviewStore } from '@/stores/overview.ts';

import { DateRange } from '@/core/datetime.ts';
import type { OverviewPeriod } from '@/core/overview_layout.ts';
import type { TransactionOverviewDataItem } from '@/models/transaction.ts';

import {
    mdiDotsVertical,
    mdiCalendarTodayOutline,
    mdiCalendarWeekOutline,
    mdiCalendarMonthOutline,
    mdiLayersTripleOutline,
    mdiListBoxOutline
} from '@mdi/js';

interface PeriodConfig {
    title: string;
    icon: string;
    dateType: DateRange;
}

const props = defineProps<{
    loading: boolean;
    period: OverviewPeriod;
    editing?: boolean
}>();

const { tt } = useI18n();
const {
    displayDateRange,
    transactionOverview,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount
} = useHomePageBase();

const overviewStore = useOverviewStore();

const allPeriodConfigMap: Record<OverviewPeriod, PeriodConfig> = {
    today: {
        title: 'Today',
        icon: mdiCalendarTodayOutline,
        dateType: DateRange.Today
    },
    thisWeek: {
        title: 'This Week',
        icon: mdiCalendarWeekOutline,
        dateType: DateRange.ThisWeek
    },
    thisMonth: {
        title: 'This Month',
        icon: mdiCalendarMonthOutline,
        dateType: DateRange.ThisMonth
    },
    thisYear: {
        title: 'This Year',
        icon: mdiLayersTripleOutline,
        dateType: DateRange.ThisYear
    }
};

const currentPeriodConfig = computed<PeriodConfig>(() => allPeriodConfigMap[props.period]);
const overviewItem = computed<TransactionOverviewDataItem | undefined>(() => transactionOverview.value[props.period]);
const displayIncomeAmount = computed<string>(() => overviewItem.value?.valid ? getDisplayIncomeAmount(overviewItem.value) : '');
const displayExpenseAmount = computed<string>(() => overviewItem.value?.valid ? getDisplayExpenseAmount(overviewItem.value) : '');

const displayDatetime = computed(() => {
    const range = displayDateRange.value[props.period];

    if (!range) {
        return '';
    }

    if ('displayTime' in range) {
        return range.displayTime || '';
    } else if ('startTime' in range && 'endTime' in range) {
        return `${range.startTime}-${range.endTime}`;
    } else {
        return '';
    }
});

const detailsUrl = computed<string>(() => {
    const dateRange: DateRange = currentPeriodConfig.value.dateType;
    return `/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: dateRange.type })}`
});
</script>
