<template>
    <v-card class="h-100" :class="{ 'disabled': loading }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="grey" size="32">
                <v-icon size="22" :icon="currentPeriodIcon" />
            </v-avatar>
            <span class="text-title-small font-weight-bold ms-2">{{ title || currentPeriodTitle }}</span>
            <v-spacer/>
            <v-btn density="comfortable" color="default" variant="text" class="ms-2" :icon="true" v-if="!editing">
                <v-icon :icon="mdiDotsVertical" />
                <v-menu activator="parent">
                    <v-list>
                        <v-list-item :prepend-icon="mdiListBoxOutline" :to="currentDetailsUrl">
                            <v-list-item-title>{{ tt('View Details') }}</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-btn>
            <div class="v-btn--size-default" style="height: var(--v-btn-height)" v-else-if="editing"></div><!-- add a button placeholder to avoid height change -->
        </v-card-text>
        <v-card-text class="py-3">
            <div class="text-truncate text-headline-small text-income me-2 mb-2" v-if="!loading || currentDisplayIncomeAmount">{{ currentDisplayIncomeAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-3 mb-6" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayIncomeAmount"></v-skeleton-loader>
            <div class="text-truncate text-title-medium text-expense pt-1" v-if="!loading || currentDisplayExpenseAmount">{{ currentDisplayExpenseAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin pt-1" style="padding-bottom: 7px" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayExpenseAmount"></v-skeleton-loader>
            <div class="text-truncate text-title-medium mt-3 mb-5" v-if="!loading && !currentDisplayIncomeAmount && !currentDisplayExpenseAmount">{{ tt('No data') }}</div>
        </v-card-text>
        <v-card-text>
            <span class="text-body-medium text-high-emphasis">{{ currentDisplayDateTime }}</span>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
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
