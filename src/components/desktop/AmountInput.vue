<template>
    <v-text-field type="text" class="text-field-with-colored-label" :class="extraClass"
                  :color="color" :base-color="color"
                  :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  :rules="enableRules ? rules : []" v-model="currentValue" v-if="!hide"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown" @paste="onPaste">
        <template #prepend-inner v-if="currency && prependText">
            <div>{{ prependText }}</div>
        </template>
        <template #append-inner v-if="currency && appendText">
            <div class="text-no-wrap">{{ appendText }}</div>
        </template>
    </v-text-field>
    <v-text-field type="password" class="text-field-with-colored-label" :class="extraClass"
                  :color="color" :base-color="color"
                  :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  :rules="enableRules ? rules : []" v-model="currentValue" v-if="hide"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown" @paste="onPaste">
        <template #prepend-inner v-if="currency && prependText">
            <div>{{ prependText }}</div>
        </template>
        <template #append-inner v-if="currency && appendText">
            <div class="text-no-wrap">{{ appendText }}</div>
        </template>
    </v-text-field>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { removeAll } from '@/lib/common.ts';
import logger from '@/lib/logger.ts';

export default {
    props: [
        'class',
        'color',
        'density',
        'currency',
        'showCurrency',
        'label',
        'placeholder',
        'persistentPlaceholder',
        'disabled',
        'readonly',
        'hide',
        'enableRules',
        'modelValue'
    ],
    emits: [
        'update:modelValue'
    ],
    data() {
        const self = this;
        const userStore = useUserStore();

        return {
            currentValue: self.getFormattedValue(userStore, self.modelValue),
            rules: [
                (v) => {
                    if (v === '') {
                        return self.$t('Amount value is not number');
                    }

                    try {
                        const val = self.$locale.parseAmount(userStore, v);

                        if (Number.isNaN(val) || !Number.isFinite(val)) {
                            return self.$t('Amount value is not number');
                        }

                        return (val >= TRANSACTION_MIN_AMOUNT && val <= TRANSACTION_MAX_AMOUNT) || self.$t('Amount value exceeds limitation');
                    } catch (ex) {
                        logger.warn('cannot parse amount in amount input, original value is ' + v, ex);
                        return self.$t('Amount value is not number');
                    }
                }
            ]
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        extraClass() {
            let finalClass = this.class;

            if (this.color) {
                finalClass += ` text-${this.color}`;
            }

            if (this.currency && this.prependText) {
                finalClass += ` has-pretend-text`;
            }

            return finalClass;
        },
        prependText() {
            if (!this.currency || !this.showCurrency) {
                return '';
            }

            const texts = this.getDisplayCurrencyPrependAndAppendText();

            if (!texts) {
                return '';
            }

            return texts.prependText;
        },
        appendText() {
            if (!this.currency || !this.showCurrency) {
                return '';
            }

            const texts = this.getDisplayCurrencyPrependAndAppendText();

            if (!texts) {
                return '';
            }

            return texts.appendText;
        }
    },
    watch: {
        'currency': function () {
            const newStringValue = this.getFormattedValue(this.userStore, this.modelValue);

            if (!(newStringValue === '0' && this.currentValue === '')) {
                this.currentValue = newStringValue;
            }
        },
        'modelValue': function (newValue) {
            const numericCurrentValue = this.$locale.parseAmount(this.userStore, this.currentValue);

            if (newValue !== numericCurrentValue) {
                const newStringValue = this.getFormattedValue(this.userStore, newValue);

                if (!(newStringValue === '0' && this.currentValue === '')) {
                    this.currentValue = newStringValue;
                }
            }
        },
        'currentValue': function (newValue) {
            let finalValue = '';

            if (newValue) {
                const decimalSeparator = this.$locale.getCurrentDecimalSeparator(this.userStore);

                for (let i = 0; i < newValue.length; i++) {
                    if (!('0' <= newValue[i] && newValue[i] <= '9') && newValue[i] !== '-' && newValue[i] !== decimalSeparator) {
                        break;
                    }

                    finalValue += newValue[i];
                }
            }

            if (finalValue !== newValue) {
                this.currentValue = finalValue;
            } else {
                this.$emit('update:modelValue', this.$locale.parseAmount(this.userStore, finalValue));
            }
        }
    },
    methods: {
        onKeyUpDown(e) {
            if (e.altKey || e.ctrlKey || e.metaKey || (e.key.indexOf('F') === 0 && (e.key.length === 2 || e.key.length === 3))
                || e.key === 'ArrowLeft' || e.key === 'ArrowRight'
                || e.key === 'Home' || e.key === 'End' || e.key === 'Tab'
                || e.key === 'Backspace' || e.key === 'Delete' || e.key === 'Del') {
                return;
            }

            const digitGroupingSymbol = this.$locale.getCurrentDigitGroupingSymbol(this.userStore);
            const decimalSeparator = this.$locale.getCurrentDecimalSeparator(this.userStore);

            if (!('0' <= e.key && e.key <= '9') && e.key !== '-' && e.key !== decimalSeparator) {
                e.preventDefault();
                return;
            }

            let str = e.target.value;

            if (str.indexOf(digitGroupingSymbol) >= 0) {
                str = removeAll(str, digitGroupingSymbol);
            }

            if (e.key === '-' && str.lastIndexOf('-') > 0) {
                const lastMinusPos = str.lastIndexOf('-');
                e.target.value = str.substring(0, lastMinusPos) + str.substring(lastMinusPos + 1, str.length);
                this.currentValue = e.target.value;
                e.preventDefault();
                return;
            }

            if (e.key === decimalSeparator && str.indexOf(decimalSeparator) !== str.lastIndexOf(decimalSeparator)) {
                const lastDecimalSeparatorPos = str.lastIndexOf(decimalSeparator);
                e.target.value = str.substring(0, lastDecimalSeparatorPos) + str.substring(lastDecimalSeparatorPos + 1, str.length);
                this.currentValue = e.target.value;
                e.preventDefault();
                return;
            }

            if (e.key === decimalSeparator && (str.indexOf(decimalSeparator) === 0 || (str.indexOf(decimalSeparator) === 1 && str.charAt(0) === '-'))) {
                const negative = str.charAt(0) === '-';

                if (negative) {
                    str = str.substring(1);
                }

                str = (negative ? '-0' : '0') + str;
                e.target.value = str;
                this.currentValue = e.target.value;
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

                e.target.value = (negative ? '-' : '') + str;
                this.currentValue = e.target.value;
                e.preventDefault();
                return;
                }

            if (decimalLength > 2) {
                e.target.value = str.substring(0, Math.min(decimalIndex + 3, str.length - 1));
                this.currentValue = e.target.value;
                e.preventDefault();
                return;
            }

            try {
                const val = this.$locale.parseAmount(this.userStore, str);
                const finalValue = this.getValidFormattedValue(val, str, decimalIndex >= 0);

                if (finalValue !== str) {
                    e.target.value = finalValue;
                    this.currentValue = finalValue;
                    e.preventDefault();
                }
            } catch (ex) {
                ex.target.value = '0';
            }
        },
        onPaste(e) {
            if (!e.clipboardData) {
                e.preventDefault();
                return;
            }

            const text = e.clipboardData.getData('Text');

            if (!text) {
                e.preventDefault();
                return;
            }

            const value = this.$locale.parseAmount(this.userStore, text);
            const textualValue = this.getFormattedValue(this.userStore, value);
            const decimalSeparator = this.$locale.getCurrentDecimalSeparator(this.userStore);
            const hasDecimalSeparator = text.indexOf(decimalSeparator) >= 0;

            this.currentValue = this.getValidFormattedValue(value, textualValue, hasDecimalSeparator);
            e.preventDefault();
        },
        getValidFormattedValue(value, textualValue, hasDecimalSeparator) {
            let maxLength = TRANSACTION_MAX_AMOUNT.toString().length;

            if (value < 0) {
                maxLength = TRANSACTION_MIN_AMOUNT.toString().length;
            }

            if (value < TRANSACTION_MIN_AMOUNT) {
                return this.getFormattedValue(this.userStore, TRANSACTION_MIN_AMOUNT);
            } else if (value > TRANSACTION_MAX_AMOUNT) {
                return this.getFormattedValue(this.userStore, TRANSACTION_MAX_AMOUNT);
            }

            if (!hasDecimalSeparator && textualValue.length > maxLength) {
                return textualValue.substring(0, maxLength);
            } else if (hasDecimalSeparator && textualValue.length > maxLength + 1) {
                return textualValue.substring(0, maxLength + 1);
            }

            return textualValue;
        },
        getFormattedValue(userStore, value) {
            if (!Number.isNaN(value) && Number.isFinite(value)) {
                const digitGroupingSymbol = this.$locale.getCurrentDigitGroupingSymbol(userStore);
                return removeAll(this.$locale.formatAmount(userStore, value, this.currency), digitGroupingSymbol);
            }

            return '0';
        },
        getDisplayCurrencyPrependAndAppendText() {
            const numericCurrentValue = this.$locale.parseAmount(this.userStore, this.currentValue);
            const isPlural = numericCurrentValue !== 100 && numericCurrentValue !== -100;

            return this.$locale.getAmountPrependAndAppendText(this.settingsStore, this.userStore, this.currency, isPlural);
        }
    }
}
</script>

<style>
.text-field-with-colored-label.has-pretend-text .v-field__input {
    padding-left: 0.5rem;
}
</style>
