<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="primary" size="32">
                <v-icon size="22" :icon="mdiPiggyBankOutline" />
            </v-avatar>
            <span class="text-title-small font-weight-bold ms-2">{{ title || currentPeriodTitle }}</span>
            <v-spacer />
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
            <div class="text-headline-small font-weight-medium me-2 mb-2"
                 :class="{ 'text-expense': !!currentDisplayNetIncomeAmount && currentNetIncomeAmount.isNegative(), 'text-income': !!currentDisplayNetIncomeAmount && !currentNetIncomeAmount.isNegative() }"
                 v-if="!loading || currentDisplayNetIncomeAmount">{{ currentDisplayNetIncomeAmount !== '' ? currentDisplayNetIncomeAmount : tt('No data') }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-3 mb-6" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayNetIncomeAmount"></v-skeleton-loader>
            <div class="d-flex align-center" v-if="!loading || currentDisplayNetIncomeAmount">
                <span class="text-body-medium" v-if="currentDisplaySavingsRate">{{ tt('Savings Rate') }}</span>
                <v-spacer />
                <span class="text-title-medium">{{ currentDisplaySavingsRate }}</span>
            </div>
            <v-skeleton-loader class="skeleton-no-margin" style="padding-bottom: 7px" type="text" width="120px" :loading="true" v-else-if="loading && !currentDisplayNetIncomeAmount"></v-skeleton-loader>
            <div class="d-flex text-high-emphasis mt-8" v-if="!loading || currentDisplayNetIncomeAmount">
                <span v-if="currentDisplayNetIncomeAmount">{{ tt('Income') }}</span>
                <span class="ms-1">{{ currentDisplayIncomeAmount }}</span>
                <v-spacer />
                <span v-if="currentDisplayExpenseAmount">{{ tt('Expense') }}</span>
                <span class="ms-1">{{ currentDisplayExpenseAmount }}</span>
            </div>
            <div class="d-flex text-high-emphasis mt-8" v-if="loading && !currentDisplayNetIncomeAmount">
                <span>{{ tt('Income') }}</span>
                <v-skeleton-loader class="skeleton-no-margin ms-1" type="text" width="80px" :loading="true"></v-skeleton-loader>
                <v-spacer />
                <span>{{ tt('Expense') }}</span>
                <v-skeleton-loader class="skeleton-no-margin ms-1" type="text" width="80px" :loading="true"></v-skeleton-loader>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
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
    currentNetIncomeAmount,
    currentDisplayIncomeAmount,
    currentDisplayExpenseAmount,
    currentDisplayNetIncomeAmount,
    currentDisplaySavingsRate,
    currentDetailsUrl
} = usePeriodStatisticsWidgetBase(props);
</script>
