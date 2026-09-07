<template>
    <f7-list strong inset dividers class="overview-transaction-list no-margin-top margin-bottom">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || currentPeriodTitle }}</small>
        </f7-list-item>
        <f7-list-item class="combined-item item-no-divider">
            <template #header>
                <div class="margin-bottom-half">{{ tt('Net Income') }}</div>
            </template>
            <template #title>
                <span class="overview-widget-primary-amount skeleton-text" v-if="loading">0.00 USD</span>
                <span class="overview-widget-primary-amount"
                      :class="{ 'text-income': !!currentDisplayNetIncomeAmount }"
                      v-else-if="!loading">
                    {{ currentDisplayNetIncomeAmount || tt('No data') }}
                </span>
            </template>
        </f7-list-item>
        <f7-list-item class="item-title-full-line">
            <template #title>
                <div class="width-100 display-flex justify-content-space-between">
                    <span>{{ tt('Savings Rate') }}</span>
                    <span class="skeleton-text" v-if="loading">0.00 USD</span>
                    <span class="text-color-primary" v-else-if="!loading">{{ currentDisplaySavingsRate || '-' }}</span>
                </div>
            </template>
        </f7-list-item>
        <f7-list-item class="item-title-full-line">
            <template #footer>
                <div class="overview-transaction-footer default-text-color display-flex justify-content-space-between">
                    <span>{{ tt('Income') }}</span>
                    <span class="skeleton-text" v-if="loading">0.00 USD</span>
                    <span class="text-income" v-else-if="!loading">{{ currentDisplayIncomeAmount || '-' }}</span>
                </div>
                <div class="overview-transaction-footer default-text-color display-flex justify-content-space-between">
                    <span>{{ tt('Expense') }}</span>
                    <span class="skeleton-text" v-if="loading">0.00 USD</span>
                    <span class="text-expense" v-else-if="!loading">{{ currentDisplayExpenseAmount || '-' }}</span>
                </div>
            </template>
        </f7-list-item>
    </f7-list>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

const props = defineProps<{
    loading: boolean;
    title?: string;
    showTitle: boolean;
    dateType: number;
}>();

const { tt } = useI18n();
const {
    currentPeriodTitle,
    currentDisplayIncomeAmount,
    currentDisplayExpenseAmount,
    currentDisplayNetIncomeAmount,
    currentDisplaySavingsRate
} = usePeriodStatisticsWidgetBase(props);
</script>
