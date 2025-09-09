<template>
    <v-text-field type="text" class="text-field-with-colored-label" :class="extraClass"
                  :color="color" :base-color="color"
                  :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  :rules="enableRules ? rules : []" v-model="currentValue" v-if="!hide && !formulaMode"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown" @paste="onPaste" @click="onClick">
        <template #prepend-inner v-if="currency && prependText">
            <div>{{ prependText }}</div>
        </template>
        <template #append-inner>
            <div class="text-no-wrap" v-if="currency && appendText">{{ appendText }}</div>
            <v-tooltip :text="tt('Enter formula mode')">
                <template v-slot:activator="{ props }">
                    <v-icon class="ms-2" :icon="mdiCalculatorVariantOutline"
                            @keydown.enter="enterFormulaMode" @keydown.space="enterFormulaMode" @click="enterFormulaMode"
                            v-bind="props" v-if="enableFormula && !formulaMode"></v-icon>
                </template>
            </v-tooltip>
        </template>
    </v-text-field>
    <v-text-field type="text" class="text-field-with-colored-label" :class="extraClass"
                  :color="color" :base-color="color"
                  :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  v-model="currentFormula" v-if="!hide && formulaMode"
                  @keydown.enter="calculateFormula" @click="onClick">
        <template #prepend-inner v-if="currency && prependText">
            <div>{{ prependText }}</div>
        </template>
        <template #append-inner>
            <div class="text-no-wrap" v-if="currency && appendText">{{ appendText }}</div>
            <v-tooltip :text="tt('Calculate formula result')">
                <template v-slot:activator="{ props }">
                    <v-icon class="ms-2" color="primary" :icon="mdiCheck"
                            @click="calculateFormula" v-bind="props"
                            v-if="formulaMode"></v-icon>
                </template>
            </v-tooltip>
            <v-tooltip :text="tt('Exit formula mode')">
                <template v-slot:activator="{ props }">
                    <v-icon class="ms-2" color="secondary" :icon="mdiClose"
                            @click="exitFormulaMode" v-bind="props"
                            v-if="formulaMode"></v-icon>
                </template>
            </v-tooltip>
        </template>
    </v-text-field>
    <v-text-field type="password" class="text-field-with-colored-label" :class="extraClass"
                  :color="color" :base-color="color"
                  :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  :rules="enableRules ? rules : []" v-model="currentValue" v-if="hide"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown" @paste="onPaste" @click="onClick">
        <template #prepend-inner v-if="currency && prependText">
            <div>{{ prependText }}</div>
        </template>
        <template #append-inner v-if="currency && appendText">
            <div class="text-no-wrap">{{ appendText }}</div>
        </template>
    </v-text-field>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonNumberInputProps, useCommonNumberInputBase } from '@/components/base/CommonNumberInputBase.ts';

import { NumeralSystem, DecimalSeparator } from '@/core/numeral.ts';
import type { CurrencyPrependAndAppendText } from '@/core/currency.ts';
import { DEFAULT_DECIMAL_NUMBER_COUNT } from '@/consts/numeral.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { isNumber, replaceAll } from '@/lib/common.ts';
import { evaluateExpressionToAmount } from '@/lib/evaluator.ts';
import type { ComponentDensity } from '@/lib/ui/desktop.ts';
import logger from '@/lib/logger.ts';

import {
    mdiCalculatorVariantOutline,
    mdiCheck,
    mdiClose
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

interface DesktopAmountInputProps extends CommonNumberInputProps {
    class?: string;
    color?: string;
    density?: ComponentDensity;
    currency: string;
    showCurrency?: boolean;
    persistentPlaceholder?: boolean;
    hide?: boolean;
    enableRules?: boolean;
    enableFormula?: boolean;
    flipNegative?: boolean;
}

const props = defineProps<DesktopAmountInputProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
}>();

const {
    tt,
    getCurrentNumeralSystemType,
    getCurrentDecimalSeparator,
    parseAmountFromLocalizedNumerals,
    formatAmountToLocalizedNumeralsWithoutDigitGrouping,
    getAmountPrependAndAppendText
} = useI18n();

const {
    currentValue,
    onKeyUpDown,
    onPaste
} = useCommonNumberInputBase(props, DEFAULT_DECIMAL_NUMBER_COUNT, getInitedFormattedValue(props.modelValue, props.flipNegative), parseAmountFromLocalizedNumerals, getFormattedValue, getValidFormattedValue);

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const rules = [
    (v: string) => {
        if (v === '') {
            return tt('Amount value is not number');
        }

        try {
            const val = parseAmountFromLocalizedNumerals(v);

            if (Number.isNaN(val) || !Number.isFinite(val)) {
                return tt('Amount value is not number');
            }

            return (val >= TRANSACTION_MIN_AMOUNT && val <= TRANSACTION_MAX_AMOUNT) || tt('Amount value exceeds limitation');
        } catch (ex) {
            logger.warn('cannot parse amount in amount input, original value is ' + v, ex);
            return tt('Amount value is not number');
        }
    }
];

const currentFormula = ref<string>('');
const formulaMode = ref<boolean>(false);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

const prependText = computed<string | undefined>(() => {
    if (!props.currency || !props.showCurrency) {
        return '';
    }

    const texts = getDisplayCurrencyPrependAndAppendText();

    if (!texts) {
        return '';
    }

    return texts.prependText;
});

const appendText = computed<string | undefined>(() => {
    if (!props.currency || !props.showCurrency) {
        return '';
    }

    const texts = getDisplayCurrencyPrependAndAppendText();

    if (!texts) {
        return '';
    }

    return texts.appendText;
});

const extraClass = computed<string>(() => {
    let finalClass = props.class || '';

    if (props.color) {
        finalClass += ` text-${props.color}`;
    }

    if (props.currency && prependText.value) {
        finalClass += ` has-pretend-text`;
    }

    return finalClass;
});

function enterFormulaMode(): void {
    if (!props.enableFormula) {
        return;
    }

    currentFormula.value = currentValue.value;
    formulaMode.value = true;
}

function calculateFormula(): void {
    const systemDecimalSeparator = DecimalSeparator.Dot.symbol;
    const decimalSeparator = getCurrentDecimalSeparator();
    let finalFormula = currentFormula.value;

    if (systemDecimalSeparator !== decimalSeparator && finalFormula.indexOf(systemDecimalSeparator) >= 0) {
        snackbar.value?.showMessage('Formula is invalid');
        return;
    } else if (systemDecimalSeparator !== decimalSeparator) {
        finalFormula = replaceAll(currentFormula.value, decimalSeparator, systemDecimalSeparator);
    }

    finalFormula = numeralSystem.value.replaceLocalizedDigitsToWesternArabicDigits(finalFormula);

    try {
        const calculatedAmount = evaluateExpressionToAmount(finalFormula);

        if (isNumber(calculatedAmount)) {
            const textualValue = getFormattedValue(calculatedAmount);
            const hasDecimalSeparator = textualValue.indexOf(decimalSeparator) >= 0;
            currentValue.value = getValidFormattedValue(calculatedAmount, textualValue, hasDecimalSeparator);
            formulaMode.value = false;
        } else {
            snackbar.value?.showMessage('Formula is invalid');
        }
    } catch (ex) {
        logger.error('cannot evaluate formula in amount input, original formula is ' + finalFormula, ex);

        if (ex instanceof Error) {
            snackbar.value?.showMessage(ex.message);
        }
    }
}

function exitFormulaMode(): void {
    formulaMode.value = false;
    currentFormula.value = '';
}

function onClick(e: MouseEvent): void {
    if (!props.disabled && !props.readonly && props.modelValue === 0 && e.target instanceof HTMLInputElement) {
        const input = e.target as HTMLInputElement;

        if ((!input?.selectionStart && !input?.selectionEnd) || input?.selectionStart === input?.selectionEnd) {
            input?.select();
        }
    }
}

function getValidFormattedValue(value: number, textualValue: string, hasDecimalSeparator: boolean): string {
    let maxLength = TRANSACTION_MAX_AMOUNT.toString().length;

    if (value < 0) {
        maxLength = TRANSACTION_MIN_AMOUNT.toString().length;
    }

    if (value < TRANSACTION_MIN_AMOUNT) {
        return getFormattedValue(TRANSACTION_MIN_AMOUNT);
    } else if (value > TRANSACTION_MAX_AMOUNT) {
        return getFormattedValue(TRANSACTION_MAX_AMOUNT);
    }

    if (!hasDecimalSeparator && textualValue.length > maxLength) {
        return textualValue.substring(0, maxLength);
    } else if (hasDecimalSeparator && textualValue.length > maxLength + 1) {
        return textualValue.substring(0, maxLength + 1);
    }

    return textualValue;
}

function getInitedFormattedValue(value: number, flipNegative?: boolean): string {
    if (flipNegative) {
        value = -value;
    }

    return getFormattedValue(value);
}

function getFormattedValue(value: number): string {
    if (!Number.isNaN(value) && Number.isFinite(value)) {
        return formatAmountToLocalizedNumeralsWithoutDigitGrouping(value, props.currency);
    }

    return numeralSystem.value.digitZero;
}

function getDisplayCurrencyPrependAndAppendText(): CurrencyPrependAndAppendText | null {
    const numericCurrentValue = parseAmountFromLocalizedNumerals(currentValue.value);
    const isPlural = numericCurrentValue !== 100 && numericCurrentValue !== -100;

    return getAmountPrependAndAppendText(props.currency, isPlural);
}

watch(() => props.currency, () => {
    const newStringValue = getInitedFormattedValue(props.modelValue, props.flipNegative);

    if (!(newStringValue === numeralSystem.value.digitZero && currentValue.value === '')) {
        currentValue.value = newStringValue;
    }
});

watch(() => props.flipNegative, (newValue) => {
    const newStringValue = getInitedFormattedValue(props.modelValue, newValue);

    if (!(newStringValue === numeralSystem.value.digitZero && currentValue.value === '')) {
        currentValue.value = newStringValue;
    }
});

watch(() => props.modelValue, (newValue) => {
    if (props.flipNegative) {
        newValue = -newValue;
    }

    const numericCurrentValue = parseAmountFromLocalizedNumerals(currentValue.value);

    if (newValue !== numericCurrentValue) {
        const newStringValue = getFormattedValue(newValue);

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

                finalValue += ch;
            }

            finalValue = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(finalValue);
        } else if (newValue === '-' || newValue === decimalSeparator || newValue === `-${decimalSeparator}`) {
            finalValue = newValue;
        }
    }

    if (finalValue !== newValue) {
        currentValue.value = finalValue;
    } else {
        let value: number = parseAmountFromLocalizedNumerals(finalValue);

        if (Number.isNaN(value) || !Number.isFinite(value)) {
            value = 0;
        }

        if (props.flipNegative) {
            value = -value;
        }

        emit('update:modelValue', value);
    }
});
</script>

<style>
.text-field-with-colored-label.has-pretend-text .v-field {
    grid-template-columns: max-content minmax(0, 1fr) min-content min-content;
}

.text-field-with-colored-label.has-pretend-text .v-field__input {
    padding-inline-start: 0.5rem;
}
</style>
