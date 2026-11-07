import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { type NameNumeralValue } from '@/core/base.ts';
import { type BrowserCacheStatistics, FileCacheType } from '@/core/cache.ts';

import {
    loadBrowserCacheStatistics,
    updateMapCacheExpiration,
    updateTransactionPictureCacheExpiration,
    clearFileCaches
} from '@/lib/cache.ts';

export function useAppBrowserCacheSettingPageBase() {
    const {
        tt,
        formatNumberToLocalizedNumerals,
        formatVolumeToLocalizedNumerals
    } = useI18n();

    const isSupportedFileCache: boolean = 'serviceWorker' in navigator && !!navigator.serviceWorker.controller;

    const settingsStore = useSettingsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const loading = ref<boolean>(true);
    const clearingCache = ref<boolean>(false);
    const fileCacheStatistics = ref<BrowserCacheStatistics | undefined>(undefined);
    const exchangeRatesCacheSize = ref<number | undefined>(undefined);
    const selectedFileCacheTypes = ref<FileCacheType[]>([]);
    const showClearFileCache = ref<boolean>(false);

    const allMapCacheExpirationOptions = computed<NameNumeralValue[]>(() => {
        return [
            { name: tt('Disable Cache'), value: -1 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(1) }), value: 86400 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(7) }), value: 604800 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(30) }), value: 2592000 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(90) }), value: 7776000 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(180) }), value: 15552000 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(365) }), value: 31536000 },
            { name: tt('No Expiration'), value: 0 }
        ];
    });

    const allTransactionPictureCacheExpirationOptions = allMapCacheExpirationOptions;

    const allExchangeRatesDataCacheExpirationOptions = computed<NameNumeralValue[]>(() => {
        return [
            { name: tt('Disable Cache'), value: -1 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(1) }), value: 86400 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(7) }), value: 604800 },
            { name: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(30) }), value: 2592000 },
            { name: tt('No Expiration'), value: 0 }
        ];
    });

    const mapCacheExpiration = computed<number>({
        get: () => isSupportedFileCache ? settingsStore.appSettings.mapCacheExpiration : -1,
        set: (value) => {
            settingsStore.setMapCacheExpiration(value);

            if (isSupportedFileCache) {
                updateMapCacheExpiration(value);
            }
        }
    });

    const transactionPictureCacheExpiration = computed<number>({
        get: () => isSupportedFileCache ? settingsStore.appSettings.transactionPictureCacheExpiration : -1,
        set: (value) => {
            settingsStore.setTransactionPictureCacheExpiration(value);

            if (isSupportedFileCache) {
                updateTransactionPictureCacheExpiration(value);
            }
        }
    });

    const exchangeRatesDataCacheExpiration = computed<number>({
        get: () => settingsStore.appSettings.exchangeRatesDataCacheExpiration,
        set: (value) => {
            settingsStore.setExchangeRatesDataCacheExpiration(value);
            exchangeRatesStore.removeExpiredExchangeRates(true);
            exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();
        }
    });

    function loadCacheStatistics(updateLoadingState: boolean): Promise<void> {
        return new Promise((resolve, reject) => {
            if (updateLoadingState) {
                loading.value = true;
            }

            exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();

            if (!isSupportedFileCache) {
                if (updateLoadingState) {
                    loading.value = false;
                }

                resolve();
                return;
            }

            loadBrowserCacheStatistics().then(statistics => {
                fileCacheStatistics.value = statistics;

                if (updateLoadingState) {
                    loading.value = false;
                }

                resolve();
            }).catch(error => {
                if (updateLoadingState) {
                    loading.value = false;
                }

                reject(error);
            });
        });
    }

    function getDisplayFileCacheSize(cacheType: FileCacheType): string {
        const statistics = fileCacheStatistics.value;

        if (!statistics) {
            return formatVolumeToLocalizedNumerals(0, 2);
        }

        switch (cacheType.key) {
            case FileCacheType.ApplicationCode.key:
                return formatVolumeToLocalizedNumerals(statistics.codeCacheSize, 2);
            case FileCacheType.ResourceFiles.key:
                return formatVolumeToLocalizedNumerals(statistics.assetsCacheSize, 2);
            case FileCacheType.MapData.key:
                return formatVolumeToLocalizedNumerals(statistics.mapCacheSize, 2);
            case FileCacheType.UserAvatars.key:
                return formatVolumeToLocalizedNumerals(statistics.userAvatarCacheSize, 2);
            case FileCacheType.TransactionPictures.key:
                return formatVolumeToLocalizedNumerals(statistics.transactionPictureCacheSize, 2);
            case FileCacheType.CustomIcons.key:
                return formatVolumeToLocalizedNumerals(statistics.customIconCacheSize, 2);
            default:
                return formatVolumeToLocalizedNumerals(0, 2);
        }
    }

    function openClearFileCache(): void {
        if (!isSupportedFileCache || loading.value || clearingCache.value || !fileCacheStatistics.value) {
            return;
        }

        selectedFileCacheTypes.value = [];
        showClearFileCache.value = true;
    }

    function setFileCacheSelected(cacheType: FileCacheType, selected: boolean): void {
        if (selected) {
            if (!selectedFileCacheTypes.value.some(item => item.key === cacheType.key)) {
                selectedFileCacheTypes.value.push(cacheType);
            }
        } else {
            selectedFileCacheTypes.value = selectedFileCacheTypes.value.filter(item => item.key !== cacheType.key);
        }
    }

    function clearSelectedFileCaches(): Promise<void> {
        if (!isSupportedFileCache || loading.value || clearingCache.value || !selectedFileCacheTypes.value.length) {
            return Promise.resolve();
        }

        clearingCache.value = true;

        return clearFileCaches(...selectedFileCacheTypes.value).then(() => {
            showClearFileCache.value = false;
            return loadCacheStatistics(true);
        }).finally(() => {
            clearingCache.value = false;
        });
    }

    function clearExchangeRatesDataCache(): void {
        exchangeRatesStore.resetLatestExchangeRates();
        exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();
    }

    return {
        // constants
        isSupportedFileCache,
        // states
        loading,
        clearingCache,
        fileCacheStatistics,
        exchangeRatesCacheSize,
        selectedFileCacheTypes,
        showClearFileCache,
        // computed states
        allMapCacheExpirationOptions,
        allTransactionPictureCacheExpirationOptions,
        allExchangeRatesDataCacheExpirationOptions,
        mapCacheExpiration,
        transactionPictureCacheExpiration,
        exchangeRatesDataCacheExpiration,
        // functions
        loadCacheStatistics,
        getDisplayFileCacheSize,
        openClearFileCache,
        setFileCacheSelected,
        clearSelectedFileCaches,
        clearExchangeRatesDataCache
    };
}
