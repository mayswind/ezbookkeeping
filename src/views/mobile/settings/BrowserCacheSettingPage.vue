<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Browser Cache Management')"></f7-nav-title>
            <f7-nav-right :class="{ 'disabled': loading }">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': loading || !fileCacheStatistics }" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('File Cache') }}</small>
            </f7-list-item>
            <f7-list-item title="Used storage" after="Count"></f7-list-item>
            <f7-list-item title="Application Code" after="Count"></f7-list-item>
            <f7-list-item title="Resource Files" after="Count"></f7-list-item>
            <f7-list-item title="Map Data" after="Count"></f7-list-item>
            <f7-list-item title="Others" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('Exchange Rates Data') }}</small>
            </f7-list-item>
            <f7-list-item title="Used storage" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('File Cache') }}</small>
            </f7-list-item>
            <f7-list-item :title="tt('Used storage')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.totalCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Application Code')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.codeCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Resource Files')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.assetsCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Map Data')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.mapCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Others')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.othersCacheSize, 2) : '-'"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('Exchange Rates Data') }}</small>
            </f7-list-item>
            <f7-list-item :title="tt('Used storage')" :after="formatVolumeToLocalizedNumerals(exchangeRatesCacheSize ?? 0, 2)"></f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading || !fileCacheStatistics }"
                                   @click="clearFileCache">{{ tt('Clear File Cache') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { type BrowserCacheStatistics } from '@/core/cache.ts';

import { loadBrowserCacheStatistics, clearAllBrowserCaches } from '@/lib/cache.ts';

const { tt, formatVolumeToLocalizedNumerals } = useI18n();
const { showConfirm, showToast } = useI18nUIComponents();

const exchangeRatesStore = useExchangeRatesStore();

const loading = ref<boolean>(true);
const showMoreActionSheet = ref<boolean>(false);
const fileCacheStatistics = ref<BrowserCacheStatistics | undefined>(undefined);
const exchangeRatesCacheSize = ref<number | undefined>(undefined);

function reloadCacheStatistics(): void {
    loading.value = true;

    loadBrowserCacheStatistics().then(statistics => {
        fileCacheStatistics.value = statistics;
        exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();
        loading.value = false;
    }).catch(() => {
        loading.value = false;
        showToast('Failed to load browser cache statistics');
    });
}

function clearFileCache(): void {
    showConfirm('Are you sure you want to clear file cache?', () => {
        clearAllBrowserCaches().then(() => {
            location.reload();
        });
    });
}

reloadCacheStatistics();
</script>
