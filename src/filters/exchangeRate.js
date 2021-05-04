export default function (oldRate, currentCurrency, allExchangeRates) {
    const exchangeRateMap = {};

    for (let i = 0; i < allExchangeRates.length; i++) {
        const exchangeRate = allExchangeRates[i];
        exchangeRateMap[exchangeRate.currency] = exchangeRate;
    }

    const toCurrencyExchangeRate = exchangeRateMap[currentCurrency];

    if (!toCurrencyExchangeRate) {
        return '';
    }

    const newRate = parseFloat(oldRate) / parseFloat(toCurrencyExchangeRate.rate);
    const newRateStr = newRate.toString();

    if (newRateStr.indexOf('.') < 0) {
        return newRateStr;
    } else {
        let firstNonZeroPos = 0;

        for (let i = 0; i < newRateStr.length; i++) {
            if (newRateStr.charAt(i) !== '.' && newRateStr.charAt(i) !== '0') {
                firstNonZeroPos = Math.min(i + 4, newRateStr.length);
                break;
            }
        }

        return newRateStr.substr(0, Math.max(6, Math.max(firstNonZeroPos, newRateStr.indexOf('.') + 2)));
    }
}
