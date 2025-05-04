import { type NumberFormatOptions, DecimalSeparator, DigitGroupingSymbol, DigitGroupingType } from '@/core/numeral.ts';
import { DEFAULT_DECIMAL_NUMBER_COUNT, MAX_SUPPORTED_DECIMAL_NUMBER_COUNT } from '@/consts/numeral.ts';

import {isString, isNumber, replaceAll, removeAll } from './common.ts';

export function appendDigitGroupingSymbol(value: number | string, options: NumberFormatOptions): string {
    let textualValue = '';

    if (isNumber(value)) {
        textualValue = value.toString();
    } else {
        textualValue = value;
    }

    if (!textualValue) {
        return textualValue;
    }

    if (!options) {
        options = {};
    }

    if (!isNumber(options.digitGrouping) || options.digitGrouping === DigitGroupingType.None.type) {
        return textualValue;
    }

    if (textualValue.length <= 3) {
        return textualValue;
    }

    const negative = textualValue.charAt(0) === '-';

    if (negative) {
        textualValue = textualValue.substring(1);
    }

    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

    const integerChars = [];
    let currentDecimalSeparator = '';
    let decimals = '';

    for (let i = 0; i < textualValue.length; i++) {
        const ch = textualValue.charAt(i);

        if ('0' <= ch && ch <= '9') {
            integerChars.push(ch);
        } else {
            currentDecimalSeparator = ch;
            decimals = textualValue.substring(i + 1);
            break;
        }
    }

    let newInteger = '';

    if (options.digitGrouping === DigitGroupingType.ThousandsSeparator.type) {
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

export function appendDecimalSeparator(value: number | string, options: NumberFormatOptions): string {
    let textualValue = '';

    if (isNumber(value)) {
        textualValue = value.toString();
    } else {
        textualValue = value;
    }

    if (!textualValue) {
        return textualValue;
    }

    if (!options) {
        options = {};
    }

    if (!isString(options.decimalSeparator)) {
        return textualValue;
    }

    if (textualValue.length < 1) {
        return textualValue;
    }

    const negative = textualValue.charAt(0) === '-';

    if (negative) {
        textualValue = textualValue.substring(1);
    }

    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

    let currentDecimalSeparator = '';
    let integer = '';
    let decimals = '';

    for (let i = 0; i < textualValue.length; i++) {
        const ch = textualValue.charAt(i);

        if ('0' <= ch && ch <= '9') {
            integer += ch;
        } else {
            currentDecimalSeparator = ch;
            decimals = textualValue.substring(i + 1);
            break;
        }
    }


    if (negative) {
        integer = `-${integer}`;
    }

    if (currentDecimalSeparator) {
        return `${integer}${decimalSeparator}${decimals}`;
    } else {
        return integer;
    }
}

export function parseAmount(str: string, options: NumberFormatOptions): number {
    if (!isString(str)) {
        return 0;
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

    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;

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

export function formatAmount(value: number | string, options: NumberFormatOptions): string {
    let textualValue = '';

    if (isNumber(value)) {
        textualValue = value.toString();
    } else {
        textualValue = value;
    }

    if (!textualValue) {
        return textualValue;
    }

    if (!options) {
        options = {};
    }

    const negative = textualValue.charAt(0) === '-';

    if (negative) {
        textualValue = textualValue.substring(1);
    }

    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;
    let decimalNumberCount = options.decimalNumberCount;

    if (!isNumber(decimalNumberCount) || decimalNumberCount > MAX_SUPPORTED_DECIMAL_NUMBER_COUNT) {
        decimalNumberCount = DEFAULT_DECIMAL_NUMBER_COUNT;
    }

    let integer = '0';
    let decimals = '00';

    if (textualValue.length > 2) {
        integer = textualValue.substring(0, textualValue.length - 2);
        decimals = textualValue.substring(textualValue.length - 2);
    } else if (textualValue.length === 2) {
        decimals = textualValue;
    } else if (textualValue.length === 1) {
        decimals = '0' + textualValue;
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
        textualValue = `${integer}${decimalSeparator}${decimals}`;
    } else {
        textualValue = integer;
    }

    if (negative) {
        textualValue = `-${textualValue}`;
    }

    return textualValue;
}

export function formatNumber(value: number, precision: number, options: NumberFormatOptions): string {
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.floor(value * ratio);
    const textualValue = (normalizedValue / ratio).toString();

    return appendDecimalSeparator(textualValue, options);
}

export function formatPercent(value: number, precision: number, lowPrecisionValue: string, options: NumberFormatOptions): string {
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.floor(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        const systemDecimalSeparator = DecimalSeparator.Dot.symbol;
        const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

        if (systemDecimalSeparator === decimalSeparator) {
            return lowPrecisionValue + '%';
        }

        return replaceAll(lowPrecisionValue, systemDecimalSeparator, decimalSeparator) + '%';
    }

    return formatNumber(value, precision, options) + '%';
}

export function getAmountWithDecimalNumberCount(amount: number, decimalNumberCount: number): number {
    if (decimalNumberCount === 0) {
        return Math.floor(amount / 100) * 100;
    } else if (decimalNumberCount === 1) {
        return Math.floor(amount / 10) * 10;
    }

    return amount;
}

export function formatExchangeRateAmount(exchangeRateAmount: number | string, options: NumberFormatOptions): string {
    if (!options) {
        options = {};
    }

    const rateStr = exchangeRateAmount.toString();
    const decimalSeparator = DecimalSeparator.Dot.symbol;

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

export function getAdaptiveDisplayAmountRate(amount1: number, amount2: number, fromExchangeRate: { rate: string }, toExchangeRate: { rate: string }, options: NumberFormatOptions): string | null {
    if (!amount1 || !amount2 || amount1 === amount2) {
        if (!fromExchangeRate || !fromExchangeRate.rate || !toExchangeRate || !toExchangeRate.rate) {
            return null;
        }

        amount1 = parseFloat(fromExchangeRate.rate);
        amount2 = parseFloat(toExchangeRate.rate);
    }

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

export function getExchangedAmountByRate(amount: number, fromRate: string, toRate: string): number | null {
    const exchangeRate = parseFloat(toRate) / parseFloat(fromRate);

    if (!isNumber(exchangeRate)) {
        return null;
    }

    return amount * exchangeRate;
}
