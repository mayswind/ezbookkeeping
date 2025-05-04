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
                    <v-icon class="ml-2" :icon="mdiCalculatorVariantOutline"
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
                    <v-icon class="ml-2" color="primary" :icon="mdiCheck"
                            @click="calculateFormula" v-bind="props"
                            v-if="formulaMode"></v-icon>
                </template>
            </v-tooltip>
            <v-tooltip :text="tt('Exit formula mode')">
                <template v-slot:activator="{ props }">
                    <v-icon class="ml-2" color="secondary" :icon="mdiClose"
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

import { DecimalSeparator } from '@/core/numeral.ts';
import type { CurrencyPrependAndAppendText } from '@/core/currency.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { isNumber, replaceAll, removeAll } from '@/lib/common.ts';
import { evaluateExpression } from '@/lib/evaluator.ts';
import type { ComponentDensity } from '@/lib/ui/desktop.ts';
import logger from '@/lib/logger.ts';

import {
    mdiCalculatorVariantOutline,
    mdiCheck,
    mdiClose
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    class?: string;
    color?: string;
    density?: ComponentDensity;
    currency: string;
    showCurrency?: boolean;
    label?: string;
    placeholder?: string;
    persistentPlaceholder?: boolean;
    disabled?: boolean;
    readonly?: boolean;
    hide?: boolean;
    enableRules?: boolean;
    enableFormula?: boolean;
    flipNegative?: boolean;
    modelValue: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
}>();

const {
    tt,
    getCurrentDecimalSeparator,
    getCurrentDigitGroupingSymbol,
    parseAmount,
    formatAmount,
    formatNumber,
    getAmountPrependAndAppendText
} = useI18n();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const rules = [
    (v: string) => {
        if (v === '') {
            return tt('Amount value is not number');
        }

        try {
            const val = parseAmount(v);

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

const currentValue = ref<string>(getInitedFormattedValue(props.modelValue, props.flipNegative));
const currentFormula = ref<string>('');
const formulaMode = ref<boolean>(false);

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

    const calculatedValue = evaluateExpression(finalFormula);

    if (isNumber(calculatedValue)) {
        const textualValue = formatNumber(calculatedValue, 2);
        const hasDecimalSeparator = textualValue.indexOf(decimalSeparator) >= 0;
        currentValue.value = getValidFormattedValue(calculatedValue, textualValue, hasDecimalSeparator);
        formulaMode.value = false;
    } else {
        snackbar.value?.showMessage('Formula is invalid');
    }
}

function exitFormulaMode(): void {
    formulaMode.value = false;
    currentFormula.value = '';
}

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

    const digitGroupingSymbol = getCurrentDigitGroupingSymbol();
    const decimalSeparator = getCurrentDecimalSeparator();

    if (!('0' <= e.key && e.key <= '9') && e.key !== '-' && e.key !== decimalSeparator) {
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

        str = (negative ? '-0' : '0') + str;
        target.value = str;
        currentValue.value = target.value;
        e.preventDefault();
        return;
    }

    let decimalLength = 0;
    const decimalIndex = str.indexOf(decimalSeparator);

    if (decimalIndex >= 0) {
        decimalLength = str.length - str.indexOf(decimalSeparator) - 1;
    } else if ((str.startsWith('0') && str.length >= 2) || (str.startsWith('-0') && str.length >= 3)) {
        const negative = str.charAt(0) === '-';

        if (negative) {
            str = str.substring(1);
        }

        while (str.charAt(0) === '0' && (str.length >= 2 || e.key !== '0')) {
            str = str.substring(1);
        }

        target.value = (negative ? '-' : '') + str;
        currentValue.value = target.value;
        e.preventDefault();
        return;
    }

    if (decimalLength > 2) {
        target.value = str.substring(0, Math.min(decimalIndex + 3, str.length - 1));
        currentValue.value = target.value;
        e.preventDefault();
        return;
    }

    try {
        const val = parseAmount(str);
        const finalValue = getValidFormattedValue(val, str, decimalIndex >= 0);

        if (finalValue !== str) {
            target.value = finalValue;
            currentValue.value = finalValue;
            e.preventDefault();
        }
    } catch (ex) {
        logger.warn('cannot parse amount in amount input, original value is ' + str, ex);
        target.value = '0';
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

    const value = parseAmount(text);
    const textualValue = getFormattedValue(value);
    const decimalSeparator = getCurrentDecimalSeparator();
    const hasDecimalSeparator = text.indexOf(decimalSeparator) >= 0;

    currentValue.value = getValidFormattedValue(value, textualValue, hasDecimalSeparator);
    e.preventDefault();
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
        const digitGroupingSymbol = getCurrentDigitGroupingSymbol();
        return removeAll(formatAmount(value, props.currency), digitGroupingSymbol);
    }

    return '0';
}

function getDisplayCurrencyPrependAndAppendText(): CurrencyPrependAndAppendText | null {
    const numericCurrentValue = parseAmount(currentValue.value);
    const isPlural = numericCurrentValue !== 100 && numericCurrentValue !== -100;

    return getAmountPrependAndAppendText(props.currency, isPlural);
}

watch(() => props.currency, () => {
    const newStringValue = getInitedFormattedValue(props.modelValue, props.flipNegative);

    if (!(newStringValue === '0' && currentValue.value === '')) {
        currentValue.value = newStringValue;
    }
});

watch(() => props.flipNegative, (newValue) => {
    const newStringValue = getInitedFormattedValue(props.modelValue, newValue);

    if (!(newStringValue === '0' && currentValue.value === '')) {
        currentValue.value = newStringValue;
    }
});

watch(() => props.modelValue, (newValue) => {
    if (props.flipNegative) {
        newValue = -newValue;
    }

    const numericCurrentValue = parseAmount(currentValue.value);

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
        let value: number = parseAmount(finalValue);

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
    padding-left: 0.5rem;
}
</style>
