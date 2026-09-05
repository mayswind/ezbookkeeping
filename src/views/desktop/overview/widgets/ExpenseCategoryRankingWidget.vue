<template>
    <v-card class="overview-widget expense-ranking-widget h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Expense Category Ranking')" :icon="mdiViewDashboardOutline" />
        </template>

        <v-card-text class="overview-widget__body">
            <v-list class="py-0" lines="two" v-if="rankingItems.length">
                <template :key="item.id" v-for="item in rankingItems">
                    <v-list-item class="px-0 py-1 mb-1 no-min-height" density="compact">
                        <template #prepend>
                            <router-link class="overview-widget__item-icon" :to="getTransactionItemLinkUrl(item.id)">
                                <ItemIcon size="28px" :icon-type="getCategoryIconType(item.iconType)"
                                          :icon-id="item.icon" :color="item.color" />
                            </router-link>
                        </template>
                        <router-link class="ranking-list-item link-no-color" :to="getTransactionItemLinkUrl(item.id)">
                            <div class="d-flex flex-column">
                                <div class="expense-ranking-widget__label-row">
                                    <span class="text-truncate" :title="item.name">{{ item.name }}</span>
                                    <span class="overview-widget__amount ranking-amount">{{ getDisplayAmount(item.value, rankingData.incomplete) }}</span>
                                </div>
                                <div class="expense-ranking-widget__progress">
                                    <v-progress-linear rounded :color="item.color ? getCategoryDisplayColor(item.color) : 'primary'"
                                                       :bg-opacity="0.1" :aria-label="item.name"
                                                       :model-value="item.percent" :height="5" />
                                    <small class="ranking-percent">{{ formatPercentToLocalizedNumerals(item.percent, 2, '<0.01') }}</small>
                                </div>
                            </div>
                        </router-link>
                    </v-list-item>
                </template>
            </v-list>
            <div v-if="loading && !rankingItems.length">
                <v-skeleton-loader class="skeleton-no-margin py-4 mb-1" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="overview-widget__empty" v-if="!loading && !rankingItems.length">
                <v-icon :icon="mdiViewDashboardOutline" size="32" />
                <span>{{ tt('No data') }}</span>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type CommonExpenseCategoryRankingWidgetProps,
    useExpenseCategoryRankingWidgetBase
} from '@/views/base/overview/ExpenseCategoryRankingWidgetBase.ts';

import { getCategoryIconType } from '@/lib/icon.ts';
import { getCategoryDisplayColor } from '@/lib/color.ts';

import {
    mdiViewDashboardOutline
} from '@mdi/js';

const props = defineProps<CommonExpenseCategoryRankingWidgetProps>();

const { tt, formatPercentToLocalizedNumerals } = useI18n();
const {
    rankingData,
    rankingItems,
    getDisplayAmount,
    getTransactionItemLinkUrl
} = useExpenseCategoryRankingWidgetBase(props);
</script>

<style scoped>
.ranking-list-item {
    display: block;
    width: 100%;
    overflow: hidden;
    text-decoration: none;
}

.ranking-percent {
    flex: 0 0 auto;
    font-size: 0.75rem;
    opacity: 0.7;
    margin-inline-start: 6px;
}

.ranking-amount {
    flex: 0 0 auto;
    opacity: 0.8;
}
</style>
