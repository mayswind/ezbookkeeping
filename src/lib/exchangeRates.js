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

export default {
    getExchangeRates,
    setExchangeRates,
    clearExchangeRates,
};
