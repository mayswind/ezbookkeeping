import currency from "../consts/currency.js";
import settings from "../lib/settings.js";
import utils from "../lib/utils.js";

function appendThousandsSeparator(value) {
    const finalChars = [];

    for (let i = 0; i < value.length; i++) {
        if (i % 3 === 0 && i > 0) {
            finalChars.push(',');
        }

        finalChars.push(value.charAt(value.length - 1 - i));
    }

    finalChars.reverse();
    return finalChars.join('');
}

export default function ({i18n}, value, currencyCode) {
    if (!utils.isNumber(value) && !utils.isString(value)) {
        return value;
    }

    if (utils.isNumber(value)) {
        value = value.toString();
    }

    if (value.length === 0) {
        value = '0.00';
    } else if (value.length === 1) {
        value = '0.0' + value;
    } else if (value.length === 2) {
        value = '0.' + value;
    } else {
        let integer = value.substr(0, value.length - 2);
        let decimals = value.substr(value.length - 2, 2);

        if (settings.isEnableThousandsSeparator() && integer.length > 3) {
            integer = appendThousandsSeparator(integer);
        }

        value = `${integer}.${decimals}`;
    }

    const currencyDisplayMode = settings.getCurrencyDisplayMode();

    if (currencyDisplayMode === 'symbol') {
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
    } else if (currencyDisplayMode === 'code') {
        return `${value} ${currencyCode}`;
    } else if (currencyDisplayMode === 'name') {
        const currencyName = i18n.t(`currency.${currencyCode}`);
        return `${value} ${currencyName}`;
    } else {
        return value;
    }
}
