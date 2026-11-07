<template>
    <f7-page ptr @ptr:refresh="reloadCacheStatistics">
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Browser Cache Management')"></f7-nav-title>
            <f7-nav-right :class="{ 'disabled': loading }">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': ((!isSupportedFileCache || !fileCacheStatistics) && !exchangeRatesCacheSize) || loading || clearingCache }" :aria-label="tt('More')" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical-half skeleton-text" v-if="loading && isSupportedFileCache">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('File Cache') }}</small>
            </f7-list-item>
            <f7-list-item title="Used storage" after="Count"></f7-list-item>
            <f7-list-item title="Application Code" after="Count"></f7-list-item>
            <f7-list-item title="Resource Files" after="Count"></f7-list-item>
            <f7-list-item title="Map Data" after="Count"></f7-list-item>
            <f7-list-item title="User Avatars" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Pictures" after="Count"></f7-list-item>
            <f7-list-item title="Custom Icons" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers :class="{ 'margin-vertical': isSupportedFileCache, 'margin-vertical-half': !isSupportedFileCache, 'skeleton-text': true }" v-if="loading">
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
            <f7-list-item title="Transaction Pictures" after="30 Days"></f7-list-item>
            <f7-list-item title="Exchange Rates Data" after="Disable Cache"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical-half" v-if="!loading && isSupportedFileCache">
            <f7-list-item group-title :sortable="false">
                <small>{{ tt('File Cache') }}</small>
            </f7-list-item>
            <f7-list-item :title="tt('Used storage')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.totalCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Application Code')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.codeCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Resource Files')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.assetsCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Map Data')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.mapCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('User Avatars')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.userAvatarCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Transaction Pictures')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.transactionPictureCacheSize, 2) : '-'"></f7-list-item>
            <f7-list-item :title="tt('Custom Icons')" :after="fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.customIconCacheSize, 2) : '-'"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers :class="{ 'margin-vertical': isSupportedFileCache, 'margin-vertical-half': !isSupportedFileCache }" v-if="!loading">
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
                :class="{ 'disabled': loading || !isSupportedFileCache || !fileCacheStatistics }"
                :title="tt('Transaction Pictures')"
                :after="findNameByValue(allTransactionPictureCacheExpirationOptions, transactionPictureCacheExpiration)"
                @click="showTransactionPictureCacheExpirationPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="value" value-field="value"
                                           title-field="name"
                                           :title="tt('Cache Expiration for Transaction Pictures')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Expiration Time')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTransactionPictureCacheExpirationOptions"
                                           v-model:show="showTransactionPictureCacheExpirationPopup"
                                           v-model="transactionPictureCacheExpiration">
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
                <f7-actions-button :class="{ 'disabled': !isSupportedFileCache || loading || clearingCache || !fileCacheStatistics }"
                                   @click="openClearFileCache">{{ tt('Clear File Cache') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading || clearingCache || !exchangeRatesCacheSize }"
                                   @click="clearExchangeRatesCache">{{ tt('Clear Exchange Rates Data Cache') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-sheet swipe-handler=".swipe-handler" style="height:auto"
                  :swipe-to-close="!clearingCache" :close-on-escape="!clearingCache"
                  :close-by-backdrop-click="!clearingCache" :close-by-outside-click="!clearingCache"
                  :opened="showClearFileCache" @sheet:closed="showClearFileCache = false">
            <div class="swipe-handler"></div>
            <f7-page-content class="margin-top no-padding-top" style="max-height:80vh">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ tt('Clear File Cache') }}</b></div>
                </div>
                <div class="padding-bottom padding-horizontal">
                    <f7-list class="clear-file-cache-list no-margin" dividers>
                        <f7-list-item checkbox
                                      :class="{ 'disabled': loading || clearingCache }"
                                      :disabled="loading || clearingCache"
                                      :key="item.key"
                                      :title="tt(item.name)"
                                      :after="getDisplayFileCacheSize(item)"
                                      :checked="selectedFileCacheTypes.some(selected => selected.key === item.key)"
                                      v-for="item in FileCacheType.values()"
                                      @change="setFileCacheSelected(item, $event.target.checked)">
                        </f7-list-item>
                    </f7-list>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-button large fill :class="{ 'disabled': loading || clearingCache || !selectedFileCacheTypes.length }"
                               :text="tt('Clear')" @click="clearFileCache"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': clearingCache }" @click="showClearFileCache = false" :text="tt('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import { useAppBrowserCacheSettingPageBase } from '@/views/base/settings/AppBrowserCacheSettingPageBase.ts';

import { FileCacheType } from '@/core/cache.ts';

import { findNameByValue } from '@/lib/common.ts';
import { isMapProviderUseExternalSDK } from '@/lib/map/index.ts';
import { getMapProvider, isMapDataFetchProxyEnabled } from '@/lib/server_settings.ts';
import logger from '@/lib/logger.ts';

const { tt, formatVolumeToLocalizedNumerals } = useI18n();
const { showConfirm, showToast } = useI18nUIComponents();

const {
    isSupportedFileCache,
    loading,
    clearingCache,
    fileCacheStatistics,
    exchangeRatesCacheSize,
    selectedFileCacheTypes,
    showClearFileCache,
    allMapCacheExpirationOptions,
    allTransactionPictureCacheExpirationOptions,
    allExchangeRatesDataCacheExpirationOptions,
    mapCacheExpiration,
    transactionPictureCacheExpiration,
    exchangeRatesDataCacheExpiration,
    loadCacheStatistics,
    getDisplayFileCacheSize,
    openClearFileCache,
    setFileCacheSelected,
    clearSelectedFileCaches,
    clearExchangeRatesDataCache
} = useAppBrowserCacheSettingPageBase();

const showMapDataCacheExpirationPopup = ref<boolean>(false);
const showTransactionPictureCacheExpirationPopup = ref<boolean>(false);
const showExchangeRatesDataCacheExpirationPopup = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);

function reloadCacheStatistics(done?: () => void): void {
    loadCacheStatistics(false).then(() => {
        done?.();
    }).catch(() => {
        done?.();
    });
}

function clearFileCache(): void {
    showLoading();

    clearSelectedFileCaches().then(() => {
        hideLoading();

        showToast('Cache has been cleared');
    }).catch(error => {
        hideLoading();

        logger.error('failed to clear browser cache', error);
        showToast('Unable to clear cache');
    });
}

function clearExchangeRatesCache(): void {
    showConfirm('Are you sure you want to clear exchange rates data cache?', () => {
        clearExchangeRatesDataCache();
        showToast('Cache has been cleared');
    });
}

loadCacheStatistics(true);
</script>

<style>
.clear-file-cache-list.list > ul > li > .item-content {
    padding-inline-start: calc(var(--f7-list-item-padding-horizontal) / 2);
}
</style>
