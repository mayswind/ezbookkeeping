<template>
    <asset-summary-widget :loading="loading" :title="widgetTitle"
                          v-if="widget.type === OverviewWidgetType.AssetSummary" />

    <account-balance-list-widget :loading="loading" :editing="editing" :title="widgetTitle"
                                 :account-categories="widget.settings['accountCategories'] as number[]"
                                 :item-count="widget.settings['itemCount'] as number"
                                 :sort-by="widget.settings['sortBy'] as string"
                                 v-else-if="widget.type === OverviewWidgetType.AccountBalanceList" />

    <monthly-expense-overview-widget :loading="loading"
                                     @refresh="$emit('refresh')"
                                     v-else-if="widget.type === OverviewWidgetType.CurrentMonthOverview" />

    <monthly-expense-progress-widget :loading="loading" :title="widgetTitle"
                                     v-else-if="widget.type === OverviewWidgetType.CurrentMonthExpenseProgress" />

    <period-income-expense-widget :loading="loading" :editing="editing" :title="widgetTitle"
                                  :date-type="widget.settings['dateRange'] as number"
                                  v-else-if="widget.type === OverviewWidgetType.PeriodIncomeExpense" />

    <period-net-income-and-savings-rate-widget :loading="loading" :editing="editing" :title="widgetTitle"
                                               :date-type="widget.settings['dateRange'] as number"
                                               v-else-if="widget.type === OverviewWidgetType.PeriodNetIncomeAndSavingsRate" />

    <income-expense-trend-widget :loading="loading" :title="widgetTitle"
                                 :months="widget.settings['months'] as number"
                                 :show-x-axis-labels="widget.settings['showXAxisLabels'] as boolean"
                                 :show-legend="widget.settings['showLegend'] as boolean"
                                 v-else-if="widget.type === OverviewWidgetType.IncomeExpenseTrend" />

    <net-assets-trend-widget :loading="loading" :title="widgetTitle" :months="widget.settings['months'] as number"
                             :show-x-axis-labels="widget.settings['showXAxisLabels'] as boolean"
                             :show-legend="widget.settings['showLegend'] as boolean"
                             v-else-if="widget.type === OverviewWidgetType.NetAssetsTrend" />

    <expense-category-ranking-widget :loading="loading" :title="widgetTitle"
                                      :date-type="widget.settings['dateRange'] as number"
                                      :category-level="widget.settings['categoryLevel'] as string"
                                      :item-count="widget.settings['itemCount'] as number"
                                      v-else-if="widget.type === OverviewWidgetType.ExpenseCategoryRanking" />

    <recent-transactions-widget :loading="loading" :editing="editing" :title="widgetTitle"
                                :item-count="widget.settings['itemCount'] as number"
                                @refresh="$emit('refresh')"
                                v-else-if="widget.type === OverviewWidgetType.RecentTransactions" />

    <transaction-calendar-heatmap-widget :loading="loading" :editing="editing" :title="widgetTitle"
                                         :transaction-type="widget.settings['transactionType'] as TransactionType"
                                         :months="widget.settings['months'] as number"
                                         v-else-if="widget.type === OverviewWidgetType.TransactionCalendarHeatmap" />
</template>

<script setup lang="ts">
import { computed } from 'vue';

import MonthlyExpenseOverviewWidget from './widgets/MonthlyExpenseOverviewWidget.vue';
import MonthlyExpenseProgressWidget from './widgets/MonthlyExpenseProgressWidget.vue';
import AssetSummaryWidget from './widgets/AssetSummaryWidget.vue';
import PeriodIncomeExpenseWidget from './widgets/PeriodIncomeExpenseWidget.vue';
import PeriodNetIncomeAndSavingsRateWidget from './widgets/PeriodNetIncomeAndSavingsRateWidget.vue';
import IncomeExpenseTrendWidget from './widgets/IncomeExpenseTrendWidget.vue';
import NetAssetsTrendWidget from './widgets/NetAssetsTrendWidget.vue';
import AccountBalanceListWidget from './widgets/AccountBalanceListWidget.vue';
import ExpenseCategoryRankingWidget from './widgets/ExpenseCategoryRankingWidget.vue';
import RecentTransactionsWidget from './widgets/RecentTransactionsWidget.vue';
import TransactionCalendarHeatmapWidget from './widgets/TransactionCalendarHeatmapWidget.vue';

import { TransactionType } from '@/core/transaction.ts';
import {
    type DesktopOverviewWidgetLayout,
    OverviewWidgetType
} from '@/core/overview_layout.ts';

const props = defineProps<{
    widget: DesktopOverviewWidgetLayout;
    loading: boolean;
    editing?: boolean
}>();

const widgetTitle = computed<string>(() => {
    const title = props.widget.settings['title'];
    return typeof title === 'string' ? title.trim() : '';
});

defineEmits<{
    (e: 'refresh'): void
}>();
</script>
