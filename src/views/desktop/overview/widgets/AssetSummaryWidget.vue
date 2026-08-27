<template>
    <v-card class="overview-widget asset-summary-widget h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Asset Summary')" :icon="mdiBankOutline" />
        </template>

        <v-card-text class="overview-widget__body mt-3">
            <div class="overview-widget__caption">
                <span class="text-body-medium" v-if="!loading || (allAccounts && allAccounts.length)">{{ tt('format.misc.youHaveAccounts', { count: displayAccountCount }) }}</span>
                <v-skeleton-loader class="skeleton-no-margin py-1" width="200px" type="text" v-else-if="loading && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
            </div>

            <div class="asset-summary-widget__metrics mt-5">
                <div class="asset-summary-widget__metric" :class="{ 'asset-summary-widget__metric--primary': item.color === 'primary' }"
                     :key="item.title" v-for="item in summaryItems">
                    <v-avatar rounded="lg" variant="tonal" :color="item.color" size="40">
                        <v-icon size="24" :icon="item.icon" />
                    </v-avatar>
                    <div class="d-flex flex-column gap-1 text-truncate">
                        <span class="overview-widget__caption text-body-medium">{{ tt(item.title) }}</span>
                        <span class="overview-widget__amount text-title-large text-truncate" :class="{ 'text-primary': item.color === 'primary' }" v-if="!loading || (allAccounts && allAccounts.length)">{{ item.value }}</span>
                        <v-skeleton-loader class="skeleton-no-margin pb-2" style="margin-top: 6px" width="120px" type="text" :loading="true" v-else-if="loading && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
                    </div>
                </div>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAssetSummaryWidgetBase } from '@/views/base/overview/AssetSummaryWidgetBase.ts';

import {
    mdiBankOutline,
    mdiCreditCardOutline,
    mdiPiggyBankOutline
} from '@mdi/js';

interface SummaryItem {
    title: string;
    value: string;
    icon: string;
    color: string;
}

defineProps<{
    loading: boolean;
    title?: string;
}>();

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const {
    allAccounts,
    netAssets,
    totalAssets,
    totalLiabilities
} = useAssetSummaryWidgetBase();

const displayAccountCount = computed<string>(() => formatNumberToLocalizedNumerals(allAccounts.value?.length ?? 0));

const summaryItems = computed<SummaryItem[]>(() => [
    {
        title: 'Total assets',
        value: totalAssets.value,
        icon: mdiBankOutline,
        color: 'grey'
    },
    {
        title: 'Total liabilities',
        value: totalLiabilities.value,
        icon: mdiCreditCardOutline,
        color: 'expense'
    },
    {
        title: 'Net assets',
        value: netAssets.value,
        icon: mdiPiggyBankOutline,
        color: 'primary'
    }
]);
</script>
