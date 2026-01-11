import { computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type CommonNumberInputProps,
    type CommonNumberInputEmits,
    useCommonNumberInputBase
} from '@/components/base/CommonNumberInputBase.ts';

import { isDefined, isNumber, replaceAll, removeAll } from '@/lib/common.ts';
import { NumeralSystem } from '@/core/numeral.ts';

export interface NumberInputProps extends CommonNumberInputProps {
    minValue?: number;
    maxValue?: number;
    maxDecimalCount?: number;
}

export type NumberInputEmits = CommonNumberInputEmits;

export function useNumberInputBase(props: NumberInputProps, emit: CommonNumberInputEmits) {
    const {
        getCurrentNumeralSystemType,
        getCurrentDecimalSeparator,
        getCurrentDigitGroupingSymbol
    } = useI18n();

    const {
        currentValue,
        onKeyUpDown,
        onPaste
    } = useCommonNumberInputBase(props, emit, props.maxDecimalCount ?? -1, getFormattedValue(props.modelValue), parseNumber, getFormattedValue, getValidFormattedValue);

    const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

    function parseNumber(value: string): number {
        if (!value) {
            return 0;
        }

        const decimalSeparator = getCurrentDecimalSeparator();

        let finalValue = '';
        for (let i = 0; i < value.length; i++) {
            const ch = value.charAt(i);

            if (!NumeralSystem.WesternArabicNumerals.isDigit(ch) && !numeralSystem.value.isDigit(ch) && ch !== '-' && ch !== decimalSeparator) {
                break;
            }

            finalValue += ch;
        }

        finalValue = numeralSystem.value.replaceLocalizedDigitsToWesternArabicDigits(finalValue);

        if (decimalSeparator !== '.') {
            finalValue = replaceAll(finalValue, decimalSeparator, '.');
        }

        return parseFloat(finalValue);
    }

    function getValidFormattedValue(value: number, textualValue: string): string {
        if (isNumber(props.minValue) && value < props.minValue) {
            return getFormattedValue(props.minValue, numeralSystem.value);
        }

        if (isNumber(props.maxValue) && value > props.maxValue) {
            return getFormattedValue(props.maxValue, numeralSystem.value);
        }

        const decimalSeparator = getCurrentDecimalSeparator();
        const digitGroupingSymbol = getCurrentDigitGroupingSymbol();
        return replaceAll(removeAll(textualValue, digitGroupingSymbol), '.', decimalSeparator);
    }

    function getFormattedValue(value: number, customNumeralSystem?: NumeralSystem): string {
        if (!isDefined(customNumeralSystem)) {
            customNumeralSystem = getCurrentNumeralSystemType();
        }

        if (!Number.isNaN(value) && Number.isFinite(value)) {
            const decimalSeparator = getCurrentDecimalSeparator();

            if (isNumber(props.maxDecimalCount) && props.maxDecimalCount >= 0) {
                return replaceAll(customNumeralSystem.replaceWesternArabicDigitsToLocalizedDigits(value.toFixed(props.maxDecimalCount)), '.', decimalSeparator);
            } else {
                return replaceAll(customNumeralSystem.replaceWesternArabicDigitsToLocalizedDigits(value.toString(10)), '.', decimalSeparator);
            }
        }

        return customNumeralSystem.digitZero;
    }

    watch(() => props.modelValue, (newValue) => {
        const numericCurrentValue = parseNumber(currentValue.value);

        if (newValue !== numericCurrentValue) {
            const newStringValue = getFormattedValue(newValue, numeralSystem.value);

            if (!(newStringValue === numeralSystem.value.digitZero && currentValue.value === '')) {
                currentValue.value = newStringValue;
            }
        }
    });

    watch(currentValue, (newValue) => {
        let actualNumeralSystem: NumeralSystem | undefined = undefined;
        let finalValue = '';

        if (newValue) {
            const decimalSeparator = getCurrentDecimalSeparator();

            if (newValue[0] === '-' || newValue[0] === decimalSeparator) {
                actualNumeralSystem = NumeralSystem.detect(newValue.charAt(1));
            } else {
                actualNumeralSystem = NumeralSystem.detect(newValue.charAt(0));
            }

            if (actualNumeralSystem && (actualNumeralSystem.type === NumeralSystem.WesternArabicNumerals.type || actualNumeralSystem.type === numeralSystem.value.type)) {
                for (let i = 0; i < newValue.length; i++) {
                    const ch = newValue.charAt(i);

                    if (!NumeralSystem.WesternArabicNumerals.isDigit(ch) && !numeralSystem.value.isDigit(ch) && ch !== '-' && ch !== decimalSeparator) {
                        break;
                    }

                    finalValue += newValue[i];
                }

                finalValue = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(finalValue);
            } else if (newValue === '-' || newValue === decimalSeparator || newValue === `-${decimalSeparator}`) {
                finalValue = newValue;
            }
        }

        if (finalValue !== newValue) {
            currentValue.value = finalValue;
        } else {
            let value: number = parseNumber(finalValue);

            if (Number.isNaN(value) || !Number.isFinite(value)) {
                value = 0;
            }

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
