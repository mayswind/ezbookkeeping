<template>
    <f7-list strong inset dividers class="overview-widget-list expense-ranking-widget no-margin-top margin-bottom">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || tt('Expense Category Ranking') }}</small>
        </f7-list-item>

        <template v-if="loading && !rankingItems.length">
            <f7-list-item class="statistics-list-item item-no-divider skeleton-text" link="#" :key="itemIdx" v-for="itemIdx in itemCount">
                <template #media>
                    <div class="display-flex no-padding-horizontal">
                        <div class="display-flex align-items-center statistics-icon">
                            <f7-icon f7="app_fill"></f7-icon>
                        </div>
                    </div>
                </template>
                <template #title>
                    <div class="statistics-list-item-text">
                        <span>Category Name</span>
                        <small class="statistics-percent">33.33</small>
                    </div>
                </template>
                <template #after>
                    <span>0.00 USD</span>
                </template>
                <template #inner-end>
                    <div class="statistics-item-end">
                        <div class="statistics-percent-line">
                            <f7-progressbar></f7-progressbar>
                        </div>
                    </div>
                </template>
            </f7-list-item>
        </template>

        <f7-list-item :title="tt('No data')" v-else-if="!loading && !rankingItems.length"></f7-list-item>

        <template v-else-if="rankingItems.length">
            <f7-list-item class="statistics-list-item item-no-divider"
                          :link="getTransactionItemLinkUrl(item.id)"
                          :key="item.id"
                          v-for="item in rankingItems"
            >
                <template #media>
                    <div class="display-flex no-padding-horizontal">
                        <div class="display-flex align-items-center statistics-icon">
                            <ItemIcon :icon-type="getCategoryIconType(item.iconType)" :icon-id="item.icon" :color="item.color" v-if="item.icon" />
                            <f7-icon f7="pencil_ellipsis_rectangle" v-else-if="!item.icon"></f7-icon>
                        </div>
                    </div>
                </template>

                <template #title>
                    <div class="statistics-list-item-text">
                        <span>{{ item.name }}</span>
                        <small class="statistics-percent" v-if="item.percent >= 0 && item.value.isPositiveOrZero()">{{ formatPercentToLocalizedNumerals(item.percent, 2, '<0.01') }}</small>
                    </div>
                </template>

                <template #after>
                    <span>{{ getDisplayAmount(item.value, rankingData.incomplete) }}</span>
                </template>

                <template #inner-end>
                    <div class="statistics-item-end">
                        <div class="statistics-percent-line">
                            <f7-progressbar :progress="item.percent >= 0 ? item.percent : 0" :style="{ '--f7-progressbar-progress-color': (item.color ? getCategoryDisplayColor(item.color) : '') } "></f7-progressbar>
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
    type CommonExpenseCategoryRankingWidgetProps,
    useExpenseCategoryRankingWidgetBase
} from '@/views/base/overview/ExpenseCategoryRankingWidgetBase.ts';

import { getCategoryIconType } from '@/lib/icon.ts';
import { getCategoryDisplayColor } from '@/lib/color.ts';

interface MobileExpenseCategoryRankingWidgetProps extends CommonExpenseCategoryRankingWidgetProps {
    showTitle: boolean;
}

const props = defineProps<MobileExpenseCategoryRankingWidgetProps>();

const { tt, formatPercentToLocalizedNumerals } = useI18n();
const {
    rankingData,
    rankingItems,
    getDisplayAmount,
    getTransactionItemLinkUrl
} = useExpenseCategoryRankingWidgetBase(props);
</script>

<style>
.expense-ranking-widget {
    .statistics-list-item .item-content {
        margin-top: 0;
    }
}
</style>
