import { defineStore } from 'pinia';

import services from '@/lib/services.js';
import logger from '@/lib/logger.ts';
import { isEquals } from '@/lib/common.ts';
import { getCurrentUnixTime, formatUnixTime } from '@/lib/datetime.ts';
import { getExchangedAmount } from '@/lib/numeral.ts';

const exchangeRatesLocalStorageKey = 'ebk_app_exchange_rates';

function getExchangeRatesFromLocalStorage() {
    const storageData = localStorage.getItem(exchangeRatesLocalStorageKey) || '{}';
    return JSON.parse(storageData);
}

function setExchangeRatesToLocalStorage(value) {
    const storageData = JSON.stringify(value);
    localStorage.setItem(exchangeRatesLocalStorageKey, storageData);
}

function clearExchangeRatesFromLocalStorage() {
    localStorage.removeItem(exchangeRatesLocalStorageKey);
}

export const useExchangeRatesStore = defineStore('exchangeRates', {
    state: () => ({
        latestExchangeRates: getExchangeRatesFromLocalStorage()
    }),
    getters: {
        exchangeRatesLastUpdateTime(state) {
            const exchangeRates = state.latestExchangeRates || {};
            return exchangeRates && exchangeRates.data ? exchangeRates.data.updateTime : null;
        },
        latestExchangeRateMap(state) {
            const exchangeRateMap = {};

            if (!state.latestExchangeRates || !state.latestExchangeRates.data || !state.latestExchangeRates.data.exchangeRates) {
                return exchangeRateMap;
            }

            for (let i = 0; i < state.latestExchangeRates.data.exchangeRates.length; i++) {
                const exchangeRate = state.latestExchangeRates.data.exchangeRates[i];
                exchangeRateMap[exchangeRate.currency] = exchangeRate;
            }

            return exchangeRateMap;
        }
    },
    actions: {
        resetLatestExchangeRates() {
            this.latestExchangeRates = {};
            clearExchangeRatesFromLocalStorage();
        },
        getLatestExchangeRates({ silent, force }) {
            const self = this;
            const currentExchangeRateData = self.latestExchangeRates;
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
                        reject({ message: 'Exchange rates data is up to date' });
                        return;
                    }

                    this.latestExchangeRates = {
                        time: now,
                        data: data.result
                    };
                    setExchangeRatesToLocalStorage(this.latestExchangeRates);

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
        },
        getExchangedAmount(amount, fromCurrency, toCurrency) {
            if (amount === 0) {
                return 0;
            }

            if (!this.latestExchangeRates || !this.latestExchangeRates.data || !this.latestExchangeRates.data.exchangeRates) {
                return null;
            }

            const exchangeRates = this.latestExchangeRates.data.exchangeRates;
            const exchangeRateMap = {};

            for (let i = 0; i < exchangeRates.length; i++) {
                const exchangeRate = exchangeRates[i];
                exchangeRateMap[exchangeRate.currency] = exchangeRate;
            }

            const fromCurrencyExchangeRate = exchangeRateMap[fromCurrency];
            const toCurrencyExchangeRate = exchangeRateMap[toCurrency];

            if (!fromCurrencyExchangeRate || !toCurrencyExchangeRate) {
                return null;
            }

            return getExchangedAmount(amount, fromCurrencyExchangeRate.rate, toCurrencyExchangeRate.rate);
        }
    }
});
