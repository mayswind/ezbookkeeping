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

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.ts';

import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { isString, isNumber, removeAll } from '@/lib/common.ts';

export default {
    props: [
        'modelValue',
        'minValue',
        'maxValue',
        'currency',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;
        const userStore = useUserStore();

        return {
            previousValue: '',
            currentSymbol: '',
            currentValue: self.getStringValue(userStore, self.modelValue)
        }
    },
    computed: {
        ...mapStores(useUserStore),
        decimalSeparator() {
            return this.$locale.getCurrentDecimalSeparator(this.userStore);
        },
        supportDecimalSeparator() {
            if (!this.currency || !ALL_CURRENCIES[this.currency] || !isNumber(ALL_CURRENCIES[this.currency].fraction)) {
                return true;
            }

            return ALL_CURRENCIES[this.currency].fraction > 0;
        },
        currentDisplay() {
            const previousValue = this.$locale.appendDigitGroupingSymbol(this.userStore, this.previousValue);
            const currentValue = this.$locale.appendDigitGroupingSymbol(this.userStore, this.currentValue);

            if (this.currentSymbol) {
                return `${previousValue} ${this.currentSymbol} ${currentValue}`;
            } else {
                return currentValue;
            }
        },
        currentDisplayNumClass() {
            const currentDisplay = this.currentDisplay || '';

            if (currentDisplay.length >= 24) {
                return 'numpad-value-small';
            } else if (currentDisplay.length >= 16) {
                return 'numpad-value-normal';
            } else {
                return 'numpad-value-large';
            }
        },
        confirmText() {
            if (this.currentSymbol) {
                return '=';
            } else {
                return this.$t('OK');
            }
        }
    },
    methods: {
        getStringValue(userStore, value) {
            if (!isNumber(value) && !isString(value)) {
                return '';
            }

            let str = this.$locale.formatAmount(userStore, value, this.currency);

            const digitGroupingSymbol = this.$locale.getCurrentDigitGroupingSymbol(userStore);

            if (str.indexOf(digitGroupingSymbol) >= 0) {
                str = removeAll(str, digitGroupingSymbol);
            }

            const decimalSeparator = this.$locale.getCurrentDecimalSeparator(userStore);
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
        },
        inputNum(num) {
            if (!this.previousValue && this.currentSymbol === '−') {
                this.currentValue = '-' + this.currentValue;
                this.currentSymbol = '';
            }

            if (this.currentValue === '0') {
                this.currentValue = num.toString();
                return;
            } else if (this.currentValue === '-0') {
                this.currentValue = '-' + num.toString();
                return;
            }

            const decimalSeparatorPos = this.currentValue.indexOf(this.decimalSeparator);

            if (decimalSeparatorPos >= 0 && this.currentValue.substring(decimalSeparatorPos + 1, this.currentValue.length).length >= 2) {
                return;
            }

            const newValue = this.currentValue + num.toString();

            if (isNumber(this.minValue)) {
                const current = this.$locale.parseAmount(this.userStore, newValue);

                if (current < this.minValue) {
                    return;
                }
            }

            if (isNumber(this.maxValue)) {
                const current = this.$locale.parseAmount(this.userStore, newValue);

                if (current > this.maxValue) {
                    return;
                }
            }

            this.currentValue = newValue;
        },
        inputDoubleNum(num) {
            this.inputNum(num);
            this.inputNum(num);
        },
        inputDecimalSeparator() {
            if (this.currentValue.indexOf(this.decimalSeparator) >= 0) {
                return;
            }

            if (!this.previousValue && this.currentSymbol === '−') {
                this.currentValue = '-' + this.currentValue;
                this.currentSymbol = '';
            }

            if (this.currentValue.length < 1) {
                this.currentValue = '0';
            } else if (this.currentValue === '-') {
                this.currentValue = '-0';
            }

            this.currentValue = this.currentValue + this.decimalSeparator;
        },
        setSymbol(symbol) {
            if (this.currentValue) {
                if (this.currentSymbol) {
                    const lastFormulaCalcResult = this.confirm();

                    if (!lastFormulaCalcResult) {
                        return;
                    }
                }

                this.previousValue = this.currentValue;
                this.currentValue = '';
            }

            this.currentSymbol = symbol;
        },
        backspace() {
            if (!this.currentValue || this.currentValue.length < 1) {
                if (this.currentSymbol) {
                    this.currentValue = this.previousValue;
                    this.previousValue = '';
                    this.currentSymbol = '';
                }

                return;
            }

            this.currentValue = this.currentValue.substring(0, this.currentValue.length - 1);
        },
        clear() {
            this.currentValue = '';
            this.previousValue = '';
            this.currentSymbol = '';
        },
        confirm() {
            if (this.currentSymbol && this.currentValue.length >= 1) {
                const previousValue = this.$locale.parseAmount(this.userStore, this.previousValue);
                const currentValue = this.$locale.parseAmount(this.userStore, this.currentValue);
                let finalValue = 0;

                switch (this.currentSymbol) {
                    case '+':
                        finalValue = previousValue + currentValue;
                        break;
                    case '−':
                        finalValue = previousValue - currentValue;
                        break;
                    case '×':
                        finalValue = Math.round(previousValue * currentValue / 100);
                        break;
                    default:
                        finalValue = previousValue;
                }

                if (isNumber(this.minValue)) {
                    if (finalValue < this.minValue) {
                        this.$toast('Numeric Overflow');
                        return false;
                    }
                }

                if (isNumber(this.maxValue)) {
                    if (finalValue > this.maxValue) {
                        this.$toast('Numeric Overflow');
                        return false;
                    }
                }

                this.currentValue = this.getStringValue(this.userStore, finalValue);
                this.previousValue = '';
                this.currentSymbol = '';

                return true;
            } else if (this.currentSymbol && this.currentValue.length < 1) {
                this.currentValue = this.previousValue;
                this.previousValue = '';
                this.currentSymbol = '';

                return true;
            } else {
                const value = this.$locale.parseAmount(this.userStore, this.currentValue);

                this.$emit('update:modelValue', value);
                this.close();

                return true;
            }
        },
        close() {
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            this.currentValue = this.getStringValue(this.userStore, this.modelValue);
        },
        onSheetClosed() {
            this.close();
        }
    }
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
