<template>
    <asset-summary-widget :loading="loading" :height="widget.settings['height'] as number"
                          :light-background-color="widget.settings['lightBackgroundColor'] as ColorValue"
                          :dark-background-color="widget.settings['darkBackgroundColor'] as ColorValue"
                          v-if="widget.type === OverviewWidgetType.AssetSummary" />

    <account-balance-list-widget :loading="loading" :title="widgetTitle"
                                 :show-title="widget.settings['showTitle'] as boolean"
                                 :account-ids="widget.settings['accountIds'] as string[]"
                                 :item-count="widget.settings['itemCount'] as number"
                                 :sort-by="widget.settings['sortBy'] as string"
                                 :always-show-amount="widget.settings['alwaysShowAmount'] as boolean"
                                 v-else-if="widget.type === OverviewWidgetType.AccountBalanceList" />

    <monthly-expense-overview-widget :loading="loading" :height="widget.settings['height'] as number"
                                     :light-background-color="widget.settings['lightBackgroundColor'] as ColorValue"
                                     :dark-background-color="widget.settings['darkBackgroundColor'] as ColorValue"
                                     v-else-if="widget.type === OverviewWidgetType.CurrentMonthOverview" />

    <monthly-expense-progress-widget :loading="loading" :title="widgetTitle"
                                     :show-title="widget.settings['showTitle'] as boolean"
                                     v-else-if="widget.type === OverviewWidgetType.CurrentMonthExpenseProgress" />

    <period-income-expense-widget :loading="loading" :date-ranges="widget.settings['dateRanges'] as number[]"
                                  v-else-if="widget.type === OverviewWidgetType.PeriodIncomeExpense" />

    <period-net-income-and-savings-rate-widget :loading="loading" :title="widgetTitle"
                                               :show-title="widget.settings['showTitle'] as boolean"
                                               :date-type="widget.settings['dateRange'] as number"
                                               v-else-if="widget.type === OverviewWidgetType.PeriodNetIncomeAndSavingsRate" />
</template>

<script setup lang="ts">
import { computed } from 'vue';

import AssetSummaryWidget from './widgets/AssetSummaryWidget.vue';
import AccountBalanceListWidget from './widgets/AccountBalanceListWidget.vue';
import MonthlyExpenseOverviewWidget from './widgets/MonthlyExpenseOverviewWidget.vue';
import MonthlyExpenseProgressWidget from './widgets/MonthlyExpenseProgressWidget.vue';
import PeriodIncomeExpenseWidget from './widgets/PeriodIncomeExpenseWidget.vue';
import PeriodNetIncomeAndSavingsRateWidget from './widgets/PeriodNetIncomeAndSavingsRateWidget.vue';

import type { ColorValue } from '@/core/color.ts';
import { type MobileOverviewWidgetLayout, OverviewWidgetType } from '@/core/overview_layout.ts';

const props = defineProps<{
    widget: MobileOverviewWidgetLayout;
    loading: boolean;
}>();

const widgetTitle = computed<string>(() => {
    const title = props.widget.settings['title'];
    return typeof title === 'string' ? title.trim() : '';
});
</script>

<style>
.overview-widget-list .item-after {
    max-width: 50%;
}

.overview-widget-primary-amount {
    margin-top: 2px;
    margin-bottom: 6px;
    font-size: 1.5em;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
