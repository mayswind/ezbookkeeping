import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { NumeralSystem } from '@/core/numeral.ts';

import { removeAll } from '@/lib/common.ts';
import logger from '@/lib/logger.ts';

export interface CommonNumberInputProps {
    label?: string;
    placeholder?: string;
    disabled?: boolean;
    readonly?: boolean;
    modelValue: number;
}

export type ParseNumberFunction = (value: string) => number;
export type FormatNumberFunction = (value: number) => string;
export type GetValidFormattedValueFunction = (value: number, textualValue: string, hasDecimalSeparator: boolean) => string;

export function useCommonNumberInputBase(props: CommonNumberInputProps, maxDecimalCount: number, initValue: string, parseNumber: ParseNumberFunction, formatNumber: FormatNumberFunction, getValidFormattedValue: GetValidFormattedValueFunction) {
    const {
        getCurrentNumeralSystemType,
        getCurrentDecimalSeparator,
        getCurrentDigitGroupingSymbol
    } = useI18n();

    const currentValue = ref<string>(initValue);

    function onKeyUpDown(e: KeyboardEvent): void {
        if (e.altKey || e.ctrlKey || e.metaKey || (e.key.indexOf('F') === 0 && (e.key.length === 2 || e.key.length === 3))
            || e.key === 'ArrowLeft' || e.key === 'ArrowRight'
            || e.key === 'Home' || e.key === 'End' || e.key === 'Tab'
            || e.key === 'Backspace' || e.key === 'Delete' || e.key === 'Del') {
            return;
        }

        if (props.readonly || props.disabled) {
            e.preventDefault();
            return;
        }

        const numeralSystem = getCurrentNumeralSystemType();
        const digitGroupingSymbol = getCurrentDigitGroupingSymbol();
        const decimalSeparator = getCurrentDecimalSeparator();

        if (!NumeralSystem.WesternArabicNumerals.isDigit(e.key) && !numeralSystem.isDigit(e.key) && e.key !== '-' && e.key !== decimalSeparator) {
            e.preventDefault();
            return;
        }

        if (maxDecimalCount === 0 && e.key === decimalSeparator) {
            e.preventDefault();
            return;
        }

        if (!e.target) {
            return;
        }

        const target = e.target as HTMLInputElement;

        let str = target.value;

        if (str.indexOf(digitGroupingSymbol) >= 0) {
            str = removeAll(str, digitGroupingSymbol);
        }

        if (e.key === '-' && str.lastIndexOf('-') > 0) {
            const lastMinusPos = str.lastIndexOf('-');
            target.value = str.substring(0, lastMinusPos) + str.substring(lastMinusPos + 1, str.length);
            currentValue.value = target.value;
            e.preventDefault();
            return;
        }

        if (e.key === decimalSeparator && str.indexOf(decimalSeparator) !== str.lastIndexOf(decimalSeparator)) {
            const lastDecimalSeparatorPos = str.lastIndexOf(decimalSeparator);
            target.value = str.substring(0, lastDecimalSeparatorPos) + str.substring(lastDecimalSeparatorPos + 1, str.length);
            currentValue.value = target.value;
            e.preventDefault();
            return;
        }

        if (e.key === decimalSeparator && (str.indexOf(decimalSeparator) === 0 || (str.indexOf(decimalSeparator) === 1 && str.charAt(0) === '-'))) {
            const negative = str.charAt(0) === '-';

            if (negative) {
                str = str.substring(1);
            }

            str = (negative ? `-${numeralSystem.digitZero}` : numeralSystem.digitZero) + str;
            target.value = str;
            currentValue.value = target.value;
            e.preventDefault();
            return;
        }

        let decimalLength = 0;
        const decimalIndex = str.indexOf(decimalSeparator);

        if (decimalIndex >= 0) {
            decimalLength = str.length - str.indexOf(decimalSeparator) - 1;
        } else if ((str.startsWith(numeralSystem.digitZero) && str.length >= 2) || (str.startsWith(`-${numeralSystem.digitZero}`) && str.length >= 3)) {
            const negative = str.charAt(0) === '-';

            if (negative) {
                str = str.substring(1);
            }

            while (str.charAt(0) === numeralSystem.digitZero && (str.length >= 2 || e.key !== numeralSystem.digitZero)) {
                str = str.substring(1);
            }

            target.value = (negative ? '-' : '') + str;
            currentValue.value = target.value;
            e.preventDefault();
            return;
        }

        if (maxDecimalCount > 0 && decimalLength > maxDecimalCount) {
            target.value = str.substring(0, Math.min(decimalIndex + maxDecimalCount + 1, str.length - 1));
            currentValue.value = target.value;
            e.preventDefault();
            return;
        } else if (maxDecimalCount === 0 && decimalIndex >= 0) {
            target.value = str.substring(0, decimalIndex);
            currentValue.value = target.value;
            e.preventDefault();
            return;
        }

        try {
            const val = parseNumber(str);
            const finalValue = getValidFormattedValue(val, str, decimalIndex >= 0);

            if (finalValue !== str) {
                target.value = finalValue;
                currentValue.value = finalValue;
                e.preventDefault();
            }
        } catch (ex) {
            logger.warn('cannot parse input number, original value is ' + str, ex);
            target.value = numeralSystem.digitZero;
        }
    }

    function onPaste(e: ClipboardEvent): void {
        if (!e.clipboardData || props.readonly || props.disabled) {
            e.preventDefault();
            return;
        }

        const text = e.clipboardData.getData('Text');

        if (!text) {
            e.preventDefault();
            return;
        }

        const value = parseNumber(text);
        const textualValue = formatNumber(value);
        const decimalSeparator = getCurrentDecimalSeparator();
        const hasDecimalSeparator = text.indexOf(decimalSeparator) >= 0;

        currentValue.value = getValidFormattedValue(value, textualValue, hasDecimalSeparator);
        e.preventDefault();
    }

    return {
        // states
        currentValue,
        // functions
        onKeyUpDown,
        onPaste
    };
}
