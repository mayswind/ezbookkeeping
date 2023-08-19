<template>
    <v-text-field class="text-field-with-colored-label"
                  :type="hide ? 'password' : 'number'" :class="extraClass"
                  :color="color" :base-color="color"
                  :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  :rules="enableRules ? rules : []" v-model="currentValue"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown">
        <template #prepend-inner v-if="currency && prependText">
            <div style="margin-top: 2px">{{ prependText }}</div>
        </template>
        <template #append-inner v-if="currency && appendText">
            <div class="text-no-wrap">{{ appendText }}</div>
        </template>
    </v-text-field>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';

import transactionConstants from '@/consts/transaction.js';
import { numericCurrencyToString, stringCurrencyToNumeric } from '@/lib/currency.js';

export default {
    props: [
        'class',
        'color',
        'density',
        'currency',
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

        return {
            currentValue: numericCurrencyToString(self.modelValue),
            rules: [
                (v) => {
                    if (v === '') {
                        return self.$t('Amount value is not number');
                    }

                    try {
                        const val = parseFloat(v);
                        return (val >= transactionConstants.minAmount && val <= transactionConstants.maxAmount) || self.$t('Amount value exceeds limitation');
                    } catch (e) {
                        return self.$t('Amount value is not number');
                    }
                }
            ]
        }
    },
    computed: {
        ...mapStores(useSettingsStore),
        extraClass() {
            let finalClass = this.class;

            if (this.color) {
                finalClass += ` text-${this.color}`;
            }

            return finalClass;
        },
        prependText() {
            if (!this.currency) {
                return '';
            }

            const texts = this.getDisplayCurrencyPrependAndAppendText();

            if (!texts) {
                return '';
            }

            return texts.prependText;
        },
        appendText() {
            if (!this.currency) {
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
        'modelValue': function (newValue) {
            const numericCurrentValue = stringCurrencyToNumeric(this.currentValue);

            if (newValue !== numericCurrentValue) {
                const newStringValue = numericCurrencyToString(newValue, false, true);

                if (!(newStringValue === '0' && this.currentValue === '')) {
                    this.currentValue = newStringValue;
                }
            }
        },
        'currentValue': function (newValue) {
            this.$emit('update:modelValue', stringCurrencyToNumeric(newValue));
        }
    },
    methods: {
        onKeyUpDown(e) {
            if (e.target.value === '' || e.target.value === 0) {
                return;
            }

            let decimalLength = 0;
            let decimalIndex = e.target.value.indexOf('.');

            if (decimalIndex >= 0) {
                decimalLength = e.target.value.length - e.target.value.indexOf('.') - 1;
            }

            if (decimalLength > 2) {
                e.target.value = e.target.value.substring(0, Math.min(decimalIndex + 3, e.target.value.length - 1));
                this.currentValue = e.target.value;
                e.preventDefault();
                return;
            }

            try {

                const val = parseFloat(e.target.value);
                let maxLength = transactionConstants.maxAmount.toString().length;

                if (val < 0) {
                    maxLength = transactionConstants.minAmount.toString().length;
                }

                if (val < transactionConstants.minAmount) {
                    e.target.value = transactionConstants.minAmount;
                    this.currentValue = e.target.value;
                    e.preventDefault();
                } else if (val > transactionConstants.maxAmount) {
                    e.target.value = transactionConstants.maxAmount;
                    this.currentValue = e.target.value;
                    e.preventDefault();
                } else if (e.target.value.length > maxLength) {
                    e.target.value = e.target.value.substring(0, maxLength);
                    this.currentValue = e.target.value;
                    e.preventDefault();
                }
            } catch (e) {
                e.target.value = 0;
            }
        },
        getDisplayCurrencyPrependAndAppendText() {
            return this.$locale.getDisplayCurrencyPrependAndAppendText(this.currency, this.settingsStore.appSettings.currencyDisplayMode);
        }
    }
}
</script>
