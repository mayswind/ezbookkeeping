<template>
    <monthly-expense-overview-widget :loading="loading"
                                     @refresh="$emit('refresh')"
                                     v-if="widget.type === OverviewWidgetType.CurrentMonthOverview" />

    <asset-summary-widget :loading="loading"
                          v-else-if="widget.type === OverviewWidgetType.AssetSummary" />

    <period-income-expense-widget :loading="loading" :editing="editing"
                                  :period="widget.settings['dateRange'] as OverviewPeriod"
                                  v-else-if="widget.type === OverviewWidgetType.PeriodIncomeExpense" />

    <income-expense-trend-widget :loading="loading"
                                 :months="widget.settings['months'] as number"
                                 v-else-if="widget.type === OverviewWidgetType.IncomeExpenseTrend" />
</template>

<script setup lang="ts">
import MonthlyExpenseOverviewWidget from './widgets/MonthlyExpenseOverviewWidget.vue';
import AssetSummaryWidget from './widgets/AssetSummaryWidget.vue';
import PeriodIncomeExpenseWidget from './widgets/PeriodIncomeExpenseWidget.vue';
import IncomeExpenseTrendWidget from './widgets/IncomeExpenseTrendWidget.vue';

import {
    type OverviewPeriod,
    type DesktopOverviewWidgetLayout,
    OverviewWidgetType
} from '@/core/overview_layout.ts';

defineProps<{
    widget: DesktopOverviewWidgetLayout;
    loading: boolean;
    editing?: boolean
}>();

defineEmits<{
    (e: 'refresh'): void
}>();
</script>
