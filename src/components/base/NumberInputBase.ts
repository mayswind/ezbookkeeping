import { watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonNumberInputProps, useCommonNumberInputBase } from '@/components/base/CommonNumberInputBase.ts';

import { isNumber, replaceAll, removeAll } from '@/lib/common.ts';

export interface NumberInputProps extends CommonNumberInputProps {
    minValue?: number;
    maxValue?: number;
    maxDecimalCount?: number;
}

export interface NumberInputEmits {
    (e: 'update:modelValue', value: number): void;
}

export function useNumberInputBase(props: NumberInputProps, emit: NumberInputEmits) {
    const {
        getCurrentDecimalSeparator,
        getCurrentDigitGroupingSymbol
    } = useI18n();

    const {
        currentValue,
        onKeyUpDown,
        onPaste
    } = useCommonNumberInputBase(props, props.maxDecimalCount ?? -1, getFormattedValue(props.modelValue), parseNumber, getFormattedValue, getValidFormattedValue);

    function parseNumber(value: string): number {
        if (!value) {
            return 0;
        }

        const decimalSeparator = getCurrentDecimalSeparator();

        let finalValue = '';
        for (let i = 0; i < value.length; i++) {
            if (!('0' <= value[i] && value[i] <= '9') && value[i] !== '-' && value[i] !== decimalSeparator) {
                break;
            }

            finalValue += value[i];
        }

        if (decimalSeparator !== '.') {
            finalValue = replaceAll(finalValue, decimalSeparator, '.');
        }

        return parseFloat(finalValue);
    }

    function getValidFormattedValue(value: number, textualValue: string): string {
        if (isNumber(props.minValue) && value < props.minValue) {
            return getFormattedValue(props.minValue);
        }

        if (isNumber(props.maxValue) && value > props.maxValue) {
            return getFormattedValue(props.maxValue);
        }

        const decimalSeparator = getCurrentDecimalSeparator();
        const digitGroupingSymbol = getCurrentDigitGroupingSymbol();
        return replaceAll(removeAll(textualValue, digitGroupingSymbol), '.', decimalSeparator);
    }

    function getFormattedValue(value: number): string {
        if (!Number.isNaN(value) && Number.isFinite(value)) {
            const decimalSeparator = getCurrentDecimalSeparator();

            if (isNumber(props.maxDecimalCount) && props.maxDecimalCount >= 0) {
                return replaceAll(value.toFixed(props.maxDecimalCount), '.', decimalSeparator);
            } else {
                return replaceAll(value.toString(), '.', decimalSeparator);
            }
        }

        return '0';
    }

    watch(() => props.modelValue, (newValue) => {
        const numericCurrentValue = parseNumber(currentValue.value);

        if (newValue !== numericCurrentValue) {
            const newStringValue = getFormattedValue(newValue);

            if (!(newStringValue === '0' && currentValue.value === '')) {
                currentValue.value = newStringValue;
            }
        }
    });

    watch(currentValue, (newValue) => {
        let finalValue = '';

        if (newValue) {
            const decimalSeparator = getCurrentDecimalSeparator();

            for (let i = 0; i < newValue.length; i++) {
                if (!('0' <= newValue[i] && newValue[i] <= '9') && newValue[i] !== '-' && newValue[i] !== decimalSeparator) {
                    break;
                }

                finalValue += newValue[i];
            }
        }

        if (finalValue !== newValue) {
            currentValue.value = finalValue;
        } else {
            const value: number = parseNumber(finalValue);
            emit('update:modelValue', value);
        }
    });

    return {
        // states
        currentValue,
        // functions
        onKeyUpDown,
        onPaste
    }
}
