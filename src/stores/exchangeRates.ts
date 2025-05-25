import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type { BeforeResolveFunction } from '@/core/base.ts';

import type {
    UserCustomExchangeRateUpdateResponse,
    LatestExchangeRate,
    LatestExchangeRateResponse
} from '@/models/exchange_rate.ts';

import { isEquals } from '@/lib/common.ts';
import { getCurrentUnixTime, formatUnixTime } from '@/lib/datetime.ts';
import { getExchangedAmountByRate } from '@/lib/numeral.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

const exchangeRatesLocalStorageKey = 'ebk_app_exchange_rates';
const userDataSourceType = 'user_custom';

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

    const isUserCustomExchangeRates = computed((): boolean => {
        if (!latestExchangeRates.value || !latestExchangeRates.value.data) {
            return false;
        }

        return latestExchangeRates.value.data.dataSource === userDataSourceType;
    });

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

    function updateExchangeRateToLatestExchangeRateList(exchangeRate: LatestExchangeRate, updateTime: number): void {
        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return;
        }

        const exchangeRates = latestExchangeRates.value.data.exchangeRates;
        let changed = false;

        for (let i = 0; i < exchangeRates.length; i++) {
            if (exchangeRates[i].currency === exchangeRate.currency) {
                exchangeRates.splice(i, 1, exchangeRate);
                changed = true;
                break;
            }
        }

        if (!changed) {
            exchangeRates.push(exchangeRate);
            changed = true;
        }

        latestExchangeRates.value.data.updateTime = updateTime;

        if (changed) {
            setExchangeRatesToLocalStorage(latestExchangeRates.value);
        }
    }

    function removeExchangeRateFromLatestExchangeRateList(currency: string): void {
        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return;
        }

        const exchangeRates = latestExchangeRates.value.data.exchangeRates;
        let changed = false;

        for (let i = 0; i < exchangeRates.length; i++) {
            if (exchangeRates[i].currency === currency) {
                exchangeRates.splice(i, 1);
                changed = true;
                break;
            }
        }

        if (changed) {
            setExchangeRatesToLocalStorage(latestExchangeRates.value);
        }
    }

    function resetLatestExchangeRates(): void {
        latestExchangeRates.value = {};
        clearExchangeRatesFromLocalStorage();
    }

    function getLatestExchangeRates({ silent, force }: { silent: boolean, force: boolean }): Promise<LatestExchangeRateResponse> {
        const currentExchangeRateData = latestExchangeRates.value;
        const now = getCurrentUnixTime();

        if (!force) {
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
                ignoreError: silent
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve exchange rates data' });
                    return;
                }

                const currentData = getExchangeRatesFromLocalStorage();

                if (force && currentData && currentData.data && isEquals(currentData.data, data.result)) {
                    reject({ message: 'Exchange rates data is up to date', isUpToDate: true });
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

    function updateUserCustomExchangeRate({ currency, rate }: { currency: string, rate: number }): Promise<UserCustomExchangeRateUpdateResponse> {
        return new Promise((resolve, reject) => {
            services.updateUserCustomExchangeRate({
                currency: currency,
                rate: rate.toString()
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to update user custom exchange rate' });
                    return;
                }

                const exchangeRate: LatestExchangeRate = {
                    currency: data.result.currency,
                    rate: data.result.rate
                };

                updateExchangeRateToLatestExchangeRateList(exchangeRate, data.result.updateTime);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to update user custom exchange rate', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to update user custom exchange rate' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteUserCustomExchangeRate({ currency, beforeResolve }: { currency: string, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteUserCustomExchangeRate({
                currency: currency
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this user custom exchange rate' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeExchangeRateFromLatestExchangeRateList(currency);
                    });
                } else {
                    removeExchangeRateFromLatestExchangeRateList(currency);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete user custom exchange rate', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this user custom exchange rate' });
                } else {
                    reject(error);
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
        isUserCustomExchangeRates,
        exchangeRatesLastUpdateTime,
        latestExchangeRateMap,
        // functions
        resetLatestExchangeRates,
        getLatestExchangeRates,
        updateUserCustomExchangeRate,
        deleteUserCustomExchangeRate,
        getExchangedAmount
    };
});
