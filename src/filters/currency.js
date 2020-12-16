import currency from "../consts/currency.js";
import settings from "../lib/settings.js";
import utils from "../lib/utils.js";

export default function ({i18n}, value, currencyCode) {
    if (!utils.isNumber(value) && !utils.isString(value)) {
        return value;
    }

    if (utils.isNumber(value)) {
        value = value.toString();
    }

    value = utils.numericCurrencyToString(value);

    const currencyDisplayMode = settings.getCurrencyDisplayMode();

    if (currencyCode && currencyDisplayMode === 'symbol') {
        const currencyInfo = currency.all[currencyCode];
        let currencySymbol = currency.defaultCurrencySymbol;

        if (currencyInfo && currencyInfo.symbol) {
            currencySymbol = currencyInfo.symbol;
        } else if (currencyInfo && currencyInfo.code) {
            currencySymbol = currencyInfo.code;
        }

        return i18n.t('format.currency.symbol', {
            amount: value,
            symbol: currencySymbol
        });
    } else if (currencyCode && currencyDisplayMode === 'code') {
        return `${value} ${currencyCode}`;
    } else if (currencyCode && currencyDisplayMode === 'name') {
        const currencyName = i18n.t(`currency.${currencyCode}`);
        return `${value} ${currencyName}`;
    } else {
        return value;
    }
}
