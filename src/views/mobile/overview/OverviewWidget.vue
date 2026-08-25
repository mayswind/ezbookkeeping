<template>
    <asset-summary-widget :loading="loading" :height="widget.settings['height'] as number"
                          :light-background-color="widget.settings['lightBackgroundColor'] as ColorValue"
                          :dark-background-color="widget.settings['darkBackgroundColor'] as ColorValue"
                          v-if="widget.type === OverviewWidgetType.AssetSummary" />

    <monthly-expense-overview-widget :loading="loading" :height="widget.settings['height'] as number"
                                     :light-background-color="widget.settings['lightBackgroundColor'] as ColorValue"
                                     :dark-background-color="widget.settings['darkBackgroundColor'] as ColorValue"
                                     v-else-if="widget.type === OverviewWidgetType.CurrentMonthOverview" />

    <period-income-expense-widget :loading="loading" :date-ranges="widget.settings['dateRanges'] as number[]"
                                  v-else-if="widget.type === OverviewWidgetType.PeriodIncomeExpense" />
</template>

<script setup lang="ts">
import AssetSummaryWidget from './widgets/AssetSummaryWidget.vue';
import MonthlyExpenseOverviewWidget from './widgets/MonthlyExpenseOverviewWidget.vue';
import PeriodIncomeExpenseWidget from './widgets/PeriodIncomeExpenseWidget.vue';

import type { ColorValue } from '@/core/color.ts';
import { type MobileOverviewWidgetLayout, OverviewWidgetType } from '@/core/overview_layout.ts';

defineProps<{
    widget: MobileOverviewWidgetLayout;
    loading: boolean;
}>();
</script>
