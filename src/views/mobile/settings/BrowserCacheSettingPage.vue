<template>
    <f7-page ptr @ptr:refresh="reloadCacheStatistics">
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Browser Cache Management')"></f7-nav-title>
            <f7-nav-right :class="{ 'disabled': loading }">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': loading || ((!isSupportedFileCache || !fileCacheStatistics) && !exchangeRatesCacheSize) }" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading && isSupportedFileCache">
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

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('Cache Expiration Time') }}</small>
            </f7-list-item>
            <f7-list-item title="Map Data" after="Disable Cache" v-if="getMapProvider()"></f7-list-item>
            <f7-list-item title="Exchange Rates Data" after="Disable Cache"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading && isSupportedFileCache">
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

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('Cache Expiration Time') }}</small>
            </f7-list-item>
            <f7-list-item
                link="#"
                :class="{ 'disabled': loading || !isSupportedFileCache || !fileCacheStatistics || isMapProviderUseExternalSDK() || !isMapDataFetchProxyEnabled() }"
                :title="tt('Map Data')"
                :after="findNameByValue(allMapCacheExpirationOptions, mapCacheExpiration)"
                @click="showMapDataCacheExpirationPopup = true"
                v-if="getMapProvider()"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="value" value-field="value"
                                           title-field="name"
                                           :title="tt('Cache Expiration for Map Data')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Expiration Time')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allMapCacheExpirationOptions"
                                           v-model:show="showMapDataCacheExpirationPopup"
                                           v-model="mapCacheExpiration">
                </list-item-selection-popup>
            </f7-list-item>
            <f7-list-item
                link="#"
                :class="{ 'disabled': loading }"
                :title="tt('Exchange Rates Data')"
                :after="findNameByValue(allExchangeRatesDataCacheExpirationOptions, exchangeRatesDataCacheExpiration)"
                @click="showExchangeRatesDataCacheExpirationPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="value" value-field="value"
                                           title-field="name"
                                           :title="tt('Cache Expiration for Exchange Rates Data')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Expiration Time')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allExchangeRatesDataCacheExpirationOptions"
                                           v-model:show="showExchangeRatesDataCacheExpirationPopup"
                                           v-model="exchangeRatesDataCacheExpiration">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group v-if="isSupportedFileCache && fileCacheStatistics">
                <f7-actions-button :class="{ 'disabled': loading || !isSupportedFileCache || !fileCacheStatistics }"
                                   @click="clearMapCache">{{ tt('Clear Map Data Cache') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': loading || !isSupportedFileCache || !fileCacheStatistics }"
                                   @click="clearAllFileCache">{{ tt('Clear All File Cache') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading || !exchangeRatesCacheSize }"
                                   @click="clearExchangeRatesCache">{{ tt('Clear Exchange Rates Data Cache') }}</f7-actions-button>
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

import { useAppBrowserCacheSettingPageBase } from '@/views/base/settings/AppBrowserCacheSettingPageBase.ts';

import { findNameByValue } from '@/lib/common.ts';
import { isMapProviderUseExternalSDK } from '@/lib/map/index.ts';
import { getMapProvider, isMapDataFetchProxyEnabled } from '@/lib/server_settings.ts';

const { tt, formatVolumeToLocalizedNumerals } = useI18n();
const { showConfirm } = useI18nUIComponents();

const {
    isSupportedFileCache,
    loading,
    fileCacheStatistics,
    exchangeRatesCacheSize,
    allMapCacheExpirationOptions,
    allExchangeRatesDataCacheExpirationOptions,
    mapCacheExpiration,
    exchangeRatesDataCacheExpiration,
    loadCacheStatistics,
    clearMapDataCache,
    clearAllBrowserCaches,
    clearExchangeRatesDataCache
} = useAppBrowserCacheSettingPageBase();

const showMapDataCacheExpirationPopup = ref<boolean>(false);
const showExchangeRatesDataCacheExpirationPopup = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);

function reloadCacheStatistics(done?: () => void): void {
    loadCacheStatistics(false).then(() => {
        done?.();
    }).catch(() => {
        done?.();
    });
}

function clearMapCache(): void {
    showConfirm('Are you sure you want to clear map data cache?', () => {
        clearMapDataCache().then(() => {
            loadCacheStatistics(true);
        });
    });
}

function clearAllFileCache(): void {
    showConfirm('Are you sure you want to clear all file cache?', () => {
        clearAllBrowserCaches().then(() => {
            location.reload();
        });
    });
}

function clearExchangeRatesCache(): void {
    showConfirm('Are you sure you want to clear exchange rates data cache?', () => {
        clearExchangeRatesDataCache();
    });
}

loadCacheStatistics(true);
</script>
