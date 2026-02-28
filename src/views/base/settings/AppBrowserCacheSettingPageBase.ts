import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { type NameNumeralValue } from '@/core/base.ts';
import { type BrowserCacheStatistics } from '@/core/cache.ts';

import {
    loadBrowserCacheStatistics,
    updateMapCacheExpiration,
    clearMapDataCache,
    clearAllBrowserCaches
} from '@/lib/cache.ts';

export function useAppBrowserCacheSettingPageBase() {
    const { tt, formatNumberToLocalizedNumerals } = useI18n();

    const isSupportedFileCache: boolean = 'serviceWorker' in navigator && !!navigator.serviceWorker.controller;

    const settingsStore = useSettingsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const loading = ref<boolean>(true);
    const fileCacheStatistics = ref<BrowserCacheStatistics | undefined>(undefined);
    const exchangeRatesCacheSize = ref<number | undefined>(undefined);

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

    const exchangeRatesDataCacheExpiration = computed<number>({
        get: () => settingsStore.appSettings.exchangeRatesDataCacheExpiration,
        set: (value) => {
            settingsStore.setExchangeRatesDataCacheExpiration(value);
            exchangeRatesStore.removeExpiredExchangeRates(true);
            exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();
        }
    });

    function loadCacheStatistics(): Promise<void> {
        return new Promise((resolve, reject) => {
            loading.value = true;

            loadBrowserCacheStatistics().then(statistics => {
                fileCacheStatistics.value = statistics;
                exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();
                loading.value = false;
                resolve();
            }).catch(error => {
                loading.value = false;
                reject(error);
            });
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
        fileCacheStatistics,
        exchangeRatesCacheSize,
        // computed states
        allMapCacheExpirationOptions,
        allExchangeRatesDataCacheExpirationOptions,
        mapCacheExpiration,
        exchangeRatesDataCacheExpiration,
        // functions
        loadCacheStatistics,
        clearMapDataCache,
        clearAllBrowserCaches,
        clearExchangeRatesDataCache
    };
}
