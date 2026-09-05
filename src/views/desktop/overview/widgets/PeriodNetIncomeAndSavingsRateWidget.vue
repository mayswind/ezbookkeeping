<template>
    <v-card class="overview-widget savings-rate-widget h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || currentPeriodTitle" :icon="mdiPiggyBankOutline">
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
        <v-card-text class="overview-widget__body savings-rate-widget__body">
            <div class="overview-widget__caption">{{ tt('Net Income') }}</div>
            <div class="overview-widget__amount text-headline-small font-weight-medium text-truncate"
                 :class="{ 'text-income': !!currentDisplayNetIncomeAmount, 'text-medium-emphasis': !currentDisplayNetIncomeAmount }"
                 v-if="!loading || currentDisplayNetIncomeAmount">{{ currentDisplayNetIncomeAmount !== '' ? currentDisplayNetIncomeAmount : tt('No data') }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-3 mb-6" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayNetIncomeAmount"></v-skeleton-loader>
            <div class="d-flex justify-space-between align-center my-3">
                <span class="text-body-medium">{{ tt('Savings Rate') }}</span>
                <span class="text-title-medium text-primary" v-if="!loading || currentDisplayNetIncomeAmount">{{ currentDisplaySavingsRate || '-' }}</span>
                <v-skeleton-loader class="skeleton-no-margin pb-3" type="text" width="80px" :loading="true" v-else-if="loading && !currentDisplayNetIncomeAmount"></v-skeleton-loader>
            </div>
            <div class="savings-rate-widget__breakdown" v-if="!loading || currentDisplayNetIncomeAmount">
                <div>
                    <span class="overview-widget__caption">{{ tt('Income') }}</span>
                    <span class="overview-widget__amount text-income text-truncate">{{ currentDisplayIncomeAmount ? currentDisplayIncomeAmount : '-' }}</span>
                </div>
                <div>
                    <span class="overview-widget__caption">{{ tt('Expense') }}</span>
                    <span class="overview-widget__amount text-expense text-truncate">{{ currentDisplayExpenseAmount ? currentDisplayExpenseAmount : '-' }}</span>
                </div>
            </div>
            <div class="savings-rate-widget__breakdown" v-if="loading && !currentDisplayNetIncomeAmount">
                <div>
                    <span>{{ tt('Income') }}</span>
                    <v-skeleton-loader class="skeleton-no-margin" type="text" width="80px" :loading="true"></v-skeleton-loader>
                </div>
                <div>
                    <span>{{ tt('Expense') }}</span>
                    <v-skeleton-loader class="skeleton-no-margin" type="text" width="80px" :loading="true"></v-skeleton-loader>
                </div>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonPeriodStatisticsWidgetProps, usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

import {
    mdiPiggyBankOutline,
    mdiDotsVertical,
    mdiListBoxOutline
} from '@mdi/js';

interface PeriodNetIncomeAndSavingsRateWidgetProps extends CommonPeriodStatisticsWidgetProps {
    loading?: boolean;
    editing?: boolean;
    title?: string;
}

const props = defineProps<PeriodNetIncomeAndSavingsRateWidgetProps>();

const { tt } = useI18n();
const {
    currentPeriodTitle,
    currentDisplayIncomeAmount,
    currentDisplayExpenseAmount,
    currentDisplayNetIncomeAmount,
    currentDisplaySavingsRate,
    currentDetailsUrl
} = usePeriodStatisticsWidgetBase(props);
</script>
