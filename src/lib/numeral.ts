import { Decimal } from 'decimal.js';

import {
    type BigDecimal,
    type HiddenAmount,
    type NumberFormatOptions,
    NumeralSystem,
    DecimalSeparator,
    DigitGroupingSymbol
} from '@/core/numeral.ts';

import { AMOUNT_FACTOR } from '@/consts/numeral.ts';

import { DEFAULT_DECIMAL_NUMBER_COUNT, MAX_SUPPORTED_DECIMAL_NUMBER_COUNT, DISPLAY_HIDDEN_AMOUNT } from '@/consts/numeral.ts';

import { isDefined, isString, isNumber, replaceAll, removeAll } from './common.ts';
import logger from './logger.ts';

class DecimalJSBigDecimal implements BigDecimal {
    public static readonly ZERO = new DecimalJSBigDecimal(0);
    public static readonly ONE = new DecimalJSBigDecimal(1);
    public static readonly NEGATIVE_ONE = new DecimalJSBigDecimal(-1);
    public static readonly NaN = new DecimalJSBigDecimal(Number.NaN);
    public static readonly POSITIVE_INFINITY = new DecimalJSBigDecimal(Number.POSITIVE_INFINITY);
    public static readonly NEGATIVE_INFINITY = new DecimalJSBigDecimal(Number.NEGATIVE_INFINITY);

    private readonly value: Decimal;

    private constructor(value: string | number | Decimal) {
        this.value = new Decimal(value);
    }

    public isZero(): boolean {
        return this.value.isZero();
    }

    public isFinite(): boolean {
        return this.value.isFinite();
    }

    public isPositive(): boolean {
        return this.value.greaterThan(0);
    }

    public isNegative(): boolean {
        return this.value.lessThan(0);
    }

    public isPositiveOrZero(): boolean {
        return this.value.greaterThanOrEqualTo(0);
    }

    public isNegativeOrZero(): boolean {
        return this.value.lessThanOrEqualTo(0);
    }

    public isPositiveInfinity(): boolean {
        return this.value.isPositive() && !this.value.isFinite();
    }

    public isNegativeInfinity(): boolean {
        return this.value.isNegative() && !this.value.isFinite();
    }

    public isNaN(): boolean {
        return this.value.isNaN();
    }

    public equals(other: BigDecimal | number | undefined): boolean {
        if (other === undefined) {
            return false;
        } else if (isBigDecimal(other)) {
            return this.value.equals((other as DecimalJSBigDecimal).value);
        } else {
            return this.value.equals(other);
        }
    }

    public notEquals(other: BigDecimal | number | undefined): boolean {
        if (other === undefined) {
            return true;
        } else if (isBigDecimal(other)) {
            return !this.value.equals((other as DecimalJSBigDecimal).value);
        } else {
            return !this.value.equals(other);
        }
    }

    public compareTo(other: BigDecimal | number): number {
        if (isBigDecimal(other)) {
            return this.value.comparedTo((other as DecimalJSBigDecimal).value);
        } else {
            return this.value.comparedTo(other);
        }
    }

    public greaterThan(other: BigDecimal | number): boolean {
        if (isBigDecimal(other)) {
            return this.value.greaterThan((other as DecimalJSBigDecimal).value);
        } else {
            return this.value.greaterThan(other);
        }
    }

    public greaterThanOrEqual(other: BigDecimal | number): boolean {
        if (isBigDecimal(other)) {
            return this.value.greaterThanOrEqualTo((other as DecimalJSBigDecimal).value);
        } else {
            return this.value.greaterThanOrEqualTo(other);
        }
    }

    public lessThan(other: BigDecimal | number): boolean {
        if (isBigDecimal(other)) {
            return this.value.lessThan((other as DecimalJSBigDecimal).value);
        } else {
            return this.value.lessThan(other);
        }
    }

    public lessThanOrEqual(other: BigDecimal | number): boolean {
        if (isBigDecimal(other)) {
            return this.value.lessThanOrEqualTo((other as DecimalJSBigDecimal).value);
        } else {
            return this.value.lessThanOrEqualTo(other);
        }
    }

    public between(min: BigDecimal | number, max: BigDecimal | number): boolean {
        const minValue = isBigDecimal(min) ? (min as DecimalJSBigDecimal).value : min;
        const maxValue = isBigDecimal(max) ? (max as DecimalJSBigDecimal).value : max;
        return this.value.greaterThanOrEqualTo(minValue) && this.value.lessThanOrEqualTo(maxValue);
    }

    public add(other: BigDecimal | number): BigDecimal {
        if (isBigDecimal(other)) {
            return new DecimalJSBigDecimal(this.value.add((other as DecimalJSBigDecimal).value));
        } else {
            return new DecimalJSBigDecimal(this.value.add(other));
        }
    }

    public subtract(other: BigDecimal | number): BigDecimal {
        if (isBigDecimal(other)) {
            return new DecimalJSBigDecimal(this.value.sub((other as DecimalJSBigDecimal).value));
        } else {
            return new DecimalJSBigDecimal(this.value.sub(other));
        }
    }

    public multiply(other: BigDecimal | number): BigDecimal {
        if (isBigDecimal(other)) {
            return new DecimalJSBigDecimal(this.value.times((other as DecimalJSBigDecimal).value));
        } else {
            return new DecimalJSBigDecimal(this.value.times(other));
        }
    }

    public divide(other: BigDecimal | number): BigDecimal {
        if (isBigDecimal(other)) {
            return new DecimalJSBigDecimal(this.value.dividedBy((other as DecimalJSBigDecimal).value));
        } else {
            return new DecimalJSBigDecimal(this.value.dividedBy(other));
        }
    }

    public negate(): BigDecimal {
        return new DecimalJSBigDecimal(this.value.negated());
    }

    public sign(): BigDecimal {
        if (this.value.isPositive()) {
            return DecimalJSBigDecimal.ONE;
        } else if (this.value.isNegative()) {
            return DecimalJSBigDecimal.NEGATIVE_ONE;
        } else {
            return DecimalJSBigDecimal.ZERO;
        }
    }

    public pow(exponent: BigDecimal | number): BigDecimal {
        if (isBigDecimal(exponent)) {
            return new DecimalJSBigDecimal(this.value.pow((exponent as DecimalJSBigDecimal).value));
        } else {
            return new DecimalJSBigDecimal(this.value.pow(exponent));
        }
    }

    public sqrt(): BigDecimal {
        return new DecimalJSBigDecimal(this.value.sqrt());
    }

    public log(): BigDecimal {
        return new DecimalJSBigDecimal(this.value.ln());
    }

    public exp(): BigDecimal {
        return new DecimalJSBigDecimal(this.value.exp());
    }

    public abs(): BigDecimal {
        return new DecimalJSBigDecimal(this.value.abs());
    }

    public truncate(): BigDecimal {
        return new DecimalJSBigDecimal(this.value.trunc());
    }

    public toSafeIntegerNumber(): number {
        const num = this.value.toNumber();

        if (!Number.isSafeInteger(num)) {
            throw new Error(`big decimal value \"${num}\" cannot be converted to a safe integer number`);
        }

        return num;
    }

    public toDoubleNumber(): number {
        return this.value.toNumber();
    }

    public toString(): string {
        return this.value.toFixed();
    }

    public static of(value: number): BigDecimal {
        return new DecimalJSBigDecimal(value);
    }

    public static parse(value: string): BigDecimal {
        return new DecimalJSBigDecimal(value);
    }
}

export const BIG_DECIMAL_ZERO: BigDecimal = DecimalJSBigDecimal.ZERO;
export const BIG_DECIMAL_POSITIVE_INFINITY: BigDecimal = DecimalJSBigDecimal.POSITIVE_INFINITY;
export const BIG_DECIMAL_NEGATIVE_INFINITY: BigDecimal = DecimalJSBigDecimal.NEGATIVE_INFINITY;

export function parseBigDecimal(value: string | number): BigDecimal {
    if (typeof value === 'number') {
        return DecimalJSBigDecimal.of(value);
    } else if (typeof value === 'string') {
        return DecimalJSBigDecimal.parse(value);
    } else {
        return DecimalJSBigDecimal.NaN;
    }
}

export function isBigDecimal(val: unknown): val is BigDecimal {
    return val instanceof DecimalJSBigDecimal;
}

export function appendDigitGroupingSymbolAndDecimalSeparator(textualNumber: string, options: NumberFormatOptions): string {
    if (!textualNumber) {
        return textualNumber;
    }

    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const digitGroupingType = options.digitGrouping;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

    const negative = textualNumber.charAt(0) === '-';

    if (negative) {
        textualNumber = textualNumber.substring(1);
    }

    const integerChars: string[] = [];
    const decimalChars: string[] = [];
    let currentDecimalSeparator = '';

    if (textualNumber === DISPLAY_HIDDEN_AMOUNT) {
        for (let i = 0; i < textualNumber.length - 2; i++) {
            integerChars.push(textualNumber.charAt(i));
        }

        const decimalStartIndex = Math.max(0, textualNumber.length - 2);

        for (let i = decimalStartIndex; i < textualNumber.length; i++) {
            decimalChars.push(textualNumber.charAt(i));
        }
    } else {
        for (let i = 0; i < textualNumber.length; i++) {
            const ch = textualNumber.charAt(i);

            if (!currentDecimalSeparator) {
                if (numeralSystem.isDigit(ch)) {
                    integerChars.push(ch);
                } else {
                    currentDecimalSeparator = ch;
                }
            } else {
                if (numeralSystem.isDigit(ch)) {
                    decimalChars.push(ch);
                } else {
                    throw new Error('Number \"' + textualNumber + '\" is not a valid textual number');
                }
            }
        }
    }

    let integer = '';

    if (digitGroupingType) {
        integer = digitGroupingType.format(integerChars, digitGroupingSymbol);
    } else {
        integer = integerChars.join('');
    }

    const decimals = decimalChars.join('');

    if (decimals) {
        textualNumber = `${integer}${decimalSeparator}${decimals}`;
    } else {
        textualNumber = integer;
    }

    if (negative) {
        textualNumber = `-${textualNumber}`;
    }

    return textualNumber;
}

export function parseAmount(str: string, options: NumberFormatOptions): number {
    if (!isString(str)) {
        return 0;
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

    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;

    if (str.indexOf(digitGroupingSymbol) >= 0) {
        str = removeAll(str, digitGroupingSymbol);
    }

    let decimalSeparatorPos = str.indexOf(decimalSeparator);

    if (decimalSeparatorPos < 0) {
        return sign * numeralSystem.parseInt(str) * AMOUNT_FACTOR;
    } else if (decimalSeparatorPos === 0) {
        str = numeralSystem.digitZero + str;
        decimalSeparatorPos++;
    }

    const integer = str.substring(0, decimalSeparatorPos);
    const decimals = str.substring(decimalSeparatorPos + 1, str.length);

    if (decimals.length < 1) {
        return sign * numeralSystem.parseInt(integer) * AMOUNT_FACTOR;
    } else if (decimals.length === 1) {
        return sign * numeralSystem.parseInt(integer) * AMOUNT_FACTOR + sign * numeralSystem.parseInt(decimals) * AMOUNT_FACTOR / 10;
    } else if (decimals.length === 2) {
        return sign * numeralSystem.parseInt(integer) * AMOUNT_FACTOR + sign * numeralSystem.parseInt(decimals);
    } else {
        return sign * numeralSystem.parseInt(integer) * AMOUNT_FACTOR + sign * numeralSystem.parseInt(decimals.substring(0, 2));
    }
}

export function formatAmount(value: BigDecimal, options: NumberFormatOptions): string {
    if (!value) {
        throw new Error('big decimal \"' + value + '\" is not valid');
    }

    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    let textualNumber = numeralSystem.formatBigDecimal(value.truncate());

    if (!textualNumber) {
        return textualNumber;
    }

    const negative = textualNumber.charAt(0) === '-';

    if (negative) {
        textualNumber = textualNumber.substring(1);
    }

    const digitGroupingType = options.digitGrouping;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;
    let decimalNumberCount = options.decimalNumberCount;

    if (!isNumber(decimalNumberCount) || decimalNumberCount > MAX_SUPPORTED_DECIMAL_NUMBER_COUNT) {
        decimalNumberCount = DEFAULT_DECIMAL_NUMBER_COUNT;
    }

    let integer = numeralSystem.digitZero;
    let decimals = numeralSystem.doubleDigitZero;

    if (textualNumber.length > 2) {
        integer = textualNumber.substring(0, textualNumber.length - 2);
        decimals = textualNumber.substring(textualNumber.length - 2);
    } else if (textualNumber.length === 2) {
        decimals = textualNumber;
    } else if (textualNumber.length === 1) {
        decimals = numeralSystem.digitZero + textualNumber;
    }

    if (decimalNumberCount === 0) {
        if (decimals === numeralSystem.doubleDigitZero) {
            decimals = '';
        } else if (decimals.charAt(1) === numeralSystem.digitZero) {
            decimals = decimals.charAt(0);
        }
    } else if (decimalNumberCount === 1) {
        if (decimals.charAt(1) === numeralSystem.digitZero) {
            decimals = decimals.charAt(0);
        }
    }

    if (options.trimTailZero) {
        if (decimals.charAt(0) === numeralSystem.digitZero && decimals.charAt(1) === numeralSystem.digitZero) {
            decimals = '';
        } else if (decimals.charAt(0) !== numeralSystem.digitZero && decimals.charAt(1) === numeralSystem.digitZero) {
            decimals = decimals.charAt(0);
        }
    }

    if (integer && integer.length > 1 && digitGroupingType) {
        integer = digitGroupingType.format(integer.split(''), digitGroupingSymbol);
    }

    if (decimals) {
        textualNumber = `${integer}${decimalSeparator}${decimals}`;
    } else {
        textualNumber = integer;
    }

    if (negative) {
        textualNumber = `-${textualNumber}`;
    }

    return textualNumber;
}

export function formatHiddenAmount(value: HiddenAmount, options: NumberFormatOptions): string {
    return appendDigitGroupingSymbolAndDecimalSeparator(value, options);
}

export function formatNumber(value: number, options: NumberFormatOptions, precision?: number): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;

    if (isDefined(precision)) {
        const ratio = Math.pow(10, precision);
        const normalizedValue = Math.trunc(value * ratio);
        const textualValue = numeralSystem.formatNumber(normalizedValue / ratio);
        return appendDigitGroupingSymbolAndDecimalSeparator(textualValue, options);
    } else {
        const textualValue = numeralSystem.formatNumber(value);
        return appendDigitGroupingSymbolAndDecimalSeparator(textualValue, options);
    }
}

export function formatBigDecimal(value: BigDecimal, options: NumberFormatOptions, precision?: number): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;

    if (isDefined(precision)) {
        const ratio = Math.pow(10, precision);
        const normalizedValue = value.multiply(ratio).truncate();
        const textualValue = numeralSystem.formatBigDecimal(normalizedValue.divide(ratio));
        return appendDigitGroupingSymbolAndDecimalSeparator(textualValue, options);
    } else {
        const textualValue = numeralSystem.formatBigDecimal(value);
        return appendDigitGroupingSymbolAndDecimalSeparator(textualValue, options);
    }
}

export function formatPercent(value: number, precision: number, lowPrecisionValue: string, options: NumberFormatOptions): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.trunc(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        const systemDecimalSeparator = DecimalSeparator.Dot.symbol;
        const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

        lowPrecisionValue = numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(lowPrecisionValue);

        if (systemDecimalSeparator === decimalSeparator) {
            return lowPrecisionValue + '%';
        }

        return replaceAll(lowPrecisionValue, systemDecimalSeparator, decimalSeparator) + '%';
    }

    return formatNumber(value, options, precision) + '%';
}

export function getAmountWithDecimalNumberCount(amount: BigDecimal, decimalNumberCount: number): BigDecimal {
    if (decimalNumberCount === 0) {
        return amount.divide(AMOUNT_FACTOR).truncate().multiply(AMOUNT_FACTOR);
    } else if (decimalNumberCount === 1) {
        const factor = AMOUNT_FACTOR / 10;
        return amount.divide(factor).truncate().multiply(factor);
    }

    return amount;
}

export function formatExchangeRateAmount(exchangeRateAmount: BigDecimal, options: NumberFormatOptions): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const rateStr = numeralSystem.formatBigDecimal(exchangeRateAmount);
    const decimalSeparator = DecimalSeparator.Dot.symbol;

    if (rateStr.indexOf(decimalSeparator) < 0) {
        return appendDigitGroupingSymbolAndDecimalSeparator(rateStr, options);
    } else {
        let firstNonZeroPos = 0;

        for (let i = 0; i < rateStr.length; i++) {
            if (rateStr.charAt(i) !== decimalSeparator && rateStr.charAt(i) !== numeralSystem.digitZero) {
                firstNonZeroPos = Math.min(i + 4, rateStr.length);
                break;
            }
        }

        const trimmedRateStr = rateStr.substring(0, Math.max(6, Math.max(firstNonZeroPos, rateStr.indexOf(decimalSeparator) + 2)));
        return appendDigitGroupingSymbolAndDecimalSeparator(trimmedRateStr, options);
    }
}

export function getAdaptiveDisplayAmountRate(amount1: number, amount2: number, options: NumberFormatOptions, fromExchangeRate?: { rate: string }, toExchangeRate?: { rate: string }): string | null {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    let finalAmount1: BigDecimal = parseBigDecimal(amount1);
    let finalAmount2: BigDecimal = parseBigDecimal(amount2);

    if (!amount1 || !amount2 || amount1 === amount2) {
        if (!fromExchangeRate || !fromExchangeRate.rate || !toExchangeRate || !toExchangeRate.rate) {
            return null;
        }

        try {
            finalAmount1 = parseBigDecimal(fromExchangeRate.rate);
            finalAmount2 = parseBigDecimal(toExchangeRate.rate);
        } catch (ex) {
            logger.warn(`failed to parse exchange rate: ${ex}`);
            return null;
        }
    }

    if (finalAmount1.greaterThan(finalAmount2)) {
        const rate: BigDecimal = finalAmount1.divide(finalAmount2);
        const displayRateStr = formatExchangeRateAmount(rate, options);
        return `${displayRateStr} : ${numeralSystem.getLocalizedDigit(1)}`;
    } else {
        const rate: BigDecimal = finalAmount2.divide(finalAmount1);
        const displayRateStr = formatExchangeRateAmount(rate, options);
        return `${numeralSystem.getLocalizedDigit(1)} : ${displayRateStr}`;
    }
}

export function getExchangedAmountByRate(amount: BigDecimal, fromRate: string, toRate: string): BigDecimal | null {
    try {
        const exchangeRate = parseBigDecimal(toRate).divide(parseBigDecimal(fromRate));
        return amount.multiply(exchangeRate);
    } catch (ex) {
        logger.warn(`failed to get exchanged amount by rate: ${ex}`);
        return null;
    }
}
