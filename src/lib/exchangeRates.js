import utils from './utils.js';

const exchangeRatesLocalStorageKey = 'lab_exchange_rates';

function getExchangeRates() {
    const storageData = localStorage.getItem(exchangeRatesLocalStorageKey) || '{}';
    return JSON.parse(storageData);
}

function setExchangeRates(value) {
    const storageData = JSON.stringify(value);
    localStorage.setItem(exchangeRatesLocalStorageKey, storageData);
}

function clearExchangeRates() {
    localStorage.removeItem(exchangeRatesLocalStorageKey);
}

function getExchangeRate(fromCurrency, toCurrency) {
    const exchangeRates = getExchangeRates().exchangeRates;
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

    return parseFloat(toCurrencyExchangeRate.rate) / parseFloat(fromCurrencyExchangeRate.rate);
}

function getOtherCurrencyAmount(amount, fromCurrency, toCurrency) {
    const exchangeRate = getExchangeRate(fromCurrency, toCurrency);

    if (!utils.isNumber(exchangeRate)) {
        return null;
    }

    return amount * exchangeRate;
}

export default {
    getExchangeRates,
    setExchangeRates,
    clearExchangeRates,
    getExchangeRate,
    getOtherCurrencyAmount,
};
