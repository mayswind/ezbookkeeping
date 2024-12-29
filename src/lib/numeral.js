import numeralConstants from '@/consts/numeral.js';

import { isString, isNumber, removeAll } from './common.ts';

export function appendDigitGroupingSymbol(value, options) {
    if (isNumber(value)) {
        value = value.toString();
    }

    if (!isString(value)) {
        return value;
    }

    if (!options) {
        options = {};
    }

    if (!isNumber(options.digitGrouping) || options.digitGrouping === numeralConstants.allDigitGroupingType.None.type) {
        return value;
    }

    if (value.length <= 3) {
        return value;
    }

    const negative = value.charAt(0) === '-';

    if (negative) {
        value = value.substring(1);
    }

    const digitGroupingSymbol = options.digitGroupingSymbol || numeralConstants.defaultDigitGroupingSymbol.symbol;
    const decimalSeparator = options.decimalSeparator || numeralConstants.defaultDecimalSeparator.symbol;

    let integerChars = [];
    let currentDecimalSeparator = '';
    let decimals = '';

    for (let i = 0; i < value.length; i++) {
        const ch = value.charAt(i);

        if ('0' <= ch && ch <= '9') {
            integerChars.push(ch);
        } else {
            currentDecimalSeparator = ch;
            decimals = value.substring(i + 1);
            break;
        }
    }

    let newInteger = '';

    if (options.digitGrouping === numeralConstants.allDigitGroupingType.ThousandsSeparator.type) {
        for (let i = integerChars.length - 1, j = 0; i >= 0; i--, j++) {
            if (j % 3 === 0 && j > 0) {
                newInteger = digitGroupingSymbol + newInteger;
            }

            newInteger = integerChars[i] + newInteger;
        }
    }

    if (negative) {
        newInteger = `-${newInteger}`;
    }

    if (currentDecimalSeparator) {
        return `${newInteger}${decimalSeparator}${decimals}`;
    } else {
        return newInteger;
    }
}

export function parseAmount(str, options) {
    if (!isString(str)) {
        return str;
    }

    if (!options) {
        options = {};
    }

    if (!str || str.length < 1) {
        return 0;
    }

    const negative = str.charAt(0) === '-';

    if (negative) {
        str = str.substring(1);
    }

    if (!str || str.length < 1) {
        return 0;
    }

    const sign = negative ? -1 : 1;

    const decimalSeparator = options.decimalSeparator || numeralConstants.defaultDecimalSeparator.symbol;
    const digitGroupingSymbol = options.digitGroupingSymbol || numeralConstants.defaultDigitGroupingSymbol.symbol;

    if (str.indexOf(digitGroupingSymbol) >= 0) {
        str = removeAll(str, digitGroupingSymbol);
    }

    let decimalSeparatorPos = str.indexOf(decimalSeparator);

    if (decimalSeparatorPos < 0) {
        return sign * parseInt(str) * 100;
    } else if (decimalSeparatorPos === 0) {
        str = '0' + str;
        decimalSeparatorPos++;
    }

    const integer = str.substring(0, decimalSeparatorPos);
    const decimals = str.substring(decimalSeparatorPos + 1, str.length);

    if (decimals.length < 1) {
        return sign * parseInt(integer) * 100;
    } else if (decimals.length === 1) {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals) * 10;
    } else if (decimals.length === 2) {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals);
    } else {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals.substring(0, 2));
    }
}

export function formatAmount(value, options) {
    if (isNumber(value)) {
        value = value.toString();
    }

    if (!isString(value)) {
        return value;
    }

    if (!options) {
        options = {};
    }

    const negative = value.charAt(0) === '-';

    if (negative) {
        value = value.substring(1);
    }

    const decimalSeparator = options.decimalSeparator || numeralConstants.defaultDecimalSeparator.symbol;
    let decimalNumberCount = options.decimalNumberCount;

    if (!isNumber(decimalNumberCount) || decimalNumberCount > numeralConstants.maxSupportedDecimalNumberCount) {
        decimalNumberCount = numeralConstants.defaultDecimalNumberCount;
    }

    let integer = '0';
    let decimals = '00';

    if (value.length > 2) {
        integer = value.substring(0, value.length - 2);
        decimals = value.substring(value.length - 2);
    } else if (value.length === 2) {
        decimals = value;
    } else if (value.length === 1) {
        decimals = '0' + value;
    }

    if (decimalNumberCount === 0) {
        if (decimals === '00') {
            decimals = '';
        } else if (decimals.charAt(1) === '0') {
            decimals = decimals.charAt(0);
        }
    } else if (decimalNumberCount === 1) {
        if (decimals.charAt(1) === '0') {
            decimals = decimals.charAt(0);
        }
    }

    if (options.trimTailZero) {
        if (decimals.charAt(0) === '0' && decimals.charAt(1) === '0') {
            decimals = '';
        } else if (decimals.charAt(0) !== '0' && decimals.charAt(1) === '0') {
            decimals = decimals.charAt(0);
        }
    }

    integer = appendDigitGroupingSymbol(integer, options);

    if (decimals !== '') {
        value = `${integer}${decimalSeparator}${decimals}`;
    } else {
        value = integer;
    }

    if (negative) {
        value = `-${value}`;
    }

    return value;
}

export function formatPercent(value, precision, lowPrecisionValue) {
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.floor(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        return lowPrecisionValue + '%';
    }

    const result = normalizedValue / ratio;
    return result + '%';
}

export function getAmountWithDecimalNumberCount(amount, decimalNumberCount) {
    if (decimalNumberCount === 0) {
        return Math.floor(amount / 100) * 100;
    } else if (decimalNumberCount === 1) {
        return Math.floor(amount / 10) * 10;
    }

    return amount;
}

export function formatExchangeRateAmount(exchangeRateAmount, options) {
    if (!options) {
        options = {};
    }

    const rateStr = exchangeRateAmount.toString();
    const decimalSeparator = numeralConstants.allDecimalSeparator.Dot.symbol;

    if (rateStr.indexOf(decimalSeparator) < 0) {
        return appendDigitGroupingSymbol(rateStr, options);
    } else {
        let firstNonZeroPos = 0;

        for (let i = 0; i < rateStr.length; i++) {
            if (rateStr.charAt(i) !== decimalSeparator && rateStr.charAt(i) !== '0') {
                firstNonZeroPos = Math.min(i + 4, rateStr.length);
                break;
            }
        }

        const trimmedRateStr = rateStr.substring(0, Math.max(6, Math.max(firstNonZeroPos, rateStr.indexOf(decimalSeparator) + 2)));
        return appendDigitGroupingSymbol(trimmedRateStr, options);
    }
}

export function getAdaptiveDisplayAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, options) {
    if (!amount1 || !amount2 || amount1 === amount2) {
        if (!fromExchangeRate || !fromExchangeRate.rate || !toExchangeRate || !toExchangeRate.rate) {
            return null;
        }

        amount1 = fromExchangeRate.rate;
        amount2 = toExchangeRate.rate;
    }

    amount1 = parseFloat(amount1);
    amount2 = parseFloat(amount2);

    if (amount1 > amount2) {
        const rateStr = (amount1 / amount2).toString();
        const displayRateStr = formatExchangeRateAmount(rateStr, options);
        return `${displayRateStr} : 1`;
    } else {
        const rateStr = (amount2 / amount1).toString();
        const displayRateStr = formatExchangeRateAmount(rateStr, options);
        return `1 : ${displayRateStr}`;
    }
}

export function getExchangedAmount(amount, fromRate, toRate) {
    const exchangeRate = parseFloat(toRate) / parseFloat(fromRate);

    if (!isNumber(exchangeRate)) {
        return null;
    }

    return amount * exchangeRate;
}

export function getConvertedAmount(baseAmount, fromExchangeRate, toExchangeRate) {
    if (!fromExchangeRate || !toExchangeRate) {
        return '';
    }

    if (baseAmount === '') {
        return 0;
    }

    return getExchangedAmount(baseAmount, fromExchangeRate.rate, toExchangeRate.rate);
}
