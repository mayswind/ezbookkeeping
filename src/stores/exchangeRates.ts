import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type { LatestExchangeRate, LatestExchangeRateResponse } from '@/models/exchange_rate.ts';

import { isEquals } from '@/lib/common.ts';
import { getCurrentUnixTime, formatUnixTime } from '@/lib/datetime.ts';
import { getExchangedAmountByRate } from '@/lib/numeral.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

const exchangeRatesLocalStorageKey = 'ebk_app_exchange_rates';

interface LatestExchangeRates {
    readonly time?: number;
    readonly data?: LatestExchangeRateResponse;
}

function getExchangeRatesFromLocalStorage(): LatestExchangeRates {
    const storageData = localStorage.getItem(exchangeRatesLocalStorageKey) || '{}';
    return JSON.parse(storageData) as LatestExchangeRates;
}

function setExchangeRatesToLocalStorage(value: LatestExchangeRates): void {
    const storageData = JSON.stringify(value);
    localStorage.setItem(exchangeRatesLocalStorageKey, storageData);
}

function clearExchangeRatesFromLocalStorage(): void {
    localStorage.removeItem(exchangeRatesLocalStorageKey);
}

export const useExchangeRatesStore = defineStore('exchangeRates', () => {
    const latestExchangeRates = ref<LatestExchangeRates>(getExchangeRatesFromLocalStorage());

    const exchangeRatesLastUpdateTime = computed<number | null>(() => {
        const exchangeRates = latestExchangeRates.value || {};
        return exchangeRates && exchangeRates.data ? exchangeRates.data.updateTime : null;
    });

    const latestExchangeRateMap = computed<Record<string, LatestExchangeRate>>(() => {
        const exchangeRateMap: Record<string, LatestExchangeRate> = {};

        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return exchangeRateMap;
        }

        for (let i = 0; i < latestExchangeRates.value.data.exchangeRates.length; i++) {
            const exchangeRate = latestExchangeRates.value.data.exchangeRates[i];
            exchangeRateMap[exchangeRate.currency] = exchangeRate;
        }

        return exchangeRateMap;
    });

    function resetLatestExchangeRates(): void {
        latestExchangeRates.value = {};
        clearExchangeRatesFromLocalStorage();
    }

    function getLatestExchangeRates(req: { silent: boolean, force: boolean }): Promise<LatestExchangeRateResponse> {
        const currentExchangeRateData = latestExchangeRates.value;
        const now = getCurrentUnixTime();

        if (!req.force) {
            if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data &&
                formatUnixTime(currentExchangeRateData.data.updateTime, 'YYYY-MM-DD') === formatUnixTime(now, 'YYYY-MM-DD')) {
                return Promise.resolve(currentExchangeRateData.data);
            }

            if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data &&
                formatUnixTime(currentExchangeRateData.time, 'YYYY-MM-DD HH') === formatUnixTime(now, 'YYYY-MM-DD HH')) {
                return Promise.resolve(currentExchangeRateData.data);
            }
        }

        return new Promise((resolve, reject) => {
            services.getLatestExchangeRates({
                ignoreError: req.silent
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve exchange rates data' });
                    return;
                }

                const currentData = getExchangeRatesFromLocalStorage();

                if (req.force && currentData && currentData.data && isEquals(currentData.data, data.result)) {
                    reject({ message: 'Exchange rates data is up to date' });
                    return;
                }

                latestExchangeRates.value = {
                    time: now,
                    data: data.result
                };
                setExchangeRatesToLocalStorage(latestExchangeRates.value);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve latest exchange rates data', error);

                if (error && error.processed) {
                    reject(error);
                } else if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to retrieve exchange rates data' });
                }
            });
        });
    }

    function getExchangedAmount(amount: number, fromCurrency: string, toCurrency: string): number | null {
        if (amount === 0) {
            return 0;
        }

        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return null;
        }

        const exchangeRates = latestExchangeRates.value.data.exchangeRates;
        const exchangeRateMap: Record<string, LatestExchangeRate> = {};

        for (let i = 0; i < exchangeRates.length; i++) {
            const exchangeRate = exchangeRates[i];
            exchangeRateMap[exchangeRate.currency] = exchangeRate;
        }

        const fromCurrencyExchangeRate = exchangeRateMap[fromCurrency];
        const toCurrencyExchangeRate = exchangeRateMap[toCurrency];

        if (!fromCurrencyExchangeRate || !toCurrencyExchangeRate) {
            return null;
        }

        return getExchangedAmountByRate(amount, fromCurrencyExchangeRate.rate, toCurrencyExchangeRate.rate);
    }

    return {
        // states
        latestExchangeRates,
        // computed states
        exchangeRatesLastUpdateTime,
        latestExchangeRateMap,
        // functions
        resetLatestExchangeRates,
        getLatestExchangeRates,
        getExchangedAmount
    };
});
