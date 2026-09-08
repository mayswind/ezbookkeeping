<template>
    <f7-card class="asset-summary-widget no-margin-top margin-bottom" :class="{ 'skeleton-text': loading }">
        <f7-card-content class="padding-horizontal padding-vertical">
            <div class="asset-summary-widget__header display-flex align-items-baseline justify-content-space-between">
                <span class="asset-summary-widget__title font-weight-bold">{{ tt('Asset Summary') }}</span>
                <f7-link class="margin-inline-start-half text-color-gray"
                       :aria-label="showAmountInHomePage ? tt('Hide Amount') : tt('Show Amount')"
                       @click="showAmountInHomePage = !showAmountInHomePage">
                    <f7-icon :f7="showAmountInHomePage ? 'eye_slash' : 'eye'" size="20"></f7-icon>
                </f7-link>
            </div>

            <div class="asset-summary-widget__caption margin-top-half text-color-gray">
                <span v-if="!loading || (allAccounts && allAccounts.length)">{{ tt('format.misc.youHaveAccounts', { count: displayAccountCount }) }}</span>
                <span v-else>Loading...</span>
            </div>

            <div class="asset-summary-widget__metrics margin-top">
                    <div class="asset-summary-widget__metric">
                        <div class="asset-summary-widget__metric-icon text-color-gray" style="background-color: rgba(128, 128, 128, 0.15);">
                            <f7-icon f7="briefcase" size="24"></f7-icon>
                        </div>
                        <div class="asset-summary-widget__metric-text display-flex flex-direction-column">
                            <span class="asset-summary-widget__metric-title">{{ tt('Total assets') }}</span>
                            <span class="asset-summary-widget__metric-amount font-weight-bold" v-if="!loading || (allAccounts && allAccounts.length)">{{ totalAssets }}</span>
                            <span class="asset-summary-widget__metric-amount font-weight-bold" v-else>0.00 USD</span>
                        </div>
                    </div>

                    <div class="asset-summary-widget__metric">
                        <div class="asset-summary-widget__metric-icon text-color-red" style="background-color: rgba(255, 59, 48, 0.1);">
                            <f7-icon f7="creditcard" size="24"></f7-icon>
                        </div>
                        <div class="asset-summary-widget__metric-text display-flex flex-direction-column">
                            <span class="asset-summary-widget__metric-title">{{ tt('Total liabilities') }}</span>
                            <span class="asset-summary-widget__metric-amount font-weight-bold" v-if="!loading || (allAccounts && allAccounts.length)">{{ totalLiabilities }}</span>
                            <span class="asset-summary-widget__metric-amount font-weight-bold" v-else>0.00 USD</span>
                        </div>
                    </div>
                    
                    <div class="asset-summary-widget__metric">
                        <div class="asset-summary-widget__metric-icon text-color-primary" style="background-color: rgba(var(--f7-theme-color-rgb), 0.1);">
                            <f7-icon f7="money_dollar_circle" size="24"></f7-icon>
                        </div>
                        <div class="asset-summary-widget__metric-text display-flex flex-direction-column">
                            <span class="asset-summary-widget__metric-title">{{ tt('Net assets') }}</span>
                            <span class="asset-summary-widget__metric-amount font-weight-bold text-color-primary" v-if="!loading || (allAccounts && allAccounts.length)">{{ netAssets }}</span>
                            <span class="asset-summary-widget__metric-amount font-weight-bold text-color-primary" v-else>0.00 USD</span>
                        </div>
                    </div>
            </div>
        </f7-card-content>
    </f7-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { useAssetSummaryWidgetBase } from '@/views/base/overview/AssetSummaryWidgetBase.ts';

defineProps<{
    loading: boolean;
}>();

const { tt, formatNumberToLocalizedNumerals } = useI18n();
const {
    showAmountInHomePage,
    allAccounts,
    netAssets,
    totalAssets,
    totalLiabilities
} = useAssetSummaryWidgetBase();

const displayAccountCount = computed<string>(() => formatNumberToLocalizedNumerals(allAccounts.value?.length ?? 0));


</script>

<style scoped>
.asset-summary-widget__title {
    font-size: 1.25rem;
}

.asset-summary-widget__caption {
    font-size: 0.85rem;
}

.asset-summary-widget__metrics {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
}

.asset-summary-widget__metric {
    display: flex;
    align-items: center;
    gap: 0.75rem;
}

.asset-summary-widget__metric-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.asset-summary-widget__metric-text {
    flex: 1;
    min-width: 0;
    gap: 2px;
}

.asset-summary-widget__metric-title {
    font-size: 0.8rem;
}

.asset-summary-widget__metric-amount {
    font-size: 1.15rem;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
