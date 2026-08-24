<template>
    <f7-card class="account-overview-card no-margin-top margin-bottom" :class="{ 'skeleton-text': loading }">
        <f7-card-header class="display-block" :style="style">
            <p class="no-margin">
                <small class="card-header-content" v-if="loading">Net assets</small>
                <small class="card-header-content" v-else-if="!loading">{{ tt('Net assets') }}</small>
            </p>
            <p class="no-margin">
                <span class="net-assets" v-if="loading">0.00 USD</span>
                <span class="net-assets" v-else-if="!loading">{{ netAssets }}</span>
            </p>
            <p class="no-margin">
                <small class="account-overview-info" v-if="loading">
                    <span>Total assets | Total liabilities</span>
                </small>
                <small class="account-overview-info" v-else-if="!loading">
                    <span>{{ tt('Total assets') }}</span>
                    <span>{{ totalAssets }}</span>
                    <span>|</span>
                    <span>{{ tt('Total liabilities') }}</span>
                    <span>{{ totalLiabilities }}</span>
                </small>
            </p>
        </f7-card-header>
    </f7-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAssetSummaryWidgetBase } from '@/views/base/overview/AssetSummaryWidgetBase.ts';

const props = defineProps<{
    loading: boolean;
    height: number;
}>();

const style = computed<Record<string, string>>(() => {
    const finalStyle: Record<string, string> = {};

    if (props.height === 1) {
        finalStyle['padding-top'] = '10px';
    } else if (props.height === 2) {
        finalStyle['padding-top'] = '60px';
    } else {
        finalStyle['padding-top'] = '120px';
    }

    return finalStyle;
});

const { tt } = useI18n();

const {
    netAssets,
    totalAssets,
    totalLiabilities
} = useAssetSummaryWidgetBase();
</script>
