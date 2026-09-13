<template>
    <f7-list strong inset dividers class="overview-widget-list no-margin-top margin-bottom">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || tt('Category Budgets') }}</small>
        </f7-list-item>

        <template v-if="loading && !budgetItems.length">
            <f7-list-item class="statistics-list-item item-no-divider skeleton-text" link="#"
                          :key="itemIdx" v-for="itemIdx in itemCount"
                          title="Category Name" after="0.00 / 0.00" />
        </template>

        <f7-list-item :title="tt('No category budgets')" v-else-if="!loading && !budgetItems.length" />

        <template v-else>
            <f7-list-item class="item-no-divider"
                          :header="tt('Spent')"
                          :title="`${getSummaryAmountText(budgetSummary.spent)} / ${getSummaryAmountText(budgetSummary.limit)}`">
                <template #inner-end>
                    <div class="statistics-item-end">
                        <div class="statistics-percent-line">
                            <f7-progressbar :progress="Math.max(0, Math.min(100, budgetSummary.percent))"
                                            :style="{ '--f7-progressbar-progress-color': budgetSummary.percent > 100 ? 'var(--f7-color-red)' : 'var(--f7-theme-color)' }" />
                        </div>
                    </div>
                </template>
            </f7-list-item>
            <f7-list-item class="statistics-list-item item-no-divider"
                          :link="getCategoryPageLink(item, true)"
                          :key="item.category.id" v-for="item in budgetItems">
                <template #media>
                    <div class="display-flex no-padding-horizontal">
                        <div class="display-flex align-items-center statistics-icon">
                            <ItemIcon :icon-type="getCategoryIconType(item.category.iconType)"
                                      :icon-id="item.category.icon" :color="item.category.color" />
                        </div>
                    </div>
                </template>
                <template #title>
                    <div class="statistics-list-item-text">
                        <span>{{ item.category.name }}</span>
                        <small class="statistics-percent" :class="{ 'text-color-red': item.overrun.isPositive() }">
                            {{ formatPercentToLocalizedNumerals(item.percent, 1, '<0.1') }}
                        </small>
                    </div>
                </template>
                <template #after>
                    <span>{{ getItemAmountText(item) }}</span>
                </template>
                <template #inner-end>
                    <div class="statistics-item-end">
                        <div class="statistics-percent-line">
                            <f7-progressbar :progress="Math.max(0, Math.min(100, item.percent))"
                                            :style="{ '--f7-progressbar-progress-color': item.overrun.isPositive() ? 'var(--f7-color-red)' : getCategoryDisplayColor(item.category.color) }" />
                        </div>
                    </div>
                </template>
            </f7-list-item>
        </template>
    </f7-list>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import {
    type CommonCategoryBudgetWidgetProps,
    useCategoryBudgetWidgetBase
} from '@/views/base/overview/CategoryBudgetWidgetBase.ts';

import { getCategoryIconType } from '@/lib/icon.ts';
import { getCategoryDisplayColor } from '@/lib/color.ts';

interface MobileCategoryBudgetWidgetProps extends CommonCategoryBudgetWidgetProps {
    showTitle: boolean;
}

const props = defineProps<MobileCategoryBudgetWidgetProps>();
const { tt, formatPercentToLocalizedNumerals } = useI18n();
const {
    budgetItems,
    budgetSummary,
    getSummaryAmountText,
    getItemAmountText,
    getCategoryPageLink
} = useCategoryBudgetWidgetBase(props);
</script>
