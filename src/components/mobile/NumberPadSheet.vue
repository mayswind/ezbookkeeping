<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="numpad-sheet" style="height: auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="numpad-values">
                <span class="numpad-value" :class="currentDisplayNumClass">{{ currentDisplay }}</span>
            </div>
            <div class="numpad-buttons">
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(7)">
                    <span class="numpad-button-text numpad-button-text-normal">7</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(8)">
                    <span class="numpad-button-text numpad-button-text-normal">8</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(9)">
                    <span class="numpad-button-text numpad-button-text-normal">9</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-function no-right-border" @click="setSymbol('×')">
                    <span class="numpad-button-text numpad-button-text-normal">&times;</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(4)">
                    <span class="numpad-button-text numpad-button-text-normal">4</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(5)">
                    <span class="numpad-button-text numpad-button-text-normal">5</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(6)">
                    <span class="numpad-button-text numpad-button-text-normal">6</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-function no-right-border" @click="setSymbol('−')">
                    <span class="numpad-button-text numpad-button-text-normal">&minus;</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(1)">
                    <span class="numpad-button-text numpad-button-text-normal">1</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(2)">
                    <span class="numpad-button-text numpad-button-text-normal">2</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(3)">
                    <span class="numpad-button-text numpad-button-text-normal">3</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-function no-right-border" @click="setSymbol('+')">
                    <span class="numpad-button-text numpad-button-text-normal">&plus;</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" v-if="supportDecimalSeparator" @click="inputDecimalSeparator()">
                    <span class="numpad-button-text numpad-button-text-normal">{{ decimalSeparator }}</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" v-if="!supportDecimalSeparator" @click="inputDoubleNum(0)">
                    <span class="numpad-button-text numpad-button-text-normal">00</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="inputNum(0)">
                    <span class="numpad-button-text numpad-button-text-normal">0</span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-num" @click="backspace" @taphold="clear()">
                <span class="numpad-button-text numpad-button-text-normal">
                    <f7-icon f7="delete_left"></f7-icon>
                </span>
                </f7-button>
                <f7-button class="numpad-button numpad-button-confirm no-right-border no-bottom-border" fill @click="confirm()">
                    <span :class="{ 'numpad-button-text': true, 'numpad-button-text-confirm': !currentSymbol }">{{ confirmText }}</span>
                </f7-button>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';

import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { isString, isNumber, removeAll } from '@/lib/common.ts';

const {
    tt,
    getCurrentDecimalSeparator,
    getCurrentDigitGroupingSymbol,
    appendDigitGroupingSymbol,
    parseAmount,
    formatAmount
} = useI18n();
const { showToast } = useI18nUIComponents();

const props = defineProps<{
    modelValue: number | string;
    minValue?: number;
    maxValue?: number;
    currency?: string;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
    (e: 'update:show', value: boolean): void;
}>();

const previousValue = ref<string>('');
const currentSymbol = ref<string>('');
const currentValue = ref<string>(getStringValue(props.modelValue));

const decimalSeparator = computed<string>(() => getCurrentDecimalSeparator());

const supportDecimalSeparator = computed<boolean>(() => {
    if (!props.currency || !ALL_CURRENCIES[props.currency] || !isNumber(ALL_CURRENCIES[props.currency].fraction)) {
        return true;
    }

    return (ALL_CURRENCIES[props.currency].fraction as number) > 0;
});

const currentDisplay = computed<string>(() => {
    const finalPreviousValue = appendDigitGroupingSymbol(previousValue.value);
    const finalCurrentValue = appendDigitGroupingSymbol(currentValue.value);

    if (currentSymbol.value) {
        return `${finalPreviousValue} ${currentSymbol.value} ${finalCurrentValue}`;
    } else {
        return finalCurrentValue;
    }
});

const currentDisplayNumClass = computed<string>(() => {
    if (currentDisplay.value && currentDisplay.value.length >= 24) {
        return 'numpad-value-small';
    } else if (currentDisplay.value && currentDisplay.value.length >= 16) {
        return 'numpad-value-normal';
    } else {
        return 'numpad-value-large';
    }
});

const confirmText = computed<string>(() => {
    if (currentSymbol.value) {
        return '=';
    } else {
        return tt('OK');
    }
});

function getStringValue(value: number | string): string {
    if (!isNumber(value) && !isString(value)) {
        return '';
    }

    let str = formatAmount(value, props.currency);

    const digitGroupingSymbol = getCurrentDigitGroupingSymbol();

    if (str.indexOf(digitGroupingSymbol) >= 0) {
        str = removeAll(str, digitGroupingSymbol);
    }

    const decimalSeparator = getCurrentDecimalSeparator();
    const decimalSeparatorPos = str.indexOf(decimalSeparator);

    if (decimalSeparatorPos < 0) {
        if (str === '0') {
            return '';
        }

        return str;
    }

    const integer = str.substring(0, decimalSeparatorPos);
    const decimals = str.substring(decimalSeparatorPos + 1, str.length);
    let newDecimals = '';

    for (let i = decimals.length - 1; i >= 0; i--) {
        if (decimals[i] !== '0' || newDecimals.length > 0) {
            newDecimals = decimals[i] + newDecimals;
        }
    }

    if (newDecimals.length < 1) {
        if (integer === '0') {
            return '';
        }

        return integer;
    }

    return `${integer}${decimalSeparator}${newDecimals}`;
}

function inputNum(num: number): void {
    if (!previousValue.value && currentSymbol.value === '−') {
        currentValue.value = '-' + currentValue.value;
        currentSymbol.value = '';
    }

    if (currentValue.value === '0') {
        currentValue.value = num.toString();
        return;
    } else if (currentValue.value === '-0') {
        currentValue.value = '-' + num.toString();
        return;
    }

    const decimalSeparatorPos = currentValue.value.indexOf(decimalSeparator.value);

    if (decimalSeparatorPos >= 0 && currentValue.value.substring(decimalSeparatorPos + 1, currentValue.value.length).length >= 2) {
        return;
    }

    const newValue = currentValue.value + num.toString();

    if (isNumber(props.minValue)) {
        const current = parseAmount(newValue);

        if (current < (props.minValue as number)) {
            return;
        }
    }

    if (isNumber(props.maxValue)) {
        const current = parseAmount(newValue);

        if (current > (props.maxValue as number)) {
            return;
        }
    }

    currentValue.value = newValue;
}

function inputDoubleNum(num: number): void {
    inputNum(num);
    inputNum(num);
}

function inputDecimalSeparator(): void {
    if (currentValue.value.indexOf(decimalSeparator.value) >= 0) {
        return;
    }

    if (!previousValue.value && currentSymbol.value === '−') {
        currentValue.value = '-' + currentValue.value;
        currentSymbol.value = '';
    }

    if (currentValue.value.length < 1) {
        currentValue.value = '0';
    } else if (currentValue.value === '-') {
        currentValue.value = '-0';
    }

    currentValue.value = currentValue.value + decimalSeparator.value;
}

function setSymbol(symbol: string): void {
    if (currentValue.value) {
        if (currentSymbol.value) {
            const lastFormulaCalcResult = confirm();

            if (!lastFormulaCalcResult) {
                return;
            }
        }

        previousValue.value = currentValue.value;
        currentValue.value = '';
    }

    currentSymbol.value = symbol;
}

function backspace(): void {
    if (!currentValue.value || currentValue.value.length < 1) {
        if (currentSymbol.value) {
            currentValue.value = previousValue.value;
            previousValue.value = '';
            currentSymbol.value = '';
        }

        return;
    }

    currentValue.value = currentValue.value.substring(0, currentValue.value.length - 1);
}

function clear(): void {
    currentValue.value = '';
    previousValue.value = '';
    currentSymbol.value = '';
}

function confirm(): boolean {
    if (currentSymbol.value && currentValue.value.length >= 1) {
        const previous = parseAmount(previousValue.value);
        const current = parseAmount(currentValue.value);
        let finalValue = 0;

        switch (currentSymbol.value) {
            case '+':
                finalValue = previous + current;
                break;
            case '−':
                finalValue = previous - current;
                break;
            case '×':
                finalValue = Math.round(previous * current / 100);
                break;
            default:
                finalValue = previous;
        }

        if (isNumber(props.minValue)) {
            if (finalValue < (props.minValue as number)) {
                showToast('Numeric Overflow');
                return false;
            }
        }

        if (isNumber(props.maxValue)) {
            if (finalValue > (props.maxValue as number)) {
                showToast('Numeric Overflow');
                return false;
            }
        }

        currentValue.value = getStringValue(finalValue);
        previousValue.value = '';
        currentSymbol.value = '';

        return true;
    } else if (currentSymbol.value && currentValue.value.length < 1) {
        currentValue.value = previousValue.value;
        previousValue.value = '';
        currentSymbol.value = '';

        return true;
    } else {
        const value = parseAmount(currentValue.value);

        emit('update:modelValue', value);
        close();

        return true;
    }
}

function close(): void {
    emit('update:show', false);
}

function onSheetOpen(): void {
    currentValue.value = getStringValue(props.modelValue);
}

function onSheetClosed(): void {
    close();
}
</script>

<style>
.numpad-sheet {
    height: auto;
}

.numpad-values {
    border-bottom: 1px solid var(--f7-page-bg-color);
}

.numpad-value {
    display: flex;
    position: relative;
    padding-left: 16px;
    line-height: 1;
    height: var(--ebk-numpad-value-height);
    align-items: center;
    box-sizing: border-box;
    user-select: none;
}

.numpad-value-small {
    font-size: var(--ebk-numpad-value-small-font-size);
}

.numpad-value-normal {
    font-size: var(--ebk-numpad-value-normal-font-size);
}

.numpad-value-large {
    font-size: var(--ebk-numpad-value-large-font-size);
}

.numpad-buttons {
    display: flex;
    flex-wrap: wrap;
}

.numpad-button {
    display: flex;
    position: relative;
    text-align: center;
    border-radius: 0;
    border-right: 1px solid var(--f7-page-bg-color);
    border-bottom: 1px solid var(--f7-page-bg-color);
    height: 60px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    box-sizing: border-box;
    user-select: none;
    touch-action: none;
}

.numpad-button-num {
    width: calc(80% / 3);
}

.numpad-button-function, .numpad-button-confirm {
    width: 20%;
}

.numpad-button-num.active-state, .numpad-button-function.active-state {
    background-color: rgba(var(--f7-color-black-rgb), .15);
}

.numpad-button-text {
    display: block;
    font-size: var(--ebk-numpad-normal-button-font-size);
    font-weight: normal;
    line-height: 1;
}

.numpad-button-text-normal {
    color: var(--f7-color-black);
}

.dark .numpad-button-text-normal {
    color: var(--f7-color-white);
}

.numpad-button-text-confirm {
    font-size: var(--ebk-numpad-confirm-button-font-size);
}
</style>
