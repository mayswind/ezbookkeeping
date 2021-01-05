import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utils from '../lib/utils.js';

import {
    STORE_LATEST_EXCHANGE_RATES
} from './mutations.js';

const exchangeRatesLocalStorageKey = 'lab_app_exchange_rates';

function getLatestExchangeRates(context, { silent, force }) {
    const currentExchangeRateData = context.state.latestExchangeRates;
    const now = new Date();

    if (!force) {
        if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data &&
            currentExchangeRateData.data.date === utils.formatDate(now, 'YYYY-MM-DD')) {
            return currentExchangeRateData.data;
        }

        if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data &&
            utils.formatUnixTime(currentExchangeRateData.time, 'YYYY-MM-DD HH') === utils.formatDate(now, 'YYYY-MM-DD HH')) {
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

            context.commit(STORE_LATEST_EXCHANGE_RATES, {
                time: utils.getUnixTime(now),
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

function exchangeRatesLastUpdateDate(state) {
    const exchangeRates = state.latestExchangeRates || {};
    return exchangeRates && exchangeRates.data ? exchangeRates.data.date : null;
}

function getExchangedAmount(state) {
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

        const exchangeRate = parseFloat(toCurrencyExchangeRate.rate) / parseFloat(fromCurrencyExchangeRate.rate);

        if (!utils.isNumber(exchangeRate)) {
            return null;
        }

        return amount * exchangeRate;
    };
}

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

export default {
    getLatestExchangeRates,
    exchangeRatesLastUpdateDate,
    getExchangedAmount,
    getExchangeRatesFromLocalStorage,
    setExchangeRatesToLocalStorage,
    clearExchangeRatesFromLocalStorage,
}
