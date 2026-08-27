<template>
    <v-card class="overview-widget period-income-expense-widget h-100" :class="{ 'disabled': loading }">
        <template #title>
            <overview-widget-header :title="title || currentPeriodTitle" :icon="currentPeriodIcon">
                <v-btn density="compact" color="default" variant="text" :icon="true" :aria-label="tt('More')" v-if="!editing">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiListBoxOutline" :to="currentDetailsUrl">
                                <v-list-item-title>{{ tt('View Details') }}</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </overview-widget-header>
        </template>
        <v-card-text class="overview-widget__body period-income-expense-widget__body d-flex flex-column">
            <div class="period-income-expense-widget__content">
                <div class="period-income-expense-widget__metrics">
                    <div class="overview-widget__stat">
                        <span class="overview-widget__caption"><span class="overview-widget__dot bg-income"></span>{{ tt('Income') }}</span>
                        <div class="overview-widget__amount text-headline-small text-income text-truncate" v-if="!loading || currentDisplayIncomeAmount">{{ currentDisplayIncomeAmount || '-' }}</div>
                        <v-skeleton-loader class="skeleton-no-margin my-2" type="text" width="120px" :loading="true" v-else></v-skeleton-loader>
                    </div>
                    <div class="overview-widget__stat">
                        <span class="overview-widget__caption"><span class="overview-widget__dot bg-expense"></span>{{ tt('Expense') }}</span>
                        <div class="overview-widget__amount text-headline-small text-expense text-truncate" v-if="!loading || currentDisplayExpenseAmount">{{ currentDisplayExpenseAmount || '-' }}</div>
                        <v-skeleton-loader class="skeleton-no-margin my-2" type="text" width="120px" :loading="true" v-else></v-skeleton-loader>
                    </div>
                    <div class="overview-widget__caption" v-if="!loading && !currentDisplayIncomeAmount && !currentDisplayExpenseAmount">{{ tt('No data') }}</div>
                </div>
                <svg class="period-income-expense-widget__illustration img-with-direction mt-2" viewBox="0 0 144 112"
                     fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true" focusable="false">
                    <circle cx="78" cy="57" r="45" fill="currentColor" opacity="0.055" />
                    <circle cx="124" cy="20" r="3" fill="currentColor" opacity="0.2" />
                    <circle cx="18" cy="89" r="2" fill="currentColor" opacity="0.2" />
                    <rect class="period-income-expense-widget__paper" x="40" y="17" width="65" height="80" rx="10" />
                    <path d="M54 34H76M54 44H90M54 54H80" stroke="currentColor" stroke-width="3" stroke-linecap="round" opacity="0.18" />
                    <rect x="53" y="69" width="39" height="14" rx="5" fill="currentColor" opacity="0.08" />
                    <path d="M60 76H76" stroke="currentColor" stroke-width="2" stroke-linecap="round" opacity="0.3" />
                    <g class="period-income-expense-widget__flow">
                        <circle cx="34" cy="44" r="17" />
                        <path d="M26 44H42M36 38L42 44L36 50" />
                        <circle cx="108" cy="79" r="17" />
                        <path d="M100 79H116M106 73L100 79L106 85" />
                    </g>
                </svg>
            </div>
            <div class="overview-widget__footer">{{ currentDisplayDateTime }}</div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonPeriodStatisticsWidgetProps, usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

import { DateRange } from '@/core/datetime.ts';

import { isInteger } from '@/lib/common.ts';

import {
    mdiDotsVertical,
    mdiCalendarTodayOutline,
    mdiCalendarWeekOutline,
    mdiCalendarMonthOutline,
    mdiLayersTripleOutline,
    mdiListBoxOutline
} from '@mdi/js';

interface PeriodIncomeExpenseWidgetProps extends CommonPeriodStatisticsWidgetProps {
    loading?: boolean;
    editing?: boolean;
    title?: string;
}

const props = defineProps<PeriodIncomeExpenseWidgetProps>();

const { tt } = useI18n();
const {
    currentPeriodTitle,
    currentDisplayDateTime,
    currentDisplayIncomeAmount,
    currentDisplayExpenseAmount,
    currentDetailsUrl
} = usePeriodStatisticsWidgetBase(props);

const allSupportedPeriodIcons: Record<number, string> = {
    [DateRange.Today.type]: mdiCalendarTodayOutline,
    [DateRange.ThisWeek.type]: mdiCalendarWeekOutline,
    [DateRange.ThisMonth.type]: mdiCalendarMonthOutline,
    [DateRange.ThisYear.type]: mdiLayersTripleOutline
};

const currentPeriodIcon = computed<string>(() => isInteger(props.dateType) ? (allSupportedPeriodIcons[props.dateType] || mdiCalendarTodayOutline) : mdiCalendarTodayOutline);
</script>
