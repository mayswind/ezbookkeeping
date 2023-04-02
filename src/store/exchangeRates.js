import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utils from '../lib/utils.js';

import {
    STORE_LATEST_EXCHANGE_RATES
} from './mutations.js';

const exchangeRatesLocalStorageKey = 'ebk_app_exchange_rates';

export function getLatestExchangeRates(context, { silent, force }) {
    const currentExchangeRateData = context.state.latestExchangeRates;
    const now = utils.getCurrentUnixTime();

    if (!force) {
        if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data &&
            utils.formatUnixTime(currentExchangeRateData.data.updateTime, 'YYYY-MM-DD') === utils.formatUnixTime(now, 'YYYY-MM-DD')) {
            return currentExchangeRateData.data;
        }

        if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data &&
            utils.formatUnixTime(currentExchangeRateData.time, 'YYYY-MM-DD HH') === utils.formatUnixTime(now, 'YYYY-MM-DD HH')) {
            return currentExchangeRateData.data;
        }
    }

    return new Promise((resolve, reject) => {
        services.getLatestExchangeRates({
            ignoreError: silent
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get exchange rates data' });
                return;
            }

            const currentData = getExchangeRatesFromLocalStorage();

            if (currentData && currentData.data && utils.isEquals(currentData.data, data.result)) {
                reject({ message: 'Exchange rates data is up to date' });
                return;
            }

            context.commit(STORE_LATEST_EXCHANGE_RATES, {
                time: now,
                data: data.result
            });

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to get latest exchange rates data', error);

            if (error && error.processed) {
                reject(error);
            } else if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else {
                reject({ message: 'Unable to get exchange rates data' });
            }
        });
    });
}

export function exchangeRatesLastUpdateTime(state) {
    const exchangeRates = state.latestExchangeRates || {};
    return exchangeRates && exchangeRates.data ? exchangeRates.data.updateTime : null;
}

export function getExchangedAmount(state) {
    return (amount, fromCurrency, toCurrency) => {
        if (!state.latestExchangeRates || !state.latestExchangeRates.data || !state.latestExchangeRates.data.exchangeRates) {
            return null;
        }

        const exchangeRates = state.latestExchangeRates.data.exchangeRates;
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

        return utils.getExchangedAmount(amount, fromCurrencyExchangeRate.rate, toCurrencyExchangeRate.rate)
    };
}

export function getExchangeRatesFromLocalStorage() {
    const storageData = localStorage.getItem(exchangeRatesLocalStorageKey) || '{}';
    return JSON.parse(storageData);
}

export function setExchangeRatesToLocalStorage(value) {
    const storageData = JSON.stringify(value);
    localStorage.setItem(exchangeRatesLocalStorageKey, storageData);
}

export function clearExchangeRatesFromLocalStorage() {
    localStorage.removeItem(exchangeRatesLocalStorageKey);
}
