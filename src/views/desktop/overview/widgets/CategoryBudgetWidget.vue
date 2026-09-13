<template>
    <v-card class="overview-widget category-budget-widget h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Category Budgets')" :icon="mdiWalletOutline" />
        </template>

        <v-card-text class="overview-widget__body">
            <div v-if="budgetItems.length">
                <div class="category-budget-widget__summary mb-3">
                    <div class="d-flex align-center justify-space-between text-body-medium">
                        <span>{{ tt('Spent') }}</span>
                        <span>{{ getSummaryAmountText(budgetSummary.spent) }} / {{ getSummaryAmountText(budgetSummary.limit) }}</span>
                    </div>
                    <v-progress-linear class="mt-2" rounded :color="budgetSummary.percent > 100 ? 'error' : 'primary'" :bg-opacity="0.1" height="7"
                                       :model-value="Math.max(0, Math.min(100, budgetSummary.percent))" />
                </div>

                <v-list class="py-0" lines="two">
                    <v-list-item class="px-0 py-1 no-min-height" density="compact"
                                 :key="item.category.id" v-for="item in budgetItems">
                        <template #prepend>
                            <router-link class="overview-widget__item-icon" :to="getCategoryPageLink(item, false)">
                                <ItemIcon size="28px" :icon-type="getCategoryIconType(item.category.iconType)"
                                          :icon-id="item.category.icon" :color="item.category.color" />
                            </router-link>
                        </template>
                        <router-link class="category-budget-widget__item link-no-color" :to="getCategoryPageLink(item, false)">
                            <div class="d-flex align-center justify-space-between">
                                <span class="text-truncate" :title="item.category.name">{{ item.category.name }}</span>
                                <span class="category-budget-widget__amount text-no-wrap">{{ getItemAmountText(item) }}</span>
                            </div>
                            <div class="d-flex align-center mt-1">
                                <v-progress-linear rounded
                                                   :color="item.overrun.isPositive() ? 'error' : getCategoryDisplayColor(item.category.color)"
                                                   :bg-opacity="0.1" height="5"
                                                   :model-value="Math.max(0, Math.min(100, item.percent))" />
                                <small class="category-budget-widget__percent ms-2"
                                       :class="{ 'text-error': item.overrun.isPositive() }">
                                    {{ formatPercentToLocalizedNumerals(item.percent, 1, '<0.1') }}
                                </small>
                            </div>
                        </router-link>
                    </v-list-item>
                </v-list>
            </div>

            <div v-if="loading && !budgetItems.length">
                <v-skeleton-loader class="skeleton-no-margin py-4 mb-1" type="text"
                                   :key="idx" :loading="true" v-for="idx in itemCount" />
            </div>
            <div class="overview-widget__empty" v-if="!loading && !budgetItems.length">
                <v-icon :icon="mdiWalletPlusOutline" size="32" />
                <span>{{ tt('No category budgets') }}</span>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type CommonCategoryBudgetWidgetProps,
    useCategoryBudgetWidgetBase
} from '@/views/base/overview/CategoryBudgetWidgetBase.ts';

import { getCategoryIconType } from '@/lib/icon.ts';
import { getCategoryDisplayColor } from '@/lib/color.ts';

import { mdiWalletOutline, mdiWalletPlusOutline } from '@mdi/js';

const props = defineProps<CommonCategoryBudgetWidgetProps>();
const { tt, formatPercentToLocalizedNumerals } = useI18n();
const {
    budgetItems,
    budgetSummary,
    getSummaryAmountText,
    getItemAmountText,
    getCategoryPageLink
} = useCategoryBudgetWidgetBase(props);
</script>

<style scoped>
.category-budget-widget__item {
    display: block;
    width: 100%;
    overflow: hidden;
    text-decoration: none;
}

.category-budget-widget__amount,
.category-budget-widget__percent {
    flex: 0 0 auto;
    font-size: 0.75rem;
    opacity: 0.8;
}
</style>
