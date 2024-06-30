import currencyConstants from '@/consts/currency.js';

import { isString, isNumber } from './common.js';

export function appendCurrencySymbol(value, currencyDisplayType, currencyCode, currencyName) {
    if (!currencyDisplayType) {
        return value;
    }

    if (isNumber(value)) {
        value = value.toString();
    }

    if (!isString(value)) {
        return value;
    }

    let symbol = '';
    let separator = currencyDisplayType.separator || '';

    if (currencyDisplayType.symbol === currencyConstants.allCurrencyDisplaySymbol.Symbol) {
        const currencyInfo = currencyConstants.all[currencyCode];

        if (currencyInfo && currencyInfo.symbol) {
            symbol = currencyInfo.symbol;
        } else if (currencyInfo && currencyInfo.code) {
            symbol = currencyInfo.code;
        }

        if (!symbol) {
            symbol = currencyConstants.defaultCurrencySymbol;
        }
    } else if (currencyDisplayType.symbol === currencyConstants.allCurrencyDisplaySymbol.Code) {
        symbol = currencyCode;
    }else if (currencyDisplayType.symbol === currencyConstants.allCurrencyDisplaySymbol.Name) {
        symbol = currencyName;
    }

    if (currencyDisplayType.location === currencyConstants.allCurrencyDisplayLocation.BeforeAmount) {
        return `${symbol}${separator}${value}`;
    } else if (currencyDisplayType.location === currencyConstants.allCurrencyDisplayLocation.AfterAmount) {
        return `${value}${separator}${symbol}`;
    } else {
        return value;
    }
}
